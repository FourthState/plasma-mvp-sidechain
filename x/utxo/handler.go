package utxo

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

// Return the next position for handler to store newly created UTXOs
// Secondary is true if NextPosition is meant to return secondary output positions for a single multioutput transaction
// If false, NextPosition will increment position to accomadate outputs for a new transaction
type NextPosition func(ctx sdk.Context, secondary bool) Position

// User-defined fee update function
type FeeUpdater func([]Output) sdk.Error

// Handler handles spends of arbitrary utxo implementation
func NewSpendHandler(um Mapper, nextPos NextPosition) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		spendMsg, ok := msg.(SpendMsg)
		if !ok {
			panic("Msg does not implement SpendMsg")
		}

		var inputKeys [][]byte
		inputs := spendMsg.Inputs()
		for _, in := range inputs {
			inKey := um.ConstructKey(in.Owner, in.Position)
			inputKeys = append(inputKeys, inKey)
		}

		txHash := tmhash.Sum(ctx.TxBytes())

		var spenderKeys [][]byte
		// Add outputs from store
		for i, o := range spendMsg.Outputs() {
			var next Position
			if i == 0 {
				next = nextPos(ctx, false)
			} else {
				next = nextPos(ctx, true)
			}
			spenderKeys = append(spenderKeys, um.ConstructKey(o.Owner, next))
			utxo := NewUTXOwithInputs(o.Owner, o.Amount, o.Denom, next, txHash, inputKeys)
			um.ReceiveUTXO(ctx, utxo)
		}

		// Spend inputs from store
		for _, i := range spendMsg.Inputs() {
			err := um.SpendUTXO(ctx, i.Owner, i.Position, spenderKeys)
			if err != nil {
				return err.Result()
			}
		}

		return sdk.Result{}
	}
}

// This function should be called within the antehandler
// Checks that the inputs = outputs + fee and handles fee collection
func AnteHelper(ctx sdk.Context, um Mapper, tx sdk.Tx, simulate bool, feeUpdater FeeUpdater) sdk.Error {
	msg := tx.GetMsgs()[0]
	spendMsg, ok := msg.(SpendMsg)
	if !ok {
		panic("Msg does not implement SpendMsg")
	}

	// Add up all inputs
	totalInput := map[string]uint64{}
	for _, i := range spendMsg.Inputs() {
		utxo := um.GetUTXO(ctx, i.Owner, i.Position)
		totalInput[utxo.Denom] += utxo.Amount
	}

	// Add up all outputs and fee
	totalOutput := map[string]uint64{}
	for _, o := range spendMsg.Outputs() {
		totalOutput[o.Denom] += o.Amount
	}
	for _, fee := range spendMsg.Fee() {
		totalOutput[fee.Denom] += fee.Amount
	}

	for denom, _ := range totalInput {
		if totalInput[denom] != totalOutput[denom] {
			return ErrInvalidTransaction(2, "Inputs do not equal Outputs")
		}
	}

	// Only update fee when we are actually delivering tx
	if !ctx.IsCheckTx() && !simulate {
		err := feeUpdater(spendMsg.Fee())
		return err
	}
	return nil
}
