package handlers

import (
	"github.com/FourthState/plasma-mvp-sidechain/store"
	cosmosStore "github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
)

func setup() (sdk.Context, store.UTXOStore, store.PlasmaStore, store.PresenceClaimStore) {
	db := db.NewMemDB()
	ms := cosmosStore.NewCommitMultiStore(db)

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())

	plasmaStoreKey := sdk.NewKVStoreKey("plasma")
  presenceClaimStoreKey := sdk.NewKVStoreKey("presenceClaim")
	utxoStoreKey := sdk.NewKVStoreKey("utxo")

	ms.MountStoreWithDB(plasmaStoreKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(utxoStoreKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()

	plasmaStore := store.NewPlasmaStore(plasmaStoreKey)
	utxoStore := store.NewUTXOStore(utxoStoreKey)
	claimStore := store.NewPresenceClaimStore(presenceClaimStoreKey)

	return ctx, utxoStore, plasmaStore, claimStore
}
