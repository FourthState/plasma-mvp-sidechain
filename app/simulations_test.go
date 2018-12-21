package app

import (
	"crypto/ecdsa"
	"math/big"
	"testing"
	"time"

	contracts "github.com/FourthState/plasma-mvp-sidechain/contracts/wrappers"
	//"github.com/FourthState/plasma-mvp-sidechain/eth"
	"github.com/FourthState/plasma-mvp-sidechain/types"
	utils "github.com/FourthState/plasma-mvp-sidechain/utils"
	"github.com/FourthState/plasma-mvp-sidechain/x/utxo"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethcmn "github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	rlp "github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
)

const (
	minExitBond = 10000
)

// create a new session with plasma contract
func newSession(privateKey *ecdsa.PrivateKey, nodeUrl string) (*contracts.PlasmaMVPSession, error) {
	c, err := rpc.Dial(nodeUrl)
	if err != nil {
		return nil, err
	}
	ec := ethclient.NewClient(c)

	plasmaContract, err := contracts.NewPlasmaMVP(ethcmn.HexToAddress(plasmaContractAddr), ec)
	if err != nil {
		return nil, err
	}

	// Create a session with the contract and operator account
	auth := bind.NewKeyedTransactor(privateKey)
	plasmaSession := &contracts.PlasmaMVPSession{
		Contract: plasmaContract,
		CallOpts: bind.CallOpts{
			Pending: true,
		},
		TransactOpts: bind.TransactOpts{
			From:     auth.From,
			Signer:   auth.Signer,
			GasLimit: 3141592, // aribitrary
		},
	}
	return plasmaSession, nil
}

// deposit 100 with passed in address
func deposit() (*contracts.PlasmaMVPSession, *ChildChain, ethcmn.Address, uint64, *ecdsa.PrivateKey, error) {
	privKey, _ := ethcrypto.HexToECDSA(privkey)
	session, err := newSession(privKey, nodeURL)
	if err != nil {
		return nil, nil, ethcmn.Address{}, 0, nil, err
	}

	cc := newChildChain()

	cc.BeginBlock(abci.RequestBeginBlock{})
	cc.EndBlock(abci.RequestEndBlock{})
	cc.Commit()

	nonce, err := session.DepositNonce()
	if err != nil {
		return nil, nil, ethcmn.Address{}, 0, nil, err
	}

	// Deposit 100 eth from the validator
	session.TransactOpts.Value = big.NewInt(100)
	addr := ethcrypto.PubkeyToAddress(privKey.PublicKey)
	if _, err = session.Deposit(addr); err != nil {
		return nil, nil, ethcmn.Address{}, 0, nil, err
	}

	time.Sleep(500 * time.Millisecond)
	return session, cc, addr, uint64(nonce.Int64()), privKey, nil
}

// AddrA deposits and spend to AddrB
func TestDepositAndSpend(t *testing.T) {
	_, cc, addrA, nonce, privKey, err := deposit()
	require.NoError(t, err)

	// Spend Deposit
	privKeyB, _ := ethcrypto.GenerateKey()
	addrB := utils.PrivKeyToAddress(privKeyB)
	msg := GenerateSimpleMsg(addrA, addrB, [4]uint64{0, 0, 0, nonce}, 100, 0)
	tx := GetTx(msg, privKey, nil, false)
	txBytes, _ := rlp.EncodeToBytes(tx)

	// Simulate a block
	cc.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 1}})

	// Deliver tx, updates states
	dres := cc.DeliverTx(txBytes)

	require.Equal(t, sdk.CodeType(0), sdk.CodeType(dres.Code), dres.Log)

	// Create context
	ctx := cc.NewContext(false, abci.Header{})

	// Retrieve UTXO from context
	position := types.NewPlasmaPosition(1, 0, 0, 0)
	res := cc.utxoMapper.GetUTXO(ctx, addrB.Bytes(), position)

	inputKey := cc.utxoMapper.ConstructKey(addrA.Bytes(), types.NewPlasmaPosition(0, 0, 0, nonce))
	txHash := tmhash.Sum(txBytes)
	expected := utxo.NewUTXOwithInputs(addrB.Bytes(), 100, "Ether", position, txHash, [][]byte{inputKey})

	require.Equal(t, expected, res, "UTXO did not get added to store correctly")
}

// AddrA deposits and exits
// spend attempt fails
func TestDepositExit(t *testing.T) {
	session, cc, addrA, nonce, privKey, err := deposit()
	require.NoError(t, err)

	// exit deposit
	session.TransactOpts.Value = big.NewInt(minExitBond)
	_, err = session.StartDepositExit(big.NewInt(int64(nonce)), big.NewInt(0))
	require.NoError(t, err)
	time.Sleep(500 * time.Millisecond)

	// attempt spend
	privKeyB, _ := ethcrypto.GenerateKey()
	addrB := utils.PrivKeyToAddress(privKeyB)
	msg := GenerateSimpleMsg(addrA, addrB, [4]uint64{0, 0, 0, nonce}, 100, 0)
	tx := GetTx(msg, privKey, nil, false)
	txBytes, _ := rlp.EncodeToBytes(tx)

	// Simulate a block
	cc.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: 1}})

	// Deliver tx, updates states
	dres := cc.DeliverTx(txBytes)

	require.Equal(t, sdk.CodeType(204), sdk.CodeType(dres.Code), dres.Log)
}

// deposit, spend deposit, exit utxo
// attempt spend of utxo, assert that this fails
func TestUTXOExitSpend(t *testing.T) {
	session, cc, addrA, nonce, privKey, err := deposit()
	require.NoError(t, err)

	// Spend Deposit
	msg := GenerateSimpleMsg(addrA, addrA, [4]uint64{0, 0, 0, nonce}, 100, 0)
	tx := GetTx(msg, privKey, nil, false)
	txBytes, _ := rlp.EncodeToBytes(tx)

	lastCommittedBlock, _ := session.LastCommittedBlock()
	blknum := lastCommittedBlock.Int64() + 1
	// Simulate a block
	cc.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: blknum}})

	// Deliver tx, updates states
	dres := cc.DeliverTx(txBytes)
	require.Equal(t, sdk.CodeType(0), sdk.CodeType(dres.Code), dres.Log)

	cc.EndBlock(abci.RequestEndBlock{})
	cc.Commit()
	cc.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: blknum + 1}})

	// submit block header
	// merkle root == block hash
	var blockHash [32]byte
	txHash := tmhash.Sum(txBytes)
	copy(blockHash[:], txHash[:])

	session.TransactOpts.Value = nil
	_, err = session.SubmitBlock(
		[][32]byte{blockHash},
		[]*big.Int{big.NewInt(1)},
		[]*big.Int{big.NewInt(0)},
		big.NewInt(blknum),
	)
	require.NoError(t, err)

	// exit utxo
	hash := tmhash.Sum(append(txHash, txHash[:]...))

	confirmSigs := CreateConfirmSig(hash, privKey, &ecdsa.PrivateKey{}, false)

	session.TransactOpts.Value = big.NewInt(minExitBond)
	_, err = session.StartTransactionExit([3]*big.Int{big.NewInt(blknum), big.NewInt(0), big.NewInt(0)}, txBytes, []byte{}, confirmSigs[0][:], big.NewInt(0))
	require.NoError(t, err)
	time.Sleep(500 * time.Millisecond)

	msg = GenerateSimpleMsg(addrA, addrA, [4]uint64{uint64(blknum), 0, 0, 0}, 100, 0)
	msg.Input0ConfirmSigs = confirmSigs

	tx = GetTx(msg, privKey, nil, false)
	txBytes, _ = rlp.EncodeToBytes(tx)

	// Deliver tx, updates states
	dres = cc.DeliverTx(txBytes)

	require.Equal(t, sdk.CodeType(204), sdk.CodeType(dres.Code), dres.Log)

}
