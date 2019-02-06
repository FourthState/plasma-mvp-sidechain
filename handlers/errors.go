package handlers

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = "handlers"

	CodeInsufficientFee              sdk.CodeType = 1
	CodeExitedInput                  sdk.CodeType = 2
	CodeSignatureVerificationFailure sdk.CodeType = 3
)

func ErrInsufficientFee(codespace sdk.CodespaceType, msg string, args ...interface{}) sdk.Error {
	return sdk.NewError(codespace, CodeInsufficientFee, msg, args)
}

func ErrExitedInput(codespace sdk.CodespaceType, msg string, args ...interface{}) sdk.Error {
	return sdk.NewError(codespace, CodeExitedInput, msg, args)
}

func ErrSignatureVerificationFailure(codespace sdk.CodespaceType, msg string, args ...interface{}) sdk.Error {
	return sdk.NewError(codespace, CodeSignatureVerificationFailure, msg, args)
}
