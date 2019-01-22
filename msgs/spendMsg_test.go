package msgs

import (
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestSpendMsgSerialization(t *testing.T) {
	msg := SpendMsg{
		Transaction: plasma.Transaction{
			Input0:  plasma.NewInput(plasma.NewPosition(nil, 1, 0, nil), [65]byte{}, nil),
			Input1:  plasma.NewInput(plasma.NewPosition(utils.Big1, 1, 1, nil), [65]byte{}, nil),
			Output0: plasma.NewOutput(common.HexToAddress("1"), utils.Big1),
			Output1: plasma.NewOutput(common.Address{}, nil),
			Fee:     utils.Big1,
		},
	}

	bytes, err := rlp.EncodeToBytes(&msg)
	require.NoError(t, err, "serialization error")

	recoveredMsg := SpendMsg{}
	err = rlp.DecodeBytes(bytes, &recoveredMsg)
	require.NoError(t, err, "deserialization error")

	require.True(t, reflect.DeepEqual(msg, recoveredMsg), "serialized and deserialized msgs not equal")
}
