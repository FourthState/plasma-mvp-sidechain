package auth

import (
	"fmt"
	types "github.com/FourthState/plasma-mvp-sidechain/types"
	utils "github.com/FourthState/plasma-mvp-sidechain/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"reflect"

	"github.com/FourthState/plasma-mvp-sidechain/x/utxo"
)

// NewAnteHandler returns an AnteHandler that checks signatures,
// confirm signatures, and increments the feeAmount
func NewAnteHandler(utxoMapper utxo.Mapper, feeUpdater utxo.FeeUpdater) sdk.AnteHandler {
	return func(
		ctx sdk.Context, tx sdk.Tx,
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

		if len(sigs) != len(signerAddrs) {
			return ctx,
				sdk.ErrUnauthorized("wrong number of signers").Result(),
				true
		}

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
		posSignBytes := position0.GetSignBytes()

		// Verify that confirmation signature
		res = processConfirmSig(ctx, utxoMapper, position0, addr0, spendMsg.ConfirmSigs0, posSignBytes)
		if !res.IsOK() {
			return ctx, res, true
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

			posSignBytes = position1.GetSignBytes()

			res = processConfirmSig(ctx, utxoMapper, position1, addr1, spendMsg.ConfirmSigs1, posSignBytes)
			if !res.IsOK() {
				return ctx, res, true
			}
		}

		balanceErr := utxo.AnteHelper(ctx, utxoMapper, tx, feeUpdater)
		if balanceErr != nil {
			return ctx, balanceErr.Result(), true
		}

		// TODO: tx tags (?)
		return ctx, sdk.Result{}, false // continue...
	}
}

func processSig(
	addr common.Address, sig types.Signature, signBytes []byte) (
	res sdk.Result) {

	hash := ethcrypto.Keccak256(signBytes)
	pubKey, err := ethcrypto.SigToPub(hash, sig.Bytes())

	if err != nil || !reflect.DeepEqual(ethcrypto.PubkeyToAddress(*pubKey).Bytes(), addr.Bytes()) {
		return sdk.ErrUnauthorized(fmt.Sprintf("signature verification failed for: %X", addr.Bytes())).Result()
	}

	return sdk.Result{}
}

func processConfirmSig(
	ctx sdk.Context, utxoMapper utxo.Mapper,
	position types.PlasmaPosition, addr common.Address, sigs [2]types.Signature, signBytes []byte) (
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

	hash := ethcrypto.Keccak256(signBytes)

	pubKey0, err0 := ethcrypto.SigToPub(hash, sigs[0].Bytes())
	if err0 != nil || !reflect.DeepEqual(ethcrypto.PubkeyToAddress(*pubKey0).Bytes(), inputAddresses[0].Bytes()) {
		return sdk.ErrUnauthorized("confirm signature 0 verification failed").Result()
	}

	if utils.ValidAddress(inputAddresses[1]) {
		pubKey1, err1 := ethcrypto.SigToPub(hash, sigs[1].Bytes())
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
