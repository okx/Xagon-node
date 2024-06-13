// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package feijoapolygonzkevm

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

// PolygonRollupBaseFeijoaBlobData is an auto generated low-level Go binding around an user-defined struct.
type PolygonRollupBaseFeijoaBlobData struct {
	BlobType       uint8
	BlobTypeParams []byte
}

// FeijoapolygonzkevmMetaData contains all meta data concerning the Feijoapolygonzkevm contract.
var FeijoapolygonzkevmMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIPolygonZkEVMGlobalExitRootV2\",\"name\":\"_globalExitRootManager\",\"type\":\"address\"},{\"internalType\":\"contractIERC20Upgradeable\",\"name\":\"_pol\",\"type\":\"address\"},{\"internalType\":\"contractIPolygonZkEVMBridgeV2\",\"name\":\"_bridgeAddress\",\"type\":\"address\"},{\"internalType\":\"contractPolygonRollupManager\",\"name\":\"_rollupManager\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"BlobHashNotFound\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BlobTypeNotSupported\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FinalAccInputHashDoesNotMatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ForceBlobNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ForceBlobTimeoutNotExpired\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ForceBlobsAlreadyActive\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ForceBlobsDecentralized\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ForceBlobsNotAllowedOnEmergencyState\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ForceBlobsOverflow\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ForcedDataDoesNotMatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasTokenNetworkMustBeZeroOnEther\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"HaltTimeoutNotExpiredAfterEmergencyState\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"HugeTokenMetadataNotSupported\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidCommitmentAndProofLength\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitializeTransaction\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRangeForceBlobTimeout\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"Invalidl1InfoLeafIndex\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MaxTimestampSequenceInvalid\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotEnoughMaticAmount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotEnoughPOLAmount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyPendingAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyRollupManager\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyTrustedSequencer\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PointEvalutionPrecompiledFail\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SequenceZeroBlobs\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SequencedTimestampBelowForcedTimestamp\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TransactionsLengthAboveMax\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AcceptAdminRole\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"forceBlobNum\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"lastGlobalExitRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sequencer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"zkGasLimit\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"transactions\",\"type\":\"bytes\"}],\"name\":\"ForceBlob\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"transactions\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"lastGlobalExitRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sequencer\",\"type\":\"address\"}],\"name\":\"InitialSequenceBlobs\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"lastBlobSequenced\",\"type\":\"uint64\"}],\"name\":\"SequenceBlobs\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"numBlob\",\"type\":\"uint64\"}],\"name\":\"SequenceForceBlobs\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newForceBlobAddress\",\"type\":\"address\"}],\"name\":\"SetForceBlobAddress\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"newforceBlobTimeout\",\"type\":\"uint64\"}],\"name\":\"SetForceBlobTimeout\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"newNetworkName\",\"type\":\"string\"}],\"name\":\"SetNetworkName\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newTrustedSequencer\",\"type\":\"address\"}],\"name\":\"SetTrustedSequencer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"newTrustedSequencerURL\",\"type\":\"string\"}],\"name\":\"SetTrustedSequencerURL\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newPendingAdmin\",\"type\":\"address\"}],\"name\":\"TransferAdminRole\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"sequneceNum\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"aggregator\",\"type\":\"address\"}],\"name\":\"VerifyBlobs\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"GLOBAL_EXIT_ROOT_MANAGER_L2\",\"outputs\":[{\"internalType\":\"contractIBasePolygonZkEVMGlobalExitRoot\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"INITIALIZE_TX_BRIDGE_LIST_LEN_LEN\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"INITIALIZE_TX_BRIDGE_PARAMS\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"INITIALIZE_TX_BRIDGE_PARAMS_AFTER_BRIDGE_ADDRESS\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"INITIALIZE_TX_BRIDGE_PARAMS_AFTER_BRIDGE_ADDRESS_EMPTY_METADATA\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"INITIALIZE_TX_CONSTANT_BYTES\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"INITIALIZE_TX_CONSTANT_BYTES_EMPTY_METADATA\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"INITIALIZE_TX_DATA_LEN_EMPTY_METADATA\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"INITIALIZE_TX_EFFECTIVE_PERCENTAGE\",\"outputs\":[{\"internalType\":\"bytes1\",\"name\":\"\",\"type\":\"bytes1\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"MAX_SEQUENCE_TIMESTAMP_FORCED\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"POINT_EVALUATION_PRECOMPILE_ADDRESS\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"SIGNATURE_INITIALIZE_TX_R\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"SIGNATURE_INITIALIZE_TX_S\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"SIGNATURE_INITIALIZE_TX_V\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"TIMESTAMP_RANGE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ZK_GAS_LIMIT_BATCH\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptAdminRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"admin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"bridgeAddress\",\"outputs\":[{\"internalType\":\"contractIPolygonZkEVMBridgeV2\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"calculatePolPerForcedZkGas\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"blobData\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"polAmount\",\"type\":\"uint256\"}],\"name\":\"forceBlob\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"forceBlobAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"forceBlobTimeout\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"forcedBlobs\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"hashedForcedBlobData\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"forcedTimestamp\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"gasTokenAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"gasTokenNetwork\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"networkID\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"_gasTokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"_gasTokenNetwork\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"_gasTokenMetadata\",\"type\":\"bytes\"}],\"name\":\"generateInitializeTransaction\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"globalExitRootManager\",\"outputs\":[{\"internalType\":\"contractIPolygonZkEVMGlobalExitRootV2\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_admin\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"sequencer\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"networkID\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"_gasTokenAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"sequencerURL\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_networkName\",\"type\":\"string\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastAccInputHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastForceBlob\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastForceBlobSequenced\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"networkName\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"lastVerifiedSequenceNum\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"aggregator\",\"type\":\"address\"}],\"name\":\"onVerifySequences\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pendingAdmin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pol\",\"outputs\":[{\"internalType\":\"contractIERC20Upgradeable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rollupManager\",\"outputs\":[{\"internalType\":\"contractPolygonRollupManager\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint8\",\"name\":\"blobType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"blobTypeParams\",\"type\":\"bytes\"}],\"internalType\":\"structPolygonRollupBaseFeijoa.BlobData[]\",\"name\":\"blobs\",\"type\":\"tuple[]\"},{\"internalType\":\"address\",\"name\":\"l2Coinbase\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"finalAccInputHash\",\"type\":\"bytes32\"}],\"name\":\"sequenceBlobs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint8\",\"name\":\"blobType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"blobTypeParams\",\"type\":\"bytes\"}],\"internalType\":\"structPolygonRollupBaseFeijoa.BlobData[]\",\"name\":\"blobs\",\"type\":\"tuple[]\"}],\"name\":\"sequenceForceBlobs\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newForceBlobAddress\",\"type\":\"address\"}],\"name\":\"setForceBlobAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"newforceBlobTimeout\",\"type\":\"uint64\"}],\"name\":\"setForceBlobTimeout\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"newNetworkName\",\"type\":\"string\"}],\"name\":\"setNetworkName\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newTrustedSequencer\",\"type\":\"address\"}],\"name\":\"setTrustedSequencer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"newTrustedSequencerURL\",\"type\":\"string\"}],\"name\":\"setTrustedSequencerURL\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newPendingAdmin\",\"type\":\"address\"}],\"name\":\"transferAdminRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"trustedSequencer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"trustedSequencerURL\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x61010060405234801562000011575f80fd5b5060405162004b0438038062004b0483398101604081905262000034916200006f565b6001600160a01b0393841660a052918316608052821660c0521660e052620000d4565b6001600160a01b03811681146200006c575f80fd5b50565b5f805f806080858703121562000083575f80fd5b8451620000908162000057565b6020860151909450620000a38162000057565b6040860151909350620000b68162000057565b6060860151909250620000c98162000057565b939692955090935050565b60805160a05160c05160e051614946620001be5f395f818161055c0152818161133a015281816113fe015281816114200152818161154d0152818161184001528181611e2f015281816121030152818161226e01528181612489015281816128b7015281816129900152818161307d015261315101525f81816106f701528181610a8001528181611a3401528181611b0901528181612b390152612c4101525f81816107be01528181610ca801528181610ea401528181611c8901526132a301525f81816108030152818161089901528181611383015281816114cc015261327801526149465ff3fe608060405234801561000f575f80fd5b5060043610610304575f3560e01c80637160c5f71161019d578063b0afe154116100e8578063d02103ca11610093578063e46761c41161006e578063e46761c4146107fe578063f35dda4714610825578063f851a4401461082d575f80fd5b8063d02103ca146107b9578063d2a679b7146107e0578063d7bc90ff146107f3575f80fd5b8063c7fffd4b116100c3578063c7fffd4b1461077e578063c89e42df14610786578063cfa8ed4714610799575f80fd5b8063b0afe1541461073f578063b45bd7f91461074b578063c0cad3021461076b575f80fd5b806393932a9111610148578063a3c573eb11610123578063a3c573eb146106f2578063a652f26c14610719578063ada8f9191461072c575f80fd5b806393932a91146106b05780639b0e35a5146106c35780639e001877146106d7575f80fd5b8063838a250311610178578063838a25031461068a578063889cfd7a146106955780638c3d7301146106a8575f80fd5b80637160c5f71461062c578063730c8e211461063b5780637a5460c51461064e575f80fd5b80633e41062e1161025d578063542028d5116102085780636e05d2cd116101e35780636e05d2cd146105fd5780636ff512cc146106065780637125702214610619575f80fd5b8063542028d5146105cd57806366e7bb1a146105d5578063676870d2146105f5575f80fd5b806349b7b8021161023857806349b7b802146105575780634bd410651461057e57806352bdeb6d14610591575f80fd5b80633e41062e146104ef57806340b5de6c146104f757806342308fab1461054f575f80fd5b806326782247116102bd57806338793b4f1161029857806338793b4f1461047d5780633c351e10146104925780633cbc795b146104b2575f80fd5b806326782247146103a95780632a6688ee146103ee5780632c2251db1461043c575f80fd5b806305835f37116102ed57806305835f371461033e578063107bf28c1461038757806311e892d41461038f575f80fd5b80630350896314610308578063042b0f0614610328575b5f80fd5b610310602081565b60405161ffff90911681526020015b60405180910390f35b610330610852565b60405190815260200161031f565b61037a6040518060400160405280600881526020017f80808401c9c3809400000000000000000000000000000000000000000000000081525081565b60405161031f9190613991565b61037a610966565b61039760f981565b60405160ff909116815260200161031f565b6001546103c99073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161031f565b61041e6103fc3660046139c2565b60066020525f90815260409020805460019091015467ffffffffffffffff1682565b6040805192835267ffffffffffffffff90911660208301520161031f565b60075461046490700100000000000000000000000000000000900467ffffffffffffffff1681565b60405167ffffffffffffffff909116815260200161031f565b61049061048b366004613a46565b6109f2565b005b6009546103c99073ffffffffffffffffffffffffffffffffffffffff1681565b6009546104da9074010000000000000000000000000000000000000000900463ffffffff1681565b60405163ffffffff909116815260200161031f565b6103c9600a81565b61051e7fff0000000000000000000000000000000000000000000000000000000000000081565b6040517fff00000000000000000000000000000000000000000000000000000000000000909116815260200161031f565b610330602481565b6103c97f000000000000000000000000000000000000000000000000000000000000000081565b61049061058c366004613a9f565b611649565b61037a6040518060400160405280600281526020017f80b800000000000000000000000000000000000000000000000000000000000081525081565b61037a611768565b6008546103c99073ffffffffffffffffffffffffffffffffffffffff1681565b610310601f81565b61033060055481565b610490610614366004613a9f565b611775565b610490610627366004613bde565b61183e565b61046467ffffffffffffffff81565b6104906106493660046139c2565b612064565b61037a6040518060400160405280600281526020017f80b900000000000000000000000000000000000000000000000000000000000081525081565b6104646305f5e10081565b6104906106a3366004613c85565b61226c565b61049061233b565b6104906106be366004613cc4565b61240d565b6007546104649067ffffffffffffffff1681565b6103c973a40d5f56745a118d0906a34e69aec8c0db1cb8fa81565b6103c97f000000000000000000000000000000000000000000000000000000000000000081565b61037a610727366004613d03565b612a3b565b61049061073a366004613a9f565b612e19565b6103306405ca1ab1e081565b6007546104649068010000000000000000900467ffffffffffffffff1681565b610490610779366004613d74565b612ee2565b61039760e481565b610490610794366004613d74565b612f74565b6002546103c99073ffffffffffffffffffffffffffffffffffffffff1681565b6103c97f000000000000000000000000000000000000000000000000000000000000000081565b6104906107ee366004613da6565b613006565b610330635ca1ab1e81565b6103c97f000000000000000000000000000000000000000000000000000000000000000081565b610397601b81565b5f546103c99062010000900473ffffffffffffffffffffffffffffffffffffffff1681565b6040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201525f90819073ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016906370a0823190602401602060405180830381865afa1580156108de573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906109029190613e17565b6007549091505f9061092c9067ffffffffffffffff68010000000000000000820481169116613e5b565b67ffffffffffffffff169050805f03610947575f9250505090565b6109556305f5e10082613e83565b61095f9083613ea0565b9250505090565b6004805461097390613ed8565b80601f016020809104026020016040519081016040528092919081815260200182805461099f90613ed8565b80156109ea5780601f106109c1576101008083540402835291602001916109ea565b820191905f5260205f20905b8154815290600101906020018083116109cd57829003601f168201915b505050505081565b60025473ffffffffffffffffffffffffffffffffffffffff163314610a43576040517f11e7be1500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b825f819003610a7e576040517fc8ea63df00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166379e2cf976040518163ffffffff1660e01b81526004015f604051808303815f87803b158015610ae3575f80fd5b505af1158015610af5573d5f803e3d5ffd5b50506007546005546801000000000000000090910467ffffffffffffffff1692509050815f805b858110156112a957368a8a83818110610b3757610b37613f29565b9050602002810190610b499190613f56565b90506002610b5a6020830183613fa7565b60ff161115610b95576040517f1d29ea1400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610ba26020820182613fa7565b60ff165f03610dcc57885f808080610bbd6020870187613fc0565b810190610bca9190614021565b9350935093509350602442610bdf919061404f565b8467ffffffffffffffff161115610c22576040517f0a00feb300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6201d4c081511115610c60576040517fa29a6c7c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b805160208201205f63ffffffff841615610d5f576040517f25eaabaf00000000000000000000000000000000000000000000000000000000815263ffffffff851660048201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906325eaabaf90602401602060405180830381865afa158015610d02573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610d269190613e17565b905080610d5f576040517f6a80570500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8b8482888a89610d7260208f018f613fa7565b604051610d9097969594939291905f9081908c908290602001614062565b604051602081830303815290604052805190602001209b508467ffffffffffffffff168a610dbe919061404f565b9950505050505050506112a0565b610dd96020820182613fa7565b60ff1660010361112557885f808080808080610df860208a018a613fc0565b810190610e059190614154565b9650965096509650965096509650602442610e20919061404f565b8767ffffffffffffffff161115610e63576040517f0a00feb300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f63ffffffff861615610f5b576040517f25eaabaf00000000000000000000000000000000000000000000000000000000815263ffffffff871660048201527f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16906325eaabaf90602401602060405180830381865afa158015610efe573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190610f229190613e17565b905080610f5b576040517f6a80570500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8151606003610f96576040517fbdb8fa9200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b844980610fcf576040517fec3601b300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f600a73ffffffffffffffffffffffffffffffffffffffff1682878787604051602001610fff94939291906141e0565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529082905261103791614213565b5f60405180830381855afa9150503d805f811461106f576040519150601f19603f3d011682016040523d82523d5f602084013e611074565b606091505b50509050806110af576040517f6df0d0e500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b50508d86828a8c8b8f5f0160208101906110c99190613fa7565b6040516110e797969594939291908c908c905f908190602001614062565b604051602081830303815290604052805190602001209d508667ffffffffffffffff168c611115919061404f565b9b505050505050505050506112a0565b5f806111346020840184613fc0565b8101906111419190614224565b91509150878061115090614244565b9850505f8282604051602001611170929190918252602082015260400190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152918152815160209283012067ffffffffffffffff8c165f908152600690935291205490915081146111f8576040517fce3d755e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60065f8a67ffffffffffffffff1667ffffffffffffffff1681526020019081526020015f205f8082015f9055600182015f6101000a81549067ffffffffffffffff02191690555050875f805f1b67ffffffffffffffff8f6305f5e100895f0160208101906112669190613fa7565b60405161128497969594939291905f9081908d908d90602001614062565b6040516020818303038152906040528051906020012097505050505b50600101610b1c565b5060075467ffffffffffffffff90811690851611156112f4576040517ff32726dd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60058390555f67ffffffffffffffff858116908416146113eb575f6113198487613e5b565b90506113296305f5e1008261426a565b67ffffffffffffffff1691506113aa7f000000000000000000000000000000000000000000000000000000000000000083611362610852565b61136c9190613e83565b73ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169190613573565b50600780547fffffffffffffffffffffffffffffffff0000000000000000ffffffffffffffff166801000000000000000067ffffffffffffffff8816021790555b5f6113f6828461404f565b90506114f4337f0000000000000000000000000000000000000000000000000000000000000000837f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663f4174a176040518163ffffffff1660e01b8152600401602060405180830381865afa158015611487573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906114ab9190613e17565b6114b59190613e83565b73ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001692919061364c565b6040517ffe01d89e0000000000000000000000000000000000000000000000000000000081526fffffffffffffffffffffffffffffffff8216600482015267ffffffffffffffff88166024820152604481018690525f907f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff169063fe01d89e906064016020604051808303815f875af11580156115a8573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906115cc9190614296565b9050888614611607576040517fda5bceb900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60405167ffffffffffffffff8216907f470f4ca4b003755c839b80ab00c3efbeb69d6eafec00e1a3677482933ec1fd0c905f90a2505050505050505050505050565b5f5462010000900473ffffffffffffffffffffffffffffffffffffffff16331461169f576040517f4755657900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60085473ffffffffffffffffffffffffffffffffffffffff166116ee576040517f6958969600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600880547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83169081179091556040519081527f2261b2af55eeb3b995b5e300659fa8e59827ff8fc99ff3a5baf5af0835aab9dd906020015b60405180910390a150565b6003805461097390613ed8565b5f5462010000900473ffffffffffffffffffffffffffffffffffffffff1633146117cb576040517f4755657900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600280547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83169081179091556040519081527ff54144f9611984021529f814a1cb6a41e22c58351510a0d9f7e822618abb9cc09060200161175d565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1633146118ad576040517fb9b3a2c800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f54610100900460ff16158080156118cb57505f54600160ff909116105b806118e45750303b1580156118e457505f5460ff166001145b611975576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201527f647920696e697469616c697a656400000000000000000000000000000000000060648201526084015b60405180910390fd5b5f80547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016600117905580156119d1575f80547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff166101001790555b606073ffffffffffffffffffffffffffffffffffffffff851615611c2e576040517fc00f14ab00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff86811660048301527f0000000000000000000000000000000000000000000000000000000000000000169063c00f14ab906024015f60405180830381865afa158015611a78573d5f803e3d5ffd5b505050506040513d5f823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168201604052611abd91908101906142b1565b6040517f318aee3d00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff87811660048301529192505f9182917f00000000000000000000000000000000000000000000000000000000000000009091169063318aee3d906024016040805180830381865afa158015611b4f573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190611b739190614323565b915091508163ffffffff165f14611bea576009805463ffffffff841674010000000000000000000000000000000000000000027fffffffffffffffff00000000000000000000000000000000000000000000000090911673ffffffffffffffffffffffffffffffffffffffff841617179055611c2b565b600980547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff89161790555b50505b6009545f90611c7590889073ffffffffffffffffffffffffffffffffffffffff81169074010000000000000000000000000000000000000000900463ffffffff1685612a3b565b90505f818051906020012090505f4290505f7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16633ed691ef6040518163ffffffff1660e01b8152600401602060405180830381865afa158015611cf0573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190611d149190613e17565b90505f80808067ffffffffffffffff8f6305f5e100600284808c8b8d611d3b60014361435b565b40604051602001611d849392919092835260c09190911b7fffffffffffffffff000000000000000000000000000000000000000000000000166020830152602882015260480190565b60405160208183030381529060405280519060200120604051602001611db49b9a99989796959493929190614062565b604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe001815290829052805160209091012060058190557ffe01d89e0000000000000000000000000000000000000000000000000000000082526305f5e1006004830152600160248301526044820181905291507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff169063fe01d89e906064016020604051808303815f875af1158015611e8a573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190611eae9190614296565b508c5f60026101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508b60025f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508860039081611f3e91906143b9565b506004611f4b89826143b9565b508c60085f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555062069780600760106101000a81548167ffffffffffffffff021916908367ffffffffffffffff1602179055507ffa56300f6f91d53e1c1283e56307c169d72b14a75380df3ecbb5b31b498d3d1e85838e604051611feb939291906144d5565b60405180910390a1505050505050801561205b575f80547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b50505050505050565b5f5462010000900473ffffffffffffffffffffffffffffffffffffffff1633146120ba576040517f4755657900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b62093a8067ffffffffffffffff82161115612101576040517fd2438ff800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166315064c966040518163ffffffff1660e01b8152600401602060405180830381865afa15801561216a573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061218e9190614513565b6121ef5760075467ffffffffffffffff7001000000000000000000000000000000009091048116908216106121ef576040517fd2438ff800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600780547fffffffffffffffff0000000000000000ffffffffffffffffffffffffffffffff1670010000000000000000000000000000000067ffffffffffffffff8416908102919091179091556040519081527fa6db492cb43063288b0b5d7c271f8df34607c41fc9347c0664e1ce325cc728e89060200161175d565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1633146122db576040517fb9b3a2c800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8073ffffffffffffffffffffffffffffffffffffffff168367ffffffffffffffff167fb19baa6f6271636400b99e9e5b3289ec1e0d74e6204a27f296cc4715ff9ded558460405161232e91815260200190565b60405180910390a3505050565b60015473ffffffffffffffffffffffffffffffffffffffff16331461238c576040517fd1ec4b2300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001545f80547fffffffffffffffffffff0000000000000000000000000000000000000000ffff1673ffffffffffffffffffffffffffffffffffffffff9092166201000081029290921790556040519081527f056dc487bbf0795d0bbb1b4f0af523a855503cff740bfb4d5475f7a90c091e8e9060200160405180910390a1565b60085473ffffffffffffffffffffffffffffffffffffffff16801580159061244b575073ffffffffffffffffffffffffffffffffffffffff81163314155b15612482576040517f59c46bd200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b4262093a807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166330c27dde6040518163ffffffff1660e01b8152600401602060405180830381865afa1580156124f0573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906125149190614296565b61251e9190614532565b67ffffffffffffffff161115612560576040517f3d49ed4c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b815f81900361259b576040517fc8ea63df00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60075467ffffffffffffffff808216916125c39184916801000000000000000090041661404f565b11156125fb576040517ff32726dd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6007546005546801000000000000000090910467ffffffffffffffff16905f5b838110156128a0575f87878381811061263657612636613f29565b90506020028101906126489190613f56565b61265190614553565b90508361265d81614244565b945050805f015160ff166002146126a0576040517f1d29ea1400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f8082602001518060200190518101906126ba91906145c2565b9150915085806126c990614244565b9650505f82826040516020016126e9929190918252602082015260400190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152918152815160209283012067ffffffffffffffff8a165f90815260069093529120549091508114612771576040517fce3d755e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61277c60018961435b565b85036128085760075467ffffffffffffffff8881165f9081526006602052604090206001015442926127c69270010000000000000000000000000000000090910481169116614532565b67ffffffffffffffff161115612808576040517fc643d3d400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff8088165f90815260066020908152604080832083815560010180547fffffffffffffffffffffffffffffffffffffffffffffffff000000000000000016905587519051612877948b94938493919233926305f5e100929091869182918e918e9101614062565b60405160208183030381529060405280519060200120955050505050808060010191505061261b565b505f6128b06305f5e10085613e83565b90506128df7f000000000000000000000000000000000000000000000000000000000000000082611362610852565b60058290556007805467ffffffffffffffff85811668010000000000000000027fffffffffffffffffffffffffffffffff0000000000000000ffffffffffffffff909216919091179091556040517ffe01d89e0000000000000000000000000000000000000000000000000000000081526fffffffffffffffffffffffffffffffff831660048201529085166024820152604481018390525f9073ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169063fe01d89e906064016020604051808303815f875af11580156129d6573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906129fa9190614296565b60405190915067ffffffffffffffff8216907f049b259b0b684f32f1d8b43d76cf6cb3c674b94697bda3290f6ec63252cfe892905f90a25050505050505050565b60605f85858573a40d5f56745a118d0906a34e69aec8c0db1cb8fa5f87604051602401612a6d969594939291906145e4565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167ff811bff70000000000000000000000000000000000000000000000000000000017905283519091506060905f03612bbd5760f9601f8351612b019190614646565b6040518060400160405280600881526020017f80808401c9c380940000000000000000000000000000000000000000000000008152507f00000000000000000000000000000000000000000000000000000000000000006040518060400160405280600281526020017f80b800000000000000000000000000000000000000000000000000000000000081525060e487604051602001612ba79796959493929190614661565b6040516020818303038152906040529050612cc1565b815161ffff1015612bfa576040517f248b8f8200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b815160f9612c09602083614646565b6040518060400160405280600881526020017f80808401c9c380940000000000000000000000000000000000000000000000008152507f00000000000000000000000000000000000000000000000000000000000000006040518060400160405280600281526020017f80b90000000000000000000000000000000000000000000000000000000000008152508588604051602001612cae9796959493929190614743565b6040516020818303038152906040529150505b8051602080830191909120604080515f80825293810180835292909252601b908201526405ca1ab1e06060820152635ca1ab1e608082015260019060a0016020604051602081039080840390855afa158015612d1f573d5f803e3d5ffd5b50506040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0015191505073ffffffffffffffffffffffffffffffffffffffff8116612d97576040517fcd16196600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040515f90612ddc9084906405ca1ab1e090635ca1ab1e90601b907fff0000000000000000000000000000000000000000000000000000000000000090602001614825565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190529450505050505b949350505050565b5f5462010000900473ffffffffffffffffffffffffffffffffffffffff163314612e6f576040517f4755657900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83169081179091556040519081527fa5b56b7906fd0a20e3f35120dd8343db1e12e037a6c90111c7e42885e82a1ce69060200161175d565b5f5462010000900473ffffffffffffffffffffffffffffffffffffffff163314612f38576040517f4755657900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6004612f4482826143b9565b507fcc3b37f0de47ea5ce245c3502f0d4e414c34664023b8463db2fe451fee5e69928160405161175d9190613991565b5f5462010000900473ffffffffffffffffffffffffffffffffffffffff163314612fca576040517f4755657900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6003612fd682826143b9565b507f6b8f723a4c7a5335cafae8a598a0aa0301be1387c037dccc085b62add6448b208160405161175d9190613991565b60085473ffffffffffffffffffffffffffffffffffffffff168015801590613044575073ffffffffffffffffffffffffffffffffffffffff81163314155b1561307b576040517f59c46bd200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166315064c966040518163ffffffff1660e01b8152600401602060405180830381865afa1580156130e4573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906131089190614513565b1561313f576040517f65afbc4900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f6305f5e10067ffffffffffffffff167f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166302f3fa606040518163ffffffff1660e01b8152600401602060405180830381865afa1580156131b8573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906131dc9190613e17565b6131e69190613e83565b905082811115613222576040517f2354600f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b61138884111561325e576040517fa29a6c7c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6132a073ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001633308461364c565b5f7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16633ed691ef6040518163ffffffff1660e01b8152600401602060405180830381865afa15801561330a573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019061332e9190613e17565b6007805491925067ffffffffffffffff909116905f61334c83614244565b91906101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550505f8142600143613383919061435b565b406040516020016133cc9392919092835260c09190911b7fffffffffffffffff000000000000000000000000000000000000000000000000166020830152602882015260480190565b604051602081830303815290604052805190602001209050604051806040016040528088886040516133ff929190614880565b6040805191829003822060208301528101849052606001604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0018152918152815160209283012083524267ffffffffffffffff9081169383019390935260075483165f908152600683522083518155920151600190920180547fffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000169290911691909117905532330361351657600754604080518481523360208201526305f5e100818301526080606082018190525f90820152905167ffffffffffffffff909216917fb18d758550a6ed34847584be90f0a34b261d8b65bb790891103d5e255aced8b29181900360a00190a261205b565b60075460405167ffffffffffffffff909116907fb18d758550a6ed34847584be90f0a34b261d8b65bb790891103d5e255aced8b29061356290859033906305f5e100908d908d9061488f565b60405180910390a250505050505050565b60405173ffffffffffffffffffffffffffffffffffffffff83166024820152604481018290526136479084907fa9059cbb00000000000000000000000000000000000000000000000000000000906064015b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff00000000000000000000000000000000000000000000000000000000909316929092179091526136b0565b505050565b60405173ffffffffffffffffffffffffffffffffffffffff808516602483015283166044820152606481018290526136aa9085907f23b872dd00000000000000000000000000000000000000000000000000000000906084016135c5565b50505050565b5f613711826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65648152508573ffffffffffffffffffffffffffffffffffffffff166137bb9092919063ffffffff16565b805190915015613647578080602001905181019061372f9190614513565b613647576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f74207375636365656400000000000000000000000000000000000000000000606482015260840161196c565b6060612e1184845f85855f808673ffffffffffffffffffffffffffffffffffffffff1685876040516137ed9190614213565b5f6040518083038185875af1925050503d805f8114613827576040519150601f19603f3d011682016040523d82523d5f602084013e61382c565b606091505b509150915061383d87838387613848565b979650505050505050565b606083156138dd5782515f036138d65773ffffffffffffffffffffffffffffffffffffffff85163b6138d6576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e7472616374000000604482015260640161196c565b5081612e11565b612e1183838151156138f25781518083602001fd5b806040517f08c379a000000000000000000000000000000000000000000000000000000000815260040161196c9190613991565b5f5b83811015613940578181015183820152602001613928565b50505f910152565b5f815180845261395f816020860160208601613926565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b602081525f6139a36020830184613948565b9392505050565b67ffffffffffffffff811681146139bf575f80fd5b50565b5f602082840312156139d2575f80fd5b81356139a3816139aa565b5f8083601f8401126139ed575f80fd5b50813567ffffffffffffffff811115613a04575f80fd5b6020830191508360208260051b8501011115613a1e575f80fd5b9250929050565b73ffffffffffffffffffffffffffffffffffffffff811681146139bf575f80fd5b5f805f8060608587031215613a59575f80fd5b843567ffffffffffffffff811115613a6f575f80fd5b613a7b878288016139dd565b9095509350506020850135613a8f81613a25565b9396929550929360400135925050565b5f60208284031215613aaf575f80fd5b81356139a381613a25565b63ffffffff811681146139bf575f80fd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715613b3f57613b3f613acb565b604052919050565b5f67ffffffffffffffff821115613b6057613b60613acb565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b5f82601f830112613b9b575f80fd5b8135613bae613ba982613b47565b613af8565b818152846020838601011115613bc2575f80fd5b816020850160208301375f918101602001919091529392505050565b5f805f805f8060c08789031215613bf3575f80fd5b8635613bfe81613a25565b95506020870135613c0e81613a25565b94506040870135613c1e81613aba565b93506060870135613c2e81613a25565b9250608087013567ffffffffffffffff80821115613c4a575f80fd5b613c568a838b01613b8c565b935060a0890135915080821115613c6b575f80fd5b50613c7889828a01613b8c565b9150509295509295509295565b5f805f60608486031215613c97575f80fd5b8335613ca2816139aa565b9250602084013591506040840135613cb981613a25565b809150509250925092565b5f8060208385031215613cd5575f80fd5b823567ffffffffffffffff811115613ceb575f80fd5b613cf7858286016139dd565b90969095509350505050565b5f805f8060808587031215613d16575f80fd5b8435613d2181613aba565b93506020850135613d3181613a25565b92506040850135613d4181613aba565b9150606085013567ffffffffffffffff811115613d5c575f80fd5b613d6887828801613b8c565b91505092959194509250565b5f60208284031215613d84575f80fd5b813567ffffffffffffffff811115613d9a575f80fd5b612e1184828501613b8c565b5f805f60408486031215613db8575f80fd5b833567ffffffffffffffff80821115613dcf575f80fd5b818601915086601f830112613de2575f80fd5b813581811115613df0575f80fd5b876020828501011115613e01575f80fd5b6020928301989097509590910135949350505050565b5f60208284031215613e27575f80fd5b5051919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b67ffffffffffffffff828116828216039080821115613e7c57613e7c613e2e565b5092915050565b8082028115828204841417613e9a57613e9a613e2e565b92915050565b5f82613ed3577f4e487b71000000000000000000000000000000000000000000000000000000005f52601260045260245ffd5b500490565b600181811c90821680613eec57607f821691505b602082108103613f23577f4e487b71000000000000000000000000000000000000000000000000000000005f52602260045260245ffd5b50919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b5f82357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc1833603018112613f88575f80fd5b9190910192915050565b803560ff81168114613fa2575f80fd5b919050565b5f60208284031215613fb7575f80fd5b6139a382613f92565b5f8083357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1843603018112613ff3575f80fd5b83018035915067ffffffffffffffff82111561400d575f80fd5b602001915036819003821315613a1e575f80fd5b5f805f8060808587031215614034575f80fd5b843561403f816139aa565b93506020850135613d31816139aa565b80820180821115613e9a57613e9a613e2e565b8b81527fffffffff000000000000000000000000000000000000000000000000000000008b60e01b1660208201528960248201525f7fffffffffffffffff000000000000000000000000000000000000000000000000808b60c01b1660448401527fffffffffffffffffffffffffffffffffffffffff0000000000000000000000008a60601b16604c840152808960c01b1660608401525061412b606883018860f81b7fff00000000000000000000000000000000000000000000000000000000000000169052565b506069810194909452608984019290925260a983015260c982015260e901979650505050505050565b5f805f805f805f60e0888a03121561416a575f80fd5b8735614175816139aa565b96506020880135614185816139aa565b9550604088013561419581613aba565b9450606088013593506080880135925060a0880135915060c088013567ffffffffffffffff8111156141c5575f80fd5b6141d18a828b01613b8c565b91505092959891949750929550565b8481528360208201528260408201525f8251614203816060850160208701613926565b9190910160600195945050505050565b5f8251613f88818460208701613926565b5f8060408385031215614235575f80fd5b50508035926020909101359150565b5f67ffffffffffffffff80831681810361426057614260613e2e565b6001019392505050565b67ffffffffffffffff81811683821602808216919082811461428e5761428e613e2e565b505092915050565b5f602082840312156142a6575f80fd5b81516139a3816139aa565b5f602082840312156142c1575f80fd5b815167ffffffffffffffff8111156142d7575f80fd5b8201601f810184136142e7575f80fd5b80516142f5613ba982613b47565b818152856020838501011115614309575f80fd5b61431a826020830160208601613926565b95945050505050565b5f8060408385031215614334575f80fd5b825161433f81613aba565b602084015190925061435081613a25565b809150509250929050565b81810381811115613e9a57613e9a613e2e565b601f82111561364757805f5260205f20601f840160051c810160208510156143935750805b601f840160051c820191505b818110156143b2575f815560010161439f565b5050505050565b815167ffffffffffffffff8111156143d3576143d3613acb565b6143e7816143e18454613ed8565b8461436e565b602080601f831160018114614439575f84156144035750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b1785556144cd565b5f858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b8281101561448557888601518255948401946001909101908401614466565b50858210156144c157878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b505060018460011b0185555b505050505050565b606081525f6144e76060830186613948565b905083602083015273ffffffffffffffffffffffffffffffffffffffff83166040830152949350505050565b5f60208284031215614523575f80fd5b815180151581146139a3575f80fd5b67ffffffffffffffff818116838216019080821115613e7c57613e7c613e2e565b5f60408236031215614563575f80fd5b6040516040810167ffffffffffffffff828210818311171561458757614587613acb565b8160405261459485613f92565b835260208501359150808211156145a9575f80fd5b506145b636828601613b8c565b60208301525092915050565b5f80604083850312156145d3575f80fd5b505080516020909101519092909150565b5f63ffffffff808916835273ffffffffffffffffffffffffffffffffffffffff8089166020850152818816604085015280871660608501528086166080850152505060c060a083015261463a60c0830184613948565b98975050505050505050565b61ffff818116838216019080821115613e7c57613e7c613e2e565b5f7fff00000000000000000000000000000000000000000000000000000000000000808a60f81b1683527fffff0000000000000000000000000000000000000000000000000000000000008960f01b16600184015287516146c9816003860160208c01613926565b80840190507fffffffffffffffffffffffffffffffffffffffff0000000000000000000000008860601b166003820152865161470c816017840160208b01613926565b808201915050818660f81b16601782015284519150614732826018830160208801613926565b016018019998505050505050505050565b7fff000000000000000000000000000000000000000000000000000000000000008860f81b1681525f7fffff000000000000000000000000000000000000000000000000000000000000808960f01b16600184015287516147ab816003860160208c01613926565b80840190507fffffffffffffffffffffffffffffffffffffffff0000000000000000000000008860601b16600382015286516147ee816017840160208b01613926565b808201915050818660f01b16601782015284519150614814826019830160208801613926565b016019019998505050505050505050565b5f8651614836818460208b01613926565b9190910194855250602084019290925260f81b7fff000000000000000000000000000000000000000000000000000000000000009081166040840152166041820152604201919050565b818382375f9101908152919050565b85815273ffffffffffffffffffffffffffffffffffffffff8516602082015267ffffffffffffffff8416604082015260806060820152816080820152818360a08301375f81830160a090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016010194935050505056fea264697066735822122025d7a53cd17e05b58a25b9f2da1cce80bffc8d569bb7bbff01a94442980576d064736f6c63430008180033",
}

// FeijoapolygonzkevmABI is the input ABI used to generate the binding from.
// Deprecated: Use FeijoapolygonzkevmMetaData.ABI instead.
var FeijoapolygonzkevmABI = FeijoapolygonzkevmMetaData.ABI

// FeijoapolygonzkevmBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use FeijoapolygonzkevmMetaData.Bin instead.
var FeijoapolygonzkevmBin = FeijoapolygonzkevmMetaData.Bin

// DeployFeijoapolygonzkevm deploys a new Ethereum contract, binding an instance of Feijoapolygonzkevm to it.
func DeployFeijoapolygonzkevm(auth *bind.TransactOpts, backend bind.ContractBackend, _globalExitRootManager common.Address, _pol common.Address, _bridgeAddress common.Address, _rollupManager common.Address) (common.Address, *types.Transaction, *Feijoapolygonzkevm, error) {
	parsed, err := FeijoapolygonzkevmMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(FeijoapolygonzkevmBin), backend, _globalExitRootManager, _pol, _bridgeAddress, _rollupManager)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Feijoapolygonzkevm{FeijoapolygonzkevmCaller: FeijoapolygonzkevmCaller{contract: contract}, FeijoapolygonzkevmTransactor: FeijoapolygonzkevmTransactor{contract: contract}, FeijoapolygonzkevmFilterer: FeijoapolygonzkevmFilterer{contract: contract}}, nil
}

// Feijoapolygonzkevm is an auto generated Go binding around an Ethereum contract.
type Feijoapolygonzkevm struct {
	FeijoapolygonzkevmCaller     // Read-only binding to the contract
	FeijoapolygonzkevmTransactor // Write-only binding to the contract
	FeijoapolygonzkevmFilterer   // Log filterer for contract events
}

// FeijoapolygonzkevmCaller is an auto generated read-only Go binding around an Ethereum contract.
type FeijoapolygonzkevmCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FeijoapolygonzkevmTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FeijoapolygonzkevmTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FeijoapolygonzkevmFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FeijoapolygonzkevmFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FeijoapolygonzkevmSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FeijoapolygonzkevmSession struct {
	Contract     *Feijoapolygonzkevm // Generic contract binding to set the session for
	CallOpts     bind.CallOpts       // Call options to use throughout this session
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// FeijoapolygonzkevmCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FeijoapolygonzkevmCallerSession struct {
	Contract *FeijoapolygonzkevmCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts             // Call options to use throughout this session
}

// FeijoapolygonzkevmTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FeijoapolygonzkevmTransactorSession struct {
	Contract     *FeijoapolygonzkevmTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts             // Transaction auth options to use throughout this session
}

// FeijoapolygonzkevmRaw is an auto generated low-level Go binding around an Ethereum contract.
type FeijoapolygonzkevmRaw struct {
	Contract *Feijoapolygonzkevm // Generic contract binding to access the raw methods on
}

// FeijoapolygonzkevmCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FeijoapolygonzkevmCallerRaw struct {
	Contract *FeijoapolygonzkevmCaller // Generic read-only contract binding to access the raw methods on
}

// FeijoapolygonzkevmTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FeijoapolygonzkevmTransactorRaw struct {
	Contract *FeijoapolygonzkevmTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFeijoapolygonzkevm creates a new instance of Feijoapolygonzkevm, bound to a specific deployed contract.
func NewFeijoapolygonzkevm(address common.Address, backend bind.ContractBackend) (*Feijoapolygonzkevm, error) {
	contract, err := bindFeijoapolygonzkevm(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Feijoapolygonzkevm{FeijoapolygonzkevmCaller: FeijoapolygonzkevmCaller{contract: contract}, FeijoapolygonzkevmTransactor: FeijoapolygonzkevmTransactor{contract: contract}, FeijoapolygonzkevmFilterer: FeijoapolygonzkevmFilterer{contract: contract}}, nil
}

// NewFeijoapolygonzkevmCaller creates a new read-only instance of Feijoapolygonzkevm, bound to a specific deployed contract.
func NewFeijoapolygonzkevmCaller(address common.Address, caller bind.ContractCaller) (*FeijoapolygonzkevmCaller, error) {
	contract, err := bindFeijoapolygonzkevm(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FeijoapolygonzkevmCaller{contract: contract}, nil
}

// NewFeijoapolygonzkevmTransactor creates a new write-only instance of Feijoapolygonzkevm, bound to a specific deployed contract.
func NewFeijoapolygonzkevmTransactor(address common.Address, transactor bind.ContractTransactor) (*FeijoapolygonzkevmTransactor, error) {
	contract, err := bindFeijoapolygonzkevm(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FeijoapolygonzkevmTransactor{contract: contract}, nil
}

// NewFeijoapolygonzkevmFilterer creates a new log filterer instance of Feijoapolygonzkevm, bound to a specific deployed contract.
func NewFeijoapolygonzkevmFilterer(address common.Address, filterer bind.ContractFilterer) (*FeijoapolygonzkevmFilterer, error) {
	contract, err := bindFeijoapolygonzkevm(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FeijoapolygonzkevmFilterer{contract: contract}, nil
}

// bindFeijoapolygonzkevm binds a generic wrapper to an already deployed contract.
func bindFeijoapolygonzkevm(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FeijoapolygonzkevmMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Feijoapolygonzkevm *FeijoapolygonzkevmRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Feijoapolygonzkevm.Contract.FeijoapolygonzkevmCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Feijoapolygonzkevm *FeijoapolygonzkevmRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Feijoapolygonzkevm.Contract.FeijoapolygonzkevmTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Feijoapolygonzkevm *FeijoapolygonzkevmRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Feijoapolygonzkevm.Contract.FeijoapolygonzkevmTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Feijoapolygonzkevm.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Feijoapolygonzkevm *FeijoapolygonzkevmTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Feijoapolygonzkevm.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Feijoapolygonzkevm *FeijoapolygonzkevmTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Feijoapolygonzkevm.Contract.contract.Transact(opts, method, params...)
}

// GLOBALEXITROOTMANAGERL2 is a free data retrieval call binding the contract method 0x9e001877.
//
// Solidity: function GLOBAL_EXIT_ROOT_MANAGER_L2() view returns(address)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCaller) GLOBALEXITROOTMANAGERL2(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Feijoapolygonzkevm.contract.Call(opts, &out, "GLOBAL_EXIT_ROOT_MANAGER_L2")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GLOBALEXITROOTMANAGERL2 is a free data retrieval call binding the contract method 0x9e001877.
//
// Solidity: function GLOBAL_EXIT_ROOT_MANAGER_L2() view returns(address)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmSession) GLOBALEXITROOTMANAGERL2() (common.Address, error) {
	return _Feijoapolygonzkevm.Contract.GLOBALEXITROOTMANAGERL2(&_Feijoapolygonzkevm.CallOpts)
}

// GLOBALEXITROOTMANAGERL2 is a free data retrieval call binding the contract method 0x9e001877.
//
// Solidity: function GLOBAL_EXIT_ROOT_MANAGER_L2() view returns(address)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCallerSession) GLOBALEXITROOTMANAGERL2() (common.Address, error) {
	return _Feijoapolygonzkevm.Contract.GLOBALEXITROOTMANAGERL2(&_Feijoapolygonzkevm.CallOpts)
}

// INITIALIZETXBRIDGELISTLENLEN is a free data retrieval call binding the contract method 0x11e892d4.
//
// Solidity: function INITIALIZE_TX_BRIDGE_LIST_LEN_LEN() view returns(uint8)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCaller) INITIALIZETXBRIDGELISTLENLEN(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Feijoapolygonzkevm.contract.Call(opts, &out, "INITIALIZE_TX_BRIDGE_LIST_LEN_LEN")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// INITIALIZETXBRIDGELISTLENLEN is a free data retrieval call binding the contract method 0x11e892d4.
//
// Solidity: function INITIALIZE_TX_BRIDGE_LIST_LEN_LEN() view returns(uint8)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmSession) INITIALIZETXBRIDGELISTLENLEN() (uint8, error) {
	return _Feijoapolygonzkevm.Contract.INITIALIZETXBRIDGELISTLENLEN(&_Feijoapolygonzkevm.CallOpts)
}

// INITIALIZETXBRIDGELISTLENLEN is a free data retrieval call binding the contract method 0x11e892d4.
//
// Solidity: function INITIALIZE_TX_BRIDGE_LIST_LEN_LEN() view returns(uint8)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCallerSession) INITIALIZETXBRIDGELISTLENLEN() (uint8, error) {
	return _Feijoapolygonzkevm.Contract.INITIALIZETXBRIDGELISTLENLEN(&_Feijoapolygonzkevm.CallOpts)
}

// INITIALIZETXBRIDGEPARAMS is a free data retrieval call binding the contract method 0x05835f37.
//
// Solidity: function INITIALIZE_TX_BRIDGE_PARAMS() view returns(bytes)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCaller) INITIALIZETXBRIDGEPARAMS(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _Feijoapolygonzkevm.contract.Call(opts, &out, "INITIALIZE_TX_BRIDGE_PARAMS")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// INITIALIZETXBRIDGEPARAMS is a free data retrieval call binding the contract method 0x05835f37.
//
// Solidity: function INITIALIZE_TX_BRIDGE_PARAMS() view returns(bytes)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmSession) INITIALIZETXBRIDGEPARAMS() ([]byte, error) {
	return _Feijoapolygonzkevm.Contract.INITIALIZETXBRIDGEPARAMS(&_Feijoapolygonzkevm.CallOpts)
}

// INITIALIZETXBRIDGEPARAMS is a free data retrieval call binding the contract method 0x05835f37.
//
// Solidity: function INITIALIZE_TX_BRIDGE_PARAMS() view returns(bytes)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCallerSession) INITIALIZETXBRIDGEPARAMS() ([]byte, error) {
	return _Feijoapolygonzkevm.Contract.INITIALIZETXBRIDGEPARAMS(&_Feijoapolygonzkevm.CallOpts)
}

// INITIALIZETXBRIDGEPARAMSAFTERBRIDGEADDRESS is a free data retrieval call binding the contract method 0x7a5460c5.
//
// Solidity: function INITIALIZE_TX_BRIDGE_PARAMS_AFTER_BRIDGE_ADDRESS() view returns(bytes)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCaller) INITIALIZETXBRIDGEPARAMSAFTERBRIDGEADDRESS(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _Feijoapolygonzkevm.contract.Call(opts, &out, "INITIALIZE_TX_BRIDGE_PARAMS_AFTER_BRIDGE_ADDRESS")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// INITIALIZETXBRIDGEPARAMSAFTERBRIDGEADDRESS is a free data retrieval call binding the contract method 0x7a5460c5.
//
// Solidity: function INITIALIZE_TX_BRIDGE_PARAMS_AFTER_BRIDGE_ADDRESS() view returns(bytes)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmSession) INITIALIZETXBRIDGEPARAMSAFTERBRIDGEADDRESS() ([]byte, error) {
	return _Feijoapolygonzkevm.Contract.INITIALIZETXBRIDGEPARAMSAFTERBRIDGEADDRESS(&_Feijoapolygonzkevm.CallOpts)
}

// INITIALIZETXBRIDGEPARAMSAFTERBRIDGEADDRESS is a free data retrieval call binding the contract method 0x7a5460c5.
//
// Solidity: function INITIALIZE_TX_BRIDGE_PARAMS_AFTER_BRIDGE_ADDRESS() view returns(bytes)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCallerSession) INITIALIZETXBRIDGEPARAMSAFTERBRIDGEADDRESS() ([]byte, error) {
	return _Feijoapolygonzkevm.Contract.INITIALIZETXBRIDGEPARAMSAFTERBRIDGEADDRESS(&_Feijoapolygonzkevm.CallOpts)
}

// INITIALIZETXBRIDGEPARAMSAFTERBRIDGEADDRESSEMPTYMETADATA is a free data retrieval call binding the contract method 0x52bdeb6d.
//
// Solidity: function INITIALIZE_TX_BRIDGE_PARAMS_AFTER_BRIDGE_ADDRESS_EMPTY_METADATA() view returns(bytes)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCaller) INITIALIZETXBRIDGEPARAMSAFTERBRIDGEADDRESSEMPTYMETADATA(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _Feijoapolygonzkevm.contract.Call(opts, &out, "INITIALIZE_TX_BRIDGE_PARAMS_AFTER_BRIDGE_ADDRESS_EMPTY_METADATA")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// INITIALIZETXBRIDGEPARAMSAFTERBRIDGEADDRESSEMPTYMETADATA is a free data retrieval call binding the contract method 0x52bdeb6d.
//
// Solidity: function INITIALIZE_TX_BRIDGE_PARAMS_AFTER_BRIDGE_ADDRESS_EMPTY_METADATA() view returns(bytes)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmSession) INITIALIZETXBRIDGEPARAMSAFTERBRIDGEADDRESSEMPTYMETADATA() ([]byte, error) {
	return _Feijoapolygonzkevm.Contract.INITIALIZETXBRIDGEPARAMSAFTERBRIDGEADDRESSEMPTYMETADATA(&_Feijoapolygonzkevm.CallOpts)
}

// INITIALIZETXBRIDGEPARAMSAFTERBRIDGEADDRESSEMPTYMETADATA is a free data retrieval call binding the contract method 0x52bdeb6d.
//
// Solidity: function INITIALIZE_TX_BRIDGE_PARAMS_AFTER_BRIDGE_ADDRESS_EMPTY_METADATA() view returns(bytes)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCallerSession) INITIALIZETXBRIDGEPARAMSAFTERBRIDGEADDRESSEMPTYMETADATA() ([]byte, error) {
	return _Feijoapolygonzkevm.Contract.INITIALIZETXBRIDGEPARAMSAFTERBRIDGEADDRESSEMPTYMETADATA(&_Feijoapolygonzkevm.CallOpts)
}

// INITIALIZETXCONSTANTBYTES is a free data retrieval call binding the contract method 0x03508963.
//
// Solidity: function INITIALIZE_TX_CONSTANT_BYTES() view returns(uint16)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCaller) INITIALIZETXCONSTANTBYTES(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _Feijoapolygonzkevm.contract.Call(opts, &out, "INITIALIZE_TX_CONSTANT_BYTES")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// INITIALIZETXCONSTANTBYTES is a free data retrieval call binding the contract method 0x03508963.
//
// Solidity: function INITIALIZE_TX_CONSTANT_BYTES() view returns(uint16)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmSession) INITIALIZETXCONSTANTBYTES() (uint16, error) {
	return _Feijoapolygonzkevm.Contract.INITIALIZETXCONSTANTBYTES(&_Feijoapolygonzkevm.CallOpts)
}

// INITIALIZETXCONSTANTBYTES is a free data retrieval call binding the contract method 0x03508963.
//
// Solidity: function INITIALIZE_TX_CONSTANT_BYTES() view returns(uint16)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCallerSession) INITIALIZETXCONSTANTBYTES() (uint16, error) {
	return _Feijoapolygonzkevm.Contract.INITIALIZETXCONSTANTBYTES(&_Feijoapolygonzkevm.CallOpts)
}

// INITIALIZETXCONSTANTBYTESEMPTYMETADATA is a free data retrieval call binding the contract method 0x676870d2.
//
// Solidity: function INITIALIZE_TX_CONSTANT_BYTES_EMPTY_METADATA() view returns(uint16)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCaller) INITIALIZETXCONSTANTBYTESEMPTYMETADATA(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _Feijoapolygonzkevm.contract.Call(opts, &out, "INITIALIZE_TX_CONSTANT_BYTES_EMPTY_METADATA")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// INITIALIZETXCONSTANTBYTESEMPTYMETADATA is a free data retrieval call binding the contract method 0x676870d2.
//
// Solidity: function INITIALIZE_TX_CONSTANT_BYTES_EMPTY_METADATA() view returns(uint16)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmSession) INITIALIZETXCONSTANTBYTESEMPTYMETADATA() (uint16, error) {
	return _Feijoapolygonzkevm.Contract.INITIALIZETXCONSTANTBYTESEMPTYMETADATA(&_Feijoapolygonzkevm.CallOpts)
}

// INITIALIZETXCONSTANTBYTESEMPTYMETADATA is a free data retrieval call binding the contract method 0x676870d2.
//
// Solidity: function INITIALIZE_TX_CONSTANT_BYTES_EMPTY_METADATA() view returns(uint16)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCallerSession) INITIALIZETXCONSTANTBYTESEMPTYMETADATA() (uint16, error) {
	return _Feijoapolygonzkevm.Contract.INITIALIZETXCONSTANTBYTESEMPTYMETADATA(&_Feijoapolygonzkevm.CallOpts)
}

// INITIALIZETXDATALENEMPTYMETADATA is a free data retrieval call binding the contract method 0xc7fffd4b.
//
// Solidity: function INITIALIZE_TX_DATA_LEN_EMPTY_METADATA() view returns(uint8)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCaller) INITIALIZETXDATALENEMPTYMETADATA(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Feijoapolygonzkevm.contract.Call(opts, &out, "INITIALIZE_TX_DATA_LEN_EMPTY_METADATA")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// INITIALIZETXDATALENEMPTYMETADATA is a free data retrieval call binding the contract method 0xc7fffd4b.
//
// Solidity: function INITIALIZE_TX_DATA_LEN_EMPTY_METADATA() view returns(uint8)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmSession) INITIALIZETXDATALENEMPTYMETADATA() (uint8, error) {
	return _Feijoapolygonzkevm.Contract.INITIALIZETXDATALENEMPTYMETADATA(&_Feijoapolygonzkevm.CallOpts)
}

// INITIALIZETXDATALENEMPTYMETADATA is a free data retrieval call binding the contract method 0xc7fffd4b.
//
// Solidity: function INITIALIZE_TX_DATA_LEN_EMPTY_METADATA() view returns(uint8)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCallerSession) INITIALIZETXDATALENEMPTYMETADATA() (uint8, error) {
	return _Feijoapolygonzkevm.Contract.INITIALIZETXDATALENEMPTYMETADATA(&_Feijoapolygonzkevm.CallOpts)
}

// INITIALIZETXEFFECTIVEPERCENTAGE is a free data retrieval call binding the contract method 0x40b5de6c.
//
// Solidity: function INITIALIZE_TX_EFFECTIVE_PERCENTAGE() view returns(bytes1)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCaller) INITIALIZETXEFFECTIVEPERCENTAGE(opts *bind.CallOpts) ([1]byte, error) {
	var out []interface{}
	err := _Feijoapolygonzkevm.contract.Call(opts, &out, "INITIALIZE_TX_EFFECTIVE_PERCENTAGE")

	if err != nil {
		return *new([1]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([1]byte)).(*[1]byte)

	return out0, err

}

// INITIALIZETXEFFECTIVEPERCENTAGE is a free data retrieval call binding the contract method 0x40b5de6c.
//
// Solidity: function INITIALIZE_TX_EFFECTIVE_PERCENTAGE() view returns(bytes1)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmSession) INITIALIZETXEFFECTIVEPERCENTAGE() ([1]byte, error) {
	return _Feijoapolygonzkevm.Contract.INITIALIZETXEFFECTIVEPERCENTAGE(&_Feijoapolygonzkevm.CallOpts)
}

// INITIALIZETXEFFECTIVEPERCENTAGE is a free data retrieval call binding the contract method 0x40b5de6c.
//
// Solidity: function INITIALIZE_TX_EFFECTIVE_PERCENTAGE() view returns(bytes1)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCallerSession) INITIALIZETXEFFECTIVEPERCENTAGE() ([1]byte, error) {
	return _Feijoapolygonzkevm.Contract.INITIALIZETXEFFECTIVEPERCENTAGE(&_Feijoapolygonzkevm.CallOpts)
}

// MAXSEQUENCETIMESTAMPFORCED is a free data retrieval call binding the contract method 0x7160c5f7.
//
// Solidity: function MAX_SEQUENCE_TIMESTAMP_FORCED() view returns(uint64)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCaller) MAXSEQUENCETIMESTAMPFORCED(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Feijoapolygonzkevm.contract.Call(opts, &out, "MAX_SEQUENCE_TIMESTAMP_FORCED")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// MAXSEQUENCETIMESTAMPFORCED is a free data retrieval call binding the contract method 0x7160c5f7.
//
// Solidity: function MAX_SEQUENCE_TIMESTAMP_FORCED() view returns(uint64)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmSession) MAXSEQUENCETIMESTAMPFORCED() (uint64, error) {
	return _Feijoapolygonzkevm.Contract.MAXSEQUENCETIMESTAMPFORCED(&_Feijoapolygonzkevm.CallOpts)
}

// MAXSEQUENCETIMESTAMPFORCED is a free data retrieval call binding the contract method 0x7160c5f7.
//
// Solidity: function MAX_SEQUENCE_TIMESTAMP_FORCED() view returns(uint64)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCallerSession) MAXSEQUENCETIMESTAMPFORCED() (uint64, error) {
	return _Feijoapolygonzkevm.Contract.MAXSEQUENCETIMESTAMPFORCED(&_Feijoapolygonzkevm.CallOpts)
}

// POINTEVALUATIONPRECOMPILEADDRESS is a free data retrieval call binding the contract method 0x3e41062e.
//
// Solidity: function POINT_EVALUATION_PRECOMPILE_ADDRESS() view returns(address)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCaller) POINTEVALUATIONPRECOMPILEADDRESS(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Feijoapolygonzkevm.contract.Call(opts, &out, "POINT_EVALUATION_PRECOMPILE_ADDRESS")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// POINTEVALUATIONPRECOMPILEADDRESS is a free data retrieval call binding the contract method 0x3e41062e.
//
// Solidity: function POINT_EVALUATION_PRECOMPILE_ADDRESS() view returns(address)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmSession) POINTEVALUATIONPRECOMPILEADDRESS() (common.Address, error) {
	return _Feijoapolygonzkevm.Contract.POINTEVALUATIONPRECOMPILEADDRESS(&_Feijoapolygonzkevm.CallOpts)
}

// POINTEVALUATIONPRECOMPILEADDRESS is a free data retrieval call binding the contract method 0x3e41062e.
//
// Solidity: function POINT_EVALUATION_PRECOMPILE_ADDRESS() view returns(address)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCallerSession) POINTEVALUATIONPRECOMPILEADDRESS() (common.Address, error) {
	return _Feijoapolygonzkevm.Contract.POINTEVALUATIONPRECOMPILEADDRESS(&_Feijoapolygonzkevm.CallOpts)
}

// SIGNATUREINITIALIZETXR is a free data retrieval call binding the contract method 0xb0afe154.
//
// Solidity: function SIGNATURE_INITIALIZE_TX_R() view returns(bytes32)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCaller) SIGNATUREINITIALIZETXR(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Feijoapolygonzkevm.contract.Call(opts, &out, "SIGNATURE_INITIALIZE_TX_R")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// SIGNATUREINITIALIZETXR is a free data retrieval call binding the contract method 0xb0afe154.
//
// Solidity: function SIGNATURE_INITIALIZE_TX_R() view returns(bytes32)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmSession) SIGNATUREINITIALIZETXR() ([32]byte, error) {
	return _Feijoapolygonzkevm.Contract.SIGNATUREINITIALIZETXR(&_Feijoapolygonzkevm.CallOpts)
}

// SIGNATUREINITIALIZETXR is a free data retrieval call binding the contract method 0xb0afe154.
//
// Solidity: function SIGNATURE_INITIALIZE_TX_R() view returns(bytes32)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCallerSession) SIGNATUREINITIALIZETXR() ([32]byte, error) {
	return _Feijoapolygonzkevm.Contract.SIGNATUREINITIALIZETXR(&_Feijoapolygonzkevm.CallOpts)
}

// SIGNATUREINITIALIZETXS is a free data retrieval call binding the contract method 0xd7bc90ff.
//
// Solidity: function SIGNATURE_INITIALIZE_TX_S() view returns(bytes32)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCaller) SIGNATUREINITIALIZETXS(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Feijoapolygonzkevm.contract.Call(opts, &out, "SIGNATURE_INITIALIZE_TX_S")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// SIGNATUREINITIALIZETXS is a free data retrieval call binding the contract method 0xd7bc90ff.
//
// Solidity: function SIGNATURE_INITIALIZE_TX_S() view returns(bytes32)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmSession) SIGNATUREINITIALIZETXS() ([32]byte, error) {
	return _Feijoapolygonzkevm.Contract.SIGNATUREINITIALIZETXS(&_Feijoapolygonzkevm.CallOpts)
}

// SIGNATUREINITIALIZETXS is a free data retrieval call binding the contract method 0xd7bc90ff.
//
// Solidity: function SIGNATURE_INITIALIZE_TX_S() view returns(bytes32)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCallerSession) SIGNATUREINITIALIZETXS() ([32]byte, error) {
	return _Feijoapolygonzkevm.Contract.SIGNATUREINITIALIZETXS(&_Feijoapolygonzkevm.CallOpts)
}

// SIGNATUREINITIALIZETXV is a free data retrieval call binding the contract method 0xf35dda47.
//
// Solidity: function SIGNATURE_INITIALIZE_TX_V() view returns(uint8)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCaller) SIGNATUREINITIALIZETXV(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Feijoapolygonzkevm.contract.Call(opts, &out, "SIGNATURE_INITIALIZE_TX_V")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// SIGNATUREINITIALIZETXV is a free data retrieval call binding the contract method 0xf35dda47.
//
// Solidity: function SIGNATURE_INITIALIZE_TX_V() view returns(uint8)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmSession) SIGNATUREINITIALIZETXV() (uint8, error) {
	return _Feijoapolygonzkevm.Contract.SIGNATUREINITIALIZETXV(&_Feijoapolygonzkevm.CallOpts)
}

// SIGNATUREINITIALIZETXV is a free data retrieval call binding the contract method 0xf35dda47.
//
// Solidity: function SIGNATURE_INITIALIZE_TX_V() view returns(uint8)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCallerSession) SIGNATUREINITIALIZETXV() (uint8, error) {
	return _Feijoapolygonzkevm.Contract.SIGNATUREINITIALIZETXV(&_Feijoapolygonzkevm.CallOpts)
}

// TIMESTAMPRANGE is a free data retrieval call binding the contract method 0x42308fab.
//
// Solidity: function TIMESTAMP_RANGE() view returns(uint256)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCaller) TIMESTAMPRANGE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Feijoapolygonzkevm.contract.Call(opts, &out, "TIMESTAMP_RANGE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TIMESTAMPRANGE is a free data retrieval call binding the contract method 0x42308fab.
//
// Solidity: function TIMESTAMP_RANGE() view returns(uint256)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmSession) TIMESTAMPRANGE() (*big.Int, error) {
	return _Feijoapolygonzkevm.Contract.TIMESTAMPRANGE(&_Feijoapolygonzkevm.CallOpts)
}

// TIMESTAMPRANGE is a free data retrieval call binding the contract method 0x42308fab.
//
// Solidity: function TIMESTAMP_RANGE() view returns(uint256)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCallerSession) TIMESTAMPRANGE() (*big.Int, error) {
	return _Feijoapolygonzkevm.Contract.TIMESTAMPRANGE(&_Feijoapolygonzkevm.CallOpts)
}

// ZKGASLIMITBATCH is a free data retrieval call binding the contract method 0x838a2503.
//
// Solidity: function ZK_GAS_LIMIT_BATCH() view returns(uint64)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCaller) ZKGASLIMITBATCH(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Feijoapolygonzkevm.contract.Call(opts, &out, "ZK_GAS_LIMIT_BATCH")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// ZKGASLIMITBATCH is a free data retrieval call binding the contract method 0x838a2503.
//
// Solidity: function ZK_GAS_LIMIT_BATCH() view returns(uint64)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmSession) ZKGASLIMITBATCH() (uint64, error) {
	return _Feijoapolygonzkevm.Contract.ZKGASLIMITBATCH(&_Feijoapolygonzkevm.CallOpts)
}

// ZKGASLIMITBATCH is a free data retrieval call binding the contract method 0x838a2503.
//
// Solidity: function ZK_GAS_LIMIT_BATCH() view returns(uint64)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCallerSession) ZKGASLIMITBATCH() (uint64, error) {
	return _Feijoapolygonzkevm.Contract.ZKGASLIMITBATCH(&_Feijoapolygonzkevm.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCaller) Admin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Feijoapolygonzkevm.contract.Call(opts, &out, "admin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmSession) Admin() (common.Address, error) {
	return _Feijoapolygonzkevm.Contract.Admin(&_Feijoapolygonzkevm.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCallerSession) Admin() (common.Address, error) {
	return _Feijoapolygonzkevm.Contract.Admin(&_Feijoapolygonzkevm.CallOpts)
}

// BridgeAddress is a free data retrieval call binding the contract method 0xa3c573eb.
//
// Solidity: function bridgeAddress() view returns(address)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCaller) BridgeAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Feijoapolygonzkevm.contract.Call(opts, &out, "bridgeAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BridgeAddress is a free data retrieval call binding the contract method 0xa3c573eb.
//
// Solidity: function bridgeAddress() view returns(address)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmSession) BridgeAddress() (common.Address, error) {
	return _Feijoapolygonzkevm.Contract.BridgeAddress(&_Feijoapolygonzkevm.CallOpts)
}

// BridgeAddress is a free data retrieval call binding the contract method 0xa3c573eb.
//
// Solidity: function bridgeAddress() view returns(address)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCallerSession) BridgeAddress() (common.Address, error) {
	return _Feijoapolygonzkevm.Contract.BridgeAddress(&_Feijoapolygonzkevm.CallOpts)
}

// CalculatePolPerForcedZkGas is a free data retrieval call binding the contract method 0x042b0f06.
//
// Solidity: function calculatePolPerForcedZkGas() view returns(uint256)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCaller) CalculatePolPerForcedZkGas(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Feijoapolygonzkevm.contract.Call(opts, &out, "calculatePolPerForcedZkGas")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CalculatePolPerForcedZkGas is a free data retrieval call binding the contract method 0x042b0f06.
//
// Solidity: function calculatePolPerForcedZkGas() view returns(uint256)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmSession) CalculatePolPerForcedZkGas() (*big.Int, error) {
	return _Feijoapolygonzkevm.Contract.CalculatePolPerForcedZkGas(&_Feijoapolygonzkevm.CallOpts)
}

// CalculatePolPerForcedZkGas is a free data retrieval call binding the contract method 0x042b0f06.
//
// Solidity: function calculatePolPerForcedZkGas() view returns(uint256)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCallerSession) CalculatePolPerForcedZkGas() (*big.Int, error) {
	return _Feijoapolygonzkevm.Contract.CalculatePolPerForcedZkGas(&_Feijoapolygonzkevm.CallOpts)
}

// ForceBlobAddress is a free data retrieval call binding the contract method 0x66e7bb1a.
//
// Solidity: function forceBlobAddress() view returns(address)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCaller) ForceBlobAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Feijoapolygonzkevm.contract.Call(opts, &out, "forceBlobAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ForceBlobAddress is a free data retrieval call binding the contract method 0x66e7bb1a.
//
// Solidity: function forceBlobAddress() view returns(address)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmSession) ForceBlobAddress() (common.Address, error) {
	return _Feijoapolygonzkevm.Contract.ForceBlobAddress(&_Feijoapolygonzkevm.CallOpts)
}

// ForceBlobAddress is a free data retrieval call binding the contract method 0x66e7bb1a.
//
// Solidity: function forceBlobAddress() view returns(address)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCallerSession) ForceBlobAddress() (common.Address, error) {
	return _Feijoapolygonzkevm.Contract.ForceBlobAddress(&_Feijoapolygonzkevm.CallOpts)
}

// ForceBlobTimeout is a free data retrieval call binding the contract method 0x2c2251db.
//
// Solidity: function forceBlobTimeout() view returns(uint64)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCaller) ForceBlobTimeout(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Feijoapolygonzkevm.contract.Call(opts, &out, "forceBlobTimeout")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// ForceBlobTimeout is a free data retrieval call binding the contract method 0x2c2251db.
//
// Solidity: function forceBlobTimeout() view returns(uint64)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmSession) ForceBlobTimeout() (uint64, error) {
	return _Feijoapolygonzkevm.Contract.ForceBlobTimeout(&_Feijoapolygonzkevm.CallOpts)
}

// ForceBlobTimeout is a free data retrieval call binding the contract method 0x2c2251db.
//
// Solidity: function forceBlobTimeout() view returns(uint64)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCallerSession) ForceBlobTimeout() (uint64, error) {
	return _Feijoapolygonzkevm.Contract.ForceBlobTimeout(&_Feijoapolygonzkevm.CallOpts)
}

// ForcedBlobs is a free data retrieval call binding the contract method 0x2a6688ee.
//
// Solidity: function forcedBlobs(uint64 ) view returns(bytes32 hashedForcedBlobData, uint64 forcedTimestamp)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCaller) ForcedBlobs(opts *bind.CallOpts, arg0 uint64) (struct {
	HashedForcedBlobData [32]byte
	ForcedTimestamp      uint64
}, error) {
	var out []interface{}
	err := _Feijoapolygonzkevm.contract.Call(opts, &out, "forcedBlobs", arg0)

	outstruct := new(struct {
		HashedForcedBlobData [32]byte
		ForcedTimestamp      uint64
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.HashedForcedBlobData = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.ForcedTimestamp = *abi.ConvertType(out[1], new(uint64)).(*uint64)

	return *outstruct, err

}

// ForcedBlobs is a free data retrieval call binding the contract method 0x2a6688ee.
//
// Solidity: function forcedBlobs(uint64 ) view returns(bytes32 hashedForcedBlobData, uint64 forcedTimestamp)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmSession) ForcedBlobs(arg0 uint64) (struct {
	HashedForcedBlobData [32]byte
	ForcedTimestamp      uint64
}, error) {
	return _Feijoapolygonzkevm.Contract.ForcedBlobs(&_Feijoapolygonzkevm.CallOpts, arg0)
}

// ForcedBlobs is a free data retrieval call binding the contract method 0x2a6688ee.
//
// Solidity: function forcedBlobs(uint64 ) view returns(bytes32 hashedForcedBlobData, uint64 forcedTimestamp)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCallerSession) ForcedBlobs(arg0 uint64) (struct {
	HashedForcedBlobData [32]byte
	ForcedTimestamp      uint64
}, error) {
	return _Feijoapolygonzkevm.Contract.ForcedBlobs(&_Feijoapolygonzkevm.CallOpts, arg0)
}

// GasTokenAddress is a free data retrieval call binding the contract method 0x3c351e10.
//
// Solidity: function gasTokenAddress() view returns(address)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCaller) GasTokenAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Feijoapolygonzkevm.contract.Call(opts, &out, "gasTokenAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GasTokenAddress is a free data retrieval call binding the contract method 0x3c351e10.
//
// Solidity: function gasTokenAddress() view returns(address)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmSession) GasTokenAddress() (common.Address, error) {
	return _Feijoapolygonzkevm.Contract.GasTokenAddress(&_Feijoapolygonzkevm.CallOpts)
}

// GasTokenAddress is a free data retrieval call binding the contract method 0x3c351e10.
//
// Solidity: function gasTokenAddress() view returns(address)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCallerSession) GasTokenAddress() (common.Address, error) {
	return _Feijoapolygonzkevm.Contract.GasTokenAddress(&_Feijoapolygonzkevm.CallOpts)
}

// GasTokenNetwork is a free data retrieval call binding the contract method 0x3cbc795b.
//
// Solidity: function gasTokenNetwork() view returns(uint32)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCaller) GasTokenNetwork(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _Feijoapolygonzkevm.contract.Call(opts, &out, "gasTokenNetwork")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// GasTokenNetwork is a free data retrieval call binding the contract method 0x3cbc795b.
//
// Solidity: function gasTokenNetwork() view returns(uint32)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmSession) GasTokenNetwork() (uint32, error) {
	return _Feijoapolygonzkevm.Contract.GasTokenNetwork(&_Feijoapolygonzkevm.CallOpts)
}

// GasTokenNetwork is a free data retrieval call binding the contract method 0x3cbc795b.
//
// Solidity: function gasTokenNetwork() view returns(uint32)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCallerSession) GasTokenNetwork() (uint32, error) {
	return _Feijoapolygonzkevm.Contract.GasTokenNetwork(&_Feijoapolygonzkevm.CallOpts)
}

// GenerateInitializeTransaction is a free data retrieval call binding the contract method 0xa652f26c.
//
// Solidity: function generateInitializeTransaction(uint32 networkID, address _gasTokenAddress, uint32 _gasTokenNetwork, bytes _gasTokenMetadata) view returns(bytes)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCaller) GenerateInitializeTransaction(opts *bind.CallOpts, networkID uint32, _gasTokenAddress common.Address, _gasTokenNetwork uint32, _gasTokenMetadata []byte) ([]byte, error) {
	var out []interface{}
	err := _Feijoapolygonzkevm.contract.Call(opts, &out, "generateInitializeTransaction", networkID, _gasTokenAddress, _gasTokenNetwork, _gasTokenMetadata)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GenerateInitializeTransaction is a free data retrieval call binding the contract method 0xa652f26c.
//
// Solidity: function generateInitializeTransaction(uint32 networkID, address _gasTokenAddress, uint32 _gasTokenNetwork, bytes _gasTokenMetadata) view returns(bytes)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmSession) GenerateInitializeTransaction(networkID uint32, _gasTokenAddress common.Address, _gasTokenNetwork uint32, _gasTokenMetadata []byte) ([]byte, error) {
	return _Feijoapolygonzkevm.Contract.GenerateInitializeTransaction(&_Feijoapolygonzkevm.CallOpts, networkID, _gasTokenAddress, _gasTokenNetwork, _gasTokenMetadata)
}

// GenerateInitializeTransaction is a free data retrieval call binding the contract method 0xa652f26c.
//
// Solidity: function generateInitializeTransaction(uint32 networkID, address _gasTokenAddress, uint32 _gasTokenNetwork, bytes _gasTokenMetadata) view returns(bytes)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCallerSession) GenerateInitializeTransaction(networkID uint32, _gasTokenAddress common.Address, _gasTokenNetwork uint32, _gasTokenMetadata []byte) ([]byte, error) {
	return _Feijoapolygonzkevm.Contract.GenerateInitializeTransaction(&_Feijoapolygonzkevm.CallOpts, networkID, _gasTokenAddress, _gasTokenNetwork, _gasTokenMetadata)
}

// GlobalExitRootManager is a free data retrieval call binding the contract method 0xd02103ca.
//
// Solidity: function globalExitRootManager() view returns(address)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCaller) GlobalExitRootManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Feijoapolygonzkevm.contract.Call(opts, &out, "globalExitRootManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GlobalExitRootManager is a free data retrieval call binding the contract method 0xd02103ca.
//
// Solidity: function globalExitRootManager() view returns(address)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmSession) GlobalExitRootManager() (common.Address, error) {
	return _Feijoapolygonzkevm.Contract.GlobalExitRootManager(&_Feijoapolygonzkevm.CallOpts)
}

// GlobalExitRootManager is a free data retrieval call binding the contract method 0xd02103ca.
//
// Solidity: function globalExitRootManager() view returns(address)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCallerSession) GlobalExitRootManager() (common.Address, error) {
	return _Feijoapolygonzkevm.Contract.GlobalExitRootManager(&_Feijoapolygonzkevm.CallOpts)
}

// LastAccInputHash is a free data retrieval call binding the contract method 0x6e05d2cd.
//
// Solidity: function lastAccInputHash() view returns(bytes32)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCaller) LastAccInputHash(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Feijoapolygonzkevm.contract.Call(opts, &out, "lastAccInputHash")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// LastAccInputHash is a free data retrieval call binding the contract method 0x6e05d2cd.
//
// Solidity: function lastAccInputHash() view returns(bytes32)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmSession) LastAccInputHash() ([32]byte, error) {
	return _Feijoapolygonzkevm.Contract.LastAccInputHash(&_Feijoapolygonzkevm.CallOpts)
}

// LastAccInputHash is a free data retrieval call binding the contract method 0x6e05d2cd.
//
// Solidity: function lastAccInputHash() view returns(bytes32)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCallerSession) LastAccInputHash() ([32]byte, error) {
	return _Feijoapolygonzkevm.Contract.LastAccInputHash(&_Feijoapolygonzkevm.CallOpts)
}

// LastForceBlob is a free data retrieval call binding the contract method 0x9b0e35a5.
//
// Solidity: function lastForceBlob() view returns(uint64)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCaller) LastForceBlob(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Feijoapolygonzkevm.contract.Call(opts, &out, "lastForceBlob")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// LastForceBlob is a free data retrieval call binding the contract method 0x9b0e35a5.
//
// Solidity: function lastForceBlob() view returns(uint64)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmSession) LastForceBlob() (uint64, error) {
	return _Feijoapolygonzkevm.Contract.LastForceBlob(&_Feijoapolygonzkevm.CallOpts)
}

// LastForceBlob is a free data retrieval call binding the contract method 0x9b0e35a5.
//
// Solidity: function lastForceBlob() view returns(uint64)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCallerSession) LastForceBlob() (uint64, error) {
	return _Feijoapolygonzkevm.Contract.LastForceBlob(&_Feijoapolygonzkevm.CallOpts)
}

// LastForceBlobSequenced is a free data retrieval call binding the contract method 0xb45bd7f9.
//
// Solidity: function lastForceBlobSequenced() view returns(uint64)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCaller) LastForceBlobSequenced(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Feijoapolygonzkevm.contract.Call(opts, &out, "lastForceBlobSequenced")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// LastForceBlobSequenced is a free data retrieval call binding the contract method 0xb45bd7f9.
//
// Solidity: function lastForceBlobSequenced() view returns(uint64)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmSession) LastForceBlobSequenced() (uint64, error) {
	return _Feijoapolygonzkevm.Contract.LastForceBlobSequenced(&_Feijoapolygonzkevm.CallOpts)
}

// LastForceBlobSequenced is a free data retrieval call binding the contract method 0xb45bd7f9.
//
// Solidity: function lastForceBlobSequenced() view returns(uint64)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCallerSession) LastForceBlobSequenced() (uint64, error) {
	return _Feijoapolygonzkevm.Contract.LastForceBlobSequenced(&_Feijoapolygonzkevm.CallOpts)
}

// NetworkName is a free data retrieval call binding the contract method 0x107bf28c.
//
// Solidity: function networkName() view returns(string)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCaller) NetworkName(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Feijoapolygonzkevm.contract.Call(opts, &out, "networkName")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// NetworkName is a free data retrieval call binding the contract method 0x107bf28c.
//
// Solidity: function networkName() view returns(string)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmSession) NetworkName() (string, error) {
	return _Feijoapolygonzkevm.Contract.NetworkName(&_Feijoapolygonzkevm.CallOpts)
}

// NetworkName is a free data retrieval call binding the contract method 0x107bf28c.
//
// Solidity: function networkName() view returns(string)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCallerSession) NetworkName() (string, error) {
	return _Feijoapolygonzkevm.Contract.NetworkName(&_Feijoapolygonzkevm.CallOpts)
}

// PendingAdmin is a free data retrieval call binding the contract method 0x26782247.
//
// Solidity: function pendingAdmin() view returns(address)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCaller) PendingAdmin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Feijoapolygonzkevm.contract.Call(opts, &out, "pendingAdmin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingAdmin is a free data retrieval call binding the contract method 0x26782247.
//
// Solidity: function pendingAdmin() view returns(address)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmSession) PendingAdmin() (common.Address, error) {
	return _Feijoapolygonzkevm.Contract.PendingAdmin(&_Feijoapolygonzkevm.CallOpts)
}

// PendingAdmin is a free data retrieval call binding the contract method 0x26782247.
//
// Solidity: function pendingAdmin() view returns(address)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCallerSession) PendingAdmin() (common.Address, error) {
	return _Feijoapolygonzkevm.Contract.PendingAdmin(&_Feijoapolygonzkevm.CallOpts)
}

// Pol is a free data retrieval call binding the contract method 0xe46761c4.
//
// Solidity: function pol() view returns(address)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCaller) Pol(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Feijoapolygonzkevm.contract.Call(opts, &out, "pol")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Pol is a free data retrieval call binding the contract method 0xe46761c4.
//
// Solidity: function pol() view returns(address)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmSession) Pol() (common.Address, error) {
	return _Feijoapolygonzkevm.Contract.Pol(&_Feijoapolygonzkevm.CallOpts)
}

// Pol is a free data retrieval call binding the contract method 0xe46761c4.
//
// Solidity: function pol() view returns(address)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCallerSession) Pol() (common.Address, error) {
	return _Feijoapolygonzkevm.Contract.Pol(&_Feijoapolygonzkevm.CallOpts)
}

// RollupManager is a free data retrieval call binding the contract method 0x49b7b802.
//
// Solidity: function rollupManager() view returns(address)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCaller) RollupManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Feijoapolygonzkevm.contract.Call(opts, &out, "rollupManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RollupManager is a free data retrieval call binding the contract method 0x49b7b802.
//
// Solidity: function rollupManager() view returns(address)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmSession) RollupManager() (common.Address, error) {
	return _Feijoapolygonzkevm.Contract.RollupManager(&_Feijoapolygonzkevm.CallOpts)
}

// RollupManager is a free data retrieval call binding the contract method 0x49b7b802.
//
// Solidity: function rollupManager() view returns(address)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCallerSession) RollupManager() (common.Address, error) {
	return _Feijoapolygonzkevm.Contract.RollupManager(&_Feijoapolygonzkevm.CallOpts)
}

// TrustedSequencer is a free data retrieval call binding the contract method 0xcfa8ed47.
//
// Solidity: function trustedSequencer() view returns(address)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCaller) TrustedSequencer(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Feijoapolygonzkevm.contract.Call(opts, &out, "trustedSequencer")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TrustedSequencer is a free data retrieval call binding the contract method 0xcfa8ed47.
//
// Solidity: function trustedSequencer() view returns(address)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmSession) TrustedSequencer() (common.Address, error) {
	return _Feijoapolygonzkevm.Contract.TrustedSequencer(&_Feijoapolygonzkevm.CallOpts)
}

// TrustedSequencer is a free data retrieval call binding the contract method 0xcfa8ed47.
//
// Solidity: function trustedSequencer() view returns(address)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCallerSession) TrustedSequencer() (common.Address, error) {
	return _Feijoapolygonzkevm.Contract.TrustedSequencer(&_Feijoapolygonzkevm.CallOpts)
}

// TrustedSequencerURL is a free data retrieval call binding the contract method 0x542028d5.
//
// Solidity: function trustedSequencerURL() view returns(string)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCaller) TrustedSequencerURL(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Feijoapolygonzkevm.contract.Call(opts, &out, "trustedSequencerURL")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TrustedSequencerURL is a free data retrieval call binding the contract method 0x542028d5.
//
// Solidity: function trustedSequencerURL() view returns(string)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmSession) TrustedSequencerURL() (string, error) {
	return _Feijoapolygonzkevm.Contract.TrustedSequencerURL(&_Feijoapolygonzkevm.CallOpts)
}

// TrustedSequencerURL is a free data retrieval call binding the contract method 0x542028d5.
//
// Solidity: function trustedSequencerURL() view returns(string)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmCallerSession) TrustedSequencerURL() (string, error) {
	return _Feijoapolygonzkevm.Contract.TrustedSequencerURL(&_Feijoapolygonzkevm.CallOpts)
}

// AcceptAdminRole is a paid mutator transaction binding the contract method 0x8c3d7301.
//
// Solidity: function acceptAdminRole() returns()
func (_Feijoapolygonzkevm *FeijoapolygonzkevmTransactor) AcceptAdminRole(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Feijoapolygonzkevm.contract.Transact(opts, "acceptAdminRole")
}

// AcceptAdminRole is a paid mutator transaction binding the contract method 0x8c3d7301.
//
// Solidity: function acceptAdminRole() returns()
func (_Feijoapolygonzkevm *FeijoapolygonzkevmSession) AcceptAdminRole() (*types.Transaction, error) {
	return _Feijoapolygonzkevm.Contract.AcceptAdminRole(&_Feijoapolygonzkevm.TransactOpts)
}

// AcceptAdminRole is a paid mutator transaction binding the contract method 0x8c3d7301.
//
// Solidity: function acceptAdminRole() returns()
func (_Feijoapolygonzkevm *FeijoapolygonzkevmTransactorSession) AcceptAdminRole() (*types.Transaction, error) {
	return _Feijoapolygonzkevm.Contract.AcceptAdminRole(&_Feijoapolygonzkevm.TransactOpts)
}

// ForceBlob is a paid mutator transaction binding the contract method 0xd2a679b7.
//
// Solidity: function forceBlob(bytes blobData, uint256 polAmount) returns()
func (_Feijoapolygonzkevm *FeijoapolygonzkevmTransactor) ForceBlob(opts *bind.TransactOpts, blobData []byte, polAmount *big.Int) (*types.Transaction, error) {
	return _Feijoapolygonzkevm.contract.Transact(opts, "forceBlob", blobData, polAmount)
}

// ForceBlob is a paid mutator transaction binding the contract method 0xd2a679b7.
//
// Solidity: function forceBlob(bytes blobData, uint256 polAmount) returns()
func (_Feijoapolygonzkevm *FeijoapolygonzkevmSession) ForceBlob(blobData []byte, polAmount *big.Int) (*types.Transaction, error) {
	return _Feijoapolygonzkevm.Contract.ForceBlob(&_Feijoapolygonzkevm.TransactOpts, blobData, polAmount)
}

// ForceBlob is a paid mutator transaction binding the contract method 0xd2a679b7.
//
// Solidity: function forceBlob(bytes blobData, uint256 polAmount) returns()
func (_Feijoapolygonzkevm *FeijoapolygonzkevmTransactorSession) ForceBlob(blobData []byte, polAmount *big.Int) (*types.Transaction, error) {
	return _Feijoapolygonzkevm.Contract.ForceBlob(&_Feijoapolygonzkevm.TransactOpts, blobData, polAmount)
}

// Initialize is a paid mutator transaction binding the contract method 0x71257022.
//
// Solidity: function initialize(address _admin, address sequencer, uint32 networkID, address _gasTokenAddress, string sequencerURL, string _networkName) returns()
func (_Feijoapolygonzkevm *FeijoapolygonzkevmTransactor) Initialize(opts *bind.TransactOpts, _admin common.Address, sequencer common.Address, networkID uint32, _gasTokenAddress common.Address, sequencerURL string, _networkName string) (*types.Transaction, error) {
	return _Feijoapolygonzkevm.contract.Transact(opts, "initialize", _admin, sequencer, networkID, _gasTokenAddress, sequencerURL, _networkName)
}

// Initialize is a paid mutator transaction binding the contract method 0x71257022.
//
// Solidity: function initialize(address _admin, address sequencer, uint32 networkID, address _gasTokenAddress, string sequencerURL, string _networkName) returns()
func (_Feijoapolygonzkevm *FeijoapolygonzkevmSession) Initialize(_admin common.Address, sequencer common.Address, networkID uint32, _gasTokenAddress common.Address, sequencerURL string, _networkName string) (*types.Transaction, error) {
	return _Feijoapolygonzkevm.Contract.Initialize(&_Feijoapolygonzkevm.TransactOpts, _admin, sequencer, networkID, _gasTokenAddress, sequencerURL, _networkName)
}

// Initialize is a paid mutator transaction binding the contract method 0x71257022.
//
// Solidity: function initialize(address _admin, address sequencer, uint32 networkID, address _gasTokenAddress, string sequencerURL, string _networkName) returns()
func (_Feijoapolygonzkevm *FeijoapolygonzkevmTransactorSession) Initialize(_admin common.Address, sequencer common.Address, networkID uint32, _gasTokenAddress common.Address, sequencerURL string, _networkName string) (*types.Transaction, error) {
	return _Feijoapolygonzkevm.Contract.Initialize(&_Feijoapolygonzkevm.TransactOpts, _admin, sequencer, networkID, _gasTokenAddress, sequencerURL, _networkName)
}

// OnVerifySequences is a paid mutator transaction binding the contract method 0x889cfd7a.
//
// Solidity: function onVerifySequences(uint64 lastVerifiedSequenceNum, bytes32 newStateRoot, address aggregator) returns()
func (_Feijoapolygonzkevm *FeijoapolygonzkevmTransactor) OnVerifySequences(opts *bind.TransactOpts, lastVerifiedSequenceNum uint64, newStateRoot [32]byte, aggregator common.Address) (*types.Transaction, error) {
	return _Feijoapolygonzkevm.contract.Transact(opts, "onVerifySequences", lastVerifiedSequenceNum, newStateRoot, aggregator)
}

// OnVerifySequences is a paid mutator transaction binding the contract method 0x889cfd7a.
//
// Solidity: function onVerifySequences(uint64 lastVerifiedSequenceNum, bytes32 newStateRoot, address aggregator) returns()
func (_Feijoapolygonzkevm *FeijoapolygonzkevmSession) OnVerifySequences(lastVerifiedSequenceNum uint64, newStateRoot [32]byte, aggregator common.Address) (*types.Transaction, error) {
	return _Feijoapolygonzkevm.Contract.OnVerifySequences(&_Feijoapolygonzkevm.TransactOpts, lastVerifiedSequenceNum, newStateRoot, aggregator)
}

// OnVerifySequences is a paid mutator transaction binding the contract method 0x889cfd7a.
//
// Solidity: function onVerifySequences(uint64 lastVerifiedSequenceNum, bytes32 newStateRoot, address aggregator) returns()
func (_Feijoapolygonzkevm *FeijoapolygonzkevmTransactorSession) OnVerifySequences(lastVerifiedSequenceNum uint64, newStateRoot [32]byte, aggregator common.Address) (*types.Transaction, error) {
	return _Feijoapolygonzkevm.Contract.OnVerifySequences(&_Feijoapolygonzkevm.TransactOpts, lastVerifiedSequenceNum, newStateRoot, aggregator)
}

// SequenceBlobs is a paid mutator transaction binding the contract method 0x38793b4f.
//
// Solidity: function sequenceBlobs((uint8,bytes)[] blobs, address l2Coinbase, bytes32 finalAccInputHash) returns()
func (_Feijoapolygonzkevm *FeijoapolygonzkevmTransactor) SequenceBlobs(opts *bind.TransactOpts, blobs []PolygonRollupBaseFeijoaBlobData, l2Coinbase common.Address, finalAccInputHash [32]byte) (*types.Transaction, error) {
	return _Feijoapolygonzkevm.contract.Transact(opts, "sequenceBlobs", blobs, l2Coinbase, finalAccInputHash)
}

// SequenceBlobs is a paid mutator transaction binding the contract method 0x38793b4f.
//
// Solidity: function sequenceBlobs((uint8,bytes)[] blobs, address l2Coinbase, bytes32 finalAccInputHash) returns()
func (_Feijoapolygonzkevm *FeijoapolygonzkevmSession) SequenceBlobs(blobs []PolygonRollupBaseFeijoaBlobData, l2Coinbase common.Address, finalAccInputHash [32]byte) (*types.Transaction, error) {
	return _Feijoapolygonzkevm.Contract.SequenceBlobs(&_Feijoapolygonzkevm.TransactOpts, blobs, l2Coinbase, finalAccInputHash)
}

// SequenceBlobs is a paid mutator transaction binding the contract method 0x38793b4f.
//
// Solidity: function sequenceBlobs((uint8,bytes)[] blobs, address l2Coinbase, bytes32 finalAccInputHash) returns()
func (_Feijoapolygonzkevm *FeijoapolygonzkevmTransactorSession) SequenceBlobs(blobs []PolygonRollupBaseFeijoaBlobData, l2Coinbase common.Address, finalAccInputHash [32]byte) (*types.Transaction, error) {
	return _Feijoapolygonzkevm.Contract.SequenceBlobs(&_Feijoapolygonzkevm.TransactOpts, blobs, l2Coinbase, finalAccInputHash)
}

// SequenceForceBlobs is a paid mutator transaction binding the contract method 0x93932a91.
//
// Solidity: function sequenceForceBlobs((uint8,bytes)[] blobs) returns()
func (_Feijoapolygonzkevm *FeijoapolygonzkevmTransactor) SequenceForceBlobs(opts *bind.TransactOpts, blobs []PolygonRollupBaseFeijoaBlobData) (*types.Transaction, error) {
	return _Feijoapolygonzkevm.contract.Transact(opts, "sequenceForceBlobs", blobs)
}

// SequenceForceBlobs is a paid mutator transaction binding the contract method 0x93932a91.
//
// Solidity: function sequenceForceBlobs((uint8,bytes)[] blobs) returns()
func (_Feijoapolygonzkevm *FeijoapolygonzkevmSession) SequenceForceBlobs(blobs []PolygonRollupBaseFeijoaBlobData) (*types.Transaction, error) {
	return _Feijoapolygonzkevm.Contract.SequenceForceBlobs(&_Feijoapolygonzkevm.TransactOpts, blobs)
}

// SequenceForceBlobs is a paid mutator transaction binding the contract method 0x93932a91.
//
// Solidity: function sequenceForceBlobs((uint8,bytes)[] blobs) returns()
func (_Feijoapolygonzkevm *FeijoapolygonzkevmTransactorSession) SequenceForceBlobs(blobs []PolygonRollupBaseFeijoaBlobData) (*types.Transaction, error) {
	return _Feijoapolygonzkevm.Contract.SequenceForceBlobs(&_Feijoapolygonzkevm.TransactOpts, blobs)
}

// SetForceBlobAddress is a paid mutator transaction binding the contract method 0x4bd41065.
//
// Solidity: function setForceBlobAddress(address newForceBlobAddress) returns()
func (_Feijoapolygonzkevm *FeijoapolygonzkevmTransactor) SetForceBlobAddress(opts *bind.TransactOpts, newForceBlobAddress common.Address) (*types.Transaction, error) {
	return _Feijoapolygonzkevm.contract.Transact(opts, "setForceBlobAddress", newForceBlobAddress)
}

// SetForceBlobAddress is a paid mutator transaction binding the contract method 0x4bd41065.
//
// Solidity: function setForceBlobAddress(address newForceBlobAddress) returns()
func (_Feijoapolygonzkevm *FeijoapolygonzkevmSession) SetForceBlobAddress(newForceBlobAddress common.Address) (*types.Transaction, error) {
	return _Feijoapolygonzkevm.Contract.SetForceBlobAddress(&_Feijoapolygonzkevm.TransactOpts, newForceBlobAddress)
}

// SetForceBlobAddress is a paid mutator transaction binding the contract method 0x4bd41065.
//
// Solidity: function setForceBlobAddress(address newForceBlobAddress) returns()
func (_Feijoapolygonzkevm *FeijoapolygonzkevmTransactorSession) SetForceBlobAddress(newForceBlobAddress common.Address) (*types.Transaction, error) {
	return _Feijoapolygonzkevm.Contract.SetForceBlobAddress(&_Feijoapolygonzkevm.TransactOpts, newForceBlobAddress)
}

// SetForceBlobTimeout is a paid mutator transaction binding the contract method 0x730c8e21.
//
// Solidity: function setForceBlobTimeout(uint64 newforceBlobTimeout) returns()
func (_Feijoapolygonzkevm *FeijoapolygonzkevmTransactor) SetForceBlobTimeout(opts *bind.TransactOpts, newforceBlobTimeout uint64) (*types.Transaction, error) {
	return _Feijoapolygonzkevm.contract.Transact(opts, "setForceBlobTimeout", newforceBlobTimeout)
}

// SetForceBlobTimeout is a paid mutator transaction binding the contract method 0x730c8e21.
//
// Solidity: function setForceBlobTimeout(uint64 newforceBlobTimeout) returns()
func (_Feijoapolygonzkevm *FeijoapolygonzkevmSession) SetForceBlobTimeout(newforceBlobTimeout uint64) (*types.Transaction, error) {
	return _Feijoapolygonzkevm.Contract.SetForceBlobTimeout(&_Feijoapolygonzkevm.TransactOpts, newforceBlobTimeout)
}

// SetForceBlobTimeout is a paid mutator transaction binding the contract method 0x730c8e21.
//
// Solidity: function setForceBlobTimeout(uint64 newforceBlobTimeout) returns()
func (_Feijoapolygonzkevm *FeijoapolygonzkevmTransactorSession) SetForceBlobTimeout(newforceBlobTimeout uint64) (*types.Transaction, error) {
	return _Feijoapolygonzkevm.Contract.SetForceBlobTimeout(&_Feijoapolygonzkevm.TransactOpts, newforceBlobTimeout)
}

// SetNetworkName is a paid mutator transaction binding the contract method 0xc0cad302.
//
// Solidity: function setNetworkName(string newNetworkName) returns()
func (_Feijoapolygonzkevm *FeijoapolygonzkevmTransactor) SetNetworkName(opts *bind.TransactOpts, newNetworkName string) (*types.Transaction, error) {
	return _Feijoapolygonzkevm.contract.Transact(opts, "setNetworkName", newNetworkName)
}

// SetNetworkName is a paid mutator transaction binding the contract method 0xc0cad302.
//
// Solidity: function setNetworkName(string newNetworkName) returns()
func (_Feijoapolygonzkevm *FeijoapolygonzkevmSession) SetNetworkName(newNetworkName string) (*types.Transaction, error) {
	return _Feijoapolygonzkevm.Contract.SetNetworkName(&_Feijoapolygonzkevm.TransactOpts, newNetworkName)
}

// SetNetworkName is a paid mutator transaction binding the contract method 0xc0cad302.
//
// Solidity: function setNetworkName(string newNetworkName) returns()
func (_Feijoapolygonzkevm *FeijoapolygonzkevmTransactorSession) SetNetworkName(newNetworkName string) (*types.Transaction, error) {
	return _Feijoapolygonzkevm.Contract.SetNetworkName(&_Feijoapolygonzkevm.TransactOpts, newNetworkName)
}

// SetTrustedSequencer is a paid mutator transaction binding the contract method 0x6ff512cc.
//
// Solidity: function setTrustedSequencer(address newTrustedSequencer) returns()
func (_Feijoapolygonzkevm *FeijoapolygonzkevmTransactor) SetTrustedSequencer(opts *bind.TransactOpts, newTrustedSequencer common.Address) (*types.Transaction, error) {
	return _Feijoapolygonzkevm.contract.Transact(opts, "setTrustedSequencer", newTrustedSequencer)
}

// SetTrustedSequencer is a paid mutator transaction binding the contract method 0x6ff512cc.
//
// Solidity: function setTrustedSequencer(address newTrustedSequencer) returns()
func (_Feijoapolygonzkevm *FeijoapolygonzkevmSession) SetTrustedSequencer(newTrustedSequencer common.Address) (*types.Transaction, error) {
	return _Feijoapolygonzkevm.Contract.SetTrustedSequencer(&_Feijoapolygonzkevm.TransactOpts, newTrustedSequencer)
}

// SetTrustedSequencer is a paid mutator transaction binding the contract method 0x6ff512cc.
//
// Solidity: function setTrustedSequencer(address newTrustedSequencer) returns()
func (_Feijoapolygonzkevm *FeijoapolygonzkevmTransactorSession) SetTrustedSequencer(newTrustedSequencer common.Address) (*types.Transaction, error) {
	return _Feijoapolygonzkevm.Contract.SetTrustedSequencer(&_Feijoapolygonzkevm.TransactOpts, newTrustedSequencer)
}

// SetTrustedSequencerURL is a paid mutator transaction binding the contract method 0xc89e42df.
//
// Solidity: function setTrustedSequencerURL(string newTrustedSequencerURL) returns()
func (_Feijoapolygonzkevm *FeijoapolygonzkevmTransactor) SetTrustedSequencerURL(opts *bind.TransactOpts, newTrustedSequencerURL string) (*types.Transaction, error) {
	return _Feijoapolygonzkevm.contract.Transact(opts, "setTrustedSequencerURL", newTrustedSequencerURL)
}

// SetTrustedSequencerURL is a paid mutator transaction binding the contract method 0xc89e42df.
//
// Solidity: function setTrustedSequencerURL(string newTrustedSequencerURL) returns()
func (_Feijoapolygonzkevm *FeijoapolygonzkevmSession) SetTrustedSequencerURL(newTrustedSequencerURL string) (*types.Transaction, error) {
	return _Feijoapolygonzkevm.Contract.SetTrustedSequencerURL(&_Feijoapolygonzkevm.TransactOpts, newTrustedSequencerURL)
}

// SetTrustedSequencerURL is a paid mutator transaction binding the contract method 0xc89e42df.
//
// Solidity: function setTrustedSequencerURL(string newTrustedSequencerURL) returns()
func (_Feijoapolygonzkevm *FeijoapolygonzkevmTransactorSession) SetTrustedSequencerURL(newTrustedSequencerURL string) (*types.Transaction, error) {
	return _Feijoapolygonzkevm.Contract.SetTrustedSequencerURL(&_Feijoapolygonzkevm.TransactOpts, newTrustedSequencerURL)
}

// TransferAdminRole is a paid mutator transaction binding the contract method 0xada8f919.
//
// Solidity: function transferAdminRole(address newPendingAdmin) returns()
func (_Feijoapolygonzkevm *FeijoapolygonzkevmTransactor) TransferAdminRole(opts *bind.TransactOpts, newPendingAdmin common.Address) (*types.Transaction, error) {
	return _Feijoapolygonzkevm.contract.Transact(opts, "transferAdminRole", newPendingAdmin)
}

// TransferAdminRole is a paid mutator transaction binding the contract method 0xada8f919.
//
// Solidity: function transferAdminRole(address newPendingAdmin) returns()
func (_Feijoapolygonzkevm *FeijoapolygonzkevmSession) TransferAdminRole(newPendingAdmin common.Address) (*types.Transaction, error) {
	return _Feijoapolygonzkevm.Contract.TransferAdminRole(&_Feijoapolygonzkevm.TransactOpts, newPendingAdmin)
}

// TransferAdminRole is a paid mutator transaction binding the contract method 0xada8f919.
//
// Solidity: function transferAdminRole(address newPendingAdmin) returns()
func (_Feijoapolygonzkevm *FeijoapolygonzkevmTransactorSession) TransferAdminRole(newPendingAdmin common.Address) (*types.Transaction, error) {
	return _Feijoapolygonzkevm.Contract.TransferAdminRole(&_Feijoapolygonzkevm.TransactOpts, newPendingAdmin)
}

// FeijoapolygonzkevmAcceptAdminRoleIterator is returned from FilterAcceptAdminRole and is used to iterate over the raw logs and unpacked data for AcceptAdminRole events raised by the Feijoapolygonzkevm contract.
type FeijoapolygonzkevmAcceptAdminRoleIterator struct {
	Event *FeijoapolygonzkevmAcceptAdminRole // Event containing the contract specifics and raw log

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
func (it *FeijoapolygonzkevmAcceptAdminRoleIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeijoapolygonzkevmAcceptAdminRole)
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
		it.Event = new(FeijoapolygonzkevmAcceptAdminRole)
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
func (it *FeijoapolygonzkevmAcceptAdminRoleIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeijoapolygonzkevmAcceptAdminRoleIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeijoapolygonzkevmAcceptAdminRole represents a AcceptAdminRole event raised by the Feijoapolygonzkevm contract.
type FeijoapolygonzkevmAcceptAdminRole struct {
	NewAdmin common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterAcceptAdminRole is a free log retrieval operation binding the contract event 0x056dc487bbf0795d0bbb1b4f0af523a855503cff740bfb4d5475f7a90c091e8e.
//
// Solidity: event AcceptAdminRole(address newAdmin)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmFilterer) FilterAcceptAdminRole(opts *bind.FilterOpts) (*FeijoapolygonzkevmAcceptAdminRoleIterator, error) {

	logs, sub, err := _Feijoapolygonzkevm.contract.FilterLogs(opts, "AcceptAdminRole")
	if err != nil {
		return nil, err
	}
	return &FeijoapolygonzkevmAcceptAdminRoleIterator{contract: _Feijoapolygonzkevm.contract, event: "AcceptAdminRole", logs: logs, sub: sub}, nil
}

// WatchAcceptAdminRole is a free log subscription operation binding the contract event 0x056dc487bbf0795d0bbb1b4f0af523a855503cff740bfb4d5475f7a90c091e8e.
//
// Solidity: event AcceptAdminRole(address newAdmin)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmFilterer) WatchAcceptAdminRole(opts *bind.WatchOpts, sink chan<- *FeijoapolygonzkevmAcceptAdminRole) (event.Subscription, error) {

	logs, sub, err := _Feijoapolygonzkevm.contract.WatchLogs(opts, "AcceptAdminRole")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeijoapolygonzkevmAcceptAdminRole)
				if err := _Feijoapolygonzkevm.contract.UnpackLog(event, "AcceptAdminRole", log); err != nil {
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
func (_Feijoapolygonzkevm *FeijoapolygonzkevmFilterer) ParseAcceptAdminRole(log types.Log) (*FeijoapolygonzkevmAcceptAdminRole, error) {
	event := new(FeijoapolygonzkevmAcceptAdminRole)
	if err := _Feijoapolygonzkevm.contract.UnpackLog(event, "AcceptAdminRole", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeijoapolygonzkevmForceBlobIterator is returned from FilterForceBlob and is used to iterate over the raw logs and unpacked data for ForceBlob events raised by the Feijoapolygonzkevm contract.
type FeijoapolygonzkevmForceBlobIterator struct {
	Event *FeijoapolygonzkevmForceBlob // Event containing the contract specifics and raw log

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
func (it *FeijoapolygonzkevmForceBlobIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeijoapolygonzkevmForceBlob)
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
		it.Event = new(FeijoapolygonzkevmForceBlob)
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
func (it *FeijoapolygonzkevmForceBlobIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeijoapolygonzkevmForceBlobIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeijoapolygonzkevmForceBlob represents a ForceBlob event raised by the Feijoapolygonzkevm contract.
type FeijoapolygonzkevmForceBlob struct {
	ForceBlobNum       uint64
	LastGlobalExitRoot [32]byte
	Sequencer          common.Address
	ZkGasLimit         uint64
	Transactions       []byte
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterForceBlob is a free log retrieval operation binding the contract event 0xb18d758550a6ed34847584be90f0a34b261d8b65bb790891103d5e255aced8b2.
//
// Solidity: event ForceBlob(uint64 indexed forceBlobNum, bytes32 lastGlobalExitRoot, address sequencer, uint64 zkGasLimit, bytes transactions)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmFilterer) FilterForceBlob(opts *bind.FilterOpts, forceBlobNum []uint64) (*FeijoapolygonzkevmForceBlobIterator, error) {

	var forceBlobNumRule []interface{}
	for _, forceBlobNumItem := range forceBlobNum {
		forceBlobNumRule = append(forceBlobNumRule, forceBlobNumItem)
	}

	logs, sub, err := _Feijoapolygonzkevm.contract.FilterLogs(opts, "ForceBlob", forceBlobNumRule)
	if err != nil {
		return nil, err
	}
	return &FeijoapolygonzkevmForceBlobIterator{contract: _Feijoapolygonzkevm.contract, event: "ForceBlob", logs: logs, sub: sub}, nil
}

// WatchForceBlob is a free log subscription operation binding the contract event 0xb18d758550a6ed34847584be90f0a34b261d8b65bb790891103d5e255aced8b2.
//
// Solidity: event ForceBlob(uint64 indexed forceBlobNum, bytes32 lastGlobalExitRoot, address sequencer, uint64 zkGasLimit, bytes transactions)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmFilterer) WatchForceBlob(opts *bind.WatchOpts, sink chan<- *FeijoapolygonzkevmForceBlob, forceBlobNum []uint64) (event.Subscription, error) {

	var forceBlobNumRule []interface{}
	for _, forceBlobNumItem := range forceBlobNum {
		forceBlobNumRule = append(forceBlobNumRule, forceBlobNumItem)
	}

	logs, sub, err := _Feijoapolygonzkevm.contract.WatchLogs(opts, "ForceBlob", forceBlobNumRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeijoapolygonzkevmForceBlob)
				if err := _Feijoapolygonzkevm.contract.UnpackLog(event, "ForceBlob", log); err != nil {
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

// ParseForceBlob is a log parse operation binding the contract event 0xb18d758550a6ed34847584be90f0a34b261d8b65bb790891103d5e255aced8b2.
//
// Solidity: event ForceBlob(uint64 indexed forceBlobNum, bytes32 lastGlobalExitRoot, address sequencer, uint64 zkGasLimit, bytes transactions)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmFilterer) ParseForceBlob(log types.Log) (*FeijoapolygonzkevmForceBlob, error) {
	event := new(FeijoapolygonzkevmForceBlob)
	if err := _Feijoapolygonzkevm.contract.UnpackLog(event, "ForceBlob", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeijoapolygonzkevmInitialSequenceBlobsIterator is returned from FilterInitialSequenceBlobs and is used to iterate over the raw logs and unpacked data for InitialSequenceBlobs events raised by the Feijoapolygonzkevm contract.
type FeijoapolygonzkevmInitialSequenceBlobsIterator struct {
	Event *FeijoapolygonzkevmInitialSequenceBlobs // Event containing the contract specifics and raw log

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
func (it *FeijoapolygonzkevmInitialSequenceBlobsIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeijoapolygonzkevmInitialSequenceBlobs)
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
		it.Event = new(FeijoapolygonzkevmInitialSequenceBlobs)
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
func (it *FeijoapolygonzkevmInitialSequenceBlobsIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeijoapolygonzkevmInitialSequenceBlobsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeijoapolygonzkevmInitialSequenceBlobs represents a InitialSequenceBlobs event raised by the Feijoapolygonzkevm contract.
type FeijoapolygonzkevmInitialSequenceBlobs struct {
	Transactions       []byte
	LastGlobalExitRoot [32]byte
	Sequencer          common.Address
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterInitialSequenceBlobs is a free log retrieval operation binding the contract event 0xfa56300f6f91d53e1c1283e56307c169d72b14a75380df3ecbb5b31b498d3d1e.
//
// Solidity: event InitialSequenceBlobs(bytes transactions, bytes32 lastGlobalExitRoot, address sequencer)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmFilterer) FilterInitialSequenceBlobs(opts *bind.FilterOpts) (*FeijoapolygonzkevmInitialSequenceBlobsIterator, error) {

	logs, sub, err := _Feijoapolygonzkevm.contract.FilterLogs(opts, "InitialSequenceBlobs")
	if err != nil {
		return nil, err
	}
	return &FeijoapolygonzkevmInitialSequenceBlobsIterator{contract: _Feijoapolygonzkevm.contract, event: "InitialSequenceBlobs", logs: logs, sub: sub}, nil
}

// WatchInitialSequenceBlobs is a free log subscription operation binding the contract event 0xfa56300f6f91d53e1c1283e56307c169d72b14a75380df3ecbb5b31b498d3d1e.
//
// Solidity: event InitialSequenceBlobs(bytes transactions, bytes32 lastGlobalExitRoot, address sequencer)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmFilterer) WatchInitialSequenceBlobs(opts *bind.WatchOpts, sink chan<- *FeijoapolygonzkevmInitialSequenceBlobs) (event.Subscription, error) {

	logs, sub, err := _Feijoapolygonzkevm.contract.WatchLogs(opts, "InitialSequenceBlobs")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeijoapolygonzkevmInitialSequenceBlobs)
				if err := _Feijoapolygonzkevm.contract.UnpackLog(event, "InitialSequenceBlobs", log); err != nil {
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

// ParseInitialSequenceBlobs is a log parse operation binding the contract event 0xfa56300f6f91d53e1c1283e56307c169d72b14a75380df3ecbb5b31b498d3d1e.
//
// Solidity: event InitialSequenceBlobs(bytes transactions, bytes32 lastGlobalExitRoot, address sequencer)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmFilterer) ParseInitialSequenceBlobs(log types.Log) (*FeijoapolygonzkevmInitialSequenceBlobs, error) {
	event := new(FeijoapolygonzkevmInitialSequenceBlobs)
	if err := _Feijoapolygonzkevm.contract.UnpackLog(event, "InitialSequenceBlobs", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeijoapolygonzkevmInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Feijoapolygonzkevm contract.
type FeijoapolygonzkevmInitializedIterator struct {
	Event *FeijoapolygonzkevmInitialized // Event containing the contract specifics and raw log

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
func (it *FeijoapolygonzkevmInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeijoapolygonzkevmInitialized)
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
		it.Event = new(FeijoapolygonzkevmInitialized)
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
func (it *FeijoapolygonzkevmInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeijoapolygonzkevmInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeijoapolygonzkevmInitialized represents a Initialized event raised by the Feijoapolygonzkevm contract.
type FeijoapolygonzkevmInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmFilterer) FilterInitialized(opts *bind.FilterOpts) (*FeijoapolygonzkevmInitializedIterator, error) {

	logs, sub, err := _Feijoapolygonzkevm.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &FeijoapolygonzkevmInitializedIterator{contract: _Feijoapolygonzkevm.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *FeijoapolygonzkevmInitialized) (event.Subscription, error) {

	logs, sub, err := _Feijoapolygonzkevm.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeijoapolygonzkevmInitialized)
				if err := _Feijoapolygonzkevm.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Feijoapolygonzkevm *FeijoapolygonzkevmFilterer) ParseInitialized(log types.Log) (*FeijoapolygonzkevmInitialized, error) {
	event := new(FeijoapolygonzkevmInitialized)
	if err := _Feijoapolygonzkevm.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeijoapolygonzkevmSequenceBlobsIterator is returned from FilterSequenceBlobs and is used to iterate over the raw logs and unpacked data for SequenceBlobs events raised by the Feijoapolygonzkevm contract.
type FeijoapolygonzkevmSequenceBlobsIterator struct {
	Event *FeijoapolygonzkevmSequenceBlobs // Event containing the contract specifics and raw log

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
func (it *FeijoapolygonzkevmSequenceBlobsIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeijoapolygonzkevmSequenceBlobs)
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
		it.Event = new(FeijoapolygonzkevmSequenceBlobs)
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
func (it *FeijoapolygonzkevmSequenceBlobsIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeijoapolygonzkevmSequenceBlobsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeijoapolygonzkevmSequenceBlobs represents a SequenceBlobs event raised by the Feijoapolygonzkevm contract.
type FeijoapolygonzkevmSequenceBlobs struct {
	LastBlobSequenced uint64
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterSequenceBlobs is a free log retrieval operation binding the contract event 0x470f4ca4b003755c839b80ab00c3efbeb69d6eafec00e1a3677482933ec1fd0c.
//
// Solidity: event SequenceBlobs(uint64 indexed lastBlobSequenced)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmFilterer) FilterSequenceBlobs(opts *bind.FilterOpts, lastBlobSequenced []uint64) (*FeijoapolygonzkevmSequenceBlobsIterator, error) {

	var lastBlobSequencedRule []interface{}
	for _, lastBlobSequencedItem := range lastBlobSequenced {
		lastBlobSequencedRule = append(lastBlobSequencedRule, lastBlobSequencedItem)
	}

	logs, sub, err := _Feijoapolygonzkevm.contract.FilterLogs(opts, "SequenceBlobs", lastBlobSequencedRule)
	if err != nil {
		return nil, err
	}
	return &FeijoapolygonzkevmSequenceBlobsIterator{contract: _Feijoapolygonzkevm.contract, event: "SequenceBlobs", logs: logs, sub: sub}, nil
}

// WatchSequenceBlobs is a free log subscription operation binding the contract event 0x470f4ca4b003755c839b80ab00c3efbeb69d6eafec00e1a3677482933ec1fd0c.
//
// Solidity: event SequenceBlobs(uint64 indexed lastBlobSequenced)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmFilterer) WatchSequenceBlobs(opts *bind.WatchOpts, sink chan<- *FeijoapolygonzkevmSequenceBlobs, lastBlobSequenced []uint64) (event.Subscription, error) {

	var lastBlobSequencedRule []interface{}
	for _, lastBlobSequencedItem := range lastBlobSequenced {
		lastBlobSequencedRule = append(lastBlobSequencedRule, lastBlobSequencedItem)
	}

	logs, sub, err := _Feijoapolygonzkevm.contract.WatchLogs(opts, "SequenceBlobs", lastBlobSequencedRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeijoapolygonzkevmSequenceBlobs)
				if err := _Feijoapolygonzkevm.contract.UnpackLog(event, "SequenceBlobs", log); err != nil {
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

// ParseSequenceBlobs is a log parse operation binding the contract event 0x470f4ca4b003755c839b80ab00c3efbeb69d6eafec00e1a3677482933ec1fd0c.
//
// Solidity: event SequenceBlobs(uint64 indexed lastBlobSequenced)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmFilterer) ParseSequenceBlobs(log types.Log) (*FeijoapolygonzkevmSequenceBlobs, error) {
	event := new(FeijoapolygonzkevmSequenceBlobs)
	if err := _Feijoapolygonzkevm.contract.UnpackLog(event, "SequenceBlobs", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeijoapolygonzkevmSequenceForceBlobsIterator is returned from FilterSequenceForceBlobs and is used to iterate over the raw logs and unpacked data for SequenceForceBlobs events raised by the Feijoapolygonzkevm contract.
type FeijoapolygonzkevmSequenceForceBlobsIterator struct {
	Event *FeijoapolygonzkevmSequenceForceBlobs // Event containing the contract specifics and raw log

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
func (it *FeijoapolygonzkevmSequenceForceBlobsIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeijoapolygonzkevmSequenceForceBlobs)
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
		it.Event = new(FeijoapolygonzkevmSequenceForceBlobs)
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
func (it *FeijoapolygonzkevmSequenceForceBlobsIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeijoapolygonzkevmSequenceForceBlobsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeijoapolygonzkevmSequenceForceBlobs represents a SequenceForceBlobs event raised by the Feijoapolygonzkevm contract.
type FeijoapolygonzkevmSequenceForceBlobs struct {
	NumBlob uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterSequenceForceBlobs is a free log retrieval operation binding the contract event 0x049b259b0b684f32f1d8b43d76cf6cb3c674b94697bda3290f6ec63252cfe892.
//
// Solidity: event SequenceForceBlobs(uint64 indexed numBlob)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmFilterer) FilterSequenceForceBlobs(opts *bind.FilterOpts, numBlob []uint64) (*FeijoapolygonzkevmSequenceForceBlobsIterator, error) {

	var numBlobRule []interface{}
	for _, numBlobItem := range numBlob {
		numBlobRule = append(numBlobRule, numBlobItem)
	}

	logs, sub, err := _Feijoapolygonzkevm.contract.FilterLogs(opts, "SequenceForceBlobs", numBlobRule)
	if err != nil {
		return nil, err
	}
	return &FeijoapolygonzkevmSequenceForceBlobsIterator{contract: _Feijoapolygonzkevm.contract, event: "SequenceForceBlobs", logs: logs, sub: sub}, nil
}

// WatchSequenceForceBlobs is a free log subscription operation binding the contract event 0x049b259b0b684f32f1d8b43d76cf6cb3c674b94697bda3290f6ec63252cfe892.
//
// Solidity: event SequenceForceBlobs(uint64 indexed numBlob)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmFilterer) WatchSequenceForceBlobs(opts *bind.WatchOpts, sink chan<- *FeijoapolygonzkevmSequenceForceBlobs, numBlob []uint64) (event.Subscription, error) {

	var numBlobRule []interface{}
	for _, numBlobItem := range numBlob {
		numBlobRule = append(numBlobRule, numBlobItem)
	}

	logs, sub, err := _Feijoapolygonzkevm.contract.WatchLogs(opts, "SequenceForceBlobs", numBlobRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeijoapolygonzkevmSequenceForceBlobs)
				if err := _Feijoapolygonzkevm.contract.UnpackLog(event, "SequenceForceBlobs", log); err != nil {
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

// ParseSequenceForceBlobs is a log parse operation binding the contract event 0x049b259b0b684f32f1d8b43d76cf6cb3c674b94697bda3290f6ec63252cfe892.
//
// Solidity: event SequenceForceBlobs(uint64 indexed numBlob)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmFilterer) ParseSequenceForceBlobs(log types.Log) (*FeijoapolygonzkevmSequenceForceBlobs, error) {
	event := new(FeijoapolygonzkevmSequenceForceBlobs)
	if err := _Feijoapolygonzkevm.contract.UnpackLog(event, "SequenceForceBlobs", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeijoapolygonzkevmSetForceBlobAddressIterator is returned from FilterSetForceBlobAddress and is used to iterate over the raw logs and unpacked data for SetForceBlobAddress events raised by the Feijoapolygonzkevm contract.
type FeijoapolygonzkevmSetForceBlobAddressIterator struct {
	Event *FeijoapolygonzkevmSetForceBlobAddress // Event containing the contract specifics and raw log

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
func (it *FeijoapolygonzkevmSetForceBlobAddressIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeijoapolygonzkevmSetForceBlobAddress)
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
		it.Event = new(FeijoapolygonzkevmSetForceBlobAddress)
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
func (it *FeijoapolygonzkevmSetForceBlobAddressIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeijoapolygonzkevmSetForceBlobAddressIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeijoapolygonzkevmSetForceBlobAddress represents a SetForceBlobAddress event raised by the Feijoapolygonzkevm contract.
type FeijoapolygonzkevmSetForceBlobAddress struct {
	NewForceBlobAddress common.Address
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterSetForceBlobAddress is a free log retrieval operation binding the contract event 0x2261b2af55eeb3b995b5e300659fa8e59827ff8fc99ff3a5baf5af0835aab9dd.
//
// Solidity: event SetForceBlobAddress(address newForceBlobAddress)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmFilterer) FilterSetForceBlobAddress(opts *bind.FilterOpts) (*FeijoapolygonzkevmSetForceBlobAddressIterator, error) {

	logs, sub, err := _Feijoapolygonzkevm.contract.FilterLogs(opts, "SetForceBlobAddress")
	if err != nil {
		return nil, err
	}
	return &FeijoapolygonzkevmSetForceBlobAddressIterator{contract: _Feijoapolygonzkevm.contract, event: "SetForceBlobAddress", logs: logs, sub: sub}, nil
}

// WatchSetForceBlobAddress is a free log subscription operation binding the contract event 0x2261b2af55eeb3b995b5e300659fa8e59827ff8fc99ff3a5baf5af0835aab9dd.
//
// Solidity: event SetForceBlobAddress(address newForceBlobAddress)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmFilterer) WatchSetForceBlobAddress(opts *bind.WatchOpts, sink chan<- *FeijoapolygonzkevmSetForceBlobAddress) (event.Subscription, error) {

	logs, sub, err := _Feijoapolygonzkevm.contract.WatchLogs(opts, "SetForceBlobAddress")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeijoapolygonzkevmSetForceBlobAddress)
				if err := _Feijoapolygonzkevm.contract.UnpackLog(event, "SetForceBlobAddress", log); err != nil {
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

// ParseSetForceBlobAddress is a log parse operation binding the contract event 0x2261b2af55eeb3b995b5e300659fa8e59827ff8fc99ff3a5baf5af0835aab9dd.
//
// Solidity: event SetForceBlobAddress(address newForceBlobAddress)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmFilterer) ParseSetForceBlobAddress(log types.Log) (*FeijoapolygonzkevmSetForceBlobAddress, error) {
	event := new(FeijoapolygonzkevmSetForceBlobAddress)
	if err := _Feijoapolygonzkevm.contract.UnpackLog(event, "SetForceBlobAddress", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeijoapolygonzkevmSetForceBlobTimeoutIterator is returned from FilterSetForceBlobTimeout and is used to iterate over the raw logs and unpacked data for SetForceBlobTimeout events raised by the Feijoapolygonzkevm contract.
type FeijoapolygonzkevmSetForceBlobTimeoutIterator struct {
	Event *FeijoapolygonzkevmSetForceBlobTimeout // Event containing the contract specifics and raw log

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
func (it *FeijoapolygonzkevmSetForceBlobTimeoutIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeijoapolygonzkevmSetForceBlobTimeout)
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
		it.Event = new(FeijoapolygonzkevmSetForceBlobTimeout)
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
func (it *FeijoapolygonzkevmSetForceBlobTimeoutIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeijoapolygonzkevmSetForceBlobTimeoutIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeijoapolygonzkevmSetForceBlobTimeout represents a SetForceBlobTimeout event raised by the Feijoapolygonzkevm contract.
type FeijoapolygonzkevmSetForceBlobTimeout struct {
	NewforceBlobTimeout uint64
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterSetForceBlobTimeout is a free log retrieval operation binding the contract event 0xa6db492cb43063288b0b5d7c271f8df34607c41fc9347c0664e1ce325cc728e8.
//
// Solidity: event SetForceBlobTimeout(uint64 newforceBlobTimeout)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmFilterer) FilterSetForceBlobTimeout(opts *bind.FilterOpts) (*FeijoapolygonzkevmSetForceBlobTimeoutIterator, error) {

	logs, sub, err := _Feijoapolygonzkevm.contract.FilterLogs(opts, "SetForceBlobTimeout")
	if err != nil {
		return nil, err
	}
	return &FeijoapolygonzkevmSetForceBlobTimeoutIterator{contract: _Feijoapolygonzkevm.contract, event: "SetForceBlobTimeout", logs: logs, sub: sub}, nil
}

// WatchSetForceBlobTimeout is a free log subscription operation binding the contract event 0xa6db492cb43063288b0b5d7c271f8df34607c41fc9347c0664e1ce325cc728e8.
//
// Solidity: event SetForceBlobTimeout(uint64 newforceBlobTimeout)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmFilterer) WatchSetForceBlobTimeout(opts *bind.WatchOpts, sink chan<- *FeijoapolygonzkevmSetForceBlobTimeout) (event.Subscription, error) {

	logs, sub, err := _Feijoapolygonzkevm.contract.WatchLogs(opts, "SetForceBlobTimeout")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeijoapolygonzkevmSetForceBlobTimeout)
				if err := _Feijoapolygonzkevm.contract.UnpackLog(event, "SetForceBlobTimeout", log); err != nil {
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

// ParseSetForceBlobTimeout is a log parse operation binding the contract event 0xa6db492cb43063288b0b5d7c271f8df34607c41fc9347c0664e1ce325cc728e8.
//
// Solidity: event SetForceBlobTimeout(uint64 newforceBlobTimeout)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmFilterer) ParseSetForceBlobTimeout(log types.Log) (*FeijoapolygonzkevmSetForceBlobTimeout, error) {
	event := new(FeijoapolygonzkevmSetForceBlobTimeout)
	if err := _Feijoapolygonzkevm.contract.UnpackLog(event, "SetForceBlobTimeout", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeijoapolygonzkevmSetNetworkNameIterator is returned from FilterSetNetworkName and is used to iterate over the raw logs and unpacked data for SetNetworkName events raised by the Feijoapolygonzkevm contract.
type FeijoapolygonzkevmSetNetworkNameIterator struct {
	Event *FeijoapolygonzkevmSetNetworkName // Event containing the contract specifics and raw log

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
func (it *FeijoapolygonzkevmSetNetworkNameIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeijoapolygonzkevmSetNetworkName)
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
		it.Event = new(FeijoapolygonzkevmSetNetworkName)
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
func (it *FeijoapolygonzkevmSetNetworkNameIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeijoapolygonzkevmSetNetworkNameIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeijoapolygonzkevmSetNetworkName represents a SetNetworkName event raised by the Feijoapolygonzkevm contract.
type FeijoapolygonzkevmSetNetworkName struct {
	NewNetworkName string
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterSetNetworkName is a free log retrieval operation binding the contract event 0xcc3b37f0de47ea5ce245c3502f0d4e414c34664023b8463db2fe451fee5e6992.
//
// Solidity: event SetNetworkName(string newNetworkName)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmFilterer) FilterSetNetworkName(opts *bind.FilterOpts) (*FeijoapolygonzkevmSetNetworkNameIterator, error) {

	logs, sub, err := _Feijoapolygonzkevm.contract.FilterLogs(opts, "SetNetworkName")
	if err != nil {
		return nil, err
	}
	return &FeijoapolygonzkevmSetNetworkNameIterator{contract: _Feijoapolygonzkevm.contract, event: "SetNetworkName", logs: logs, sub: sub}, nil
}

// WatchSetNetworkName is a free log subscription operation binding the contract event 0xcc3b37f0de47ea5ce245c3502f0d4e414c34664023b8463db2fe451fee5e6992.
//
// Solidity: event SetNetworkName(string newNetworkName)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmFilterer) WatchSetNetworkName(opts *bind.WatchOpts, sink chan<- *FeijoapolygonzkevmSetNetworkName) (event.Subscription, error) {

	logs, sub, err := _Feijoapolygonzkevm.contract.WatchLogs(opts, "SetNetworkName")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeijoapolygonzkevmSetNetworkName)
				if err := _Feijoapolygonzkevm.contract.UnpackLog(event, "SetNetworkName", log); err != nil {
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

// ParseSetNetworkName is a log parse operation binding the contract event 0xcc3b37f0de47ea5ce245c3502f0d4e414c34664023b8463db2fe451fee5e6992.
//
// Solidity: event SetNetworkName(string newNetworkName)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmFilterer) ParseSetNetworkName(log types.Log) (*FeijoapolygonzkevmSetNetworkName, error) {
	event := new(FeijoapolygonzkevmSetNetworkName)
	if err := _Feijoapolygonzkevm.contract.UnpackLog(event, "SetNetworkName", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeijoapolygonzkevmSetTrustedSequencerIterator is returned from FilterSetTrustedSequencer and is used to iterate over the raw logs and unpacked data for SetTrustedSequencer events raised by the Feijoapolygonzkevm contract.
type FeijoapolygonzkevmSetTrustedSequencerIterator struct {
	Event *FeijoapolygonzkevmSetTrustedSequencer // Event containing the contract specifics and raw log

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
func (it *FeijoapolygonzkevmSetTrustedSequencerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeijoapolygonzkevmSetTrustedSequencer)
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
		it.Event = new(FeijoapolygonzkevmSetTrustedSequencer)
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
func (it *FeijoapolygonzkevmSetTrustedSequencerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeijoapolygonzkevmSetTrustedSequencerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeijoapolygonzkevmSetTrustedSequencer represents a SetTrustedSequencer event raised by the Feijoapolygonzkevm contract.
type FeijoapolygonzkevmSetTrustedSequencer struct {
	NewTrustedSequencer common.Address
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterSetTrustedSequencer is a free log retrieval operation binding the contract event 0xf54144f9611984021529f814a1cb6a41e22c58351510a0d9f7e822618abb9cc0.
//
// Solidity: event SetTrustedSequencer(address newTrustedSequencer)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmFilterer) FilterSetTrustedSequencer(opts *bind.FilterOpts) (*FeijoapolygonzkevmSetTrustedSequencerIterator, error) {

	logs, sub, err := _Feijoapolygonzkevm.contract.FilterLogs(opts, "SetTrustedSequencer")
	if err != nil {
		return nil, err
	}
	return &FeijoapolygonzkevmSetTrustedSequencerIterator{contract: _Feijoapolygonzkevm.contract, event: "SetTrustedSequencer", logs: logs, sub: sub}, nil
}

// WatchSetTrustedSequencer is a free log subscription operation binding the contract event 0xf54144f9611984021529f814a1cb6a41e22c58351510a0d9f7e822618abb9cc0.
//
// Solidity: event SetTrustedSequencer(address newTrustedSequencer)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmFilterer) WatchSetTrustedSequencer(opts *bind.WatchOpts, sink chan<- *FeijoapolygonzkevmSetTrustedSequencer) (event.Subscription, error) {

	logs, sub, err := _Feijoapolygonzkevm.contract.WatchLogs(opts, "SetTrustedSequencer")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeijoapolygonzkevmSetTrustedSequencer)
				if err := _Feijoapolygonzkevm.contract.UnpackLog(event, "SetTrustedSequencer", log); err != nil {
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

// ParseSetTrustedSequencer is a log parse operation binding the contract event 0xf54144f9611984021529f814a1cb6a41e22c58351510a0d9f7e822618abb9cc0.
//
// Solidity: event SetTrustedSequencer(address newTrustedSequencer)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmFilterer) ParseSetTrustedSequencer(log types.Log) (*FeijoapolygonzkevmSetTrustedSequencer, error) {
	event := new(FeijoapolygonzkevmSetTrustedSequencer)
	if err := _Feijoapolygonzkevm.contract.UnpackLog(event, "SetTrustedSequencer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeijoapolygonzkevmSetTrustedSequencerURLIterator is returned from FilterSetTrustedSequencerURL and is used to iterate over the raw logs and unpacked data for SetTrustedSequencerURL events raised by the Feijoapolygonzkevm contract.
type FeijoapolygonzkevmSetTrustedSequencerURLIterator struct {
	Event *FeijoapolygonzkevmSetTrustedSequencerURL // Event containing the contract specifics and raw log

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
func (it *FeijoapolygonzkevmSetTrustedSequencerURLIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeijoapolygonzkevmSetTrustedSequencerURL)
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
		it.Event = new(FeijoapolygonzkevmSetTrustedSequencerURL)
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
func (it *FeijoapolygonzkevmSetTrustedSequencerURLIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeijoapolygonzkevmSetTrustedSequencerURLIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeijoapolygonzkevmSetTrustedSequencerURL represents a SetTrustedSequencerURL event raised by the Feijoapolygonzkevm contract.
type FeijoapolygonzkevmSetTrustedSequencerURL struct {
	NewTrustedSequencerURL string
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterSetTrustedSequencerURL is a free log retrieval operation binding the contract event 0x6b8f723a4c7a5335cafae8a598a0aa0301be1387c037dccc085b62add6448b20.
//
// Solidity: event SetTrustedSequencerURL(string newTrustedSequencerURL)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmFilterer) FilterSetTrustedSequencerURL(opts *bind.FilterOpts) (*FeijoapolygonzkevmSetTrustedSequencerURLIterator, error) {

	logs, sub, err := _Feijoapolygonzkevm.contract.FilterLogs(opts, "SetTrustedSequencerURL")
	if err != nil {
		return nil, err
	}
	return &FeijoapolygonzkevmSetTrustedSequencerURLIterator{contract: _Feijoapolygonzkevm.contract, event: "SetTrustedSequencerURL", logs: logs, sub: sub}, nil
}

// WatchSetTrustedSequencerURL is a free log subscription operation binding the contract event 0x6b8f723a4c7a5335cafae8a598a0aa0301be1387c037dccc085b62add6448b20.
//
// Solidity: event SetTrustedSequencerURL(string newTrustedSequencerURL)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmFilterer) WatchSetTrustedSequencerURL(opts *bind.WatchOpts, sink chan<- *FeijoapolygonzkevmSetTrustedSequencerURL) (event.Subscription, error) {

	logs, sub, err := _Feijoapolygonzkevm.contract.WatchLogs(opts, "SetTrustedSequencerURL")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeijoapolygonzkevmSetTrustedSequencerURL)
				if err := _Feijoapolygonzkevm.contract.UnpackLog(event, "SetTrustedSequencerURL", log); err != nil {
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

// ParseSetTrustedSequencerURL is a log parse operation binding the contract event 0x6b8f723a4c7a5335cafae8a598a0aa0301be1387c037dccc085b62add6448b20.
//
// Solidity: event SetTrustedSequencerURL(string newTrustedSequencerURL)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmFilterer) ParseSetTrustedSequencerURL(log types.Log) (*FeijoapolygonzkevmSetTrustedSequencerURL, error) {
	event := new(FeijoapolygonzkevmSetTrustedSequencerURL)
	if err := _Feijoapolygonzkevm.contract.UnpackLog(event, "SetTrustedSequencerURL", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeijoapolygonzkevmTransferAdminRoleIterator is returned from FilterTransferAdminRole and is used to iterate over the raw logs and unpacked data for TransferAdminRole events raised by the Feijoapolygonzkevm contract.
type FeijoapolygonzkevmTransferAdminRoleIterator struct {
	Event *FeijoapolygonzkevmTransferAdminRole // Event containing the contract specifics and raw log

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
func (it *FeijoapolygonzkevmTransferAdminRoleIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeijoapolygonzkevmTransferAdminRole)
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
		it.Event = new(FeijoapolygonzkevmTransferAdminRole)
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
func (it *FeijoapolygonzkevmTransferAdminRoleIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeijoapolygonzkevmTransferAdminRoleIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeijoapolygonzkevmTransferAdminRole represents a TransferAdminRole event raised by the Feijoapolygonzkevm contract.
type FeijoapolygonzkevmTransferAdminRole struct {
	NewPendingAdmin common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterTransferAdminRole is a free log retrieval operation binding the contract event 0xa5b56b7906fd0a20e3f35120dd8343db1e12e037a6c90111c7e42885e82a1ce6.
//
// Solidity: event TransferAdminRole(address newPendingAdmin)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmFilterer) FilterTransferAdminRole(opts *bind.FilterOpts) (*FeijoapolygonzkevmTransferAdminRoleIterator, error) {

	logs, sub, err := _Feijoapolygonzkevm.contract.FilterLogs(opts, "TransferAdminRole")
	if err != nil {
		return nil, err
	}
	return &FeijoapolygonzkevmTransferAdminRoleIterator{contract: _Feijoapolygonzkevm.contract, event: "TransferAdminRole", logs: logs, sub: sub}, nil
}

// WatchTransferAdminRole is a free log subscription operation binding the contract event 0xa5b56b7906fd0a20e3f35120dd8343db1e12e037a6c90111c7e42885e82a1ce6.
//
// Solidity: event TransferAdminRole(address newPendingAdmin)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmFilterer) WatchTransferAdminRole(opts *bind.WatchOpts, sink chan<- *FeijoapolygonzkevmTransferAdminRole) (event.Subscription, error) {

	logs, sub, err := _Feijoapolygonzkevm.contract.WatchLogs(opts, "TransferAdminRole")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeijoapolygonzkevmTransferAdminRole)
				if err := _Feijoapolygonzkevm.contract.UnpackLog(event, "TransferAdminRole", log); err != nil {
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
func (_Feijoapolygonzkevm *FeijoapolygonzkevmFilterer) ParseTransferAdminRole(log types.Log) (*FeijoapolygonzkevmTransferAdminRole, error) {
	event := new(FeijoapolygonzkevmTransferAdminRole)
	if err := _Feijoapolygonzkevm.contract.UnpackLog(event, "TransferAdminRole", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeijoapolygonzkevmVerifyBlobsIterator is returned from FilterVerifyBlobs and is used to iterate over the raw logs and unpacked data for VerifyBlobs events raised by the Feijoapolygonzkevm contract.
type FeijoapolygonzkevmVerifyBlobsIterator struct {
	Event *FeijoapolygonzkevmVerifyBlobs // Event containing the contract specifics and raw log

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
func (it *FeijoapolygonzkevmVerifyBlobsIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeijoapolygonzkevmVerifyBlobs)
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
		it.Event = new(FeijoapolygonzkevmVerifyBlobs)
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
func (it *FeijoapolygonzkevmVerifyBlobsIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeijoapolygonzkevmVerifyBlobsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeijoapolygonzkevmVerifyBlobs represents a VerifyBlobs event raised by the Feijoapolygonzkevm contract.
type FeijoapolygonzkevmVerifyBlobs struct {
	SequneceNum uint64
	StateRoot   [32]byte
	Aggregator  common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterVerifyBlobs is a free log retrieval operation binding the contract event 0xb19baa6f6271636400b99e9e5b3289ec1e0d74e6204a27f296cc4715ff9ded55.
//
// Solidity: event VerifyBlobs(uint64 indexed sequneceNum, bytes32 stateRoot, address indexed aggregator)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmFilterer) FilterVerifyBlobs(opts *bind.FilterOpts, sequneceNum []uint64, aggregator []common.Address) (*FeijoapolygonzkevmVerifyBlobsIterator, error) {

	var sequneceNumRule []interface{}
	for _, sequneceNumItem := range sequneceNum {
		sequneceNumRule = append(sequneceNumRule, sequneceNumItem)
	}

	var aggregatorRule []interface{}
	for _, aggregatorItem := range aggregator {
		aggregatorRule = append(aggregatorRule, aggregatorItem)
	}

	logs, sub, err := _Feijoapolygonzkevm.contract.FilterLogs(opts, "VerifyBlobs", sequneceNumRule, aggregatorRule)
	if err != nil {
		return nil, err
	}
	return &FeijoapolygonzkevmVerifyBlobsIterator{contract: _Feijoapolygonzkevm.contract, event: "VerifyBlobs", logs: logs, sub: sub}, nil
}

// WatchVerifyBlobs is a free log subscription operation binding the contract event 0xb19baa6f6271636400b99e9e5b3289ec1e0d74e6204a27f296cc4715ff9ded55.
//
// Solidity: event VerifyBlobs(uint64 indexed sequneceNum, bytes32 stateRoot, address indexed aggregator)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmFilterer) WatchVerifyBlobs(opts *bind.WatchOpts, sink chan<- *FeijoapolygonzkevmVerifyBlobs, sequneceNum []uint64, aggregator []common.Address) (event.Subscription, error) {

	var sequneceNumRule []interface{}
	for _, sequneceNumItem := range sequneceNum {
		sequneceNumRule = append(sequneceNumRule, sequneceNumItem)
	}

	var aggregatorRule []interface{}
	for _, aggregatorItem := range aggregator {
		aggregatorRule = append(aggregatorRule, aggregatorItem)
	}

	logs, sub, err := _Feijoapolygonzkevm.contract.WatchLogs(opts, "VerifyBlobs", sequneceNumRule, aggregatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeijoapolygonzkevmVerifyBlobs)
				if err := _Feijoapolygonzkevm.contract.UnpackLog(event, "VerifyBlobs", log); err != nil {
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

// ParseVerifyBlobs is a log parse operation binding the contract event 0xb19baa6f6271636400b99e9e5b3289ec1e0d74e6204a27f296cc4715ff9ded55.
//
// Solidity: event VerifyBlobs(uint64 indexed sequneceNum, bytes32 stateRoot, address indexed aggregator)
func (_Feijoapolygonzkevm *FeijoapolygonzkevmFilterer) ParseVerifyBlobs(log types.Log) (*FeijoapolygonzkevmVerifyBlobs, error) {
	event := new(FeijoapolygonzkevmVerifyBlobs)
	if err := _Feijoapolygonzkevm.contract.UnpackLog(event, "VerifyBlobs", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
