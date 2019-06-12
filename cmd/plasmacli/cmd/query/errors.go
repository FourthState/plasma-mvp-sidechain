package query

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = "query"

	CodeInvalidArg    = 1
	CodeSerialization = 2
)

func ErrInvalidArg(msg string, args ...interface{}) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeInvalidArg, msg, args)
}

func ErrSerialization(msg string, args ...interface{}) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeSerialization, msg, args)
}
