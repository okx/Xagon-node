// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package elderberrypolygonzkevm

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

// PolygonRollupBaseEtrogBatchData is an auto generated low-level Go binding around an user-defined struct.
type PolygonRollupBaseEtrogBatchData struct {
	Transactions         []byte
	ForcedGlobalExitRoot [32]byte
	ForcedTimestamp      uint64
	ForcedBlockHashL1    [32]byte
}

// ElderberrypolygonzkevmMetaData contains all meta data concerning the Elderberrypolygonzkevm contract.
var ElderberrypolygonzkevmMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIPolygonZkEVMGlobalExitRootV2\",\"name\":\"_globalExitRootManager\",\"type\":\"address\"},{\"internalType\":\"contractIERC20Upgradeable\",\"name\":\"_pol\",\"type\":\"address\"},{\"internalType\":\"contractIPolygonZkEVMBridgeV2\",\"name\":\"_bridgeAddress\",\"type\":\"address\"},{\"internalType\":\"contractPolygonRollupManager\",\"name\":\"_rollupManager\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"BatchAlreadyVerified\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BatchNotSequencedOrNotSequenceEnd\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ExceedMaxVerifyBatches\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FinalNumBatchBelowLastVerifiedBatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FinalNumBatchDoesNotMatchPendingState\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FinalPendingStateNumInvalid\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ForceBatchNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ForceBatchTimeoutNotExpired\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ForceBatchesAlreadyActive\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ForceBatchesDecentralized\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ForceBatchesNotAllowedOnEmergencyState\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ForceBatchesOverflow\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ForcedDataDoesNotMatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GasTokenNetworkMustBeZeroOnEther\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"GlobalExitRootNotExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"HaltTimeoutNotExpired\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"HaltTimeoutNotExpiredAfterEmergencyState\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"HugeTokenMetadataNotSupported\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InitNumBatchAboveLastVerifiedBatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InitNumBatchDoesNotMatchPendingState\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitializeTransaction\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidProof\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRangeBatchTimeTarget\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRangeForceBatchTimeout\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRangeMultiplierBatchFee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NewAccInputHashDoesNotExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NewPendingStateTimeoutMustBeLower\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NewStateRootNotInsidePrime\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NewTrustedAggregatorTimeoutMustBeLower\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotEnoughMaticAmount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotEnoughPOLAmount\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OldAccInputHashDoesNotExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OldStateRootDoesNotExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyPendingAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyRollupManager\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyTrustedAggregator\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyTrustedSequencer\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PendingStateDoesNotExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PendingStateInvalid\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PendingStateNotConsolidable\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PendingStateTimeoutExceedHaltAggregationTimeout\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SequenceZeroBatches\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SequencedTimestampBelowForcedTimestamp\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SequencedTimestampInvalid\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"StoredRootMustBeDifferentThanNewRoot\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TransactionsLengthAboveMax\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TrustedAggregatorTimeoutExceedHaltAggregationTimeout\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TrustedAggregatorTimeoutNotExpired\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AcceptAdminRole\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"forceBatchNum\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"lastGlobalExitRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sequencer\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"transactions\",\"type\":\"bytes\"}],\"name\":\"ForceBatch\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"transactions\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"lastGlobalExitRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sequencer\",\"type\":\"address\"}],\"name\":\"InitialSequenceBatches\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"numBatch\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"l1InfoRoot\",\"type\":\"bytes32\"}],\"name\":\"SequenceBatches\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"numBatch\",\"type\":\"uint64\"}],\"name\":\"SequenceForceBatches\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newForceBatchAddress\",\"type\":\"address\"}],\"name\":\"SetForceBatchAddress\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"newforceBatchTimeout\",\"type\":\"uint64\"}],\"name\":\"SetForceBatchTimeout\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newTrustedSequencer\",\"type\":\"address\"}],\"name\":\"SetTrustedSequencer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"newTrustedSequencerURL\",\"type\":\"string\"}],\"name\":\"SetTrustedSequencerURL\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newPendingAdmin\",\"type\":\"address\"}],\"name\":\"TransferAdminRole\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"numBatch\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"transactions\",\"type\":\"bytes\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"lastGlobalExitRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"sequencer\",\"type\":\"address\"}],\"name\":\"UpdateEtrogSequence\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint64\",\"name\":\"numBatch\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"aggregator\",\"type\":\"address\"}],\"name\":\"VerifyBatches\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"GLOBAL_EXIT_ROOT_MANAGER_L2\",\"outputs\":[{\"internalType\":\"contractIBasePolygonZkEVMGlobalExitRoot\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"INITIALIZE_TX_BRIDGE_LIST_LEN_LEN\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"INITIALIZE_TX_BRIDGE_PARAMS\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"INITIALIZE_TX_BRIDGE_PARAMS_AFTER_BRIDGE_ADDRESS\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"INITIALIZE_TX_BRIDGE_PARAMS_AFTER_BRIDGE_ADDRESS_EMPTY_METADATA\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"INITIALIZE_TX_CONSTANT_BYTES\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"INITIALIZE_TX_CONSTANT_BYTES_EMPTY_METADATA\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"INITIALIZE_TX_DATA_LEN_EMPTY_METADATA\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"INITIALIZE_TX_EFFECTIVE_PERCENTAGE\",\"outputs\":[{\"internalType\":\"bytes1\",\"name\":\"\",\"type\":\"bytes1\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"SET_UP_ETROG_TX\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"SIGNATURE_INITIALIZE_TX_R\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"SIGNATURE_INITIALIZE_TX_S\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"SIGNATURE_INITIALIZE_TX_V\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"acceptAdminRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"admin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"bridgeAddress\",\"outputs\":[{\"internalType\":\"contractIPolygonZkEVMBridgeV2\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"calculatePolPerForceBatch\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"transactions\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"polAmount\",\"type\":\"uint256\"}],\"name\":\"forceBatch\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"forceBatchAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"forceBatchTimeout\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"forcedBatches\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"gasTokenAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"gasTokenNetwork\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"networkID\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"_gasTokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"_gasTokenNetwork\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"_gasTokenMetadata\",\"type\":\"bytes\"}],\"name\":\"generateInitializeTransaction\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"globalExitRootManager\",\"outputs\":[{\"internalType\":\"contractIPolygonZkEVMGlobalExitRootV2\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_admin\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"sequencer\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"networkID\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"_gasTokenAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"sequencerURL\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_networkName\",\"type\":\"string\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_admin\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_trustedSequencer\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"_trustedSequencerURL\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_networkName\",\"type\":\"string\"},{\"internalType\":\"bytes32\",\"name\":\"_lastAccInputHash\",\"type\":\"bytes32\"}],\"name\":\"initializeUpgrade\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastAccInputHash\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastForceBatch\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastForceBatchSequenced\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"networkName\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"lastVerifiedBatch\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"aggregator\",\"type\":\"address\"}],\"name\":\"onVerifyBatches\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pendingAdmin\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pol\",\"outputs\":[{\"internalType\":\"contractIERC20Upgradeable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rollupManager\",\"outputs\":[{\"internalType\":\"contractPolygonRollupManager\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"transactions\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"forcedGlobalExitRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"forcedTimestamp\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"forcedBlockHashL1\",\"type\":\"bytes32\"}],\"internalType\":\"structPolygonRollupBaseEtrog.BatchData[]\",\"name\":\"batches\",\"type\":\"tuple[]\"},{\"internalType\":\"address\",\"name\":\"l2Coinbase\",\"type\":\"address\"}],\"name\":\"sequenceBatches\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes\",\"name\":\"transactions\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"forcedGlobalExitRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"forcedTimestamp\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"forcedBlockHashL1\",\"type\":\"bytes32\"}],\"internalType\":\"structPolygonRollupBaseEtrog.BatchData[]\",\"name\":\"batches\",\"type\":\"tuple[]\"}],\"name\":\"sequenceForceBatches\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newForceBatchAddress\",\"type\":\"address\"}],\"name\":\"setForceBatchAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"newforceBatchTimeout\",\"type\":\"uint64\"}],\"name\":\"setForceBatchTimeout\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newTrustedSequencer\",\"type\":\"address\"}],\"name\":\"setTrustedSequencer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"newTrustedSequencerURL\",\"type\":\"string\"}],\"name\":\"setTrustedSequencerURL\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newPendingAdmin\",\"type\":\"address\"}],\"name\":\"transferAdminRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"trustedSequencer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"trustedSequencerURL\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x6101006040523480156200001257600080fd5b5060405162004a6138038062004a61833981016040819052620000359162000071565b6001600160a01b0393841660a052918316608052821660c0521660e052620000d9565b6001600160a01b03811681146200006e57600080fd5b50565b600080600080608085870312156200008857600080fd5b8451620000958162000058565b6020860151909450620000a88162000058565b6040860151909350620000bb8162000058565b6060860151909250620000ce8162000058565b939692955090935050565b60805160a05160c05160e05161488a620001d7600039600081816105060152818161098f01528181610afc01528181610c7d01528181610fd4015281816112de015281816118c101528181611d65015281816121bc015281816122b20152818161295701528181612a1f01528181613342015281816133bb015281816133dd01526134f5015260008181610673015281816114d2015281816115ac015281816124830152818161258b0152612e7c01526000818161073701528181610e450152818161173501528181612b6a0152612efe0152600081816107690152818161083b0152818161220501528181612b3e015261348b015261488a6000f3fe608060405234801561001057600080fd5b50600436106102f35760003560e01c80637a5460c511610191578063c7fffd4b116100e3578063e46761c411610097578063ecef3f9911610071578063ecef3f99146107b2578063f35dda47146107c5578063f851a440146107cd57600080fd5b8063e46761c414610764578063e7a7ed021461078b578063eaeb077b1461079f57600080fd5b8063cfa8ed47116100c8578063cfa8ed4714610712578063d02103ca14610732578063d7bc90ff1461075957600080fd5b8063c7fffd4b146106f7578063c89e42df146106ff57600080fd5b8063a3c573eb11610145578063af7f3e021161011f578063af7f3e02146106bb578063b0afe154146106c3578063c754c7ed146106cf57600080fd5b8063a3c573eb1461066e578063a652f26c14610695578063ada8f919146106a857600080fd5b806391cafe321161017657806391cafe321461062d5780639e001877146106405780639f26f8401461065b57600080fd5b80637a5460c5146105e95780638c3d73011461062557600080fd5b8063456052671161024a5780635d6717a5116101fe5780636e05d2cd116101d85780636e05d2cd146105ba5780636ff512cc146105c357806371257022146105d657600080fd5b80635d6717a51461057f578063676870d2146105925780636b8616ce1461059a57600080fd5b80634e4877061161022f5780634e4877061461052857806352bdeb6d1461053b578063542028d51461057757600080fd5b806345605267146104c857806349b7b8021461050157600080fd5b806326782247116102ac5780633c351e10116102865780633c351e10146104135780633cbc795b1461043357806340b5de6c1461047057600080fd5b806326782247146103995780632c111c06146103de57806332c2d153146103fe57600080fd5b806305835f37116102dd57806305835f371461032e578063107bf28c1461037757806311e892d41461037f57600080fd5b8062d0295d146102f85780630350896314610313575b600080fd5b6103006107f3565b6040519081526020015b60405180910390f35b61031b602081565b60405161ffff909116815260200161030a565b61036a6040518060400160405280600881526020017f80808401c9c3809400000000000000000000000000000000000000000000000081525081565b60405161030a91906139f2565b61036a6108ff565b61038760f981565b60405160ff909116815260200161030a565b6001546103b99073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161030a565b6008546103b99073ffffffffffffffffffffffffffffffffffffffff1681565b61041161040c366004613a47565b61098d565b005b6009546103b99073ffffffffffffffffffffffffffffffffffffffff1681565b60095461045b9074010000000000000000000000000000000000000000900463ffffffff1681565b60405163ffffffff909116815260200161030a565b6104977fff0000000000000000000000000000000000000000000000000000000000000081565b6040517fff00000000000000000000000000000000000000000000000000000000000000909116815260200161030a565b6007546104e89068010000000000000000900467ffffffffffffffff1681565b60405167ffffffffffffffff909116815260200161030a565b6103b97f000000000000000000000000000000000000000000000000000000000000000081565b610411610536366004613a89565b610a5c565b61036a6040518060400160405280600281526020017f80b800000000000000000000000000000000000000000000000000000000000081525081565b61036a610c6e565b61041161058d366004613bc0565b610c7b565b61031b601f81565b6103006105a8366004613a89565b60066020526000908152604090205481565b61030060055481565b6104116105d1366004613c51565b611212565b6104116105e4366004613c80565b6112dc565b61036a6040518060400160405280600281526020017f80b900000000000000000000000000000000000000000000000000000000000081525081565b610411611afd565b61041161063b366004613c51565b611bd0565b6103b973a40d5f56745a118d0906a34e69aec8c0db1cb8fa81565b610411610669366004613d79565b611ce9565b6103b97f000000000000000000000000000000000000000000000000000000000000000081565b61036a6106a3366004613dbb565b612382565b6104116106b6366004613c51565b612767565b61036a612831565b6103006405ca1ab1e081565b6007546104e890700100000000000000000000000000000000900467ffffffffffffffff1681565b61038760e481565b61041161070d366004613e30565b61284d565b6002546103b99073ffffffffffffffffffffffffffffffffffffffff1681565b6103b97f000000000000000000000000000000000000000000000000000000000000000081565b610300635ca1ab1e81565b6103b97f000000000000000000000000000000000000000000000000000000000000000081565b6007546104e89067ffffffffffffffff1681565b6104116107ad366004613e65565b6128e0565b6104116107c0366004613edd565b612db1565b610387601b81565b6000546103b99062010000900473ffffffffffffffffffffffffffffffffffffffff1681565b6040517f70a08231000000000000000000000000000000000000000000000000000000008152306004820152600090819073ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016906370a0823190602401602060405180830381865afa158015610882573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906108a69190613f29565b6007549091506000906108d19067ffffffffffffffff68010000000000000000820481169116613f71565b67ffffffffffffffff169050806000036108ee5760009250505090565b6108f88183613f99565b9250505090565b6004805461090c90613fd4565b80601f016020809104026020016040519081016040528092919081815260200182805461093890613fd4565b80156109855780601f1061095a57610100808354040283529160200191610985565b820191906000526020600020905b81548152906001019060200180831161096857829003601f168201915b505050505081565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1633146109fc576040517fb9b3a2c800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8073ffffffffffffffffffffffffffffffffffffffff168367ffffffffffffffff167f9c72852172521097ba7e1482e6b44b351323df0155f97f4ea18fcec28e1f596684604051610a4f91815260200190565b60405180910390a3505050565b60005462010000900473ffffffffffffffffffffffffffffffffffffffff163314610ab3576040517f4755657900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b62093a8067ffffffffffffffff82161115610afa576040517ff5e37f2f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166315064c966040518163ffffffff1660e01b8152600401602060405180830381865afa158015610b65573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610b899190614027565b610bea5760075467ffffffffffffffff700100000000000000000000000000000000909104811690821610610bea576040517ff5e37f2f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600780547fffffffffffffffff0000000000000000ffffffffffffffffffffffffffffffff1670010000000000000000000000000000000067ffffffffffffffff8416908102919091179091556040519081527fa7eb6cb8a613eb4e8bddc1ac3d61ec6cf10898760f0b187bcca794c6ca6fa40b906020015b60405180910390a150565b6003805461090c90613fd4565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff163314610cea576040517fb9b3a2c800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600054610100900460ff1615808015610d0a5750600054600160ff909116105b80610d245750303b158015610d24575060005460ff166001145b610db5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201527f647920696e697469616c697a656400000000000000000000000000000000000060648201526084015b60405180910390fd5b600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011790558015610e1357600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff166101001790555b60006040518060a00160405280606281526020016147f3606291399050600081805190602001209050600042905060007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16633ed691ef6040518163ffffffff1660e01b8152600401602060405180830381865afa158015610eae573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610ed29190613f29565b90506000868483858d610ee6600143614049565b60408051602081019790975286019490945260608086019390935260c09190911b7fffffffffffffffff000000000000000000000000000000000000000000000000166080850152901b7fffffffffffffffffffffffffffffffffffffffff00000000000000000000000016608883015240609c82015260bc01604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe001815290829052805160209091012060058190557f9a908e73000000000000000000000000000000000000000000000000000000008252600160048301526024820181905291506000907f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1690639a908e73906044016020604051808303816000875af1158015611032573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906110569190614062565b90508b600060026101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508a600260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555089600390816110e991906140c5565b5060046110f68a826140c5565b508b600860006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555062069780600760106101000a81548167ffffffffffffffff021916908367ffffffffffffffff1602179055507fd2c80353fc15ef62c6affc7cd6b7ab5b42c43290c50be3372e55ae552cecd19c8187858e60405161119994939291906141df565b60405180910390a1505050505050801561120a57600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b505050505050565b60005462010000900473ffffffffffffffffffffffffffffffffffffffff163314611269576040517f4755657900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600280547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83169081179091556040519081527ff54144f9611984021529f814a1cb6a41e22c58351510a0d9f7e822618abb9cc090602001610c63565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16331461134b576040517fb9b3a2c800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600054610100900460ff161580801561136b5750600054600160ff909116105b806113855750303b158015611385575060005460ff166001145b611411576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201527f647920696e697469616c697a65640000000000000000000000000000000000006064820152608401610dac565b600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00166001179055801561146f57600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff166101001790555b606073ffffffffffffffffffffffffffffffffffffffff8516156116d6576040517fc00f14ab00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff86811660048301527f0000000000000000000000000000000000000000000000000000000000000000169063c00f14ab90602401600060405180830381865afa158015611519573d6000803e3d6000fd5b505050506040513d6000823e601f3d9081017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016820160405261155f919081019061422f565b6040517f318aee3d00000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff878116600483015291925060009182917f00000000000000000000000000000000000000000000000000000000000000009091169063318aee3d9060240160408051808303816000875af11580156115f6573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061161a91906142a6565b915091508163ffffffff16600014611692576009805463ffffffff841674010000000000000000000000000000000000000000027fffffffffffffffff00000000000000000000000000000000000000000000000090911673ffffffffffffffffffffffffffffffffffffffff8416171790556116d3565b600980547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff89161790555b50505b60095460009061171e90889073ffffffffffffffffffffffffffffffffffffffff81169074010000000000000000000000000000000000000000900463ffffffff1685612382565b9050600081805190602001209050600042905060007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16633ed691ef6040518163ffffffff1660e01b8152600401602060405180830381865afa15801561179e573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906117c29190613f29565b90506000808483858f6117d6600143614049565b60408051602081019790975286019490945260608086019390935260c09190911b7fffffffffffffffff000000000000000000000000000000000000000000000000166080850152901b7fffffffffffffffffffffffffffffffffffffffff00000000000000000000000016608883015240609c82015260bc01604080518083037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe001815290829052805160209091012060058190557f9a908e73000000000000000000000000000000000000000000000000000000008252600160048301526024820181905291507f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1690639a908e73906044016020604051808303816000875af115801561191f573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906119439190614062565b508c600060026101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508b600260006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555088600390816119d591906140c5565b5060046119e289826140c5565b508c600860006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555062069780600760106101000a81548167ffffffffffffffff021916908367ffffffffffffffff1602179055507f060116213bcbf54ca19fd649dc84b59ab2bbd200ab199770e4d923e222a28e7f85838e604051611a83939291906142e0565b60405180910390a15050505050508015611af457600080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b50505050505050565b60015473ffffffffffffffffffffffffffffffffffffffff163314611b4e576040517fd1ec4b2300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600154600080547fffffffffffffffffffff0000000000000000000000000000000000000000ffff1673ffffffffffffffffffffffffffffffffffffffff9092166201000081029290921790556040519081527f056dc487bbf0795d0bbb1b4f0af523a855503cff740bfb4d5475f7a90c091e8e9060200160405180910390a1565b60005462010000900473ffffffffffffffffffffffffffffffffffffffff163314611c27576040517f4755657900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60085473ffffffffffffffffffffffffffffffffffffffff16611c76576040517fc89374d800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600880547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83169081179091556040519081527f5fbd7dd171301c4a1611a84aac4ba86d119478560557755f7927595b082634fb90602001610c63565b60085473ffffffffffffffffffffffffffffffffffffffff168015801590611d27575073ffffffffffffffffffffffffffffffffffffffff81163314155b15611d5e576040517f24eff8c300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b4262093a807f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166330c27dde6040518163ffffffff1660e01b8152600401602060405180830381865afa158015611dce573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611df29190614062565b611dfc919061431f565b67ffffffffffffffff161115611e3e576040517f3d49ed4c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b816000819003611e7a576040517fcb591a5f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6103e8811115611eb6576040517fb59f753a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60075467ffffffffffffffff80821691611ede91849168010000000000000000900416614340565b1115611f16576040517fc630a00d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6007546005546801000000000000000090910467ffffffffffffffff169060005b838110156121b6576000878783818110611f5357611f53614353565b9050602002810190611f659190614382565b611f6e906143c0565b905083611f7a81614449565b825180516020918201208185015160408087015160608801519151959a50929550600094611fe7948794929101938452602084019290925260c01b7fffffffffffffffff000000000000000000000000000000000000000000000000166040830152604882015260680190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152918152815160209283012067ffffffffffffffff8916600090815260069093529120549091508114612070576040517fce3d755e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff8616600090815260066020526040812055612095600188614049565b84036121045742600760109054906101000a900467ffffffffffffffff1684604001516120c2919061431f565b67ffffffffffffffff161115612104576040517fc44a082100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60208381015160408086015160608088015183519586018b90529285018790528481019390935260c01b7fffffffffffffffff0000000000000000000000000000000000000000000000001660808401523390911b7fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166088830152609c82015260bc0160405160208183030381529060405280519060200120945050505080806121ae90614470565b915050611f37565b5061222c7f0000000000000000000000000000000000000000000000000000000000000000846121e46107f3565b6121ee91906144a8565b73ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001691906135ca565b60058190556007805467ffffffffffffffff841668010000000000000000027fffffffffffffffffffffffffffffffff0000000000000000ffffffffffffffff9091161790556040517f9a908e7300000000000000000000000000000000000000000000000000000000815260009073ffffffffffffffffffffffffffffffffffffffff7f00000000000000000000000000000000000000000000000000000000000000001690639a908e73906122fe908790869060040167ffffffffffffffff929092168252602082015260400190565b6020604051808303816000875af115801561231d573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906123419190614062565b60405190915067ffffffffffffffff8216907f648a61dd2438f072f5a1960939abd30f37aea80d2e94c9792ad142d3e0a490a490600090a250505050505050565b6060600085858573a40d5f56745a118d0906a34e69aec8c0db1cb8fa6000876040516024016123b6969594939291906144bf565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167ff811bff70000000000000000000000000000000000000000000000000000000017905283519091506060906000036125075760f9601f835161244b9190614522565b6040518060400160405280600881526020017f80808401c9c380940000000000000000000000000000000000000000000000008152507f00000000000000000000000000000000000000000000000000000000000000006040518060400160405280600281526020017f80b800000000000000000000000000000000000000000000000000000000000081525060e4876040516020016124f1979695949392919061453d565b604051602081830303815290604052905061260b565b815161ffff1015612544576040517f248b8f8200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b815160f9612553602083614522565b6040518060400160405280600881526020017f80808401c9c380940000000000000000000000000000000000000000000000008152507f00000000000000000000000000000000000000000000000000000000000000006040518060400160405280600281526020017f80b900000000000000000000000000000000000000000000000000000000000081525085886040516020016125f89796959493929190614620565b6040516020818303038152906040529150505b805160208083019190912060408051600080825293810180835292909252601b908201526405ca1ab1e06060820152635ca1ab1e608082015260019060a0016020604051602081039080840390855afa15801561266c573d6000803e3d6000fd5b50506040517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0015191505073ffffffffffffffffffffffffffffffffffffffff81166126e4576040517fcd16196600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60405160009061272a9084906405ca1ab1e090635ca1ab1e90601b907fff0000000000000000000000000000000000000000000000000000000000000090602001614703565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190529450505050505b949350505050565b60005462010000900473ffffffffffffffffffffffffffffffffffffffff1633146127be576040517f4755657900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600180547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff83169081179091556040519081527fa5b56b7906fd0a20e3f35120dd8343db1e12e037a6c90111c7e42885e82a1ce690602001610c63565b6040518060a00160405280606281526020016147f36062913981565b60005462010000900473ffffffffffffffffffffffffffffffffffffffff1633146128a4576040517f4755657900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60036128b082826140c5565b507f6b8f723a4c7a5335cafae8a598a0aa0301be1387c037dccc085b62add6448b2081604051610c6391906139f2565b60085473ffffffffffffffffffffffffffffffffffffffff16801580159061291e575073ffffffffffffffffffffffffffffffffffffffff81163314155b15612955576040517f24eff8c300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166315064c966040518163ffffffff1660e01b8152600401602060405180830381865afa1580156129c0573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906129e49190614027565b15612a1b576040517f39258d1800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663604691696040518163ffffffff1660e01b8152600401602060405180830381865afa158015612a88573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612aac9190613f29565b905082811115612ae8576040517f2354600f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b611388841115612b24576040517fa29a6c7c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b612b6673ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000163330846136a3565b60007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16633ed691ef6040518163ffffffff1660e01b8152600401602060405180830381865afa158015612bd3573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612bf79190613f29565b6007805491925067ffffffffffffffff909116906000612c1683614449565b91906101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550508585604051612c4d92919061475f565b6040519081900390208142612c63600143614049565b60408051602081019590955284019290925260c01b7fffffffffffffffff000000000000000000000000000000000000000000000000166060830152406068820152608801604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152918152815160209283012060075467ffffffffffffffff1660009081526006909352912055323303612d5b57600754604080518381523360208201526060818301819052600090820152905167ffffffffffffffff909216917ff94bb37db835f1ab585ee00041849a09b12cd081d77fa15ca070757619cbc9319181900360800190a261120a565b60075460405167ffffffffffffffff909116907ff94bb37db835f1ab585ee00041849a09b12cd081d77fa15ca070757619cbc93190612da190849033908b908b9061476f565b60405180910390a2505050505050565b60025473ffffffffffffffffffffffffffffffffffffffff163314612e02576040517f11e7be1500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b816000819003612e3e576040517fcb591a5f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6103e8811115612e7a576040517fb59f753a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166379e2cf976040518163ffffffff1660e01b8152600401600060405180830381600087803b158015612ee257600080fd5b505af1158015612ef6573d6000803e3d6000fd5b5050505060007f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16635ca1e1656040518163ffffffff1660e01b8152600401602060405180830381865afa158015612f67573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190612f8b9190613f29565b60075460055491925042916801000000000000000090910467ffffffffffffffff16908160005b868110156132b45760008a8a83818110612fce57612fce614353565b9050602002810190612fe09190614382565b612fe9906143c0565b8051805160209091012060408201519192509067ffffffffffffffff16156131ce578561301581614449565b9650506000818360200151846040015185606001516040516020016130789493929190938452602084019290925260c01b7fffffffffffffffff000000000000000000000000000000000000000000000000166040830152604882015260680190565b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0818403018152918152815160209283012067ffffffffffffffff8a16600090815260069093529120549091508114613101576040517fce3d755e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60208381015160408086015160608088015183519586018c90529285018790528481019390935260c01b7fffffffffffffffff000000000000000000000000000000000000000000000000166080840152908d901b7fffffffffffffffffffffffffffffffffffffffff000000000000000000000000166088830152609c82015260bc01604051602081830303815290604052805190602001209550600660008867ffffffffffffffff1667ffffffffffffffff168152602001908152602001600020600090555061329f565b8151516201d4c0101561320d576040517fa29a6c7c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b604080516020810187905290810182905260608082018a905260c089901b7fffffffffffffffff0000000000000000000000000000000000000000000000001660808301528b901b7fffffffffffffffffffffffffffffffffffffffff0000000000000000000000001660888201526000609c82015260bc016040516020818303038152906040528051906020012094505b505080806132ac90614470565b915050612fb2565b5060075467ffffffffffffffff90811690841611156132ff576040517fc630a00d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60058290558567ffffffffffffffff848116908316146133b55760006133258386613f71565b905061333b67ffffffffffffffff821683614049565b91506133747f00000000000000000000000000000000000000000000000000000000000000008267ffffffffffffffff166121e46107f3565b50600780547fffffffffffffffffffffffffffffffff0000000000000000ffffffffffffffff166801000000000000000067ffffffffffffffff8716021790555b6134b3337f0000000000000000000000000000000000000000000000000000000000000000837f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1663477fa2706040518163ffffffff1660e01b8152600401602060405180830381865afa158015613446573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061346a9190613f29565b61347491906144a8565b73ffffffffffffffffffffffffffffffffffffffff7f0000000000000000000000000000000000000000000000000000000000000000169291906136a3565b6040517f9a908e7300000000000000000000000000000000000000000000000000000000815267ffffffffffffffff88166004820152602481018490526000907f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1690639a908e73906044016020604051808303816000875af1158015613553573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906135779190614062565b90508067ffffffffffffffff167f3e54d0825ed78523037d00a81759237eb436ce774bd546993ee67a1b67b6e766886040516135b591815260200190565b60405180910390a25050505050505050505050565b60405173ffffffffffffffffffffffffffffffffffffffff831660248201526044810182905261369e9084907fa9059cbb00000000000000000000000000000000000000000000000000000000906064015b604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08184030181529190526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fffffffff0000000000000000000000000000000000000000000000000000000090931692909217909152613707565b505050565b60405173ffffffffffffffffffffffffffffffffffffffff808516602483015283166044820152606481018290526137019085907f23b872dd000000000000000000000000000000000000000000000000000000009060840161361c565b50505050565b6000613769826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65648152508573ffffffffffffffffffffffffffffffffffffffff166138139092919063ffffffff16565b80519091501561369e57808060200190518101906137879190614027565b61369e576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f742073756363656564000000000000000000000000000000000000000000006064820152608401610dac565b606061275f8484600085856000808673ffffffffffffffffffffffffffffffffffffffff16858760405161384791906147e0565b60006040518083038185875af1925050503d8060008114613884576040519150601f19603f3d011682016040523d82523d6000602084013e613889565b606091505b509150915061389a878383876138a5565b979650505050505050565b6060831561393b5782516000036139345773ffffffffffffffffffffffffffffffffffffffff85163b613934576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e74726163740000006044820152606401610dac565b508161275f565b61275f83838151156139505781518083602001fd5b806040517f08c379a0000000000000000000000000000000000000000000000000000000008152600401610dac91906139f2565b60005b8381101561399f578181015183820152602001613987565b50506000910152565b600081518084526139c0816020860160208601613984565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b602081526000613a0560208301846139a8565b9392505050565b67ffffffffffffffff81168114613a2257600080fd5b50565b73ffffffffffffffffffffffffffffffffffffffff81168114613a2257600080fd5b600080600060608486031215613a5c57600080fd5b8335613a6781613a0c565b9250602084013591506040840135613a7e81613a25565b809150509250925092565b600060208284031215613a9b57600080fd5b8135613a0581613a0c565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016810167ffffffffffffffff81118282101715613b1c57613b1c613aa6565b604052919050565b600067ffffffffffffffff821115613b3e57613b3e613aa6565b50601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01660200190565b600082601f830112613b7b57600080fd5b8135613b8e613b8982613b24565b613ad5565b818152846020838601011115613ba357600080fd5b816020850160208301376000918101602001919091529392505050565b600080600080600060a08688031215613bd857600080fd5b8535613be381613a25565b94506020860135613bf381613a25565b9350604086013567ffffffffffffffff80821115613c1057600080fd5b613c1c89838a01613b6a565b94506060880135915080821115613c3257600080fd5b50613c3f88828901613b6a565b95989497509295608001359392505050565b600060208284031215613c6357600080fd5b8135613a0581613a25565b63ffffffff81168114613a2257600080fd5b60008060008060008060c08789031215613c9957600080fd5b8635613ca481613a25565b95506020870135613cb481613a25565b94506040870135613cc481613c6e565b93506060870135613cd481613a25565b9250608087013567ffffffffffffffff80821115613cf157600080fd5b613cfd8a838b01613b6a565b935060a0890135915080821115613d1357600080fd5b50613d2089828a01613b6a565b9150509295509295509295565b60008083601f840112613d3f57600080fd5b50813567ffffffffffffffff811115613d5757600080fd5b6020830191508360208260051b8501011115613d7257600080fd5b9250929050565b60008060208385031215613d8c57600080fd5b823567ffffffffffffffff811115613da357600080fd5b613daf85828601613d2d565b90969095509350505050565b60008060008060808587031215613dd157600080fd5b8435613ddc81613c6e565b93506020850135613dec81613a25565b92506040850135613dfc81613c6e565b9150606085013567ffffffffffffffff811115613e1857600080fd5b613e2487828801613b6a565b91505092959194509250565b600060208284031215613e4257600080fd5b813567ffffffffffffffff811115613e5957600080fd5b61275f84828501613b6a565b600080600060408486031215613e7a57600080fd5b833567ffffffffffffffff80821115613e9257600080fd5b818601915086601f830112613ea657600080fd5b813581811115613eb557600080fd5b876020828501011115613ec757600080fd5b6020928301989097509590910135949350505050565b600080600060408486031215613ef257600080fd5b833567ffffffffffffffff811115613f0957600080fd5b613f1586828701613d2d565b9094509250506020840135613a7e81613a25565b600060208284031215613f3b57600080fd5b5051919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b67ffffffffffffffff828116828216039080821115613f9257613f92613f42565b5092915050565b600082613fcf577f4e487b7100000000000000000000000000000000000000000000000000000000600052601260045260246000fd5b500490565b600181811c90821680613fe857607f821691505b602082108103614021577f4e487b7100000000000000000000000000000000000000000000000000000000600052602260045260246000fd5b50919050565b60006020828403121561403957600080fd5b81518015158114613a0557600080fd5b8181038181111561405c5761405c613f42565b92915050565b60006020828403121561407457600080fd5b8151613a0581613a0c565b601f82111561369e57600081815260208120601f850160051c810160208610156140a65750805b601f850160051c820191505b8181101561120a578281556001016140b2565b815167ffffffffffffffff8111156140df576140df613aa6565b6140f3816140ed8454613fd4565b8461407f565b602080601f83116001811461414657600084156141105750858301515b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600386901b1c1916600185901b17855561120a565b6000858152602081207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08616915b8281101561419357888601518255948401946001909101908401614174565b50858210156141cf57878501517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff600388901b60f8161c191681555b5050505050600190811b01905550565b67ffffffffffffffff8516815260806020820152600061420260808301866139a8565b905083604083015273ffffffffffffffffffffffffffffffffffffffff8316606083015295945050505050565b60006020828403121561424157600080fd5b815167ffffffffffffffff81111561425857600080fd5b8201601f8101841361426957600080fd5b8051614277613b8982613b24565b81815285602083850101111561428c57600080fd5b61429d826020830160208601613984565b95945050505050565b600080604083850312156142b957600080fd5b82516142c481613c6e565b60208401519092506142d581613a25565b809150509250929050565b6060815260006142f360608301866139a8565b905083602083015273ffffffffffffffffffffffffffffffffffffffff83166040830152949350505050565b67ffffffffffffffff818116838216019080821115613f9257613f92613f42565b8082018082111561405c5761405c613f42565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600082357fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff818336030181126143b657600080fd5b9190910192915050565b6000608082360312156143d257600080fd5b6040516080810167ffffffffffffffff82821081831117156143f6576143f6613aa6565b81604052843591508082111561440b57600080fd5b5061441836828601613b6a565b82525060208301356020820152604083013561443381613a0c565b6040820152606092830135928101929092525090565b600067ffffffffffffffff80831681810361446657614466613f42565b6001019392505050565b60007fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036144a1576144a1613f42565b5060010190565b808202811582820484141761405c5761405c613f42565b600063ffffffff808916835273ffffffffffffffffffffffffffffffffffffffff8089166020850152818816604085015280871660608501528086166080850152505060c060a083015261451660c08301846139a8565b98975050505050505050565b61ffff818116838216019080821115613f9257613f92613f42565b60007fff00000000000000000000000000000000000000000000000000000000000000808a60f81b1683527fffff0000000000000000000000000000000000000000000000000000000000008960f01b16600184015287516145a6816003860160208c01613984565b80840190507fffffffffffffffffffffffffffffffffffffffff0000000000000000000000008860601b16600382015286516145e9816017840160208b01613984565b808201915050818660f81b1660178201528451915061460f826018830160208801613984565b016018019998505050505050505050565b7fff000000000000000000000000000000000000000000000000000000000000008860f81b16815260007fffff000000000000000000000000000000000000000000000000000000000000808960f01b1660018401528751614689816003860160208c01613984565b80840190507fffffffffffffffffffffffffffffffffffffffff0000000000000000000000008860601b16600382015286516146cc816017840160208b01613984565b808201915050818660f01b166017820152845191506146f2826019830160208801613984565b016019019998505050505050505050565b60008651614715818460208b01613984565b9190910194855250602084019290925260f81b7fff000000000000000000000000000000000000000000000000000000000000009081166040840152166041820152604201919050565b8183823760009101908152919050565b84815273ffffffffffffffffffffffffffffffffffffffff8416602082015260606040820152816060820152818360808301376000818301608090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe01601019392505050565b600082516143b681846020870161398456fedf2a8080944d5cf5032b2a844602278b01199ed191a86c93ff8080821092808000000000000000000000000000000000000000000000000000000005ca1ab1e000000000000000000000000000000000000000000000000000000005ca1ab1e01bffa26469706673582212201221389ded8ea187a66f83d3bd052755e28647dbf3bc616c9e91e0a8b7ecf74364736f6c63430008140033",
}

// ElderberrypolygonzkevmABI is the input ABI used to generate the binding from.
// Deprecated: Use ElderberrypolygonzkevmMetaData.ABI instead.
var ElderberrypolygonzkevmABI = ElderberrypolygonzkevmMetaData.ABI

// ElderberrypolygonzkevmBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ElderberrypolygonzkevmMetaData.Bin instead.
var ElderberrypolygonzkevmBin = ElderberrypolygonzkevmMetaData.Bin

// DeployElderberrypolygonzkevm deploys a new Ethereum contract, binding an instance of Elderberrypolygonzkevm to it.
func DeployElderberrypolygonzkevm(auth *bind.TransactOpts, backend bind.ContractBackend, _globalExitRootManager common.Address, _pol common.Address, _bridgeAddress common.Address, _rollupManager common.Address) (common.Address, *types.Transaction, *Elderberrypolygonzkevm, error) {
	parsed, err := ElderberrypolygonzkevmMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ElderberrypolygonzkevmBin), backend, _globalExitRootManager, _pol, _bridgeAddress, _rollupManager)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Elderberrypolygonzkevm{ElderberrypolygonzkevmCaller: ElderberrypolygonzkevmCaller{contract: contract}, ElderberrypolygonzkevmTransactor: ElderberrypolygonzkevmTransactor{contract: contract}, ElderberrypolygonzkevmFilterer: ElderberrypolygonzkevmFilterer{contract: contract}}, nil
}

// Elderberrypolygonzkevm is an auto generated Go binding around an Ethereum contract.
type Elderberrypolygonzkevm struct {
	ElderberrypolygonzkevmCaller     // Read-only binding to the contract
	ElderberrypolygonzkevmTransactor // Write-only binding to the contract
	ElderberrypolygonzkevmFilterer   // Log filterer for contract events
}

// ElderberrypolygonzkevmCaller is an auto generated read-only Go binding around an Ethereum contract.
type ElderberrypolygonzkevmCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ElderberrypolygonzkevmTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ElderberrypolygonzkevmTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ElderberrypolygonzkevmFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ElderberrypolygonzkevmFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ElderberrypolygonzkevmSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ElderberrypolygonzkevmSession struct {
	Contract     *Elderberrypolygonzkevm // Generic contract binding to set the session for
	CallOpts     bind.CallOpts           // Call options to use throughout this session
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// ElderberrypolygonzkevmCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ElderberrypolygonzkevmCallerSession struct {
	Contract *ElderberrypolygonzkevmCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                 // Call options to use throughout this session
}

// ElderberrypolygonzkevmTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ElderberrypolygonzkevmTransactorSession struct {
	Contract     *ElderberrypolygonzkevmTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                 // Transaction auth options to use throughout this session
}

// ElderberrypolygonzkevmRaw is an auto generated low-level Go binding around an Ethereum contract.
type ElderberrypolygonzkevmRaw struct {
	Contract *Elderberrypolygonzkevm // Generic contract binding to access the raw methods on
}

// ElderberrypolygonzkevmCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ElderberrypolygonzkevmCallerRaw struct {
	Contract *ElderberrypolygonzkevmCaller // Generic read-only contract binding to access the raw methods on
}

// ElderberrypolygonzkevmTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ElderberrypolygonzkevmTransactorRaw struct {
	Contract *ElderberrypolygonzkevmTransactor // Generic write-only contract binding to access the raw methods on
}

// NewElderberrypolygonzkevm creates a new instance of Elderberrypolygonzkevm, bound to a specific deployed contract.
func NewElderberrypolygonzkevm(address common.Address, backend bind.ContractBackend) (*Elderberrypolygonzkevm, error) {
	contract, err := bindElderberrypolygonzkevm(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Elderberrypolygonzkevm{ElderberrypolygonzkevmCaller: ElderberrypolygonzkevmCaller{contract: contract}, ElderberrypolygonzkevmTransactor: ElderberrypolygonzkevmTransactor{contract: contract}, ElderberrypolygonzkevmFilterer: ElderberrypolygonzkevmFilterer{contract: contract}}, nil
}

// NewElderberrypolygonzkevmCaller creates a new read-only instance of Elderberrypolygonzkevm, bound to a specific deployed contract.
func NewElderberrypolygonzkevmCaller(address common.Address, caller bind.ContractCaller) (*ElderberrypolygonzkevmCaller, error) {
	contract, err := bindElderberrypolygonzkevm(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ElderberrypolygonzkevmCaller{contract: contract}, nil
}

// NewElderberrypolygonzkevmTransactor creates a new write-only instance of Elderberrypolygonzkevm, bound to a specific deployed contract.
func NewElderberrypolygonzkevmTransactor(address common.Address, transactor bind.ContractTransactor) (*ElderberrypolygonzkevmTransactor, error) {
	contract, err := bindElderberrypolygonzkevm(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ElderberrypolygonzkevmTransactor{contract: contract}, nil
}

// NewElderberrypolygonzkevmFilterer creates a new log filterer instance of Elderberrypolygonzkevm, bound to a specific deployed contract.
func NewElderberrypolygonzkevmFilterer(address common.Address, filterer bind.ContractFilterer) (*ElderberrypolygonzkevmFilterer, error) {
	contract, err := bindElderberrypolygonzkevm(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ElderberrypolygonzkevmFilterer{contract: contract}, nil
}

// bindElderberrypolygonzkevm binds a generic wrapper to an already deployed contract.
func bindElderberrypolygonzkevm(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := ElderberrypolygonzkevmMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Elderberrypolygonzkevm.Contract.ElderberrypolygonzkevmCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Elderberrypolygonzkevm.Contract.ElderberrypolygonzkevmTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Elderberrypolygonzkevm.Contract.ElderberrypolygonzkevmTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Elderberrypolygonzkevm.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Elderberrypolygonzkevm.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Elderberrypolygonzkevm.Contract.contract.Transact(opts, method, params...)
}

// GLOBALEXITROOTMANAGERL2 is a free data retrieval call binding the contract method 0x9e001877.
//
// Solidity: function GLOBAL_EXIT_ROOT_MANAGER_L2() view returns(address)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCaller) GLOBALEXITROOTMANAGERL2(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Elderberrypolygonzkevm.contract.Call(opts, &out, "GLOBAL_EXIT_ROOT_MANAGER_L2")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GLOBALEXITROOTMANAGERL2 is a free data retrieval call binding the contract method 0x9e001877.
//
// Solidity: function GLOBAL_EXIT_ROOT_MANAGER_L2() view returns(address)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmSession) GLOBALEXITROOTMANAGERL2() (common.Address, error) {
	return _Elderberrypolygonzkevm.Contract.GLOBALEXITROOTMANAGERL2(&_Elderberrypolygonzkevm.CallOpts)
}

// GLOBALEXITROOTMANAGERL2 is a free data retrieval call binding the contract method 0x9e001877.
//
// Solidity: function GLOBAL_EXIT_ROOT_MANAGER_L2() view returns(address)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCallerSession) GLOBALEXITROOTMANAGERL2() (common.Address, error) {
	return _Elderberrypolygonzkevm.Contract.GLOBALEXITROOTMANAGERL2(&_Elderberrypolygonzkevm.CallOpts)
}

// INITIALIZETXBRIDGELISTLENLEN is a free data retrieval call binding the contract method 0x11e892d4.
//
// Solidity: function INITIALIZE_TX_BRIDGE_LIST_LEN_LEN() view returns(uint8)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCaller) INITIALIZETXBRIDGELISTLENLEN(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Elderberrypolygonzkevm.contract.Call(opts, &out, "INITIALIZE_TX_BRIDGE_LIST_LEN_LEN")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// INITIALIZETXBRIDGELISTLENLEN is a free data retrieval call binding the contract method 0x11e892d4.
//
// Solidity: function INITIALIZE_TX_BRIDGE_LIST_LEN_LEN() view returns(uint8)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmSession) INITIALIZETXBRIDGELISTLENLEN() (uint8, error) {
	return _Elderberrypolygonzkevm.Contract.INITIALIZETXBRIDGELISTLENLEN(&_Elderberrypolygonzkevm.CallOpts)
}

// INITIALIZETXBRIDGELISTLENLEN is a free data retrieval call binding the contract method 0x11e892d4.
//
// Solidity: function INITIALIZE_TX_BRIDGE_LIST_LEN_LEN() view returns(uint8)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCallerSession) INITIALIZETXBRIDGELISTLENLEN() (uint8, error) {
	return _Elderberrypolygonzkevm.Contract.INITIALIZETXBRIDGELISTLENLEN(&_Elderberrypolygonzkevm.CallOpts)
}

// INITIALIZETXBRIDGEPARAMS is a free data retrieval call binding the contract method 0x05835f37.
//
// Solidity: function INITIALIZE_TX_BRIDGE_PARAMS() view returns(bytes)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCaller) INITIALIZETXBRIDGEPARAMS(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _Elderberrypolygonzkevm.contract.Call(opts, &out, "INITIALIZE_TX_BRIDGE_PARAMS")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// INITIALIZETXBRIDGEPARAMS is a free data retrieval call binding the contract method 0x05835f37.
//
// Solidity: function INITIALIZE_TX_BRIDGE_PARAMS() view returns(bytes)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmSession) INITIALIZETXBRIDGEPARAMS() ([]byte, error) {
	return _Elderberrypolygonzkevm.Contract.INITIALIZETXBRIDGEPARAMS(&_Elderberrypolygonzkevm.CallOpts)
}

// INITIALIZETXBRIDGEPARAMS is a free data retrieval call binding the contract method 0x05835f37.
//
// Solidity: function INITIALIZE_TX_BRIDGE_PARAMS() view returns(bytes)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCallerSession) INITIALIZETXBRIDGEPARAMS() ([]byte, error) {
	return _Elderberrypolygonzkevm.Contract.INITIALIZETXBRIDGEPARAMS(&_Elderberrypolygonzkevm.CallOpts)
}

// INITIALIZETXBRIDGEPARAMSAFTERBRIDGEADDRESS is a free data retrieval call binding the contract method 0x7a5460c5.
//
// Solidity: function INITIALIZE_TX_BRIDGE_PARAMS_AFTER_BRIDGE_ADDRESS() view returns(bytes)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCaller) INITIALIZETXBRIDGEPARAMSAFTERBRIDGEADDRESS(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _Elderberrypolygonzkevm.contract.Call(opts, &out, "INITIALIZE_TX_BRIDGE_PARAMS_AFTER_BRIDGE_ADDRESS")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// INITIALIZETXBRIDGEPARAMSAFTERBRIDGEADDRESS is a free data retrieval call binding the contract method 0x7a5460c5.
//
// Solidity: function INITIALIZE_TX_BRIDGE_PARAMS_AFTER_BRIDGE_ADDRESS() view returns(bytes)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmSession) INITIALIZETXBRIDGEPARAMSAFTERBRIDGEADDRESS() ([]byte, error) {
	return _Elderberrypolygonzkevm.Contract.INITIALIZETXBRIDGEPARAMSAFTERBRIDGEADDRESS(&_Elderberrypolygonzkevm.CallOpts)
}

// INITIALIZETXBRIDGEPARAMSAFTERBRIDGEADDRESS is a free data retrieval call binding the contract method 0x7a5460c5.
//
// Solidity: function INITIALIZE_TX_BRIDGE_PARAMS_AFTER_BRIDGE_ADDRESS() view returns(bytes)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCallerSession) INITIALIZETXBRIDGEPARAMSAFTERBRIDGEADDRESS() ([]byte, error) {
	return _Elderberrypolygonzkevm.Contract.INITIALIZETXBRIDGEPARAMSAFTERBRIDGEADDRESS(&_Elderberrypolygonzkevm.CallOpts)
}

// INITIALIZETXBRIDGEPARAMSAFTERBRIDGEADDRESSEMPTYMETADATA is a free data retrieval call binding the contract method 0x52bdeb6d.
//
// Solidity: function INITIALIZE_TX_BRIDGE_PARAMS_AFTER_BRIDGE_ADDRESS_EMPTY_METADATA() view returns(bytes)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCaller) INITIALIZETXBRIDGEPARAMSAFTERBRIDGEADDRESSEMPTYMETADATA(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _Elderberrypolygonzkevm.contract.Call(opts, &out, "INITIALIZE_TX_BRIDGE_PARAMS_AFTER_BRIDGE_ADDRESS_EMPTY_METADATA")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// INITIALIZETXBRIDGEPARAMSAFTERBRIDGEADDRESSEMPTYMETADATA is a free data retrieval call binding the contract method 0x52bdeb6d.
//
// Solidity: function INITIALIZE_TX_BRIDGE_PARAMS_AFTER_BRIDGE_ADDRESS_EMPTY_METADATA() view returns(bytes)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmSession) INITIALIZETXBRIDGEPARAMSAFTERBRIDGEADDRESSEMPTYMETADATA() ([]byte, error) {
	return _Elderberrypolygonzkevm.Contract.INITIALIZETXBRIDGEPARAMSAFTERBRIDGEADDRESSEMPTYMETADATA(&_Elderberrypolygonzkevm.CallOpts)
}

// INITIALIZETXBRIDGEPARAMSAFTERBRIDGEADDRESSEMPTYMETADATA is a free data retrieval call binding the contract method 0x52bdeb6d.
//
// Solidity: function INITIALIZE_TX_BRIDGE_PARAMS_AFTER_BRIDGE_ADDRESS_EMPTY_METADATA() view returns(bytes)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCallerSession) INITIALIZETXBRIDGEPARAMSAFTERBRIDGEADDRESSEMPTYMETADATA() ([]byte, error) {
	return _Elderberrypolygonzkevm.Contract.INITIALIZETXBRIDGEPARAMSAFTERBRIDGEADDRESSEMPTYMETADATA(&_Elderberrypolygonzkevm.CallOpts)
}

// INITIALIZETXCONSTANTBYTES is a free data retrieval call binding the contract method 0x03508963.
//
// Solidity: function INITIALIZE_TX_CONSTANT_BYTES() view returns(uint16)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCaller) INITIALIZETXCONSTANTBYTES(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _Elderberrypolygonzkevm.contract.Call(opts, &out, "INITIALIZE_TX_CONSTANT_BYTES")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// INITIALIZETXCONSTANTBYTES is a free data retrieval call binding the contract method 0x03508963.
//
// Solidity: function INITIALIZE_TX_CONSTANT_BYTES() view returns(uint16)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmSession) INITIALIZETXCONSTANTBYTES() (uint16, error) {
	return _Elderberrypolygonzkevm.Contract.INITIALIZETXCONSTANTBYTES(&_Elderberrypolygonzkevm.CallOpts)
}

// INITIALIZETXCONSTANTBYTES is a free data retrieval call binding the contract method 0x03508963.
//
// Solidity: function INITIALIZE_TX_CONSTANT_BYTES() view returns(uint16)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCallerSession) INITIALIZETXCONSTANTBYTES() (uint16, error) {
	return _Elderberrypolygonzkevm.Contract.INITIALIZETXCONSTANTBYTES(&_Elderberrypolygonzkevm.CallOpts)
}

// INITIALIZETXCONSTANTBYTESEMPTYMETADATA is a free data retrieval call binding the contract method 0x676870d2.
//
// Solidity: function INITIALIZE_TX_CONSTANT_BYTES_EMPTY_METADATA() view returns(uint16)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCaller) INITIALIZETXCONSTANTBYTESEMPTYMETADATA(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _Elderberrypolygonzkevm.contract.Call(opts, &out, "INITIALIZE_TX_CONSTANT_BYTES_EMPTY_METADATA")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// INITIALIZETXCONSTANTBYTESEMPTYMETADATA is a free data retrieval call binding the contract method 0x676870d2.
//
// Solidity: function INITIALIZE_TX_CONSTANT_BYTES_EMPTY_METADATA() view returns(uint16)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmSession) INITIALIZETXCONSTANTBYTESEMPTYMETADATA() (uint16, error) {
	return _Elderberrypolygonzkevm.Contract.INITIALIZETXCONSTANTBYTESEMPTYMETADATA(&_Elderberrypolygonzkevm.CallOpts)
}

// INITIALIZETXCONSTANTBYTESEMPTYMETADATA is a free data retrieval call binding the contract method 0x676870d2.
//
// Solidity: function INITIALIZE_TX_CONSTANT_BYTES_EMPTY_METADATA() view returns(uint16)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCallerSession) INITIALIZETXCONSTANTBYTESEMPTYMETADATA() (uint16, error) {
	return _Elderberrypolygonzkevm.Contract.INITIALIZETXCONSTANTBYTESEMPTYMETADATA(&_Elderberrypolygonzkevm.CallOpts)
}

// INITIALIZETXDATALENEMPTYMETADATA is a free data retrieval call binding the contract method 0xc7fffd4b.
//
// Solidity: function INITIALIZE_TX_DATA_LEN_EMPTY_METADATA() view returns(uint8)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCaller) INITIALIZETXDATALENEMPTYMETADATA(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Elderberrypolygonzkevm.contract.Call(opts, &out, "INITIALIZE_TX_DATA_LEN_EMPTY_METADATA")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// INITIALIZETXDATALENEMPTYMETADATA is a free data retrieval call binding the contract method 0xc7fffd4b.
//
// Solidity: function INITIALIZE_TX_DATA_LEN_EMPTY_METADATA() view returns(uint8)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmSession) INITIALIZETXDATALENEMPTYMETADATA() (uint8, error) {
	return _Elderberrypolygonzkevm.Contract.INITIALIZETXDATALENEMPTYMETADATA(&_Elderberrypolygonzkevm.CallOpts)
}

// INITIALIZETXDATALENEMPTYMETADATA is a free data retrieval call binding the contract method 0xc7fffd4b.
//
// Solidity: function INITIALIZE_TX_DATA_LEN_EMPTY_METADATA() view returns(uint8)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCallerSession) INITIALIZETXDATALENEMPTYMETADATA() (uint8, error) {
	return _Elderberrypolygonzkevm.Contract.INITIALIZETXDATALENEMPTYMETADATA(&_Elderberrypolygonzkevm.CallOpts)
}

// INITIALIZETXEFFECTIVEPERCENTAGE is a free data retrieval call binding the contract method 0x40b5de6c.
//
// Solidity: function INITIALIZE_TX_EFFECTIVE_PERCENTAGE() view returns(bytes1)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCaller) INITIALIZETXEFFECTIVEPERCENTAGE(opts *bind.CallOpts) ([1]byte, error) {
	var out []interface{}
	err := _Elderberrypolygonzkevm.contract.Call(opts, &out, "INITIALIZE_TX_EFFECTIVE_PERCENTAGE")

	if err != nil {
		return *new([1]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([1]byte)).(*[1]byte)

	return out0, err

}

// INITIALIZETXEFFECTIVEPERCENTAGE is a free data retrieval call binding the contract method 0x40b5de6c.
//
// Solidity: function INITIALIZE_TX_EFFECTIVE_PERCENTAGE() view returns(bytes1)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmSession) INITIALIZETXEFFECTIVEPERCENTAGE() ([1]byte, error) {
	return _Elderberrypolygonzkevm.Contract.INITIALIZETXEFFECTIVEPERCENTAGE(&_Elderberrypolygonzkevm.CallOpts)
}

// INITIALIZETXEFFECTIVEPERCENTAGE is a free data retrieval call binding the contract method 0x40b5de6c.
//
// Solidity: function INITIALIZE_TX_EFFECTIVE_PERCENTAGE() view returns(bytes1)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCallerSession) INITIALIZETXEFFECTIVEPERCENTAGE() ([1]byte, error) {
	return _Elderberrypolygonzkevm.Contract.INITIALIZETXEFFECTIVEPERCENTAGE(&_Elderberrypolygonzkevm.CallOpts)
}

// SETUPETROGTX is a free data retrieval call binding the contract method 0xaf7f3e02.
//
// Solidity: function SET_UP_ETROG_TX() view returns(bytes)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCaller) SETUPETROGTX(opts *bind.CallOpts) ([]byte, error) {
	var out []interface{}
	err := _Elderberrypolygonzkevm.contract.Call(opts, &out, "SET_UP_ETROG_TX")

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// SETUPETROGTX is a free data retrieval call binding the contract method 0xaf7f3e02.
//
// Solidity: function SET_UP_ETROG_TX() view returns(bytes)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmSession) SETUPETROGTX() ([]byte, error) {
	return _Elderberrypolygonzkevm.Contract.SETUPETROGTX(&_Elderberrypolygonzkevm.CallOpts)
}

// SETUPETROGTX is a free data retrieval call binding the contract method 0xaf7f3e02.
//
// Solidity: function SET_UP_ETROG_TX() view returns(bytes)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCallerSession) SETUPETROGTX() ([]byte, error) {
	return _Elderberrypolygonzkevm.Contract.SETUPETROGTX(&_Elderberrypolygonzkevm.CallOpts)
}

// SIGNATUREINITIALIZETXR is a free data retrieval call binding the contract method 0xb0afe154.
//
// Solidity: function SIGNATURE_INITIALIZE_TX_R() view returns(bytes32)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCaller) SIGNATUREINITIALIZETXR(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Elderberrypolygonzkevm.contract.Call(opts, &out, "SIGNATURE_INITIALIZE_TX_R")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// SIGNATUREINITIALIZETXR is a free data retrieval call binding the contract method 0xb0afe154.
//
// Solidity: function SIGNATURE_INITIALIZE_TX_R() view returns(bytes32)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmSession) SIGNATUREINITIALIZETXR() ([32]byte, error) {
	return _Elderberrypolygonzkevm.Contract.SIGNATUREINITIALIZETXR(&_Elderberrypolygonzkevm.CallOpts)
}

// SIGNATUREINITIALIZETXR is a free data retrieval call binding the contract method 0xb0afe154.
//
// Solidity: function SIGNATURE_INITIALIZE_TX_R() view returns(bytes32)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCallerSession) SIGNATUREINITIALIZETXR() ([32]byte, error) {
	return _Elderberrypolygonzkevm.Contract.SIGNATUREINITIALIZETXR(&_Elderberrypolygonzkevm.CallOpts)
}

// SIGNATUREINITIALIZETXS is a free data retrieval call binding the contract method 0xd7bc90ff.
//
// Solidity: function SIGNATURE_INITIALIZE_TX_S() view returns(bytes32)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCaller) SIGNATUREINITIALIZETXS(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Elderberrypolygonzkevm.contract.Call(opts, &out, "SIGNATURE_INITIALIZE_TX_S")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// SIGNATUREINITIALIZETXS is a free data retrieval call binding the contract method 0xd7bc90ff.
//
// Solidity: function SIGNATURE_INITIALIZE_TX_S() view returns(bytes32)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmSession) SIGNATUREINITIALIZETXS() ([32]byte, error) {
	return _Elderberrypolygonzkevm.Contract.SIGNATUREINITIALIZETXS(&_Elderberrypolygonzkevm.CallOpts)
}

// SIGNATUREINITIALIZETXS is a free data retrieval call binding the contract method 0xd7bc90ff.
//
// Solidity: function SIGNATURE_INITIALIZE_TX_S() view returns(bytes32)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCallerSession) SIGNATUREINITIALIZETXS() ([32]byte, error) {
	return _Elderberrypolygonzkevm.Contract.SIGNATUREINITIALIZETXS(&_Elderberrypolygonzkevm.CallOpts)
}

// SIGNATUREINITIALIZETXV is a free data retrieval call binding the contract method 0xf35dda47.
//
// Solidity: function SIGNATURE_INITIALIZE_TX_V() view returns(uint8)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCaller) SIGNATUREINITIALIZETXV(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Elderberrypolygonzkevm.contract.Call(opts, &out, "SIGNATURE_INITIALIZE_TX_V")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// SIGNATUREINITIALIZETXV is a free data retrieval call binding the contract method 0xf35dda47.
//
// Solidity: function SIGNATURE_INITIALIZE_TX_V() view returns(uint8)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmSession) SIGNATUREINITIALIZETXV() (uint8, error) {
	return _Elderberrypolygonzkevm.Contract.SIGNATUREINITIALIZETXV(&_Elderberrypolygonzkevm.CallOpts)
}

// SIGNATUREINITIALIZETXV is a free data retrieval call binding the contract method 0xf35dda47.
//
// Solidity: function SIGNATURE_INITIALIZE_TX_V() view returns(uint8)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCallerSession) SIGNATUREINITIALIZETXV() (uint8, error) {
	return _Elderberrypolygonzkevm.Contract.SIGNATUREINITIALIZETXV(&_Elderberrypolygonzkevm.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCaller) Admin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Elderberrypolygonzkevm.contract.Call(opts, &out, "admin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmSession) Admin() (common.Address, error) {
	return _Elderberrypolygonzkevm.Contract.Admin(&_Elderberrypolygonzkevm.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCallerSession) Admin() (common.Address, error) {
	return _Elderberrypolygonzkevm.Contract.Admin(&_Elderberrypolygonzkevm.CallOpts)
}

// BridgeAddress is a free data retrieval call binding the contract method 0xa3c573eb.
//
// Solidity: function bridgeAddress() view returns(address)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCaller) BridgeAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Elderberrypolygonzkevm.contract.Call(opts, &out, "bridgeAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BridgeAddress is a free data retrieval call binding the contract method 0xa3c573eb.
//
// Solidity: function bridgeAddress() view returns(address)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmSession) BridgeAddress() (common.Address, error) {
	return _Elderberrypolygonzkevm.Contract.BridgeAddress(&_Elderberrypolygonzkevm.CallOpts)
}

// BridgeAddress is a free data retrieval call binding the contract method 0xa3c573eb.
//
// Solidity: function bridgeAddress() view returns(address)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCallerSession) BridgeAddress() (common.Address, error) {
	return _Elderberrypolygonzkevm.Contract.BridgeAddress(&_Elderberrypolygonzkevm.CallOpts)
}

// CalculatePolPerForceBatch is a free data retrieval call binding the contract method 0x00d0295d.
//
// Solidity: function calculatePolPerForceBatch() view returns(uint256)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCaller) CalculatePolPerForceBatch(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Elderberrypolygonzkevm.contract.Call(opts, &out, "calculatePolPerForceBatch")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CalculatePolPerForceBatch is a free data retrieval call binding the contract method 0x00d0295d.
//
// Solidity: function calculatePolPerForceBatch() view returns(uint256)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmSession) CalculatePolPerForceBatch() (*big.Int, error) {
	return _Elderberrypolygonzkevm.Contract.CalculatePolPerForceBatch(&_Elderberrypolygonzkevm.CallOpts)
}

// CalculatePolPerForceBatch is a free data retrieval call binding the contract method 0x00d0295d.
//
// Solidity: function calculatePolPerForceBatch() view returns(uint256)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCallerSession) CalculatePolPerForceBatch() (*big.Int, error) {
	return _Elderberrypolygonzkevm.Contract.CalculatePolPerForceBatch(&_Elderberrypolygonzkevm.CallOpts)
}

// ForceBatchAddress is a free data retrieval call binding the contract method 0x2c111c06.
//
// Solidity: function forceBatchAddress() view returns(address)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCaller) ForceBatchAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Elderberrypolygonzkevm.contract.Call(opts, &out, "forceBatchAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ForceBatchAddress is a free data retrieval call binding the contract method 0x2c111c06.
//
// Solidity: function forceBatchAddress() view returns(address)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmSession) ForceBatchAddress() (common.Address, error) {
	return _Elderberrypolygonzkevm.Contract.ForceBatchAddress(&_Elderberrypolygonzkevm.CallOpts)
}

// ForceBatchAddress is a free data retrieval call binding the contract method 0x2c111c06.
//
// Solidity: function forceBatchAddress() view returns(address)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCallerSession) ForceBatchAddress() (common.Address, error) {
	return _Elderberrypolygonzkevm.Contract.ForceBatchAddress(&_Elderberrypolygonzkevm.CallOpts)
}

// ForceBatchTimeout is a free data retrieval call binding the contract method 0xc754c7ed.
//
// Solidity: function forceBatchTimeout() view returns(uint64)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCaller) ForceBatchTimeout(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Elderberrypolygonzkevm.contract.Call(opts, &out, "forceBatchTimeout")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// ForceBatchTimeout is a free data retrieval call binding the contract method 0xc754c7ed.
//
// Solidity: function forceBatchTimeout() view returns(uint64)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmSession) ForceBatchTimeout() (uint64, error) {
	return _Elderberrypolygonzkevm.Contract.ForceBatchTimeout(&_Elderberrypolygonzkevm.CallOpts)
}

// ForceBatchTimeout is a free data retrieval call binding the contract method 0xc754c7ed.
//
// Solidity: function forceBatchTimeout() view returns(uint64)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCallerSession) ForceBatchTimeout() (uint64, error) {
	return _Elderberrypolygonzkevm.Contract.ForceBatchTimeout(&_Elderberrypolygonzkevm.CallOpts)
}

// ForcedBatches is a free data retrieval call binding the contract method 0x6b8616ce.
//
// Solidity: function forcedBatches(uint64 ) view returns(bytes32)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCaller) ForcedBatches(opts *bind.CallOpts, arg0 uint64) ([32]byte, error) {
	var out []interface{}
	err := _Elderberrypolygonzkevm.contract.Call(opts, &out, "forcedBatches", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ForcedBatches is a free data retrieval call binding the contract method 0x6b8616ce.
//
// Solidity: function forcedBatches(uint64 ) view returns(bytes32)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmSession) ForcedBatches(arg0 uint64) ([32]byte, error) {
	return _Elderberrypolygonzkevm.Contract.ForcedBatches(&_Elderberrypolygonzkevm.CallOpts, arg0)
}

// ForcedBatches is a free data retrieval call binding the contract method 0x6b8616ce.
//
// Solidity: function forcedBatches(uint64 ) view returns(bytes32)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCallerSession) ForcedBatches(arg0 uint64) ([32]byte, error) {
	return _Elderberrypolygonzkevm.Contract.ForcedBatches(&_Elderberrypolygonzkevm.CallOpts, arg0)
}

// GasTokenAddress is a free data retrieval call binding the contract method 0x3c351e10.
//
// Solidity: function gasTokenAddress() view returns(address)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCaller) GasTokenAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Elderberrypolygonzkevm.contract.Call(opts, &out, "gasTokenAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GasTokenAddress is a free data retrieval call binding the contract method 0x3c351e10.
//
// Solidity: function gasTokenAddress() view returns(address)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmSession) GasTokenAddress() (common.Address, error) {
	return _Elderberrypolygonzkevm.Contract.GasTokenAddress(&_Elderberrypolygonzkevm.CallOpts)
}

// GasTokenAddress is a free data retrieval call binding the contract method 0x3c351e10.
//
// Solidity: function gasTokenAddress() view returns(address)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCallerSession) GasTokenAddress() (common.Address, error) {
	return _Elderberrypolygonzkevm.Contract.GasTokenAddress(&_Elderberrypolygonzkevm.CallOpts)
}

// GasTokenNetwork is a free data retrieval call binding the contract method 0x3cbc795b.
//
// Solidity: function gasTokenNetwork() view returns(uint32)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCaller) GasTokenNetwork(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _Elderberrypolygonzkevm.contract.Call(opts, &out, "gasTokenNetwork")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// GasTokenNetwork is a free data retrieval call binding the contract method 0x3cbc795b.
//
// Solidity: function gasTokenNetwork() view returns(uint32)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmSession) GasTokenNetwork() (uint32, error) {
	return _Elderberrypolygonzkevm.Contract.GasTokenNetwork(&_Elderberrypolygonzkevm.CallOpts)
}

// GasTokenNetwork is a free data retrieval call binding the contract method 0x3cbc795b.
//
// Solidity: function gasTokenNetwork() view returns(uint32)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCallerSession) GasTokenNetwork() (uint32, error) {
	return _Elderberrypolygonzkevm.Contract.GasTokenNetwork(&_Elderberrypolygonzkevm.CallOpts)
}

// GenerateInitializeTransaction is a free data retrieval call binding the contract method 0xa652f26c.
//
// Solidity: function generateInitializeTransaction(uint32 networkID, address _gasTokenAddress, uint32 _gasTokenNetwork, bytes _gasTokenMetadata) view returns(bytes)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCaller) GenerateInitializeTransaction(opts *bind.CallOpts, networkID uint32, _gasTokenAddress common.Address, _gasTokenNetwork uint32, _gasTokenMetadata []byte) ([]byte, error) {
	var out []interface{}
	err := _Elderberrypolygonzkevm.contract.Call(opts, &out, "generateInitializeTransaction", networkID, _gasTokenAddress, _gasTokenNetwork, _gasTokenMetadata)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GenerateInitializeTransaction is a free data retrieval call binding the contract method 0xa652f26c.
//
// Solidity: function generateInitializeTransaction(uint32 networkID, address _gasTokenAddress, uint32 _gasTokenNetwork, bytes _gasTokenMetadata) view returns(bytes)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmSession) GenerateInitializeTransaction(networkID uint32, _gasTokenAddress common.Address, _gasTokenNetwork uint32, _gasTokenMetadata []byte) ([]byte, error) {
	return _Elderberrypolygonzkevm.Contract.GenerateInitializeTransaction(&_Elderberrypolygonzkevm.CallOpts, networkID, _gasTokenAddress, _gasTokenNetwork, _gasTokenMetadata)
}

// GenerateInitializeTransaction is a free data retrieval call binding the contract method 0xa652f26c.
//
// Solidity: function generateInitializeTransaction(uint32 networkID, address _gasTokenAddress, uint32 _gasTokenNetwork, bytes _gasTokenMetadata) view returns(bytes)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCallerSession) GenerateInitializeTransaction(networkID uint32, _gasTokenAddress common.Address, _gasTokenNetwork uint32, _gasTokenMetadata []byte) ([]byte, error) {
	return _Elderberrypolygonzkevm.Contract.GenerateInitializeTransaction(&_Elderberrypolygonzkevm.CallOpts, networkID, _gasTokenAddress, _gasTokenNetwork, _gasTokenMetadata)
}

// GlobalExitRootManager is a free data retrieval call binding the contract method 0xd02103ca.
//
// Solidity: function globalExitRootManager() view returns(address)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCaller) GlobalExitRootManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Elderberrypolygonzkevm.contract.Call(opts, &out, "globalExitRootManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GlobalExitRootManager is a free data retrieval call binding the contract method 0xd02103ca.
//
// Solidity: function globalExitRootManager() view returns(address)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmSession) GlobalExitRootManager() (common.Address, error) {
	return _Elderberrypolygonzkevm.Contract.GlobalExitRootManager(&_Elderberrypolygonzkevm.CallOpts)
}

// GlobalExitRootManager is a free data retrieval call binding the contract method 0xd02103ca.
//
// Solidity: function globalExitRootManager() view returns(address)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCallerSession) GlobalExitRootManager() (common.Address, error) {
	return _Elderberrypolygonzkevm.Contract.GlobalExitRootManager(&_Elderberrypolygonzkevm.CallOpts)
}

// LastAccInputHash is a free data retrieval call binding the contract method 0x6e05d2cd.
//
// Solidity: function lastAccInputHash() view returns(bytes32)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCaller) LastAccInputHash(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Elderberrypolygonzkevm.contract.Call(opts, &out, "lastAccInputHash")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// LastAccInputHash is a free data retrieval call binding the contract method 0x6e05d2cd.
//
// Solidity: function lastAccInputHash() view returns(bytes32)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmSession) LastAccInputHash() ([32]byte, error) {
	return _Elderberrypolygonzkevm.Contract.LastAccInputHash(&_Elderberrypolygonzkevm.CallOpts)
}

// LastAccInputHash is a free data retrieval call binding the contract method 0x6e05d2cd.
//
// Solidity: function lastAccInputHash() view returns(bytes32)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCallerSession) LastAccInputHash() ([32]byte, error) {
	return _Elderberrypolygonzkevm.Contract.LastAccInputHash(&_Elderberrypolygonzkevm.CallOpts)
}

// LastForceBatch is a free data retrieval call binding the contract method 0xe7a7ed02.
//
// Solidity: function lastForceBatch() view returns(uint64)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCaller) LastForceBatch(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Elderberrypolygonzkevm.contract.Call(opts, &out, "lastForceBatch")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// LastForceBatch is a free data retrieval call binding the contract method 0xe7a7ed02.
//
// Solidity: function lastForceBatch() view returns(uint64)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmSession) LastForceBatch() (uint64, error) {
	return _Elderberrypolygonzkevm.Contract.LastForceBatch(&_Elderberrypolygonzkevm.CallOpts)
}

// LastForceBatch is a free data retrieval call binding the contract method 0xe7a7ed02.
//
// Solidity: function lastForceBatch() view returns(uint64)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCallerSession) LastForceBatch() (uint64, error) {
	return _Elderberrypolygonzkevm.Contract.LastForceBatch(&_Elderberrypolygonzkevm.CallOpts)
}

// LastForceBatchSequenced is a free data retrieval call binding the contract method 0x45605267.
//
// Solidity: function lastForceBatchSequenced() view returns(uint64)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCaller) LastForceBatchSequenced(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Elderberrypolygonzkevm.contract.Call(opts, &out, "lastForceBatchSequenced")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// LastForceBatchSequenced is a free data retrieval call binding the contract method 0x45605267.
//
// Solidity: function lastForceBatchSequenced() view returns(uint64)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmSession) LastForceBatchSequenced() (uint64, error) {
	return _Elderberrypolygonzkevm.Contract.LastForceBatchSequenced(&_Elderberrypolygonzkevm.CallOpts)
}

// LastForceBatchSequenced is a free data retrieval call binding the contract method 0x45605267.
//
// Solidity: function lastForceBatchSequenced() view returns(uint64)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCallerSession) LastForceBatchSequenced() (uint64, error) {
	return _Elderberrypolygonzkevm.Contract.LastForceBatchSequenced(&_Elderberrypolygonzkevm.CallOpts)
}

// NetworkName is a free data retrieval call binding the contract method 0x107bf28c.
//
// Solidity: function networkName() view returns(string)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCaller) NetworkName(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Elderberrypolygonzkevm.contract.Call(opts, &out, "networkName")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// NetworkName is a free data retrieval call binding the contract method 0x107bf28c.
//
// Solidity: function networkName() view returns(string)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmSession) NetworkName() (string, error) {
	return _Elderberrypolygonzkevm.Contract.NetworkName(&_Elderberrypolygonzkevm.CallOpts)
}

// NetworkName is a free data retrieval call binding the contract method 0x107bf28c.
//
// Solidity: function networkName() view returns(string)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCallerSession) NetworkName() (string, error) {
	return _Elderberrypolygonzkevm.Contract.NetworkName(&_Elderberrypolygonzkevm.CallOpts)
}

// PendingAdmin is a free data retrieval call binding the contract method 0x26782247.
//
// Solidity: function pendingAdmin() view returns(address)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCaller) PendingAdmin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Elderberrypolygonzkevm.contract.Call(opts, &out, "pendingAdmin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PendingAdmin is a free data retrieval call binding the contract method 0x26782247.
//
// Solidity: function pendingAdmin() view returns(address)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmSession) PendingAdmin() (common.Address, error) {
	return _Elderberrypolygonzkevm.Contract.PendingAdmin(&_Elderberrypolygonzkevm.CallOpts)
}

// PendingAdmin is a free data retrieval call binding the contract method 0x26782247.
//
// Solidity: function pendingAdmin() view returns(address)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCallerSession) PendingAdmin() (common.Address, error) {
	return _Elderberrypolygonzkevm.Contract.PendingAdmin(&_Elderberrypolygonzkevm.CallOpts)
}

// Pol is a free data retrieval call binding the contract method 0xe46761c4.
//
// Solidity: function pol() view returns(address)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCaller) Pol(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Elderberrypolygonzkevm.contract.Call(opts, &out, "pol")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Pol is a free data retrieval call binding the contract method 0xe46761c4.
//
// Solidity: function pol() view returns(address)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmSession) Pol() (common.Address, error) {
	return _Elderberrypolygonzkevm.Contract.Pol(&_Elderberrypolygonzkevm.CallOpts)
}

// Pol is a free data retrieval call binding the contract method 0xe46761c4.
//
// Solidity: function pol() view returns(address)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCallerSession) Pol() (common.Address, error) {
	return _Elderberrypolygonzkevm.Contract.Pol(&_Elderberrypolygonzkevm.CallOpts)
}

// RollupManager is a free data retrieval call binding the contract method 0x49b7b802.
//
// Solidity: function rollupManager() view returns(address)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCaller) RollupManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Elderberrypolygonzkevm.contract.Call(opts, &out, "rollupManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// RollupManager is a free data retrieval call binding the contract method 0x49b7b802.
//
// Solidity: function rollupManager() view returns(address)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmSession) RollupManager() (common.Address, error) {
	return _Elderberrypolygonzkevm.Contract.RollupManager(&_Elderberrypolygonzkevm.CallOpts)
}

// RollupManager is a free data retrieval call binding the contract method 0x49b7b802.
//
// Solidity: function rollupManager() view returns(address)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCallerSession) RollupManager() (common.Address, error) {
	return _Elderberrypolygonzkevm.Contract.RollupManager(&_Elderberrypolygonzkevm.CallOpts)
}

// TrustedSequencer is a free data retrieval call binding the contract method 0xcfa8ed47.
//
// Solidity: function trustedSequencer() view returns(address)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCaller) TrustedSequencer(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Elderberrypolygonzkevm.contract.Call(opts, &out, "trustedSequencer")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TrustedSequencer is a free data retrieval call binding the contract method 0xcfa8ed47.
//
// Solidity: function trustedSequencer() view returns(address)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmSession) TrustedSequencer() (common.Address, error) {
	return _Elderberrypolygonzkevm.Contract.TrustedSequencer(&_Elderberrypolygonzkevm.CallOpts)
}

// TrustedSequencer is a free data retrieval call binding the contract method 0xcfa8ed47.
//
// Solidity: function trustedSequencer() view returns(address)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCallerSession) TrustedSequencer() (common.Address, error) {
	return _Elderberrypolygonzkevm.Contract.TrustedSequencer(&_Elderberrypolygonzkevm.CallOpts)
}

// TrustedSequencerURL is a free data retrieval call binding the contract method 0x542028d5.
//
// Solidity: function trustedSequencerURL() view returns(string)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCaller) TrustedSequencerURL(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Elderberrypolygonzkevm.contract.Call(opts, &out, "trustedSequencerURL")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TrustedSequencerURL is a free data retrieval call binding the contract method 0x542028d5.
//
// Solidity: function trustedSequencerURL() view returns(string)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmSession) TrustedSequencerURL() (string, error) {
	return _Elderberrypolygonzkevm.Contract.TrustedSequencerURL(&_Elderberrypolygonzkevm.CallOpts)
}

// TrustedSequencerURL is a free data retrieval call binding the contract method 0x542028d5.
//
// Solidity: function trustedSequencerURL() view returns(string)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmCallerSession) TrustedSequencerURL() (string, error) {
	return _Elderberrypolygonzkevm.Contract.TrustedSequencerURL(&_Elderberrypolygonzkevm.CallOpts)
}

// AcceptAdminRole is a paid mutator transaction binding the contract method 0x8c3d7301.
//
// Solidity: function acceptAdminRole() returns()
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmTransactor) AcceptAdminRole(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Elderberrypolygonzkevm.contract.Transact(opts, "acceptAdminRole")
}

// AcceptAdminRole is a paid mutator transaction binding the contract method 0x8c3d7301.
//
// Solidity: function acceptAdminRole() returns()
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmSession) AcceptAdminRole() (*types.Transaction, error) {
	return _Elderberrypolygonzkevm.Contract.AcceptAdminRole(&_Elderberrypolygonzkevm.TransactOpts)
}

// AcceptAdminRole is a paid mutator transaction binding the contract method 0x8c3d7301.
//
// Solidity: function acceptAdminRole() returns()
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmTransactorSession) AcceptAdminRole() (*types.Transaction, error) {
	return _Elderberrypolygonzkevm.Contract.AcceptAdminRole(&_Elderberrypolygonzkevm.TransactOpts)
}

// ForceBatch is a paid mutator transaction binding the contract method 0xeaeb077b.
//
// Solidity: function forceBatch(bytes transactions, uint256 polAmount) returns()
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmTransactor) ForceBatch(opts *bind.TransactOpts, transactions []byte, polAmount *big.Int) (*types.Transaction, error) {
	return _Elderberrypolygonzkevm.contract.Transact(opts, "forceBatch", transactions, polAmount)
}

// ForceBatch is a paid mutator transaction binding the contract method 0xeaeb077b.
//
// Solidity: function forceBatch(bytes transactions, uint256 polAmount) returns()
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmSession) ForceBatch(transactions []byte, polAmount *big.Int) (*types.Transaction, error) {
	return _Elderberrypolygonzkevm.Contract.ForceBatch(&_Elderberrypolygonzkevm.TransactOpts, transactions, polAmount)
}

// ForceBatch is a paid mutator transaction binding the contract method 0xeaeb077b.
//
// Solidity: function forceBatch(bytes transactions, uint256 polAmount) returns()
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmTransactorSession) ForceBatch(transactions []byte, polAmount *big.Int) (*types.Transaction, error) {
	return _Elderberrypolygonzkevm.Contract.ForceBatch(&_Elderberrypolygonzkevm.TransactOpts, transactions, polAmount)
}

// Initialize is a paid mutator transaction binding the contract method 0x71257022.
//
// Solidity: function initialize(address _admin, address sequencer, uint32 networkID, address _gasTokenAddress, string sequencerURL, string _networkName) returns()
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmTransactor) Initialize(opts *bind.TransactOpts, _admin common.Address, sequencer common.Address, networkID uint32, _gasTokenAddress common.Address, sequencerURL string, _networkName string) (*types.Transaction, error) {
	return _Elderberrypolygonzkevm.contract.Transact(opts, "initialize", _admin, sequencer, networkID, _gasTokenAddress, sequencerURL, _networkName)
}

// Initialize is a paid mutator transaction binding the contract method 0x71257022.
//
// Solidity: function initialize(address _admin, address sequencer, uint32 networkID, address _gasTokenAddress, string sequencerURL, string _networkName) returns()
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmSession) Initialize(_admin common.Address, sequencer common.Address, networkID uint32, _gasTokenAddress common.Address, sequencerURL string, _networkName string) (*types.Transaction, error) {
	return _Elderberrypolygonzkevm.Contract.Initialize(&_Elderberrypolygonzkevm.TransactOpts, _admin, sequencer, networkID, _gasTokenAddress, sequencerURL, _networkName)
}

// Initialize is a paid mutator transaction binding the contract method 0x71257022.
//
// Solidity: function initialize(address _admin, address sequencer, uint32 networkID, address _gasTokenAddress, string sequencerURL, string _networkName) returns()
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmTransactorSession) Initialize(_admin common.Address, sequencer common.Address, networkID uint32, _gasTokenAddress common.Address, sequencerURL string, _networkName string) (*types.Transaction, error) {
	return _Elderberrypolygonzkevm.Contract.Initialize(&_Elderberrypolygonzkevm.TransactOpts, _admin, sequencer, networkID, _gasTokenAddress, sequencerURL, _networkName)
}

// InitializeUpgrade is a paid mutator transaction binding the contract method 0x5d6717a5.
//
// Solidity: function initializeUpgrade(address _admin, address _trustedSequencer, string _trustedSequencerURL, string _networkName, bytes32 _lastAccInputHash) returns()
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmTransactor) InitializeUpgrade(opts *bind.TransactOpts, _admin common.Address, _trustedSequencer common.Address, _trustedSequencerURL string, _networkName string, _lastAccInputHash [32]byte) (*types.Transaction, error) {
	return _Elderberrypolygonzkevm.contract.Transact(opts, "initializeUpgrade", _admin, _trustedSequencer, _trustedSequencerURL, _networkName, _lastAccInputHash)
}

// InitializeUpgrade is a paid mutator transaction binding the contract method 0x5d6717a5.
//
// Solidity: function initializeUpgrade(address _admin, address _trustedSequencer, string _trustedSequencerURL, string _networkName, bytes32 _lastAccInputHash) returns()
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmSession) InitializeUpgrade(_admin common.Address, _trustedSequencer common.Address, _trustedSequencerURL string, _networkName string, _lastAccInputHash [32]byte) (*types.Transaction, error) {
	return _Elderberrypolygonzkevm.Contract.InitializeUpgrade(&_Elderberrypolygonzkevm.TransactOpts, _admin, _trustedSequencer, _trustedSequencerURL, _networkName, _lastAccInputHash)
}

// InitializeUpgrade is a paid mutator transaction binding the contract method 0x5d6717a5.
//
// Solidity: function initializeUpgrade(address _admin, address _trustedSequencer, string _trustedSequencerURL, string _networkName, bytes32 _lastAccInputHash) returns()
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmTransactorSession) InitializeUpgrade(_admin common.Address, _trustedSequencer common.Address, _trustedSequencerURL string, _networkName string, _lastAccInputHash [32]byte) (*types.Transaction, error) {
	return _Elderberrypolygonzkevm.Contract.InitializeUpgrade(&_Elderberrypolygonzkevm.TransactOpts, _admin, _trustedSequencer, _trustedSequencerURL, _networkName, _lastAccInputHash)
}

// OnVerifyBatches is a paid mutator transaction binding the contract method 0x32c2d153.
//
// Solidity: function onVerifyBatches(uint64 lastVerifiedBatch, bytes32 newStateRoot, address aggregator) returns()
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmTransactor) OnVerifyBatches(opts *bind.TransactOpts, lastVerifiedBatch uint64, newStateRoot [32]byte, aggregator common.Address) (*types.Transaction, error) {
	return _Elderberrypolygonzkevm.contract.Transact(opts, "onVerifyBatches", lastVerifiedBatch, newStateRoot, aggregator)
}

// OnVerifyBatches is a paid mutator transaction binding the contract method 0x32c2d153.
//
// Solidity: function onVerifyBatches(uint64 lastVerifiedBatch, bytes32 newStateRoot, address aggregator) returns()
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmSession) OnVerifyBatches(lastVerifiedBatch uint64, newStateRoot [32]byte, aggregator common.Address) (*types.Transaction, error) {
	return _Elderberrypolygonzkevm.Contract.OnVerifyBatches(&_Elderberrypolygonzkevm.TransactOpts, lastVerifiedBatch, newStateRoot, aggregator)
}

// OnVerifyBatches is a paid mutator transaction binding the contract method 0x32c2d153.
//
// Solidity: function onVerifyBatches(uint64 lastVerifiedBatch, bytes32 newStateRoot, address aggregator) returns()
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmTransactorSession) OnVerifyBatches(lastVerifiedBatch uint64, newStateRoot [32]byte, aggregator common.Address) (*types.Transaction, error) {
	return _Elderberrypolygonzkevm.Contract.OnVerifyBatches(&_Elderberrypolygonzkevm.TransactOpts, lastVerifiedBatch, newStateRoot, aggregator)
}

// SequenceBatches is a paid mutator transaction binding the contract method 0xecef3f99.
//
// Solidity: function sequenceBatches((bytes,bytes32,uint64,bytes32)[] batches, address l2Coinbase) returns()
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmTransactor) SequenceBatches(opts *bind.TransactOpts, batches []PolygonRollupBaseEtrogBatchData, l2Coinbase common.Address) (*types.Transaction, error) {
	return _Elderberrypolygonzkevm.contract.Transact(opts, "sequenceBatches", batches, l2Coinbase)
}

// SequenceBatches is a paid mutator transaction binding the contract method 0xecef3f99.
//
// Solidity: function sequenceBatches((bytes,bytes32,uint64,bytes32)[] batches, address l2Coinbase) returns()
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmSession) SequenceBatches(batches []PolygonRollupBaseEtrogBatchData, l2Coinbase common.Address) (*types.Transaction, error) {
	return _Elderberrypolygonzkevm.Contract.SequenceBatches(&_Elderberrypolygonzkevm.TransactOpts, batches, l2Coinbase)
}

// SequenceBatches is a paid mutator transaction binding the contract method 0xecef3f99.
//
// Solidity: function sequenceBatches((bytes,bytes32,uint64,bytes32)[] batches, address l2Coinbase) returns()
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmTransactorSession) SequenceBatches(batches []PolygonRollupBaseEtrogBatchData, l2Coinbase common.Address) (*types.Transaction, error) {
	return _Elderberrypolygonzkevm.Contract.SequenceBatches(&_Elderberrypolygonzkevm.TransactOpts, batches, l2Coinbase)
}

// SequenceForceBatches is a paid mutator transaction binding the contract method 0x9f26f840.
//
// Solidity: function sequenceForceBatches((bytes,bytes32,uint64,bytes32)[] batches) returns()
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmTransactor) SequenceForceBatches(opts *bind.TransactOpts, batches []PolygonRollupBaseEtrogBatchData) (*types.Transaction, error) {
	return _Elderberrypolygonzkevm.contract.Transact(opts, "sequenceForceBatches", batches)
}

// SequenceForceBatches is a paid mutator transaction binding the contract method 0x9f26f840.
//
// Solidity: function sequenceForceBatches((bytes,bytes32,uint64,bytes32)[] batches) returns()
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmSession) SequenceForceBatches(batches []PolygonRollupBaseEtrogBatchData) (*types.Transaction, error) {
	return _Elderberrypolygonzkevm.Contract.SequenceForceBatches(&_Elderberrypolygonzkevm.TransactOpts, batches)
}

// SequenceForceBatches is a paid mutator transaction binding the contract method 0x9f26f840.
//
// Solidity: function sequenceForceBatches((bytes,bytes32,uint64,bytes32)[] batches) returns()
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmTransactorSession) SequenceForceBatches(batches []PolygonRollupBaseEtrogBatchData) (*types.Transaction, error) {
	return _Elderberrypolygonzkevm.Contract.SequenceForceBatches(&_Elderberrypolygonzkevm.TransactOpts, batches)
}

// SetForceBatchAddress is a paid mutator transaction binding the contract method 0x91cafe32.
//
// Solidity: function setForceBatchAddress(address newForceBatchAddress) returns()
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmTransactor) SetForceBatchAddress(opts *bind.TransactOpts, newForceBatchAddress common.Address) (*types.Transaction, error) {
	return _Elderberrypolygonzkevm.contract.Transact(opts, "setForceBatchAddress", newForceBatchAddress)
}

// SetForceBatchAddress is a paid mutator transaction binding the contract method 0x91cafe32.
//
// Solidity: function setForceBatchAddress(address newForceBatchAddress) returns()
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmSession) SetForceBatchAddress(newForceBatchAddress common.Address) (*types.Transaction, error) {
	return _Elderberrypolygonzkevm.Contract.SetForceBatchAddress(&_Elderberrypolygonzkevm.TransactOpts, newForceBatchAddress)
}

// SetForceBatchAddress is a paid mutator transaction binding the contract method 0x91cafe32.
//
// Solidity: function setForceBatchAddress(address newForceBatchAddress) returns()
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmTransactorSession) SetForceBatchAddress(newForceBatchAddress common.Address) (*types.Transaction, error) {
	return _Elderberrypolygonzkevm.Contract.SetForceBatchAddress(&_Elderberrypolygonzkevm.TransactOpts, newForceBatchAddress)
}

// SetForceBatchTimeout is a paid mutator transaction binding the contract method 0x4e487706.
//
// Solidity: function setForceBatchTimeout(uint64 newforceBatchTimeout) returns()
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmTransactor) SetForceBatchTimeout(opts *bind.TransactOpts, newforceBatchTimeout uint64) (*types.Transaction, error) {
	return _Elderberrypolygonzkevm.contract.Transact(opts, "setForceBatchTimeout", newforceBatchTimeout)
}

// SetForceBatchTimeout is a paid mutator transaction binding the contract method 0x4e487706.
//
// Solidity: function setForceBatchTimeout(uint64 newforceBatchTimeout) returns()
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmSession) SetForceBatchTimeout(newforceBatchTimeout uint64) (*types.Transaction, error) {
	return _Elderberrypolygonzkevm.Contract.SetForceBatchTimeout(&_Elderberrypolygonzkevm.TransactOpts, newforceBatchTimeout)
}

// SetForceBatchTimeout is a paid mutator transaction binding the contract method 0x4e487706.
//
// Solidity: function setForceBatchTimeout(uint64 newforceBatchTimeout) returns()
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmTransactorSession) SetForceBatchTimeout(newforceBatchTimeout uint64) (*types.Transaction, error) {
	return _Elderberrypolygonzkevm.Contract.SetForceBatchTimeout(&_Elderberrypolygonzkevm.TransactOpts, newforceBatchTimeout)
}

// SetTrustedSequencer is a paid mutator transaction binding the contract method 0x6ff512cc.
//
// Solidity: function setTrustedSequencer(address newTrustedSequencer) returns()
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmTransactor) SetTrustedSequencer(opts *bind.TransactOpts, newTrustedSequencer common.Address) (*types.Transaction, error) {
	return _Elderberrypolygonzkevm.contract.Transact(opts, "setTrustedSequencer", newTrustedSequencer)
}

// SetTrustedSequencer is a paid mutator transaction binding the contract method 0x6ff512cc.
//
// Solidity: function setTrustedSequencer(address newTrustedSequencer) returns()
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmSession) SetTrustedSequencer(newTrustedSequencer common.Address) (*types.Transaction, error) {
	return _Elderberrypolygonzkevm.Contract.SetTrustedSequencer(&_Elderberrypolygonzkevm.TransactOpts, newTrustedSequencer)
}

// SetTrustedSequencer is a paid mutator transaction binding the contract method 0x6ff512cc.
//
// Solidity: function setTrustedSequencer(address newTrustedSequencer) returns()
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmTransactorSession) SetTrustedSequencer(newTrustedSequencer common.Address) (*types.Transaction, error) {
	return _Elderberrypolygonzkevm.Contract.SetTrustedSequencer(&_Elderberrypolygonzkevm.TransactOpts, newTrustedSequencer)
}

// SetTrustedSequencerURL is a paid mutator transaction binding the contract method 0xc89e42df.
//
// Solidity: function setTrustedSequencerURL(string newTrustedSequencerURL) returns()
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmTransactor) SetTrustedSequencerURL(opts *bind.TransactOpts, newTrustedSequencerURL string) (*types.Transaction, error) {
	return _Elderberrypolygonzkevm.contract.Transact(opts, "setTrustedSequencerURL", newTrustedSequencerURL)
}

// SetTrustedSequencerURL is a paid mutator transaction binding the contract method 0xc89e42df.
//
// Solidity: function setTrustedSequencerURL(string newTrustedSequencerURL) returns()
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmSession) SetTrustedSequencerURL(newTrustedSequencerURL string) (*types.Transaction, error) {
	return _Elderberrypolygonzkevm.Contract.SetTrustedSequencerURL(&_Elderberrypolygonzkevm.TransactOpts, newTrustedSequencerURL)
}

// SetTrustedSequencerURL is a paid mutator transaction binding the contract method 0xc89e42df.
//
// Solidity: function setTrustedSequencerURL(string newTrustedSequencerURL) returns()
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmTransactorSession) SetTrustedSequencerURL(newTrustedSequencerURL string) (*types.Transaction, error) {
	return _Elderberrypolygonzkevm.Contract.SetTrustedSequencerURL(&_Elderberrypolygonzkevm.TransactOpts, newTrustedSequencerURL)
}

// TransferAdminRole is a paid mutator transaction binding the contract method 0xada8f919.
//
// Solidity: function transferAdminRole(address newPendingAdmin) returns()
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmTransactor) TransferAdminRole(opts *bind.TransactOpts, newPendingAdmin common.Address) (*types.Transaction, error) {
	return _Elderberrypolygonzkevm.contract.Transact(opts, "transferAdminRole", newPendingAdmin)
}

// TransferAdminRole is a paid mutator transaction binding the contract method 0xada8f919.
//
// Solidity: function transferAdminRole(address newPendingAdmin) returns()
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmSession) TransferAdminRole(newPendingAdmin common.Address) (*types.Transaction, error) {
	return _Elderberrypolygonzkevm.Contract.TransferAdminRole(&_Elderberrypolygonzkevm.TransactOpts, newPendingAdmin)
}

// TransferAdminRole is a paid mutator transaction binding the contract method 0xada8f919.
//
// Solidity: function transferAdminRole(address newPendingAdmin) returns()
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmTransactorSession) TransferAdminRole(newPendingAdmin common.Address) (*types.Transaction, error) {
	return _Elderberrypolygonzkevm.Contract.TransferAdminRole(&_Elderberrypolygonzkevm.TransactOpts, newPendingAdmin)
}

// ElderberrypolygonzkevmAcceptAdminRoleIterator is returned from FilterAcceptAdminRole and is used to iterate over the raw logs and unpacked data for AcceptAdminRole events raised by the Elderberrypolygonzkevm contract.
type ElderberrypolygonzkevmAcceptAdminRoleIterator struct {
	Event *ElderberrypolygonzkevmAcceptAdminRole // Event containing the contract specifics and raw log

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
func (it *ElderberrypolygonzkevmAcceptAdminRoleIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ElderberrypolygonzkevmAcceptAdminRole)
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
		it.Event = new(ElderberrypolygonzkevmAcceptAdminRole)
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
func (it *ElderberrypolygonzkevmAcceptAdminRoleIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ElderberrypolygonzkevmAcceptAdminRoleIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ElderberrypolygonzkevmAcceptAdminRole represents a AcceptAdminRole event raised by the Elderberrypolygonzkevm contract.
type ElderberrypolygonzkevmAcceptAdminRole struct {
	NewAdmin common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterAcceptAdminRole is a free log retrieval operation binding the contract event 0x056dc487bbf0795d0bbb1b4f0af523a855503cff740bfb4d5475f7a90c091e8e.
//
// Solidity: event AcceptAdminRole(address newAdmin)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmFilterer) FilterAcceptAdminRole(opts *bind.FilterOpts) (*ElderberrypolygonzkevmAcceptAdminRoleIterator, error) {

	logs, sub, err := _Elderberrypolygonzkevm.contract.FilterLogs(opts, "AcceptAdminRole")
	if err != nil {
		return nil, err
	}
	return &ElderberrypolygonzkevmAcceptAdminRoleIterator{contract: _Elderberrypolygonzkevm.contract, event: "AcceptAdminRole", logs: logs, sub: sub}, nil
}

// WatchAcceptAdminRole is a free log subscription operation binding the contract event 0x056dc487bbf0795d0bbb1b4f0af523a855503cff740bfb4d5475f7a90c091e8e.
//
// Solidity: event AcceptAdminRole(address newAdmin)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmFilterer) WatchAcceptAdminRole(opts *bind.WatchOpts, sink chan<- *ElderberrypolygonzkevmAcceptAdminRole) (event.Subscription, error) {

	logs, sub, err := _Elderberrypolygonzkevm.contract.WatchLogs(opts, "AcceptAdminRole")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ElderberrypolygonzkevmAcceptAdminRole)
				if err := _Elderberrypolygonzkevm.contract.UnpackLog(event, "AcceptAdminRole", log); err != nil {
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
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmFilterer) ParseAcceptAdminRole(log types.Log) (*ElderberrypolygonzkevmAcceptAdminRole, error) {
	event := new(ElderberrypolygonzkevmAcceptAdminRole)
	if err := _Elderberrypolygonzkevm.contract.UnpackLog(event, "AcceptAdminRole", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ElderberrypolygonzkevmForceBatchIterator is returned from FilterForceBatch and is used to iterate over the raw logs and unpacked data for ForceBatch events raised by the Elderberrypolygonzkevm contract.
type ElderberrypolygonzkevmForceBatchIterator struct {
	Event *ElderberrypolygonzkevmForceBatch // Event containing the contract specifics and raw log

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
func (it *ElderberrypolygonzkevmForceBatchIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ElderberrypolygonzkevmForceBatch)
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
		it.Event = new(ElderberrypolygonzkevmForceBatch)
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
func (it *ElderberrypolygonzkevmForceBatchIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ElderberrypolygonzkevmForceBatchIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ElderberrypolygonzkevmForceBatch represents a ForceBatch event raised by the Elderberrypolygonzkevm contract.
type ElderberrypolygonzkevmForceBatch struct {
	ForceBatchNum      uint64
	LastGlobalExitRoot [32]byte
	Sequencer          common.Address
	Transactions       []byte
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterForceBatch is a free log retrieval operation binding the contract event 0xf94bb37db835f1ab585ee00041849a09b12cd081d77fa15ca070757619cbc931.
//
// Solidity: event ForceBatch(uint64 indexed forceBatchNum, bytes32 lastGlobalExitRoot, address sequencer, bytes transactions)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmFilterer) FilterForceBatch(opts *bind.FilterOpts, forceBatchNum []uint64) (*ElderberrypolygonzkevmForceBatchIterator, error) {

	var forceBatchNumRule []interface{}
	for _, forceBatchNumItem := range forceBatchNum {
		forceBatchNumRule = append(forceBatchNumRule, forceBatchNumItem)
	}

	logs, sub, err := _Elderberrypolygonzkevm.contract.FilterLogs(opts, "ForceBatch", forceBatchNumRule)
	if err != nil {
		return nil, err
	}
	return &ElderberrypolygonzkevmForceBatchIterator{contract: _Elderberrypolygonzkevm.contract, event: "ForceBatch", logs: logs, sub: sub}, nil
}

// WatchForceBatch is a free log subscription operation binding the contract event 0xf94bb37db835f1ab585ee00041849a09b12cd081d77fa15ca070757619cbc931.
//
// Solidity: event ForceBatch(uint64 indexed forceBatchNum, bytes32 lastGlobalExitRoot, address sequencer, bytes transactions)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmFilterer) WatchForceBatch(opts *bind.WatchOpts, sink chan<- *ElderberrypolygonzkevmForceBatch, forceBatchNum []uint64) (event.Subscription, error) {

	var forceBatchNumRule []interface{}
	for _, forceBatchNumItem := range forceBatchNum {
		forceBatchNumRule = append(forceBatchNumRule, forceBatchNumItem)
	}

	logs, sub, err := _Elderberrypolygonzkevm.contract.WatchLogs(opts, "ForceBatch", forceBatchNumRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ElderberrypolygonzkevmForceBatch)
				if err := _Elderberrypolygonzkevm.contract.UnpackLog(event, "ForceBatch", log); err != nil {
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

// ParseForceBatch is a log parse operation binding the contract event 0xf94bb37db835f1ab585ee00041849a09b12cd081d77fa15ca070757619cbc931.
//
// Solidity: event ForceBatch(uint64 indexed forceBatchNum, bytes32 lastGlobalExitRoot, address sequencer, bytes transactions)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmFilterer) ParseForceBatch(log types.Log) (*ElderberrypolygonzkevmForceBatch, error) {
	event := new(ElderberrypolygonzkevmForceBatch)
	if err := _Elderberrypolygonzkevm.contract.UnpackLog(event, "ForceBatch", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ElderberrypolygonzkevmInitialSequenceBatchesIterator is returned from FilterInitialSequenceBatches and is used to iterate over the raw logs and unpacked data for InitialSequenceBatches events raised by the Elderberrypolygonzkevm contract.
type ElderberrypolygonzkevmInitialSequenceBatchesIterator struct {
	Event *ElderberrypolygonzkevmInitialSequenceBatches // Event containing the contract specifics and raw log

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
func (it *ElderberrypolygonzkevmInitialSequenceBatchesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ElderberrypolygonzkevmInitialSequenceBatches)
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
		it.Event = new(ElderberrypolygonzkevmInitialSequenceBatches)
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
func (it *ElderberrypolygonzkevmInitialSequenceBatchesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ElderberrypolygonzkevmInitialSequenceBatchesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ElderberrypolygonzkevmInitialSequenceBatches represents a InitialSequenceBatches event raised by the Elderberrypolygonzkevm contract.
type ElderberrypolygonzkevmInitialSequenceBatches struct {
	Transactions       []byte
	LastGlobalExitRoot [32]byte
	Sequencer          common.Address
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterInitialSequenceBatches is a free log retrieval operation binding the contract event 0x060116213bcbf54ca19fd649dc84b59ab2bbd200ab199770e4d923e222a28e7f.
//
// Solidity: event InitialSequenceBatches(bytes transactions, bytes32 lastGlobalExitRoot, address sequencer)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmFilterer) FilterInitialSequenceBatches(opts *bind.FilterOpts) (*ElderberrypolygonzkevmInitialSequenceBatchesIterator, error) {

	logs, sub, err := _Elderberrypolygonzkevm.contract.FilterLogs(opts, "InitialSequenceBatches")
	if err != nil {
		return nil, err
	}
	return &ElderberrypolygonzkevmInitialSequenceBatchesIterator{contract: _Elderberrypolygonzkevm.contract, event: "InitialSequenceBatches", logs: logs, sub: sub}, nil
}

// WatchInitialSequenceBatches is a free log subscription operation binding the contract event 0x060116213bcbf54ca19fd649dc84b59ab2bbd200ab199770e4d923e222a28e7f.
//
// Solidity: event InitialSequenceBatches(bytes transactions, bytes32 lastGlobalExitRoot, address sequencer)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmFilterer) WatchInitialSequenceBatches(opts *bind.WatchOpts, sink chan<- *ElderberrypolygonzkevmInitialSequenceBatches) (event.Subscription, error) {

	logs, sub, err := _Elderberrypolygonzkevm.contract.WatchLogs(opts, "InitialSequenceBatches")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ElderberrypolygonzkevmInitialSequenceBatches)
				if err := _Elderberrypolygonzkevm.contract.UnpackLog(event, "InitialSequenceBatches", log); err != nil {
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

// ParseInitialSequenceBatches is a log parse operation binding the contract event 0x060116213bcbf54ca19fd649dc84b59ab2bbd200ab199770e4d923e222a28e7f.
//
// Solidity: event InitialSequenceBatches(bytes transactions, bytes32 lastGlobalExitRoot, address sequencer)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmFilterer) ParseInitialSequenceBatches(log types.Log) (*ElderberrypolygonzkevmInitialSequenceBatches, error) {
	event := new(ElderberrypolygonzkevmInitialSequenceBatches)
	if err := _Elderberrypolygonzkevm.contract.UnpackLog(event, "InitialSequenceBatches", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ElderberrypolygonzkevmInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Elderberrypolygonzkevm contract.
type ElderberrypolygonzkevmInitializedIterator struct {
	Event *ElderberrypolygonzkevmInitialized // Event containing the contract specifics and raw log

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
func (it *ElderberrypolygonzkevmInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ElderberrypolygonzkevmInitialized)
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
		it.Event = new(ElderberrypolygonzkevmInitialized)
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
func (it *ElderberrypolygonzkevmInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ElderberrypolygonzkevmInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ElderberrypolygonzkevmInitialized represents a Initialized event raised by the Elderberrypolygonzkevm contract.
type ElderberrypolygonzkevmInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmFilterer) FilterInitialized(opts *bind.FilterOpts) (*ElderberrypolygonzkevmInitializedIterator, error) {

	logs, sub, err := _Elderberrypolygonzkevm.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &ElderberrypolygonzkevmInitializedIterator{contract: _Elderberrypolygonzkevm.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *ElderberrypolygonzkevmInitialized) (event.Subscription, error) {

	logs, sub, err := _Elderberrypolygonzkevm.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ElderberrypolygonzkevmInitialized)
				if err := _Elderberrypolygonzkevm.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmFilterer) ParseInitialized(log types.Log) (*ElderberrypolygonzkevmInitialized, error) {
	event := new(ElderberrypolygonzkevmInitialized)
	if err := _Elderberrypolygonzkevm.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ElderberrypolygonzkevmSequenceBatchesIterator is returned from FilterSequenceBatches and is used to iterate over the raw logs and unpacked data for SequenceBatches events raised by the Elderberrypolygonzkevm contract.
type ElderberrypolygonzkevmSequenceBatchesIterator struct {
	Event *ElderberrypolygonzkevmSequenceBatches // Event containing the contract specifics and raw log

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
func (it *ElderberrypolygonzkevmSequenceBatchesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ElderberrypolygonzkevmSequenceBatches)
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
		it.Event = new(ElderberrypolygonzkevmSequenceBatches)
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
func (it *ElderberrypolygonzkevmSequenceBatchesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ElderberrypolygonzkevmSequenceBatchesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ElderberrypolygonzkevmSequenceBatches represents a SequenceBatches event raised by the Elderberrypolygonzkevm contract.
type ElderberrypolygonzkevmSequenceBatches struct {
	NumBatch   uint64
	L1InfoRoot [32]byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterSequenceBatches is a free log retrieval operation binding the contract event 0x3e54d0825ed78523037d00a81759237eb436ce774bd546993ee67a1b67b6e766.
//
// Solidity: event SequenceBatches(uint64 indexed numBatch, bytes32 l1InfoRoot)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmFilterer) FilterSequenceBatches(opts *bind.FilterOpts, numBatch []uint64) (*ElderberrypolygonzkevmSequenceBatchesIterator, error) {

	var numBatchRule []interface{}
	for _, numBatchItem := range numBatch {
		numBatchRule = append(numBatchRule, numBatchItem)
	}

	logs, sub, err := _Elderberrypolygonzkevm.contract.FilterLogs(opts, "SequenceBatches", numBatchRule)
	if err != nil {
		return nil, err
	}
	return &ElderberrypolygonzkevmSequenceBatchesIterator{contract: _Elderberrypolygonzkevm.contract, event: "SequenceBatches", logs: logs, sub: sub}, nil
}

// WatchSequenceBatches is a free log subscription operation binding the contract event 0x3e54d0825ed78523037d00a81759237eb436ce774bd546993ee67a1b67b6e766.
//
// Solidity: event SequenceBatches(uint64 indexed numBatch, bytes32 l1InfoRoot)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmFilterer) WatchSequenceBatches(opts *bind.WatchOpts, sink chan<- *ElderberrypolygonzkevmSequenceBatches, numBatch []uint64) (event.Subscription, error) {

	var numBatchRule []interface{}
	for _, numBatchItem := range numBatch {
		numBatchRule = append(numBatchRule, numBatchItem)
	}

	logs, sub, err := _Elderberrypolygonzkevm.contract.WatchLogs(opts, "SequenceBatches", numBatchRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ElderberrypolygonzkevmSequenceBatches)
				if err := _Elderberrypolygonzkevm.contract.UnpackLog(event, "SequenceBatches", log); err != nil {
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

// ParseSequenceBatches is a log parse operation binding the contract event 0x3e54d0825ed78523037d00a81759237eb436ce774bd546993ee67a1b67b6e766.
//
// Solidity: event SequenceBatches(uint64 indexed numBatch, bytes32 l1InfoRoot)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmFilterer) ParseSequenceBatches(log types.Log) (*ElderberrypolygonzkevmSequenceBatches, error) {
	event := new(ElderberrypolygonzkevmSequenceBatches)
	if err := _Elderberrypolygonzkevm.contract.UnpackLog(event, "SequenceBatches", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ElderberrypolygonzkevmSequenceForceBatchesIterator is returned from FilterSequenceForceBatches and is used to iterate over the raw logs and unpacked data for SequenceForceBatches events raised by the Elderberrypolygonzkevm contract.
type ElderberrypolygonzkevmSequenceForceBatchesIterator struct {
	Event *ElderberrypolygonzkevmSequenceForceBatches // Event containing the contract specifics and raw log

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
func (it *ElderberrypolygonzkevmSequenceForceBatchesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ElderberrypolygonzkevmSequenceForceBatches)
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
		it.Event = new(ElderberrypolygonzkevmSequenceForceBatches)
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
func (it *ElderberrypolygonzkevmSequenceForceBatchesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ElderberrypolygonzkevmSequenceForceBatchesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ElderberrypolygonzkevmSequenceForceBatches represents a SequenceForceBatches event raised by the Elderberrypolygonzkevm contract.
type ElderberrypolygonzkevmSequenceForceBatches struct {
	NumBatch uint64
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterSequenceForceBatches is a free log retrieval operation binding the contract event 0x648a61dd2438f072f5a1960939abd30f37aea80d2e94c9792ad142d3e0a490a4.
//
// Solidity: event SequenceForceBatches(uint64 indexed numBatch)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmFilterer) FilterSequenceForceBatches(opts *bind.FilterOpts, numBatch []uint64) (*ElderberrypolygonzkevmSequenceForceBatchesIterator, error) {

	var numBatchRule []interface{}
	for _, numBatchItem := range numBatch {
		numBatchRule = append(numBatchRule, numBatchItem)
	}

	logs, sub, err := _Elderberrypolygonzkevm.contract.FilterLogs(opts, "SequenceForceBatches", numBatchRule)
	if err != nil {
		return nil, err
	}
	return &ElderberrypolygonzkevmSequenceForceBatchesIterator{contract: _Elderberrypolygonzkevm.contract, event: "SequenceForceBatches", logs: logs, sub: sub}, nil
}

// WatchSequenceForceBatches is a free log subscription operation binding the contract event 0x648a61dd2438f072f5a1960939abd30f37aea80d2e94c9792ad142d3e0a490a4.
//
// Solidity: event SequenceForceBatches(uint64 indexed numBatch)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmFilterer) WatchSequenceForceBatches(opts *bind.WatchOpts, sink chan<- *ElderberrypolygonzkevmSequenceForceBatches, numBatch []uint64) (event.Subscription, error) {

	var numBatchRule []interface{}
	for _, numBatchItem := range numBatch {
		numBatchRule = append(numBatchRule, numBatchItem)
	}

	logs, sub, err := _Elderberrypolygonzkevm.contract.WatchLogs(opts, "SequenceForceBatches", numBatchRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ElderberrypolygonzkevmSequenceForceBatches)
				if err := _Elderberrypolygonzkevm.contract.UnpackLog(event, "SequenceForceBatches", log); err != nil {
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

// ParseSequenceForceBatches is a log parse operation binding the contract event 0x648a61dd2438f072f5a1960939abd30f37aea80d2e94c9792ad142d3e0a490a4.
//
// Solidity: event SequenceForceBatches(uint64 indexed numBatch)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmFilterer) ParseSequenceForceBatches(log types.Log) (*ElderberrypolygonzkevmSequenceForceBatches, error) {
	event := new(ElderberrypolygonzkevmSequenceForceBatches)
	if err := _Elderberrypolygonzkevm.contract.UnpackLog(event, "SequenceForceBatches", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ElderberrypolygonzkevmSetForceBatchAddressIterator is returned from FilterSetForceBatchAddress and is used to iterate over the raw logs and unpacked data for SetForceBatchAddress events raised by the Elderberrypolygonzkevm contract.
type ElderberrypolygonzkevmSetForceBatchAddressIterator struct {
	Event *ElderberrypolygonzkevmSetForceBatchAddress // Event containing the contract specifics and raw log

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
func (it *ElderberrypolygonzkevmSetForceBatchAddressIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ElderberrypolygonzkevmSetForceBatchAddress)
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
		it.Event = new(ElderberrypolygonzkevmSetForceBatchAddress)
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
func (it *ElderberrypolygonzkevmSetForceBatchAddressIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ElderberrypolygonzkevmSetForceBatchAddressIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ElderberrypolygonzkevmSetForceBatchAddress represents a SetForceBatchAddress event raised by the Elderberrypolygonzkevm contract.
type ElderberrypolygonzkevmSetForceBatchAddress struct {
	NewForceBatchAddress common.Address
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterSetForceBatchAddress is a free log retrieval operation binding the contract event 0x5fbd7dd171301c4a1611a84aac4ba86d119478560557755f7927595b082634fb.
//
// Solidity: event SetForceBatchAddress(address newForceBatchAddress)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmFilterer) FilterSetForceBatchAddress(opts *bind.FilterOpts) (*ElderberrypolygonzkevmSetForceBatchAddressIterator, error) {

	logs, sub, err := _Elderberrypolygonzkevm.contract.FilterLogs(opts, "SetForceBatchAddress")
	if err != nil {
		return nil, err
	}
	return &ElderberrypolygonzkevmSetForceBatchAddressIterator{contract: _Elderberrypolygonzkevm.contract, event: "SetForceBatchAddress", logs: logs, sub: sub}, nil
}

// WatchSetForceBatchAddress is a free log subscription operation binding the contract event 0x5fbd7dd171301c4a1611a84aac4ba86d119478560557755f7927595b082634fb.
//
// Solidity: event SetForceBatchAddress(address newForceBatchAddress)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmFilterer) WatchSetForceBatchAddress(opts *bind.WatchOpts, sink chan<- *ElderberrypolygonzkevmSetForceBatchAddress) (event.Subscription, error) {

	logs, sub, err := _Elderberrypolygonzkevm.contract.WatchLogs(opts, "SetForceBatchAddress")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ElderberrypolygonzkevmSetForceBatchAddress)
				if err := _Elderberrypolygonzkevm.contract.UnpackLog(event, "SetForceBatchAddress", log); err != nil {
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

// ParseSetForceBatchAddress is a log parse operation binding the contract event 0x5fbd7dd171301c4a1611a84aac4ba86d119478560557755f7927595b082634fb.
//
// Solidity: event SetForceBatchAddress(address newForceBatchAddress)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmFilterer) ParseSetForceBatchAddress(log types.Log) (*ElderberrypolygonzkevmSetForceBatchAddress, error) {
	event := new(ElderberrypolygonzkevmSetForceBatchAddress)
	if err := _Elderberrypolygonzkevm.contract.UnpackLog(event, "SetForceBatchAddress", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ElderberrypolygonzkevmSetForceBatchTimeoutIterator is returned from FilterSetForceBatchTimeout and is used to iterate over the raw logs and unpacked data for SetForceBatchTimeout events raised by the Elderberrypolygonzkevm contract.
type ElderberrypolygonzkevmSetForceBatchTimeoutIterator struct {
	Event *ElderberrypolygonzkevmSetForceBatchTimeout // Event containing the contract specifics and raw log

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
func (it *ElderberrypolygonzkevmSetForceBatchTimeoutIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ElderberrypolygonzkevmSetForceBatchTimeout)
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
		it.Event = new(ElderberrypolygonzkevmSetForceBatchTimeout)
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
func (it *ElderberrypolygonzkevmSetForceBatchTimeoutIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ElderberrypolygonzkevmSetForceBatchTimeoutIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ElderberrypolygonzkevmSetForceBatchTimeout represents a SetForceBatchTimeout event raised by the Elderberrypolygonzkevm contract.
type ElderberrypolygonzkevmSetForceBatchTimeout struct {
	NewforceBatchTimeout uint64
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterSetForceBatchTimeout is a free log retrieval operation binding the contract event 0xa7eb6cb8a613eb4e8bddc1ac3d61ec6cf10898760f0b187bcca794c6ca6fa40b.
//
// Solidity: event SetForceBatchTimeout(uint64 newforceBatchTimeout)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmFilterer) FilterSetForceBatchTimeout(opts *bind.FilterOpts) (*ElderberrypolygonzkevmSetForceBatchTimeoutIterator, error) {

	logs, sub, err := _Elderberrypolygonzkevm.contract.FilterLogs(opts, "SetForceBatchTimeout")
	if err != nil {
		return nil, err
	}
	return &ElderberrypolygonzkevmSetForceBatchTimeoutIterator{contract: _Elderberrypolygonzkevm.contract, event: "SetForceBatchTimeout", logs: logs, sub: sub}, nil
}

// WatchSetForceBatchTimeout is a free log subscription operation binding the contract event 0xa7eb6cb8a613eb4e8bddc1ac3d61ec6cf10898760f0b187bcca794c6ca6fa40b.
//
// Solidity: event SetForceBatchTimeout(uint64 newforceBatchTimeout)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmFilterer) WatchSetForceBatchTimeout(opts *bind.WatchOpts, sink chan<- *ElderberrypolygonzkevmSetForceBatchTimeout) (event.Subscription, error) {

	logs, sub, err := _Elderberrypolygonzkevm.contract.WatchLogs(opts, "SetForceBatchTimeout")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ElderberrypolygonzkevmSetForceBatchTimeout)
				if err := _Elderberrypolygonzkevm.contract.UnpackLog(event, "SetForceBatchTimeout", log); err != nil {
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

// ParseSetForceBatchTimeout is a log parse operation binding the contract event 0xa7eb6cb8a613eb4e8bddc1ac3d61ec6cf10898760f0b187bcca794c6ca6fa40b.
//
// Solidity: event SetForceBatchTimeout(uint64 newforceBatchTimeout)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmFilterer) ParseSetForceBatchTimeout(log types.Log) (*ElderberrypolygonzkevmSetForceBatchTimeout, error) {
	event := new(ElderberrypolygonzkevmSetForceBatchTimeout)
	if err := _Elderberrypolygonzkevm.contract.UnpackLog(event, "SetForceBatchTimeout", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ElderberrypolygonzkevmSetTrustedSequencerIterator is returned from FilterSetTrustedSequencer and is used to iterate over the raw logs and unpacked data for SetTrustedSequencer events raised by the Elderberrypolygonzkevm contract.
type ElderberrypolygonzkevmSetTrustedSequencerIterator struct {
	Event *ElderberrypolygonzkevmSetTrustedSequencer // Event containing the contract specifics and raw log

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
func (it *ElderberrypolygonzkevmSetTrustedSequencerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ElderberrypolygonzkevmSetTrustedSequencer)
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
		it.Event = new(ElderberrypolygonzkevmSetTrustedSequencer)
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
func (it *ElderberrypolygonzkevmSetTrustedSequencerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ElderberrypolygonzkevmSetTrustedSequencerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ElderberrypolygonzkevmSetTrustedSequencer represents a SetTrustedSequencer event raised by the Elderberrypolygonzkevm contract.
type ElderberrypolygonzkevmSetTrustedSequencer struct {
	NewTrustedSequencer common.Address
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterSetTrustedSequencer is a free log retrieval operation binding the contract event 0xf54144f9611984021529f814a1cb6a41e22c58351510a0d9f7e822618abb9cc0.
//
// Solidity: event SetTrustedSequencer(address newTrustedSequencer)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmFilterer) FilterSetTrustedSequencer(opts *bind.FilterOpts) (*ElderberrypolygonzkevmSetTrustedSequencerIterator, error) {

	logs, sub, err := _Elderberrypolygonzkevm.contract.FilterLogs(opts, "SetTrustedSequencer")
	if err != nil {
		return nil, err
	}
	return &ElderberrypolygonzkevmSetTrustedSequencerIterator{contract: _Elderberrypolygonzkevm.contract, event: "SetTrustedSequencer", logs: logs, sub: sub}, nil
}

// WatchSetTrustedSequencer is a free log subscription operation binding the contract event 0xf54144f9611984021529f814a1cb6a41e22c58351510a0d9f7e822618abb9cc0.
//
// Solidity: event SetTrustedSequencer(address newTrustedSequencer)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmFilterer) WatchSetTrustedSequencer(opts *bind.WatchOpts, sink chan<- *ElderberrypolygonzkevmSetTrustedSequencer) (event.Subscription, error) {

	logs, sub, err := _Elderberrypolygonzkevm.contract.WatchLogs(opts, "SetTrustedSequencer")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ElderberrypolygonzkevmSetTrustedSequencer)
				if err := _Elderberrypolygonzkevm.contract.UnpackLog(event, "SetTrustedSequencer", log); err != nil {
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
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmFilterer) ParseSetTrustedSequencer(log types.Log) (*ElderberrypolygonzkevmSetTrustedSequencer, error) {
	event := new(ElderberrypolygonzkevmSetTrustedSequencer)
	if err := _Elderberrypolygonzkevm.contract.UnpackLog(event, "SetTrustedSequencer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ElderberrypolygonzkevmSetTrustedSequencerURLIterator is returned from FilterSetTrustedSequencerURL and is used to iterate over the raw logs and unpacked data for SetTrustedSequencerURL events raised by the Elderberrypolygonzkevm contract.
type ElderberrypolygonzkevmSetTrustedSequencerURLIterator struct {
	Event *ElderberrypolygonzkevmSetTrustedSequencerURL // Event containing the contract specifics and raw log

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
func (it *ElderberrypolygonzkevmSetTrustedSequencerURLIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ElderberrypolygonzkevmSetTrustedSequencerURL)
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
		it.Event = new(ElderberrypolygonzkevmSetTrustedSequencerURL)
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
func (it *ElderberrypolygonzkevmSetTrustedSequencerURLIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ElderberrypolygonzkevmSetTrustedSequencerURLIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ElderberrypolygonzkevmSetTrustedSequencerURL represents a SetTrustedSequencerURL event raised by the Elderberrypolygonzkevm contract.
type ElderberrypolygonzkevmSetTrustedSequencerURL struct {
	NewTrustedSequencerURL string
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterSetTrustedSequencerURL is a free log retrieval operation binding the contract event 0x6b8f723a4c7a5335cafae8a598a0aa0301be1387c037dccc085b62add6448b20.
//
// Solidity: event SetTrustedSequencerURL(string newTrustedSequencerURL)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmFilterer) FilterSetTrustedSequencerURL(opts *bind.FilterOpts) (*ElderberrypolygonzkevmSetTrustedSequencerURLIterator, error) {

	logs, sub, err := _Elderberrypolygonzkevm.contract.FilterLogs(opts, "SetTrustedSequencerURL")
	if err != nil {
		return nil, err
	}
	return &ElderberrypolygonzkevmSetTrustedSequencerURLIterator{contract: _Elderberrypolygonzkevm.contract, event: "SetTrustedSequencerURL", logs: logs, sub: sub}, nil
}

// WatchSetTrustedSequencerURL is a free log subscription operation binding the contract event 0x6b8f723a4c7a5335cafae8a598a0aa0301be1387c037dccc085b62add6448b20.
//
// Solidity: event SetTrustedSequencerURL(string newTrustedSequencerURL)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmFilterer) WatchSetTrustedSequencerURL(opts *bind.WatchOpts, sink chan<- *ElderberrypolygonzkevmSetTrustedSequencerURL) (event.Subscription, error) {

	logs, sub, err := _Elderberrypolygonzkevm.contract.WatchLogs(opts, "SetTrustedSequencerURL")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ElderberrypolygonzkevmSetTrustedSequencerURL)
				if err := _Elderberrypolygonzkevm.contract.UnpackLog(event, "SetTrustedSequencerURL", log); err != nil {
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
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmFilterer) ParseSetTrustedSequencerURL(log types.Log) (*ElderberrypolygonzkevmSetTrustedSequencerURL, error) {
	event := new(ElderberrypolygonzkevmSetTrustedSequencerURL)
	if err := _Elderberrypolygonzkevm.contract.UnpackLog(event, "SetTrustedSequencerURL", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ElderberrypolygonzkevmTransferAdminRoleIterator is returned from FilterTransferAdminRole and is used to iterate over the raw logs and unpacked data for TransferAdminRole events raised by the Elderberrypolygonzkevm contract.
type ElderberrypolygonzkevmTransferAdminRoleIterator struct {
	Event *ElderberrypolygonzkevmTransferAdminRole // Event containing the contract specifics and raw log

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
func (it *ElderberrypolygonzkevmTransferAdminRoleIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ElderberrypolygonzkevmTransferAdminRole)
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
		it.Event = new(ElderberrypolygonzkevmTransferAdminRole)
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
func (it *ElderberrypolygonzkevmTransferAdminRoleIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ElderberrypolygonzkevmTransferAdminRoleIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ElderberrypolygonzkevmTransferAdminRole represents a TransferAdminRole event raised by the Elderberrypolygonzkevm contract.
type ElderberrypolygonzkevmTransferAdminRole struct {
	NewPendingAdmin common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterTransferAdminRole is a free log retrieval operation binding the contract event 0xa5b56b7906fd0a20e3f35120dd8343db1e12e037a6c90111c7e42885e82a1ce6.
//
// Solidity: event TransferAdminRole(address newPendingAdmin)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmFilterer) FilterTransferAdminRole(opts *bind.FilterOpts) (*ElderberrypolygonzkevmTransferAdminRoleIterator, error) {

	logs, sub, err := _Elderberrypolygonzkevm.contract.FilterLogs(opts, "TransferAdminRole")
	if err != nil {
		return nil, err
	}
	return &ElderberrypolygonzkevmTransferAdminRoleIterator{contract: _Elderberrypolygonzkevm.contract, event: "TransferAdminRole", logs: logs, sub: sub}, nil
}

// WatchTransferAdminRole is a free log subscription operation binding the contract event 0xa5b56b7906fd0a20e3f35120dd8343db1e12e037a6c90111c7e42885e82a1ce6.
//
// Solidity: event TransferAdminRole(address newPendingAdmin)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmFilterer) WatchTransferAdminRole(opts *bind.WatchOpts, sink chan<- *ElderberrypolygonzkevmTransferAdminRole) (event.Subscription, error) {

	logs, sub, err := _Elderberrypolygonzkevm.contract.WatchLogs(opts, "TransferAdminRole")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ElderberrypolygonzkevmTransferAdminRole)
				if err := _Elderberrypolygonzkevm.contract.UnpackLog(event, "TransferAdminRole", log); err != nil {
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
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmFilterer) ParseTransferAdminRole(log types.Log) (*ElderberrypolygonzkevmTransferAdminRole, error) {
	event := new(ElderberrypolygonzkevmTransferAdminRole)
	if err := _Elderberrypolygonzkevm.contract.UnpackLog(event, "TransferAdminRole", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ElderberrypolygonzkevmUpdateEtrogSequenceIterator is returned from FilterUpdateEtrogSequence and is used to iterate over the raw logs and unpacked data for UpdateEtrogSequence events raised by the Elderberrypolygonzkevm contract.
type ElderberrypolygonzkevmUpdateEtrogSequenceIterator struct {
	Event *ElderberrypolygonzkevmUpdateEtrogSequence // Event containing the contract specifics and raw log

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
func (it *ElderberrypolygonzkevmUpdateEtrogSequenceIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ElderberrypolygonzkevmUpdateEtrogSequence)
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
		it.Event = new(ElderberrypolygonzkevmUpdateEtrogSequence)
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
func (it *ElderberrypolygonzkevmUpdateEtrogSequenceIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ElderberrypolygonzkevmUpdateEtrogSequenceIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ElderberrypolygonzkevmUpdateEtrogSequence represents a UpdateEtrogSequence event raised by the Elderberrypolygonzkevm contract.
type ElderberrypolygonzkevmUpdateEtrogSequence struct {
	NumBatch           uint64
	Transactions       []byte
	LastGlobalExitRoot [32]byte
	Sequencer          common.Address
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterUpdateEtrogSequence is a free log retrieval operation binding the contract event 0xd2c80353fc15ef62c6affc7cd6b7ab5b42c43290c50be3372e55ae552cecd19c.
//
// Solidity: event UpdateEtrogSequence(uint64 numBatch, bytes transactions, bytes32 lastGlobalExitRoot, address sequencer)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmFilterer) FilterUpdateEtrogSequence(opts *bind.FilterOpts) (*ElderberrypolygonzkevmUpdateEtrogSequenceIterator, error) {

	logs, sub, err := _Elderberrypolygonzkevm.contract.FilterLogs(opts, "UpdateEtrogSequence")
	if err != nil {
		return nil, err
	}
	return &ElderberrypolygonzkevmUpdateEtrogSequenceIterator{contract: _Elderberrypolygonzkevm.contract, event: "UpdateEtrogSequence", logs: logs, sub: sub}, nil
}

// WatchUpdateEtrogSequence is a free log subscription operation binding the contract event 0xd2c80353fc15ef62c6affc7cd6b7ab5b42c43290c50be3372e55ae552cecd19c.
//
// Solidity: event UpdateEtrogSequence(uint64 numBatch, bytes transactions, bytes32 lastGlobalExitRoot, address sequencer)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmFilterer) WatchUpdateEtrogSequence(opts *bind.WatchOpts, sink chan<- *ElderberrypolygonzkevmUpdateEtrogSequence) (event.Subscription, error) {

	logs, sub, err := _Elderberrypolygonzkevm.contract.WatchLogs(opts, "UpdateEtrogSequence")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ElderberrypolygonzkevmUpdateEtrogSequence)
				if err := _Elderberrypolygonzkevm.contract.UnpackLog(event, "UpdateEtrogSequence", log); err != nil {
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

// ParseUpdateEtrogSequence is a log parse operation binding the contract event 0xd2c80353fc15ef62c6affc7cd6b7ab5b42c43290c50be3372e55ae552cecd19c.
//
// Solidity: event UpdateEtrogSequence(uint64 numBatch, bytes transactions, bytes32 lastGlobalExitRoot, address sequencer)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmFilterer) ParseUpdateEtrogSequence(log types.Log) (*ElderberrypolygonzkevmUpdateEtrogSequence, error) {
	event := new(ElderberrypolygonzkevmUpdateEtrogSequence)
	if err := _Elderberrypolygonzkevm.contract.UnpackLog(event, "UpdateEtrogSequence", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// ElderberrypolygonzkevmVerifyBatchesIterator is returned from FilterVerifyBatches and is used to iterate over the raw logs and unpacked data for VerifyBatches events raised by the Elderberrypolygonzkevm contract.
type ElderberrypolygonzkevmVerifyBatchesIterator struct {
	Event *ElderberrypolygonzkevmVerifyBatches // Event containing the contract specifics and raw log

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
func (it *ElderberrypolygonzkevmVerifyBatchesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ElderberrypolygonzkevmVerifyBatches)
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
		it.Event = new(ElderberrypolygonzkevmVerifyBatches)
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
func (it *ElderberrypolygonzkevmVerifyBatchesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ElderberrypolygonzkevmVerifyBatchesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ElderberrypolygonzkevmVerifyBatches represents a VerifyBatches event raised by the Elderberrypolygonzkevm contract.
type ElderberrypolygonzkevmVerifyBatches struct {
	NumBatch   uint64
	StateRoot  [32]byte
	Aggregator common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterVerifyBatches is a free log retrieval operation binding the contract event 0x9c72852172521097ba7e1482e6b44b351323df0155f97f4ea18fcec28e1f5966.
//
// Solidity: event VerifyBatches(uint64 indexed numBatch, bytes32 stateRoot, address indexed aggregator)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmFilterer) FilterVerifyBatches(opts *bind.FilterOpts, numBatch []uint64, aggregator []common.Address) (*ElderberrypolygonzkevmVerifyBatchesIterator, error) {

	var numBatchRule []interface{}
	for _, numBatchItem := range numBatch {
		numBatchRule = append(numBatchRule, numBatchItem)
	}

	var aggregatorRule []interface{}
	for _, aggregatorItem := range aggregator {
		aggregatorRule = append(aggregatorRule, aggregatorItem)
	}

	logs, sub, err := _Elderberrypolygonzkevm.contract.FilterLogs(opts, "VerifyBatches", numBatchRule, aggregatorRule)
	if err != nil {
		return nil, err
	}
	return &ElderberrypolygonzkevmVerifyBatchesIterator{contract: _Elderberrypolygonzkevm.contract, event: "VerifyBatches", logs: logs, sub: sub}, nil
}

// WatchVerifyBatches is a free log subscription operation binding the contract event 0x9c72852172521097ba7e1482e6b44b351323df0155f97f4ea18fcec28e1f5966.
//
// Solidity: event VerifyBatches(uint64 indexed numBatch, bytes32 stateRoot, address indexed aggregator)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmFilterer) WatchVerifyBatches(opts *bind.WatchOpts, sink chan<- *ElderberrypolygonzkevmVerifyBatches, numBatch []uint64, aggregator []common.Address) (event.Subscription, error) {

	var numBatchRule []interface{}
	for _, numBatchItem := range numBatch {
		numBatchRule = append(numBatchRule, numBatchItem)
	}

	var aggregatorRule []interface{}
	for _, aggregatorItem := range aggregator {
		aggregatorRule = append(aggregatorRule, aggregatorItem)
	}

	logs, sub, err := _Elderberrypolygonzkevm.contract.WatchLogs(opts, "VerifyBatches", numBatchRule, aggregatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ElderberrypolygonzkevmVerifyBatches)
				if err := _Elderberrypolygonzkevm.contract.UnpackLog(event, "VerifyBatches", log); err != nil {
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

// ParseVerifyBatches is a log parse operation binding the contract event 0x9c72852172521097ba7e1482e6b44b351323df0155f97f4ea18fcec28e1f5966.
//
// Solidity: event VerifyBatches(uint64 indexed numBatch, bytes32 stateRoot, address indexed aggregator)
func (_Elderberrypolygonzkevm *ElderberrypolygonzkevmFilterer) ParseVerifyBatches(log types.Log) (*ElderberrypolygonzkevmVerifyBatches, error) {
	event := new(ElderberrypolygonzkevmVerifyBatches)
	if err := _Elderberrypolygonzkevm.contract.UnpackLog(event, "VerifyBatches", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
