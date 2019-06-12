package store

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

type UTXO struct {
	InputKeys        [][]byte // keys to retrieve the inputs of this output
	SpenderKeys      [][]byte // keys to retrieve the spenders of this output
	ConfirmationHash []byte   // confirmation hash of the input transaction
	MerkleHash       []byte   // merkle hash of the input transaction

	Output   plasma.Output
	Spent    bool
	Position plasma.Position
}

/* UTXO helper functions */

func (u UTXO) InputAddresses() []common.Address {
	var result []common.Address
	for _, key := range u.InputKeys {
		addr := key[:common.AddressLength]
		result = append(result, common.BytesToAddress(addr[:]))
	}

	return result
}

func (u UTXO) InputPositions() []plasma.Position {
	var result []plasma.Position
	for _, key := range u.InputKeys {
		bytes := key[common.AddressLength:]
		pos := plasma.Position{}
		if err := rlp.DecodeBytes(bytes, &pos); err != nil {
			panic(fmt.Errorf("utxo store corrupted %s", err))
		}

		result = append(result, pos)
	}

	return result
}

func (u UTXO) SpenderAddresses() []common.Address {
	var result []common.Address
	for _, key := range u.SpenderKeys {
		addr := key[:common.AddressLength]
		result = append(result, common.BytesToAddress(addr[:]))
	}

	return result
}

func (u UTXO) SpenderPositions() []plasma.Position {
	var result []plasma.Position
	for _, key := range u.SpenderKeys {
		bytes := key[common.AddressLength:]
		pos := plasma.Position{}
		if err := rlp.DecodeBytes(bytes, &pos); err != nil {
			panic(fmt.Errorf("utxo store corrupted %s", err))
		}

		result = append(result, pos)
	}

	return result
}

/* Store */

type UTXOStore struct {
	kvStore
}

func NewUTXOStore(ctxKey sdk.StoreKey) UTXOStore {
	return UTXOStore{
		kvStore: NewKVStore(ctxKey),
	}
}

func (store UTXOStore) GetUTXOWithKey(ctx sdk.Context, key []byte) (UTXO, bool) {
	data := store.Get(ctx, key)
	if data == nil {
		return UTXO{}, false
	}

	var utxo UTXO
	if err := rlp.DecodeBytes(data, &utxo); err != nil {
		panic(fmt.Sprintf("utxo store corrupted: %s", err))
	}

	return utxo, true
}

func (store UTXOStore) GetUTXO(ctx sdk.Context, addr common.Address, pos plasma.Position) (UTXO, bool) {
	key := GetUTXOStoreKey(addr, pos)
	return store.GetUTXOWithKey(ctx, key)
}

func (store UTXOStore) GetUTXOSet(ctx sdk.Context, addr common.Address) []UTXO {
	var utxos []UTXO
	iter := sdk.KVStorePrefixIterator(store.KVStore(ctx), addr.Bytes())
	for ; iter.Valid(); iter.Next() {
		key := iter.Key()
		utxo, ok := store.GetUTXOWithKey(ctx, key)
		if !ok {
			panic(fmt.Sprintf("utxo store corrupted: non-existent key in set: 0x%x", key))
		}

		utxos = append(utxos, utxo)
	}

	return utxos
}

func (store UTXOStore) HasUTXO(ctx sdk.Context, addr common.Address, pos plasma.Position) bool {
	key := GetUTXOStoreKey(addr, pos)
	return store.Has(ctx, key)
}

func (store UTXOStore) StoreUTXO(ctx sdk.Context, utxo UTXO) {
	key := GetStoreKey(utxo)
	data, err := rlp.EncodeToBytes(&utxo)
	if err != nil {
		panic(fmt.Sprintf("Error marshaling utxo: %s", err))
	}

	store.Set(ctx, key, data)
}

func (store UTXOStore) SpendUTXO(ctx sdk.Context, addr common.Address, pos plasma.Position, spenderKeys [][]byte) sdk.Result {
	utxo, ok := store.GetUTXO(ctx, addr, pos)
	if !ok {
		return sdk.ErrUnknownRequest(fmt.Sprintf("output with address 0x%x and position %v does not exist", addr, pos)).Result()
	} else if utxo.Spent {
		return sdk.ErrUnauthorized(fmt.Sprintf("output with address 0x%x and position %v is already spent", addr, pos)).Result()
	}

	utxo.Spent = true
	utxo.SpenderKeys = spenderKeys

	store.StoreUTXO(ctx, utxo)

	return sdk.Result{}
}
