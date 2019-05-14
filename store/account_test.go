package store

import (
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/require"
	"math/big"
	"reflect"
	"testing"
)

// Test that an account can be serialized and deserialized
func TestAccountSerialization(t *testing.T) {
	// Construct Account
	acc := Account{
		Balance: big.NewInt(234578),
		Unspent: []plasma.Position{GetPosition("(8745.1239.1.0)"), GetPosition("(23409.12456.0.0)"), GetPosition("(894301.1.1.0)"), GetPosition("(0.0.0.540124)")},
		Spent:   []plasma.Position{GetPosition("0.0.0.3"), GetPosition("7.734.1.3")},
	}

	// RLP Encode
	bytes, err := rlp.EncodeToBytes(&acc)
	require.NoError(t, err)

	// RLP Decode
	recoveredAcc := Account{}
	err = rlp.DecodeBytes(bytes, &recoveredAcc)
	require.NoError(t, err)

	require.True(t, reflect.DeepEqual(acc, recoveredAcc), "mismatch in serialized and deserialized account")
}
