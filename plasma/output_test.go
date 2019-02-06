package plasma

import (
	"github.com/FourthState/plasma-mvp-sidechain/utils"
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

func TestOutputValidation(t *testing.T) {
	type validationCase struct {
		reason string
		Output
	}

	invalidOutputs := []validationCase{
		validationCase{
			reason: "output spending funds to the nil address",
			Output: NewOutput(utils.ZeroAddress, utils.Big1),
		},
		validationCase{
			reason: "output with spending no funds",
			Output: NewOutput(common.HexToAddress("1"), utils.Big0),
		},
	}
	for _, output := range invalidOutputs {
		err := output.ValidateBasic()
		require.Error(t, err, "did not catch: %s", output.reason)
	}

	// valid output
	output := NewOutput(common.HexToAddress("1"), utils.Big1)
	err := output.ValidateBasic()
	require.NoError(t, err, "marked output as invalid: %s", err)
}
