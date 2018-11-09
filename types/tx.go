package types

import (
	"fmt"
	utils "github.com/FourthState/plasma-mvp-sidechain/utils"
	utxo "github.com/FourthState/plasma-mvp-sidechain/x/utxo"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	rlp "github.com/ethereum/go-ethereum/rlp"
)

var _ utxo.SpendMsg = SpendMsg{}

type SpendMsg struct {
	Blknum0      uint64
	Txindex0     uint16
	Oindex0      uint8
	DepositNum0  uint64
	Owner0       common.Address
	ConfirmSigs0 [2]Signature
	Blknum1      uint64
	Txindex1     uint16
	Oindex1      uint8
	DepositNum1  uint64
	Owner1       common.Address
	ConfirmSigs1 [2]Signature
	Newowner0    common.Address
	Amount0      uint64
	Newowner1    common.Address
	Amount1      uint64
	FeeAmount    uint64
}

// Implements Msg. Improve later
func (msg SpendMsg) Type() string { return "spend_utxo" }

// Implements Msg.
func (msg SpendMsg) Route() string { return "spend" }

// Implements Msg.
func (msg SpendMsg) ValidateBasic() sdk.Error {
	if !utils.ValidAddress(msg.Owner0) {
		return ErrInvalidAddress(DefaultCodespace, "input owner must have a valid address", msg.Owner0)
	}
	if !utils.ValidAddress(msg.Newowner0) {
		return ErrInvalidAddress(DefaultCodespace, "no recipients of transaction")
	}
	if msg.Blknum0 == msg.Blknum1 && msg.Txindex0 == msg.Txindex1 && msg.Oindex0 == msg.Oindex1 && msg.DepositNum0 == msg.DepositNum1 {
		return ErrInvalidTransaction(DefaultCodespace, fmt.Sprintf("cannot spend same position twice: (%d, %d, %d, %d)", msg.Blknum0, msg.Txindex0, msg.Oindex0, msg.DepositNum0))

	}

	switch {

	case msg.Oindex0 != 0 && msg.Oindex0 != 1:
		return ErrInvalidOIndex(DefaultCodespace, "output index 0 must be either 0 or 1")

	case msg.DepositNum0 != 0 && (msg.Blknum0 != 0 || msg.Txindex0 != 0 || msg.Oindex0 != 0):
		return ErrInvalidTransaction(DefaultCodespace, "first input is malformed. Deposit's position must be 0, 0, 0")

	case msg.DepositNum1 != 0 && (msg.Blknum1 != 0 || msg.Txindex1 != 0 || msg.Oindex1 != 0):
		return ErrInvalidTransaction(DefaultCodespace, "second input is malformed. Deposit's position must be 0, 0, 0")

	case msg.Blknum1 != 0 && msg.Oindex1 != 0 && msg.Oindex1 != 1:
		return ErrInvalidOIndex(DefaultCodespace, "output index 1 must be either 0 or 1")

	case msg.Amount0 == 0:
		return ErrInvalidAmount(DefaultCodespace, "first amount must be positive")
	}

	return nil
}

// Implements Msg.
func (msg SpendMsg) GetSignBytes() []byte {
	b, err := rlp.EncodeToBytes(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// Implements Msg.
func (msg SpendMsg) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, 1)
	addrs[0] = sdk.AccAddress(msg.Owner0.Bytes())
	if utils.ValidAddress(msg.Owner1) {
		addrs = append(addrs, sdk.AccAddress(msg.Owner1.Bytes()))
	}
	return addrs
}

func (msg SpendMsg) Inputs() []utxo.Input {
	inputs := []utxo.Input{utxo.Input{
		Owner:    msg.Owner0.Bytes(),
		Position: NewPlasmaPosition(msg.Blknum0, msg.Txindex0, msg.Oindex0, msg.DepositNum0),
	}}
	if NewPlasmaPosition(msg.Blknum1, msg.Txindex1, msg.Oindex1, msg.DepositNum1).IsValid() {
		// Add valid second input
		inputs = append(inputs, utxo.Input{
			Owner:    msg.Owner1.Bytes(),
			Position: NewPlasmaPosition(msg.Blknum1, msg.Txindex1, msg.Oindex1, msg.DepositNum1),
		})
	}
	return inputs
}

func (msg SpendMsg) Outputs() []utxo.Output {
	outputs := []utxo.Output{utxo.Output{msg.Newowner0.Bytes(), Denom, msg.Amount0}}
	if msg.Amount1 != 0 {
		outputs = append(outputs, utxo.Output{msg.Newowner1.Bytes(), Denom, msg.Amount1})
	}
	return outputs
}

func (msg SpendMsg) Fee() []utxo.Output {
	return []utxo.Output{utxo.Output{
		Denom:  Denom,
		Amount: msg.FeeAmount,
	}}
}

//----------------------------------------
// BaseTx
var _ sdk.Tx = BaseTx{}

type BaseTx struct {
	Msg        SpendMsg
	Signatures []Signature
}

func NewBaseTx(msg SpendMsg, sigs []Signature) BaseTx {
	return BaseTx{
		Msg:        msg,
		Signatures: sigs,
	}
}

func (tx BaseTx) GetMsgs() []sdk.Msg         { return []sdk.Msg{tx.Msg} }
func (tx BaseTx) GetSignatures() []Signature { return tx.Signatures }

//-----------------------------------------
// Wrapper for signature byte arrays
type Signature struct {
	Sig []byte
}

func (s Signature) Bytes() []byte {
	return s.Sig
}
