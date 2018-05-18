package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Reserve errors 100 ~ 199
const (
	DefaultCodespace sdk.CodespaceType = 2

	CodeInvalidAddress 		sdk.CodeType = 101
	CodeInvalidOIndex 		sdk.CodeType = 102
	CodeInvalidDenomination sdk.CodeType = 103
	CodeInvalidIOF 			sdk.CodeType = 104
	CodeInvalidUTXO			sdk.CodeType = 105
)

func codeToDefaultMsg(code sdk.CodeType) string {
	switch code {
	default:
		return sdk.CodeToDefaultMsg(code)
	}
}

//----------------------------------------
// Error constructors

func ErrInvalidAddress(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAddress, msg)
}

func ErrInvalidOIndex(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidOIndex, msg)
}

func ErrInvalidDenom(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidDenomination, msg)
}

func ErrInvalidIOF(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidIOF, msg)
}

func ErrInvalidUTXO(codespace sdk.CodespaceType, msg string) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidUTXO, msg)
}