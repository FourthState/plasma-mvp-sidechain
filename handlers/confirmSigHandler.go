package handlers

import (
	"github.com/FourthState/plasma-mvp-sidechain/msgs"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewSpendHandler sets the inputs of a spend msg to spent and creates new
// outputs that are added to the data store.
func NewConfirmSigHandler(ds store.DataStore) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		confirmSigMsg, ok := msg.(msgs.ConfirmSigMsg)
		if !ok {
			panic("Msg does not implement SpendMsg")
		}

		/* Get tx, update tx data with confirm sigs, and store updated tx object */
		// TODO: Validation on correctness of confirm sigs?
		tx1, ok := ds.GetTxWithPosition(ctx, confirmSigMsg.Input1.Position)
		if !ok {
			panic("no transaction exists for the position provided")
		}

		for i := 0; i < len(tx1.Transaction.Inputs); i++ {
			if tx1.Transaction.Inputs[i].Position.TxIndex != confirmSigMsg.Input1.Position.TxIndex {
				continue
			}

			// TODO: How to update ConfirmSignatures? Only 1 confirm sig held in msg, tx.inputs[i].confirmsignatures supports array
			tx1.Transaction.Inputs[i].ConfirmSignatures = confirmSigMsg.Input1.ConfirmSignatures
		}

		ds.StoreTx(ctx, tx1)

		tx2, ok := ds.GetTxWithPosition(ctx, confirmSigMsg.Input2.Position)
		if !ok {
			panic("no transaction exists for the position provided")
		}

		for i := 0; i < len(tx2.Transaction.Inputs); i++ {
			if tx2.Transaction.Inputs[i].Position.TxIndex != confirmSigMsg.Input2.Position.TxIndex {
				continue
			}

			// TODO: How to update ConfirmSignatures? Only 1 confirm sig held in msg, tx.inputs[i].confirmsignatures supports array
			tx2.Transaction.Inputs[i].ConfirmSignatures = confirmSigMsg.Input2.ConfirmSignatures
		}

		ds.StoreTx(ctx, tx2)

		return sdk.Result{}
	}
}
