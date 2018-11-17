package utxo

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	rlp "github.com/ethereum/go-ethereum/rlp"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"

	utils "github.com/FourthState/plasma-mvp-sidechain/utils"
)

/*  The following structs are used to do testing on the handler.
They are not fully implemented and should not be used
besides for testing
*/

type testApp struct {
	txindex uint64
	oindex  uint64
}

func (a *testApp) testNextPosition(ctx sdk.Context, secondary bool) Position {
	if !secondary {
		a.txindex++
	}
	a.oindex++
	return newTestPosition([]uint64{uint64(ctx.BlockHeight()), uint64(a.txindex - 1), uint64(a.oindex - 1)})
}

var _ SpendMsg = testSpendMsg{}

type testSpendMsg struct {
	Input  []Input
	Output []Output
	Fees   Output
}

func (msg testSpendMsg) Type() string { return "spend_utxo" }

func (msg testSpendMsg) Route() string { return "spend" }

func (msg testSpendMsg) ValidateBasic() sdk.Error {
	return nil
}

func (msg testSpendMsg) GetSignBytes() []byte {
	b, err := rlp.EncodeToBytes(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg testSpendMsg) GetSigners() []sdk.AccAddress {
	addrs := make([]sdk.AccAddress, 1)
	addrs[0] = sdk.AccAddress(msg.Input[0].Owner)
	return addrs
}

func (msg testSpendMsg) Inputs() []Input {
	return msg.Input
}

func (msg testSpendMsg) Outputs() []Output {
	return msg.Output
}

func (msg testSpendMsg) Fee() []Output {
	return []Output{msg.Fees}
}

func TestHandleSpendMessage(t *testing.T) {
	const len int = 10 // number of addresses avaliable
	var keys [len]*ecdsa.PrivateKey
	var addrs []common.Address
	for i := 0; i < len; i++ {
		keys[i], _ = ethcrypto.GenerateKey()
		addrs = append(addrs, utils.PrivKeyToAddress(keys[i]))
	}

	cases := []struct {
		inputNum  int
		outputNum int
	}{
		// Test Case 0: 1 input 1 ouput
		{1, 1},
		// Test Case 1: 1 input multiple outputs
		{1, 10},
		// Test Case 2: multiple inputs 1 output
		{10, 1},
		// Test Case 3: multiple inputs multiple outputs
		{10, 10},
	}

	for index, tc := range cases {
		ms, capKey, _ := SetupMultiStore()

		cdc := MakeCodec()
		cdc.RegisterConcrete(&testUTXO{}, "x/utxo/testUTXO", nil)
		cdc.RegisterConcrete(&testPosition{}, "x/utxo/testPosition", nil)
		mapper := NewBaseMapper(capKey, cdc)
		app := testApp{0, 0}
		handler := NewSpendHandler(mapper, app.testNextPosition, testProtoUTXO)

		ctx := sdk.NewContext(ms, abci.Header{Height: 6}, false, log.NewNopLogger())
		var inputs []Input
		var outputs []Output

		// Add utxo's that will be spent
		for i := 0; i < tc.inputNum; i++ {
			position := newTestPosition([]uint64{5, uint64(i), 0})
			utxo := newTestUTXO(addrs[i].Bytes(), 100, position)
			mapper.AddUTXO(ctx, utxo)

			utxo = mapper.GetUTXO(ctx, addrs[i].Bytes(), position)
			require.NotNil(t, utxo)

			inputs = append(inputs, Input{addrs[i].Bytes(), position})
		}

		for i := 0; i < tc.outputNum; i++ {
			outputs = append(outputs, Output{addrs[(i+1)%len].Bytes(), "Ether", uint64((100 * tc.inputNum) / tc.outputNum)})
		}

		// Create spend msg
		msg := testSpendMsg{
			Input:  inputs,
			Output: outputs,
			Fees:   Output{[]byte{}, "Ether", 100},
		}

		res := handler(ctx, msg)
		require.Equal(t, sdk.CodeType(0), sdk.CodeType(res.Code), res.Log)

		// Delete inputs
		for _, in := range msg.Inputs() {
			mapper.DeleteUTXO(ctx, in.Owner, in.Position)
			utxo := mapper.GetUTXO(ctx, in.Owner, in.Position)
			require.Nil(t, utxo)
		}

		// Check that outputs were created and are valid
		// Then delete the outputs
		for i, o := range msg.Outputs() {
			position := newTestPosition([]uint64{6, 0, uint64(i)})
			utxo := mapper.GetUTXO(ctx, o.Owner, position)
			require.NotNil(t, utxo, fmt.Sprintf("test case %d, output %d", index, i))

			require.Equal(t, uint64((tc.inputNum*100)/tc.outputNum), utxo.GetAmount())
			require.EqualValues(t, addrs[(i+1)%len].Bytes(), utxo.GetAddress())

			mapper.DeleteUTXO(ctx, o.Owner, position)
			utxo = mapper.GetUTXO(ctx, o.Owner, position)
			require.Nil(t, utxo)
		}
	}
}
