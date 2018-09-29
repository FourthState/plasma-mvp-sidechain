package auth

import (
	"crypto/ecdsa"
	"fmt"
	types "github.com/FourthState/plasma-mvp-sidechain/types"
	utils "github.com/FourthState/plasma-mvp-sidechain/utils"
	"github.com/FourthState/plasma-mvp-sidechain/x/utxo"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	"testing"
)

func setup() (sdk.Context, utxo.Mapper, utxo.FeeUpdater) {
	ms, capKey := utxo.SetupMultiStore()

	ctx := sdk.NewContext(ms, abci.Header{}, false, log.NewNopLogger())
	cdc := utxo.MakeCodec()
	types.RegisterAmino(cdc)
	mapper := utxo.NewBaseMapper(capKey, cdc)

	return ctx, mapper, feeUpdater
}

// should be modified when fees are implemented
func feeUpdater(outputs []utxo.Output) sdk.Error {
	return nil
}

func GenSpendMsg() types.SpendMsg {
	// Creates Basic Spend Msg with owners and recipients
	confirmSigs := [2]types.Signature{types.Signature{}, types.Signature{}}
	privKeyA, _ := ethcrypto.GenerateKey()
	privKeyB, _ := ethcrypto.GenerateKey()

	return types.SpendMsg{
		Blknum0:      1,
		Txindex0:     0,
		Oindex0:      0,
		DepositNum0:  0,
		Owner0:       utils.PrivKeyToAddress(privKeyA),
		ConfirmSigs0: confirmSigs,
		Blknum1:      1,
		Txindex1:     1,
		Oindex1:      0,
		DepositNum1:  0,
		Owner1:       utils.PrivKeyToAddress(privKeyA),
		ConfirmSigs1: confirmSigs,
		Newowner0:    utils.PrivKeyToAddress(privKeyB),
		Amount0:      150,
		Newowner1:    utils.PrivKeyToAddress(privKeyB),
		Amount1:      50,
		FeeAmount:    0,
	}
}

// Returns a confirmsig array signed by privKey0 and privKey1
func CreateConfirmSig(position types.PlasmaPosition, privKey0, privKey1 *ecdsa.PrivateKey, two_inputs bool) (confirmSigs [2]types.Signature) {
	confirmBytes := position.GetSignBytes()
	hash := ethcrypto.Keccak256(confirmBytes)
	confirmSig, _ := ethcrypto.Sign(hash, privKey0)

	var confirmSig1 []byte
	if two_inputs {
		confirmSig1, _ = ethcrypto.Sign(hash, privKey1)
	}
	confirmSigs = [2]types.Signature{types.Signature{confirmSig}, types.Signature{confirmSig1}}
	return confirmSigs
}

// helper for constructing single or double input tx
func GetTx(msg types.SpendMsg, privKey0, privKey1 *ecdsa.PrivateKey, two_sigs bool) (tx types.BaseTx) {
	hash := ethcrypto.Keccak256(msg.GetSignBytes())
	sig0, _ := ethcrypto.Sign(hash, privKey0)

	tx = types.NewBaseTx(msg, []types.Signature{{
		Sig: sig0,
	}})

	if two_sigs {
		sig1, _ := ethcrypto.Sign(hash, privKey1)
		tx.Signatures = append(tx.Signatures, types.Signature{sig1})
	}

	return tx
}

// helper for constructing input addresses
func getInputAddr(addr0, addr1 common.Address, two bool) [2]common.Address {
	if two {
		return [2]common.Address{addr0, addr1}
	} else {
		return [2]common.Address{addr0, common.Address{}}
	}
}

// No signatures are provided
func TestNoSigs(t *testing.T) {
	ctx, mapper, feeUpdater := setup()

	var msg = GenSpendMsg()
	tx := types.NewBaseTx(msg, []types.Signature{})

	handler := NewAnteHandler(mapper, feeUpdater)
	_, res, abort := handler(ctx, tx)

	assert.Equal(t, true, abort, "did not abort with no signatures")
	require.Equal(t, sdk.ToABCICode(sdk.CodespaceType(1), sdk.CodeType(4)), res.Code, "tx had processed with no signatures")
}

// The wrong amount of signatures are provided
func TestNotEnoughSigs(t *testing.T) {
	ctx, mapper, feeUpdater := setup()

	var msg = GenSpendMsg()
	priv, _ := ethcrypto.GenerateKey()
	hash := ethcrypto.Keccak256(msg.GetSignBytes())
	sig, _ := ethcrypto.Sign(hash, priv)
	tx := types.NewBaseTx(msg, []types.Signature{types.Signature{sig}})

	handler := NewAnteHandler(mapper, feeUpdater)
	_, res, abort := handler(ctx, tx)

	assert.Equal(t, true, abort, "did not abort with incorrect number of signatures")
	require.Equal(t, sdk.ToABCICode(sdk.CodespaceType(1), sdk.CodeType(4)), res.Code, "tx had processed with incorrect number of signatures")
}

// helper struct for readability
type input struct {
	owner_index  int64
	addr         common.Address
	position     types.PlasmaPosition
	input_index0 int64
	input_index1 int64
}

// Tests a different cases
func TestDifferentCases(t *testing.T) {
	ctx, mapper, feeUpdater := setup()

	var keys [6]*ecdsa.PrivateKey
	var addrs []common.Address

	for i := 0; i < 6; i++ {
		keys[i], _ = ethcrypto.GenerateKey()
		addrs = append(addrs, utils.PrivKeyToAddress(keys[i]))
	}

	cases := []struct {
		input0    input
		input1    input
		newowner0 common.Address
		amount0   uint64
		newowner1 common.Address
		amount1   uint64
		abort     bool
	}{
		// Test Case 0: Tx signed by the wrong address
		{
			input{1, addrs[0], types.NewPlasmaPosition(2, 0, 0, 0), 1, -1}, // first input
			input{-1, common.Address{}, types.PlasmaPosition{}, -1, -1},    // second input
			addrs[1], 1000, // first output
			addrs[2], 1000, // second output
			true,
		},

		// Test Case 1: Inputs != Outputs + Fee
		{
			input{0, addrs[0], types.NewPlasmaPosition(3, 0, 0, 0), 1, -1},
			input{-1, common.Address{}, types.PlasmaPosition{}, -1, -1},
			addrs[1], 2000,
			addrs[2], 1000,
			true,
		},

		// Test Case 2: 1 input 2 output
		{
			input{0, addrs[0], types.NewPlasmaPosition(4, 0, 0, 0), 1, -1},
			input{-1, common.Address{}, types.PlasmaPosition{}, -1, -1},
			addrs[1], 1000,
			addrs[2], 1000,
			false,
		},

		// Test Case 3: 2 input 2 output
		{
			input{1, addrs[1], types.NewPlasmaPosition(5, 0, 0, 0), 0, -1},
			input{2, addrs[2], types.NewPlasmaPosition(5, 0, 1, 0), 0, -1},
			addrs[3], 2500,
			addrs[4], 1500,
			false,
		},
	}

	for index, tc := range cases {
		input0_index1 := utils.GetIndex(tc.input0.input_index1)
		input1_index0 := utils.GetIndex(tc.input1.input_index0)
		input1_index1 := utils.GetIndex(tc.input1.input_index1)
		var msg = types.SpendMsg{
			Blknum0:      tc.input0.position.Blknum,
			Txindex0:     tc.input0.position.TxIndex,
			Oindex0:      tc.input0.position.Oindex,
			DepositNum0:  tc.input0.position.DepositNum,
			Owner0:       tc.input0.addr,
			ConfirmSigs0: CreateConfirmSig(tc.input0.position, keys[tc.input0.input_index0], keys[input0_index1], tc.input0.input_index1 != -1),
			Blknum1:      tc.input1.position.Blknum,
			Txindex1:     tc.input1.position.TxIndex,
			Oindex1:      tc.input1.position.Oindex,
			DepositNum1:  tc.input1.position.DepositNum,
			Owner1:       tc.input1.addr,
			ConfirmSigs1: CreateConfirmSig(tc.input1.position, keys[input1_index0], keys[input1_index1], tc.input1.input_index1 != -1),
			Newowner0:    tc.newowner0,
			Amount0:      tc.amount0,
			Newowner1:    tc.newowner1,
			Amount1:      tc.amount1 - 5,
			FeeAmount:    5,
		}

		owner_index1 := utils.GetIndex(tc.input1.owner_index)
		tx := GetTx(msg, keys[tc.input0.owner_index], keys[owner_index1], tc.input1.owner_index != -1)

		handler := NewAnteHandler(mapper, feeUpdater)
		_, res, abort := handler(ctx, tx)

		assert.Equal(t, true, abort, fmt.Sprintf("did not abort on utxo that does not exist. Case: %d", index))
		require.Equal(t, sdk.ToABCICode(sdk.CodespaceType(1), sdk.CodeType(6)), res.Code, res.Log)

		inputAddr := getInputAddr(addrs[tc.input0.input_index0], addrs[input0_index1], tc.input0.input_index1 != -1)
		utxo0 := types.NewBaseUTXO(tc.input0.addr, inputAddr, 2000, "Ether", tc.input0.position)

		var utxo1 utxo.UTXO
		if tc.input1.owner_index != -1 {
			inputAddr = getInputAddr(addrs[input1_index0], addrs[input1_index1], tc.input0.input_index1 != -1)
			utxo1 = types.NewBaseUTXO(tc.input1.addr, inputAddr, 2000, "Ether", tc.input1.position)
		}

		mapper.AddUTXO(ctx, utxo0)
		if tc.input1.owner_index != -1 {
			mapper.AddUTXO(ctx, utxo1)
		}
		_, res, abort = handler(ctx, tx)

		assert.Equal(t, tc.abort, abort, fmt.Sprintf("aborted on case: %d", index))
		if tc.abort == false {
			require.Equal(t, sdk.ToABCICode(sdk.CodespaceType(1), sdk.CodeType(0)), res.Code, res.Log)
		} else {
			require.NotEqual(t, sdk.ToABCICode(sdk.CodespaceType(1), sdk.CodeType(0)), res.Code, res.Log)
		}
	}
}
