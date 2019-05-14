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

// Returns the amount avaliable to be spent
func (acc Account) GetBalance() *big.Int {
	return acc.Balance
}
