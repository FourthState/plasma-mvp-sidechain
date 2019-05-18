package handlers

import (
	"github.com/FourthState/plasma-mvp-sidechain/msgs"
	//"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	//"github.com/FourthState/plasma-mvp-sidechain/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	//"github.com/ethereum/go-ethereum/crypto"
	//"math/big"
	//"fmt"
)

func PostLogsHandler(claimStore store.PresenceClaimStore, utxoStore store.UTXOStore, nextTxIndex NextTxIndex, client plasmaConn) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		postLogsMsg, ok := msg.(msgs.PostLogsMsg)
		if !ok {
			panic("Msg does not implement InitiatePresenceClaimMsg")
		}

		claim, ok := claimStore.GetPresenceClaim(ctx, postLogsMsg.ClaimID)

		if !ok {
			msgs.ErrInvalidTransaction(DefaultCodespace, "No claim found with claimID").Result()
		}

		claim.LogsHash = &(postLogsMsg.LogsHash)

		claimStore.StorePresenceClaim(ctx, claim)

		return sdk.Result{}
	}
}
