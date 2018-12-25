package handlers

import (
	"github.com/FourthState/plasma-mvp-sidechain/msgs"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type UtxoStore interface {
	// methods it expects the store to implement
}

func NewSpendHandler() sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		spendMsg, ok := msg.(msgs.SpendMsg)
		// TODO: Is it okay to panic here or do we return an error?
		if !ok {
			panic("Msg does not implement SpendMsg")
		}

		// continue on here

		return sdk.Result{}
	}
}
