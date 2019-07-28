package handlers

import (
	"github.com/FourthState/plasma-mvp-sidechain/msgs"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewDepositHandler(ds store.DataStore, nextTxIndex NextTxIndex, client plasmaConn) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		depositMsg, ok := msg.(msgs.IncludeDepositMsg)
		if !ok {
			panic("Msg does not implement IncludeDepositMsg")
		}

		// Increment txIndex so that it doesn't collide with SpendMsg
		nextTxIndex()

		deposit, _, _ := client.GetDeposit(ds.PlasmaBlockHeight(ctx), depositMsg.DepositNonce)

		ds.StoreDeposit(ctx, depositMsg.DepositNonce, deposit)

		return sdk.Result{}
	}
}
