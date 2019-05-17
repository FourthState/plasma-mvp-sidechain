package handlers

import (
	"github.com/FourthState/plasma-mvp-sidechain/msgs"
	//"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
	//"math/big"
	"fmt"
)

func InitiatePresenceClaimHandler(claimStore store.PresenceClaimStore, nextTxIndex NextTxIndex, client plasmaConn) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		fmt.Println("InitiatePresenceClaimHandler begin")
		claimMsg, ok := msg.(msgs.InitiatePresenceClaimMsg)
		if !ok {
			panic("Msg does not implement InitiatePresenceClaimMsg")
		}

		txHash := utils.ToEthSignedMessageHash(claimMsg.TxHash())
		pubKey, err := crypto.SigToPub(txHash, *(claimMsg.Signature))

		if err != nil {
			msgs.ErrInvalidTransaction(DefaultCodespace, "failed recovering signers").Result()
		}

		burnerAddress := crypto.PubkeyToAddress(*pubKey)

		var logsHash []byte
		claim := store.PresenceClaim{
			ZoneID:       claimMsg.ZoneID,
			UTXOPosition: claimMsg.UTXOPosition,
			UserAddress:  burnerAddress,
			LogsHash:     &logsHash,
		}
		fmt.Println("InitiatePresenceClaimHandler store")

		claimStore.StorePresenceClaim(ctx, claim)

		fmt.Println("InitiatePresenceClaimHandler end")

		return sdk.Result{}
	}
}
