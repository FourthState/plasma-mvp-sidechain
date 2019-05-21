package handlers

import (
	"github.com/FourthState/plasma-mvp-sidechain/msgs"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CreateZoneHandler(zoneStore store.ZoneStore) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		createZoneMsg, ok := msg.(msgs.CreateZoneMsg)
		if !ok {
			panic("Msg does not implement InitiatePresenceClaimMsg")
		}

		zone := store.Zone{
			ZoneID:  createZoneMsg.ZoneID,
			Beacons: createZoneMsg.Beacons,
			Geohash: createZoneMsg.Geohash,
		}
		zoneStore.StoreZone(ctx, zone)

		return sdk.Result{}
	}
}
