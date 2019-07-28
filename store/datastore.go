package store

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DataSTo
type DataStore struct {
	contextStoreKey sdk.StoreKey
}

func NewDataStore(contextStoreKey sdk.StoreKey) DataStore {
	return DataStore{
		contextStoreKey: contextStoreKey,
	}
}

func (ds DataStore) Set(ctx sdk.Context, key []byte, value []byte) {
	store := ctx.KVStore(ds.contextStoreKey)
	store.Set(key, value)
}

func (ds DataStore) Get(ctx sdk.Context, key []byte) []byte {
	store := ctx.KVStore(ds.contextStoreKey)
	return store.Get(key)
}

func (ds DataStore) Delete(ctx sdk.Context, key []byte) {
	store := ctx.KVStore(ds.contextStoreKey)
	store.Delete(key)
}

func (ds DataStore) Has(ctx sdk.Context, key []byte) bool {
	store := ctx.KVStore(ds.contextStoreKey)
	return store.Has(key)
}

func (ds DataStore) KVStore(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(ds.contextStoreKey)
}
