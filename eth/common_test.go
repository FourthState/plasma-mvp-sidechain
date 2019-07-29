package eth

import (
	"github.com/FourthState/plasma-mvp-sidechain/store"
	cosmosStore "github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

func setup() (sdk.Context, store.DataStore) {
	db := db.NewMemDB()
	ms := cosmosStore.NewCommitMultiStore(db)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())

	dsKey := sdk.NewKVStoreKey(store.DataStoreName)

	ms.MountStoreWithDB(dsKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	ds := store.NewDataStore(dsKey)

	return ctx, ds
}
