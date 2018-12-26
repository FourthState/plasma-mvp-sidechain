package plasma

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/require"
	"math/big"
	"reflect"
	"testing"
)

func TestOutputSerialization(t *testing.T) {
	output := NewOutput(common.HexToAddress("69"), big.NewInt(10))

	data, err := rlp.EncodeToBytes(&output)
	require.NoError(t, err, "error serializing output")

	recoveredOutput := Output{}
	err = rlp.DecodeBytes(data, &recoveredOutput)
	require.NoError(t, err, "error deserializing output")

	require.True(t, reflect.DeepEqual(output, recoveredOutput), "serialized and deserialized output are not deeply equal")
}
