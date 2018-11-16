package types

import (
	"errors"
	"fmt"
	rlp "github.com/ethereum/go-ethereum/rlp"
	amino "github.com/tendermint/go-amino"

	utils "github.com/FourthState/plasma-mvp-sidechain/utils"
	"github.com/FourthState/plasma-mvp-sidechain/x/utxo"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"

	"github.com/tendermint/tendermint/crypto/tmhash"
)

const (
	// Only allowed Denomination on this plasma chain
	Denom = "Ether"
)

var _ utxo.UTXO = &BaseUTXO{}

// Implements UTXO interface
type BaseUTXO struct {
	MsgHash        []byte
	InputAddresses [2]common.Address
	Address        common.Address
	Amount         uint64
	Denom          string
	Position       PlasmaPosition
	TxHash         []byte
}

func ProtoUTXO(ctx sdk.Context, msg sdk.Msg) utxo.UTXO {
	spendmsg, ok := msg.(SpendMsg)
	if !ok {
		return nil
	}

	msgHash := ethcrypto.Keccak256(spendmsg.GetSignBytes())

	return &BaseUTXO{
		InputAddresses: [2]common.Address{spendmsg.Owner0, spendmsg.Owner1},
		TxHash:         tmhash.Sum(ctx.TxBytes()),
		MsgHash:        msgHash,
	}
}

func NewBaseUTXO(addr common.Address, inputaddr [2]common.Address, amount uint64,
	denom string, position PlasmaPosition) *BaseUTXO {
	return &BaseUTXO{
		InputAddresses: inputaddr,
		Address:        addr,
		Amount:         amount,
		Denom:          denom,
		Position:       position,
	}
}

func (baseutxo BaseUTXO) GetMsgHash() []byte {
	return baseutxo.MsgHash
}

//Implements UTXO
func (baseutxo BaseUTXO) GetAddress() []byte {
	return baseutxo.Address.Bytes()
}

//Implements UTXO
func (baseutxo *BaseUTXO) SetAddress(addr []byte) error {
	if !utils.ZeroAddress(baseutxo.Address) {
		return fmt.Errorf("address already set to: %X", baseutxo.Address)
	}
	address := common.BytesToAddress(addr)
	if utils.ZeroAddress(address) {
		return fmt.Errorf("invalid address provided: %X", address)
	}
	baseutxo.Address = address
	return nil
}

//Implements UTXO
func (baseutxo *BaseUTXO) SetInputAddresses(addrs [2]common.Address) error {
	if !utils.ZeroAddress(baseutxo.InputAddresses[0]) {
		return fmt.Errorf("input addresses already set to: %X, %X", baseutxo.InputAddresses[0], baseutxo.InputAddresses[1])
	}
	if utils.ZeroAddress(addrs[0]) {
		return fmt.Errorf("invalid address provided: %X", addrs[0])
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
		return fmt.Errorf("amount already set to: %d", baseutxo.Amount)
	}
	baseutxo.Amount = amount
	return nil
}

func (baseutxo BaseUTXO) GetPosition() utxo.Position {
	return baseutxo.Position
}

func (baseutxo *BaseUTXO) SetPosition(position utxo.Position) error {
	if baseutxo.Position.IsValid() {
		return fmt.Errorf("position already set to: %v", baseutxo.Position)
	} else if !position.IsValid() {
		return errors.New("invalid position provided")
	}

	plasmaposition, ok := position.(PlasmaPosition)
	if !ok {
		return errors.New("position must be of type PlasmaPosition")
	}
	baseutxo.Position = plasmaposition
	return nil
}

func (baseutxo BaseUTXO) GetDenom() string {
	return Denom
}

func (baseutxo *BaseUTXO) SetDenom(denom string) error {
	return nil
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

func NewPlasmaPosition(blknum uint64, txIndex uint16, oIndex uint8, depositNum uint64) PlasmaPosition {
	return PlasmaPosition{
		Blknum:     blknum,
		TxIndex:    txIndex,
		Oindex:     oIndex,
		DepositNum: depositNum,
	}
}

func (position PlasmaPosition) Get() []sdk.Uint {
	return []sdk.Uint{sdk.NewUint(position.Blknum), sdk.NewUint(uint64(position.TxIndex)), sdk.NewUint(uint64(position.Oindex)), sdk.NewUint(position.DepositNum)}
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
	cdc.RegisterConcrete(&BaseUTXO{}, "types/BaseUTXO", nil)
	cdc.RegisterConcrete(&PlasmaPosition{}, "types/PlasmaPosition", nil)
	cdc.RegisterConcrete(BaseTx{}, "types/BaseTX", nil)
	cdc.RegisterConcrete(SpendMsg{}, "types/SpendMsg", nil)
}
