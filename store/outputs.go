package store

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"math/big"
)

// -----------------------------------------------------------------------------
/* Getters */

// GetWallet returns the wallet at the associated address.
func (store DataStore) GetWallet(ctx sdk.Context, addr common.Address) (Wallet, bool) {
	key := GetWalletKey(addr)
	data := store.Get(ctx, key)
	if data == nil {
		return Wallet{}, false
	}

	var wallet Wallet
	if err := rlp.DecodeBytes(data, &wallet); err != nil {
		panic(fmt.Sprintf("transaction store corrupted: %s", err))
	}

	return wallet, true
}

// GetDeposit returns the deposit at the given nonce.
func (store DataStore) GetDeposit(ctx sdk.Context, nonce *big.Int) (Deposit, bool) {
	key := GetDepositKey(nonce)
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

// GetFee returns the fee at the given position.
func (store DataStore) GetFee(ctx sdk.Context, pos plasma.Position) (Output, bool) {
	key := GetFeeKey(pos)
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

// GetOutput returns the output at the given position.
func (store DataStore) GetOutput(ctx sdk.Context, pos plasma.Position) (Output, bool) {
	// allow deposits/fees to returned as an output
	if pos.IsDeposit() {
		return store.depositToOutput(ctx, pos.DepositNonce)
	}

	if pos.IsFee() {
		fee, ok := store.GetFee(ctx, pos)
		if !ok {
			return Output{}, ok
		}
		return fee, ok
	}

	key := GetOutputKey(pos)
	hash := store.Get(ctx, key)

	tx, ok := store.GetTx(ctx, hash)
	if !ok {
		return Output{}, ok
	}

	output := Output{
		Output:    tx.Transaction.Outputs[pos.OutputIndex],
		Spent:     tx.Spent[pos.OutputIndex],
		SpenderTx: tx.SpenderTxs[pos.OutputIndex],
	}

	return output, ok
}

// GetTx returns the transaction with the provided transaction hash.
func (store DataStore) GetTx(ctx sdk.Context, hash []byte) (Transaction, bool) {
	key := GetTxKey(hash)
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

// GetTxWithPosition returns the transaction that contains the provided
// position as an output.
func (store DataStore) GetTxWithPosition(ctx sdk.Context, pos plasma.Position) (Transaction, bool) {
	key := GetOutputKey(pos)
	hash := store.Get(ctx, key)
	return store.GetTx(ctx, hash)
}

// -----------------------------------------------------------------------------
/* Has */

// HasWallet returns whether an wallet at the given address exists.
func (store DataStore) HasWallet(ctx sdk.Context, addr common.Address) bool {
	key := GetWalletKey(addr)
	return store.Has(ctx, key)
}

// HasDeposit returns whether a deposit with the given nonce exists.
func (store DataStore) HasDeposit(ctx sdk.Context, nonce *big.Int) bool {
	key := GetDepositKey(nonce)
	return store.Has(ctx, key)
}

// HasFee returns whether a fee with the given position exists.
func (store DataStore) HasFee(ctx sdk.Context, pos plasma.Position) bool {
	key := GetFeeKey(pos)
	return store.Has(ctx, key)
}

// HasOutput returns whether an output with the given position exists.
func (store DataStore) HasOutput(ctx sdk.Context, pos plasma.Position) bool {
	key := GetOutputKey(pos)
	hash := store.Get(ctx, key)

	return store.HasTx(ctx, hash)
}

// HasTx returns whether a transaction with the given transaction hash
// exists.
func (store DataStore) HasTx(ctx sdk.Context, hash []byte) bool {
	key := GetTxKey(hash)
	return store.Has(ctx, key)
}

// -----------------------------------------------------------------------------
/* Set */

// setWallet overwrites the wallet stored at the given address.
func (store DataStore) setWallet(ctx sdk.Context, addr common.Address, wallet Wallet) {
	key := GetWalletKey(addr)
	data, err := rlp.EncodeToBytes(&wallet)
	if err != nil {
		panic(fmt.Sprintf("error marshaling wallet with address %s: %s", addr, err))
	}

	store.Set(ctx, key, data)
}

// setDeposit overwrites the deposit stored at the given nonce.
func (store DataStore) setDeposit(ctx sdk.Context, nonce *big.Int, deposit Deposit) {
	data, err := rlp.EncodeToBytes(&deposit)
	if err != nil {
		panic(fmt.Sprintf("error marshaling deposit with nonce %s: %s", nonce, err))
	}

	key := GetDepositKey(nonce)
	store.Set(ctx, key, data)
}

// setFee overwrites the fee stored at the given position.
func (store DataStore) setFee(ctx sdk.Context, pos plasma.Position, fee Output) {
	data, err := rlp.EncodeToBytes(&fee)
	if err != nil {
		panic(fmt.Sprintf("error marshaling fee with position %s: %s", pos, err))
	}

	key := GetFeeKey(pos)
	store.Set(ctx, key, data)
}

// setOutput adds a mapping from position to transaction hash.
func (store DataStore) setOutput(ctx sdk.Context, pos plasma.Position, hash []byte) {
	key := GetOutputKey(pos)
	store.Set(ctx, key, hash)
}

// setTx overwrites the mapping from transaction hash to transaction.
func (store DataStore) setTx(ctx sdk.Context, tx Transaction) {
	data, err := rlp.EncodeToBytes(&tx)
	if err != nil {
		panic(fmt.Sprintf("error marshaling transaction: %s", err))
	}

	key := GetTxKey(tx.Transaction.TxHash())
	store.Set(ctx, key, data)
}

// -----------------------------------------------------------------------------
/* Store */

// StoreDeposit adds an unspent deposit and updates the deposit owner's
// wallet.
func (store DataStore) StoreDeposit(ctx sdk.Context, nonce *big.Int, deposit plasma.Deposit) {
	store.setDeposit(ctx, nonce, Deposit{deposit, false, make([]byte, 0)})
	store.addToWallet(ctx, deposit.Owner, deposit.Amount, plasma.NewPosition(big.NewInt(0), 0, 0, nonce))
}

// StoreFee adds an unspent fee and updates the fee owner's wallet.
func (store DataStore) StoreFee(ctx sdk.Context, blockNum *big.Int, output plasma.Output) {
	pos := plasma.NewPosition(blockNum, 1<<16-1, 0, big.NewInt(0))
	store.setFee(ctx, pos, Output{output, false, make([]byte, 0)})
	store.addToWallet(ctx, output.Owner, output.Amount, pos)
}

// StoreTx adds the transaction.
func (store DataStore) StoreTx(ctx sdk.Context, tx Transaction) {
	store.setTx(ctx, tx)
}

// StoreOutputs adds new Output UTXO's to respective owner wallets.
func (store DataStore) StoreOutputs(ctx sdk.Context, tx Transaction) {
	for i, output := range tx.Transaction.Outputs {
		store.addToWallet(ctx, output.Owner, output.Amount, plasma.NewPosition(tx.Position.BlockNum, tx.Position.TxIndex, uint8(i), big.NewInt(0)))
		store.setOutput(ctx, plasma.NewPosition(tx.Position.BlockNum, tx.Position.TxIndex, uint8(i), big.NewInt(0)), tx.Transaction.TxHash())
	}
}

// -----------------------------------------------------------------------------
/* Spend */

// SpendDeposit changes the deposit to be spent and updates the wallet of
// the deposit owner.
func (store DataStore) SpendDeposit(ctx sdk.Context, nonce *big.Int, spenderTx []byte) sdk.Result {
	deposit, ok := store.GetDeposit(ctx, nonce)
	if !ok {
		return ErrDNE(fmt.Sprintf("deposit with nonce %s does not exist", nonce)).Result()
	} else if deposit.Spent {
		return ErrOutputSpent(fmt.Sprintf("deposit with nonce %s is already spent", nonce)).Result()
	}

	deposit.Spent = true
	deposit.SpenderTx = spenderTx

	store.setDeposit(ctx, nonce, deposit)
	store.subtractFromWallet(ctx, deposit.Deposit.Owner, deposit.Deposit.Amount, plasma.NewPosition(big.NewInt(0), 0, 0, nonce))

	return sdk.Result{}
}

// SpendFee changes the fee to be spent and updates the wallet of the fee
// owner.
func (store DataStore) SpendFee(ctx sdk.Context, pos plasma.Position, spenderTx []byte) sdk.Result {
	fee, ok := store.GetFee(ctx, pos)
	if !ok {
		return ErrDNE(fmt.Sprintf("fee with position %s does not exist", pos)).Result()
	} else if fee.Spent {
		return ErrOutputSpent(fmt.Sprintf("fee with position %s is already spent", pos)).Result()
	}

	fee.Spent = true
	fee.SpenderTx = spenderTx

	store.setFee(ctx, pos, fee)
	store.subtractFromWallet(ctx, fee.Output.Owner, fee.Output.Amount, pos)

	return sdk.Result{}
}

// SpendOutput changes the output to be spent and updates the wallet of the
// output owner.
func (store DataStore) SpendOutput(ctx sdk.Context, pos plasma.Position, spenderTx []byte) sdk.Result {
	key := GetOutputKey(pos)
	hash := store.Get(ctx, key)

	tx, ok := store.GetTx(ctx, hash)
	if !ok {
		return ErrDNE(fmt.Sprintf("output with index %x and transaction hash 0x%x does not exist", pos.OutputIndex, hash)).Result()
	} else if tx.Spent[pos.OutputIndex] {
		return ErrOutputSpent(fmt.Sprintf("output with index %x and transaction hash 0x%x is already spent", pos.OutputIndex, hash)).Result()
	}

	tx.Spent[pos.OutputIndex] = true
	tx.SpenderTxs[pos.OutputIndex] = spenderTx

	store.setTx(ctx, tx)
	store.subtractFromWallet(ctx, tx.Transaction.Outputs[pos.OutputIndex].Owner, tx.Transaction.Outputs[pos.OutputIndex].Amount, pos)

	return sdk.Result{}
}

// -----------------------------------------------------------------------------
/* Helpers */

// GetUnspentForWallet returns the unspent outputs that belong to the given
// wallet. Returns the struct TxOutput so user has access to the
// transactional information related to the output.
func (store DataStore) GetUnspentForWallet(ctx sdk.Context, wallet Wallet) (utxos []TxOutput) {
	for _, p := range wallet.Unspent {
		output, ok := store.GetOutput(ctx, p)
		if !ok {
			panic(fmt.Sprintf("Corrupted store: Wallet contains unspent position (%v) that doesn't exist in store", p))
		}
		tx, ok := store.GetTxWithPosition(ctx, p)
		if !ok {
			panic(fmt.Sprintf("Corrupted store: Wallet contains unspent position (%v) that doesn't have corresponding tx", p))
		}

		txo := NewTxOutput(output.Output, p, tx.ConfirmationHash, tx.Transaction.TxHash(), output.Spent, output.SpenderTx)
		utxos = append(utxos, txo)
	}
	return utxos
}

// depositToOutput retrieves the deposit with the given nonce, and returns
// it as an output.
func (store DataStore) depositToOutput(ctx sdk.Context, nonce *big.Int) (Output, bool) {
	deposit, ok := store.GetDeposit(ctx, nonce)
	if !ok {
		return Output{}, ok
	}
	output := Output{
		Output:    plasma.NewOutput(deposit.Deposit.Owner, deposit.Deposit.Amount),
		Spent:     deposit.Spent,
		SpenderTx: deposit.SpenderTx,
	}
	return output, ok
}

// addToWallet adds the passed in amount to the wallet with the given
// address and adds the position provided to the list of unspent positions
// within the wallet.
func (store DataStore) addToWallet(ctx sdk.Context, addr common.Address, amount *big.Int, pos plasma.Position) {
	wallet, ok := store.GetWallet(ctx, addr)
	if !ok {
		wallet = Wallet{big.NewInt(0), make([]plasma.Position, 0), make([]plasma.Position, 0)}
	}

	wallet.Balance = new(big.Int).Add(wallet.Balance, amount)
	wallet.Unspent = append(wallet.Unspent, pos)
	store.setWallet(ctx, addr, wallet)
}

// subtractFromWallet subtracts the passed in amount from the wallet with
// the given address and moves the provided position from the unspent list
// to the spent list.
func (store DataStore) subtractFromWallet(ctx sdk.Context, addr common.Address, amount *big.Int, pos plasma.Position) {
	wallet, ok := store.GetWallet(ctx, addr)
	if !ok {
		panic(fmt.Sprintf("output store has been corrupted"))
	}

	// Update Wallet
	wallet.Balance = new(big.Int).Sub(wallet.Balance, amount)
	if wallet.Balance.Sign() == -1 {
		panic(fmt.Sprintf("wallet with address 0x%x has a negative balance", addr))
	}

	wallet.Unspent = removePosition(wallet.Unspent, pos)
	wallet.Spent = append(wallet.Spent, pos)
	store.setWallet(ctx, addr, wallet)
}

// helper function to remove a position from the unspent list.
func removePosition(positions []plasma.Position, pos plasma.Position) []plasma.Position {
	for i, p := range positions {
		if p.String() == pos.String() {
			positions = append(positions[:i], positions[i+1:]...)
		}
	}
	return positions
}
