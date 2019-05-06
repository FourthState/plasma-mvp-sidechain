package eth

import (
	"bytes"
	//plasmaTypes "github.com/FourthState/plasma-mvp-sidechain/plasma"
	//"github.com/FourthState/plasma-mvp-sidechain/utils"
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
	plasmaContractAddr = "5cae340fb2c2bb0a2f194a95cda8a1ffdc9d2f85"
	operatorPrivKey    = "9cd69f009ac86203e54ec50e3686de95ff6126d3b30a19f926a0fe9323c17181"

	minExitBond = 200000
)

var (
	commitmentRate, _ = time.ParseDuration("1m")
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

/* Test needs to be changed to simulate the sdk context and plasma store.
func TestSubmitBlock(t *testing.T) {
	logger := log.NewNopLogger()
	client, _ := InitEthConn(clientAddr, logger)

	privKey, _ := crypto.HexToECDSA(operatorPrivKey)
	plasma, _ := InitPlasma(common.HexToAddress(plasmaContractAddr), client, 1, commitmentRate, logger, true, privKey)

	var header [32]byte
	copy(header[:], crypto.Keccak256([]byte("blah")))
	block := plasmaTypes.Block{
		Header:    header,
		TxnCount:  1,
		FeeAmount: utils.Big0,
	}
	err := plasma.SubmitBlock(block)
	require.NoError(t, err, "block submission error")

	blockNum, err := plasma.contract.LastCommittedBlock(nil)
	require.NoError(t, err, "failed to query for the last committed block")

	result, err := plasma.contract.PlasmaChain(nil, blockNum)
	require.NoError(t, err, "failed contract plasma chain query")

	require.Truef(t, bytes.Compare(result.Header[:], header[:]) == 0,
		"Mismatch in block headers. Got: %x. Expected: %x", result, header)
}
*/

func TestDepositFinalityBound(t *testing.T) {
	logger := log.NewNopLogger()
	client, _ := InitEthConn(clientAddr, logger)

	privKey, _ := crypto.HexToECDSA(operatorPrivKey)
	// finality bound of 2 ethereum blocks
	plasma, _ := InitPlasma(common.HexToAddress(plasmaContractAddr), client, 2, commitmentRate, logger, true, privKey)

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

	_, threshold, ok := plasma.GetDeposit(big.NewInt(1), nonce)
	require.False(t, ok, "retrieved a deposit that occurred after pegged block")
	require.Equal(t, big.NewInt(2), threshold, "Finality threshold calculated incorrectly. Should still need to wait two more blocks")

	/* Mine 3 blocks for finality bound */
	for i := 0; i < 3; i++ {
		// mine another block so that the deposit falls outside the finality bound
		err = client.rpc.Call(nil, "evm_mine")
		require.NoError(t, err, "error mining a block")
		time.Sleep(1 * time.Second)
	}

	deposit, threshold, ok := plasma.GetDeposit(big.NewInt(1), nonce)
	require.True(t, ok, "could not retrieve a deposit that was deemed final")

	require.Equal(t, uint64(10), deposit.Amount.Uint64(), "deposit amount mismatch")
	require.True(t, bytes.Equal(operatorAddress[:], deposit.Owner[:]), "deposit owner mismatch")
	require.True(t, threshold.Sign() <= 0, "Finality threshold not calculated correctly. Deposit should be final with threshold = 0")
}
