package store

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = "store"

	CodeDNE         sdk.CodeType = 1
	CodeOutputSpent sdk.CodeType = 2
	CodeInvalidPath sdk.CodeType = 3
)

func ErrDNE(msg string, args ...interface{}) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeDNE, msg, args)
}

func ErrOutputSpent(msg string, args ...interface{}) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeOutputSpent, msg, args)
}

func ErrInvalidPath(msg string, args ...interface{}) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeInvalidPath, msg, args)
}
