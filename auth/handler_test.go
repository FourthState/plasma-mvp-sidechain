package auth

import (
	"crypto/ecdsa"
	"github.com/stretchr/testify/require"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	abci "github.com/tendermint/abci/types"
	"github.com/tendermint/tmlibs/log"

	db "github.com/FourthState/plasma-mvp-sidechain/db"
	types "github.com/FourthState/plasma-mvp-sidechain/types"
	utils "github.com/FourthState/plasma-mvp-sidechain/utils"
)

/// @param privA confirmSig Address
/// @param privB owner address
func NewUTXO(privA *ecdsa.PrivateKey, privB *ecdsa.PrivateKey, position types.Position) types.UTXO {
	addrA := utils.PrivKeyToAddress(privA)
	addrB := utils.PrivKeyToAddress(privB)
	confirmAddr := [2]common.Address{addrA, addrA}
	return types.NewBaseUTXO(addrB, confirmAddr, 100, position)
}

// Tests a valid spendmsg
// 2 different inputs and 2 different outputs
func TestHandleSpendMessage(t *testing.T) {
	ms, capKey := db.SetupMultiStore()

	ctx := sdk.NewContext(ms, abci.Header{Height: 2}, false, nil, log.NewNopLogger())
	mapper := db.NewUTXOMapper(capKey, db.MakeCodec())
	keeper := db.NewUTXOKeeper(mapper)
	txIndex := new(uint16)
	handler := NewHandler(keeper, txIndex)

	// Add in 2 parentUTXO
	privA, _ := ethcrypto.GenerateKey()
	privB, _ := ethcrypto.GenerateKey()
	privC, _ := ethcrypto.GenerateKey()
	positionB := types.Position{1000, 0, 0, 0}
	positionC := types.Position{1000, 1, 0, 0}
	utxo1 := NewUTXO(privA, privB, positionB)
	utxo2 := NewUTXO(privA, privC, positionC)
	mapper.AddUTXO(ctx, utxo1)
	mapper.AddUTXO(ctx, utxo2)
	utxo1 = mapper.GetUTXO(ctx, positionB)
	utxo2 = mapper.GetUTXO(ctx, positionC)
	require.NotNil(t, utxo1)
	require.NotNil(t, utxo2)

	newownerA := utils.GenerateAddress()
	newownerB := utils.GenerateAddress()
	confirmSigs := [2]types.Signature{types.Signature{}, types.Signature{}}

	// Add in SpendMsg,
	var msg = types.SpendMsg{
		Blknum1:      1000,
		Txindex1:     0,
		Oindex1:      0,
		DepositNum1:  0,
		Owner1:       utils.PrivKeyToAddress(privB),
		ConfirmSigs1: confirmSigs,
		Blknum2:      1000,
		Txindex2:     1,
		Oindex2:      0,
		DepositNum2:  0,
		Owner2:       utils.PrivKeyToAddress(privC),
		ConfirmSigs2: confirmSigs,
		Newowner1:    newownerA,
		Denom1:       150,
		Newowner2:    newownerB,
		Denom2:       50,
		Fee:          0,
	}

	res := handler(ctx, msg)
	require.Equal(t, sdk.CodeType(0), sdk.CodeType(res.Code), res.Log)

	require.Equal(t, uint16(1), *txIndex) // txIndex incremented

	//Check that inputs were deleted
	utxo := mapper.GetUTXO(ctx, positionB)
	require.Nil(t, utxo)
	utxo = mapper.GetUTXO(ctx, positionC)
	require.Nil(t, utxo)

	// Check to see if outputs were added
	require.Equal(t, int64(2), ctx.BlockHeight())
	positionD := types.Position{2, 0, 0, 0}
	positionE := types.Position{2, 0, 1, 0}
	utxo1 = mapper.GetUTXO(ctx, positionD)
	require.NotNil(t, utxo1)
	utxo2 = mapper.GetUTXO(ctx, positionE)
	require.NotNil(t, utxo2)

	// Check that outputs are valid
	inputAddresses := [2]common.Address{utils.PrivKeyToAddress(privB), utils.PrivKeyToAddress(privC)}
	require.Equal(t, uint64(150), utxo1.GetDenom())
	require.Equal(t, uint64(50), utxo2.GetDenom())
	require.EqualValues(t, newownerA, utxo1.GetAddress())
	require.EqualValues(t, newownerB, utxo2.GetAddress())
	require.EqualValues(t, inputAddresses, utxo1.GetInputAddresses())
	require.EqualValues(t, inputAddresses, utxo2.GetInputAddresses())
}

// Tests a valid spendmsg
// 1 input and 2 different outputs
func TestOneInput(t *testing.T) {
	ms, capKey := db.SetupMultiStore()

	ctx := sdk.NewContext(ms, abci.Header{Height: 2}, false, nil, log.NewNopLogger())
	mapper := db.NewUTXOMapper(capKey, db.MakeCodec())
	keeper := db.NewUTXOKeeper(mapper)
	txIndex := new(uint16)
	handler := NewHandler(keeper, txIndex)

	// Add in 2 parentUTXO
	privA, _ := ethcrypto.GenerateKey()
	privB, _ := ethcrypto.GenerateKey()
	positionB := types.Position{1000, 0, 0, 0}
	utxo1 := NewUTXO(privA, privB, positionB)
	mapper.AddUTXO(ctx, utxo1)
	utxo1 = mapper.GetUTXO(ctx, positionB)
	require.NotNil(t, utxo1)

	newownerA := utils.GenerateAddress()
	newownerB := utils.GenerateAddress()
	confirmSigs := [2]types.Signature{types.Signature{}, types.Signature{}}

	// Add in SpendMsg,
	var msg = types.SpendMsg{
		Blknum1:      1000,
		Txindex1:     0,
		Oindex1:      0,
		DepositNum1:  0,
		Owner1:       utils.PrivKeyToAddress(privB),
		ConfirmSigs1: confirmSigs,
		Blknum2:      0,
		Txindex2:     0,
		Oindex2:      0,
		DepositNum2:  0,
		Owner2:       common.Address{},
		ConfirmSigs2: confirmSigs,
		Newowner1:    newownerA,
		Denom1:       25,
		Newowner2:    newownerB,
		Denom2:       75,
		Fee:          0,
	}

	res := handler(ctx, msg)
	require.Equal(t, sdk.CodeType(0), sdk.CodeType(res.Code), res.Log)

	require.Equal(t, uint16(1), *txIndex) // txIndex incremented

	//Check that inputs were deleted
	utxo := mapper.GetUTXO(ctx, positionB)
	require.Nil(t, utxo)

	// Check to see if outputs were added
	require.Equal(t, int64(2), ctx.BlockHeight())
	positionD := types.Position{2, 0, 0, 0}
	positionE := types.Position{2, 0, 1, 0}
	utxo1 = mapper.GetUTXO(ctx, positionD)
	require.NotNil(t, utxo1)
	utxo2 := mapper.GetUTXO(ctx, positionE)
	require.NotNil(t, utxo2)

	// Check that outputs are valid
	require.Equal(t, uint64(25), utxo1.GetDenom())
	require.Equal(t, uint64(75), utxo2.GetDenom())
	require.EqualValues(t, newownerA, utxo1.GetAddress())
	require.EqualValues(t, newownerB, utxo2.GetAddress())
}
