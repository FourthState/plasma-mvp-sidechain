package db

import (
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/cosmos/cosmos-sdk/store"
	types "plasma-mvp-sidechain/types"
	"github.com/tendermint/go-amino"
	sdk "github.com/cosmos/cosmos-sdk/types"
	crypto "github.com/tendermint/go-crypto"
)

func SetupMultiStore() (sdk.MultiStore, *sdk.KVStoreKey) {
	db := dbm.NewMemDB()
	capKey := sdk.NewKVStoreKey("capkey")
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(capKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()
	return ms, capKey
}

func MakeCodec() *amino.Codec {
	cdc := amino.NewCodec()
	types.RegisterAmino(cdc)
	crypto.RegisterAmino(cdc)
	return cdc
}
