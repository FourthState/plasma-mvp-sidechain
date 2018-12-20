package types

import (
	"github.com/FourthState/plasma-mvp-sidechain/x/utxo"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	amino "github.com/tendermint/go-amino"
)

const (
	// Only allowed Denomination on this plasma chain
	Denom = "Ether"
)

//----------------------------------------
// Position

type PlasmaPosition struct {
	Blknum     uint64
	TxIndex    uint16
	Oindex     uint8
	DepositNum uint64
}

func ProtoPosition() utxo.Position {
	return &PlasmaPosition{}
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

func (position PlasmaPosition) IsDeposit() bool {
	if !position.IsValid() {
		return false
	}
	return position.DepositNum != 0
}

type Deposit struct {
	Owner    common.Address
	Amount   sdk.Uint
	BlockNum sdk.Uint
}

//-------------------------------------------------------
// misc
func RegisterAmino(cdc *amino.Codec) {
	cdc.RegisterConcrete(PlasmaPosition{}, "types/PlasmaPosition", nil)
	cdc.RegisterConcrete(BaseTx{}, "types/BaseTX", nil)
	cdc.RegisterConcrete(SpendMsg{}, "types/SpendMsg", nil)
}
