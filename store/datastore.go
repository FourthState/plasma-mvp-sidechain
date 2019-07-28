package store

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DataSTo
type DataStore struct {
	contextStoreKey sdk.StoreKey
}

func NewDataStore(contextStoreKey sdk.StoreKey) kvStore {
	return kvStore{
		contextStoreKey: contextStoreKey,
	}
}

func (kv kvStore) Set(ctx sdk.Context, key []byte, value []byte) {
	store := ctx.KVStore(kv.contextStoreKey)
	store.Set(key, value)
}

func (kv kvStore) Get(ctx sdk.Context, key []byte) []byte {
	store := ctx.KVStore(kv.contextStoreKey)
	return store.Get(key)
}

func (kv kvStore) Delete(ctx sdk.Context, key []byte) {
	store := ctx.KVStore(kv.contextStoreKey)
	store.Delete(key)
}

func (kv kvStore) Has(ctx sdk.Context, key []byte) bool {
	store := ctx.KVStore(kv.contextStoreKey)
	return store.Has(key)
}

func (kv kvStore) KVStore(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(kv.contextStoreKey)
}
