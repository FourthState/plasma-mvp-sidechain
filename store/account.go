package store

import (
	"math/big"
)

/* Account */
type Account struct {
	Balance *big.Int // total amount avaliable to be spent
	//	Unspent [][32]utxoCache // hashes of unspent transaction outputs
	//	Spent   [][32]utxoCache // hashes of spent transaction outputs
}

//func (a *Account) EncodeRLP(w io.writer) error {

//}

//func (a *Account) DecodeRLP(s *rlp.Stream) error {

//}

/* Account helper functions */
// Ret balance
// Get unspent utxos
