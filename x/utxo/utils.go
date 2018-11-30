package utxo

import (
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/go-amino"
	cryptoAmino "github.com/tendermint/tendermint/crypto/encoding/amino"
	dbm "github.com/tendermint/tendermint/libs/db"
)

func SetupMultiStore() (sdk.MultiStore, *sdk.KVStoreKey, *sdk.KVStoreKey) {
	db := dbm.NewMemDB()
	capKey := sdk.NewKVStoreKey("capkey")
	capKey2 := sdk.NewKVStoreKey("capkey2")
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(capKey, sdk.StoreTypeIAVL, db)
	ms.MountStoreWithDB(capKey2, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()
	return ms, capKey, capKey2
}

func MakeCodec() *amino.Codec {
	cdc := amino.NewCodec()
	RegisterAmino(cdc)
	cryptoAmino.RegisterAmino(cdc)
	return cdc
}

func RegisterAmino(cdc *amino.Codec) {
	cdc.RegisterInterface((*Position)(nil), nil)
}
