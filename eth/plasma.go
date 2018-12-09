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
)

// Contains the binded wrapper and keys of the operator
type Plasma struct {
	session       *rootchain.RootChainSession
	client        *Client
	logger        log.Logger
	memdb         *memdb.DB
	db            *leveldb.DB
	blockNum      sdk.Uint
	ethBlockNum   *big.Int
	finalityBound uint64
}

// InitPlasma binds the go wrapper to the deployed contract. This private key provides authentication
// for the operator
func InitPlasma(contractAddr string, privateKey *ecdsa.PrivateKey, client *Client, logger log.Logger, isValidator bool, finalityBound uint64) (*Plasma, error) {
	plasmaContract, err := rootchain.NewRootChain(common.HexToAddress(contractAddr), client.ec)
	if err != nil {
		return nil, err
	}

	var plasmaSession *rootchain.RootChainSession
	if isValidator {
		// Create a session with the contract
		auth := bind.NewKeyedTransactor(privateKey)
		plasmaSession = &rootchain.RootChainSession{
			Contract: plasmaContract,
			CallOpts: bind.CallOpts{
				Pending: true,
			},
			TransactOpts: bind.TransactOpts{
				From:     auth.From,
				Nonce:    big.NewInt(3),
				Signer:   auth.Signer,
				GasLimit: 3141592, // aribitrary. TODO: check this
			},
		}

	} else {
		// Create a session with the contract
		plasmaSession = &rootchain.RootChainSession{
			Contract: plasmaContract,
			CallOpts: bind.CallOpts{
				Pending: true,
			},
			TransactOpts: bind.TransactOpts{
				Nonce:    big.NewInt(3),
				GasLimit: 3141592, // aribitrary. TODO: check this
			},
		}
	}

	plasma := &Plasma{
		session: plasmaSession,
		client:  client,
		// capacity argument is advisory and not enforced in the memdb implementation
		// TODO: flush the in-memory DB to a local one to bound memory consumption
		memdb:         memdb.New(comparer.DefaultComparer, 1),
		logger:        logger,
		blockNum:      sdk.ZeroUint(),
		finalityBound: finalityBound,
	}

	ethCh, err := plasma.client.SubscribeToHeads()
	if err != nil {
		plasma.logger.Error("Could not successfully subscribe to heads: %v", err)
		return nil, err
	}

	go trackEthBLocks(plasma, ethCh)
	go plasma.watchDeposits()
	go plasma.watchExits()

	return plasma, nil
}

func trackEthBLocks(plasma *Plasma, ch <-chan *types.Header) {
	for {
		header := <-ch
		plasma.ethBlockNum = header.Number
	}
}

// SubmitBlock proxy. TODO: handle batching with a timmer interrupt
func (plasma *Plasma) SubmitBlock(header []byte, numTxns sdk.Uint, fee sdk.Uint) (*types.Transaction, error) {
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
>>>>>>> 0069dc98085d16bf1e8f9d194defcaceddc9b1b8
	key := prefixKey(depositPrefix, nonce.BigInt().Bytes())
	data, err := plasma.memdb.Get(key)

	var deposit plasmaTypes.Deposit

	// if entry exists, only continue if we can decode successfully
	if err == nil {
		// try to decode and return
		err := rlp.DecodeBytes(data, &deposit)
		if err != nil {
			plasma.memdb.Delete(key)
			plasma.logger.Error("Error decoding cached deposit: %x", data)
		} else if new(big.Int).Sub(plasma.ethBlockNum, deposit.BlockNum.BigInt()).Uint64() >= plasma.finalityBound {
			return &deposit, nil
		} else {
			return nil, errors.New("deposit not finalized")
		}
	}

	// conduct a contract call if the deposit does not exist in the cache or decoding failed
	d, err := plasma.session.Deposits(nonce.BigInt())
	if err != nil {
		plasma.logger.Error("Contract call, GetDeposit, failed %v", err)
		return utxo.Deposit{}, err
	}

	// deposit does not existed if the timestamp is the default value
	if d.CreatedAt.Sign() == 0 {
		return nil, errors.New("deposit does not exist")
	}

	deposit = plasmaTypes.Deposit{
		Owner:    d.Owner,
		Amount:   sdk.NewUintFromBigInt(d.Amount),
		BlockNum: sdk.NewUintFromBigInt(d.EthBlocknum),
	}

	data, err = rlp.EncodeToBytes(deposit)
	if err != nil {
		plasma.logger.Error("Error encoding: %v. Will not be cached", deposit)
	} else { // cache only if we can encode successfully
		plasma.memdb.Put(key, data)
	}

	// check finality bound for the deposit

	if new(big.Int).Sub(plasma.ethBlockNum, d.EthBlocknum).Uint64() >= plasma.finalityBound {
		return &deposit, nil
	} else {
		return nil, errors.New("deposit not finalized")
	}
}

// CheckTransaction indicates if the position has every been exited
func (plasma *Plasma) HasTXBeenExited(position [4]sdk.Uint) (bool, error) {
	var key []byte
	if position[3].Sign() == 0 { // utxo exit
		pos := [3]*big.Int{position[0].BigInt(), position[1].BigInt(), position[3].BigInt()}
		priority := calcPriority(pos).Bytes()
		key = prefixKey(transactionExitPrefix, priority)
	} else { // deposit exit
		key = prefixKey(depositExitPrefix, position[3].BigInt().Bytes())
	}

	return plasma.memdb.Contains(key), nil
}

func (plasma *Plasma) watchDeposits() {
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
		val, err := rlp.EncodeToBytes(plasmaTypes.Deposit{
			Owner:    deposit.Depositor,
			Amount:   sdk.NewUintFromBigInt(deposit.Amount),
			BlockNum: sdk.NewUintFromBigInt(deposit.EthBlockNum),
		})

		if err != nil {
			plasma.logger.Error("Error encoding deposit event from contract: %v", deposit)
		} else {
			plasma.memdb.Put(key, val)
		}
	}
}

func (plasma *Plasma) watchExits() {
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
	}()

	go func() {
		for transactionExit := range startedTransactionExits {
			priority := calcPriority(transactionExit.Position).Bytes()
			key := prefixKey(transactionExitPrefix, priority)
			plasma.memdb.Put(key, nil)
		}
	}()
}
