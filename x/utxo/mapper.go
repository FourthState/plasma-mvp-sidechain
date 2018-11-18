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
	AddUTXO(ctx sdk.Context, utxo UTXO)
	DeleteUTXO(ctx sdk.Context, addr []byte, position Position)
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
	key := um.constructKey(addr, position)
	bz := store.Get(key)

	if bz == nil {
		return nil
	}

	utxo := um.decodeUTXO(bz)
	return utxo
}

// Returns all the UTXOs owned by an address.
// Returns empty slice if no UTXO exists for the address.
func (um baseMapper) GetUTXOsForAddress(ctx sdk.Context, addr []byte) []UTXO {
	store := ctx.KVStore(um.contextKey)
	iterator := sdk.KVStorePrefixIterator(store, addr)
	utxos := make([]UTXO, 0)

	for ; iterator.Valid(); iterator.Next() {
		utxo := um.decodeUTXO(iterator.Value())
		utxos = append(utxos, utxo)
	}
	iterator.Close()

	return utxos
}

// Adds the UTXO to the mapper
func (um baseMapper) AddUTXO(ctx sdk.Context, utxo UTXO) {
	position := utxo.GetPosition()
	address := utxo.GetAddress()
	store := ctx.KVStore(um.contextKey)

	key := um.constructKey(address, position)
	bz := um.encodeUTXO(utxo)
	store.Set(key, bz)
}

// Deletes UTXO corresponding to address + position from mapping
func (um baseMapper) DeleteUTXO(ctx sdk.Context, addr []byte, position Position) {
	store := ctx.KVStore(um.contextKey)
	key := um.constructKey(addr, position)
	store.Delete(key)
}

// (<address> + <encoded position>) forms the unique key that maps to an UTXO.
func (um baseMapper) constructKey(address []byte, position Position) []byte {
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
