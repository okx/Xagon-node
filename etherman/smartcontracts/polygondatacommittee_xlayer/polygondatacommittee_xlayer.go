// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package polygondatacommittee_xlayer

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

// PolygondatacommitteeMetaData contains all meta data concerning the Polygondatacommittee contract.
var PolygondatacommitteeMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"CommitteeAddressDoesntExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"EmptyURLNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TooManyRequiredSignatures\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnexpectedAddrsAndSignaturesSize\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnexpectedAddrsBytesLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnexpectedCommitteeHash\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"WrongAddrOrder\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"committeeHash\",\"type\":\"bytes32\"}],\"name\":\"CommitteeUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"committeeHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAmountOfMembers\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getProcotolName\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"members\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"url\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"requiredAmountOfSignatures\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_requiredAmountOfSignatures\",\"type\":\"uint256\"},{\"internalType\":\"string[]\",\"name\":\"urls\",\"type\":\"string[]\"},{\"internalType\":\"bytes\",\"name\":\"addrsBytes\",\"type\":\"bytes\"}],\"name\":\"setupCommittee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"signedHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"signaturesAndAddrs\",\"type\":\"bytes\"}],\"name\":\"verifyMessage\",\"outputs\":[],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50611646806100206000396000f3fe608060405234801561001057600080fd5b50600436106100c95760003560e01c8063715018a611610081578063dce1e2b61161005b578063dce1e2b614610178578063e4f1712014610180578063f2fde38b146101bf57600080fd5b8063715018a6146101405780638129fc1c146101485780638da5cb5b1461015057600080fd5b80635daf08ca116100b25780635daf08ca146100f6578063609d4544146101205780636beedd391461013757600080fd5b8063078fba2a146100ce5780633b51be4b146100e3575b600080fd5b6100e16100dc366004610fe9565b6101d2565b005b6100e16100f1366004611094565b6104d4565b6101096101043660046110e0565b61071f565b60405161011792919061115d565b60405180910390f35b61012960665481565b604051908152602001610117565b61012960655481565b6100e16107f1565b6100e1610805565b60335460405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610117565b606754610129565b604080518082018252601981527f44617461417661696c6162696c697479436f6d6d697474656500000000000000602082015290516101179190611195565b6100e16101cd3660046111af565b61099c565b6101da610a50565b8285811015610215576040517f2e7dcd6e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610220601482611214565b8214610258576040517f2ab6a12900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61026460676000610ef5565b6000805b8281101561047857600061027d601483611214565b9050600086828761028f60148361122b565b9261029c9392919061123e565b6102a591611268565b60601c90508888848181106102bc576102bc6112b0565b90506020028101906102ce91906112df565b9050600003610309576040517fb54b70e400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8073ffffffffffffffffffffffffffffffffffffffff168473ffffffffffffffffffffffffffffffffffffffff161061036e576040517fd53cfbe000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b606760405180604001604052808b8b8781811061038d5761038d6112b0565b905060200281019061039f91906112df565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092018290525093855250505073ffffffffffffffffffffffffffffffffffffffff851660209283015283546001810185559381522081519192600202019081906104139082611415565b5060209190910151600190910180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff90921691909117905592508190506104708161152f565b915050610268565b508383604051610489929190611567565b6040519081900381206066819055606589905581527f831403fd381b3e6ac875d912ec2eee0e0203d0d29f7b3e0c96fc8f582d6db6579060200160405180910390a150505050505050565b60655460006104e4826041611214565b905080831080610508575060146104fb8285611577565b61050591906115b9565b15155b1561053f576040517f6b8eec4600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60665461054e8483818861123e565b60405161055c929190611567565b60405180910390201461059b576040517f6b156b2800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60008060146105aa8487611577565b6105b491906115cd565b905060005b848110156107155760006105ce604183611214565b9050600061062b8a8a848b6105e460418361122b565b926105f19392919061123e565b8080601f016020809104026020016040519081016040528093929190818152602001838380828437600092019190915250610ad192505050565b90506000855b858110156106c7576000610646601483611214565b610650908a61122b565b905060008c828d61066260148361122b565b9261066f9392919061123e565b61067891611268565b60601c905073ffffffffffffffffffffffffffffffffffffffff851681036106b2576106a583600161122b565b98506001935050506106c7565b505080806106bf9061152f565b915050610631565b50806106ff576040517f8431721300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b505050808061070d9061152f565b9150506105b9565b5050505050505050565b6067818154811061072f57600080fd5b906000526020600020906002020160009150905080600001805461075290611373565b80601f016020809104026020016040519081016040528092919081815260200182805461077e90611373565b80156107cb5780601f106107a0576101008083540402835291602001916107cb565b820191906000526020600020905b8154815290600101906020018083116107ae57829003601f168201915b5050506001909301549192505073ffffffffffffffffffffffffffffffffffffffff1682565b6107f9610a50565b6108036000610af7565b565b600054610100900460ff16158080156108255750600054600160ff909116105b8061083f5750303b15801561083f575060005460ff166001145b6108d0576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201527f647920696e697469616c697a656400000000000000000000000000000000000060648201526084015b60405180910390fd5b600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055801561092e57600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff166101001790555b610936610b6e565b801561099957600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b50565b6109a4610a50565b73ffffffffffffffffffffffffffffffffffffffff8116610a47576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201527f646472657373000000000000000000000000000000000000000000000000000060648201526084016108c7565b61099981610af7565b60335473ffffffffffffffffffffffffffffffffffffffff163314610803576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064016108c7565b6000806000610ae08585610c0e565b91509150610aed81610c53565b5090505b92915050565b6033805473ffffffffffffffffffffffffffffffffffffffff8381167fffffffffffffffffffffffff0000000000000000000000000000000000000000831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b600054610100900460ff16610c05576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602b60248201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960448201527f6e697469616c697a696e6700000000000000000000000000000000000000000060648201526084016108c7565b61080333610af7565b6000808251604103610c445760208301516040840151606085015160001a610c3887828585610e06565b94509450505050610c4c565b506000905060025b9250929050565b6000816004811115610c6757610c676115e1565b03610c6f5750565b6001816004811115610c8357610c836115e1565b03610cea576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601860248201527f45434453413a20696e76616c6964207369676e6174757265000000000000000060448201526064016108c7565b6002816004811115610cfe57610cfe6115e1565b03610d65576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601f60248201527f45434453413a20696e76616c6964207369676e6174757265206c656e6774680060448201526064016108c7565b6003816004811115610d7957610d796115e1565b03610999576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602260248201527f45434453413a20696e76616c6964207369676e6174757265202773272076616c60448201527f756500000000000000000000000000000000000000000000000000000000000060648201526084016108c7565b6000807f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a0831115610e3d5750600090506003610eec565b6040805160008082526020820180845289905260ff881692820192909252606081018690526080810185905260019060a0016020604051602081039080840390855afa158015610e91573d6000803e3d6000fd5b50506040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0015191505073ffffffffffffffffffffffffffffffffffffffff8116610ee557600060019250925050610eec565b9150600090505b94509492505050565b508054600082556002029060005260206000209081019061099991905b80821115610f59576000610f268282610f5d565b506001810180547fffffffffffffffffffffffff0000000000000000000000000000000000000000169055600201610f12565b5090565b508054610f6990611373565b6000825580601f10610f79575050565b601f01602090049060005260206000209081019061099991905b80821115610f595760008155600101610f93565b60008083601f840112610fb957600080fd5b50813567ffffffffffffffff811115610fd157600080fd5b602083019150836020828501011115610c4c57600080fd5b60008060008060006060868803121561100157600080fd5b85359450602086013567ffffffffffffffff8082111561102057600080fd5b818801915088601f83011261103457600080fd5b81358181111561104357600080fd5b8960208260051b850101111561105857600080fd5b60208301965080955050604088013591508082111561107657600080fd5b5061108388828901610fa7565b969995985093965092949392505050565b6000806000604084860312156110a957600080fd5b83359250602084013567ffffffffffffffff8111156110c757600080fd5b6110d386828701610fa7565b9497909650939450505050565b6000602082840312156110f257600080fd5b5035919050565b6000815180845260005b8181101561111f57602081850181015186830182015201611103565b5060006020828601015260207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f83011685010191505092915050565b60408152600061117060408301856110f9565b905073ffffffffffffffffffffffffffffffffffffffff831660208301529392505050565b6020815260006111a860208301846110f9565b9392505050565b6000602082840312156111c157600080fd5b813573ffffffffffffffffffffffffffffffffffffffff811681146111a857600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b8082028115828204841417610af157610af16111e5565b80820180821115610af157610af16111e5565b6000808585111561124e57600080fd5b8386111561125b57600080fd5b5050820193919092039150565b7fffffffffffffffffffffffffffffffffffffffff00000000000000000000000081358181169160148510156112a85780818660140360031b1b83161692505b505092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60008083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe184360301811261131457600080fd5b83018035915067ffffffffffffffff82111561132f57600080fd5b602001915036819003821315610c4c57600080fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b600181811c9082168061138757607f821691505b6020821081036113c0577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b601f82111561141057600081815260208120601f850160051c810160208610156113ed5750805b601f850160051c820191505b8181101561140c578281556001016113f9565b5050505b505050565b815167ffffffffffffffff81111561142f5761142f611344565b6114438161143d8454611373565b846113c6565b602080601f83116001811461149657600084156114605750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b17855561140c565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b828110156114e3578886015182559484019460019091019084016114c4565b508582101561151f57878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203611560576115606111e5565b5060010190565b8183823760009101908152919050565b81810381811115610af157610af16111e5565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b6000826115c8576115c861158a565b500690565b6000826115dc576115dc61158a565b500490565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052602160045260246000fdfea2646970667358221220e34b71e7c7c23d67a42aa345fc1c3d9c57287ac1c2a2024084974dcc23e4088864736f6c63430008140033",
}

// PolygondatacommitteeABI is the input ABI used to generate the binding from.
// Deprecated: Use PolygondatacommitteeMetaData.ABI instead.
var PolygondatacommitteeABI = PolygondatacommitteeMetaData.ABI

// PolygondatacommitteeBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use PolygondatacommitteeMetaData.Bin instead.
var PolygondatacommitteeBin = PolygondatacommitteeMetaData.Bin

// DeployPolygondatacommittee deploys a new Ethereum contract, binding an instance of Polygondatacommittee to it.
func DeployPolygondatacommittee(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Polygondatacommittee, error) {
	parsed, err := PolygondatacommitteeMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(PolygondatacommitteeBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Polygondatacommittee{PolygondatacommitteeCaller: PolygondatacommitteeCaller{contract: contract}, PolygondatacommitteeTransactor: PolygondatacommitteeTransactor{contract: contract}, PolygondatacommitteeFilterer: PolygondatacommitteeFilterer{contract: contract}}, nil
}

// Polygondatacommittee is an auto generated Go binding around an Ethereum contract.
type Polygondatacommittee struct {
	PolygondatacommitteeCaller     // Read-only binding to the contract
	PolygondatacommitteeTransactor // Write-only binding to the contract
	PolygondatacommitteeFilterer   // Log filterer for contract events
}

// PolygondatacommitteeCaller is an auto generated read-only Go binding around an Ethereum contract.
type PolygondatacommitteeCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PolygondatacommitteeTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PolygondatacommitteeTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PolygondatacommitteeFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PolygondatacommitteeFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PolygondatacommitteeSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PolygondatacommitteeSession struct {
	Contract     *Polygondatacommittee // Generic contract binding to set the session for
	CallOpts     bind.CallOpts         // Call options to use throughout this session
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// PolygondatacommitteeCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PolygondatacommitteeCallerSession struct {
	Contract *PolygondatacommitteeCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts               // Call options to use throughout this session
}

// PolygondatacommitteeTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PolygondatacommitteeTransactorSession struct {
	Contract     *PolygondatacommitteeTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts               // Transaction auth options to use throughout this session
}

// PolygondatacommitteeRaw is an auto generated low-level Go binding around an Ethereum contract.
type PolygondatacommitteeRaw struct {
	Contract *Polygondatacommittee // Generic contract binding to access the raw methods on
}

// PolygondatacommitteeCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PolygondatacommitteeCallerRaw struct {
	Contract *PolygondatacommitteeCaller // Generic read-only contract binding to access the raw methods on
}

// PolygondatacommitteeTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PolygondatacommitteeTransactorRaw struct {
	Contract *PolygondatacommitteeTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPolygondatacommittee creates a new instance of Polygondatacommittee, bound to a specific deployed contract.
func NewPolygondatacommittee(address common.Address, backend bind.ContractBackend) (*Polygondatacommittee, error) {
	contract, err := bindPolygondatacommittee(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Polygondatacommittee{PolygondatacommitteeCaller: PolygondatacommitteeCaller{contract: contract}, PolygondatacommitteeTransactor: PolygondatacommitteeTransactor{contract: contract}, PolygondatacommitteeFilterer: PolygondatacommitteeFilterer{contract: contract}}, nil
}

// NewPolygondatacommitteeCaller creates a new read-only instance of Polygondatacommittee, bound to a specific deployed contract.
func NewPolygondatacommitteeCaller(address common.Address, caller bind.ContractCaller) (*PolygondatacommitteeCaller, error) {
	contract, err := bindPolygondatacommittee(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PolygondatacommitteeCaller{contract: contract}, nil
}

// NewPolygondatacommitteeTransactor creates a new write-only instance of Polygondatacommittee, bound to a specific deployed contract.
func NewPolygondatacommitteeTransactor(address common.Address, transactor bind.ContractTransactor) (*PolygondatacommitteeTransactor, error) {
	contract, err := bindPolygondatacommittee(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PolygondatacommitteeTransactor{contract: contract}, nil
}

// NewPolygondatacommitteeFilterer creates a new log filterer instance of Polygondatacommittee, bound to a specific deployed contract.
func NewPolygondatacommitteeFilterer(address common.Address, filterer bind.ContractFilterer) (*PolygondatacommitteeFilterer, error) {
	contract, err := bindPolygondatacommittee(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PolygondatacommitteeFilterer{contract: contract}, nil
}

// bindPolygondatacommittee binds a generic wrapper to an already deployed contract.
func bindPolygondatacommittee(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PolygondatacommitteeMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Polygondatacommittee *PolygondatacommitteeRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Polygondatacommittee.Contract.PolygondatacommitteeCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Polygondatacommittee *PolygondatacommitteeRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Polygondatacommittee.Contract.PolygondatacommitteeTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Polygondatacommittee *PolygondatacommitteeRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Polygondatacommittee.Contract.PolygondatacommitteeTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Polygondatacommittee *PolygondatacommitteeCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Polygondatacommittee.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Polygondatacommittee *PolygondatacommitteeTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Polygondatacommittee.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Polygondatacommittee *PolygondatacommitteeTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Polygondatacommittee.Contract.contract.Transact(opts, method, params...)
}

// CommitteeHash is a free data retrieval call binding the contract method 0x609d4544.
//
// Solidity: function committeeHash() view returns(bytes32)
func (_Polygondatacommittee *PolygondatacommitteeCaller) CommitteeHash(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Polygondatacommittee.contract.Call(opts, &out, "committeeHash")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// CommitteeHash is a free data retrieval call binding the contract method 0x609d4544.
//
// Solidity: function committeeHash() view returns(bytes32)
func (_Polygondatacommittee *PolygondatacommitteeSession) CommitteeHash() ([32]byte, error) {
	return _Polygondatacommittee.Contract.CommitteeHash(&_Polygondatacommittee.CallOpts)
}

// CommitteeHash is a free data retrieval call binding the contract method 0x609d4544.
//
// Solidity: function committeeHash() view returns(bytes32)
func (_Polygondatacommittee *PolygondatacommitteeCallerSession) CommitteeHash() ([32]byte, error) {
	return _Polygondatacommittee.Contract.CommitteeHash(&_Polygondatacommittee.CallOpts)
}

// GetAmountOfMembers is a free data retrieval call binding the contract method 0xdce1e2b6.
//
// Solidity: function getAmountOfMembers() view returns(uint256)
func (_Polygondatacommittee *PolygondatacommitteeCaller) GetAmountOfMembers(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Polygondatacommittee.contract.Call(opts, &out, "getAmountOfMembers")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAmountOfMembers is a free data retrieval call binding the contract method 0xdce1e2b6.
//
// Solidity: function getAmountOfMembers() view returns(uint256)
func (_Polygondatacommittee *PolygondatacommitteeSession) GetAmountOfMembers() (*big.Int, error) {
	return _Polygondatacommittee.Contract.GetAmountOfMembers(&_Polygondatacommittee.CallOpts)
}

// GetAmountOfMembers is a free data retrieval call binding the contract method 0xdce1e2b6.
//
// Solidity: function getAmountOfMembers() view returns(uint256)
func (_Polygondatacommittee *PolygondatacommitteeCallerSession) GetAmountOfMembers() (*big.Int, error) {
	return _Polygondatacommittee.Contract.GetAmountOfMembers(&_Polygondatacommittee.CallOpts)
}

// GetProcotolName is a free data retrieval call binding the contract method 0xe4f17120.
//
// Solidity: function getProcotolName() pure returns(string)
func (_Polygondatacommittee *PolygondatacommitteeCaller) GetProcotolName(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Polygondatacommittee.contract.Call(opts, &out, "getProcotolName")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetProcotolName is a free data retrieval call binding the contract method 0xe4f17120.
//
// Solidity: function getProcotolName() pure returns(string)
func (_Polygondatacommittee *PolygondatacommitteeSession) GetProcotolName() (string, error) {
	return _Polygondatacommittee.Contract.GetProcotolName(&_Polygondatacommittee.CallOpts)
}

// GetProcotolName is a free data retrieval call binding the contract method 0xe4f17120.
//
// Solidity: function getProcotolName() pure returns(string)
func (_Polygondatacommittee *PolygondatacommitteeCallerSession) GetProcotolName() (string, error) {
	return _Polygondatacommittee.Contract.GetProcotolName(&_Polygondatacommittee.CallOpts)
}

// Members is a free data retrieval call binding the contract method 0x5daf08ca.
//
// Solidity: function members(uint256 ) view returns(string url, address addr)
func (_Polygondatacommittee *PolygondatacommitteeCaller) Members(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Url  string
	Addr common.Address
}, error) {
	var out []interface{}
	err := _Polygondatacommittee.contract.Call(opts, &out, "members", arg0)

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
func (_Polygondatacommittee *PolygondatacommitteeSession) Members(arg0 *big.Int) (struct {
	Url  string
	Addr common.Address
}, error) {
	return _Polygondatacommittee.Contract.Members(&_Polygondatacommittee.CallOpts, arg0)
}

// Members is a free data retrieval call binding the contract method 0x5daf08ca.
//
// Solidity: function members(uint256 ) view returns(string url, address addr)
func (_Polygondatacommittee *PolygondatacommitteeCallerSession) Members(arg0 *big.Int) (struct {
	Url  string
	Addr common.Address
}, error) {
	return _Polygondatacommittee.Contract.Members(&_Polygondatacommittee.CallOpts, arg0)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Polygondatacommittee *PolygondatacommitteeCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Polygondatacommittee.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Polygondatacommittee *PolygondatacommitteeSession) Owner() (common.Address, error) {
	return _Polygondatacommittee.Contract.Owner(&_Polygondatacommittee.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Polygondatacommittee *PolygondatacommitteeCallerSession) Owner() (common.Address, error) {
	return _Polygondatacommittee.Contract.Owner(&_Polygondatacommittee.CallOpts)
}

// RequiredAmountOfSignatures is a free data retrieval call binding the contract method 0x6beedd39.
//
// Solidity: function requiredAmountOfSignatures() view returns(uint256)
func (_Polygondatacommittee *PolygondatacommitteeCaller) RequiredAmountOfSignatures(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Polygondatacommittee.contract.Call(opts, &out, "requiredAmountOfSignatures")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RequiredAmountOfSignatures is a free data retrieval call binding the contract method 0x6beedd39.
//
// Solidity: function requiredAmountOfSignatures() view returns(uint256)
func (_Polygondatacommittee *PolygondatacommitteeSession) RequiredAmountOfSignatures() (*big.Int, error) {
	return _Polygondatacommittee.Contract.RequiredAmountOfSignatures(&_Polygondatacommittee.CallOpts)
}

// RequiredAmountOfSignatures is a free data retrieval call binding the contract method 0x6beedd39.
//
// Solidity: function requiredAmountOfSignatures() view returns(uint256)
func (_Polygondatacommittee *PolygondatacommitteeCallerSession) RequiredAmountOfSignatures() (*big.Int, error) {
	return _Polygondatacommittee.Contract.RequiredAmountOfSignatures(&_Polygondatacommittee.CallOpts)
}

// VerifyMessage is a free data retrieval call binding the contract method 0x3b51be4b.
//
// Solidity: function verifyMessage(bytes32 signedHash, bytes signaturesAndAddrs) view returns()
func (_Polygondatacommittee *PolygondatacommitteeCaller) VerifyMessage(opts *bind.CallOpts, signedHash [32]byte, signaturesAndAddrs []byte) error {
	var out []interface{}
	err := _Polygondatacommittee.contract.Call(opts, &out, "verifyMessage", signedHash, signaturesAndAddrs)

	if err != nil {
		return err
	}

	return err

}

// VerifyMessage is a free data retrieval call binding the contract method 0x3b51be4b.
//
// Solidity: function verifyMessage(bytes32 signedHash, bytes signaturesAndAddrs) view returns()
func (_Polygondatacommittee *PolygondatacommitteeSession) VerifyMessage(signedHash [32]byte, signaturesAndAddrs []byte) error {
	return _Polygondatacommittee.Contract.VerifyMessage(&_Polygondatacommittee.CallOpts, signedHash, signaturesAndAddrs)
}

// VerifyMessage is a free data retrieval call binding the contract method 0x3b51be4b.
//
// Solidity: function verifyMessage(bytes32 signedHash, bytes signaturesAndAddrs) view returns()
func (_Polygondatacommittee *PolygondatacommitteeCallerSession) VerifyMessage(signedHash [32]byte, signaturesAndAddrs []byte) error {
	return _Polygondatacommittee.Contract.VerifyMessage(&_Polygondatacommittee.CallOpts, signedHash, signaturesAndAddrs)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_Polygondatacommittee *PolygondatacommitteeTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Polygondatacommittee.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_Polygondatacommittee *PolygondatacommitteeSession) Initialize() (*types.Transaction, error) {
	return _Polygondatacommittee.Contract.Initialize(&_Polygondatacommittee.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_Polygondatacommittee *PolygondatacommitteeTransactorSession) Initialize() (*types.Transaction, error) {
	return _Polygondatacommittee.Contract.Initialize(&_Polygondatacommittee.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Polygondatacommittee *PolygondatacommitteeTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Polygondatacommittee.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Polygondatacommittee *PolygondatacommitteeSession) RenounceOwnership() (*types.Transaction, error) {
	return _Polygondatacommittee.Contract.RenounceOwnership(&_Polygondatacommittee.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Polygondatacommittee *PolygondatacommitteeTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Polygondatacommittee.Contract.RenounceOwnership(&_Polygondatacommittee.TransactOpts)
}

// SetupCommittee is a paid mutator transaction binding the contract method 0x078fba2a.
//
// Solidity: function setupCommittee(uint256 _requiredAmountOfSignatures, string[] urls, bytes addrsBytes) returns()
func (_Polygondatacommittee *PolygondatacommitteeTransactor) SetupCommittee(opts *bind.TransactOpts, _requiredAmountOfSignatures *big.Int, urls []string, addrsBytes []byte) (*types.Transaction, error) {
	return _Polygondatacommittee.contract.Transact(opts, "setupCommittee", _requiredAmountOfSignatures, urls, addrsBytes)
}

// SetupCommittee is a paid mutator transaction binding the contract method 0x078fba2a.
//
// Solidity: function setupCommittee(uint256 _requiredAmountOfSignatures, string[] urls, bytes addrsBytes) returns()
func (_Polygondatacommittee *PolygondatacommitteeSession) SetupCommittee(_requiredAmountOfSignatures *big.Int, urls []string, addrsBytes []byte) (*types.Transaction, error) {
	return _Polygondatacommittee.Contract.SetupCommittee(&_Polygondatacommittee.TransactOpts, _requiredAmountOfSignatures, urls, addrsBytes)
}

// SetupCommittee is a paid mutator transaction binding the contract method 0x078fba2a.
//
// Solidity: function setupCommittee(uint256 _requiredAmountOfSignatures, string[] urls, bytes addrsBytes) returns()
func (_Polygondatacommittee *PolygondatacommitteeTransactorSession) SetupCommittee(_requiredAmountOfSignatures *big.Int, urls []string, addrsBytes []byte) (*types.Transaction, error) {
	return _Polygondatacommittee.Contract.SetupCommittee(&_Polygondatacommittee.TransactOpts, _requiredAmountOfSignatures, urls, addrsBytes)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Polygondatacommittee *PolygondatacommitteeTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Polygondatacommittee.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Polygondatacommittee *PolygondatacommitteeSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Polygondatacommittee.Contract.TransferOwnership(&_Polygondatacommittee.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Polygondatacommittee *PolygondatacommitteeTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Polygondatacommittee.Contract.TransferOwnership(&_Polygondatacommittee.TransactOpts, newOwner)
}

// PolygondatacommitteeCommitteeUpdatedIterator is returned from FilterCommitteeUpdated and is used to iterate over the raw logs and unpacked data for CommitteeUpdated events raised by the Polygondatacommittee contract.
type PolygondatacommitteeCommitteeUpdatedIterator struct {
	Event *PolygondatacommitteeCommitteeUpdated // Event containing the contract specifics and raw log

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
func (it *PolygondatacommitteeCommitteeUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolygondatacommitteeCommitteeUpdated)
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
		it.Event = new(PolygondatacommitteeCommitteeUpdated)
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
func (it *PolygondatacommitteeCommitteeUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolygondatacommitteeCommitteeUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolygondatacommitteeCommitteeUpdated represents a CommitteeUpdated event raised by the Polygondatacommittee contract.
type PolygondatacommitteeCommitteeUpdated struct {
	CommitteeHash [32]byte
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterCommitteeUpdated is a free log retrieval operation binding the contract event 0x831403fd381b3e6ac875d912ec2eee0e0203d0d29f7b3e0c96fc8f582d6db657.
//
// Solidity: event CommitteeUpdated(bytes32 committeeHash)
func (_Polygondatacommittee *PolygondatacommitteeFilterer) FilterCommitteeUpdated(opts *bind.FilterOpts) (*PolygondatacommitteeCommitteeUpdatedIterator, error) {

	logs, sub, err := _Polygondatacommittee.contract.FilterLogs(opts, "CommitteeUpdated")
	if err != nil {
		return nil, err
	}
	return &PolygondatacommitteeCommitteeUpdatedIterator{contract: _Polygondatacommittee.contract, event: "CommitteeUpdated", logs: logs, sub: sub}, nil
}

// WatchCommitteeUpdated is a free log subscription operation binding the contract event 0x831403fd381b3e6ac875d912ec2eee0e0203d0d29f7b3e0c96fc8f582d6db657.
//
// Solidity: event CommitteeUpdated(bytes32 committeeHash)
func (_Polygondatacommittee *PolygondatacommitteeFilterer) WatchCommitteeUpdated(opts *bind.WatchOpts, sink chan<- *PolygondatacommitteeCommitteeUpdated) (event.Subscription, error) {

	logs, sub, err := _Polygondatacommittee.contract.WatchLogs(opts, "CommitteeUpdated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolygondatacommitteeCommitteeUpdated)
				if err := _Polygondatacommittee.contract.UnpackLog(event, "CommitteeUpdated", log); err != nil {
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
func (_Polygondatacommittee *PolygondatacommitteeFilterer) ParseCommitteeUpdated(log types.Log) (*PolygondatacommitteeCommitteeUpdated, error) {
	event := new(PolygondatacommitteeCommitteeUpdated)
	if err := _Polygondatacommittee.contract.UnpackLog(event, "CommitteeUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolygondatacommitteeInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Polygondatacommittee contract.
type PolygondatacommitteeInitializedIterator struct {
	Event *PolygondatacommitteeInitialized // Event containing the contract specifics and raw log

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
func (it *PolygondatacommitteeInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolygondatacommitteeInitialized)
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
		it.Event = new(PolygondatacommitteeInitialized)
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
func (it *PolygondatacommitteeInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolygondatacommitteeInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolygondatacommitteeInitialized represents a Initialized event raised by the Polygondatacommittee contract.
type PolygondatacommitteeInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Polygondatacommittee *PolygondatacommitteeFilterer) FilterInitialized(opts *bind.FilterOpts) (*PolygondatacommitteeInitializedIterator, error) {

	logs, sub, err := _Polygondatacommittee.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &PolygondatacommitteeInitializedIterator{contract: _Polygondatacommittee.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Polygondatacommittee *PolygondatacommitteeFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *PolygondatacommitteeInitialized) (event.Subscription, error) {

	logs, sub, err := _Polygondatacommittee.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolygondatacommitteeInitialized)
				if err := _Polygondatacommittee.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Polygondatacommittee *PolygondatacommitteeFilterer) ParseInitialized(log types.Log) (*PolygondatacommitteeInitialized, error) {
	event := new(PolygondatacommitteeInitialized)
	if err := _Polygondatacommittee.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PolygondatacommitteeOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Polygondatacommittee contract.
type PolygondatacommitteeOwnershipTransferredIterator struct {
	Event *PolygondatacommitteeOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *PolygondatacommitteeOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PolygondatacommitteeOwnershipTransferred)
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
		it.Event = new(PolygondatacommitteeOwnershipTransferred)
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
func (it *PolygondatacommitteeOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PolygondatacommitteeOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PolygondatacommitteeOwnershipTransferred represents a OwnershipTransferred event raised by the Polygondatacommittee contract.
type PolygondatacommitteeOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Polygondatacommittee *PolygondatacommitteeFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*PolygondatacommitteeOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Polygondatacommittee.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &PolygondatacommitteeOwnershipTransferredIterator{contract: _Polygondatacommittee.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Polygondatacommittee *PolygondatacommitteeFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *PolygondatacommitteeOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Polygondatacommittee.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PolygondatacommitteeOwnershipTransferred)
				if err := _Polygondatacommittee.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_Polygondatacommittee *PolygondatacommitteeFilterer) ParseOwnershipTransferred(log types.Log) (*PolygondatacommitteeOwnershipTransferred, error) {
	event := new(PolygondatacommitteeOwnershipTransferred)
	if err := _Polygondatacommittee.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
