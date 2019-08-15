package handlers

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Error codes for handlers
const (
	DefaultCodespace sdk.CodespaceType = "handlers"

	CodeInsufficientFee              sdk.CodeType = 1
	CodeExitedInput                  sdk.CodeType = 2
	CodeSignatureVerificationFailure sdk.CodeType = 3
	CodeInvalidTransaction           sdk.CodeType = 4
	CodeInvalidSignature             sdk.CodeType = 5
	CodeInvalidInput                 sdk.CodeType = 6
)

// ErrInsufficientFee error for an insufficient fee
func ErrInsufficientFee(msg string, args ...interface{}) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeInsufficientFee, msg, args)
}

// ErrExitedInput error for if the input has already exited
func ErrExitedInput(msg string, args ...interface{}) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeExitedInput, msg, args)
}

// ErrSignatureVerficiationFailure error for signature verifcation failing to
// complete
func ErrSignatureVerificationFailure(msg string, args ...interface{}) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeSignatureVerificationFailure, msg, args)
}

// ErrInvalidTransaction error for an invalid transaction
func ErrInvalidTransaction(msg string, args ...interface{}) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeInvalidTransaction, msg, args)
}

// ErrInvalidSignature error for an incorrect signature
func ErrInvalidSignature(msg string, args ...interface{}) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeInvalidSignature, msg, args)
}

// ErrInvalidInput
func ErrInvalidInput(msg string, args ...interface{}) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeInvalidInput, msg, args)
}
