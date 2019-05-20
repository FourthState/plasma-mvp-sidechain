package store

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

type Zone struct {
	ZoneID  []byte           `json:"zoneID"`
	Beacons []common.Address `json:"beacons"`
	Geohash []string         `json:"geohash"`
}

type ZoneStore struct {
	KVStore
}

func NewZoneStore(ctxKey sdk.StoreKey) ZoneStore {
	return ZoneStore{
		KVStore: NewKVStore(ctxKey),
	}
}

func (store ZoneStore) GetZoneByID(ctx sdk.Context, key []byte) (Zone, bool) {
	data := store.Get(ctx, key)
	if data == nil {
		return Zone{}, false
	}

	var zone Zone
	if err := rlp.DecodeBytes(data, &zone); err != nil {
		panic(fmt.Sprintf("Zone store corrupted: %s", err))
	}

	return zone, true
}

func (store ZoneStore) GetZonesByAddress(ctx sdk.Context, key []byte) (Zone, bool) {
	data := store.Get(ctx, key)
	if data == nil {
		return Zone{}, false
	}

	var zone Zone
	if err := rlp.DecodeBytes(data, &zone); err != nil {
		panic(fmt.Sprintf("Zone store corrupted: %s", err))
	}

	return zone, true
}

func (store ZoneStore) StoreZone(ctx sdk.Context, zone Zone) {

	data, err := rlp.EncodeToBytes(&zone)
	if err != nil {
		panic(fmt.Sprintf("Error marshaling zone: %s", err))
	}

	store.Set(ctx, zone.ZoneID, data)

	for _, a := range zone.Beacons {
		store.Set(ctx, a.Bytes(), data)
	}
}
