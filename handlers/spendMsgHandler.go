package handlers

import (
	"crypto/sha256"
	"github.com/FourthState/plasma-mvp-sidechain/msgs"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"math/big"
)

// returns the next tx index in the current block
type NextTxIndex func() uint16

func NewSpendHandler(utxoStore store.UTXOStore, nextTxIndex NextTxIndex) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		spendMsg, ok := msg.(msgs.SpendMsg)
		if !ok {
			panic("Msg does not implement SpendMsg")
		}

		txIndex := nextTxIndex()

		// construct the confirmation hash
		merkleHash := spendMsg.MerkleHash()
		header := ctx.BlockHeader().DataHash
		confirmationHash := sha256.Sum256(append(merkleHash[:], header...))

		var spenderKeys [][]byte
		spenderKeys = append(spenderKeys, spendMsg.Output0.Owner[:])
		if spendMsg.HasSecondOutput() {
			spenderKeys = append(spenderKeys, spendMsg.Output1.Owner[:])
		}

		var inputKeys [][]byte
		inputKeys = append(inputKeys, spendMsg.Input0.Owner[:])
		if spendMsg.HasSecondInput() {
			inputKeys = append(inputKeys, spendMsg.Input1.Owner[:])
		}

		// create new outputs
		for i, _ := range spenderKeys {
			position := plasma.NewPosition(big.NewInt(ctx.BlockHeight()), txIndex, uint8(i), nil)

			// Hacky solution. Keys should only be constructed within the store module.
			spenderKeys[i] = append(spenderKeys[i], position.Bytes()...)

			utxo := store.UTXO{
				InputKeys:        inputKeys,
				ConfirmationHash: confirmationHash,
				Output:           spendMsg.OutputAt(uint8(i)),
				Spent:            false,
				Position:         position,
			}

			utxoStore.StoreUTXO(ctx, utxo)
		}

		// spend the inputs
		utxoStore.SpendUTXO(ctx, spendMsg.Input0.Owner, spendMsg.Input0.Position, spenderKeys)
		if spendMsg.HasSecondInput() {
			utxoStore.SpendUTXO(ctx, spendMsg.Input1.Owner, spendMsg.Input1.Position, spenderKeys)
		}

		return sdk.Result{}
	}
}
