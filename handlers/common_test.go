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

func setup() (sdk.Context, store.OutputStore, store.BlockStore) {
	db := db.NewMemDB()
	ms := cosmosStore.NewCommitMultiStore(db)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())

	blockStoreKey := sdk.NewKVStoreKey("block")
	outputStoreKey := sdk.NewKVStoreKey("output")

	ms.MountStoreWithDB(blockStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(outputStoreKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	blockStore := store.NewBlockStore(blockStoreKey)
	outputStore := store.NewOutputStore(outputStoreKey)

	return ctx, outputStore, blockStore
}

func getPosition(posStr string) plasma.Position {
	pos, _ := plasma.FromPositionString(posStr)
	return pos
}
