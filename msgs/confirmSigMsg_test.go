package msgs

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestConfirmSigMsgValidate(t *testing.T) {
	type confirmSigCase struct {
		input1  plasma.Input
		input2  plasma.Input
	}

	invalidCases := []confirmSigCase{}
		//{plasma.NewInput(), plasma.NewInput()},
		//{plasma.NewInput(), plasma.NewInput()},
		//{plasma.NewInput(), plasma.NewInput()},

	for i, c := range invalidCases {
		confirmSigMsg := ConfirmSigMsg{
			Input1: c.input1,
			Input2: c.input2,
		}
		require.NotNil(t, confirmSigMsg.ValidateBasic(), fmt.Sprintf("Testcase %d failed", i))
	}
}

func TestConfirmSigMsgSerialization(t *testing.T) {
	msg := ConfirmSigMsg{
		//Input1: plasma.NewInput(),
		//Input2: plasma.NewInput(),
	}

	bytes, err := rlp.EncodeToBytes(&msg)
	require.NoError(t, err, "serialization error")

	tx, err := TxDecoder(bytes)

	require.NoError(t, err, "deserialization error")

	require.True(t, reflect.DeepEqual(msg, tx), "serialized and deserialized msgs not equal")
}