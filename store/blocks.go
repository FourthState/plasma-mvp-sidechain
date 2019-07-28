package store

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
)

// GetBlock - retrive a block at the specified height
func (store DataStore) GetBlock(ctx sdk.Context, blockHeight *big.Int) (Block, bool) {
	key := GetBlockKey(blockHeight)
	data := store.Get(ctx, key)
	if data == nil {
		return Block{}, false
	}

	block := Block{}
	if err := rlp.DecodeBytes(data, &block); err != nil {
		panic(fmt.Sprintf("block store corrupted: %s", err))
	}

	return block, true
}

// StoreBlock will store the plasma block and return the plasma block number
// in which it was stored at.
func (store DataStore) StoreBlock(ctx sdk.Context, tmBlockHeight uint64, block plasma.Block) *big.Int {
	blockHeight := store.PlasmaBlockHeight(ctx)
	if blockHeight == nil {
		// first block to store
		blockHeight = big.NewInt(1)
	} else {
		// increment the height counter
		blockHeight = blockHeight.Add(blockHeight, utils.Big1)
	}

	blockKey := GetBlockKey(blockHeight)
	blockData, err := rlp.EncodeToBytes(&Block{block, tmBlockHeight})
	if err != nil {
		panic(fmt.Sprintf("error rlp encoding block: %s", err))
	}

	// store the block and updated the height counter
	store.Set(ctx, blockKey, blockData)
	store.Set(ctx, GetBlockHeightKey(), blockHeight.Bytes())

	return blockHeight
}

// PlasmaBlockHeight returns the current plasma block height.
func (store DataStore) PlasmaBlockHeight(ctx sdk.Context) *big.Int {
	var plasmaBlockNum *big.Int
	data := store.Get(ctx, GetBlockHeightKey())
	if data == nil {
		return nil
	} else {
		plasmaBlockNum = new(big.Int).SetBytes(data)
	}

	return plasmaBlockNum
}
