package query

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = "store/query"

	CodeSerialization sdk.CodeType = 4
	CodeInvalidPath   sdk.CodeType = 5
	CodeTxDNE         sdk.CodeType = 6
	CodeWalletDNE     sdk.CodeType = 7
	CodeOutputDNE     sdk.CodeType = 8
)

func ErrSerialization(msg string, args ...interface{}) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeSerialization, msg, args)
}

func ErrInvalidPath(msg string, args ...interface{}) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeInvalidPath, msg, args)
}

// TODO: refactor DNE errors into one
func ErrOutputDNE(msg string, args ...interface{}) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeOutputDNE, msg, args)
}

func ErrTxDNE(msg string, args ...interface{}) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeTxDNE, msg, args)
}

func ErrWalletDNE(msg string, args ...interface{}) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeWalletDNE, msg, args)
}
