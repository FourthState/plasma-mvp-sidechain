package handlers

import (
	"github.com/FourthState/plasma-mvp-sidechain/msgs"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewDepositHandler(depositStore store.DepositStore, txStore store.TxStore, blockStore store.BlockStore, nextTxIndex NextTxIndex, client plasmaConn) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		depositMsg, ok := msg.(msgs.IncludeDepositMsg)
		if !ok {
			panic("Msg does not implement IncludeDepositMsg")
		}

		// Increment txIndex so that it doesn't collide with SpendMsg
		nextTxIndex()

		deposit, _, _ := client.GetDeposit(blockStore.CurrentPlasmaBlockNum(ctx), depositMsg.DepositNonce)

		dep := store.Deposit{
			Deposit: deposit,
			Spent:   false,
			Spender: nil,
		}
		depositStore.StoreDeposit(ctx, depositMsg.DepositNonce, dep)
		txStore.StoreDepositWithAccount(ctx, depositMsg.DepositNonce, deposit)

		return sdk.Result{}
	}
}
