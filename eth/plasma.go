package eth

import (
	"context"
	"crypto/ecdsa"
	"errors"
	rootchain "github.com/FourthState/plasma-mvp-sidechain/contracts/wrappers"
	plasmaTypes "github.com/FourthState/plasma-mvp-sidechain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/comparer"
	"github.com/syndtr/goleveldb/leveldb/memdb"
	"github.com/tendermint/tendermint/libs/log"
	"math/big"
	"sync"
)

// Plasma holds related unexported members
type Plasma struct {
	session *rootchain.RootChainSession
	client  *Client
	logger  log.Logger

	memdb *memdb.DB
	db    *leveldb.DB

	blockNum      sdk.Uint
	ethBlockNum   *big.Int
	finalityBound int64

	lock *sync.Mutex
}

type serializedDeposit struct {
	owner    [20]byte
	amount   []byte
	blocknum []byte
}

// InitPlasma binds the go wrapper to the deployed contract. This private key provides authentication for the operator
func InitPlasma(contractAddr string, privateKey *ecdsa.PrivateKey, client *Client, logger log.Logger, finalityBound int64) (*Plasma, error) {
	plasmaContract, err := rootchain.NewRootChain(common.HexToAddress(contractAddr), client.ec)
	if err != nil {
		return nil, err
	}

	// Create a session with the contract and operator account
	auth := bind.NewKeyedTransactor(privateKey)
	plasmaSession := &rootchain.RootChainSession{
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

	plasma := &Plasma{
		session: plasmaSession,
		client:  client,
		// capacity argument is advisory and not enforced in the memdb implementation
		// TODO: flush the in-memory DB to a local one to bound memory consumption
		memdb:  memdb.New(comparer.DefaultComparer, 1),
		logger: logger,

		ethBlockNum:   big.NewInt(0),
		blockNum:      sdk.ZeroUint(),
		finalityBound: finalityBound,

		lock: &sync.Mutex{},
	}

	ethCh, err := plasma.client.SubscribeToHeads()
	if err != nil {
		plasma.logger.Error("Could not successfully subscribe to heads: %v", err)
		return nil, err
	}

	// start listeners
	go watchEthBlocks(plasma, ethCh)
	go watchDeposits(plasma)
	go watchExits(plasma)

	return plasma, nil
}

// SubmitBlock proxy. TODO: handle batching with a timmer interrupt
func (plasma *Plasma) SubmitBlock(header []byte, numTxns sdk.Uint, fee sdk.Uint) (*types.Transaction, error) {
	plasma.blockNum = plasma.blockNum.Add(sdk.OneUint())

	tx, err := plasma.session.SubmitBlock(
		header,
		[]*big.Int{numTxns.BigInt()},
		[]*big.Int{fee.BigInt()},
		plasma.blockNum.BigInt())

	if err != nil {
		return nil, err
	}

	return tx, nil
}

// GetDeposit checks the existence of a deposit nonce
func (plasma *Plasma) GetDeposit(nonce sdk.Uint) (*plasmaTypes.Deposit, error) {
	key := prefixKey(depositPrefix, nonce.BigInt().Bytes())
	data, err := plasma.memdb.Get(key)

	var deposit plasmaTypes.Deposit
	// check against the contract if the deposit is not in the cache or decoding fail
	if err != nil && deserializeDeposit(data, &deposit) != nil {
		if plasma.memdb.Contains(key) {
			plasma.logger.Info("corrupted deposit found within db")
			plasma.memdb.Delete(key)
		}

		d, err := plasma.session.Deposits(nonce.BigInt())
		if err != nil {
			plasma.logger.Error("contract call, deposits, failed")
			return nil, err
		}

		if d.CreatedAt.Sign() == 0 {
			return nil, errors.New("deposit does not exist")
		}

		deposit = plasmaTypes.Deposit{
			Owner:    d.Owner,
			Amount:   sdk.NewUintFromBigInt(d.Amount),
			BlockNum: sdk.NewUintFromBigInt(d.EthBlocknum),
		}
	}

	// save to the db
	data, err = rlp.EncodeToBytes(serializedDeposit{
		owner:    deposit.Owner,
		amount:   deposit.Amount.BigInt().Bytes(),
		blocknum: deposit.BlockNum.BigInt().Bytes(),
	})
	if err != nil {
		plasma.logger.Error("error encoding deposit. will not be cached")
	} else {
		plasma.memdb.Put(key, data)
	}

	// check finality bound for the deposit
	plasma.lock.Lock()
	ethBlockNum := plasma.ethBlockNum
	plasma.lock.Unlock()
	if new(big.Int).Sub(ethBlockNum, deposit.BlockNum.BigInt()).Int64() >= plasma.finalityBound {
		return &deposit, nil
	} else {
		return nil, errors.New("deposit not finalized")
	}
}

// HasTXBeenExited indicates if the position has ever been exited
func (plasma *Plasma) HasTXBeenExited(position [4]sdk.Uint) bool {
	var key []byte
	if position[3].Sign() == 0 { // utxo exit
		pos := [3]*big.Int{position[0].BigInt(), position[1].BigInt(), position[3].BigInt()}
		priority := calcPriority(pos).Bytes()
		key = prefixKey(transactionExitPrefix, priority)
	} else { // deposit exit
		key = prefixKey(depositExitPrefix, position[3].BigInt().Bytes())
	}

	return plasma.memdb.Contains(key)
}

func watchDeposits(plasma *Plasma) {
	// suscribe to future deposits
	deposits := make(chan *rootchain.RootChainDeposit)
	opts := &bind.WatchOpts{
		Start:   nil, // latest block
		Context: context.Background(),
	}
	plasma.session.Contract.WatchDeposit(opts, deposits)

	for deposit := range deposits {
		key := prefixKey(depositPrefix, deposit.DepositNonce.Bytes())

		// remove the nonce, encode, and store
		val, err := rlp.EncodeToBytes(serializedDeposit{
			owner:    deposit.Depositor,
			amount:   deposit.Amount.Bytes(),
			blocknum: deposit.EthBlockNum.Bytes(),
		})

		if err != nil {
			plasma.logger.Error("Error encoding deposit event from contract -", deposit)
		} else {
			plasma.memdb.Put(key, val)
		}
	}
}

func watchExits(plasma *Plasma) {
	startedDepositExits := make(chan *rootchain.RootChainStartedDepositExit)
	startedTransactionExits := make(chan *rootchain.RootChainStartedTransactionExit)

	opts := &bind.WatchOpts{
		Start:   nil, // latest block
		Context: context.Background(),
	}
	plasma.session.Contract.WatchStartedDepositExit(opts, startedDepositExits)
	plasma.session.Contract.WatchStartedTransactionExit(opts, startedTransactionExits)

	go func() {
		for depositExit := range startedDepositExits {
			nonce := depositExit.Nonce.Bytes()
			key := prefixKey(depositExitPrefix, nonce)
			plasma.memdb.Put(key, nil)
		}

		plasma.logger.Info("stopped watching deposit exits")
	}()

	go func() {
		for transactionExit := range startedTransactionExits {
			priority := calcPriority(transactionExit.Position).Bytes()
			key := prefixKey(transactionExitPrefix, priority)
			plasma.memdb.Put(key, nil)
		}

		plasma.logger.Info("stopped watching transaction exits")
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

func deserializeDeposit(data []byte, deposit *plasmaTypes.Deposit) error {
	var dep serializedDeposit
	if err := rlp.DecodeBytes(data, &dep); err != nil {
		return err
	}

	deposit.Owner = dep.owner
	deposit.Amount = sdk.NewUintFromBigInt(new(big.Int).SetBytes(dep.amount))
	deposit.Amount = sdk.NewUintFromBigInt(new(big.Int).SetBytes(dep.amount))
	deposit.BlockNum = sdk.NewUintFromBigInt(new(big.Int).SetBytes(dep.blocknum))

	return nil
}
