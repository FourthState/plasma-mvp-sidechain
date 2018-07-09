package auth

import (
	types "github.com/FourthState/plasma-mvp-sidechain/types"
	utils "github.com/FourthState/plasma-mvp-sidechain/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"reflect"
)

// NewAnteHandler returns an AnteHandler that checks signatures,
// confirm signatures, and increments the feeAmount
func NewAnteHandler(utxoMapper types.UTXOMapper, txIndex *uint16, feeAmount *uint64) sdk.AnteHandler {
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

		msg := baseTx.GetMsg()

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
		addr1 := common.BytesToAddress(signerAddrs[0].Bytes())
		position1 := types.Position{spendMsg.Blknum1, spendMsg.Txindex1, spendMsg.Oindex1, spendMsg.DepositNum1}

		indenom1, res := checkUTXO(ctx, utxoMapper, position1, addr1)
		if !res.IsOK() {
			return ctx, res, true
		}

		res = processSig(addr1, sigs[0], signBytes)

		if !res.IsOK() {
			return ctx, res, true
		}
		posSignBytes := position1.GetSignBytes()

		// Verify that confirmation signature
		res = processConfirmSig(ctx, utxoMapper, position1, addr1, spendMsg.ConfirmSigs1, posSignBytes)
		if !res.IsOK() {
			return ctx, res, true
		}

		inputs := indenom1

		// Verify the second input
		if utils.ValidAddress(spendMsg.Owner2) {
			addr2 := common.BytesToAddress(signerAddrs[1].Bytes())
			position2 := types.Position{spendMsg.Blknum2, spendMsg.Txindex2, spendMsg.Oindex2, spendMsg.DepositNum2}

			indenom2, res := checkUTXO(ctx, utxoMapper, position2, addr2)
			if !res.IsOK() {
				return ctx, res, true
			}

			res = processSig(addr2, sigs[1], signBytes)

			if !res.IsOK() {
				return ctx, res, true
			}

			posSignBytes = position2.GetSignBytes()

			res = processConfirmSig(ctx, utxoMapper, position2, addr2, spendMsg.ConfirmSigs2, posSignBytes)
			if !res.IsOK() {
				return ctx, res, true
			}

			inputs += indenom2
		}

		if inputs != (spendMsg.Denom1 + spendMsg.Denom2 + spendMsg.Fee) {
			return ctx, sdk.ErrUnauthorized("inputs are not equal to outputs plus the fee").Result(), true
		}

		// Increment amount of fees collected
		if !ctx.IsCheckTx() {
			(*feeAmount) += spendMsg.Fee
		}

		// TODO: tx tags (?)
		return ctx, sdk.Result{}, false // continue...
	}
}

func processSig(
	addr common.Address, sig types.Signature, signBytes []byte) (
	res sdk.Result) {

	hash := ethcrypto.Keccak256(signBytes)
	pubKey1, err1 := ethcrypto.SigToPub(hash, sig.Bytes())

	if err1 != nil || !reflect.DeepEqual(ethcrypto.PubkeyToAddress(*pubKey1).Bytes(), addr.Bytes()) {
		return sdk.ErrUnauthorized("signature verification failed").Result()
	}

	return sdk.Result{}
}

func processConfirmSig(
	ctx sdk.Context, utxoMapper types.UTXOMapper,
	position types.Position, addr common.Address, sigs [2]types.Signature, signBytes []byte) (
	res sdk.Result) {

	// Verify utxo exists
	utxo := utxoMapper.GetUTXO(ctx, addr, position)
	if utxo == nil {
		return sdk.ErrUnknownRequest("Confirm Sig verification failed: UTXO trying to be spent, does not exist").Result()
	}
	inputAddresses := utxo.GetInputAddresses()

	hash := ethcrypto.Keccak256(signBytes)

	pubKey1, err1 := ethcrypto.SigToPub(hash, sigs[0].Bytes())
	if err1 != nil || !reflect.DeepEqual(ethcrypto.PubkeyToAddress(*pubKey1).Bytes(), inputAddresses[0].Bytes()) {
		return sdk.ErrUnauthorized("confirm signature 1 verification failed").Result()
	}

	if utils.ValidAddress(inputAddresses[1]) {
		pubKey2, err2 := ethcrypto.SigToPub(hash, sigs[1].Bytes())
		if err2 != nil || !reflect.DeepEqual(ethcrypto.PubkeyToAddress(*pubKey2).Bytes(), inputAddresses[1].Bytes()) {
			return sdk.ErrUnauthorized("confirm signature 2 verification failed").Result()
		}
	}

	return sdk.Result{}
}

// Checks that utxo at the position specified exists, matches the address in the SpendMsg
// and returns the denomination associated with the utxo
func checkUTXO(ctx sdk.Context, mapper types.UTXOMapper, position types.Position, addr common.Address) (indenom uint64, res sdk.Result) {
	utxo := mapper.GetUTXO(ctx, addr, position)
	if utxo == nil {
		return 0, sdk.ErrUnknownRequest("UTXO trying to be spend, does not exist").Result()
	}

	// Verify that utxo owner equals input address in the transaction
	if !reflect.DeepEqual(utxo.GetAddress().Bytes(), addr.Bytes()) {
		return 0, sdk.ErrUnauthorized("signer does not match utxo owner").Result()
	}

	return utxo.GetDenom(), sdk.Result{}
}
