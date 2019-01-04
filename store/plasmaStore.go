package store

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
)

type PlasmaStore struct {
	KVStore
}

const (
	confirmSigPrefix = "confirmSignature"
	deposit          = "deposit"
	blockKey         = "block"
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

func (store PlasmaStore) StoreBlock(ctx sdk.Context, blockHeight *big.Int, block plasma.Block) {
	key := prefixKey(blockKey, blockHeight.Bytes())
	data, err := rlp.EncodeToBytes(&block)
	if err != nil {
		panic(fmt.Sprintf("error rlp encoding block: %s", err))
	}

	store.Set(ctx, key, data)
}

func (store PlasmaStore) StoreConfirmSignatures(ctx sdk.Context, position plasma.Position, confirmSignatures [][65]byte) {
	key := prefixKey(confirmSigPrefix, position.Bytes())

	var sigs []byte
	sigs = append(sigs, confirmSignatures[0][:]...)
	if len(confirmSignatures) == 2 {
		sigs = append(sigs, confirmSignatures[1][:]...)
	}

	store.Set(ctx, key, sigs)
}
