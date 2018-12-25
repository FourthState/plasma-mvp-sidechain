package handlers

import (
	"bytes"
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/eth"
	"github.com/FourthState/plasma-mvp-sidechain/msgs"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

// FeeUpdater updates the aggregate fee amount in a block
type FeeUpdater func(amt *big.Int) sdk.Error

func NewAnteHandler(utxoStore store.UTXOStore, plasmaStore store.PlasmaStore, client *eth.Plasma, feeUpdater FeeUpdater) sdk.AnteHandler {
	return func(ctx sdk.Context, tx sdk.Tx, simulate bool) (newCtx sdk.Context, res sdk.Result, abort bool) {
		spendMsg, ok := tx.(msgs.SpendMsg)
		if !ok {
			return ctx, sdk.ErrInternal("tx must in the form of a spendMsg").Result(), true
		}

		txHash := spendMsg.TxHash()

		var totalInputAmt, totalOutputAmt *big.Int

		/* validate the first input */
		amt, res := validateInput(ctx, spendMsg.Input0, true, spendMsg.Fee, utxoStore, client)
		if !res.IsOK() {
			return ctx, res, true
		}
		res = validateSignatures(ctx, spendMsg.Input0, txHash, utxoStore)
		if !res.IsOK() {
			return ctx, res, true
		}
		if client.HasTxBeenExited(spendMsg.Input0.Position) {
			return ctx, ErrExitedInput(DefaultCodespace, "first input utxo has exited").Result(), true
		}

		totalInputAmt = amt

		// store confirm signatures
		if !spendMsg.Input0.Position.IsDeposit() && spendMsg.Input0.TxIndex < 1<<16-1 {
			plasmaStore.StoreConfirmSignatures(ctx, spendMsg.Input0.Position, spendMsg.Input0.ConfirmSignatures)
		}

		/* validate second input if applicable */
		if spendMsg.HasSecondInput() {
			amt, res = validateInput(ctx, spendMsg.Input1, false, nil, utxoStore, client)
			if !res.IsOK() {
				return ctx, res, true
			}
			res = validateSignatures(ctx, spendMsg.Input1, txHash, utxoStore)
			if !res.IsOK() {
				return ctx, res, true
			}
			if client.HasTxBeenExited(spendMsg.Input1.Position) {
				return ctx, ErrExitedInput(DefaultCodespace, "second input utxo has exited").Result(), true
			}

			// store confirm signature
			if !spendMsg.Input1.Position.IsDeposit() && spendMsg.Input1.TxIndex < 1<<16-1 {
				plasmaStore.StoreConfirmSignatures(ctx, spendMsg.Input1.Position, spendMsg.Input1.ConfirmSignatures)
			}

			totalInputAmt = totalInputAmt.Add(totalInputAmt, amt)
		}

		// input0 + input1 (totalInputAmt) == output0 + output1 + Fee (totalOutputAmt)
		totalOutputAmt = totalOutputAmt.Add(totalOutputAmt.Add(spendMsg.Output0.Amount, spendMsg.Output1.Amount), spendMsg.Fee)
		if totalInputAmt.Cmp(totalOutputAmt) != 0 {
			return ctx, msgs.ErrInvalidTransaction(DefaultCodespace, "Inputs do not equal Outputs (+ fee)").Result(), true
		}

		// only update fee when we are actually delivering tx
		if !ctx.IsCheckTx() && !simulate {
			err := feeUpdater(spendMsg.Fee)
			if err != nil {
				return ctx, err.Result(), true
			}
		}

		return ctx, sdk.Result{}, false
	}
}

// validates the inputs against the utxo store and returns the amount of the respective input
func validateInput(ctx sdk.Context, input plasma.Input, firstInput bool, feeAmount *big.Int,
	utxoStore store.UTXOStore, client *eth.Plasma) (*big.Int, sdk.Result) {
	var amt *big.Int
	if input.IsDeposit() {
		// TODO: add utxo to the app state?
		deposit, ok := client.GetDeposit(input.DepositNonce)
		if !ok {
			return nil, msgs.ErrInvalidTransaction(DefaultCodespace, "deposit input does not exist").Result()
		}
		if !bytes.Equal(deposit.Owner[:], input.Owner[:]) {
			return nil, sdk.ErrUnauthorized(fmt.Sprintf("signer does not own the deposit: Signer: %x, Owner: %x", deposit.Owner, input.Owner)).Result()
		}

		amt = deposit.Amount
	} else {
		inputUTXO, ok := utxoStore.GetUTXO(ctx, input.Owner, input.Position)
		if !ok {
			return nil, msgs.ErrInvalidTransaction(DefaultCodespace, "first input does not exist").Result()
		}
		if !bytes.Equal(inputUTXO.Output.Owner[:], input.Owner[:]) {
			return nil, sdk.ErrUnauthorized(fmt.Sprintf("signer does not own the input: Signer: %x, Owner: %x", inputUTXO.Output.Owner, input.Owner)).Result()
		}

		amt = inputUTXO.Output.Amount
	}

	// first input must pay the fee
	if firstInput && amt.Cmp(feeAmount) < 0 {
		return nil, ErrInsufficientFee(DefaultCodespace, "first input has an insufficient amount to pay the fee").Result()
	}

	return amt, sdk.Result{}
}

// validates the input's signature and confirm signatures
func validateSignatures(ctx sdk.Context, input plasma.Input, txHash [32]byte, utxoStore store.UTXOStore) sdk.Result {
	/* check transaction signatures */
	pubKey, err := crypto.SigToPub(txHash[:], input.Signature[:])
	if err != nil {
		return sdk.ErrInternal("Error recovering address from signature").Result()
	}

	signer := crypto.PubkeyToAddress(*pubKey)
	if !bytes.Equal(signer[:], input.Owner[:]) {
		return sdk.ErrUnauthorized(fmt.Sprintf("Signature mismatch. Signer: %x, Owner: %x", signer, input.Owner)).Result()
	}

	/* check input confirm signatures if the input is not a deposit nor fee utxo */
	if !input.IsDeposit() && input.TxIndex != 1<<16-1 {
		// `validateInput` ensures the output exists
		inputUTXO, _ := utxoStore.GetUTXO(ctx, input.Owner, input.Position)

		if len(inputUTXO.InputKeys) != len(input.ConfirmSignatures) {
			return msgs.ErrInvalidTransaction(DefaultCodespace, "incorrect number of confirm signatures").Result()
		}

		confirmationHash := inputUTXO.ConfirmationHash[:]
		for i, key := range inputUTXO.InputKeys {
			address := key[:common.AddressLength]

			pubKey, _ := crypto.SigToPub(confirmationHash, input.ConfirmSignatures[i][:])
			signer = crypto.PubkeyToAddress(*pubKey)
			if !bytes.Equal(signer[:], address) {
				return sdk.ErrUnauthorized(fmt.Sprintf("Confirm signature not signed by the correct address. Got: %x. Expected %x", signer, address)).Result()
			}
		}
	}

	return sdk.Result{}
}
