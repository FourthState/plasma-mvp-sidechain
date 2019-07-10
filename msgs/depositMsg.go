package msgs

import (
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

const (
	IncludeDepositMsgRoute = "include"
)

var _ sdk.Tx = IncludeDepositMsg{}

// IncludeDepositMsg implements sdk.Msg and sdk.Tx interfaces since
// no authentication is happening.
type IncludeDepositMsg struct {
	DepositNonce *big.Int
	Owner        common.Address
	ReplayNonce  uint64 // to get around tx cache issues when resubmitting
}

// Type returns the message type.
func (msg IncludeDepositMsg) Type() string { return "include_deposit" }

// Route returns the route for this message.
func (msg IncludeDepositMsg) Route() string { return IncludeDepositMsgRoute }

// GetSigners returns nil since no signers necessary on IncludeDepositMsg.
func (msg IncludeDepositMsg) GetSigners() []sdk.AccAddress {
	return nil
}

// GetSignBytes returns nil since no signature validation required for
// IncludeDepositMsg.
func (msg IncludeDepositMsg) GetSignBytes() []byte {
	return nil
}

// ValidateBasic asserts that the DepositNonce is positive and that the
// Owner field is not the zero address.
func (msg IncludeDepositMsg) ValidateBasic() sdk.Error {
	if msg.DepositNonce.Sign() != 1 {
		return ErrInvalidIncludeDepositMsg(DefaultCodespace, "DepositNonce must be greater than 0")
	}
	if utils.IsZeroAddress(msg.Owner) {
		return ErrInvalidIncludeDepositMsg(DefaultCodespace, "Owner must have non-zero address")
	}
	return nil
}

// GetMsgs implements the sdk.Tx interface
func (msg IncludeDepositMsg) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}
