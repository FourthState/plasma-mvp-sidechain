package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	rlp "github.com/ethereum/go-ethereum/rlp"
	"github.com/tendermint/go-amino"
	crypto "github.com/tendermint/go-crypto"
	//"fmt"
)

// Consider correct types to use
// ConfirmSig1 has confirm signatures from spenders of transaction located at [Blknum1, Txindex1, Oindex1] with signatures in order
// ConfirmSig2 has confirm signatures from spenders of transaction located at [Blknum2, Txindex2, Oindex2] with signatures in order
type SpendMsg struct {
	Blknum1      uint64
	Txindex1     uint16
	Oindex1      uint8
	DepositNum1  uint8
	Owner1       crypto.Address
	ConfirmSigs1 [2]crypto.Signature
	Blknum2      uint64
	Txindex2     uint16
	Oindex2      uint8
	DepositNum2  uint8
	Owner2       crypto.Address
	ConfirmSigs2 [2]crypto.Signature
	Newowner1    crypto.Address
	Denom1       uint64
	Newowner2    crypto.Address
	Denom2       uint64
	Fee          uint64
}

func NewSpendMsg(blknum1 uint64, txindex1 uint16, oindex1 uint8,
	depositnum1 uint8, owner1 crypto.Address, confirmSigs1 [2]crypto.Signature,
	blknum2 uint64, txindex2 uint16, oindex2 uint8,
	depositnum2 uint8, owner2 crypto.Address, confirmSigs2 [2]crypto.Signature,
	newowner1 crypto.Address, denom1 uint64,
	newowner2 crypto.Address, denom2 uint64, fee uint64) SpendMsg {
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
func (msg SpendMsg) Type() string { return "txs" } // TODO: decide on something better

// Implements Msg.
func (msg SpendMsg) ValidateBasic() sdk.Error {
	if !ValidAddress(msg.Owner1) {
		return ErrInvalidAddress(DefaultCodespace, "Input owner must have a valid address")
	}
	if !ValidAddress(msg.Newowner1) {
		return ErrInvalidAddress(DefaultCodespace, "No recipients of transaction")
	}

	switch {

	case msg.Oindex1 != 0 && msg.Oindex1 != 1:
		return ErrInvalidOIndex(DefaultCodespace, "Output index 1 must be either 0 or 1")

	case msg.DepositNum1 != 0 && (msg.Blknum1 != 0 || msg.Txindex1 != 0 || msg.Oindex1 != 0):
		return ErrInvalidTransaction(DefaultCodespace, "First input is malformed. Deposit's position must be 0, 0, 0")

	case msg.DepositNum2 != 0 && (msg.Blknum2 != 0 || msg.Txindex2 != 0 || msg.Oindex2 != 0):
		return ErrInvalidTransaction(DefaultCodespace, "Second input is malformed. Deposit's position must be 0, 0, 0")

	case msg.Blknum2 != 0 && msg.Oindex2 != 0 && msg.Oindex2 != 1:
		return ErrInvalidOIndex(DefaultCodespace, "Output index 2 must be either 0 or 1")

	case msg.Blknum2 != 0 && msg.Denom2 == 0:
		return ErrInvalidDenom(DefaultCodespace, "Second denomination must be positive")

	case msg.Denom1 == 0:
		return ErrInvalidDenom(DefaultCodespace, "First denomination must be positive")
	}

	return nil
}

// Implements Msg.
func (msg SpendMsg) String() string {
	return "Spend" // TODO: Issue #3
}

// Implements Msg.
func (msg SpendMsg) Get(key interface{}) (value interface{}) {
	return nil // TODO: Implement
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
	addrs[0] = crypto.Address(msg.Owner1)
	if ValidAddress(msg.Owner2) {
		addrs = append(addrs, crypto.Address(msg.Owner2))
	}
	return addrs
}

//----------------------------------------
// BaseTx (Transaction wrapper for depositmsg and spendmsg)

type BaseTx struct {
	sdk.Msg
	Signatures []sdk.StdSignature
}

func NewBaseTx(msg SpendMsg, sigs []sdk.StdSignature) BaseTx {
	return BaseTx{
		Msg:        msg,
		Signatures: sigs,
	}
}

func (tx BaseTx) GetMsg() sdk.Msg                   { return tx.Msg }
func (tx BaseTx) GetFeePayer() crypto.Address       { return tx.Signatures[0].PubKey.Address() }
func (tx BaseTx) GetSignatures() []sdk.StdSignature { return tx.Signatures }

func RegisterAmino(cdc *amino.Codec) {
	// TODO: include option to always include prefix bytes
	cdc.RegisterInterface((*UTXO)(nil), nil)
	cdc.RegisterConcrete(BaseUTXO{}, "types/BaseUTXO", nil)
	cdc.RegisterConcrete(Position{}, "types/Position", nil)
	cdc.RegisterConcrete(SpendMsg{}, "plasma-mvp-sidechain/SpendMsg", nil)
	cdc.RegisterConcrete(BaseTx{}, "plasma-mvp-sidechain/BaseTx", nil)
}
