package msgs

import (
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	// ConfirmSigMsgRoute is used for routing this message.
	ConfirmSigMsgRoute = "confirmSig"
)

// SpendMsg implements the RLP interface through `Transaction`
type ConfirmSigMsg struct {
	Input1    plasma.Input
	Input2    plasma.Input
}

// Type implements the sdk.Msg interface.
func (msg ConfirmSigMsg) Type() string { return "confirm_sig" }

// Route implements the sdk.Msg interface.
func (msg ConfirmSigMsg) Route() string { return ConfirmSigMsgRoute }

// GetSigners will attempt to retrieve the signers of the message.
// CONTRACT: a nil slice is returned if recovery fails
func (msg ConfirmSigMsg) GetSigners() []sdk.AccAddress {
	positionHash1 := utils.ToEthSignedMessageHash(msg.Input1.Position.Bytes())
	var addrs []sdk.AccAddress

	// recover first owner
	pubKey, err := crypto.SigToPub(positionHash1, msg.Input1.Signature[:])
	if err != nil {
		return nil
	}
	addrs = append(addrs, sdk.AccAddress(crypto.PubkeyToAddress(*pubKey).Bytes()))

	if msg.Input2.Signature != [65]byte{} {
		// recover the second owner
		positionHash2 := utils.ToEthSignedMessageHash(msg.Input1.Position.Bytes())
		pubKey, err = crypto.SigToPub(positionHash2, msg.Input2.Signature[:])
		if err != nil {
			return nil
		}
		addrs = append(addrs, sdk.AccAddress(crypto.PubkeyToAddress(*pubKey).Bytes()))
	}

	return addrs
}

// GetSignBytes returns the Keccak256 hash of the transaction.
func (msg ConfirmSigMsg) GetSignBytes() []byte {
	return nil
}

// ValidateBasic verifies that the transaction is valid.
func (msg ConfirmSigMsg) ValidateBasic() sdk.Error {
	if err := msg.Input1.ValidateBasic(); err != nil {
		return ErrInvalidConfirmSigMsg(DefaultCodespace, err.Error())
	}

	if err := msg.Input2.ValidateBasic(); err != nil {
		return ErrInvalidConfirmSigMsg(DefaultCodespace, err.Error())
	}

	return nil
}

// GetMsgs implements the sdk.Tx interface
func (msg ConfirmSigMsg) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}