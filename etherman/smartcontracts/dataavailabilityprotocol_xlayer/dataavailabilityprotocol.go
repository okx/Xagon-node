// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package dataavailabilityprotocol_xlayer

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

// DataavailabilityprotocolMetaData contains all meta data concerning the Dataavailabilityprotocol contract.
var DataavailabilityprotocolMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"getProcotolName\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"hash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"dataAvailabilityMessage\",\"type\":\"bytes\"}],\"name\":\"verifyMessage\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// DataavailabilityprotocolABI is the input ABI used to generate the binding from.
// Deprecated: Use DataavailabilityprotocolMetaData.ABI instead.
var DataavailabilityprotocolABI = DataavailabilityprotocolMetaData.ABI

// Dataavailabilityprotocol is an auto generated Go binding around an Ethereum contract.
type Dataavailabilityprotocol struct {
	DataavailabilityprotocolCaller     // Read-only binding to the contract
	DataavailabilityprotocolTransactor // Write-only binding to the contract
	DataavailabilityprotocolFilterer   // Log filterer for contract events
}

// DataavailabilityprotocolCaller is an auto generated read-only Go binding around an Ethereum contract.
type DataavailabilityprotocolCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DataavailabilityprotocolTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DataavailabilityprotocolTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DataavailabilityprotocolFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DataavailabilityprotocolFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DataavailabilityprotocolSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DataavailabilityprotocolSession struct {
	Contract     *Dataavailabilityprotocol // Generic contract binding to set the session for
	CallOpts     bind.CallOpts             // Call options to use throughout this session
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// DataavailabilityprotocolCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DataavailabilityprotocolCallerSession struct {
	Contract *DataavailabilityprotocolCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                   // Call options to use throughout this session
}

// DataavailabilityprotocolTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DataavailabilityprotocolTransactorSession struct {
	Contract     *DataavailabilityprotocolTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                   // Transaction auth options to use throughout this session
}

// DataavailabilityprotocolRaw is an auto generated low-level Go binding around an Ethereum contract.
type DataavailabilityprotocolRaw struct {
	Contract *Dataavailabilityprotocol // Generic contract binding to access the raw methods on
}

// DataavailabilityprotocolCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DataavailabilityprotocolCallerRaw struct {
	Contract *DataavailabilityprotocolCaller // Generic read-only contract binding to access the raw methods on
}

// DataavailabilityprotocolTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DataavailabilityprotocolTransactorRaw struct {
	Contract *DataavailabilityprotocolTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDataavailabilityprotocol creates a new instance of Dataavailabilityprotocol, bound to a specific deployed contract.
func NewDataavailabilityprotocol(address common.Address, backend bind.ContractBackend) (*Dataavailabilityprotocol, error) {
	contract, err := bindDataavailabilityprotocol(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Dataavailabilityprotocol{DataavailabilityprotocolCaller: DataavailabilityprotocolCaller{contract: contract}, DataavailabilityprotocolTransactor: DataavailabilityprotocolTransactor{contract: contract}, DataavailabilityprotocolFilterer: DataavailabilityprotocolFilterer{contract: contract}}, nil
}

// NewDataavailabilityprotocolCaller creates a new read-only instance of Dataavailabilityprotocol, bound to a specific deployed contract.
func NewDataavailabilityprotocolCaller(address common.Address, caller bind.ContractCaller) (*DataavailabilityprotocolCaller, error) {
	contract, err := bindDataavailabilityprotocol(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DataavailabilityprotocolCaller{contract: contract}, nil
}

// NewDataavailabilityprotocolTransactor creates a new write-only instance of Dataavailabilityprotocol, bound to a specific deployed contract.
func NewDataavailabilityprotocolTransactor(address common.Address, transactor bind.ContractTransactor) (*DataavailabilityprotocolTransactor, error) {
	contract, err := bindDataavailabilityprotocol(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DataavailabilityprotocolTransactor{contract: contract}, nil
}

// NewDataavailabilityprotocolFilterer creates a new log filterer instance of Dataavailabilityprotocol, bound to a specific deployed contract.
func NewDataavailabilityprotocolFilterer(address common.Address, filterer bind.ContractFilterer) (*DataavailabilityprotocolFilterer, error) {
	contract, err := bindDataavailabilityprotocol(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DataavailabilityprotocolFilterer{contract: contract}, nil
}

// bindDataavailabilityprotocol binds a generic wrapper to an already deployed contract.
func bindDataavailabilityprotocol(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DataavailabilityprotocolMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Dataavailabilityprotocol *DataavailabilityprotocolRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Dataavailabilityprotocol.Contract.DataavailabilityprotocolCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Dataavailabilityprotocol *DataavailabilityprotocolRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Dataavailabilityprotocol.Contract.DataavailabilityprotocolTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Dataavailabilityprotocol *DataavailabilityprotocolRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Dataavailabilityprotocol.Contract.DataavailabilityprotocolTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Dataavailabilityprotocol *DataavailabilityprotocolCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Dataavailabilityprotocol.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Dataavailabilityprotocol *DataavailabilityprotocolTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Dataavailabilityprotocol.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Dataavailabilityprotocol *DataavailabilityprotocolTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Dataavailabilityprotocol.Contract.contract.Transact(opts, method, params...)
}

// GetProcotolName is a free data retrieval call binding the contract method 0xe4f17120.
//
// Solidity: function getProcotolName() pure returns(string)
func (_Dataavailabilityprotocol *DataavailabilityprotocolCaller) GetProcotolName(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Dataavailabilityprotocol.contract.Call(opts, &out, "getProcotolName")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetProcotolName is a free data retrieval call binding the contract method 0xe4f17120.
//
// Solidity: function getProcotolName() pure returns(string)
func (_Dataavailabilityprotocol *DataavailabilityprotocolSession) GetProcotolName() (string, error) {
	return _Dataavailabilityprotocol.Contract.GetProcotolName(&_Dataavailabilityprotocol.CallOpts)
}

// GetProcotolName is a free data retrieval call binding the contract method 0xe4f17120.
//
// Solidity: function getProcotolName() pure returns(string)
func (_Dataavailabilityprotocol *DataavailabilityprotocolCallerSession) GetProcotolName() (string, error) {
	return _Dataavailabilityprotocol.Contract.GetProcotolName(&_Dataavailabilityprotocol.CallOpts)
}

// VerifyMessage is a free data retrieval call binding the contract method 0x3b51be4b.
//
// Solidity: function verifyMessage(bytes32 hash, bytes dataAvailabilityMessage) view returns()
func (_Dataavailabilityprotocol *DataavailabilityprotocolCaller) VerifyMessage(opts *bind.CallOpts, hash [32]byte, dataAvailabilityMessage []byte) error {
	var out []interface{}
	err := _Dataavailabilityprotocol.contract.Call(opts, &out, "verifyMessage", hash, dataAvailabilityMessage)

	if err != nil {
		return err
	}

	return err

}

// VerifyMessage is a free data retrieval call binding the contract method 0x3b51be4b.
//
// Solidity: function verifyMessage(bytes32 hash, bytes dataAvailabilityMessage) view returns()
func (_Dataavailabilityprotocol *DataavailabilityprotocolSession) VerifyMessage(hash [32]byte, dataAvailabilityMessage []byte) error {
	return _Dataavailabilityprotocol.Contract.VerifyMessage(&_Dataavailabilityprotocol.CallOpts, hash, dataAvailabilityMessage)
}

// VerifyMessage is a free data retrieval call binding the contract method 0x3b51be4b.
//
// Solidity: function verifyMessage(bytes32 hash, bytes dataAvailabilityMessage) view returns()
func (_Dataavailabilityprotocol *DataavailabilityprotocolCallerSession) VerifyMessage(hash [32]byte, dataAvailabilityMessage []byte) error {
	return _Dataavailabilityprotocol.Contract.VerifyMessage(&_Dataavailabilityprotocol.CallOpts, hash, dataAvailabilityMessage)
}
