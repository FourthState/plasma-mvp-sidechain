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
	blknum2   uint
	txindex2  uint
	oindex2   uint
	newowner1 crypto.Address
	denom1    uint
	newowner2  crypto.Address
	denom2    uint
	fee       uint
}

func NewSpendMsg(blknum1 uint, txindex1 uint, oindex1 uint, blknum2 uint, txindex2 uint, oindex2 uint,
				newowner1 crypto.Address, denom1 uint, newowner2 crypto.Address, denom2 uint, fee uint) SpendMsg {
	return SpendMsg{
		blknum1: 	blknum1,
		txindex1: 	txindex1,
		oindex1:	oindex1,
		blknum2:	blknum2,
		txindex2:	txindex2,
		oindex2:	oindex2,
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
	// TODO
}

// Implements Msg.
func (msg SpendMsg) GetSigners() []crypto.Address {
	// TODO
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

}

// Implements Msg. 
func (msg DepositMsg) String() string {
	return "Spend" // TODO: Implement so contents of Msg are returned
}

// Implements Msg.
func (msg DepositMsg) Get(key interface{}) (value interface{}) {
	return nil // TODO: Implement 
}

// Implements Msg.
func (msg DepositMsg) GetSignBytes() []byte {
	// TODO
}

// Implements Msg.
func (msg DepositMsg) GetSigners() []crypto.Address {
	// TODO
}