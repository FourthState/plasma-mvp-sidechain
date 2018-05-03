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
	Fee          uint64
}

func NewSpendMsg(blknum1 uint64, txindex1 uint16, oindex1 uint8,
	indenom1 uint64, owner1 crypto.Address, confirmSigs1 [2]crypto.Signature,
	blknum2 uint64, txindex2 uint16, oindex2 uint8,
	indenom2 uint64, owner2 crypto.Address, confirmSigs2 [2]crypto.Signature,
	newowner1 crypto.Address, denom1 uint64,
	newowner2 crypto.Address, denom2 uint64, fee uint64) SpendMsg {
	return SpendMsg{
		Blknum1:      blknum1,
		Txindex1:     txindex1,
		Oindex1:      oindex1,
		Indenom1:     indenom1,
		Owner1:	 	  owner1,
		ConfirmSigs1: confirmSigs1, 
		Blknum2:      blknum2,
		Txindex2:     txindex2,
		Oindex2:      oindex2,
		Indenom2:     indenom2,
		Owner2:	 	  owner2,
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
		return sdk.NewError(100, "First owner must be filled")
	}
	if !ValidAddress(msg.Newowner1) {
		return sdk.NewError(100, "No recipients of transaction")
	}
	switch {
	case msg.Oindex1 != 0 && msg.Oindex1 != 1:
		return sdk.NewError(100, "Output index 1 must be either 0 or 1")
	case msg.Blknum2 != 0:
		if msg.Oindex2 != 0 && msg.Oindex2 != 1 {
			return sdk.NewError(100, "Output index 2 must be either 0 or 1")
		}
		if msg.Denom2 == 0 {
			return sdk.NewError(100, "Second denomination must be positive")
		}
	case msg.Indenom1 == 0:
		return sdk.NewError(100, "First input denomination must be positive.")
	case msg.Denom1 == 0:
		return sdk.NewError(100, "First denomination must be positive")
	case msg.Indenom1 + msg.Indenom2 == msg.Denom1 + msg.Denom2 + msg.Fee:
		return sdk.NewError(100, "Inputs do not equal outputs")
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
