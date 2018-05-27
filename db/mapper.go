package db

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	amino "github.com/tendermint/go-amino"
	types "plasma-mvp-sidechain/types"
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

// Returns the UTXO corresponding to the go amino encoded Position struct
// Returns nil if no UTXO exists at that position
func (um utxoMapper) GetUTXO(ctx sdk.Context, position types.Position) types.UTXO {
	store := ctx.KVStore(um.contextKey)
	pos := um.encodePosition(position)
	bz := store.Get(pos)

	if bz == nil {
		return nil
	}

	utxo := um.decodeUTXO(bz)
	return utxo
}

// Adds the UTXO to the mapper
func (um utxoMapper) AddUTXO(ctx sdk.Context, utxo types.UTXO) {
	position := utxo.GetPosition()
	pos := um.encodePosition(position)

	store := ctx.KVStore(um.contextKey)
	bz := um.encodeUTXO(utxo)
	store.Set(pos, bz)
}

// Deletes UTXO corresponding to the position from mapping
func (um utxoMapper) DeleteUTXO(ctx sdk.Context, position types.Position) {
	store := ctx.KVStore(um.contextKey)
	pos := um.encodePosition(position)
	bz := store.Get(pos)
	// NOTE: For testing, this should never happen
	if bz == nil {
		fmt.Println("Tried to Delete a UTXO that does not exist") // for testing
		return
	}

	store.Delete(pos)
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

func (um utxoMapper) decodePosition(bz []byte) types.Position {
	pos := types.Position{}
	err := um.cdc.UnmarshalBinary(bz, pos)
	if err != nil {
		panic(err)
	}
	return pos
}
