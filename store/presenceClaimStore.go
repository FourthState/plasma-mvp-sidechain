package store

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/msgs"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"io"
)

type PresenceClaim struct {
	ZoneID       [32]byte        `json:"zoneID`
	UTXOPosition plasma.Position `json:"utxoPosition"`
	UserAddress  common.Address  `json:"userAddress"`
	LogsHash     []byte          `json:"logsHash"`
}

func (claim *PresenceClaim) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, &claim)
}

func (claim *PresenceClaim) DecodeRLP(s *rlp.Stream) error {
	return s.Decode(&claim)

}

type PresenceClaimStore struct {
	KVStore
}

func NewPresenceClaimStore(ctxKey sdk.StoreKey) PresenceClaimStore {
	return PresenceClaimStore{
		KVStore: NewKVStore(ctxKey),
	}
}

func (store PresenceClaimStore) GetPresenceClaim(ctx sdk.Context, key []byte) (PresenceClaim, bool) {
	data := store.Get(ctx, key)
	if data == nil {
		return PresenceClaim{}, false
	}

	var claim PresenceClaim
	if err := rlp.DecodeBytes(data, &claim); err != nil {
		panic(fmt.Sprintf("PresenceClaim store corrupted: %s", err))
	}

	return claim, true
}

func (store PresenceClaimStore) HasPresenceClaim(ctx sdk.Context, key []byte) bool {
	return store.Has(ctx, key)
}

func (store PresenceClaimStore) StorePresenceClaim(ctx sdk.Context, claim PresenceClaim) {

	messageNoSig := msgs.InitiatePresenceClaimMsg{}
	messageNoSig.ZoneID = claim.ZoneID
	messageNoSig.UTXOPosition = claim.UTXOPosition

	claimHash := messageNoSig.TxHash()
	data, err := rlp.EncodeToBytes(&claim)
	if err != nil {
		panic(fmt.Sprintf("Error marshaling utxo: %s", err))
	}

	store.Set(ctx, claimHash, data)
}
