// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package Sha

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

// ShaMetaData contains all meta data concerning the Sha contract.
var ShaMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"hash\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5060dd8061001f6000396000f3fe6080604052348015600f57600080fd5b506004361060285760003560e01c806309bd5a6014602d575b600080fd5b60336035565b005b6040516a1a195b1b1bc81ddbdc9b1960aa1b8152600290600b01602060405180830381855afa158015606b573d6000803e3d6000fd5b5050506040513d601f19601f82011682018060405250810190608c9190608f565b50565b60006020828403121560a057600080fd5b505191905056fea2646970667358221220fe50f9cdaf095e844b048d17690c0862e54c7d7b550c4e2d9045215198e571d564736f6c634300080c0033",
}

// ShaABI is the input ABI used to generate the binding from.
// Deprecated: Use ShaMetaData.ABI instead.
var ShaABI = ShaMetaData.ABI

// ShaBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ShaMetaData.Bin instead.
var ShaBin = ShaMetaData.Bin

// DeploySha deploys a new Ethereum contract, binding an instance of Sha to it.
func DeploySha(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Sha, error) {
	parsed, err := ShaMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ShaBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Sha{ShaCaller: ShaCaller{contract: contract}, ShaTransactor: ShaTransactor{contract: contract}, ShaFilterer: ShaFilterer{contract: contract}}, nil
}

// Sha is an auto generated Go binding around an Ethereum contract.
type Sha struct {
	ShaCaller     // Read-only binding to the contract
	ShaTransactor // Write-only binding to the contract
	ShaFilterer   // Log filterer for contract events
}

// ShaCaller is an auto generated read-only Go binding around an Ethereum contract.
type ShaCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ShaTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ShaTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ShaFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ShaFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ShaSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ShaSession struct {
	Contract     *Sha              // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ShaCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ShaCallerSession struct {
	Contract *ShaCaller    // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// ShaTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ShaTransactorSession struct {
	Contract     *ShaTransactor    // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ShaRaw is an auto generated low-level Go binding around an Ethereum contract.
type ShaRaw struct {
	Contract *Sha // Generic contract binding to access the raw methods on
}

// ShaCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ShaCallerRaw struct {
	Contract *ShaCaller // Generic read-only contract binding to access the raw methods on
}

// ShaTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ShaTransactorRaw struct {
	Contract *ShaTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSha creates a new instance of Sha, bound to a specific deployed contract.
func NewSha(address common.Address, backend bind.ContractBackend) (*Sha, error) {
	contract, err := bindSha(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Sha{ShaCaller: ShaCaller{contract: contract}, ShaTransactor: ShaTransactor{contract: contract}, ShaFilterer: ShaFilterer{contract: contract}}, nil
}

// NewShaCaller creates a new read-only instance of Sha, bound to a specific deployed contract.
func NewShaCaller(address common.Address, caller bind.ContractCaller) (*ShaCaller, error) {
	contract, err := bindSha(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ShaCaller{contract: contract}, nil
}

// NewShaTransactor creates a new write-only instance of Sha, bound to a specific deployed contract.
func NewShaTransactor(address common.Address, transactor bind.ContractTransactor) (*ShaTransactor, error) {
	contract, err := bindSha(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ShaTransactor{contract: contract}, nil
}

// NewShaFilterer creates a new log filterer instance of Sha, bound to a specific deployed contract.
func NewShaFilterer(address common.Address, filterer bind.ContractFilterer) (*ShaFilterer, error) {
	contract, err := bindSha(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ShaFilterer{contract: contract}, nil
}

// bindSha binds a generic wrapper to an already deployed contract.
func bindSha(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ShaMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Sha *ShaRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Sha.Contract.ShaCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Sha *ShaRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Sha.Contract.ShaTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Sha *ShaRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Sha.Contract.ShaTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Sha *ShaCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Sha.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Sha *ShaTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Sha.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Sha *ShaTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Sha.Contract.contract.Transact(opts, method, params...)
}

// Hash is a paid mutator transaction binding the contract method 0x09bd5a60.
//
// Solidity: function hash() returns()
func (_Sha *ShaTransactor) Hash(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Sha.contract.Transact(opts, "hash")
}

// Hash is a paid mutator transaction binding the contract method 0x09bd5a60.
//
// Solidity: function hash() returns()
func (_Sha *ShaSession) Hash() (*types.Transaction, error) {
	return _Sha.Contract.Hash(&_Sha.TransactOpts)
}

// Hash is a paid mutator transaction binding the contract method 0x09bd5a60.
//
// Solidity: function hash() returns()
func (_Sha *ShaTransactorSession) Hash() (*types.Transaction, error) {
	return _Sha.Contract.Hash(&_Sha.TransactOpts)
}
