package utxo

import (
	"github.com/stretchr/testify/require"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"

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

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
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

	utxo = mapper.GetUTXO(ctx, addrB, positionB)
	require.NotNil(t, utxo)

	mapper.DeleteUTXO(ctx, addrB, positionB)
	utxo = mapper.GetUTXO(ctx, addrB, positionB)
	require.Nil(t, utxo)
}

/*
	Basic test of Multiple Additions and Deletes in the same block
	Creates a valid UTXOs and adds them to the uxto mapping.
	Then deletes the UTXO's from the mapping.
*/

func TestMultiUTXOAddDeleteSameBlock(t *testing.T) {
	ms, capKey := SetupMultiStore()

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
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

		utxo = mapper.GetUTXO(ctx, addrB, positionB)
		require.NotNil(t, utxo)
	}

	for i := 0; i < 10; i++ {
		position := types.Position{1000, uint16(i), 0, 0}

		utxo := mapper.GetUTXO(ctx, addrB, position)
		require.NotNil(t, utxo)
		mapper.DeleteUTXO(ctx, addrB, position)
		utxo = mapper.GetUTXO(ctx, addrB, position)
		require.Nil(t, utxo)
	}

}

func TestInvalidAddress(t *testing.T) {
	ms, capKey := SetupMultiStore()

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
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

	// GetUTXO with correct position but wrong address
	utxo = mapper.GetUTXO(ctx, addrA, positionB)
	require.Nil(t, utxo)

	utxo = mapper.GetUTXO(ctx, addrB, positionB)
	require.NotNil(t, utxo)

	// DeleteUTXO with correct position but wrong address
	mapper.DeleteUTXO(ctx, addrA, positionB)
	utxo = mapper.GetUTXO(ctx, addrB, positionB)
	require.NotNil(t, utxo)

	mapper.DeleteUTXO(ctx, addrB, positionB)
	utxo = mapper.GetUTXO(ctx, addrB, positionB)
	require.Nil(t, utxo)
}

/*
	Basic test of Multiple Additions and Deletes in the different block
	Creates a valid UTXOs and adds them to the uxto mapping.
	Then deletes the UTXO's from the mapping.
*/

func TestMultiUTXOAddDeleteDifferentBlock(t *testing.T) {
	ms, capKey := SetupMultiStore()

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
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

		utxo = mapper.GetUTXO(ctx, addrB, positionB)
		require.NotNil(t, utxo)
	}

	for i := 0; i < 10; i++ {
		position := types.Position{uint64(i), 0, 0, 0}

		utxo := mapper.GetUTXO(ctx, addrB, position)
		require.NotNil(t, utxo)
		mapper.DeleteUTXO(ctx, addrB, position)
		utxo = mapper.GetUTXO(ctx, addrB, position)
		require.Nil(t, utxo)
	}

}

/*
	Test getting all UTXOs for an Address.
*/

func TestGetUTXOsForAddress(t *testing.T) {
	ms, capKey := SetupMultiStore()

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	mapper := NewUTXOMapper(capKey, MakeCodec())

	privA, _ := ethcrypto.GenerateKey()
	addrA := utils.PrivKeyToAddress(privA)

	privB, _ := ethcrypto.GenerateKey()
	addrB := utils.PrivKeyToAddress(privB)

	privC, _ := ethcrypto.GenerateKey()
	addrC := utils.PrivKeyToAddress(privC)

	positionB1 := types.Position{1000, 0, 0, 0}
	positionB2 := types.Position{1001, 1, 0, 0}
	positionB3 := types.Position{1002, 2, 1, 0}
	confirmAddr := [2]common.Address{addrA, addrA}

	utxo1 := types.NewBaseUTXO(addrB, confirmAddr, 100, positionB1)
	utxo2 := types.NewBaseUTXO(addrB, confirmAddr, 200, positionB2)
	utxo3 := types.NewBaseUTXO(addrB, confirmAddr, 300, positionB3)

	mapper.AddUTXO(ctx, utxo1)
	mapper.AddUTXO(ctx, utxo2)
	mapper.AddUTXO(ctx, utxo3)

	utxosForAddressB := mapper.GetUTXOsForAddress(ctx, addrB)
	require.NotNil(t, utxosForAddressB)
	require.Equal(t, 3, len(utxosForAddressB))
	require.Equal(t, utxo1, utxosForAddressB[0])
	require.Equal(t, utxo2, utxosForAddressB[1])
	require.Equal(t, utxo3, utxosForAddressB[2])

	positionC1 := types.Position{1002, 3, 0, 0}
	utxo4 := types.NewBaseUTXO(addrC, confirmAddr, 300, positionC1)
	mapper.AddUTXO(ctx, utxo4)
	utxosForAddressC := mapper.GetUTXOsForAddress(ctx, addrC)
	require.NotNil(t, utxosForAddressC)
	require.Equal(t, 1, len(utxosForAddressC))
	require.Equal(t, utxo4, utxosForAddressC[0])

	// check returns empty slice if no UTXOs exist for address
	utxosForAddressA := mapper.GetUTXOsForAddress(ctx, addrA)
	require.Empty(t, utxosForAddressA)
}
