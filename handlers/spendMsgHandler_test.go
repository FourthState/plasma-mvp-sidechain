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

var nextTxIndex = func() uint16 { return 0 }

func TestSpend(t *testing.T) {
	// ctx is at block height 0
	ctx, utxoStore, _ := setup()

	spendHandler := NewSpendHandler(utxoStore, nextTxIndex)

	// store inputs to be spent
	nilAddress := common.Address{}
	pos := plasma.NewPosition(utils.Big1, 0, 0, nil)
	utxo := store.UTXO{
		Output: plasma.NewOutput(nilAddress, big.NewInt(10)),
		// position (1,0,0)
		Position: pos,
		Spent:    false,
	}
	utxoStore.StoreUTXO(ctx, utxo)

	// create a msg that spends the first input and creates two outputs
	msg := msgs.SpendMsg{
		Transaction: plasma.Transaction{
			Input0:  plasma.NewInput(pos, common.Address{}, [65]byte{}, [][65]byte{}),
			Input1:  plasma.NewInput(plasma.NewPosition(nil, 0, 0, nil), nilAddress, [65]byte{}, [][65]byte{}),
			Output0: plasma.NewOutput(nilAddress, big.NewInt(10)),
			Output1: plasma.NewOutput(common.HexToAddress("1"), big.NewInt(10)),
			Fee:     utils.Big0,
		},
	}

	res := spendHandler(ctx, msg)
	require.True(t, res.IsOK(), "failed to handle spend")

	// check that the utxo store reflects the spend
	utxo, ok := utxoStore.GetUTXO(ctx, nilAddress, pos)
	require.True(t, ok, "input to the msg does not exist in the store")
	require.True(t, utxo.Spent, "input not marked as spent after the handler")

	// check that the new outputs were created

	// new first output was created at BlockHeight 0 and txIndex 0 and outputIndex 0
	pos = plasma.NewPosition(utils.Big0, 0, 0, nil)
	utxo, ok = utxoStore.GetUTXO(ctx, nilAddress, pos)
	require.True(t, ok, "new output was not created")
	require.False(t, utxo.Spent, "new output marked as spent")
	require.Equal(t, utxo.Output.Amount, big.NewInt(10), "new output has incorrect amount")

	// new second output was created at BlockHeight 0 and txIndex 0 and outputIndex 1
	pos = plasma.NewPosition(utils.Big0, 0, 1, nil)
	utxo, ok = utxoStore.GetUTXO(ctx, common.HexToAddress("1"), pos)
	require.True(t, ok, "new output was not created")
	require.False(t, utxo.Spent, "new output marked as spent")
	require.Equal(t, utxo.Output.Amount, big.NewInt(10), "new output has incorrect amount")
}
