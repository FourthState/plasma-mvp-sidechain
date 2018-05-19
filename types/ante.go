package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	crypto "github.com/tendermint/go-crypto"
	//"fmt"
	//abci "github.com/tendermint/abci/types"
	//"github.com/spf13/viper"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"reflect"
)

// NewAnteHandler returns an AnteHandler that checks signatures,
// and deducts fees from the first signer.
func NewAnteHandler(utxoMapper UTXOMapper, txIndex *uint16) sdk.AnteHandler {
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
		// GetSigners returns list of crypto.Address
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

		position1 := Position{spendMsg.Blknum1, spendMsg.Txindex1, spendMsg.Oindex1, spendMsg.DepositNum1}
		res := processSig(ctx, utxoMapper, position1, sigs[0], signBytes)
		if !res.IsOK() {
			return ctx, res, true
		}

		signBytes = position1.GetSignBytes()

		res = processConfirmSig(ctx, utxoMapper, position1, spendMsg.ConfirmSigs1, signBytes)
		if !res.IsOK() {
			return ctx, res, true
		}

		position2 := Position{spendMsg.Blknum2, spendMsg.Txindex2, spendMsg.Oindex2, spendMsg.DepositNum2}
		res = processSig(ctx, utxoMapper, position2, sigs[1], signBytes)
		if !res.IsOK() {
			return ctx, res, true
		}

		signBytes = position2.GetSignBytes()

		res = processConfirmSig(ctx, utxoMapper, position2, spendMsg.ConfirmSigs2, signBytes)
		if !res.IsOK() {
			return ctx, res, true
		}

		// If DeliverTx() update fee
		// Rough outline of dealing with fees
		if !ctx.IsCheckTx() {
			header := ctx.BlockHeader()
			feeTxIndex := uint16(header.GetNumTxs()) - 1
			feePosition := Position{uint64(ctx.BlockHeight()) * 1000, feeTxIndex, 0, 0}
			feeUTXO := utxoMapper.GetUTXO(ctx, feePosition)
			// change 0 to txindex
			if *txIndex == feeTxIndex { //is fee msg
				if feeUTXO.GetDenom() != spendMsg.Denom1 {
					return ctx, sdk.ErrUnauthorized("Fees collected does not match fees reported").Result(), true
				}
				utxoMapper.DeleteUTXO(ctx, feePosition)
			} else {
				// Is not fee Msg
				fee := spendMsg.Fee
				// first transaction in a block
				if feeUTXO != nil {
					fee = fee + feeUTXO.GetDenom()
					utxoMapper.DeleteUTXO(ctx, feePosition)
				}

				feeUTXO = NewBaseUTXO(crypto.Address([]byte("")), [2]crypto.Address{crypto.Address([]byte("")),
					crypto.Address([]byte(""))}, fee, feePosition)

				utxoMapper.AddUTXO(ctx, feeUTXO)
			}
		}

		// TODO: tx tags (?)

		return ctx, sdk.Result{}, false // continue...
	}
}

func processSig(
	ctx sdk.Context, um UTXOMapper,
	position Position, sig sdk.StdSignature, signBytes []byte) (
	res sdk.Result) {

	// Get the utxo.
	utxo := um.GetUTXO(ctx, position)
	if utxo == nil {
		return sdk.ErrUnknownRequest("UTXO trying to be spent, does not exist").Result()
	}

	if !sig.PubKey.VerifyBytes(signBytes, sig.Signature) {
		return sdk.ErrUnauthorized("signature verification failed").Result()
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
		return sdk.ErrUnauthorized("signature verification failed").Result()
	}

	if ValidAddress(inputAddresses[1]) {
		pubKey2, err2 := ethcrypto.SigToPub(hash, ethsigs[1].Bytes())
		if err2 != nil || !reflect.DeepEqual(ethcrypto.PubkeyToAddress(*pubKey2).Bytes(), inputAddresses[1].Bytes()) {
			return sdk.ErrUnauthorized("signature verification failed").Result()
		}
	}

	return sdk.Result{}
}
