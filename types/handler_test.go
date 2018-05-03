package types 

import (
	"testing"
	"github.com/stretchr/testify/assert"

	abci "github.com/tendermint/abci/types"
	crypto "github.com/tendermint/go-crypto"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewUTXO(privA crypto.PrivKey, privB crypto.PrivKey, position Position) UTXO {
	addrA := privA.PubKey().Address()
	addrB := privB.PubKey().Address()
	confirmAddr := [2]crypto.Address{addrA, addrA}
	return NewBaseUTXO(addrB, confirmAddr, 100, position)
}

/*
	Tests a valid spendmsg 
	2 different inputs and 2 different outputs
	Inputs are from the same block
*/
func TestHandleSpendMessage(t *testing.T) {
	ms, capKey := setupMultiStore()

	ctx := sdk.NewContext(ms, abci.Header{Height: 2}, false, nil)
	mapper := NewUTXOMapper(capKey, MakeCodec())
	keeper := NewUTXOKeeper(mapper)
	txIndex := new(uint16)
	handler := NewHandler(keeper, txIndex)

	// Add in 2 parentUTXO
	privA := crypto.GenPrivKeySecp256k1()
	privB := crypto.GenPrivKeySecp256k1()
	privC := crypto.GenPrivKeySecp256k1()
	positionB := Position{1000, 0, 0}
	positionC := Position{1000, 1, 0}
	utxo1 := NewUTXO(privA, privB, positionB)
	utxo2 := NewUTXO(privA, privC, positionC)
	mapper.AddUTXO(ctx, utxo1)
	mapper.AddUTXO(ctx, utxo2)
	utxo1 = mapper.GetUTXO(ctx, positionB)
	utxo2 = mapper.GetUTXO(ctx, positionC)
	assert.NotNil(t, utxo1)
	assert.NotNil(t, utxo2)

	newownerA := crypto.Address([]byte("newownerA"))
	newownerB := crypto.Address([]byte("newownerB"))
	confrimSigs := [2]crypto.Signature{crypto.SignatureSecp256k1{}, crypto.SignatureSecp256k1{}}

	// Add in SpendMsg,
	var msg = SpendMsg{
		Blknum1: 		1000,
		Txindex1: 		0,
		Oindex1: 		0,
		Indenom1: 		100,
		Owner1: 		privB.PubKey().Address(),
		ConfirmSigs1: 	confrimSigs,
		Blknum2:		1000,
		Txindex2: 		1,
		Oindex2: 		0,
		Indenom2: 		100,
		Owner2: 		privC.PubKey().Address(),
		ConfirmSigs2: 	confrimSigs,
		Newowner1: 		newownerA,
		Denom1: 		150,
		Newowner2: 		newownerB,
		Denom2: 		50,
		Fee: 			0,
	}

	res := handler(ctx, msg)
	assert.Equal(t, sdk.CodeType(0), sdk.CodeType(res.Code), res.Log)

	assert.Equal(t, uint16(1), *txIndex) // txIndex incremented
	
	//Check that inputs were deleted
	utxo := mapper.GetUTXO(ctx, positionB)
	assert.Nil(t, utxo)
	utxo = mapper.GetUTXO(ctx, positionC)
	assert.Nil(t, utxo)

	// Check to see if outputs were added
	assert.Equal(t, int64(2), ctx.BlockHeight())
	positionD := Position{2000, 0, 0}
	positionE := Position{2000, 0, 1}
	utxo1 = mapper.GetUTXO(ctx, positionD)
	assert.NotNil(t, utxo1)
	utxo2 = mapper.GetUTXO(ctx, positionE)
	assert.NotNil(t, utxo2)

	// Check that outputs are valid
	inputAddresses := [2]crypto.Address{privB.PubKey().Address(), privC.PubKey().Address()}
	assert.Equal(t, uint64(150), utxo1.GetDenom())
	assert.Equal(t, uint64(50), utxo2.GetDenom())
	assert.EqualValues(t, newownerA, utxo1.GetAddress())
	assert.EqualValues(t, newownerB, utxo2.GetAddress())
	assert.EqualValues(t, inputAddresses, utxo1.GetInputAddresses())
	assert.EqualValues(t, inputAddresses, utxo2.GetInputAddresses())
}

/*
	Tests a valid spendmsg 
	1 input and 2 different outputs
*/
func TestOneInput(t *testing.T) {
	ms, capKey := setupMultiStore()

	ctx := sdk.NewContext(ms, abci.Header{Height: 2}, false, nil)
	mapper := NewUTXOMapper(capKey, MakeCodec())
	keeper := NewUTXOKeeper(mapper)
	txIndex := new(uint16)
	handler := NewHandler(keeper, txIndex)

	// Add in 2 parentUTXO
	privA := crypto.GenPrivKeySecp256k1()
	privB := crypto.GenPrivKeySecp256k1()
	positionB := Position{1000, 0, 0}
	utxo1 := NewUTXO(privA, privB, positionB)
	mapper.AddUTXO(ctx, utxo1)
	utxo1 = mapper.GetUTXO(ctx, positionB)
	assert.NotNil(t, utxo1)
	
	newownerA := crypto.Address([]byte("newownerA"))
	newownerB := crypto.Address([]byte("newownerB"))
	confrimSigs := [2]crypto.Signature{crypto.SignatureSecp256k1{}, crypto.SignatureSecp256k1{}}

	// Add in SpendMsg,
	var msg = SpendMsg{
		Blknum1: 		1000,
		Txindex1: 		0,
		Oindex1: 		0,
		Indenom1: 		100,
		Owner1: 		privB.PubKey().Address(),
		ConfirmSigs1: 	confrimSigs,
		Blknum2:		0,
		Txindex2: 		0,
		Oindex2: 		0,
		Indenom2: 		0,
		Owner2: 		crypto.Address([]byte("")),
		ConfirmSigs2: 	confrimSigs,
		Newowner1: 		newownerA,
		Denom1: 		25,
		Newowner2: 		newownerB,
		Denom2: 		75,
		Fee: 			0,
	}

	res := handler(ctx, msg)
	assert.Equal(t, sdk.CodeType(0), sdk.CodeType(res.Code), res.Log)

	assert.Equal(t, uint16(1), *txIndex) // txIndex incremented
	
	//Check that inputs were deleted
	utxo := mapper.GetUTXO(ctx, positionB)
	assert.Nil(t, utxo)

	// Check to see if outputs were added
	assert.Equal(t, int64(2), ctx.BlockHeight())
	positionD := Position{2000, 0, 0}
	positionE := Position{2000, 0, 1}
	utxo1 = mapper.GetUTXO(ctx, positionD)
	assert.NotNil(t, utxo1)
	utxo2 := mapper.GetUTXO(ctx, positionE)
	assert.NotNil(t, utxo2)

	// Check that outputs are valid
	assert.Equal(t, uint64(25), utxo1.GetDenom())
	assert.Equal(t, uint64(75), utxo2.GetDenom())
	assert.EqualValues(t, newownerA, utxo1.GetAddress())
	assert.EqualValues(t, newownerB, utxo2.GetAddress())
}