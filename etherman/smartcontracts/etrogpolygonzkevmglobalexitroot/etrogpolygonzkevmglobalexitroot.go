// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package etrogpolygonzkevmglobalexitroot

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

// EtrogpolygonzkevmglobalexitrootMetaData contains all meta data concerning the Etrogpolygonzkevmglobalexitroot contract.
var EtrogpolygonzkevmglobalexitrootMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_rollupManager\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_bridgeAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"MerkleTreeFull\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyAllowedContracts\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"mainnetExitRoot\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"rollupExitRoot\",\"type\":\"bytes32\"}],\"name\":\"UpdateL1InfoTree\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"bridgeAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"leafHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[32]\",\"name\":\"smtProof\",\"type\":\"bytes32[32]\"},{\"internalType\":\"uint32\",\"name\":\"index\",\"type\":\"uint32\"}],\"name\":\"calculateRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"depositCount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getLastGlobalExitRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"newGlobalExitRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"lastBlockHash\",\"type\":\"uint256\"},{\"internalType\":\"uint64\",\"name\":\"timestamp\",\"type\":\"uint64\"}],\"name\":\"getLeafValue\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"globalExitRootMap\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastMainnetExitRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastRollupExitRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rollupManager\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"newRoot\",\"type\":\"bytes32\"}],\"name\":\"updateExitRoot\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"leafHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[32]\",\"name\":\"smtProof\",\"type\":\"bytes32[32]\"},{\"internalType\":\"uint32\",\"name\":\"index\",\"type\":\"uint32\"},{\"internalType\":\"bytes32\",\"name\":\"root\",\"type\":\"bytes32\"}],\"name\":\"verifyMerkleProof\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x60c060405234801561001057600080fd5b50604051610b3c380380610b3c83398101604081905261002f91610062565b6001600160a01b0391821660a05216608052610095565b80516001600160a01b038116811461005d57600080fd5b919050565b6000806040838503121561007557600080fd5b61007e83610046565b915061008c60208401610046565b90509250929050565b60805160a051610a746100c86000396000818161014901526102c401526000818161021801526102770152610a746000f3fe608060405234801561001057600080fd5b50600436106100d45760003560e01c806349b7b8021161008157806383f244031161005b57806383f2440314610200578063a3c573eb14610213578063fb5708341461023a57600080fd5b806349b7b802146101445780635ca1e165146101905780635d8105011461019857600080fd5b8063319cf735116100b2578063319cf7351461011e57806333d6247d146101275780633ed691ef1461013c57600080fd5b806301fd9044146100d9578063257b3632146100f55780632dfdf0b514610115575b600080fd5b6100e260005481565b6040519081526020015b60405180910390f35b6100e2610103366004610722565b60026020526000908152604090205481565b6100e260235481565b6100e260015481565b61013a610135366004610722565b61025d565b005b6100e2610406565b61016b7f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff90911681526020016100ec565b6100e261041b565b6100e26101a636600461073b565b604080516020808201959095528082019390935260c09190911b7fffffffffffffffff0000000000000000000000000000000000000000000000001660608301528051604881840301815260689092019052805191012090565b6100e261020e3660046107ac565b610425565b61016b7f000000000000000000000000000000000000000000000000000000000000000081565b61024d6102483660046107eb565b6104fb565b60405190151581526020016100ec565b60008073ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001633036102ad57505060018190556000548161032d565b73ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001633036102fb5750506000819055600154819061032d565b6040517fb49365dd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006103398284610513565b6000818152600260205260408120549192500361040057600061035d600143610862565b60008381526002602090815260409182902092409283905581518082018690528083018490527fffffffffffffffff0000000000000000000000000000000000000000000000004260c01b16606082015282518082036048018152606890910190925281519101209091506103d190610542565b604051849084907fda61aa7823fcd807e37b95aabcbe17f03a6f3efd514176444dae191d27fd66b390600090a3505b50505050565b6000610416600154600054610513565b905090565b6000610416610645565b600083815b60208110156104f257600163ffffffff8516821c811690036104955784816020811061045857610458610875565b602002013582604051602001610478929190918252602082015260400190565b6040516020818303038152906040528051906020012091506104e0565b818582602081106104a8576104a8610875565b60200201356040516020016104c7929190918252602082015260400190565b6040516020818303038152906040528051906020012091505b806104ea816108a4565b91505061042a565b50949350505050565b600081610509868686610425565b1495945050505050565b604080516020808201859052818301849052825180830384018152606090920190925280519101205b92915050565b806001610551602060026109fc565b61055b9190610862565b60235410610595576040517fef5ccf6600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60006023600081546105a6906108a4565b9182905550905060005b6020811015610637578082901c6001166001036105e35782600382602081106105db576105db610875565b015550505050565b600381602081106105f6576105f6610875565b01546040805160208101929092528101849052606001604051602081830303815290604052805190602001209250808061062f906108a4565b9150506105b0565b50610640610a0f565b505050565b602354600090819081805b6020811015610719578083901c6001166001036106ad576003816020811061067a5761067a610875565b015460408051602081019290925281018590526060016040516020818303038152906040528051906020012093506106da565b60408051602081018690529081018390526060016040516020818303038152906040528051906020012093505b60408051602081018490529081018390526060016040516020818303038152906040528051906020012091508080610711906108a4565b915050610650565b50919392505050565b60006020828403121561073457600080fd5b5035919050565b60008060006060848603121561075057600080fd5b8335925060208401359150604084013567ffffffffffffffff8116811461077657600080fd5b809150509250925092565b80610400810183101561053c57600080fd5b803563ffffffff811681146107a757600080fd5b919050565b600080600061044084860312156107c257600080fd5b833592506107d38560208601610781565b91506107e26104208501610793565b90509250925092565b600080600080610460858703121561080257600080fd5b843593506108138660208701610781565b92506108226104208601610793565b939692955092936104400135925050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b8181038181111561053c5761053c610833565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036108d5576108d5610833565b5060010190565b600181815b8085111561093557817fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0482111561091b5761091b610833565b8085161561092857918102915b93841c93908002906108e1565b509250929050565b60008261094c5750600161053c565b816109595750600061053c565b816001811461096f576002811461097957610995565b600191505061053c565b60ff84111561098a5761098a610833565b50506001821b61053c565b5060208310610133831016604e8410600b84101617156109b8575081810a61053c565b6109c283836108dc565b807fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff048211156109f4576109f4610833565b029392505050565b6000610a08838361093d565b9392505050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052600160045260246000fdfea2646970667358221220fc07ebcb1bf3607eb76c734998833eef05f4a3c59de6fc9a8c736d9a5464407464736f6c63430008140033",
}

// EtrogpolygonzkevmglobalexitrootABI is the input ABI used to generate the binding from.
// Deprecated: Use EtrogpolygonzkevmglobalexitrootMetaData.ABI instead.
var EtrogpolygonzkevmglobalexitrootABI = EtrogpolygonzkevmglobalexitrootMetaData.ABI

// EtrogpolygonzkevmglobalexitrootBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use EtrogpolygonzkevmglobalexitrootMetaData.Bin instead.
var EtrogpolygonzkevmglobalexitrootBin = EtrogpolygonzkevmglobalexitrootMetaData.Bin

// DeployEtrogpolygonzkevmglobalexitroot deploys a new Ethereum contract, binding an instance of Etrogpolygonzkevmglobalexitroot to it.
func DeployEtrogpolygonzkevmglobalexitroot(auth *bind.TransactOpts, backend bind.ContractBackend, _rollupManager common.Address, _bridgeAddress common.Address) (common.Address, *types.Transaction, *Etrogpolygonzkevmglobalexitroot, error) {
	parsed, err := EtrogpolygonzkevmglobalexitrootMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(EtrogpolygonzkevmglobalexitrootBin), backend, _rollupManager, _bridgeAddress)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Etrogpolygonzkevmglobalexitroot{EtrogpolygonzkevmglobalexitrootCaller: EtrogpolygonzkevmglobalexitrootCaller{contract: contract}, EtrogpolygonzkevmglobalexitrootTransactor: EtrogpolygonzkevmglobalexitrootTransactor{contract: contract}, EtrogpolygonzkevmglobalexitrootFilterer: EtrogpolygonzkevmglobalexitrootFilterer{contract: contract}}, nil
}

// Etrogpolygonzkevmglobalexitroot is an auto generated Go binding around an Ethereum contract.
type Etrogpolygonzkevmglobalexitroot struct {
	EtrogpolygonzkevmglobalexitrootCaller     // Read-only binding to the contract
	EtrogpolygonzkevmglobalexitrootTransactor // Write-only binding to the contract
	EtrogpolygonzkevmglobalexitrootFilterer   // Log filterer for contract events
}

// EtrogpolygonzkevmglobalexitrootCaller is an auto generated read-only Go binding around an Ethereum contract.
type EtrogpolygonzkevmglobalexitrootCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EtrogpolygonzkevmglobalexitrootTransactor is an auto generated write-only Go binding around an Ethereum contract.
type EtrogpolygonzkevmglobalexitrootTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EtrogpolygonzkevmglobalexitrootFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type EtrogpolygonzkevmglobalexitrootFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EtrogpolygonzkevmglobalexitrootSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EtrogpolygonzkevmglobalexitrootSession struct {
	Contract     *Etrogpolygonzkevmglobalexitroot // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                    // Call options to use throughout this session
	TransactOpts bind.TransactOpts                // Transaction auth options to use throughout this session
}

// EtrogpolygonzkevmglobalexitrootCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EtrogpolygonzkevmglobalexitrootCallerSession struct {
	Contract *EtrogpolygonzkevmglobalexitrootCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                          // Call options to use throughout this session
}

// EtrogpolygonzkevmglobalexitrootTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EtrogpolygonzkevmglobalexitrootTransactorSession struct {
	Contract     *EtrogpolygonzkevmglobalexitrootTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                          // Transaction auth options to use throughout this session
}

// EtrogpolygonzkevmglobalexitrootRaw is an auto generated low-level Go binding around an Ethereum contract.
type EtrogpolygonzkevmglobalexitrootRaw struct {
	Contract *Etrogpolygonzkevmglobalexitroot // Generic contract binding to access the raw methods on
}

// EtrogpolygonzkevmglobalexitrootCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type EtrogpolygonzkevmglobalexitrootCallerRaw struct {
	Contract *EtrogpolygonzkevmglobalexitrootCaller // Generic read-only contract binding to access the raw methods on
}

// EtrogpolygonzkevmglobalexitrootTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type EtrogpolygonzkevmglobalexitrootTransactorRaw struct {
	Contract *EtrogpolygonzkevmglobalexitrootTransactor // Generic write-only contract binding to access the raw methods on
}

// NewEtrogpolygonzkevmglobalexitroot creates a new instance of Etrogpolygonzkevmglobalexitroot, bound to a specific deployed contract.
func NewEtrogpolygonzkevmglobalexitroot(address common.Address, backend bind.ContractBackend) (*Etrogpolygonzkevmglobalexitroot, error) {
	contract, err := bindEtrogpolygonzkevmglobalexitroot(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Etrogpolygonzkevmglobalexitroot{EtrogpolygonzkevmglobalexitrootCaller: EtrogpolygonzkevmglobalexitrootCaller{contract: contract}, EtrogpolygonzkevmglobalexitrootTransactor: EtrogpolygonzkevmglobalexitrootTransactor{contract: contract}, EtrogpolygonzkevmglobalexitrootFilterer: EtrogpolygonzkevmglobalexitrootFilterer{contract: contract}}, nil
}

// NewEtrogpolygonzkevmglobalexitrootCaller creates a new read-only instance of Etrogpolygonzkevmglobalexitroot, bound to a specific deployed contract.
func NewEtrogpolygonzkevmglobalexitrootCaller(address common.Address, caller bind.ContractCaller) (*EtrogpolygonzkevmglobalexitrootCaller, error) {
	contract, err := bindEtrogpolygonzkevmglobalexitroot(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EtrogpolygonzkevmglobalexitrootCaller{contract: contract}, nil
}

// NewEtrogpolygonzkevmglobalexitrootTransactor creates a new write-only instance of Etrogpolygonzkevmglobalexitroot, bound to a specific deployed contract.
func NewEtrogpolygonzkevmglobalexitrootTransactor(address common.Address, transactor bind.ContractTransactor) (*EtrogpolygonzkevmglobalexitrootTransactor, error) {
	contract, err := bindEtrogpolygonzkevmglobalexitroot(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EtrogpolygonzkevmglobalexitrootTransactor{contract: contract}, nil
}

// NewEtrogpolygonzkevmglobalexitrootFilterer creates a new log filterer instance of Etrogpolygonzkevmglobalexitroot, bound to a specific deployed contract.
func NewEtrogpolygonzkevmglobalexitrootFilterer(address common.Address, filterer bind.ContractFilterer) (*EtrogpolygonzkevmglobalexitrootFilterer, error) {
	contract, err := bindEtrogpolygonzkevmglobalexitroot(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EtrogpolygonzkevmglobalexitrootFilterer{contract: contract}, nil
}

// bindEtrogpolygonzkevmglobalexitroot binds a generic wrapper to an already deployed contract.
func bindEtrogpolygonzkevmglobalexitroot(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := EtrogpolygonzkevmglobalexitrootMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Etrogpolygonzkevmglobalexitroot *EtrogpolygonzkevmglobalexitrootRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Etrogpolygonzkevmglobalexitroot.Contract.EtrogpolygonzkevmglobalexitrootCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Etrogpolygonzkevmglobalexitroot *EtrogpolygonzkevmglobalexitrootRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Etrogpolygonzkevmglobalexitroot.Contract.EtrogpolygonzkevmglobalexitrootTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Etrogpolygonzkevmglobalexitroot *EtrogpolygonzkevmglobalexitrootRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Etrogpolygonzkevmglobalexitroot.Contract.EtrogpolygonzkevmglobalexitrootTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Etrogpolygonzkevmglobalexitroot *EtrogpolygonzkevmglobalexitrootCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Etrogpolygonzkevmglobalexitroot.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Etrogpolygonzkevmglobalexitroot *EtrogpolygonzkevmglobalexitrootTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Etrogpolygonzkevmglobalexitroot.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Etrogpolygonzkevmglobalexitroot *EtrogpolygonzkevmglobalexitrootTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Etrogpolygonzkevmglobalexitroot.Contract.contract.Transact(opts, method, params...)
}

// BridgeAddress is a free data retrieval call binding the contract method 0xa3c573eb.
//
// Solidity: function bridgeAddress() view returns(address)
func (_Etrogpolygonzkevmglobalexitroot *EtrogpolygonzkevmglobalexitrootCaller) BridgeAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Etrogpolygonzkevmglobalexitroot.contract.Call(opts, &out, "bridgeAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BridgeAddress is a free data retrieval call binding the contract method 0xa3c573eb.
//
// Solidity: function bridgeAddress() view returns(address)
func (_Etrogpolygonzkevmglobalexitroot *EtrogpolygonzkevmglobalexitrootSession) BridgeAddress() (common.Address, error) {
	return _Etrogpolygonzkevmglobalexitroot.Contract.BridgeAddress(&_Etrogpolygonzkevmglobalexitroot.CallOpts)
}

// BridgeAddress is a free data retrieval call binding the contract method 0xa3c573eb.
//
// Solidity: function bridgeAddress() view returns(address)
func (_Etrogpolygonzkevmglobalexitroot *EtrogpolygonzkevmglobalexitrootCallerSession) BridgeAddress() (common.Address, error) {
	return _Etrogpolygonzkevmglobalexitroot.Contract.BridgeAddress(&_Etrogpolygonzkevmglobalexitroot.CallOpts)
}

// CalculateRoot is a free data retrieval call binding the contract method 0x83f24403.
//
// Solidity: function calculateRoot(bytes32 leafHash, bytes32[32] smtProof, uint32 index) pure returns(bytes32)
func (_Etrogpolygonzkevmglobalexitroot *EtrogpolygonzkevmglobalexitrootCaller) CalculateRoot(opts *bind.CallOpts, leafHash [32]byte, smtProof [32][32]byte, index uint32) ([32]byte, error) {
	var out []interface{}
	err := _Etrogpolygonzkevmglobalexitroot.contract.Call(opts, &out, "calculateRoot", leafHash, smtProof, index)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// CalculateRoot is a free data retrieval call binding the contract method 0x83f24403.
//
// Solidity: function calculateRoot(bytes32 leafHash, bytes32[32] smtProof, uint32 index) pure returns(bytes32)
func (_Etrogpolygonzkevmglobalexitroot *EtrogpolygonzkevmglobalexitrootSession) CalculateRoot(leafHash [32]byte, smtProof [32][32]byte, index uint32) ([32]byte, error) {
	return _Etrogpolygonzkevmglobalexitroot.Contract.CalculateRoot(&_Etrogpolygonzkevmglobalexitroot.CallOpts, leafHash, smtProof, index)
}

// CalculateRoot is a free data retrieval call binding the contract method 0x83f24403.
//
// Solidity: function calculateRoot(bytes32 leafHash, bytes32[32] smtProof, uint32 index) pure returns(bytes32)
func (_Etrogpolygonzkevmglobalexitroot *EtrogpolygonzkevmglobalexitrootCallerSession) CalculateRoot(leafHash [32]byte, smtProof [32][32]byte, index uint32) ([32]byte, error) {
	return _Etrogpolygonzkevmglobalexitroot.Contract.CalculateRoot(&_Etrogpolygonzkevmglobalexitroot.CallOpts, leafHash, smtProof, index)
}

// DepositCount is a free data retrieval call binding the contract method 0x2dfdf0b5.
//
// Solidity: function depositCount() view returns(uint256)
func (_Etrogpolygonzkevmglobalexitroot *EtrogpolygonzkevmglobalexitrootCaller) DepositCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Etrogpolygonzkevmglobalexitroot.contract.Call(opts, &out, "depositCount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DepositCount is a free data retrieval call binding the contract method 0x2dfdf0b5.
//
// Solidity: function depositCount() view returns(uint256)
func (_Etrogpolygonzkevmglobalexitroot *EtrogpolygonzkevmglobalexitrootSession) DepositCount() (*big.Int, error) {
	return _Etrogpolygonzkevmglobalexitroot.Contract.DepositCount(&_Etrogpolygonzkevmglobalexitroot.CallOpts)
}

// DepositCount is a free data retrieval call binding the contract method 0x2dfdf0b5.
//
// Solidity: function depositCount() view returns(uint256)
func (_Etrogpolygonzkevmglobalexitroot *EtrogpolygonzkevmglobalexitrootCallerSession) DepositCount() (*big.Int, error) {
	return _Etrogpolygonzkevmglobalexitroot.Contract.DepositCount(&_Etrogpolygonzkevmglobalexitroot.CallOpts)
}

// GetLastGlobalExitRoot is a free data retrieval call binding the contract method 0x3ed691ef.
//
// Solidity: function getLastGlobalExitRoot() view returns(bytes32)
func (_Etrogpolygonzkevmglobalexitroot *EtrogpolygonzkevmglobalexitrootCaller) GetLastGlobalExitRoot(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Etrogpolygonzkevmglobalexitroot.contract.Call(opts, &out, "getLastGlobalExitRoot")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetLastGlobalExitRoot is a free data retrieval call binding the contract method 0x3ed691ef.
//
// Solidity: function getLastGlobalExitRoot() view returns(bytes32)
func (_Etrogpolygonzkevmglobalexitroot *EtrogpolygonzkevmglobalexitrootSession) GetLastGlobalExitRoot() ([32]byte, error) {
	return _Etrogpolygonzkevmglobalexitroot.Contract.GetLastGlobalExitRoot(&_Etrogpolygonzkevmglobalexitroot.CallOpts)
}

// GetLastGlobalExitRoot is a free data retrieval call binding the contract method 0x3ed691ef.
//
// Solidity: function getLastGlobalExitRoot() view returns(bytes32)
func (_Etrogpolygonzkevmglobalexitroot *EtrogpolygonzkevmglobalexitrootCallerSession) GetLastGlobalExitRoot() ([32]byte, error) {
	return _Etrogpolygonzkevmglobalexitroot.Contract.GetLastGlobalExitRoot(&_Etrogpolygonzkevmglobalexitroot.CallOpts)
}

// GetLeafValue is a free data retrieval call binding the contract method 0x5d810501.
//
// Solidity: function getLeafValue(bytes32 newGlobalExitRoot, uint256 lastBlockHash, uint64 timestamp) pure returns(bytes32)
func (_Etrogpolygonzkevmglobalexitroot *EtrogpolygonzkevmglobalexitrootCaller) GetLeafValue(opts *bind.CallOpts, newGlobalExitRoot [32]byte, lastBlockHash *big.Int, timestamp uint64) ([32]byte, error) {
	var out []interface{}
	err := _Etrogpolygonzkevmglobalexitroot.contract.Call(opts, &out, "getLeafValue", newGlobalExitRoot, lastBlockHash, timestamp)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetLeafValue is a free data retrieval call binding the contract method 0x5d810501.
//
// Solidity: function getLeafValue(bytes32 newGlobalExitRoot, uint256 lastBlockHash, uint64 timestamp) pure returns(bytes32)
func (_Etrogpolygonzkevmglobalexitroot *EtrogpolygonzkevmglobalexitrootSession) GetLeafValue(newGlobalExitRoot [32]byte, lastBlockHash *big.Int, timestamp uint64) ([32]byte, error) {
	return _Etrogpolygonzkevmglobalexitroot.Contract.GetLeafValue(&_Etrogpolygonzkevmglobalexitroot.CallOpts, newGlobalExitRoot, lastBlockHash, timestamp)
}

// GetLeafValue is a free data retrieval call binding the contract method 0x5d810501.
//
// Solidity: function getLeafValue(bytes32 newGlobalExitRoot, uint256 lastBlockHash, uint64 timestamp) pure returns(bytes32)
func (_Etrogpolygonzkevmglobalexitroot *EtrogpolygonzkevmglobalexitrootCallerSession) GetLeafValue(newGlobalExitRoot [32]byte, lastBlockHash *big.Int, timestamp uint64) ([32]byte, error) {
	return _Etrogpolygonzkevmglobalexitroot.Contract.GetLeafValue(&_Etrogpolygonzkevmglobalexitroot.CallOpts, newGlobalExitRoot, lastBlockHash, timestamp)
}

// GetRoot is a free data retrieval call binding the contract method 0x5ca1e165.
//
// Solidity: function getRoot() view returns(bytes32)
func (_Etrogpolygonzkevmglobalexitroot *EtrogpolygonzkevmglobalexitrootCaller) GetRoot(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Etrogpolygonzkevmglobalexitroot.contract.Call(opts, &out, "getRoot")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoot is a free data retrieval call binding the contract method 0x5ca1e165.
//
// Solidity: function getRoot() view returns(bytes32)
func (_Etrogpolygonzkevmglobalexitroot *EtrogpolygonzkevmglobalexitrootSession) GetRoot() ([32]byte, error) {
	return _Etrogpolygonzkevmglobalexitroot.Contract.GetRoot(&_Etrogpolygonzkevmglobalexitroot.CallOpts)
}

// GetRoot is a free data retrieval call binding the contract method 0x5ca1e165.
//
// Solidity: function getRoot() view returns(bytes32)
func (_Etrogpolygonzkevmglobalexitroot *EtrogpolygonzkevmglobalexitrootCallerSession) GetRoot() ([32]byte, error) {
	return _Etrogpolygonzkevmglobalexitroot.Contract.GetRoot(&_Etrogpolygonzkevmglobalexitroot.CallOpts)
}

// GlobalExitRootMap is a free data retrieval call binding the contract method 0x257b3632.
//
// Solidity: function globalExitRootMap(bytes32 ) view returns(uint256)
func (_Etrogpolygonzkevmglobalexitroot *EtrogpolygonzkevmglobalexitrootCaller) GlobalExitRootMap(opts *bind.CallOpts, arg0 [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _Etrogpolygonzkevmglobalexitroot.contract.Call(opts, &out, "globalExitRootMap", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GlobalExitRootMap is a free data retrieval call binding the contract method 0x257b3632.
//
// Solidity: function globalExitRootMap(bytes32 ) view returns(uint256)
func (_Etrogpolygonzkevmglobalexitroot *EtrogpolygonzkevmglobalexitrootSession) GlobalExitRootMap(arg0 [32]byte) (*big.Int, error) {
	return _Etrogpolygonzkevmglobalexitroot.Contract.GlobalExitRootMap(&_Etrogpolygonzkevmglobalexitroot.CallOpts, arg0)
}

// GlobalExitRootMap is a free data retrieval call binding the contract method 0x257b3632.
//
// Solidity: function globalExitRootMap(bytes32 ) view returns(uint256)
func (_Etrogpolygonzkevmglobalexitroot *EtrogpolygonzkevmglobalexitrootCallerSession) GlobalExitRootMap(arg0 [32]byte) (*big.Int, error) {
	return _Etrogpolygonzkevmglobalexitroot.Contract.GlobalExitRootMap(&_Etrogpolygonzkevmglobalexitroot.CallOpts, arg0)
}

// LastMainnetExitRoot is a free data retrieval call binding the contract method 0x319cf735.
//
// Solidity: function lastMainnetExitRoot() view returns(bytes32)
func (_Etrogpolygonzkevmglobalexitroot *EtrogpolygonzkevmglobalexitrootCaller) LastMainnetExitRoot(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Etrogpolygonzkevmglobalexitroot.contract.Call(opts, &out, "lastMainnetExitRoot")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// LastMainnetExitRoot is a free data retrieval call binding the contract method 0x319cf735.
//
// Solidity: function lastMainnetExitRoot() view returns(bytes32)
func (_Etrogpolygonzkevmglobalexitroot *EtrogpolygonzkevmglobalexitrootSession) LastMainnetExitRoot() ([32]byte, error) {
	return _Etrogpolygonzkevmglobalexitroot.Contract.LastMainnetExitRoot(&_Etrogpolygonzkevmglobalexitroot.CallOpts)
}

// LastMainnetExitRoot is a free data retrieval call binding the contract method 0x319cf735.
//
// Solidity: function lastMainnetExitRoot() view returns(bytes32)
func (_Etrogpolygonzkevmglobalexitroot *EtrogpolygonzkevmglobalexitrootCallerSession) LastMainnetExitRoot() ([32]byte, error) {
	return _Etrogpolygonzkevmglobalexitroot.Contract.LastMainnetExitRoot(&_Etrogpolygonzkevmglobalexitroot.CallOpts)
}

// LastRollupExitRoot is a free data retrieval call binding the contract method 0x01fd9044.
//
// Solidity: function lastRollupExitRoot() view returns(bytes32)
func (_Etrogpolygonzkevmglobalexitroot *EtrogpolygonzkevmglobalexitrootCaller) LastRollupExitRoot(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Etrogpolygonzkevmglobalexitroot.contract.Call(opts, &out, "lastRollupExitRoot")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// LastRollupExitRoot is a free data retrieval call binding the contract method 0x01fd9044.
//
// Solidity: function lastRollupExitRoot() view returns(bytes32)
func (_Etrogpolygonzkevmglobalexitroot *EtrogpolygonzkevmglobalexitrootSession) LastRollupExitRoot() ([32]byte, error) {
	return _Etrogpolygonzkevmglobalexitroot.Contract.LastRollupExitRoot(&_Etrogpolygonzkevmglobalexitroot.CallOpts)
}

// LastRollupExitRoot is a free data retrieval call binding the contract method 0x01fd9044.
//
// Solidity: function lastRollupExitRoot() view returns(bytes32)
func (_Etrogpolygonzkevmglobalexitroot *EtrogpolygonzkevmglobalexitrootCallerSession) LastRollupExitRoot() ([32]byte, error) {
	return _Etrogpolygonzkevmglobalexitroot.Contract.LastRollupExitRoot(&_Etrogpolygonzkevmglobalexitroot.CallOpts)
}

// RollupManager is a free data retrieval call binding the contract method 0x49b7b802.
//
// Solidity: function rollupManager() view returns(address)
func (_Etrogpolygonzkevmglobalexitroot *EtrogpolygonzkevmglobalexitrootCaller) RollupManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Etrogpolygonzkevmglobalexitroot.contract.Call(opts, &out, "rollupManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RollupManager is a free data retrieval call binding the contract method 0x49b7b802.
//
// Solidity: function rollupManager() view returns(address)
func (_Etrogpolygonzkevmglobalexitroot *EtrogpolygonzkevmglobalexitrootSession) RollupManager() (common.Address, error) {
	return _Etrogpolygonzkevmglobalexitroot.Contract.RollupManager(&_Etrogpolygonzkevmglobalexitroot.CallOpts)
}

// RollupManager is a free data retrieval call binding the contract method 0x49b7b802.
//
// Solidity: function rollupManager() view returns(address)
func (_Etrogpolygonzkevmglobalexitroot *EtrogpolygonzkevmglobalexitrootCallerSession) RollupManager() (common.Address, error) {
	return _Etrogpolygonzkevmglobalexitroot.Contract.RollupManager(&_Etrogpolygonzkevmglobalexitroot.CallOpts)
}

// VerifyMerkleProof is a free data retrieval call binding the contract method 0xfb570834.
//
// Solidity: function verifyMerkleProof(bytes32 leafHash, bytes32[32] smtProof, uint32 index, bytes32 root) pure returns(bool)
func (_Etrogpolygonzkevmglobalexitroot *EtrogpolygonzkevmglobalexitrootCaller) VerifyMerkleProof(opts *bind.CallOpts, leafHash [32]byte, smtProof [32][32]byte, index uint32, root [32]byte) (bool, error) {
	var out []interface{}
	err := _Etrogpolygonzkevmglobalexitroot.contract.Call(opts, &out, "verifyMerkleProof", leafHash, smtProof, index, root)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// VerifyMerkleProof is a free data retrieval call binding the contract method 0xfb570834.
//
// Solidity: function verifyMerkleProof(bytes32 leafHash, bytes32[32] smtProof, uint32 index, bytes32 root) pure returns(bool)
func (_Etrogpolygonzkevmglobalexitroot *EtrogpolygonzkevmglobalexitrootSession) VerifyMerkleProof(leafHash [32]byte, smtProof [32][32]byte, index uint32, root [32]byte) (bool, error) {
	return _Etrogpolygonzkevmglobalexitroot.Contract.VerifyMerkleProof(&_Etrogpolygonzkevmglobalexitroot.CallOpts, leafHash, smtProof, index, root)
}

// VerifyMerkleProof is a free data retrieval call binding the contract method 0xfb570834.
//
// Solidity: function verifyMerkleProof(bytes32 leafHash, bytes32[32] smtProof, uint32 index, bytes32 root) pure returns(bool)
func (_Etrogpolygonzkevmglobalexitroot *EtrogpolygonzkevmglobalexitrootCallerSession) VerifyMerkleProof(leafHash [32]byte, smtProof [32][32]byte, index uint32, root [32]byte) (bool, error) {
	return _Etrogpolygonzkevmglobalexitroot.Contract.VerifyMerkleProof(&_Etrogpolygonzkevmglobalexitroot.CallOpts, leafHash, smtProof, index, root)
}

// UpdateExitRoot is a paid mutator transaction binding the contract method 0x33d6247d.
//
// Solidity: function updateExitRoot(bytes32 newRoot) returns()
func (_Etrogpolygonzkevmglobalexitroot *EtrogpolygonzkevmglobalexitrootTransactor) UpdateExitRoot(opts *bind.TransactOpts, newRoot [32]byte) (*types.Transaction, error) {
	return _Etrogpolygonzkevmglobalexitroot.contract.Transact(opts, "updateExitRoot", newRoot)
}

// UpdateExitRoot is a paid mutator transaction binding the contract method 0x33d6247d.
//
// Solidity: function updateExitRoot(bytes32 newRoot) returns()
func (_Etrogpolygonzkevmglobalexitroot *EtrogpolygonzkevmglobalexitrootSession) UpdateExitRoot(newRoot [32]byte) (*types.Transaction, error) {
	return _Etrogpolygonzkevmglobalexitroot.Contract.UpdateExitRoot(&_Etrogpolygonzkevmglobalexitroot.TransactOpts, newRoot)
}

// UpdateExitRoot is a paid mutator transaction binding the contract method 0x33d6247d.
//
// Solidity: function updateExitRoot(bytes32 newRoot) returns()
func (_Etrogpolygonzkevmglobalexitroot *EtrogpolygonzkevmglobalexitrootTransactorSession) UpdateExitRoot(newRoot [32]byte) (*types.Transaction, error) {
	return _Etrogpolygonzkevmglobalexitroot.Contract.UpdateExitRoot(&_Etrogpolygonzkevmglobalexitroot.TransactOpts, newRoot)
}

// EtrogpolygonzkevmglobalexitrootUpdateL1InfoTreeIterator is returned from FilterUpdateL1InfoTree and is used to iterate over the raw logs and unpacked data for UpdateL1InfoTree events raised by the Etrogpolygonzkevmglobalexitroot contract.
type EtrogpolygonzkevmglobalexitrootUpdateL1InfoTreeIterator struct {
	Event *EtrogpolygonzkevmglobalexitrootUpdateL1InfoTree // Event containing the contract specifics and raw log

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
func (it *EtrogpolygonzkevmglobalexitrootUpdateL1InfoTreeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EtrogpolygonzkevmglobalexitrootUpdateL1InfoTree)
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
		it.Event = new(EtrogpolygonzkevmglobalexitrootUpdateL1InfoTree)
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
func (it *EtrogpolygonzkevmglobalexitrootUpdateL1InfoTreeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EtrogpolygonzkevmglobalexitrootUpdateL1InfoTreeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EtrogpolygonzkevmglobalexitrootUpdateL1InfoTree represents a UpdateL1InfoTree event raised by the Etrogpolygonzkevmglobalexitroot contract.
type EtrogpolygonzkevmglobalexitrootUpdateL1InfoTree struct {
	MainnetExitRoot [32]byte
	RollupExitRoot  [32]byte
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterUpdateL1InfoTree is a free log retrieval operation binding the contract event 0xda61aa7823fcd807e37b95aabcbe17f03a6f3efd514176444dae191d27fd66b3.
//
// Solidity: event UpdateL1InfoTree(bytes32 indexed mainnetExitRoot, bytes32 indexed rollupExitRoot)
func (_Etrogpolygonzkevmglobalexitroot *EtrogpolygonzkevmglobalexitrootFilterer) FilterUpdateL1InfoTree(opts *bind.FilterOpts, mainnetExitRoot [][32]byte, rollupExitRoot [][32]byte) (*EtrogpolygonzkevmglobalexitrootUpdateL1InfoTreeIterator, error) {

	var mainnetExitRootRule []interface{}
	for _, mainnetExitRootItem := range mainnetExitRoot {
		mainnetExitRootRule = append(mainnetExitRootRule, mainnetExitRootItem)
	}
	var rollupExitRootRule []interface{}
	for _, rollupExitRootItem := range rollupExitRoot {
		rollupExitRootRule = append(rollupExitRootRule, rollupExitRootItem)
	}

	logs, sub, err := _Etrogpolygonzkevmglobalexitroot.contract.FilterLogs(opts, "UpdateL1InfoTree", mainnetExitRootRule, rollupExitRootRule)
	if err != nil {
		return nil, err
	}
	return &EtrogpolygonzkevmglobalexitrootUpdateL1InfoTreeIterator{contract: _Etrogpolygonzkevmglobalexitroot.contract, event: "UpdateL1InfoTree", logs: logs, sub: sub}, nil
}

// WatchUpdateL1InfoTree is a free log subscription operation binding the contract event 0xda61aa7823fcd807e37b95aabcbe17f03a6f3efd514176444dae191d27fd66b3.
//
// Solidity: event UpdateL1InfoTree(bytes32 indexed mainnetExitRoot, bytes32 indexed rollupExitRoot)
func (_Etrogpolygonzkevmglobalexitroot *EtrogpolygonzkevmglobalexitrootFilterer) WatchUpdateL1InfoTree(opts *bind.WatchOpts, sink chan<- *EtrogpolygonzkevmglobalexitrootUpdateL1InfoTree, mainnetExitRoot [][32]byte, rollupExitRoot [][32]byte) (event.Subscription, error) {

	var mainnetExitRootRule []interface{}
	for _, mainnetExitRootItem := range mainnetExitRoot {
		mainnetExitRootRule = append(mainnetExitRootRule, mainnetExitRootItem)
	}
	var rollupExitRootRule []interface{}
	for _, rollupExitRootItem := range rollupExitRoot {
		rollupExitRootRule = append(rollupExitRootRule, rollupExitRootItem)
	}

	logs, sub, err := _Etrogpolygonzkevmglobalexitroot.contract.WatchLogs(opts, "UpdateL1InfoTree", mainnetExitRootRule, rollupExitRootRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EtrogpolygonzkevmglobalexitrootUpdateL1InfoTree)
				if err := _Etrogpolygonzkevmglobalexitroot.contract.UnpackLog(event, "UpdateL1InfoTree", log); err != nil {
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

// ParseUpdateL1InfoTree is a log parse operation binding the contract event 0xda61aa7823fcd807e37b95aabcbe17f03a6f3efd514176444dae191d27fd66b3.
//
// Solidity: event UpdateL1InfoTree(bytes32 indexed mainnetExitRoot, bytes32 indexed rollupExitRoot)
func (_Etrogpolygonzkevmglobalexitroot *EtrogpolygonzkevmglobalexitrootFilterer) ParseUpdateL1InfoTree(log types.Log) (*EtrogpolygonzkevmglobalexitrootUpdateL1InfoTree, error) {
	event := new(EtrogpolygonzkevmglobalexitrootUpdateL1InfoTree)
	if err := _Etrogpolygonzkevmglobalexitroot.contract.UnpackLog(event, "UpdateL1InfoTree", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
