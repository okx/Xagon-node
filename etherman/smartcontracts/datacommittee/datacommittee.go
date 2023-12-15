// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package datacommittee

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

// DatacommitteeMetaData contains all meta data concerning the Datacommittee contract.
var DatacommitteeMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"CommitteeAddressDoesntExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptyURLNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyRequiredSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnexpectedAddrsAndSignaturesSize\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnexpectedAddrsBytesLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnexpectedCommitteeHash\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WrongAddrOrder\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"committeeHash\",\"type\":\"bytes32\"}],\"name\":\"CommitteeUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"committeeHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAmountOfMembers\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"members\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"url\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"requiredAmountOfSignatures\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_requiredAmountOfSignatures\",\"type\":\"uint256\"},{\"internalType\":\"string[]\",\"name\":\"urls\",\"type\":\"string[]\"},{\"internalType\":\"bytes\",\"name\":\"addrsBytes\",\"type\":\"bytes\"}],\"name\":\"setupCommittee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"signedHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"signaturesAndAddrs\",\"type\":\"bytes\"}],\"name\":\"verifySignatures\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506115e0806100206000396000f3fe608060405234801561001057600080fd5b50600436106100be5760003560e01c80638129fc1c11610076578063c7a823e01161005b578063c7a823e01461015a578063dce1e2b61461016d578063f2fde38b1461017557600080fd5b80638129fc1c1461012a5780638da5cb5b1461013257600080fd5b8063609d4544116100a7578063609d4544146101025780636beedd3914610119578063715018a61461012257600080fd5b8063078fba2a146100c35780635daf08ca146100d8575b600080fd5b6100d66100d1366004610fa9565b610188565b005b6100eb6100e6366004611054565b61048c565b6040516100f992919061106d565b60405180910390f35b61010b60665481565b6040519081526020016100f9565b61010b60655481565b6100d661055e565b6100d6610572565b60335460405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100f9565b6100d66101683660046110f6565b610709565b60675461010b565b6100d6610183366004611142565b61095c565b610190610a10565b82858110156101cb576040517f2e7dcd6e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6101d66014826111ae565b821461020e576040517f2ab6a12900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61021a60676000610eb5565b6000805b828110156104305760006102336014836111ae565b905060008682876102456014836111c5565b92610252939291906111d8565b61025b91611202565b60601c90508888848181106102725761027261124a565b90506020028101906102849190611279565b90506000036102bf576040517fb54b70e400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8073ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff1610610324576040517fd53cfbe000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b809350606760405180604001604052808b8b878181106103465761034661124a565b90506020028101906103589190611279565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092018290525093855250505073ffffffffffffffffffffffffffffffffffffffff851660209283015283546001810185559381522081519192600202019081906103cc90826113af565b5060209190910151600190910180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff90921691909117905550819050610428816114c9565b91505061021e565b508383604051610441929190611501565b6040519081900381206066819055606589905581527f831403fd381b3e6ac875d912ec2eee0e0203d0d29f7b3e0c96fc8f582d6db6579060200160405180910390a150505050505050565b6067818154811061049c57600080fd5b90600052602060002090600202016000915090508060000180546104bf9061130d565b80601f01602080910402602001604051908101604052809291908181526020018280546104eb9061130d565b80156105385780601f1061050d57610100808354040283529160200191610538565b820191906000526020600020905b81548152906001019060200180831161051b57829003601f168201915b5050506001909301549192505073ffffffffffffffffffffffffffffffffffffffff1682565b610566610a10565b6105706000610a91565b565b600054610100900460ff16158080156105925750600054600160ff909116105b806105ac5750303b1580156105ac575060005460ff166001145b61063d576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201527f647920696e697469616c697a656400000000000000000000000000000000000060648201526084015b60405180910390fd5b600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055801561069b57600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff166101001790555b6106a3610b08565b801561070657600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b50565b6000606554604161071a91906111ae565b90508082108061073e575060146107318284611511565b61073b9190611553565b15155b15610775576040517f6b8eec4600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b606654610784838381876111d8565b604051610792929190611501565b6040518091039020146107d1576040517f6b156b2800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008060146107e08486611511565b6107ea9190611567565b905060005b60655481101561095357600061086a88888861080c6041876111ae565b90604161081981896111ae565b61082391906111c5565b92610830939291906111d8565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250610ba892505050565b90506000845b848110156109065760006108856014836111ae565b61088f90896111c5565b905060008a828b6108a16014836111c5565b926108ae939291906111d8565b6108b791611202565b60601c905073ffffffffffffffffffffffffffffffffffffffff851681036108f1576108e48360016111c5565b9750600193505050610906565b505080806108fe906114c9565b915050610870565b508061093e576040517f8431721300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5050808061094b906114c9565b9150506107ef565b50505050505050565b610964610a10565b73ffffffffffffffffffffffffffffffffffffffff8116610a07576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201527f64647265737300000000000000000000000000000000000000000000000000006064820152608401610634565b61070681610a91565b60335473ffffffffffffffffffffffffffffffffffffffff163314610570576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152606401610634565b6033805473ffffffffffffffffffffffffffffffffffffffff8381167fffffffffffffffffffffffff0000000000000000000000000000000000000000831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b600054610100900460ff16610b9f576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602b60248201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960448201527f6e697469616c697a696e670000000000000000000000000000000000000000006064820152608401610634565b61057033610a91565b6000806000610bb78585610bce565b91509150610bc481610c13565b5090505b92915050565b6000808251604103610c045760208301516040840151606085015160001a610bf887828585610dc6565b94509450505050610c0c565b506000905060025b9250929050565b6000816004811115610c2757610c2761157b565b03610c2f5750565b6001816004811115610c4357610c4361157b565b03610caa576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f45434453413a20696e76616c6964207369676e617475726500000000000000006044820152606401610634565b6002816004811115610cbe57610cbe61157b565b03610d25576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601f60248201527f45434453413a20696e76616c6964207369676e6174757265206c656e677468006044820152606401610634565b6003816004811115610d3957610d3961157b565b03610706576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602260248201527f45434453413a20696e76616c6964207369676e6174757265202773272076616c60448201527f75650000000000000000000000000000000000000000000000000000000000006064820152608401610634565b6000807f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a0831115610dfd5750600090506003610eac565b6040805160008082526020820180845289905260ff881692820192909252606081018690526080810185905260019060a0016020604051602081039080840390855afa158015610e51573d6000803e3d6000fd5b50506040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0015191505073ffffffffffffffffffffffffffffffffffffffff8116610ea557600060019250925050610eac565b9150600090505b94509492505050565b508054600082556002029060005260206000209081019061070691905b80821115610f19576000610ee68282610f1d565b506001810180547fffffffffffffffffffffffff0000000000000000000000000000000000000000169055600201610ed2565b5090565b508054610f299061130d565b6000825580601f10610f39575050565b601f01602090049060005260206000209081019061070691905b80821115610f195760008155600101610f53565b60008083601f840112610f7957600080fd5b50813567ffffffffffffffff811115610f9157600080fd5b602083019150836020828501011115610c0c57600080fd5b600080600080600060608688031215610fc157600080fd5b85359450602086013567ffffffffffffffff80821115610fe057600080fd5b818801915088601f830112610ff457600080fd5b81358181111561100357600080fd5b8960208260051b850101111561101857600080fd5b60208301965080955050604088013591508082111561103657600080fd5b5061104388828901610f67565b969995985093965092949392505050565b60006020828403121561106657600080fd5b5035919050565b604081526000835180604084015260005b8181101561109b576020818701810151606086840101520161107e565b5060006060828501015260607fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011684010191505073ffffffffffffffffffffffffffffffffffffffff831660208301529392505050565b60008060006040848603121561110b57600080fd5b83359250602084013567ffffffffffffffff81111561112957600080fd5b61113586828701610f67565b9497909650939450505050565b60006020828403121561115457600080fd5b813573ffffffffffffffffffffffffffffffffffffffff8116811461117857600080fd5b9392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b8082028115828204841417610bc857610bc861117f565b80820180821115610bc857610bc861117f565b600080858511156111e857600080fd5b838611156111f557600080fd5b5050820193919092039150565b7fffffffffffffffffffffffffffffffffffffffff00000000000000000000000081358181169160148510156112425780818660140360031b1b83161692505b505092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe18436030181126112ae57600080fd5b83018035915067ffffffffffffffff8211156112c957600080fd5b602001915036819003821315610c0c57600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600181811c9082168061132157607f821691505b60208210810361135a577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b601f8211156113aa57600081815260208120601f850160051c810160208610156113875750805b601f850160051c820191505b818110156113a657828155600101611393565b5050505b505050565b815167ffffffffffffffff8111156113c9576113c96112de565b6113dd816113d7845461130d565b84611360565b602080601f83116001811461143057600084156113fa5750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b1785556113a6565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b8281101561147d5788860151825594840194600190910190840161145e565b50858210156114b957878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036114fa576114fa61117f565b5060010190565b8183823760009101908152919050565b81810381811115610bc857610bc861117f565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b60008261156257611562611524565b500690565b60008261157657611576611524565b500490565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fdfea264697066735822122062f6d4f079b53519f8ddf0f0e97eef9555bd97448836b76cbbb1071cb28b4c5264736f6c63430008140033",
}

// DatacommitteeABI is the input ABI used to generate the binding from.
// Deprecated: Use DatacommitteeMetaData.ABI instead.
var DatacommitteeABI = DatacommitteeMetaData.ABI

// DatacommitteeBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use DatacommitteeMetaData.Bin instead.
var DatacommitteeBin = DatacommitteeMetaData.Bin

// DeployDatacommittee deploys a new Ethereum contract, binding an instance of Datacommittee to it.
func DeployDatacommittee(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Datacommittee, error) {
	parsed, err := DatacommitteeMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DatacommitteeBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Datacommittee{DatacommitteeCaller: DatacommitteeCaller{contract: contract}, DatacommitteeTransactor: DatacommitteeTransactor{contract: contract}, DatacommitteeFilterer: DatacommitteeFilterer{contract: contract}}, nil
}

// Datacommittee is an auto generated Go binding around an Ethereum contract.
type Datacommittee struct {
	DatacommitteeCaller     // Read-only binding to the contract
	DatacommitteeTransactor // Write-only binding to the contract
	DatacommitteeFilterer   // Log filterer for contract events
}

// DatacommitteeCaller is an auto generated read-only Go binding around an Ethereum contract.
type DatacommitteeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DatacommitteeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DatacommitteeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DatacommitteeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DatacommitteeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DatacommitteeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DatacommitteeSession struct {
	Contract     *Datacommittee    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DatacommitteeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DatacommitteeCallerSession struct {
	Contract *DatacommitteeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// DatacommitteeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DatacommitteeTransactorSession struct {
	Contract     *DatacommitteeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// DatacommitteeRaw is an auto generated low-level Go binding around an Ethereum contract.
type DatacommitteeRaw struct {
	Contract *Datacommittee // Generic contract binding to access the raw methods on
}

// DatacommitteeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DatacommitteeCallerRaw struct {
	Contract *DatacommitteeCaller // Generic read-only contract binding to access the raw methods on
}

// DatacommitteeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DatacommitteeTransactorRaw struct {
	Contract *DatacommitteeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDatacommittee creates a new instance of Datacommittee, bound to a specific deployed contract.
func NewDatacommittee(address common.Address, backend bind.ContractBackend) (*Datacommittee, error) {
	contract, err := bindDatacommittee(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Datacommittee{DatacommitteeCaller: DatacommitteeCaller{contract: contract}, DatacommitteeTransactor: DatacommitteeTransactor{contract: contract}, DatacommitteeFilterer: DatacommitteeFilterer{contract: contract}}, nil
}

// NewDatacommitteeCaller creates a new read-only instance of Datacommittee, bound to a specific deployed contract.
func NewDatacommitteeCaller(address common.Address, caller bind.ContractCaller) (*DatacommitteeCaller, error) {
	contract, err := bindDatacommittee(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DatacommitteeCaller{contract: contract}, nil
}

// NewDatacommitteeTransactor creates a new write-only instance of Datacommittee, bound to a specific deployed contract.
func NewDatacommitteeTransactor(address common.Address, transactor bind.ContractTransactor) (*DatacommitteeTransactor, error) {
	contract, err := bindDatacommittee(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DatacommitteeTransactor{contract: contract}, nil
}

// NewDatacommitteeFilterer creates a new log filterer instance of Datacommittee, bound to a specific deployed contract.
func NewDatacommitteeFilterer(address common.Address, filterer bind.ContractFilterer) (*DatacommitteeFilterer, error) {
	contract, err := bindDatacommittee(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DatacommitteeFilterer{contract: contract}, nil
}

// bindDatacommittee binds a generic wrapper to an already deployed contract.
func bindDatacommittee(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DatacommitteeMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Datacommittee *DatacommitteeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Datacommittee.Contract.DatacommitteeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Datacommittee *DatacommitteeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Datacommittee.Contract.DatacommitteeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Datacommittee *DatacommitteeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Datacommittee.Contract.DatacommitteeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Datacommittee *DatacommitteeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Datacommittee.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Datacommittee *DatacommitteeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Datacommittee.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Datacommittee *DatacommitteeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Datacommittee.Contract.contract.Transact(opts, method, params...)
}

// CommitteeHash is a free data retrieval call binding the contract method 0x609d4544.
//
// Solidity: function committeeHash() view returns(bytes32)
func (_Datacommittee *DatacommitteeCaller) CommitteeHash(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Datacommittee.contract.Call(opts, &out, "committeeHash")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// CommitteeHash is a free data retrieval call binding the contract method 0x609d4544.
//
// Solidity: function committeeHash() view returns(bytes32)
func (_Datacommittee *DatacommitteeSession) CommitteeHash() ([32]byte, error) {
	return _Datacommittee.Contract.CommitteeHash(&_Datacommittee.CallOpts)
}

// CommitteeHash is a free data retrieval call binding the contract method 0x609d4544.
//
// Solidity: function committeeHash() view returns(bytes32)
func (_Datacommittee *DatacommitteeCallerSession) CommitteeHash() ([32]byte, error) {
	return _Datacommittee.Contract.CommitteeHash(&_Datacommittee.CallOpts)
}

// GetAmountOfMembers is a free data retrieval call binding the contract method 0xdce1e2b6.
//
// Solidity: function getAmountOfMembers() view returns(uint256)
func (_Datacommittee *DatacommitteeCaller) GetAmountOfMembers(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Datacommittee.contract.Call(opts, &out, "getAmountOfMembers")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAmountOfMembers is a free data retrieval call binding the contract method 0xdce1e2b6.
//
// Solidity: function getAmountOfMembers() view returns(uint256)
func (_Datacommittee *DatacommitteeSession) GetAmountOfMembers() (*big.Int, error) {
	return _Datacommittee.Contract.GetAmountOfMembers(&_Datacommittee.CallOpts)
}

// GetAmountOfMembers is a free data retrieval call binding the contract method 0xdce1e2b6.
//
// Solidity: function getAmountOfMembers() view returns(uint256)
func (_Datacommittee *DatacommitteeCallerSession) GetAmountOfMembers() (*big.Int, error) {
	return _Datacommittee.Contract.GetAmountOfMembers(&_Datacommittee.CallOpts)
}

// Members is a free data retrieval call binding the contract method 0x5daf08ca.
//
// Solidity: function members(uint256 ) view returns(string url, address addr)
func (_Datacommittee *DatacommitteeCaller) Members(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Url  string
	Addr common.Address
}, error) {
	var out []interface{}
	err := _Datacommittee.contract.Call(opts, &out, "members", arg0)

	outstruct := new(struct {
		Url  string
		Addr common.Address
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Url = *abi.ConvertType(out[0], new(string)).(*string)
	outstruct.Addr = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)

	return *outstruct, err

}

// Members is a free data retrieval call binding the contract method 0x5daf08ca.
//
// Solidity: function members(uint256 ) view returns(string url, address addr)
func (_Datacommittee *DatacommitteeSession) Members(arg0 *big.Int) (struct {
	Url  string
	Addr common.Address
}, error) {
	return _Datacommittee.Contract.Members(&_Datacommittee.CallOpts, arg0)
}

// Members is a free data retrieval call binding the contract method 0x5daf08ca.
//
// Solidity: function members(uint256 ) view returns(string url, address addr)
func (_Datacommittee *DatacommitteeCallerSession) Members(arg0 *big.Int) (struct {
	Url  string
	Addr common.Address
}, error) {
	return _Datacommittee.Contract.Members(&_Datacommittee.CallOpts, arg0)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Datacommittee *DatacommitteeCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Datacommittee.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Datacommittee *DatacommitteeSession) Owner() (common.Address, error) {
	return _Datacommittee.Contract.Owner(&_Datacommittee.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Datacommittee *DatacommitteeCallerSession) Owner() (common.Address, error) {
	return _Datacommittee.Contract.Owner(&_Datacommittee.CallOpts)
}

// RequiredAmountOfSignatures is a free data retrieval call binding the contract method 0x6beedd39.
//
// Solidity: function requiredAmountOfSignatures() view returns(uint256)
func (_Datacommittee *DatacommitteeCaller) RequiredAmountOfSignatures(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Datacommittee.contract.Call(opts, &out, "requiredAmountOfSignatures")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RequiredAmountOfSignatures is a free data retrieval call binding the contract method 0x6beedd39.
//
// Solidity: function requiredAmountOfSignatures() view returns(uint256)
func (_Datacommittee *DatacommitteeSession) RequiredAmountOfSignatures() (*big.Int, error) {
	return _Datacommittee.Contract.RequiredAmountOfSignatures(&_Datacommittee.CallOpts)
}

// RequiredAmountOfSignatures is a free data retrieval call binding the contract method 0x6beedd39.
//
// Solidity: function requiredAmountOfSignatures() view returns(uint256)
func (_Datacommittee *DatacommitteeCallerSession) RequiredAmountOfSignatures() (*big.Int, error) {
	return _Datacommittee.Contract.RequiredAmountOfSignatures(&_Datacommittee.CallOpts)
}

// VerifySignatures is a free data retrieval call binding the contract method 0xc7a823e0.
//
// Solidity: function verifySignatures(bytes32 signedHash, bytes signaturesAndAddrs) view returns()
func (_Datacommittee *DatacommitteeCaller) VerifySignatures(opts *bind.CallOpts, signedHash [32]byte, signaturesAndAddrs []byte) error {
	var out []interface{}
	err := _Datacommittee.contract.Call(opts, &out, "verifySignatures", signedHash, signaturesAndAddrs)

	if err != nil {
		return err
	}

	return err

}

// VerifySignatures is a free data retrieval call binding the contract method 0xc7a823e0.
//
// Solidity: function verifySignatures(bytes32 signedHash, bytes signaturesAndAddrs) view returns()
func (_Datacommittee *DatacommitteeSession) VerifySignatures(signedHash [32]byte, signaturesAndAddrs []byte) error {
	return _Datacommittee.Contract.VerifySignatures(&_Datacommittee.CallOpts, signedHash, signaturesAndAddrs)
}

// VerifySignatures is a free data retrieval call binding the contract method 0xc7a823e0.
//
// Solidity: function verifySignatures(bytes32 signedHash, bytes signaturesAndAddrs) view returns()
func (_Datacommittee *DatacommitteeCallerSession) VerifySignatures(signedHash [32]byte, signaturesAndAddrs []byte) error {
	return _Datacommittee.Contract.VerifySignatures(&_Datacommittee.CallOpts, signedHash, signaturesAndAddrs)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_Datacommittee *DatacommitteeTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Datacommittee.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_Datacommittee *DatacommitteeSession) Initialize() (*types.Transaction, error) {
	return _Datacommittee.Contract.Initialize(&_Datacommittee.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_Datacommittee *DatacommitteeTransactorSession) Initialize() (*types.Transaction, error) {
	return _Datacommittee.Contract.Initialize(&_Datacommittee.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Datacommittee *DatacommitteeTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Datacommittee.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Datacommittee *DatacommitteeSession) RenounceOwnership() (*types.Transaction, error) {
	return _Datacommittee.Contract.RenounceOwnership(&_Datacommittee.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Datacommittee *DatacommitteeTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Datacommittee.Contract.RenounceOwnership(&_Datacommittee.TransactOpts)
}

// SetupCommittee is a paid mutator transaction binding the contract method 0x078fba2a.
//
// Solidity: function setupCommittee(uint256 _requiredAmountOfSignatures, string[] urls, bytes addrsBytes) returns()
func (_Datacommittee *DatacommitteeTransactor) SetupCommittee(opts *bind.TransactOpts, _requiredAmountOfSignatures *big.Int, urls []string, addrsBytes []byte) (*types.Transaction, error) {
	return _Datacommittee.contract.Transact(opts, "setupCommittee", _requiredAmountOfSignatures, urls, addrsBytes)
}

// SetupCommittee is a paid mutator transaction binding the contract method 0x078fba2a.
//
// Solidity: function setupCommittee(uint256 _requiredAmountOfSignatures, string[] urls, bytes addrsBytes) returns()
func (_Datacommittee *DatacommitteeSession) SetupCommittee(_requiredAmountOfSignatures *big.Int, urls []string, addrsBytes []byte) (*types.Transaction, error) {
	return _Datacommittee.Contract.SetupCommittee(&_Datacommittee.TransactOpts, _requiredAmountOfSignatures, urls, addrsBytes)
}

// SetupCommittee is a paid mutator transaction binding the contract method 0x078fba2a.
//
// Solidity: function setupCommittee(uint256 _requiredAmountOfSignatures, string[] urls, bytes addrsBytes) returns()
func (_Datacommittee *DatacommitteeTransactorSession) SetupCommittee(_requiredAmountOfSignatures *big.Int, urls []string, addrsBytes []byte) (*types.Transaction, error) {
	return _Datacommittee.Contract.SetupCommittee(&_Datacommittee.TransactOpts, _requiredAmountOfSignatures, urls, addrsBytes)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Datacommittee *DatacommitteeTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Datacommittee.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Datacommittee *DatacommitteeSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Datacommittee.Contract.TransferOwnership(&_Datacommittee.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Datacommittee *DatacommitteeTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Datacommittee.Contract.TransferOwnership(&_Datacommittee.TransactOpts, newOwner)
}

// DatacommitteeCommitteeUpdatedIterator is returned from FilterCommitteeUpdated and is used to iterate over the raw logs and unpacked data for CommitteeUpdated events raised by the Datacommittee contract.
type DatacommitteeCommitteeUpdatedIterator struct {
	Event *DatacommitteeCommitteeUpdated // Event containing the contract specifics and raw log

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
func (it *DatacommitteeCommitteeUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DatacommitteeCommitteeUpdated)
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
		it.Event = new(DatacommitteeCommitteeUpdated)
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
func (it *DatacommitteeCommitteeUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DatacommitteeCommitteeUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DatacommitteeCommitteeUpdated represents a CommitteeUpdated event raised by the Datacommittee contract.
type DatacommitteeCommitteeUpdated struct {
	CommitteeHash [32]byte
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterCommitteeUpdated is a free log retrieval operation binding the contract event 0x831403fd381b3e6ac875d912ec2eee0e0203d0d29f7b3e0c96fc8f582d6db657.
//
// Solidity: event CommitteeUpdated(bytes32 committeeHash)
func (_Datacommittee *DatacommitteeFilterer) FilterCommitteeUpdated(opts *bind.FilterOpts) (*DatacommitteeCommitteeUpdatedIterator, error) {

	logs, sub, err := _Datacommittee.contract.FilterLogs(opts, "CommitteeUpdated")
	if err != nil {
		return nil, err
	}
	return &DatacommitteeCommitteeUpdatedIterator{contract: _Datacommittee.contract, event: "CommitteeUpdated", logs: logs, sub: sub}, nil
}

// WatchCommitteeUpdated is a free log subscription operation binding the contract event 0x831403fd381b3e6ac875d912ec2eee0e0203d0d29f7b3e0c96fc8f582d6db657.
//
// Solidity: event CommitteeUpdated(bytes32 committeeHash)
func (_Datacommittee *DatacommitteeFilterer) WatchCommitteeUpdated(opts *bind.WatchOpts, sink chan<- *DatacommitteeCommitteeUpdated) (event.Subscription, error) {

	logs, sub, err := _Datacommittee.contract.WatchLogs(opts, "CommitteeUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DatacommitteeCommitteeUpdated)
				if err := _Datacommittee.contract.UnpackLog(event, "CommitteeUpdated", log); err != nil {
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

// ParseCommitteeUpdated is a log parse operation binding the contract event 0x831403fd381b3e6ac875d912ec2eee0e0203d0d29f7b3e0c96fc8f582d6db657.
//
// Solidity: event CommitteeUpdated(bytes32 committeeHash)
func (_Datacommittee *DatacommitteeFilterer) ParseCommitteeUpdated(log types.Log) (*DatacommitteeCommitteeUpdated, error) {
	event := new(DatacommitteeCommitteeUpdated)
	if err := _Datacommittee.contract.UnpackLog(event, "CommitteeUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DatacommitteeInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Datacommittee contract.
type DatacommitteeInitializedIterator struct {
	Event *DatacommitteeInitialized // Event containing the contract specifics and raw log

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
func (it *DatacommitteeInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DatacommitteeInitialized)
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
		it.Event = new(DatacommitteeInitialized)
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
func (it *DatacommitteeInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DatacommitteeInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DatacommitteeInitialized represents a Initialized event raised by the Datacommittee contract.
type DatacommitteeInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Datacommittee *DatacommitteeFilterer) FilterInitialized(opts *bind.FilterOpts) (*DatacommitteeInitializedIterator, error) {

	logs, sub, err := _Datacommittee.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &DatacommitteeInitializedIterator{contract: _Datacommittee.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Datacommittee *DatacommitteeFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *DatacommitteeInitialized) (event.Subscription, error) {

	logs, sub, err := _Datacommittee.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DatacommitteeInitialized)
				if err := _Datacommittee.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Datacommittee *DatacommitteeFilterer) ParseInitialized(log types.Log) (*DatacommitteeInitialized, error) {
	event := new(DatacommitteeInitialized)
	if err := _Datacommittee.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DatacommitteeOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Datacommittee contract.
type DatacommitteeOwnershipTransferredIterator struct {
	Event *DatacommitteeOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *DatacommitteeOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DatacommitteeOwnershipTransferred)
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
		it.Event = new(DatacommitteeOwnershipTransferred)
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
func (it *DatacommitteeOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DatacommitteeOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DatacommitteeOwnershipTransferred represents a OwnershipTransferred event raised by the Datacommittee contract.
type DatacommitteeOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Datacommittee *DatacommitteeFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*DatacommitteeOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Datacommittee.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &DatacommitteeOwnershipTransferredIterator{contract: _Datacommittee.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Datacommittee *DatacommitteeFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *DatacommitteeOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Datacommittee.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DatacommitteeOwnershipTransferred)
				if err := _Datacommittee.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Datacommittee *DatacommitteeFilterer) ParseOwnershipTransferred(log types.Log) (*DatacommitteeOwnershipTransferred, error) {
	event := new(DatacommitteeOwnershipTransferred)
	if err := _Datacommittee.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
