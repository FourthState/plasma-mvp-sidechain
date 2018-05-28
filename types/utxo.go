package types

import (
	"errors"
	rlp "github.com/ethereum/go-ethereum/rlp"
	amino "github.com/tendermint/go-amino"
	crypto "github.com/tendermint/go-crypto"
	utils "github.com/FourthState/plasma-mvp-sidechain/utils"
)

// UTXO is a standard unspent transaction output
type UTXO interface {
	// Address to which UTXO's are sent
	GetAddress() crypto.Address
	SetAddress(crypto.Address) error // errors if already set

	GetInputAddresses() [2]crypto.Address
	SetInputAddresses([2]crypto.Address) error

	GetDenom() uint64
	SetDenom(uint64) error //errors if already set

	GetPosition() Position
	SetPosition(uint64, uint16, uint8, uint8) error

	Get(key interface{}) (value interface{}, err error)
	Set(key interface{}, value interface{}) error
}

// Implements UTXO interface
type BaseUTXO struct {
	InputAddresses [2]crypto.Address
	Address        crypto.Address
	Denom          uint64
	Position       Position
}

func NewBaseUTXO(addr crypto.Address, inputaddr [2]crypto.Address, denom uint64,
	position Position) UTXO {
	return BaseUTXO{
		InputAddresses: inputaddr,
		Address:        addr,
		Denom:          denom,
		Position:       position,
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
	if addr == nil || utils.ZeroAddress(addr) {
		return errors.New("address provided is nil")
	}
	utxo.Address = addr
	return nil
}

func (utxo BaseUTXO) SetInputAddresses(addrs [2]crypto.Address) error {
	if utxo.InputAddresses[0] != nil {
		return errors.New("cannot override BaseUTXO Address")
	}
	if addrs[0] == nil || utils.ZeroAddress(addrs[0]) {
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

func (utxo BaseUTXO) SetPosition(blockNum uint64, txIndex uint16, oIndex uint8, depositNum uint8) error {
	if utxo.Position.Blknum != 0 {
		return errors.New("Cannot override BaseUTXO Position")
	}
	utxo.Position = Position{blockNum, txIndex, oIndex, depositNum}
	return nil
}

//----------------------------------------
// misc

type Position struct {
	Blknum     uint64
	TxIndex    uint16
	Oindex     uint8
	DepositNum uint8
}

func NewPosition(blknum uint64, txIndex uint16, oIndex uint8, depositNum uint8) Position {
	return Position{
		Blknum:     blknum,
		TxIndex:    txIndex,
		Oindex:     oIndex,
		DepositNum: depositNum,
	}
}

// Used to determine Sign Bytes for confirm signatures
func (position Position) GetSignBytes() []byte {
	b, err := rlp.EncodeToBytes(position)
	if err != nil {
		panic(err)
	}
	return b
}

//-------------------------------------------------------

func RegisterAmino(cdc *amino.Codec) {
	cdc.RegisterInterface((*UTXO)(nil), nil)
	cdc.RegisterConcrete(BaseUTXO{}, "types/BaseUTXO", nil)
	cdc.RegisterConcrete(Position{}, "types/Position", nil)
	cdc.RegisterConcrete(BaseTx{}, "types/BaseTX", nil)
	cdc.RegisterConcrete(SpendMsg{}, "types/SpendMsg", nil)
}
