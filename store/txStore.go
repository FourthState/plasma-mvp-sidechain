package store

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
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
	Spenders         [][]byte
	Position         plasma.Position
}

/* Wrap plasma output with spend information */
type Output struct {
	Output  plasma.Output
	Spent   bool
	Spender []byte
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

// Return the account at the associated address
// Returns nothing if the account does no exist
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

// Return the transaction with the provided transaction hash
func (store TxStore) GetTx(ctx sdk.Context, hash []byte) (Transaction, bool) {
	key := prefixKey(hashKey, hash)
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

// Return the transaction that contains the provided position as an output
func (store TxStore) GetTxWithPosition(ctx sdk.Context, pos plasma.Position) (Transaction, bool) {
	key := prefixKey(positionKey, pos.Bytes())
	hash := store.Get(ctx, key)
	return store.GetTx(ctx, hash)
}

// Return the output at the specified position
// along with if it has been spent and what transaction spent it
func (store TxStore) GetUTXO(ctx sdk.Context, pos plasma.Position) (Output, bool) {
	key := prefixKey(positionKey, pos.Bytes())
	hash := store.Get(ctx, key)

	tx, ok := store.GetTx(ctx, hash)
	if !ok {
		return Output{}, ok
	}

	output := Output{
		Output:  tx.Transaction.Outputs[pos.OutputIndex],
		Spent:   tx.Spent[pos.OutputIndex],
		Spender: tx.Spenders[pos.OutputIndex],
	}

	return output, ok
}

// Checks if a transaction exists using the transaction hash provided
func (store TxStore) HasTx(ctx sdk.Context, hash []byte) bool {
	key := prefixKey(hashKey, hash)
	return store.Has(ctx, key)
}

// Checks if an account exists for the provided address
func (store TxStore) HasAccount(ctx sdk.Context, addr common.Address) bool {
	key := prefixKey(accountKey, addr.Bytes())
	return store.Has(ctx, key)
}

// Checks if the utxo exists using its position
func (store TxStore) HasUTXO(ctx sdk.Context, pos plasma.Position) bool {
	key := prefixKey(positionKey, pos.Bytes())
	hash := store.Get(ctx, key)

	return store.HasTx(ctx, hash)
}

// Store the given Account
func (store TxStore) StoreAccount(ctx sdk.Context, addr common.Address, acc Account) {
	key := prefixKey(accountKey, addr.Bytes())
	data, err := rlp.EncodeToBytes(&acc)
	if err != nil {
		panic(fmt.Sprintf("error marshaling transaction: %s", err))
	}

	store.Set(ctx, key, data)
}

// Add a mapping from transaction hash to transaction
func (store TxStore) StoreTx(ctx sdk.Context, tx Transaction) {
	data, err := rlp.EncodeToBytes(&tx)
	if err != nil {
		panic(fmt.Sprintf("error marshaling transaction: %s", err))
	}

	key := prefixKey(hashKey, tx.Transaction.TxHash())
	store.Set(ctx, key, data)
	store.storeUTXOsWithAccount(ctx, tx)
}

// Add a mapping from position to transaction hash
func (store TxStore) StoreUTXO(ctx sdk.Context, pos plasma.Position, hash []byte) {
	key := prefixKey(positionKey, pos.Bytes())
	store.Set(ctx, key, hash)
}

// Updates Spent, Spender fields and Account associated with this utxo
func (store TxStore) SpendUTXO(ctx sdk.Context, pos plasma.Position, spender []byte) sdk.Result {
	key := prefixKey(positionKey, pos.Bytes())
	hash := store.Get(ctx, key)

	tx, ok := store.GetTx(ctx, hash)
	if !ok {
		return ErrOutputDNE(DefaultCodespace, fmt.Sprintf("output with index %x and transaction hash 0x%x does not exist", pos.OutputIndex, hash)).Result()
	} else if tx.Spent[pos.OutputIndex] {
		return ErrOutputSpent(DefaultCodespace, fmt.Sprintf("output with index %x and transaction hash 0x%x is already spent", pos.OutputIndex, hash)).Result()
	}

	tx.Spent[pos.OutputIndex] = true
	tx.Spenders[pos.OutputIndex] = spender

	store.StoreTx(ctx, tx)
	store.spendUTXOWithAccount(ctx, pos, tx.Transaction)

	return sdk.Result{}
}

/* Helpers */

func (store TxStore) GetUnspentForAccount(ctx sdk.Context, acc Account) (utxos []Output) {
	for _, p := range acc.Unspent {
		utxo, ok := store.GetUTXO(ctx, p)
		if ok {
			utxos = append(utxos, utxo)
		}
	}
	return utxos
}

func (store TxStore) StoreDepositWithAccount(ctx sdk.Context, nonce *big.Int, deposit plasma.Deposit) {
	store.addToAccount(ctx, deposit.Owner, deposit.Amount, plasma.NewPosition(big.NewInt(0), 0, 0, nonce))
}

func (store TxStore) storeUTXOsWithAccount(ctx sdk.Context, tx Transaction) {
	for i, output := range tx.Transaction.Outputs {
		store.addToAccount(ctx, output.Owner, output.Amount, plasma.NewPosition(tx.Position.BlockNum, tx.Position.TxIndex, uint8(i), big.NewInt(0)))
	}
}

func (store TxStore) spendDepositWithAccount(ctx sdk.Context, nonce *big.Int, deposit plasma.Deposit) {
	store.subtractFromAccount(ctx, deposit.Owner, deposit.Amount, plasma.NewPosition(big.NewInt(0), 0, 0, nonce))
}

func (store TxStore) spendUTXOWithAccount(ctx sdk.Context, pos plasma.Position, plasmaTx plasma.Transaction) {
	store.subtractFromAccount(ctx, plasmaTx.Outputs[pos.OutputIndex].Owner, plasmaTx.Outputs[pos.OutputIndex].Amount, pos)
}

func (store TxStore) addToAccount(ctx sdk.Context, addr common.Address, amount *big.Int, pos plasma.Position) {
	acc, ok := store.GetAccount(ctx, addr)
	if !ok {
		acc = Account{big.NewInt(0), make([]plasma.Position, 0), make([]plasma.Position, 0)}
	}

	acc.Balance = new(big.Int).Add(acc.Balance, amount)
	acc.Unspent = append(acc.Unspent, pos)
	store.StoreAccount(ctx, addr, acc)
}

func (store TxStore) subtractFromAccount(ctx sdk.Context, addr common.Address, amount *big.Int, pos plasma.Position) {
	acc, ok := store.GetAccount(ctx, addr)
	if !ok {
		panic(fmt.Sprintf("transaction store has been corrupted"))
	}

	// Update Account
	acc.Balance = new(big.Int).Sub(acc.Balance, amount)
	if acc.Balance.Sign() == -1 {
		panic(fmt.Sprintf("account with address 0x%x has a negative balance", addr))
	}

	removePosition(acc.Unspent, pos)
	acc.Spent = append(acc.Spent, pos)
	store.StoreAccount(ctx, addr, acc)
}

func removePosition(positions []plasma.Position, pos plasma.Position) []plasma.Position {
	for i, p := range positions {
		if p.String() == pos.String() {
			positions = append(positions[:i], positions[i+1:]...)
		}
	}
	return positions
}
