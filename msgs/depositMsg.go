package msgs

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
)

const (
	IncludeDepositMsgRoute = "include"
)

var _ sdk.Tx = IncludeDepositMsg{}

// Implements sdk.Msg and sdk.Tx interfaces
// since no authentication is happening
type IncludeDepositMsg struct {
	DepositNonce *big.Int
	Owner       common.Address
	ReplayNonce uint64 // to get around tx cache issues when resubmitting
}

func (msg IncludeDepositMsg) Type() string { return "include_deposit" }

func (msg IncludeDepositMsg) Route() string { return IncludeDepositMsgRoute }

// No signers necessary on IncludeDepositMsg
func (msg IncludeDepositMsg) GetSigners() []sdk.AccAddress {
	return nil
}

func (msg IncludeDepositMsg) GetSignBytes() []byte {
	return nil
}

func (msg IncludeDepositMsg) ValidateBasic() sdk.Error {
	if msg.DepositNonce.Sign() != 1 {
		return ErrInvalidTransaction(DefaultCodespace, "DepositNonce must be greater than 0")
	}
	if (utils.IsZeroAddress(msg.Owner)) {
		return ErrInvalidTransaction(DefaultCodespace, "Owner must have non-zero address")
	}
	return nil
}

// Also satisfy the sdk.Tx interface
func (msg IncludeDepositMsg) GetMsgs() []sdk.Msg {
	return []sdk.Msg{msg}
}
