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

func TestDepositIndication(t *testing.T) {
	position := NewPosition(big.NewInt(1), 0, 0, big.NewInt(0))
	require.False(t, position.IsDeposit(), "position [1,0,0,0] marked as a deposit")

	position = NewPosition(big.NewInt(0), 0, 0, big.NewInt(1))
	require.True(t, position.IsDeposit(), "position [0,0,0,1] not marked as a deposit")
}
