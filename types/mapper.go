package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	crypto "github.com/tendermint/go-crypto"
	"utxo"
	"reflect"
	wire "github.com/tendermint/go-wire"
)

// Implements UTXOMapper

type utxoMapper struct {
	
	// The key used to access the store from the Context.
	key sdk.StoreKey

	// The prototypical UTXO concrete type
	proto UTXO

	// The wire codec for binary encoding/decoding of utxo's 
	cdc *wire.Codec 
}

// NewUTXOMapper returns a new utxoMapper that
// uses go-wire to (binary) encode and decode concrete UTXO
func NewUTXOMapper(key sdk.StoreKey, proto UTXO) utxoMapper {
	cdc := wire.NewCodec()
	return utxoMapper {
		key:   key,
		proto: proto,
		cdc:   cdc, 
	}
}

// Create and return a sealed utxo mapper. Not sure if necessary
func NewUTXOMapperSealed(key sdk.StoreKey, proto UTXO) sealedUTXOMapper {
	cdc := wire.NewCodec()
	// um is utxo mapper
	um := utxoMapper {
		key:   key,	
		proto: proto,
		cdc:   cdc,
	}
	// ISSUE: something about register wire here?

	// make accountMapper's WireCodec() inaccessible, return
	return um.Seal()
}

// Returns the go-wire codec.
func (um utxoMapper) WireCodec() *wire.Codec {
	return um.cdc
}

// Returns a "sealed" utxoMapper
// the codec is not accessible from a sealedUTXOMapper
func (um utxoMapper) Seal() sealedUTXOMapper {
	return sealedUTXOMapper{um}
}

// Implements UTXO
func (um utxoMapper) GetUTXO(ctx Context, pubkey crypto.PubKey) UTXO {
	store := ctx.KVStore(um.key)
	bz := store.Get(pubkey)
	if bz == nil {
		return nil
	}
	utxo := um.decodeUTXO(bz)
	return utxo
}

//Implements UTXO
func (um utxoMapper) CreateUTXO(ctx Context, utxo UTXO) {
	pk := utxo.GetPubKey()
	store := ctx.KVStore(um.key)
	bz := um.encodeUTXO(utxo)
	store.Set(pk, bz)
}

//Implements UTXO
func (um utxoMapper) DestroyUTXO(ctx Context, utxo UTXO) {
	pk := utxo.GetPubKey()
	store := ctx.KVStore(um.key)
	bz := store.Get(pk)
	store.Delete(bz)
}

//----------------------------------------
// sealedUTXOMapper

func sealedUTXOMapper struct {
	utxoMapper
}

// There's no way for external modules to mutate the 
// sum.utxoMapper.ctx from here, even with reflection
func (sum sealedUTXOMapper) WireCodec() *wire.Codec {
	panic("utxoMapper is sealed")
}

//----------------------------------------
// misc.

func (um utxoMapper) clonePrototypePtr() interface{} {
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
func (um utxoMapper) clonePrototype() sdk.Account {
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

func (um utxoMapper) encodeUTXO(utxo UTXO) []byte {
	bz, err := um.cdc.MarshalBinary(utxo)
	if err != nil {
		panic(err)
	}
	return bz
}

func (um utxoMapper) decodeUTXO(bz []byte) UTXO {
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