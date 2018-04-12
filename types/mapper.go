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
	cdc.RegisterInterface((*UTXOHolder)(nil), nil)
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

func (um UTXOMapper) GetUTXO(ctx sdk.Context, addr crypto.Address, position [3]uint) UTXO {
	store := ctx.KVStore(um.contextKey) // Get the utxo store
	bz := store.Get(addr)               // Gets the encoded bytes at the address addr
	// Checks to see if there is a utxo at that address
	if bz == nil {
		return nil
	}
	utxoHolder := um.decodeUTXOHolder(bz)   // Decode the go-amino encoded utxoHolder
	utxo, _ := utxoHolder.GetUTXO(position) // Get the utxo from utxoHolder
	return utxo
}

func (um UTXOMapper) CreateUTXO(ctx sdk.Context, utxo UTXO) {
	addr := utxo.GetAddress()           // Get the address of the utxo
	store := ctx.KVStore(um.contextKey) // Get the utxo store
	bz := store.Get(addr)
	var utxoHolder UTXOHolder
	if bz != nil {
		// Holder already exists
		utxoHolder = um.decodeUTXOHolder(bz)
	} else {
		// Holder does not exist and needs to be created
		utxoHolder = NewUTXOHolder() // change
	}
	utxoHolder.AddUTXO(utxo)             //Add the utxo to the utxoHolder
	bz = um.encodeUTXOHolder(utxoHolder) // Encode the utxoHolder
	store.Set(addr, bz)                  // Add the encoded utxo to the utxo store at address addr
}

func (um UTXOMapper) DestroyUTXO(ctx sdk.Context, utxo UTXO) {
	addr := utxo.GetAddress()
	store := ctx.KVStore(um.contextKey) // Get the utxo store
	bz := store.Get(addr)
	if bz == nil {
		// Add error messages
		return
	}
	utxoHolder := um.decodeUTXOHolder(bz)
	utxoHolder.DeleteUTXO(utxo)
	if utxoHolder.GetLength() == 0 {
		store.Delete(addr)
	} else {
		bz = um.encodeUTXOHolder(utxoHolder)
		store.Set(addr, bz)
	}

}

func (um UTXOMapper) FinalizeUTXO(ctx sdk.Context, addr crypto.Address, denom uint64, position [3]uint, sigs [2]crypto.Signature) sdk.Error {
	store := ctx.KVStore(um.contextKey)
	bz := store.Get(addr)
	if bz == nil {
		return sdk.NewError(100, "No store associated with address")
	}
	utxoHolder := um.decodeUTXOHolder(bz)
	err := utxoHolder.FinalizeUTXO(denom, sigs, position)
	if err == nil {
		return sdk.NewError(100, err.Error())
	}
	bz = um.encodeUTXOHolder(utxoHolder)
	store.Set(addr, bz)
	sigStore := ctx.KVStore(um.sigKey)
	key := []byte{byte(position[0]), byte(position[1]), byte(position[2])}
	sz := store.Get(key)
	if sz != nil {
		return sdk.NewError(100, "Signatures for given position already present in store")
	}
	encodedSigs, encErr := um.cdc.MarshalBinary(sigs)
	if encErr != nil {
		return sdk.NewError(100, "Encoding of signatures failed")
	}
	sigStore.Set(key, encodedSigs)
	return nil
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


func (um UTXOMapper) encodeUTXOHolder(uh UTXOHolder) []byte {
	bz, err := um.cdc.MarshalBinary(uh)
	if err != nil {
		panic(err)
	}
	return bz
}

func (um UTXOMapper) decodeUTXOHolder(bz []byte) UTXOHolder {
	utxoHolder := &BaseUTXOHolder{}
	err := um.cdc.UnmarshalBinary(bz, utxoHolder)
	if err != nil {
		panic(err)
	}
	return utxoHolder
}

//----------------------------------------
// UTXOKeeper

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
	utxo := uk.um.GetUTXO(ctx, addr, position) // Get the utxo that should be spent
	// Check to see if utxo exists
	if utxo == nil {
		return sdk.NewError(101, "Unrecognized UTXO. Does not exist.")
	}
	fmt.Println(utxo)
	uk.um.DestroyUTXO(ctx, utxo) // Delete utxo from utxo store
	return nil
}

// Creates a new utxo and adds it to the utxo store
func (uk UTXOKeeper) RecieveUTXO(ctx sdk.Context, addr crypto.Address, denom uint64) sdk.Error {
	utxo := NewBaseUTXO(addr, denom) // Creates new utxo
	uk.um.CreateUTXO(ctx, utxo)      // Adds utxo to utxo store
	return nil
}

func (uk UTXOKeeper) FinalizeUTXO(ctx sdk.Context, addr crypto.Address, denom uint64, position [3]uint, sigs [2]crypto.Signature) sdk.Error {
	return uk.um.FinalizeUTXO(ctx, addr, denom, position, sigs)
}
