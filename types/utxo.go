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

	// Address corresponding to confirm signature
	GetCSAddress() [2]crypto.Address
	SetCSAddress([2]crypto.Address) error // errors if already set

	// Public Key for UTXO
	GetPubKey() crypto.PubKey 
	SetPubKey(crypto.PubKey) error 

	// Public Key for confirm Sigs
	GetCSPubKey() [2]crypto.PubKey 
	SetCSPubKey([2]crypto.PubKey) error

	//Get and Set denomination of the utxo. Is uint64 appropriate type?
	GetDenom() uint64
	SetDenom(uint64) error //errors if already set

	GetPosition() Position
	SetPosition(uint64, uint16, uint8) error

	Get(key interface{}) (value interface{}, err error)
	Set(key interface{}, value interface{}) error

	//TODO: ADD SUPPORT FOR DIFFERENT COINS

}

// BaseUTXO must have all confirm signatures in order of most recent up until the signatures of the original depsosits.
type BaseUTXO struct {
	Address     crypto.Address
	CSAddress 	[2]crypto.Address
	PubKey 		crypto.PubKey
	CSPubKey 	[2]crypto.PubKey
	Denom       uint64
	Position    Position
}

func NewBaseUTXO(addr crypto.Address, csaddr [2]crypto.Address, pubkey crypto.PubKey, 
	cspubkey [2]crypto.PubKey, denom uint64, 
	position Position) BaseUTXO {
	return BaseUTXO{
		Address:     addr,
		CSAddress:	 csaddr,
		PubKey:		 pubkey,
		CSPubKey: 	 cspubkey,
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

// Implements UTXO
func (utxo BaseUTXO) GetCSAddress() [2]crypto.Address {
	return utxo.CSAddress
}

//Implements UTXO
func (utxo BaseUTXO) SetCSAddress(csaddr [2]crypto.Address) error {
	if utxo.Address != nil {
		return errors.New("cannot override BaseUTXO Confirm Signature Address")
	}
	if csaddr[0] == nil || ZeroAddress(csaddr[0]) {
		return errors.New("address provided is nil")
	}
	utxo.CSAddress = csaddr
	return nil
}

// Implements UTXO
func (utxo BaseUTXO) GetPubKey() crypto.PubKey {
	return utxo.PubKey
}

// Implements UTXO
func (utxo BaseUTXO) SetPubKey(pubkey crypto.PubKey) error {
	if utxo.PubKey != nil {
		return errors.New("cannot override BaseUTXO PubKey")
	}
	if pubkey == nil {
		return errors.New("pubkey provided is nil")
	}
	utxo.PubKey = pubkey
	return nil
}

// Implements UTXO
func (utxo BaseUTXO) GetCSPubKey() [2]crypto.PubKey {
	return utxo.CSPubKey
}

// Implements UTXO
func (utxo BaseUTXO) SetCSPubKey(cspubkey [2]crypto.PubKey) error {
	if utxo.CSPubKey[0] != nil {
		return errors.New("cannot override BaseUTXO confirm sig PubKey")
	}
	if cspubkey[0] == nil {
		return errors.New("pubkey provided is nil")
	}
	utxo.CSPubKey = cspubkey
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