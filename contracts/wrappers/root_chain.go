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

// RootChainABI is the input ABI used to generate the binding from.
const RootChainABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"txIndexFactor\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"balances\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"lastCommittedBlock\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"txExits\",\"outputs\":[{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"createdAt\",\"type\":\"uint256\"},{\"name\":\"owner\",\"type\":\"address\"},{\"name\":\"state\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"blockIndexFactor\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"deposits\",\"outputs\":[{\"name\":\"owner\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"createdAt\",\"type\":\"uint256\"},{\"name\":\"ethBlocknum\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalWithdrawBalance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"depositExits\",\"outputs\":[{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"createdAt\",\"type\":\"uint256\"},{\"name\":\"owner\",\"type\":\"address\"},{\"name\":\"state\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"depositNonce\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"childChain\",\"outputs\":[{\"name\":\"root\",\"type\":\"bytes32\"},{\"name\":\"numTxns\",\"type\":\"uint256\"},{\"name\":\"feeAmount\",\"type\":\"uint256\"},{\"name\":\"createdAt\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"AddedToBalances\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"root\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"numTxns\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"feeAmount\",\"type\":\"uint256\"}],\"name\":\"BlockSubmitted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"depositor\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"depositNonce\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"ethBlockNum\",\"type\":\"uint256\"}],\"name\":\"Deposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"position\",\"type\":\"uint256[3]\"},{\"indexed\":false,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"confirmSignatures\",\"type\":\"bytes\"}],\"name\":\"StartedTransactionExit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"nonce\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"StartedDepositExit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"position\",\"type\":\"uint256[4]\"},{\"indexed\":false,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"ChallengedExit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"position\",\"type\":\"uint256[4]\"},{\"indexed\":false,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"FinalizedExit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"name\":\"blocks\",\"type\":\"bytes\"},{\"name\":\"txnsPerBlock\",\"type\":\"uint256[]\"},{\"name\":\"feesPerBlock\",\"type\":\"uint256[]\"},{\"name\":\"blockNum\",\"type\":\"uint256\"}],\"name\":\"submitBlock\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"deposit\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"nonce\",\"type\":\"uint256\"}],\"name\":\"startDepositExit\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"txPos\",\"type\":\"uint256[3]\"},{\"name\":\"txBytes\",\"type\":\"bytes\"},{\"name\":\"proof\",\"type\":\"bytes\"},{\"name\":\"confirmSignatures\",\"type\":\"bytes\"}],\"name\":\"startTransactionExit\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"blockNumber\",\"type\":\"uint256\"}],\"name\":\"startFeeExit\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"nonce\",\"type\":\"uint256\"},{\"name\":\"newTxPos\",\"type\":\"uint256[3]\"},{\"name\":\"txBytes\",\"type\":\"bytes\"},{\"name\":\"proof\",\"type\":\"bytes\"},{\"name\":\"confirmSignature\",\"type\":\"bytes\"}],\"name\":\"challengeDepositExit\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"exitingTxPos\",\"type\":\"uint256[3]\"},{\"name\":\"challengingTxPos\",\"type\":\"uint256[3]\"},{\"name\":\"txBytes\",\"type\":\"bytes\"},{\"name\":\"proof\",\"type\":\"bytes\"},{\"name\":\"confirmSignature\",\"type\":\"bytes\"}],\"name\":\"challengeTransactionExit\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"finalizeDepositExits\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"finalizeTransactionExits\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"withdraw\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"childChainBalance\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_address\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// RootChain is an auto generated Go binding around an Ethereum contract.
type RootChain struct {
	RootChainCaller     // Read-only binding to the contract
	RootChainTransactor // Write-only binding to the contract
	RootChainFilterer   // Log filterer for contract events
}

// RootChainCaller is an auto generated read-only Go binding around an Ethereum contract.
type RootChainCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RootChainTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RootChainTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RootChainFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RootChainFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RootChainSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RootChainSession struct {
	Contract     *RootChain        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RootChainCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RootChainCallerSession struct {
	Contract *RootChainCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// RootChainTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RootChainTransactorSession struct {
	Contract     *RootChainTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// RootChainRaw is an auto generated low-level Go binding around an Ethereum contract.
type RootChainRaw struct {
	Contract *RootChain // Generic contract binding to access the raw methods on
}

// RootChainCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RootChainCallerRaw struct {
	Contract *RootChainCaller // Generic read-only contract binding to access the raw methods on
}

// RootChainTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RootChainTransactorRaw struct {
	Contract *RootChainTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRootChain creates a new instance of RootChain, bound to a specific deployed contract.
func NewRootChain(address common.Address, backend bind.ContractBackend) (*RootChain, error) {
	contract, err := bindRootChain(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &RootChain{RootChainCaller: RootChainCaller{contract: contract}, RootChainTransactor: RootChainTransactor{contract: contract}, RootChainFilterer: RootChainFilterer{contract: contract}}, nil
}

// NewRootChainCaller creates a new read-only instance of RootChain, bound to a specific deployed contract.
func NewRootChainCaller(address common.Address, caller bind.ContractCaller) (*RootChainCaller, error) {
	contract, err := bindRootChain(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RootChainCaller{contract: contract}, nil
}

// NewRootChainTransactor creates a new write-only instance of RootChain, bound to a specific deployed contract.
func NewRootChainTransactor(address common.Address, transactor bind.ContractTransactor) (*RootChainTransactor, error) {
	contract, err := bindRootChain(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RootChainTransactor{contract: contract}, nil
}

// NewRootChainFilterer creates a new log filterer instance of RootChain, bound to a specific deployed contract.
func NewRootChainFilterer(address common.Address, filterer bind.ContractFilterer) (*RootChainFilterer, error) {
	contract, err := bindRootChain(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RootChainFilterer{contract: contract}, nil
}

// bindRootChain binds a generic wrapper to an already deployed contract.
func bindRootChain(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(RootChainABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RootChain *RootChainRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _RootChain.Contract.RootChainCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RootChain *RootChainRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RootChain.Contract.RootChainTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RootChain *RootChainRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RootChain.Contract.RootChainTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_RootChain *RootChainCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _RootChain.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_RootChain *RootChainTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RootChain.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_RootChain *RootChainTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _RootChain.Contract.contract.Transact(opts, method, params...)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_address address) constant returns(uint256)
func (_RootChain *RootChainCaller) BalanceOf(opts *bind.CallOpts, _address common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _RootChain.contract.Call(opts, out, "balanceOf", _address)
	return *ret0, err
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_address address) constant returns(uint256)
func (_RootChain *RootChainSession) BalanceOf(_address common.Address) (*big.Int, error) {
	return _RootChain.Contract.BalanceOf(&_RootChain.CallOpts, _address)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(_address address) constant returns(uint256)
func (_RootChain *RootChainCallerSession) BalanceOf(_address common.Address) (*big.Int, error) {
	return _RootChain.Contract.BalanceOf(&_RootChain.CallOpts, _address)
}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances( address) constant returns(uint256)
func (_RootChain *RootChainCaller) Balances(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _RootChain.contract.Call(opts, out, "balances", arg0)
	return *ret0, err
}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances( address) constant returns(uint256)
func (_RootChain *RootChainSession) Balances(arg0 common.Address) (*big.Int, error) {
	return _RootChain.Contract.Balances(&_RootChain.CallOpts, arg0)
}

// Balances is a free data retrieval call binding the contract method 0x27e235e3.
//
// Solidity: function balances( address) constant returns(uint256)
func (_RootChain *RootChainCallerSession) Balances(arg0 common.Address) (*big.Int, error) {
	return _RootChain.Contract.Balances(&_RootChain.CallOpts, arg0)
}

// BlockIndexFactor is a free data retrieval call binding the contract method 0x89609149.
//
// Solidity: function blockIndexFactor() constant returns(uint256)
func (_RootChain *RootChainCaller) BlockIndexFactor(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _RootChain.contract.Call(opts, out, "blockIndexFactor")
	return *ret0, err
}

// BlockIndexFactor is a free data retrieval call binding the contract method 0x89609149.
//
// Solidity: function blockIndexFactor() constant returns(uint256)
func (_RootChain *RootChainSession) BlockIndexFactor() (*big.Int, error) {
	return _RootChain.Contract.BlockIndexFactor(&_RootChain.CallOpts)
}

// BlockIndexFactor is a free data retrieval call binding the contract method 0x89609149.
//
// Solidity: function blockIndexFactor() constant returns(uint256)
func (_RootChain *RootChainCallerSession) BlockIndexFactor() (*big.Int, error) {
	return _RootChain.Contract.BlockIndexFactor(&_RootChain.CallOpts)
}

// ChildChain is a free data retrieval call binding the contract method 0xf95643b1.
//
// Solidity: function childChain( uint256) constant returns(root bytes32, numTxns uint256, feeAmount uint256, createdAt uint256)
func (_RootChain *RootChainCaller) ChildChain(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Root      [32]byte
	NumTxns   *big.Int
	FeeAmount *big.Int
	CreatedAt *big.Int
}, error) {
	ret := new(struct {
		Root      [32]byte
		NumTxns   *big.Int
		FeeAmount *big.Int
		CreatedAt *big.Int
	})
	out := ret
	err := _RootChain.contract.Call(opts, out, "childChain", arg0)
	return *ret, err
}

// ChildChain is a free data retrieval call binding the contract method 0xf95643b1.
//
// Solidity: function childChain( uint256) constant returns(root bytes32, numTxns uint256, feeAmount uint256, createdAt uint256)
func (_RootChain *RootChainSession) ChildChain(arg0 *big.Int) (struct {
	Root      [32]byte
	NumTxns   *big.Int
	FeeAmount *big.Int
	CreatedAt *big.Int
}, error) {
	return _RootChain.Contract.ChildChain(&_RootChain.CallOpts, arg0)
}

// ChildChain is a free data retrieval call binding the contract method 0xf95643b1.
//
// Solidity: function childChain( uint256) constant returns(root bytes32, numTxns uint256, feeAmount uint256, createdAt uint256)
func (_RootChain *RootChainCallerSession) ChildChain(arg0 *big.Int) (struct {
	Root      [32]byte
	NumTxns   *big.Int
	FeeAmount *big.Int
	CreatedAt *big.Int
}, error) {
	return _RootChain.Contract.ChildChain(&_RootChain.CallOpts, arg0)
}

// ChildChainBalance is a free data retrieval call binding the contract method 0x385e2fd3.
//
// Solidity: function childChainBalance() constant returns(uint256)
func (_RootChain *RootChainCaller) ChildChainBalance(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _RootChain.contract.Call(opts, out, "childChainBalance")
	return *ret0, err
}

// ChildChainBalance is a free data retrieval call binding the contract method 0x385e2fd3.
//
// Solidity: function childChainBalance() constant returns(uint256)
func (_RootChain *RootChainSession) ChildChainBalance() (*big.Int, error) {
	return _RootChain.Contract.ChildChainBalance(&_RootChain.CallOpts)
}

// ChildChainBalance is a free data retrieval call binding the contract method 0x385e2fd3.
//
// Solidity: function childChainBalance() constant returns(uint256)
func (_RootChain *RootChainCallerSession) ChildChainBalance() (*big.Int, error) {
	return _RootChain.Contract.ChildChainBalance(&_RootChain.CallOpts)
}

// DepositExits is a free data retrieval call binding the contract method 0xce84f906.
//
// Solidity: function depositExits( uint256) constant returns(amount uint256, createdAt uint256, owner address, state uint8)
func (_RootChain *RootChainCaller) DepositExits(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Amount    *big.Int
	CreatedAt *big.Int
	Owner     common.Address
	State     uint8
}, error) {
	ret := new(struct {
		Amount    *big.Int
		CreatedAt *big.Int
		Owner     common.Address
		State     uint8
	})
	out := ret
	err := _RootChain.contract.Call(opts, out, "depositExits", arg0)
	return *ret, err
}

// DepositExits is a free data retrieval call binding the contract method 0xce84f906.
//
// Solidity: function depositExits( uint256) constant returns(amount uint256, createdAt uint256, owner address, state uint8)
func (_RootChain *RootChainSession) DepositExits(arg0 *big.Int) (struct {
	Amount    *big.Int
	CreatedAt *big.Int
	Owner     common.Address
	State     uint8
}, error) {
	return _RootChain.Contract.DepositExits(&_RootChain.CallOpts, arg0)
}

// DepositExits is a free data retrieval call binding the contract method 0xce84f906.
//
// Solidity: function depositExits( uint256) constant returns(amount uint256, createdAt uint256, owner address, state uint8)
func (_RootChain *RootChainCallerSession) DepositExits(arg0 *big.Int) (struct {
	Amount    *big.Int
	CreatedAt *big.Int
	Owner     common.Address
	State     uint8
}, error) {
	return _RootChain.Contract.DepositExits(&_RootChain.CallOpts, arg0)
}

// DepositNonce is a free data retrieval call binding the contract method 0xde35f5cb.
//
// Solidity: function depositNonce() constant returns(uint256)
func (_RootChain *RootChainCaller) DepositNonce(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _RootChain.contract.Call(opts, out, "depositNonce")
	return *ret0, err
}

// DepositNonce is a free data retrieval call binding the contract method 0xde35f5cb.
//
// Solidity: function depositNonce() constant returns(uint256)
func (_RootChain *RootChainSession) DepositNonce() (*big.Int, error) {
	return _RootChain.Contract.DepositNonce(&_RootChain.CallOpts)
}

// DepositNonce is a free data retrieval call binding the contract method 0xde35f5cb.
//
// Solidity: function depositNonce() constant returns(uint256)
func (_RootChain *RootChainCallerSession) DepositNonce() (*big.Int, error) {
	return _RootChain.Contract.DepositNonce(&_RootChain.CallOpts)
}

// Deposits is a free data retrieval call binding the contract method 0xb02c43d0.
//
// Solidity: function deposits( uint256) constant returns(owner address, amount uint256, createdAt uint256, ethBlocknum uint256)
func (_RootChain *RootChainCaller) Deposits(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Owner       common.Address
	Amount      *big.Int
	CreatedAt   *big.Int
	EthBlocknum *big.Int
}, error) {
	ret := new(struct {
		Owner       common.Address
		Amount      *big.Int
		CreatedAt   *big.Int
		EthBlocknum *big.Int
	})
	out := ret
	err := _RootChain.contract.Call(opts, out, "deposits", arg0)
	return *ret, err
}

// Deposits is a free data retrieval call binding the contract method 0xb02c43d0.
//
// Solidity: function deposits( uint256) constant returns(owner address, amount uint256, createdAt uint256, ethBlocknum uint256)
func (_RootChain *RootChainSession) Deposits(arg0 *big.Int) (struct {
	Owner       common.Address
	Amount      *big.Int
	CreatedAt   *big.Int
	EthBlocknum *big.Int
}, error) {
	return _RootChain.Contract.Deposits(&_RootChain.CallOpts, arg0)
}

// Deposits is a free data retrieval call binding the contract method 0xb02c43d0.
//
// Solidity: function deposits( uint256) constant returns(owner address, amount uint256, createdAt uint256, ethBlocknum uint256)
func (_RootChain *RootChainCallerSession) Deposits(arg0 *big.Int) (struct {
	Owner       common.Address
	Amount      *big.Int
	CreatedAt   *big.Int
	EthBlocknum *big.Int
}, error) {
	return _RootChain.Contract.Deposits(&_RootChain.CallOpts, arg0)
}

// LastCommittedBlock is a free data retrieval call binding the contract method 0x3acb097a.
//
// Solidity: function lastCommittedBlock() constant returns(uint256)
func (_RootChain *RootChainCaller) LastCommittedBlock(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _RootChain.contract.Call(opts, out, "lastCommittedBlock")
	return *ret0, err
}

// LastCommittedBlock is a free data retrieval call binding the contract method 0x3acb097a.
//
// Solidity: function lastCommittedBlock() constant returns(uint256)
func (_RootChain *RootChainSession) LastCommittedBlock() (*big.Int, error) {
	return _RootChain.Contract.LastCommittedBlock(&_RootChain.CallOpts)
}

// LastCommittedBlock is a free data retrieval call binding the contract method 0x3acb097a.
//
// Solidity: function lastCommittedBlock() constant returns(uint256)
func (_RootChain *RootChainCallerSession) LastCommittedBlock() (*big.Int, error) {
	return _RootChain.Contract.LastCommittedBlock(&_RootChain.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_RootChain *RootChainCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _RootChain.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_RootChain *RootChainSession) Owner() (common.Address, error) {
	return _RootChain.Contract.Owner(&_RootChain.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_RootChain *RootChainCallerSession) Owner() (common.Address, error) {
	return _RootChain.Contract.Owner(&_RootChain.CallOpts)
}

// TotalWithdrawBalance is a free data retrieval call binding the contract method 0xc430c438.
//
// Solidity: function totalWithdrawBalance() constant returns(uint256)
func (_RootChain *RootChainCaller) TotalWithdrawBalance(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _RootChain.contract.Call(opts, out, "totalWithdrawBalance")
	return *ret0, err
}

// TotalWithdrawBalance is a free data retrieval call binding the contract method 0xc430c438.
//
// Solidity: function totalWithdrawBalance() constant returns(uint256)
func (_RootChain *RootChainSession) TotalWithdrawBalance() (*big.Int, error) {
	return _RootChain.Contract.TotalWithdrawBalance(&_RootChain.CallOpts)
}

// TotalWithdrawBalance is a free data retrieval call binding the contract method 0xc430c438.
//
// Solidity: function totalWithdrawBalance() constant returns(uint256)
func (_RootChain *RootChainCallerSession) TotalWithdrawBalance() (*big.Int, error) {
	return _RootChain.Contract.TotalWithdrawBalance(&_RootChain.CallOpts)
}

// TxExits is a free data retrieval call binding the contract method 0x6d3d8b1a.
//
// Solidity: function txExits( uint256) constant returns(amount uint256, createdAt uint256, owner address, state uint8)
func (_RootChain *RootChainCaller) TxExits(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Amount    *big.Int
	CreatedAt *big.Int
	Owner     common.Address
	State     uint8
}, error) {
	ret := new(struct {
		Amount    *big.Int
		CreatedAt *big.Int
		Owner     common.Address
		State     uint8
	})
	out := ret
	err := _RootChain.contract.Call(opts, out, "txExits", arg0)
	return *ret, err
}

// TxExits is a free data retrieval call binding the contract method 0x6d3d8b1a.
//
// Solidity: function txExits( uint256) constant returns(amount uint256, createdAt uint256, owner address, state uint8)
func (_RootChain *RootChainSession) TxExits(arg0 *big.Int) (struct {
	Amount    *big.Int
	CreatedAt *big.Int
	Owner     common.Address
	State     uint8
}, error) {
	return _RootChain.Contract.TxExits(&_RootChain.CallOpts, arg0)
}

// TxExits is a free data retrieval call binding the contract method 0x6d3d8b1a.
//
// Solidity: function txExits( uint256) constant returns(amount uint256, createdAt uint256, owner address, state uint8)
func (_RootChain *RootChainCallerSession) TxExits(arg0 *big.Int) (struct {
	Amount    *big.Int
	CreatedAt *big.Int
	Owner     common.Address
	State     uint8
}, error) {
	return _RootChain.Contract.TxExits(&_RootChain.CallOpts, arg0)
}

// TxIndexFactor is a free data retrieval call binding the contract method 0x00d2980a.
//
// Solidity: function txIndexFactor() constant returns(uint256)
func (_RootChain *RootChainCaller) TxIndexFactor(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _RootChain.contract.Call(opts, out, "txIndexFactor")
	return *ret0, err
}

// TxIndexFactor is a free data retrieval call binding the contract method 0x00d2980a.
//
// Solidity: function txIndexFactor() constant returns(uint256)
func (_RootChain *RootChainSession) TxIndexFactor() (*big.Int, error) {
	return _RootChain.Contract.TxIndexFactor(&_RootChain.CallOpts)
}

// TxIndexFactor is a free data retrieval call binding the contract method 0x00d2980a.
//
// Solidity: function txIndexFactor() constant returns(uint256)
func (_RootChain *RootChainCallerSession) TxIndexFactor() (*big.Int, error) {
	return _RootChain.Contract.TxIndexFactor(&_RootChain.CallOpts)
}

// ChallengeDepositExit is a paid mutator transaction binding the contract method 0x28d71e1f.
//
// Solidity: function challengeDepositExit(nonce uint256, newTxPos uint256[3], txBytes bytes, proof bytes, confirmSignature bytes) returns()
func (_RootChain *RootChainTransactor) ChallengeDepositExit(opts *bind.TransactOpts, nonce *big.Int, newTxPos [3]*big.Int, txBytes []byte, proof []byte, confirmSignature []byte) (*types.Transaction, error) {
	return _RootChain.contract.Transact(opts, "challengeDepositExit", nonce, newTxPos, txBytes, proof, confirmSignature)
}

// ChallengeDepositExit is a paid mutator transaction binding the contract method 0x28d71e1f.
//
// Solidity: function challengeDepositExit(nonce uint256, newTxPos uint256[3], txBytes bytes, proof bytes, confirmSignature bytes) returns()
func (_RootChain *RootChainSession) ChallengeDepositExit(nonce *big.Int, newTxPos [3]*big.Int, txBytes []byte, proof []byte, confirmSignature []byte) (*types.Transaction, error) {
	return _RootChain.Contract.ChallengeDepositExit(&_RootChain.TransactOpts, nonce, newTxPos, txBytes, proof, confirmSignature)
}

// ChallengeDepositExit is a paid mutator transaction binding the contract method 0x28d71e1f.
//
// Solidity: function challengeDepositExit(nonce uint256, newTxPos uint256[3], txBytes bytes, proof bytes, confirmSignature bytes) returns()
func (_RootChain *RootChainTransactorSession) ChallengeDepositExit(nonce *big.Int, newTxPos [3]*big.Int, txBytes []byte, proof []byte, confirmSignature []byte) (*types.Transaction, error) {
	return _RootChain.Contract.ChallengeDepositExit(&_RootChain.TransactOpts, nonce, newTxPos, txBytes, proof, confirmSignature)
}

// ChallengeTransactionExit is a paid mutator transaction binding the contract method 0x5e3f945b.
//
// Solidity: function challengeTransactionExit(exitingTxPos uint256[3], challengingTxPos uint256[3], txBytes bytes, proof bytes, confirmSignature bytes) returns()
func (_RootChain *RootChainTransactor) ChallengeTransactionExit(opts *bind.TransactOpts, exitingTxPos [3]*big.Int, challengingTxPos [3]*big.Int, txBytes []byte, proof []byte, confirmSignature []byte) (*types.Transaction, error) {
	return _RootChain.contract.Transact(opts, "challengeTransactionExit", exitingTxPos, challengingTxPos, txBytes, proof, confirmSignature)
}

// ChallengeTransactionExit is a paid mutator transaction binding the contract method 0x5e3f945b.
//
// Solidity: function challengeTransactionExit(exitingTxPos uint256[3], challengingTxPos uint256[3], txBytes bytes, proof bytes, confirmSignature bytes) returns()
func (_RootChain *RootChainSession) ChallengeTransactionExit(exitingTxPos [3]*big.Int, challengingTxPos [3]*big.Int, txBytes []byte, proof []byte, confirmSignature []byte) (*types.Transaction, error) {
	return _RootChain.Contract.ChallengeTransactionExit(&_RootChain.TransactOpts, exitingTxPos, challengingTxPos, txBytes, proof, confirmSignature)
}

// ChallengeTransactionExit is a paid mutator transaction binding the contract method 0x5e3f945b.
//
// Solidity: function challengeTransactionExit(exitingTxPos uint256[3], challengingTxPos uint256[3], txBytes bytes, proof bytes, confirmSignature bytes) returns()
func (_RootChain *RootChainTransactorSession) ChallengeTransactionExit(exitingTxPos [3]*big.Int, challengingTxPos [3]*big.Int, txBytes []byte, proof []byte, confirmSignature []byte) (*types.Transaction, error) {
	return _RootChain.Contract.ChallengeTransactionExit(&_RootChain.TransactOpts, exitingTxPos, challengingTxPos, txBytes, proof, confirmSignature)
}

// Deposit is a paid mutator transaction binding the contract method 0xf340fa01.
//
// Solidity: function deposit(owner address) returns()
func (_RootChain *RootChainTransactor) Deposit(opts *bind.TransactOpts, owner common.Address) (*types.Transaction, error) {
	return _RootChain.contract.Transact(opts, "deposit", owner)
}

// Deposit is a paid mutator transaction binding the contract method 0xf340fa01.
//
// Solidity: function deposit(owner address) returns()
func (_RootChain *RootChainSession) Deposit(owner common.Address) (*types.Transaction, error) {
	return _RootChain.Contract.Deposit(&_RootChain.TransactOpts, owner)
}

// Deposit is a paid mutator transaction binding the contract method 0xf340fa01.
//
// Solidity: function deposit(owner address) returns()
func (_RootChain *RootChainTransactorSession) Deposit(owner common.Address) (*types.Transaction, error) {
	return _RootChain.Contract.Deposit(&_RootChain.TransactOpts, owner)
}

// FinalizeDepositExits is a paid mutator transaction binding the contract method 0xfcf5f9eb.
//
// Solidity: function finalizeDepositExits() returns()
func (_RootChain *RootChainTransactor) FinalizeDepositExits(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RootChain.contract.Transact(opts, "finalizeDepositExits")
}

// FinalizeDepositExits is a paid mutator transaction binding the contract method 0xfcf5f9eb.
//
// Solidity: function finalizeDepositExits() returns()
func (_RootChain *RootChainSession) FinalizeDepositExits() (*types.Transaction, error) {
	return _RootChain.Contract.FinalizeDepositExits(&_RootChain.TransactOpts)
}

// FinalizeDepositExits is a paid mutator transaction binding the contract method 0xfcf5f9eb.
//
// Solidity: function finalizeDepositExits() returns()
func (_RootChain *RootChainTransactorSession) FinalizeDepositExits() (*types.Transaction, error) {
	return _RootChain.Contract.FinalizeDepositExits(&_RootChain.TransactOpts)
}

// FinalizeTransactionExits is a paid mutator transaction binding the contract method 0x884fc7d6.
//
// Solidity: function finalizeTransactionExits() returns()
func (_RootChain *RootChainTransactor) FinalizeTransactionExits(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RootChain.contract.Transact(opts, "finalizeTransactionExits")
}

// FinalizeTransactionExits is a paid mutator transaction binding the contract method 0x884fc7d6.
//
// Solidity: function finalizeTransactionExits() returns()
func (_RootChain *RootChainSession) FinalizeTransactionExits() (*types.Transaction, error) {
	return _RootChain.Contract.FinalizeTransactionExits(&_RootChain.TransactOpts)
}

// FinalizeTransactionExits is a paid mutator transaction binding the contract method 0x884fc7d6.
//
// Solidity: function finalizeTransactionExits() returns()
func (_RootChain *RootChainTransactorSession) FinalizeTransactionExits() (*types.Transaction, error) {
	return _RootChain.Contract.FinalizeTransactionExits(&_RootChain.TransactOpts)
}

// StartDepositExit is a paid mutator transaction binding the contract method 0xf5239f64.
//
// Solidity: function startDepositExit(nonce uint256) returns()
func (_RootChain *RootChainTransactor) StartDepositExit(opts *bind.TransactOpts, nonce *big.Int) (*types.Transaction, error) {
	return _RootChain.contract.Transact(opts, "startDepositExit", nonce)
}

// StartDepositExit is a paid mutator transaction binding the contract method 0xf5239f64.
//
// Solidity: function startDepositExit(nonce uint256) returns()
func (_RootChain *RootChainSession) StartDepositExit(nonce *big.Int) (*types.Transaction, error) {
	return _RootChain.Contract.StartDepositExit(&_RootChain.TransactOpts, nonce)
}

// StartDepositExit is a paid mutator transaction binding the contract method 0xf5239f64.
//
// Solidity: function startDepositExit(nonce uint256) returns()
func (_RootChain *RootChainTransactorSession) StartDepositExit(nonce *big.Int) (*types.Transaction, error) {
	return _RootChain.Contract.StartDepositExit(&_RootChain.TransactOpts, nonce)
}

// StartFeeExit is a paid mutator transaction binding the contract method 0xae80d8c8.
//
// Solidity: function startFeeExit(blockNumber uint256) returns()
func (_RootChain *RootChainTransactor) StartFeeExit(opts *bind.TransactOpts, blockNumber *big.Int) (*types.Transaction, error) {
	return _RootChain.contract.Transact(opts, "startFeeExit", blockNumber)
}

// StartFeeExit is a paid mutator transaction binding the contract method 0xae80d8c8.
//
// Solidity: function startFeeExit(blockNumber uint256) returns()
func (_RootChain *RootChainSession) StartFeeExit(blockNumber *big.Int) (*types.Transaction, error) {
	return _RootChain.Contract.StartFeeExit(&_RootChain.TransactOpts, blockNumber)
}

// StartFeeExit is a paid mutator transaction binding the contract method 0xae80d8c8.
//
// Solidity: function startFeeExit(blockNumber uint256) returns()
func (_RootChain *RootChainTransactorSession) StartFeeExit(blockNumber *big.Int) (*types.Transaction, error) {
	return _RootChain.Contract.StartFeeExit(&_RootChain.TransactOpts, blockNumber)
}

// StartTransactionExit is a paid mutator transaction binding the contract method 0x6621cb3f.
//
// Solidity: function startTransactionExit(txPos uint256[3], txBytes bytes, proof bytes, confirmSignatures bytes) returns()
func (_RootChain *RootChainTransactor) StartTransactionExit(opts *bind.TransactOpts, txPos [3]*big.Int, txBytes []byte, proof []byte, confirmSignatures []byte) (*types.Transaction, error) {
	return _RootChain.contract.Transact(opts, "startTransactionExit", txPos, txBytes, proof, confirmSignatures)
}

// StartTransactionExit is a paid mutator transaction binding the contract method 0x6621cb3f.
//
// Solidity: function startTransactionExit(txPos uint256[3], txBytes bytes, proof bytes, confirmSignatures bytes) returns()
func (_RootChain *RootChainSession) StartTransactionExit(txPos [3]*big.Int, txBytes []byte, proof []byte, confirmSignatures []byte) (*types.Transaction, error) {
	return _RootChain.Contract.StartTransactionExit(&_RootChain.TransactOpts, txPos, txBytes, proof, confirmSignatures)
}

// StartTransactionExit is a paid mutator transaction binding the contract method 0x6621cb3f.
//
// Solidity: function startTransactionExit(txPos uint256[3], txBytes bytes, proof bytes, confirmSignatures bytes) returns()
func (_RootChain *RootChainTransactorSession) StartTransactionExit(txPos [3]*big.Int, txBytes []byte, proof []byte, confirmSignatures []byte) (*types.Transaction, error) {
	return _RootChain.Contract.StartTransactionExit(&_RootChain.TransactOpts, txPos, txBytes, proof, confirmSignatures)
}

// SubmitBlock is a paid mutator transaction binding the contract method 0xb42f000a.
//
// Solidity: function submitBlock(blocks bytes, txnsPerBlock uint256[], feesPerBlock uint256[], blockNum uint256) returns()
func (_RootChain *RootChainTransactor) SubmitBlock(opts *bind.TransactOpts, blocks []byte, txnsPerBlock []*big.Int, feesPerBlock []*big.Int, blockNum *big.Int) (*types.Transaction, error) {
	return _RootChain.contract.Transact(opts, "submitBlock", blocks, txnsPerBlock, feesPerBlock, blockNum)
}

// SubmitBlock is a paid mutator transaction binding the contract method 0xb42f000a.
//
// Solidity: function submitBlock(blocks bytes, txnsPerBlock uint256[], feesPerBlock uint256[], blockNum uint256) returns()
func (_RootChain *RootChainSession) SubmitBlock(blocks []byte, txnsPerBlock []*big.Int, feesPerBlock []*big.Int, blockNum *big.Int) (*types.Transaction, error) {
	return _RootChain.Contract.SubmitBlock(&_RootChain.TransactOpts, blocks, txnsPerBlock, feesPerBlock, blockNum)
}

// SubmitBlock is a paid mutator transaction binding the contract method 0xb42f000a.
//
// Solidity: function submitBlock(blocks bytes, txnsPerBlock uint256[], feesPerBlock uint256[], blockNum uint256) returns()
func (_RootChain *RootChainTransactorSession) SubmitBlock(blocks []byte, txnsPerBlock []*big.Int, feesPerBlock []*big.Int, blockNum *big.Int) (*types.Transaction, error) {
	return _RootChain.Contract.SubmitBlock(&_RootChain.TransactOpts, blocks, txnsPerBlock, feesPerBlock, blockNum)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(newOwner address) returns()
func (_RootChain *RootChainTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _RootChain.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(newOwner address) returns()
func (_RootChain *RootChainSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _RootChain.Contract.TransferOwnership(&_RootChain.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(newOwner address) returns()
func (_RootChain *RootChainTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _RootChain.Contract.TransferOwnership(&_RootChain.TransactOpts, newOwner)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns(uint256)
func (_RootChain *RootChainTransactor) Withdraw(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _RootChain.contract.Transact(opts, "withdraw")
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns(uint256)
func (_RootChain *RootChainSession) Withdraw() (*types.Transaction, error) {
	return _RootChain.Contract.Withdraw(&_RootChain.TransactOpts)
}

// Withdraw is a paid mutator transaction binding the contract method 0x3ccfd60b.
//
// Solidity: function withdraw() returns(uint256)
func (_RootChain *RootChainTransactorSession) Withdraw() (*types.Transaction, error) {
	return _RootChain.Contract.Withdraw(&_RootChain.TransactOpts)
}

// RootChainAddedToBalancesIterator is returned from FilterAddedToBalances and is used to iterate over the raw logs and unpacked data for AddedToBalances events raised by the RootChain contract.
type RootChainAddedToBalancesIterator struct {
	Event *RootChainAddedToBalances // Event containing the contract specifics and raw log

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
func (it *RootChainAddedToBalancesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootChainAddedToBalances)
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
		it.Event = new(RootChainAddedToBalances)
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
func (it *RootChainAddedToBalancesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootChainAddedToBalancesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootChainAddedToBalances represents a AddedToBalances event raised by the RootChain contract.
type RootChainAddedToBalances struct {
	Owner  common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterAddedToBalances is a free log retrieval operation binding the contract event 0xf8552a24c7d58fd05114f6fc9db7b3a354db64d5fc758184af1696ccd8f158f3.
//
// Solidity: e AddedToBalances(owner address, amount uint256)
func (_RootChain *RootChainFilterer) FilterAddedToBalances(opts *bind.FilterOpts) (*RootChainAddedToBalancesIterator, error) {

	logs, sub, err := _RootChain.contract.FilterLogs(opts, "AddedToBalances")
	if err != nil {
		return nil, err
	}
	return &RootChainAddedToBalancesIterator{contract: _RootChain.contract, event: "AddedToBalances", logs: logs, sub: sub}, nil
}

// WatchAddedToBalances is a free log subscription operation binding the contract event 0xf8552a24c7d58fd05114f6fc9db7b3a354db64d5fc758184af1696ccd8f158f3.
//
// Solidity: e AddedToBalances(owner address, amount uint256)
func (_RootChain *RootChainFilterer) WatchAddedToBalances(opts *bind.WatchOpts, sink chan<- *RootChainAddedToBalances) (event.Subscription, error) {

	logs, sub, err := _RootChain.contract.WatchLogs(opts, "AddedToBalances")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootChainAddedToBalances)
				if err := _RootChain.contract.UnpackLog(event, "AddedToBalances", log); err != nil {
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

// RootChainBlockSubmittedIterator is returned from FilterBlockSubmitted and is used to iterate over the raw logs and unpacked data for BlockSubmitted events raised by the RootChain contract.
type RootChainBlockSubmittedIterator struct {
	Event *RootChainBlockSubmitted // Event containing the contract specifics and raw log

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
func (it *RootChainBlockSubmittedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootChainBlockSubmitted)
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
		it.Event = new(RootChainBlockSubmitted)
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
func (it *RootChainBlockSubmittedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootChainBlockSubmittedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootChainBlockSubmitted represents a BlockSubmitted event raised by the RootChain contract.
type RootChainBlockSubmitted struct {
	Root        [32]byte
	BlockNumber *big.Int
	NumTxns     *big.Int
	FeeAmount   *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterBlockSubmitted is a free log retrieval operation binding the contract event 0x044ff3798f9b3ad55d1155cea9a40508c71b4c64335f5dae87e8e11551515a06.
//
// Solidity: e BlockSubmitted(root bytes32, blockNumber uint256, numTxns uint256, feeAmount uint256)
func (_RootChain *RootChainFilterer) FilterBlockSubmitted(opts *bind.FilterOpts) (*RootChainBlockSubmittedIterator, error) {

	logs, sub, err := _RootChain.contract.FilterLogs(opts, "BlockSubmitted")
	if err != nil {
		return nil, err
	}
	return &RootChainBlockSubmittedIterator{contract: _RootChain.contract, event: "BlockSubmitted", logs: logs, sub: sub}, nil
}

// WatchBlockSubmitted is a free log subscription operation binding the contract event 0x044ff3798f9b3ad55d1155cea9a40508c71b4c64335f5dae87e8e11551515a06.
//
// Solidity: e BlockSubmitted(root bytes32, blockNumber uint256, numTxns uint256, feeAmount uint256)
func (_RootChain *RootChainFilterer) WatchBlockSubmitted(opts *bind.WatchOpts, sink chan<- *RootChainBlockSubmitted) (event.Subscription, error) {

	logs, sub, err := _RootChain.contract.WatchLogs(opts, "BlockSubmitted")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootChainBlockSubmitted)
				if err := _RootChain.contract.UnpackLog(event, "BlockSubmitted", log); err != nil {
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

// RootChainChallengedExitIterator is returned from FilterChallengedExit and is used to iterate over the raw logs and unpacked data for ChallengedExit events raised by the RootChain contract.
type RootChainChallengedExitIterator struct {
	Event *RootChainChallengedExit // Event containing the contract specifics and raw log

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
func (it *RootChainChallengedExitIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootChainChallengedExit)
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
		it.Event = new(RootChainChallengedExit)
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
func (it *RootChainChallengedExitIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootChainChallengedExitIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootChainChallengedExit represents a ChallengedExit event raised by the RootChain contract.
type RootChainChallengedExit struct {
	Position [4]*big.Int
	Owner    common.Address
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterChallengedExit is a free log retrieval operation binding the contract event 0xe1289dafb1083e540206bcd7d95a9705ba2590d6a9229c35a1c4c4c5efbda901.
//
// Solidity: e ChallengedExit(position uint256[4], owner address, amount uint256)
func (_RootChain *RootChainFilterer) FilterChallengedExit(opts *bind.FilterOpts) (*RootChainChallengedExitIterator, error) {

	logs, sub, err := _RootChain.contract.FilterLogs(opts, "ChallengedExit")
	if err != nil {
		return nil, err
	}
	return &RootChainChallengedExitIterator{contract: _RootChain.contract, event: "ChallengedExit", logs: logs, sub: sub}, nil
}

// WatchChallengedExit is a free log subscription operation binding the contract event 0xe1289dafb1083e540206bcd7d95a9705ba2590d6a9229c35a1c4c4c5efbda901.
//
// Solidity: e ChallengedExit(position uint256[4], owner address, amount uint256)
func (_RootChain *RootChainFilterer) WatchChallengedExit(opts *bind.WatchOpts, sink chan<- *RootChainChallengedExit) (event.Subscription, error) {

	logs, sub, err := _RootChain.contract.WatchLogs(opts, "ChallengedExit")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootChainChallengedExit)
				if err := _RootChain.contract.UnpackLog(event, "ChallengedExit", log); err != nil {
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

// RootChainDepositIterator is returned from FilterDeposit and is used to iterate over the raw logs and unpacked data for Deposit events raised by the RootChain contract.
type RootChainDepositIterator struct {
	Event *RootChainDeposit // Event containing the contract specifics and raw log

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
func (it *RootChainDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootChainDeposit)
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
		it.Event = new(RootChainDeposit)
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
func (it *RootChainDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootChainDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootChainDeposit represents a Deposit event raised by the RootChain contract.
type RootChainDeposit struct {
	Depositor    common.Address
	Amount       *big.Int
	DepositNonce *big.Int
	EthBlockNum  *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterDeposit is a free log retrieval operation binding the contract event 0x36af321ec8d3c75236829c5317affd40ddb308863a1236d2d277a4025cccee1e.
//
// Solidity: e Deposit(depositor address, amount uint256, depositNonce uint256, ethBlockNum uint256)
func (_RootChain *RootChainFilterer) FilterDeposit(opts *bind.FilterOpts) (*RootChainDepositIterator, error) {

	logs, sub, err := _RootChain.contract.FilterLogs(opts, "Deposit")
	if err != nil {
		return nil, err
	}
	return &RootChainDepositIterator{contract: _RootChain.contract, event: "Deposit", logs: logs, sub: sub}, nil
}

// WatchDeposit is a free log subscription operation binding the contract event 0x36af321ec8d3c75236829c5317affd40ddb308863a1236d2d277a4025cccee1e.
//
// Solidity: e Deposit(depositor address, amount uint256, depositNonce uint256, ethBlockNum uint256)
func (_RootChain *RootChainFilterer) WatchDeposit(opts *bind.WatchOpts, sink chan<- *RootChainDeposit) (event.Subscription, error) {

	logs, sub, err := _RootChain.contract.WatchLogs(opts, "Deposit")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootChainDeposit)
				if err := _RootChain.contract.UnpackLog(event, "Deposit", log); err != nil {
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

// RootChainFinalizedExitIterator is returned from FilterFinalizedExit and is used to iterate over the raw logs and unpacked data for FinalizedExit events raised by the RootChain contract.
type RootChainFinalizedExitIterator struct {
	Event *RootChainFinalizedExit // Event containing the contract specifics and raw log

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
func (it *RootChainFinalizedExitIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootChainFinalizedExit)
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
		it.Event = new(RootChainFinalizedExit)
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
func (it *RootChainFinalizedExitIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootChainFinalizedExitIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootChainFinalizedExit represents a FinalizedExit event raised by the RootChain contract.
type RootChainFinalizedExit struct {
	Position [4]*big.Int
	Owner    common.Address
	Amount   *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterFinalizedExit is a free log retrieval operation binding the contract event 0xb5083a27a38f8a9aa999efb3306b7be96dc3f42010a968dd86627880ba7fdbe2.
//
// Solidity: e FinalizedExit(position uint256[4], owner address, amount uint256)
func (_RootChain *RootChainFilterer) FilterFinalizedExit(opts *bind.FilterOpts) (*RootChainFinalizedExitIterator, error) {

	logs, sub, err := _RootChain.contract.FilterLogs(opts, "FinalizedExit")
	if err != nil {
		return nil, err
	}
	return &RootChainFinalizedExitIterator{contract: _RootChain.contract, event: "FinalizedExit", logs: logs, sub: sub}, nil
}

// WatchFinalizedExit is a free log subscription operation binding the contract event 0xb5083a27a38f8a9aa999efb3306b7be96dc3f42010a968dd86627880ba7fdbe2.
//
// Solidity: e FinalizedExit(position uint256[4], owner address, amount uint256)
func (_RootChain *RootChainFilterer) WatchFinalizedExit(opts *bind.WatchOpts, sink chan<- *RootChainFinalizedExit) (event.Subscription, error) {

	logs, sub, err := _RootChain.contract.WatchLogs(opts, "FinalizedExit")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootChainFinalizedExit)
				if err := _RootChain.contract.UnpackLog(event, "FinalizedExit", log); err != nil {
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

// RootChainOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the RootChain contract.
type RootChainOwnershipTransferredIterator struct {
	Event *RootChainOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *RootChainOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootChainOwnershipTransferred)
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
		it.Event = new(RootChainOwnershipTransferred)
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
func (it *RootChainOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootChainOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootChainOwnershipTransferred represents a OwnershipTransferred event raised by the RootChain contract.
type RootChainOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: e OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_RootChain *RootChainFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*RootChainOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _RootChain.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &RootChainOwnershipTransferredIterator{contract: _RootChain.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: e OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_RootChain *RootChainFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *RootChainOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _RootChain.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootChainOwnershipTransferred)
				if err := _RootChain.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// RootChainStartedDepositExitIterator is returned from FilterStartedDepositExit and is used to iterate over the raw logs and unpacked data for StartedDepositExit events raised by the RootChain contract.
type RootChainStartedDepositExitIterator struct {
	Event *RootChainStartedDepositExit // Event containing the contract specifics and raw log

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
func (it *RootChainStartedDepositExitIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootChainStartedDepositExit)
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
		it.Event = new(RootChainStartedDepositExit)
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
func (it *RootChainStartedDepositExitIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootChainStartedDepositExitIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootChainStartedDepositExit represents a StartedDepositExit event raised by the RootChain contract.
type RootChainStartedDepositExit struct {
	Nonce  *big.Int
	Owner  common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterStartedDepositExit is a free log retrieval operation binding the contract event 0x0bdfdd54dc0a51ef460d31ddf95470493780afed2eee6046199b65c2b1d66b91.
//
// Solidity: e StartedDepositExit(nonce uint256, owner address, amount uint256)
func (_RootChain *RootChainFilterer) FilterStartedDepositExit(opts *bind.FilterOpts) (*RootChainStartedDepositExitIterator, error) {

	logs, sub, err := _RootChain.contract.FilterLogs(opts, "StartedDepositExit")
	if err != nil {
		return nil, err
	}
	return &RootChainStartedDepositExitIterator{contract: _RootChain.contract, event: "StartedDepositExit", logs: logs, sub: sub}, nil
}

// WatchStartedDepositExit is a free log subscription operation binding the contract event 0x0bdfdd54dc0a51ef460d31ddf95470493780afed2eee6046199b65c2b1d66b91.
//
// Solidity: e StartedDepositExit(nonce uint256, owner address, amount uint256)
func (_RootChain *RootChainFilterer) WatchStartedDepositExit(opts *bind.WatchOpts, sink chan<- *RootChainStartedDepositExit) (event.Subscription, error) {

	logs, sub, err := _RootChain.contract.WatchLogs(opts, "StartedDepositExit")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootChainStartedDepositExit)
				if err := _RootChain.contract.UnpackLog(event, "StartedDepositExit", log); err != nil {
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

// RootChainStartedTransactionExitIterator is returned from FilterStartedTransactionExit and is used to iterate over the raw logs and unpacked data for StartedTransactionExit events raised by the RootChain contract.
type RootChainStartedTransactionExitIterator struct {
	Event *RootChainStartedTransactionExit // Event containing the contract specifics and raw log

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
func (it *RootChainStartedTransactionExitIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootChainStartedTransactionExit)
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
		it.Event = new(RootChainStartedTransactionExit)
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
func (it *RootChainStartedTransactionExitIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootChainStartedTransactionExitIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootChainStartedTransactionExit represents a StartedTransactionExit event raised by the RootChain contract.
type RootChainStartedTransactionExit struct {
	Position          [3]*big.Int
	Owner             common.Address
	Amount            *big.Int
	ConfirmSignatures []byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterStartedTransactionExit is a free log retrieval operation binding the contract event 0xb9a20182985ca044806e684b82bcbff288e57bdfffdc1c1148e8de595665eb2c.
//
// Solidity: e StartedTransactionExit(position uint256[3], owner address, amount uint256, confirmSignatures bytes)
func (_RootChain *RootChainFilterer) FilterStartedTransactionExit(opts *bind.FilterOpts) (*RootChainStartedTransactionExitIterator, error) {

	logs, sub, err := _RootChain.contract.FilterLogs(opts, "StartedTransactionExit")
	if err != nil {
		return nil, err
	}
	return &RootChainStartedTransactionExitIterator{contract: _RootChain.contract, event: "StartedTransactionExit", logs: logs, sub: sub}, nil
}

// WatchStartedTransactionExit is a free log subscription operation binding the contract event 0xb9a20182985ca044806e684b82bcbff288e57bdfffdc1c1148e8de595665eb2c.
//
// Solidity: e StartedTransactionExit(position uint256[3], owner address, amount uint256, confirmSignatures bytes)
func (_RootChain *RootChainFilterer) WatchStartedTransactionExit(opts *bind.WatchOpts, sink chan<- *RootChainStartedTransactionExit) (event.Subscription, error) {

	logs, sub, err := _RootChain.contract.WatchLogs(opts, "StartedTransactionExit")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootChainStartedTransactionExit)
				if err := _RootChain.contract.UnpackLog(event, "StartedTransactionExit", log); err != nil {
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
