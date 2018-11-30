package utxo

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	amino "github.com/tendermint/go-amino"
)

// Mapper stores and retrieves UTXO's from stores
// retrieved from the context.
type Mapper interface {
	GetUTXO(ctx sdk.Context, addr []byte, position Position) UTXO
	GetUTXOsForAddress(ctx sdk.Context, addr []byte) []UTXO
	ConstructKey(addr []byte, position Position) []byte
	ReceiveUTXO(sdk.Context, UTXO)
	ValidateUTXO(sdk.Context, UTXO) sdk.Error
	InvalidateUTXO(sdk.Context, UTXO)
	SpendUTXO(ctx sdk.Context, addr []byte, position Position, spenderKeys [][]byte) sdk.Error
}

// Maps Address+Position to UTXO
// Uses go-amino encoding/decoding library
// Implements Mapper
type baseMapper struct {
	// The contextKey used to access the store from the Context.
	contextKey sdk.StoreKey

	// The Amino codec for binary encoding/decoding
	cdc *amino.Codec
}

func NewBaseMapper(contextKey sdk.StoreKey, cdc *amino.Codec) Mapper {
	return baseMapper{
		contextKey: contextKey,
		cdc:        cdc,
	}
}

// Returns the UTXO corresponding to the address + go amino encoded Position struct
// Returns nil if no UTXO exists at that position
func (um baseMapper) GetUTXO(ctx sdk.Context, addr []byte, position Position) UTXO {
	store := ctx.KVStore(um.contextKey)
	key := um.ConstructKey(addr, position)
	bz := store.Get(key)

	if bz == nil {
		return UTXO{}
	}

	utxo := um.decodeUTXO(bz)
	return utxo
}

// Returns all the valid UTXOs owned by an address.
// Returns empty slice if no UTXO exists for the address.
func (um baseMapper) GetUTXOsForAddress(ctx sdk.Context, addr []byte) []UTXO {
	store := ctx.KVStore(um.contextKey)
	iterator := sdk.KVStorePrefixIterator(store, addr)
	utxos := make([]UTXO, 0)

	for ; iterator.Valid(); iterator.Next() {
		utxo := um.decodeUTXO(iterator.Value())
		// Only append if UTXO is valid
		if utxo.Valid {
			utxos = append(utxos, utxo)
		}
	}
	iterator.Close()

	return utxos
}

// Adds the UTXO to the mapper
func (um baseMapper) ReceiveUTXO(ctx sdk.Context, utxo UTXO) {
	store := ctx.KVStore(um.contextKey)

	key := utxo.StoreKey(um.cdc)
	bz := um.encodeUTXO(utxo)
	store.Set(key, bz)
}

// Deletes UTXO corresponding to address + position from mapping
func (um baseMapper) SpendUTXO(ctx sdk.Context, addr []byte, position Position, spenderKeys [][]byte) sdk.Error {
	store := ctx.KVStore(um.contextKey)
	key := um.ConstructKey(addr, position)
	utxo := um.GetUTXO(ctx, addr, position)
	if !utxo.Valid {
		return sdk.ErrUnauthorized("UTXO is not valid for spend")
	}
	utxo.Valid = false
	utxo.SpenderKeys = spenderKeys
	encodedUTXO := um.encodeUTXO(utxo)
	store.Set(key, encodedUTXO)
	return nil
}

// Validates UTXO only if it not spent already
func (um baseMapper) ValidateUTXO(ctx sdk.Context, utxo UTXO) sdk.Error {
	store := ctx.KVStore(um.contextKey)
	key := utxo.StoreKey(um.cdc)
	if utxo.SpenderKeys != nil {
		return sdk.ErrUnauthorized("Cannot validate spent UTXO")
	}
	utxo.Valid = true
	encodedUTXO := um.encodeUTXO(utxo)
	store.Set(key, encodedUTXO)
	return nil
}

// Invalidates UTXO
func (um baseMapper) InvalidateUTXO(ctx sdk.Context, utxo UTXO) {
	store := ctx.KVStore(um.contextKey)
	key := utxo.StoreKey(um.cdc)
	utxo.Valid = false
	encodedUTXO := um.encodeUTXO(utxo)
	store.Set(key, encodedUTXO)
}

// (<address> + <encoded position>) forms the unique key that maps to an UTXO.
func (um baseMapper) ConstructKey(address []byte, position Position) []byte {
	posBytes, err := um.cdc.MarshalBinaryBare(position)
	if err != nil {
		panic(err)
	}
	key := append(address, posBytes...)
	return key
}

func (um baseMapper) encodeUTXO(utxo UTXO) []byte {
	bz, err := um.cdc.MarshalBinaryBare(utxo)
	if err != nil {
		panic(err)
	}
	return bz
}

func (um baseMapper) decodeUTXO(bz []byte) (utxo UTXO) {
	err := um.cdc.UnmarshalBinaryBare(bz, &utxo)
	if err != nil {
		panic(err)
	}
	return utxo
}
