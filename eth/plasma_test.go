package eth

import (
	"bytes"
	"crypto/sha256"
	plasmaTypes "github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/tendermint/tendermint/libs/log"
	"math/big"
	"os"
	"reflect"
	"testing"
	"time"
)

// private/public keys using the `plasma` mnemonic with ganache-cli
// `ganache-cli -m=plasma`
// plasmaContractAddr will be deterministic. `truffle migrate` immediately after `ganache-cli -m=plasma`
const (
	clientAddr         = "ws://127.0.0.1:8545"
	plasmaContractAddr = "5cae340fb2c2bb0a2f194a95cda8a1ffdc9d2f85"
	operatorPrivKey    = "9cd69f009ac86203e54ec50e3686de95ff6126d3b30a19f926a0fe9323c17181"

	minExitBond = 10000
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

	privKey, _ := crypto.HexToECDSA(operatorPrivKey)
	_, err = InitPlasma(common.HexToAddress(plasmaContractAddr), privKey, client, logger, 1)
	if err != nil {
		t.Fatal("Could not bind contract -", err)
	}
}

func TestSubmitBlock(t *testing.T) {
	logger := log.NewTMLogger(os.Stderr)
	client, _ := InitEthConn(clientAddr, logger)

	privKey, _ := crypto.HexToECDSA(operatorPrivKey)
	plasma, _ := InitPlasma(common.HexToAddress(plasmaContractAddr), privKey, client, logger, 1)

	header := crypto.Keccak256([]byte("blah"))
	var root [32]byte
	copy(root[:], header)
	_, err := plasma.SubmitBlock(root, big.NewInt(0), big.NewInt(0))
	if err != nil {
		t.Fatal("Failed block submission -", err)
	}

	blockNum, err := plasma.session.LastCommittedBlock()
	if err != nil {
		t.Fatal("Failed query for the last committed block -", err)
	}

	result, err := plasma.session.ChildChain(blockNum)
	if err != nil {
		t.Fatal("Failed query for the child chain - ", err)
	}

	if bytes.Compare(result.Root[:], header) != 0 {
		t.Errorf("Mismatch in block headers. Got: %x. Expected: %x", result, header)
	}
}

func TestEthBlockWatching(t *testing.T) {
	logger := log.NewTMLogger(os.Stderr)
	client, _ := InitEthConn(clientAddr, logger)

	privKey, _ := crypto.HexToECDSA(operatorPrivKey)
	plasma, _ := InitPlasma(common.HexToAddress(plasmaContractAddr), privKey, client, logger, 1)

	// mine a block so that `ethBlockNum` within plasma gets set
	// sleep after an rpc call to deal with the asynchrony
	if err := client.rpc.Call(nil, "evm_mine"); err != nil {
		t.Fatal("Could not mine a block -", err)
	}
	time.Sleep(1 * time.Second)

	plasma.lock.Lock()
	lastEthBlockNum := plasma.ethBlockNum.Uint64()
	plasma.lock.Unlock()

	// mine another block that should get caught
	if err := client.rpc.Call(nil, "evm_mine"); err != nil {
		t.Fatal("Could not mine a block -", err)
	}
	time.Sleep(1 * time.Second)

	plasma.lock.Lock()
	currEthBlockNum := plasma.ethBlockNum.Uint64()
	plasma.lock.Unlock()
	if currEthBlockNum != lastEthBlockNum+1 {
		t.Fatalf("EthBlockNum not incremented. Expected: %d, Got: %d",
			lastEthBlockNum+1, currEthBlockNum)
	}
}

func TestDepositWatching(t *testing.T) {
	logger := log.NewTMLogger(os.Stderr)
	client, _ := InitEthConn(clientAddr, logger)

	privKey, _ := crypto.HexToECDSA(operatorPrivKey)
	plasma, _ := InitPlasma(common.HexToAddress(plasmaContractAddr), privKey, client, logger, 0)

	nonce, err := plasma.session.DepositNonce()
	if err != nil {
		t.Fatalf("Could not query for the next deposit nonce")
	}

	// Deposit 10 eth from the operator
	plasma.session.TransactOpts.Value = big.NewInt(10)
	operatorAddress := crypto.PubkeyToAddress(privKey.PublicKey)
	_, err = plasma.session.Deposit(operatorAddress)
	if err != nil {
		t.Fatalf("Error sending a deposit tx")
	}
	time.Sleep(500 * time.Millisecond)

	deposit, ok := plasma.GetDeposit(nonce)
	if !ok {
		t.Fatal("Deposit not caught")
	}

	if deposit.Amount.Uint64() != uint64(10) {
		t.Errorf("Deposit amount incorrect. Expected 10, Got %d", deposit.Amount.Uint64())
	}

	if !bytes.Equal(operatorAddress[:], deposit.Owner[:]) {
		t.Errorf("Deposit owner incorrect. Expected %x, Got %x", operatorAddress, deposit.Owner)
	}

	// check persistence in the db
	key := prefixKey(depositPrefix, nonce.Bytes())
	data, err := plasma.memdb.Get(key)
	if err != nil {
		t.Fatalf("Deposit not persisted - %s", err)
	}

	var d plasmaTypes.Deposit
	err = rlp.DecodeBytes(data, &d)
	if err != nil {
		t.Fatalf("Error unmarshaling cache'd deposit - %s", err)
	}

	if !reflect.DeepEqual(deposit, &d) {
		t.Fatalf("Mismatch in the persisted deposit and `GetDeposit`")
	}
}

func TestDepositExitWatching(t *testing.T) {
	logger := log.NewTMLogger(os.Stderr)
	client, _ := InitEthConn(clientAddr, logger)

	privKey, _ := crypto.HexToECDSA(operatorPrivKey)
	plasma, _ := InitPlasma(common.HexToAddress(plasmaContractAddr), privKey, client, logger, 0)

	// deposit and exit
	nonce, _ := plasma.session.DepositNonce()
	plasma.session.TransactOpts.Value = big.NewInt(10)
	_, err := plasma.session.Deposit(crypto.PubkeyToAddress(privKey.PublicKey))
	if err != nil {
		t.Fatal("Failed deposit -", err)
	}

	plasma.session.TransactOpts.Value = big.NewInt(minExitBond)
	_, err = plasma.session.StartDepositExit(nonce, big.NewInt(0))
	if err != nil {
		t.Fatal("Error starting deposit exit -", err)
	}
	time.Sleep(500 * time.Millisecond)

	zero := big.NewInt(0)
	position := plasmaTypes.NewPosition(zero, 0, 0, nonce)
	exited := plasma.HasTXBeenExited(position)

	if !exited {
		t.Errorf("Deposit nonce not marked as exited")
	}
}

type SpendMsg struct {
	Blknum0           uint64
	Txindex0          uint16
	Oindex0           uint8
	DepositNum0       uint64
	Owner0            common.Address
	Input0ConfirmSigs [][65]byte
	Blknum1           uint64
	Txindex1          uint16
	Oindex1           uint8
	DepositNum1       uint64
	OWner1            common.Address
	Input1ConfirmSigs [][65]byte
	Newowner0         common.Address
	Amount0           uint64
	Newowner1         common.Address
	Amount1           uint64
	FeeAmount         uint64
}

type tx struct {
	TxList SpendMsg
	Sigs   [2][]byte
}

func toEthSignedMessageHash(msg []byte) []byte {
	buffer := new(bytes.Buffer)
	buffer.Write([]byte("\x19Ethereum Signed Message:\n32"))
	buffer.Write(msg)
	return crypto.Keccak256(buffer.Bytes())
}

func TestTxExitWatchingAndChallenge(t *testing.T) {
	logger := log.NewTMLogger(os.Stderr)
	client, _ := InitEthConn(clientAddr, logger)

	privKey, _ := crypto.HexToECDSA(operatorPrivKey)
	plasma, _ := InitPlasma(common.HexToAddress(plasmaContractAddr), privKey, client, logger, 0)
	zero := big.NewInt(0)

	// deposit and spend
	nonce, _ := plasma.session.DepositNonce()
	plasma.session.TransactOpts.Value = big.NewInt(10)
	_, err := plasma.session.Deposit(crypto.PubkeyToAddress(privKey.PublicKey))
	if err != nil {
		t.Fatal("Failed deposit -", err)
	}

	// generate tx
	var msg SpendMsg
	msg.DepositNum0 = nonce.Uint64()
	msg.Newowner0 = crypto.PubkeyToAddress(privKey.PublicKey)
	msg.Amount0 = 10
	txList, _ := rlp.EncodeToBytes(msg)

	txHash := crypto.Keccak256(txList)
	txHash = toEthSignedMessageHash(txHash)
	sig0, _ := crypto.Sign(txHash, privKey)
	sigs := [2][]byte{sig0, make([]byte, 65)}

	txBytes, _ := rlp.EncodeToBytes(tx{msg, sigs})

	// submit header. header == merklehash
	header := sha256.Sum256(txBytes)
	plasma.session.TransactOpts.Value = nil
	_, err = plasma.SubmitBlock(header, big.NewInt(1), zero)
	if err != nil {
		t.Fatal("Error submitting block -", err)
	}

	time.Sleep(500 * time.Millisecond)

	// merkleHash == header
	var data []byte
	data = append(data, header[:]...)
	data = append(data, header[:]...)
	confirmationHash := sha256.Sum256(data)
	confHash := toEthSignedMessageHash(confirmationHash[:])
	confirmSignature, _ := crypto.Sign(confHash, privKey)

	plasma.session.TransactOpts.Value = big.NewInt(minExitBond)
	_, err = plasma.session.StartTransactionExit([3]*big.Int{plasma.blockNum, zero, zero}, txBytes, []byte{}, confirmSignature, zero)
	if err != nil {
		t.Fatal("Error starting tx exit -", err)
	}
	time.Sleep(500 * time.Millisecond)

	txPos := plasmaTypes.NewPosition(plasma.blockNum, 0, 0, zero)
	exited := plasma.HasTXBeenExited(txPos)
	if !exited {
		t.Errorf("Transaction not marked as exited")
	}

	// attempt to exit the deposit & challenge
	depositPos := plasmaTypes.NewPosition(zero, 0, 0, nonce)
	_, err = plasma.session.StartDepositExit(nonce, zero)
	if err != nil {
		t.Fatal("Error exiting deposit -", err)
	}
	time.Sleep(500 * time.Millisecond)

	exited = plasma.HasTXBeenExited(depositPos)
	if !exited {
		t.Errorf("Deposit not marked as exited after exiting")
	}
	plasma.session.TransactOpts.Value = nil
	_, err = plasma.session.ChallengeExit(depositPos.ToBigIntArray(), [2]*big.Int{plasma.blockNum, zero}, txBytes, []byte{}, confirmSignature)
	if err != nil {
		t.Fatal("Error challenging exit -", err)
	}
	time.Sleep(500 * time.Millisecond)

	exited = plasma.HasTXBeenExited(depositPos)
	if exited {
		t.Errorf("Deposit marked as exited after being challenged")
	}
}
