package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	rlp "github.com/ethereum/go-ethereum/rlp"
	"github.com/tendermint/go-amino"
	crypto "github.com/tendermint/go-crypto"
)

// Consider correct types to use
// ConfirmSig1 has confirm signatures from spenders of transaction located at [Blknum1, Txindex1, Oindex1] with signatures in order
// ConfirmSig2 has confirm signatures from spenders of transaction located at [Blknum2, Txindex2, Oindex2] with signatures in order
type SpendMsg struct {
	Blknum1      uint64
	Txindex1     uint16
	Oindex1      uint8
	Indenom1	 uint64
	Owner1		 crypto.Address
	ConfirmSigs1 [2]crypto.Signature
	Blknum2      uint64
	Txindex2     uint16
	Oindex2      uint8
	Indenom2	 uint64
	Owner2 		 crypto.Address
	ConfirmSigs2 [2]crypto.Signature
	Newowner1    crypto.Address
	Denom1       uint64
	Newowner2    crypto.Address
	Denom2       uint64
	Fee          uint
}

func NewSpendMsg(blknum1 uint, txindex1 uint, oindex1 uint,
	indenom1 uint64, owner1 crypto.Address, confirmSigs1 [2]crypto.Signature,
	blknum2 uint, txindex2 uint, oindex2 uint,
	indenom2 uint64, owner2 crypto.Address, confirmSigs2 [2]crypto.Signature,
	newowner1 crypto.Address, denom1 uint64,
	newowner2 crypto.Address, denom2 uint64, fee uint) SpendMsg {
	return SpendMsg{
		Blknum1:      blknum1,
		Txindex1:     txindex1,
		Oindex1:      oindex1,
		Indenom1:     indenom1,
		Owner1:	 	  owner1,
		ConfirmSigs1: [2]confirmSigs1, 
		Blknum2:      blknum2,
		Txindex2:     txindex2,
		Oindex2:      oindex2,
		Indenom2:     indenom2,
		Owner2:	 	  owner2,
		ConfirmSigs2: [2]confirmSigs2, 
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
	if msg.Newowner1 == nil && msg.Newowner2 == nil {
		return sdk.NewError(100, "No recipients of transaction")
	}
	switch {
	case ZeroAddress(msg.Newowner1): // address is 0x0
		return sdk.NewError(100, "Must provide address in Owner1 field")
	case msg.Blknum1 == 0:
		return msg.validateDepositMsg()
	}
	return msg.validateSpendMsg()
}

func (msg SpendMsg) validateDepositMsg() sdk.Error {
	switch {
	case !ZeroAddress(msg.Owner1) || !ZeroAddress(msg.Owner2):
		return sdk.NewError(100, "Deposit message malformed")
	case msg.Indenom1 == 0 && msg.Indenom2 == 0:
		return sdk.NewError(100, "Deposit message malformed")
	case msg.Indenom1 < 0 || msg.Indenom2 < 0:
		return sdk.NewError(100, "Deposit message malformed")
	case msg.Txindex1 != 0 || msg.Txindex2 != 0:
		return sdk.NewError(100, "Deposit message malformed")
	case msg.Oindex1 != 0 || msg.Oindex2 != 0:
		return sdk.NewError(100, "Deposit message malformed")
	case msg.Blknum2 != 0:
		return sdk.NewError(100, "Deposit message malformed")
	case msg.Denom1 <= 0:
		return sdk.NewError(100, "First denomination must be positive")
	case msg.Denom2 != 0:
		return sdk.NewError(100, "Deposit message malformed")
	case msg.Fee < 0:
		return sdk.NewError(100, "Fee cannot be negative")
	case !msg.ConfirmSigs1[0].IsZero() || !msg.ConfirmSigs1[1].IsZero():
		return sdk.NewError(100, "Deposit must have zero-bytes as confirm sigs")
	case !msg.ConfirmSigs2[0].IsZero() || !msg.ConfirmSigs2[1].IsZero():
		return sdk.NewError(100, "Deposit must have zero-bytes as confirm sigs")
	}
	return nil
}

// CA: Do we need the checks for < 0 if the type is uint?
func (msg SpendMsg) validateSpendMsg() sdk.Error {
	switch {
	case msg.Txindex1 < 0:
		return sdk.NewError(100, "Transaction index cannot be negative")
	case msg.Oindex1 != 0 && msg.Oindex1 != 1:
		return sdk.NewError(100, "Output index 1 must be either 0 or 1")
	case msg.Blknum2 != 0:
		if msg.Txindex2 < 0 {
			return sdk.NewError(100, "Transaction index cannot be negative")
		}
		if msg.Oindex2 != 0 && msg.Oindex2 != 1 {
			return sdk.NewError(100, "Output index 2 must be either 0 or 1")
		}
		if msg.Denom2 <= 0 {
			return sdk.NewError(100, "Second denomination must be positive")
		}
	case msg.Denom1 <= 0:
		return sdk.NewError(100, "First denomination must be positive")
	case msg.Fee < 0:
		return sdk.NewError(100, "Fee cannot be negative")
	case msg.Indenom1 < 0 || msg.Indenom2 < 0:
		return sdk.NewError(100, "Input denominations must be positive")
	}
	return nil
}

func (msg SpendMsg) IsDeposit() bool {
	return msg.Blknum1 == 0
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
	if !ZeroAddress(msg.Owner2) {
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
	cdc.RegisterConcrete(SpendMsg{}, "plasma-mvp-sidechain/SpendMsg", nil)
	cdc.RegisterConcrete(BaseTx{}, "plasma-mvp-sidechain/BaseTx", nil)
}
