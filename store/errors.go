package store

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = "store"

	CodeOutputDNE   sdk.CodeType = 1
	CodeOutputSpent sdk.CodeType = 2
)

func ErrOutputSpent(msg string, args ...interface{}) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeOutputSpent, msg, args)
}

func ErrOutputDNE(msg string, args ...interface{}) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeOutputDNE, msg, args)
}
