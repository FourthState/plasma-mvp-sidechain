package handlers

import (
	"github.com/FourthState/plasma-mvp-sidechain/msgs"
	//"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	//"math/big"
)

func InitiatePresenceClaimHandler(utxoStore store.UTXOStore, nextTxIndex NextTxIndex, client plasmaConn) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		_, ok := msg.(msgs.InitiatePresenceClaimMsg)
		if !ok {
			panic("Msg does not implement InitiatePresenceClaimMsg")
		}
		return sdk.Result{}
	}
}
