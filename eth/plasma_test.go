package eth

import (
	"bytes"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tendermint/tendermint/libs/log"
	"math/big"
	"os"
	"testing"
)

// private/public keys using the `plasma` mnemonic with ganache-cli
// `ganache-cli -m=plasma`
const (
	clientAddr         = "ws://127.0.0.1:8545"
	plasmaContractAddr = "31e491fc70cdb231774c61b7f46d94699dace664"
	operatorPrivKey    = "9cd69f009ac86203e54ec50e3686de95ff6126d3b30a19f926a0fe9323c17181"
	sampleAccount      = "66b2e0a229d38764cea81dc99bfbd1eb85354b33"
)

func TestConnection(t *testing.T) {
	logger := log.NewTMLogger(os.Stderr)

	t.Logf("Connecting to remote client: %s", clientAddr)
	client, err := InitEthConn(clientAddr, logger)
	if err != nil {
		t.Fatal("Connection Error -", err)
	}

	_, err = client.accounts()
	if err != nil {
		t.Error("Error Retrieving Accounts -", err)
	}
}

func TestPlasmaInit(t *testing.T) {
	logger := log.NewTMLogger(os.Stderr)
	client, err := InitEthConn(clientAddr, logger)

	t.Logf("Binding go wrapper to deployed contract: %s", plasmaContractAddr)
	privKey, err := crypto.HexToECDSA(operatorPrivKey)
	if err != nil {
		t.Fatal("Could not convert hex private key")
	}

	plasma, err := InitPlasma(plasmaContractAddr, privKey, client, logger, 1)
	if err != nil {
		t.Fatal("Could not bind contract -", err)
	}

	// sample call
	balance, err := plasma.session.ChildChainBalance()
	if err != nil {
		t.Fatal("Could not query contract's plasma balance -", err)
	}

	if balance.Sign() != 0 {
		t.Error("Incorrectly result for contract's plasma balance. Expected 0")
	}
}

func TestSubmitBlock(t *testing.T) {
	logger := log.NewTMLogger(os.Stderr)
	client, _ := InitEthConn(clientAddr, logger)

	privKey, _ := crypto.HexToECDSA(operatorPrivKey)
	plasma, _ := InitPlasma(plasmaContractAddr, privKey, client, logger, 1)

	blockNum, err := plasma.session.LastCommittedBlock()
	if err != nil {
		t.Fatal("Could not query for the last committed block -", err)
	}

	blockNum = blockNum.Add(blockNum, big.NewInt(1))

	header := crypto.Keccak256([]byte("blah"))
	_, err = plasma.SubmitBlock(header, sdk.NewUint(0), sdk.NewUint(0))
	if err != nil {
		t.Fatal("Could not submit block -", err)
	}

	// mine a block
	err = client.rpc.Call(nil, "evm_mine")
	if err != nil {
		t.Fatal("Failed rpc call -", err)
	}

	result, err := plasma.session.ChildChain(blockNum)
	if err != nil {
		t.Fatal("Could not query the child chain - ", err)
	}

	if bytes.Compare(result.Root[:], header) != 0 {
		t.Errorf("Mismatch in block headers.\nGot: %x\nExpected: %x", result, header)
	}
}

func TestEthBlockWatching(t *testing.T) {
	logger := log.NewTMLogger(os.Stderr)
	client, _ := InitEthConn(clientAddr, logger)

	privKey, _ := crypto.HexToECDSA(operatorPrivKey)
	plasma, _ := InitPlasma(plasmaContractAddr, privKey, client, logger, 1)

	lastBlockNum := plasma.ethBlockNum
	t.Log(lastBlockNum)

	// mine a block
	err := client.rpc.Call(nil, "evm_mine")
	if err != nil {
		t.Fatal("Could not mine a block -", err)
	}

	t.Log(plasma.ethBlockNum)

	if plasma.ethBlockNum.Cmp(lastBlockNum.Add(lastBlockNum, big.NewInt(1))) != 0 {
		t.Error("Client did not catch the minted block and update correctly")
	}
}
