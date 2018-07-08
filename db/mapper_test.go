package db

import (
	"github.com/stretchr/testify/require"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	abci "github.com/tendermint/abci/types"
	"github.com/tendermint/tmlibs/log"

	"github.com/FourthState/plasma-mvp-sidechain/types"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
)

/*
	Basic test of Get, Add, Delete
	Creates a valid UTXO and adds it to the uxto mapping.
	Checks to make sure UTXO isn't nil after adding to mapping.
	Then deletes the UTXO from the mapping
*/

func TestUTXOGetAddDelete(t *testing.T) {
	ms, capKey := SetupMultiStore()

	ctx := sdk.NewContext(ms, abci.Header{}, false, nil, log.NewNopLogger())
	mapper := NewUTXOMapper(capKey, MakeCodec())

	privA, _ := ethcrypto.GenerateKey()
	addrA := utils.PrivKeyToAddress(privA)

	privB, _ := ethcrypto.GenerateKey()
	addrB := utils.PrivKeyToAddress(privB)

	positionB := types.Position{1000, 0, 0, 0}
	confirmAddr := [2]common.Address{addrA, addrA}

	utxo := types.NewBaseUTXO(addrB, confirmAddr, 100, positionB)
	require.NotNil(t, utxo)
	require.Equal(t, addrB, utxo.GetAddress())
	require.EqualValues(t, confirmAddr, utxo.GetInputAddresses())
	require.EqualValues(t, positionB, utxo.GetPosition())

	mapper.AddUTXO(ctx, utxo)
	utxo = mapper.GetUTXO(ctx, positionB)
	require.NotNil(t, utxo)

	mapper.DeleteUTXO(ctx, positionB)
	utxo = mapper.GetUTXO(ctx, positionB)
	require.Nil(t, utxo)
}

/*
	Basic test of Multiple Additions and Deletes in the same block
	Creates a valid UTXOs and adds them to the uxto mapping.
	Then deletes the UTXO's from the mapping.
*/

func TestMultiUTXOAddDeleteSameBlock(t *testing.T) {
	ms, capKey := SetupMultiStore()

	ctx := sdk.NewContext(ms, abci.Header{}, false, nil, log.NewNopLogger())
	mapper := NewUTXOMapper(capKey, MakeCodec())

	// These are not being tested
	privA, _ := ethcrypto.GenerateKey()
	addrA := utils.PrivKeyToAddress(privA)

	privB, _ := ethcrypto.GenerateKey()
	addrB := utils.PrivKeyToAddress(privB)

	confirmAddr := [2]common.Address{addrA, addrA}

	// Main part being tested
	for i := 0; i < 10; i++ {
		positionB := types.Position{1000, uint16(i), 0, 0}
		utxo := types.NewBaseUTXO(addrB, confirmAddr, 100, positionB)
		mapper.AddUTXO(ctx, utxo)
		utxo = mapper.GetUTXO(ctx, positionB)
		require.NotNil(t, utxo)
	}

	for i := 0; i < 10; i++ {
		position := types.Position{1000, uint16(i), 0, 0}
		utxo := mapper.GetUTXO(ctx, position)
		require.NotNil(t, utxo)
		mapper.DeleteUTXO(ctx, position)
		utxo = mapper.GetUTXO(ctx, position)
		require.Nil(t, utxo)
	}

}

/*
	Basic test of Multiple Additions and Deletes in the different block
	Creates a valid UTXOs and adds them to the uxto mapping.
	Then deletes the UTXO's from the mapping.
*/

func TestMultiUTXOAddDeleteDifferentBlock(t *testing.T) {
	ms, capKey := SetupMultiStore()

	ctx := sdk.NewContext(ms, abci.Header{}, false, nil, log.NewNopLogger())
	mapper := NewUTXOMapper(capKey, MakeCodec())

	// These are not being tested
	privA, _ := ethcrypto.GenerateKey()
	addrA := utils.PrivKeyToAddress(privA)

	privB, _ := ethcrypto.GenerateKey()
	addrB := utils.PrivKeyToAddress(privB)

	confirmAddr := [2]common.Address{addrA, addrA}

	// Main part being tested
	for i := 0; i < 10; i++ {
		positionB := types.Position{uint64(i), 0, 0, 0}
		utxo := types.NewBaseUTXO(addrB, confirmAddr, 100, positionB)
		mapper.AddUTXO(ctx, utxo)
		utxo = mapper.GetUTXO(ctx, positionB)
		require.NotNil(t, utxo)
	}

	for i := 0; i < 10; i++ {
		position := types.Position{uint64(i), 0, 0, 0}
		utxo := mapper.GetUTXO(ctx, position)
		require.NotNil(t, utxo)
		mapper.DeleteUTXO(ctx, position)
		utxo = mapper.GetUTXO(ctx, position)
		require.Nil(t, utxo)
	}

}
