package db

import (
	types "github.com/FourthState/plasma-mvp-sidechain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	amino "github.com/tendermint/go-amino"
)

// Maps Position struct to UTXO
// Uses go-amino encoding/decoding library
// Implements UTXOMapper
type utxoMapper struct {

	// The contextKey used to access the store from the Context.
	contextKey sdk.StoreKey

	// The Amino codec for binary encoding/decoding
	cdc *amino.Codec
}

func NewUTXOMapper(contextKey sdk.StoreKey, cdc *amino.Codec) types.UTXOMapper {
	return utxoMapper{
		contextKey: contextKey,
		cdc:        cdc,
	}

}

// Returns the UTXO corresponding to the address + go amino encoded Position struct
// Returns nil if no UTXO exists at that position
func (um utxoMapper) GetUTXO(ctx sdk.Context, addr common.Address, position types.Position) types.UTXO {
	store := ctx.KVStore(um.contextKey)
	key := um.getKeyFromAddressAndPosition(addr, position)

	bz := store.Get(key)

	if bz == nil {
		return nil
	}

	utxo := um.decodeUTXO(bz)
	return utxo
}

// Returns all the UTXOs owned by an address.
// Returns empty slice if no UTXO exists for the address.
func (um utxoMapper) GetAllUTXOsForAddress(ctx sdk.Context, addr common.Address) []types.UTXO {
	store := ctx.KVStore(um.contextKey)
	iterator := sdk.KVStorePrefixIterator(store, addr.Bytes())
	utxos := make([]types.UTXO, 0)

	for ; iterator.Valid(); iterator.Next() {
		utxo := um.decodeUTXO(iterator.Value())
		utxos = append(utxos, utxo)
	}
	iterator.Close()

	return utxos
}

// Adds the UTXO to the mapper
func (um utxoMapper) AddUTXO(ctx sdk.Context, utxo types.UTXO) {
	position := utxo.GetPosition()
	address := utxo.GetAddress()
	store := ctx.KVStore(um.contextKey)

	key := um.getKeyFromAddressAndPosition(address, position)
	bz := um.encodeUTXO(utxo)
	store.Set(key, bz)
}

// Deletes UTXO corresponding to address + position from mapping
func (um utxoMapper) DeleteUTXO(ctx sdk.Context, addr common.Address, position types.Position) {
	store := ctx.KVStore(um.contextKey)
	key := um.getKeyFromAddressAndPosition(addr, position)
	store.Delete(key)
}

// (<address> + <encoded position>) forms the unique key that maps to an UTXO.
func (um utxoMapper) getKeyFromAddressAndPosition(address common.Address, position types.Position) []byte {
	pos := um.encodePosition(position)
	key := append(address.Bytes(), pos...)
	return key
}

func (um utxoMapper) encodeUTXO(utxo types.UTXO) []byte {
	bz, err := um.cdc.MarshalBinary(utxo)
	if err != nil {
		panic(err)
	}
	return bz
}

func (um utxoMapper) decodeUTXO(bz []byte) types.UTXO {
	utxo := &types.BaseUTXO{}
	err := um.cdc.UnmarshalBinary(bz, utxo)
	if err != nil {
		panic(err)
	}
	return utxo
}

func (um utxoMapper) encodePosition(pos types.Position) []byte {
	bz, err := um.cdc.MarshalBinary(pos)
	if err != nil {
		panic(err)
	}
	return bz
}
