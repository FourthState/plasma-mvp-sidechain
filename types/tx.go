package types

import (
	utils "github.com/FourthState/plasma-mvp-sidechain/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	rlp "github.com/ethereum/go-ethereum/rlp"
	crypto "github.com/tendermint/go-crypto"
)

var _ sdk.Msg = SpendMsg{}

type SpendMsg struct {
	Blknum1      uint64
	Txindex1     uint16
	Oindex1      uint8
	DepositNum1  uint64
	Owner1       common.Address
	ConfirmSigs1 [2]Signature
	Blknum2      uint64
	Txindex2     uint16
	Oindex2      uint8
	DepositNum2  uint64
	Owner2       common.Address
	ConfirmSigs2 [2]Signature
	Newowner1    common.Address
	Denom1       uint64
	Newowner2    common.Address
	Denom2       uint64
	Fee          uint64
}

func NewSpendMsg(blknum1 uint64, txindex1 uint16, oindex1 uint8,
	depositnum1 uint64, owner1 common.Address, confirmSigs1 [2]Signature,
	blknum2 uint64, txindex2 uint16, oindex2 uint8,
	depositnum2 uint64, owner2 common.Address, confirmSigs2 [2]Signature,
	newowner1 common.Address, denom1 uint64,
	newowner2 common.Address, denom2 uint64, fee uint64) SpendMsg {
	return SpendMsg{
		Blknum1:      blknum1,
		Txindex1:     txindex1,
		Oindex1:      oindex1,
		DepositNum1:  depositnum1,
		Owner1:       owner1,
		ConfirmSigs1: confirmSigs1,
		Blknum2:      blknum2,
		Txindex2:     txindex2,
		Oindex2:      oindex2,
		DepositNum2:  depositnum2,
		Owner2:       owner2,
		ConfirmSigs2: confirmSigs2,
		Newowner1:    newowner1,
		Denom1:       denom1,
		Newowner2:    newowner2,
		Denom2:       denom2,
		Fee:          fee,
	}
}

// Implements Msg.
func (msg SpendMsg) Type() string { return "spend" }

// Implements Msg.
func (msg SpendMsg) ValidateBasic() sdk.Error {
	if !utils.ValidAddress(msg.Owner1) {
		return ErrInvalidAddress(DefaultCodespace, "input owner must have a valid address")
	}
	if !utils.ValidAddress(msg.Newowner1) {
		return ErrInvalidAddress(DefaultCodespace, "no recipients of transaction")
	}
	if msg.Blknum1 == msg.Blknum2 && msg.Txindex1 == msg.Txindex2 && msg.Oindex1 == msg.Oindex2 && msg.DepositNum1 == msg.DepositNum2 {
		return ErrInvalidTransaction(DefaultCodespace, "Cannot spend same position twice")
	}

	switch {

	case msg.Oindex1 != 0 && msg.Oindex1 != 1:
		return ErrInvalidOIndex(DefaultCodespace, "output index 1 must be either 0 or 1")

	case msg.DepositNum1 != 0 && (msg.Blknum1 != 0 || msg.Txindex1 != 0 || msg.Oindex1 != 0):
		return ErrInvalidTransaction(DefaultCodespace, "first input is malformed. Deposit's position must be 0, 0, 0")

	case msg.DepositNum2 != 0 && (msg.Blknum2 != 0 || msg.Txindex2 != 0 || msg.Oindex2 != 0):
		return ErrInvalidTransaction(DefaultCodespace, "second input is malformed. Deposit's position must be 0, 0, 0")

	case msg.Blknum2 != 0 && msg.Oindex2 != 0 && msg.Oindex2 != 1:
		return ErrInvalidOIndex(DefaultCodespace, "output index 2 must be either 0 or 1")

	case msg.Denom1 == 0:
		return ErrInvalidDenom(DefaultCodespace, "first denomination must be positive")
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
func (msg SpendMsg) GetSigners() []crypto.Address {
	addrs := make([]crypto.Address, 1)
	addrs[0] = crypto.Address(msg.Owner1.Bytes())
	if utils.ValidAddress(msg.Owner2) {
		addrs = append(addrs, crypto.Address(msg.Owner2.Bytes()))
	}
	return addrs
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

func (tx BaseTx) GetMsg() sdk.Msg            { return tx.Msg }
func (tx BaseTx) GetSignatures() []Signature { return tx.Signatures }

//-----------------------------------------
// Wrapper for signature byte arrays
type Signature struct {
	Sig []byte
}

func (s Signature) Bytes() []byte {
	return s.Sig
}
