package store

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/rlp"
	"io"
	"math/big"
)

type BlockStore struct {
	kvStore
}

type Block struct {
	plasma.Block
	TMBlockHeight uint64
}

type block struct {
	PlasmaBlock   plasma.Block
	TMBlockHeight uint64
}

// EncodeRLP RLP encodes a Block struct
func (b *Block) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, &block{b.Block, b.TMBlockHeight})
}

// DecodeRLP decodes the byte stream into a Block
func (b *Block) DecodeRLP(s *rlp.Stream) error {
	var block block
	if err := s.Decode(&block); err != nil {
		return err
	}

	b.Block = block.PlasmaBlock
	b.TMBlockHeight = block.TMBlockHeight
	return nil
}

// keys
var (
	blockKey          = []byte{0x0}
	plasmaBlockNumKey = []byte{0x1}
)

// NewBlockStore is a constructor function for BlockStore
func NewBlockStore(ctxKey sdk.StoreKey) BlockStore {
	return BlockStore{
		kvStore: NewKVStore(ctxKey),
	}
}

// GetBlock returns the plasma block at the provided height
func (store BlockStore) GetBlock(ctx sdk.Context, blockHeight *big.Int) (Block, bool) {
	key := prefixKey(blockKey, blockHeight.Bytes())
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

// StoreBlock will store the plasma block and return the plasma block number in which it was stored under
func (store BlockStore) StoreBlock(ctx sdk.Context, tmBlockHeight uint64, block plasma.Block) *big.Int {
	plasmaBlockNum := store.NextPlasmaBlockNum(ctx)

	plasmaBlockKey := prefixKey(blockKey, plasmaBlockNum.Bytes())
	plasmaBlockData, err := rlp.EncodeToBytes(&Block{block, tmBlockHeight})
	if err != nil {
		panic(fmt.Sprintf("error rlp encoding block: %s", err))
	}

	// store the block
	store.Set(ctx, plasmaBlockKey, plasmaBlockData)

	// latest plasma block number
	store.Set(ctx, []byte(plasmaBlockNumKey), plasmaBlockNum.Bytes())

	return plasmaBlockNum
}

// PlasmaBlockHeight returns the current plasma block height
func (store BlockStore) PlasmaBlockHeight(ctx sdk.Context) *big.Int {
	var plasmaBlockNum *big.Int
	data := store.Get(ctx, []byte(plasmaBlockNumKey))
	if data == nil {
		plasmaBlockNum = big.NewInt(1)
	} else {
		plasmaBlockNum = new(big.Int).SetBytes(data)
	}

	return plasmaBlockNum
}

// NextPlasmaBlockNum returns the next plasma block number to be used
func (store BlockStore) NextPlasmaBlockNum(ctx sdk.Context) *big.Int {
	var plasmaBlockNum *big.Int
	data := store.Get(ctx, []byte(plasmaBlockNumKey))
	if data == nil {
		plasmaBlockNum = big.NewInt(1)
	} else {
		plasmaBlockNum = new(big.Int).SetBytes(data)

		// increment the block number
		plasmaBlockNum = plasmaBlockNum.Add(plasmaBlockNum, utils.Big1)
	}

	return plasmaBlockNum
}
