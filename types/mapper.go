package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	amino "github.com/tendermint/go-amino"
	crypto "github.com/tendermint/go-crypto"
)

// Maps Position struct to UTXO
// Uses go-amino encoding/decoding library
type UTXOMapper struct {

	// The contextKey used to access the store from the Context.
	contextKey sdk.StoreKey

	// The Amino codec for binary encoding/decoding
	cdc *amino.Codec
}

func NewUTXOMapper(contextKey sdk.StoreKey) UTXOMapper {
	cdc := amino.NewCodec()
	Register(cdc)
	return UTXOMapper{
		contextKey: contextKey,
		cdc:        cdc,
	}

}

// Create and return a sealed utxo mapper. Not sure if necessary
func NewUTXOMapperSealed(contextKey sdk.StoreKey) sealedUTXOMapper {
	cdc := amino.NewCodec()
	um := UTXOMapper{
		contextKey: contextKey,
		cdc:        cdc,
	}
	// Register for amino encoding/decoding
	Register(cdc)

	// make accountMapper's AminoCodec() inaccessible
	return um.Seal()
}

// Registers all interfaces and concrete types necessary for encoding/decoding
func Register(cdc *amino.Codec) {
	crypto.RegisterAmino(cdc)
	cdc.RegisterInterface((*UTXO)(nil), nil)
	cdc.RegisterConcrete(BaseUTXO{}, "types/BaseUTXO", nil)
	cdc.RegisterConcrete(Position{}, "types/Position", nil)
}

// Returns the go-Amino codec.
func (um UTXOMapper) AminoCodec() *amino.Codec {
	return um.cdc
}

// Returns a "sealed" utxoMapper
// the codec is not accessible from a sealedUTXOMapper
func (um UTXOMapper) Seal() sealedUTXOMapper {
	return sealedUTXOMapper{um}
}

// Returns the UTXO corresponding to the go amino encoded Position struct
// Returns nil if no UTXO exists at that position
func (um UTXOMapper) GetUTXO(ctx sdk.Context, position Position) UTXO {
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
func (um UTXOMapper) AddUTXO(ctx sdk.Context, utxo UTXO) {
	position := utxo.GetPosition() 
	pos := um.encodePosition(position)

	store := ctx.KVStore(um.contextKey) 
	bz := um.encodeUTXO(utxo) 			
	store.Set(pos, bz)                  
}

// Deletes UTXO corresponding to the position from mapping
func (um UTXOMapper) DeleteUTXO(ctx sdk.Context, position Position) {
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

//----------------------------------------
// sealedUTXOMapper

type sealedUTXOMapper struct {
	UTXOMapper
}

// There's no way for external modules to mutate the
// sum.utxoMapper.ctx from here, even with reflection
func (sum sealedUTXOMapper) AminoCodec() *amino.Codec {
	panic("utxoMapper is sealed")
}


func (um UTXOMapper) encodeUTXO(utxo UTXO) []byte {
	bz, err := um.cdc.MarshalBinary(utxo)
	if err != nil {
		panic(err)
	}
	return bz
}

func (um UTXOMapper) decodeUTXO(bz []byte) UTXO {
	utxo := &BaseUTXO{}
	err := um.cdc.UnmarshalBinary(bz, utxo)
	if err != nil {
		panic(err)
	}
	return utxo
}

func (um UTXOMapper) encodePosition(pos Position) []byte {
	bz, err := um.cdc.MarshalBinary(pos)
	if err != nil {
		panic(err)
	}
	return bz
}

func (um UTXOMapper) decodePosition(bz []byte) Position {
	pos := Position{}
	err := um.cdc.UnmarshalBinary(bz, pos)
	if err != nil {
		panic(err)
	}
	return pos
}

//----------------------------------------
// UTXOKeeper
// Unnecessary?
// UTXOKeeper manages spending and recieving inputs/outputs
type UTXOKeeper struct {
	um UTXOMapper
}

// NewUTXOKeeper returns a new UTXOKeeper
func NewUTXOKeeper(um UTXOMapper) UTXOKeeper {
	return UTXOKeeper{um: um}
}

// Delete's utxo from utxo store
func (uk UTXOKeeper) SpendUTXO(ctx sdk.Context, addr crypto.Address, position Position) sdk.Error {
	
	utxo := uk.um.GetUTXO(ctx, position) // Get the utxo that should be spent
	// Check to see if utxo exists
	if utxo == nil {
		return sdk.NewError(101, "Unrecognized UTXO. Does not exist.")
	}
	uk.um.DeleteUTXO(ctx, position) // Delete utxo from utxo store
	return nil
}

// Creates a new utxo and adds it to the utxo store
func (uk UTXOKeeper) RecieveUTXO(ctx sdk.Context, addr crypto.Address, denom uint64,
	oldutxos [2]UTXO, oindex uint8) sdk.Error {
	inputAddresses := [2]crypto.Address{oldutxos[0].GetAddress(), oldutxos[1].GetAddress()}
	position := Position{uint64(ctx.BlockHeight()) * 1000, GetTxIndex(ctx), oindex}
	utxo := NewBaseUTXO(addr, inputAddresses, denom, position) 
	uk.um.AddUTXO(ctx, utxo)      // Adds utxo to utxo store
	return nil
}