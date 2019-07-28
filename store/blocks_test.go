package store

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/require"
	"math/big"
	"reflect"
	"testing"
)

// test that a block can be serialized and deserialized
func TestBlockSerialization(t *testing.T) {
	// Construct Block
	plasmaBlock := plasma.Block{}
	plasmaBlock.Header[0] = byte(10)
	plasmaBlock.TxnCount = 3
	plasmaBlock.FeeAmount = utils.Big2

	block := Block{
		Block:         plasmaBlock,
		TMBlockHeight: 2,
	}

	// RLP Encode
	bytes, err := rlp.EncodeToBytes(&block)
	require.NoError(t, err)

	// RLP Decode
	recoveredBlock := Block{}
	err = rlp.DecodeBytes(bytes, &recoveredBlock)
	require.NoError(t, err)

	require.True(t, reflect.DeepEqual(block, recoveredBlock), "mismatch in serialized and deserialized block")
}

// test that the plasma block number increments correctly
func TestPlasmaBlockStorage(t *testing.T) {
	ctx, key := setup()
	store := NewDataStore(key)

	height := store.PlasmaBlockHeight(ctx)
	require.Nil(t, height, "non nil height with an empty block store")

	for i := int64(1); i <= 10; i++ {
		// Retrieve nonexistent blocks
		recoveredBlock, ok := store.GetBlock(ctx, big.NewInt(i))
		require.Empty(t, recoveredBlock, "did not return empty struct for nonexistent block")
		require.False(t, ok, "did not return error on nonexistent block")

		nextHeight := store.NextPlasmaBlockHeight(ctx)
		require.Equal(t, nextHeight, big.NewInt(i), "next block height calculated incorrectly")

		// Create and store new block
		var header [32]byte
		hash := crypto.Keccak256([]byte("a plasma block header"))
		copy(header[:], hash[:])

		plasmaBlock := plasma.NewBlock(header, uint16(i*13), big.NewInt(i*123))
		block := Block{plasmaBlock, uint64(i * 1123)}
		blockNum := store.StoreBlock(ctx, uint64(i*1123), plasmaBlock)

		// check the height
		height := store.PlasmaBlockHeight(ctx)
		require.Equal(t, height, blockNum, "block height does not reflect the changes")

		recoveredBlock, ok = store.GetBlock(ctx, blockNum)
		require.True(t, ok, "error when retrieving block")
		require.True(t, reflect.DeepEqual(block, recoveredBlock), fmt.Sprintf("mismatch in stored block and retrieved block, iteration %d", i))
	}
}
