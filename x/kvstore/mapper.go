package kvstore

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type KVStore struct {
	contextKey sdk.StoreKey
}

func NewKVStore(contextKey sdk.StoreKey) KVStore {
	return KVStore{
		contextKey: contextKey,
	}
}

func (kvstore KVStore) Set(ctx sdk.Context, key []byte, value []byte) {
	store := ctx.KVStore(kvstore.contextKey)
	store.Set(key, value)
}

func (kvstore KVStore) Get(ctx sdk.Context, key []byte) []byte {
	store := ctx.KVStore(kvstore.contextKey)
	return store.Get(key)
}

func (kvstore KVStore) Delete(ctx sdk.Context, key []byte) {
	store := ctx.KVStore(kvstore.contextKey)
	store.Delete(key)
}
