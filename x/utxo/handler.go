package utxo

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Return the next position for handler to store newly created UTXOs
type NextPosition func(ctx sdk.Context) Position

// Proto function to create application's UTXO implementation
type ProtoUTXO func() UTXO

// User-defined fee update function
type FeeUpdater func([]Output) sdk.Error

// Handler handles spends of arbitrary utxo implementation
func NewSpendHandler(um Mapper, nextPos NextPosition, proto ProtoUTXO) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		spendMsg, ok := msg.(SpendMsg)
		if !ok {
			panic("Msg does not implement SpendMsg")
		}

		// Delete inputs from store
		for _, i := range spendMsg.Inputs() {
			um.DeleteUTXO(ctx, i.Owner, i.Position)
		}

		// Add outputs from store
		for _, o := range spendMsg.Outputs() {
			next := nextPos(ctx)
			utxo := proto()
			utxo.SetPosition(next)
			utxo.SetAddress(o.Owner)
			utxo.SetDenom(o.Denom)
			utxo.SetAmount(o.Amount)
			um.AddUTXO(ctx, utxo)
		}

		return sdk.Result{}
	}
}

// This function should be called within the antehandler
// Checks that the inputs = outputs + fee and handles fee collection
func AnteHelper(ctx sdk.Context, um Mapper, tx sdk.Tx, feeUpdater FeeUpdater) sdk.Error {
	msg := tx.GetMsgs()[0]
	spendMsg, ok := msg.(SpendMsg)
	if !ok {
		panic("Msg does not implement SpendMsg")
	}

	// Add up all inputs
	totalInput := map[string]uint64{}
	for _, i := range spendMsg.Inputs() {
		utxo := um.GetUTXO(ctx, i.Owner, i.Position)
		totalInput[utxo.GetDenom()] += utxo.GetAmount()
	}

	// Add up all outputs and fee
	totalOutput := map[string]uint64{}
	for _, o := range spendMsg.Outputs() {
		totalOutput[o.Denom] += o.Amount
	}
	for _, fee := range spendMsg.Fee() {
		totalOutput[fee.Denom] += fee.Amount
	}
	
	for denom, amount := range totalInput {
		if totalInput[denom] != totalOutput[denom] {
			return ErrInvalidTransaction(2, "Inputs do not equal Outputs")
		}
	}

	err := feeUpdater(spendMsg.Fee())
	return err
}
