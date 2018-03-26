package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	crypto "github.com/tendermint/go-crypto"
)

// Consider correct types to use
// Also consider changing to have input/output structs. Not sure how 
// would work with rootContract
type SpendMsg struct {
	blknum1   uint
	txindex1  uint
	oindex1   uint
	owner1	  crypto.Address
	blknum2   uint
	txindex2  uint
	oindex2   uint
	owner2	  crypto.Address
	newowner1 crypto.Address
	denom1    uint
	newowner2  crypto.Address
	denom2    uint
	fee       uint
}

func NewSpendMsg(blknum1 uint, txindex1 uint, oindex1 uint, owner1 crypto.Address, blknum2 uint, txindex2 uint, oindex2 uint, owner2 crypto.Address,
				newowner1 crypto.Address, denom1 uint, newowner2 crypto.Address, denom2 uint, fee uint) SpendMsg {
	return SpendMsg{
		blknum1: 	blknum1,
		txindex1: 	txindex1,
		oindex1:	oindex1,
		owner1:		owner1,
		blknum2:	blknum2,
		txindex2:	txindex2,
		oindex2:	oindex2,
		owner2:		owner2,
		newowner1:	newowner1,
		denom1:		denom1,
		newowner2:	newowner2,
		denom2:		denom2,
		fee:		fee,
	}
}

// Implements Msg.
func (msg SpendMsg) Type() string { return "txs" } // TODO: decide on something better

// Implements Msg.
func (msg SpendMsg) ValidateBasic() sdk.Error {
	// this just ensures everything is correctly formatted
	// Add more checks?
	if msg.newowner1 == nil && msg.newowner2 == nil {
		return sdk.NewError(10,"No recipient of transaction")
	}
	return nil
}

// Implements Msg. 
func (msg SpendMsg) String() string {
	return "Spend" // TODO: Implement so contents of Msg are returned
}

// Implements Msg.
func (msg SpendMsg) Get(key interface{}) (value interface{}) {
	return nil // TODO: Implement 
}

// Implements Msg.
func (msg SpendMsg) GetSignBytes() []byte {
	// TODO: Implement
	return nil
}

// Implements Msg.
func (msg SpendMsg) GetSigners() []crypto.Address {
	// TODO
	addrs := make([]crypto.Address, 2)
	addrs[0] = msg.owner1
	addrs[1] = msg.owner2
	return addrs
}
//----------------------------------------
// DepositMsg

//Consider changing rootchain contract
type DepositMsg struct {
	owner 	crypto.Address
	denom 	uint
}

// Implements Msg.
func (msg DepositMsg) Type() string { return "txs" } // TODO: decide on something better

// Implements Msg.
func (msg DepositMsg) ValidateBasic() sdk.Error {
	// this just ensures everything is correctly formatted
	// TODO: Implement
	return nil
}

// Implements Msg. 
func (msg DepositMsg) String() string {
	return "Deposit" // TODO: Implement so contents of Msg are returned
}

// Implements Msg.
func (msg DepositMsg) Get(key interface{}) (value interface{}) {
	return nil // TODO: Implement 
}

// Implements Msg.
func (msg DepositMsg) GetSignBytes() []byte {
	// TODO: Implement
	return nil
}

// Implements Msg.
func (msg DepositMsg) GetSigners() []crypto.Address {
	// TODO: Implement
	return nil
}

//----------------------------------------
// BaseTx (Transaction wrapper for depositmsg and spendmsg)

type BaseTx struct {
	sdk.Msg
	Signatures []sdk.StdSignature
}

func NewStdTx(msg sdk.Msg, sigs []sdk.StdSignature) BaseTx {
	return BaseTx{
		Msg: 		msg,
		Signatures: sigs,
	}
}

func (tx BaseTx) GetMsg() sdk.Msg 					{ return tx.Msg }
func (tx BaseTx) GetFeePayer() crypto.Address	{ return tx.Signatures[0].PubKey.Address() }
func (tx BaseTx) GetSignatures() []sdk.StdSignature { return tx.Signatures }