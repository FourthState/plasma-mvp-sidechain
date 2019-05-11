package store

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = "store"

	CodeOutputDNE   sdk.CodeType = 1
	CodeOutputSpent sdk.CodeType = 2
)

func ErrOutputDNE(codespace sdk.CodespaceType, msg string, args ...interface{}) sdk.Error {
	return sdk.NewError(codespace, CodeOutputDNE, msg, args)
}

func ErrOutputSpent(codespace sdk.CodespaceType, msg string, args ...interface{}) sdk.Error {
	return sdk.NewError(codespace, CodeOutputSpent, msg, args)
}
