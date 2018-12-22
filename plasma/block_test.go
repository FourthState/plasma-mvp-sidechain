package plasma

import (
	"crypto/sha256"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/require"
	"math/big"
	"reflect"
	"testing"
)

func TestBlockSerialization(t *testing.T) {
	header := sha256.Sum256([]byte("header"))
	block := NewBlock(header, 10, big.NewInt(1))

	bytes, err := rlp.EncodeToBytes(block)
	require.NoError(t, err, "Error serializing block")

	recoveredBlock := &Block{}
	err = rlp.DecodeBytes(bytes, recoveredBlock)
	require.NoError(t, err, "Error deserializing block")

	require.True(t, reflect.DeepEqual(block, recoveredBlock), "serialized and deserialized objects not deeply equal")
}
