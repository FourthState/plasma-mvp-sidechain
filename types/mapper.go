package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	crypto "github.com/tendermint/go-crypto"
	"reflect"
	amino "github.com/tendermint/go-amino"
	"errors"
	"fmt"
)

// Implements UTXOMapper

type UtxoMapper struct {
	
	// The key used to access the store from the Context.
	key sdk.StoreKey

	// The prototypical UTXO concrete type
	proto UTXO

	// The Amino codec for binary encoding/decoding of utxo's 
	cdc *amino.Codec 
}

// NewUTXOMapper returns a new utxoMapper that
// uses go-Amino to (binary) encode and decode concrete UTXO
func NewUTXOMapper(key sdk.StoreKey, proto UTXO) UtxoMapper {
	cdc := amino.NewCodec()
	return UtxoMapper {
		key:   key,
		proto: proto,
		cdc:   cdc, 
	}
}

// Create and return a sealed utxo mapper. Not sure if necessary
func NewUTXOMapperSealed(key sdk.StoreKey, proto UTXO) sealedUTXOMapper {
	cdc := amino.NewCodec()
	// um is utxo mapper
	um := UtxoMapper {
		key:   key,	
		proto: proto,
		cdc:   cdc,
	}
	// ISSUE: something about register Amino here?

	// make accountMapper's AminoCodec() inaccessible, return
	return um.Seal()
}

// Returns the go-Amino codec.
func (um UtxoMapper) AminoCodec() *amino.Codec {
	return um.cdc
}

// Returns a "sealed" utxoMapper
// the codec is not accessible from a sealedUTXOMapper
func (um UtxoMapper) Seal() sealedUTXOMapper {
	return sealedUTXOMapper{um}
}

// Implements UTXO
func (um UtxoMapper) GetUTXO(ctx sdk.Context, addr crypto.Address) UTXO {
	store := ctx.KVStore(um.key)
	bz := store.Get(addr)
	if bz == nil {
		return nil
	}
	utxo := um.decodeUTXO(bz)
	return utxo
}

//Implements UTXO
func (um UtxoMapper) CreateUTXO(ctx sdk.Context, utxo UTXO) {
	addr := utxo.GetAddress()
	store := ctx.KVStore(um.key)
	bz := um.encodeUTXO(utxo)
	store.Set(addr, bz)
}

//Implements UTXO
func (um UtxoMapper) DestroyUTXO(ctx sdk.Context, utxo UTXO) {
	addr := utxo.GetAddress()
	store := ctx.KVStore(um.key)
	bz := store.Get(addr)
	store.Delete(bz)
}

//----------------------------------------
// sealedUTXOMapper

type sealedUTXOMapper struct {
	UtxoMapper
}

// There's no way for external modules to mutate the 
// sum.utxoMapper.ctx from here, even with reflection
func (sum sealedUTXOMapper) AminoCodec() *amino.Codec {
	panic("utxoMapper is sealed")
}

//----------------------------------------
// misc.

func (um UtxoMapper) clonePrototypePtr() interface{} {
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
func (um UtxoMapper) clonePrototype() UTXO {
	protoRt := reflect.TypeOf(um.proto)
	if protoRt.Kind() == reflect.Ptr {
		protoCrt := protoRt.Elem()
		if protoCrt.Kind() != reflect.Struct {
			panic("utxoMapper requires a struct proto UTXO, or a pointer to one")
		}
		protoRv := reflect.New(protoCrt)
		clone, ok := protoRv.Interface().(UTXO)
		if !ok {
			panic(fmt.Sprintf("utxoMapper requires a proto UTXO, but %v doesn't implement UTXO", protoRt))
		}
		return clone
	} else {
		protoRv := reflect.New(protoRt).Elem()
		clone, ok := protoRv.Interface().(UTXO)
		if !ok {
			panic(fmt.Sprintf("utxoMapper requires a proto UTXO, but %v doesn't implement UTXO", protoRt))
		}
		return clone
	}
}

func (um UtxoMapper) encodeUTXO(utxo UTXO) []byte {
	bz, err := um.cdc.MarshalBinary(utxo)
	if err != nil {
		panic(err)
	}
	return bz
}

func (um UtxoMapper) decodeUTXO(bz []byte) UTXO {
	utxoPtr := um.clonePrototypePtr()
	err := um.cdc.UnmarshalBinary(bz, utxoPtr)
	if err != nil {
		panic(err)
	}
	if reflect.ValueOf(um.proto).Kind() == reflect.Ptr {
		return reflect.ValueOf(utxoPtr).Interface().(UTXO)
	} else {
		return reflect.ValueOf(utxoPtr).Elem().Interface().(UTXO)
	}
}

//----------------------------------------
// UTXOKeeper

// UTXOKeeper manages spending and recieving inputs/outputs
type UTXOKeeper struct {
	um UtxoMapper
}

// NewUTXOKeeper returns a new UTXOKeeper
func NewUTXOKeeper(um UtxoMapper) UTXOKeeper {
	return UTXOKeeper{um: um}
}

//May need to add check for valid context
func (uk UTXOKeeper) SpendUTXO(ctx sdk.Context, addr crypto.Address) error {
	utxo := uk.um.GetUTXO(ctx, addr)
	if utxo == nil {
		return errors.New("UTXO does not exist")
	}
	uk.um.DestroyUTXO(ctx, utxo)
	return nil
}

//May need more checks for error
func (uk UTXOKeeper) RecieveUTXO(ctx sdk.Context, addr crypto.Address, denom uint64) error {
	utxo := NewBaseUTXO(addr, denom)
	uk.um.CreateUTXO(ctx, utxo)
	return nil
}
