// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package CounterAndBlock

import (
	"errors"
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
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// CounterAndBlockMetaData contains all meta data concerning the CounterAndBlock contract.
var CounterAndBlockMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"count\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"increment\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5060eb8061001f6000396000f3fe6080604052348015600f57600080fd5b5060043610603c5760003560e01c806306661abd146041578063a87d942c14605c578063d09de08a146071575b600080fd5b604960005481565b6040519081526020015b60405180910390f35b60005460408051918252426020830152016053565b60776079565b005b6001600080828254608991906090565b9091555050565b6000821982111560b057634e487b7160e01b600052601160045260246000fd5b50019056fea26469706673582212205aa9aebefdfb857d27d7bdc8475c08138617cc37e78c2e6bd98acb9a1484994964736f6c634300080c0033",
}

// CounterAndBlockABI is the input ABI used to generate the binding from.
// Deprecated: Use CounterAndBlockMetaData.ABI instead.
var CounterAndBlockABI = CounterAndBlockMetaData.ABI

// CounterAndBlockBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use CounterAndBlockMetaData.Bin instead.
var CounterAndBlockBin = CounterAndBlockMetaData.Bin

// DeployCounterAndBlock deploys a new Ethereum contract, binding an instance of CounterAndBlock to it.
func DeployCounterAndBlock(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *CounterAndBlock, error) {
	parsed, err := CounterAndBlockMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(CounterAndBlockBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &CounterAndBlock{CounterAndBlockCaller: CounterAndBlockCaller{contract: contract}, CounterAndBlockTransactor: CounterAndBlockTransactor{contract: contract}, CounterAndBlockFilterer: CounterAndBlockFilterer{contract: contract}}, nil
}

// CounterAndBlock is an auto generated Go binding around an Ethereum contract.
type CounterAndBlock struct {
	CounterAndBlockCaller     // Read-only binding to the contract
	CounterAndBlockTransactor // Write-only binding to the contract
	CounterAndBlockFilterer   // Log filterer for contract events
}

// CounterAndBlockCaller is an auto generated read-only Go binding around an Ethereum contract.
type CounterAndBlockCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CounterAndBlockTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CounterAndBlockTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CounterAndBlockFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CounterAndBlockFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CounterAndBlockSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CounterAndBlockSession struct {
	Contract     *CounterAndBlock  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CounterAndBlockCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CounterAndBlockCallerSession struct {
	Contract *CounterAndBlockCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// CounterAndBlockTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CounterAndBlockTransactorSession struct {
	Contract     *CounterAndBlockTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// CounterAndBlockRaw is an auto generated low-level Go binding around an Ethereum contract.
type CounterAndBlockRaw struct {
	Contract *CounterAndBlock // Generic contract binding to access the raw methods on
}

// CounterAndBlockCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CounterAndBlockCallerRaw struct {
	Contract *CounterAndBlockCaller // Generic read-only contract binding to access the raw methods on
}

// CounterAndBlockTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CounterAndBlockTransactorRaw struct {
	Contract *CounterAndBlockTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCounterAndBlock creates a new instance of CounterAndBlock, bound to a specific deployed contract.
func NewCounterAndBlock(address common.Address, backend bind.ContractBackend) (*CounterAndBlock, error) {
	contract, err := bindCounterAndBlock(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CounterAndBlock{CounterAndBlockCaller: CounterAndBlockCaller{contract: contract}, CounterAndBlockTransactor: CounterAndBlockTransactor{contract: contract}, CounterAndBlockFilterer: CounterAndBlockFilterer{contract: contract}}, nil
}

// NewCounterAndBlockCaller creates a new read-only instance of CounterAndBlock, bound to a specific deployed contract.
func NewCounterAndBlockCaller(address common.Address, caller bind.ContractCaller) (*CounterAndBlockCaller, error) {
	contract, err := bindCounterAndBlock(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CounterAndBlockCaller{contract: contract}, nil
}

// NewCounterAndBlockTransactor creates a new write-only instance of CounterAndBlock, bound to a specific deployed contract.
func NewCounterAndBlockTransactor(address common.Address, transactor bind.ContractTransactor) (*CounterAndBlockTransactor, error) {
	contract, err := bindCounterAndBlock(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CounterAndBlockTransactor{contract: contract}, nil
}

// NewCounterAndBlockFilterer creates a new log filterer instance of CounterAndBlock, bound to a specific deployed contract.
func NewCounterAndBlockFilterer(address common.Address, filterer bind.ContractFilterer) (*CounterAndBlockFilterer, error) {
	contract, err := bindCounterAndBlock(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CounterAndBlockFilterer{contract: contract}, nil
}

// bindCounterAndBlock binds a generic wrapper to an already deployed contract.
func bindCounterAndBlock(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CounterAndBlockMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CounterAndBlock *CounterAndBlockRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CounterAndBlock.Contract.CounterAndBlockCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CounterAndBlock *CounterAndBlockRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CounterAndBlock.Contract.CounterAndBlockTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CounterAndBlock *CounterAndBlockRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CounterAndBlock.Contract.CounterAndBlockTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CounterAndBlock *CounterAndBlockCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CounterAndBlock.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CounterAndBlock *CounterAndBlockTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CounterAndBlock.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CounterAndBlock *CounterAndBlockTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CounterAndBlock.Contract.contract.Transact(opts, method, params...)
}

// Count is a free data retrieval call binding the contract method 0x06661abd.
//
// Solidity: function count() view returns(uint256)
func (_CounterAndBlock *CounterAndBlockCaller) Count(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _CounterAndBlock.contract.Call(opts, &out, "count")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Count is a free data retrieval call binding the contract method 0x06661abd.
//
// Solidity: function count() view returns(uint256)
func (_CounterAndBlock *CounterAndBlockSession) Count() (*big.Int, error) {
	return _CounterAndBlock.Contract.Count(&_CounterAndBlock.CallOpts)
}

// Count is a free data retrieval call binding the contract method 0x06661abd.
//
// Solidity: function count() view returns(uint256)
func (_CounterAndBlock *CounterAndBlockCallerSession) Count() (*big.Int, error) {
	return _CounterAndBlock.Contract.Count(&_CounterAndBlock.CallOpts)
}

// GetCount is a free data retrieval call binding the contract method 0xa87d942c.
//
// Solidity: function getCount() view returns(uint256, uint256)
func (_CounterAndBlock *CounterAndBlockCaller) GetCount(opts *bind.CallOpts) (*big.Int, *big.Int, error) {
	var out []interface{}
	err := _CounterAndBlock.contract.Call(opts, &out, "getCount")

	if err != nil {
		return *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

// GetCount is a free data retrieval call binding the contract method 0xa87d942c.
//
// Solidity: function getCount() view returns(uint256, uint256)
func (_CounterAndBlock *CounterAndBlockSession) GetCount() (*big.Int, *big.Int, error) {
	return _CounterAndBlock.Contract.GetCount(&_CounterAndBlock.CallOpts)
}

// GetCount is a free data retrieval call binding the contract method 0xa87d942c.
//
// Solidity: function getCount() view returns(uint256, uint256)
func (_CounterAndBlock *CounterAndBlockCallerSession) GetCount() (*big.Int, *big.Int, error) {
	return _CounterAndBlock.Contract.GetCount(&_CounterAndBlock.CallOpts)
}

// Increment is a paid mutator transaction binding the contract method 0xd09de08a.
//
// Solidity: function increment() returns()
func (_CounterAndBlock *CounterAndBlockTransactor) Increment(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CounterAndBlock.contract.Transact(opts, "increment")
}

// Increment is a paid mutator transaction binding the contract method 0xd09de08a.
//
// Solidity: function increment() returns()
func (_CounterAndBlock *CounterAndBlockSession) Increment() (*types.Transaction, error) {
	return _CounterAndBlock.Contract.Increment(&_CounterAndBlock.TransactOpts)
}

// Increment is a paid mutator transaction binding the contract method 0xd09de08a.
//
// Solidity: function increment() returns()
func (_CounterAndBlock *CounterAndBlockTransactorSession) Increment() (*types.Transaction, error) {
	return _CounterAndBlock.Contract.Increment(&_CounterAndBlock.TransactOpts)
}
