package store

import (
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/require"
	"math/big"
	"reflect"
	"testing"
)

// Test that an account can be serialized and deserialized
func TestAccountSerialization(t *testing.T) {
	// Construct Account
	acc := Account{
		Balance: big.NewInt(234578),
		Unspent: []plasma.Position{getPosition("(8745.1239.1.0)"), getPosition("(23409.12456.0.0)"), getPosition("(894301.1.1.0)"), getPosition("(0.0.0.540124)")},
		Spent:   []plasma.Position{getPosition("0.0.0.3"), getPosition("7.734.1.3")},
	}

	bytes, err := rlp.EncodeToBytes(&acc)
	require.NoError(t, err)

	recoveredAcc := Account{}
	err = rlp.DecodeBytes(bytes, &recoveredAcc)
	require.NoError(t, err)

	require.True(t, reflect.DeepEqual(acc, recoveredAcc), "mismatch in serialized and deserialized account")
}

// Test that the Deposit can be serialized and deserialized without loss of information
func TestDepositSerialization(t *testing.T) {
	// Construct deposit
	plasmaDeposit := plasma.Deposit{
		Owner:       common.BytesToAddress([]byte("an ethereum address")),
		Amount:      big.NewInt(12312310),
		EthBlockNum: big.NewInt(100123123),
	}

	deposit := Deposit{
		Deposit:   plasmaDeposit,
		Spent:     true,
		SpenderTx: []byte{},
	}

	bytes, err := rlp.EncodeToBytes(&deposit)
	require.NoError(t, err)

	recoveredDeposit := Deposit{}
	err = rlp.DecodeBytes(bytes, &recoveredDeposit)
	require.NoError(t, err)

	require.True(t, reflect.DeepEqual(deposit, recoveredDeposit), "mismatch in serialized and deserialized deposit")
}

// Test that Transaction can be serialized and deserialized without loss of information
func TestTxSerialization(t *testing.T) {
	hashes := []byte("fourthstate")
	var sigs [65]byte

	// Construct Transaction
	transaction := plasma.Transaction{
		Inputs:  []plasma.Input{plasma.NewInput(getPosition("(1.15.1.0)"), sigs, [][65]byte{}), plasma.NewInput(getPosition("(0.0.0.1)"), sigs, [][65]byte{})},
		Outputs: []plasma.Output{plasma.NewOutput(common.HexToAddress("1"), utils.Big1), plasma.NewOutput(common.HexToAddress("2"), utils.Big2)},
		Fee:     utils.Big1,
	}

	tx := Transaction{
		Transaction:      transaction,
		Spent:            []bool{false, false},
		SpenderTxs:       [][]byte{},
		ConfirmationHash: hashes,
	}

	bytes, err := rlp.EncodeToBytes(&tx)
	require.NoError(t, err)

	recoveredTx := Transaction{}
	err = rlp.DecodeBytes(bytes, &recoveredTx)
	require.NoError(t, err)

	require.True(t, reflect.DeepEqual(tx, recoveredTx), "mismatch in serialized and deserialized transactions")
}
