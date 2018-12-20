package eth

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	contracts "github.com/FourthState/plasma-mvp-sidechain/contracts/wrappers"
	plasmaTypes "github.com/FourthState/plasma-mvp-sidechain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/comparer"
	"github.com/syndtr/goleveldb/leveldb/memdb"
	"github.com/tendermint/tendermint/libs/log"
	"math/big"
	"sync"
)

// Plasma holds related unexported members
type Plasma struct {
	session *contracts.PlasmaMVPSession
	client  *Client
	logger  log.Logger

	memdb *memdb.DB
	db    *leveldb.DB

	blockNum      *big.Int
	ethBlockNum   *big.Int
	finalityBound uint64

	lock *sync.Mutex
}

// InitPlasma binds the go wrapper to the deployed contract. This private key provides authentication for the operator
func InitPlasma(contractAddr common.Address, privateKey *ecdsa.PrivateKey, client *Client, logger log.Logger, finalityBound uint64) (*Plasma, error) {
	plasmaContract, err := contracts.NewPlasmaMVP(contractAddr, client.ec)
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
			GasLimit: 3141592, // aribitrary. TODO: check this
		},
	}

	// TODO: deal with syncing issues
	lastCommittedBlock, err := plasmaSession.LastCommittedBlock()
	if err != nil {
		return nil, fmt.Errorf("Contract session not correctly established - %s", err)
	}

	plasma := &Plasma{
		session: plasmaSession,
		client:  client,
		// capacity argument is advisory and not enforced in the memdb implementation
		// TODO: flush the in-memory DB to a local one to bound memory consumption
		memdb:  memdb.New(comparer.DefaultComparer, 1),
		logger: logger,

		ethBlockNum:   big.NewInt(-1),
		blockNum:      lastCommittedBlock,
		finalityBound: finalityBound,

		lock: &sync.Mutex{},
	}

	// listen to new ethereum block headers
	ethCh, err := client.SubscribeToHeads()
	if err != nil {
		logger.Error("Could not successfully subscribe to heads: %v", err)
		return nil, err
	}

	// start listeners
	go watchEthBlocks(plasma, ethCh)
	go watchDeposits(plasma)
	go watchExits(plasma)

	return plasma, nil
}

// SubmitBlock proxy. TODO: handle batching with a timmer interrupt
func (plasma *Plasma) SubmitBlock(header [32]byte, numTxns *big.Int, fee *big.Int) (*types.Transaction, error) {
	plasma.blockNum = plasma.blockNum.Add(plasma.blockNum, big.NewInt(1))

	tx, err := plasma.session.SubmitBlock(
		[][32]byte{header},
		[]*big.Int{numTxns},
		[]*big.Int{fee},
		plasma.blockNum)

	if err != nil {
		return nil, err
	}

	return tx, nil
}

// GetDeposit checks the existence of a deposit nonce
func (plasma *Plasma) GetDeposit(nonce *big.Int) (*plasmaTypes.Deposit, error) {
	key := prefixKey(depositPrefix, nonce.Bytes())
	data, err := plasma.memdb.Get(key)

	var deposit plasmaTypes.Deposit
	// check against the contract if the deposit is not in the cache or decoding fails
	if err != nil || json.Unmarshal(data, &deposit) != nil {
		if plasma.memdb.Contains(key) {
			plasma.logger.Info("corrupted deposit found within db")
			plasma.memdb.Delete(key)
		}

		d, err := plasma.session.Deposits(nonce)
		if err != nil {
			plasma.logger.Error("contract call, deposits, failed")
			return nil, err
		}

		if d.CreatedAt.Sign() == 0 {
			return nil, fmt.Errorf("deposit does not exist")
		}

		deposit = plasmaTypes.Deposit{
			Owner:    d.Owner,
			Amount:   sdk.NewUintFromBigInt(d.Amount),
			BlockNum: sdk.NewUintFromBigInt(d.EthBlockNum),
		}

		// save to the db
		data, err = json.Marshal(deposit)
		if err != nil {
			plasma.logger.Error("error encoding deposit. will not be cached")
		} else {
			plasma.memdb.Put(key, data)
		}
	}

	// check finality bound for the deposit
	plasma.lock.Lock()
	ethBlockNum := plasma.ethBlockNum
	plasma.lock.Unlock()
	if ethBlockNum.Sign() < 0 {
		return nil, fmt.Errorf("not subscribed to ethereum block headers")
	}

	if new(big.Int).Sub(ethBlockNum, deposit.BlockNum.BigInt()).Uint64() < plasma.finalityBound {
		return nil, fmt.Errorf("deposit not finalized")
	}

	return &deposit, nil
}

// HasTXBeenExited indicates if the position has ever been exited
func (plasma *Plasma) HasTXBeenExited(position [4]*big.Int) bool {
	var key []byte
	var priority *big.Int
	if position[3].Sign() == 0 { // utxo exit
		txPos := [3]*big.Int{position[0], position[1], position[3]}
		priority = calcPriority(txPos)
		key = prefixKey(transactionExitPrefix, priority.Bytes())
	} else { // deposit exit
		priority = position[3]
		key = prefixKey(depositExitPrefix, priority.Bytes())
	}

	type exit struct {
		Amount       *big.Int
		CommittedFee *big.Int
		CreatedAt    *big.Int
		Owner        common.Address
		State        uint8
	}

	if !plasma.memdb.Contains(key) {
		var e exit
		var err error
		if position[3].Sign() == 0 {
			e, err = plasma.session.TxExits(priority)
		} else {
			e, err = plasma.session.DepositExits(priority)
		}

		// default to true if the contract cannot be queried. Nothing should be spent
		if err != nil {
			plasma.logger.Error(fmt.Sprintf("Error querying contract %s", err))
			return true
		}

		if e.State == 1 || e.State == 3 {
			return true
		} else {
			return false
		}
	}

	return true
}

func watchDeposits(plasma *Plasma) {
	// suscribe to future deposits
	deposits := make(chan *contracts.PlasmaMVPDeposit)
	opts := &bind.WatchOpts{
		Start:   nil, // latest block
		Context: context.Background(),
	}
	plasma.session.Contract.WatchDeposit(opts, deposits)

	for deposit := range deposits {
		key := prefixKey(depositPrefix, deposit.DepositNonce.Bytes())

		// remove the nonce, encode, and store
		data, err := json.Marshal(plasmaTypes.Deposit{
			Owner:    deposit.Depositor,
			Amount:   sdk.NewUintFromBigInt(deposit.Amount),
			BlockNum: sdk.NewUintFromBigInt(deposit.EthBlockNum),
		})

		if err != nil {
			plasma.logger.Error("Error encoding deposit event from contract -", deposit)
		} else {
			plasma.memdb.Put(key, data)
		}
	}
}

func watchExits(plasma *Plasma) {
	startedDepositExits := make(chan *contracts.PlasmaMVPStartedDepositExit)
	startedTransactionExits := make(chan *contracts.PlasmaMVPStartedTransactionExit)
	challengedExits := make(chan *contracts.PlasmaMVPChallengedExit)

	opts := &bind.WatchOpts{
		Start:   nil, // latest block
		Context: context.Background(),
	}
	plasma.session.Contract.WatchStartedDepositExit(opts, startedDepositExits)
	plasma.session.Contract.WatchStartedTransactionExit(opts, startedTransactionExits)
	plasma.session.Contract.WatchChallengedExit(opts, challengedExits)

	go func() {
		for depositExit := range startedDepositExits {
			nonce := depositExit.Nonce.Bytes()
			key := prefixKey(depositExitPrefix, nonce)
			plasma.memdb.Put(key, nil)
		}

		plasma.logger.Info("stopped watching for deposit exits")
	}()

	go func() {
		for transactionExit := range startedTransactionExits {
			priority := calcPriority(transactionExit.Position).Bytes()
			key := prefixKey(transactionExitPrefix, priority)
			plasma.memdb.Put(key, nil)
		}

		plasma.logger.Info("stopped watching for transaction exits")
	}()

	go func() {
		for challengedExit := range challengedExits {
			if challengedExit.Position[3].Sign() == 0 {
				position := [3]*big.Int{challengedExit.Position[0], challengedExit.Position[1], challengedExit.Position[2]}
				key := prefixKey(transactionExitPrefix, calcPriority(position).Bytes())
				plasma.memdb.Delete(key)
			} else {
				key := prefixKey(depositExitPrefix, challengedExit.Position[3].Bytes())
				plasma.memdb.Delete(key)
			}
		}

		plasma.logger.Info("stopped watching for challenged exit")
	}()
}

func watchEthBlocks(plasma *Plasma, ch <-chan *types.Header) {
	for header := range ch {
		plasma.lock.Lock()
		plasma.ethBlockNum = header.Number
		plasma.lock.Unlock()
	}

	plasma.logger.Info("Block subscription closed.")
}
