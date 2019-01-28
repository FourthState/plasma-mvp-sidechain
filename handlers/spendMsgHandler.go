package handlers

import (
	"crypto/sha256"
	"github.com/FourthState/plasma-mvp-sidechain/msgs"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
)

// returns the next tx index in the current block
type NextTxIndex func() uint16

func NewSpendHandler(utxoStore store.UTXOStore, plasmaStore store.PlasmaStore, nextTxIndex NextTxIndex) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		spendMsg, ok := msg.(msgs.SpendMsg)
		if !ok {
			panic("Msg does not implement SpendMsg")
		}

		txIndex := nextTxIndex()
		blockHeight := plasmaStore.NextPlasmaBlockNum(ctx)

		// construct the confirmation hash
		merkleHash := spendMsg.MerkleHash()
		header := ctx.BlockHeader().DataHash
		confirmationHash := sha256.Sum256(append(merkleHash, header...))

		// positional part of these keys addeded when the new positions are created
		var spenderKeys [][]byte
		var positions []plasma.Position
		positions = append(positions, plasma.NewPosition(blockHeight, txIndex, 0, nil))
		spenderKeys = append(spenderKeys, append(spendMsg.Output0.Owner[:], positions[0].Bytes()...))
		if spendMsg.HasSecondOutput() {
			positions = append(positions, plasma.NewPosition(blockHeight, txIndex, 1, nil))
			spenderKeys = append(spenderKeys, append(spendMsg.Output1.Owner[:], positions[1].Bytes()...))
		}

		var inputKeys [][]byte
		for i, key := range spendMsg.GetSigners() {
			inputKeys = append(inputKeys, append(key[:], spendMsg.InputAt(uint8(i)).Position.Bytes()...))
		}

		// try to spend the inputs. Abort if the inputs don't exist or have been spent
		res := utxoStore.SpendUTXO(ctx, common.BytesToAddress(inputKeys[0][:common.AddressLength]), spendMsg.Input0.Position, spenderKeys)
		if !res.IsOK() {
			return res
		}
		if spendMsg.HasSecondInput() {
			res := utxoStore.SpendUTXO(ctx, common.BytesToAddress(inputKeys[1][:common.AddressLength]), spendMsg.Input1.Position, spenderKeys)
			if !res.IsOK() {
				return res
			}
		}

		// create new outputs
		for i, _ := range spenderKeys {
			utxo := store.UTXO{
				InputKeys:        inputKeys,
				ConfirmationHash: confirmationHash[:],
				MerkleHash:       merkleHash,
				Output:           spendMsg.OutputAt(uint8(i)),
				Spent:            false,
				Position:         positions[i],
			}

			utxoStore.StoreUTXO(ctx, utxo)
		}

		return sdk.Result{}
	}
}
