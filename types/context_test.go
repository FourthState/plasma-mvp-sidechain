package types

import (
	"testing"
	"github.com/stretchr/testify/assert"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/abci/types"

)

func TestContextWithTxIndexIncrement(t *testing.T) {
	ms, _ := setupMultiStore() // setupMultiStore() is in mapper_test.go
	ctx := sdk.NewContext(ms, abci.Header{}, false, nil)

	txIndex := GetTxIndex(ctx)
	assert.Equal(t, uint16(0), txIndex)

	ctx = WithTxIndex(ctx, txIndex)

	for i := 0; i < 10; i++ {
		txIndex = uint16(txIndex + 1)
		ctx = WithTxIndex(ctx, txIndex)
	}
	
	txIndex = GetTxIndex(ctx)
	assert.Equal(t, uint16(10), txIndex)
}


func TestContextWithTxIndex(t *testing.T) {
	ms, _ := setupMultiStore() 
	ctx := sdk.NewContext(ms, abci.Header{ChainID: "testingchain"}, false, nil)

	txIndex := GetTxIndex(ctx)
	assert.Equal(t, uint16(0), txIndex)

	ctx2 := WithTxIndex(ctx, txIndex + 1)
	assert.Equal(t, uint16(0), GetTxIndex(ctx))
	assert.Equal(t, uint16(1), GetTxIndex(ctx2))

}