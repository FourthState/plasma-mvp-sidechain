package store

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
)

type PlasmaStore struct {
	KVStore
}

const (
	confirmSigKey     = "confirmSignature"
	blockKey          = "block"
	plasmaBlockNumKey = "plasmaBlockNum"
	plasmaToTmKey     = "plasmatotm"
)

func NewPlasmaStore(ctxKey sdk.StoreKey) PlasmaStore {
	return PlasmaStore{
		KVStore: NewKVStore(ctxKey),
	}
}

func (store PlasmaStore) GetBlock(ctx sdk.Context, blockHeight *big.Int) (plasma.Block, bool) {
	key := prefixKey(blockKey, blockHeight.Bytes())
	data := store.Get(ctx, key)
	if data == nil {
		return plasma.Block{}, false
	}

	block := plasma.Block{}
	if err := rlp.DecodeBytes(data, &block); err != nil {
		panic(fmt.Sprintf("plasma store corrupted: %s", err))
	}

	return block, true
}

// StoreBlock will store the plasma block and return the plasma block number in which it was stored under
func (store PlasmaStore) StoreBlock(ctx sdk.Context, tmBlockHeight *big.Int, block plasma.Block) *big.Int {
	plasmaBlockNum := store.NextPlasmaBlockNum(ctx)

	plasmaBlockKey := prefixKey(blockKey, plasmaBlockNum.Bytes())
	plasmaBlockData, err := rlp.EncodeToBytes(&block)
	if err != nil {
		panic(fmt.Sprintf("error rlp encoding block: %s", err))
	}

	// store the block
	store.Set(ctx, plasmaBlockKey, plasmaBlockData)

	// latest plasma block number
	store.Set(ctx, []byte(plasmaBlockNumKey), plasmaBlockNum.Bytes())

	// plasma block number => tendermint block numper
	store.Set(ctx, prefixKey(plasmaToTmKey, plasmaBlockNum.Bytes()), tmBlockHeight.Bytes())

	return plasmaBlockNum
}

func (store PlasmaStore) StoreConfirmSignatures(ctx sdk.Context, position plasma.Position, confirmSignatures [][65]byte) {
	key := prefixKey(confirmSigKey, position.Bytes())

	var sigs []byte
	sigs = append(sigs, confirmSignatures[0][:]...)
	if len(confirmSignatures) == 2 {
		sigs = append(sigs, confirmSignatures[1][:]...)
	}

	store.Set(ctx, key, sigs)
}

func (store PlasmaStore) NextPlasmaBlockNum(ctx sdk.Context) *big.Int {
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

func (store PlasmaStore) CurrentPlasmaBlockNum(ctx sdk.Context) *big.Int {
	var plasmaBlockNum *big.Int
	data := store.Get(ctx, []byte(plasmaBlockNumKey))
	if data == nil {
		plasmaBlockNum = big.NewInt(1)
	} else {
		plasmaBlockNum = new(big.Int).SetBytes(data)
	}

	return plasmaBlockNum
}
