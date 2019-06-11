package eth

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	plasmaTypes "github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	"math/big"
	"testing"
	"time"
)

// private/public keys using the `plasma` mnemonic with ganache-cli
// `ganache-cli -m=plasma`
// plasmaContractAddr will be deterministic. `truffle migrate` immediately after `ganache-cli -m=plasma`
const (
	clientAddr         = "http://127.0.0.1:8545"
	plasmaContractAddr = "31E491FC70cDb231774c61B7F46d94699dacE664"
	operatorPrivKey    = "9cd69f009ac86203e54ec50e3686de95ff6126d3b30a19f926a0fe9323c17181"

	minExitBond = 200000
)

var (
	commitmentRate, _ = time.ParseDuration("1s")
)

func TestConnection(t *testing.T) {
	logger := log.NewNopLogger()

	t.Logf("Connecting to remote client: %s", clientAddr)
	client, err := InitEthConn(clientAddr, logger)
	require.NoError(t, err, "connection error")

	_, err = client.accounts()
	require.NoError(t, err, "error retrieving accounts")
}

func TestLatestBlockNum(t *testing.T) {
	logger := log.NewNopLogger()
	client, _ := InitEthConn(clientAddr, logger)
	_, err := client.LatestBlockNum()
	require.NoError(t, err)
}

func TestPlasmaInit(t *testing.T) {
	logger := log.NewNopLogger()
	client, err := InitEthConn(clientAddr, logger)

	_, err = InitPlasma(common.HexToAddress(plasmaContractAddr), client, 1, commitmentRate, logger, false, nil)
	require.NoError(t, err, "error binding to contract")
}

// Test needs to be changed to simulate the sdk context and plasma store.
func TestSubmitBlock(t *testing.T) {
	logger := log.NewNopLogger()
	client, _ := InitEthConn(clientAddr, logger)

	privKey, _ := crypto.HexToECDSA(operatorPrivKey)
	plasma, _ := InitPlasma(common.HexToAddress(plasmaContractAddr), client, 1, commitmentRate, logger, true, privKey)

	// Setup context and plasma store
	ctx, plasmaStore := setup()

	time.Sleep(2 * time.Second)

	// Submit 2 blocks
	var expectedBlocks []plasmaTypes.Block
	for i := 1; i < 3; i++ {
		block := plasmaTypes.Block{
			Header:    sha256.Sum256([]byte(fmt.Sprintf("Block: %d", i))),
			TxnCount:  uint16(i + 1),
			FeeAmount: big.NewInt(int64(i + 2)),
		}
		expectedBlocks = append(expectedBlocks, block)
		plasmaStore.StoreBlock(ctx, big.NewInt(int64(i)), block)
	}

	err := plasma.CommitPlasmaHeaders(ctx, plasmaStore)

	require.NoError(t, err, "block submission error")

	blockNum, err := plasma.contract.LastCommittedBlock(nil)
	require.NoError(t, err, "failed to query for the last committed block")
	require.Equal(t, big.NewInt(2), blockNum, "Did not submit both blocks correctly")

	for j := 0; j < 2; j++ {
		result, err := plasma.contract.PlasmaChain(nil, big.NewInt(int64(j+1)))
		require.NoError(t, err, "failed contract plasma chain query")

		require.Truef(t, bytes.Compare(expectedBlocks[j].Header[:], result.Header[:]) == 0,
			"Mismatch in block headers for submitted block %d. Got: %x. Expected: %x", j, result.Header[:], expectedBlocks[j].Header[:])

		require.Equal(t, big.NewInt(int64(expectedBlocks[j].TxnCount)), result.NumTxns, fmt.Sprintf("Wrong number of tx's for submitted block: %d", j))

		require.Equal(t, expectedBlocks[j].FeeAmount, result.FeeAmount, fmt.Sprintf("Wrong Fee amount for submitted block: %d", j))

	}
}

func TestDepositFinalityBound(t *testing.T) {
	logger := log.NewNopLogger()
	client, _ := InitEthConn(clientAddr, logger)

	privKey, _ := crypto.HexToECDSA(operatorPrivKey)
	// finality bound of 2 ethereum blocks
	plasma, _ := InitPlasma(common.HexToAddress(plasmaContractAddr), client, 1, commitmentRate, logger, true, privKey)

	// mine a block so that the headers channel is filled with a block
	err := client.rpc.Call(nil, "evm_mine")
	require.NoError(t, err, "error mining a block")
	time.Sleep(1 * time.Second)

	nonce, err := plasma.contract.DepositNonce(nil)
	require.NoError(t, err, "error querying for the deposit nonce")

	// Deposit 10 eth from the operator
	plasma.operatorSession.TransactOpts.Value = big.NewInt(10)
	operatorAddress := crypto.PubkeyToAddress(privKey.PublicKey)
	_, err = plasma.operatorSession.Deposit(operatorAddress)
	require.NoError(t, err, "error sending a deposit tx")

	// Setup context and plasma store
	ctx, plasmaStore := setup()

	// Reset operatorSession
	plasma.operatorSession.TransactOpts.Value = nil

	var block plasmaTypes.Block
	// Must restore old blocks since we're using fresh plasmaStore but using old contract
	// that already has submitted blocks. Store blocks 1-3 to get plasmaConn to submit new block 3
	for i := 1; i < 4; i++ {
		block = plasmaTypes.Block{
			Header:    sha256.Sum256([]byte(fmt.Sprintf("Block: %d", i))),
			TxnCount:  uint16(i + 1),
			FeeAmount: big.NewInt(int64(i + 2)),
		}
		plasmaStore.StoreBlock(ctx, big.NewInt(int64(i)), block)
	}

	err = plasma.CommitPlasmaHeaders(ctx, plasmaStore)

	require.NoError(t, err, "block submission error")

	err = plasma.CommitPlasmaHeaders(ctx, plasmaStore)
	require.NoError(t, err, "block submission error")

	// Try to retrieve deposit from before peg
	_, threshold, ok := plasma.GetDeposit(big.NewInt(2), nonce)
	require.False(t, ok, "retrieved a deposit that occurred after pegged block")
	require.Equal(t, big.NewInt(3), threshold, "Finality threshold calculated incorrectly. Should still need to wait two more blocks")

	/* Mine 3 blocks for finality bound */
	for i := 0; i < 3; i++ {
		// mine another block so that the deposit falls outside the finality bound
		err = client.rpc.Call(nil, "evm_mine")
		require.NoError(t, err, "error mining a block")
		time.Sleep(1 * time.Second)
	}

	/* Submit block to advance peg */
	block = plasmaTypes.Block{
		Header:    sha256.Sum256([]byte("Block: 4")),
		TxnCount:  uint16(2),
		FeeAmount: big.NewInt(3),
	}
	plasmaStore.StoreBlock(ctx, big.NewInt(4), block)

	err = plasma.CommitPlasmaHeaders(ctx, plasmaStore)
	require.NoError(t, err, "block submission error")

	// Try to retrieve deposit once peg has advanced AND finality bound reached.
	deposit, threshold, ok := plasma.GetDeposit(big.NewInt(4), nonce)
	require.True(t, ok, "could not retrieve a deposit that was deemed final")

	require.Equal(t, uint64(10), deposit.Amount.Uint64(), "deposit amount mismatch")
	require.True(t, bytes.Equal(operatorAddress[:], deposit.Owner[:]), "deposit owner mismatch")
	require.True(t, threshold.Sign() == 0, "Finality threshold not calculated correctly. Deposit should be final with threshold = 0")
}
