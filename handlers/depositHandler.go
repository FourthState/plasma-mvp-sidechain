package handlers

import (
	"github.com/FourthState/plasma-mvp-sidechain/msgs"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewDepositHandler(outputStore store.OutputStore, blockStore store.BlockStore, nextTxIndex NextTxIndex, client plasmaConn) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		depositMsg, ok := msg.(msgs.IncludeDepositMsg)
		if !ok {
			panic("Msg does not implement IncludeDepositMsg")
		}

		// Increment txIndex so that it doesn't collide with SpendMsg
		nextTxIndex()

		deposit, _, _ := client.GetDeposit(blockStore.CurrentPlasmaBlockNum(ctx), depositMsg.DepositNonce)

		outputStore.StoreDeposit(ctx, depositMsg.DepositNonce, deposit)

		return sdk.Result{}
	}
}
