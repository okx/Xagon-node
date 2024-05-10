// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package triggerErrors

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

// TriggerErrorsMetaData contains all meta data concerning the TriggerErrors contract.
var TriggerErrorsMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"count\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"outOfCountersKeccaks\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"test\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"outOfCountersPoseidon\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"outOfCountersSteps\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"outOfGas\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60806040526000805534801561001457600080fd5b5061016c806100246000396000f3fe608060405234801561001057600080fd5b50600436106100575760003560e01c806306661abd1461005c5780632621002a1461007757806331fe52e8146100835780638bd7b5381461008d578063cb4e8cd114610095575b600080fd5b61006560005481565b60405190815260200160405180910390f35b620f4240600020610065565b61008b61009d565b005b61008b6100c3565b61008b6100e9565b60005b60648110156100c0578060005580806100b89061010d565b9150506100a0565b50565b60005b620186a08110156100c0576104d2600052806100e18161010d565b9150506100c6565b60005b61c3508110156100c0578060005580806101059061010d565b9150506100ec565b600060001982141561012f57634e487b7160e01b600052601160045260246000fd5b506001019056fea26469706673582212208f01c5dc055b1f376f5da5deb33e2c96ee776174bf48874c5ebba0f606de2ac564736f6c634300080c0033",
}

// TriggerErrorsABI is the input ABI used to generate the binding from.
// Deprecated: Use TriggerErrorsMetaData.ABI instead.
var TriggerErrorsABI = TriggerErrorsMetaData.ABI

// TriggerErrorsBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use TriggerErrorsMetaData.Bin instead.
var TriggerErrorsBin = TriggerErrorsMetaData.Bin

// DeployTriggerErrors deploys a new Ethereum contract, binding an instance of TriggerErrors to it.
func DeployTriggerErrors(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *TriggerErrors, error) {
	parsed, err := TriggerErrorsMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(TriggerErrorsBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &TriggerErrors{TriggerErrorsCaller: TriggerErrorsCaller{contract: contract}, TriggerErrorsTransactor: TriggerErrorsTransactor{contract: contract}, TriggerErrorsFilterer: TriggerErrorsFilterer{contract: contract}}, nil
}

// TriggerErrors is an auto generated Go binding around an Ethereum contract.
type TriggerErrors struct {
	TriggerErrorsCaller     // Read-only binding to the contract
	TriggerErrorsTransactor // Write-only binding to the contract
	TriggerErrorsFilterer   // Log filterer for contract events
}

// TriggerErrorsCaller is an auto generated read-only Go binding around an Ethereum contract.
type TriggerErrorsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TriggerErrorsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type TriggerErrorsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TriggerErrorsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type TriggerErrorsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// TriggerErrorsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type TriggerErrorsSession struct {
	Contract     *TriggerErrors    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// TriggerErrorsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type TriggerErrorsCallerSession struct {
	Contract *TriggerErrorsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// TriggerErrorsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type TriggerErrorsTransactorSession struct {
	Contract     *TriggerErrorsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// TriggerErrorsRaw is an auto generated low-level Go binding around an Ethereum contract.
type TriggerErrorsRaw struct {
	Contract *TriggerErrors // Generic contract binding to access the raw methods on
}

// TriggerErrorsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type TriggerErrorsCallerRaw struct {
	Contract *TriggerErrorsCaller // Generic read-only contract binding to access the raw methods on
}

// TriggerErrorsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type TriggerErrorsTransactorRaw struct {
	Contract *TriggerErrorsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewTriggerErrors creates a new instance of TriggerErrors, bound to a specific deployed contract.
func NewTriggerErrors(address common.Address, backend bind.ContractBackend) (*TriggerErrors, error) {
	contract, err := bindTriggerErrors(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &TriggerErrors{TriggerErrorsCaller: TriggerErrorsCaller{contract: contract}, TriggerErrorsTransactor: TriggerErrorsTransactor{contract: contract}, TriggerErrorsFilterer: TriggerErrorsFilterer{contract: contract}}, nil
}

// NewTriggerErrorsCaller creates a new read-only instance of TriggerErrors, bound to a specific deployed contract.
func NewTriggerErrorsCaller(address common.Address, caller bind.ContractCaller) (*TriggerErrorsCaller, error) {
	contract, err := bindTriggerErrors(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &TriggerErrorsCaller{contract: contract}, nil
}

// NewTriggerErrorsTransactor creates a new write-only instance of TriggerErrors, bound to a specific deployed contract.
func NewTriggerErrorsTransactor(address common.Address, transactor bind.ContractTransactor) (*TriggerErrorsTransactor, error) {
	contract, err := bindTriggerErrors(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &TriggerErrorsTransactor{contract: contract}, nil
}

// NewTriggerErrorsFilterer creates a new log filterer instance of TriggerErrors, bound to a specific deployed contract.
func NewTriggerErrorsFilterer(address common.Address, filterer bind.ContractFilterer) (*TriggerErrorsFilterer, error) {
	contract, err := bindTriggerErrors(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &TriggerErrorsFilterer{contract: contract}, nil
}

// bindTriggerErrors binds a generic wrapper to an already deployed contract.
func bindTriggerErrors(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := TriggerErrorsMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TriggerErrors *TriggerErrorsRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TriggerErrors.Contract.TriggerErrorsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TriggerErrors *TriggerErrorsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TriggerErrors.Contract.TriggerErrorsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TriggerErrors *TriggerErrorsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TriggerErrors.Contract.TriggerErrorsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_TriggerErrors *TriggerErrorsCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _TriggerErrors.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_TriggerErrors *TriggerErrorsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TriggerErrors.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_TriggerErrors *TriggerErrorsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _TriggerErrors.Contract.contract.Transact(opts, method, params...)
}

// Count is a free data retrieval call binding the contract method 0x06661abd.
//
// Solidity: function count() view returns(uint256)
func (_TriggerErrors *TriggerErrorsCaller) Count(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _TriggerErrors.contract.Call(opts, &out, "count")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Count is a free data retrieval call binding the contract method 0x06661abd.
//
// Solidity: function count() view returns(uint256)
func (_TriggerErrors *TriggerErrorsSession) Count() (*big.Int, error) {
	return _TriggerErrors.Contract.Count(&_TriggerErrors.CallOpts)
}

// Count is a free data retrieval call binding the contract method 0x06661abd.
//
// Solidity: function count() view returns(uint256)
func (_TriggerErrors *TriggerErrorsCallerSession) Count() (*big.Int, error) {
	return _TriggerErrors.Contract.Count(&_TriggerErrors.CallOpts)
}

// OutOfCountersKeccaks is a free data retrieval call binding the contract method 0x2621002a.
//
// Solidity: function outOfCountersKeccaks() pure returns(bytes32 test)
func (_TriggerErrors *TriggerErrorsCaller) OutOfCountersKeccaks(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _TriggerErrors.contract.Call(opts, &out, "outOfCountersKeccaks")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// OutOfCountersKeccaks is a free data retrieval call binding the contract method 0x2621002a.
//
// Solidity: function outOfCountersKeccaks() pure returns(bytes32 test)
func (_TriggerErrors *TriggerErrorsSession) OutOfCountersKeccaks() ([32]byte, error) {
	return _TriggerErrors.Contract.OutOfCountersKeccaks(&_TriggerErrors.CallOpts)
}

// OutOfCountersKeccaks is a free data retrieval call binding the contract method 0x2621002a.
//
// Solidity: function outOfCountersKeccaks() pure returns(bytes32 test)
func (_TriggerErrors *TriggerErrorsCallerSession) OutOfCountersKeccaks() ([32]byte, error) {
	return _TriggerErrors.Contract.OutOfCountersKeccaks(&_TriggerErrors.CallOpts)
}

// OutOfCountersSteps is a free data retrieval call binding the contract method 0x8bd7b538.
//
// Solidity: function outOfCountersSteps() pure returns()
func (_TriggerErrors *TriggerErrorsCaller) OutOfCountersSteps(opts *bind.CallOpts) error {
	var out []interface{}
	err := _TriggerErrors.contract.Call(opts, &out, "outOfCountersSteps")

	if err != nil {
		return err
	}

	return err

}

// OutOfCountersSteps is a free data retrieval call binding the contract method 0x8bd7b538.
//
// Solidity: function outOfCountersSteps() pure returns()
func (_TriggerErrors *TriggerErrorsSession) OutOfCountersSteps() error {
	return _TriggerErrors.Contract.OutOfCountersSteps(&_TriggerErrors.CallOpts)
}

// OutOfCountersSteps is a free data retrieval call binding the contract method 0x8bd7b538.
//
// Solidity: function outOfCountersSteps() pure returns()
func (_TriggerErrors *TriggerErrorsCallerSession) OutOfCountersSteps() error {
	return _TriggerErrors.Contract.OutOfCountersSteps(&_TriggerErrors.CallOpts)
}

// OutOfCountersPoseidon is a paid mutator transaction binding the contract method 0xcb4e8cd1.
//
// Solidity: function outOfCountersPoseidon() returns()
func (_TriggerErrors *TriggerErrorsTransactor) OutOfCountersPoseidon(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TriggerErrors.contract.Transact(opts, "outOfCountersPoseidon")
}

// OutOfCountersPoseidon is a paid mutator transaction binding the contract method 0xcb4e8cd1.
//
// Solidity: function outOfCountersPoseidon() returns()
func (_TriggerErrors *TriggerErrorsSession) OutOfCountersPoseidon() (*types.Transaction, error) {
	return _TriggerErrors.Contract.OutOfCountersPoseidon(&_TriggerErrors.TransactOpts)
}

// OutOfCountersPoseidon is a paid mutator transaction binding the contract method 0xcb4e8cd1.
//
// Solidity: function outOfCountersPoseidon() returns()
func (_TriggerErrors *TriggerErrorsTransactorSession) OutOfCountersPoseidon() (*types.Transaction, error) {
	return _TriggerErrors.Contract.OutOfCountersPoseidon(&_TriggerErrors.TransactOpts)
}

// OutOfGas is a paid mutator transaction binding the contract method 0x31fe52e8.
//
// Solidity: function outOfGas() returns()
func (_TriggerErrors *TriggerErrorsTransactor) OutOfGas(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _TriggerErrors.contract.Transact(opts, "outOfGas")
}

// OutOfGas is a paid mutator transaction binding the contract method 0x31fe52e8.
//
// Solidity: function outOfGas() returns()
func (_TriggerErrors *TriggerErrorsSession) OutOfGas() (*types.Transaction, error) {
	return _TriggerErrors.Contract.OutOfGas(&_TriggerErrors.TransactOpts)
}

// OutOfGas is a paid mutator transaction binding the contract method 0x31fe52e8.
//
// Solidity: function outOfGas() returns()
func (_TriggerErrors *TriggerErrorsTransactorSession) OutOfGas() (*types.Transaction, error) {
	return _TriggerErrors.Contract.OutOfGas(&_TriggerErrors.TransactOpts)
}
