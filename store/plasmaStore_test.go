package store

import (
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestBlockSerialization(t *testing.T) {
	plasmaBlock := plasma.Block{}
	plasmaBlock.Header[0] = byte(10)
	plasmaBlock.TxnCount = 3
	plasmaBlock.FeeAmount = utils.Big2

	block := Block{
		Block:         plasmaBlock,
		TMBlockHeight: 2,
	}

	bytes, err := rlp.EncodeToBytes(&block)
	require.NoError(t, err)

	recoveredBlock := Block{}
	err = rlp.DecodeBytes(bytes, &recoveredBlock)
	require.NoError(t, err)

	require.True(t, reflect.DeepEqual(block, recoveredBlock), "mismatch in serialized and deserialized block")
}
