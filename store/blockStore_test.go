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

// Test that a block can be serialized and deserialized
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

// Test that the plasma block number increments correctly
func TestPlasmaBlockStorage(t *testing.T) {
	ctx, key := setup()
	blockStore := NewBlockStore(key)

	for i := int64(1); i <= 10; i++ {
		// Retrieve nonexistent blocks
		recoveredBlock, ok := blockStore.GetBlock(ctx, big.NewInt(i))
		require.Empty(t, recoveredBlock, "did not return empty struct for nonexistent block")
		require.False(t, ok, "did not return error on nonexistent block")

		// Check increment
		plasmaBlockNum := blockStore.NextPlasmaBlockNum(ctx)
		require.Equal(t, plasmaBlockNum, big.NewInt(i), fmt.Sprintf("plasma block increment returned %d on iteration %d", plasmaBlockNum, i))

		// Create and store new block
		var header [32]byte
		hash := crypto.Keccak256([]byte("a plasma block header"))
		copy(header[:], hash[:])

		plasmaBlock := plasma.NewBlock(header, uint16(i*13), big.NewInt(i*123))
		block := Block{plasmaBlock, uint64(i * 1123)}
		blockNum := blockStore.StoreBlock(ctx, uint64(i*1123), plasmaBlock)
		require.Equal(t, blockNum, plasmaBlockNum, "inconsistency in plasma block number after storing new block")

		recoveredBlock, ok = blockStore.GetBlock(ctx, blockNum)
		require.True(t, ok, "error when retrieving block")
		require.True(t, reflect.DeepEqual(block, recoveredBlock), fmt.Sprintf("mismatch in stored block and retrieved block, iteration %d", i))

		currPlasmaBlock := blockStore.PlasmaBlockHeight(ctx)
		require.Equal(t, currPlasmaBlock, blockNum, fmt.Sprintf("stored block number returned %d, current plasma block returned %d", blockNum, currPlasmaBlock))

	}
}
