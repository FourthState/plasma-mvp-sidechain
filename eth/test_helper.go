package eth

import (
	"github.com/FourthState/plasma-mvp-sidechain/store"
	cosmosStore "github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

func setup() (sdk.Context, store.PlasmaStore) {
	db := db.NewMemDB()
	ms := cosmosStore.NewCommitMultiStore(db)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())

	plasmaStoreKey := sdk.NewKVStoreKey("plasma")

	ms.MountStoreWithDB(plasmaStoreKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	plasmaStore := store.NewPlasmaStore(plasmaStoreKey)

	return ctx, plasmaStore
}
