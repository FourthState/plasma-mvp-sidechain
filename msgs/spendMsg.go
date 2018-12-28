package msgs

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type SpendMsg struct {
	plasma.Transaction
}

// Implement the sdk.Msg interface

func (msg SpendMsg) Type() string { return "spend_utxo" }

func (msg SpendMsg) Route() string { return "spend" }

func (msg SpendMsg) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, 1)
	addrs[0] = sdk.AccAddress(msg.Input0.Owner.Bytes())

	if !utils.IsZeroAddress(msg.Input1.Owner) {
		addrs = append(addrs, sdk.AccAddress(msg.Input1.Owner.Bytes()))
	}

	return addrs
}

func (msg SpendMsg) GetSignBytes() []byte {
	txHash := msg.TxHash()
	return txHash[:]
}

func (msg SpendMsg) ValidateBasic() sdk.Error {
	if msg.Input0.BlockNum.Cmp(msg.Input1.BlockNum) == 0 &&
		msg.Input0.TxIndex == msg.Input1.TxIndex &&
		msg.Input0.OutputIndex == msg.Input1.OutputIndex {
		return ErrInvalidTransaction(DefaultCodespace,
			fmt.Sprintf("cannot spend same position twice: (%d, %d, %d, %d)",
				msg.Input0.BlockNum, msg.Input0.TxIndex, msg.Input0.OutputIndex, msg.Input0.DepositNonce))
	}

	if utils.IsZeroAddress(msg.Input0.Owner) {
		return ErrInvalidAddress(DefaultCodespace, "first input owner must have a valid address: %x", msg.Input0.Owner)
	}
	if utils.IsZeroAddress(msg.Output0.Owner) {
		return ErrInvalidAddress(DefaultCodespace, "no recipients of transaction")
	}
	if utils.IsZeroAddress(msg.Output1.Owner) && msg.Output1.Amount.Sign() != 0 {
		return ErrInvalidTransaction(DefaultCodespace, "Cannot send funds to the nil address")
	}

	if msg.Output0.Amount.Sign() == 0 {
		return ErrInvalidAmount(DefaultCodespace, "first output must have a positive amount")
	}

	/* First input/output validation */
	if msg.Input0.OutputIndex != 0 && msg.Input0.OutputIndex != 1 {
		return ErrInvalidOIndex(DefaultCodespace, "first output index 0 must be either 0 or 1")
	}
	if msg.Input0.Position.IsDeposit() && (msg.Input0.BlockNum.Sign() != 0 || msg.Input0.TxIndex != 0 || msg.Input0.OutputIndex != 0) {
		return ErrInvalidTransaction(DefaultCodespace, "first input is malformed. cannot specify a deposit nonce and utxo position simultaneously")
	}
	if !msg.Input0.Position.IsDeposit() && len(msg.Input0.ConfirmSignatures) == 0 {
		return ErrInvalidTransaction(DefaultCodespace, "first input must include confirm signatures")
	}

	/* Second input/output validation if applicable */
	if msg.HasSecondInput() {
		if msg.Input1.IsDeposit() && (msg.Input1.BlockNum.Sign() != 0 || msg.Input1.TxIndex != 0 || msg.Input1.OutputIndex != 0) {
			return ErrInvalidTransaction(DefaultCodespace, "second input is malformed. cannot specify a deposit nonce and utxo position simultaneously")
		}
		if msg.Input1.IsDeposit() && (msg.Input1.OutputIndex != 0 && msg.Input1.OutputIndex != 1) {
			return ErrInvalidOIndex(DefaultCodespace, "second output index 0 must be either 0 or 1")
		}
		if !msg.Input1.IsDeposit() && len(msg.Input1.ConfirmSignatures) == 0 {
			return ErrInvalidTransaction(DefaultCodespace, "second input must include confirm signatures")
		}
	}

	return nil
}

// Also satisfy the sdk.Tx interface

func (msg SpendMsg) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}
