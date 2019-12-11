package msgs

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Codes for msgs errors
const (
	DefaultCodespace sdk.CodespaceType = "msgs"

	CodeInvalidSpendMsg          sdk.CodeType = 1
	CodeInvalidIncludeDepositMsg sdk.CodeType = 2
)

// ErrInvalidSpendMsg error for an invalid spend msg
func ErrInvalidSpendMsg(codespace sdk.CodespaceType, msg string, args ...interface{}) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidSpendMsg, msg, args...)
}

// ErrInvalidIncludeDepositMsg error for an invalid include deposit msg
func ErrInvalidIncludeDepositMsg(codespace sdk.CodespaceType, msg string, args ...interface{}) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidIncludeDepositMsg, msg, args...)
}
