package handlers

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = "handlers"

	CodeInsufficientFee              sdk.CodeType = 1
	CodeExitedInput                  sdk.CodeType = 2
	CodeSignatureVerificationFailure sdk.CodeType = 3
	CodeInvalidTransaction           sdk.CodeType = 4
	CodeInvalidSignature             sdk.CodeType = 5
	CodeInvalidInput                 sdk.CodeType = 6
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

func ErrInvalidTransaction(codespace sdk.CodespaceType, msg string, args ...interface{}) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidTransaction, msg, args)
}

func ErrInvalidSignature(codespace sdk.CodespaceType, msg string, args ...interface{}) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidSignature, msg, args)
}

func ErrInvalidInput(codespace sdk.CodespaceType, msg string, args ...interface{}) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInput, msg, args)
}
