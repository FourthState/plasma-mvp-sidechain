package store

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = "store"

	CodeOutputDNE   sdk.CodeType = 1
	CodeOutputSpent sdk.CodeType = 2
	CodeAccountDNE  sdk.CodeType = 3
	CodeTxDNE       sdk.CodeType = 4
)

func ErrOutputDNE(msg string, args ...interface{}) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeOutputDNE, msg, args)
}

func ErrOutputSpent(msg string, args ...interface{}) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeOutputSpent, msg, args)
}

func ErrAccountDNE(msg string, args ...interface{}) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeAccountDNE, msg, args)
}

func ErrTxDNE(msg string, args ...interface{}) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeTxDNE, msg, args)
}
