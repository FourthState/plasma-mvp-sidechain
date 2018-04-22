package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type contextKey int 

const (
	contextKeyTxIndex contextKey = iota
)

func WithTxIndex(ctx sdk.Context, txIndex uint16) sdk.Context {
	return ctx.WithValue(contextKeyTxIndex, txIndex)
}

func GetTxIndex(ctx sdk.Context) uint16 {
	v := ctx.Value(contextKeyTxIndex)
	if v == nil {
		return 0 // bug
	}
	return v.(uint16)
}