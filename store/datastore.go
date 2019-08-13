package store

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// Name - store name
	DataStoreName = "data"
)

// DataStore
type DataStore struct {
	contextStoreKey sdk.StoreKey
}

// NewDataStore returns a new data store object
func NewDataStore(contextStoreKey sdk.StoreKey) DataStore {
	return DataStore{
		contextStoreKey: contextStoreKey,
	}
}

// Set sets the key value pair in the store
func (ds DataStore) Set(ctx sdk.Context, key []byte, value []byte) {
	store := ctx.KVStore(ds.contextStoreKey)
	store.Set(key, value)
}

// Get returns the value for the provided key from the store
func (ds DataStore) Get(ctx sdk.Context, key []byte) []byte {
	store := ctx.KVStore(ds.contextStoreKey)
	return store.Get(key)
}

// Delete removes the provided key value pair from the store
func (ds DataStore) Delete(ctx sdk.Context, key []byte) {
	store := ctx.KVStore(ds.contextStoreKey)
	store.Delete(key)
}

// Has returns whether the key exists in the store
func (ds DataStore) Has(ctx sdk.Context, key []byte) bool {
	store := ctx.KVStore(ds.contextStoreKey)
	return store.Has(key)
}

// KVStore returns the key value store
func (ds DataStore) KVStore(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(ds.contextStoreKey)
}
