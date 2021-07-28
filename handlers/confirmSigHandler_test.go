package handlers

import (
	"github.com/FourthState/plasma-mvp-sidechain/msgs"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestIncludeConfirmSig(t *testing.T) {
	// blockStore is at next block height 1
	ctx, ds := setup()

	// Give deposit a cooked connection that will always provide deposit with given position
	confirmSigHandler := NewConfirmSigHandler(ds)

	// Add tx with empty confirm sig into ds
	pos1 := plasma.NewPosition(utils.Big1, 1, 2, utils.Big0)
	pos2 := plasma.NewPosition(utils.Big1, 2, 3, utils.Big0)
	newOwner := common.HexToAddress("1")
	transaction := plasma.Transaction{
		Inputs:  []plasma.Input{plasma.NewInput(pos1, [65]byte{}, nil), plasma.NewInput(pos2, [65]byte{}, nil)},
		Outputs: []plasma.Output{plasma.NewOutput(newOwner, big.NewInt(10)), plasma.NewOutput(newOwner, big.NewInt(10))},
		Fee:     utils.Big0,
	}
	storeTransaction1 := store.Transaction{
		Transaction: transaction,
		ConfirmationHash: []byte{1234},
		Spent: []bool{false},
		SpenderTxs: [][]byte{},
		Position: pos1,
	}
	storeTransaction2 := store.Transaction{
		Transaction: transaction,
		ConfirmationHash: []byte{5678},
		Spent: []bool{false},
		SpenderTxs: [][]byte{},
		Position: pos2,
	}

	ds.StoreTx(ctx, storeTransaction1)
	ds.StoreTx(ctx, storeTransaction2)

	// create a msg that populates confirm sig fields in ds
	input1 := plasma.Input{Position: pos1, Signature: [65]byte{1}, ConfirmSignatures: [][65]byte{{1}}}
	input2 := plasma.Input{Position: pos2, Signature: [65]byte{1}, ConfirmSignatures: [][65]byte{{1}}}
	msg := msgs.ConfirmSigMsg{
		Input1: input1,
		Input2: input2,
	}

	confirmSigHandler(ctx, msg)
	tx1, ok := ds.GetTx(ctx, []byte{1234})
	tx2, ok := ds.GetTx(ctx, []byte{5678})

	require.True(t, ok, "tx does not exist in store")
	require.Equal(t, input1.ConfirmSignatures, tx1.Transaction.Inputs[0].ConfirmSignatures, "confirm sigs not populated")
	require.Equal(t, input2.ConfirmSignatures, tx2.Transaction.Inputs[1].ConfirmSignatures, "confirm sigs not populated")
}