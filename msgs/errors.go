package msgs

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Reserve errors 100 ~ 199
const (
	DefaultCodespace sdk.CodespaceType = "msgs"

	CodeInvalidTransaction sdk.CodeType = 1
)

func ErrInvalidTransaction(codespace sdk.CodespaceType, msg string, args ...interface{}) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidTransaction, msg, args...)
}
