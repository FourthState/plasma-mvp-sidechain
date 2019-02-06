package handlers

import (
	"bytes"
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/msgs"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

// the reason for an interface is to allow the connection object
// to be cooked when testing the ante handler
type plasmaConn interface {
	GetDeposit(*big.Int) (plasma.Deposit, bool)
	HasTxBeenExited(plasma.Position) bool
}

func NewAnteHandler(utxoStore store.UTXOStore, plasmaStore store.PlasmaStore, client plasmaConn) sdk.AnteHandler {
	return func(ctx sdk.Context, tx sdk.Tx, simulate bool) (newCtx sdk.Context, res sdk.Result, abort bool) {
		spendMsg, ok := tx.(msgs.SpendMsg)
		if !ok {
			return ctx, sdk.ErrInternal("tx must in the form of a spendMsg").Result(), true
		}

		var totalInputAmt, totalOutputAmt *big.Int

		// attempt to recover signers
		signers := spendMsg.GetSigners()
		if len(signers) == 0 {
			return ctx, msgs.ErrInvalidTransaction(DefaultCodespace, "failed recovering signers").Result(), true
		}

		/* validate the first input */
		amt, res := validateInput(ctx, spendMsg.Input0, common.BytesToAddress(signers[0]), utxoStore, client)
		if !res.IsOK() {
			return ctx, res, true
		}
		if client.HasTxBeenExited(spendMsg.Input0.Position) {
			return ctx, ErrExitedInput(DefaultCodespace, "first input utxo has exited").Result(), true
		}
		// must cover the fee
		if amt.Cmp(spendMsg.Fee) < 0 {
			return ctx, ErrInsufficientFee(DefaultCodespace, "first input has an insufficient amount to pay the fee").Result(), true
		}

		totalInputAmt = amt

		// store confirm signatures
		if !spendMsg.Input0.Position.IsDeposit() && spendMsg.Input0.TxIndex < 1<<16-1 {
			plasmaStore.StoreConfirmSignatures(ctx, spendMsg.Input0.Position, spendMsg.Input0.ConfirmSignatures)
		}

		/* validate second input if applicable */
		if spendMsg.HasSecondInput() {
			if len(signers) != 2 {
				return ctx, msgs.ErrInvalidTransaction(DefaultCodespace, "second signature not present").Result(), true
			}
			amt, res = validateInput(ctx, spendMsg.Input1, common.BytesToAddress(signers[1]), utxoStore, client)
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
		totalOutputAmt = spendMsg.Output0.Amount
		if spendMsg.HasSecondOutput() {
			totalOutputAmt = totalOutputAmt.Add(totalOutputAmt.Add(totalOutputAmt, spendMsg.Output1.Amount), spendMsg.Fee)
		} else {
			totalOutputAmt = totalOutputAmt.Add(totalOutputAmt, spendMsg.Fee)
		}

		if totalInputAmt.Cmp(totalOutputAmt) != 0 {
			return ctx, msgs.ErrInvalidTransaction(DefaultCodespace, "inputs do not equal Outputs (+ fee)").Result(), true
		}

		return ctx, sdk.Result{}, false
	}
}

// validates the inputs against the utxo store and returns the amount of the respective input
func validateInput(ctx sdk.Context, input plasma.Input, signer common.Address, utxoStore store.UTXOStore, client plasmaConn) (*big.Int, sdk.Result) {
	var amt *big.Int

	if input.IsDeposit() {
		deposit, ok := client.GetDeposit(input.DepositNonce)
		if !ok {
			return nil, msgs.ErrInvalidTransaction(DefaultCodespace, "deposit, %s, does not exist", input.DepositNonce.String()).Result()
		}

		// add deposit to app state if non existent
		if !utxoStore.HasUTXO(ctx, deposit.Owner, input.Position) {
			utxo := store.UTXO{
				Output:   plasma.NewOutput(deposit.Owner, deposit.Amount),
				Position: input.Position,
				Spent:    false,
			}

			utxoStore.StoreUTXO(ctx, utxo)
		}

		// the owner of the deposit might not equal the signer so we must explicity
		// check for a match

		if !bytes.Equal(deposit.Owner[:], signer[:]) {
			return nil, sdk.ErrUnauthorized(fmt.Sprintf("signer does not own the deposit: Signer: %x, Owner: %x", signer, deposit.Owner)).Result()
		}

		amt = deposit.Amount
	} else {
		// inputUTXO must be owned by the signer due to the prefix so we do not need to
		// check the owner of the position
		inputUTXO, ok := utxoStore.GetUTXO(ctx, signer, input.Position)
		if !ok {
			return nil, msgs.ErrInvalidTransaction(DefaultCodespace, "input, %s, does not exist by owner %x", inputUTXO.Position, signer).Result()
		}
		if inputUTXO.Spent {
			return nil, msgs.ErrInvalidTransaction(DefaultCodespace, "input, %s, already spent", inputUTXO.Position).Result()
		}

		// validate confirm signatures if not a fee utxo
		if input.TxIndex < 1<<16-1 {
			res := validateConfirmSignatures(ctx, input, inputUTXO)
			if !res.IsOK() {
				return nil, res
			}
		}

		// check if the parent utxo has exited
		for _, key := range inputUTXO.InputKeys {
			utxo, _ := utxoStore.GetUTXOWithKey(ctx, key)
			if client.HasTxBeenExited(utxo.Position) {
				return nil, sdk.ErrUnauthorized(fmt.Sprintf("a parent of the input has exited. Owner: %x, Position: %s", utxo.Output.Owner, utxo.Position)).Result()
			}
		}

		amt = inputUTXO.Output.Amount
	}

	return amt, sdk.Result{}
}

// validates the input's confirm signatures
func validateConfirmSignatures(ctx sdk.Context, input plasma.Input, inputUTXO store.UTXO) sdk.Result {
	if len(inputUTXO.InputKeys) != len(input.ConfirmSignatures) {
		return msgs.ErrInvalidTransaction(DefaultCodespace, "incorrect number of confirm signatures").Result()
	}

	confirmationHash := utils.ToEthSignedMessageHash(inputUTXO.ConfirmationHash[:])[:]
	for i, key := range inputUTXO.InputKeys {
		address := key[:common.AddressLength]

		pubKey, _ := crypto.SigToPub(confirmationHash, input.ConfirmSignatures[i][:])
		signer := crypto.PubkeyToAddress(*pubKey)
		if !bytes.Equal(signer[:], address) {
			return sdk.ErrUnauthorized(fmt.Sprintf("confirm signature not signed by the correct address. Got: %x. Expected %x", signer, address)).Result()
		}
	}

	return sdk.Result{}
}
