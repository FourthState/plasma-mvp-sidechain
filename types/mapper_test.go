package types

import (
	"testing"
	"github.com/stretchr/testify/assert"

	dbm "github.com/tendermint/tmlibs/db"
	abci "github.com/tendermint/abci/types"
	crypto "github.com/tendermint/go-crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/store"

)

func setupMultiStore() (sdk.MultiStore, *sdk.KVStoreKey) {
	db := dbm.NewMemDB()
	capKey := sdk.NewKVStoreKey("capkey")
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(capKey, sdk.StoreTypeIAVL, db)
	ms.LoadLatestVersion()
	return ms, capKey
}

/*
	Basic test of Get and Add
	Creates a valid UTXO and adds it to the uxto mapping.
	Checks to make sure UTXO isn't nil after adding to mapping.
*/

func TestUTXOMapperGetAdd(t *testing.T) {
	ms, capKey := setupMultiStore()

	ctx := sdk.NewContext(ms, abci.Header{}, false, nil)
	mapper := NewUTXOMapper(capKey)

	privA := crypto.GenPrivKeySecp256k1()
	pubKeyA := privA.PubKey()
	addrA := pubKeyA.Address()

	privB := crypto.GenPrivKeySecp256k1()
	pubKeyB := privB.PubKey()
	addrB := pubKeyB.Address()

	positionB := Position{1000, 0, 0}
	confirmAddr := [2]crypto.Address{addrA, addrA}
	confirmPubKey := [2]crypto.PubKey{pubKeyA, pubKeyA}

	// These lines of code error. Why?
	//utxo := mapper.GetUXTO(ctx, positionB)
	//assert.Nil(t, utxo)

	utxo := NewBaseUTXO(addrB, confirmAddr, pubKeyB, confirmPubKey, 100, positionB)
	assert.NotNil(t, utxo)
	assert.Equal(t, addrB, utxo.GetAddress())
	assert.Equal(t, pubKeyB, utxo.GetPubKey())
	assert.EqualValues(t, confirmAddr, utxo.GetCSAddress())
	assert.EqualValues(t, confirmPubKey, utxo.GetCSPubKey())
	assert.EqualValues(t, positionB, utxo.GetPosition())

	mapper.AddUTXO(ctx, utxo)
	utxo = mapper.GetUTXO(ctx, positionB)
	assert.NotNil(t, utxo)
}