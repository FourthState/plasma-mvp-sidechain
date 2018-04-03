package types

import (
	"reflect"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	crypto "github.com/tendermint/go-crypto"
	amino "github.com/tendermint/go-amino"
)

// Manages utxo's in existence
// Uses go-amino encoding/decoding library
// Does not need to be changed to RLP
type UTXOMapper struct {
	
	// The key used to access the store from the Context.
	key sdk.StoreKey

	// The prototypical UTXO concrete type
	proto UTXOHolder

	// The Amino codec for binary encoding/decoding of utxo's 
	cdc *amino.Codec 
}

// NewUTXOMapper returns a new UtxoMapper that
// uses go-Amino to (binary) encode and decode concrete UTXO
func NewUTXOMapper(key sdk.StoreKey, proto UTXOHolder) UTXOMapper {
	cdc := amino.NewCodec()
	Register(cdc)
	return UTXOMapper {
		key:   key,
		proto: proto,
		cdc:   cdc, 
	}

}

// Create and return a sealed utxo mapper. Not sure if necessary
func NewUTXOMapperSealed(key sdk.StoreKey, proto UTXOHolder) sealedUTXOMapper {
	cdc := amino.NewCodec()
	um := UTXOMapper {
		key:   key,	
		proto: proto,
		cdc:   cdc,
	}
	// Register for amino encoding/decoding
	Register(cdc)

	// make accountMapper's AminoCodec() inaccessible
	return um.Seal()
}

// Register all crypto interfaces and concrete types necessary
func Register(cdc *amino.Codec) {
	cdc.RegisterConcrete(crypto.PubKey{}, "go-crypto/PubKey", nil)
	cdc.RegisterConcrete(crypto.PrivKey{}, "go-crypto/PrivKey", nil)
	cdc.RegisterConcrete(crypto.Signature{}, "go-crypto/Signature", nil)
	cdc.RegisterConcrete(sdk.StdSignature{}, "sdk/StdSignature", nil)
	cdc.RegisterInterface((*crypto.PubKeyInner)(nil), nil)
	cdc.RegisterConcrete(crypto.PubKeyEd25519{}, "go-crypto/PubKeyEd25519", nil)
	cdc.RegisterConcrete(crypto.SignatureEd25519{}, "go-crypto/SignatureEd25519", nil)
	cdc.RegisterInterface((*crypto.SignatureInner)(nil), nil)
	cdc.RegisterInterface((*UTXOHolder)(nil), nil)
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
	store := ctx.KVStore(um.key) // Get the utxo store
	bz := store.Get(addr) // Gets the encoded bytes at the address addr
	// Checks to see if there is a utxo at that address
	if bz == nil {
		return nil
	}
	utxoHolder := um.decodeUTXOHolder(bz) // Decode the go-amino encoded utxoHolder
	utxo, _ := utxoHolder.GetUTXO(position) // Get the utxo from utxoHolder
	return utxo
}

func (um UTXOMapper) CreateUTXO(ctx sdk.Context, utxo UTXO) {
	addr := utxo.GetAddress() // Get the address of the utxo
	store := ctx.KVStore(um.key) // Get the utxo store 
	bz := store.Get(addr)
	var utxoHolder UTXOHolder
	if bz != nil {
		// Holder already exists
		utxoHolder = um.decodeUTXOHolder(bz)
	} else {
		// Holder does not exist and needs to be created
		utxoHolder = NewUTXOHolder() // change
	}
	utxoHolder.AddUTXO(utxo) //Add the utxo to the utxoHolder
	bz = um.encodeUTXOHolder(utxoHolder) // Encode the utxoHolder
	store.Set(addr, bz) // Add the encoded utxo to the utxo store at address addr
}

func (um UTXOMapper) DestroyUTXO(ctx sdk.Context, utxo UTXO) {
	addr := utxo.GetAddress()
	store := ctx.KVStore(um.key) // Get the utxo store
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

func (um UTXOMapper) FinalizeUTXO(ctx sdk.Context, addr crypto.Address, denom uint64, position [3]uint, sigs []sdk.StdSignature) sdk.Error {
	store := ctx.KVStore(um.key)
	bz := store.Get(addr)
	if (bz == nil) {
		return sdk.NewError(100, "No store associated with address")
	}
	utxoHolder := um.decodeUTXOHolder(bz)
	err := utxoHolder.FinalizeUTXO(denom, sigs, position)
	if err == nil {
		return sdk.NewError(100, err.Error())
	}
	bz = um.encodeUTXOHolder(utxoHolder)
	store.Set(addr, bz)
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

//----------------------------------------
// misc.

func (um UTXOMapper) clonePrototypePtr() interface{} {
	protoRt := reflect.TypeOf(um.proto)
	if protoRt.Kind() == reflect.Ptr {
		protoErt := protoRt.Elem()
		if protoErt.Kind() != reflect.Struct {
			panic("utxoMapper requires a struct proto UTXO, or a pointer to one")
		}
		protoRv := reflect.New(protoErt)
		return protoRv.Interface()
	} else {
		protoRv := reflect.New(protoRt)
		return protoRv.Interface()
	}
}

// Creates a new struct (or pointer to struct) from um.proto.
func (um UTXOMapper) clonePrototype() UTXOHolder {
	protoRt := reflect.TypeOf(um.proto)
	if protoRt.Kind() == reflect.Ptr {
		protoCrt := protoRt.Elem()
		if protoCrt.Kind() != reflect.Struct {
			panic("utxoMapper requires a struct proto UTXOHolder, or a pointer to one")
		}
		protoRv := reflect.New(protoCrt)
		clone, ok := protoRv.Interface().(UTXOHolder)
		if !ok {
			panic(fmt.Sprintf("utxoMapper requires a proto UTXO, but %v doesn't implement UTXO", protoRt))
		}
		return clone
	} else {
		protoRv := reflect.New(protoRt).Elem()
		clone, ok := protoRv.Interface().(UTXOHolder)
		if !ok {
			panic(fmt.Sprintf("utxoMapper requires a proto UTXOHolder, but %v doesn't implement UTXO", protoRt))
		}
		return clone
	}
}

func (um UTXOMapper) encodeUTXOHolder(uh UTXOHolder) []byte {
	bz, err := um.cdc.MarshalBinary(uh)
	if err != nil {
		panic(err)
	}
	return bz
}

func (um UTXOMapper) decodeUTXOHolder(bz []byte) UTXOHolder {
	uhPtr := um.clonePrototypePtr()
	err := um.cdc.UnmarshalBinary(bz, uhPtr)
	if err != nil {
		panic(err)
	}
	if reflect.ValueOf(um.proto).Kind() == reflect.Ptr {
		return reflect.ValueOf(uhPtr).Interface().(UTXOHolder)
	} else {
		return reflect.ValueOf(uhPtr).Elem().Interface().(UTXOHolder)
	}
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
	uk.um.DestroyUTXO(ctx, utxo) // Delete utxo from utxo store
	return nil
}

// Creates a new utxo and adds it to the utxo store
func (uk UTXOKeeper) RecieveUTXO(ctx sdk.Context, addr crypto.Address, denom uint64) sdk.Error {
	utxo := NewBaseUTXO(addr, denom) // Creates new utxo
	uk.um.CreateUTXO(ctx, utxo) // Adds utxo to utxo store
	return nil
}

func (uk UTXOKeeper) FinalizeUTXO(ctx sdk.Context, addr crypto.Address, denom uint64, position [3]uint, sigs []sdk.StdSignature) sdk.Error {
	return uk.um.FinalizeUTXO(ctx, addr, denom, position, sigs)
}
