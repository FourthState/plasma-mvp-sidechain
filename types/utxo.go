package types

import (
	"errors"
	rlp "github.com/ethereum/go-ethereum/rlp"
	amino "github.com/tendermint/go-amino"

	utils "github.com/FourthState/plasma-mvp-sidechain/utils"
	"github.com/ethereum/go-ethereum/common"
)

// UTXO is a standard unspent transaction output
type UTXO interface {
	// Address that owns UTXO
	GetAddress() common.Address
	SetAddress(common.Address) error // errors if already set

	GetInputAddresses() [2]common.Address
	SetInputAddresses([2]common.Address) error

	GetDenom() uint64
	SetDenom(uint64) error //errors if already set

	GetPosition() Position
	SetPosition(uint64, uint16, uint8, uint64) error
}

// Implements UTXO interface
type BaseUTXO struct {
	InputAddresses [2]common.Address
	Address        common.Address
	Denom          uint64
	Position       Position
}

func NewBaseUTXO(addr common.Address, inputaddr [2]common.Address, denom uint64,
	position Position) UTXO {
	return &BaseUTXO{
		InputAddresses: inputaddr,
		Address:        addr,
		Denom:          denom,
		Position:       position,
	}
}

//Implements UTXO
func (utxo *BaseUTXO) GetAddress() common.Address {
	return utxo.Address
}

//Implements UTXO
func (utxo *BaseUTXO) SetAddress(addr common.Address) error {
	if utils.ZeroAddress(utxo.Address) {
		return errors.New("cannot override BaseUTXO Address")
	}
	if utils.ZeroAddress(addr) {
		return errors.New("address provided is nil")
	}
	utxo.Address = addr
	return nil
}

func (utxo *BaseUTXO) SetInputAddresses(addrs [2]common.Address) error {
	if utils.ZeroAddress(utxo.InputAddresses[0]) {
		return errors.New("cannot override BaseUTXO Address")
	}
	if utils.ZeroAddress(addrs[0]) {
		return errors.New("address provided is nil")
	}
	utxo.InputAddresses = addrs
	return nil
}

func (utxo *BaseUTXO) GetInputAddresses() [2]common.Address {
	return utxo.InputAddresses
}

//Implements UTXO
func (utxo *BaseUTXO) GetDenom() uint64 {
	return utxo.Denom
}

//Implements UTXO
func (utxo *BaseUTXO) SetDenom(denom uint64) error {
	if utxo.Denom != 0 {
		return errors.New("Cannot override BaseUTXO denomination")
	}
	utxo.Denom = denom
	return nil
}

func (utxo *BaseUTXO) GetPosition() Position {
	return utxo.Position
}

func (utxo *BaseUTXO) SetPosition(blockNum uint64, txIndex uint16, oIndex uint8, depositNum uint64) error {
	if utxo.Position.Blknum != 0 {
		return errors.New("cannot override BaseUTXO Position")
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
	DepositNum uint64
}

func NewPosition(blknum uint64, txIndex uint16, oIndex uint8, depositNum uint64) Position {
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
