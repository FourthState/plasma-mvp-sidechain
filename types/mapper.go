package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	amino "github.com/tendermint/go-amino"
	crypto "github.com/tendermint/go-crypto"
)

// Manages utxo's in existence
// Uses go-amino encoding/decoding library
type UTXOMapper struct {

	// The contextKey used to access the store from the Context.
	contextKey sdk.StoreKey

	// The Key required to access store with all confirmSigs. Persists throughout application
	sigKey sdk.StoreKey

	// The Amino codec for binary encoding/decoding of utxo's
	cdc *amino.Codec
}

// NewUTXOMapper returns a new UtxoMapper that
// uses go-Amino to (binary) encode and decode concrete UTXO
func NewUTXOMapper(contextKey sdk.StoreKey, sigKey sdk.StoreKey) UTXOMapper {
	cdc := amino.NewCodec()
	Register(cdc)
	return UTXOMapper{
		contextKey: contextKey,
		sigKey:     sigKey,
		cdc:        cdc,
	}

}

// Create and return a sealed utxo mapper. Not sure if necessary
func NewUTXOMapperSealed(contextKey sdk.StoreKey, sigKey sdk.StoreKey) sealedUTXOMapper {
	cdc := amino.NewCodec()
	um := UTXOMapper{
		contextKey: contextKey,
		sigKey:     sigKey,
		cdc:        cdc,
	}
	// Register for amino encoding/decoding
	Register(cdc)

	// make accountMapper's AminoCodec() inaccessible
	return um.Seal()
}

// Register all crypto interfaces and concrete types necessary
func Register(cdc *amino.Codec) {
	crypto.RegisterAmino(cdc)
	cdc.RegisterInterface((*UTXO)(nil), nil)
	cdc.RegisterConcrete(BaseUTXOHolder{}, "types/BaseUTXOHolder", nil)
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

func (um UTXOMapper) GetUTXO(ctx sdk.Context, position [3]uint) UTXO {
	store := ctx.KVStore(um.contextKey) // Get the utxo store
	pos := position[0] * 100000 + position[1] * 10 + position[2]
	bz := store.Get(pos)               // Gets the encoded bytes at the address addr
	// Checks to see if there is a utxo at that address
	if bz == nil {
		return nil
	}
	utxo := um.decodeUTXO(bz)   // Decode the go-amino encoded utxo
	return utxo
}

func (um UTXOMapper) AddUTXO(ctx sdk.Context, utxo UTXO) {
	position := utxo.GetPosition()       // Get the position of the utxo
	store := ctx.KVStore(um.contextKey) // Get the utxo store
	pos := position[0] * 100000 + position[1] * 10 + position[2]
	bz = um.encodeUTXO(utxo) 			// Encode the utxoHolder
	store.Set(pos, bz)                  // Add the encoded utxo to the utxo store at address addr
}

func (um UTXOMapper) DeleteUTXO(ctx sdk.Context, position [3]uint) {
	store := ctx.KVStore(um.contextKey) // Get the utxo store
	pos := position[0] * 100000 + position[1] * 10 + position[2]
	bz := store.Get(pos)
	if bz == nil {
		// Add error messages
		return
	}
	store.Set(pos, nil)
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


func (um UTXOMapper) encodeUTXO(uh UTXO) []byte {
	bz, err := um.cdc.MarshalBinary(uh)
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
func (uk UTXOKeeper) SpendUTXO(ctx sdk.Context, addr crypto.Address, position [3]uint) sdk.Error {
	
	utxo := uk.um.GetUTXO(ctx, position) // Get the utxo that should be spent
	// Check to see if utxo exists
	if utxo == nil {
		return sdk.NewError(101, "Unrecognized UTXO. Does not exist.")
	}
	uk.um.DeleteUTXO(ctx, utxo) // Delete utxo from utxo store
	return nil
}

// Creates a new utxo and adds it to the utxo store
func (uk UTXOKeeper) RecieveUTXO(ctx sdk.Context, addr crypto.Address, denom uint64,
	 oldutxo UTXO, oindex uint) sdk.Error {

	position := [3]uint{ctx.BlockHeight(), GetTxIndex(ctx), oindex}
	utxo := NewBaseUTXO(addr, oldutxo.GetCSAddress(), nil, oldutxo.GetCSPubKey(), denom, position) 
	uk.um.AddUTXO(ctx, utxo)      // Adds utxo to utxo store
	return nil
}