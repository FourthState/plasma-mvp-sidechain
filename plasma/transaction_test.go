package plasma

import (
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/require"
	"math/big"
	"reflect"
	"testing"
)

func TestTransactionSerialization(t *testing.T) {
	one := big.NewInt(1)
	zero := big.NewInt(0)

	// contstruct a transaction
	tx := &Transaction{}
	pos, _ := FromPositionString("(1.10000.1.0)")
	tx.Input0 = NewInput(pos, [65]byte{}, nil)
	tx.Input0.Signature[1] = byte(1)
	pos, _ = FromPositionString("(0.0.0.1)")
	tx.Input1 = NewInput(pos, [65]byte{}, nil)
	tx.Output0 = NewOutput(common.HexToAddress("1"), one)
	tx.Output1 = NewOutput(common.HexToAddress("0"), zero)
	tx.Fee = big.NewInt(1)

	bytes, err := rlp.EncodeToBytes(tx)
	require.NoError(t, err, "error serializing transaction")
	require.Equal(t, 811, len(bytes), "encoded bytes should sum to 811")

	recoveredTx := &Transaction{}
	err = rlp.DecodeBytes(bytes, recoveredTx)
	require.NoError(t, err, "error deserializing transaction")

	require.EqualValues(t, tx, recoveredTx, "serialized and deserialized transaction not deeply equal")
	require.True(t, reflect.DeepEqual(tx, recoveredTx), "serialized and deserialized transactions not deeply equal")
}

func GetPosition(posStr string) Position {
	pos, _ := FromPositionString(posStr)
	return pos
}

func TestTransactionValidation(t *testing.T) {
	privKey, _ := crypto.GenerateKey()
	addr := crypto.PubkeyToAddress(privKey.PublicKey)

	emptySig := [65]byte{}

	sampleSig := [65]byte{}
	sampleSig[0] = byte(1)
	sampleConfirmSig := [][65]byte{sampleSig}

	type validationCase struct {
		reason string
		Transaction
	}

	invalidTxs := []validationCase{
		validationCase{
			reason: "tx with an empty first input",
			Transaction: Transaction{
				Input0:  NewInput(GetPosition("(0.0.0.0)"), emptySig, nil),
				Input1:  NewInput(GetPosition("(0.0.0.0)"), emptySig, nil),
				Output0: NewOutput(utils.ZeroAddress, utils.Big0),
				Output1: NewOutput(utils.ZeroAddress, utils.Big0),
				Fee:     utils.Big0,
			},
		},
		validationCase{
			reason: "tx with no recipient",
			Transaction: Transaction{
				Input0:  NewInput(GetPosition("(0.0.0.1)"), sampleSig, sampleConfirmSig),
				Input1:  NewInput(GetPosition("(0.0.0.0)"), emptySig, nil),
				Output0: NewOutput(utils.ZeroAddress, utils.Big1),
				Output1: NewOutput(utils.ZeroAddress, utils.Big0),
				Fee:     utils.Big0,
			},
		},
		validationCase{
			reason: "tx with no output amount",
			Transaction: Transaction{
				Input0:  NewInput(GetPosition("(0.0.0.1)"), sampleSig, sampleConfirmSig),
				Input1:  NewInput(GetPosition("(0.0.0.0)"), emptySig, nil),
				Output0: NewOutput(addr, utils.Big0),
				Output1: NewOutput(utils.ZeroAddress, utils.Big0),
				Fee:     utils.Big0,
			},
		},
		validationCase{
			reason: "tx with the same position for both inputs",
			Transaction: Transaction{
				Input0:  NewInput(GetPosition("(0.0.0.1)"), sampleSig, sampleConfirmSig),
				Input1:  NewInput(GetPosition("(0.0.0.1)"), sampleSig, sampleConfirmSig),
				Output0: NewOutput(addr, utils.Big0),
				Output1: NewOutput(utils.ZeroAddress, utils.Big0),
				Fee:     utils.Big0,
			},
		},
	}

	for _, tx := range invalidTxs {
		err := tx.ValidateBasic()
		require.Error(t, err, tx.reason)
	}

	validTxs := []validationCase{
		validationCase{
			reason: "tx with one input and one output",
			Transaction: Transaction{
				Input0:  NewInput(GetPosition("(0.0.0.1)"), sampleSig, nil),
				Input1:  NewInput(GetPosition("(0.0.0.0)"), emptySig, nil),
				Output0: NewOutput(addr, utils.Big1),
				Output1: NewOutput(utils.ZeroAddress, utils.Big0),
				Fee:     utils.Big0,
			},
		},
		validationCase{
			reason: "tx with one input and two output",
			Transaction: Transaction{
				Input0:  NewInput(GetPosition("(0.0.0.1)"), sampleSig, nil),
				Input1:  NewInput(GetPosition("(0.0.0.0)"), emptySig, nil),
				Output0: NewOutput(addr, utils.Big1),
				Output1: NewOutput(addr, utils.Big1),
				Fee:     utils.Big0,
			},
		},
		validationCase{
			reason: "tx with two input and one output",
			Transaction: Transaction{
				Input0:  NewInput(GetPosition("(0.0.0.1)"), sampleSig, nil),
				Input1:  NewInput(GetPosition("(1.0.1.0)"), sampleSig, sampleConfirmSig),
				Output0: NewOutput(addr, utils.Big1),
				Output1: NewOutput(utils.ZeroAddress, utils.Big0),
				Fee:     utils.Big0,
			},
		},
		validationCase{
			reason: "tx with two input and two outputs",
			Transaction: Transaction{
				Input0:  NewInput(GetPosition("(0.0.0.1)"), sampleSig, nil),
				Input1:  NewInput(GetPosition("(1.0.1.0)"), sampleSig, sampleConfirmSig),
				Output0: NewOutput(addr, utils.Big1),
				Output1: NewOutput(addr, utils.Big1),
				Fee:     utils.Big0,
			},
		},
	}

	for _, tx := range validTxs {
		err := tx.ValidateBasic()
		require.NoError(t, err, tx.reason)
	}
}
