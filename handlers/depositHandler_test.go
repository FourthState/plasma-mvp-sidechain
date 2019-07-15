package handlers

import (
	"github.com/FourthState/plasma-mvp-sidechain/msgs"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestIncludeDeposit(t *testing.T) {
	// plasmaStore is at next block height 1
	ctx, utxoStore, _, _ := setup()

	// Give deposit a cooked connection that will always provide deposit with given position
	depositHandler := NewDepositHandler(utxoStore, nextTxIndex, conn{})

	// create a msg that spends the first input and creates two outputs
	msg := msgs.IncludeDepositMsg{
		DepositNonce: big.NewInt(5),
		Owner:        addr,
	}

	depositHandler(ctx, msg)

	plasmaPosition := plasma.NewPosition(nil, 0, 0, big.NewInt(5))
	utxo, ok := utxoStore.GetUTXO(ctx, addr, plasmaPosition)

	require.True(t, ok, "UTXO does not exist in store")
	require.Equal(t, addr, utxo.Output.Owner, "UTXO has wrong owner")
	require.Equal(t, big.NewInt(10), utxo.Output.Amount, "UTXO has wrong amount")
	require.False(t, utxo.Spent, "Deposit UTXO is incorrectly marked as spent")
	require.Equal(t, [][]byte{}, utxo.InputKeys, "Deposit UTXO has input keys set to non-nil value")
}
