// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package eigendaverifier

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

// BN254G1Point is an auto generated low-level Go binding around an user-defined struct.
type BN254G1Point struct {
	X *big.Int
	Y *big.Int
}

// EigenDARollupUtilsBlobVerificationProof is an auto generated low-level Go binding around an user-defined struct.
type EigenDARollupUtilsBlobVerificationProof struct {
	BatchId        uint32
	BlobIndex      uint32
	BatchMetadata  IEigenDAServiceManagerBatchMetadata
	InclusionProof []byte
	QuorumIndices  []byte
}

// EigenDAVerifierBlobData is an auto generated low-level Go binding around an user-defined struct.
type EigenDAVerifierBlobData struct {
	BlobHeader            IEigenDAServiceManagerBlobHeader
	BlobVerificationProof EigenDARollupUtilsBlobVerificationProof
}

// IEigenDAServiceManagerBatchHeader is an auto generated low-level Go binding around an user-defined struct.
type IEigenDAServiceManagerBatchHeader struct {
	BlobHeadersRoot       [32]byte
	QuorumNumbers         []byte
	SignedStakeForQuorums []byte
	ReferenceBlockNumber  uint32
}

// IEigenDAServiceManagerBatchMetadata is an auto generated low-level Go binding around an user-defined struct.
type IEigenDAServiceManagerBatchMetadata struct {
	BatchHeader             IEigenDAServiceManagerBatchHeader
	SignatoryRecordHash     [32]byte
	ConfirmationBlockNumber uint32
}

// IEigenDAServiceManagerBlobHeader is an auto generated low-level Go binding around an user-defined struct.
type IEigenDAServiceManagerBlobHeader struct {
	Commitment       BN254G1Point
	DataLength       uint32
	QuorumBlobParams []IEigenDAServiceManagerQuorumBlobParam
}

// IEigenDAServiceManagerQuorumBlobParam is an auto generated low-level Go binding around an user-defined struct.
type IEigenDAServiceManagerQuorumBlobParam struct {
	QuorumNumber                    uint8
	AdversaryThresholdPercentage    uint8
	ConfirmationThresholdPercentage uint8
	ChunkLength                     uint32
}

// EigendaverifierMetaData contains all meta data concerning the Eigendaverifier contract.
var EigendaverifierMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_admin\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_eigenDAServiceManager\",\"type\":\"address\",\"internalType\":\"contractIEigenDAServiceManager\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"acceptAdminRole\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"admin\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"decodeBlobData\",\"inputs\":[{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"blobData\",\"type\":\"tuple\",\"internalType\":\"structEigenDAVerifier.BlobData\",\"components\":[{\"name\":\"blobHeader\",\"type\":\"tuple\",\"internalType\":\"structIEigenDAServiceManager.BlobHeader\",\"components\":[{\"name\":\"commitment\",\"type\":\"tuple\",\"internalType\":\"structBN254.G1Point\",\"components\":[{\"name\":\"X\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"Y\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"name\":\"dataLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"quorumBlobParams\",\"type\":\"tuple[]\",\"internalType\":\"structIEigenDAServiceManager.QuorumBlobParam[]\",\"components\":[{\"name\":\"quorumNumber\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"adversaryThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"confirmationThresholdPercentage\",\"type\":\"uint8\",\"internalType\":\"uint8\"},{\"name\":\"chunkLength\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]}]},{\"name\":\"blobVerificationProof\",\"type\":\"tuple\",\"internalType\":\"structEigenDARollupUtils.BlobVerificationProof\",\"components\":[{\"name\":\"batchId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"blobIndex\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"batchMetadata\",\"type\":\"tuple\",\"internalType\":\"structIEigenDAServiceManager.BatchMetadata\",\"components\":[{\"name\":\"batchHeader\",\"type\":\"tuple\",\"internalType\":\"structIEigenDAServiceManager.BatchHeader\",\"components\":[{\"name\":\"blobHeadersRoot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"quorumNumbers\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"signedStakeForQuorums\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"referenceBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"signatoryRecordHash\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"confirmationBlockNumber\",\"type\":\"uint32\",\"internalType\":\"uint32\"}]},{\"name\":\"inclusionProof\",\"type\":\"bytes\",\"internalType\":\"bytes\"},{\"name\":\"quorumIndices\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}]}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"getDataAvailabilityProtocol\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getProcotolName\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"pendingAdmin\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setDataAvailabilityProtocol\",\"inputs\":[{\"name\":\"newDataAvailabilityProtocol\",\"type\":\"address\",\"internalType\":\"contractIEigenDAServiceManager\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferAdminRole\",\"inputs\":[{\"name\":\"newPendingAdmin\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"verifyMessage\",\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"AcceptAdminRole\",\"inputs\":[{\"name\":\"newAdmin\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SetDataAvailabilityProtocol\",\"inputs\":[{\"name\":\"newTrustedSequencer\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"contractIEigenDAServiceManager\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"TransferAdminRole\",\"inputs\":[{\"name\":\"newPendingAdmin\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"BatchAlreadyVerified\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"BatchNotSequencedOrNotSequenceEnd\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ExceedMaxVerifyBatches\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FinalNumBatchBelowLastVerifiedBatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FinalNumBatchDoesNotMatchPendingState\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FinalPendingStateNumInvalid\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ForceBatchNotAllowed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ForceBatchTimeoutNotExpired\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ForceBatchesAlreadyActive\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ForceBatchesDecentralized\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ForceBatchesNotAllowedOnEmergencyState\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ForceBatchesOverflow\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ForcedDataDoesNotMatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"GasTokenNetworkMustBeZeroOnEther\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"GlobalExitRootNotExist\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"HaltTimeoutNotExpired\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"HaltTimeoutNotExpiredAfterEmergencyState\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"HugeTokenMetadataNotSupported\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InitNumBatchAboveLastVerifiedBatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InitNumBatchDoesNotMatchPendingState\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InitSequencedBatchDoesNotMatch\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidInitializeTransaction\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidProof\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidRangeBatchTimeTarget\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidRangeForceBatchTimeout\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidRangeMultiplierBatchFee\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MaxTimestampSequenceInvalid\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NewAccInputHashDoesNotExist\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NewPendingStateTimeoutMustBeLower\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NewStateRootNotInsidePrime\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NewTrustedAggregatorTimeoutMustBeLower\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotEnoughMaticAmount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotEnoughPOLAmount\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OldAccInputHashDoesNotExist\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OldStateRootDoesNotExist\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyAdmin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyPendingAdmin\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyRollupManager\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyTrustedAggregator\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OnlyTrustedSequencer\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PendingStateDoesNotExist\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PendingStateInvalid\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PendingStateNotConsolidable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"PendingStateTimeoutExceedHaltAggregationTimeout\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SequenceZeroBatches\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SequencedTimestampBelowForcedTimestamp\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SequencedTimestampInvalid\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"StoredRootMustBeDifferentThanNewRoot\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TransactionsLengthAboveMax\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TrustedAggregatorTimeoutExceedHaltAggregationTimeout\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TrustedAggregatorTimeoutNotExpired\",\"inputs\":[]}]",
	Bin: "0x608060405234801561001057600080fd5b50604051610ddc380380610ddc83398101604081905261002f91610078565b600180546001600160a01b039384166001600160a01b031991821617909155600080549290931691161790556100b2565b6001600160a01b038116811461007557600080fd5b50565b6000806040838503121561008b57600080fd5b825161009681610060565b60208401519092506100a781610060565b809150509250929050565b610d1b806100c16000396000f3fe608060405234801561001057600080fd5b50600436106100935760003560e01c8063aba4c80d11610066578063aba4c80d146100f8578063ada8f91914610118578063b9c67c331461012b578063e4f171201461013c578063f851a4401461016557600080fd5b806326782247146100985780633b51be4b146100c85780637cd76b8b146100dd5780638c3d7301146100f0575b600080fd5b6002546100ab906001600160a01b031681565b6040516001600160a01b0390911681526020015b60405180910390f35b6100db6100d6366004610484565b610178565b005b6100db6100eb3660046104e8565b61020a565b6100db61028a565b61010b610106366004610505565b61030b565b6040516100bf9190610660565b6100db6101263660046104e8565b610326565b6000546001600160a01b03166100ab565b6040805180820182526007815266456967656e444160c81b602082015290516100bf919061073d565b6001546100ab906001600160a01b031681565b6000610184838361030b565b8051600054602083015160405163219460e160e21b815293945073__$399f9ce8dd33a06d144ce1eb24d845e280$__936386518384936101d49390926001600160a01b0390911691600401610750565b60006040518083038186803b1580156101ec57600080fd5b505af4158015610200573d6000803e3d6000fd5b5050505050505050565b6001546001600160a01b0316331461023557604051634755657960e01b815260040160405180910390fd5b600080546001600160a01b0319166001600160a01b0383169081179091556040519081527fd331bd4c4cd1afecb94a225184bded161ff3213624ba4fb58c4f30c5a861144a906020015b60405180910390a150565b6002546001600160a01b031633146102b55760405163d1ec4b2360e01b815260040160405180910390fd5b600254600180546001600160a01b0319166001600160a01b0390921691821790556040519081527f056dc487bbf0795d0bbb1b4f0af523a855503cff740bfb4d5475f7a90c091e8e9060200160405180910390a1565b61031361039f565b61031f82840184610b27565b9392505050565b6001546001600160a01b0316331461035157604051634755657960e01b815260040160405180910390fd5b600280546001600160a01b0319166001600160a01b0383169081179091556040519081527fa5b56b7906fd0a20e3f35120dd8343db1e12e037a6c90111c7e42885e82a1ce69060200161027f565b6040805160e081018252600060a0820181815260c083018290529282019283526060808301919091526080820152908152602081016104366040805160a0808201835260008083526020808401829052845160e0810186526060808201848152608083018290529482015260c081018390529283528201819052818401529091820190815260200160608152602001606081525090565b905290565b60008083601f84011261044d57600080fd5b50813567ffffffffffffffff81111561046557600080fd5b60208301915083602082850101111561047d57600080fd5b9250929050565b60008060006040848603121561049957600080fd5b83359250602084013567ffffffffffffffff8111156104b757600080fd5b6104c38682870161043b565b9497909650939450505050565b6001600160a01b03811681146104e557600080fd5b50565b6000602082840312156104fa57600080fd5b813561031f816104d0565b6000806020838503121561051857600080fd5b823567ffffffffffffffff81111561052f57600080fd5b61053b8582860161043b565b90969095509350505050565b6000815180845260005b8181101561056d57602081850181015186830182015201610551565b506000602082860101526020601f19601f83011685010191505092915050565b600063ffffffff808351168452806020840151166020850152604083015160a060408601528051606060a08701528051610100870152602081015160806101208801526105de610180880182610547565b9050604082015160ff19888303016101408901526105fc8282610547565b91505083606083015116610160880152602083015160c08801528360408401511660e088015260608601519350868103606088015261063b8185610547565b9350505050608083015184820360808601526106578282610547565b95945050505050565b60006020808352835160408285015260e0840161068b60608601835180518252602090810151910152565b8183015163ffffffff1660a0860152604090910151608060c0860181905281519283905290830191600091906101008701905b808410156107155761070182865160ff815116825260ff602082015116602083015260ff604082015116604083015263ffffffff60608201511660608301525050565b9385019360019390930192908201906106be565b5093870151868503601f1901604088015293610731818661058d565b98975050505050505050565b60208152600061031f6020830184610547565b60608152600060e0820161077260608401875180518252602090810151910152565b60208681015163ffffffff1660a08501526040870151608060c0860181905281519384905290820192600091906101008701905b808410156107fd576107e982875160ff815116825260ff602082015116602083015260ff604082015116604083015263ffffffff60608201511660608301525050565b9484019460019390930192908201906107a6565b506001600160a01b03891687850152868103604088015261081e818961058d565b9a9950505050505050505050565b634e487b7160e01b600052604160045260246000fd5b6040516060810167ffffffffffffffff811182821017156108655761086561082c565b60405290565b6040516080810167ffffffffffffffff811182821017156108655761086561082c565b60405160a0810167ffffffffffffffff811182821017156108655761086561082c565b6040805190810167ffffffffffffffff811182821017156108655761086561082c565b604051601f8201601f1916810167ffffffffffffffff811182821017156108fd576108fd61082c565b604052919050565b803563ffffffff8116811461091957600080fd5b919050565b803560ff8116811461091957600080fd5b600082601f83011261094057600080fd5b813567ffffffffffffffff81111561095a5761095a61082c565b61096d601f8201601f19166020016108d4565b81815284602083860101111561098257600080fd5b816020850160208301376000918101602001919091529392505050565b6000606082840312156109b157600080fd5b6109b9610842565b9050813567ffffffffffffffff808211156109d357600080fd5b90830190608082860312156109e757600080fd5b6109ef61086b565b82358152602083013582811115610a0557600080fd5b610a118782860161092f565b602083015250604083013582811115610a2957600080fd5b610a358782860161092f565b604083015250610a4760608401610905565b60608201528352505060208281013590820152610a6660408301610905565b604082015292915050565b600060a08284031215610a8357600080fd5b610a8b61088e565b9050610a9682610905565b8152610aa460208301610905565b6020820152604082013567ffffffffffffffff80821115610ac457600080fd5b610ad08583860161099f565b60408401526060840135915080821115610ae957600080fd5b610af58583860161092f565b60608401526080840135915080821115610b0e57600080fd5b50610b1b8482850161092f565b60808301525092915050565b60006020808385031215610b3a57600080fd5b823567ffffffffffffffff80821115610b5257600080fd5b81850191506040808388031215610b6857600080fd5b610b706108b1565b833583811115610b7f57600080fd5b84018089036080811215610b9257600080fd5b610b9a610842565b84821215610ba757600080fd5b610baf6108b1565b9150823582528783013588830152818152610bcb858401610905565b88820152606091508183013586811115610be457600080fd5b8084019350508a601f840112610bf957600080fd5b823586811115610c0b57610c0b61082c565b610c19898260051b016108d4565b81815260079190911b8401890190898101908d831115610c3857600080fd5b948a01945b82861015610ca8576080868f031215610c565760008081fd5b610c5e61086b565b610c678761091e565b8152610c748c880161091e565b8c820152610c8389880161091e565b89820152610c92868801610905565b81870152825260809590950194908a0190610c3d565b96830196909652508352505083850135915082821115610cc757600080fd5b610cd388838601610a71565b8582015280955050505050509291505056fea264697066735822122074ae60b83f74a760fb2c661671d91cbe3b1757e03077b2559d96b421d2e9525164736f6c63430008140033",
}

// EigendaverifierABI is the input ABI used to generate the binding from.
// Deprecated: Use EigendaverifierMetaData.ABI instead.
var EigendaverifierABI = EigendaverifierMetaData.ABI

// EigendaverifierBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use EigendaverifierMetaData.Bin instead.
var EigendaverifierBin = EigendaverifierMetaData.Bin

// DeployEigendaverifier deploys a new Ethereum contract, binding an instance of Eigendaverifier to it.
func DeployEigendaverifier(auth *bind.TransactOpts, backend bind.ContractBackend, _admin common.Address, _eigenDAServiceManager common.Address) (common.Address, *types.Transaction, *Eigendaverifier, error) {
	parsed, err := EigendaverifierMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(EigendaverifierBin), backend, _admin, _eigenDAServiceManager)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Eigendaverifier{EigendaverifierCaller: EigendaverifierCaller{contract: contract}, EigendaverifierTransactor: EigendaverifierTransactor{contract: contract}, EigendaverifierFilterer: EigendaverifierFilterer{contract: contract}}, nil
}

// Eigendaverifier is an auto generated Go binding around an Ethereum contract.
type Eigendaverifier struct {
	EigendaverifierCaller     // Read-only binding to the contract
	EigendaverifierTransactor // Write-only binding to the contract
	EigendaverifierFilterer   // Log filterer for contract events
}

// EigendaverifierCaller is an auto generated read-only Go binding around an Ethereum contract.
type EigendaverifierCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EigendaverifierTransactor is an auto generated write-only Go binding around an Ethereum contract.
type EigendaverifierTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EigendaverifierFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type EigendaverifierFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EigendaverifierSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EigendaverifierSession struct {
	Contract     *Eigendaverifier  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EigendaverifierCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EigendaverifierCallerSession struct {
	Contract *EigendaverifierCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// EigendaverifierTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EigendaverifierTransactorSession struct {
	Contract     *EigendaverifierTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// EigendaverifierRaw is an auto generated low-level Go binding around an Ethereum contract.
type EigendaverifierRaw struct {
	Contract *Eigendaverifier // Generic contract binding to access the raw methods on
}

// EigendaverifierCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type EigendaverifierCallerRaw struct {
	Contract *EigendaverifierCaller // Generic read-only contract binding to access the raw methods on
}

// EigendaverifierTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type EigendaverifierTransactorRaw struct {
	Contract *EigendaverifierTransactor // Generic write-only contract binding to access the raw methods on
}

// NewEigendaverifier creates a new instance of Eigendaverifier, bound to a specific deployed contract.
func NewEigendaverifier(address common.Address, backend bind.ContractBackend) (*Eigendaverifier, error) {
	contract, err := bindEigendaverifier(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Eigendaverifier{EigendaverifierCaller: EigendaverifierCaller{contract: contract}, EigendaverifierTransactor: EigendaverifierTransactor{contract: contract}, EigendaverifierFilterer: EigendaverifierFilterer{contract: contract}}, nil
}

// NewEigendaverifierCaller creates a new read-only instance of Eigendaverifier, bound to a specific deployed contract.
func NewEigendaverifierCaller(address common.Address, caller bind.ContractCaller) (*EigendaverifierCaller, error) {
	contract, err := bindEigendaverifier(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EigendaverifierCaller{contract: contract}, nil
}

// NewEigendaverifierTransactor creates a new write-only instance of Eigendaverifier, bound to a specific deployed contract.
func NewEigendaverifierTransactor(address common.Address, transactor bind.ContractTransactor) (*EigendaverifierTransactor, error) {
	contract, err := bindEigendaverifier(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EigendaverifierTransactor{contract: contract}, nil
}

// NewEigendaverifierFilterer creates a new log filterer instance of Eigendaverifier, bound to a specific deployed contract.
func NewEigendaverifierFilterer(address common.Address, filterer bind.ContractFilterer) (*EigendaverifierFilterer, error) {
	contract, err := bindEigendaverifier(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EigendaverifierFilterer{contract: contract}, nil
}

// bindEigendaverifier binds a generic wrapper to an already deployed contract.
func bindEigendaverifier(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := EigendaverifierMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Eigendaverifier *EigendaverifierRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Eigendaverifier.Contract.EigendaverifierCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Eigendaverifier *EigendaverifierRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Eigendaverifier.Contract.EigendaverifierTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Eigendaverifier *EigendaverifierRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Eigendaverifier.Contract.EigendaverifierTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Eigendaverifier *EigendaverifierCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Eigendaverifier.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Eigendaverifier *EigendaverifierTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Eigendaverifier.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Eigendaverifier *EigendaverifierTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Eigendaverifier.Contract.contract.Transact(opts, method, params...)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_Eigendaverifier *EigendaverifierCaller) Admin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Eigendaverifier.contract.Call(opts, &out, "admin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_Eigendaverifier *EigendaverifierSession) Admin() (common.Address, error) {
	return _Eigendaverifier.Contract.Admin(&_Eigendaverifier.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_Eigendaverifier *EigendaverifierCallerSession) Admin() (common.Address, error) {
	return _Eigendaverifier.Contract.Admin(&_Eigendaverifier.CallOpts)
}

// DecodeBlobData is a free data retrieval call binding the contract method 0xaba4c80d.
//
// Solidity: function decodeBlobData(bytes data) pure returns((((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]),(uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes)) blobData)
func (_Eigendaverifier *EigendaverifierCaller) DecodeBlobData(opts *bind.CallOpts, data []byte) (EigenDAVerifierBlobData, error) {
	var out []interface{}
	err := _Eigendaverifier.contract.Call(opts, &out, "decodeBlobData", data)

	if err != nil {
		return *new(EigenDAVerifierBlobData), err
	}

	out0 := *abi.ConvertType(out[0], new(EigenDAVerifierBlobData)).(*EigenDAVerifierBlobData)

	return out0, err

}

// DecodeBlobData is a free data retrieval call binding the contract method 0xaba4c80d.
//
// Solidity: function decodeBlobData(bytes data) pure returns((((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]),(uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes)) blobData)
func (_Eigendaverifier *EigendaverifierSession) DecodeBlobData(data []byte) (EigenDAVerifierBlobData, error) {
	return _Eigendaverifier.Contract.DecodeBlobData(&_Eigendaverifier.CallOpts, data)
}

// DecodeBlobData is a free data retrieval call binding the contract method 0xaba4c80d.
//
// Solidity: function decodeBlobData(bytes data) pure returns((((uint256,uint256),uint32,(uint8,uint8,uint8,uint32)[]),(uint32,uint32,((bytes32,bytes,bytes,uint32),bytes32,uint32),bytes,bytes)) blobData)
func (_Eigendaverifier *EigendaverifierCallerSession) DecodeBlobData(data []byte) (EigenDAVerifierBlobData, error) {
	return _Eigendaverifier.Contract.DecodeBlobData(&_Eigendaverifier.CallOpts, data)
}

// GetDataAvailabilityProtocol is a free data retrieval call binding the contract method 0xb9c67c33.
//
// Solidity: function getDataAvailabilityProtocol() view returns(address)
func (_Eigendaverifier *EigendaverifierCaller) GetDataAvailabilityProtocol(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Eigendaverifier.contract.Call(opts, &out, "getDataAvailabilityProtocol")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetDataAvailabilityProtocol is a free data retrieval call binding the contract method 0xb9c67c33.
//
// Solidity: function getDataAvailabilityProtocol() view returns(address)
func (_Eigendaverifier *EigendaverifierSession) GetDataAvailabilityProtocol() (common.Address, error) {
	return _Eigendaverifier.Contract.GetDataAvailabilityProtocol(&_Eigendaverifier.CallOpts)
}

// GetDataAvailabilityProtocol is a free data retrieval call binding the contract method 0xb9c67c33.
//
// Solidity: function getDataAvailabilityProtocol() view returns(address)
func (_Eigendaverifier *EigendaverifierCallerSession) GetDataAvailabilityProtocol() (common.Address, error) {
	return _Eigendaverifier.Contract.GetDataAvailabilityProtocol(&_Eigendaverifier.CallOpts)
}

// GetProcotolName is a free data retrieval call binding the contract method 0xe4f17120.
//
// Solidity: function getProcotolName() pure returns(string)
func (_Eigendaverifier *EigendaverifierCaller) GetProcotolName(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Eigendaverifier.contract.Call(opts, &out, "getProcotolName")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetProcotolName is a free data retrieval call binding the contract method 0xe4f17120.
//
// Solidity: function getProcotolName() pure returns(string)
func (_Eigendaverifier *EigendaverifierSession) GetProcotolName() (string, error) {
	return _Eigendaverifier.Contract.GetProcotolName(&_Eigendaverifier.CallOpts)
}

// GetProcotolName is a free data retrieval call binding the contract method 0xe4f17120.
//
// Solidity: function getProcotolName() pure returns(string)
func (_Eigendaverifier *EigendaverifierCallerSession) GetProcotolName() (string, error) {
	return _Eigendaverifier.Contract.GetProcotolName(&_Eigendaverifier.CallOpts)
}

// PendingAdmin is a free data retrieval call binding the contract method 0x26782247.
//
// Solidity: function pendingAdmin() view returns(address)
func (_Eigendaverifier *EigendaverifierCaller) PendingAdmin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Eigendaverifier.contract.Call(opts, &out, "pendingAdmin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingAdmin is a free data retrieval call binding the contract method 0x26782247.
//
// Solidity: function pendingAdmin() view returns(address)
func (_Eigendaverifier *EigendaverifierSession) PendingAdmin() (common.Address, error) {
	return _Eigendaverifier.Contract.PendingAdmin(&_Eigendaverifier.CallOpts)
}

// PendingAdmin is a free data retrieval call binding the contract method 0x26782247.
//
// Solidity: function pendingAdmin() view returns(address)
func (_Eigendaverifier *EigendaverifierCallerSession) PendingAdmin() (common.Address, error) {
	return _Eigendaverifier.Contract.PendingAdmin(&_Eigendaverifier.CallOpts)
}

// VerifyMessage is a free data retrieval call binding the contract method 0x3b51be4b.
//
// Solidity: function verifyMessage(bytes32 , bytes data) view returns()
func (_Eigendaverifier *EigendaverifierCaller) VerifyMessage(opts *bind.CallOpts, arg0 [32]byte, data []byte) error {
	var out []interface{}
	err := _Eigendaverifier.contract.Call(opts, &out, "verifyMessage", arg0, data)

	if err != nil {
		return err
	}

	return err

}

// VerifyMessage is a free data retrieval call binding the contract method 0x3b51be4b.
//
// Solidity: function verifyMessage(bytes32 , bytes data) view returns()
func (_Eigendaverifier *EigendaverifierSession) VerifyMessage(arg0 [32]byte, data []byte) error {
	return _Eigendaverifier.Contract.VerifyMessage(&_Eigendaverifier.CallOpts, arg0, data)
}

// VerifyMessage is a free data retrieval call binding the contract method 0x3b51be4b.
//
// Solidity: function verifyMessage(bytes32 , bytes data) view returns()
func (_Eigendaverifier *EigendaverifierCallerSession) VerifyMessage(arg0 [32]byte, data []byte) error {
	return _Eigendaverifier.Contract.VerifyMessage(&_Eigendaverifier.CallOpts, arg0, data)
}

// AcceptAdminRole is a paid mutator transaction binding the contract method 0x8c3d7301.
//
// Solidity: function acceptAdminRole() returns()
func (_Eigendaverifier *EigendaverifierTransactor) AcceptAdminRole(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Eigendaverifier.contract.Transact(opts, "acceptAdminRole")
}

// AcceptAdminRole is a paid mutator transaction binding the contract method 0x8c3d7301.
//
// Solidity: function acceptAdminRole() returns()
func (_Eigendaverifier *EigendaverifierSession) AcceptAdminRole() (*types.Transaction, error) {
	return _Eigendaverifier.Contract.AcceptAdminRole(&_Eigendaverifier.TransactOpts)
}

// AcceptAdminRole is a paid mutator transaction binding the contract method 0x8c3d7301.
//
// Solidity: function acceptAdminRole() returns()
func (_Eigendaverifier *EigendaverifierTransactorSession) AcceptAdminRole() (*types.Transaction, error) {
	return _Eigendaverifier.Contract.AcceptAdminRole(&_Eigendaverifier.TransactOpts)
}

// SetDataAvailabilityProtocol is a paid mutator transaction binding the contract method 0x7cd76b8b.
//
// Solidity: function setDataAvailabilityProtocol(address newDataAvailabilityProtocol) returns()
func (_Eigendaverifier *EigendaverifierTransactor) SetDataAvailabilityProtocol(opts *bind.TransactOpts, newDataAvailabilityProtocol common.Address) (*types.Transaction, error) {
	return _Eigendaverifier.contract.Transact(opts, "setDataAvailabilityProtocol", newDataAvailabilityProtocol)
}

// SetDataAvailabilityProtocol is a paid mutator transaction binding the contract method 0x7cd76b8b.
//
// Solidity: function setDataAvailabilityProtocol(address newDataAvailabilityProtocol) returns()
func (_Eigendaverifier *EigendaverifierSession) SetDataAvailabilityProtocol(newDataAvailabilityProtocol common.Address) (*types.Transaction, error) {
	return _Eigendaverifier.Contract.SetDataAvailabilityProtocol(&_Eigendaverifier.TransactOpts, newDataAvailabilityProtocol)
}

// SetDataAvailabilityProtocol is a paid mutator transaction binding the contract method 0x7cd76b8b.
//
// Solidity: function setDataAvailabilityProtocol(address newDataAvailabilityProtocol) returns()
func (_Eigendaverifier *EigendaverifierTransactorSession) SetDataAvailabilityProtocol(newDataAvailabilityProtocol common.Address) (*types.Transaction, error) {
	return _Eigendaverifier.Contract.SetDataAvailabilityProtocol(&_Eigendaverifier.TransactOpts, newDataAvailabilityProtocol)
}

// TransferAdminRole is a paid mutator transaction binding the contract method 0xada8f919.
//
// Solidity: function transferAdminRole(address newPendingAdmin) returns()
func (_Eigendaverifier *EigendaverifierTransactor) TransferAdminRole(opts *bind.TransactOpts, newPendingAdmin common.Address) (*types.Transaction, error) {
	return _Eigendaverifier.contract.Transact(opts, "transferAdminRole", newPendingAdmin)
}

// TransferAdminRole is a paid mutator transaction binding the contract method 0xada8f919.
//
// Solidity: function transferAdminRole(address newPendingAdmin) returns()
func (_Eigendaverifier *EigendaverifierSession) TransferAdminRole(newPendingAdmin common.Address) (*types.Transaction, error) {
	return _Eigendaverifier.Contract.TransferAdminRole(&_Eigendaverifier.TransactOpts, newPendingAdmin)
}

// TransferAdminRole is a paid mutator transaction binding the contract method 0xada8f919.
//
// Solidity: function transferAdminRole(address newPendingAdmin) returns()
func (_Eigendaverifier *EigendaverifierTransactorSession) TransferAdminRole(newPendingAdmin common.Address) (*types.Transaction, error) {
	return _Eigendaverifier.Contract.TransferAdminRole(&_Eigendaverifier.TransactOpts, newPendingAdmin)
}

// EigendaverifierAcceptAdminRoleIterator is returned from FilterAcceptAdminRole and is used to iterate over the raw logs and unpacked data for AcceptAdminRole events raised by the Eigendaverifier contract.
type EigendaverifierAcceptAdminRoleIterator struct {
	Event *EigendaverifierAcceptAdminRole // Event containing the contract specifics and raw log

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
func (it *EigendaverifierAcceptAdminRoleIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EigendaverifierAcceptAdminRole)
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
		it.Event = new(EigendaverifierAcceptAdminRole)
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
func (it *EigendaverifierAcceptAdminRoleIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EigendaverifierAcceptAdminRoleIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EigendaverifierAcceptAdminRole represents a AcceptAdminRole event raised by the Eigendaverifier contract.
type EigendaverifierAcceptAdminRole struct {
	NewAdmin common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterAcceptAdminRole is a free log retrieval operation binding the contract event 0x056dc487bbf0795d0bbb1b4f0af523a855503cff740bfb4d5475f7a90c091e8e.
//
// Solidity: event AcceptAdminRole(address newAdmin)
func (_Eigendaverifier *EigendaverifierFilterer) FilterAcceptAdminRole(opts *bind.FilterOpts) (*EigendaverifierAcceptAdminRoleIterator, error) {

	logs, sub, err := _Eigendaverifier.contract.FilterLogs(opts, "AcceptAdminRole")
	if err != nil {
		return nil, err
	}
	return &EigendaverifierAcceptAdminRoleIterator{contract: _Eigendaverifier.contract, event: "AcceptAdminRole", logs: logs, sub: sub}, nil
}

// WatchAcceptAdminRole is a free log subscription operation binding the contract event 0x056dc487bbf0795d0bbb1b4f0af523a855503cff740bfb4d5475f7a90c091e8e.
//
// Solidity: event AcceptAdminRole(address newAdmin)
func (_Eigendaverifier *EigendaverifierFilterer) WatchAcceptAdminRole(opts *bind.WatchOpts, sink chan<- *EigendaverifierAcceptAdminRole) (event.Subscription, error) {

	logs, sub, err := _Eigendaverifier.contract.WatchLogs(opts, "AcceptAdminRole")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EigendaverifierAcceptAdminRole)
				if err := _Eigendaverifier.contract.UnpackLog(event, "AcceptAdminRole", log); err != nil {
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

// ParseAcceptAdminRole is a log parse operation binding the contract event 0x056dc487bbf0795d0bbb1b4f0af523a855503cff740bfb4d5475f7a90c091e8e.
//
// Solidity: event AcceptAdminRole(address newAdmin)
func (_Eigendaverifier *EigendaverifierFilterer) ParseAcceptAdminRole(log types.Log) (*EigendaverifierAcceptAdminRole, error) {
	event := new(EigendaverifierAcceptAdminRole)
	if err := _Eigendaverifier.contract.UnpackLog(event, "AcceptAdminRole", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EigendaverifierSetDataAvailabilityProtocolIterator is returned from FilterSetDataAvailabilityProtocol and is used to iterate over the raw logs and unpacked data for SetDataAvailabilityProtocol events raised by the Eigendaverifier contract.
type EigendaverifierSetDataAvailabilityProtocolIterator struct {
	Event *EigendaverifierSetDataAvailabilityProtocol // Event containing the contract specifics and raw log

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
func (it *EigendaverifierSetDataAvailabilityProtocolIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EigendaverifierSetDataAvailabilityProtocol)
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
		it.Event = new(EigendaverifierSetDataAvailabilityProtocol)
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
func (it *EigendaverifierSetDataAvailabilityProtocolIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EigendaverifierSetDataAvailabilityProtocolIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EigendaverifierSetDataAvailabilityProtocol represents a SetDataAvailabilityProtocol event raised by the Eigendaverifier contract.
type EigendaverifierSetDataAvailabilityProtocol struct {
	NewTrustedSequencer common.Address
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterSetDataAvailabilityProtocol is a free log retrieval operation binding the contract event 0xd331bd4c4cd1afecb94a225184bded161ff3213624ba4fb58c4f30c5a861144a.
//
// Solidity: event SetDataAvailabilityProtocol(address newTrustedSequencer)
func (_Eigendaverifier *EigendaverifierFilterer) FilterSetDataAvailabilityProtocol(opts *bind.FilterOpts) (*EigendaverifierSetDataAvailabilityProtocolIterator, error) {

	logs, sub, err := _Eigendaverifier.contract.FilterLogs(opts, "SetDataAvailabilityProtocol")
	if err != nil {
		return nil, err
	}
	return &EigendaverifierSetDataAvailabilityProtocolIterator{contract: _Eigendaverifier.contract, event: "SetDataAvailabilityProtocol", logs: logs, sub: sub}, nil
}

// WatchSetDataAvailabilityProtocol is a free log subscription operation binding the contract event 0xd331bd4c4cd1afecb94a225184bded161ff3213624ba4fb58c4f30c5a861144a.
//
// Solidity: event SetDataAvailabilityProtocol(address newTrustedSequencer)
func (_Eigendaverifier *EigendaverifierFilterer) WatchSetDataAvailabilityProtocol(opts *bind.WatchOpts, sink chan<- *EigendaverifierSetDataAvailabilityProtocol) (event.Subscription, error) {

	logs, sub, err := _Eigendaverifier.contract.WatchLogs(opts, "SetDataAvailabilityProtocol")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EigendaverifierSetDataAvailabilityProtocol)
				if err := _Eigendaverifier.contract.UnpackLog(event, "SetDataAvailabilityProtocol", log); err != nil {
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

// ParseSetDataAvailabilityProtocol is a log parse operation binding the contract event 0xd331bd4c4cd1afecb94a225184bded161ff3213624ba4fb58c4f30c5a861144a.
//
// Solidity: event SetDataAvailabilityProtocol(address newTrustedSequencer)
func (_Eigendaverifier *EigendaverifierFilterer) ParseSetDataAvailabilityProtocol(log types.Log) (*EigendaverifierSetDataAvailabilityProtocol, error) {
	event := new(EigendaverifierSetDataAvailabilityProtocol)
	if err := _Eigendaverifier.contract.UnpackLog(event, "SetDataAvailabilityProtocol", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EigendaverifierTransferAdminRoleIterator is returned from FilterTransferAdminRole and is used to iterate over the raw logs and unpacked data for TransferAdminRole events raised by the Eigendaverifier contract.
type EigendaverifierTransferAdminRoleIterator struct {
	Event *EigendaverifierTransferAdminRole // Event containing the contract specifics and raw log

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
func (it *EigendaverifierTransferAdminRoleIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EigendaverifierTransferAdminRole)
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
		it.Event = new(EigendaverifierTransferAdminRole)
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
func (it *EigendaverifierTransferAdminRoleIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EigendaverifierTransferAdminRoleIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EigendaverifierTransferAdminRole represents a TransferAdminRole event raised by the Eigendaverifier contract.
type EigendaverifierTransferAdminRole struct {
	NewPendingAdmin common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterTransferAdminRole is a free log retrieval operation binding the contract event 0xa5b56b7906fd0a20e3f35120dd8343db1e12e037a6c90111c7e42885e82a1ce6.
//
// Solidity: event TransferAdminRole(address newPendingAdmin)
func (_Eigendaverifier *EigendaverifierFilterer) FilterTransferAdminRole(opts *bind.FilterOpts) (*EigendaverifierTransferAdminRoleIterator, error) {

	logs, sub, err := _Eigendaverifier.contract.FilterLogs(opts, "TransferAdminRole")
	if err != nil {
		return nil, err
	}
	return &EigendaverifierTransferAdminRoleIterator{contract: _Eigendaverifier.contract, event: "TransferAdminRole", logs: logs, sub: sub}, nil
}

// WatchTransferAdminRole is a free log subscription operation binding the contract event 0xa5b56b7906fd0a20e3f35120dd8343db1e12e037a6c90111c7e42885e82a1ce6.
//
// Solidity: event TransferAdminRole(address newPendingAdmin)
func (_Eigendaverifier *EigendaverifierFilterer) WatchTransferAdminRole(opts *bind.WatchOpts, sink chan<- *EigendaverifierTransferAdminRole) (event.Subscription, error) {

	logs, sub, err := _Eigendaverifier.contract.WatchLogs(opts, "TransferAdminRole")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EigendaverifierTransferAdminRole)
				if err := _Eigendaverifier.contract.UnpackLog(event, "TransferAdminRole", log); err != nil {
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

// ParseTransferAdminRole is a log parse operation binding the contract event 0xa5b56b7906fd0a20e3f35120dd8343db1e12e037a6c90111c7e42885e82a1ce6.
//
// Solidity: event TransferAdminRole(address newPendingAdmin)
func (_Eigendaverifier *EigendaverifierFilterer) ParseTransferAdminRole(log types.Log) (*EigendaverifierTransferAdminRole, error) {
	event := new(EigendaverifierTransferAdminRole)
	if err := _Eigendaverifier.contract.UnpackLog(event, "TransferAdminRole", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
