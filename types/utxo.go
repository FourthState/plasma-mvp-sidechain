package types //might need to be adjusted

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

	//Get and Set denomination of the utxo. Is uint64 appropriate type?
	GetDenom() uint64
	SetDenom(uint64) error //errors if already set

	GetPosition() [3]uint
	SetPosition(uint, uint, uint) error

	GetConfirmSigs() [2]crypto.Signature
	SetConfirmSigs([2]crypto.Signature) error

	Get(key interface{}) (value interface{}, err error)
	Set(key interface{}, value interface{}) error

	//TODO: ADD SUPPORT FOR DIFFERENT COINS

}

type UTXOHolder interface {
	GetUTXO(position [3]uint) (UTXO, int)
	DeleteUTXO(utxo UTXO) error
	AddUTXO(utxo UTXO) error
	FinalizeUTXO(denom uint64, sigs [2]crypto.Signature, position [3]uint) error
	GetLength() int
}

// BaseUTXO must have all confirm signatures in order of most recent up until the signatures of the original depsosits.
type BaseUTXO struct {
	Address     crypto.Address
	Denom       uint64
	Position    [3]uint
	ConfirmSigs [2]crypto.Signature
}

func NewBaseUTXO(addr crypto.Address, denom uint64) BaseUTXO {
	return BaseUTXO{
		Address:     addr,
		Denom:       denom,
		Position:    [3]uint{0, 0, 0},
		ConfirmSigs: emptySignatures(),
	}
}

// Implements UTXO
// Not sure what this is supposed to achieve. Modeled from baseaccount
func (utxo BaseUTXO) Get(key interface{}) (value interface{}, err error) {
	panic("not implemented yet")
}

// Implements UTXO
// Not sure what this is supposed to achieve. Modeled from baseaccount
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
	if addr == nil {
		return errors.New("address provided is nil")
	}
	utxo.Address = addr
	return nil
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

func (utxo BaseUTXO) GetPosition() [3]uint {
	return utxo.Position
}

func (utxo BaseUTXO) SetPosition(blockNum uint, txIndex uint, oIndex uint) error {
	if utxo.Position[0] != 0 {
		return errors.New("Cannot override BaseUTXO Position")
	}
	utxo.Position = [3]uint{blockNum, txIndex, oIndex}
	return nil
}

func (utxo BaseUTXO) GetConfirmSigs() [2]crypto.Signature {
	return utxo.ConfirmSigs
}

func (utxo BaseUTXO) SetConfirmSigs(sigs [2]crypto.Signature) error {
	if utxo.GetConfirmSigs() != emptySignatures() {
		return errors.New("Confirm Sigs already set")
	}

	utxo.ConfirmSigs = sigs
	return nil
}

func emptySignatures() [2]crypto.Signature {
	return [2]crypto.Signature{crypto.SignatureSecp256k1{}, crypto.SignatureSecp256k1{}}
}

//----------------------------------------
// UTXOHolder

// Holds a list of UTXO's
// All utxo's have same address, but possibly different denominations
type BaseUTXOHolder struct {
	UtxoList []UTXO
}

// Creates a new UTXOHolder
// UtxoList is a slice initialized with length 1 and capacity 10
func NewUTXOHolder() BaseUTXOHolder {
	return BaseUTXOHolder{
		UtxoList: make([]UTXO, 1, 10),
	}
}

// Gets the utxo from the UtxoList
func (uh BaseUTXOHolder) GetUTXO(position [3]uint) (UTXO, int) {
	for index, elem := range uh.UtxoList {
		if elem.GetPosition() == position {
			return elem, index
		}
	}
	return BaseUTXO{}, 0 // utxo is not in the list
}

// Delete utxo from utxoList
func (uh BaseUTXOHolder) DeleteUTXO(utxo UTXO) error {
	for index, elem := range uh.UtxoList {
		// If two utxo's are identical in the list it will delete the first one
		if elem.GetPosition() == utxo.GetPosition() {
			uh.UtxoList = append(uh.UtxoList[:index], uh.UtxoList[index+1:]...)
			return nil
		}
	}
	return errors.New("utxo does not exist in utxoList")
}

// Apends a utxo to the utxoList
func (uh BaseUTXOHolder) AddUTXO(utxo UTXO) error {
	uh.UtxoList = append(uh.UtxoList, utxo)
	return nil
}

func (uh BaseUTXOHolder) FinalizeUTXO(denom uint64, sigs [2]crypto.Signature, position [3]uint) error {
	for _, elem := range uh.UtxoList {
		// Find first unfinalized UTXO with same position and finalize with position
		if elem.GetDenom() == denom && elem.GetPosition()[0] == 0 {
			elem.SetPosition(position[0], position[1], position[2])
			elem.SetConfirmSigs(sigs)
			return nil
		}
	}
	return errors.New("Unfinalized UTXO with given position and denom does not exist")
}

func (uh BaseUTXOHolder) GetLength() int {
	return len(uh.UtxoList)
}
