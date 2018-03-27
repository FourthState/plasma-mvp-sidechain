package types //might need to be adjusted

import (
	crypto "github.com/tendermint/go-crypto"
	"errors"
	//sdk "github.com/cosmos/cosmos-sdk/types"

)

// UTXO is a standard unspent transaction output
// Has pubkey for authentication
type UTXO interface {
	// Decided to use address instead of public keys
	// since each utxo newly created will be created from an
	// Ethereum address
	GetAddress() crypto.Address
	SetAddress(crypto.Address) error // errors if already set

	//Get and Set denomination of the utxo. Is uint64 appropriate type?
	GetDenom() uint64
	SetDenom(uint64) error //errors if already set

	Get(key interface{}) (value interface{}, err error)
	Set(key interface{}, value interface{}) error
	
	//TODO: ADD SUPPORT FOR DIFFERENT COINS

}

type UTXOHolder interface {
	GetUTXO(denom uint64) (UTXO, int) 
	DeleteUTXO(utxo UTXO) error
	AddUTXO(utxo UTXO) error 
	GetLength() int 
	
}

// Consider moving BaseUTXO and AppUTXO to another file. Are they necessary?
// Currently being used a prototype
type BaseUTXO struct {
	Address crypto.Address
	Denom uint64
}

func NewBaseUTXO(addr crypto.Address, denom uint64) BaseUTXO {
	return BaseUTXO {
		Address: addr,
		Denom: denom,
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
func (utxo BaseUTXO) SetAddress(addr crypto.Address) error{
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

//----------------------------------------
// UTXOHolder

// Holds a list of UTXO's
// All utxo's have same address, but possibly different denominations 
type BaseUTXOHolder struct {
	utxoList []UTXO
}

// Creates a new UTXOHolder 
// utxoList is a slice initialized with length 1 and capacity 10
func NewUTXOHolder() BaseUTXOHolder {
	return BaseUTXOHolder {
		utxoList: 	make([]UTXO, 1, 10),
	}
}

// Gets the utxo from the utxoList
func (uh BaseUTXOHolder) GetUTXO(denom uint64) (UTXO, int) {
	for index, elem := range uh.utxoList {
		if elem.GetDenom() == denom {
			return elem, index
		}
	}
	return BaseUTXO{}, 0 //utxo is not in the list
}

// Delete utxo from utxoList
func (uh BaseUTXOHolder) DeleteUTXO(utxo UTXO) error {
	for index, elem := range uh.utxoList {
		// If two utxo's are identical in the list it will delete the first one
		if elem.GetDenom() == utxo.GetDenom() {
			uh.utxoList = append(uh.utxoList[:index], uh.utxoList[index + 1:]...)
			return nil
		}
	}
	return errors.New("utxo does not exist in utxoList")
}

// Apends a utxo to the utxoList
func (uh BaseUTXOHolder) AddUTXO(utxo UTXO) error {
	uh.utxoList = append(uh.utxoList, utxo)
	return nil
}

func (uh BaseUTXOHolder) GetLength() int {
	return len(uh.utxoList)
}
