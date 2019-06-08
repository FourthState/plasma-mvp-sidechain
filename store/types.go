package store

import (
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"math/big"
)

// Wallet holds reference to the total balance, unspent, and spent outputs
// at a given address
type Wallet struct {
	Balance *big.Int          // total amount avaliable to be spent
	Unspent []plasma.Position // position of unspent transaction outputs
	Spent   []plasma.Position // position of spent transaction outputs
}

/* Wrap plasma deposit with spend information */
type Deposit struct {
	Deposit   plasma.Deposit
	Spent     bool
	SpenderTx []byte // transaction hash that spends this deposit
}

/* Wrap plasma output with spend information */
type Output struct {
	Output    plasma.Output
	Spent     bool
	SpenderTx []byte // transaction that spends this output
}

// TODO: remove
type OutputInfo struct {
	Output Output
	Tx     Transaction
}

/* Wrap plasma transaction with spend information */
type Transaction struct {
	Transaction      plasma.Transaction
	ConfirmationHash []byte
	Spent            []bool
	SpenderTxs       [][]byte // transaction hashes that spend the outputs of this transaction
	Position         plasma.Position
}

// TransactionOutput holds all transactional information related to an output
// It is used to return output information from the store
type TransactionOutput struct {
	plasma.Output
	Position       plasma.Position
	Spent          bool
	SpenderTx      []byte
	InputAddresses []ethcmn.Address
	InputPositions []plasma.Position
}

// NewTransactionOutput is a constructor function for TransactionOutput
func NewTransactionOutput(output plasma.Output, pos plasma.Position, spent bool, spenderTx []byte, inputAddresses []ethcmn.Address, inputPosition []plasma.Position) TransactionOutput {
	return TransactionOutput{
		Output:         output,
		Position:       pos,
		Spent:          spent,
		SpenderTx:      spenderTx,
		InputAddresses: inputAddresses,
		InputPositions: inputPositions,
	}
}
