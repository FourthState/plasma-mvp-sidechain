package metadata

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MetadataMapper struct {
	contextKey sdk.StoreKey
}

func NewMetadataMapper(contextKey sdk.StoreKey) MetadataMapper {
	return MetadataMapper{
		contextKey: contextKey,
	}
}

func (mm MetadataMapper) StoreMetadata(ctx sdk.Context, key []byte, metadata []byte) {
	store := ctx.KVStore(mm.contextKey)
	store.Set(key, metadata)
}

func (mm MetadataMapper) GetMetadata(ctx sdk.Context, key []byte) []byte {
	store := ctx.KVStore(mm.contextKey)
	return store.Get(key)
}

func (mm MetadataMapper) DeleteMetadata(ctx sdk.Context, key []byte) {
	store := ctx.KVStore(mm.contextKey)
	store.Delete(key)
}
