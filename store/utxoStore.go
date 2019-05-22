package store

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"io"
)

type UTXO struct {
	InputKeys        [][]byte `json:"inputKeys"`        // keys to retrieve the inputs of this output
	SpenderKeys      [][]byte `json:"spenderKeys"`      // keys to retrieve the spenders of this output
	ConfirmationHash []byte   `json:"confirmationHash"` // confirmation hash of the input transaction
	MerkleHash       []byte   `json:"merkleHash`        // merkle hash of the input transaction

	Output   plasma.Output   `json:"output"`
	Spent    bool            `json:"spent"`
	Position plasma.Position `json:"position"`
}

type utxo struct {
	InputKeys        [][]byte
	SpenderKeys      [][]byte
	ConfirmationHash []byte
	MerkleHash       []byte
	Output           []byte
	Spent            bool
	Position         []byte
}

func (u *UTXO) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, &utxo{
		u.InputKeys,
		u.SpenderKeys,
		u.ConfirmationHash,
		u.MerkleHash,
		u.Output.Bytes(),
		u.Spent,
		u.Position.Bytes(),
	})
}

func (u *UTXO) DecodeRLP(s *rlp.Stream) error {
	utxo := utxo{}
	if err := s.Decode(&utxo); err != nil {
		return err
	}
	if err := rlp.DecodeBytes(utxo.Output, &u.Output); err != nil {
		return err
	}
	if err := rlp.DecodeBytes(utxo.Position, &u.Position); err != nil {
		return err
	}

	u.InputKeys = utxo.InputKeys
	u.SpenderKeys = utxo.SpenderKeys
	u.ConfirmationHash = utxo.ConfirmationHash
	u.MerkleHash = utxo.MerkleHash
	u.Spent = utxo.Spent

	return nil
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
	KVStore
}

func NewUTXOStore(ctxKey sdk.StoreKey) UTXOStore {
	return UTXOStore{
		KVStore: NewKVStore(ctxKey),
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
		return sdk.ErrUnknownRequest("output does not exist").Result()
	} else if utxo.Spent {
		return sdk.ErrUnauthorized("output already spent").Result()
	}

	utxo.Spent = true
	utxo.SpenderKeys = spenderKeys

	store.StoreUTXO(ctx, utxo)

	return sdk.Result{}
}
