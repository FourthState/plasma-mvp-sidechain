package eth

import (
	"crypto/ecdsa"
	"fmt"
	contracts "github.com/FourthState/plasma-mvp-sidechain/contracts/wrappers"
	plasmaTypes "github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/tendermint/tendermint/libs/log"
	"math/big"
	"sync"
)

// Plasma holds related unexported members
type Plasma struct {
	*sync.Mutex
	operatorSession *contracts.PlasmaMVPSession
	contract        *contracts.PlasmaMVP

	client Client
	logger log.Logger

	blockNum      *big.Int
	ethBlockNum   *big.Int
	finalityBound uint64
}

// TODO: synching issues when rebooting a full node that contains plasma headers that have not been committed

// InitPlasma binds the go wrapper to the deployed contract. This private key provides authentication for the operator
func InitPlasma(contractAddr common.Address, client Client, finalityBound uint64, logger log.Logger, isOperator bool, operatorPrivKey *ecdsa.PrivateKey) (*Plasma, error) {
	logger.Info(fmt.Sprintf("binding to contract address 0x%x", contractAddr))
	plasmaContract, err := contracts.NewPlasmaMVP(contractAddr, client.ec)
	if err != nil {
		return nil, err
	}

	// Create a session with the contract and operator account
	var operatorSession *contracts.PlasmaMVPSession = nil
	if isOperator {
		auth := bind.NewKeyedTransactor(operatorPrivKey)
		operatorSession = &contracts.PlasmaMVPSession{
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
	}

	ethBlockNum, err := client.LatestBlockNum()
	if err != nil {
		return nil, fmt.Errorf("error retrieving latest ethereum block number { %s }", err)
	}

	lastCommittedBlock, err := plasmaContract.LastCommittedBlock(nil)
	if err != nil {
		return nil, fmt.Errorf("contract connection not correctly established { %s }", err)
	}

	plasma := &Plasma{
		operatorSession: operatorSession,
		contract:        plasmaContract,
		client:          client,
		logger:          logger,

		ethBlockNum:   ethBlockNum,
		blockNum:      lastCommittedBlock,
		finalityBound: finalityBound,

		Mutex: &sync.Mutex{},
	}

	// listen to new ethereum block headers
	ethCh, err := client.SubscribeToHeads()
	if err != nil {
		return nil, fmt.Errorf("could not successfully subscribe to heads: { %s }", err)
	}

	go watchEthBlocks(plasma, ethCh)

	return plasma, nil
}

func (plasma *Plasma) OperatorAddress() (common.Address, error) {
	return plasma.contract.Operator(nil)
}

// SubmitBlock proxy. TODO: handle batching with a timmer interrupt
func (plasma *Plasma) SubmitBlock(block plasmaTypes.Block) error {
	// only the contract operator can submit blocks
	if plasma.operatorSession == nil {
		return nil
	}

	plasma.blockNum = plasma.blockNum.Add(plasma.blockNum, big.NewInt(1))

	_, err := plasma.operatorSession.SubmitBlock(
		[][32]byte{block.Header},
		[]*big.Int{big.NewInt(int64(block.TxnCount))},
		[]*big.Int{block.FeeAmount},
		plasma.blockNum)

	return err
}

// GetDeposit checks the existence of a deposit nonce
func (plasma *Plasma) GetDeposit(nonce *big.Int) (plasmaTypes.Deposit, bool) {
	deposit, err := plasma.contract.Deposits(nil, nonce)
	if err != nil {
		// TODO: log the error
		return plasmaTypes.Deposit{}, false
	}

	if deposit.CreatedAt.Sign() == 0 {
		return plasmaTypes.Deposit{}, false
	}

	// check the finality bound
	plasma.Lock()
	ethBlockNum := plasma.ethBlockNum
	plasma.Unlock()
	if ethBlockNum.Sign() < 0 {
		plasma.logger.Error(fmt.Sprintf("failed to retreive information about deposit %s", nonce))
		return plasmaTypes.Deposit{}, false
	}

	if new(big.Int).Sub(ethBlockNum, deposit.EthBlockNum).Uint64() < plasma.finalityBound {
		return plasmaTypes.Deposit{}, false
	}

	return plasmaTypes.Deposit{
		Owner:       deposit.Owner,
		Amount:      deposit.Amount,
		EthBlockNum: deposit.EthBlockNum,
	}, true
}

// HasTXBeenExited indicates if the position has ever been exited
func (plasma *Plasma) HasTxBeenExited(position plasmaTypes.Position) bool {
	type exit struct {
		Amount       *big.Int
		CommittedFee *big.Int
		CreatedAt    *big.Int
		Owner        common.Address
		State        uint8
	}

	var (
		e   exit
		err error
	)

	priority := position.Priority()
	if position.IsDeposit() {
		e, err = plasma.contract.DepositExits(nil, priority)
	} else {
		e, err = plasma.contract.TxExits(nil, priority)
	}

	// censor spends until the error is fixed
	if err != nil {
		plasma.logger.Error(fmt.Sprintf("failed to retrieve exit information about position %s { %s }", position, err))
		return true
	}

	return e.State == 1 || e.State == 3
}

// keep track of the probabilistic latest ethereum state (reorgs) for finality purposes
func watchEthBlocks(plasma *Plasma, ch <-chan *types.Header) {
	plasma.logger.Info("listening to ethereum block headers")

	for header := range ch {
		plasma.Lock()
		plasma.ethBlockNum = header.Number
		plasma.Unlock()
	}

	// block header subsciprtion should never close unless the daemon is shut off
	plasma.logger.Error("etheruem block header subscription closed")
}
