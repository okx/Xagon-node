// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package customModExp

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

// CustomModExpMetaData contains all meta data concerning the CustomModExp contract.
var CustomModExpMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"input\",\"type\":\"bytes\"}],\"name\":\"modExpGeneric\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610208806100206000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c8063d5665d6f14610030575b600080fd5b61004361003e3660046100e2565b610045565b005b61004d6100ad565b6101408183516020850160055afa60009081555b600a8110156100a8578181600a811061007c5761007c610193565b6020020151600482600a811061009457610094610193565b0155806100a0816101a9565b915050610061565b505050565b604051806101400160405280600a906020820280368337509192915050565b634e487b7160e01b600052604160045260246000fd5b6000602082840312156100f457600080fd5b813567ffffffffffffffff8082111561010c57600080fd5b818401915084601f83011261012057600080fd5b813581811115610132576101326100cc565b604051601f8201601f19908116603f0116810190838211818310171561015a5761015a6100cc565b8160405282815287602084870101111561017357600080fd5b826020860160208301376000928101602001929092525095945050505050565b634e487b7160e01b600052603260045260246000fd5b60006000198214156101cb57634e487b7160e01b600052601160045260246000fd5b506001019056fea26469706673582212206c4940b4c9a7086754420734c8b4921cdb547ec8b31fc3bf8cd884ad9778a5b364736f6c634300080c0033",
}

// CustomModExpABI is the input ABI used to generate the binding from.
// Deprecated: Use CustomModExpMetaData.ABI instead.
var CustomModExpABI = CustomModExpMetaData.ABI

// CustomModExpBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use CustomModExpMetaData.Bin instead.
var CustomModExpBin = CustomModExpMetaData.Bin

// DeployCustomModExp deploys a new Ethereum contract, binding an instance of CustomModExp to it.
func DeployCustomModExp(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *CustomModExp, error) {
	parsed, err := CustomModExpMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(CustomModExpBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &CustomModExp{CustomModExpCaller: CustomModExpCaller{contract: contract}, CustomModExpTransactor: CustomModExpTransactor{contract: contract}, CustomModExpFilterer: CustomModExpFilterer{contract: contract}}, nil
}

// CustomModExp is an auto generated Go binding around an Ethereum contract.
type CustomModExp struct {
	CustomModExpCaller     // Read-only binding to the contract
	CustomModExpTransactor // Write-only binding to the contract
	CustomModExpFilterer   // Log filterer for contract events
}

// CustomModExpCaller is an auto generated read-only Go binding around an Ethereum contract.
type CustomModExpCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CustomModExpTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CustomModExpTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CustomModExpFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CustomModExpFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CustomModExpSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CustomModExpSession struct {
	Contract     *CustomModExp     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CustomModExpCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CustomModExpCallerSession struct {
	Contract *CustomModExpCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// CustomModExpTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CustomModExpTransactorSession struct {
	Contract     *CustomModExpTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// CustomModExpRaw is an auto generated low-level Go binding around an Ethereum contract.
type CustomModExpRaw struct {
	Contract *CustomModExp // Generic contract binding to access the raw methods on
}

// CustomModExpCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CustomModExpCallerRaw struct {
	Contract *CustomModExpCaller // Generic read-only contract binding to access the raw methods on
}

// CustomModExpTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CustomModExpTransactorRaw struct {
	Contract *CustomModExpTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCustomModExp creates a new instance of CustomModExp, bound to a specific deployed contract.
func NewCustomModExp(address common.Address, backend bind.ContractBackend) (*CustomModExp, error) {
	contract, err := bindCustomModExp(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &CustomModExp{CustomModExpCaller: CustomModExpCaller{contract: contract}, CustomModExpTransactor: CustomModExpTransactor{contract: contract}, CustomModExpFilterer: CustomModExpFilterer{contract: contract}}, nil
}

// NewCustomModExpCaller creates a new read-only instance of CustomModExp, bound to a specific deployed contract.
func NewCustomModExpCaller(address common.Address, caller bind.ContractCaller) (*CustomModExpCaller, error) {
	contract, err := bindCustomModExp(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CustomModExpCaller{contract: contract}, nil
}

// NewCustomModExpTransactor creates a new write-only instance of CustomModExp, bound to a specific deployed contract.
func NewCustomModExpTransactor(address common.Address, transactor bind.ContractTransactor) (*CustomModExpTransactor, error) {
	contract, err := bindCustomModExp(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CustomModExpTransactor{contract: contract}, nil
}

// NewCustomModExpFilterer creates a new log filterer instance of CustomModExp, bound to a specific deployed contract.
func NewCustomModExpFilterer(address common.Address, filterer bind.ContractFilterer) (*CustomModExpFilterer, error) {
	contract, err := bindCustomModExp(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CustomModExpFilterer{contract: contract}, nil
}

// bindCustomModExp binds a generic wrapper to an already deployed contract.
func bindCustomModExp(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CustomModExpMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CustomModExp *CustomModExpRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CustomModExp.Contract.CustomModExpCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CustomModExp *CustomModExpRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CustomModExp.Contract.CustomModExpTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CustomModExp *CustomModExpRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CustomModExp.Contract.CustomModExpTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_CustomModExp *CustomModExpCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _CustomModExp.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_CustomModExp *CustomModExpTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _CustomModExp.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_CustomModExp *CustomModExpTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _CustomModExp.Contract.contract.Transact(opts, method, params...)
}

// ModExpGeneric is a paid mutator transaction binding the contract method 0xd5665d6f.
//
// Solidity: function modExpGeneric(bytes input) returns()
func (_CustomModExp *CustomModExpTransactor) ModExpGeneric(opts *bind.TransactOpts, input []byte) (*types.Transaction, error) {
	return _CustomModExp.contract.Transact(opts, "modExpGeneric", input)
}

// ModExpGeneric is a paid mutator transaction binding the contract method 0xd5665d6f.
//
// Solidity: function modExpGeneric(bytes input) returns()
func (_CustomModExp *CustomModExpSession) ModExpGeneric(input []byte) (*types.Transaction, error) {
	return _CustomModExp.Contract.ModExpGeneric(&_CustomModExp.TransactOpts, input)
}

// ModExpGeneric is a paid mutator transaction binding the contract method 0xd5665d6f.
//
// Solidity: function modExpGeneric(bytes input) returns()
func (_CustomModExp *CustomModExpTransactorSession) ModExpGeneric(input []byte) (*types.Transaction, error) {
	return _CustomModExp.Contract.ModExpGeneric(&_CustomModExp.TransactOpts, input)
}
