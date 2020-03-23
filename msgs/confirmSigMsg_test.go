package msgs

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
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

	invalidCases := []confirmSigCase{
		// nil position non-nil sig
		{plasma.NewInput(plasma.Position{}, [65]byte{100}, nil), plasma.NewInput(plasma.Position{}, [65]byte{}, nil)},
		// nil position non-nil confirm sig
		{plasma.NewInput(plasma.Position{}, [65]byte{}, nil), plasma.NewInput(plasma.Position{}, [65]byte{}, [][65]byte{{1}})},
		// invalid position
		{plasma.NewInput(plasma.Position{ BlockNum: utils.Big1, TxIndex: 1, OutputIndex: 1, DepositNonce: utils.Big1}, [65]byte{1}, [][65]byte{{1}}),
			plasma.NewInput(plasma.Position{BlockNum: utils.Big1, TxIndex: 1, OutputIndex: 1, DepositNonce: utils.Big1}, [65]byte{}, nil)},
		// valid deposit empty sig
		{plasma.NewInput(plasma.Position{BlockNum: utils.Big0, TxIndex: 0, OutputIndex: 0, DepositNonce:utils.Big1}, [65]byte{}, nil),
			plasma.NewInput(plasma.Position{BlockNum: utils.Big0, TxIndex: 0, OutputIndex: 0, DepositNonce:utils.Big2}, [65]byte{}, nil)},
		// valid deposit non-nil confirm sig
		{plasma.NewInput(plasma.Position{BlockNum: utils.Big0, TxIndex: 0, OutputIndex: 0, DepositNonce:utils.Big1}, [65]byte{1}, [][65]byte{{1}}),
			plasma.NewInput(plasma.Position{BlockNum: utils.Big0, TxIndex: 0, OutputIndex: 0, DepositNonce:utils.Big2}, [65]byte{1}, [][65]byte{{1}})},
		// valid output empty confirm sig
		{plasma.NewInput(plasma.Position{BlockNum: utils.Big1, TxIndex: 1, OutputIndex: 1, DepositNonce:utils.Big0}, [65]byte{1}, [][65]byte{{1}}),
			plasma.NewInput(plasma.Position{BlockNum: utils.Big1, TxIndex: 2, OutputIndex: 1, DepositNonce:utils.Big0}, [65]byte{1}, nil)},
	}


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