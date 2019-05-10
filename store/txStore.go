package store

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

const (
	accountKey  = "acc"
	hashKey     = "hash"
	positionKey = "pos"
)

/* Wrap plasma transaction with spend information */
type Transaction struct {
	Transaction      plasma.Transaction
	ConfirmationHash []byte
	Spent            []bool
	Spenders         [][32]byte
}

/* Wrap plasma output with spend information */
type Output struct {
	Output  plasma.Output
	Spent   bool
	Spender [32]byte
}

/* Transaction Store */
type TxStore struct {
	kvStore
}

func NewTxStore(ctxKey sdk.StoreKey) TxStore {
	return TxStore{
		kvStore: NewKVStore(ctxKey),
	}
}

func (store TxStore) GetAccount(ctx sdk.Context, addr common.Address) (Account, bool) {
	key := prefixKey(accountKey, addr.Bytes())
	data := store.Get(ctx, key)
	if data == nil {
		return Account{}, false
	}

	var acc Account
	if err := rlp.DecodeBytes(data, &acc); err != nil {
		panic(fmt.Sprintf("transaction store corrupted: %s", err))
	}

	return acc, true
}

func (store TxStore) GetTx(ctx sdk.Context, hash [32]byte) (Transaction, bool) {
	key := prefixKey(hashKey, hash[:])
	data := store.Get(ctx, key)
	if data == nil {
		return Transaction{}, false
	}

	var tx Transaction
	if err := rlp.DecodeBytes(data, &tx); err != nil {
		panic(fmt.Sprintf("transaction store corrupted: %s", err))
	}

	return tx, true
}

// Return the output at the specified position
// along with if it has been spent and what transaction spent it
func (store TxStore) GetUTXO(ctx sdk.Context, pos plasma.Position) (Output, bool) {
	key := prefixKey(positionKey, pos.Bytes())
	data := store.Get(ctx, key)
	var hash [32]byte
	copy(data[:32], hash[:])

	tx, ok := store.GetTx(ctx, hash)
	if !ok {
		return Output{}, ok
	}

	output := Output{
		Output:  tx.Transaction.OutputAt(pos.OutputIndex),
		Spent:   tx.Spent[pos.OutputIndex],
		Spender: tx.Spenders[pos.OutputIndex],
	}

	return output, ok
}

func (store TxStore) HasTx(ctx sdk.Context, hash [32]byte) bool {
	key := prefixKey(hashKey, hash[:])
	return store.Has(ctx, key)
}

func (store TxStore) HasAccount(ctx sdk.Context, addr common.Address) bool {
	key := prefixKey(accountKey, addr.Bytes())
	return store.Has(ctx, key)
}

func (store TxStore) HasUTXO(ctx sdk.Context, pos plasma.Position) bool {
	key := prefixKey(positionKey, pos.Bytes())
	data := store.Get(ctx, key)
	var hash [32]byte
	copy(data[:32], hash[:])

	return store.HasTx(ctx, hash)
}

func (store TxStore) StoreTx(ctx sdk.Context, tx Transaction) {
	data, err := rlp.EncodeToBytes(&tx)
	if err != nil {
		panic(fmt.Sprintf("error marshaling transaction: %s", err))
	}

	key := prefixKey(hashKey, tx.Transaction.TxHash())
	store.Set(ctx, key, data)
}

func (store TxStore) StoreUTXO(ctx sdk.Context, pos plasma.Position, hash [32]byte) {
	data, err := rlp.EncodeToBytes(&pos)
	if err != nil {
		panic(fmt.Sprintf("error marshaling position %s: %s", pos, err))
	}

	key := prefixKey(positionKey, pos.Bytes())
	store.Set(ctx, key, data)
}

func (store TxStore) SpendUTXO(ctx sdk.Context, hash [32]byte, outputIndex int, spenderKey [32]byte) sdk.Result {
	tx, ok := store.GetTx(ctx, hash)
	if !ok {
		return sdk.ErrUnknownRequest(fmt.Sprintf("output with index %x and transaction hash 0x%x does not exist", outputIndex, hash)).Result()
	} else if tx.Spent[outputIndex] {
		return sdk.ErrUnauthorized(fmt.Sprintf("output with index %x and transaction hash 0x%x is already spent", outputIndex, hash)).Result()
	}

	tx.Spent[outputIndex] = true
	tx.Spenders[outputIndex] = spenderKey

	store.StoreTx(ctx, tx)

	return sdk.Result{}
}
