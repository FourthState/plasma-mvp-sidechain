package types 

import (
	"errors"
	crypto "github.com/tendermint/go-crypto"
)

// UTXO is a standard unspent transaction output
// Has pubkey for authentication
type UTXO interface {
	// Address to which UTXO's are sent
	GetAddress() crypto.Address
	SetAddress(crypto.Address) error // errors if already set

	GetInputAddresses() [2]crypto.Address
	SetInputAddresses([2]crypto.Address) error

	//Get and Set denomination of the utxo. Is uint64 appropriate type?
	GetDenom() uint64
	SetDenom(uint64) error //errors if already set

	GetPosition() Position
	SetPosition(uint64, uint16, uint8) error

	Get(key interface{}) (value interface{}, err error)
	Set(key interface{}, value interface{}) error

	//TODO: ADD SUPPORT FOR DIFFERENT COINS
	//Will possibly have to add support for input positions
	//Do not need to add Pubkey to struct. Instead client ecrecovers confirm signature to get Pubkey. Then call PubKey.Address() to verify and verify against InputAddresses.

}

// BaseUTXO must have all confirm signatures in order of most recent up until the signatures of the original depsosits.
type BaseUTXO struct {
	InputAddresses [2]crypto.Address
	Address     crypto.Address
	Denom       uint64
	Position    Position
}

func NewBaseUTXO(addr crypto.Address, inputaddr [2]crypto.Address, denom uint64, 
	position Position) UTXO {
	return BaseUTXO{
		InputAddresses:	 inputaddr,
		Address: addr,
		Denom:       denom,
		Position:    position,
	}
}

// Implements UTXO
func (utxo BaseUTXO) Get(key interface{}) (value interface{}, err error) {
	panic("not implemented yet")
}

// Implements UTXO
func (utxo BaseUTXO) Set(key interface{}, value interface{}) error {
	panic("not implemented yet")
}

//Implements UTXO
func (utxo BaseUTXO) GetAddress() crypto.Address {
	return utxo.Address
}

//Implements UTXO
func (utxo BaseUTXO) SetAddress(addr crypto.Address) error {
	if utxo.Address != nil {
		return errors.New("cannot override BaseUTXO Address")
	}
	if addr == nil || ZeroAddress(addr) {
		return errors.New("address provided is nil")
	}
	utxo.Address = addr
	return nil
}

func (utxo BaseUTXO) SetInputAddresses(addrs [2]crypto.Address) error {
	if utxo.InputAddresses[0] != nil {
		return errors.New("cannot override BaseUTXO Address")
	}
	if addrs[0] == nil || ZeroAddress(addrs[0]) {
		return errors.New("address provided is nil")
	}
	utxo.InputAddresses = addrs
	return nil
}

func (utxo BaseUTXO) GetInputAddresses() [2]crypto.Address {
	return utxo.InputAddresses
}

//Implements UTXO
func (utxo BaseUTXO) GetDenom() uint64 {
	return utxo.Denom
}

//Implements UTXO
func (utxo BaseUTXO) SetDenom(denom uint64) error {
	if utxo.Denom != 0 {
		return errors.New("Cannot override BaseUTXO denomination")
	}
	utxo.Denom = denom
	return nil
}

func (utxo BaseUTXO) GetPosition() Position {
	return utxo.Position
}

func (utxo BaseUTXO) SetPosition(blockNum uint64, txIndex uint16, oIndex uint8) error {
	if utxo.Position.Blknum != 0 {
		return errors.New("Cannot override BaseUTXO Position")
	}
	utxo.Position = Position{blockNum, txIndex, oIndex}
	return nil
}

//----------------------------------------
// misc
// total position = Position.Blknum * 1000000 + Position.TxIndex * 10 + Position.Oindex

type Position struct {
	Blknum 		uint64
	TxIndex		uint16
	Oindex 		uint8
}

func NewPosition(blknum uint64, txIndex uint16, oIndex uint8) Position {
	return Position{
		Blknum: 	blknum,
		TxIndex: 	txIndex,
		Oindex: 	oIndex,
	}
}