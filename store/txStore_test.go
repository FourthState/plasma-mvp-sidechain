package store

import (
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestTxSerialization(t *testing.T) {
	hashes := []byte("allam")
	var sigs [65]byte

	// no default attributes
	transaction := plasma.Transaction{
		Input0:  plasma.NewInput(plasma.NewPosition(utils.Big1, 15, 1, utils.Big0), sigs, [][65]byte{}),
		Input1:  plasma.NewInput(plasma.NewPosition(utils.Big0, 0, 0, utils.Big1), sigs, [][65]byte{}),
		Output0: plasma.NewOutput(common.HexToAddress("1"), utils.Big1),
		Output1: plasma.NewOutput(common.HexToAddress("2"), utils.Big2),
		Fee:     utils.Big1,
	}

	tx := Transaction{
		Transaction:      transaction,
		Spent:            []bool{false, false},
		Spenders:         [][32]byte{},
		ConfirmationHash: hashes,
	}

	bytes, err := rlp.EncodeToBytes(&tx)
	require.NoError(t, err)

	recoveredTx := Transaction{}
	err = rlp.DecodeBytes(bytes, &recoveredTx)
	require.NoError(t, err)

	require.True(t, reflect.DeepEqual(tx, recoveredTx), "mismatch in serialized and deserialized transactions")
}
