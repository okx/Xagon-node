// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package feijoapolygonzkevmglobalexitroot

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

// FeijoapolygonzkevmglobalexitrootMetaData contains all meta data concerning the Feijoapolygonzkevmglobalexitroot contract.
var FeijoapolygonzkevmglobalexitrootMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_rollupManager\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_bridgeAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"MerkleTreeFull\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyAllowedContracts\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"mainnetExitRoot\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"rollupExitRoot\",\"type\":\"bytes32\"}],\"name\":\"UpdateL1InfoTreeRecursive\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"bridgeAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"leafHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[32]\",\"name\":\"smtProof\",\"type\":\"bytes32[32]\"},{\"internalType\":\"uint32\",\"name\":\"index\",\"type\":\"uint32\"}],\"name\":\"calculateRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"depositCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"newGlobalExitRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"lastBlockHash\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"timestamp\",\"type\":\"uint64\"}],\"name\":\"getL1InfoTreeHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLastGlobalExitRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"l1InfoRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"l1InfoTreeHash\",\"type\":\"bytes32\"}],\"name\":\"getLeafValue\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"globalExitRootMap\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"depositCount\",\"type\":\"uint256\"}],\"name\":\"l1InfoLeafMap\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"l1InfoLeafHash\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastMainnetExitRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastRollupExitRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rollupManager\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"newRoot\",\"type\":\"bytes32\"}],\"name\":\"updateExitRoot\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"leafHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[32]\",\"name\":\"smtProof\",\"type\":\"bytes32[32]\"},{\"internalType\":\"uint32\",\"name\":\"index\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"root\",\"type\":\"bytes32\"}],\"name\":\"verifyMerkleProof\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x60c060405234801561000f575f80fd5b50604051610e6e380380610e6e83398101604081905261002e91610060565b6001600160a01b0391821660a05216608052610091565b80516001600160a01b038116811461005b575f80fd5b919050565b5f8060408385031215610071575f80fd5b61007a83610045565b915061008860208401610045565b90509250929050565b60805160a051610dae6100c05f395f8181610181015261031501525f818161026b01526102c90152610dae5ff3fe608060405234801561000f575f80fd5b50600436106100f0575f3560e01c806349b7b802116100935780638129fc1c116100635780638129fc1c1461024b57806383f2440314610253578063a3c573eb14610266578063fb5708341461028d575f80fd5b806349b7b8021461017c5780635ca1e165146101c85780635e0bd481146101d057806365f438d0146101e3575f80fd5b80632dfdf0b5116100ce5780632dfdf0b51461014d578063319cf7351461015657806333d6247d1461015f5780633ed691ef14610174575f80fd5b806301fd9044146100f4578063257b36321461010f57806325eaabaf1461012e575b5f80fd5b6100fc5f5481565b6040519081526020015b60405180910390f35b6100fc61011d366004610a54565b60026020525f908152604090205481565b6100fc61013c366004610a54565b602f6020525f908152604090205481565b6100fc60235481565b6100fc60015481565b61017261016d366004610a54565b6102b0565b005b6100fc610474565b6101a37f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001610106565b6100fc610487565b6100fc6101de366004610a6b565b610490565b6100fc6101f1366004610a8b565b604080516020808201959095528082019390935260c09190911b7fffffffffffffffff0000000000000000000000000000000000000000000000001660608301528051604881840301815260689092019052805191012090565b6101726104bf565b6100fc610261366004610af6565b61076c565b6101a37f000000000000000000000000000000000000000000000000000000000000000081565b6102a061029b366004610b32565b610837565b6040519015158152602001610106565b5f8073ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001633036102fe57505060018190555f548161037d565b73ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016330361034b5750505f819055600154819061037d565b6040517fb49365dd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f6103888284610490565b5f818152600260205260408120549192500361046e575f6103aa600143610ba4565b5f83815260026020908152604080832093409384905580518083018790528082018590527fffffffffffffffff0000000000000000000000000000000000000000000000004260c01b166060820152815180820360480181526068909101909152805191012091925090610421905b6101de610487565b6023545f908152602f60205260409020819055905061043f8161084e565b604051859085907f99d6f7ca42aad690b3768da4a5166fda058e4d023aea6eb922a08295c46360c4905f90a350505b50505050565b5f6104826001545f54610490565b905090565b5f610482610963565b604080516020808201859052818301849052825180830384018152606090920190925280519101205b92915050565b602e54610100900460ff16158080156104df5750602e54600160ff909116105b806104f95750303b1580156104f95750602e5460ff166001145b610589576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201527f647920696e697469616c697a6564000000000000000000000000000000000000606482015260840160405180910390fd5b602e80547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016600117905580156105e757602e80547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff166101001790555b602354156106fd575f5b6020811015610619576003816020811061060d5761060d610bb7565b5f9101556001016105f1565b505f6023819055610628610474565b90505f610636600143610ba4565b4090505f61069d610419848442604080516020808201959095528082019390935260c09190911b7fffffffffffffffff0000000000000000000000000000000000000000000000001660608301528051604881840301815260689092019052805191012090565b90506106a85f61084e565b6023545f908152602f602052604090208190556106c48161084e565b5f546001547f99d6f7ca42aad690b3768da4a5166fda058e4d023aea6eb922a08295c46360c460405160405180910390a3505050610706565b6107065f61084e565b801561076957602e80547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b50565b5f83815b602081101561082e57600163ffffffff8516821c811690036107db5784816020811061079e5761079e610bb7565b6020020135826040516020016107be929190918252602082015260400190565b604051602081830303815290604052805190602001209150610826565b818582602081106107ee576107ee610bb7565b602002013560405160200161080d929190918252602082015260400190565b6040516020818303038152906040528051906020012091505b600101610770565b50949350505050565b5f8161084486868661076c565b1495945050505050565b80600161085d60206002610d02565b6108679190610ba4565b602354106108a1576040517fef5ccf6600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f60235f81546108b090610d14565b918290555090505f5b6020811015610955578082901c6001166001036108ec5782600382602081106108e4576108e4610bb7565b015550505050565b600381602081106108ff576108ff610bb7565b01546040805160208101929092528101849052606001604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152919052805160209091012092506001016108b9565b5061095e610d4b565b505050565b6023545f90819081805b6020811015610a4b578083901c6001166001036109ca576003816020811061099757610997610bb7565b015460408051602081019290925281018590526060016040516020818303038152906040528051906020012093506109f7565b60408051602081018690529081018390526060016040516020818303038152906040528051906020012093505b6040805160208101849052908101839052606001604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190528051602090910120915060010161096d565b50919392505050565b5f60208284031215610a64575f80fd5b5035919050565b5f8060408385031215610a7c575f80fd5b50508035926020909101359150565b5f805f60608486031215610a9d575f80fd5b8335925060208401359150604084013567ffffffffffffffff81168114610ac2575f80fd5b809150509250925092565b8061040081018310156104b9575f80fd5b803563ffffffff81168114610af1575f80fd5b919050565b5f805f6104408486031215610b09575f80fd5b83359250610b1a8560208601610acd565b9150610b296104208501610ade565b90509250925092565b5f805f806104608587031215610b46575f80fd5b84359350610b578660208701610acd565b9250610b666104208601610ade565b939692955092936104400135925050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b818103818111156104b9576104b9610b77565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b600181815b80851115610c3d57817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04821115610c2357610c23610b77565b80851615610c3057918102915b93841c9390800290610be9565b509250929050565b5f82610c53575060016104b9565b81610c5f57505f6104b9565b8160018114610c755760028114610c7f57610c9b565b60019150506104b9565b60ff841115610c9057610c90610b77565b50506001821b6104b9565b5060208310610133831016604e8410600b8410161715610cbe575081810a6104b9565b610cc88383610be4565b807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff04821115610cfa57610cfa610b77565b029392505050565b5f610d0d8383610c45565b9392505050565b5f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff8203610d4457610d44610b77565b5060010190565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52600160045260245ffdfea2646970667358221220871578336aa6be75bc24afcd4fe1fa219065615aec285b318915fb767e5e011c64736f6c63430008180033",
}

// FeijoapolygonzkevmglobalexitrootABI is the input ABI used to generate the binding from.
// Deprecated: Use FeijoapolygonzkevmglobalexitrootMetaData.ABI instead.
var FeijoapolygonzkevmglobalexitrootABI = FeijoapolygonzkevmglobalexitrootMetaData.ABI

// FeijoapolygonzkevmglobalexitrootBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use FeijoapolygonzkevmglobalexitrootMetaData.Bin instead.
var FeijoapolygonzkevmglobalexitrootBin = FeijoapolygonzkevmglobalexitrootMetaData.Bin

// DeployFeijoapolygonzkevmglobalexitroot deploys a new Ethereum contract, binding an instance of Feijoapolygonzkevmglobalexitroot to it.
func DeployFeijoapolygonzkevmglobalexitroot(auth *bind.TransactOpts, backend bind.ContractBackend, _rollupManager common.Address, _bridgeAddress common.Address) (common.Address, *types.Transaction, *Feijoapolygonzkevmglobalexitroot, error) {
	parsed, err := FeijoapolygonzkevmglobalexitrootMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(FeijoapolygonzkevmglobalexitrootBin), backend, _rollupManager, _bridgeAddress)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Feijoapolygonzkevmglobalexitroot{FeijoapolygonzkevmglobalexitrootCaller: FeijoapolygonzkevmglobalexitrootCaller{contract: contract}, FeijoapolygonzkevmglobalexitrootTransactor: FeijoapolygonzkevmglobalexitrootTransactor{contract: contract}, FeijoapolygonzkevmglobalexitrootFilterer: FeijoapolygonzkevmglobalexitrootFilterer{contract: contract}}, nil
}

// Feijoapolygonzkevmglobalexitroot is an auto generated Go binding around an Ethereum contract.
type Feijoapolygonzkevmglobalexitroot struct {
	FeijoapolygonzkevmglobalexitrootCaller     // Read-only binding to the contract
	FeijoapolygonzkevmglobalexitrootTransactor // Write-only binding to the contract
	FeijoapolygonzkevmglobalexitrootFilterer   // Log filterer for contract events
}

// FeijoapolygonzkevmglobalexitrootCaller is an auto generated read-only Go binding around an Ethereum contract.
type FeijoapolygonzkevmglobalexitrootCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FeijoapolygonzkevmglobalexitrootTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FeijoapolygonzkevmglobalexitrootTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FeijoapolygonzkevmglobalexitrootFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FeijoapolygonzkevmglobalexitrootFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FeijoapolygonzkevmglobalexitrootSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FeijoapolygonzkevmglobalexitrootSession struct {
	Contract     *Feijoapolygonzkevmglobalexitroot // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                     // Call options to use throughout this session
	TransactOpts bind.TransactOpts                 // Transaction auth options to use throughout this session
}

// FeijoapolygonzkevmglobalexitrootCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FeijoapolygonzkevmglobalexitrootCallerSession struct {
	Contract *FeijoapolygonzkevmglobalexitrootCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                           // Call options to use throughout this session
}

// FeijoapolygonzkevmglobalexitrootTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FeijoapolygonzkevmglobalexitrootTransactorSession struct {
	Contract     *FeijoapolygonzkevmglobalexitrootTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                           // Transaction auth options to use throughout this session
}

// FeijoapolygonzkevmglobalexitrootRaw is an auto generated low-level Go binding around an Ethereum contract.
type FeijoapolygonzkevmglobalexitrootRaw struct {
	Contract *Feijoapolygonzkevmglobalexitroot // Generic contract binding to access the raw methods on
}

// FeijoapolygonzkevmglobalexitrootCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FeijoapolygonzkevmglobalexitrootCallerRaw struct {
	Contract *FeijoapolygonzkevmglobalexitrootCaller // Generic read-only contract binding to access the raw methods on
}

// FeijoapolygonzkevmglobalexitrootTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FeijoapolygonzkevmglobalexitrootTransactorRaw struct {
	Contract *FeijoapolygonzkevmglobalexitrootTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFeijoapolygonzkevmglobalexitroot creates a new instance of Feijoapolygonzkevmglobalexitroot, bound to a specific deployed contract.
func NewFeijoapolygonzkevmglobalexitroot(address common.Address, backend bind.ContractBackend) (*Feijoapolygonzkevmglobalexitroot, error) {
	contract, err := bindFeijoapolygonzkevmglobalexitroot(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Feijoapolygonzkevmglobalexitroot{FeijoapolygonzkevmglobalexitrootCaller: FeijoapolygonzkevmglobalexitrootCaller{contract: contract}, FeijoapolygonzkevmglobalexitrootTransactor: FeijoapolygonzkevmglobalexitrootTransactor{contract: contract}, FeijoapolygonzkevmglobalexitrootFilterer: FeijoapolygonzkevmglobalexitrootFilterer{contract: contract}}, nil
}

// NewFeijoapolygonzkevmglobalexitrootCaller creates a new read-only instance of Feijoapolygonzkevmglobalexitroot, bound to a specific deployed contract.
func NewFeijoapolygonzkevmglobalexitrootCaller(address common.Address, caller bind.ContractCaller) (*FeijoapolygonzkevmglobalexitrootCaller, error) {
	contract, err := bindFeijoapolygonzkevmglobalexitroot(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FeijoapolygonzkevmglobalexitrootCaller{contract: contract}, nil
}

// NewFeijoapolygonzkevmglobalexitrootTransactor creates a new write-only instance of Feijoapolygonzkevmglobalexitroot, bound to a specific deployed contract.
func NewFeijoapolygonzkevmglobalexitrootTransactor(address common.Address, transactor bind.ContractTransactor) (*FeijoapolygonzkevmglobalexitrootTransactor, error) {
	contract, err := bindFeijoapolygonzkevmglobalexitroot(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FeijoapolygonzkevmglobalexitrootTransactor{contract: contract}, nil
}

// NewFeijoapolygonzkevmglobalexitrootFilterer creates a new log filterer instance of Feijoapolygonzkevmglobalexitroot, bound to a specific deployed contract.
func NewFeijoapolygonzkevmglobalexitrootFilterer(address common.Address, filterer bind.ContractFilterer) (*FeijoapolygonzkevmglobalexitrootFilterer, error) {
	contract, err := bindFeijoapolygonzkevmglobalexitroot(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FeijoapolygonzkevmglobalexitrootFilterer{contract: contract}, nil
}

// bindFeijoapolygonzkevmglobalexitroot binds a generic wrapper to an already deployed contract.
func bindFeijoapolygonzkevmglobalexitroot(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FeijoapolygonzkevmglobalexitrootMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Feijoapolygonzkevmglobalexitroot.Contract.FeijoapolygonzkevmglobalexitrootCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Feijoapolygonzkevmglobalexitroot.Contract.FeijoapolygonzkevmglobalexitrootTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Feijoapolygonzkevmglobalexitroot.Contract.FeijoapolygonzkevmglobalexitrootTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Feijoapolygonzkevmglobalexitroot.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Feijoapolygonzkevmglobalexitroot.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Feijoapolygonzkevmglobalexitroot.Contract.contract.Transact(opts, method, params...)
}

// BridgeAddress is a free data retrieval call binding the contract method 0xa3c573eb.
//
// Solidity: function bridgeAddress() view returns(address)
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootCaller) BridgeAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Feijoapolygonzkevmglobalexitroot.contract.Call(opts, &out, "bridgeAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BridgeAddress is a free data retrieval call binding the contract method 0xa3c573eb.
//
// Solidity: function bridgeAddress() view returns(address)
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootSession) BridgeAddress() (common.Address, error) {
	return _Feijoapolygonzkevmglobalexitroot.Contract.BridgeAddress(&_Feijoapolygonzkevmglobalexitroot.CallOpts)
}

// BridgeAddress is a free data retrieval call binding the contract method 0xa3c573eb.
//
// Solidity: function bridgeAddress() view returns(address)
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootCallerSession) BridgeAddress() (common.Address, error) {
	return _Feijoapolygonzkevmglobalexitroot.Contract.BridgeAddress(&_Feijoapolygonzkevmglobalexitroot.CallOpts)
}

// CalculateRoot is a free data retrieval call binding the contract method 0x83f24403.
//
// Solidity: function calculateRoot(bytes32 leafHash, bytes32[32] smtProof, uint32 index) pure returns(bytes32)
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootCaller) CalculateRoot(opts *bind.CallOpts, leafHash [32]byte, smtProof [32][32]byte, index uint32) ([32]byte, error) {
	var out []interface{}
	err := _Feijoapolygonzkevmglobalexitroot.contract.Call(opts, &out, "calculateRoot", leafHash, smtProof, index)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// CalculateRoot is a free data retrieval call binding the contract method 0x83f24403.
//
// Solidity: function calculateRoot(bytes32 leafHash, bytes32[32] smtProof, uint32 index) pure returns(bytes32)
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootSession) CalculateRoot(leafHash [32]byte, smtProof [32][32]byte, index uint32) ([32]byte, error) {
	return _Feijoapolygonzkevmglobalexitroot.Contract.CalculateRoot(&_Feijoapolygonzkevmglobalexitroot.CallOpts, leafHash, smtProof, index)
}

// CalculateRoot is a free data retrieval call binding the contract method 0x83f24403.
//
// Solidity: function calculateRoot(bytes32 leafHash, bytes32[32] smtProof, uint32 index) pure returns(bytes32)
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootCallerSession) CalculateRoot(leafHash [32]byte, smtProof [32][32]byte, index uint32) ([32]byte, error) {
	return _Feijoapolygonzkevmglobalexitroot.Contract.CalculateRoot(&_Feijoapolygonzkevmglobalexitroot.CallOpts, leafHash, smtProof, index)
}

// DepositCount is a free data retrieval call binding the contract method 0x2dfdf0b5.
//
// Solidity: function depositCount() view returns(uint256)
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootCaller) DepositCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Feijoapolygonzkevmglobalexitroot.contract.Call(opts, &out, "depositCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DepositCount is a free data retrieval call binding the contract method 0x2dfdf0b5.
//
// Solidity: function depositCount() view returns(uint256)
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootSession) DepositCount() (*big.Int, error) {
	return _Feijoapolygonzkevmglobalexitroot.Contract.DepositCount(&_Feijoapolygonzkevmglobalexitroot.CallOpts)
}

// DepositCount is a free data retrieval call binding the contract method 0x2dfdf0b5.
//
// Solidity: function depositCount() view returns(uint256)
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootCallerSession) DepositCount() (*big.Int, error) {
	return _Feijoapolygonzkevmglobalexitroot.Contract.DepositCount(&_Feijoapolygonzkevmglobalexitroot.CallOpts)
}

// GetL1InfoTreeHash is a free data retrieval call binding the contract method 0x65f438d0.
//
// Solidity: function getL1InfoTreeHash(bytes32 newGlobalExitRoot, uint256 lastBlockHash, uint64 timestamp) pure returns(bytes32)
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootCaller) GetL1InfoTreeHash(opts *bind.CallOpts, newGlobalExitRoot [32]byte, lastBlockHash *big.Int, timestamp uint64) ([32]byte, error) {
	var out []interface{}
	err := _Feijoapolygonzkevmglobalexitroot.contract.Call(opts, &out, "getL1InfoTreeHash", newGlobalExitRoot, lastBlockHash, timestamp)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetL1InfoTreeHash is a free data retrieval call binding the contract method 0x65f438d0.
//
// Solidity: function getL1InfoTreeHash(bytes32 newGlobalExitRoot, uint256 lastBlockHash, uint64 timestamp) pure returns(bytes32)
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootSession) GetL1InfoTreeHash(newGlobalExitRoot [32]byte, lastBlockHash *big.Int, timestamp uint64) ([32]byte, error) {
	return _Feijoapolygonzkevmglobalexitroot.Contract.GetL1InfoTreeHash(&_Feijoapolygonzkevmglobalexitroot.CallOpts, newGlobalExitRoot, lastBlockHash, timestamp)
}

// GetL1InfoTreeHash is a free data retrieval call binding the contract method 0x65f438d0.
//
// Solidity: function getL1InfoTreeHash(bytes32 newGlobalExitRoot, uint256 lastBlockHash, uint64 timestamp) pure returns(bytes32)
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootCallerSession) GetL1InfoTreeHash(newGlobalExitRoot [32]byte, lastBlockHash *big.Int, timestamp uint64) ([32]byte, error) {
	return _Feijoapolygonzkevmglobalexitroot.Contract.GetL1InfoTreeHash(&_Feijoapolygonzkevmglobalexitroot.CallOpts, newGlobalExitRoot, lastBlockHash, timestamp)
}

// GetLastGlobalExitRoot is a free data retrieval call binding the contract method 0x3ed691ef.
//
// Solidity: function getLastGlobalExitRoot() view returns(bytes32)
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootCaller) GetLastGlobalExitRoot(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Feijoapolygonzkevmglobalexitroot.contract.Call(opts, &out, "getLastGlobalExitRoot")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetLastGlobalExitRoot is a free data retrieval call binding the contract method 0x3ed691ef.
//
// Solidity: function getLastGlobalExitRoot() view returns(bytes32)
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootSession) GetLastGlobalExitRoot() ([32]byte, error) {
	return _Feijoapolygonzkevmglobalexitroot.Contract.GetLastGlobalExitRoot(&_Feijoapolygonzkevmglobalexitroot.CallOpts)
}

// GetLastGlobalExitRoot is a free data retrieval call binding the contract method 0x3ed691ef.
//
// Solidity: function getLastGlobalExitRoot() view returns(bytes32)
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootCallerSession) GetLastGlobalExitRoot() ([32]byte, error) {
	return _Feijoapolygonzkevmglobalexitroot.Contract.GetLastGlobalExitRoot(&_Feijoapolygonzkevmglobalexitroot.CallOpts)
}

// GetLeafValue is a free data retrieval call binding the contract method 0x5e0bd481.
//
// Solidity: function getLeafValue(bytes32 l1InfoRoot, bytes32 l1InfoTreeHash) pure returns(bytes32)
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootCaller) GetLeafValue(opts *bind.CallOpts, l1InfoRoot [32]byte, l1InfoTreeHash [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _Feijoapolygonzkevmglobalexitroot.contract.Call(opts, &out, "getLeafValue", l1InfoRoot, l1InfoTreeHash)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetLeafValue is a free data retrieval call binding the contract method 0x5e0bd481.
//
// Solidity: function getLeafValue(bytes32 l1InfoRoot, bytes32 l1InfoTreeHash) pure returns(bytes32)
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootSession) GetLeafValue(l1InfoRoot [32]byte, l1InfoTreeHash [32]byte) ([32]byte, error) {
	return _Feijoapolygonzkevmglobalexitroot.Contract.GetLeafValue(&_Feijoapolygonzkevmglobalexitroot.CallOpts, l1InfoRoot, l1InfoTreeHash)
}

// GetLeafValue is a free data retrieval call binding the contract method 0x5e0bd481.
//
// Solidity: function getLeafValue(bytes32 l1InfoRoot, bytes32 l1InfoTreeHash) pure returns(bytes32)
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootCallerSession) GetLeafValue(l1InfoRoot [32]byte, l1InfoTreeHash [32]byte) ([32]byte, error) {
	return _Feijoapolygonzkevmglobalexitroot.Contract.GetLeafValue(&_Feijoapolygonzkevmglobalexitroot.CallOpts, l1InfoRoot, l1InfoTreeHash)
}

// GetRoot is a free data retrieval call binding the contract method 0x5ca1e165.
//
// Solidity: function getRoot() view returns(bytes32)
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootCaller) GetRoot(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Feijoapolygonzkevmglobalexitroot.contract.Call(opts, &out, "getRoot")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoot is a free data retrieval call binding the contract method 0x5ca1e165.
//
// Solidity: function getRoot() view returns(bytes32)
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootSession) GetRoot() ([32]byte, error) {
	return _Feijoapolygonzkevmglobalexitroot.Contract.GetRoot(&_Feijoapolygonzkevmglobalexitroot.CallOpts)
}

// GetRoot is a free data retrieval call binding the contract method 0x5ca1e165.
//
// Solidity: function getRoot() view returns(bytes32)
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootCallerSession) GetRoot() ([32]byte, error) {
	return _Feijoapolygonzkevmglobalexitroot.Contract.GetRoot(&_Feijoapolygonzkevmglobalexitroot.CallOpts)
}

// GlobalExitRootMap is a free data retrieval call binding the contract method 0x257b3632.
//
// Solidity: function globalExitRootMap(bytes32 ) view returns(uint256)
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootCaller) GlobalExitRootMap(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _Feijoapolygonzkevmglobalexitroot.contract.Call(opts, &out, "globalExitRootMap", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GlobalExitRootMap is a free data retrieval call binding the contract method 0x257b3632.
//
// Solidity: function globalExitRootMap(bytes32 ) view returns(uint256)
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootSession) GlobalExitRootMap(arg0 [32]byte) (*big.Int, error) {
	return _Feijoapolygonzkevmglobalexitroot.Contract.GlobalExitRootMap(&_Feijoapolygonzkevmglobalexitroot.CallOpts, arg0)
}

// GlobalExitRootMap is a free data retrieval call binding the contract method 0x257b3632.
//
// Solidity: function globalExitRootMap(bytes32 ) view returns(uint256)
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootCallerSession) GlobalExitRootMap(arg0 [32]byte) (*big.Int, error) {
	return _Feijoapolygonzkevmglobalexitroot.Contract.GlobalExitRootMap(&_Feijoapolygonzkevmglobalexitroot.CallOpts, arg0)
}

// L1InfoLeafMap is a free data retrieval call binding the contract method 0x25eaabaf.
//
// Solidity: function l1InfoLeafMap(uint256 depositCount) view returns(bytes32 l1InfoLeafHash)
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootCaller) L1InfoLeafMap(opts *bind.CallOpts, depositCount *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _Feijoapolygonzkevmglobalexitroot.contract.Call(opts, &out, "l1InfoLeafMap", depositCount)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// L1InfoLeafMap is a free data retrieval call binding the contract method 0x25eaabaf.
//
// Solidity: function l1InfoLeafMap(uint256 depositCount) view returns(bytes32 l1InfoLeafHash)
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootSession) L1InfoLeafMap(depositCount *big.Int) ([32]byte, error) {
	return _Feijoapolygonzkevmglobalexitroot.Contract.L1InfoLeafMap(&_Feijoapolygonzkevmglobalexitroot.CallOpts, depositCount)
}

// L1InfoLeafMap is a free data retrieval call binding the contract method 0x25eaabaf.
//
// Solidity: function l1InfoLeafMap(uint256 depositCount) view returns(bytes32 l1InfoLeafHash)
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootCallerSession) L1InfoLeafMap(depositCount *big.Int) ([32]byte, error) {
	return _Feijoapolygonzkevmglobalexitroot.Contract.L1InfoLeafMap(&_Feijoapolygonzkevmglobalexitroot.CallOpts, depositCount)
}

// LastMainnetExitRoot is a free data retrieval call binding the contract method 0x319cf735.
//
// Solidity: function lastMainnetExitRoot() view returns(bytes32)
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootCaller) LastMainnetExitRoot(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Feijoapolygonzkevmglobalexitroot.contract.Call(opts, &out, "lastMainnetExitRoot")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// LastMainnetExitRoot is a free data retrieval call binding the contract method 0x319cf735.
//
// Solidity: function lastMainnetExitRoot() view returns(bytes32)
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootSession) LastMainnetExitRoot() ([32]byte, error) {
	return _Feijoapolygonzkevmglobalexitroot.Contract.LastMainnetExitRoot(&_Feijoapolygonzkevmglobalexitroot.CallOpts)
}

// LastMainnetExitRoot is a free data retrieval call binding the contract method 0x319cf735.
//
// Solidity: function lastMainnetExitRoot() view returns(bytes32)
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootCallerSession) LastMainnetExitRoot() ([32]byte, error) {
	return _Feijoapolygonzkevmglobalexitroot.Contract.LastMainnetExitRoot(&_Feijoapolygonzkevmglobalexitroot.CallOpts)
}

// LastRollupExitRoot is a free data retrieval call binding the contract method 0x01fd9044.
//
// Solidity: function lastRollupExitRoot() view returns(bytes32)
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootCaller) LastRollupExitRoot(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Feijoapolygonzkevmglobalexitroot.contract.Call(opts, &out, "lastRollupExitRoot")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// LastRollupExitRoot is a free data retrieval call binding the contract method 0x01fd9044.
//
// Solidity: function lastRollupExitRoot() view returns(bytes32)
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootSession) LastRollupExitRoot() ([32]byte, error) {
	return _Feijoapolygonzkevmglobalexitroot.Contract.LastRollupExitRoot(&_Feijoapolygonzkevmglobalexitroot.CallOpts)
}

// LastRollupExitRoot is a free data retrieval call binding the contract method 0x01fd9044.
//
// Solidity: function lastRollupExitRoot() view returns(bytes32)
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootCallerSession) LastRollupExitRoot() ([32]byte, error) {
	return _Feijoapolygonzkevmglobalexitroot.Contract.LastRollupExitRoot(&_Feijoapolygonzkevmglobalexitroot.CallOpts)
}

// RollupManager is a free data retrieval call binding the contract method 0x49b7b802.
//
// Solidity: function rollupManager() view returns(address)
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootCaller) RollupManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Feijoapolygonzkevmglobalexitroot.contract.Call(opts, &out, "rollupManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RollupManager is a free data retrieval call binding the contract method 0x49b7b802.
//
// Solidity: function rollupManager() view returns(address)
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootSession) RollupManager() (common.Address, error) {
	return _Feijoapolygonzkevmglobalexitroot.Contract.RollupManager(&_Feijoapolygonzkevmglobalexitroot.CallOpts)
}

// RollupManager is a free data retrieval call binding the contract method 0x49b7b802.
//
// Solidity: function rollupManager() view returns(address)
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootCallerSession) RollupManager() (common.Address, error) {
	return _Feijoapolygonzkevmglobalexitroot.Contract.RollupManager(&_Feijoapolygonzkevmglobalexitroot.CallOpts)
}

// VerifyMerkleProof is a free data retrieval call binding the contract method 0xfb570834.
//
// Solidity: function verifyMerkleProof(bytes32 leafHash, bytes32[32] smtProof, uint32 index, bytes32 root) pure returns(bool)
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootCaller) VerifyMerkleProof(opts *bind.CallOpts, leafHash [32]byte, smtProof [32][32]byte, index uint32, root [32]byte) (bool, error) {
	var out []interface{}
	err := _Feijoapolygonzkevmglobalexitroot.contract.Call(opts, &out, "verifyMerkleProof", leafHash, smtProof, index, root)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// VerifyMerkleProof is a free data retrieval call binding the contract method 0xfb570834.
//
// Solidity: function verifyMerkleProof(bytes32 leafHash, bytes32[32] smtProof, uint32 index, bytes32 root) pure returns(bool)
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootSession) VerifyMerkleProof(leafHash [32]byte, smtProof [32][32]byte, index uint32, root [32]byte) (bool, error) {
	return _Feijoapolygonzkevmglobalexitroot.Contract.VerifyMerkleProof(&_Feijoapolygonzkevmglobalexitroot.CallOpts, leafHash, smtProof, index, root)
}

// VerifyMerkleProof is a free data retrieval call binding the contract method 0xfb570834.
//
// Solidity: function verifyMerkleProof(bytes32 leafHash, bytes32[32] smtProof, uint32 index, bytes32 root) pure returns(bool)
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootCallerSession) VerifyMerkleProof(leafHash [32]byte, smtProof [32][32]byte, index uint32, root [32]byte) (bool, error) {
	return _Feijoapolygonzkevmglobalexitroot.Contract.VerifyMerkleProof(&_Feijoapolygonzkevmglobalexitroot.CallOpts, leafHash, smtProof, index, root)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Feijoapolygonzkevmglobalexitroot.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootSession) Initialize() (*types.Transaction, error) {
	return _Feijoapolygonzkevmglobalexitroot.Contract.Initialize(&_Feijoapolygonzkevmglobalexitroot.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootTransactorSession) Initialize() (*types.Transaction, error) {
	return _Feijoapolygonzkevmglobalexitroot.Contract.Initialize(&_Feijoapolygonzkevmglobalexitroot.TransactOpts)
}

// UpdateExitRoot is a paid mutator transaction binding the contract method 0x33d6247d.
//
// Solidity: function updateExitRoot(bytes32 newRoot) returns()
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootTransactor) UpdateExitRoot(opts *bind.TransactOpts, newRoot [32]byte) (*types.Transaction, error) {
	return _Feijoapolygonzkevmglobalexitroot.contract.Transact(opts, "updateExitRoot", newRoot)
}

// UpdateExitRoot is a paid mutator transaction binding the contract method 0x33d6247d.
//
// Solidity: function updateExitRoot(bytes32 newRoot) returns()
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootSession) UpdateExitRoot(newRoot [32]byte) (*types.Transaction, error) {
	return _Feijoapolygonzkevmglobalexitroot.Contract.UpdateExitRoot(&_Feijoapolygonzkevmglobalexitroot.TransactOpts, newRoot)
}

// UpdateExitRoot is a paid mutator transaction binding the contract method 0x33d6247d.
//
// Solidity: function updateExitRoot(bytes32 newRoot) returns()
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootTransactorSession) UpdateExitRoot(newRoot [32]byte) (*types.Transaction, error) {
	return _Feijoapolygonzkevmglobalexitroot.Contract.UpdateExitRoot(&_Feijoapolygonzkevmglobalexitroot.TransactOpts, newRoot)
}

// FeijoapolygonzkevmglobalexitrootInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Feijoapolygonzkevmglobalexitroot contract.
type FeijoapolygonzkevmglobalexitrootInitializedIterator struct {
	Event *FeijoapolygonzkevmglobalexitrootInitialized // Event containing the contract specifics and raw log

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
func (it *FeijoapolygonzkevmglobalexitrootInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeijoapolygonzkevmglobalexitrootInitialized)
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
		it.Event = new(FeijoapolygonzkevmglobalexitrootInitialized)
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
func (it *FeijoapolygonzkevmglobalexitrootInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeijoapolygonzkevmglobalexitrootInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeijoapolygonzkevmglobalexitrootInitialized represents a Initialized event raised by the Feijoapolygonzkevmglobalexitroot contract.
type FeijoapolygonzkevmglobalexitrootInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootFilterer) FilterInitialized(opts *bind.FilterOpts) (*FeijoapolygonzkevmglobalexitrootInitializedIterator, error) {

	logs, sub, err := _Feijoapolygonzkevmglobalexitroot.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &FeijoapolygonzkevmglobalexitrootInitializedIterator{contract: _Feijoapolygonzkevmglobalexitroot.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *FeijoapolygonzkevmglobalexitrootInitialized) (event.Subscription, error) {

	logs, sub, err := _Feijoapolygonzkevmglobalexitroot.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeijoapolygonzkevmglobalexitrootInitialized)
				if err := _Feijoapolygonzkevmglobalexitroot.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootFilterer) ParseInitialized(log types.Log) (*FeijoapolygonzkevmglobalexitrootInitialized, error) {
	event := new(FeijoapolygonzkevmglobalexitrootInitialized)
	if err := _Feijoapolygonzkevmglobalexitroot.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeijoapolygonzkevmglobalexitrootUpdateL1InfoTreeRecursiveIterator is returned from FilterUpdateL1InfoTreeRecursive and is used to iterate over the raw logs and unpacked data for UpdateL1InfoTreeRecursive events raised by the Feijoapolygonzkevmglobalexitroot contract.
type FeijoapolygonzkevmglobalexitrootUpdateL1InfoTreeRecursiveIterator struct {
	Event *FeijoapolygonzkevmglobalexitrootUpdateL1InfoTreeRecursive // Event containing the contract specifics and raw log

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
func (it *FeijoapolygonzkevmglobalexitrootUpdateL1InfoTreeRecursiveIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeijoapolygonzkevmglobalexitrootUpdateL1InfoTreeRecursive)
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
		it.Event = new(FeijoapolygonzkevmglobalexitrootUpdateL1InfoTreeRecursive)
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
func (it *FeijoapolygonzkevmglobalexitrootUpdateL1InfoTreeRecursiveIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeijoapolygonzkevmglobalexitrootUpdateL1InfoTreeRecursiveIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeijoapolygonzkevmglobalexitrootUpdateL1InfoTreeRecursive represents a UpdateL1InfoTreeRecursive event raised by the Feijoapolygonzkevmglobalexitroot contract.
type FeijoapolygonzkevmglobalexitrootUpdateL1InfoTreeRecursive struct {
	MainnetExitRoot [32]byte
	RollupExitRoot  [32]byte
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterUpdateL1InfoTreeRecursive is a free log retrieval operation binding the contract event 0x99d6f7ca42aad690b3768da4a5166fda058e4d023aea6eb922a08295c46360c4.
//
// Solidity: event UpdateL1InfoTreeRecursive(bytes32 indexed mainnetExitRoot, bytes32 indexed rollupExitRoot)
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootFilterer) FilterUpdateL1InfoTreeRecursive(opts *bind.FilterOpts, mainnetExitRoot [][32]byte, rollupExitRoot [][32]byte) (*FeijoapolygonzkevmglobalexitrootUpdateL1InfoTreeRecursiveIterator, error) {

	var mainnetExitRootRule []interface{}
	for _, mainnetExitRootItem := range mainnetExitRoot {
		mainnetExitRootRule = append(mainnetExitRootRule, mainnetExitRootItem)
	}
	var rollupExitRootRule []interface{}
	for _, rollupExitRootItem := range rollupExitRoot {
		rollupExitRootRule = append(rollupExitRootRule, rollupExitRootItem)
	}

	logs, sub, err := _Feijoapolygonzkevmglobalexitroot.contract.FilterLogs(opts, "UpdateL1InfoTreeRecursive", mainnetExitRootRule, rollupExitRootRule)
	if err != nil {
		return nil, err
	}
	return &FeijoapolygonzkevmglobalexitrootUpdateL1InfoTreeRecursiveIterator{contract: _Feijoapolygonzkevmglobalexitroot.contract, event: "UpdateL1InfoTreeRecursive", logs: logs, sub: sub}, nil
}

// WatchUpdateL1InfoTreeRecursive is a free log subscription operation binding the contract event 0x99d6f7ca42aad690b3768da4a5166fda058e4d023aea6eb922a08295c46360c4.
//
// Solidity: event UpdateL1InfoTreeRecursive(bytes32 indexed mainnetExitRoot, bytes32 indexed rollupExitRoot)
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootFilterer) WatchUpdateL1InfoTreeRecursive(opts *bind.WatchOpts, sink chan<- *FeijoapolygonzkevmglobalexitrootUpdateL1InfoTreeRecursive, mainnetExitRoot [][32]byte, rollupExitRoot [][32]byte) (event.Subscription, error) {

	var mainnetExitRootRule []interface{}
	for _, mainnetExitRootItem := range mainnetExitRoot {
		mainnetExitRootRule = append(mainnetExitRootRule, mainnetExitRootItem)
	}
	var rollupExitRootRule []interface{}
	for _, rollupExitRootItem := range rollupExitRoot {
		rollupExitRootRule = append(rollupExitRootRule, rollupExitRootItem)
	}

	logs, sub, err := _Feijoapolygonzkevmglobalexitroot.contract.WatchLogs(opts, "UpdateL1InfoTreeRecursive", mainnetExitRootRule, rollupExitRootRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeijoapolygonzkevmglobalexitrootUpdateL1InfoTreeRecursive)
				if err := _Feijoapolygonzkevmglobalexitroot.contract.UnpackLog(event, "UpdateL1InfoTreeRecursive", log); err != nil {
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

// ParseUpdateL1InfoTreeRecursive is a log parse operation binding the contract event 0x99d6f7ca42aad690b3768da4a5166fda058e4d023aea6eb922a08295c46360c4.
//
// Solidity: event UpdateL1InfoTreeRecursive(bytes32 indexed mainnetExitRoot, bytes32 indexed rollupExitRoot)
func (_Feijoapolygonzkevmglobalexitroot *FeijoapolygonzkevmglobalexitrootFilterer) ParseUpdateL1InfoTreeRecursive(log types.Log) (*FeijoapolygonzkevmglobalexitrootUpdateL1InfoTreeRecursive, error) {
	event := new(FeijoapolygonzkevmglobalexitrootUpdateL1InfoTreeRecursive)
	if err := _Feijoapolygonzkevmglobalexitroot.contract.UnpackLog(event, "UpdateL1InfoTreeRecursive", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
