// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package feijoapolygonrollupmanager

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

// PolygonRollupManagerPendingStateSequenceBased is an auto generated low-level Go binding around an user-defined struct.
type PolygonRollupManagerPendingStateSequenceBased struct {
	Timestamp            uint64
	LastVerifiedSequence uint64
	ExitRoot             [32]byte
	StateRoot            [32]byte
}

// PolygonRollupManagerSequencedData is an auto generated low-level Go binding around an user-defined struct.
type PolygonRollupManagerSequencedData struct {
	AccInputHash       [32]byte
	SequencedTimestamp uint64
	CurrentBlobNum     uint64
	AccZkGasLimit      *big.Int
}

// PolygonRollupManagerVerifySequenceData is an auto generated low-level Go binding around an user-defined struct.
type PolygonRollupManagerVerifySequenceData struct {
	RollupID         uint32
	PendingStateNum  uint64
	InitSequenceNum  uint64
	FinalSequenceNum uint64
	NewLocalExitRoot [32]byte
	NewStateRoot     [32]byte
}

// FeijoapolygonrollupmanagerMetaData contains all meta data concerning the Feijoapolygonrollupmanager contract.
var FeijoapolygonrollupmanagerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIPolygonZkEVMGlobalExitRootV2\",\"name\":\"_globalExitRootManager\",\"type\":\"address\"},{\"internalType\":\"contractIERC20Upgradeable\",\"name\":\"_pol\",\"type\":\"address\"},{\"internalType\":\"contractIPolygonZkEVMBridge\",\"name\":\"_bridgeAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessControlOnlyCanRenounceRolesForSelf\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AddressDoNotHaveRequiredRole\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AllSequencedMustBeVerified\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AllzkEVMSequencedBatchesMustBeVerified\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BatchFeeOutOfRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CannotUpdateWithUnconsolidatedPendingState\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ChainIDAlreadyExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ChainIDOutOfRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ExceedMaxVerifyBatches\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FinalNumBatchBelowLastVerifiedBatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FinalNumBatchDoesNotMatchPendingState\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FinalNumSequenceBelowLastVerifiedSequence\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FinalNumSequenceDoesNotMatchPendingState\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FinalPendingStateNumInvalid\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"HaltTimeoutNotExpired\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InitBatchMustMatchCurrentForkID\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InitNumBatchAboveLastVerifiedBatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InitNumBatchDoesNotMatchPendingState\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InitSequenceMustMatchCurrentForkID\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InitSequenceNumDoesNotMatchPendingState\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidProof\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRangeBatchTimeTarget\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRangeMultiplierBatchFee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRangeMultiplierZkGasPrice\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRangeSequenceTimeTarget\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustSequenceSomeBatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustSequenceSomeBlob\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NewAccInputHashDoesNotExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NewPendingStateTimeoutMustBeLower\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NewStateRootNotInsidePrime\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NewTrustedAggregatorTimeoutMustBeLower\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OldAccInputHashDoesNotExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OldStateRootDoesNotExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyEmergencyState\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyNotEmergencyState\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyRollupAdmin\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PendingStateDoesNotExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PendingStateInvalid\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PendingStateNotConsolidable\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RollupAddressAlreadyExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RollupIDNotAscendingOrder\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RollupMustExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RollupTypeDoesNotExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RollupTypeObsolete\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SenderMustBeRollup\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"StoredRootMustBeDifferentThanNewRoot\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TrustedAggregatorTimeoutNotExpired\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpdateNotCompatible\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpdateToSameRollupTypeID\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"zkGasPriceOfRange\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"forkID\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"rollupAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"chainID\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"rollupCompatibilityID\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"lastVerifiedSequenceBeforeUpgrade\",\"type\":\"uint64\"}],\"name\":\"AddExistingRollup\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"rollupTypeID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consensusImplementation\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"verifier\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"forkID\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"rollupCompatibilityID\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"genesis\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"description\",\"type\":\"string\"}],\"name\":\"AddNewRollupType\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"numSequence\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"exitRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"pendingStateNum\",\"type\":\"uint64\"}],\"name\":\"ConsolidatePendingState\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"rollupTypeID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"rollupAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"chainID\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"gasTokenAddress\",\"type\":\"address\"}],\"name\":\"CreateNewRollup\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"EmergencyStateActivated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"EmergencyStateDeactivated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"rollupTypeID\",\"type\":\"uint32\"}],\"name\":\"ObsoleteRollupType\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint128\",\"name\":\"zkGasLimit\",\"type\":\"uint128\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"blobsSequenced\",\"type\":\"uint64\"}],\"name\":\"OnSequence\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"numSequence\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"exitRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"aggregator\",\"type\":\"address\"}],\"name\":\"OverridePendingState\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"storedStateRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"provedStateRoot\",\"type\":\"bytes32\"}],\"name\":\"ProveNonDeterministicPendingState\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"contractIVerifierRollup\",\"name\":\"aggregateRollupVerifier\",\"type\":\"address\"}],\"name\":\"SetAggregateRollupVerifier\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"newMultiplierSequenceFee\",\"type\":\"uint16\"}],\"name\":\"SetMultiplierZkGasPrice\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"newPendingStateTimeout\",\"type\":\"uint64\"}],\"name\":\"SetPendingStateTimeout\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newSequenceFee\",\"type\":\"uint256\"}],\"name\":\"SetSequenceFee\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newTrustedAggregator\",\"type\":\"address\"}],\"name\":\"SetTrustedAggregator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"newTrustedAggregatorTimeout\",\"type\":\"uint64\"}],\"name\":\"SetTrustedAggregatorTimeout\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"newVerifySequenceTimeTarget\",\"type\":\"uint64\"}],\"name\":\"SetVerifySequenceTimeTarget\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"newRollupTypeID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"lastVerifiedSequenceBeforeUpgrade\",\"type\":\"uint64\"}],\"name\":\"UpdateRollup\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"sequenceNum\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"exitRoot\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"aggregator\",\"type\":\"address\"}],\"name\":\"VerifySequences\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"aggregator\",\"type\":\"address\"}],\"name\":\"VerifySequencesMultiProof\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"numSequence\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"exitRoot\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"aggregator\",\"type\":\"address\"}],\"name\":\"VerifySequencesTrustedAggregator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"aggregator\",\"type\":\"address\"}],\"name\":\"VerifySequencesTrustedAggregatorMultiProof\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ZK_GAS_LIMIT_BATCH\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"activateEmergencyState\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIPolygonRollupBaseFeijoa\",\"name\":\"rollupAddress\",\"type\":\"address\"},{\"internalType\":\"contractIVerifierRollup\",\"name\":\"verifier\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"forkID\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"chainID\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"genesis\",\"type\":\"bytes32\"},{\"internalType\":\"uint8\",\"name\":\"rollupCompatibilityID\",\"type\":\"uint8\"}],\"name\":\"addExistingRollup\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"consensusImplementation\",\"type\":\"address\"},{\"internalType\":\"contractIVerifierRollup\",\"name\":\"verifier\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"forkID\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"rollupCompatibilityID\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"genesis\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"description\",\"type\":\"string\"}],\"name\":\"addNewRollupType\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"aggregateRollupVerifier\",\"outputs\":[{\"internalType\":\"contractIVerifierRollup\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"bridgeAddress\",\"outputs\":[{\"internalType\":\"contractIPolygonZkEVMBridge\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"calculateRewardPerZkGas\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainID\",\"type\":\"uint64\"}],\"name\":\"chainIDToRollupID\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"pendingStateNum\",\"type\":\"uint64\"}],\"name\":\"consolidatePendingState\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupTypeID\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"chainID\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"sequencer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"gasTokenAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"sequencerURL\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"networkName\",\"type\":\"string\"}],\"name\":\"createNewRollup\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"deactivateEmergencyState\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getForcedZkGasPrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"}],\"name\":\"getLastVerifiedSequence\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRollupExitRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNum\",\"type\":\"uint64\"}],\"name\":\"getRollupPendingStateTransitions\",\"outputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"timestamp\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"lastVerifiedSequence\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"exitRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"}],\"internalType\":\"structPolygonRollupManager.PendingStateSequenceBased\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNum\",\"type\":\"uint64\"}],\"name\":\"getRollupSequencedSequences\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"accInputHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sequencedTimestamp\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"currentBlobNum\",\"type\":\"uint64\"},{\"internalType\":\"uint128\",\"name\":\"accZkGasLimit\",\"type\":\"uint128\"}],\"internalType\":\"structPolygonRollupManager.SequencedData\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"sequenceNum\",\"type\":\"uint64\"}],\"name\":\"getRollupsequenceNumToStateRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getZkGasPrice\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"globalExitRootManager\",\"outputs\":[{\"internalType\":\"contractIPolygonZkEVMGlobalExitRootV2\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isEmergencyState\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"pendingStateNum\",\"type\":\"uint64\"}],\"name\":\"isPendingStateConsolidable\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastAggregationTimestamp\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastDeactivatedEmergencyStateTimestamp\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"multiplierZkGasPrice\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupTypeID\",\"type\":\"uint32\"}],\"name\":\"obsoleteRollupType\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint128\",\"name\":\"zkGasLimitSequenced\",\"type\":\"uint128\"},{\"internalType\":\"uint64\",\"name\":\"blobsSequenced\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"newAccInputHash\",\"type\":\"bytes32\"}],\"name\":\"onSequence\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"initPendingStateNum\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"finalPendingStateNum\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"initSequenceNum\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"finalSequenceNum\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"newLocalExitRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[24]\",\"name\":\"proof\",\"type\":\"bytes32[24]\"}],\"name\":\"overridePendingState\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pendingStateTimeout\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pol\",\"outputs\":[{\"internalType\":\"contractIERC20Upgradeable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"initPendingStateNum\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"finalPendingStateNum\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"initSequenceNum\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"finalSequenceNum\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"newLocalExitRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[24]\",\"name\":\"proof\",\"type\":\"bytes32[24]\"}],\"name\":\"proveNonDeterministicPendingState\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"rollupAddress\",\"type\":\"address\"}],\"name\":\"rollupAddressToID\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rollupCount\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"}],\"name\":\"rollupIDToRollupData\",\"outputs\":[{\"internalType\":\"contractIPolygonRollupBaseFeijoa\",\"name\":\"rollupContract\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"chainID\",\"type\":\"uint64\"},{\"internalType\":\"contractIVerifierRollup\",\"name\":\"verifier\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"forkID\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"lastLocalExitRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"lastSequenceNum\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"lastVerifiedSequenceNum\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"lastPendingState\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"lastPendingStateConsolidated\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"lastVerifiedSequenceBeforeUpgrade\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"rollupTypeID\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"rollupCompatibilityID\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rollupTypeCount\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupTypeID\",\"type\":\"uint32\"}],\"name\":\"rollupTypeMap\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"consensusImplementation\",\"type\":\"address\"},{\"internalType\":\"contractIVerifierRollup\",\"name\":\"verifier\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"forkID\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"rollupCompatibilityID\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"obsolete\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"genesis\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIVerifierRollup\",\"name\":\"newAggregateRollupVerifier\",\"type\":\"address\"}],\"name\":\"setAggregateRollupVerifier\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"newMultiplierZkGasPrice\",\"type\":\"uint16\"}],\"name\":\"setMultiplierZkGasPrice\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"newPendingStateTimeout\",\"type\":\"uint64\"}],\"name\":\"setPendingStateTimeout\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"newTrustedAggregatorTimeout\",\"type\":\"uint64\"}],\"name\":\"setTrustedAggregatorTimeout\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"newVerifySequenceTimeTarget\",\"type\":\"uint64\"}],\"name\":\"setVerifySequenceTimeTarget\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newZkGasPrice\",\"type\":\"uint256\"}],\"name\":\"setZkGasPrice\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalVerifiedZkGasLimit\",\"outputs\":[{\"internalType\":\"uint128\",\"name\":\"\",\"type\":\"uint128\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalZkGasLimit\",\"outputs\":[{\"internalType\":\"uint128\",\"name\":\"\",\"type\":\"uint128\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"trustedAggregatorTimeout\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractITransparentUpgradeableProxy\",\"name\":\"rollupContract\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"newRollupTypeID\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"upgradeData\",\"type\":\"bytes\"}],\"name\":\"updateRollup\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractITransparentUpgradeableProxy\",\"name\":\"rollupContract\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"newRollupTypeID\",\"type\":\"uint32\"}],\"name\":\"updateRollupByRollupAdmin\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"verifySequenceTimeTarget\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"pendingStateNum\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"initSequenceNum\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"finalSequenceNum\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"newLocalExitRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"}],\"internalType\":\"structPolygonRollupManager.VerifySequenceData[]\",\"name\":\"verifySequencesData\",\"type\":\"tuple[]\"},{\"internalType\":\"address\",\"name\":\"beneficiary\",\"type\":\"address\"},{\"internalType\":\"bytes32[24]\",\"name\":\"proof\",\"type\":\"bytes32[24]\"}],\"name\":\"verifySequencesMultiProof\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"pendingStateNum\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"initSequenceNum\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"finalSequenceNum\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"newLocalExitRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"}],\"internalType\":\"structPolygonRollupManager.VerifySequenceData[]\",\"name\":\"verifySequencesData\",\"type\":\"tuple[]\"},{\"internalType\":\"address\",\"name\":\"beneficiary\",\"type\":\"address\"},{\"internalType\":\"bytes32[24]\",\"name\":\"proof\",\"type\":\"bytes32[24]\"}],\"name\":\"verifySequencesTrustedAggregatorMultiProof\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60e060405234801562000010575f80fd5b506040516200600c3803806200600c833981016040819052620000339162000136565b6001600160a01b0380841660805282811660c052811660a052620000566200005f565b50505062000187565b5f54610100900460ff1615620000cb5760405162461bcd60e51b815260206004820152602760248201527f496e697469616c697a61626c653a20636f6e747261637420697320696e697469604482015266616c697a696e6760c81b606482015260840160405180910390fd5b5f5460ff90811610156200011c575f805460ff191660ff9081179091556040519081527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b565b6001600160a01b038116811462000133575f80fd5b50565b5f805f6060848603121562000149575f80fd5b835162000156816200011e565b602085015190935062000169816200011e565b60408501519092506200017c816200011e565b809150509250925092565b60805160a05160c051615e22620001ea5f395f8181610ae501528181611bf60152613d7701525f818161078d0152818161229601526135ec01525f8181610a2c01528181610dda0152818161250c01528181612bae01526134e10152615e225ff3fe608060405234801562000010575f80fd5b5060043610620003b0575f3560e01c80639ff22cb511620001ef578063d939b3151162000113578063eb142b4011620000ab578063f4174a171162000083578063f4174a171462000b78578063f4e926751462000b81578063f9c4c2ae1462000b92578063fe01d89e1462000ca8575f80fd5b8063eb142b401462000b07578063f00bdaa41462000b4a578063f34eb8eb1462000b61575f80fd5b8063dfdb8c5e11620000eb578063dfdb8c5e1462000a9a578063e0bfd3d21462000ab1578063e2bfe8b31462000ac8578063e46761c41462000adf575f80fd5b8063d939b3151462000a65578063dbc169761462000a79578063de7948501462000a83575f80fd5b8063b99d0ad71162000187578063c4c928c2116200015f578063c4c928c214620009e7578063ceee281d14620009fe578063d02103ca1462000a26578063d547741f1462000a4e575f80fd5b8063b99d0ad714620008d7578063ba988cef14620009b1578063c1acbc3414620009cc575f80fd5b8063a2967d9911620001c7578063a2967d99146200077d578063a3c573eb1462000787578063a9a7703114620007c8578063b739753614620008bc575f80fd5b80639ff22cb51462000734578063a1094df3146200075e578063a217fddf1462000775575f80fd5b806365c0504d11620002d75780638185f9d3116200026f5780638bd4f07111620002475780638bd4f07114620006c157806390031d5c14620006d857806391d1485414620006e25780639c9f3dfe146200071d575f80fd5b80638185f9d31462000683578063838a2503146200069a578063841b24d714620006a6575f80fd5b8063727885e911620002af578063727885e914620006235780637ec31def146200063a5780637fb6e76a14620006515780638129fc1c1462000679575f80fd5b806365c0504d146200054a5780636c6be9eb14620005f85780637222020f146200060c575f80fd5b80632072f6c5116200034b5780632f2ff15d11620003235780632f2ff15d14620004f157806330c27dde146200050857806336568abe146200051c578063394218e91462000533575f80fd5b80632072f6c5146200048e578063248a9ca3146200049857806327696c5e14620004bd575f80fd5b806312b86e19116200038b57806312b86e19146200042957806315064c9614620004425780631608859c14620004505780631796a1ae1462000467575f80fd5b806302f3fa6014620003b4578063080b311114620003d15780630a7eef7a14620003f9575b5f80fd5b620003be62000cbf565b6040519081526020015b60405180910390f35b620003e8620003e236600462004a5d565b62000cd6565b6040519015158152602001620003c8565b620004106200040a36600462004a93565b62000cff565b6040516001600160401b039091168152602001620003c8565b620004406200043a36600462004ac1565b62000d1e565b005b606f54620003e89060ff1681565b620004406200046136600462004a5d565b62000ed9565b607e54620004789063ffffffff1681565b60405163ffffffff9091168152602001620003c8565b6200044062000f83565b620003be620004a936600462004b53565b5f9081526034602052604090206001015490565b608954620004d890600160801b90046001600160801b031681565b6040516001600160801b039091168152602001620003c8565b620004406200050236600462004b80565b6200105f565b60875462000410906001600160401b031681565b620004406200052d36600462004b80565b62001087565b620004406200054436600462004bb1565b620010c1565b620005ae6200055b36600462004a93565b607f6020525f90815260409020805460018201546002909201546001600160a01b0391821692918216916001600160401b03600160a01b8204169160ff600160e01b8304811692600160e81b9004169086565b604080516001600160a01b0397881681529690951660208701526001600160401b039093169385019390935260ff166060840152901515608083015260a082015260c001620003c8565b608954620004d8906001600160801b031681565b620004406200061d36600462004a93565b6200118d565b620004406200063436600462004c7b565b62001287565b620004406200064b36600462004b53565b62001713565b620004786200066236600462004bb1565b60836020525f908152604090205463ffffffff1681565b62000440620017ab565b620004406200069436600462004bb1565b62001a68565b620004106305f5e10081565b6084546200041090600160c01b90046001600160401b031681565b62000440620006d236600462004ac1565b62001b1e565b620003be62001bd5565b620003e8620006f336600462004b80565b5f9182526034602090815260408084206001600160a01b0393909316845291905290205460ff1690565b620004406200072e36600462004bb1565b62001cb6565b6085546200074a90600160801b900461ffff1681565b60405161ffff9091168152602001620003c8565b620004406200076f36600462004d41565b62001d6d565b620003be5f81565b620003be62001e30565b620007af7f000000000000000000000000000000000000000000000000000000000000000081565b6040516001600160a01b039091168152602001620003c8565b62000871620007d936600462004a5d565b604080516080810182525f8082526020820181905291810182905260608101919091525063ffffffff919091165f9081526088602090815260408083206001600160401b03948516845260030182529182902082516080810184528154815260019091015480851692820192909252600160401b820490931691830191909152600160801b90046001600160801b0316606082015290565b60408051825181526020808401516001600160401b03908116918301919091528383015116918101919091526060918201516001600160801b031691810191909152608001620003c8565b6085546200041090600160401b90046001600160401b031681565b6200096c620008e836600462004a5d565b60408051608080820183525f8083526020808401829052838501829052606093840182905263ffffffff969096168152608886528381206001600160401b03958616825260040186528390208351918201845280548086168352600160401b9004909416948101949094526001830154918401919091526002909101549082015290565b604051620003c891905f6080820190506001600160401b0380845116835280602085015116602084015250604083015160408301526060830151606083015292915050565b608754620007af90600160401b90046001600160a01b031681565b6084546200041090600160801b90046001600160401b031681565b62000440620009f836600462004d6b565b620021d5565b6200047862000a0f36600462004de2565b60826020525f908152604090205463ffffffff1681565b620007af7f000000000000000000000000000000000000000000000000000000000000000081565b6200044062000a5f36600462004b80565b62002214565b60855462000410906001600160401b031681565b620004406200223c565b6200044062000a9436600462004e00565b62002308565b6200044062000aab36600462004e98565b6200258d565b6200044062000ac236600462004ed8565b620026de565b6200044062000ad936600462004de2565b620027b1565b620007af7f000000000000000000000000000000000000000000000000000000000000000081565b620003be62000b1836600462004a5d565b63ffffffff82165f9081526088602090815260408083206001600160401b038516845260020190915290205492915050565b6200044062000b5b36600462004e00565b6200284d565b6200044062000b7236600462004f50565b62002c39565b608a54620003be565b608054620004789063ffffffff1681565b62000c2862000ba336600462004a93565b60886020525f9081526040902080546001820154600583015460068401546007909401546001600160a01b0380851695600160a01b958690046001600160401b039081169692861695929092048216939282821692600160401b808404821693600160801b808204841694600160c01b90920484169380831693830416910460ff168c565b604080516001600160a01b039d8e1681526001600160401b039c8d1660208201529c909a16998c019990995296891660608b015260808a019590955292871660a089015290861660c0880152851660e0870152841661010086015283166101208501529190911661014083015260ff1661016082015261018001620003c8565b6200041062000cb936600462004fe2565b62002e2d565b5f608a54606462000cd1919062005040565b905090565b63ffffffff82165f90815260886020526040812062000cf69083620030b5565b90505b92915050565b63ffffffff81165f90815260886020526040812062000cf990620030f9565b7f084e94f375e9d647f87f5b2ceffba1e062c70f6009fdbcf80291e803b5c9edd462000d4a8162003168565b63ffffffff89165f90815260886020526040902062000d70818a8a8a8a8a8a8a62003174565b60068101805467ffffffffffffffff60401b1916600160401b6001600160401b038981169182029290921783555f9081526002840160205260409020869055600583018790559054600160801b9004161562000dd8576006810180546001600160801b031690555b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166333d6247d62000e1162001e30565b6040518263ffffffff1660e01b815260040162000e3091815260200190565b5f604051808303815f87803b15801562000e48575f80fd5b505af115801562000e5b573d5f803e3d5ffd5b5050608480546001600160c01b031661127560c71b1790555050604080516001600160401b03881681526020810186905290810186905233606082015263ffffffff8b16907f3182bd6e6f74fc1fdc88b60f3a4f4c7f79db6ae6f5b88a1b3f5a1e28ec210d5e9060800160405180910390a250505050505050505050565b63ffffffff82165f9081526088602090815260408083203384527fc17b14a573f65366cdad721c7c0a0f76536bb4a86b935cdac44610e4f010b52a9092529091205460ff1662000f7257606f5460ff161562000f4857604051630bc011ff60e21b815260040160405180910390fd5b62000f548183620030b5565b62000f7257604051630674f25160e11b815260040160405180910390fd5b62000f7e8183620033fe565b505050565b335f9081527f8875b94af5657a2903def9906d67a3f42d8a836d24b5602c00f00fc855339fcd602052604090205460ff166200105357608454600160801b90046001600160401b03161580620010045750608454429062000ff99062093a8090600160801b90046001600160401b03166200505a565b6001600160401b0316115b806200103457506087544290620010299062093a80906001600160401b03166200505a565b6001600160401b0316115b15620010535760405163692baaad60e11b815260040160405180910390fd5b6200105d620035ea565b565b5f828152603460205260409020600101546200107b8162003168565b62000f7e838362003664565b6001600160a01b0381163314620010b157604051630b4ad1cd60e31b815260040160405180910390fd5b620010bd8282620036e8565b5050565b7fa5c5790f581d443ed43873ab47cfb8c5d66a6db268e58b5971bb33fc66e07db1620010ed8162003168565b606f5460ff166200112f576084546001600160401b03600160c01b9091048116908316106200112f5760405163401636df60e01b815260040160405180910390fd5b608480546001600160c01b0316600160c01b6001600160401b038516908102919091179091556040519081527f1f4fa24c2e4bad19a7f3ec5c5485f70d46c798461c2e684f55bbd0fc661373a1906020015b60405180910390a15050565b7fab66e11c4f712cd06ab11bf9339b48bef39e12d4a22eeef71d2860a0c90482bd620011b98162003168565b63ffffffff82161580620011d85750607e5463ffffffff908116908316115b15620011f757604051637512e5cb60e01b815260040160405180910390fd5b63ffffffff82165f908152607f60205260409020600180820154600160e81b900460ff16151590036200123d57604051633b8d3d9960e01b815260040160405180910390fd5b60018101805460ff60e81b1916600160e81b17905560405163ffffffff8416907f4710d2ee567ef1ed6eb2f651dde4589524bcf7cebc62147a99b281cc836e7e44905f90a2505050565b7fa0fab074aba36a6fa69f1a83ee86e5abfb8433966eb57efb13dc2fc2f24ddd08620012b38162003168565b63ffffffff88161580620012d25750607e5463ffffffff908116908916115b15620012f157604051637512e5cb60e01b815260040160405180910390fd5b63ffffffff88165f908152607f60205260409020600180820154600160e81b900460ff16151590036200133757604051633b8d3d9960e01b815260040160405180910390fd5b63ffffffff6001600160401b03891611156200136657604051634c753f5760e01b815260040160405180910390fd5b6001600160401b0388165f9081526083602052604090205463ffffffff1615620013a3576040516337c8fe0960e11b815260040160405180910390fd5b608080545f91908290620013bd9063ffffffff1662005084565b825463ffffffff8281166101009490940a93840293021916919091179091558254604080515f80825260208201928390529394506001600160a01b039092169130916200140a9062004a24565b6200141893929190620050fa565b604051809103905ff08015801562001432573d5f803e3d5ffd5b5090508160835f8c6001600160401b03166001600160401b031681526020019081526020015f205f6101000a81548163ffffffff021916908363ffffffff1602179055508160825f836001600160a01b03166001600160a01b031681526020019081526020015f205f6101000a81548163ffffffff021916908363ffffffff1602179055505f60885f8463ffffffff1663ffffffff1681526020019081526020015f20905081815f015f6101000a8154816001600160a01b0302191690836001600160a01b031602179055508360010160149054906101000a90046001600160401b03168160010160146101000a8154816001600160401b0302191690836001600160401b03160217905550836001015f9054906101000a90046001600160a01b0316816001015f6101000a8154816001600160a01b0302191690836001600160a01b031602179055508a815f0160146101000a8154816001600160401b0302191690836001600160401b031602179055508360020154816002015f806001600160401b031681526020019081526020015f20819055508b63ffffffff168160070160086101000a8154816001600160401b0302191690836001600160401b0316021790555083600101601c9054906101000a900460ff168160070160106101000a81548160ff021916908360ff1602179055508263ffffffff167f194c983456df6701c6a50830b90fe80e72b823411d0d524970c9590dc277a6418d848e8c60405162001696949392919063ffffffff9490941684526001600160a01b0392831660208501526001600160401b0391909116604084015216606082015260800190565b60405180910390a2604051633892b81160e11b81526001600160a01b03831690637125702290620016d6908d908d9088908e908e908e9060040162005130565b5f604051808303815f87803b158015620016ee575f80fd5b505af115801562001701573d5f803e3d5ffd5b50505050505050505050505050505050565b7f8cf807f6970720f8e2c208c7c5037595982c7bd9ed93c380d09df743d0dcc3fb6200173f8162003168565b670de0b6b3a7640000821180620017565750600182105b156200177557604051630c0bbd2760e01b815260040160405180910390fd5b608a8290556040518281527f13b1c630ad78354572e9ad473455d51831407e164b79dda20732f5acac5033829060200162001181565b5f54600390610100900460ff16158015620017cc57505f5460ff8083169116105b620018445760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201527f647920696e697469616c697a656400000000000000000000000000000000000060648201526084015b60405180910390fd5b5f805461ffff191660ff83161761010017905560015b60805463ffffffff16811162001a245763ffffffff81165f9081526081602090815260408083206088909252909120815481546001600160401b03600160a01b92839004811683027fffffffff0000000000000000ffffffffffffffffffffffffffffffffffffffff90921691909117835560018085018054918501805473ffffffffffffffffffffffffffffffffffffffff1981166001600160a01b039094169384178255915485900484169094026001600160e01b03199091169091171790915560058084015490830155600780840180549184018054600160401b938490048516840267ffffffffffffffff60401b19821681178355925460ff600160801b91829004160270ff000000000000000000000000000000001990931670ffffffffffffffffff000000000000000019909116179190911790556006840154908104821691168114620019ac575f80fd5b6001600160401b0381165f81815260028086016020908152604080842054848052928701825280842092909255928252600380870184528183205483805290860190935290205560705462001a07906305f5e10090620051a6565b608a555082915062001a1b905081620051bc565b9150506200185a565b505f805461ff001916905560405160ff821681527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a150565b7fa5c5790f581d443ed43873ab47cfb8c5d66a6db268e58b5971bb33fc66e07db162001a948162003168565b62015180826001600160401b0316111562001ac257604051633812d75d60e21b815260040160405180910390fd5b6085805467ffffffffffffffff60401b1916600160401b6001600160401b038516908102919091179091556040519081527fe84eacb10b29a9cd283d1c48f59cd87da8c2f99c554576228566d69aeba740cd9060200162001181565b606f5460ff161562001b4357604051630bc011ff60e21b815260040160405180910390fd5b63ffffffff88165f90815260886020526040902062001b69818989898989898962003174565b6001600160401b0387165f9081526004820160209081526040918290206002015482519081529081018590527f1f44c21118c4603cfb4e1b621dbcfa2b73efcececee2b99b620b2953d33a7010910160405180910390a162001bca620035ea565b505050505050505050565b6040516370a0823160e01b81523060048201525f9081906001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016906370a0823190602401602060405180830381865afa15801562001c3c573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019062001c629190620051d7565b6089549091505f9062001c88906001600160801b03600160801b820481169116620051ef565b6001600160801b03169050805f0362001ca3575f9250505090565b62001caf8183620051a6565b9250505090565b7fa5c5790f581d443ed43873ab47cfb8c5d66a6db268e58b5971bb33fc66e07db162001ce28162003168565b606f5460ff1662001d1d576085546001600160401b039081169083161062001d1d5760405163048a05a960e41b815260040160405180910390fd5b6085805467ffffffffffffffff19166001600160401b0384169081179091556040519081527fc4121f4e22c69632ebb7cf1f462be0511dc034f999b52013eddfb24aab765c759060200162001181565b7fa5c5790f581d443ed43873ab47cfb8c5d66a6db268e58b5971bb33fc66e07db162001d998162003168565b6103e88261ffff16108062001db357506103ff8261ffff16115b1562001dd2576040516344ceee7360e01b815260040160405180910390fd5b6085805471ffff000000000000000000000000000000001916600160801b61ffff8516908102919091179091556040519081527f5c8a9e64670a8ec12a8004aa047cbb455403a6c4f2d2ad4e52328400dc8142659060200162001181565b6080545f9063ffffffff1680820362001e4a57505f919050565b5f816001600160401b0381111562001e665762001e6662004bcd565b60405190808252806020026020018201604052801562001e90578160200160208202803683370190505b5090505f5b8281101562001ef45760885f62001eae83600162005212565b63ffffffff1663ffffffff1681526020019081526020015f206005015482828151811062001ee05762001ee062005228565b602090810291909101015260010162001e95565b505f60205b8360011462002140575f62001f106002866200523c565b62001f1d600287620051a6565b62001f29919062005212565b90505f816001600160401b0381111562001f475762001f4762004bcd565b60405190808252806020026020018201604052801562001f71578160200160208202803683370190505b5090505f5b82811015620020ec5762001f8c60018462005252565b8114801562001fa7575062001fa36002886200523c565b6001145b156200202f578562001fbb82600262005040565b8151811062001fce5762001fce62005228565b60200260200101518560405160200162001ff2929190918252602082015260400190565b604051602081830303815290604052805190602001208282815181106200201d576200201d62005228565b602002602001018181525050620020e3565b856200203d82600262005040565b8151811062002050576200205062005228565b60200260200101518682600262002068919062005040565b6200207590600162005212565b8151811062002088576200208862005228565b6020026020010151604051602001620020ab929190918252602082015260400190565b60405160208183030381529060405280519060200120828281518110620020d657620020d662005228565b6020026020010181815250505b60010162001f76565b50809450819550838460405160200162002110929190918252602082015260400190565b6040516020818303038152906040528051906020012093508280620021359062005268565b935050505062001ef9565b5f835f8151811062002156576200215662005228565b602002602001015190505f5b82811015620021cb57604080516020810184905290810185905260600160408051601f198184030181528282528051602091820120908301879052908201869052925060600160408051601f198184030181529190528051602090910120935060010162002162565b5095945050505050565b7f66156603fe29d13f97c6f3e3dff4ef71919f9aa61c555be0182d954e94221aac620022018162003168565b6200220e8484846200376a565b50505050565b5f82815260346020526040902060010154620022308162003168565b62000f7e8383620036e8565b7f62ba6ba2ffed8cfe316b583325ea41ac6e7ba9e5864d2bc6fabba7ac26d2f0f4620022688162003168565b6087805467ffffffffffffffff1916426001600160401b031617905560408051636de0b4bb60e11b815290517f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03169163dbc16976916004808301925f92919082900301818387803b158015620022e4575f80fd5b505af1158015620022f7573d5f803e3d5ffd5b505050506200230562003a57565b50565b7f084e94f375e9d647f87f5b2ceffba1e062c70f6009fdbcf80291e803b5c9edd4620023348162003168565b620023428585858562003aaf565b5f5b8481101562002509575f86868381811062002363576200236362005228565b905060c002018036038101906200237b919062005280565b805163ffffffff165f908152608860209081526040808320606085015160068201805467ffffffffffffffff60401b1916600160401b6001600160401b0393841690810291909117825560a0880151908752600284019095529290942092909255608084015160058301555492935091600160801b900416156200240b576006810180546001600160801b031690555b815163ffffffff165f908152608860205260409081902054606084015160a0850151925163444e7ebd60e11b81526001600160401b03909116600482015260248101929092523360448301526001600160a01b03169063889cfd7a906064015f604051808303815f87803b15801562002482575f80fd5b505af115801562002495573d5f803e3d5ffd5b5050835160608086015160a08701516080880151604080516001600160401b03909416845260208401929092529082015233945063ffffffff90921692507fba7fad50a32b4eb9847ff1f56dd7528178eae3cd0b008c7a798e0d5375de88da910160405180910390a3505060010162002344565b507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166333d6247d6200254362001e30565b6040518263ffffffff1660e01b81526004016200256291815260200190565b5f604051808303815f87803b1580156200257a575f80fd5b505af115801562001bca573d5f803e3d5ffd5b336001600160a01b0316826001600160a01b031663f851a4406040518163ffffffff1660e01b81526004016020604051808303815f875af1158015620025d5573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190620025fb91906200531d565b6001600160a01b031614620026235760405163696072e960e01b815260040160405180910390fd5b6001600160a01b0382165f9081526082602090815260408083205463ffffffff1683526088909152902060068101546001600160401b03808216600160401b9092041614620026855760405163664316a560e11b815260040160405180910390fd5b600781015463ffffffff8316600160401b9091046001600160401b031610620026c157604051634f61d51960e01b815260040160405180910390fd5b604080515f81526020810190915262000f7e90849084906200376a565b7f3dfe277d2a2c04b75fb2eb3743fa00005ae3678a20c299e65fdf4df76517f68e6200270a8162003168565b6001600160401b0384165f9081526083602052604090205463ffffffff161562002747576040516337c8fe0960e11b815260040160405180910390fd5b6001600160a01b0387165f9081526082602052604090205463ffffffff16156200278457604051630d409b9360e41b815260040160405180910390fd5b5f62002794888888888762003df2565b5f8080526002909101602052604090209390935550505050505050565b7f3dfe277d2a2c04b75fb2eb3743fa00005ae3678a20c299e65fdf4df76517f68e620027dd8162003168565b608780547fffffffff0000000000000000000000000000000000000000ffffffffffffffff16600160401b6001600160a01b038516908102919091179091556040519081527f53ab89ca5f00e99098ada1782f593e3f76b5489459ece48450e554c2928daa5e9060200162001181565b606f5460ff16156200287257604051630bc011ff60e21b815260040160405180910390fd5b620028808484848462003aaf565b5f5b8381101562002bab575f858583818110620028a157620028a162005228565b905060c00201803603810190620028b9919062005280565b805163ffffffff165f90815260886020908152604080832060845460608601516001600160401b039081168652600383019094529190932060010154939450919242926200291392600160c01b909104811691166200505a565b6001600160401b031611156200293c57604051638a0704d360e01b815260040160405180910390fd5b6200294c81836060015162004012565b6085546001600160401b03165f03620029cc57606082015160068201805467ffffffffffffffff19166001600160401b03928316908117825560a08501515f918252600285016020526040909120556080840151600584015554600160801b90041615620029c6576006810180546001600160801b031690555b62002aad565b620029d78162004260565b600681018054600160801b90046001600160401b0316906010620029fb836200533b565b82546001600160401b039182166101009390930a92830292820219169190911790915560408051608080820183524284168252606080880151851660208085019182529289015184860190815260a08a01519285019283526006890154600160801b900487165f90815260048a01909452949092209251835492518616600160401b026fffffffffffffffffffffffffffffffff19909316951694909417178155905160018201559051600290910155505b815163ffffffff165f908152608860205260409081902054606084015160a0850151925163444e7ebd60e11b81526001600160401b03909116600482015260248101929092523360448301526001600160a01b03169063889cfd7a906064015f604051808303815f87803b15801562002b24575f80fd5b505af115801562002b37573d5f803e3d5ffd5b5050835160608086015160a08701516080880151604080516001600160401b03909416845260208401929092529082015233945063ffffffff90921692507f716b8543c1c3c328a13d34cd51e064a780149a2d06455e44097de219b150e8b4910160405180910390a3505060010162002882565b507f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166333d6247d62002be562001e30565b6040518263ffffffff1660e01b815260040162002c0491815260200190565b5f604051808303815f87803b15801562002c1c575f80fd5b505af115801562002c2f573d5f803e3d5ffd5b5050505050505050565b7fac75d24dbb35ea80e25fab167da4dea46c1915260426570db84f184891f5f59062002c658162003168565b607e80545f9190829062002c7f9063ffffffff1662005084565b91906101000a81548163ffffffff021916908363ffffffff160217905590506040518060c00160405280896001600160a01b03168152602001886001600160a01b03168152602001876001600160401b031681526020018660ff1681526020015f1515815260200185815250607f5f8363ffffffff1663ffffffff1681526020019081526020015f205f820151815f015f6101000a8154816001600160a01b0302191690836001600160a01b031602179055506020820151816001015f6101000a8154816001600160a01b0302191690836001600160a01b0316021790555060408201518160010160146101000a8154816001600160401b0302191690836001600160401b03160217905550606082015181600101601c6101000a81548160ff021916908360ff160217905550608082015181600101601d6101000a81548160ff02191690831515021790555060a082015181600201559050508063ffffffff167fa2970448b3bd66ba7e524e7b2a5b9cf94fa29e32488fb942afdfe70dd4b77b5289898989898960405162002e1b9695949392919062005359565b60405180910390a25050505050505050565b606f545f9060ff161562002e5457604051630bc011ff60e21b815260040160405180910390fd5b335f9081526082602052604081205463ffffffff169081900362002e8b576040516371653c1560e01b815260040160405180910390fd5b836001600160401b03165f0362002eb55760405163158aa4dd60e21b815260040160405180910390fd5b63ffffffff81165f908152608860205260408120608980549192889262002ee79084906001600160801b0316620053b0565b82546001600160801b039182166101009390930a92830291909202199091161790555060068101546001600160401b03165f62002f268260016200505a565b6001600160401b0383165f9081526003850160205260408120600101549192509062002f64908a90600160801b90046001600160801b0316620053b0565b6001600160401b038085165f9081526003870160205260408120600101549293509162002f9b918b91600160401b9004166200505a565b6006860180546001600160401b0380871667ffffffffffffffff199092168217909255604080516080810182528c815242841660208083019182528587168385019081526001600160801b03808b16606086019081525f97885260038f01909352949095209251835590516001929092018054945191518416600160801b02918616600160401b026fffffffffffffffffffffffffffffffff19909516929095169190911792909217161790559050620030558562004260565b604080516001600160801b038c1681526001600160401b038b16602082015263ffffffff8816917fd3104eaeb2b51fc52b7d354a19bf146d10ed8d047b43764be8f78cbb3ffd8be4910160405180910390a2509098975050505050505050565b6085546001600160401b038281165f90815260048501602052604081205490924292620030e79291811691166200505a565b6001600160401b031611159392505050565b60068101545f90600160801b90046001600160401b0316156200314b575060068101546001600160401b03600160801b90910481165f9081526004909201602052604090912054600160401b90041690565b5060060154600160401b90046001600160401b031690565b919050565b62002305813362004329565b5f620031828989886200436c565b60068a01549091506001600160401b03600160801b90910481169088161180620031be5750876001600160401b0316876001600160401b031611155b80620031e2575060068901546001600160401b03600160c01b909104811690881611155b15620032015760405163bfa7079f60e01b815260040160405180910390fd5b6001600160401b038781165f90815260048b016020526040902054600160401b9004811690861614620032475760405163b7d5b4a360e01b815260040160405180910390fd5b60605f806200325a610100601462005212565b90506040519250806040840101604052808352602083019150620032848c8a8a8a888b8862004484565b3360601b815291505f7f30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f0000001600285604051620032c09190620053d3565b602060405180830381855afa158015620032dc573d5f803e3d5ffd5b5050506040513d601f19601f82011682018060405250810190620033019190620051d7565b6200330d91906200523c565b60018e0154604080516020810182528381529051634890ed4560e11b81529293506001600160a01b0390911691639121da8a9162003351918a9190600401620053f0565b602060405180830381865afa1580156200336d573d5f803e3d5ffd5b505050506040513d601f19601f820116820180604052508101906200339391906200542c565b620033b1576040516309bde33960e01b815260040160405180910390fd5b6001600160401b038b165f90815260048e016020526040902060020154879003620033ef5760405163a47276bd60e01b815260040160405180910390fd5b50505050505050505050505050565b60068201546001600160401b03600160c01b90910481169082161115806200343d575060068201546001600160401b03600160801b9091048116908216115b156200345c5760405163d086b70b60e01b815260040160405180910390fd5b6001600160401b038181165f8181526004850160209081526040808320805460068901805467ffffffffffffffff60401b1916600160401b92839004909816918202979097178755600280830154828752908a0190945291909320919091556001820154600587015583546001600160c01b0316600160c01b909302929092179092557f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03166333d6247d6200351862001e30565b6040518263ffffffff1660e01b81526004016200353791815260200190565b5f604051808303815f87803b1580156200354f575f80fd5b505af115801562003562573d5f803e3d5ffd5b505085546001600160a01b03165f90815260826020908152604091829020546002870154600188015484516001600160401b03898116825294810192909252818501529188166060830152915163ffffffff90921693507f581910eb7a27738945c2f00a91f2284b2d6de9d4e472b12f901c2b0df045e21b925081900360800190a250505050565b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b0316632072f6c56040518163ffffffff1660e01b81526004015f604051808303815f87803b15801562003643575f80fd5b505af115801562003656573d5f803e3d5ffd5b505050506200105d620045d9565b5f8281526034602090815260408083206001600160a01b038516845290915290205460ff16620010bd575f8281526034602090815260408083206001600160a01b0385168085529252808320805460ff1916600117905551339285917f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d9190a45050565b5f8281526034602090815260408083206001600160a01b038516845290915290205460ff1615620010bd575f8281526034602090815260408083206001600160a01b0385168085529252808320805460ff1916905551339285917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a45050565b63ffffffff82161580620037895750607e5463ffffffff908116908316115b15620037a857604051637512e5cb60e01b815260040160405180910390fd5b6001600160a01b0383165f9081526082602052604081205463ffffffff1690819003620037e8576040516374a086a360e01b815260040160405180910390fd5b63ffffffff8181165f908152608860205260409020600781015490918516600160401b9091046001600160401b0316036200383657604051634f61d51960e01b815260040160405180910390fd5b63ffffffff84165f908152607f60205260409020600180820154600160e81b900460ff16151590036200387c57604051633b8d3d9960e01b815260040160405180910390fd5b60018101546007830154600160801b900460ff908116600160e01b9092041614620038ba57604051635aa0d5f160e11b815260040160405180910390fd5b6001818101805491840180546001600160a01b0390931673ffffffffffffffffffffffffffffffffffffffff1984168117825591546001600160e01b0319909316909117600160a01b928390046001600160401b0390811690930217905560078301805467ffffffffffffffff60401b191663ffffffff8816600160401b021790556006830154600160c01b81048216600160801b909104909116146200397457604051639d59507b60e01b815260040160405180910390fd5b5f620039808462000cff565b60078401805467ffffffffffffffff19166001600160401b038316179055825460405163278f794360e11b81529192506001600160a01b0389811692634f1ef28692620039d492169089906004016200544d565b5f604051808303815f87803b158015620039ec575f80fd5b505af1158015620039ff573d5f803e3d5ffd5b50506040805163ffffffff8a811682526001600160401b0386166020830152881693507ff585e04c05d396901170247783d3e5f0ee9c1df23072985b50af089f5e48b19d92500160405180910390a250505050505050565b606f5460ff1662003a7b57604051635386698160e01b815260040160405180910390fd5b606f805460ff191690556040517f1e5e34eea33501aecf2ebec9fe0e884a40804275ea7fe10b2ba084c8374308b3905f90a1565b60605f8062003ac18661010062005040565b62003ace90601462005212565b905060405192508060408401016040528083526020830191505f805f5b8881101562003baa575f8a8a8381811062003b0a5762003b0a62005228565b62003b2292602060c090920201908101915062004a93565b90508363ffffffff168163ffffffff161162003b51576040516328fe7b1560e11b815260040160405180910390fd5b8093505f62003b8d8c8c8581811062003b6e5762003b6e62005228565b905060c0020180360381019062003b86919062005280565b8862004635565b9750905062003b9d8185620053b0565b9350505060010162003aeb565b503360601b84525f7f30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f000000160028760405162003be59190620053d3565b602060405180830381855afa15801562003c01573d5f803e3d5ffd5b5050506040513d601f19601f8201168201806040525081019062003c269190620051d7565b62003c3291906200523c565b90505f60018a900362003c995760885f8c8c5f81811062003c575762003c5762005228565b62003c6f92602060c090920201908101915062004a93565b63ffffffff16815260208101919091526040015f20600101546001600160a01b0316905062003cae565b50608754600160401b90046001600160a01b03165b604080516020810182528381529051634890ed4560e11b81526001600160a01b03831691639121da8a9162003ce8918c91600401620053f0565b602060405180830381865afa15801562003d04573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019062003d2a91906200542c565b62003d48576040516309bde33960e01b815260040160405180910390fd5b62003d9f89846001600160801b031662003d6162001bd5565b62003d6d919062005040565b6001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016919062004732565b62003dab8380620053b0565b5050608480547fffffffffffffffff0000000000000000ffffffffffffffffffffffffffffffff16600160801b426001600160401b03160217905550505050505050505050565b608080545f918291829062003e0d9063ffffffff1662005084565b91906101000a81548163ffffffff021916908363ffffffff160217905590508060835f866001600160401b03166001600160401b031681526020019081526020015f205f6101000a81548163ffffffff021916908363ffffffff1602179055508060825f896001600160a01b03166001600160a01b031681526020019081526020015f205f6101000a81548163ffffffff021916908363ffffffff16021790555060885f8263ffffffff1663ffffffff1681526020019081526020015f20915086825f015f6101000a8154816001600160a01b0302191690836001600160a01b03160217905550848260010160146101000a8154816001600160401b0302191690836001600160401b0316021790555085826001015f6101000a8154816001600160a01b0302191690836001600160a01b0316021790555083825f0160146101000a8154816001600160401b0302191690836001600160401b03160217905550828260070160106101000a81548160ff021916908360ff1602179055508063ffffffff167fadfc7d56f7e39b08b321534f14bfb135ad27698f7d2f5ad0edc2356ea9a3f850868987875f604051620040009594939291906001600160401b0395861681526001600160a01b03949094166020850152918416604084015260ff166060830152909116608082015260a00190565b60405180910390a25095945050505050565b5f6200401e83620030f9565b6001600160401b038082165f9081526003860160205260408082206001908101549387168352908220015492935084929091829162004076916001600160801b03600160801b918290048116929190910416620051ef565b6085546001600160801b039190911691505f90620040a590600160401b90046001600160401b03164262005252565b90505b846001600160401b0316846001600160401b03161462004132576001600160401b038085165f908152600389016020526040902060018101549091168210156200410157620040f960018662005470565b94506200412b565b60018101546200412290600160801b90046001600160801b03168462005252565b93505062004132565b50620040a8565b5f6200413f848462005252565b905080841015620041a3576305f5e10084820304600c811162004163578062004166565b600c5b9050806103e80a81608560109054906101000a900461ffff1661ffff160a608a54028162004198576200419862005192565b04608a555062004220565b6305f5e10081850304600c8111620041bc5780620041bf565b600c5b90505f816103e80a82608560109054906101000a900461ffff1661ffff160a670de0b6b3a76400000281620041f857620041f862005192565b04905080608a54670de0b6b3a7640000028162004219576200421962005192565b04608a5550505b670de0b6b3a7640000608a5411156200424557670de0b6b3a7640000608a5562002c2f565b6001608a54101562002c2f576001608a555050505050505050565b60068101546001600160401b03600160c01b82048116600160801b909204161115620023055760068101545f90620042aa90600160c01b90046001600160401b031660016200505a565b9050620042b88282620030b5565b15620010bd5760068201545f90600290620042e5908490600160801b90046001600160401b031662005470565b620042f1919062005493565b620042fd90836200505a565b90506200430b8382620030b5565b156200431d5762000f7e8382620033fe565b62000f7e8383620033fe565b5f8281526034602090815260408083206001600160a01b038516845290915290205460ff16620010bd57604051637615be1f60e11b815260040160405180910390fd5b60078301545f906001600160401b039081169083161015620043a15760405163f5f2eb1360e01b815260040160405180910390fd5b5f6001600160401b03841615620044425760068501546001600160401b03600160801b90910481169085161115620043ec5760405163bb14c20560e01b815260040160405180910390fd5b506001600160401b038084165f9081526004860160205260409020600281015481549092858116600160401b90920416146200443b5760405163686446b160e01b815260040160405180910390fd5b506200447c565b506001600160401b0382165f908152600285016020526040902054806200447c576040516324cbdcc360e11b815260040160405180910390fd5b949350505050565b6001600160401b038087165f81815260038a01602052604080822054938916825281205490929115801590620044b8575081155b15620044d75760405163340c614f60e11b815260040160405180910390fd5b80620044f6576040516366385b5160e01b815260040160405180910390fd5b62004501856200479b565b6200451f576040516305dae44f60e21b815260040160405180910390fd5b6001600160401b039889165f90815260038b01602090815260408083206001908101549b909c1683528083208c01549887528682018390528601939093527fffffffffffffffff000000000000000000000000000000000000000000000000600160401b998a900460c090811b821660608801528c54851b60688801529a909b015490921b6070850152607884019490945260988301525060b881019190915292900490921b90921660d883015260e08201526101000190565b606f5460ff1615620045fe57604051630bc011ff60e21b815260040160405180910390fd5b606f805460ff191660011790556040517f2261efe5aef6fedc1fd1550b25facc9181745623049c7901287030b9ad1a5497905f90a1565b815163ffffffff165f90815260886020908152604080832091850151908501518392918391620046679184916200436c565b90505f6200467583620030f9565b9050806001600160401b031687606001516001600160401b031611620046ae576040516321798fc960e11b815260040160405180910390fd5b5f620046d08489604001518a606001518b60800151878d60a001518d62004484565b6001600160401b038084165f90815260038701602052604080822060019081015460608e015190941683529120015491925062004725916001600160801b03600160801b9283900481169290910416620051ef565b9890975095505050505050565b604080516001600160a01b038416602482015260448082018490528251808303909101815260649091019091526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff1663a9059cbb60e01b17905262000f7e9084906200481f565b5f67ffffffff000000016001600160401b038316108015620047d1575067ffffffff00000001604083901c6001600160401b0316105b8015620047f2575067ffffffff00000001608083901c6001600160401b0316105b80156200480a575067ffffffff0000000160c083901c105b156200481857506001919050565b505f919050565b5f62004875826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c6564815250856001600160a01b0316620048f79092919063ffffffff16565b80519091501562000f7e57808060200190518101906200489691906200542c565b62000f7e5760405162461bcd60e51b815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e6044820152691bdd081cdd58d8d9595960b21b60648201526084016200183b565b60606200447c84845f85855f80866001600160a01b031685876040516200491f9190620053d3565b5f6040518083038185875af1925050503d805f81146200495b576040519150601f19603f3d011682016040523d82523d5f602084013e62004960565b606091505b509150915062004973878383876200497e565b979650505050505050565b60608315620049f15782515f03620049e9576001600160a01b0385163b620049e95760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e747261637400000060448201526064016200183b565b50816200447c565b6200447c838381511562004a085781518083602001fd5b8060405162461bcd60e51b81526004016200183b9190620054bb565b61091d80620054d083390190565b803563ffffffff8116811462003163575f80fd5b80356001600160401b038116811462003163575f80fd5b5f806040838503121562004a6f575f80fd5b62004a7a8362004a32565b915062004a8a6020840162004a46565b90509250929050565b5f6020828403121562004aa4575f80fd5b62000cf68262004a32565b80610300810183101562000cf9575f80fd5b5f805f805f805f806103e0898b03121562004ada575f80fd5b62004ae58962004a32565b975062004af560208a0162004a46565b965062004b0560408a0162004a46565b955062004b1560608a0162004a46565b945062004b2560808a0162004a46565b935060a0890135925060c0890135915062004b448a60e08b0162004aaf565b90509295985092959890939650565b5f6020828403121562004b64575f80fd5b5035919050565b6001600160a01b038116811462002305575f80fd5b5f806040838503121562004b92575f80fd5b82359150602083013562004ba68162004b6b565b809150509250929050565b5f6020828403121562004bc2575f80fd5b62000cf68262004a46565b634e487b7160e01b5f52604160045260245ffd5b5f6001600160401b038084111562004bfd5762004bfd62004bcd565b604051601f8501601f19908116603f0116810190828211818310171562004c285762004c2862004bcd565b8160405280935085815286868601111562004c41575f80fd5b858560208301375f602087830101525050509392505050565b5f82601f83011262004c6a575f80fd5b62000cf68383356020850162004be1565b5f805f805f805f60e0888a03121562004c92575f80fd5b62004c9d8862004a32565b965062004cad6020890162004a46565b9550604088013562004cbf8162004b6b565b9450606088013562004cd18162004b6b565b9350608088013562004ce38162004b6b565b925060a08801356001600160401b038082111562004cff575f80fd5b62004d0d8b838c0162004c5a565b935060c08a013591508082111562004d23575f80fd5b5062004d328a828b0162004c5a565b91505092959891949750929550565b5f6020828403121562004d52575f80fd5b813561ffff8116811462004d64575f80fd5b9392505050565b5f805f6060848603121562004d7e575f80fd5b833562004d8b8162004b6b565b925062004d9b6020850162004a32565b915060408401356001600160401b0381111562004db6575f80fd5b8401601f8101861362004dc7575f80fd5b62004dd88682356020840162004be1565b9150509250925092565b5f6020828403121562004df3575f80fd5b813562004d648162004b6b565b5f805f80610340858703121562004e15575f80fd5b84356001600160401b038082111562004e2c575f80fd5b818701915087601f83011262004e40575f80fd5b81358181111562004e4f575f80fd5b88602060c08302850101111562004e64575f80fd5b6020928301965094505085013562004e7c8162004b6b565b915062004e8d866040870162004aaf565b905092959194509250565b5f806040838503121562004eaa575f80fd5b823562004eb78162004b6b565b915062004a8a6020840162004a32565b803560ff8116811462003163575f80fd5b5f805f805f8060c0878903121562004eee575f80fd5b863562004efb8162004b6b565b9550602087013562004f0d8162004b6b565b945062004f1d6040880162004a46565b935062004f2d6060880162004a46565b92506080870135915062004f4460a0880162004ec7565b90509295509295509295565b5f805f805f8060c0878903121562004f66575f80fd5b863562004f738162004b6b565b9550602087013562004f858162004b6b565b945062004f956040880162004a46565b935062004fa56060880162004ec7565b92506080870135915060a08701356001600160401b0381111562004fc7575f80fd5b62004fd589828a0162004c5a565b9150509295509295509295565b5f805f6060848603121562004ff5575f80fd5b83356001600160801b03811681146200500c575f80fd5b92506200501c6020850162004a46565b9150604084013590509250925092565b634e487b7160e01b5f52601160045260245ffd5b808202811582820484141762000cf95762000cf96200502c565b6001600160401b038181168382160190808211156200507d576200507d6200502c565b5092915050565b5f63ffffffff8083168181036200509f576200509f6200502c565b6001019392505050565b5f5b83811015620050c5578181015183820152602001620050ab565b50505f910152565b5f8151808452620050e6816020860160208601620050a9565b601f01601f19169290920160200192915050565b5f6001600160a01b03808616835280851660208401525060606040830152620051276060830184620050cd565b95945050505050565b5f6001600160a01b038089168352808816602084015263ffffffff8716604084015280861660608401525060c060808301526200517160c0830185620050cd565b82810360a0840152620051858185620050cd565b9998505050505050505050565b634e487b7160e01b5f52601260045260245ffd5b5f82620051b757620051b762005192565b500490565b5f60018201620051d057620051d06200502c565b5060010190565b5f60208284031215620051e8575f80fd5b5051919050565b6001600160801b038281168282160390808211156200507d576200507d6200502c565b8082018082111562000cf95762000cf96200502c565b634e487b7160e01b5f52603260045260245ffd5b5f826200524d576200524d62005192565b500690565b8181038181111562000cf95762000cf96200502c565b5f816200527957620052796200502c565b505f190190565b5f60c0828403121562005291575f80fd5b60405160c081018181106001600160401b0382111715620052b657620052b662004bcd565b604052620052c48362004a32565b8152620052d46020840162004a46565b6020820152620052e76040840162004a46565b6040820152620052fa6060840162004a46565b60608201526080830135608082015260a083013560a08201528091505092915050565b5f602082840312156200532e575f80fd5b815162004d648162004b6b565b5f6001600160401b038083168181036200509f576200509f6200502c565b5f6001600160a01b0380891683528088166020840152506001600160401b038616604083015260ff8516606083015283608083015260c060a0830152620053a460c0830184620050cd565b98975050505050505050565b6001600160801b038181168382160190808211156200507d576200507d6200502c565b5f8251620053e6818460208701620050a9565b9190910192915050565b6103208101610300808584378201835f5b60018110156200542257815183526020928301929091019060010162005401565b5050509392505050565b5f602082840312156200543d575f80fd5b8151801515811462004d64575f80fd5b6001600160a01b0383168152604060208201525f6200447c6040830184620050cd565b6001600160401b038281168282160390808211156200507d576200507d6200502c565b5f6001600160401b0380841680620054af57620054af62005192565b92169190910492915050565b602081525f62000cf66020830184620050cd56fe60a06040526040516200091d3803806200091d833981016040819052620000269162000375565b828162000034828262000060565b50506001600160a01b038216608052620000576200005160805190565b620000c5565b5050506200046c565b6200006b8262000136565b6040516001600160a01b038316907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b905f90a2805115620000b757620000b28282620001b5565b505050565b620000c16200022e565b5050565b7f7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f620001065f80516020620008fd833981519152546001600160a01b031690565b604080516001600160a01b03928316815291841660208301520160405180910390a1620001338162000250565b50565b806001600160a01b03163b5f036200017157604051634c9c8ce360e01b81526001600160a01b03821660048201526024015b60405180910390fd5b807f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc5b80546001600160a01b0319166001600160a01b039290921691909117905550565b60605f80846001600160a01b031684604051620001d391906200044f565b5f60405180830381855af49150503d805f81146200020d576040519150601f19603f3d011682016040523d82523d5f602084013e62000212565b606091505b5090925090506200022585838362000291565b95945050505050565b34156200024e5760405163b398979f60e01b815260040160405180910390fd5b565b6001600160a01b0381166200027b57604051633173bdd160e11b81525f600482015260240162000168565b805f80516020620008fd83398151915262000194565b606082620002aa57620002a482620002f7565b620002f0565b8151158015620002c257506001600160a01b0384163b155b15620002ed57604051639996b31560e01b81526001600160a01b038516600482015260240162000168565b50805b9392505050565b805115620003085780518082602001fd5b604051630a12f52160e11b815260040160405180910390fd5b80516001600160a01b038116811462000338575f80fd5b919050565b634e487b7160e01b5f52604160045260245ffd5b5f5b838110156200036d57818101518382015260200162000353565b50505f910152565b5f805f6060848603121562000388575f80fd5b620003938462000321565b9250620003a36020850162000321565b60408501519092506001600160401b0380821115620003c0575f80fd5b818601915086601f830112620003d4575f80fd5b815181811115620003e957620003e96200033d565b604051601f8201601f19908116603f011681019083821181831017156200041457620004146200033d565b816040528281528960208487010111156200042d575f80fd5b6200044083602083016020880162000351565b80955050505050509250925092565b5f82516200046281846020870162000351565b9190910192915050565b608051610479620004845f395f601001526104795ff3fe608060405261000c61000e565b005b7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03163303610081575f357fffffffff000000000000000000000000000000000000000000000000000000001663278f794360e11b1461007957610077610085565b565b610077610095565b6100775b6100776100906100c3565b6100fa565b5f806100a43660048184610313565b8101906100b1919061034e565b915091506100bf8282610118565b5050565b5f6100f57f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc546001600160a01b031690565b905090565b365f80375f80365f845af43d5f803e808015610114573d5ff35b3d5ffd5b61012182610172565b6040516001600160a01b038316907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b905f90a280511561016a5761016582826101fa565b505050565b6100bf61026c565b806001600160a01b03163b5f036101ac57604051634c9c8ce360e01b81526001600160a01b03821660048201526024015b60405180910390fd5b7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc805473ffffffffffffffffffffffffffffffffffffffff19166001600160a01b0392909216919091179055565b60605f80846001600160a01b0316846040516102169190610417565b5f60405180830381855af49150503d805f811461024e576040519150601f19603f3d011682016040523d82523d5f602084013e610253565b606091505b509150915061026385838361028b565b95945050505050565b34156100775760405163b398979f60e01b815260040160405180910390fd5b6060826102a05761029b826102ea565b6102e3565b81511580156102b757506001600160a01b0384163b155b156102e057604051639996b31560e01b81526001600160a01b03851660048201526024016101a3565b50805b9392505050565b8051156102fa5780518082602001fd5b604051630a12f52160e11b815260040160405180910390fd5b5f8085851115610321575f80fd5b8386111561032d575f80fd5b5050820193919092039150565b634e487b7160e01b5f52604160045260245ffd5b5f806040838503121561035f575f80fd5b82356001600160a01b0381168114610375575f80fd5b9150602083013567ffffffffffffffff80821115610391575f80fd5b818501915085601f8301126103a4575f80fd5b8135818111156103b6576103b661033a565b604051601f8201601f19908116603f011681019083821181831017156103de576103de61033a565b816040528281528860208487010111156103f6575f80fd5b826020860160208301375f6020848301015280955050505050509250929050565b5f82515f5b81811015610436576020818601810151858301520161041c565b505f92019182525091905056fea2646970667358221220cdb50aeb657f43ff038a16fe0b1c0e5f0d88a7122cf9eee7cfe0167fe21044db64736f6c63430008180033b53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103a2646970667358221220f04deec52f95bd17f45524f0387264f4a659e8c8f0679527fc3ccc9b491277c064736f6c63430008180033",
}

// FeijoapolygonrollupmanagerABI is the input ABI used to generate the binding from.
// Deprecated: Use FeijoapolygonrollupmanagerMetaData.ABI instead.
var FeijoapolygonrollupmanagerABI = FeijoapolygonrollupmanagerMetaData.ABI

// FeijoapolygonrollupmanagerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use FeijoapolygonrollupmanagerMetaData.Bin instead.
var FeijoapolygonrollupmanagerBin = FeijoapolygonrollupmanagerMetaData.Bin

// DeployFeijoapolygonrollupmanager deploys a new Ethereum contract, binding an instance of Feijoapolygonrollupmanager to it.
func DeployFeijoapolygonrollupmanager(auth *bind.TransactOpts, backend bind.ContractBackend, _globalExitRootManager common.Address, _pol common.Address, _bridgeAddress common.Address) (common.Address, *types.Transaction, *Feijoapolygonrollupmanager, error) {
	parsed, err := FeijoapolygonrollupmanagerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(FeijoapolygonrollupmanagerBin), backend, _globalExitRootManager, _pol, _bridgeAddress)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Feijoapolygonrollupmanager{FeijoapolygonrollupmanagerCaller: FeijoapolygonrollupmanagerCaller{contract: contract}, FeijoapolygonrollupmanagerTransactor: FeijoapolygonrollupmanagerTransactor{contract: contract}, FeijoapolygonrollupmanagerFilterer: FeijoapolygonrollupmanagerFilterer{contract: contract}}, nil
}

// Feijoapolygonrollupmanager is an auto generated Go binding around an Ethereum contract.
type Feijoapolygonrollupmanager struct {
	FeijoapolygonrollupmanagerCaller     // Read-only binding to the contract
	FeijoapolygonrollupmanagerTransactor // Write-only binding to the contract
	FeijoapolygonrollupmanagerFilterer   // Log filterer for contract events
}

// FeijoapolygonrollupmanagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type FeijoapolygonrollupmanagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FeijoapolygonrollupmanagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FeijoapolygonrollupmanagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FeijoapolygonrollupmanagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FeijoapolygonrollupmanagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FeijoapolygonrollupmanagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FeijoapolygonrollupmanagerSession struct {
	Contract     *Feijoapolygonrollupmanager // Generic contract binding to set the session for
	CallOpts     bind.CallOpts               // Call options to use throughout this session
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// FeijoapolygonrollupmanagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FeijoapolygonrollupmanagerCallerSession struct {
	Contract *FeijoapolygonrollupmanagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                     // Call options to use throughout this session
}

// FeijoapolygonrollupmanagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FeijoapolygonrollupmanagerTransactorSession struct {
	Contract     *FeijoapolygonrollupmanagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                     // Transaction auth options to use throughout this session
}

// FeijoapolygonrollupmanagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type FeijoapolygonrollupmanagerRaw struct {
	Contract *Feijoapolygonrollupmanager // Generic contract binding to access the raw methods on
}

// FeijoapolygonrollupmanagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FeijoapolygonrollupmanagerCallerRaw struct {
	Contract *FeijoapolygonrollupmanagerCaller // Generic read-only contract binding to access the raw methods on
}

// FeijoapolygonrollupmanagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FeijoapolygonrollupmanagerTransactorRaw struct {
	Contract *FeijoapolygonrollupmanagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFeijoapolygonrollupmanager creates a new instance of Feijoapolygonrollupmanager, bound to a specific deployed contract.
func NewFeijoapolygonrollupmanager(address common.Address, backend bind.ContractBackend) (*Feijoapolygonrollupmanager, error) {
	contract, err := bindFeijoapolygonrollupmanager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Feijoapolygonrollupmanager{FeijoapolygonrollupmanagerCaller: FeijoapolygonrollupmanagerCaller{contract: contract}, FeijoapolygonrollupmanagerTransactor: FeijoapolygonrollupmanagerTransactor{contract: contract}, FeijoapolygonrollupmanagerFilterer: FeijoapolygonrollupmanagerFilterer{contract: contract}}, nil
}

// NewFeijoapolygonrollupmanagerCaller creates a new read-only instance of Feijoapolygonrollupmanager, bound to a specific deployed contract.
func NewFeijoapolygonrollupmanagerCaller(address common.Address, caller bind.ContractCaller) (*FeijoapolygonrollupmanagerCaller, error) {
	contract, err := bindFeijoapolygonrollupmanager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FeijoapolygonrollupmanagerCaller{contract: contract}, nil
}

// NewFeijoapolygonrollupmanagerTransactor creates a new write-only instance of Feijoapolygonrollupmanager, bound to a specific deployed contract.
func NewFeijoapolygonrollupmanagerTransactor(address common.Address, transactor bind.ContractTransactor) (*FeijoapolygonrollupmanagerTransactor, error) {
	contract, err := bindFeijoapolygonrollupmanager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FeijoapolygonrollupmanagerTransactor{contract: contract}, nil
}

// NewFeijoapolygonrollupmanagerFilterer creates a new log filterer instance of Feijoapolygonrollupmanager, bound to a specific deployed contract.
func NewFeijoapolygonrollupmanagerFilterer(address common.Address, filterer bind.ContractFilterer) (*FeijoapolygonrollupmanagerFilterer, error) {
	contract, err := bindFeijoapolygonrollupmanager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FeijoapolygonrollupmanagerFilterer{contract: contract}, nil
}

// bindFeijoapolygonrollupmanager binds a generic wrapper to an already deployed contract.
func bindFeijoapolygonrollupmanager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FeijoapolygonrollupmanagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Feijoapolygonrollupmanager.Contract.FeijoapolygonrollupmanagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.Contract.FeijoapolygonrollupmanagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.Contract.FeijoapolygonrollupmanagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Feijoapolygonrollupmanager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Feijoapolygonrollupmanager.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _Feijoapolygonrollupmanager.Contract.DEFAULTADMINROLE(&_Feijoapolygonrollupmanager.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _Feijoapolygonrollupmanager.Contract.DEFAULTADMINROLE(&_Feijoapolygonrollupmanager.CallOpts)
}

// ZKGASLIMITBATCH is a free data retrieval call binding the contract method 0x838a2503.
//
// Solidity: function ZK_GAS_LIMIT_BATCH() view returns(uint64)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCaller) ZKGASLIMITBATCH(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Feijoapolygonrollupmanager.contract.Call(opts, &out, "ZK_GAS_LIMIT_BATCH")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// ZKGASLIMITBATCH is a free data retrieval call binding the contract method 0x838a2503.
//
// Solidity: function ZK_GAS_LIMIT_BATCH() view returns(uint64)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerSession) ZKGASLIMITBATCH() (uint64, error) {
	return _Feijoapolygonrollupmanager.Contract.ZKGASLIMITBATCH(&_Feijoapolygonrollupmanager.CallOpts)
}

// ZKGASLIMITBATCH is a free data retrieval call binding the contract method 0x838a2503.
//
// Solidity: function ZK_GAS_LIMIT_BATCH() view returns(uint64)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCallerSession) ZKGASLIMITBATCH() (uint64, error) {
	return _Feijoapolygonrollupmanager.Contract.ZKGASLIMITBATCH(&_Feijoapolygonrollupmanager.CallOpts)
}

// AggregateRollupVerifier is a free data retrieval call binding the contract method 0xba988cef.
//
// Solidity: function aggregateRollupVerifier() view returns(address)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCaller) AggregateRollupVerifier(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Feijoapolygonrollupmanager.contract.Call(opts, &out, "aggregateRollupVerifier")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AggregateRollupVerifier is a free data retrieval call binding the contract method 0xba988cef.
//
// Solidity: function aggregateRollupVerifier() view returns(address)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerSession) AggregateRollupVerifier() (common.Address, error) {
	return _Feijoapolygonrollupmanager.Contract.AggregateRollupVerifier(&_Feijoapolygonrollupmanager.CallOpts)
}

// AggregateRollupVerifier is a free data retrieval call binding the contract method 0xba988cef.
//
// Solidity: function aggregateRollupVerifier() view returns(address)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCallerSession) AggregateRollupVerifier() (common.Address, error) {
	return _Feijoapolygonrollupmanager.Contract.AggregateRollupVerifier(&_Feijoapolygonrollupmanager.CallOpts)
}

// BridgeAddress is a free data retrieval call binding the contract method 0xa3c573eb.
//
// Solidity: function bridgeAddress() view returns(address)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCaller) BridgeAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Feijoapolygonrollupmanager.contract.Call(opts, &out, "bridgeAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BridgeAddress is a free data retrieval call binding the contract method 0xa3c573eb.
//
// Solidity: function bridgeAddress() view returns(address)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerSession) BridgeAddress() (common.Address, error) {
	return _Feijoapolygonrollupmanager.Contract.BridgeAddress(&_Feijoapolygonrollupmanager.CallOpts)
}

// BridgeAddress is a free data retrieval call binding the contract method 0xa3c573eb.
//
// Solidity: function bridgeAddress() view returns(address)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCallerSession) BridgeAddress() (common.Address, error) {
	return _Feijoapolygonrollupmanager.Contract.BridgeAddress(&_Feijoapolygonrollupmanager.CallOpts)
}

// CalculateRewardPerZkGas is a free data retrieval call binding the contract method 0x90031d5c.
//
// Solidity: function calculateRewardPerZkGas() view returns(uint256)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCaller) CalculateRewardPerZkGas(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Feijoapolygonrollupmanager.contract.Call(opts, &out, "calculateRewardPerZkGas")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CalculateRewardPerZkGas is a free data retrieval call binding the contract method 0x90031d5c.
//
// Solidity: function calculateRewardPerZkGas() view returns(uint256)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerSession) CalculateRewardPerZkGas() (*big.Int, error) {
	return _Feijoapolygonrollupmanager.Contract.CalculateRewardPerZkGas(&_Feijoapolygonrollupmanager.CallOpts)
}

// CalculateRewardPerZkGas is a free data retrieval call binding the contract method 0x90031d5c.
//
// Solidity: function calculateRewardPerZkGas() view returns(uint256)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCallerSession) CalculateRewardPerZkGas() (*big.Int, error) {
	return _Feijoapolygonrollupmanager.Contract.CalculateRewardPerZkGas(&_Feijoapolygonrollupmanager.CallOpts)
}

// ChainIDToRollupID is a free data retrieval call binding the contract method 0x7fb6e76a.
//
// Solidity: function chainIDToRollupID(uint64 chainID) view returns(uint32 rollupID)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCaller) ChainIDToRollupID(opts *bind.CallOpts, chainID uint64) (uint32, error) {
	var out []interface{}
	err := _Feijoapolygonrollupmanager.contract.Call(opts, &out, "chainIDToRollupID", chainID)

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// ChainIDToRollupID is a free data retrieval call binding the contract method 0x7fb6e76a.
//
// Solidity: function chainIDToRollupID(uint64 chainID) view returns(uint32 rollupID)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerSession) ChainIDToRollupID(chainID uint64) (uint32, error) {
	return _Feijoapolygonrollupmanager.Contract.ChainIDToRollupID(&_Feijoapolygonrollupmanager.CallOpts, chainID)
}

// ChainIDToRollupID is a free data retrieval call binding the contract method 0x7fb6e76a.
//
// Solidity: function chainIDToRollupID(uint64 chainID) view returns(uint32 rollupID)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCallerSession) ChainIDToRollupID(chainID uint64) (uint32, error) {
	return _Feijoapolygonrollupmanager.Contract.ChainIDToRollupID(&_Feijoapolygonrollupmanager.CallOpts, chainID)
}

// GetForcedZkGasPrice is a free data retrieval call binding the contract method 0x02f3fa60.
//
// Solidity: function getForcedZkGasPrice() view returns(uint256)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCaller) GetForcedZkGasPrice(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Feijoapolygonrollupmanager.contract.Call(opts, &out, "getForcedZkGasPrice")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetForcedZkGasPrice is a free data retrieval call binding the contract method 0x02f3fa60.
//
// Solidity: function getForcedZkGasPrice() view returns(uint256)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerSession) GetForcedZkGasPrice() (*big.Int, error) {
	return _Feijoapolygonrollupmanager.Contract.GetForcedZkGasPrice(&_Feijoapolygonrollupmanager.CallOpts)
}

// GetForcedZkGasPrice is a free data retrieval call binding the contract method 0x02f3fa60.
//
// Solidity: function getForcedZkGasPrice() view returns(uint256)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCallerSession) GetForcedZkGasPrice() (*big.Int, error) {
	return _Feijoapolygonrollupmanager.Contract.GetForcedZkGasPrice(&_Feijoapolygonrollupmanager.CallOpts)
}

// GetLastVerifiedSequence is a free data retrieval call binding the contract method 0x0a7eef7a.
//
// Solidity: function getLastVerifiedSequence(uint32 rollupID) view returns(uint64)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCaller) GetLastVerifiedSequence(opts *bind.CallOpts, rollupID uint32) (uint64, error) {
	var out []interface{}
	err := _Feijoapolygonrollupmanager.contract.Call(opts, &out, "getLastVerifiedSequence", rollupID)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// GetLastVerifiedSequence is a free data retrieval call binding the contract method 0x0a7eef7a.
//
// Solidity: function getLastVerifiedSequence(uint32 rollupID) view returns(uint64)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerSession) GetLastVerifiedSequence(rollupID uint32) (uint64, error) {
	return _Feijoapolygonrollupmanager.Contract.GetLastVerifiedSequence(&_Feijoapolygonrollupmanager.CallOpts, rollupID)
}

// GetLastVerifiedSequence is a free data retrieval call binding the contract method 0x0a7eef7a.
//
// Solidity: function getLastVerifiedSequence(uint32 rollupID) view returns(uint64)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCallerSession) GetLastVerifiedSequence(rollupID uint32) (uint64, error) {
	return _Feijoapolygonrollupmanager.Contract.GetLastVerifiedSequence(&_Feijoapolygonrollupmanager.CallOpts, rollupID)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _Feijoapolygonrollupmanager.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _Feijoapolygonrollupmanager.Contract.GetRoleAdmin(&_Feijoapolygonrollupmanager.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _Feijoapolygonrollupmanager.Contract.GetRoleAdmin(&_Feijoapolygonrollupmanager.CallOpts, role)
}

// GetRollupExitRoot is a free data retrieval call binding the contract method 0xa2967d99.
//
// Solidity: function getRollupExitRoot() view returns(bytes32)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCaller) GetRollupExitRoot(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Feijoapolygonrollupmanager.contract.Call(opts, &out, "getRollupExitRoot")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRollupExitRoot is a free data retrieval call binding the contract method 0xa2967d99.
//
// Solidity: function getRollupExitRoot() view returns(bytes32)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerSession) GetRollupExitRoot() ([32]byte, error) {
	return _Feijoapolygonrollupmanager.Contract.GetRollupExitRoot(&_Feijoapolygonrollupmanager.CallOpts)
}

// GetRollupExitRoot is a free data retrieval call binding the contract method 0xa2967d99.
//
// Solidity: function getRollupExitRoot() view returns(bytes32)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCallerSession) GetRollupExitRoot() ([32]byte, error) {
	return _Feijoapolygonrollupmanager.Contract.GetRollupExitRoot(&_Feijoapolygonrollupmanager.CallOpts)
}

// GetRollupPendingStateTransitions is a free data retrieval call binding the contract method 0xb99d0ad7.
//
// Solidity: function getRollupPendingStateTransitions(uint32 rollupID, uint64 sequenceNum) view returns((uint64,uint64,bytes32,bytes32))
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCaller) GetRollupPendingStateTransitions(opts *bind.CallOpts, rollupID uint32, sequenceNum uint64) (PolygonRollupManagerPendingStateSequenceBased, error) {
	var out []interface{}
	err := _Feijoapolygonrollupmanager.contract.Call(opts, &out, "getRollupPendingStateTransitions", rollupID, sequenceNum)

	if err != nil {
		return *new(PolygonRollupManagerPendingStateSequenceBased), err
	}

	out0 := *abi.ConvertType(out[0], new(PolygonRollupManagerPendingStateSequenceBased)).(*PolygonRollupManagerPendingStateSequenceBased)

	return out0, err

}

// GetRollupPendingStateTransitions is a free data retrieval call binding the contract method 0xb99d0ad7.
//
// Solidity: function getRollupPendingStateTransitions(uint32 rollupID, uint64 sequenceNum) view returns((uint64,uint64,bytes32,bytes32))
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerSession) GetRollupPendingStateTransitions(rollupID uint32, sequenceNum uint64) (PolygonRollupManagerPendingStateSequenceBased, error) {
	return _Feijoapolygonrollupmanager.Contract.GetRollupPendingStateTransitions(&_Feijoapolygonrollupmanager.CallOpts, rollupID, sequenceNum)
}

// GetRollupPendingStateTransitions is a free data retrieval call binding the contract method 0xb99d0ad7.
//
// Solidity: function getRollupPendingStateTransitions(uint32 rollupID, uint64 sequenceNum) view returns((uint64,uint64,bytes32,bytes32))
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCallerSession) GetRollupPendingStateTransitions(rollupID uint32, sequenceNum uint64) (PolygonRollupManagerPendingStateSequenceBased, error) {
	return _Feijoapolygonrollupmanager.Contract.GetRollupPendingStateTransitions(&_Feijoapolygonrollupmanager.CallOpts, rollupID, sequenceNum)
}

// GetRollupSequencedSequences is a free data retrieval call binding the contract method 0xa9a77031.
//
// Solidity: function getRollupSequencedSequences(uint32 rollupID, uint64 sequenceNum) view returns((bytes32,uint64,uint64,uint128))
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCaller) GetRollupSequencedSequences(opts *bind.CallOpts, rollupID uint32, sequenceNum uint64) (PolygonRollupManagerSequencedData, error) {
	var out []interface{}
	err := _Feijoapolygonrollupmanager.contract.Call(opts, &out, "getRollupSequencedSequences", rollupID, sequenceNum)

	if err != nil {
		return *new(PolygonRollupManagerSequencedData), err
	}

	out0 := *abi.ConvertType(out[0], new(PolygonRollupManagerSequencedData)).(*PolygonRollupManagerSequencedData)

	return out0, err

}

// GetRollupSequencedSequences is a free data retrieval call binding the contract method 0xa9a77031.
//
// Solidity: function getRollupSequencedSequences(uint32 rollupID, uint64 sequenceNum) view returns((bytes32,uint64,uint64,uint128))
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerSession) GetRollupSequencedSequences(rollupID uint32, sequenceNum uint64) (PolygonRollupManagerSequencedData, error) {
	return _Feijoapolygonrollupmanager.Contract.GetRollupSequencedSequences(&_Feijoapolygonrollupmanager.CallOpts, rollupID, sequenceNum)
}

// GetRollupSequencedSequences is a free data retrieval call binding the contract method 0xa9a77031.
//
// Solidity: function getRollupSequencedSequences(uint32 rollupID, uint64 sequenceNum) view returns((bytes32,uint64,uint64,uint128))
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCallerSession) GetRollupSequencedSequences(rollupID uint32, sequenceNum uint64) (PolygonRollupManagerSequencedData, error) {
	return _Feijoapolygonrollupmanager.Contract.GetRollupSequencedSequences(&_Feijoapolygonrollupmanager.CallOpts, rollupID, sequenceNum)
}

// GetRollupsequenceNumToStateRoot is a free data retrieval call binding the contract method 0xeb142b40.
//
// Solidity: function getRollupsequenceNumToStateRoot(uint32 rollupID, uint64 sequenceNum) view returns(bytes32)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCaller) GetRollupsequenceNumToStateRoot(opts *bind.CallOpts, rollupID uint32, sequenceNum uint64) ([32]byte, error) {
	var out []interface{}
	err := _Feijoapolygonrollupmanager.contract.Call(opts, &out, "getRollupsequenceNumToStateRoot", rollupID, sequenceNum)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRollupsequenceNumToStateRoot is a free data retrieval call binding the contract method 0xeb142b40.
//
// Solidity: function getRollupsequenceNumToStateRoot(uint32 rollupID, uint64 sequenceNum) view returns(bytes32)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerSession) GetRollupsequenceNumToStateRoot(rollupID uint32, sequenceNum uint64) ([32]byte, error) {
	return _Feijoapolygonrollupmanager.Contract.GetRollupsequenceNumToStateRoot(&_Feijoapolygonrollupmanager.CallOpts, rollupID, sequenceNum)
}

// GetRollupsequenceNumToStateRoot is a free data retrieval call binding the contract method 0xeb142b40.
//
// Solidity: function getRollupsequenceNumToStateRoot(uint32 rollupID, uint64 sequenceNum) view returns(bytes32)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCallerSession) GetRollupsequenceNumToStateRoot(rollupID uint32, sequenceNum uint64) ([32]byte, error) {
	return _Feijoapolygonrollupmanager.Contract.GetRollupsequenceNumToStateRoot(&_Feijoapolygonrollupmanager.CallOpts, rollupID, sequenceNum)
}

// GetZkGasPrice is a free data retrieval call binding the contract method 0xf4174a17.
//
// Solidity: function getZkGasPrice() view returns(uint256)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCaller) GetZkGasPrice(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Feijoapolygonrollupmanager.contract.Call(opts, &out, "getZkGasPrice")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetZkGasPrice is a free data retrieval call binding the contract method 0xf4174a17.
//
// Solidity: function getZkGasPrice() view returns(uint256)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerSession) GetZkGasPrice() (*big.Int, error) {
	return _Feijoapolygonrollupmanager.Contract.GetZkGasPrice(&_Feijoapolygonrollupmanager.CallOpts)
}

// GetZkGasPrice is a free data retrieval call binding the contract method 0xf4174a17.
//
// Solidity: function getZkGasPrice() view returns(uint256)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCallerSession) GetZkGasPrice() (*big.Int, error) {
	return _Feijoapolygonrollupmanager.Contract.GetZkGasPrice(&_Feijoapolygonrollupmanager.CallOpts)
}

// GlobalExitRootManager is a free data retrieval call binding the contract method 0xd02103ca.
//
// Solidity: function globalExitRootManager() view returns(address)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCaller) GlobalExitRootManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Feijoapolygonrollupmanager.contract.Call(opts, &out, "globalExitRootManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GlobalExitRootManager is a free data retrieval call binding the contract method 0xd02103ca.
//
// Solidity: function globalExitRootManager() view returns(address)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerSession) GlobalExitRootManager() (common.Address, error) {
	return _Feijoapolygonrollupmanager.Contract.GlobalExitRootManager(&_Feijoapolygonrollupmanager.CallOpts)
}

// GlobalExitRootManager is a free data retrieval call binding the contract method 0xd02103ca.
//
// Solidity: function globalExitRootManager() view returns(address)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCallerSession) GlobalExitRootManager() (common.Address, error) {
	return _Feijoapolygonrollupmanager.Contract.GlobalExitRootManager(&_Feijoapolygonrollupmanager.CallOpts)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _Feijoapolygonrollupmanager.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _Feijoapolygonrollupmanager.Contract.HasRole(&_Feijoapolygonrollupmanager.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _Feijoapolygonrollupmanager.Contract.HasRole(&_Feijoapolygonrollupmanager.CallOpts, role, account)
}

// IsEmergencyState is a free data retrieval call binding the contract method 0x15064c96.
//
// Solidity: function isEmergencyState() view returns(bool)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCaller) IsEmergencyState(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Feijoapolygonrollupmanager.contract.Call(opts, &out, "isEmergencyState")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsEmergencyState is a free data retrieval call binding the contract method 0x15064c96.
//
// Solidity: function isEmergencyState() view returns(bool)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerSession) IsEmergencyState() (bool, error) {
	return _Feijoapolygonrollupmanager.Contract.IsEmergencyState(&_Feijoapolygonrollupmanager.CallOpts)
}

// IsEmergencyState is a free data retrieval call binding the contract method 0x15064c96.
//
// Solidity: function isEmergencyState() view returns(bool)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCallerSession) IsEmergencyState() (bool, error) {
	return _Feijoapolygonrollupmanager.Contract.IsEmergencyState(&_Feijoapolygonrollupmanager.CallOpts)
}

// IsPendingStateConsolidable is a free data retrieval call binding the contract method 0x080b3111.
//
// Solidity: function isPendingStateConsolidable(uint32 rollupID, uint64 pendingStateNum) view returns(bool)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCaller) IsPendingStateConsolidable(opts *bind.CallOpts, rollupID uint32, pendingStateNum uint64) (bool, error) {
	var out []interface{}
	err := _Feijoapolygonrollupmanager.contract.Call(opts, &out, "isPendingStateConsolidable", rollupID, pendingStateNum)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsPendingStateConsolidable is a free data retrieval call binding the contract method 0x080b3111.
//
// Solidity: function isPendingStateConsolidable(uint32 rollupID, uint64 pendingStateNum) view returns(bool)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerSession) IsPendingStateConsolidable(rollupID uint32, pendingStateNum uint64) (bool, error) {
	return _Feijoapolygonrollupmanager.Contract.IsPendingStateConsolidable(&_Feijoapolygonrollupmanager.CallOpts, rollupID, pendingStateNum)
}

// IsPendingStateConsolidable is a free data retrieval call binding the contract method 0x080b3111.
//
// Solidity: function isPendingStateConsolidable(uint32 rollupID, uint64 pendingStateNum) view returns(bool)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCallerSession) IsPendingStateConsolidable(rollupID uint32, pendingStateNum uint64) (bool, error) {
	return _Feijoapolygonrollupmanager.Contract.IsPendingStateConsolidable(&_Feijoapolygonrollupmanager.CallOpts, rollupID, pendingStateNum)
}

// LastAggregationTimestamp is a free data retrieval call binding the contract method 0xc1acbc34.
//
// Solidity: function lastAggregationTimestamp() view returns(uint64)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCaller) LastAggregationTimestamp(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Feijoapolygonrollupmanager.contract.Call(opts, &out, "lastAggregationTimestamp")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// LastAggregationTimestamp is a free data retrieval call binding the contract method 0xc1acbc34.
//
// Solidity: function lastAggregationTimestamp() view returns(uint64)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerSession) LastAggregationTimestamp() (uint64, error) {
	return _Feijoapolygonrollupmanager.Contract.LastAggregationTimestamp(&_Feijoapolygonrollupmanager.CallOpts)
}

// LastAggregationTimestamp is a free data retrieval call binding the contract method 0xc1acbc34.
//
// Solidity: function lastAggregationTimestamp() view returns(uint64)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCallerSession) LastAggregationTimestamp() (uint64, error) {
	return _Feijoapolygonrollupmanager.Contract.LastAggregationTimestamp(&_Feijoapolygonrollupmanager.CallOpts)
}

// LastDeactivatedEmergencyStateTimestamp is a free data retrieval call binding the contract method 0x30c27dde.
//
// Solidity: function lastDeactivatedEmergencyStateTimestamp() view returns(uint64)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCaller) LastDeactivatedEmergencyStateTimestamp(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Feijoapolygonrollupmanager.contract.Call(opts, &out, "lastDeactivatedEmergencyStateTimestamp")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// LastDeactivatedEmergencyStateTimestamp is a free data retrieval call binding the contract method 0x30c27dde.
//
// Solidity: function lastDeactivatedEmergencyStateTimestamp() view returns(uint64)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerSession) LastDeactivatedEmergencyStateTimestamp() (uint64, error) {
	return _Feijoapolygonrollupmanager.Contract.LastDeactivatedEmergencyStateTimestamp(&_Feijoapolygonrollupmanager.CallOpts)
}

// LastDeactivatedEmergencyStateTimestamp is a free data retrieval call binding the contract method 0x30c27dde.
//
// Solidity: function lastDeactivatedEmergencyStateTimestamp() view returns(uint64)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCallerSession) LastDeactivatedEmergencyStateTimestamp() (uint64, error) {
	return _Feijoapolygonrollupmanager.Contract.LastDeactivatedEmergencyStateTimestamp(&_Feijoapolygonrollupmanager.CallOpts)
}

// MultiplierZkGasPrice is a free data retrieval call binding the contract method 0x9ff22cb5.
//
// Solidity: function multiplierZkGasPrice() view returns(uint16)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCaller) MultiplierZkGasPrice(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _Feijoapolygonrollupmanager.contract.Call(opts, &out, "multiplierZkGasPrice")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// MultiplierZkGasPrice is a free data retrieval call binding the contract method 0x9ff22cb5.
//
// Solidity: function multiplierZkGasPrice() view returns(uint16)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerSession) MultiplierZkGasPrice() (uint16, error) {
	return _Feijoapolygonrollupmanager.Contract.MultiplierZkGasPrice(&_Feijoapolygonrollupmanager.CallOpts)
}

// MultiplierZkGasPrice is a free data retrieval call binding the contract method 0x9ff22cb5.
//
// Solidity: function multiplierZkGasPrice() view returns(uint16)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCallerSession) MultiplierZkGasPrice() (uint16, error) {
	return _Feijoapolygonrollupmanager.Contract.MultiplierZkGasPrice(&_Feijoapolygonrollupmanager.CallOpts)
}

// PendingStateTimeout is a free data retrieval call binding the contract method 0xd939b315.
//
// Solidity: function pendingStateTimeout() view returns(uint64)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCaller) PendingStateTimeout(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Feijoapolygonrollupmanager.contract.Call(opts, &out, "pendingStateTimeout")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// PendingStateTimeout is a free data retrieval call binding the contract method 0xd939b315.
//
// Solidity: function pendingStateTimeout() view returns(uint64)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerSession) PendingStateTimeout() (uint64, error) {
	return _Feijoapolygonrollupmanager.Contract.PendingStateTimeout(&_Feijoapolygonrollupmanager.CallOpts)
}

// PendingStateTimeout is a free data retrieval call binding the contract method 0xd939b315.
//
// Solidity: function pendingStateTimeout() view returns(uint64)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCallerSession) PendingStateTimeout() (uint64, error) {
	return _Feijoapolygonrollupmanager.Contract.PendingStateTimeout(&_Feijoapolygonrollupmanager.CallOpts)
}

// Pol is a free data retrieval call binding the contract method 0xe46761c4.
//
// Solidity: function pol() view returns(address)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCaller) Pol(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Feijoapolygonrollupmanager.contract.Call(opts, &out, "pol")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Pol is a free data retrieval call binding the contract method 0xe46761c4.
//
// Solidity: function pol() view returns(address)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerSession) Pol() (common.Address, error) {
	return _Feijoapolygonrollupmanager.Contract.Pol(&_Feijoapolygonrollupmanager.CallOpts)
}

// Pol is a free data retrieval call binding the contract method 0xe46761c4.
//
// Solidity: function pol() view returns(address)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCallerSession) Pol() (common.Address, error) {
	return _Feijoapolygonrollupmanager.Contract.Pol(&_Feijoapolygonrollupmanager.CallOpts)
}

// RollupAddressToID is a free data retrieval call binding the contract method 0xceee281d.
//
// Solidity: function rollupAddressToID(address rollupAddress) view returns(uint32 rollupID)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCaller) RollupAddressToID(opts *bind.CallOpts, rollupAddress common.Address) (uint32, error) {
	var out []interface{}
	err := _Feijoapolygonrollupmanager.contract.Call(opts, &out, "rollupAddressToID", rollupAddress)

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// RollupAddressToID is a free data retrieval call binding the contract method 0xceee281d.
//
// Solidity: function rollupAddressToID(address rollupAddress) view returns(uint32 rollupID)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerSession) RollupAddressToID(rollupAddress common.Address) (uint32, error) {
	return _Feijoapolygonrollupmanager.Contract.RollupAddressToID(&_Feijoapolygonrollupmanager.CallOpts, rollupAddress)
}

// RollupAddressToID is a free data retrieval call binding the contract method 0xceee281d.
//
// Solidity: function rollupAddressToID(address rollupAddress) view returns(uint32 rollupID)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCallerSession) RollupAddressToID(rollupAddress common.Address) (uint32, error) {
	return _Feijoapolygonrollupmanager.Contract.RollupAddressToID(&_Feijoapolygonrollupmanager.CallOpts, rollupAddress)
}

// RollupCount is a free data retrieval call binding the contract method 0xf4e92675.
//
// Solidity: function rollupCount() view returns(uint32)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCaller) RollupCount(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _Feijoapolygonrollupmanager.contract.Call(opts, &out, "rollupCount")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// RollupCount is a free data retrieval call binding the contract method 0xf4e92675.
//
// Solidity: function rollupCount() view returns(uint32)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerSession) RollupCount() (uint32, error) {
	return _Feijoapolygonrollupmanager.Contract.RollupCount(&_Feijoapolygonrollupmanager.CallOpts)
}

// RollupCount is a free data retrieval call binding the contract method 0xf4e92675.
//
// Solidity: function rollupCount() view returns(uint32)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCallerSession) RollupCount() (uint32, error) {
	return _Feijoapolygonrollupmanager.Contract.RollupCount(&_Feijoapolygonrollupmanager.CallOpts)
}

// RollupIDToRollupData is a free data retrieval call binding the contract method 0xf9c4c2ae.
//
// Solidity: function rollupIDToRollupData(uint32 rollupID) view returns(address rollupContract, uint64 chainID, address verifier, uint64 forkID, bytes32 lastLocalExitRoot, uint64 lastSequenceNum, uint64 lastVerifiedSequenceNum, uint64 lastPendingState, uint64 lastPendingStateConsolidated, uint64 lastVerifiedSequenceBeforeUpgrade, uint64 rollupTypeID, uint8 rollupCompatibilityID)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCaller) RollupIDToRollupData(opts *bind.CallOpts, rollupID uint32) (struct {
	RollupContract                    common.Address
	ChainID                           uint64
	Verifier                          common.Address
	ForkID                            uint64
	LastLocalExitRoot                 [32]byte
	LastSequenceNum                   uint64
	LastVerifiedSequenceNum           uint64
	LastPendingState                  uint64
	LastPendingStateConsolidated      uint64
	LastVerifiedSequenceBeforeUpgrade uint64
	RollupTypeID                      uint64
	RollupCompatibilityID             uint8
}, error) {
	var out []interface{}
	err := _Feijoapolygonrollupmanager.contract.Call(opts, &out, "rollupIDToRollupData", rollupID)

	outstruct := new(struct {
		RollupContract                    common.Address
		ChainID                           uint64
		Verifier                          common.Address
		ForkID                            uint64
		LastLocalExitRoot                 [32]byte
		LastSequenceNum                   uint64
		LastVerifiedSequenceNum           uint64
		LastPendingState                  uint64
		LastPendingStateConsolidated      uint64
		LastVerifiedSequenceBeforeUpgrade uint64
		RollupTypeID                      uint64
		RollupCompatibilityID             uint8
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.RollupContract = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.ChainID = *abi.ConvertType(out[1], new(uint64)).(*uint64)
	outstruct.Verifier = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)
	outstruct.ForkID = *abi.ConvertType(out[3], new(uint64)).(*uint64)
	outstruct.LastLocalExitRoot = *abi.ConvertType(out[4], new([32]byte)).(*[32]byte)
	outstruct.LastSequenceNum = *abi.ConvertType(out[5], new(uint64)).(*uint64)
	outstruct.LastVerifiedSequenceNum = *abi.ConvertType(out[6], new(uint64)).(*uint64)
	outstruct.LastPendingState = *abi.ConvertType(out[7], new(uint64)).(*uint64)
	outstruct.LastPendingStateConsolidated = *abi.ConvertType(out[8], new(uint64)).(*uint64)
	outstruct.LastVerifiedSequenceBeforeUpgrade = *abi.ConvertType(out[9], new(uint64)).(*uint64)
	outstruct.RollupTypeID = *abi.ConvertType(out[10], new(uint64)).(*uint64)
	outstruct.RollupCompatibilityID = *abi.ConvertType(out[11], new(uint8)).(*uint8)

	return *outstruct, err

}

// RollupIDToRollupData is a free data retrieval call binding the contract method 0xf9c4c2ae.
//
// Solidity: function rollupIDToRollupData(uint32 rollupID) view returns(address rollupContract, uint64 chainID, address verifier, uint64 forkID, bytes32 lastLocalExitRoot, uint64 lastSequenceNum, uint64 lastVerifiedSequenceNum, uint64 lastPendingState, uint64 lastPendingStateConsolidated, uint64 lastVerifiedSequenceBeforeUpgrade, uint64 rollupTypeID, uint8 rollupCompatibilityID)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerSession) RollupIDToRollupData(rollupID uint32) (struct {
	RollupContract                    common.Address
	ChainID                           uint64
	Verifier                          common.Address
	ForkID                            uint64
	LastLocalExitRoot                 [32]byte
	LastSequenceNum                   uint64
	LastVerifiedSequenceNum           uint64
	LastPendingState                  uint64
	LastPendingStateConsolidated      uint64
	LastVerifiedSequenceBeforeUpgrade uint64
	RollupTypeID                      uint64
	RollupCompatibilityID             uint8
}, error) {
	return _Feijoapolygonrollupmanager.Contract.RollupIDToRollupData(&_Feijoapolygonrollupmanager.CallOpts, rollupID)
}

// RollupIDToRollupData is a free data retrieval call binding the contract method 0xf9c4c2ae.
//
// Solidity: function rollupIDToRollupData(uint32 rollupID) view returns(address rollupContract, uint64 chainID, address verifier, uint64 forkID, bytes32 lastLocalExitRoot, uint64 lastSequenceNum, uint64 lastVerifiedSequenceNum, uint64 lastPendingState, uint64 lastPendingStateConsolidated, uint64 lastVerifiedSequenceBeforeUpgrade, uint64 rollupTypeID, uint8 rollupCompatibilityID)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCallerSession) RollupIDToRollupData(rollupID uint32) (struct {
	RollupContract                    common.Address
	ChainID                           uint64
	Verifier                          common.Address
	ForkID                            uint64
	LastLocalExitRoot                 [32]byte
	LastSequenceNum                   uint64
	LastVerifiedSequenceNum           uint64
	LastPendingState                  uint64
	LastPendingStateConsolidated      uint64
	LastVerifiedSequenceBeforeUpgrade uint64
	RollupTypeID                      uint64
	RollupCompatibilityID             uint8
}, error) {
	return _Feijoapolygonrollupmanager.Contract.RollupIDToRollupData(&_Feijoapolygonrollupmanager.CallOpts, rollupID)
}

// RollupTypeCount is a free data retrieval call binding the contract method 0x1796a1ae.
//
// Solidity: function rollupTypeCount() view returns(uint32)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCaller) RollupTypeCount(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _Feijoapolygonrollupmanager.contract.Call(opts, &out, "rollupTypeCount")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// RollupTypeCount is a free data retrieval call binding the contract method 0x1796a1ae.
//
// Solidity: function rollupTypeCount() view returns(uint32)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerSession) RollupTypeCount() (uint32, error) {
	return _Feijoapolygonrollupmanager.Contract.RollupTypeCount(&_Feijoapolygonrollupmanager.CallOpts)
}

// RollupTypeCount is a free data retrieval call binding the contract method 0x1796a1ae.
//
// Solidity: function rollupTypeCount() view returns(uint32)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCallerSession) RollupTypeCount() (uint32, error) {
	return _Feijoapolygonrollupmanager.Contract.RollupTypeCount(&_Feijoapolygonrollupmanager.CallOpts)
}

// RollupTypeMap is a free data retrieval call binding the contract method 0x65c0504d.
//
// Solidity: function rollupTypeMap(uint32 rollupTypeID) view returns(address consensusImplementation, address verifier, uint64 forkID, uint8 rollupCompatibilityID, bool obsolete, bytes32 genesis)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCaller) RollupTypeMap(opts *bind.CallOpts, rollupTypeID uint32) (struct {
	ConsensusImplementation common.Address
	Verifier                common.Address
	ForkID                  uint64
	RollupCompatibilityID   uint8
	Obsolete                bool
	Genesis                 [32]byte
}, error) {
	var out []interface{}
	err := _Feijoapolygonrollupmanager.contract.Call(opts, &out, "rollupTypeMap", rollupTypeID)

	outstruct := new(struct {
		ConsensusImplementation common.Address
		Verifier                common.Address
		ForkID                  uint64
		RollupCompatibilityID   uint8
		Obsolete                bool
		Genesis                 [32]byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.ConsensusImplementation = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Verifier = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.ForkID = *abi.ConvertType(out[2], new(uint64)).(*uint64)
	outstruct.RollupCompatibilityID = *abi.ConvertType(out[3], new(uint8)).(*uint8)
	outstruct.Obsolete = *abi.ConvertType(out[4], new(bool)).(*bool)
	outstruct.Genesis = *abi.ConvertType(out[5], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

// RollupTypeMap is a free data retrieval call binding the contract method 0x65c0504d.
//
// Solidity: function rollupTypeMap(uint32 rollupTypeID) view returns(address consensusImplementation, address verifier, uint64 forkID, uint8 rollupCompatibilityID, bool obsolete, bytes32 genesis)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerSession) RollupTypeMap(rollupTypeID uint32) (struct {
	ConsensusImplementation common.Address
	Verifier                common.Address
	ForkID                  uint64
	RollupCompatibilityID   uint8
	Obsolete                bool
	Genesis                 [32]byte
}, error) {
	return _Feijoapolygonrollupmanager.Contract.RollupTypeMap(&_Feijoapolygonrollupmanager.CallOpts, rollupTypeID)
}

// RollupTypeMap is a free data retrieval call binding the contract method 0x65c0504d.
//
// Solidity: function rollupTypeMap(uint32 rollupTypeID) view returns(address consensusImplementation, address verifier, uint64 forkID, uint8 rollupCompatibilityID, bool obsolete, bytes32 genesis)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCallerSession) RollupTypeMap(rollupTypeID uint32) (struct {
	ConsensusImplementation common.Address
	Verifier                common.Address
	ForkID                  uint64
	RollupCompatibilityID   uint8
	Obsolete                bool
	Genesis                 [32]byte
}, error) {
	return _Feijoapolygonrollupmanager.Contract.RollupTypeMap(&_Feijoapolygonrollupmanager.CallOpts, rollupTypeID)
}

// TotalVerifiedZkGasLimit is a free data retrieval call binding the contract method 0x27696c5e.
//
// Solidity: function totalVerifiedZkGasLimit() view returns(uint128)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCaller) TotalVerifiedZkGasLimit(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Feijoapolygonrollupmanager.contract.Call(opts, &out, "totalVerifiedZkGasLimit")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalVerifiedZkGasLimit is a free data retrieval call binding the contract method 0x27696c5e.
//
// Solidity: function totalVerifiedZkGasLimit() view returns(uint128)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerSession) TotalVerifiedZkGasLimit() (*big.Int, error) {
	return _Feijoapolygonrollupmanager.Contract.TotalVerifiedZkGasLimit(&_Feijoapolygonrollupmanager.CallOpts)
}

// TotalVerifiedZkGasLimit is a free data retrieval call binding the contract method 0x27696c5e.
//
// Solidity: function totalVerifiedZkGasLimit() view returns(uint128)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCallerSession) TotalVerifiedZkGasLimit() (*big.Int, error) {
	return _Feijoapolygonrollupmanager.Contract.TotalVerifiedZkGasLimit(&_Feijoapolygonrollupmanager.CallOpts)
}

// TotalZkGasLimit is a free data retrieval call binding the contract method 0x6c6be9eb.
//
// Solidity: function totalZkGasLimit() view returns(uint128)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCaller) TotalZkGasLimit(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Feijoapolygonrollupmanager.contract.Call(opts, &out, "totalZkGasLimit")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalZkGasLimit is a free data retrieval call binding the contract method 0x6c6be9eb.
//
// Solidity: function totalZkGasLimit() view returns(uint128)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerSession) TotalZkGasLimit() (*big.Int, error) {
	return _Feijoapolygonrollupmanager.Contract.TotalZkGasLimit(&_Feijoapolygonrollupmanager.CallOpts)
}

// TotalZkGasLimit is a free data retrieval call binding the contract method 0x6c6be9eb.
//
// Solidity: function totalZkGasLimit() view returns(uint128)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCallerSession) TotalZkGasLimit() (*big.Int, error) {
	return _Feijoapolygonrollupmanager.Contract.TotalZkGasLimit(&_Feijoapolygonrollupmanager.CallOpts)
}

// TrustedAggregatorTimeout is a free data retrieval call binding the contract method 0x841b24d7.
//
// Solidity: function trustedAggregatorTimeout() view returns(uint64)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCaller) TrustedAggregatorTimeout(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Feijoapolygonrollupmanager.contract.Call(opts, &out, "trustedAggregatorTimeout")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// TrustedAggregatorTimeout is a free data retrieval call binding the contract method 0x841b24d7.
//
// Solidity: function trustedAggregatorTimeout() view returns(uint64)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerSession) TrustedAggregatorTimeout() (uint64, error) {
	return _Feijoapolygonrollupmanager.Contract.TrustedAggregatorTimeout(&_Feijoapolygonrollupmanager.CallOpts)
}

// TrustedAggregatorTimeout is a free data retrieval call binding the contract method 0x841b24d7.
//
// Solidity: function trustedAggregatorTimeout() view returns(uint64)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCallerSession) TrustedAggregatorTimeout() (uint64, error) {
	return _Feijoapolygonrollupmanager.Contract.TrustedAggregatorTimeout(&_Feijoapolygonrollupmanager.CallOpts)
}

// VerifySequenceTimeTarget is a free data retrieval call binding the contract method 0xb7397536.
//
// Solidity: function verifySequenceTimeTarget() view returns(uint64)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCaller) VerifySequenceTimeTarget(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Feijoapolygonrollupmanager.contract.Call(opts, &out, "verifySequenceTimeTarget")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// VerifySequenceTimeTarget is a free data retrieval call binding the contract method 0xb7397536.
//
// Solidity: function verifySequenceTimeTarget() view returns(uint64)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerSession) VerifySequenceTimeTarget() (uint64, error) {
	return _Feijoapolygonrollupmanager.Contract.VerifySequenceTimeTarget(&_Feijoapolygonrollupmanager.CallOpts)
}

// VerifySequenceTimeTarget is a free data retrieval call binding the contract method 0xb7397536.
//
// Solidity: function verifySequenceTimeTarget() view returns(uint64)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerCallerSession) VerifySequenceTimeTarget() (uint64, error) {
	return _Feijoapolygonrollupmanager.Contract.VerifySequenceTimeTarget(&_Feijoapolygonrollupmanager.CallOpts)
}

// ActivateEmergencyState is a paid mutator transaction binding the contract method 0x2072f6c5.
//
// Solidity: function activateEmergencyState() returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerTransactor) ActivateEmergencyState(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.contract.Transact(opts, "activateEmergencyState")
}

// ActivateEmergencyState is a paid mutator transaction binding the contract method 0x2072f6c5.
//
// Solidity: function activateEmergencyState() returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerSession) ActivateEmergencyState() (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.Contract.ActivateEmergencyState(&_Feijoapolygonrollupmanager.TransactOpts)
}

// ActivateEmergencyState is a paid mutator transaction binding the contract method 0x2072f6c5.
//
// Solidity: function activateEmergencyState() returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerTransactorSession) ActivateEmergencyState() (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.Contract.ActivateEmergencyState(&_Feijoapolygonrollupmanager.TransactOpts)
}

// AddExistingRollup is a paid mutator transaction binding the contract method 0xe0bfd3d2.
//
// Solidity: function addExistingRollup(address rollupAddress, address verifier, uint64 forkID, uint64 chainID, bytes32 genesis, uint8 rollupCompatibilityID) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerTransactor) AddExistingRollup(opts *bind.TransactOpts, rollupAddress common.Address, verifier common.Address, forkID uint64, chainID uint64, genesis [32]byte, rollupCompatibilityID uint8) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.contract.Transact(opts, "addExistingRollup", rollupAddress, verifier, forkID, chainID, genesis, rollupCompatibilityID)
}

// AddExistingRollup is a paid mutator transaction binding the contract method 0xe0bfd3d2.
//
// Solidity: function addExistingRollup(address rollupAddress, address verifier, uint64 forkID, uint64 chainID, bytes32 genesis, uint8 rollupCompatibilityID) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerSession) AddExistingRollup(rollupAddress common.Address, verifier common.Address, forkID uint64, chainID uint64, genesis [32]byte, rollupCompatibilityID uint8) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.Contract.AddExistingRollup(&_Feijoapolygonrollupmanager.TransactOpts, rollupAddress, verifier, forkID, chainID, genesis, rollupCompatibilityID)
}

// AddExistingRollup is a paid mutator transaction binding the contract method 0xe0bfd3d2.
//
// Solidity: function addExistingRollup(address rollupAddress, address verifier, uint64 forkID, uint64 chainID, bytes32 genesis, uint8 rollupCompatibilityID) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerTransactorSession) AddExistingRollup(rollupAddress common.Address, verifier common.Address, forkID uint64, chainID uint64, genesis [32]byte, rollupCompatibilityID uint8) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.Contract.AddExistingRollup(&_Feijoapolygonrollupmanager.TransactOpts, rollupAddress, verifier, forkID, chainID, genesis, rollupCompatibilityID)
}

// AddNewRollupType is a paid mutator transaction binding the contract method 0xf34eb8eb.
//
// Solidity: function addNewRollupType(address consensusImplementation, address verifier, uint64 forkID, uint8 rollupCompatibilityID, bytes32 genesis, string description) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerTransactor) AddNewRollupType(opts *bind.TransactOpts, consensusImplementation common.Address, verifier common.Address, forkID uint64, rollupCompatibilityID uint8, genesis [32]byte, description string) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.contract.Transact(opts, "addNewRollupType", consensusImplementation, verifier, forkID, rollupCompatibilityID, genesis, description)
}

// AddNewRollupType is a paid mutator transaction binding the contract method 0xf34eb8eb.
//
// Solidity: function addNewRollupType(address consensusImplementation, address verifier, uint64 forkID, uint8 rollupCompatibilityID, bytes32 genesis, string description) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerSession) AddNewRollupType(consensusImplementation common.Address, verifier common.Address, forkID uint64, rollupCompatibilityID uint8, genesis [32]byte, description string) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.Contract.AddNewRollupType(&_Feijoapolygonrollupmanager.TransactOpts, consensusImplementation, verifier, forkID, rollupCompatibilityID, genesis, description)
}

// AddNewRollupType is a paid mutator transaction binding the contract method 0xf34eb8eb.
//
// Solidity: function addNewRollupType(address consensusImplementation, address verifier, uint64 forkID, uint8 rollupCompatibilityID, bytes32 genesis, string description) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerTransactorSession) AddNewRollupType(consensusImplementation common.Address, verifier common.Address, forkID uint64, rollupCompatibilityID uint8, genesis [32]byte, description string) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.Contract.AddNewRollupType(&_Feijoapolygonrollupmanager.TransactOpts, consensusImplementation, verifier, forkID, rollupCompatibilityID, genesis, description)
}

// ConsolidatePendingState is a paid mutator transaction binding the contract method 0x1608859c.
//
// Solidity: function consolidatePendingState(uint32 rollupID, uint64 pendingStateNum) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerTransactor) ConsolidatePendingState(opts *bind.TransactOpts, rollupID uint32, pendingStateNum uint64) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.contract.Transact(opts, "consolidatePendingState", rollupID, pendingStateNum)
}

// ConsolidatePendingState is a paid mutator transaction binding the contract method 0x1608859c.
//
// Solidity: function consolidatePendingState(uint32 rollupID, uint64 pendingStateNum) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerSession) ConsolidatePendingState(rollupID uint32, pendingStateNum uint64) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.Contract.ConsolidatePendingState(&_Feijoapolygonrollupmanager.TransactOpts, rollupID, pendingStateNum)
}

// ConsolidatePendingState is a paid mutator transaction binding the contract method 0x1608859c.
//
// Solidity: function consolidatePendingState(uint32 rollupID, uint64 pendingStateNum) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerTransactorSession) ConsolidatePendingState(rollupID uint32, pendingStateNum uint64) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.Contract.ConsolidatePendingState(&_Feijoapolygonrollupmanager.TransactOpts, rollupID, pendingStateNum)
}

// CreateNewRollup is a paid mutator transaction binding the contract method 0x727885e9.
//
// Solidity: function createNewRollup(uint32 rollupTypeID, uint64 chainID, address admin, address sequencer, address gasTokenAddress, string sequencerURL, string networkName) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerTransactor) CreateNewRollup(opts *bind.TransactOpts, rollupTypeID uint32, chainID uint64, admin common.Address, sequencer common.Address, gasTokenAddress common.Address, sequencerURL string, networkName string) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.contract.Transact(opts, "createNewRollup", rollupTypeID, chainID, admin, sequencer, gasTokenAddress, sequencerURL, networkName)
}

// CreateNewRollup is a paid mutator transaction binding the contract method 0x727885e9.
//
// Solidity: function createNewRollup(uint32 rollupTypeID, uint64 chainID, address admin, address sequencer, address gasTokenAddress, string sequencerURL, string networkName) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerSession) CreateNewRollup(rollupTypeID uint32, chainID uint64, admin common.Address, sequencer common.Address, gasTokenAddress common.Address, sequencerURL string, networkName string) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.Contract.CreateNewRollup(&_Feijoapolygonrollupmanager.TransactOpts, rollupTypeID, chainID, admin, sequencer, gasTokenAddress, sequencerURL, networkName)
}

// CreateNewRollup is a paid mutator transaction binding the contract method 0x727885e9.
//
// Solidity: function createNewRollup(uint32 rollupTypeID, uint64 chainID, address admin, address sequencer, address gasTokenAddress, string sequencerURL, string networkName) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerTransactorSession) CreateNewRollup(rollupTypeID uint32, chainID uint64, admin common.Address, sequencer common.Address, gasTokenAddress common.Address, sequencerURL string, networkName string) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.Contract.CreateNewRollup(&_Feijoapolygonrollupmanager.TransactOpts, rollupTypeID, chainID, admin, sequencer, gasTokenAddress, sequencerURL, networkName)
}

// DeactivateEmergencyState is a paid mutator transaction binding the contract method 0xdbc16976.
//
// Solidity: function deactivateEmergencyState() returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerTransactor) DeactivateEmergencyState(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.contract.Transact(opts, "deactivateEmergencyState")
}

// DeactivateEmergencyState is a paid mutator transaction binding the contract method 0xdbc16976.
//
// Solidity: function deactivateEmergencyState() returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerSession) DeactivateEmergencyState() (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.Contract.DeactivateEmergencyState(&_Feijoapolygonrollupmanager.TransactOpts)
}

// DeactivateEmergencyState is a paid mutator transaction binding the contract method 0xdbc16976.
//
// Solidity: function deactivateEmergencyState() returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerTransactorSession) DeactivateEmergencyState() (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.Contract.DeactivateEmergencyState(&_Feijoapolygonrollupmanager.TransactOpts)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.Contract.GrantRole(&_Feijoapolygonrollupmanager.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.Contract.GrantRole(&_Feijoapolygonrollupmanager.TransactOpts, role, account)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerSession) Initialize() (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.Contract.Initialize(&_Feijoapolygonrollupmanager.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerTransactorSession) Initialize() (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.Contract.Initialize(&_Feijoapolygonrollupmanager.TransactOpts)
}

// ObsoleteRollupType is a paid mutator transaction binding the contract method 0x7222020f.
//
// Solidity: function obsoleteRollupType(uint32 rollupTypeID) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerTransactor) ObsoleteRollupType(opts *bind.TransactOpts, rollupTypeID uint32) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.contract.Transact(opts, "obsoleteRollupType", rollupTypeID)
}

// ObsoleteRollupType is a paid mutator transaction binding the contract method 0x7222020f.
//
// Solidity: function obsoleteRollupType(uint32 rollupTypeID) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerSession) ObsoleteRollupType(rollupTypeID uint32) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.Contract.ObsoleteRollupType(&_Feijoapolygonrollupmanager.TransactOpts, rollupTypeID)
}

// ObsoleteRollupType is a paid mutator transaction binding the contract method 0x7222020f.
//
// Solidity: function obsoleteRollupType(uint32 rollupTypeID) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerTransactorSession) ObsoleteRollupType(rollupTypeID uint32) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.Contract.ObsoleteRollupType(&_Feijoapolygonrollupmanager.TransactOpts, rollupTypeID)
}

// OnSequence is a paid mutator transaction binding the contract method 0xfe01d89e.
//
// Solidity: function onSequence(uint128 zkGasLimitSequenced, uint64 blobsSequenced, bytes32 newAccInputHash) returns(uint64)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerTransactor) OnSequence(opts *bind.TransactOpts, zkGasLimitSequenced *big.Int, blobsSequenced uint64, newAccInputHash [32]byte) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.contract.Transact(opts, "onSequence", zkGasLimitSequenced, blobsSequenced, newAccInputHash)
}

// OnSequence is a paid mutator transaction binding the contract method 0xfe01d89e.
//
// Solidity: function onSequence(uint128 zkGasLimitSequenced, uint64 blobsSequenced, bytes32 newAccInputHash) returns(uint64)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerSession) OnSequence(zkGasLimitSequenced *big.Int, blobsSequenced uint64, newAccInputHash [32]byte) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.Contract.OnSequence(&_Feijoapolygonrollupmanager.TransactOpts, zkGasLimitSequenced, blobsSequenced, newAccInputHash)
}

// OnSequence is a paid mutator transaction binding the contract method 0xfe01d89e.
//
// Solidity: function onSequence(uint128 zkGasLimitSequenced, uint64 blobsSequenced, bytes32 newAccInputHash) returns(uint64)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerTransactorSession) OnSequence(zkGasLimitSequenced *big.Int, blobsSequenced uint64, newAccInputHash [32]byte) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.Contract.OnSequence(&_Feijoapolygonrollupmanager.TransactOpts, zkGasLimitSequenced, blobsSequenced, newAccInputHash)
}

// OverridePendingState is a paid mutator transaction binding the contract method 0x12b86e19.
//
// Solidity: function overridePendingState(uint32 rollupID, uint64 initPendingStateNum, uint64 finalPendingStateNum, uint64 initSequenceNum, uint64 finalSequenceNum, bytes32 newLocalExitRoot, bytes32 newStateRoot, bytes32[24] proof) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerTransactor) OverridePendingState(opts *bind.TransactOpts, rollupID uint32, initPendingStateNum uint64, finalPendingStateNum uint64, initSequenceNum uint64, finalSequenceNum uint64, newLocalExitRoot [32]byte, newStateRoot [32]byte, proof [24][32]byte) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.contract.Transact(opts, "overridePendingState", rollupID, initPendingStateNum, finalPendingStateNum, initSequenceNum, finalSequenceNum, newLocalExitRoot, newStateRoot, proof)
}

// OverridePendingState is a paid mutator transaction binding the contract method 0x12b86e19.
//
// Solidity: function overridePendingState(uint32 rollupID, uint64 initPendingStateNum, uint64 finalPendingStateNum, uint64 initSequenceNum, uint64 finalSequenceNum, bytes32 newLocalExitRoot, bytes32 newStateRoot, bytes32[24] proof) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerSession) OverridePendingState(rollupID uint32, initPendingStateNum uint64, finalPendingStateNum uint64, initSequenceNum uint64, finalSequenceNum uint64, newLocalExitRoot [32]byte, newStateRoot [32]byte, proof [24][32]byte) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.Contract.OverridePendingState(&_Feijoapolygonrollupmanager.TransactOpts, rollupID, initPendingStateNum, finalPendingStateNum, initSequenceNum, finalSequenceNum, newLocalExitRoot, newStateRoot, proof)
}

// OverridePendingState is a paid mutator transaction binding the contract method 0x12b86e19.
//
// Solidity: function overridePendingState(uint32 rollupID, uint64 initPendingStateNum, uint64 finalPendingStateNum, uint64 initSequenceNum, uint64 finalSequenceNum, bytes32 newLocalExitRoot, bytes32 newStateRoot, bytes32[24] proof) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerTransactorSession) OverridePendingState(rollupID uint32, initPendingStateNum uint64, finalPendingStateNum uint64, initSequenceNum uint64, finalSequenceNum uint64, newLocalExitRoot [32]byte, newStateRoot [32]byte, proof [24][32]byte) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.Contract.OverridePendingState(&_Feijoapolygonrollupmanager.TransactOpts, rollupID, initPendingStateNum, finalPendingStateNum, initSequenceNum, finalSequenceNum, newLocalExitRoot, newStateRoot, proof)
}

// ProveNonDeterministicPendingState is a paid mutator transaction binding the contract method 0x8bd4f071.
//
// Solidity: function proveNonDeterministicPendingState(uint32 rollupID, uint64 initPendingStateNum, uint64 finalPendingStateNum, uint64 initSequenceNum, uint64 finalSequenceNum, bytes32 newLocalExitRoot, bytes32 newStateRoot, bytes32[24] proof) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerTransactor) ProveNonDeterministicPendingState(opts *bind.TransactOpts, rollupID uint32, initPendingStateNum uint64, finalPendingStateNum uint64, initSequenceNum uint64, finalSequenceNum uint64, newLocalExitRoot [32]byte, newStateRoot [32]byte, proof [24][32]byte) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.contract.Transact(opts, "proveNonDeterministicPendingState", rollupID, initPendingStateNum, finalPendingStateNum, initSequenceNum, finalSequenceNum, newLocalExitRoot, newStateRoot, proof)
}

// ProveNonDeterministicPendingState is a paid mutator transaction binding the contract method 0x8bd4f071.
//
// Solidity: function proveNonDeterministicPendingState(uint32 rollupID, uint64 initPendingStateNum, uint64 finalPendingStateNum, uint64 initSequenceNum, uint64 finalSequenceNum, bytes32 newLocalExitRoot, bytes32 newStateRoot, bytes32[24] proof) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerSession) ProveNonDeterministicPendingState(rollupID uint32, initPendingStateNum uint64, finalPendingStateNum uint64, initSequenceNum uint64, finalSequenceNum uint64, newLocalExitRoot [32]byte, newStateRoot [32]byte, proof [24][32]byte) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.Contract.ProveNonDeterministicPendingState(&_Feijoapolygonrollupmanager.TransactOpts, rollupID, initPendingStateNum, finalPendingStateNum, initSequenceNum, finalSequenceNum, newLocalExitRoot, newStateRoot, proof)
}

// ProveNonDeterministicPendingState is a paid mutator transaction binding the contract method 0x8bd4f071.
//
// Solidity: function proveNonDeterministicPendingState(uint32 rollupID, uint64 initPendingStateNum, uint64 finalPendingStateNum, uint64 initSequenceNum, uint64 finalSequenceNum, bytes32 newLocalExitRoot, bytes32 newStateRoot, bytes32[24] proof) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerTransactorSession) ProveNonDeterministicPendingState(rollupID uint32, initPendingStateNum uint64, finalPendingStateNum uint64, initSequenceNum uint64, finalSequenceNum uint64, newLocalExitRoot [32]byte, newStateRoot [32]byte, proof [24][32]byte) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.Contract.ProveNonDeterministicPendingState(&_Feijoapolygonrollupmanager.TransactOpts, rollupID, initPendingStateNum, finalPendingStateNum, initSequenceNum, finalSequenceNum, newLocalExitRoot, newStateRoot, proof)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.contract.Transact(opts, "renounceRole", role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.Contract.RenounceRole(&_Feijoapolygonrollupmanager.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerTransactorSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.Contract.RenounceRole(&_Feijoapolygonrollupmanager.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.Contract.RevokeRole(&_Feijoapolygonrollupmanager.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.Contract.RevokeRole(&_Feijoapolygonrollupmanager.TransactOpts, role, account)
}

// SetAggregateRollupVerifier is a paid mutator transaction binding the contract method 0xe2bfe8b3.
//
// Solidity: function setAggregateRollupVerifier(address newAggregateRollupVerifier) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerTransactor) SetAggregateRollupVerifier(opts *bind.TransactOpts, newAggregateRollupVerifier common.Address) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.contract.Transact(opts, "setAggregateRollupVerifier", newAggregateRollupVerifier)
}

// SetAggregateRollupVerifier is a paid mutator transaction binding the contract method 0xe2bfe8b3.
//
// Solidity: function setAggregateRollupVerifier(address newAggregateRollupVerifier) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerSession) SetAggregateRollupVerifier(newAggregateRollupVerifier common.Address) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.Contract.SetAggregateRollupVerifier(&_Feijoapolygonrollupmanager.TransactOpts, newAggregateRollupVerifier)
}

// SetAggregateRollupVerifier is a paid mutator transaction binding the contract method 0xe2bfe8b3.
//
// Solidity: function setAggregateRollupVerifier(address newAggregateRollupVerifier) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerTransactorSession) SetAggregateRollupVerifier(newAggregateRollupVerifier common.Address) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.Contract.SetAggregateRollupVerifier(&_Feijoapolygonrollupmanager.TransactOpts, newAggregateRollupVerifier)
}

// SetMultiplierZkGasPrice is a paid mutator transaction binding the contract method 0xa1094df3.
//
// Solidity: function setMultiplierZkGasPrice(uint16 newMultiplierZkGasPrice) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerTransactor) SetMultiplierZkGasPrice(opts *bind.TransactOpts, newMultiplierZkGasPrice uint16) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.contract.Transact(opts, "setMultiplierZkGasPrice", newMultiplierZkGasPrice)
}

// SetMultiplierZkGasPrice is a paid mutator transaction binding the contract method 0xa1094df3.
//
// Solidity: function setMultiplierZkGasPrice(uint16 newMultiplierZkGasPrice) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerSession) SetMultiplierZkGasPrice(newMultiplierZkGasPrice uint16) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.Contract.SetMultiplierZkGasPrice(&_Feijoapolygonrollupmanager.TransactOpts, newMultiplierZkGasPrice)
}

// SetMultiplierZkGasPrice is a paid mutator transaction binding the contract method 0xa1094df3.
//
// Solidity: function setMultiplierZkGasPrice(uint16 newMultiplierZkGasPrice) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerTransactorSession) SetMultiplierZkGasPrice(newMultiplierZkGasPrice uint16) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.Contract.SetMultiplierZkGasPrice(&_Feijoapolygonrollupmanager.TransactOpts, newMultiplierZkGasPrice)
}

// SetPendingStateTimeout is a paid mutator transaction binding the contract method 0x9c9f3dfe.
//
// Solidity: function setPendingStateTimeout(uint64 newPendingStateTimeout) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerTransactor) SetPendingStateTimeout(opts *bind.TransactOpts, newPendingStateTimeout uint64) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.contract.Transact(opts, "setPendingStateTimeout", newPendingStateTimeout)
}

// SetPendingStateTimeout is a paid mutator transaction binding the contract method 0x9c9f3dfe.
//
// Solidity: function setPendingStateTimeout(uint64 newPendingStateTimeout) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerSession) SetPendingStateTimeout(newPendingStateTimeout uint64) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.Contract.SetPendingStateTimeout(&_Feijoapolygonrollupmanager.TransactOpts, newPendingStateTimeout)
}

// SetPendingStateTimeout is a paid mutator transaction binding the contract method 0x9c9f3dfe.
//
// Solidity: function setPendingStateTimeout(uint64 newPendingStateTimeout) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerTransactorSession) SetPendingStateTimeout(newPendingStateTimeout uint64) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.Contract.SetPendingStateTimeout(&_Feijoapolygonrollupmanager.TransactOpts, newPendingStateTimeout)
}

// SetTrustedAggregatorTimeout is a paid mutator transaction binding the contract method 0x394218e9.
//
// Solidity: function setTrustedAggregatorTimeout(uint64 newTrustedAggregatorTimeout) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerTransactor) SetTrustedAggregatorTimeout(opts *bind.TransactOpts, newTrustedAggregatorTimeout uint64) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.contract.Transact(opts, "setTrustedAggregatorTimeout", newTrustedAggregatorTimeout)
}

// SetTrustedAggregatorTimeout is a paid mutator transaction binding the contract method 0x394218e9.
//
// Solidity: function setTrustedAggregatorTimeout(uint64 newTrustedAggregatorTimeout) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerSession) SetTrustedAggregatorTimeout(newTrustedAggregatorTimeout uint64) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.Contract.SetTrustedAggregatorTimeout(&_Feijoapolygonrollupmanager.TransactOpts, newTrustedAggregatorTimeout)
}

// SetTrustedAggregatorTimeout is a paid mutator transaction binding the contract method 0x394218e9.
//
// Solidity: function setTrustedAggregatorTimeout(uint64 newTrustedAggregatorTimeout) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerTransactorSession) SetTrustedAggregatorTimeout(newTrustedAggregatorTimeout uint64) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.Contract.SetTrustedAggregatorTimeout(&_Feijoapolygonrollupmanager.TransactOpts, newTrustedAggregatorTimeout)
}

// SetVerifySequenceTimeTarget is a paid mutator transaction binding the contract method 0x8185f9d3.
//
// Solidity: function setVerifySequenceTimeTarget(uint64 newVerifySequenceTimeTarget) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerTransactor) SetVerifySequenceTimeTarget(opts *bind.TransactOpts, newVerifySequenceTimeTarget uint64) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.contract.Transact(opts, "setVerifySequenceTimeTarget", newVerifySequenceTimeTarget)
}

// SetVerifySequenceTimeTarget is a paid mutator transaction binding the contract method 0x8185f9d3.
//
// Solidity: function setVerifySequenceTimeTarget(uint64 newVerifySequenceTimeTarget) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerSession) SetVerifySequenceTimeTarget(newVerifySequenceTimeTarget uint64) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.Contract.SetVerifySequenceTimeTarget(&_Feijoapolygonrollupmanager.TransactOpts, newVerifySequenceTimeTarget)
}

// SetVerifySequenceTimeTarget is a paid mutator transaction binding the contract method 0x8185f9d3.
//
// Solidity: function setVerifySequenceTimeTarget(uint64 newVerifySequenceTimeTarget) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerTransactorSession) SetVerifySequenceTimeTarget(newVerifySequenceTimeTarget uint64) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.Contract.SetVerifySequenceTimeTarget(&_Feijoapolygonrollupmanager.TransactOpts, newVerifySequenceTimeTarget)
}

// SetZkGasPrice is a paid mutator transaction binding the contract method 0x7ec31def.
//
// Solidity: function setZkGasPrice(uint256 newZkGasPrice) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerTransactor) SetZkGasPrice(opts *bind.TransactOpts, newZkGasPrice *big.Int) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.contract.Transact(opts, "setZkGasPrice", newZkGasPrice)
}

// SetZkGasPrice is a paid mutator transaction binding the contract method 0x7ec31def.
//
// Solidity: function setZkGasPrice(uint256 newZkGasPrice) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerSession) SetZkGasPrice(newZkGasPrice *big.Int) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.Contract.SetZkGasPrice(&_Feijoapolygonrollupmanager.TransactOpts, newZkGasPrice)
}

// SetZkGasPrice is a paid mutator transaction binding the contract method 0x7ec31def.
//
// Solidity: function setZkGasPrice(uint256 newZkGasPrice) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerTransactorSession) SetZkGasPrice(newZkGasPrice *big.Int) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.Contract.SetZkGasPrice(&_Feijoapolygonrollupmanager.TransactOpts, newZkGasPrice)
}

// UpdateRollup is a paid mutator transaction binding the contract method 0xc4c928c2.
//
// Solidity: function updateRollup(address rollupContract, uint32 newRollupTypeID, bytes upgradeData) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerTransactor) UpdateRollup(opts *bind.TransactOpts, rollupContract common.Address, newRollupTypeID uint32, upgradeData []byte) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.contract.Transact(opts, "updateRollup", rollupContract, newRollupTypeID, upgradeData)
}

// UpdateRollup is a paid mutator transaction binding the contract method 0xc4c928c2.
//
// Solidity: function updateRollup(address rollupContract, uint32 newRollupTypeID, bytes upgradeData) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerSession) UpdateRollup(rollupContract common.Address, newRollupTypeID uint32, upgradeData []byte) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.Contract.UpdateRollup(&_Feijoapolygonrollupmanager.TransactOpts, rollupContract, newRollupTypeID, upgradeData)
}

// UpdateRollup is a paid mutator transaction binding the contract method 0xc4c928c2.
//
// Solidity: function updateRollup(address rollupContract, uint32 newRollupTypeID, bytes upgradeData) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerTransactorSession) UpdateRollup(rollupContract common.Address, newRollupTypeID uint32, upgradeData []byte) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.Contract.UpdateRollup(&_Feijoapolygonrollupmanager.TransactOpts, rollupContract, newRollupTypeID, upgradeData)
}

// UpdateRollupByRollupAdmin is a paid mutator transaction binding the contract method 0xdfdb8c5e.
//
// Solidity: function updateRollupByRollupAdmin(address rollupContract, uint32 newRollupTypeID) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerTransactor) UpdateRollupByRollupAdmin(opts *bind.TransactOpts, rollupContract common.Address, newRollupTypeID uint32) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.contract.Transact(opts, "updateRollupByRollupAdmin", rollupContract, newRollupTypeID)
}

// UpdateRollupByRollupAdmin is a paid mutator transaction binding the contract method 0xdfdb8c5e.
//
// Solidity: function updateRollupByRollupAdmin(address rollupContract, uint32 newRollupTypeID) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerSession) UpdateRollupByRollupAdmin(rollupContract common.Address, newRollupTypeID uint32) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.Contract.UpdateRollupByRollupAdmin(&_Feijoapolygonrollupmanager.TransactOpts, rollupContract, newRollupTypeID)
}

// UpdateRollupByRollupAdmin is a paid mutator transaction binding the contract method 0xdfdb8c5e.
//
// Solidity: function updateRollupByRollupAdmin(address rollupContract, uint32 newRollupTypeID) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerTransactorSession) UpdateRollupByRollupAdmin(rollupContract common.Address, newRollupTypeID uint32) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.Contract.UpdateRollupByRollupAdmin(&_Feijoapolygonrollupmanager.TransactOpts, rollupContract, newRollupTypeID)
}

// VerifySequencesMultiProof is a paid mutator transaction binding the contract method 0xf00bdaa4.
//
// Solidity: function verifySequencesMultiProof((uint32,uint64,uint64,uint64,bytes32,bytes32)[] verifySequencesData, address beneficiary, bytes32[24] proof) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerTransactor) VerifySequencesMultiProof(opts *bind.TransactOpts, verifySequencesData []PolygonRollupManagerVerifySequenceData, beneficiary common.Address, proof [24][32]byte) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.contract.Transact(opts, "verifySequencesMultiProof", verifySequencesData, beneficiary, proof)
}

// VerifySequencesMultiProof is a paid mutator transaction binding the contract method 0xf00bdaa4.
//
// Solidity: function verifySequencesMultiProof((uint32,uint64,uint64,uint64,bytes32,bytes32)[] verifySequencesData, address beneficiary, bytes32[24] proof) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerSession) VerifySequencesMultiProof(verifySequencesData []PolygonRollupManagerVerifySequenceData, beneficiary common.Address, proof [24][32]byte) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.Contract.VerifySequencesMultiProof(&_Feijoapolygonrollupmanager.TransactOpts, verifySequencesData, beneficiary, proof)
}

// VerifySequencesMultiProof is a paid mutator transaction binding the contract method 0xf00bdaa4.
//
// Solidity: function verifySequencesMultiProof((uint32,uint64,uint64,uint64,bytes32,bytes32)[] verifySequencesData, address beneficiary, bytes32[24] proof) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerTransactorSession) VerifySequencesMultiProof(verifySequencesData []PolygonRollupManagerVerifySequenceData, beneficiary common.Address, proof [24][32]byte) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.Contract.VerifySequencesMultiProof(&_Feijoapolygonrollupmanager.TransactOpts, verifySequencesData, beneficiary, proof)
}

// VerifySequencesTrustedAggregatorMultiProof is a paid mutator transaction binding the contract method 0xde794850.
//
// Solidity: function verifySequencesTrustedAggregatorMultiProof((uint32,uint64,uint64,uint64,bytes32,bytes32)[] verifySequencesData, address beneficiary, bytes32[24] proof) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerTransactor) VerifySequencesTrustedAggregatorMultiProof(opts *bind.TransactOpts, verifySequencesData []PolygonRollupManagerVerifySequenceData, beneficiary common.Address, proof [24][32]byte) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.contract.Transact(opts, "verifySequencesTrustedAggregatorMultiProof", verifySequencesData, beneficiary, proof)
}

// VerifySequencesTrustedAggregatorMultiProof is a paid mutator transaction binding the contract method 0xde794850.
//
// Solidity: function verifySequencesTrustedAggregatorMultiProof((uint32,uint64,uint64,uint64,bytes32,bytes32)[] verifySequencesData, address beneficiary, bytes32[24] proof) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerSession) VerifySequencesTrustedAggregatorMultiProof(verifySequencesData []PolygonRollupManagerVerifySequenceData, beneficiary common.Address, proof [24][32]byte) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.Contract.VerifySequencesTrustedAggregatorMultiProof(&_Feijoapolygonrollupmanager.TransactOpts, verifySequencesData, beneficiary, proof)
}

// VerifySequencesTrustedAggregatorMultiProof is a paid mutator transaction binding the contract method 0xde794850.
//
// Solidity: function verifySequencesTrustedAggregatorMultiProof((uint32,uint64,uint64,uint64,bytes32,bytes32)[] verifySequencesData, address beneficiary, bytes32[24] proof) returns()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerTransactorSession) VerifySequencesTrustedAggregatorMultiProof(verifySequencesData []PolygonRollupManagerVerifySequenceData, beneficiary common.Address, proof [24][32]byte) (*types.Transaction, error) {
	return _Feijoapolygonrollupmanager.Contract.VerifySequencesTrustedAggregatorMultiProof(&_Feijoapolygonrollupmanager.TransactOpts, verifySequencesData, beneficiary, proof)
}

// FeijoapolygonrollupmanagerAddExistingRollupIterator is returned from FilterAddExistingRollup and is used to iterate over the raw logs and unpacked data for AddExistingRollup events raised by the Feijoapolygonrollupmanager contract.
type FeijoapolygonrollupmanagerAddExistingRollupIterator struct {
	Event *FeijoapolygonrollupmanagerAddExistingRollup // Event containing the contract specifics and raw log

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
func (it *FeijoapolygonrollupmanagerAddExistingRollupIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeijoapolygonrollupmanagerAddExistingRollup)
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
		it.Event = new(FeijoapolygonrollupmanagerAddExistingRollup)
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
func (it *FeijoapolygonrollupmanagerAddExistingRollupIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeijoapolygonrollupmanagerAddExistingRollupIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeijoapolygonrollupmanagerAddExistingRollup represents a AddExistingRollup event raised by the Feijoapolygonrollupmanager contract.
type FeijoapolygonrollupmanagerAddExistingRollup struct {
	RollupID                          uint32
	ForkID                            uint64
	RollupAddress                     common.Address
	ChainID                           uint64
	RollupCompatibilityID             uint8
	LastVerifiedSequenceBeforeUpgrade uint64
	Raw                               types.Log // Blockchain specific contextual infos
}

// FilterAddExistingRollup is a free log retrieval operation binding the contract event 0xadfc7d56f7e39b08b321534f14bfb135ad27698f7d2f5ad0edc2356ea9a3f850.
//
// Solidity: event AddExistingRollup(uint32 indexed rollupID, uint64 forkID, address rollupAddress, uint64 chainID, uint8 rollupCompatibilityID, uint64 lastVerifiedSequenceBeforeUpgrade)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) FilterAddExistingRollup(opts *bind.FilterOpts, rollupID []uint32) (*FeijoapolygonrollupmanagerAddExistingRollupIterator, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	logs, sub, err := _Feijoapolygonrollupmanager.contract.FilterLogs(opts, "AddExistingRollup", rollupIDRule)
	if err != nil {
		return nil, err
	}
	return &FeijoapolygonrollupmanagerAddExistingRollupIterator{contract: _Feijoapolygonrollupmanager.contract, event: "AddExistingRollup", logs: logs, sub: sub}, nil
}

// WatchAddExistingRollup is a free log subscription operation binding the contract event 0xadfc7d56f7e39b08b321534f14bfb135ad27698f7d2f5ad0edc2356ea9a3f850.
//
// Solidity: event AddExistingRollup(uint32 indexed rollupID, uint64 forkID, address rollupAddress, uint64 chainID, uint8 rollupCompatibilityID, uint64 lastVerifiedSequenceBeforeUpgrade)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) WatchAddExistingRollup(opts *bind.WatchOpts, sink chan<- *FeijoapolygonrollupmanagerAddExistingRollup, rollupID []uint32) (event.Subscription, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	logs, sub, err := _Feijoapolygonrollupmanager.contract.WatchLogs(opts, "AddExistingRollup", rollupIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeijoapolygonrollupmanagerAddExistingRollup)
				if err := _Feijoapolygonrollupmanager.contract.UnpackLog(event, "AddExistingRollup", log); err != nil {
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

// ParseAddExistingRollup is a log parse operation binding the contract event 0xadfc7d56f7e39b08b321534f14bfb135ad27698f7d2f5ad0edc2356ea9a3f850.
//
// Solidity: event AddExistingRollup(uint32 indexed rollupID, uint64 forkID, address rollupAddress, uint64 chainID, uint8 rollupCompatibilityID, uint64 lastVerifiedSequenceBeforeUpgrade)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) ParseAddExistingRollup(log types.Log) (*FeijoapolygonrollupmanagerAddExistingRollup, error) {
	event := new(FeijoapolygonrollupmanagerAddExistingRollup)
	if err := _Feijoapolygonrollupmanager.contract.UnpackLog(event, "AddExistingRollup", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeijoapolygonrollupmanagerAddNewRollupTypeIterator is returned from FilterAddNewRollupType and is used to iterate over the raw logs and unpacked data for AddNewRollupType events raised by the Feijoapolygonrollupmanager contract.
type FeijoapolygonrollupmanagerAddNewRollupTypeIterator struct {
	Event *FeijoapolygonrollupmanagerAddNewRollupType // Event containing the contract specifics and raw log

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
func (it *FeijoapolygonrollupmanagerAddNewRollupTypeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeijoapolygonrollupmanagerAddNewRollupType)
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
		it.Event = new(FeijoapolygonrollupmanagerAddNewRollupType)
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
func (it *FeijoapolygonrollupmanagerAddNewRollupTypeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeijoapolygonrollupmanagerAddNewRollupTypeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeijoapolygonrollupmanagerAddNewRollupType represents a AddNewRollupType event raised by the Feijoapolygonrollupmanager contract.
type FeijoapolygonrollupmanagerAddNewRollupType struct {
	RollupTypeID            uint32
	ConsensusImplementation common.Address
	Verifier                common.Address
	ForkID                  uint64
	RollupCompatibilityID   uint8
	Genesis                 [32]byte
	Description             string
	Raw                     types.Log // Blockchain specific contextual infos
}

// FilterAddNewRollupType is a free log retrieval operation binding the contract event 0xa2970448b3bd66ba7e524e7b2a5b9cf94fa29e32488fb942afdfe70dd4b77b52.
//
// Solidity: event AddNewRollupType(uint32 indexed rollupTypeID, address consensusImplementation, address verifier, uint64 forkID, uint8 rollupCompatibilityID, bytes32 genesis, string description)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) FilterAddNewRollupType(opts *bind.FilterOpts, rollupTypeID []uint32) (*FeijoapolygonrollupmanagerAddNewRollupTypeIterator, error) {

	var rollupTypeIDRule []interface{}
	for _, rollupTypeIDItem := range rollupTypeID {
		rollupTypeIDRule = append(rollupTypeIDRule, rollupTypeIDItem)
	}

	logs, sub, err := _Feijoapolygonrollupmanager.contract.FilterLogs(opts, "AddNewRollupType", rollupTypeIDRule)
	if err != nil {
		return nil, err
	}
	return &FeijoapolygonrollupmanagerAddNewRollupTypeIterator{contract: _Feijoapolygonrollupmanager.contract, event: "AddNewRollupType", logs: logs, sub: sub}, nil
}

// WatchAddNewRollupType is a free log subscription operation binding the contract event 0xa2970448b3bd66ba7e524e7b2a5b9cf94fa29e32488fb942afdfe70dd4b77b52.
//
// Solidity: event AddNewRollupType(uint32 indexed rollupTypeID, address consensusImplementation, address verifier, uint64 forkID, uint8 rollupCompatibilityID, bytes32 genesis, string description)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) WatchAddNewRollupType(opts *bind.WatchOpts, sink chan<- *FeijoapolygonrollupmanagerAddNewRollupType, rollupTypeID []uint32) (event.Subscription, error) {

	var rollupTypeIDRule []interface{}
	for _, rollupTypeIDItem := range rollupTypeID {
		rollupTypeIDRule = append(rollupTypeIDRule, rollupTypeIDItem)
	}

	logs, sub, err := _Feijoapolygonrollupmanager.contract.WatchLogs(opts, "AddNewRollupType", rollupTypeIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeijoapolygonrollupmanagerAddNewRollupType)
				if err := _Feijoapolygonrollupmanager.contract.UnpackLog(event, "AddNewRollupType", log); err != nil {
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

// ParseAddNewRollupType is a log parse operation binding the contract event 0xa2970448b3bd66ba7e524e7b2a5b9cf94fa29e32488fb942afdfe70dd4b77b52.
//
// Solidity: event AddNewRollupType(uint32 indexed rollupTypeID, address consensusImplementation, address verifier, uint64 forkID, uint8 rollupCompatibilityID, bytes32 genesis, string description)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) ParseAddNewRollupType(log types.Log) (*FeijoapolygonrollupmanagerAddNewRollupType, error) {
	event := new(FeijoapolygonrollupmanagerAddNewRollupType)
	if err := _Feijoapolygonrollupmanager.contract.UnpackLog(event, "AddNewRollupType", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeijoapolygonrollupmanagerConsolidatePendingStateIterator is returned from FilterConsolidatePendingState and is used to iterate over the raw logs and unpacked data for ConsolidatePendingState events raised by the Feijoapolygonrollupmanager contract.
type FeijoapolygonrollupmanagerConsolidatePendingStateIterator struct {
	Event *FeijoapolygonrollupmanagerConsolidatePendingState // Event containing the contract specifics and raw log

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
func (it *FeijoapolygonrollupmanagerConsolidatePendingStateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeijoapolygonrollupmanagerConsolidatePendingState)
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
		it.Event = new(FeijoapolygonrollupmanagerConsolidatePendingState)
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
func (it *FeijoapolygonrollupmanagerConsolidatePendingStateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeijoapolygonrollupmanagerConsolidatePendingStateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeijoapolygonrollupmanagerConsolidatePendingState represents a ConsolidatePendingState event raised by the Feijoapolygonrollupmanager contract.
type FeijoapolygonrollupmanagerConsolidatePendingState struct {
	RollupID        uint32
	NumSequence     uint64
	StateRoot       [32]byte
	ExitRoot        [32]byte
	PendingStateNum uint64
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterConsolidatePendingState is a free log retrieval operation binding the contract event 0x581910eb7a27738945c2f00a91f2284b2d6de9d4e472b12f901c2b0df045e21b.
//
// Solidity: event ConsolidatePendingState(uint32 indexed rollupID, uint64 numSequence, bytes32 stateRoot, bytes32 exitRoot, uint64 pendingStateNum)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) FilterConsolidatePendingState(opts *bind.FilterOpts, rollupID []uint32) (*FeijoapolygonrollupmanagerConsolidatePendingStateIterator, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	logs, sub, err := _Feijoapolygonrollupmanager.contract.FilterLogs(opts, "ConsolidatePendingState", rollupIDRule)
	if err != nil {
		return nil, err
	}
	return &FeijoapolygonrollupmanagerConsolidatePendingStateIterator{contract: _Feijoapolygonrollupmanager.contract, event: "ConsolidatePendingState", logs: logs, sub: sub}, nil
}

// WatchConsolidatePendingState is a free log subscription operation binding the contract event 0x581910eb7a27738945c2f00a91f2284b2d6de9d4e472b12f901c2b0df045e21b.
//
// Solidity: event ConsolidatePendingState(uint32 indexed rollupID, uint64 numSequence, bytes32 stateRoot, bytes32 exitRoot, uint64 pendingStateNum)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) WatchConsolidatePendingState(opts *bind.WatchOpts, sink chan<- *FeijoapolygonrollupmanagerConsolidatePendingState, rollupID []uint32) (event.Subscription, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	logs, sub, err := _Feijoapolygonrollupmanager.contract.WatchLogs(opts, "ConsolidatePendingState", rollupIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeijoapolygonrollupmanagerConsolidatePendingState)
				if err := _Feijoapolygonrollupmanager.contract.UnpackLog(event, "ConsolidatePendingState", log); err != nil {
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

// ParseConsolidatePendingState is a log parse operation binding the contract event 0x581910eb7a27738945c2f00a91f2284b2d6de9d4e472b12f901c2b0df045e21b.
//
// Solidity: event ConsolidatePendingState(uint32 indexed rollupID, uint64 numSequence, bytes32 stateRoot, bytes32 exitRoot, uint64 pendingStateNum)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) ParseConsolidatePendingState(log types.Log) (*FeijoapolygonrollupmanagerConsolidatePendingState, error) {
	event := new(FeijoapolygonrollupmanagerConsolidatePendingState)
	if err := _Feijoapolygonrollupmanager.contract.UnpackLog(event, "ConsolidatePendingState", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeijoapolygonrollupmanagerCreateNewRollupIterator is returned from FilterCreateNewRollup and is used to iterate over the raw logs and unpacked data for CreateNewRollup events raised by the Feijoapolygonrollupmanager contract.
type FeijoapolygonrollupmanagerCreateNewRollupIterator struct {
	Event *FeijoapolygonrollupmanagerCreateNewRollup // Event containing the contract specifics and raw log

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
func (it *FeijoapolygonrollupmanagerCreateNewRollupIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeijoapolygonrollupmanagerCreateNewRollup)
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
		it.Event = new(FeijoapolygonrollupmanagerCreateNewRollup)
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
func (it *FeijoapolygonrollupmanagerCreateNewRollupIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeijoapolygonrollupmanagerCreateNewRollupIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeijoapolygonrollupmanagerCreateNewRollup represents a CreateNewRollup event raised by the Feijoapolygonrollupmanager contract.
type FeijoapolygonrollupmanagerCreateNewRollup struct {
	RollupID        uint32
	RollupTypeID    uint32
	RollupAddress   common.Address
	ChainID         uint64
	GasTokenAddress common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterCreateNewRollup is a free log retrieval operation binding the contract event 0x194c983456df6701c6a50830b90fe80e72b823411d0d524970c9590dc277a641.
//
// Solidity: event CreateNewRollup(uint32 indexed rollupID, uint32 rollupTypeID, address rollupAddress, uint64 chainID, address gasTokenAddress)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) FilterCreateNewRollup(opts *bind.FilterOpts, rollupID []uint32) (*FeijoapolygonrollupmanagerCreateNewRollupIterator, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	logs, sub, err := _Feijoapolygonrollupmanager.contract.FilterLogs(opts, "CreateNewRollup", rollupIDRule)
	if err != nil {
		return nil, err
	}
	return &FeijoapolygonrollupmanagerCreateNewRollupIterator{contract: _Feijoapolygonrollupmanager.contract, event: "CreateNewRollup", logs: logs, sub: sub}, nil
}

// WatchCreateNewRollup is a free log subscription operation binding the contract event 0x194c983456df6701c6a50830b90fe80e72b823411d0d524970c9590dc277a641.
//
// Solidity: event CreateNewRollup(uint32 indexed rollupID, uint32 rollupTypeID, address rollupAddress, uint64 chainID, address gasTokenAddress)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) WatchCreateNewRollup(opts *bind.WatchOpts, sink chan<- *FeijoapolygonrollupmanagerCreateNewRollup, rollupID []uint32) (event.Subscription, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	logs, sub, err := _Feijoapolygonrollupmanager.contract.WatchLogs(opts, "CreateNewRollup", rollupIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeijoapolygonrollupmanagerCreateNewRollup)
				if err := _Feijoapolygonrollupmanager.contract.UnpackLog(event, "CreateNewRollup", log); err != nil {
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

// ParseCreateNewRollup is a log parse operation binding the contract event 0x194c983456df6701c6a50830b90fe80e72b823411d0d524970c9590dc277a641.
//
// Solidity: event CreateNewRollup(uint32 indexed rollupID, uint32 rollupTypeID, address rollupAddress, uint64 chainID, address gasTokenAddress)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) ParseCreateNewRollup(log types.Log) (*FeijoapolygonrollupmanagerCreateNewRollup, error) {
	event := new(FeijoapolygonrollupmanagerCreateNewRollup)
	if err := _Feijoapolygonrollupmanager.contract.UnpackLog(event, "CreateNewRollup", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeijoapolygonrollupmanagerEmergencyStateActivatedIterator is returned from FilterEmergencyStateActivated and is used to iterate over the raw logs and unpacked data for EmergencyStateActivated events raised by the Feijoapolygonrollupmanager contract.
type FeijoapolygonrollupmanagerEmergencyStateActivatedIterator struct {
	Event *FeijoapolygonrollupmanagerEmergencyStateActivated // Event containing the contract specifics and raw log

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
func (it *FeijoapolygonrollupmanagerEmergencyStateActivatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeijoapolygonrollupmanagerEmergencyStateActivated)
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
		it.Event = new(FeijoapolygonrollupmanagerEmergencyStateActivated)
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
func (it *FeijoapolygonrollupmanagerEmergencyStateActivatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeijoapolygonrollupmanagerEmergencyStateActivatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeijoapolygonrollupmanagerEmergencyStateActivated represents a EmergencyStateActivated event raised by the Feijoapolygonrollupmanager contract.
type FeijoapolygonrollupmanagerEmergencyStateActivated struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterEmergencyStateActivated is a free log retrieval operation binding the contract event 0x2261efe5aef6fedc1fd1550b25facc9181745623049c7901287030b9ad1a5497.
//
// Solidity: event EmergencyStateActivated()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) FilterEmergencyStateActivated(opts *bind.FilterOpts) (*FeijoapolygonrollupmanagerEmergencyStateActivatedIterator, error) {

	logs, sub, err := _Feijoapolygonrollupmanager.contract.FilterLogs(opts, "EmergencyStateActivated")
	if err != nil {
		return nil, err
	}
	return &FeijoapolygonrollupmanagerEmergencyStateActivatedIterator{contract: _Feijoapolygonrollupmanager.contract, event: "EmergencyStateActivated", logs: logs, sub: sub}, nil
}

// WatchEmergencyStateActivated is a free log subscription operation binding the contract event 0x2261efe5aef6fedc1fd1550b25facc9181745623049c7901287030b9ad1a5497.
//
// Solidity: event EmergencyStateActivated()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) WatchEmergencyStateActivated(opts *bind.WatchOpts, sink chan<- *FeijoapolygonrollupmanagerEmergencyStateActivated) (event.Subscription, error) {

	logs, sub, err := _Feijoapolygonrollupmanager.contract.WatchLogs(opts, "EmergencyStateActivated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeijoapolygonrollupmanagerEmergencyStateActivated)
				if err := _Feijoapolygonrollupmanager.contract.UnpackLog(event, "EmergencyStateActivated", log); err != nil {
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

// ParseEmergencyStateActivated is a log parse operation binding the contract event 0x2261efe5aef6fedc1fd1550b25facc9181745623049c7901287030b9ad1a5497.
//
// Solidity: event EmergencyStateActivated()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) ParseEmergencyStateActivated(log types.Log) (*FeijoapolygonrollupmanagerEmergencyStateActivated, error) {
	event := new(FeijoapolygonrollupmanagerEmergencyStateActivated)
	if err := _Feijoapolygonrollupmanager.contract.UnpackLog(event, "EmergencyStateActivated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeijoapolygonrollupmanagerEmergencyStateDeactivatedIterator is returned from FilterEmergencyStateDeactivated and is used to iterate over the raw logs and unpacked data for EmergencyStateDeactivated events raised by the Feijoapolygonrollupmanager contract.
type FeijoapolygonrollupmanagerEmergencyStateDeactivatedIterator struct {
	Event *FeijoapolygonrollupmanagerEmergencyStateDeactivated // Event containing the contract specifics and raw log

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
func (it *FeijoapolygonrollupmanagerEmergencyStateDeactivatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeijoapolygonrollupmanagerEmergencyStateDeactivated)
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
		it.Event = new(FeijoapolygonrollupmanagerEmergencyStateDeactivated)
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
func (it *FeijoapolygonrollupmanagerEmergencyStateDeactivatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeijoapolygonrollupmanagerEmergencyStateDeactivatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeijoapolygonrollupmanagerEmergencyStateDeactivated represents a EmergencyStateDeactivated event raised by the Feijoapolygonrollupmanager contract.
type FeijoapolygonrollupmanagerEmergencyStateDeactivated struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterEmergencyStateDeactivated is a free log retrieval operation binding the contract event 0x1e5e34eea33501aecf2ebec9fe0e884a40804275ea7fe10b2ba084c8374308b3.
//
// Solidity: event EmergencyStateDeactivated()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) FilterEmergencyStateDeactivated(opts *bind.FilterOpts) (*FeijoapolygonrollupmanagerEmergencyStateDeactivatedIterator, error) {

	logs, sub, err := _Feijoapolygonrollupmanager.contract.FilterLogs(opts, "EmergencyStateDeactivated")
	if err != nil {
		return nil, err
	}
	return &FeijoapolygonrollupmanagerEmergencyStateDeactivatedIterator{contract: _Feijoapolygonrollupmanager.contract, event: "EmergencyStateDeactivated", logs: logs, sub: sub}, nil
}

// WatchEmergencyStateDeactivated is a free log subscription operation binding the contract event 0x1e5e34eea33501aecf2ebec9fe0e884a40804275ea7fe10b2ba084c8374308b3.
//
// Solidity: event EmergencyStateDeactivated()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) WatchEmergencyStateDeactivated(opts *bind.WatchOpts, sink chan<- *FeijoapolygonrollupmanagerEmergencyStateDeactivated) (event.Subscription, error) {

	logs, sub, err := _Feijoapolygonrollupmanager.contract.WatchLogs(opts, "EmergencyStateDeactivated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeijoapolygonrollupmanagerEmergencyStateDeactivated)
				if err := _Feijoapolygonrollupmanager.contract.UnpackLog(event, "EmergencyStateDeactivated", log); err != nil {
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

// ParseEmergencyStateDeactivated is a log parse operation binding the contract event 0x1e5e34eea33501aecf2ebec9fe0e884a40804275ea7fe10b2ba084c8374308b3.
//
// Solidity: event EmergencyStateDeactivated()
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) ParseEmergencyStateDeactivated(log types.Log) (*FeijoapolygonrollupmanagerEmergencyStateDeactivated, error) {
	event := new(FeijoapolygonrollupmanagerEmergencyStateDeactivated)
	if err := _Feijoapolygonrollupmanager.contract.UnpackLog(event, "EmergencyStateDeactivated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeijoapolygonrollupmanagerInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Feijoapolygonrollupmanager contract.
type FeijoapolygonrollupmanagerInitializedIterator struct {
	Event *FeijoapolygonrollupmanagerInitialized // Event containing the contract specifics and raw log

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
func (it *FeijoapolygonrollupmanagerInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeijoapolygonrollupmanagerInitialized)
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
		it.Event = new(FeijoapolygonrollupmanagerInitialized)
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
func (it *FeijoapolygonrollupmanagerInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeijoapolygonrollupmanagerInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeijoapolygonrollupmanagerInitialized represents a Initialized event raised by the Feijoapolygonrollupmanager contract.
type FeijoapolygonrollupmanagerInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) FilterInitialized(opts *bind.FilterOpts) (*FeijoapolygonrollupmanagerInitializedIterator, error) {

	logs, sub, err := _Feijoapolygonrollupmanager.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &FeijoapolygonrollupmanagerInitializedIterator{contract: _Feijoapolygonrollupmanager.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *FeijoapolygonrollupmanagerInitialized) (event.Subscription, error) {

	logs, sub, err := _Feijoapolygonrollupmanager.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeijoapolygonrollupmanagerInitialized)
				if err := _Feijoapolygonrollupmanager.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) ParseInitialized(log types.Log) (*FeijoapolygonrollupmanagerInitialized, error) {
	event := new(FeijoapolygonrollupmanagerInitialized)
	if err := _Feijoapolygonrollupmanager.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeijoapolygonrollupmanagerObsoleteRollupTypeIterator is returned from FilterObsoleteRollupType and is used to iterate over the raw logs and unpacked data for ObsoleteRollupType events raised by the Feijoapolygonrollupmanager contract.
type FeijoapolygonrollupmanagerObsoleteRollupTypeIterator struct {
	Event *FeijoapolygonrollupmanagerObsoleteRollupType // Event containing the contract specifics and raw log

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
func (it *FeijoapolygonrollupmanagerObsoleteRollupTypeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeijoapolygonrollupmanagerObsoleteRollupType)
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
		it.Event = new(FeijoapolygonrollupmanagerObsoleteRollupType)
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
func (it *FeijoapolygonrollupmanagerObsoleteRollupTypeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeijoapolygonrollupmanagerObsoleteRollupTypeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeijoapolygonrollupmanagerObsoleteRollupType represents a ObsoleteRollupType event raised by the Feijoapolygonrollupmanager contract.
type FeijoapolygonrollupmanagerObsoleteRollupType struct {
	RollupTypeID uint32
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterObsoleteRollupType is a free log retrieval operation binding the contract event 0x4710d2ee567ef1ed6eb2f651dde4589524bcf7cebc62147a99b281cc836e7e44.
//
// Solidity: event ObsoleteRollupType(uint32 indexed rollupTypeID)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) FilterObsoleteRollupType(opts *bind.FilterOpts, rollupTypeID []uint32) (*FeijoapolygonrollupmanagerObsoleteRollupTypeIterator, error) {

	var rollupTypeIDRule []interface{}
	for _, rollupTypeIDItem := range rollupTypeID {
		rollupTypeIDRule = append(rollupTypeIDRule, rollupTypeIDItem)
	}

	logs, sub, err := _Feijoapolygonrollupmanager.contract.FilterLogs(opts, "ObsoleteRollupType", rollupTypeIDRule)
	if err != nil {
		return nil, err
	}
	return &FeijoapolygonrollupmanagerObsoleteRollupTypeIterator{contract: _Feijoapolygonrollupmanager.contract, event: "ObsoleteRollupType", logs: logs, sub: sub}, nil
}

// WatchObsoleteRollupType is a free log subscription operation binding the contract event 0x4710d2ee567ef1ed6eb2f651dde4589524bcf7cebc62147a99b281cc836e7e44.
//
// Solidity: event ObsoleteRollupType(uint32 indexed rollupTypeID)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) WatchObsoleteRollupType(opts *bind.WatchOpts, sink chan<- *FeijoapolygonrollupmanagerObsoleteRollupType, rollupTypeID []uint32) (event.Subscription, error) {

	var rollupTypeIDRule []interface{}
	for _, rollupTypeIDItem := range rollupTypeID {
		rollupTypeIDRule = append(rollupTypeIDRule, rollupTypeIDItem)
	}

	logs, sub, err := _Feijoapolygonrollupmanager.contract.WatchLogs(opts, "ObsoleteRollupType", rollupTypeIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeijoapolygonrollupmanagerObsoleteRollupType)
				if err := _Feijoapolygonrollupmanager.contract.UnpackLog(event, "ObsoleteRollupType", log); err != nil {
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

// ParseObsoleteRollupType is a log parse operation binding the contract event 0x4710d2ee567ef1ed6eb2f651dde4589524bcf7cebc62147a99b281cc836e7e44.
//
// Solidity: event ObsoleteRollupType(uint32 indexed rollupTypeID)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) ParseObsoleteRollupType(log types.Log) (*FeijoapolygonrollupmanagerObsoleteRollupType, error) {
	event := new(FeijoapolygonrollupmanagerObsoleteRollupType)
	if err := _Feijoapolygonrollupmanager.contract.UnpackLog(event, "ObsoleteRollupType", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeijoapolygonrollupmanagerOnSequenceIterator is returned from FilterOnSequence and is used to iterate over the raw logs and unpacked data for OnSequence events raised by the Feijoapolygonrollupmanager contract.
type FeijoapolygonrollupmanagerOnSequenceIterator struct {
	Event *FeijoapolygonrollupmanagerOnSequence // Event containing the contract specifics and raw log

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
func (it *FeijoapolygonrollupmanagerOnSequenceIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeijoapolygonrollupmanagerOnSequence)
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
		it.Event = new(FeijoapolygonrollupmanagerOnSequence)
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
func (it *FeijoapolygonrollupmanagerOnSequenceIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeijoapolygonrollupmanagerOnSequenceIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeijoapolygonrollupmanagerOnSequence represents a OnSequence event raised by the Feijoapolygonrollupmanager contract.
type FeijoapolygonrollupmanagerOnSequence struct {
	RollupID       uint32
	ZkGasLimit     *big.Int
	BlobsSequenced uint64
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterOnSequence is a free log retrieval operation binding the contract event 0xd3104eaeb2b51fc52b7d354a19bf146d10ed8d047b43764be8f78cbb3ffd8be4.
//
// Solidity: event OnSequence(uint32 indexed rollupID, uint128 zkGasLimit, uint64 blobsSequenced)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) FilterOnSequence(opts *bind.FilterOpts, rollupID []uint32) (*FeijoapolygonrollupmanagerOnSequenceIterator, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	logs, sub, err := _Feijoapolygonrollupmanager.contract.FilterLogs(opts, "OnSequence", rollupIDRule)
	if err != nil {
		return nil, err
	}
	return &FeijoapolygonrollupmanagerOnSequenceIterator{contract: _Feijoapolygonrollupmanager.contract, event: "OnSequence", logs: logs, sub: sub}, nil
}

// WatchOnSequence is a free log subscription operation binding the contract event 0xd3104eaeb2b51fc52b7d354a19bf146d10ed8d047b43764be8f78cbb3ffd8be4.
//
// Solidity: event OnSequence(uint32 indexed rollupID, uint128 zkGasLimit, uint64 blobsSequenced)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) WatchOnSequence(opts *bind.WatchOpts, sink chan<- *FeijoapolygonrollupmanagerOnSequence, rollupID []uint32) (event.Subscription, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	logs, sub, err := _Feijoapolygonrollupmanager.contract.WatchLogs(opts, "OnSequence", rollupIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeijoapolygonrollupmanagerOnSequence)
				if err := _Feijoapolygonrollupmanager.contract.UnpackLog(event, "OnSequence", log); err != nil {
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

// ParseOnSequence is a log parse operation binding the contract event 0xd3104eaeb2b51fc52b7d354a19bf146d10ed8d047b43764be8f78cbb3ffd8be4.
//
// Solidity: event OnSequence(uint32 indexed rollupID, uint128 zkGasLimit, uint64 blobsSequenced)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) ParseOnSequence(log types.Log) (*FeijoapolygonrollupmanagerOnSequence, error) {
	event := new(FeijoapolygonrollupmanagerOnSequence)
	if err := _Feijoapolygonrollupmanager.contract.UnpackLog(event, "OnSequence", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeijoapolygonrollupmanagerOverridePendingStateIterator is returned from FilterOverridePendingState and is used to iterate over the raw logs and unpacked data for OverridePendingState events raised by the Feijoapolygonrollupmanager contract.
type FeijoapolygonrollupmanagerOverridePendingStateIterator struct {
	Event *FeijoapolygonrollupmanagerOverridePendingState // Event containing the contract specifics and raw log

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
func (it *FeijoapolygonrollupmanagerOverridePendingStateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeijoapolygonrollupmanagerOverridePendingState)
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
		it.Event = new(FeijoapolygonrollupmanagerOverridePendingState)
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
func (it *FeijoapolygonrollupmanagerOverridePendingStateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeijoapolygonrollupmanagerOverridePendingStateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeijoapolygonrollupmanagerOverridePendingState represents a OverridePendingState event raised by the Feijoapolygonrollupmanager contract.
type FeijoapolygonrollupmanagerOverridePendingState struct {
	RollupID    uint32
	NumSequence uint64
	StateRoot   [32]byte
	ExitRoot    [32]byte
	Aggregator  common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterOverridePendingState is a free log retrieval operation binding the contract event 0x3182bd6e6f74fc1fdc88b60f3a4f4c7f79db6ae6f5b88a1b3f5a1e28ec210d5e.
//
// Solidity: event OverridePendingState(uint32 indexed rollupID, uint64 numSequence, bytes32 stateRoot, bytes32 exitRoot, address aggregator)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) FilterOverridePendingState(opts *bind.FilterOpts, rollupID []uint32) (*FeijoapolygonrollupmanagerOverridePendingStateIterator, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	logs, sub, err := _Feijoapolygonrollupmanager.contract.FilterLogs(opts, "OverridePendingState", rollupIDRule)
	if err != nil {
		return nil, err
	}
	return &FeijoapolygonrollupmanagerOverridePendingStateIterator{contract: _Feijoapolygonrollupmanager.contract, event: "OverridePendingState", logs: logs, sub: sub}, nil
}

// WatchOverridePendingState is a free log subscription operation binding the contract event 0x3182bd6e6f74fc1fdc88b60f3a4f4c7f79db6ae6f5b88a1b3f5a1e28ec210d5e.
//
// Solidity: event OverridePendingState(uint32 indexed rollupID, uint64 numSequence, bytes32 stateRoot, bytes32 exitRoot, address aggregator)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) WatchOverridePendingState(opts *bind.WatchOpts, sink chan<- *FeijoapolygonrollupmanagerOverridePendingState, rollupID []uint32) (event.Subscription, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	logs, sub, err := _Feijoapolygonrollupmanager.contract.WatchLogs(opts, "OverridePendingState", rollupIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeijoapolygonrollupmanagerOverridePendingState)
				if err := _Feijoapolygonrollupmanager.contract.UnpackLog(event, "OverridePendingState", log); err != nil {
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

// ParseOverridePendingState is a log parse operation binding the contract event 0x3182bd6e6f74fc1fdc88b60f3a4f4c7f79db6ae6f5b88a1b3f5a1e28ec210d5e.
//
// Solidity: event OverridePendingState(uint32 indexed rollupID, uint64 numSequence, bytes32 stateRoot, bytes32 exitRoot, address aggregator)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) ParseOverridePendingState(log types.Log) (*FeijoapolygonrollupmanagerOverridePendingState, error) {
	event := new(FeijoapolygonrollupmanagerOverridePendingState)
	if err := _Feijoapolygonrollupmanager.contract.UnpackLog(event, "OverridePendingState", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeijoapolygonrollupmanagerProveNonDeterministicPendingStateIterator is returned from FilterProveNonDeterministicPendingState and is used to iterate over the raw logs and unpacked data for ProveNonDeterministicPendingState events raised by the Feijoapolygonrollupmanager contract.
type FeijoapolygonrollupmanagerProveNonDeterministicPendingStateIterator struct {
	Event *FeijoapolygonrollupmanagerProveNonDeterministicPendingState // Event containing the contract specifics and raw log

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
func (it *FeijoapolygonrollupmanagerProveNonDeterministicPendingStateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeijoapolygonrollupmanagerProveNonDeterministicPendingState)
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
		it.Event = new(FeijoapolygonrollupmanagerProveNonDeterministicPendingState)
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
func (it *FeijoapolygonrollupmanagerProveNonDeterministicPendingStateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeijoapolygonrollupmanagerProveNonDeterministicPendingStateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeijoapolygonrollupmanagerProveNonDeterministicPendingState represents a ProveNonDeterministicPendingState event raised by the Feijoapolygonrollupmanager contract.
type FeijoapolygonrollupmanagerProveNonDeterministicPendingState struct {
	StoredStateRoot [32]byte
	ProvedStateRoot [32]byte
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterProveNonDeterministicPendingState is a free log retrieval operation binding the contract event 0x1f44c21118c4603cfb4e1b621dbcfa2b73efcececee2b99b620b2953d33a7010.
//
// Solidity: event ProveNonDeterministicPendingState(bytes32 storedStateRoot, bytes32 provedStateRoot)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) FilterProveNonDeterministicPendingState(opts *bind.FilterOpts) (*FeijoapolygonrollupmanagerProveNonDeterministicPendingStateIterator, error) {

	logs, sub, err := _Feijoapolygonrollupmanager.contract.FilterLogs(opts, "ProveNonDeterministicPendingState")
	if err != nil {
		return nil, err
	}
	return &FeijoapolygonrollupmanagerProveNonDeterministicPendingStateIterator{contract: _Feijoapolygonrollupmanager.contract, event: "ProveNonDeterministicPendingState", logs: logs, sub: sub}, nil
}

// WatchProveNonDeterministicPendingState is a free log subscription operation binding the contract event 0x1f44c21118c4603cfb4e1b621dbcfa2b73efcececee2b99b620b2953d33a7010.
//
// Solidity: event ProveNonDeterministicPendingState(bytes32 storedStateRoot, bytes32 provedStateRoot)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) WatchProveNonDeterministicPendingState(opts *bind.WatchOpts, sink chan<- *FeijoapolygonrollupmanagerProveNonDeterministicPendingState) (event.Subscription, error) {

	logs, sub, err := _Feijoapolygonrollupmanager.contract.WatchLogs(opts, "ProveNonDeterministicPendingState")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeijoapolygonrollupmanagerProveNonDeterministicPendingState)
				if err := _Feijoapolygonrollupmanager.contract.UnpackLog(event, "ProveNonDeterministicPendingState", log); err != nil {
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

// ParseProveNonDeterministicPendingState is a log parse operation binding the contract event 0x1f44c21118c4603cfb4e1b621dbcfa2b73efcececee2b99b620b2953d33a7010.
//
// Solidity: event ProveNonDeterministicPendingState(bytes32 storedStateRoot, bytes32 provedStateRoot)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) ParseProveNonDeterministicPendingState(log types.Log) (*FeijoapolygonrollupmanagerProveNonDeterministicPendingState, error) {
	event := new(FeijoapolygonrollupmanagerProveNonDeterministicPendingState)
	if err := _Feijoapolygonrollupmanager.contract.UnpackLog(event, "ProveNonDeterministicPendingState", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeijoapolygonrollupmanagerRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the Feijoapolygonrollupmanager contract.
type FeijoapolygonrollupmanagerRoleAdminChangedIterator struct {
	Event *FeijoapolygonrollupmanagerRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *FeijoapolygonrollupmanagerRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeijoapolygonrollupmanagerRoleAdminChanged)
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
		it.Event = new(FeijoapolygonrollupmanagerRoleAdminChanged)
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
func (it *FeijoapolygonrollupmanagerRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeijoapolygonrollupmanagerRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeijoapolygonrollupmanagerRoleAdminChanged represents a RoleAdminChanged event raised by the Feijoapolygonrollupmanager contract.
type FeijoapolygonrollupmanagerRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*FeijoapolygonrollupmanagerRoleAdminChangedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _Feijoapolygonrollupmanager.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &FeijoapolygonrollupmanagerRoleAdminChangedIterator{contract: _Feijoapolygonrollupmanager.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *FeijoapolygonrollupmanagerRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _Feijoapolygonrollupmanager.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeijoapolygonrollupmanagerRoleAdminChanged)
				if err := _Feijoapolygonrollupmanager.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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

// ParseRoleAdminChanged is a log parse operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) ParseRoleAdminChanged(log types.Log) (*FeijoapolygonrollupmanagerRoleAdminChanged, error) {
	event := new(FeijoapolygonrollupmanagerRoleAdminChanged)
	if err := _Feijoapolygonrollupmanager.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeijoapolygonrollupmanagerRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the Feijoapolygonrollupmanager contract.
type FeijoapolygonrollupmanagerRoleGrantedIterator struct {
	Event *FeijoapolygonrollupmanagerRoleGranted // Event containing the contract specifics and raw log

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
func (it *FeijoapolygonrollupmanagerRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeijoapolygonrollupmanagerRoleGranted)
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
		it.Event = new(FeijoapolygonrollupmanagerRoleGranted)
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
func (it *FeijoapolygonrollupmanagerRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeijoapolygonrollupmanagerRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeijoapolygonrollupmanagerRoleGranted represents a RoleGranted event raised by the Feijoapolygonrollupmanager contract.
type FeijoapolygonrollupmanagerRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*FeijoapolygonrollupmanagerRoleGrantedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Feijoapolygonrollupmanager.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &FeijoapolygonrollupmanagerRoleGrantedIterator{contract: _Feijoapolygonrollupmanager.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *FeijoapolygonrollupmanagerRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Feijoapolygonrollupmanager.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeijoapolygonrollupmanagerRoleGranted)
				if err := _Feijoapolygonrollupmanager.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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

// ParseRoleGranted is a log parse operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) ParseRoleGranted(log types.Log) (*FeijoapolygonrollupmanagerRoleGranted, error) {
	event := new(FeijoapolygonrollupmanagerRoleGranted)
	if err := _Feijoapolygonrollupmanager.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeijoapolygonrollupmanagerRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the Feijoapolygonrollupmanager contract.
type FeijoapolygonrollupmanagerRoleRevokedIterator struct {
	Event *FeijoapolygonrollupmanagerRoleRevoked // Event containing the contract specifics and raw log

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
func (it *FeijoapolygonrollupmanagerRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeijoapolygonrollupmanagerRoleRevoked)
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
		it.Event = new(FeijoapolygonrollupmanagerRoleRevoked)
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
func (it *FeijoapolygonrollupmanagerRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeijoapolygonrollupmanagerRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeijoapolygonrollupmanagerRoleRevoked represents a RoleRevoked event raised by the Feijoapolygonrollupmanager contract.
type FeijoapolygonrollupmanagerRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*FeijoapolygonrollupmanagerRoleRevokedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Feijoapolygonrollupmanager.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &FeijoapolygonrollupmanagerRoleRevokedIterator{contract: _Feijoapolygonrollupmanager.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *FeijoapolygonrollupmanagerRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _Feijoapolygonrollupmanager.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeijoapolygonrollupmanagerRoleRevoked)
				if err := _Feijoapolygonrollupmanager.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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

// ParseRoleRevoked is a log parse operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) ParseRoleRevoked(log types.Log) (*FeijoapolygonrollupmanagerRoleRevoked, error) {
	event := new(FeijoapolygonrollupmanagerRoleRevoked)
	if err := _Feijoapolygonrollupmanager.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeijoapolygonrollupmanagerSetAggregateRollupVerifierIterator is returned from FilterSetAggregateRollupVerifier and is used to iterate over the raw logs and unpacked data for SetAggregateRollupVerifier events raised by the Feijoapolygonrollupmanager contract.
type FeijoapolygonrollupmanagerSetAggregateRollupVerifierIterator struct {
	Event *FeijoapolygonrollupmanagerSetAggregateRollupVerifier // Event containing the contract specifics and raw log

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
func (it *FeijoapolygonrollupmanagerSetAggregateRollupVerifierIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeijoapolygonrollupmanagerSetAggregateRollupVerifier)
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
		it.Event = new(FeijoapolygonrollupmanagerSetAggregateRollupVerifier)
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
func (it *FeijoapolygonrollupmanagerSetAggregateRollupVerifierIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeijoapolygonrollupmanagerSetAggregateRollupVerifierIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeijoapolygonrollupmanagerSetAggregateRollupVerifier represents a SetAggregateRollupVerifier event raised by the Feijoapolygonrollupmanager contract.
type FeijoapolygonrollupmanagerSetAggregateRollupVerifier struct {
	AggregateRollupVerifier common.Address
	Raw                     types.Log // Blockchain specific contextual infos
}

// FilterSetAggregateRollupVerifier is a free log retrieval operation binding the contract event 0x53ab89ca5f00e99098ada1782f593e3f76b5489459ece48450e554c2928daa5e.
//
// Solidity: event SetAggregateRollupVerifier(address aggregateRollupVerifier)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) FilterSetAggregateRollupVerifier(opts *bind.FilterOpts) (*FeijoapolygonrollupmanagerSetAggregateRollupVerifierIterator, error) {

	logs, sub, err := _Feijoapolygonrollupmanager.contract.FilterLogs(opts, "SetAggregateRollupVerifier")
	if err != nil {
		return nil, err
	}
	return &FeijoapolygonrollupmanagerSetAggregateRollupVerifierIterator{contract: _Feijoapolygonrollupmanager.contract, event: "SetAggregateRollupVerifier", logs: logs, sub: sub}, nil
}

// WatchSetAggregateRollupVerifier is a free log subscription operation binding the contract event 0x53ab89ca5f00e99098ada1782f593e3f76b5489459ece48450e554c2928daa5e.
//
// Solidity: event SetAggregateRollupVerifier(address aggregateRollupVerifier)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) WatchSetAggregateRollupVerifier(opts *bind.WatchOpts, sink chan<- *FeijoapolygonrollupmanagerSetAggregateRollupVerifier) (event.Subscription, error) {

	logs, sub, err := _Feijoapolygonrollupmanager.contract.WatchLogs(opts, "SetAggregateRollupVerifier")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeijoapolygonrollupmanagerSetAggregateRollupVerifier)
				if err := _Feijoapolygonrollupmanager.contract.UnpackLog(event, "SetAggregateRollupVerifier", log); err != nil {
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

// ParseSetAggregateRollupVerifier is a log parse operation binding the contract event 0x53ab89ca5f00e99098ada1782f593e3f76b5489459ece48450e554c2928daa5e.
//
// Solidity: event SetAggregateRollupVerifier(address aggregateRollupVerifier)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) ParseSetAggregateRollupVerifier(log types.Log) (*FeijoapolygonrollupmanagerSetAggregateRollupVerifier, error) {
	event := new(FeijoapolygonrollupmanagerSetAggregateRollupVerifier)
	if err := _Feijoapolygonrollupmanager.contract.UnpackLog(event, "SetAggregateRollupVerifier", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeijoapolygonrollupmanagerSetMultiplierZkGasPriceIterator is returned from FilterSetMultiplierZkGasPrice and is used to iterate over the raw logs and unpacked data for SetMultiplierZkGasPrice events raised by the Feijoapolygonrollupmanager contract.
type FeijoapolygonrollupmanagerSetMultiplierZkGasPriceIterator struct {
	Event *FeijoapolygonrollupmanagerSetMultiplierZkGasPrice // Event containing the contract specifics and raw log

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
func (it *FeijoapolygonrollupmanagerSetMultiplierZkGasPriceIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeijoapolygonrollupmanagerSetMultiplierZkGasPrice)
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
		it.Event = new(FeijoapolygonrollupmanagerSetMultiplierZkGasPrice)
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
func (it *FeijoapolygonrollupmanagerSetMultiplierZkGasPriceIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeijoapolygonrollupmanagerSetMultiplierZkGasPriceIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeijoapolygonrollupmanagerSetMultiplierZkGasPrice represents a SetMultiplierZkGasPrice event raised by the Feijoapolygonrollupmanager contract.
type FeijoapolygonrollupmanagerSetMultiplierZkGasPrice struct {
	NewMultiplierSequenceFee uint16
	Raw                      types.Log // Blockchain specific contextual infos
}

// FilterSetMultiplierZkGasPrice is a free log retrieval operation binding the contract event 0x5c8a9e64670a8ec12a8004aa047cbb455403a6c4f2d2ad4e52328400dc814265.
//
// Solidity: event SetMultiplierZkGasPrice(uint16 newMultiplierSequenceFee)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) FilterSetMultiplierZkGasPrice(opts *bind.FilterOpts) (*FeijoapolygonrollupmanagerSetMultiplierZkGasPriceIterator, error) {

	logs, sub, err := _Feijoapolygonrollupmanager.contract.FilterLogs(opts, "SetMultiplierZkGasPrice")
	if err != nil {
		return nil, err
	}
	return &FeijoapolygonrollupmanagerSetMultiplierZkGasPriceIterator{contract: _Feijoapolygonrollupmanager.contract, event: "SetMultiplierZkGasPrice", logs: logs, sub: sub}, nil
}

// WatchSetMultiplierZkGasPrice is a free log subscription operation binding the contract event 0x5c8a9e64670a8ec12a8004aa047cbb455403a6c4f2d2ad4e52328400dc814265.
//
// Solidity: event SetMultiplierZkGasPrice(uint16 newMultiplierSequenceFee)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) WatchSetMultiplierZkGasPrice(opts *bind.WatchOpts, sink chan<- *FeijoapolygonrollupmanagerSetMultiplierZkGasPrice) (event.Subscription, error) {

	logs, sub, err := _Feijoapolygonrollupmanager.contract.WatchLogs(opts, "SetMultiplierZkGasPrice")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeijoapolygonrollupmanagerSetMultiplierZkGasPrice)
				if err := _Feijoapolygonrollupmanager.contract.UnpackLog(event, "SetMultiplierZkGasPrice", log); err != nil {
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

// ParseSetMultiplierZkGasPrice is a log parse operation binding the contract event 0x5c8a9e64670a8ec12a8004aa047cbb455403a6c4f2d2ad4e52328400dc814265.
//
// Solidity: event SetMultiplierZkGasPrice(uint16 newMultiplierSequenceFee)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) ParseSetMultiplierZkGasPrice(log types.Log) (*FeijoapolygonrollupmanagerSetMultiplierZkGasPrice, error) {
	event := new(FeijoapolygonrollupmanagerSetMultiplierZkGasPrice)
	if err := _Feijoapolygonrollupmanager.contract.UnpackLog(event, "SetMultiplierZkGasPrice", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeijoapolygonrollupmanagerSetPendingStateTimeoutIterator is returned from FilterSetPendingStateTimeout and is used to iterate over the raw logs and unpacked data for SetPendingStateTimeout events raised by the Feijoapolygonrollupmanager contract.
type FeijoapolygonrollupmanagerSetPendingStateTimeoutIterator struct {
	Event *FeijoapolygonrollupmanagerSetPendingStateTimeout // Event containing the contract specifics and raw log

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
func (it *FeijoapolygonrollupmanagerSetPendingStateTimeoutIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeijoapolygonrollupmanagerSetPendingStateTimeout)
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
		it.Event = new(FeijoapolygonrollupmanagerSetPendingStateTimeout)
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
func (it *FeijoapolygonrollupmanagerSetPendingStateTimeoutIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeijoapolygonrollupmanagerSetPendingStateTimeoutIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeijoapolygonrollupmanagerSetPendingStateTimeout represents a SetPendingStateTimeout event raised by the Feijoapolygonrollupmanager contract.
type FeijoapolygonrollupmanagerSetPendingStateTimeout struct {
	NewPendingStateTimeout uint64
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterSetPendingStateTimeout is a free log retrieval operation binding the contract event 0xc4121f4e22c69632ebb7cf1f462be0511dc034f999b52013eddfb24aab765c75.
//
// Solidity: event SetPendingStateTimeout(uint64 newPendingStateTimeout)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) FilterSetPendingStateTimeout(opts *bind.FilterOpts) (*FeijoapolygonrollupmanagerSetPendingStateTimeoutIterator, error) {

	logs, sub, err := _Feijoapolygonrollupmanager.contract.FilterLogs(opts, "SetPendingStateTimeout")
	if err != nil {
		return nil, err
	}
	return &FeijoapolygonrollupmanagerSetPendingStateTimeoutIterator{contract: _Feijoapolygonrollupmanager.contract, event: "SetPendingStateTimeout", logs: logs, sub: sub}, nil
}

// WatchSetPendingStateTimeout is a free log subscription operation binding the contract event 0xc4121f4e22c69632ebb7cf1f462be0511dc034f999b52013eddfb24aab765c75.
//
// Solidity: event SetPendingStateTimeout(uint64 newPendingStateTimeout)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) WatchSetPendingStateTimeout(opts *bind.WatchOpts, sink chan<- *FeijoapolygonrollupmanagerSetPendingStateTimeout) (event.Subscription, error) {

	logs, sub, err := _Feijoapolygonrollupmanager.contract.WatchLogs(opts, "SetPendingStateTimeout")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeijoapolygonrollupmanagerSetPendingStateTimeout)
				if err := _Feijoapolygonrollupmanager.contract.UnpackLog(event, "SetPendingStateTimeout", log); err != nil {
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

// ParseSetPendingStateTimeout is a log parse operation binding the contract event 0xc4121f4e22c69632ebb7cf1f462be0511dc034f999b52013eddfb24aab765c75.
//
// Solidity: event SetPendingStateTimeout(uint64 newPendingStateTimeout)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) ParseSetPendingStateTimeout(log types.Log) (*FeijoapolygonrollupmanagerSetPendingStateTimeout, error) {
	event := new(FeijoapolygonrollupmanagerSetPendingStateTimeout)
	if err := _Feijoapolygonrollupmanager.contract.UnpackLog(event, "SetPendingStateTimeout", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeijoapolygonrollupmanagerSetSequenceFeeIterator is returned from FilterSetSequenceFee and is used to iterate over the raw logs and unpacked data for SetSequenceFee events raised by the Feijoapolygonrollupmanager contract.
type FeijoapolygonrollupmanagerSetSequenceFeeIterator struct {
	Event *FeijoapolygonrollupmanagerSetSequenceFee // Event containing the contract specifics and raw log

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
func (it *FeijoapolygonrollupmanagerSetSequenceFeeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeijoapolygonrollupmanagerSetSequenceFee)
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
		it.Event = new(FeijoapolygonrollupmanagerSetSequenceFee)
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
func (it *FeijoapolygonrollupmanagerSetSequenceFeeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeijoapolygonrollupmanagerSetSequenceFeeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeijoapolygonrollupmanagerSetSequenceFee represents a SetSequenceFee event raised by the Feijoapolygonrollupmanager contract.
type FeijoapolygonrollupmanagerSetSequenceFee struct {
	NewSequenceFee *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterSetSequenceFee is a free log retrieval operation binding the contract event 0x13b1c630ad78354572e9ad473455d51831407e164b79dda20732f5acac503382.
//
// Solidity: event SetSequenceFee(uint256 newSequenceFee)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) FilterSetSequenceFee(opts *bind.FilterOpts) (*FeijoapolygonrollupmanagerSetSequenceFeeIterator, error) {

	logs, sub, err := _Feijoapolygonrollupmanager.contract.FilterLogs(opts, "SetSequenceFee")
	if err != nil {
		return nil, err
	}
	return &FeijoapolygonrollupmanagerSetSequenceFeeIterator{contract: _Feijoapolygonrollupmanager.contract, event: "SetSequenceFee", logs: logs, sub: sub}, nil
}

// WatchSetSequenceFee is a free log subscription operation binding the contract event 0x13b1c630ad78354572e9ad473455d51831407e164b79dda20732f5acac503382.
//
// Solidity: event SetSequenceFee(uint256 newSequenceFee)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) WatchSetSequenceFee(opts *bind.WatchOpts, sink chan<- *FeijoapolygonrollupmanagerSetSequenceFee) (event.Subscription, error) {

	logs, sub, err := _Feijoapolygonrollupmanager.contract.WatchLogs(opts, "SetSequenceFee")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeijoapolygonrollupmanagerSetSequenceFee)
				if err := _Feijoapolygonrollupmanager.contract.UnpackLog(event, "SetSequenceFee", log); err != nil {
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

// ParseSetSequenceFee is a log parse operation binding the contract event 0x13b1c630ad78354572e9ad473455d51831407e164b79dda20732f5acac503382.
//
// Solidity: event SetSequenceFee(uint256 newSequenceFee)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) ParseSetSequenceFee(log types.Log) (*FeijoapolygonrollupmanagerSetSequenceFee, error) {
	event := new(FeijoapolygonrollupmanagerSetSequenceFee)
	if err := _Feijoapolygonrollupmanager.contract.UnpackLog(event, "SetSequenceFee", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeijoapolygonrollupmanagerSetTrustedAggregatorIterator is returned from FilterSetTrustedAggregator and is used to iterate over the raw logs and unpacked data for SetTrustedAggregator events raised by the Feijoapolygonrollupmanager contract.
type FeijoapolygonrollupmanagerSetTrustedAggregatorIterator struct {
	Event *FeijoapolygonrollupmanagerSetTrustedAggregator // Event containing the contract specifics and raw log

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
func (it *FeijoapolygonrollupmanagerSetTrustedAggregatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeijoapolygonrollupmanagerSetTrustedAggregator)
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
		it.Event = new(FeijoapolygonrollupmanagerSetTrustedAggregator)
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
func (it *FeijoapolygonrollupmanagerSetTrustedAggregatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeijoapolygonrollupmanagerSetTrustedAggregatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeijoapolygonrollupmanagerSetTrustedAggregator represents a SetTrustedAggregator event raised by the Feijoapolygonrollupmanager contract.
type FeijoapolygonrollupmanagerSetTrustedAggregator struct {
	NewTrustedAggregator common.Address
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterSetTrustedAggregator is a free log retrieval operation binding the contract event 0x61f8fec29495a3078e9271456f05fb0707fd4e41f7661865f80fc437d06681ca.
//
// Solidity: event SetTrustedAggregator(address newTrustedAggregator)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) FilterSetTrustedAggregator(opts *bind.FilterOpts) (*FeijoapolygonrollupmanagerSetTrustedAggregatorIterator, error) {

	logs, sub, err := _Feijoapolygonrollupmanager.contract.FilterLogs(opts, "SetTrustedAggregator")
	if err != nil {
		return nil, err
	}
	return &FeijoapolygonrollupmanagerSetTrustedAggregatorIterator{contract: _Feijoapolygonrollupmanager.contract, event: "SetTrustedAggregator", logs: logs, sub: sub}, nil
}

// WatchSetTrustedAggregator is a free log subscription operation binding the contract event 0x61f8fec29495a3078e9271456f05fb0707fd4e41f7661865f80fc437d06681ca.
//
// Solidity: event SetTrustedAggregator(address newTrustedAggregator)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) WatchSetTrustedAggregator(opts *bind.WatchOpts, sink chan<- *FeijoapolygonrollupmanagerSetTrustedAggregator) (event.Subscription, error) {

	logs, sub, err := _Feijoapolygonrollupmanager.contract.WatchLogs(opts, "SetTrustedAggregator")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeijoapolygonrollupmanagerSetTrustedAggregator)
				if err := _Feijoapolygonrollupmanager.contract.UnpackLog(event, "SetTrustedAggregator", log); err != nil {
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

// ParseSetTrustedAggregator is a log parse operation binding the contract event 0x61f8fec29495a3078e9271456f05fb0707fd4e41f7661865f80fc437d06681ca.
//
// Solidity: event SetTrustedAggregator(address newTrustedAggregator)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) ParseSetTrustedAggregator(log types.Log) (*FeijoapolygonrollupmanagerSetTrustedAggregator, error) {
	event := new(FeijoapolygonrollupmanagerSetTrustedAggregator)
	if err := _Feijoapolygonrollupmanager.contract.UnpackLog(event, "SetTrustedAggregator", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeijoapolygonrollupmanagerSetTrustedAggregatorTimeoutIterator is returned from FilterSetTrustedAggregatorTimeout and is used to iterate over the raw logs and unpacked data for SetTrustedAggregatorTimeout events raised by the Feijoapolygonrollupmanager contract.
type FeijoapolygonrollupmanagerSetTrustedAggregatorTimeoutIterator struct {
	Event *FeijoapolygonrollupmanagerSetTrustedAggregatorTimeout // Event containing the contract specifics and raw log

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
func (it *FeijoapolygonrollupmanagerSetTrustedAggregatorTimeoutIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeijoapolygonrollupmanagerSetTrustedAggregatorTimeout)
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
		it.Event = new(FeijoapolygonrollupmanagerSetTrustedAggregatorTimeout)
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
func (it *FeijoapolygonrollupmanagerSetTrustedAggregatorTimeoutIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeijoapolygonrollupmanagerSetTrustedAggregatorTimeoutIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeijoapolygonrollupmanagerSetTrustedAggregatorTimeout represents a SetTrustedAggregatorTimeout event raised by the Feijoapolygonrollupmanager contract.
type FeijoapolygonrollupmanagerSetTrustedAggregatorTimeout struct {
	NewTrustedAggregatorTimeout uint64
	Raw                         types.Log // Blockchain specific contextual infos
}

// FilterSetTrustedAggregatorTimeout is a free log retrieval operation binding the contract event 0x1f4fa24c2e4bad19a7f3ec5c5485f70d46c798461c2e684f55bbd0fc661373a1.
//
// Solidity: event SetTrustedAggregatorTimeout(uint64 newTrustedAggregatorTimeout)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) FilterSetTrustedAggregatorTimeout(opts *bind.FilterOpts) (*FeijoapolygonrollupmanagerSetTrustedAggregatorTimeoutIterator, error) {

	logs, sub, err := _Feijoapolygonrollupmanager.contract.FilterLogs(opts, "SetTrustedAggregatorTimeout")
	if err != nil {
		return nil, err
	}
	return &FeijoapolygonrollupmanagerSetTrustedAggregatorTimeoutIterator{contract: _Feijoapolygonrollupmanager.contract, event: "SetTrustedAggregatorTimeout", logs: logs, sub: sub}, nil
}

// WatchSetTrustedAggregatorTimeout is a free log subscription operation binding the contract event 0x1f4fa24c2e4bad19a7f3ec5c5485f70d46c798461c2e684f55bbd0fc661373a1.
//
// Solidity: event SetTrustedAggregatorTimeout(uint64 newTrustedAggregatorTimeout)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) WatchSetTrustedAggregatorTimeout(opts *bind.WatchOpts, sink chan<- *FeijoapolygonrollupmanagerSetTrustedAggregatorTimeout) (event.Subscription, error) {

	logs, sub, err := _Feijoapolygonrollupmanager.contract.WatchLogs(opts, "SetTrustedAggregatorTimeout")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeijoapolygonrollupmanagerSetTrustedAggregatorTimeout)
				if err := _Feijoapolygonrollupmanager.contract.UnpackLog(event, "SetTrustedAggregatorTimeout", log); err != nil {
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

// ParseSetTrustedAggregatorTimeout is a log parse operation binding the contract event 0x1f4fa24c2e4bad19a7f3ec5c5485f70d46c798461c2e684f55bbd0fc661373a1.
//
// Solidity: event SetTrustedAggregatorTimeout(uint64 newTrustedAggregatorTimeout)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) ParseSetTrustedAggregatorTimeout(log types.Log) (*FeijoapolygonrollupmanagerSetTrustedAggregatorTimeout, error) {
	event := new(FeijoapolygonrollupmanagerSetTrustedAggregatorTimeout)
	if err := _Feijoapolygonrollupmanager.contract.UnpackLog(event, "SetTrustedAggregatorTimeout", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeijoapolygonrollupmanagerSetVerifySequenceTimeTargetIterator is returned from FilterSetVerifySequenceTimeTarget and is used to iterate over the raw logs and unpacked data for SetVerifySequenceTimeTarget events raised by the Feijoapolygonrollupmanager contract.
type FeijoapolygonrollupmanagerSetVerifySequenceTimeTargetIterator struct {
	Event *FeijoapolygonrollupmanagerSetVerifySequenceTimeTarget // Event containing the contract specifics and raw log

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
func (it *FeijoapolygonrollupmanagerSetVerifySequenceTimeTargetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeijoapolygonrollupmanagerSetVerifySequenceTimeTarget)
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
		it.Event = new(FeijoapolygonrollupmanagerSetVerifySequenceTimeTarget)
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
func (it *FeijoapolygonrollupmanagerSetVerifySequenceTimeTargetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeijoapolygonrollupmanagerSetVerifySequenceTimeTargetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeijoapolygonrollupmanagerSetVerifySequenceTimeTarget represents a SetVerifySequenceTimeTarget event raised by the Feijoapolygonrollupmanager contract.
type FeijoapolygonrollupmanagerSetVerifySequenceTimeTarget struct {
	NewVerifySequenceTimeTarget uint64
	Raw                         types.Log // Blockchain specific contextual infos
}

// FilterSetVerifySequenceTimeTarget is a free log retrieval operation binding the contract event 0xe84eacb10b29a9cd283d1c48f59cd87da8c2f99c554576228566d69aeba740cd.
//
// Solidity: event SetVerifySequenceTimeTarget(uint64 newVerifySequenceTimeTarget)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) FilterSetVerifySequenceTimeTarget(opts *bind.FilterOpts) (*FeijoapolygonrollupmanagerSetVerifySequenceTimeTargetIterator, error) {

	logs, sub, err := _Feijoapolygonrollupmanager.contract.FilterLogs(opts, "SetVerifySequenceTimeTarget")
	if err != nil {
		return nil, err
	}
	return &FeijoapolygonrollupmanagerSetVerifySequenceTimeTargetIterator{contract: _Feijoapolygonrollupmanager.contract, event: "SetVerifySequenceTimeTarget", logs: logs, sub: sub}, nil
}

// WatchSetVerifySequenceTimeTarget is a free log subscription operation binding the contract event 0xe84eacb10b29a9cd283d1c48f59cd87da8c2f99c554576228566d69aeba740cd.
//
// Solidity: event SetVerifySequenceTimeTarget(uint64 newVerifySequenceTimeTarget)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) WatchSetVerifySequenceTimeTarget(opts *bind.WatchOpts, sink chan<- *FeijoapolygonrollupmanagerSetVerifySequenceTimeTarget) (event.Subscription, error) {

	logs, sub, err := _Feijoapolygonrollupmanager.contract.WatchLogs(opts, "SetVerifySequenceTimeTarget")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeijoapolygonrollupmanagerSetVerifySequenceTimeTarget)
				if err := _Feijoapolygonrollupmanager.contract.UnpackLog(event, "SetVerifySequenceTimeTarget", log); err != nil {
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

// ParseSetVerifySequenceTimeTarget is a log parse operation binding the contract event 0xe84eacb10b29a9cd283d1c48f59cd87da8c2f99c554576228566d69aeba740cd.
//
// Solidity: event SetVerifySequenceTimeTarget(uint64 newVerifySequenceTimeTarget)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) ParseSetVerifySequenceTimeTarget(log types.Log) (*FeijoapolygonrollupmanagerSetVerifySequenceTimeTarget, error) {
	event := new(FeijoapolygonrollupmanagerSetVerifySequenceTimeTarget)
	if err := _Feijoapolygonrollupmanager.contract.UnpackLog(event, "SetVerifySequenceTimeTarget", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeijoapolygonrollupmanagerUpdateRollupIterator is returned from FilterUpdateRollup and is used to iterate over the raw logs and unpacked data for UpdateRollup events raised by the Feijoapolygonrollupmanager contract.
type FeijoapolygonrollupmanagerUpdateRollupIterator struct {
	Event *FeijoapolygonrollupmanagerUpdateRollup // Event containing the contract specifics and raw log

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
func (it *FeijoapolygonrollupmanagerUpdateRollupIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeijoapolygonrollupmanagerUpdateRollup)
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
		it.Event = new(FeijoapolygonrollupmanagerUpdateRollup)
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
func (it *FeijoapolygonrollupmanagerUpdateRollupIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeijoapolygonrollupmanagerUpdateRollupIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeijoapolygonrollupmanagerUpdateRollup represents a UpdateRollup event raised by the Feijoapolygonrollupmanager contract.
type FeijoapolygonrollupmanagerUpdateRollup struct {
	RollupID                          uint32
	NewRollupTypeID                   uint32
	LastVerifiedSequenceBeforeUpgrade uint64
	Raw                               types.Log // Blockchain specific contextual infos
}

// FilterUpdateRollup is a free log retrieval operation binding the contract event 0xf585e04c05d396901170247783d3e5f0ee9c1df23072985b50af089f5e48b19d.
//
// Solidity: event UpdateRollup(uint32 indexed rollupID, uint32 newRollupTypeID, uint64 lastVerifiedSequenceBeforeUpgrade)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) FilterUpdateRollup(opts *bind.FilterOpts, rollupID []uint32) (*FeijoapolygonrollupmanagerUpdateRollupIterator, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	logs, sub, err := _Feijoapolygonrollupmanager.contract.FilterLogs(opts, "UpdateRollup", rollupIDRule)
	if err != nil {
		return nil, err
	}
	return &FeijoapolygonrollupmanagerUpdateRollupIterator{contract: _Feijoapolygonrollupmanager.contract, event: "UpdateRollup", logs: logs, sub: sub}, nil
}

// WatchUpdateRollup is a free log subscription operation binding the contract event 0xf585e04c05d396901170247783d3e5f0ee9c1df23072985b50af089f5e48b19d.
//
// Solidity: event UpdateRollup(uint32 indexed rollupID, uint32 newRollupTypeID, uint64 lastVerifiedSequenceBeforeUpgrade)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) WatchUpdateRollup(opts *bind.WatchOpts, sink chan<- *FeijoapolygonrollupmanagerUpdateRollup, rollupID []uint32) (event.Subscription, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	logs, sub, err := _Feijoapolygonrollupmanager.contract.WatchLogs(opts, "UpdateRollup", rollupIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeijoapolygonrollupmanagerUpdateRollup)
				if err := _Feijoapolygonrollupmanager.contract.UnpackLog(event, "UpdateRollup", log); err != nil {
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

// ParseUpdateRollup is a log parse operation binding the contract event 0xf585e04c05d396901170247783d3e5f0ee9c1df23072985b50af089f5e48b19d.
//
// Solidity: event UpdateRollup(uint32 indexed rollupID, uint32 newRollupTypeID, uint64 lastVerifiedSequenceBeforeUpgrade)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) ParseUpdateRollup(log types.Log) (*FeijoapolygonrollupmanagerUpdateRollup, error) {
	event := new(FeijoapolygonrollupmanagerUpdateRollup)
	if err := _Feijoapolygonrollupmanager.contract.UnpackLog(event, "UpdateRollup", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeijoapolygonrollupmanagerVerifySequencesIterator is returned from FilterVerifySequences and is used to iterate over the raw logs and unpacked data for VerifySequences events raised by the Feijoapolygonrollupmanager contract.
type FeijoapolygonrollupmanagerVerifySequencesIterator struct {
	Event *FeijoapolygonrollupmanagerVerifySequences // Event containing the contract specifics and raw log

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
func (it *FeijoapolygonrollupmanagerVerifySequencesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeijoapolygonrollupmanagerVerifySequences)
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
		it.Event = new(FeijoapolygonrollupmanagerVerifySequences)
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
func (it *FeijoapolygonrollupmanagerVerifySequencesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeijoapolygonrollupmanagerVerifySequencesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeijoapolygonrollupmanagerVerifySequences represents a VerifySequences event raised by the Feijoapolygonrollupmanager contract.
type FeijoapolygonrollupmanagerVerifySequences struct {
	RollupID    uint32
	SequenceNum uint64
	StateRoot   [32]byte
	ExitRoot    [32]byte
	Aggregator  common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterVerifySequences is a free log retrieval operation binding the contract event 0x716b8543c1c3c328a13d34cd51e064a780149a2d06455e44097de219b150e8b4.
//
// Solidity: event VerifySequences(uint32 indexed rollupID, uint64 sequenceNum, bytes32 stateRoot, bytes32 exitRoot, address indexed aggregator)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) FilterVerifySequences(opts *bind.FilterOpts, rollupID []uint32, aggregator []common.Address) (*FeijoapolygonrollupmanagerVerifySequencesIterator, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	var aggregatorRule []interface{}
	for _, aggregatorItem := range aggregator {
		aggregatorRule = append(aggregatorRule, aggregatorItem)
	}

	logs, sub, err := _Feijoapolygonrollupmanager.contract.FilterLogs(opts, "VerifySequences", rollupIDRule, aggregatorRule)
	if err != nil {
		return nil, err
	}
	return &FeijoapolygonrollupmanagerVerifySequencesIterator{contract: _Feijoapolygonrollupmanager.contract, event: "VerifySequences", logs: logs, sub: sub}, nil
}

// WatchVerifySequences is a free log subscription operation binding the contract event 0x716b8543c1c3c328a13d34cd51e064a780149a2d06455e44097de219b150e8b4.
//
// Solidity: event VerifySequences(uint32 indexed rollupID, uint64 sequenceNum, bytes32 stateRoot, bytes32 exitRoot, address indexed aggregator)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) WatchVerifySequences(opts *bind.WatchOpts, sink chan<- *FeijoapolygonrollupmanagerVerifySequences, rollupID []uint32, aggregator []common.Address) (event.Subscription, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	var aggregatorRule []interface{}
	for _, aggregatorItem := range aggregator {
		aggregatorRule = append(aggregatorRule, aggregatorItem)
	}

	logs, sub, err := _Feijoapolygonrollupmanager.contract.WatchLogs(opts, "VerifySequences", rollupIDRule, aggregatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeijoapolygonrollupmanagerVerifySequences)
				if err := _Feijoapolygonrollupmanager.contract.UnpackLog(event, "VerifySequences", log); err != nil {
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

// ParseVerifySequences is a log parse operation binding the contract event 0x716b8543c1c3c328a13d34cd51e064a780149a2d06455e44097de219b150e8b4.
//
// Solidity: event VerifySequences(uint32 indexed rollupID, uint64 sequenceNum, bytes32 stateRoot, bytes32 exitRoot, address indexed aggregator)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) ParseVerifySequences(log types.Log) (*FeijoapolygonrollupmanagerVerifySequences, error) {
	event := new(FeijoapolygonrollupmanagerVerifySequences)
	if err := _Feijoapolygonrollupmanager.contract.UnpackLog(event, "VerifySequences", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeijoapolygonrollupmanagerVerifySequencesMultiProofIterator is returned from FilterVerifySequencesMultiProof and is used to iterate over the raw logs and unpacked data for VerifySequencesMultiProof events raised by the Feijoapolygonrollupmanager contract.
type FeijoapolygonrollupmanagerVerifySequencesMultiProofIterator struct {
	Event *FeijoapolygonrollupmanagerVerifySequencesMultiProof // Event containing the contract specifics and raw log

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
func (it *FeijoapolygonrollupmanagerVerifySequencesMultiProofIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeijoapolygonrollupmanagerVerifySequencesMultiProof)
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
		it.Event = new(FeijoapolygonrollupmanagerVerifySequencesMultiProof)
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
func (it *FeijoapolygonrollupmanagerVerifySequencesMultiProofIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeijoapolygonrollupmanagerVerifySequencesMultiProofIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeijoapolygonrollupmanagerVerifySequencesMultiProof represents a VerifySequencesMultiProof event raised by the Feijoapolygonrollupmanager contract.
type FeijoapolygonrollupmanagerVerifySequencesMultiProof struct {
	Aggregator common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterVerifySequencesMultiProof is a free log retrieval operation binding the contract event 0x73520b4a8035df0a5543b7c7d63fcb1c3d68d80bd9dce27299f3e03faaf4d7d6.
//
// Solidity: event VerifySequencesMultiProof(address indexed aggregator)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) FilterVerifySequencesMultiProof(opts *bind.FilterOpts, aggregator []common.Address) (*FeijoapolygonrollupmanagerVerifySequencesMultiProofIterator, error) {

	var aggregatorRule []interface{}
	for _, aggregatorItem := range aggregator {
		aggregatorRule = append(aggregatorRule, aggregatorItem)
	}

	logs, sub, err := _Feijoapolygonrollupmanager.contract.FilterLogs(opts, "VerifySequencesMultiProof", aggregatorRule)
	if err != nil {
		return nil, err
	}
	return &FeijoapolygonrollupmanagerVerifySequencesMultiProofIterator{contract: _Feijoapolygonrollupmanager.contract, event: "VerifySequencesMultiProof", logs: logs, sub: sub}, nil
}

// WatchVerifySequencesMultiProof is a free log subscription operation binding the contract event 0x73520b4a8035df0a5543b7c7d63fcb1c3d68d80bd9dce27299f3e03faaf4d7d6.
//
// Solidity: event VerifySequencesMultiProof(address indexed aggregator)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) WatchVerifySequencesMultiProof(opts *bind.WatchOpts, sink chan<- *FeijoapolygonrollupmanagerVerifySequencesMultiProof, aggregator []common.Address) (event.Subscription, error) {

	var aggregatorRule []interface{}
	for _, aggregatorItem := range aggregator {
		aggregatorRule = append(aggregatorRule, aggregatorItem)
	}

	logs, sub, err := _Feijoapolygonrollupmanager.contract.WatchLogs(opts, "VerifySequencesMultiProof", aggregatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeijoapolygonrollupmanagerVerifySequencesMultiProof)
				if err := _Feijoapolygonrollupmanager.contract.UnpackLog(event, "VerifySequencesMultiProof", log); err != nil {
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

// ParseVerifySequencesMultiProof is a log parse operation binding the contract event 0x73520b4a8035df0a5543b7c7d63fcb1c3d68d80bd9dce27299f3e03faaf4d7d6.
//
// Solidity: event VerifySequencesMultiProof(address indexed aggregator)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) ParseVerifySequencesMultiProof(log types.Log) (*FeijoapolygonrollupmanagerVerifySequencesMultiProof, error) {
	event := new(FeijoapolygonrollupmanagerVerifySequencesMultiProof)
	if err := _Feijoapolygonrollupmanager.contract.UnpackLog(event, "VerifySequencesMultiProof", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeijoapolygonrollupmanagerVerifySequencesTrustedAggregatorIterator is returned from FilterVerifySequencesTrustedAggregator and is used to iterate over the raw logs and unpacked data for VerifySequencesTrustedAggregator events raised by the Feijoapolygonrollupmanager contract.
type FeijoapolygonrollupmanagerVerifySequencesTrustedAggregatorIterator struct {
	Event *FeijoapolygonrollupmanagerVerifySequencesTrustedAggregator // Event containing the contract specifics and raw log

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
func (it *FeijoapolygonrollupmanagerVerifySequencesTrustedAggregatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeijoapolygonrollupmanagerVerifySequencesTrustedAggregator)
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
		it.Event = new(FeijoapolygonrollupmanagerVerifySequencesTrustedAggregator)
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
func (it *FeijoapolygonrollupmanagerVerifySequencesTrustedAggregatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeijoapolygonrollupmanagerVerifySequencesTrustedAggregatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeijoapolygonrollupmanagerVerifySequencesTrustedAggregator represents a VerifySequencesTrustedAggregator event raised by the Feijoapolygonrollupmanager contract.
type FeijoapolygonrollupmanagerVerifySequencesTrustedAggregator struct {
	RollupID    uint32
	NumSequence uint64
	StateRoot   [32]byte
	ExitRoot    [32]byte
	Aggregator  common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterVerifySequencesTrustedAggregator is a free log retrieval operation binding the contract event 0xba7fad50a32b4eb9847ff1f56dd7528178eae3cd0b008c7a798e0d5375de88da.
//
// Solidity: event VerifySequencesTrustedAggregator(uint32 indexed rollupID, uint64 numSequence, bytes32 stateRoot, bytes32 exitRoot, address indexed aggregator)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) FilterVerifySequencesTrustedAggregator(opts *bind.FilterOpts, rollupID []uint32, aggregator []common.Address) (*FeijoapolygonrollupmanagerVerifySequencesTrustedAggregatorIterator, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	var aggregatorRule []interface{}
	for _, aggregatorItem := range aggregator {
		aggregatorRule = append(aggregatorRule, aggregatorItem)
	}

	logs, sub, err := _Feijoapolygonrollupmanager.contract.FilterLogs(opts, "VerifySequencesTrustedAggregator", rollupIDRule, aggregatorRule)
	if err != nil {
		return nil, err
	}
	return &FeijoapolygonrollupmanagerVerifySequencesTrustedAggregatorIterator{contract: _Feijoapolygonrollupmanager.contract, event: "VerifySequencesTrustedAggregator", logs: logs, sub: sub}, nil
}

// WatchVerifySequencesTrustedAggregator is a free log subscription operation binding the contract event 0xba7fad50a32b4eb9847ff1f56dd7528178eae3cd0b008c7a798e0d5375de88da.
//
// Solidity: event VerifySequencesTrustedAggregator(uint32 indexed rollupID, uint64 numSequence, bytes32 stateRoot, bytes32 exitRoot, address indexed aggregator)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) WatchVerifySequencesTrustedAggregator(opts *bind.WatchOpts, sink chan<- *FeijoapolygonrollupmanagerVerifySequencesTrustedAggregator, rollupID []uint32, aggregator []common.Address) (event.Subscription, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	var aggregatorRule []interface{}
	for _, aggregatorItem := range aggregator {
		aggregatorRule = append(aggregatorRule, aggregatorItem)
	}

	logs, sub, err := _Feijoapolygonrollupmanager.contract.WatchLogs(opts, "VerifySequencesTrustedAggregator", rollupIDRule, aggregatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeijoapolygonrollupmanagerVerifySequencesTrustedAggregator)
				if err := _Feijoapolygonrollupmanager.contract.UnpackLog(event, "VerifySequencesTrustedAggregator", log); err != nil {
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

// ParseVerifySequencesTrustedAggregator is a log parse operation binding the contract event 0xba7fad50a32b4eb9847ff1f56dd7528178eae3cd0b008c7a798e0d5375de88da.
//
// Solidity: event VerifySequencesTrustedAggregator(uint32 indexed rollupID, uint64 numSequence, bytes32 stateRoot, bytes32 exitRoot, address indexed aggregator)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) ParseVerifySequencesTrustedAggregator(log types.Log) (*FeijoapolygonrollupmanagerVerifySequencesTrustedAggregator, error) {
	event := new(FeijoapolygonrollupmanagerVerifySequencesTrustedAggregator)
	if err := _Feijoapolygonrollupmanager.contract.UnpackLog(event, "VerifySequencesTrustedAggregator", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeijoapolygonrollupmanagerVerifySequencesTrustedAggregatorMultiProofIterator is returned from FilterVerifySequencesTrustedAggregatorMultiProof and is used to iterate over the raw logs and unpacked data for VerifySequencesTrustedAggregatorMultiProof events raised by the Feijoapolygonrollupmanager contract.
type FeijoapolygonrollupmanagerVerifySequencesTrustedAggregatorMultiProofIterator struct {
	Event *FeijoapolygonrollupmanagerVerifySequencesTrustedAggregatorMultiProof // Event containing the contract specifics and raw log

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
func (it *FeijoapolygonrollupmanagerVerifySequencesTrustedAggregatorMultiProofIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeijoapolygonrollupmanagerVerifySequencesTrustedAggregatorMultiProof)
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
		it.Event = new(FeijoapolygonrollupmanagerVerifySequencesTrustedAggregatorMultiProof)
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
func (it *FeijoapolygonrollupmanagerVerifySequencesTrustedAggregatorMultiProofIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeijoapolygonrollupmanagerVerifySequencesTrustedAggregatorMultiProofIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeijoapolygonrollupmanagerVerifySequencesTrustedAggregatorMultiProof represents a VerifySequencesTrustedAggregatorMultiProof event raised by the Feijoapolygonrollupmanager contract.
type FeijoapolygonrollupmanagerVerifySequencesTrustedAggregatorMultiProof struct {
	Aggregator common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterVerifySequencesTrustedAggregatorMultiProof is a free log retrieval operation binding the contract event 0x97437d34f2cd0d38f9d9399c49bec20084acb988d68397d2629aa8316cacd4f1.
//
// Solidity: event VerifySequencesTrustedAggregatorMultiProof(address indexed aggregator)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) FilterVerifySequencesTrustedAggregatorMultiProof(opts *bind.FilterOpts, aggregator []common.Address) (*FeijoapolygonrollupmanagerVerifySequencesTrustedAggregatorMultiProofIterator, error) {

	var aggregatorRule []interface{}
	for _, aggregatorItem := range aggregator {
		aggregatorRule = append(aggregatorRule, aggregatorItem)
	}

	logs, sub, err := _Feijoapolygonrollupmanager.contract.FilterLogs(opts, "VerifySequencesTrustedAggregatorMultiProof", aggregatorRule)
	if err != nil {
		return nil, err
	}
	return &FeijoapolygonrollupmanagerVerifySequencesTrustedAggregatorMultiProofIterator{contract: _Feijoapolygonrollupmanager.contract, event: "VerifySequencesTrustedAggregatorMultiProof", logs: logs, sub: sub}, nil
}

// WatchVerifySequencesTrustedAggregatorMultiProof is a free log subscription operation binding the contract event 0x97437d34f2cd0d38f9d9399c49bec20084acb988d68397d2629aa8316cacd4f1.
//
// Solidity: event VerifySequencesTrustedAggregatorMultiProof(address indexed aggregator)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) WatchVerifySequencesTrustedAggregatorMultiProof(opts *bind.WatchOpts, sink chan<- *FeijoapolygonrollupmanagerVerifySequencesTrustedAggregatorMultiProof, aggregator []common.Address) (event.Subscription, error) {

	var aggregatorRule []interface{}
	for _, aggregatorItem := range aggregator {
		aggregatorRule = append(aggregatorRule, aggregatorItem)
	}

	logs, sub, err := _Feijoapolygonrollupmanager.contract.WatchLogs(opts, "VerifySequencesTrustedAggregatorMultiProof", aggregatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeijoapolygonrollupmanagerVerifySequencesTrustedAggregatorMultiProof)
				if err := _Feijoapolygonrollupmanager.contract.UnpackLog(event, "VerifySequencesTrustedAggregatorMultiProof", log); err != nil {
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

// ParseVerifySequencesTrustedAggregatorMultiProof is a log parse operation binding the contract event 0x97437d34f2cd0d38f9d9399c49bec20084acb988d68397d2629aa8316cacd4f1.
//
// Solidity: event VerifySequencesTrustedAggregatorMultiProof(address indexed aggregator)
func (_Feijoapolygonrollupmanager *FeijoapolygonrollupmanagerFilterer) ParseVerifySequencesTrustedAggregatorMultiProof(log types.Log) (*FeijoapolygonrollupmanagerVerifySequencesTrustedAggregatorMultiProof, error) {
	event := new(FeijoapolygonrollupmanagerVerifySequencesTrustedAggregatorMultiProof)
	if err := _Feijoapolygonrollupmanager.contract.UnpackLog(event, "VerifySequencesTrustedAggregatorMultiProof", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
