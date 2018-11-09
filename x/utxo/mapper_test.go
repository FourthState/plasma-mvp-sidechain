package utxo

import (
	"errors"
	"fmt"
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

func testProtoUTXO(ctx sdk.Context, msg sdk.Msg) UTXO {
	return &testUTXO{}
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

var _ UTXO = &testUTXO{}

type testUTXO struct {
	Owner    []byte
	Amount   uint64
	Denom    string
	Position testPosition
}

func newTestUTXO(owner []byte, amount uint64, position testPosition) UTXO {
	return &testUTXO{
		Owner:    owner,
		Amount:   amount,
		Denom:    "Ether",
		Position: position,
	}
}

func (utxo testUTXO) GetAddress() []byte {
	return utxo.Owner
}

func (utxo *testUTXO) SetAddress(address []byte) error {
	if utxo.Owner != nil {
		return errors.New("Owner already set")
	}
	utxo.Owner = address
	return nil
}

func (utxo testUTXO) GetAmount() uint64 {
	return utxo.Amount
}

func (utxo *testUTXO) SetAmount(amount uint64) error {
	if utxo.Amount != 0 {
		return errors.New("Owner already set")
	}
	utxo.Amount = amount
	return nil
}

func (utxo testUTXO) GetDenom() string {
	return utxo.Denom
}

func (utxo *testUTXO) SetDenom(denom string) error {
	if utxo.Denom != "" {
		return errors.New("Owner already set")
	}
	utxo.Denom = denom
	return nil
}

func (utxo testUTXO) GetPosition() Position {
	return utxo.Position
}

func (utxo *testUTXO) SetPosition(pos Position) error {

	position, ok := pos.(testPosition)
	if !ok {
		fmt.Println("ah")
		return errors.New("position setting err")
	}
	utxo.Position = position
	return nil
}

/*
	Basic test of Get, Add, Delete
	Creates a valid UTXO and adds it to the uxto mapping.
	Checks to make sure UTXO isn't nil after adding to mapping.
	Then deletes the UTXO from the mapping
*/

func TestUTXOGetAddDelete(t *testing.T) {
	ms, capKey, _ := SetupMultiStore()

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	cdc := MakeCodec()
	cdc.RegisterConcrete(&testUTXO{}, "x/utxo/testUTXO", nil)
	cdc.RegisterConcrete(&testPosition{}, "x/utxo/testPosition", nil)
	mapper := NewBaseMapper(capKey, cdc)

	priv, _ := ethcrypto.GenerateKey()
	addr := utils.PrivKeyToAddress(priv)

	position := newTestPosition([]uint64{1, 0, 0})

	utxo := newTestUTXO(addr.Bytes(), 100, position)

	require.NotNil(t, utxo)
	require.Equal(t, addr.Bytes(), utxo.GetAddress())
	require.EqualValues(t, position, utxo.GetPosition())

	mapper.AddUTXO(ctx, utxo)

	utxo = mapper.GetUTXO(ctx, addr.Bytes(), position)
	require.NotNil(t, utxo)

	mapper.DeleteUTXO(ctx, addr.Bytes(), position)
	utxo = mapper.GetUTXO(ctx, addr.Bytes(), position)
	require.Nil(t, utxo)
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
	cdc.RegisterConcrete(&testUTXO{}, "x/utxo/testUTXO", nil)
	cdc.RegisterConcrete(&testPosition{}, "x/utxo/testPosition", nil)
	mapper := NewBaseMapper(capKey, cdc)

	priv, _ := ethcrypto.GenerateKey()
	addr := utils.PrivKeyToAddress(priv)

	for i := 0; i < 20; i++ {
		position := newTestPosition([]uint64{uint64(i%4) + 1, uint64(i / 4), 0})
		utxo := newTestUTXO(addr.Bytes(), 100, position)
		mapper.AddUTXO(ctx, utxo)

		utxo = mapper.GetUTXO(ctx, addr.Bytes(), position)
		require.NotNil(t, utxo)
	}

	for i := 0; i < 20; i++ {
		position := newTestPosition([]uint64{uint64(i%4) + 1, uint64(i / 4), 0})

		utxo := mapper.GetUTXO(ctx, addr.Bytes(), position)
		require.NotNil(t, utxo)
		mapper.DeleteUTXO(ctx, addr.Bytes(), position)
		utxo = mapper.GetUTXO(ctx, addr.Bytes(), position)
		require.Nil(t, utxo)
	}

}

func TestInvalidAddress(t *testing.T) {
	ms, capKey, _ := SetupMultiStore()

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())

	cdc := MakeCodec()
	cdc.RegisterConcrete(&testUTXO{}, "x/utxo/testUTXO", nil)
	cdc.RegisterConcrete(&testPosition{}, "x/utxo/testPosition", nil)
	mapper := NewBaseMapper(capKey, cdc)

	priv0, _ := ethcrypto.GenerateKey()
	addr0 := utils.PrivKeyToAddress(priv0)

	priv1, _ := ethcrypto.GenerateKey()
	addr1 := utils.PrivKeyToAddress(priv1)

	position := newTestPosition([]uint64{1, 0, 0})

	utxo := newTestUTXO(addr0.Bytes(), 100, position)

	require.NotNil(t, utxo)
	require.Equal(t, addr0.Bytes(), utxo.GetAddress())
	require.EqualValues(t, position, utxo.GetPosition())

	mapper.AddUTXO(ctx, utxo)

	// GetUTXO with correct position but wrong address
	utxo = mapper.GetUTXO(ctx, addr1.Bytes(), position)
	require.Nil(t, utxo)

	utxo = mapper.GetUTXO(ctx, addr0.Bytes(), position)
	require.NotNil(t, utxo)

	// DeleteUTXO with correct position but wrong address
	mapper.DeleteUTXO(ctx, addr1.Bytes(), position)
	utxo = mapper.GetUTXO(ctx, addr0.Bytes(), position)
	require.NotNil(t, utxo)

	mapper.DeleteUTXO(ctx, addr0.Bytes(), position)
	utxo = mapper.GetUTXO(ctx, addr0.Bytes(), position)
	require.Nil(t, utxo)
}

/*
	Test getting all UTXOs for an Address.
*/

func TestGetUTXOsForAddress(t *testing.T) {
	ms, capKey, _ := SetupMultiStore()

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())

	cdc := MakeCodec()
	cdc.RegisterConcrete(&testUTXO{}, "x/utxo/testUTXO", nil)
	cdc.RegisterConcrete(&testPosition{}, "x/utxo/testPosition", nil)
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

	utxo0 := newTestUTXO(addrB.Bytes(), 100, positionB0)
	utxo1 := newTestUTXO(addrB.Bytes(), 200, positionB1)
	utxo2 := newTestUTXO(addrB.Bytes(), 300, positionB2)

	mapper.AddUTXO(ctx, utxo0)
	mapper.AddUTXO(ctx, utxo1)
	mapper.AddUTXO(ctx, utxo2)

	utxosForAddressB := mapper.GetUTXOsForAddress(ctx, addrB.Bytes())
	require.NotNil(t, utxosForAddressB)
	require.Equal(t, 3, len(utxosForAddressB))
	require.Equal(t, utxo0, utxosForAddressB[0])
	require.Equal(t, utxo1, utxosForAddressB[1])
	require.Equal(t, utxo2, utxosForAddressB[2])

	positionC0 := newTestPosition([]uint64{2, 3, 0})
	utxo3 := newTestUTXO(addrC.Bytes(), 300, positionC0)
	mapper.AddUTXO(ctx, utxo3)
	utxosForAddressC := mapper.GetUTXOsForAddress(ctx, addrC.Bytes())
	require.NotNil(t, utxosForAddressC)
	require.Equal(t, 1, len(utxosForAddressC))
	require.Equal(t, utxo3, utxosForAddressC[0])

	// check returns empty slice if no UTXOs exist for address
	utxosForAddressA := mapper.GetUTXOsForAddress(ctx, addrA.Bytes())
	require.Empty(t, utxosForAddressA)
}
