// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package wrappers

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// PlasmaMVPABI is the input ABI used to generate the binding from.
const PlasmaMVPABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"balances\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"lastCommittedBlock\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"plasmaChain\",\"outputs\":[{\"name\":\"header\",\"type\":\"bytes32\"},{\"name\":\"numTxns\",\"type\":\"uint256\"},{\"name\":\"feeAmount\",\"type\":\"uint256\"},{\"name\":\"createdAt\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"operator\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"depositExitQueue\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"txExits\",\"outputs\":[{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"committedFee\",\"type\":\"uint256\"},{\"name\":\"createdAt\",\"type\":\"uint256\"},{\"name\":\"owner\",\"type\":\"address\"},{\"name\":\"state\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"txExitQueue\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"deposits\",\"outputs\":[{\"name\":\"owner\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"createdAt\",\"type\":\"uint256\"},{\"name\":\"ethBlockNum\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalWithdrawBalance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"depositExits\",\"outputs\":[{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"committedFee\",\"type\":\"uint256\"},{\"name\":\"createdAt\",\"type\":\"uint256\"},{\"name\":\"owner\",\"type\":\"address\"},{\"name\":\"state\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minExitBond\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"depositNonce\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"oldOperator\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"newOperator\",\"type\":\"address\"}],\"name\":\"ChangedOperator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"AddedToBalances\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"header\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"numTxns\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"feeAmount\",\"type\":\"uint256\"}],\"name\":\"BlockSubmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"depositor\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"depositNonce\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"ethBlockNum\",\"type\":\"uint256\"}],\"name\":\"Deposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"position\",\"type\":\"uint256[3]\"},{\"indexed\":false,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"confirmSignatures\",\"type\":\"bytes\"},{\"indexed\":false,\"name\":\"committedFee\",\"type\":\"uint256\"}],\"name\":\"StartedTransactionExit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"nonce\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"committedFee\",\"type\":\"uint256\"}],\"name\":\"StartedDepositExit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"position\",\"type\":\"uint256[4]\"},{\"indexed\":false,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"ChallengedExit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"position\",\"type\":\"uint256[4]\"},{\"indexed\":false,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"FinalizedExit\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOperator\",\"type\":\"address\"}],\"name\":\"changeOperator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"headers\",\"type\":\"bytes32[]\"},{\"name\":\"txnsPerBlock\",\"type\":\"uint256[]\"},{\"name\":\"feePerBlock\",\"type\":\"uint256[]\"},{\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"submitBlock\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"deposit\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"nonce\",\"type\":\"uint256\"},{\"name\":\"committedFee\",\"type\":\"uint256\"}],\"name\":\"startDepositExit\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"txPos\",\"type\":\"uint256[3]\"},{\"name\":\"txBytes\",\"type\":\"bytes\"},{\"name\":\"proof\",\"type\":\"bytes\"},{\"name\":\"confirmSignatures\",\"type\":\"bytes\"},{\"name\":\"committedFee\",\"type\":\"uint256\"}],\"name\":\"startTransactionExit\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"name\":\"committedFee\",\"type\":\"uint256\"}],\"name\":\"startFeeExit\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"exitingTxPos\",\"type\":\"uint256[4]\"},{\"name\":\"challengingTxPos\",\"type\":\"uint256[2]\"},{\"name\":\"txBytes\",\"type\":\"bytes\"},{\"name\":\"proof\",\"type\":\"bytes\"},{\"name\":\"confirmSignature\",\"type\":\"bytes\"}],\"name\":\"challengeExit\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"finalizeDepositExits\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"finalizeTransactionExits\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"withdraw\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"plasmaChainBalance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"txQueueLength\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"depositQueueLength\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// PlasmaMVP is an auto generated Go binding around an Ethereum contract.
type PlasmaMVP struct {
	PlasmaMVPCaller     // Read-only binding to the contract
	PlasmaMVPTransactor // Write-only binding to the contract
	PlasmaMVPFilterer   // Log filterer for contract events
}

// PlasmaMVPCaller is an auto generated read-only Go binding around an Ethereum contract.
type PlasmaMVPCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PlasmaMVPTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PlasmaMVPTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PlasmaMVPFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PlasmaMVPFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PlasmaMVPSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PlasmaMVPSession struct {
	Contract     *PlasmaMVP        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PlasmaMVPCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PlasmaMVPCallerSession struct {
	Contract *PlasmaMVPCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// PlasmaMVPTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PlasmaMVPTransactorSession struct {
	Contract     *PlasmaMVPTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// PlasmaMVPRaw is an auto generated low-level Go binding around an Ethereum contract.
type PlasmaMVPRaw struct {
	Contract *PlasmaMVP // Generic contract binding to access the raw methods on
}

// PlasmaMVPCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PlasmaMVPCallerRaw struct {
	Contract *PlasmaMVPCaller // Generic read-only contract binding to access the raw methods on
}

// PlasmaMVPTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PlasmaMVPTransactorRaw struct {
	Contract *PlasmaMVPTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPlasmaMVP creates a new instance of PlasmaMVP, bound to a specific deployed contract.
func NewPlasmaMVP(address common.Address, backend bind.ContractBackend) (*PlasmaMVP, error) {
	contract, err := bindPlasmaMVP(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PlasmaMVP{PlasmaMVPCaller: PlasmaMVPCaller{contract: contract}, PlasmaMVPTransactor: PlasmaMVPTransactor{contract: contract}, PlasmaMVPFilterer: PlasmaMVPFilterer{contract: contract}}, nil
}

// NewPlasmaMVPCaller creates a new read-only instance of PlasmaMVP, bound to a specific deployed contract.
func NewPlasmaMVPCaller(address common.Address, caller bind.ContractCaller) (*PlasmaMVPCaller, error) {
	contract, err := bindPlasmaMVP(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PlasmaMVPCaller{contract: contract}, nil
}

// NewPlasmaMVPTransactor creates a new write-only instance of PlasmaMVP, bound to a specific deployed contract.
func NewPlasmaMVPTransactor(address common.Address, transactor bind.ContractTransactor) (*PlasmaMVPTransactor, error) {
	contract, err := bindPlasmaMVP(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PlasmaMVPTransactor{contract: contract}, nil
}

// NewPlasmaMVPFilterer creates a new log filterer instance of PlasmaMVP, bound to a specific deployed contract.
func NewPlasmaMVPFilterer(address common.Address, filterer bind.ContractFilterer) (*PlasmaMVPFilterer, error) {
	contract, err := bindPlasmaMVP(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PlasmaMVPFilterer{contract: contract}, nil
}

// bindPlasmaMVP binds a generic wrapper to an already deployed contract.
func bindPlasmaMVP(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(PlasmaMVPABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PlasmaMVP *PlasmaMVPRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _PlasmaMVP.Contract.PlasmaMVPCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PlasmaMVP *PlasmaMVPRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PlasmaMVP.Contract.PlasmaMVPTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PlasmaMVP *PlasmaMVPRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PlasmaMVP.Contract.PlasmaMVPTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PlasmaMVP *PlasmaMVPCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _PlasmaMVP.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PlasmaMVP *PlasmaMVPTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PlasmaMVP.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PlasmaMVP *PlasmaMVPTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PlasmaMVP.Contract.contract.Transact(opts, method, params...)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _address) constant returns(uint256)
func (_PlasmaMVP *PlasmaMVPCaller) BalanceOf(opts *bind.CallOpts, _address common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _PlasmaMVP.contract.Call(opts, out, "balanceOf", _address)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _address) constant returns(uint256)
func (_PlasmaMVP *PlasmaMVPSession) BalanceOf(_address common.Address) (*big.Int, error) {
	return _PlasmaMVP.Contract.BalanceOf(&_PlasmaMVP.CallOpts, _address)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _address) constant returns(uint256)
func (_PlasmaMVP *PlasmaMVPCallerSession) BalanceOf(_address common.Address) (*big.Int, error) {
	return _PlasmaMVP.Contract.BalanceOf(&_PlasmaMVP.CallOpts, _address)
}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances(address ) constant returns(uint256)
func (_PlasmaMVP *PlasmaMVPCaller) Balances(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _PlasmaMVP.contract.Call(opts, out, "balances", arg0)
	return *ret0, err
}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances(address ) constant returns(uint256)
func (_PlasmaMVP *PlasmaMVPSession) Balances(arg0 common.Address) (*big.Int, error) {
	return _PlasmaMVP.Contract.Balances(&_PlasmaMVP.CallOpts, arg0)
}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances(address ) constant returns(uint256)
func (_PlasmaMVP *PlasmaMVPCallerSession) Balances(arg0 common.Address) (*big.Int, error) {
	return _PlasmaMVP.Contract.Balances(&_PlasmaMVP.CallOpts, arg0)
}

// DepositExitQueue is a free data retrieval call binding the contract method 0x5b3081d7.
//
// Solidity: function depositExitQueue(uint256 ) constant returns(uint256)
func (_PlasmaMVP *PlasmaMVPCaller) DepositExitQueue(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _PlasmaMVP.contract.Call(opts, out, "depositExitQueue", arg0)
	return *ret0, err
}

// DepositExitQueue is a free data retrieval call binding the contract method 0x5b3081d7.
//
// Solidity: function depositExitQueue(uint256 ) constant returns(uint256)
func (_PlasmaMVP *PlasmaMVPSession) DepositExitQueue(arg0 *big.Int) (*big.Int, error) {
	return _PlasmaMVP.Contract.DepositExitQueue(&_PlasmaMVP.CallOpts, arg0)
}

// DepositExitQueue is a free data retrieval call binding the contract method 0x5b3081d7.
//
// Solidity: function depositExitQueue(uint256 ) constant returns(uint256)
func (_PlasmaMVP *PlasmaMVPCallerSession) DepositExitQueue(arg0 *big.Int) (*big.Int, error) {
	return _PlasmaMVP.Contract.DepositExitQueue(&_PlasmaMVP.CallOpts, arg0)
}

// DepositExits is a free data retrieval call binding the contract method 0xce84f906.
//
// Solidity: function depositExits(uint256 ) constant returns(uint256 amount, uint256 committedFee, uint256 createdAt, address owner, uint8 state)
func (_PlasmaMVP *PlasmaMVPCaller) DepositExits(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Amount       *big.Int
	CommittedFee *big.Int
	CreatedAt    *big.Int
	Owner        common.Address
	State        uint8
}, error) {
	ret := new(struct {
		Amount       *big.Int
		CommittedFee *big.Int
		CreatedAt    *big.Int
		Owner        common.Address
		State        uint8
	})
	out := ret
	err := _PlasmaMVP.contract.Call(opts, out, "depositExits", arg0)
	return *ret, err
}

// DepositExits is a free data retrieval call binding the contract method 0xce84f906.
//
// Solidity: function depositExits(uint256 ) constant returns(uint256 amount, uint256 committedFee, uint256 createdAt, address owner, uint8 state)
func (_PlasmaMVP *PlasmaMVPSession) DepositExits(arg0 *big.Int) (struct {
	Amount       *big.Int
	CommittedFee *big.Int
	CreatedAt    *big.Int
	Owner        common.Address
	State        uint8
}, error) {
	return _PlasmaMVP.Contract.DepositExits(&_PlasmaMVP.CallOpts, arg0)
}

// DepositExits is a free data retrieval call binding the contract method 0xce84f906.
//
// Solidity: function depositExits(uint256 ) constant returns(uint256 amount, uint256 committedFee, uint256 createdAt, address owner, uint8 state)
func (_PlasmaMVP *PlasmaMVPCallerSession) DepositExits(arg0 *big.Int) (struct {
	Amount       *big.Int
	CommittedFee *big.Int
	CreatedAt    *big.Int
	Owner        common.Address
	State        uint8
}, error) {
	return _PlasmaMVP.Contract.DepositExits(&_PlasmaMVP.CallOpts, arg0)
}

// DepositNonce is a free data retrieval call binding the contract method 0xde35f5cb.
//
// Solidity: function depositNonce() constant returns(uint256)
func (_PlasmaMVP *PlasmaMVPCaller) DepositNonce(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _PlasmaMVP.contract.Call(opts, out, "depositNonce")
	return *ret0, err
}

// DepositNonce is a free data retrieval call binding the contract method 0xde35f5cb.
//
// Solidity: function depositNonce() constant returns(uint256)
func (_PlasmaMVP *PlasmaMVPSession) DepositNonce() (*big.Int, error) {
	return _PlasmaMVP.Contract.DepositNonce(&_PlasmaMVP.CallOpts)
}

// DepositNonce is a free data retrieval call binding the contract method 0xde35f5cb.
//
// Solidity: function depositNonce() constant returns(uint256)
func (_PlasmaMVP *PlasmaMVPCallerSession) DepositNonce() (*big.Int, error) {
	return _PlasmaMVP.Contract.DepositNonce(&_PlasmaMVP.CallOpts)
}

// DepositQueueLength is a free data retrieval call binding the contract method 0xca11dd97.
//
// Solidity: function depositQueueLength() constant returns(uint256)
func (_PlasmaMVP *PlasmaMVPCaller) DepositQueueLength(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _PlasmaMVP.contract.Call(opts, out, "depositQueueLength")
	return *ret0, err
}

// DepositQueueLength is a free data retrieval call binding the contract method 0xca11dd97.
//
// Solidity: function depositQueueLength() constant returns(uint256)
func (_PlasmaMVP *PlasmaMVPSession) DepositQueueLength() (*big.Int, error) {
	return _PlasmaMVP.Contract.DepositQueueLength(&_PlasmaMVP.CallOpts)
}

// DepositQueueLength is a free data retrieval call binding the contract method 0xca11dd97.
//
// Solidity: function depositQueueLength() constant returns(uint256)
func (_PlasmaMVP *PlasmaMVPCallerSession) DepositQueueLength() (*big.Int, error) {
	return _PlasmaMVP.Contract.DepositQueueLength(&_PlasmaMVP.CallOpts)
}

// Deposits is a free data retrieval call binding the contract method 0xb02c43d0.
//
// Solidity: function deposits(uint256 ) constant returns(address owner, uint256 amount, uint256 createdAt, uint256 ethBlockNum)
func (_PlasmaMVP *PlasmaMVPCaller) Deposits(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Owner       common.Address
	Amount      *big.Int
	CreatedAt   *big.Int
	EthBlockNum *big.Int
}, error) {
	ret := new(struct {
		Owner       common.Address
		Amount      *big.Int
		CreatedAt   *big.Int
		EthBlockNum *big.Int
	})
	out := ret
	err := _PlasmaMVP.contract.Call(opts, out, "deposits", arg0)
	return *ret, err
}

// Deposits is a free data retrieval call binding the contract method 0xb02c43d0.
//
// Solidity: function deposits(uint256 ) constant returns(address owner, uint256 amount, uint256 createdAt, uint256 ethBlockNum)
func (_PlasmaMVP *PlasmaMVPSession) Deposits(arg0 *big.Int) (struct {
	Owner       common.Address
	Amount      *big.Int
	CreatedAt   *big.Int
	EthBlockNum *big.Int
}, error) {
	return _PlasmaMVP.Contract.Deposits(&_PlasmaMVP.CallOpts, arg0)
}

// Deposits is a free data retrieval call binding the contract method 0xb02c43d0.
//
// Solidity: function deposits(uint256 ) constant returns(address owner, uint256 amount, uint256 createdAt, uint256 ethBlockNum)
func (_PlasmaMVP *PlasmaMVPCallerSession) Deposits(arg0 *big.Int) (struct {
	Owner       common.Address
	Amount      *big.Int
	CreatedAt   *big.Int
	EthBlockNum *big.Int
}, error) {
	return _PlasmaMVP.Contract.Deposits(&_PlasmaMVP.CallOpts, arg0)
}

// LastCommittedBlock is a free data retrieval call binding the contract method 0x3acb097a.
//
// Solidity: function lastCommittedBlock() constant returns(uint256)
func (_PlasmaMVP *PlasmaMVPCaller) LastCommittedBlock(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _PlasmaMVP.contract.Call(opts, out, "lastCommittedBlock")
	return *ret0, err
}

// LastCommittedBlock is a free data retrieval call binding the contract method 0x3acb097a.
//
// Solidity: function lastCommittedBlock() constant returns(uint256)
func (_PlasmaMVP *PlasmaMVPSession) LastCommittedBlock() (*big.Int, error) {
	return _PlasmaMVP.Contract.LastCommittedBlock(&_PlasmaMVP.CallOpts)
}

// LastCommittedBlock is a free data retrieval call binding the contract method 0x3acb097a.
//
// Solidity: function lastCommittedBlock() constant returns(uint256)
func (_PlasmaMVP *PlasmaMVPCallerSession) LastCommittedBlock() (*big.Int, error) {
	return _PlasmaMVP.Contract.LastCommittedBlock(&_PlasmaMVP.CallOpts)
}

// MinExitBond is a free data retrieval call binding the contract method 0xd68545a3.
//
// Solidity: function minExitBond() constant returns(uint256)
func (_PlasmaMVP *PlasmaMVPCaller) MinExitBond(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _PlasmaMVP.contract.Call(opts, out, "minExitBond")
	return *ret0, err
}

// MinExitBond is a free data retrieval call binding the contract method 0xd68545a3.
//
// Solidity: function minExitBond() constant returns(uint256)
func (_PlasmaMVP *PlasmaMVPSession) MinExitBond() (*big.Int, error) {
	return _PlasmaMVP.Contract.MinExitBond(&_PlasmaMVP.CallOpts)
}

// MinExitBond is a free data retrieval call binding the contract method 0xd68545a3.
//
// Solidity: function minExitBond() constant returns(uint256)
func (_PlasmaMVP *PlasmaMVPCallerSession) MinExitBond() (*big.Int, error) {
	return _PlasmaMVP.Contract.MinExitBond(&_PlasmaMVP.CallOpts)
}

// Operator is a free data retrieval call binding the contract method 0x570ca735.
//
// Solidity: function operator() constant returns(address)
func (_PlasmaMVP *PlasmaMVPCaller) Operator(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _PlasmaMVP.contract.Call(opts, out, "operator")
	return *ret0, err
}

// Operator is a free data retrieval call binding the contract method 0x570ca735.
//
// Solidity: function operator() constant returns(address)
func (_PlasmaMVP *PlasmaMVPSession) Operator() (common.Address, error) {
	return _PlasmaMVP.Contract.Operator(&_PlasmaMVP.CallOpts)
}

// Operator is a free data retrieval call binding the contract method 0x570ca735.
//
// Solidity: function operator() constant returns(address)
func (_PlasmaMVP *PlasmaMVPCallerSession) Operator() (common.Address, error) {
	return _PlasmaMVP.Contract.Operator(&_PlasmaMVP.CallOpts)
}

// PlasmaChain is a free data retrieval call binding the contract method 0x3cdac42c.
//
// Solidity: function plasmaChain(uint256 ) constant returns(bytes32 header, uint256 numTxns, uint256 feeAmount, uint256 createdAt)
func (_PlasmaMVP *PlasmaMVPCaller) PlasmaChain(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Header    [32]byte
	NumTxns   *big.Int
	FeeAmount *big.Int
	CreatedAt *big.Int
}, error) {
	ret := new(struct {
		Header    [32]byte
		NumTxns   *big.Int
		FeeAmount *big.Int
		CreatedAt *big.Int
	})
	out := ret
	err := _PlasmaMVP.contract.Call(opts, out, "plasmaChain", arg0)
	return *ret, err
}

// PlasmaChain is a free data retrieval call binding the contract method 0x3cdac42c.
//
// Solidity: function plasmaChain(uint256 ) constant returns(bytes32 header, uint256 numTxns, uint256 feeAmount, uint256 createdAt)
func (_PlasmaMVP *PlasmaMVPSession) PlasmaChain(arg0 *big.Int) (struct {
	Header    [32]byte
	NumTxns   *big.Int
	FeeAmount *big.Int
	CreatedAt *big.Int
}, error) {
	return _PlasmaMVP.Contract.PlasmaChain(&_PlasmaMVP.CallOpts, arg0)
}

// PlasmaChain is a free data retrieval call binding the contract method 0x3cdac42c.
//
// Solidity: function plasmaChain(uint256 ) constant returns(bytes32 header, uint256 numTxns, uint256 feeAmount, uint256 createdAt)
func (_PlasmaMVP *PlasmaMVPCallerSession) PlasmaChain(arg0 *big.Int) (struct {
	Header    [32]byte
	NumTxns   *big.Int
	FeeAmount *big.Int
	CreatedAt *big.Int
}, error) {
	return _PlasmaMVP.Contract.PlasmaChain(&_PlasmaMVP.CallOpts, arg0)
}

// PlasmaChainBalance is a free data retrieval call binding the contract method 0x45cbefa2.
//
// Solidity: function plasmaChainBalance() constant returns(uint256)
func (_PlasmaMVP *PlasmaMVPCaller) PlasmaChainBalance(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _PlasmaMVP.contract.Call(opts, out, "plasmaChainBalance")
	return *ret0, err
}

// PlasmaChainBalance is a free data retrieval call binding the contract method 0x45cbefa2.
//
// Solidity: function plasmaChainBalance() constant returns(uint256)
func (_PlasmaMVP *PlasmaMVPSession) PlasmaChainBalance() (*big.Int, error) {
	return _PlasmaMVP.Contract.PlasmaChainBalance(&_PlasmaMVP.CallOpts)
}

// PlasmaChainBalance is a free data retrieval call binding the contract method 0x45cbefa2.
//
// Solidity: function plasmaChainBalance() constant returns(uint256)
func (_PlasmaMVP *PlasmaMVPCallerSession) PlasmaChainBalance() (*big.Int, error) {
	return _PlasmaMVP.Contract.PlasmaChainBalance(&_PlasmaMVP.CallOpts)
}

// TotalWithdrawBalance is a free data retrieval call binding the contract method 0xc430c438.
//
// Solidity: function totalWithdrawBalance() constant returns(uint256)
func (_PlasmaMVP *PlasmaMVPCaller) TotalWithdrawBalance(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _PlasmaMVP.contract.Call(opts, out, "totalWithdrawBalance")
	return *ret0, err
}

// TotalWithdrawBalance is a free data retrieval call binding the contract method 0xc430c438.
//
// Solidity: function totalWithdrawBalance() constant returns(uint256)
func (_PlasmaMVP *PlasmaMVPSession) TotalWithdrawBalance() (*big.Int, error) {
	return _PlasmaMVP.Contract.TotalWithdrawBalance(&_PlasmaMVP.CallOpts)
}

// TotalWithdrawBalance is a free data retrieval call binding the contract method 0xc430c438.
//
// Solidity: function totalWithdrawBalance() constant returns(uint256)
func (_PlasmaMVP *PlasmaMVPCallerSession) TotalWithdrawBalance() (*big.Int, error) {
	return _PlasmaMVP.Contract.TotalWithdrawBalance(&_PlasmaMVP.CallOpts)
}

// TxExitQueue is a free data retrieval call binding the contract method 0x875b8ea0.
//
// Solidity: function txExitQueue(uint256 ) constant returns(uint256)
func (_PlasmaMVP *PlasmaMVPCaller) TxExitQueue(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _PlasmaMVP.contract.Call(opts, out, "txExitQueue", arg0)
	return *ret0, err
}

// TxExitQueue is a free data retrieval call binding the contract method 0x875b8ea0.
//
// Solidity: function txExitQueue(uint256 ) constant returns(uint256)
func (_PlasmaMVP *PlasmaMVPSession) TxExitQueue(arg0 *big.Int) (*big.Int, error) {
	return _PlasmaMVP.Contract.TxExitQueue(&_PlasmaMVP.CallOpts, arg0)
}

// TxExitQueue is a free data retrieval call binding the contract method 0x875b8ea0.
//
// Solidity: function txExitQueue(uint256 ) constant returns(uint256)
func (_PlasmaMVP *PlasmaMVPCallerSession) TxExitQueue(arg0 *big.Int) (*big.Int, error) {
	return _PlasmaMVP.Contract.TxExitQueue(&_PlasmaMVP.CallOpts, arg0)
}

// TxExits is a free data retrieval call binding the contract method 0x6d3d8b1a.
//
// Solidity: function txExits(uint256 ) constant returns(uint256 amount, uint256 committedFee, uint256 createdAt, address owner, uint8 state)
func (_PlasmaMVP *PlasmaMVPCaller) TxExits(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Amount       *big.Int
	CommittedFee *big.Int
	CreatedAt    *big.Int
	Owner        common.Address
	State        uint8
}, error) {
	ret := new(struct {
		Amount       *big.Int
		CommittedFee *big.Int
		CreatedAt    *big.Int
		Owner        common.Address
		State        uint8
	})
	out := ret
	err := _PlasmaMVP.contract.Call(opts, out, "txExits", arg0)
	return *ret, err
}

// TxExits is a free data retrieval call binding the contract method 0x6d3d8b1a.
//
// Solidity: function txExits(uint256 ) constant returns(uint256 amount, uint256 committedFee, uint256 createdAt, address owner, uint8 state)
func (_PlasmaMVP *PlasmaMVPSession) TxExits(arg0 *big.Int) (struct {
	Amount       *big.Int
	CommittedFee *big.Int
	CreatedAt    *big.Int
	Owner        common.Address
	State        uint8
}, error) {
	return _PlasmaMVP.Contract.TxExits(&_PlasmaMVP.CallOpts, arg0)
}

// TxExits is a free data retrieval call binding the contract method 0x6d3d8b1a.
//
// Solidity: function txExits(uint256 ) constant returns(uint256 amount, uint256 committedFee, uint256 createdAt, address owner, uint8 state)
func (_PlasmaMVP *PlasmaMVPCallerSession) TxExits(arg0 *big.Int) (struct {
	Amount       *big.Int
	CommittedFee *big.Int
	CreatedAt    *big.Int
	Owner        common.Address
	State        uint8
}, error) {
	return _PlasmaMVP.Contract.TxExits(&_PlasmaMVP.CallOpts, arg0)
}

// TxQueueLength is a free data retrieval call binding the contract method 0xb06f4777.
//
// Solidity: function txQueueLength() constant returns(uint256)
func (_PlasmaMVP *PlasmaMVPCaller) TxQueueLength(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _PlasmaMVP.contract.Call(opts, out, "txQueueLength")
	return *ret0, err
}

// TxQueueLength is a free data retrieval call binding the contract method 0xb06f4777.
//
// Solidity: function txQueueLength() constant returns(uint256)
func (_PlasmaMVP *PlasmaMVPSession) TxQueueLength() (*big.Int, error) {
	return _PlasmaMVP.Contract.TxQueueLength(&_PlasmaMVP.CallOpts)
}

// TxQueueLength is a free data retrieval call binding the contract method 0xb06f4777.
//
// Solidity: function txQueueLength() constant returns(uint256)
func (_PlasmaMVP *PlasmaMVPCallerSession) TxQueueLength() (*big.Int, error) {
	return _PlasmaMVP.Contract.TxQueueLength(&_PlasmaMVP.CallOpts)
}

// ChallengeExit is a paid mutator transaction binding the contract method 0xd344e8e4.
//
// Solidity: function challengeExit(uint256[4] exitingTxPos, uint256[2] challengingTxPos, bytes txBytes, bytes proof, bytes confirmSignature) returns()
func (_PlasmaMVP *PlasmaMVPTransactor) ChallengeExit(opts *bind.TransactOpts, exitingTxPos [4]*big.Int, challengingTxPos [2]*big.Int, txBytes []byte, proof []byte, confirmSignature []byte) (*types.Transaction, error) {
	return _PlasmaMVP.contract.Transact(opts, "challengeExit", exitingTxPos, challengingTxPos, txBytes, proof, confirmSignature)
}

// ChallengeExit is a paid mutator transaction binding the contract method 0xd344e8e4.
//
// Solidity: function challengeExit(uint256[4] exitingTxPos, uint256[2] challengingTxPos, bytes txBytes, bytes proof, bytes confirmSignature) returns()
func (_PlasmaMVP *PlasmaMVPSession) ChallengeExit(exitingTxPos [4]*big.Int, challengingTxPos [2]*big.Int, txBytes []byte, proof []byte, confirmSignature []byte) (*types.Transaction, error) {
	return _PlasmaMVP.Contract.ChallengeExit(&_PlasmaMVP.TransactOpts, exitingTxPos, challengingTxPos, txBytes, proof, confirmSignature)
}

// ChallengeExit is a paid mutator transaction binding the contract method 0xd344e8e4.
//
// Solidity: function challengeExit(uint256[4] exitingTxPos, uint256[2] challengingTxPos, bytes txBytes, bytes proof, bytes confirmSignature) returns()
func (_PlasmaMVP *PlasmaMVPTransactorSession) ChallengeExit(exitingTxPos [4]*big.Int, challengingTxPos [2]*big.Int, txBytes []byte, proof []byte, confirmSignature []byte) (*types.Transaction, error) {
	return _PlasmaMVP.Contract.ChallengeExit(&_PlasmaMVP.TransactOpts, exitingTxPos, challengingTxPos, txBytes, proof, confirmSignature)
}

// ChangeOperator is a paid mutator transaction binding the contract method 0x06394c9b.
//
// Solidity: function changeOperator(address newOperator) returns()
func (_PlasmaMVP *PlasmaMVPTransactor) ChangeOperator(opts *bind.TransactOpts, newOperator common.Address) (*types.Transaction, error) {
	return _PlasmaMVP.contract.Transact(opts, "changeOperator", newOperator)
}

// ChangeOperator is a paid mutator transaction binding the contract method 0x06394c9b.
//
// Solidity: function changeOperator(address newOperator) returns()
func (_PlasmaMVP *PlasmaMVPSession) ChangeOperator(newOperator common.Address) (*types.Transaction, error) {
	return _PlasmaMVP.Contract.ChangeOperator(&_PlasmaMVP.TransactOpts, newOperator)
}

// ChangeOperator is a paid mutator transaction binding the contract method 0x06394c9b.
//
// Solidity: function changeOperator(address newOperator) returns()
func (_PlasmaMVP *PlasmaMVPTransactorSession) ChangeOperator(newOperator common.Address) (*types.Transaction, error) {
	return _PlasmaMVP.Contract.ChangeOperator(&_PlasmaMVP.TransactOpts, newOperator)
}

// Deposit is a paid mutator transaction binding the contract method 0xf340fa01.
//
// Solidity: function deposit(address owner) returns()
func (_PlasmaMVP *PlasmaMVPTransactor) Deposit(opts *bind.TransactOpts, owner common.Address) (*types.Transaction, error) {
	return _PlasmaMVP.contract.Transact(opts, "deposit", owner)
}

// Deposit is a paid mutator transaction binding the contract method 0xf340fa01.
//
// Solidity: function deposit(address owner) returns()
func (_PlasmaMVP *PlasmaMVPSession) Deposit(owner common.Address) (*types.Transaction, error) {
	return _PlasmaMVP.Contract.Deposit(&_PlasmaMVP.TransactOpts, owner)
}

// Deposit is a paid mutator transaction binding the contract method 0xf340fa01.
//
// Solidity: function deposit(address owner) returns()
func (_PlasmaMVP *PlasmaMVPTransactorSession) Deposit(owner common.Address) (*types.Transaction, error) {
	return _PlasmaMVP.Contract.Deposit(&_PlasmaMVP.TransactOpts, owner)
}

// FinalizeDepositExits is a paid mutator transaction binding the contract method 0xfcf5f9eb.
//
// Solidity: function finalizeDepositExits() returns()
func (_PlasmaMVP *PlasmaMVPTransactor) FinalizeDepositExits(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PlasmaMVP.contract.Transact(opts, "finalizeDepositExits")
}

// FinalizeDepositExits is a paid mutator transaction binding the contract method 0xfcf5f9eb.
//
// Solidity: function finalizeDepositExits() returns()
func (_PlasmaMVP *PlasmaMVPSession) FinalizeDepositExits() (*types.Transaction, error) {
	return _PlasmaMVP.Contract.FinalizeDepositExits(&_PlasmaMVP.TransactOpts)
}

// FinalizeDepositExits is a paid mutator transaction binding the contract method 0xfcf5f9eb.
//
// Solidity: function finalizeDepositExits() returns()
func (_PlasmaMVP *PlasmaMVPTransactorSession) FinalizeDepositExits() (*types.Transaction, error) {
	return _PlasmaMVP.Contract.FinalizeDepositExits(&_PlasmaMVP.TransactOpts)
}

// FinalizeTransactionExits is a paid mutator transaction binding the contract method 0x884fc7d6.
//
// Solidity: function finalizeTransactionExits() returns()
func (_PlasmaMVP *PlasmaMVPTransactor) FinalizeTransactionExits(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PlasmaMVP.contract.Transact(opts, "finalizeTransactionExits")
}

// FinalizeTransactionExits is a paid mutator transaction binding the contract method 0x884fc7d6.
//
// Solidity: function finalizeTransactionExits() returns()
func (_PlasmaMVP *PlasmaMVPSession) FinalizeTransactionExits() (*types.Transaction, error) {
	return _PlasmaMVP.Contract.FinalizeTransactionExits(&_PlasmaMVP.TransactOpts)
}

// FinalizeTransactionExits is a paid mutator transaction binding the contract method 0x884fc7d6.
//
// Solidity: function finalizeTransactionExits() returns()
func (_PlasmaMVP *PlasmaMVPTransactorSession) FinalizeTransactionExits() (*types.Transaction, error) {
	return _PlasmaMVP.Contract.FinalizeTransactionExits(&_PlasmaMVP.TransactOpts)
}

// StartDepositExit is a paid mutator transaction binding the contract method 0x70e4abf6.
//
// Solidity: function startDepositExit(uint256 nonce, uint256 committedFee) returns()
func (_PlasmaMVP *PlasmaMVPTransactor) StartDepositExit(opts *bind.TransactOpts, nonce *big.Int, committedFee *big.Int) (*types.Transaction, error) {
	return _PlasmaMVP.contract.Transact(opts, "startDepositExit", nonce, committedFee)
}

// StartDepositExit is a paid mutator transaction binding the contract method 0x70e4abf6.
//
// Solidity: function startDepositExit(uint256 nonce, uint256 committedFee) returns()
func (_PlasmaMVP *PlasmaMVPSession) StartDepositExit(nonce *big.Int, committedFee *big.Int) (*types.Transaction, error) {
	return _PlasmaMVP.Contract.StartDepositExit(&_PlasmaMVP.TransactOpts, nonce, committedFee)
}

// StartDepositExit is a paid mutator transaction binding the contract method 0x70e4abf6.
//
// Solidity: function startDepositExit(uint256 nonce, uint256 committedFee) returns()
func (_PlasmaMVP *PlasmaMVPTransactorSession) StartDepositExit(nonce *big.Int, committedFee *big.Int) (*types.Transaction, error) {
	return _PlasmaMVP.Contract.StartDepositExit(&_PlasmaMVP.TransactOpts, nonce, committedFee)
}

// StartFeeExit is a paid mutator transaction binding the contract method 0xed1695df.
//
// Solidity: function startFeeExit(uint256 blockNumber, uint256 committedFee) returns()
func (_PlasmaMVP *PlasmaMVPTransactor) StartFeeExit(opts *bind.TransactOpts, blockNumber *big.Int, committedFee *big.Int) (*types.Transaction, error) {
	return _PlasmaMVP.contract.Transact(opts, "startFeeExit", blockNumber, committedFee)
}

// StartFeeExit is a paid mutator transaction binding the contract method 0xed1695df.
//
// Solidity: function startFeeExit(uint256 blockNumber, uint256 committedFee) returns()
func (_PlasmaMVP *PlasmaMVPSession) StartFeeExit(blockNumber *big.Int, committedFee *big.Int) (*types.Transaction, error) {
	return _PlasmaMVP.Contract.StartFeeExit(&_PlasmaMVP.TransactOpts, blockNumber, committedFee)
}

// StartFeeExit is a paid mutator transaction binding the contract method 0xed1695df.
//
// Solidity: function startFeeExit(uint256 blockNumber, uint256 committedFee) returns()
func (_PlasmaMVP *PlasmaMVPTransactorSession) StartFeeExit(blockNumber *big.Int, committedFee *big.Int) (*types.Transaction, error) {
	return _PlasmaMVP.Contract.StartFeeExit(&_PlasmaMVP.TransactOpts, blockNumber, committedFee)
}

// StartTransactionExit is a paid mutator transaction binding the contract method 0xcf024ea6.
//
// Solidity: function startTransactionExit(uint256[3] txPos, bytes txBytes, bytes proof, bytes confirmSignatures, uint256 committedFee) returns()
func (_PlasmaMVP *PlasmaMVPTransactor) StartTransactionExit(opts *bind.TransactOpts, txPos [3]*big.Int, txBytes []byte, proof []byte, confirmSignatures []byte, committedFee *big.Int) (*types.Transaction, error) {
	return _PlasmaMVP.contract.Transact(opts, "startTransactionExit", txPos, txBytes, proof, confirmSignatures, committedFee)
}

// StartTransactionExit is a paid mutator transaction binding the contract method 0xcf024ea6.
//
// Solidity: function startTransactionExit(uint256[3] txPos, bytes txBytes, bytes proof, bytes confirmSignatures, uint256 committedFee) returns()
func (_PlasmaMVP *PlasmaMVPSession) StartTransactionExit(txPos [3]*big.Int, txBytes []byte, proof []byte, confirmSignatures []byte, committedFee *big.Int) (*types.Transaction, error) {
	return _PlasmaMVP.Contract.StartTransactionExit(&_PlasmaMVP.TransactOpts, txPos, txBytes, proof, confirmSignatures, committedFee)
}

// StartTransactionExit is a paid mutator transaction binding the contract method 0xcf024ea6.
//
// Solidity: function startTransactionExit(uint256[3] txPos, bytes txBytes, bytes proof, bytes confirmSignatures, uint256 committedFee) returns()
func (_PlasmaMVP *PlasmaMVPTransactorSession) StartTransactionExit(txPos [3]*big.Int, txBytes []byte, proof []byte, confirmSignatures []byte, committedFee *big.Int) (*types.Transaction, error) {
	return _PlasmaMVP.Contract.StartTransactionExit(&_PlasmaMVP.TransactOpts, txPos, txBytes, proof, confirmSignatures, committedFee)
}

// SubmitBlock is a paid mutator transaction binding the contract method 0xd84ba62f.
//
// Solidity: function submitBlock(bytes32[] headers, uint256[] txnsPerBlock, uint256[] feePerBlock, uint256 blockNum) returns()
func (_PlasmaMVP *PlasmaMVPTransactor) SubmitBlock(opts *bind.TransactOpts, headers [][32]byte, txnsPerBlock []*big.Int, feePerBlock []*big.Int, blockNum *big.Int) (*types.Transaction, error) {
	return _PlasmaMVP.contract.Transact(opts, "submitBlock", headers, txnsPerBlock, feePerBlock, blockNum)
}

// SubmitBlock is a paid mutator transaction binding the contract method 0xd84ba62f.
//
// Solidity: function submitBlock(bytes32[] headers, uint256[] txnsPerBlock, uint256[] feePerBlock, uint256 blockNum) returns()
func (_PlasmaMVP *PlasmaMVPSession) SubmitBlock(headers [][32]byte, txnsPerBlock []*big.Int, feePerBlock []*big.Int, blockNum *big.Int) (*types.Transaction, error) {
	return _PlasmaMVP.Contract.SubmitBlock(&_PlasmaMVP.TransactOpts, headers, txnsPerBlock, feePerBlock, blockNum)
}

// SubmitBlock is a paid mutator transaction binding the contract method 0xd84ba62f.
//
// Solidity: function submitBlock(bytes32[] headers, uint256[] txnsPerBlock, uint256[] feePerBlock, uint256 blockNum) returns()
func (_PlasmaMVP *PlasmaMVPTransactorSession) SubmitBlock(headers [][32]byte, txnsPerBlock []*big.Int, feePerBlock []*big.Int, blockNum *big.Int) (*types.Transaction, error) {
	return _PlasmaMVP.Contract.SubmitBlock(&_PlasmaMVP.TransactOpts, headers, txnsPerBlock, feePerBlock, blockNum)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns(uint256)
func (_PlasmaMVP *PlasmaMVPTransactor) Withdraw(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PlasmaMVP.contract.Transact(opts, "withdraw")
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns(uint256)
func (_PlasmaMVP *PlasmaMVPSession) Withdraw() (*types.Transaction, error) {
	return _PlasmaMVP.Contract.Withdraw(&_PlasmaMVP.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns(uint256)
func (_PlasmaMVP *PlasmaMVPTransactorSession) Withdraw() (*types.Transaction, error) {
	return _PlasmaMVP.Contract.Withdraw(&_PlasmaMVP.TransactOpts)
}

// PlasmaMVPAddedToBalancesIterator is returned from FilterAddedToBalances and is used to iterate over the raw logs and unpacked data for AddedToBalances events raised by the PlasmaMVP contract.
type PlasmaMVPAddedToBalancesIterator struct {
	Event *PlasmaMVPAddedToBalances // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PlasmaMVPAddedToBalancesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PlasmaMVPAddedToBalances)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PlasmaMVPAddedToBalances)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PlasmaMVPAddedToBalancesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PlasmaMVPAddedToBalancesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PlasmaMVPAddedToBalances represents a AddedToBalances event raised by the PlasmaMVP contract.
type PlasmaMVPAddedToBalances struct {
	Owner  common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterAddedToBalances is a free log retrieval operation binding the contract event 0xf8552a24c7d58fd05114f6fc9db7b3a354db64d5fc758184af1696ccd8f158f3.
//
// Solidity: event AddedToBalances(address owner, uint256 amount)
func (_PlasmaMVP *PlasmaMVPFilterer) FilterAddedToBalances(opts *bind.FilterOpts) (*PlasmaMVPAddedToBalancesIterator, error) {

	logs, sub, err := _PlasmaMVP.contract.FilterLogs(opts, "AddedToBalances")
	if err != nil {
		return nil, err
	}
	return &PlasmaMVPAddedToBalancesIterator{contract: _PlasmaMVP.contract, event: "AddedToBalances", logs: logs, sub: sub}, nil
}

// WatchAddedToBalances is a free log subscription operation binding the contract event 0xf8552a24c7d58fd05114f6fc9db7b3a354db64d5fc758184af1696ccd8f158f3.
//
// Solidity: event AddedToBalances(address owner, uint256 amount)
func (_PlasmaMVP *PlasmaMVPFilterer) WatchAddedToBalances(opts *bind.WatchOpts, sink chan<- *PlasmaMVPAddedToBalances) (event.Subscription, error) {

	logs, sub, err := _PlasmaMVP.contract.WatchLogs(opts, "AddedToBalances")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PlasmaMVPAddedToBalances)
				if err := _PlasmaMVP.contract.UnpackLog(event, "AddedToBalances", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// PlasmaMVPBlockSubmittedIterator is returned from FilterBlockSubmitted and is used to iterate over the raw logs and unpacked data for BlockSubmitted events raised by the PlasmaMVP contract.
type PlasmaMVPBlockSubmittedIterator struct {
	Event *PlasmaMVPBlockSubmitted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PlasmaMVPBlockSubmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PlasmaMVPBlockSubmitted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PlasmaMVPBlockSubmitted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PlasmaMVPBlockSubmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PlasmaMVPBlockSubmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PlasmaMVPBlockSubmitted represents a BlockSubmitted event raised by the PlasmaMVP contract.
type PlasmaMVPBlockSubmitted struct {
	Header      [32]byte
	BlockNumber *big.Int
	NumTxns     *big.Int
	FeeAmount   *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterBlockSubmitted is a free log retrieval operation binding the contract event 0x044ff3798f9b3ad55d1155cea9a40508c71b4c64335f5dae87e8e11551515a06.
//
// Solidity: event BlockSubmitted(bytes32 header, uint256 blockNumber, uint256 numTxns, uint256 feeAmount)
func (_PlasmaMVP *PlasmaMVPFilterer) FilterBlockSubmitted(opts *bind.FilterOpts) (*PlasmaMVPBlockSubmittedIterator, error) {

	logs, sub, err := _PlasmaMVP.contract.FilterLogs(opts, "BlockSubmitted")
	if err != nil {
		return nil, err
	}
	return &PlasmaMVPBlockSubmittedIterator{contract: _PlasmaMVP.contract, event: "BlockSubmitted", logs: logs, sub: sub}, nil
}

// WatchBlockSubmitted is a free log subscription operation binding the contract event 0x044ff3798f9b3ad55d1155cea9a40508c71b4c64335f5dae87e8e11551515a06.
//
// Solidity: event BlockSubmitted(bytes32 header, uint256 blockNumber, uint256 numTxns, uint256 feeAmount)
func (_PlasmaMVP *PlasmaMVPFilterer) WatchBlockSubmitted(opts *bind.WatchOpts, sink chan<- *PlasmaMVPBlockSubmitted) (event.Subscription, error) {

	logs, sub, err := _PlasmaMVP.contract.WatchLogs(opts, "BlockSubmitted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PlasmaMVPBlockSubmitted)
				if err := _PlasmaMVP.contract.UnpackLog(event, "BlockSubmitted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// PlasmaMVPChallengedExitIterator is returned from FilterChallengedExit and is used to iterate over the raw logs and unpacked data for ChallengedExit events raised by the PlasmaMVP contract.
type PlasmaMVPChallengedExitIterator struct {
	Event *PlasmaMVPChallengedExit // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PlasmaMVPChallengedExitIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PlasmaMVPChallengedExit)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PlasmaMVPChallengedExit)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PlasmaMVPChallengedExitIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PlasmaMVPChallengedExitIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PlasmaMVPChallengedExit represents a ChallengedExit event raised by the PlasmaMVP contract.
type PlasmaMVPChallengedExit struct {
	Position [4]*big.Int
	Owner    common.Address
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterChallengedExit is a free log retrieval operation binding the contract event 0xe1289dafb1083e540206bcd7d95a9705ba2590d6a9229c35a1c4c4c5efbda901.
//
// Solidity: event ChallengedExit(uint256[4] position, address owner, uint256 amount)
func (_PlasmaMVP *PlasmaMVPFilterer) FilterChallengedExit(opts *bind.FilterOpts) (*PlasmaMVPChallengedExitIterator, error) {

	logs, sub, err := _PlasmaMVP.contract.FilterLogs(opts, "ChallengedExit")
	if err != nil {
		return nil, err
	}
	return &PlasmaMVPChallengedExitIterator{contract: _PlasmaMVP.contract, event: "ChallengedExit", logs: logs, sub: sub}, nil
}

// WatchChallengedExit is a free log subscription operation binding the contract event 0xe1289dafb1083e540206bcd7d95a9705ba2590d6a9229c35a1c4c4c5efbda901.
//
// Solidity: event ChallengedExit(uint256[4] position, address owner, uint256 amount)
func (_PlasmaMVP *PlasmaMVPFilterer) WatchChallengedExit(opts *bind.WatchOpts, sink chan<- *PlasmaMVPChallengedExit) (event.Subscription, error) {

	logs, sub, err := _PlasmaMVP.contract.WatchLogs(opts, "ChallengedExit")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PlasmaMVPChallengedExit)
				if err := _PlasmaMVP.contract.UnpackLog(event, "ChallengedExit", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// PlasmaMVPChangedOperatorIterator is returned from FilterChangedOperator and is used to iterate over the raw logs and unpacked data for ChangedOperator events raised by the PlasmaMVP contract.
type PlasmaMVPChangedOperatorIterator struct {
	Event *PlasmaMVPChangedOperator // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PlasmaMVPChangedOperatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PlasmaMVPChangedOperator)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PlasmaMVPChangedOperator)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PlasmaMVPChangedOperatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PlasmaMVPChangedOperatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PlasmaMVPChangedOperator represents a ChangedOperator event raised by the PlasmaMVP contract.
type PlasmaMVPChangedOperator struct {
	OldOperator common.Address
	NewOperator common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterChangedOperator is a free log retrieval operation binding the contract event 0x3aff32e6289f2f2a2463481071051456b2768bb391c64ae3c91f9033e208cda1.
//
// Solidity: event ChangedOperator(address oldOperator, address newOperator)
func (_PlasmaMVP *PlasmaMVPFilterer) FilterChangedOperator(opts *bind.FilterOpts) (*PlasmaMVPChangedOperatorIterator, error) {

	logs, sub, err := _PlasmaMVP.contract.FilterLogs(opts, "ChangedOperator")
	if err != nil {
		return nil, err
	}
	return &PlasmaMVPChangedOperatorIterator{contract: _PlasmaMVP.contract, event: "ChangedOperator", logs: logs, sub: sub}, nil
}

// WatchChangedOperator is a free log subscription operation binding the contract event 0x3aff32e6289f2f2a2463481071051456b2768bb391c64ae3c91f9033e208cda1.
//
// Solidity: event ChangedOperator(address oldOperator, address newOperator)
func (_PlasmaMVP *PlasmaMVPFilterer) WatchChangedOperator(opts *bind.WatchOpts, sink chan<- *PlasmaMVPChangedOperator) (event.Subscription, error) {

	logs, sub, err := _PlasmaMVP.contract.WatchLogs(opts, "ChangedOperator")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PlasmaMVPChangedOperator)
				if err := _PlasmaMVP.contract.UnpackLog(event, "ChangedOperator", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// PlasmaMVPDepositIterator is returned from FilterDeposit and is used to iterate over the raw logs and unpacked data for Deposit events raised by the PlasmaMVP contract.
type PlasmaMVPDepositIterator struct {
	Event *PlasmaMVPDeposit // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PlasmaMVPDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PlasmaMVPDeposit)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PlasmaMVPDeposit)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PlasmaMVPDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PlasmaMVPDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PlasmaMVPDeposit represents a Deposit event raised by the PlasmaMVP contract.
type PlasmaMVPDeposit struct {
	Depositor    common.Address
	Amount       *big.Int
	DepositNonce *big.Int
	EthBlockNum  *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterDeposit is a free log retrieval operation binding the contract event 0x36af321ec8d3c75236829c5317affd40ddb308863a1236d2d277a4025cccee1e.
//
// Solidity: event Deposit(address depositor, uint256 amount, uint256 depositNonce, uint256 ethBlockNum)
func (_PlasmaMVP *PlasmaMVPFilterer) FilterDeposit(opts *bind.FilterOpts) (*PlasmaMVPDepositIterator, error) {

	logs, sub, err := _PlasmaMVP.contract.FilterLogs(opts, "Deposit")
	if err != nil {
		return nil, err
	}
	return &PlasmaMVPDepositIterator{contract: _PlasmaMVP.contract, event: "Deposit", logs: logs, sub: sub}, nil
}

// WatchDeposit is a free log subscription operation binding the contract event 0x36af321ec8d3c75236829c5317affd40ddb308863a1236d2d277a4025cccee1e.
//
// Solidity: event Deposit(address depositor, uint256 amount, uint256 depositNonce, uint256 ethBlockNum)
func (_PlasmaMVP *PlasmaMVPFilterer) WatchDeposit(opts *bind.WatchOpts, sink chan<- *PlasmaMVPDeposit) (event.Subscription, error) {

	logs, sub, err := _PlasmaMVP.contract.WatchLogs(opts, "Deposit")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PlasmaMVPDeposit)
				if err := _PlasmaMVP.contract.UnpackLog(event, "Deposit", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// PlasmaMVPFinalizedExitIterator is returned from FilterFinalizedExit and is used to iterate over the raw logs and unpacked data for FinalizedExit events raised by the PlasmaMVP contract.
type PlasmaMVPFinalizedExitIterator struct {
	Event *PlasmaMVPFinalizedExit // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PlasmaMVPFinalizedExitIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PlasmaMVPFinalizedExit)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PlasmaMVPFinalizedExit)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PlasmaMVPFinalizedExitIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PlasmaMVPFinalizedExitIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PlasmaMVPFinalizedExit represents a FinalizedExit event raised by the PlasmaMVP contract.
type PlasmaMVPFinalizedExit struct {
	Position [4]*big.Int
	Owner    common.Address
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterFinalizedExit is a free log retrieval operation binding the contract event 0xb5083a27a38f8a9aa999efb3306b7be96dc3f42010a968dd86627880ba7fdbe2.
//
// Solidity: event FinalizedExit(uint256[4] position, address owner, uint256 amount)
func (_PlasmaMVP *PlasmaMVPFilterer) FilterFinalizedExit(opts *bind.FilterOpts) (*PlasmaMVPFinalizedExitIterator, error) {

	logs, sub, err := _PlasmaMVP.contract.FilterLogs(opts, "FinalizedExit")
	if err != nil {
		return nil, err
	}
	return &PlasmaMVPFinalizedExitIterator{contract: _PlasmaMVP.contract, event: "FinalizedExit", logs: logs, sub: sub}, nil
}

// WatchFinalizedExit is a free log subscription operation binding the contract event 0xb5083a27a38f8a9aa999efb3306b7be96dc3f42010a968dd86627880ba7fdbe2.
//
// Solidity: event FinalizedExit(uint256[4] position, address owner, uint256 amount)
func (_PlasmaMVP *PlasmaMVPFilterer) WatchFinalizedExit(opts *bind.WatchOpts, sink chan<- *PlasmaMVPFinalizedExit) (event.Subscription, error) {

	logs, sub, err := _PlasmaMVP.contract.WatchLogs(opts, "FinalizedExit")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PlasmaMVPFinalizedExit)
				if err := _PlasmaMVP.contract.UnpackLog(event, "FinalizedExit", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// PlasmaMVPStartedDepositExitIterator is returned from FilterStartedDepositExit and is used to iterate over the raw logs and unpacked data for StartedDepositExit events raised by the PlasmaMVP contract.
type PlasmaMVPStartedDepositExitIterator struct {
	Event *PlasmaMVPStartedDepositExit // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PlasmaMVPStartedDepositExitIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PlasmaMVPStartedDepositExit)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PlasmaMVPStartedDepositExit)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PlasmaMVPStartedDepositExitIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PlasmaMVPStartedDepositExitIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PlasmaMVPStartedDepositExit represents a StartedDepositExit event raised by the PlasmaMVP contract.
type PlasmaMVPStartedDepositExit struct {
	Nonce        *big.Int
	Owner        common.Address
	Amount       *big.Int
	CommittedFee *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterStartedDepositExit is a free log retrieval operation binding the contract event 0xe90dc7204eb622c87d2c8d67d8e27afdfd34042584591e7b3d35014873cf9cfd.
//
// Solidity: event StartedDepositExit(uint256 nonce, address owner, uint256 amount, uint256 committedFee)
func (_PlasmaMVP *PlasmaMVPFilterer) FilterStartedDepositExit(opts *bind.FilterOpts) (*PlasmaMVPStartedDepositExitIterator, error) {

	logs, sub, err := _PlasmaMVP.contract.FilterLogs(opts, "StartedDepositExit")
	if err != nil {
		return nil, err
	}
	return &PlasmaMVPStartedDepositExitIterator{contract: _PlasmaMVP.contract, event: "StartedDepositExit", logs: logs, sub: sub}, nil
}

// WatchStartedDepositExit is a free log subscription operation binding the contract event 0xe90dc7204eb622c87d2c8d67d8e27afdfd34042584591e7b3d35014873cf9cfd.
//
// Solidity: event StartedDepositExit(uint256 nonce, address owner, uint256 amount, uint256 committedFee)
func (_PlasmaMVP *PlasmaMVPFilterer) WatchStartedDepositExit(opts *bind.WatchOpts, sink chan<- *PlasmaMVPStartedDepositExit) (event.Subscription, error) {

	logs, sub, err := _PlasmaMVP.contract.WatchLogs(opts, "StartedDepositExit")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PlasmaMVPStartedDepositExit)
				if err := _PlasmaMVP.contract.UnpackLog(event, "StartedDepositExit", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// PlasmaMVPStartedTransactionExitIterator is returned from FilterStartedTransactionExit and is used to iterate over the raw logs and unpacked data for StartedTransactionExit events raised by the PlasmaMVP contract.
type PlasmaMVPStartedTransactionExitIterator struct {
	Event *PlasmaMVPStartedTransactionExit // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *PlasmaMVPStartedTransactionExitIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PlasmaMVPStartedTransactionExit)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(PlasmaMVPStartedTransactionExit)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *PlasmaMVPStartedTransactionExitIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PlasmaMVPStartedTransactionExitIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PlasmaMVPStartedTransactionExit represents a StartedTransactionExit event raised by the PlasmaMVP contract.
type PlasmaMVPStartedTransactionExit struct {
	Position          [3]*big.Int
	Owner             common.Address
	Amount            *big.Int
	ConfirmSignatures []byte
	CommittedFee      *big.Int
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterStartedTransactionExit is a free log retrieval operation binding the contract event 0x20d695720ae96d3511520c6f51d6ab23aa19a3796da77024ad027b344bb72530.
//
// Solidity: event StartedTransactionExit(uint256[3] position, address owner, uint256 amount, bytes confirmSignatures, uint256 committedFee)
func (_PlasmaMVP *PlasmaMVPFilterer) FilterStartedTransactionExit(opts *bind.FilterOpts) (*PlasmaMVPStartedTransactionExitIterator, error) {

	logs, sub, err := _PlasmaMVP.contract.FilterLogs(opts, "StartedTransactionExit")
	if err != nil {
		return nil, err
	}
	return &PlasmaMVPStartedTransactionExitIterator{contract: _PlasmaMVP.contract, event: "StartedTransactionExit", logs: logs, sub: sub}, nil
}

// WatchStartedTransactionExit is a free log subscription operation binding the contract event 0x20d695720ae96d3511520c6f51d6ab23aa19a3796da77024ad027b344bb72530.
//
// Solidity: event StartedTransactionExit(uint256[3] position, address owner, uint256 amount, bytes confirmSignatures, uint256 committedFee)
func (_PlasmaMVP *PlasmaMVPFilterer) WatchStartedTransactionExit(opts *bind.WatchOpts, sink chan<- *PlasmaMVPStartedTransactionExit) (event.Subscription, error) {

	logs, sub, err := _PlasmaMVP.contract.WatchLogs(opts, "StartedTransactionExit")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PlasmaMVPStartedTransactionExit)
				if err := _PlasmaMVP.contract.UnpackLog(event, "StartedTransactionExit", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}
