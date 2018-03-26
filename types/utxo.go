package types //might need to be adjusted

import (
	crypto "github.com/tendermint/go-crypto"
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

)

// UTXO is a standard unspent transaction output
// Has pubkey for authentication
type UTXO interface {
	GetAddress() crypto.Address
	SetAddress(crypto.Address) error // errors if already set

	//Get and Set denomination of the utxo. Is uint64 appropriate type?
	GetDenom() uint64
	SetDenom(uint64) error //errors if already set

	Get(key interface{}) (value interface{}, err error)
	Set(key interface{}, value interface{}) error
	
	//TODO: ADD SUPPORT FOR DIFFERENT COINS

}

// UTXOMapper interface which stores and retrieves utxos from stores
// retrieved from the context
// Can create and destory utxo 
// Consider Changing?
type UTXOMapper interface {
	GetUXTO(ctx sdk.Context, addr crypto.Address) UTXO
	CreateUTXO(ctx sdk.Context, utxo UTXO)
	DestroyUTXO(ctx sdk.Context, utxo UTXO)
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