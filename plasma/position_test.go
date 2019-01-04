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
