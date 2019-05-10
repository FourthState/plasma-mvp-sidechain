package plasma

import (
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/require"
	"math/big"
	"reflect"
	"testing"
)

func TestPositionSerialization(t *testing.T) {
	position := NewPosition(big.NewInt(1), 6, 9, big.NewInt(0))
	bytes, err := rlp.EncodeToBytes(&position)
	require.NoError(t, err, "error serializing position")

	recoveredPosition := Position{}
	err = rlp.DecodeBytes(bytes, &recoveredPosition)
	require.NoError(t, err, "error deserializing position")

	require.True(t, reflect.DeepEqual(&position, &recoveredPosition), "serialized and deserialized position not deeply equal")
}

func TestPositionDepositIndication(t *testing.T) {
	position := NewPosition(big.NewInt(1), 0, 0, big.NewInt(0))
	require.False(t, position.IsDeposit(), "position [1,0,0,0] marked as a deposit")

	position = NewPosition(big.NewInt(0), 0, 0, big.NewInt(1))
	require.True(t, position.IsDeposit(), "position [0,0,0,1] not marked as a deposit")
}

func TestPositionFeeIndication(t *testing.T) {
	position := NewPosition(big.NewInt(2), 0, 0, big.NewInt(0))
	require.False(t, position.IsFee(), "position [2, 0, 0, 0] marked as a fee")

	position = NewPosition(big.NewInt(1), 65535, 0, big.NewInt(0))
	require.True(t, position.IsFee(), "position [1, 65535, 0, 0] not marked as a fee")
}

func TestPositionFromString(t *testing.T) {
	posStr := "(1.0.0.0)"
	pos, err := FromPositionString(posStr)
	require.NoError(t, err, "error converting from string")
	require.Equal(t, pos.String(), posStr)

	posStr = "(0.0.5.0)"
	pos, err = FromPositionString(posStr)
	require.Error(t, err, "converted with an invalid output index")

	posStr = "(1.0.0.1)"
	pos, err = FromPositionString(posStr)
	require.Error(t, err, "converted with nonce and chain positions specified")
}

func TestPositionValidation(t *testing.T) {
	posStr := "(1.0.0.0)"
	pos, _ := FromPositionString(posStr)
	require.NoError(t, pos.ValidateBasic(), "valid position marked as an error")

	// specify both deposit nonce and chain position
	cases := []string{
		// mutual exclusivity between deposit nonce and chain position required
		"(1.0.0.5)",
		"(0.1.0.5)",
		"(0.0.1.5)",
		// chain position with block number zero
		"(0.1.1.0)",
		// invalid output index
		"(1.1.3.0)",
		// nil position is not a valid position
		"(0.0.0.0)",
	}
	for _, posStr := range cases {
		pos, _ = FromPositionString(posStr)
		require.Errorf(t, pos.ValidateBasic(), "invalid position: %s", posStr)
	}
}

func TestPositionFromKey(t *testing.T) {
	utxoKey := big.NewInt(133*blockIndexFactor + 14*txIndexFactor)
	utxo := FromExitKey(utxoKey, false)
	require.Equal(t, utxo, NewPosition(big.NewInt(133), 14, 0, big.NewInt(0)), "error retrieving correct position from exit key")

	depositKey := big.NewInt(10)
	deposit := FromExitKey(depositKey, true)
	require.Equal(t, deposit, NewPosition(big.NewInt(0), 0, 0, depositKey), "error retrieving correct position from deposit exit key")
}
