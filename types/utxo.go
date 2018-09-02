package types

import (
	"errors"
	rlp "github.com/ethereum/go-ethereum/rlp"
	amino "github.com/tendermint/go-amino"

	utils "github.com/FourthState/plasma-mvp-sidechain/utils"
	"github.com/FourthState/plasma-mvp-sidechain/x/utxo"
	"github.com/ethereum/go-ethereum/common"
)

var _ utxo.UTXO = BaseUTXO{}

// Implements UTXO interface
type BaseUTXO struct {
	InputAddresses [2]common.Address
	Address        common.Address
	Amount         uint64
	Denom          string
	Position       Position
}

func NewBaseUTXO(addr common.Address, inputaddr [2]common.Address, amount uint64,
	denom string, position Position) utxo.UTXO {
	return &BaseUTXO{
		InputAddresses: inputaddr,
		Address:        addr,
		Amount:         amount,
		Denom:          denom,
		Position:       position,
	}
}

//Implements UTXO
func (baseutxo *BaseUTXO) GetAddress() common.Address {
	return baseutxo.Address
}

//Implements UTXO
func (baseutxo *BaseUTXO) SetAddress(addr common.Address) error {
	if !utils.ZeroAddress(baseutxo.Address) {
		return errors.New("cannot override BaseUTXO Address")
	}
	if utils.ZeroAddress(addr) {
		return errors.New("address provided is nil")
	}
	baseutxo.Address = addr
}

//Implements UTXO
func (baseutxo *BaseUTXO) SetInputAddresses(addrs [2]common.Address) error {
	if !utils.ZeroAddress(baseutxo.InputAddresses[0]) {
		return errors.New("cannot override BaseUTXO Address")
	}
	if utils.ZeroAddress(addrs[0]) {
		return errors.New("address provided is nil")
	}
	baseutxo.InputAddresses = addrs
}

//Implements UTXO
func (baseutxo *BaseUTXO) GetInputAddresses() [2]common.Address {
	return baseutxo.InputAddresses
}

//Implements UTXO
func (baseutxo *BaseUTXO) GetAmount() uint64 {
	return baseutxo.Amount
}

//Implements UTXO
func (baseutxo *BaseUTXO) SetAmount(amount uint64) error {
	if baseutxo.Amount != 0 {
		return errors.New("cannot override BaseUTXO amount")
	}
	baseutxo.Amount = amount
}

func (baseutxo *BaseUTXO) GetPosition() utxo.Position {
	return baseutxo.Position
}

func (baseutxo *BaseUTXO) SetPosition(position utxo.Position) error {
	if baseutxo.Position.IsValid() {
		return errors.New("cannot override BaseUTXO position")
	}
	baseutxo.Position = position
}

//----------------------------------------
// Position

var _ utxo.Position = PlasmaPosition{}

type PlasmaPosition struct {
	Blknum     uint64
	TxIndex    uint16
	Oindex     uint8
	DepositNum uint64
}

func NewPlasmaPosition(blknum uint64, txIndex uint16, oIndex uint8, depositNum uint64) PlasmaPosition {
	return PlasmaPosition{
		Blknum:     blknum,
		TxIndex:    txIndex,
		Oindex:     oIndex,
		DepositNum: depositNum,
	}
}

// Used to determine Sign Bytes for confirm signatures
// Implements Position
func (position PlasmaPosition) GetSignBytes() []byte {
	b, err := rlp.EncodeToBytes(position)
	if err != nil {
		panic(err)
	}
	return b
}

// check that the position is formatted correctly
// Implements Position
func (position PlasmaPosition) IsValid() bool {
	if position.Blknum != 0 {
		return position.Oindex < 2 && position.DepositNum == 0
	} else {
		return position.Blknum != 0
	}
}

//-------------------------------------------------------
// misc
func RegisterAmino(cdc *amino.Codec) {
	cdc.RegisterInterface((*UTXO)(nil), nil)
	cdc.RegisterConcrete(BaseUTXO{}, "types/BaseUTXO", nil)
	cdc.RegisterConcrete(Position{}, "types/PlasmaPosition", nil)
	cdc.RegisterConcrete(BaseTx{}, "types/BaseTX", nil)
	cdc.RegisterConcrete(SpendMsg{}, "types/SpendMsg", nil)
}
