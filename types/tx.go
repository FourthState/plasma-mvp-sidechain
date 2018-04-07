package types

import (
	big "math/big"
	sdk "github.com/cosmos/cosmos-sdk/types"
	crypto "github.com/tendermint/go-crypto"
	"github.com/tendermint/go-amino"
	rlp "github.com/ethereum/go-ethereum/rlp"
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
	switch {
	case ZeroAddress(msg.Newowner1): // address is 0x0
		return sdk.NewError(100,"Must provide address in Owner1 field")
	case msg.Blknum1 == 0:
		return msg.validateDepositMsg()
	}
	return msg.validateSpendMsg()
}

func (msg SpendMsg) validateDepositMsg() sdk.Error {
	switch {
	case !ZeroAddress(msg.Owner1) || !ZeroAddress(msg.Owner2):
		return sdk.NewError(100,"Deposit message malformed")
	case msg.Txindex1 != 0 || msg.Txindex2 != 0:
		return sdk.NewError(100,"Deposit message malformed")
	case msg.Oindex1 != 0 || msg.Oindex2 != 0:
		return sdk.NewError(100,"Deposit message malformed")
	case msg.Blknum2 != 0:
		return sdk.NewError(100,"Deposit message malformed")
	case msg.Denom1 <= 0:
		return sdk.NewError(100,"First denomination must be positive")
	case msg.Denom2 != 0:
		return sdk.NewError(100,"Deposit message malformed")
	case msg.Fee < 0:
		return sdk.NewError(100,"Fee cannot be negative")
	}
	return nil
}

func (msg SpendMsg) validateSpendMsg() sdk.Error {
	switch {
	case msg.Txindex1 < 0:
		return sdk.NewError(100,"Transaction index cannot be negative")
	case msg.Oindex1 != 0 && msg.Oindex1 != 1:
		return sdk.NewError(100,"Output index 1 must be either 0 or 1")
	case msg.Blknum2 != 0:
		if (msg.Txindex2 < 0) {
			return sdk.NewError(100, "Transaction index cannot be negative")
		}
		if (msg.Oindex2 != 0 && msg.Oindex2 != 1) {
			return sdk.NewError(100,"Output index 2 must be either 0 or 1")
		}
		if (msg.Denom2 <= 0) {
			return sdk.NewError(100, "Second denomination must be positive")
		}
	case msg.Denom1 <= 0:
		return sdk.NewError(100,"First denomination must be positive")
	case msg.Fee < 0:
		return sdk.NewError(100,"Fee cannot be negative")
	}
	return nil
}

func (msg SpendMsg) IsDeposit() bool {
	return msg.Blknum1 == 0
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
	b, err := rlp.EncodeToBytes(msg)
	if err != nil {
		panic(err)
	}
	return b
}

// Implements Msg.
func (msg SpendMsg) GetSigners() []crypto.Address {
	// TODO
	addrs := make([]crypto.Address, 1)
	addrs[0] = crypto.Address(msg.Owner1)
	if new(big.Int).SetBytes(msg.Owner2.Bytes()).Sign() != 0 {
		addrs = append(addrs, crypto.Address(msg.Owner2))
	}
	return addrs
}

type FinalizeMsg struct {
	Spend SpendMsg
	Position [3]uint
	ConfirmSigs []sdk.StdSignature
}

func NewFinalizeMsg(spend SpendMsg, position [3]uint, sigs []sdk.StdSignature) FinalizeMsg {
	return FinalizeMsg{
		Spend: spend,
		Position: position,
		ConfirmSigs: sigs,
	}
}

func (msg FinalizeMsg) Type() string {
	return "txs"
}

func (msg FinalizeMsg) String() string {
	return "Finalize"
}

func (msg FinalizeMsg) ValidateBasic() sdk.Error {
	switch {
	case msg.Position[0] <= msg.Spend.Blknum1 || msg.Position[0] <= msg.Spend.Blknum2:
		return sdk.NewError(100, "New UTXO Blocknum must have greater height than inputs")
	case msg.Position[1] < 0:
		return sdk.NewError(100, "Transaction index negative")
	case msg.Position[2] != 0 && msg.Position[2] != 0:
		return sdk.NewError(100, "Output index must be either 0 or 1")
	case len(msg.ConfirmSigs) == 0:
		return sdk.NewError(100, "Msg must have at least one confirm sig")
	}
	return nil
}

func (msg FinalizeMsg) Get(key interface{}) (val interface{}) {
	return nil
}

func (msg FinalizeMsg) GetSignBytes() []byte {
	b, err := rlp.EncodeToBytes(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg FinalizeMsg) GetSigners() []crypto.Address {
	return nil
}

//----------------------------------------
// BaseTx (Transaction wrapper for depositmsg and spendmsg)

type BaseTx struct {
	sdk.Msg
	Signatures []sdk.StdSignature
}

func NewBaseTx(msg SpendMsg, sigs []sdk.StdSignature) BaseTx {
	return BaseTx{
		Msg: 		msg,
		Signatures: sigs,
	}
}

func (tx BaseTx) GetMsg() sdk.Msg					{ return tx.Msg }
func (tx BaseTx) GetFeePayer() crypto.Address		{ return tx.Signatures[0].PubKey.Address() }
func (tx BaseTx) GetSignatures() []sdk.StdSignature { return tx.Signatures }

func RegisterAmino(cdc *amino.Codec) {
	// TODO include option to always include prefix bytes.
	cdc.RegisterConcrete(SpendMsg{}, "plasma-mvp-sidechain/SpendMsg", nil)
	cdc.RegisterConcrete(BaseTx{}, "plasma-mvp-sidechain/BaseTx", nil)
}