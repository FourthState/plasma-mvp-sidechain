package eth

import (
	"crypto/ecdsa"
	"fmt"
	contracts "github.com/FourthState/plasma-mvp-sidechain/contracts/wrappers"
	plasmaTypes "github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tendermint/tendermint/libs/log"
	"math/big"
	"sync"
	"time"
)

// Plasma holds related unexported members
type Plasma struct {
	*sync.Mutex
	operatorSession *contracts.PlasmaMVPSession
	contract        *contracts.PlasmaMVP

	commitmentRate      time.Duration
	lastBlockSubmission time.Time

	client Client
	logger log.Logger

	finalityBound uint64
}

// TODO: synching issues when rebooting a full node that contains plasma headers that have not been committed

// InitPlasma binds the go wrapper to the deployed contract. This private key provides authentication for the operator
func InitPlasma(contractAddr common.Address, client Client, finalityBound uint64, commitmentRate time.Duration, logger log.Logger, isOperator bool, operatorPrivKey *ecdsa.PrivateKey) (*Plasma, error) {
	logger.Info(fmt.Sprintf("binding to contract address 0x%x", contractAddr))
	logger.Info(fmt.Sprintf("block commitment rate set to %s", commitmentRate))
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

		logger.Info("operator mode. authenticated contract session started")
	}

	plasma := &Plasma{
		operatorSession: operatorSession,
		contract:        plasmaContract,
		client:          client,
		logger:          logger,

		commitmentRate:      commitmentRate,
		lastBlockSubmission: time.Now(),

		finalityBound: finalityBound,

		Mutex: &sync.Mutex{},
	}

	return plasma, nil
}

func (plasma *Plasma) OperatorAddress() (common.Address, error) {
	return plasma.contract.Operator(nil)
}

// CommitPlasmaHeaders will commit all new non-committed headers to the smart contract.
// the commitmentRate interval must pass since the last commitment
func (plasma *Plasma) CommitPlasmaHeaders(ctx sdk.Context, blockStore store.BlockStore) error {
	// only the contract operator can submit blocks. The commitment duration must also pass
	if plasma.operatorSession == nil || time.Since(plasma.lastBlockSubmission).Seconds() < plasma.commitmentRate.Seconds() {
		return nil
	}

	plasma.logger.Info("attempting to commit plasma headers...")

	lastCommittedBlock, err := plasma.contract.LastCommittedBlock(nil)
	if err != nil {
		plasma.logger.Error("error retrieving the last committed block number")
		return err
	}

	firstBlockNum := new(big.Int).Add(lastCommittedBlock, utils.Big1)
	blockNum := lastCommittedBlock.Add(lastCommittedBlock, utils.Big1)

	var (
		headers      [][32]byte
		txnsPerBlock []*big.Int
		feesPerBlock []*big.Int
	)

	block, ok := blockStore.GetBlock(ctx, blockNum)
	if !ok { // no blocks to submit
		plasma.logger.Info("no plasma blocks to commit")
		return nil
	}

	for ok {
		headers = append(headers, block.Header)
		txnsPerBlock = append(txnsPerBlock, big.NewInt(int64(block.TxnCount)))
		feesPerBlock = append(feesPerBlock, block.FeeAmount)

		blockNum = blockNum.Add(blockNum, utils.Big1)
		block, ok = blockStore.GetBlock(ctx, blockNum)
	}

	plasma.logger.Info(fmt.Sprintf("committing %d plasma blocks. first block num: %s", len(headers), firstBlockNum))
	plasma.lastBlockSubmission = time.Now()
	_, err = plasma.operatorSession.SubmitBlock(headers, txnsPerBlock, feesPerBlock, firstBlockNum)
	if err != nil {
		plasma.logger.Error(fmt.Sprintf("error committing headers { %s }", err))
		return err
	}

	return err
}

// GetDeposit checks the existence of a deposit nonce
func (plasma *Plasma) GetDeposit(plasmaBlock *big.Int, nonce *big.Int) (plasmaTypes.Deposit, *big.Int, bool) {
	deposit, err := plasma.contract.Deposits(nil, nonce)
	if err != nil {
		// TODO: log the error
		return plasmaTypes.Deposit{}, nil, false
	}

	if deposit.CreatedAt.Sign() == 0 {
		return plasmaTypes.Deposit{}, nil, false
	}

	// check the finality bound based off pegged ETH block
	ethBlockNum, err := plasma.ethBlockPeg(plasmaBlock)
	if err != nil {
		plasma.logger.Error(fmt.Sprintf("could not get pegged ETH Block for sidechain block %d: %s", plasmaBlock.Int64(), err.Error()))
		return plasmaTypes.Deposit{}, nil, false
	}

	// how many blocks have occurred since deposit.
	// Note: Since pegged ETH block num could be before deposit's EthBlockNum, interval may be negative
	interval := new(big.Int).Sub(ethBlockNum, deposit.EthBlockNum)
	// how many more blocks need to get added for deposit to be considered final
	// Note: If deposit is finalized, threshold can be 0 or negative
	threshold := new(big.Int).Sub(big.NewInt(int64(plasma.finalityBound)), interval)
	if threshold.Sign() > 0 {
		return plasmaTypes.Deposit{}, threshold, false
	}

	return plasmaTypes.Deposit{
		Owner:       deposit.Owner,
		Amount:      deposit.Amount,
		EthBlockNum: deposit.EthBlockNum,
	}, threshold, true
}

// HasTXBeenExited indicates if the position has ever been exited
func (plasma *Plasma) HasTxBeenExited(plasmaBlock *big.Int, position plasmaTypes.Position) bool {
	type exit struct {
		Amount       *big.Int
		CommittedFee *big.Int
		CreatedAt    *big.Int
		EthBlockNum  *big.Int
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

	ethBlockNum, err := plasma.ethBlockPeg(plasmaBlock)
	if err != nil {
		plasma.logger.Error(fmt.Sprintf("could not get associated ETH Block for plasma block %d: %s", plasmaBlock.Int64(), err.Error()))
		return true
	}

	// Return true if exit is pending/finalized AND exit happened before the pegged ETH block
	return (e.State == 1 || e.State == 3) && e.EthBlockNum.Cmp(ethBlockNum) <= 0
}

// Return the Ethereum Block that sidechain should use to synchronize current block tx's with rootchain state
func (plasma *Plasma) ethBlockPeg(plasmaBlock *big.Int) (*big.Int, error) {
	lastCommittedBlock, err := plasma.contract.LastCommittedBlock(nil)
	if err != nil {
		return nil, err
	}
	// If no blocks submitted, use latestBlock as peg
	if lastCommittedBlock.Sign() == 0 {
		latestBlock, err := plasma.client.LatestBlockNum()
		if err != nil {
			return nil, err
		}
		return latestBlock, nil
	}
	var blockIndex *big.Int
	prevBlock := new(big.Int).Sub(plasmaBlock, big.NewInt(1))
	// For syncing nodes, peg to EthBlock at plasmaBlock-1 submission
	// For live nodes, peg to LastCommittedBlock
	if lastCommittedBlock.Cmp(prevBlock) == 1 {
		blockIndex = prevBlock
	} else {
		blockIndex = lastCommittedBlock
	}
	submittedBlock, err := plasma.contract.PlasmaChain(nil, blockIndex)
	if err != nil {
		return nil, err
	}
	return submittedBlock.EthBlockNum, nil
}
