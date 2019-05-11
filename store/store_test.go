package store

import (
	cosmosStore "github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

func setup() (ctx sdk.Context, key sdk.StoreKey) {
	db := db.NewMemDB()
	ms := cosmosStore.NewCommitMultiStore(db)

	ctx = sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	key = sdk.NewKVStoreKey("store")
	ms.MountStoreWithDB(key, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	return ctx, key
}
