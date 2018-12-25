package store

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	amino "github.com/tendermint/go-amino"
)

// Wrapper around
type UTXO struct {
	InputKeys        [][]byte // keys to retrieve the inputs of this output
	ConfirmationHash [32]byte // confirmation hash of the input transaction

	Output   plasma.Output
	Spent    bool
	Position plasma.Position
}

type UTXOStore struct {
	KVStore
	cdc *amino.Codec
}

func NewUTXOStore(ctxKey sdk.StoreKey) UTXOStore {
	return UTXOStore{
		KVStore: NewKVStore(ctxKey),
	}
}

func (store UTXOStore) GetUTXO(ctx sdk.Context, addr common.Address, pos plasma.Position) (UTXO, bool) {
	key := append(addr.Bytes(), pos.Bytes()...)

	data := store.Get(ctx, key)
	if data == nil {
		return UTXO{}, false
	}

	var utxo UTXO
	err := store.cdc.UnmarshalBinaryBare(data, &utxo)
	if err != nil {
		panic(fmt.Sprintf("utxo store corrupted: %s", err))
	}

	return utxo, true
}

func (store UTXOStore) StoreUTXO(ctx sdk.Context, utxo UTXO) {
	key := append(utxo.Output.Owner.Bytes(), utxo.Position.Bytes()...)
	data, err := store.cdc.MarshalBinaryBare(utxo)
	if err != nil {
		panic(fmt.Sprintf("Error marshaling utxo: %s", err))
	}

	store.Set(ctx, key, data)
}

func (store UTXOStore) SpendUTXO(ctx sdk.Context, addr common.Address, pos plasma.Position, spenderKeys [][]byte) sdk.Error {
	key := append(addr.Bytes(), pos.Bytes()...)
	utxo, ok := store.GetUTXO(ctx, addr, pos)
	if !ok {
		return sdk.ErrUnknownRequest("utxo does not exist")
	}
	if utxo.Spent {
		return sdk.ErrUnauthorized("utxo already marked as spent")
	}

	utxo.Spent = true
	utxo.InputKeys = spenderKeys

	data, err := store.cdc.MarshalBinaryBare(utxo)
	if err != nil {
		panic(fmt.Sprintf("Error marshaling utxo: %s", err))
	}

	store.Set(ctx, key, data)
	return nil
}
