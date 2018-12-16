package auth

import (
	"encoding/binary"
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/eth"
	types "github.com/FourthState/plasma-mvp-sidechain/types"
	utils "github.com/FourthState/plasma-mvp-sidechain/utils"
	"github.com/FourthState/plasma-mvp-sidechain/x/kvstore"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	amino "github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/tmhash"
	"reflect"

	"github.com/FourthState/plasma-mvp-sidechain/x/utxo"
)

// NewAnteHandler returns an AnteHandler that checks signatures,
// confirm signatures, and increments the feeAmount
func NewAnteHandler(utxoMapper utxo.Mapper, plasmaStore kvstore.KVStore, feeUpdater utxo.FeeUpdater, plasmaClient *eth.Plasma) sdk.AnteHandler {
	return func(
		ctx sdk.Context, tx sdk.Tx, simulate bool,
	) (_ sdk.Context, _ sdk.Result, abort bool) {

		baseTx, ok := tx.(types.BaseTx)
		if !ok {
			return ctx, sdk.ErrInternal("tx must be in form of BaseTx").Result(), true
		}

		sigs := baseTx.GetSignatures()

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

		res := checkUTXO(ctx, plasmaClient, utxoMapper, position0, addr0)
		if !res.IsOK() {
			return ctx, res, true
		}
		exitErr := hasTXExited(plasmaClient, position0)
		if exitErr != nil {
			return ctx, exitErr.Result(), true
		}
		if position0.IsDeposit() {
			deposit, _ := DepositExists(position0.DepositNum, plasmaClient)
			inputUTXO := utxo.NewUTXO(deposit.Owner.Bytes(), deposit.Amount.Uint64(), types.Denom, position0)
			utxoMapper.ReceiveUTXO(ctx, inputUTXO)
		}

		res = processSig(addr0, sigs[0], signBytes)

		if !res.IsOK() {
			return ctx, res, true
		}

		// verify the confirmation signature if the input is not a deposit
		if position0.DepositNum == 0 && position0.TxIndex != 1<<16-1 {
			res = processConfirmSig(ctx, utxoMapper, plasmaStore, position0, addr0, spendMsg.Input0ConfirmSigs)
			if !res.IsOK() {
				return ctx, res, true
			}

			setConfirmSigs(ctx, plasmaStore, spendMsg.Blknum0, uint64(spendMsg.Txindex0), spendMsg.Input0ConfirmSigs)
		}

		// Verify the second input
		if utils.ValidAddress(spendMsg.Owner1) {
			addr1 := common.BytesToAddress(signerAddrs[1].Bytes())
			position1 := types.PlasmaPosition{spendMsg.Blknum1, spendMsg.Txindex1, spendMsg.Oindex1, spendMsg.DepositNum1}

			exitErr := hasTXExited(plasmaClient, position1)
			if exitErr != nil {
				return ctx, exitErr.Result(), true
			}
			res := checkUTXO(ctx, plasmaClient, utxoMapper, position1, addr1)
			if !res.IsOK() {
				return ctx, res, true
			}
			if position1.IsDeposit() {
				deposit, _ := DepositExists(position1.DepositNum, plasmaClient)
				inputUTXO := utxo.NewUTXO(deposit.Owner.Bytes(), deposit.Amount.Uint64(), types.Denom, position1)
				utxoMapper.ReceiveUTXO(ctx, inputUTXO)
			}

			res = processSig(addr1, sigs[1], signBytes)

			if !res.IsOK() {
				return ctx, res, true
			}

			if position1.DepositNum == 0 && position1.TxIndex != 1<<16-1 {
				res = processConfirmSig(ctx, utxoMapper, plasmaStore, position1, addr1, spendMsg.Input1ConfirmSigs)
				if !res.IsOK() {
					return ctx, res, true
				}

				if spendMsg.Blknum0 != spendMsg.Blknum1 && spendMsg.Txindex0 != spendMsg.Txindex1 {
					setConfirmSigs(ctx, plasmaStore, spendMsg.Blknum1, uint64(spendMsg.Txindex1), spendMsg.Input1ConfirmSigs)
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
	signHash := utils.SignHash(hash)
	pubKey, err := ethcrypto.SigToPub(signHash, sig[:])

	if err != nil || !reflect.DeepEqual(ethcrypto.PubkeyToAddress(*pubKey).Bytes(), addr.Bytes()) {
		return sdk.ErrUnauthorized(fmt.Sprintf("signature verification failed for: %X", addr.Bytes())).Result()
	}

	return sdk.Result{}
}

func processConfirmSig(
	ctx sdk.Context, utxoMapper utxo.Mapper, plasmaStore kvstore.KVStore,
	position types.PlasmaPosition, addr common.Address, sigs [][65]byte) (
	res sdk.Result) {

	// Verify utxo exists
	input := utxoMapper.GetUTXO(ctx, addr.Bytes(), &position)
	if reflect.DeepEqual(input, utxo.UTXO{}) {
		return sdk.ErrUnknownRequest(fmt.Sprintf("confirm Sig verification failed: UTXO trying to be spent, does not exist: %v.", position)).Result()
	}
	// Get input addresses for input UTXO (grandfather inputs)
	inputAddresses := input.InputAddresses()
	if len(inputAddresses) != len(sigs) {
		return sdk.ErrUnauthorized("Wrong number of confirm sigs").Result()
	}

	// Get the block hash that input was created in
	blknumKey := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(blknumKey, input.Position.Get()[0].Uint64())
	key := append(utils.RootHashPrefix, blknumKey...)
	blockHash := plasmaStore.Get(ctx, key)

	// Create confirm signature hash
	hash := append(input.TxHash, blockHash...)
	confirmHash := tmhash.Sum(hash)
	signHash := utils.SignHash(confirmHash)

	for i, sig := range sigs {
		pubKey, err := ethcrypto.SigToPub(signHash, sig[:])
		if err != nil || !reflect.DeepEqual(ethcrypto.PubkeyToAddress(*pubKey).Bytes(), inputAddresses[i]) {
			return sdk.ErrUnauthorized(fmt.Sprintf("confirm signature %d verification failed", i)).Result()
		}
	}

	return sdk.Result{}
}

// Helper function for setting confirmation signatures into app state
// Key is formed as :  prefix + {Blknum, TxIndex}
// This prevents double storage of confirmation signatures
// Returns true if successful
func setConfirmSigs(ctx sdk.Context, plasmaStore kvstore.KVStore, blknum, txindex uint64, sigs [][65]byte) bool {
	cdc := amino.NewCodec()
	pos := [2]uint64{blknum, txindex}

	bz, err := cdc.MarshalBinaryBare(pos)
	if err != nil {
		return false
	}

	plasmaKey := append(utils.ConfirmSigPrefix, bz...)

	var confirmSigs []byte
	confirmSigs = append(confirmSigs, sigs[0][:]...)
	if len(sigs) == 2 {
		confirmSigs = append(confirmSigs, sigs[1][:]...)
	}
	plasmaStore.Set(ctx, plasmaKey, confirmSigs)
	return true
}

// Checks that utxo at the position specified exists, matches the address in the SpendMsg
// and returns the denomination associated with the utxo
func checkUTXO(ctx sdk.Context, plasmaClient *eth.Plasma, mapper utxo.Mapper, position types.PlasmaPosition, addr common.Address) sdk.Result {
	var inputAddress []byte
	input := mapper.GetUTXO(ctx, addr.Bytes(), &position)
	if position.IsDeposit() && reflect.DeepEqual(input, utxo.UTXO{}) {
		deposit, ok := DepositExists(position.DepositNum, plasmaClient)
		if !ok {
			return utxo.ErrInvalidUTXO(2, "Deposit UTXO does not exist yet").Result()
		}
		inputAddress = deposit.Owner.Bytes()
	} else {
		if !input.Valid {
			return sdk.ErrUnknownRequest(fmt.Sprintf("UTXO trying to be spent, is not valid: %v.", position)).Result()
		}
		inputAddress = input.Address
	}

	// Verify that utxo owner equals input address in the transaction
	if !reflect.DeepEqual(inputAddress, addr.Bytes()) {
		return sdk.ErrUnauthorized(fmt.Sprintf("signer does not match utxo owner, signer: %X  owner: %X", addr.Bytes(), inputAddress)).Result()
	}
	return sdk.Result{}
}

func DepositExists(nonce uint64, plasmaClient *eth.Plasma) (types.Deposit, bool) {
	deposit, err := plasmaClient.GetDeposit(sdk.NewUint(nonce))

	if err != nil {
		return types.Deposit{}, false
	}
	return *deposit, true
}

func hasTXExited(plasmaClient *eth.Plasma, pos types.PlasmaPosition) sdk.Error {
	if plasmaClient == nil {
		return nil
	}

	var positions [4]sdk.Uint
	for i, num := range pos.Get() {
		positions[i] = num
	}
	exited := plasmaClient.HasTXBeenExited(positions)
	if exited {
		return types.ErrInvalidTransaction(types.DefaultCodespace, "Input UTXO has already exited")
	}
	return nil
}
