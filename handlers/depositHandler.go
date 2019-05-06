package handlers

import (
	"github.com/FourthState/plasma-mvp-sidechain/msgs"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"math/big"
)

func NewDepositHandler(utxoStore store.UTXOStore, nextTxIndex NextTxIndex, client plasmaConn) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		depositMsg, ok := msg.(msgs.IncludeDepositMsg)
		if !ok {
			panic("Msg does not implement IncludeDepositMsg")
		}
		depositPosition := plasma.NewPosition(big.NewInt(0), 0, 0, depositMsg.DepositNonce)

		// Increment txIndex so that it doesn't collide with SpendMsg
		nextTxIndex()

		deposit, _, _ := client.GetDeposit(big.NewInt(ctx.BlockHeight()), depositMsg.DepositNonce)

		utxo := store.UTXO{
			Output:   plasma.NewOutput(deposit.Owner, deposit.Amount),
			Position: depositPosition,
			Spent:    false,
		}
		utxoStore.StoreUTXO(ctx, utxo)
		return sdk.Result{}
	}
}
