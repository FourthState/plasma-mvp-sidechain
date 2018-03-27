package types

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	crypto "github.com/tendermint/go-crypto"
	"github.com/tendermint/go-amino"
)

// Consider correct types to use
// Also consider changing to have input/output structs. Not sure how 
// would work with rootContract
type SpendMsg struct {
	Blknum1   uint
	Txindex1  uint
	Oindex1   uint
	Owner1	  crypto.Address
	Blknum2   uint
	Txindex2  uint
	Oindex2   uint
	Owner2	  crypto.Address
	Newowner1 crypto.Address
	Denom1    uint64
	Newowner2  crypto.Address
	Denom2    uint64
	Fee       uint
}

func NewSpendMsg(blknum1 uint, txindex1 uint, oindex1 uint, owner1 crypto.Address, blknum2 uint, txindex2 uint, oindex2 uint, owner2 crypto.Address,
				newowner1 crypto.Address, denom1 uint64, newowner2 crypto.Address, denom2 uint64, fee uint) SpendMsg {
	return SpendMsg{
		Blknum1: 	blknum1,
		Txindex1: 	txindex1,
		Oindex1:	oindex1,
		Owner1:		owner1,
		Blknum2:	blknum2,
		Txindex2:	txindex2,
		Oindex2:	oindex2,
		Owner2:		owner2,
		Newowner1:	newowner1,
		Denom1:		denom1,
		Newowner2:	newowner2,
		Denom2:		denom2,
		Fee:		fee,
	}
}

// Implements Msg.
func (msg SpendMsg) Type() string { return "txs" } // TODO: decide on something better

// Implements Msg.
func (msg SpendMsg) ValidateBasic() sdk.Error {
	// this just ensures everything is correctly formatted
	// Add more checks?
	if msg.Newowner1 == nil && msg.Newowner2 == nil {
		return sdk.NewError(100,"No recipients of transaction")
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
	// TODO: Implement with RLP encoding
	b, err := json.Marshal(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

// Implements Msg.
func (msg SpendMsg) GetSigners() []crypto.Address {
	// TODO
	addrs := make([]crypto.Address, 2)
	addrs[0] = msg.Owner1
	addrs[1] = msg.Owner2
	return addrs
}
//----------------------------------------
// DepositMsg

//Consider changing rootchain contract
type DepositMsg struct {
	Owner 	crypto.Address
	Denom 	uint
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
	// TODO: Implement with RLP encoding
	b, err := json.Marshal(msg) // XXX: ensure some canonical form
	if err != nil {
		panic(err)
	}
	return b
}

// Implements Msg.
func (msg DepositMsg) GetSigners() []crypto.Address {
	// TODO: Implement
	addrs := make([]crypto.Address, 2)
	addrs[0] = msg.Owner
	return addrs
}

//----------------------------------------
// BaseTx (Transaction wrapper for depositmsg and spendmsg)

type BaseTx struct {
	sdk.Msg
	Signatures []sdk.StdSignature
}

func NewBaseTx(msg sdk.Msg, sigs []sdk.StdSignature) BaseTx {
	return BaseTx{
		Msg: 		msg,
		Signatures: sigs,
	}
}

func (tx BaseTx) GetMsg() sdk.Msg 					{ return tx.Msg }
func (tx BaseTx) GetFeePayer() crypto.Address		{ return tx.Signatures[0].PubKey.Address() }
func (tx BaseTx) GetSignatures() []sdk.StdSignature { return tx.Signatures }

func RegisterAmino(cdc *amino.Codec) {
	// TODO include option to always include prefix bytes.
	cdc.RegisterConcrete(SpendMsg{}, "plasma-mvp-sidechain/SpendMsg", nil)
	cdc.RegisterConcrete(DepositMsg{}, "plasma-mvp-sidechain/DepositMsg", nil)
	cdc.RegisterConcrete(BaseTx{}, "plasma-mvp-sidechain/BaseTx", nil)
}