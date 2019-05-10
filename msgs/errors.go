package msgs

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = "msgs"

	CodeInvalidSpendMsg          sdk.CodeType = 1
	CodeInvalidIncludeDepositMsg sdk.CodeType = 2
)

func ErrInvalidSpendMsg(codespace sdk.CodespaceType, msg string, args ...interface{}) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidSpendMsg, msg, args...)
}

func ErrInvalidIncludeDepositMsg(codespace sdk.CodespaceType, msg string, args ...interface{}) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidIncludeDepositMsg, msg, args...)
}
