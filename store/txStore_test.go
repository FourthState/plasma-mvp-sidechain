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

func TestUTXOSerialization(t *testing.T) {
	keys := [][]byte{[]byte("hamdi")}
	hashes := []byte("allam")
	// no default attributes
	utxo := UTXO{
		InputKeys:        keys,
		SpenderKeys:      keys,
		ConfirmationHash: hashes,
		MerkleHash:       hashes,

		Output:   plasma.NewOutput(common.HexToAddress("1"), utils.Big1),
		Position: plasma.NewPosition(utils.Big1, 0, 0, utils.Big1),
		Spent:    true,
	}

	bytes, err := rlp.EncodeToBytes(&utxo)
	require.NoError(t, err)

	recoveredUTXO := UTXO{}
	err = rlp.DecodeBytes(bytes, &recoveredUTXO)
	require.NoError(t, err)

	require.True(t, reflect.DeepEqual(utxo, recoveredUTXO), "mismatch in serialized and deserialized utxos")
}
