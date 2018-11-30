package utxo

import (
	"github.com/stretchr/testify/require"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/FourthState/plasma-mvp-sidechain/utils"
)

/* structs used to test the mapper */
var _ Position = &testPosition{}

type testPosition struct {
	BlockNumber uint64
	TxIndex     uint32
	OutputIndex uint8
}

func testProtoPosition() Position {
	return &testPosition{}
}

func newTestPosition(newPos []uint64) testPosition {
	return testPosition{
		BlockNumber: newPos[0],
		TxIndex:     uint32(newPos[1]),
		OutputIndex: uint8(newPos[2]),
	}
}

func (pos testPosition) Get() []sdk.Uint {
	return []sdk.Uint{sdk.NewUint(pos.BlockNumber), sdk.NewUint(uint64(pos.TxIndex)), sdk.NewUint(uint64(pos.OutputIndex))}
}

func (pos testPosition) IsValid() bool {
	if pos.OutputIndex > 4 {
		return false
	}
	return true
}

/*
	Basic test of Get, Receive, Spend
	Creates a valid UTXO and adds it to the uxto mapping.
	Checks to make sure UTXO isn't nil after adding to mapping.
	Then deletes the UTXO from the mapping
*/

func TestUTXOGetReceiveSpend(t *testing.T) {
	ms, capKey, _ := SetupMultiStore()

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	cdc := MakeCodec()
	cdc.RegisterConcrete(testPosition{}, "x/utxo/testPosition", nil)
	mapper := NewBaseMapper(capKey, cdc)

	priv, _ := ethcrypto.GenerateKey()
	addr := utils.PrivKeyToAddress(priv)

	position := newTestPosition([]uint64{1, 0, 0})

	utxo := NewUTXO(addr.Bytes(), 100, "testEther", position)

	require.Equal(t, addr.Bytes(), utxo.Address)
	require.EqualValues(t, position, utxo.Position)

	mapper.ReceiveUTXO(ctx, utxo)

	received := mapper.GetUTXO(ctx, addr.Bytes(), position)
	require.True(t, received.Valid, "output UTXO is not valid")
	require.Equal(t, utxo, received, "not equal after receive")

	mapper.SpendUTXO(ctx, addr.Bytes(), position, [][]byte{[]byte("spenderKey")})
	utxo = mapper.GetUTXO(ctx, addr.Bytes(), position)
	require.False(t, utxo.Valid, "Spent UTXO is still valid")
	require.Equal(t, utxo.SpenderKeys, [][]byte{[]byte("spenderKey")})
}

/*
	Test multiple additions and deletions in the same block and different blocks
	Creates a valid UTXOs and adds them to the uxto mapping.
	Then deletes the UTXO's from the mapping.
*/

func TestMultiUTXOAddDeleteSameBlock(t *testing.T) {
	ms, capKey, _ := SetupMultiStore()

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	cdc := MakeCodec()
	cdc.RegisterConcrete(testPosition{}, "x/utxo/testPosition", nil)
	mapper := NewBaseMapper(capKey, cdc)

	priv, _ := ethcrypto.GenerateKey()
	addr := utils.PrivKeyToAddress(priv)

	for i := 0; i < 20; i++ {
		position := newTestPosition([]uint64{uint64(i%4) + 1, uint64(i / 4), 0})
		utxo := NewUTXO(addr.Bytes(), 100, "testEther", position)
		mapper.ReceiveUTXO(ctx, utxo)

		utxo = mapper.GetUTXO(ctx, addr.Bytes(), position)
		require.True(t, utxo.Valid, "Received UTXO is not valid")
	}

	for i := 0; i < 20; i++ {
		position := newTestPosition([]uint64{uint64(i%4) + 1, uint64(i / 4), 0})

		utxo := mapper.GetUTXO(ctx, addr.Bytes(), position)
		mapper.SpendUTXO(ctx, addr.Bytes(), position, [][]byte{[]byte("spenderKey")})
		utxo = mapper.GetUTXO(ctx, addr.Bytes(), position)
		require.False(t, utxo.Valid, "Spent UTXO is still valid")
		require.Equal(t, utxo.SpenderKeys, [][]byte{[]byte("spenderKey")})
	}

}

func TestInvalidAddress(t *testing.T) {
	ms, capKey, _ := SetupMultiStore()

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())

	cdc := MakeCodec()
	cdc.RegisterConcrete(testPosition{}, "x/utxo/testPosition", nil)
	mapper := NewBaseMapper(capKey, cdc)

	priv0, _ := ethcrypto.GenerateKey()
	addr0 := utils.PrivKeyToAddress(priv0)

	priv1, _ := ethcrypto.GenerateKey()
	addr1 := utils.PrivKeyToAddress(priv1)

	position := newTestPosition([]uint64{1, 0, 0})

	utxo := NewUTXO(addr0.Bytes(), 100, "testEther", position)

	require.Equal(t, addr0.Bytes(), utxo.Address)
	require.EqualValues(t, position, utxo.Position)

	mapper.ReceiveUTXO(ctx, utxo)

	// GetUTXO with correct position but wrong address
	utxo = mapper.GetUTXO(ctx, addr1.Bytes(), position)
	require.Equal(t, utxo, UTXO{}, "Valid UTXO in wrong location")

	utxo = mapper.GetUTXO(ctx, addr0.Bytes(), position)

	// SpendUTXO with correct position but wrong address
	mapper.SpendUTXO(ctx, addr1.Bytes(), position, [][]byte{[]byte("spenderKey")})
	utxo = mapper.GetUTXO(ctx, addr0.Bytes(), position)
	require.True(t, utxo.Valid, "UTXO invalid after invalid spend")
	require.Nil(t, utxo.SpenderKeys, "UTXO has spenderKeys set after invalid spend")

	mapper.SpendUTXO(ctx, addr0.Bytes(), position, [][]byte{[]byte("spenderKey")})
	utxo = mapper.GetUTXO(ctx, addr0.Bytes(), position)
	require.False(t, utxo.Valid, "UTXO still valid after valid spend")
	require.Equal(t, utxo.SpenderKeys, [][]byte{[]byte("spenderKey")}, "UTXO doesn't have spenderKeys set after valid spend")
}

/*
	Test getting all UTXOs for an Address.
*/

func TestGetUTXOsForAddress(t *testing.T) {
	ms, capKey, _ := SetupMultiStore()

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())

	cdc := MakeCodec()
	cdc.RegisterConcrete(testPosition{}, "x/utxo/testPosition", nil)
	mapper := NewBaseMapper(capKey, cdc)

	privA, _ := ethcrypto.GenerateKey()
	addrA := utils.PrivKeyToAddress(privA)

	privB, _ := ethcrypto.GenerateKey()
	addrB := utils.PrivKeyToAddress(privB)

	privC, _ := ethcrypto.GenerateKey()
	addrC := utils.PrivKeyToAddress(privC)

	positionB0 := newTestPosition([]uint64{1, 0, 0})
	positionB1 := newTestPosition([]uint64{2, 1, 0})
	positionB2 := newTestPosition([]uint64{3, 2, 1})

	utxo0 := NewUTXO(addrB.Bytes(), 100, "testEther", positionB0)
	utxo1 := NewUTXO(addrB.Bytes(), 200, "testEther", positionB1)
	utxo2 := NewUTXO(addrB.Bytes(), 300, "testEther", positionB2)

	mapper.ReceiveUTXO(ctx, utxo0)
	mapper.ReceiveUTXO(ctx, utxo1)
	mapper.ReceiveUTXO(ctx, utxo2)

	utxosForAddressB := mapper.GetUTXOsForAddress(ctx, addrB.Bytes())
	require.Equal(t, 3, len(utxosForAddressB))
	require.Equal(t, utxo0, utxosForAddressB[0])
	require.Equal(t, utxo1, utxosForAddressB[1])
	require.Equal(t, utxo2, utxosForAddressB[2])

	mapper.SpendUTXO(ctx, addrB.Bytes(), positionB1, [][]byte{[]byte("the aether")})

	utxosForAddressB = mapper.GetUTXOsForAddress(ctx, addrB.Bytes())
	require.Equal(t, 2, len(utxosForAddressB))
	require.Equal(t, utxo0, utxosForAddressB[0])
	require.Equal(t, utxo2, utxosForAddressB[1])

	positionC0 := newTestPosition([]uint64{2, 3, 0})
	utxo3 := NewUTXO(addrC.Bytes(), 300, "testEther", positionC0)
	mapper.ReceiveUTXO(ctx, utxo3)
	utxosForAddressC := mapper.GetUTXOsForAddress(ctx, addrC.Bytes())
	require.Equal(t, 1, len(utxosForAddressC))
	require.Equal(t, utxo3, utxosForAddressC[0])

	// check returns empty slice if no UTXOs exist for address
	utxosForAddressA := mapper.GetUTXOsForAddress(ctx, addrA.Bytes())
	require.Empty(t, utxosForAddressA)
}

func TestSpendInvalidUTXO(t *testing.T) {
	ms, capKey, _ := SetupMultiStore()

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	cdc := MakeCodec()
	cdc.RegisterConcrete(testPosition{}, "x/utxo/testPosition", nil)
	mapper := NewBaseMapper(capKey, cdc)

	priv, _ := ethcrypto.GenerateKey()
	addr := utils.PrivKeyToAddress(priv)

	position := newTestPosition([]uint64{1, 0, 0})

	utxo := NewUTXO(addr.Bytes(), 100, "testEther", position)

	require.True(t, utxo.Valid, "UTXO is not valid on creation")

	mapper.InvalidateUTXO(ctx, utxo)

	err := mapper.SpendUTXO(ctx, addr.Bytes(), position, [][]byte{[]byte("spenderKey")})

	utxo = mapper.GetUTXO(ctx, addr.Bytes(), position)
	require.NotNil(t, err, "Allowed invalid UTXO to be spent")
	require.Nil(t, utxo.SpenderKeys, "UTXO mutated after invalid spend")

	mapper.ValidateUTXO(ctx, utxo)
	err = mapper.SpendUTXO(ctx, addr.Bytes(), position, [][]byte{[]byte("spenderKey")})

	utxo = mapper.GetUTXO(ctx, addr.Bytes(), position)
	require.Nil(t, err, "Spend of valid UTXO errorred")
	require.False(t, utxo.Valid, "Spent UTXO is still valid")
	require.Equal(t, utxo.SpenderKeys, [][]byte{[]byte("spenderKey")}, "UTXO doesn't have spenderKeys set after valid spend")
}

func TestUTXOMethods(t *testing.T) {
	_, capKey, _ := SetupMultiStore()

	cdc := MakeCodec()
	cdc.RegisterConcrete(testPosition{}, "x/utxo/testPosition", nil)
	mapper := NewBaseMapper(capKey, cdc)

	addr1 := []byte("12345")
	addr2 := []byte("13579")
	addr3 := []byte("67890")

	outputPos1 := testPosition{7, 8, 9}
	outputPos2 := testPosition{3, 4, 5}

	outputKey1 := mapper.ConstructKey(addr1, outputPos1)
	outputKey2 := mapper.ConstructKey(addr2, outputPos2)

	testUTXO := NewUTXO(addr3, 100, "Ether", testPosition{0, 1, 2})
	testUTXO.SpenderKeys = [][]byte{outputKey1, outputKey2}

	addrs := testUTXO.SpenderAddresses()
	require.Equal(t, [][]byte{addr1, addr2}, addrs, "Spender addresses are not correct")

	positions := testUTXO.SpenderPositions(cdc, testProtoPosition)
	require.Equal(t, []Position{&outputPos1, &outputPos2}, positions, "Spender position not correct")
}
