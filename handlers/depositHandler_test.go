package handlers

import (
	"github.com/FourthState/plasma-mvp-sidechain/msgs"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestIncludeDeposit(t *testing.T) {
	// blockStore is at next block height 1
	ctx, outputStore, blockStore := setup()

	// Give deposit a cooked connection that will always provide deposit with given position
	depositHandler := NewDepositHandler(outputStore, blockStore, nextTxIndex, conn{})

	// create a msg that spends the first input and creates two outputs
	msg := msgs.IncludeDepositMsg{
		DepositNonce: big.NewInt(5),
		Owner:        addr,
	}

	depositHandler(ctx, msg)

	deposit, ok := outputStore.GetDeposit(ctx, big.NewInt(5))

	require.True(t, ok, "deposit does not exist in store")
	require.Equal(t, addr, deposit.Deposit.Owner, "deposit has wrong owner")
	require.Equal(t, big.NewInt(10), deposit.Deposit.Amount, "deposit has wrong amount")
	require.False(t, deposit.Spent, "Deposit is incorrectly marked as spent")
}
