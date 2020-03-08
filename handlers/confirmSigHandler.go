package handlers

import (
	"crypto/sha256"
	"github.com/FourthState/plasma-mvp-sidechain/msgs"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"math/big"
)

// NewSpendHandler sets the inputs of a spend msg to spent and creates new
// outputs that are added to the data store.
func NewConfirmSigHandler(ds store.DataStore, nextTxIndex NextTxIndex, feeUpdater FeeUpdater) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		confirmSigMsg, ok := msg.(msgs.ConfirmSigMsg)
		if !ok {
			panic("Msg does not implement SpendMsg")
		}

		/* Store Transaction and create new outputs */
		// TODO: Find best way to store sigs/inputs in DS
		// input1 := store.TxInput{}
		// input2 := store.TxInput{}
		// ds.StoreInput(ctx, input1)
		return sdk.Result{}
	}
}
