package query

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	DefaultCodespace sdk.CodespaceType = "store/query"

	CodeSerialization sdk.CodeType = 4
	CodeInvalidPath   sdk.CodeType = 5
)

func ErrSerialization(msg string, args ...interface{}) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeSerialization, msg, args)
}

func ErrInvalidPath(msg string, args ...interface{}) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeInvalidPath, msg, args)
}
