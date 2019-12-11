package handlers

import (
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	cosmosStore "github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

/* This file contains helper functions for testing */

func setup() (sdk.Context, store.DataStore) {
	db := db.NewMemDB()
	ms := cosmosStore.NewCommitMultiStore(db)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())

	dataStoreKey := sdk.NewKVStoreKey("data")

	ms.MountStoreWithDB(dataStoreKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	dataStore := store.NewDataStore(dataStoreKey)
	return ctx, dataStore
}

func getPosition(posStr string) plasma.Position {
	pos, _ := plasma.FromPositionString(posStr)
	return pos
}
