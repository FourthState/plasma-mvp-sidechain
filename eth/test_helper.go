package eth

import (
	"github.com/FourthState/plasma-mvp-sidechain/store"
	cosmosStore "github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

func setup() (sdk.Context, store.BlockStore) {
	db := db.NewMemDB()
	ms := cosmosStore.NewCommitMultiStore(db)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())

	blockStoreKey := sdk.NewKVStoreKey("block")

	ms.MountStoreWithDB(blockStoreKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	blockStore := store.NewBlockStore(blockStoreKey)

	return ctx, blockStore
}
