package types

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/abci/types"
	"github.com/tendermint/go-amino"
	crypto "github.com/tendermint/go-crypto"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
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

	ctx := sdk.NewContext(ms, abci.Header{}, false, nil, log.NewNopLogger())
	mapper := NewUTXOMapper(capKey, MakeCodec())

	privA, _ := ethcrypto.GenerateKey()
	addrA := EthPrivKeyToSDKAddress(privA)

	privB, _ := ethcrypto.GenerateKey()
	addrB := EthPrivKeyToSDKAddress(privB)

	positionB := Position{1000, 0, 0, 0}
	confirmAddr := [2]crypto.Address{addrA, addrA}

	// These lines of code error. Why?
	//utxo := mapper.GetUXTO(ctx, positionB)
	//assert.Nil(t, utxo)

	utxo := NewBaseUTXO(addrB, confirmAddr, 100, positionB)
	assert.NotNil(t, utxo)
	assert.Equal(t, addrB, utxo.GetAddress())
	assert.EqualValues(t, confirmAddr, utxo.GetInputAddresses())
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

	ctx := sdk.NewContext(ms, abci.Header{}, false, nil, log.NewNopLogger())
	mapper := NewUTXOMapper(capKey, MakeCodec())

	// These are not being tested
	privA, _ := ethcrypto.GenerateKey()
	addrA := EthPrivKeyToSDKAddress(privA)

	privB, _ := ethcrypto.GenerateKey()
	addrB := EthPrivKeyToSDKAddress(privB)


	confirmAddr := [2]crypto.Address{addrA, addrA}

	// Main part being tested
	for i := 0; i < 10; i++ {
		positionB := Position{1000, uint16(i), 0, 0}
		utxo := NewBaseUTXO(addrB, confirmAddr, 100, positionB)
		mapper.AddUTXO(ctx, utxo)
		utxo = mapper.GetUTXO(ctx, positionB)
		assert.NotNil(t, utxo)
	}

	for i := 0; i < 10; i++ {
		position := Position{1000, uint16(i), 0, 0}
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

	ctx := sdk.NewContext(ms, abci.Header{}, false, nil, log.NewNopLogger())
	mapper := NewUTXOMapper(capKey, MakeCodec())

	// These are not being tested
	privA, _ := ethcrypto.GenerateKey()
	addrA := EthPrivKeyToSDKAddress(privA)

	privB, _ := ethcrypto.GenerateKey()
	addrB := EthPrivKeyToSDKAddress(privB)


	confirmAddr := [2]crypto.Address{addrA, addrA}

	// Main part being tested
	for i := 0; i < 10; i++ {
		positionB := Position{uint64(i), 0, 0, 0}
		utxo := NewBaseUTXO(addrB, confirmAddr, 100, positionB)
		mapper.AddUTXO(ctx, utxo)
		utxo = mapper.GetUTXO(ctx, positionB)
		assert.NotNil(t, utxo)
	}

	for i := 0; i < 10; i++ {
		position := Position{uint64(i), 0, 0, 0}
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
