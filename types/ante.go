package types

import (
	"bytes"
	"reflect"
	//"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	crypto "github.com/tendermint/go-crypto"
)

// NewAnteHandler returns an AnteHandler that checks signatures,
// and deducts fees from the first signer.
func NewAnteHandler(utxoMapper UTXOMapper, txIndex *uint16, feeAmount *uint64) sdk.AnteHandler {
	return func(
		ctx sdk.Context, tx sdk.Tx,
	) (_ sdk.Context, _ sdk.Result, abort bool) {

		sigs := tx.GetSignatures()
		if len(sigs) == 0 {
			return ctx,
				sdk.ErrUnauthorized("no signers").Result(),
				true
		}

		msg := tx.GetMsg()

		_, ok := tx.(BaseTx)
		if !ok {
			return ctx, sdk.ErrInternal("tx must be in form of BaseTx").Result(), true
		}

		// Assert that number of signatures is correct.
		var signerAddrs = msg.GetSigners()
		if len(sigs) != len(signerAddrs) {
			return ctx,
				sdk.ErrUnauthorized("wrong number of signers").Result(),
				true
		}

		spendMsg, ok := msg.(SpendMsg)
		if !ok {
			return ctx, sdk.ErrInternal("Msg must be of type SpendMsg").Result(), true
		}
		signBytes := spendMsg.GetSignBytes()

		//Check that inputs are valid and having valid signatures
		position1 := Position{spendMsg.Blknum1, spendMsg.Txindex1, spendMsg.Oindex1, spendMsg.DepositNum1}
		res := processSig(ctx, utxoMapper, position1, signerAddrs[0], sigs[0], signBytes)

		if !res.IsOK() {
			return ctx, res, true
		}

		signBytes = position1.GetSignBytes()

		//Check that confirm signature is valid
		res = processConfirmSig(ctx, utxoMapper, position1, spendMsg.ConfirmSigs1, signBytes)
		if !res.IsOK() {
			return ctx, res, true
		}
		//Verify validity of second input
		if ValidAddress(spendMsg.Owner2) {
			position2 := Position{spendMsg.Blknum2, spendMsg.Txindex2, spendMsg.Oindex2, spendMsg.DepositNum2}
			res = processSig(ctx, utxoMapper, position2, signerAddrs[1], sigs[1], signBytes)
			if !res.IsOK() {
				return ctx, res, true
			}

			signBytes = position2.GetSignBytes()

			res = processConfirmSig(ctx, utxoMapper, position2, spendMsg.ConfirmSigs2, signBytes)
			if !res.IsOK() {
				return ctx, res, true
			}
		}

		// If DeliverTx() update fee
		if !ctx.IsCheckTx() {
			header := ctx.BlockHeader()
			feeTxIndex := uint16(header.GetNumTxs()) - 1
			if *txIndex == feeTxIndex {
				//Assert Denom1 contains correct amount of fees
				if *feeAmount != spendMsg.Denom1 {
					return ctx, sdk.ErrUnauthorized("Fees collected does not match fees reported").Result(), true
				}

			} else {

				//Increment amount of fees collected
				(*feeAmount) += spendMsg.Fee
			}
		}

		// TODO: tx tags (?)
		return ctx, sdk.Result{}, false // continue...
	}
}

func processSig(
	ctx sdk.Context, um UTXOMapper,
	position Position, addr crypto.Address, sig sdk.StdSignature, signBytes []byte) (
	res sdk.Result) {
	// Check UTXO is not nil.
	utxo := um.GetUTXO(ctx, position)
	if utxo == nil {
		return sdk.ErrUnknownRequest("UTXO trying to be spent, does not exist").Result()
	}

	if !bytes.Equal(sig.PubKey.Bytes(), addr.Bytes()) {
		return sdk.ErrUnauthorized("Signer Address does not match owner").Result()
	}

	hash := ethcrypto.Keccak256(signBytes)

	pubKey1, err1 := ethcrypto.SigToPub(hash, sig.Signature.Bytes())

	if err1 != nil || !reflect.DeepEqual(ethcrypto.PubkeyToAddress(*pubKey1).Bytes(), addr.Bytes()) {
		return sdk.ErrUnauthorized("signature 1 verification failed").Result()
	}

	return sdk.Result{}
}

func processConfirmSig(
	ctx sdk.Context, utxoMapper UTXOMapper,
	position Position, sig [2]crypto.Signature, signBytes []byte) (
	res sdk.Result) {

	utxo := utxoMapper.GetUTXO(ctx, position)
	if utxo == nil {
		return sdk.ErrUnknownRequest("UTXO trying to be spent, does not exist").Result()
	}
	inputAddresses := utxo.GetInputAddresses()

	ethsigs := make([]crypto.SignatureSecp256k1, 2)
	for i, s := range sig {
		ethsigs[i] = s.(crypto.SignatureSecp256k1)
	}

	hash := ethcrypto.Keccak256(signBytes)

	pubKey1, err1 := ethcrypto.SigToPub(hash, ethsigs[0].Bytes())
	if err1 != nil || !reflect.DeepEqual(ethcrypto.PubkeyToAddress(*pubKey1).Bytes(), inputAddresses[0].Bytes()) {
		return sdk.ErrUnauthorized("confirm signature 1 verification failed").Result()
	}

	if ValidAddress(inputAddresses[1]) {
		pubKey2, err2 := ethcrypto.SigToPub(hash, ethsigs[1].Bytes())
		if err2 != nil || !reflect.DeepEqual(ethcrypto.PubkeyToAddress(*pubKey2).Bytes(), inputAddresses[1].Bytes()) {
			return sdk.ErrUnauthorized("confirm signature 2 verification failed").Result()
		}
	}

	return sdk.Result{}
}
