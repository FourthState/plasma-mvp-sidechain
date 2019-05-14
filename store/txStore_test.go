package store

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func GetPosition(posStr string) plasma.Position {
	pos, _ := plasma.FromPositionString(posStr)
	return pos
}

func TestTxSerialization(t *testing.T) {
	hashes := []byte("allam")
	var sigs [65]byte

	// no default attributes
	transaction := plasma.Transaction{
		Inputs:  []plasma.Input{plasma.NewInput(plasma.NewPosition(utils.Big1, 15, 1, utils.Big0), sigs, [][65]byte{}), plasma.NewInput(plasma.NewPosition(utils.Big0, 0, 0, utils.Big1), sigs, [][65]byte{})},
		Outputs: []plasma.Output{plasma.NewOutput(common.HexToAddress("1"), utils.Big1), plasma.NewOutput(common.HexToAddress("2"), utils.Big2)},
		Fee:     utils.Big1,
	}

	tx := Transaction{
		Transaction:      transaction,
		Spent:            []bool{false, false},
		Spenders:         [][]byte{},
		ConfirmationHash: hashes,
	}

	bytes, err := rlp.EncodeToBytes(&tx)
	require.NoError(t, err)

	recoveredTx := Transaction{}
	err = rlp.DecodeBytes(bytes, &recoveredTx)
	require.NoError(t, err)

	require.True(t, reflect.DeepEqual(tx, recoveredTx), "mismatch in serialized and deserialized transactions")
}

// Test Get, Has, Store, Spend functions
func TestTransactions(t *testing.T) {
	// Setup
	ctx, key := setup()
	txStore := NewTxStore(key)

	privKey, _ := crypto.GenerateKey()
	addr := crypto.PubkeyToAddress(privKey.PublicKey)

	sampleSig := [65]byte{}
	sampleSig[0] = byte(1)
	sampleConfirmSig := [][65]byte{sampleSig}
	confirmationHash := []byte("confirmation hash")

	type validationCase struct {
		reason string
		plasma.Transaction
		plasma.Position
	}

	txs := []validationCase{
		validationCase{
			reason: "tx with two inputs and one output",
			Transaction: plasma.Transaction{
				Inputs:  []plasma.Input{plasma.NewInput(GetPosition("(122.3.1.0)"), sampleSig, nil), plasma.NewInput(GetPosition("(17622.13.5.0)"), sampleSig, nil)},
				Outputs: []plasma.Output{plasma.NewOutput(addr, utils.Big1)},
				Fee:     utils.Big0,
			},
			Position: GetPosition("(8765.6847.0.0)"),
		},
		validationCase{
			reason: "tx with one input and two output",
			Transaction: plasma.Transaction{
				Inputs:  []plasma.Input{plasma.NewInput(GetPosition("(4.1234.1.0)"), sampleSig, nil)},
				Outputs: []plasma.Output{plasma.NewOutput(addr, utils.Big1), plasma.NewOutput(addr, utils.Big1)},
				Fee:     utils.Big0,
			},
			Position: GetPosition("(4354.8765.0.0)"),
		},
		validationCase{
			reason: "tx with two input and one output",
			Transaction: plasma.Transaction{
				Inputs:  []plasma.Input{plasma.NewInput(GetPosition("(123.1.0.0)"), sampleSig, nil), plasma.NewInput(GetPosition("(10.12.1.0)"), sampleSig, sampleConfirmSig)},
				Outputs: []plasma.Output{plasma.NewOutput(addr, utils.Big1)},
				Fee:     utils.Big0,
			},
			Position: GetPosition("(11.123.0.0)"),
		},
		validationCase{
			reason: "tx with two input and two outputs",
			Transaction: plasma.Transaction{
				Inputs:  []plasma.Input{plasma.NewInput(GetPosition("(132.231.1.0)"), sampleSig, nil), plasma.NewInput(GetPosition("(635.927.1.0)"), sampleSig, sampleConfirmSig)},
				Outputs: []plasma.Output{plasma.NewOutput(addr, utils.Big1), plasma.NewOutput(addr, utils.Big1)},
				Fee:     utils.Big0,
			},
			Position: GetPosition("(1121234.12.0.0)"),
		},
	}

	for i, plasmaTx := range txs {
		pos := plasmaTx.Position
		// Retrieve/Spend nonexistent txs
		exists := txStore.HasTx(ctx, plasmaTx.Transaction.TxHash())
		require.False(t, exists, "returned true for nonexistent transaction")
		recoveredTx, ok := txStore.GetTx(ctx, plasmaTx.Transaction.TxHash())
		require.Empty(t, recoveredTx, "did not return empty struct for nonexistent transaction")
		require.False(t, ok, "did not return error on nonexistent transaction")

		// Create and store new transaction
		tx := Transaction{plasmaTx.Transaction, confirmationHash, make([]bool, len(plasmaTx.Transaction.Outputs)), make([][]byte, len(plasmaTx.Transaction.Outputs)), GetPosition("(4567.1.1.0)")}
		for i, _ := range tx.Spenders {
			tx.Spenders[i] = []byte{}
		}
		txStore.StoreTx(ctx, tx)

		// Create Outputs
		for j, _ := range plasmaTx.Transaction.Outputs {
			p := plasma.NewPosition(pos.BlockNum, pos.TxIndex, uint8(j), pos.DepositNonce)
			exists = txStore.HasUTXO(ctx, p)
			require.False(t, exists, "returned true for nonexistent output")
			res := txStore.SpendUTXO(ctx, p, plasmaTx.Transaction.MerkleHash())
			require.Equal(t, res.Code, CodeOutputDNE, "did not return that utxo does not exist")

			txStore.StoreUTXO(ctx, p, plasmaTx.Transaction.TxHash())
			exists = txStore.HasUTXO(ctx, p)
			require.True(t, exists, "returned false for stored utxo")
		}

		// Check for Tx
		exists = txStore.HasTx(ctx, plasmaTx.Transaction.TxHash())
		require.True(t, exists, "returned false for transaction that was stored")
		recoveredTx, ok = txStore.GetTx(ctx, plasmaTx.Transaction.TxHash())
		require.True(t, ok, "error when retrieving transaction")
		require.Equal(t, tx, recoveredTx, fmt.Sprintf("mismatch in stored transaction and retrieved transaction before spends on case %d", i))

		// Spend Outputs
		for j, _ := range plasmaTx.Transaction.Outputs {
			p := plasma.NewPosition(pos.BlockNum, pos.TxIndex, uint8(j), pos.DepositNonce)
			recoveredTx, ok = txStore.GetTxWithPosition(ctx, p)
			require.True(t, ok, "error when retrieving transaction")
			require.True(t, reflect.DeepEqual(tx, recoveredTx), fmt.Sprintf("mismatch in stored transaction and retrieved transaction on case %d", i))

			res := txStore.SpendUTXO(ctx, p, plasmaTx.Transaction.MerkleHash())
			require.True(t, res.IsOK(), "returned error when spending utxo")
			res = txStore.SpendUTXO(ctx, p, plasmaTx.Transaction.MerkleHash())
			require.Equal(t, res.Code, CodeOutputSpent, "allowed output to be spent twice")

			tx.Spent[j] = true
			tx.Spenders[j] = plasmaTx.Transaction.MerkleHash()
			recoveredTx, ok = txStore.GetTxWithPosition(ctx, p)
			require.True(t, ok, "error when retrieving transaction")
			require.True(t, reflect.DeepEqual(tx, recoveredTx), "mismatch in stored and retrieved transaction")
			recoveredTx, ok = txStore.GetTx(ctx, plasmaTx.Transaction.TxHash())
			require.True(t, ok, "error when retrieving transaction")
			require.True(t, reflect.DeepEqual(tx, recoveredTx), "mismatch in stored and retrieved transaction")
		}
	}
}
