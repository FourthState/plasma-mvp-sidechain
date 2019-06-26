package store

import (
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	cosmosStore "github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

/* This file contains helper functions for testing */

func setup() (ctx sdk.Context, key sdk.StoreKey) {
	db := db.NewMemDB()
	ms := cosmosStore.NewCommitMultiStore(db)

	ctx = sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	key = sdk.NewKVStoreKey("store")
	ms.MountStoreWithDB(key, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	return ctx, key
}

func getPosition(posStr string) plasma.Position {
	pos, _ := plasma.FromPositionString(posStr)
	return pos
}
