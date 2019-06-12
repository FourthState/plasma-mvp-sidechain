package store

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type kvStore struct {
	contextKey sdk.StoreKey
}

func NewKVStore(contextKey sdk.StoreKey) kvStore {
	return kvStore{
		contextKey: contextKey,
	}
}

func (kv kvStore) Set(ctx sdk.Context, key []byte, value []byte) {
	store := ctx.KVStore(kv.contextKey)
	store.Set(key, value)
}

func (kv kvStore) Get(ctx sdk.Context, key []byte) []byte {
	store := ctx.KVStore(kv.contextKey)
	return store.Get(key)
}

func (kv kvStore) Delete(ctx sdk.Context, key []byte) {
	store := ctx.KVStore(kv.contextKey)
	store.Delete(key)
}

func (kv kvStore) Has(ctx sdk.Context, key []byte) bool {
	store := ctx.KVStore(kv.contextKey)
	return store.Has(key)
}

func (kv kvStore) KVStore(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(kv.contextKey)
}
