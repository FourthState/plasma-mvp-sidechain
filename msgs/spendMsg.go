// Package msgs provides the messages that indicate either a spend of utxos or
// a deposit on the rootchain contract.
package msgs

import (
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	// SpendMsgRoute is used for routing this message.
	SpendMsgRoute = "spend"
)

// SpendMsg implements the RLP interface through `Transaction`
type SpendMsg struct {
	plasma.Transaction
}

// Type implements the sdk.Msg interface.
func (msg SpendMsg) Type() string { return "spend_utxo" }

// Route implements the sdk.Msg interface.
func (msg SpendMsg) Route() string { return SpendMsgRoute }

// GetSigners will attempt to retrieve the signers of the message.
// CONTRACT: a nil slice is returned if recovery fails
func (msg SpendMsg) GetSigners() []sdk.AccAddress {
	txHash := utils.ToEthSignedMessageHash(msg.TxHash())
	var addrs []sdk.AccAddress

	// recover first owner
	pubKey, err := crypto.SigToPub(txHash, msg.Inputs[0].Signature[:])
	if err != nil {
		return nil
	}
	addrs = append(addrs, sdk.AccAddress(crypto.PubkeyToAddress(*pubKey).Bytes()))

	if len(msg.Inputs) > 1 {
		// recover the second owner
		pubKey, err = crypto.SigToPub(txHash, msg.Inputs[1].Signature[:])
		if err != nil {
			return nil
		}
		addrs = append(addrs, sdk.AccAddress(crypto.PubkeyToAddress(*pubKey).Bytes()))
	}

	return addrs
}

// GetSignBytes returns the Keccak256 hash of the transaction.
func (msg SpendMsg) GetSignBytes() []byte {
	return msg.TxHash()
}

// ValidateBasic verifies that the transaction is valid.
func (msg SpendMsg) ValidateBasic() sdk.Error {
	if err := msg.Transaction.ValidateBasic(); err != nil {
		return ErrInvalidSpendMsg(DefaultCodespace, err.Error())
	}

	return nil
}

// GetMsgs implements the sdk.Tx interface
func (msg SpendMsg) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}
