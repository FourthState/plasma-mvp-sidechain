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

// InitPlasma binds the go wrapper to the deployed contract. This private key provides authentication for the operator
func InitPlasma(contractAddr common.Address, client Client, finalityBound uint64, logger log.Logger, isOperator bool, operatorPrivKey *ecdsa.PrivateKey) (*Plasma, error) {
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

	// TODO: deal with syncing issues
	lastCommittedBlock, err := plasmaContract.LastCommittedBlock(nil)
	if err != nil {
		return nil, fmt.Errorf("Contract connection not correctly established - %s", err)
	}

	plasma := &Plasma{
		operatorSession: operatorSession,
		contract:        plasmaContract,
		client:          client,
		logger:          logger,

		ethBlockNum:   big.NewInt(-1),
		blockNum:      lastCommittedBlock,
		finalityBound: finalityBound,

		Mutex: &sync.Mutex{},
	}

	// listen to new ethereum block headers
	ethCh, err := client.SubscribeToHeads()
	if err != nil {
		logger.Error("Could not successfully subscribe to heads: %v", err)
		return nil, err
	}

	go watchEthBlocks(plasma, ethCh)

	return plasma, nil
}

// SubmitBlock proxy. TODO: handle batching with a timmer interrupt
func (plasma *Plasma) SubmitBlock(block plasmaTypes.Block) error {
	plasma.blockNum = plasma.blockNum.Add(plasma.blockNum, big.NewInt(1))

	_, err := plasma.operatorSession.SubmitBlock(
		[][32]byte{block.Header},
		[]*big.Int{big.NewInt(int64(block.TxnCount))},
		[]*big.Int{block.TotalFee},
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
		plasma.logger.Error("failed `GetDeposit`. not subscribed to ethereum headers")
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

	if err != nil {
		// TODO: log the error
		return true // censor spends until resolved
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

	plasma.logger.Info("etheruem block header subscription closed")
}
