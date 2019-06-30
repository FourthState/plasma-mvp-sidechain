package store

import (
	"fmt"
	"github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"math/big"
	"reflect"
	"testing"
)

// Test Get, Has, Store, Spend functions for deposits
func TestDeposits(t *testing.T) {
	ctx, key := setup()
	outputStore := NewOutputStore(key)

	addr := common.BytesToAddress([]byte("asdfasdf"))
	for i := int64(1); i <= 15; i++ {
		nonce := big.NewInt(i)
		hash := []byte("hash that the deposit was spent in")

		// Retrieve/Spend nonexistent deposits
		exists := outputStore.HasDeposit(ctx, nonce)
		require.False(t, exists, "returned true for nonexistent deposit")
		recoveredDeposit, ok := outputStore.GetDeposit(ctx, nonce)
		require.Empty(t, recoveredDeposit, "did not return empty struct for nonexistent deposit")
		require.False(t, ok, "did not return error on nonexistent deposit")
		res := outputStore.SpendDeposit(ctx, nonce, hash)
		require.Equal(t, res.Code, CodeOutputDNE, "did not return that deposit does not exist")

		// Create and store new deposit
		plasmaDeposit := plasma.NewDeposit(addr, big.NewInt(i*4123), big.NewInt(i*123))
		deposit := Deposit{plasmaDeposit, false, []byte{}}
		outputStore.StoreDeposit(ctx, nonce, plasmaDeposit)

		exists = outputStore.HasDeposit(ctx, nonce)
		require.True(t, exists, "returned false for deposit that was stored")
		recoveredDeposit, ok = outputStore.GetDeposit(ctx, nonce)
		require.True(t, ok, "error when retrieving deposit")
		require.True(t, reflect.DeepEqual(deposit, recoveredDeposit), fmt.Sprintf("mismatch in stored deposit and retrieved deposit on iteration %d", i))

		// Spend Deposit
		res = outputStore.SpendDeposit(ctx, nonce, hash)
		require.True(t, res.IsOK(), "returned error when spending deposit")
		res = outputStore.SpendDeposit(ctx, nonce, hash)
		require.Equal(t, res.Code, CodeOutputSpent, "allowed output to be spent twice")

		deposit.Spent = true
		deposit.SpenderTx = hash
		recoveredDeposit, ok = outputStore.GetDeposit(ctx, nonce)
		require.True(t, ok, "error when retrieving deposit")
		require.True(t, reflect.DeepEqual(deposit, recoveredDeposit), "mismatch in stored and retrieved deposit")
	}
}

// Test Get, Has, Store, Spend functions for fees
func TestFees(t *testing.T) {
	ctx, key := setup()
	outputStore := NewOutputStore(key)

	addr := common.BytesToAddress([]byte("asdfasdf"))
	for i := int64(1); i <= 15; i++ {
		pos := plasma.NewPosition(big.NewInt(int64(i)), 1<<16-1, 0, big.NewInt(0))
		hash := []byte("hash that the fee was spent in")

		// Retrieve/Spend nonexistent fee
		exists := outputStore.HasFee(ctx, pos)
		require.False(t, exists, "returned true for nonexistent fee")
		recoveredFee, ok := outputStore.GetFee(ctx, pos)
		require.Empty(t, recoveredFee, "did not return empty struct for nonexistent fee")
		require.False(t, ok, "did not return error on nonexistent fee")
		res := outputStore.SpendFee(ctx, pos, hash)
		require.Equal(t, res.Code, CodeOutputDNE, "did not return that fee does not exist")

		// Create and store new fee
		output := plasma.NewOutput(addr, big.NewInt(int64(1000*i)))
		fee := Output{output, false, make([]byte, 0), make([]byte, 0)}
		outputStore.StoreFee(ctx, pos, output)

		exists = outputStore.HasFee(ctx, pos)
		require.True(t, exists, "returned false for fee that was stored")
		recoveredFee, ok = outputStore.GetFee(ctx, pos)
		require.True(t, ok, "error when retrieving fee")
		require.True(t, reflect.DeepEqual(fee, recoveredFee), fmt.Sprintf("mismatch in stored fee and retrieved fee on iteration %d", i))

		// Spend Fee
		res = outputStore.SpendFee(ctx, pos, hash)
		require.True(t, res.IsOK(), "returned error when spending fee")
		res = outputStore.SpendFee(ctx, pos, hash)
		require.Equal(t, res.Code, CodeOutputSpent, "allowed output to be spent twice")

		fee.Spent = true
		fee.SpenderTx = hash
		recoveredFee, ok = outputStore.GetFee(ctx, pos)
		require.True(t, ok, "error when retrieving fee")
		require.True(t, reflect.DeepEqual(fee, recoveredFee), "mismatch in stored and retrieved fee")
	}
}

// Test Get, Has, Store, Spend functions for transactions
func TestTransactions(t *testing.T) {
	// Setup
	ctx, key := setup()
	outputStore := NewOutputStore(key)

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
				Inputs:  []plasma.Input{plasma.NewInput(getPosition("(122.3.1.0)"), sampleSig, nil), plasma.NewInput(getPosition("(17622.13.5.0)"), sampleSig, nil)},
				Outputs: []plasma.Output{plasma.NewOutput(addr, utils.Big1)},
				Fee:     utils.Big0,
			},
			Position: getPosition("(8765.6847.0.0)"),
		},
		validationCase{
			reason: "tx with one input and two output",
			Transaction: plasma.Transaction{
				Inputs:  []plasma.Input{plasma.NewInput(getPosition("(4.1234.1.0)"), sampleSig, nil)},
				Outputs: []plasma.Output{plasma.NewOutput(addr, utils.Big1), plasma.NewOutput(addr, utils.Big1)},
				Fee:     utils.Big0,
			},
			Position: getPosition("(4354.8765.0.0)"),
		},
		validationCase{
			reason: "tx with two input and one output",
			Transaction: plasma.Transaction{
				Inputs:  []plasma.Input{plasma.NewInput(getPosition("(123.1.0.0)"), sampleSig, nil), plasma.NewInput(getPosition("(10.12.1.0)"), sampleSig, sampleConfirmSig)},
				Outputs: []plasma.Output{plasma.NewOutput(addr, utils.Big1)},
				Fee:     utils.Big0,
			},
			Position: getPosition("(11.123.0.0)"),
		},
		validationCase{
			reason: "tx with two input and two outputs",
			Transaction: plasma.Transaction{
				Inputs:  []plasma.Input{plasma.NewInput(getPosition("(132.231.1.0)"), sampleSig, nil), plasma.NewInput(getPosition("(635.927.1.0)"), sampleSig, sampleConfirmSig)},
				Outputs: []plasma.Output{plasma.NewOutput(addr, utils.Big1), plasma.NewOutput(addr, utils.Big1)},
				Fee:     utils.Big0,
			},
			Position: getPosition("(1121234.12.0.0)"),
		},
	}

	for i, plasmaTx := range txs {
		pos := plasmaTx.Position
		// Retrieve/Spend nonexistent txs
		exists := outputStore.HasTx(ctx, plasmaTx.Transaction.TxHash())
		require.False(t, exists, "returned true for nonexistent transaction")
		recoveredTx, ok := outputStore.GetTx(ctx, plasmaTx.Transaction.TxHash())
		require.Empty(t, recoveredTx, "did not return empty struct for nonexistent transaction")
		require.False(t, ok, "did not return error on nonexistent transaction")

		// Retrieve/Spend nonexistent Outputs
		for j, _ := range plasmaTx.Transaction.Outputs {
			p := plasma.NewPosition(pos.BlockNum, pos.TxIndex, uint8(j), pos.DepositNonce)
			exists = outputStore.HasOutput(ctx, p)
			require.False(t, exists, "returned true for nonexistent output")
			res := outputStore.SpendOutput(ctx, p, plasmaTx.Transaction.MerkleHash())
			require.Equal(t, res.Code, CodeOutputDNE, "did not return that Output does not exist")
		}

		// Create and store new transaction
		tx := Transaction{make([][]byte, len(plasmaTx.Transaction.Outputs)), plasmaTx.Transaction, confirmationHash, make([]bool, len(plasmaTx.Transaction.Outputs)), make([][]byte, len(plasmaTx.Transaction.Outputs)), plasmaTx.Position}
		for i, _ := range tx.SpenderTxs {
			tx.SpenderTxs[i] = []byte{}
			tx.InputTxs[i] = []byte{}
		}
		outputStore.StoreTx(ctx, tx)
		outputStore.StoreOutputs(ctx, tx)

		// Check for outputs
		for j, _ := range plasmaTx.Transaction.Outputs {
			p := plasma.NewPosition(pos.BlockNum, pos.TxIndex, uint8(j), pos.DepositNonce)
			exists = outputStore.HasOutput(ctx, p)
			require.True(t, exists, fmt.Sprintf("returned false for stored output with index %d on case %d", j, i))
		}

		// Check for Tx
		exists = outputStore.HasTx(ctx, plasmaTx.Transaction.TxHash())
		require.True(t, exists, "returned false for transaction that was stored")
		recoveredTx, ok = outputStore.GetTx(ctx, plasmaTx.Transaction.TxHash())
		require.True(t, ok, "error when retrieving transaction")
		require.Equal(t, tx, recoveredTx, fmt.Sprintf("mismatch in stored transaction and retrieved transaction before spends on case %d", i))

		// Spend Outputs
		for j, _ := range plasmaTx.Transaction.Outputs {
			p := plasma.NewPosition(pos.BlockNum, pos.TxIndex, uint8(j), pos.DepositNonce)
			recoveredTx, ok = outputStore.GetTxWithPosition(ctx, p)
			require.True(t, ok, "error when retrieving transaction")
			require.True(t, reflect.DeepEqual(tx, recoveredTx), fmt.Sprintf("mismatch in stored transaction and retrieved transaction on case %d", i))

			res := outputStore.SpendOutput(ctx, p, plasmaTx.Transaction.MerkleHash())
			require.True(t, res.IsOK(), "returned error when spending output")
			res = outputStore.SpendOutput(ctx, p, plasmaTx.Transaction.MerkleHash())
			require.Equal(t, res.Code, CodeOutputSpent, fmt.Sprintf("allowed output with index %d to be spent twice on case %d", j, i))

			tx.Spent[j] = true
			tx.SpenderTxs[j] = plasmaTx.Transaction.MerkleHash()
			recoveredTx, ok = outputStore.GetTxWithPosition(ctx, p)
			require.True(t, ok, "error when retrieving transaction")
			require.True(t, reflect.DeepEqual(tx, recoveredTx), fmt.Sprintf("mismatch in stored and retrieved transaction on case %d", i))
			recoveredTx, ok = outputStore.GetTx(ctx, plasmaTx.Transaction.TxHash())
			require.True(t, ok, "error when retrieving transaction")
			require.True(t, reflect.DeepEqual(tx, recoveredTx), fmt.Sprintf("mismatch in stored and retrieved transaction on case %d", i))
		}
	}
}
