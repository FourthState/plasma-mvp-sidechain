package auth

import (
	"encoding/binary"
	"fmt"
	types "github.com/FourthState/plasma-mvp-sidechain/types"
	utils "github.com/FourthState/plasma-mvp-sidechain/utils"
	"github.com/FourthState/plasma-mvp-sidechain/x/metadata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"reflect"

	"github.com/FourthState/plasma-mvp-sidechain/x/utxo"
)

// NewAnteHandler returns an AnteHandler that checks signatures,
// confirm signatures, and increments the feeAmount
func NewAnteHandler(utxoMapper utxo.Mapper, metadataMapper metadata.MetadataMapper, feeUpdater utxo.FeeUpdater) sdk.AnteHandler {
	return func(
		ctx sdk.Context, tx sdk.Tx, simulate bool,
	) (_ sdk.Context, _ sdk.Result, abort bool) {

		baseTx, ok := tx.(types.BaseTx)
		if !ok {
			return ctx, sdk.ErrInternal("tx must be in form of BaseTx").Result(), true
		}

		// Assert that there are signatures
		sigs := baseTx.GetSignatures()
		if len(sigs) == 0 {
			return ctx,
				sdk.ErrUnauthorized("no signers").Result(),
				true
		}

		// Base Tx must have only one msg
		msg := baseTx.GetMsgs()[0]

		// Assert that number of signatures is correct.
		var signerAddrs = msg.GetSigners()

		spendMsg, ok := msg.(types.SpendMsg)
		if !ok {
			return ctx, sdk.ErrInternal("msg must be of type SpendMsg").Result(), true
		}
		signBytes := spendMsg.GetSignBytes()

		// Verify the first input signature
		addr0 := common.BytesToAddress(signerAddrs[0].Bytes())
		position0 := types.PlasmaPosition{spendMsg.Blknum0, spendMsg.Txindex0, spendMsg.Oindex0, spendMsg.DepositNum0}

		res := checkUTXO(ctx, utxoMapper, position0, addr0)
		if !res.IsOK() {
			return ctx, res, true
		}

		res = processSig(addr0, sigs[0], signBytes)

		if !res.IsOK() {
			return ctx, res, true
		}

		// verify the confirmation signature if the input is not a deposit
		if position0.DepositNum == 0 && position0.TxIndex != 1<<16-1 {
			res = processConfirmSig(ctx, utxoMapper, metadataMapper, position0, addr0, spendMsg.Input0ConfirmSigs)
			if !res.IsOK() {
				return ctx, res, true
			}
		}

		// Verify the second input
		if utils.ValidAddress(spendMsg.Owner1) {
			addr1 := common.BytesToAddress(signerAddrs[1].Bytes())
			position1 := types.PlasmaPosition{spendMsg.Blknum1, spendMsg.Txindex1, spendMsg.Oindex1, spendMsg.DepositNum1}

			res := checkUTXO(ctx, utxoMapper, position1, addr1)
			if !res.IsOK() {
				return ctx, res, true
			}

			res = processSig(addr1, sigs[1], signBytes)

			if !res.IsOK() {
				return ctx, res, true
			}

			if position1.DepositNum == 0 && position1.TxIndex != 1<<16-1 {
				res = processConfirmSig(ctx, utxoMapper, metadataMapper, position1, addr1, spendMsg.Input1ConfirmSigs)
				if !res.IsOK() {
					return ctx, res, true
				}
			}
		}

		balanceErr := utxo.AnteHelper(ctx, utxoMapper, tx, simulate, feeUpdater)
		if balanceErr != nil {
			return ctx, balanceErr.Result(), true
		}

		// TODO: tx tags (?)
		return ctx, sdk.Result{}, false // continue...
	}
}

func processSig(
	addr common.Address, sig [65]byte, signBytes []byte) (
	res sdk.Result) {

	hash := ethcrypto.Keccak256(signBytes)
	pubKey, err := ethcrypto.SigToPub(hash, sig[:])

	if err != nil || !reflect.DeepEqual(ethcrypto.PubkeyToAddress(*pubKey).Bytes(), addr.Bytes()) {
		return sdk.ErrUnauthorized(fmt.Sprintf("signature verification failed for: %X", addr.Bytes())).Result()
	}

	return sdk.Result{}
}

func processConfirmSig(
	ctx sdk.Context, utxoMapper utxo.Mapper, metadataMapper metadata.MetadataMapper,
	position types.PlasmaPosition, addr common.Address, sigs [][65]byte) (
	res sdk.Result) {

	// Verify utxo exists
	utxo := utxoMapper.GetUTXO(ctx, addr.Bytes(), &position)
	if utxo == nil {
		return sdk.ErrUnknownRequest(fmt.Sprintf("confirm Sig verification failed: UTXO trying to be spent, does not exist: %v.", position)).Result()
	}
	plasmaUTXO, ok := utxo.(*types.BaseUTXO)
	if !ok {
		return sdk.ErrInternal("utxo must be of type BaseUTXO").Result()
	}
	inputAddresses := plasmaUTXO.GetInputAddresses()

	// Get the block hash
	blknumKey := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(blknumKey, plasmaUTXO.GetPosition().Get()[0].Uint64())
	blockHash := metadataMapper.GetMetadata(ctx, blknumKey)

	txHash := plasmaUTXO.GetTxHash()

	hash := append(txHash, blockHash...)
	confirmHash := ethcrypto.Keccak256(hash)

	pubKey0, err0 := ethcrypto.SigToPub(confirmHash, sigs[0][:])
	if err0 != nil || !reflect.DeepEqual(ethcrypto.PubkeyToAddress(*pubKey0).Bytes(), inputAddresses[0].Bytes()) {
		return sdk.ErrUnauthorized("confirm signature 0 verification failed").Result()
	}

	if utils.ValidAddress(inputAddresses[1]) {
		pubKey1, err1 := ethcrypto.SigToPub(confirmHash, sigs[1][:])
		if err1 != nil || !reflect.DeepEqual(ethcrypto.PubkeyToAddress(*pubKey1).Bytes(), inputAddresses[1].Bytes()) {
			return sdk.ErrUnauthorized("confirm signature 1 verification failed").Result()
		}
	}

	return sdk.Result{}
}

// Checks that utxo at the position specified exists, matches the address in the SpendMsg
// and returns the denomination associated with the utxo
func checkUTXO(ctx sdk.Context, mapper utxo.Mapper, position types.PlasmaPosition, addr common.Address) sdk.Result {
	utxo := mapper.GetUTXO(ctx, addr.Bytes(), &position)
	if utxo == nil {
		return sdk.ErrUnknownRequest(fmt.Sprintf("UTXO trying to be spent, does not exist: %v.", position)).Result()
	}

	// Verify that utxo owner equals input address in the transaction
	if !reflect.DeepEqual(utxo.GetAddress(), addr.Bytes()) {
		return sdk.ErrUnauthorized(fmt.Sprintf("signer does not match utxo owner, signer: %X  owner: %X", addr.Bytes(), utxo.GetAddress())).Result()
	}

	return sdk.Result{}
}
