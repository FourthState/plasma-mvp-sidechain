package handlers

import (
	"crypto/sha256"
	"github.com/FourthState/plasma-mvp-sidechain/msgs"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	//"github.com/FourthState/plasma-mvp-sidechain/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	//"github.com/ethereum/go-ethereum/crypto"
	//"math/big"
	"fmt"
)

func PostLogsHandler(claimStore store.PresenceClaimStore, utxoStore store.UTXOStore, plasmaStore store.PlasmaStore, nextTxIndex NextTxIndex, client plasmaConn) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		postLogsMsg, ok := msg.(msgs.PostLogsMsg)
		if !ok {
			panic("Msg does not implement InitiatePresenceClaimMsg")
		}

		claim, ok := claimStore.GetPresenceClaim(ctx, postLogsMsg.ClaimID)

		if !ok {
			msgs.ErrInvalidTransaction(DefaultCodespace, "No claim found with claimID").Result()
		}

		zeroAddress := common.HexToAddress("0x0000000000000000000000000000000000000001")
		utxo, ok := utxoStore.GetUTXO(ctx, zeroAddress, claim.UTXOPosition)

		if !ok {
			msgs.ErrInvalidTransaction(DefaultCodespace, "No claim found with claimID").Result()
		}

		claim.LogsHash = &(postLogsMsg.LogsHash)

		claimStore.StorePresenceClaim(ctx, claim)

		txIndex := nextTxIndex()
		blockHeight := plasmaStore.NextPlasmaBlockNum(ctx)

		// construct the confirmation hash

		hash := sha256.Sum256(postLogsMsg.TxHash())
		merkleHash := hash[:]
		header := ctx.BlockHeader().DataHash
		confirmationHash := sha256.Sum256(append(merkleHash, header...))

		//	txHash := utils.ToEthSignedMessageHash(postLogsMsg.TxHash())
		//	pubKey, _ := crypto.SigToPub(txHash, postLogsMsg.Signature)

		recipient := postLogsMsg.Beacons[0]
		fmt.Printf("Recipient %v", recipient.Hex())
		spenderKeys := [][]byte{append(recipient[:], claim.UTXOPosition.Bytes()...)}
		output := plasma.Output{
			Owner:  recipient,
			Amount: utxo.Output.Amount,
		}

		position := plasma.NewPosition(blockHeight, txIndex, 0, nil)
		var inputKeys [][]byte
		inputKeys = append(inputKeys, append(recipient.Bytes()[:], position.Bytes()...))

		newUtxo := store.UTXO{
			InputKeys:        inputKeys,
			SpenderKeys:      spenderKeys,
			ConfirmationHash: confirmationHash[:],
			MerkleHash:       merkleHash,
			Output:           output,
			Spent:            false,
			Position:         position,
		}

		utxoStore.StoreUTXO(ctx, newUtxo)
		oldUTXOKey := store.GetUTXOStoreKey(zeroAddress, claim.UTXOPosition)
		utxoStore.Delete(ctx, oldUTXOKey)

		return sdk.Result{}
	}
}
