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

func ErrInsufficientFee(msg string, args ...interface{}) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeInsufficientFee, msg, args)
}

func ErrExitedInput(msg string, args ...interface{}) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeExitedInput, msg, args)
}

func ErrSignatureVerificationFailure(msg string, args ...interface{}) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeSignatureVerificationFailure, msg, args)
}

func ErrInvalidTransaction(msg string, args ...interface{}) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeInvalidTransaction, msg, args)
}

func ErrInvalidSignature(msg string, args ...interface{}) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeInvalidSignature, msg, args)
}

func ErrInvalidInput(msg string, args ...interface{}) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeInvalidInput, msg, args)
}
