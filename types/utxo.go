package types

import (
	"errors"
	rlp "github.com/ethereum/go-ethereum/rlp"
	amino "github.com/tendermint/go-amino"

	utils "github.com/FourthState/plasma-mvp-sidechain/utils"
	"github.com/FourthState/plasma-mvp-sidechain/x/utxo"
	"github.com/ethereum/go-ethereum/common"
)

var _ utxo.UTXO = &BaseUTXO{}

// Implements UTXO interface
type BaseUTXO struct {
	InputAddresses [2]common.Address
	Address        common.Address
	Amount         uint64
	Denom          string
	Position       PlasmaPosition
}

func ProtoUTXO() utxo.UTXO {
	return &BaseUTXO{}
}

func NewBaseUTXO(addr common.Address, inputaddr [2]common.Address, amount uint64,
	denom string, position PlasmaPosition) utxo.UTXO {
	return &BaseUTXO{
		InputAddresses: inputaddr,
		Address:        addr,
		Amount:         amount,
		Denom:          denom,
		Position:       position,
	}
}

//Implements UTXO
func (baseutxo BaseUTXO) GetAddress() []byte {
	return baseutxo.Address.Bytes()
}

//Implements UTXO
func (baseutxo *BaseUTXO) SetAddress(addr []byte) error {
	if !utils.ZeroAddress(baseutxo.Address) {
		return errors.New("cannot override BaseUTXO Address")
	}
	address := common.BytesToAddress(addr)
	if utils.ZeroAddress(address) {
		return errors.New("address provided is nil")
	}
	baseutxo.Address = address
	return nil
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
	return nil
}

//Implements UTXO
func (baseutxo BaseUTXO) GetInputAddresses() [2]common.Address {
	return baseutxo.InputAddresses
}

//Implements UTXO
func (baseutxo BaseUTXO) GetAmount() uint64 {
	return baseutxo.Amount
}

//Implements UTXO
func (baseutxo *BaseUTXO) SetAmount(amount uint64) error {
	if baseutxo.Amount != 0 {
		return errors.New("cannot override BaseUTXO amount")
	}
	baseutxo.Amount = amount
	return nil
}

func (baseutxo BaseUTXO) GetPosition() utxo.Position {
	return &baseutxo.Position
}

func (baseutxo *BaseUTXO) SetPosition(position utxo.Position) error {
	if baseutxo.Position.IsValid() {
		return errors.New("cannot override BaseUTXO position")
	}
	plasmaposition, ok := position.(*PlasmaPosition)
	if !ok {
		return errors.New("Position must be of type PlasmaPosition")
	}
	baseutxo.Position = *plasmaposition
	return nil
}

func (baseutxo BaseUTXO) GetDenom() string {
	return "Ether"
}

func (baseutxo *BaseUTXO) SetDenom(denom string) error {
	return errors.New("Cannot set denom")
}

//----------------------------------------
// Position

var _ utxo.Position = &PlasmaPosition{}

type PlasmaPosition struct {
	Blknum     uint64
	TxIndex    uint16
	Oindex     uint8
	DepositNum uint64
}

func NewPlasmaPosition(blknum uint64, txIndex uint16, oIndex uint8, depositNum uint64) *PlasmaPosition {
	return &PlasmaPosition{
		Blknum:     blknum,
		TxIndex:    txIndex,
		Oindex:     oIndex,
		DepositNum: depositNum,
	}
}

func (position PlasmaPosition) Get() []uint64 {
	return []uint64{position.Blknum, uint64(position.TxIndex), uint64(position.Oindex), position.DepositNum}
}

func (position *PlasmaPosition) Set(fields []uint64) error {
	if position.IsValid() {
		return errors.New("Position already set")
	}
	position.Blknum = fields[0]
	position.TxIndex = uint16(fields[1])
	position.Oindex = uint8(fields[2])
	position.DepositNum = fields[3]
	return nil
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
	// If position is a regular tx, output index must be 0 or 1 and depositnum must be 0
	if position.Blknum != 0 {
		return position.Oindex < 2 && position.DepositNum == 0
	} else {
		// If position represents deposit, depositnum is not 0 and txindex and oindex are 0.
		return position.DepositNum != 0 && position.TxIndex == 0 && position.Oindex == 0
	}
}

//-------------------------------------------------------
// misc
func RegisterAmino(cdc *amino.Codec) {
	cdc.RegisterConcrete(BaseUTXO{}, "types/BaseUTXO", nil)
	cdc.RegisterConcrete(PlasmaPosition{}, "types/PlasmaPosition", nil)
	cdc.RegisterConcrete(BaseTx{}, "types/BaseTX", nil)
	cdc.RegisterConcrete(SpendMsg{}, "types/SpendMsg", nil)
}
