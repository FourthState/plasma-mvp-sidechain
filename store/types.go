package store

import (
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"math/big"
)

/* Account */
type Account struct {
	Balance *big.Int          // total amount avaliable to be spent
	Unspent []plasma.Position // position of unspent transaction outputs
	Spent   []plasma.Position // position of spent transaction outputs
}

/* Wrap plasma deposit with spend information */
type Deposit struct {
	Deposit plasma.Deposit
	Spent   bool
	Spender []byte
}

/* Wrap plasma output with spend information */
type Output struct {
	Output  plasma.Output
	Spent   bool
	Spender []byte
}

/* Wrap plasma transaction with spend information */
type Transaction struct {
	Transaction      plasma.Transaction
	ConfirmationHash []byte
	Spent            []bool
	Spenders         [][]byte
	Position         plasma.Position
}

/* Wraps Output with transaction is was created with
   this allows for input addresses to be retireved */
type QueryOutput struct {
	Output Output
	Tx     Transaction
}
