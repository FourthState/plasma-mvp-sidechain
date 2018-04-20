package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type contextKey int // not sure if needed

const (
	contextKeyTxIndex contextKey = iota
)

func WithTxIndex(ctx sdk.Context, uint16 txIndex) sdk.Context {
	return ctx.WithValue(contextKeyTxIndex, txIndex)
}

func GetTxIndex(ctx sdk.Context) uint16 {
	v := ctx.Value(contextKeyTxIndex)
	if v == nil {
		return 0 // bug
	}
	return v.(uint16)
}