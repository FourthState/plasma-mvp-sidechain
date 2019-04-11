package handlers

import (
	"crypto/sha256"
	//"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/msgs"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	//ethcmn "github.com/ethereum/go-ethereum/common"
	"math/big"
)

// returns the next tx index in the current block
type NextTxIndex func() uint16

// FeeUpdater updates the aggregate fee amount in a block
type FeeUpdater func(amt *big.Int) sdk.Error

func NewSpendHandler(utxoStore store.UTXOStore, plasmaStore store.PlasmaStore, nextTxIndex NextTxIndex, feeUpdater FeeUpdater) sdk.Handler {
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
		//fmt.Println("NewSpendHandler")
		for i, key := range spendMsg.GetSigners() {
			inputKeys = append(inputKeys, append(key[:], spendMsg.InputAt(uint8(i)).Position.Bytes()...))
		}

		// try to spend the inputs. Abort if the inputs don't exist or have been spent
		//fmt.Println("Spending first output")
		//fmt.Println(common.BytesToAddress(inputKeys[0][:common.AddressLength]).String(), spendMsg.Input0.Position, ethcmn.ToHex(spenderKeys[0]))
		res := utxoStore.SpendUTXO(ctx, common.BytesToAddress(inputKeys[0][:common.AddressLength]), spendMsg.Input0.Position, spenderKeys)
		if !res.IsOK() {
			return res
		}
		if spendMsg.HasSecondInput() {
			//fmt.Println("Spending second output")
			res := utxoStore.SpendUTXO(ctx, common.BytesToAddress(inputKeys[1][:common.AddressLength]), spendMsg.Input1.Position, spenderKeys)
			if !res.IsOK() {
				return res
			}
		}

		// update the aggregate fee amount for the block
		if err := feeUpdater(spendMsg.Fee); err != nil {
			return sdk.ErrInternal("error updating the aggregate fee").Result()
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
