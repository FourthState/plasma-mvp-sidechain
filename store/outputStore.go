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
	depositKey  = "deposit"
	feeKey      = "fee"
	hashKey     = "hash"
	positionKey = "pos"
)

/* Output Store */
type OutputStore struct {
	kvStore
}

func NewOutputStore(ctxKey sdk.StoreKey) OutputStore {
	return OutputStore{
		kvStore: NewKVStore(ctxKey),
	}
}

// -----------------------------------------------------------------------------
/* Getters */

// GetAccount returns the account at the associated address
func (store OutputStore) GetAccount(ctx sdk.Context, addr common.Address) (Account, bool) {
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

// GetDeposit returns the deposit at the given nonce
func (store OutputStore) GetDeposit(ctx sdk.Context, nonce *big.Int) (Deposit, bool) {
	key := prefixKey(depositKey, nonce.Bytes())
	data := store.Get(ctx, key)
	if data == nil {
		return Deposit{}, false
	}

	var deposit Deposit
	if err := rlp.DecodeBytes(data, &deposit); err != nil {
		panic(fmt.Sprintf("deposit store corrupted: %s", err))
	}

	return deposit, true
}

// GetFee returns the fee at the given position
func (store OutputStore) GetFee(ctx sdk.Context, pos plasma.Position) (Output, bool) {
	key := prefixKey(feeKey, pos.Bytes())
	data := store.Get(ctx, key)
	if data == nil {
		return Output{}, false
	}

	var fee Output
	if err := rlp.DecodeBytes(data, &fee); err != nil {
		panic(fmt.Sprintf("output store corrupted: %s", err))
	}

	return fee, true
}

// GetOutput returns the output at the given position
func (store OutputStore) GetOutput(ctx sdk.Context, pos plasma.Position) (Output, bool) {
	// allow deposits to returned as an output
	if pos.IsDeposit() {
		return store.depositToOutput(ctx, pos.DepositNonce)
	}

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

// GetTx returns the transaction with the provided transaction hash
func (store OutputStore) GetTx(ctx sdk.Context, hash []byte) (Transaction, bool) {
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

// GetTxWithPosition returns the transaction that contains the provided position as an output
func (store OutputStore) GetTxWithPosition(ctx sdk.Context, pos plasma.Position) (Transaction, bool) {
	key := prefixKey(positionKey, pos.Bytes())
	hash := store.Get(ctx, key)
	return store.GetTx(ctx, hash)
}

// -----------------------------------------------------------------------------
/* Has */

// HasAccount returns whether an account at the given address exists
func (store OutputStore) HasAccount(ctx sdk.Context, addr common.Address) bool {
	key := prefixKey(accountKey, addr.Bytes())
	return store.Has(ctx, key)
}

// HasDeposit returns whether a deposit with the given nonce exists
func (store OutputStore) HasDeposit(ctx sdk.Context, nonce *big.Int) bool {
	key := prefixKey(depositKey, nonce.Bytes())
	return store.Has(ctx, key)
}

// HasFee returns whether a fee with the given position exists
func (store OutputStore) HasFee(ctx sdk.Context, pos plasma.Position) bool {
	key := prefixKey(feeKey, pos.Bytes())
	return store.Has(ctx, key)
}

// HasOutput returns whether an output with the given position exists
func (store OutputStore) HasOutput(ctx sdk.Context, pos plasma.Position) bool {
	key := prefixKey(positionKey, pos.Bytes())
	hash := store.Get(ctx, key)

	return store.HasTx(ctx, hash)
}

// HasTx returns whether a transaction with the given transaction hash exists
func (store OutputStore) HasTx(ctx sdk.Context, hash []byte) bool {
	key := prefixKey(hashKey, hash)
	return store.Has(ctx, key)
}

// -----------------------------------------------------------------------------
/* Set */

// SetAccount overwrites the account stored at the given address
func (store OutputStore) setAccount(ctx sdk.Context, addr common.Address, acc Account) {
	key := prefixKey(accountKey, addr.Bytes())
	data, err := rlp.EncodeToBytes(&acc)
	if err != nil {
		panic(fmt.Sprintf("error marshaling transaction: %s", err))
	}

	store.Set(ctx, key, data)
}

// SetDeposit overwrites the deposit stored with the given nonce
func (store OutputStore) setDeposit(ctx sdk.Context, nonce *big.Int, deposit Deposit) {
	data, err := rlp.EncodeToBytes(&deposit)
	if err != nil {
		panic(fmt.Sprintf("error marshaling deposit with nonce %s: %s", nonce, err))
	}

	key := prefixKey(depositKey, nonce.Bytes())
	store.Set(ctx, key, data)
}

// setFee overwrites the fee stored with the given position
func (store OutputStore) setFee(ctx sdk.Context, pos plasma.Position, fee Output) {
	data, err := rlp.EncodeToBytes(&fee)
	if err != nil {
		panic(fmt.Sprintf("error marshaling fee with position %s: %s", pos, err))
	}

	key := prefixKey(feeKey, pos.Bytes())
	store.Set(ctx, key, data)
}

// SetOutput adds a mapping from position to transaction hash
func (store OutputStore) setOutput(ctx sdk.Context, pos plasma.Position, hash []byte) {
	key := prefixKey(positionKey, pos.Bytes())
	store.Set(ctx, key, hash)
}

// SetTx overwrites the mapping from transaction hash to transaction
func (store OutputStore) setTx(ctx sdk.Context, tx Transaction) {
	data, err := rlp.EncodeToBytes(&tx)
	if err != nil {
		panic(fmt.Sprintf("error marshaling transaction: %s", err))
	}

	key := prefixKey(hashKey, tx.Transaction.TxHash())
	store.Set(ctx, key, data)
}

// -----------------------------------------------------------------------------
/* Store */

// StoreDeposit adds an unspent deposit
// Updates the deposit owner's account
func (store OutputStore) StoreDeposit(ctx sdk.Context, nonce *big.Int, deposit plasma.Deposit) {
	store.setDeposit(ctx, nonce, Deposit{deposit, false, make([]byte, 0)})
	store.addToAccount(ctx, deposit.Owner, deposit.Amount, plasma.NewPosition(big.NewInt(0), 0, 0, nonce))
}

// StoreFee adds an unspent fee
// Updates the fee owner's account
func (store OutputStore) StoreFee(ctx sdk.Context, pos plasma.Position, output plasma.Output) {
	store.setFee(ctx, pos, Output{output, false, make([]byte, 0)})
	store.addToAccount(ctx, output.Owner, output.Amount, pos)
}

// StoreTx adds the transaction
// Updates the output owner's accounts
func (store OutputStore) StoreTx(ctx sdk.Context, tx Transaction) {
	store.setTx(ctx, tx)
	for i, output := range tx.Transaction.Outputs {
		store.addToAccount(ctx, output.Owner, output.Amount, plasma.NewPosition(tx.Position.BlockNum, tx.Position.TxIndex, uint8(i), big.NewInt(0)))
		store.setOutput(ctx, plasma.NewPosition(tx.Position.BlockNum, tx.Position.TxIndex, uint8(i), big.NewInt(0)), tx.Transaction.TxHash())
	}
}

// -----------------------------------------------------------------------------
/* Spend */

// SpendDeposit changes the deposit to be spent
// Updates the account of the deposit owner
func (store OutputStore) SpendDeposit(ctx sdk.Context, nonce *big.Int, spender []byte) sdk.Result {
	deposit, ok := store.GetDeposit(ctx, nonce)
	if !ok {
		return ErrOutputDNE(DefaultCodespace, fmt.Sprintf("deposit with nonce %s does not exist", nonce)).Result()
	} else if deposit.Spent {
		return ErrOutputSpent(DefaultCodespace, fmt.Sprintf("deposit with nonce %s is already spent", nonce)).Result()
	}

	deposit.Spent = true
	deposit.Spender = spender

	store.setDeposit(ctx, nonce, deposit)
	store.subtractFromAccount(ctx, deposit.Deposit.Owner, deposit.Deposit.Amount, plasma.NewPosition(big.NewInt(0), 0, 0, nonce))

	return sdk.Result{}
}

// SpendFee changes the fee to be spent
// Updates the account of the fee owner
func (store OutputStore) SpendFee(ctx sdk.Context, pos plasma.Position, spender []byte) sdk.Result {
	fee, ok := store.GetFee(ctx, pos)
	if !ok {
		return ErrOutputDNE(DefaultCodespace, fmt.Sprintf("fee with position %s does not exist", pos)).Result()
	} else if fee.Spent {
		return ErrOutputSpent(DefaultCodespace, fmt.Sprintf("fee with position %s is already spent", pos)).Result()
	}

	fee.Spent = true
	fee.Spender = spender

	store.setFee(ctx, pos, fee)
	store.subtractFromAccount(ctx, fee.Output.Owner, fee.Output.Amount, pos)

	return sdk.Result{}
}

// SpendOutput changes the output to be spent
// Updates the account of the output owner
func (store OutputStore) SpendOutput(ctx sdk.Context, pos plasma.Position, spender []byte) sdk.Result {
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

	store.setTx(ctx, tx)
	store.subtractFromAccount(ctx, tx.Transaction.Outputs[pos.OutputIndex].Owner, tx.Transaction.Outputs[pos.OutputIndex].Amount, pos)

	return sdk.Result{}
}

// -----------------------------------------------------------------------------
/* Helpers */

// GetUnspentForAccount returns the unspent outputs that belong to the given account
func (store OutputStore) GetUnspentForAccount(ctx sdk.Context, acc Account) (utxos []Output) {
	for _, p := range acc.Unspent {
		utxo, ok := store.GetOutput(ctx, p)
		if ok {
			utxos = append(utxos, utxo)
		}
	}
	return utxos
}

// depositToOutput retrieves the deposit with the given nonce, and returns it as an output
func (store OutputStore) depositToOutput(ctx sdk.Context, nonce *big.Int) (Output, bool) {
	deposit, ok := store.GetDeposit(ctx, nonce)
	if !ok {
		return Output{}, ok
	}
	output := Output{
		Output:  plasma.NewOutput(deposit.Deposit.Owner, deposit.Deposit.Amount),
		Spent:   deposit.Spent,
		Spender: deposit.Spender,
	}
	return output, ok
}

// addToAccount adds the passed in amount to the account with the given address
// adds the position provided to the list of unspent positions within the account
func (store OutputStore) addToAccount(ctx sdk.Context, addr common.Address, amount *big.Int, pos plasma.Position) {
	acc, ok := store.GetAccount(ctx, addr)
	if !ok {
		acc = Account{big.NewInt(0), make([]plasma.Position, 0), make([]plasma.Position, 0)}
	}

	acc.Balance = new(big.Int).Add(acc.Balance, amount)
	acc.Unspent = append(acc.Unspent, pos)
	store.setAccount(ctx, addr, acc)
}

// subtractFromAccount subtracts the passed in amount from the account with the given address
// moves the provided position from the unspent list to the spent list
func (store OutputStore) subtractFromAccount(ctx sdk.Context, addr common.Address, amount *big.Int, pos plasma.Position) {
	acc, ok := store.GetAccount(ctx, addr)
	if !ok {
		panic(fmt.Sprintf("transaction store has been corrupted"))
	}

	// Update Account
	acc.Balance = new(big.Int).Sub(acc.Balance, amount)
	if acc.Balance.Sign() == -1 {
		panic(fmt.Sprintf("account with address 0x%x has a negative balance", addr))
	}

	acc.Unspent = removePosition(acc.Unspent, pos)
	acc.Spent = append(acc.Spent, pos)
	store.setAccount(ctx, addr, acc)
}

// helper function to remove a position from the unspent list
func removePosition(positions []plasma.Position, pos plasma.Position) []plasma.Position {
	for i, p := range positions {
		if p.String() == pos.String() {
			positions = append(positions[:i], positions[i+1:]...)
		}
	}
	return positions
}
