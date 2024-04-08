// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package Log0

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

// Log0MetaData contains all meta data concerning the Log0 contract.
var Log0MetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"opLog0\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"opLog00\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"opLog01\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
	Bin: "0x6080604052348015600f57600080fd5b5060938061001e6000396000f3fe60806040526004361060305760003560e01c80633e2d0b8514603557806357e4605514603d578063ecc5544a146043575b600080fd5b603b6049565b005b603b6050565b603b6056565b601c6000a0565b600080a0565b60206000a056fea26469706673582212209aba01a729d89e6da96ac8ca0b8f1940565356ed4f7849c9af7a95f5188d22d964736f6c634300080c0033",
}

// Log0ABI is the input ABI used to generate the binding from.
// Deprecated: Use Log0MetaData.ABI instead.
var Log0ABI = Log0MetaData.ABI

// Log0Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use Log0MetaData.Bin instead.
var Log0Bin = Log0MetaData.Bin

// DeployLog0 deploys a new Ethereum contract, binding an instance of Log0 to it.
func DeployLog0(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Log0, error) {
	parsed, err := Log0MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(Log0Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Log0{Log0Caller: Log0Caller{contract: contract}, Log0Transactor: Log0Transactor{contract: contract}, Log0Filterer: Log0Filterer{contract: contract}}, nil
}

// Log0 is an auto generated Go binding around an Ethereum contract.
type Log0 struct {
	Log0Caller     // Read-only binding to the contract
	Log0Transactor // Write-only binding to the contract
	Log0Filterer   // Log filterer for contract events
}

// Log0Caller is an auto generated read-only Go binding around an Ethereum contract.
type Log0Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Log0Transactor is an auto generated write-only Go binding around an Ethereum contract.
type Log0Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Log0Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type Log0Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Log0Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type Log0Session struct {
	Contract     *Log0             // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Log0CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type Log0CallerSession struct {
	Contract *Log0Caller   // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// Log0TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type Log0TransactorSession struct {
	Contract     *Log0Transactor   // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Log0Raw is an auto generated low-level Go binding around an Ethereum contract.
type Log0Raw struct {
	Contract *Log0 // Generic contract binding to access the raw methods on
}

// Log0CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type Log0CallerRaw struct {
	Contract *Log0Caller // Generic read-only contract binding to access the raw methods on
}

// Log0TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type Log0TransactorRaw struct {
	Contract *Log0Transactor // Generic write-only contract binding to access the raw methods on
}

// NewLog0 creates a new instance of Log0, bound to a specific deployed contract.
func NewLog0(address common.Address, backend bind.ContractBackend) (*Log0, error) {
	contract, err := bindLog0(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Log0{Log0Caller: Log0Caller{contract: contract}, Log0Transactor: Log0Transactor{contract: contract}, Log0Filterer: Log0Filterer{contract: contract}}, nil
}

// NewLog0Caller creates a new read-only instance of Log0, bound to a specific deployed contract.
func NewLog0Caller(address common.Address, caller bind.ContractCaller) (*Log0Caller, error) {
	contract, err := bindLog0(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &Log0Caller{contract: contract}, nil
}

// NewLog0Transactor creates a new write-only instance of Log0, bound to a specific deployed contract.
func NewLog0Transactor(address common.Address, transactor bind.ContractTransactor) (*Log0Transactor, error) {
	contract, err := bindLog0(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &Log0Transactor{contract: contract}, nil
}

// NewLog0Filterer creates a new log filterer instance of Log0, bound to a specific deployed contract.
func NewLog0Filterer(address common.Address, filterer bind.ContractFilterer) (*Log0Filterer, error) {
	contract, err := bindLog0(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &Log0Filterer{contract: contract}, nil
}

// bindLog0 binds a generic wrapper to an already deployed contract.
func bindLog0(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := Log0MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Log0 *Log0Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Log0.Contract.Log0Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Log0 *Log0Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Log0.Contract.Log0Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Log0 *Log0Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Log0.Contract.Log0Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Log0 *Log0CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Log0.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Log0 *Log0TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Log0.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Log0 *Log0TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Log0.Contract.contract.Transact(opts, method, params...)
}

// OpLog0 is a paid mutator transaction binding the contract method 0xecc5544a.
//
// Solidity: function opLog0() payable returns()
func (_Log0 *Log0Transactor) OpLog0(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Log0.contract.Transact(opts, "opLog0")
}

// OpLog0 is a paid mutator transaction binding the contract method 0xecc5544a.
//
// Solidity: function opLog0() payable returns()
func (_Log0 *Log0Session) OpLog0() (*types.Transaction, error) {
	return _Log0.Contract.OpLog0(&_Log0.TransactOpts)
}

// OpLog0 is a paid mutator transaction binding the contract method 0xecc5544a.
//
// Solidity: function opLog0() payable returns()
func (_Log0 *Log0TransactorSession) OpLog0() (*types.Transaction, error) {
	return _Log0.Contract.OpLog0(&_Log0.TransactOpts)
}

// OpLog00 is a paid mutator transaction binding the contract method 0x57e46055.
//
// Solidity: function opLog00() payable returns()
func (_Log0 *Log0Transactor) OpLog00(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Log0.contract.Transact(opts, "opLog00")
}

// OpLog00 is a paid mutator transaction binding the contract method 0x57e46055.
//
// Solidity: function opLog00() payable returns()
func (_Log0 *Log0Session) OpLog00() (*types.Transaction, error) {
	return _Log0.Contract.OpLog00(&_Log0.TransactOpts)
}

// OpLog00 is a paid mutator transaction binding the contract method 0x57e46055.
//
// Solidity: function opLog00() payable returns()
func (_Log0 *Log0TransactorSession) OpLog00() (*types.Transaction, error) {
	return _Log0.Contract.OpLog00(&_Log0.TransactOpts)
}

// OpLog01 is a paid mutator transaction binding the contract method 0x3e2d0b85.
//
// Solidity: function opLog01() payable returns()
func (_Log0 *Log0Transactor) OpLog01(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Log0.contract.Transact(opts, "opLog01")
}

// OpLog01 is a paid mutator transaction binding the contract method 0x3e2d0b85.
//
// Solidity: function opLog01() payable returns()
func (_Log0 *Log0Session) OpLog01() (*types.Transaction, error) {
	return _Log0.Contract.OpLog01(&_Log0.TransactOpts)
}

// OpLog01 is a paid mutator transaction binding the contract method 0x3e2d0b85.
//
// Solidity: function opLog01() payable returns()
func (_Log0 *Log0TransactorSession) OpLog01() (*types.Transaction, error) {
	return _Log0.Contract.OpLog01(&_Log0.TransactOpts)
}
