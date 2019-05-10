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

	// contstruct a transaction
	tx := &Transaction{}
	pos, _ := FromPositionString("(1.10000.1.0)")
	confirmSig0 := make([][65]byte, 1)
	copy(confirmSig0[0][65-len([]byte("confirm sig")):], []byte("confirm sig"))
	tx.Inputs = append(tx.Inputs, NewInput(pos, [65]byte{}, confirmSig0))
	tx.Inputs[0].Signature[1] = byte(1)
	pos, _ = FromPositionString("(0.0.0.1)")
	confirmSig1 := make([][65]byte, 2)
	copy(confirmSig1[0][65-len([]byte("the second confirm sig")):], []byte("the second confirm sig"))
	copy(confirmSig1[1][65-len([]byte("a very long string turned into bytes")):], []byte("a very long string turned into bytes"))
	tx.Inputs = append(tx.Inputs, NewInput(pos, [65]byte{}, confirmSig1))
	tx.Outputs = append(tx.Outputs, NewOutput(common.HexToAddress("1"), one))
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
				Inputs:  []Input{NewInput(GetPosition("(0.0.0.0)"), emptySig, nil)},
				Outputs: []Output{NewOutput(utils.ZeroAddress, utils.Big0)},
				Fee:     utils.Big0,
			},
		},
		validationCase{
			reason: "tx with no recipient",
			Transaction: Transaction{
				Inputs:  []Input{NewInput(GetPosition("(0.0.0.1)"), sampleSig, sampleConfirmSig)},
				Outputs: []Output{NewOutput(utils.ZeroAddress, utils.Big1)},
				Fee:     utils.Big0,
			},
		},
		validationCase{
			reason: "tx with no output amount",
			Transaction: Transaction{
				Inputs:  []Input{NewInput(GetPosition("(0.0.0.1)"), sampleSig, sampleConfirmSig)},
				Outputs: []Output{NewOutput(addr, utils.Big0)},
				Fee:     utils.Big0,
			},
		},
		validationCase{
			reason: "tx with the same position for both inputs",
			Transaction: Transaction{
				Inputs:  []Input{NewInput(GetPosition("(0.0.0.1)"), sampleSig, sampleConfirmSig), NewInput(GetPosition("(0.0.0.1)"), sampleSig, sampleConfirmSig)},
				Outputs: []Output{NewOutput(addr, utils.Big0)},
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
				Inputs:  []Input{NewInput(GetPosition("(0.0.0.1)"), sampleSig, nil)},
				Outputs: []Output{NewOutput(addr, utils.Big1)},
				Fee:     utils.Big0,
			},
		},
		validationCase{
			reason: "tx with one input and two output",
			Transaction: Transaction{
				Inputs:  []Input{NewInput(GetPosition("(0.0.0.1)"), sampleSig, nil)},
				Outputs: []Output{NewOutput(addr, utils.Big1), NewOutput(addr, utils.Big1)},
				Fee:     utils.Big0,
			},
		},
		validationCase{
			reason: "tx with two input and one output",
			Transaction: Transaction{
				Inputs:  []Input{NewInput(GetPosition("(0.0.0.1)"), sampleSig, nil), NewInput(GetPosition("(1.0.1.0)"), sampleSig, sampleConfirmSig)},
				Outputs: []Output{NewOutput(addr, utils.Big1)},
				Fee:     utils.Big0,
			},
		},
		validationCase{
			reason: "tx with two input and two outputs",
			Transaction: Transaction{
				Inputs:  []Input{NewInput(GetPosition("(0.0.0.1)"), sampleSig, nil), NewInput(GetPosition("(1.0.1.0)"), sampleSig, sampleConfirmSig)},
				Outputs: []Output{NewOutput(addr, utils.Big1), NewOutput(addr, utils.Big1)},
				Fee:     utils.Big0,
			},
		},
	}

	for _, tx := range validTxs {
		err := tx.ValidateBasic()
		require.NoError(t, err, tx.reason)
	}
}
