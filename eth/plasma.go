package eth

import (
	"bytes"
	"crypto/ecdsa"
	"fmt"
	contracts "github.com/FourthState/plasma-mvp-sidechain/contracts/wrappers"
	plasmaTypes "github.com/FourthState/plasma-mvp-sidechain/plasma"
	"github.com/FourthState/plasma-mvp-sidechain/store"
	"github.com/FourthState/plasma-mvp-sidechain/utils"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tendermint/tendermint/libs/log"
	"math/big"
	"time"
)

var logger log.Logger = log.NewNopLogger()

// SetLogger will adjust the logger used in the eth module
func SetLogger(l log.Logger) {
	logger = l
}

// Plasma holds related unexported members
type Plasma struct {
	*contracts.PlasmaMVP // expose all the contract methods

	client          Client
	finalityBound   uint64
	operatorSession *operatorSession
}

type operatorSession struct {
	*contracts.PlasmaMVPSession

	commitmentRate      time.Duration
	lastBlockSubmission time.Time
}

// InitPlasma binds the go wrapper to the deployed contract. This private key provides authentication for the operator
func InitPlasma(contractAddr common.Address, client Client, finalityBound uint64) (*Plasma, error) {
	logger.Info(fmt.Sprintf("binding to contract address 0x%x", contractAddr))
	plasmaContract, err := contracts.NewPlasmaMVP(contractAddr, client.ec)
	if err != nil {
		return nil, err
	}

	plasma := &Plasma{
		PlasmaMVP:     plasmaContract,
		client:        client,
		finalityBound: finalityBound,
	}

	return plasma, nil
}

// WithOperatorSession will set up an operators session with the smart contract. The contract's operator public key must
// match the public key corresponding `operatorPrivKey`
func (plasma *Plasma) WithOperatorSession(operatorPrivkey *ecdsa.PrivateKey, commitmentRate time.Duration) (*Plasma, error) {
	logger.Info(fmt.Sprintf("block commitment rate set to %s", commitmentRate))

	// check that the public key matches the address of the operator
	addr := crypto.PubkeyToAddress(operatorPrivkey.PublicKey)
	operator, err := plasma.OperatorAddress()
	if err != nil {
		return plasma, err
	}
	if !bytes.Equal(operator[:], addr[:]) {
		return plasma, fmt.Errorf("operator address mismatch. Got 0x%x. Expected:0x%x", addr, operator)
	}

	auth := bind.NewKeyedTransactor(operatorPrivkey)
	contractSession := &contracts.PlasmaMVPSession{
		Contract: plasma.PlasmaMVP,
		CallOpts: bind.CallOpts{
			Pending: true,
		},
		TransactOpts: bind.TransactOpts{
			From:     auth.From,
			Signer:   auth.Signer,
			GasLimit: 3141592, // aribitrary. TODO: check this
		},
	}

	opSession := &operatorSession{
		PlasmaMVPSession:    contractSession,
		commitmentRate:      commitmentRate,
		lastBlockSubmission: time.Now(),
	}

	plasma.operatorSession = opSession
	return plasma, nil
}

// OperatorAddress will fetch the plasma operator address from the connected smart contract
func (plasma *Plasma) OperatorAddress() (common.Address, error) {
	return plasma.Operator(nil)
}

// CommitPlasmaHeaders will commit all new non-committed headers to the smart contract.
// the commitmentRate interval must pass since the last commitment
func (plasma *Plasma) CommitPlasmaHeaders(ctx sdk.Context, plasmaStore store.PlasmaStore) error {
	// only the contract operator can submit blocks. The commitment duration must also pass
	if plasma.operatorSession == nil || time.Since(plasma.operatorSession.lastBlockSubmission).Seconds() < plasma.operatorSession.commitmentRate.Seconds() {
		return nil
	}

	logger.Info("attempting to commit plasma headers...")

	lastCommittedBlock, err := plasma.LastCommittedBlock(nil)
	if err != nil {
		logger.Error("error retrieving the last committed block number")
		return err
	}

	firstBlockNum := new(big.Int).Add(lastCommittedBlock, utils.Big1)
	blockNum := lastCommittedBlock.Add(lastCommittedBlock, utils.Big1)

	var (
		headers      [][32]byte
		txnsPerBlock []*big.Int
		feesPerBlock []*big.Int
	)

	block, ok := plasmaStore.GetBlock(ctx, blockNum)
	if !ok { // no blocks to submit
		logger.Info("no plasma blocks to commit")
		return nil
	}

	for ok {
		headers = append(headers, block.Header)
		txnsPerBlock = append(txnsPerBlock, big.NewInt(int64(block.TxnCount)))
		feesPerBlock = append(feesPerBlock, block.FeeAmount)

		blockNum = blockNum.Add(blockNum, utils.Big1)
		block, ok = plasmaStore.GetBlock(ctx, blockNum)
	}

	logger.Info(fmt.Sprintf("committing %d plasma blocks. first block num: %s", len(headers), firstBlockNum))
	plasma.operatorSession.lastBlockSubmission = time.Now()
	_, err = plasma.operatorSession.SubmitBlock(headers, txnsPerBlock, feesPerBlock, firstBlockNum)
	if err != nil {
		logger.Error(fmt.Sprintf("error committing headers { %s }", err))
		return err
	}

	return err
}

// GetDeposit checks the existence of a deposit nonce. The state is synchronized with the provided `plasmaBlockHeight. The deposit
// must have occured before or at the same pegged ethereum block as `plasmaBlockHeight`.
func (plasma *Plasma) GetDeposit(plasmaBlockHeight *big.Int, nonce *big.Int) (plasmaTypes.Deposit, *big.Int, bool) {
	deposit, err := plasma.Deposits(nil, nonce)
	if err != nil {
		logger.Error(fmt.Sprintf("failed deposit retrieval: %s", err))
		return plasmaTypes.Deposit{}, nil, false
	}

	if deposit.CreatedAt.Sign() == 0 {
		return plasmaTypes.Deposit{}, nil, false
	}

	// check the finality bound based off pegged ETH block
	ethBlockNum, err := plasma.ethBlockPeg(plasmaBlockHeight)
	if err != nil {
		logger.Error(fmt.Sprintf("could not get pegged ETH Block for sidechain block %s: %s", plasmaBlockHeight, err))
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

// HasTxExited indicates if the position has ever been exited at a time less than or equal to
// the time `plasmaBlockHeight` was submitted. If nil, it is checked against the latest state
func (plasma *Plasma) HasTxExited(plasmaBlockHeight *big.Int, position plasmaTypes.Position) (bool, error) {
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
		e, err = plasma.DepositExits(nil, priority)
	} else {
		e, err = plasma.TxExits(nil, priority)
	}

	// censor spends until the error is fixed
	if err != nil {
		logger.Error(fmt.Sprintf("failed to retrieve exit information about position %s { %s }", position, err))
		return true, err
	}

	// `Pending` or `Challenged` stateee
	exited := e.State == 1 || e.State == 3

	// synchronize with the correct ethereum state
	if plasmaBlockHeight != nil {
		ethBlockNum, err := plasma.ethBlockPeg(plasmaBlockHeight)
		if err != nil {
			// censore spends until the error is fixed
			logger.Error(fmt.Sprintf("could not get associated ETH Block for plasma block %s: %s", plasmaBlockHeight, err))
			return true, err
		}

		// exited AND the exit occured before or in the pegged ethereum block
		return exited && e.EthBlockNum.Cmp(ethBlockNum) <= 0, nil
	}

	return exited, nil
}

// Return the Ethereum Block that sidechain should use to synchronize current block tx's with rootchain state
func (plasma *Plasma) ethBlockPeg(plasmaBlockHeight *big.Int) (*big.Int, error) {
	lastCommittedBlock, err := plasma.LastCommittedBlock(nil)
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
	prevBlock := new(big.Int).Sub(plasmaBlockHeight, big.NewInt(1))
	// For syncing nodes, peg to EthBlock at plasmaBlock-1 submission
	// For live nodes, peg to LastCommittedBlock
	if lastCommittedBlock.Cmp(prevBlock) == 1 {
		blockIndex = prevBlock
	} else {
		blockIndex = lastCommittedBlock
	}
	submittedBlock, err := plasma.PlasmaChain(nil, blockIndex)
	if err != nil {
		return nil, err
	}
	return submittedBlock.EthBlockNum, nil
}
