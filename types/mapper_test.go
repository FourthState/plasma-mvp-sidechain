package types

import (
	"testing"
	"github.com/stretchr/testify/assert"

	dbm "github.com/tendermint/tmlibs/db"
	abci "github.com/tendermint/abci/types"
	crypto "github.com/tendermint/go-crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/store"
	"github.com/tendermint/go-amino" 

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
	Basic test of Get, Add, Delete
	Creates a valid UTXO and adds it to the uxto mapping.
	Checks to make sure UTXO isn't nil after adding to mapping.
	Then deletes the UTXO from the mapping
*/

func TestUTXOGetAddDelete(t *testing.T) {
	ms, capKey := setupMultiStore()

	ctx := sdk.NewContext(ms, abci.Header{}, false, nil)
	mapper := NewUTXOMapper(capKey, MakeCodec())

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

	mapper.DeleteUTXO(ctx, positionB)
	utxo = mapper.GetUTXO(ctx, positionB)
	assert.Nil(t, utxo)
}



/*
	Basic test of Multiple Additions and Deletes in the same block
	Creates a valid UTXOs and adds them to the uxto mapping.
	Then deletes the UTXO's from the mapping.
*/

func TestMultiUTXOAddDeleteSameBlock(t *testing.T) {
	ms, capKey := setupMultiStore()

	ctx := sdk.NewContext(ms, abci.Header{}, false, nil)
	mapper := NewUTXOMapper(capKey, MakeCodec())

	// These are not being tested
	privA := crypto.GenPrivKeySecp256k1()
	pubKeyA := privA.PubKey()
	addrA := pubKeyA.Address()

	privB := crypto.GenPrivKeySecp256k1()
	pubKeyB := privB.PubKey()
	addrB := pubKeyB.Address()

	confirmAddr := [2]crypto.Address{addrA, addrA}
	confirmPubKey := [2]crypto.PubKey{pubKeyA, pubKeyA}

	// Main part being tested
	for i := 0; i < 10; i++ {
		positionB := Position{1000, uint16(i), 0}
		utxo := NewBaseUTXO(addrB, confirmAddr, pubKeyB, confirmPubKey, 100, positionB)
		mapper.AddUTXO(ctx, utxo)
		utxo = mapper.GetUTXO(ctx, positionB)
		assert.NotNil(t, utxo)
	}

	for i := 0; i < 10; i++ {
		position := Position{1000, uint16(i), 0}
		utxo := mapper.GetUTXO(ctx, position)
		assert.NotNil(t, utxo)
		mapper.DeleteUTXO(ctx, position)
		utxo = mapper.GetUTXO(ctx, position)
		assert.Nil(t, utxo)
	}

}

/*
	Basic test of Multiple Additions and Deletes in the different block
	Creates a valid UTXOs and adds them to the uxto mapping.
	Then deletes the UTXO's from the mapping.
*/

func TestMultiUTXOAddDeleteDifferentBlock(t *testing.T) {
	ms, capKey := setupMultiStore()

	ctx := sdk.NewContext(ms, abci.Header{}, false, nil)
	mapper := NewUTXOMapper(capKey, MakeCodec())

	// These are not being tested
	privA := crypto.GenPrivKeySecp256k1()
	pubKeyA := privA.PubKey()
	addrA := pubKeyA.Address()

	privB := crypto.GenPrivKeySecp256k1()
	pubKeyB := privB.PubKey()
	addrB := pubKeyB.Address()

	confirmAddr := [2]crypto.Address{addrA, addrA}
	confirmPubKey := [2]crypto.PubKey{pubKeyA, pubKeyA}

	// Main part being tested
	for i := 0; i < 10; i++ {
		positionB := Position{uint64(1000 * i), 0, 0}
		utxo := NewBaseUTXO(addrB, confirmAddr, pubKeyB, confirmPubKey, 100, positionB)
		mapper.AddUTXO(ctx, utxo)
		utxo = mapper.GetUTXO(ctx, positionB)
		assert.NotNil(t, utxo)
	}

	for i := 0; i < 10; i++ {
		position := Position{uint64(1000 * i), 0, 0}
		utxo := mapper.GetUTXO(ctx, position)
		assert.NotNil(t, utxo)
		mapper.DeleteUTXO(ctx, position)
		utxo = mapper.GetUTXO(ctx, position)
		assert.Nil(t, utxo)
	}

}

func MakeCodec() *amino.Codec {
	cdc := amino.NewCodec()
	cdc.RegisterInterface((*sdk.Msg)(nil), nil)
	RegisterAmino(cdc)   
	crypto.RegisterAmino(cdc)
	return cdc
}