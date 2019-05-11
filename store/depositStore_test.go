package store

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/require"
	"math/big"
	"reflect"
	"testing"
)

// Test that the Deposit can be serialized and deserialized without loss of information
func TestDepositSerialization(t *testing.T) {
	// Construct deposit
	plasmaDeposit := plasma.Deposit{
		Owner:       common.BytesToAddress([]byte("an ethereum address")),
		Amount:      big.NewInt(12312310),
		EthBlockNum: big.NewInt(100123123),
	}

	deposit := Deposit{
		Deposit: plasmaDeposit,
		Spent:   true,
		Spender: []byte{},
	}

	// RLP Encode
	bytes, err := rlp.EncodeToBytes(&deposit)
	require.NoError(t, err)

	// RLP Decode
	recoveredDeposit := Deposit{}
	err = rlp.DecodeBytes(bytes, &recoveredDeposit)
	require.NoError(t, err)

	require.True(t, reflect.DeepEqual(deposit, recoveredDeposit), "mismatch in serialized and deserialized deposits")
}

// Test Get, Has, Store, Spend functions
func TestDeposits(t *testing.T) {
	ctx, key := setup()
	depositStore := NewDepositStore(key)

	addr := common.BytesToAddress([]byte("asdfasdf"))
	for i := int64(1); i <= 15; i++ {
		nonce := big.NewInt(i)
		hash := []byte("hash that the deposit was spent in")

		// Retrieve/Spend nonexistent deposits
		exists := depositStore.HasDeposit(ctx, nonce)
		require.False(t, exists, "returned true for nonexistent deposit")
		recoveredDeposit, ok := depositStore.GetDeposit(ctx, nonce)
		require.Empty(t, recoveredDeposit, "did not return empty struct for nonexistent deposit")
		require.False(t, ok, "did not return error on nonexistent deposit")
		res := depositStore.SpendDeposit(ctx, nonce, hash)
		require.Equal(t, res.Code, CodeOutputDNE, "did not return that deposit does not exist")

		// Create and store new deposit
		plasmaDeposit := plasma.NewDeposit(addr, big.NewInt(i*4123), big.NewInt(i*123))
		deposit := Deposit{plasmaDeposit, false, []byte{}}
		depositStore.StoreDeposit(ctx, nonce, deposit)

		exists = depositStore.HasDeposit(ctx, nonce)
		require.True(t, exists, "returned false for deposit that was stored")
		recoveredDeposit, ok = depositStore.GetDeposit(ctx, nonce)
		require.True(t, ok, "error when retrieving deposit")
		require.True(t, reflect.DeepEqual(deposit, recoveredDeposit), fmt.Sprintf("mismatch in stored deposit and retrieved deposit on iteration %d", i))

		// Spend Deposit
		res = depositStore.SpendDeposit(ctx, nonce, hash)
		require.True(t, res.IsOK(), "returned error when spending deposit")
		res = depositStore.SpendDeposit(ctx, nonce, hash)
		require.Equal(t, res.Code, CodeOutputSpent, "allowed output to be spent twice")

		deposit.Spent = true
		deposit.Spender = hash
		recoveredDeposit, ok = depositStore.GetDeposit(ctx, nonce)
		require.True(t, ok, "error when retrieving deposit")
		require.True(t, reflect.DeepEqual(deposit, recoveredDeposit), "mismatch in stored and retrieved deposit")
	}
}
