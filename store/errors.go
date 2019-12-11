package store

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Error codes for data store
const (
	DefaultCodespace sdk.CodespaceType = "store"

	CodeDNE         sdk.CodeType = 1
	CodeOutputSpent sdk.CodeType = 2
	CodeInvalidPath sdk.CodeType = 3
)

// ErrDNE error for an object that does not exist
func ErrDNE(msg string, args ...interface{}) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeDNE, msg, args)
}

// ErrOutputSpent error for an output that is marked as spent
func ErrOutputSpent(msg string, args ...interface{}) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeOutputSpent, msg, args)
}

// ErrInvalidPath error for an invalid query path
func ErrInvalidPath(msg string, args ...interface{}) sdk.Error {
	return sdk.NewError(DefaultCodespace, CodeInvalidPath, msg, args)
}
