// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package mockfeijoapolygonrollupmanager

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

// LegacyZKEVMStateVariablesPendingState is an auto generated low-level Go binding around an user-defined struct.
type LegacyZKEVMStateVariablesPendingState struct {
	Timestamp         uint64
	LastVerifiedBatch uint64
	ExitRoot          [32]byte
	StateRoot         [32]byte
}

// LegacyZKEVMStateVariablesSequencedBatchData is an auto generated low-level Go binding around an user-defined struct.
type LegacyZKEVMStateVariablesSequencedBatchData struct {
	AccInputHash               [32]byte
	SequencedTimestamp         uint64
	PreviousLastBatchSequenced uint64
}

// MockfeijoapolygonrollupmanagerMetaData contains all meta data concerning the Mockfeijoapolygonrollupmanager contract.
var MockfeijoapolygonrollupmanagerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIPolygonZkEVMGlobalExitRootV2\",\"name\":\"_globalExitRootManager\",\"type\":\"address\"},{\"internalType\":\"contractIERC20Upgradeable\",\"name\":\"_pol\",\"type\":\"address\"},{\"internalType\":\"contractIPolygonZkEVMBridge\",\"name\":\"_bridgeAddress\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessControlOnlyCanRenounceRolesForSelf\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AddressDoNotHaveRequiredRole\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"AllzkEVMSequencedBatchesMustBeVerified\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"BatchFeeOutOfRange\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ChainIDAlreadyExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ExceedMaxVerifyBatches\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FinalNumBatchBelowLastVerifiedBatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FinalNumBatchDoesNotMatchPendingState\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"FinalPendingStateNumInvalid\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"HaltTimeoutNotExpired\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InitBatchMustMatchCurrentForkID\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InitNumBatchAboveLastVerifiedBatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InitNumBatchDoesNotMatchPendingState\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidProof\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRangeBatchTimeTarget\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidRangeMultiplierBatchFee\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustSequenceSomeBatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NewAccInputHashDoesNotExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NewPendingStateTimeoutMustBeLower\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NewStateRootNotInsidePrime\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NewTrustedAggregatorTimeoutMustBeLower\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OldAccInputHashDoesNotExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OldStateRootDoesNotExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyEmergencyState\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"OnlyNotEmergencyState\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PendingStateDoesNotExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PendingStateInvalid\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PendingStateNotConsolidable\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RollupAddressAlreadyExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RollupMustExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RollupTypeDoesNotExist\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"RollupTypeObsolete\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SenderMustBeRollup\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"StoredRootMustBeDifferentThanNewRoot\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TrustedAggregatorTimeoutNotExpired\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpdateNotCompatible\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UpdateToSameRollupTypeID\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"forkID\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"rollupAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"chainID\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"rollupCompatibilityID\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"lastVerifiedBatchBeforeUpgrade\",\"type\":\"uint64\"}],\"name\":\"AddExistingRollup\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"rollupTypeID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"consensusImplementation\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"verifier\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"forkID\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"rollupCompatibilityID\",\"type\":\"uint8\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"genesis\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"string\",\"name\":\"description\",\"type\":\"string\"}],\"name\":\"AddNewRollupType\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"numBatch\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"exitRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"pendingStateNum\",\"type\":\"uint64\"}],\"name\":\"ConsolidatePendingState\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"rollupTypeID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"rollupAddress\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"chainID\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"gasTokenAddress\",\"type\":\"address\"}],\"name\":\"CreateNewRollup\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"EmergencyStateActivated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"EmergencyStateDeactivated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"rollupTypeID\",\"type\":\"uint32\"}],\"name\":\"ObsoleteRollupType\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"lastBatchSequenced\",\"type\":\"uint64\"}],\"name\":\"OnSequenceBatches\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"numBatch\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"exitRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"aggregator\",\"type\":\"address\"}],\"name\":\"OverridePendingState\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"storedStateRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"provedStateRoot\",\"type\":\"bytes32\"}],\"name\":\"ProveNonDeterministicPendingState\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newBatchFee\",\"type\":\"uint256\"}],\"name\":\"SetBatchFee\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"newMultiplierBatchFee\",\"type\":\"uint16\"}],\"name\":\"SetMultiplierBatchFee\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"newPendingStateTimeout\",\"type\":\"uint64\"}],\"name\":\"SetPendingStateTimeout\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newTrustedAggregator\",\"type\":\"address\"}],\"name\":\"SetTrustedAggregator\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"newTrustedAggregatorTimeout\",\"type\":\"uint64\"}],\"name\":\"SetTrustedAggregatorTimeout\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"newVerifyBatchTimeTarget\",\"type\":\"uint64\"}],\"name\":\"SetVerifyBatchTimeTarget\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"newRollupTypeID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"lastVerifiedBatchBeforeUpgrade\",\"type\":\"uint64\"}],\"name\":\"UpdateRollup\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"numBatch\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"exitRoot\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"aggregator\",\"type\":\"address\"}],\"name\":\"VerifyBatches\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"numBatch\",\"type\":\"uint64\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"exitRoot\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"aggregator\",\"type\":\"address\"}],\"name\":\"VerifyBatchesTrustedAggregator\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"activateEmergencyState\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractIPolygonRollupBase\",\"name\":\"rollupAddress\",\"type\":\"address\"},{\"internalType\":\"contractIVerifierRollup\",\"name\":\"verifier\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"forkID\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"chainID\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"genesis\",\"type\":\"bytes32\"},{\"internalType\":\"uint8\",\"name\":\"rollupCompatibilityID\",\"type\":\"uint8\"}],\"name\":\"addExistingRollup\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"consensusImplementation\",\"type\":\"address\"},{\"internalType\":\"contractIVerifierRollup\",\"name\":\"verifier\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"forkID\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"rollupCompatibilityID\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"genesis\",\"type\":\"bytes32\"},{\"internalType\":\"string\",\"name\":\"description\",\"type\":\"string\"}],\"name\":\"addNewRollupType\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"bridgeAddress\",\"outputs\":[{\"internalType\":\"contractIPolygonZkEVMBridge\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"calculateRewardPerBatch\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"chainID\",\"type\":\"uint64\"}],\"name\":\"chainIDToRollupID\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"pendingStateNum\",\"type\":\"uint64\"}],\"name\":\"consolidatePendingState\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupTypeID\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"chainID\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"sequencer\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"gasTokenAddress\",\"type\":\"address\"},{\"internalType\":\"string\",\"name\":\"sequencerURL\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"networkName\",\"type\":\"string\"}],\"name\":\"createNewRollup\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"deactivateEmergencyState\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBatchFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getForcedBatchFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"initNumBatch\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"finalNewBatch\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"newLocalExitRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"oldStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"}],\"name\":\"getInputSnarkBytes\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"}],\"name\":\"getLastVerifiedBatch\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"batchNum\",\"type\":\"uint64\"}],\"name\":\"getRollupBatchNumToStateRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getRollupExitRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"batchNum\",\"type\":\"uint64\"}],\"name\":\"getRollupPendingStateTransitions\",\"outputs\":[{\"components\":[{\"internalType\":\"uint64\",\"name\":\"timestamp\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"lastVerifiedBatch\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"exitRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"}],\"internalType\":\"structLegacyZKEVMStateVariables.PendingState\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"batchNum\",\"type\":\"uint64\"}],\"name\":\"getRollupSequencedBatches\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"accInputHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"sequencedTimestamp\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"previousLastBatchSequenced\",\"type\":\"uint64\"}],\"internalType\":\"structLegacyZKEVMStateVariables.SequencedBatchData\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"globalExitRootManager\",\"outputs\":[{\"internalType\":\"contractIPolygonZkEVMGlobalExitRootV2\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"trustedAggregator\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"_pendingStateTimeout\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"_trustedAggregatorTimeout\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"timelock\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"emergencyCouncil\",\"type\":\"address\"},{\"internalType\":\"contractPolygonZkEVMExistentEtrog\",\"name\":\"polygonZkEVM\",\"type\":\"address\"},{\"internalType\":\"contractIVerifierRollup\",\"name\":\"zkEVMVerifier\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"zkEVMForkID\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"zkEVMChainID\",\"type\":\"uint64\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"isEmergencyState\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"pendingStateNum\",\"type\":\"uint64\"}],\"name\":\"isPendingStateConsolidable\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastAggregationTimestamp\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"lastDeactivatedEmergencyStateTimestamp\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"multiplierBatchFee\",\"outputs\":[{\"internalType\":\"uint16\",\"name\":\"\",\"type\":\"uint16\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupTypeID\",\"type\":\"uint32\"}],\"name\":\"obsoleteRollupType\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"newSequencedBatches\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"newAccInputHash\",\"type\":\"bytes32\"}],\"name\":\"onSequenceBatches\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"initPendingStateNum\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"finalPendingStateNum\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"initNumBatch\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"finalNewBatch\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"newLocalExitRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[24]\",\"name\":\"proof\",\"type\":\"bytes32[24]\"}],\"name\":\"overridePendingState\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pendingStateTimeout\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pol\",\"outputs\":[{\"internalType\":\"contractIERC20Upgradeable\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"initPendingStateNum\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"finalPendingStateNum\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"initNumBatch\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"finalNewBatch\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"newLocalExitRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32[24]\",\"name\":\"proof\",\"type\":\"bytes32[24]\"}],\"name\":\"proveNonDeterministicPendingState\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"rollupAddress\",\"type\":\"address\"}],\"name\":\"rollupAddressToID\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rollupCount\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"}],\"name\":\"rollupIDToRollupData\",\"outputs\":[{\"internalType\":\"contractIPolygonRollupBase\",\"name\":\"rollupContract\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"chainID\",\"type\":\"uint64\"},{\"internalType\":\"contractIVerifierRollup\",\"name\":\"verifier\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"forkID\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"lastLocalExitRoot\",\"type\":\"bytes32\"},{\"internalType\":\"uint64\",\"name\":\"lastBatchSequenced\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"lastVerifiedBatch\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"lastPendingState\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"lastPendingStateConsolidated\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"lastVerifiedBatchBeforeUpgrade\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"rollupTypeID\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"rollupCompatibilityID\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rollupTypeCount\",\"outputs\":[{\"internalType\":\"uint32\",\"name\":\"\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupTypeID\",\"type\":\"uint32\"}],\"name\":\"rollupTypeMap\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"consensusImplementation\",\"type\":\"address\"},{\"internalType\":\"contractIVerifierRollup\",\"name\":\"verifier\",\"type\":\"address\"},{\"internalType\":\"uint64\",\"name\":\"forkID\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"rollupCompatibilityID\",\"type\":\"uint8\"},{\"internalType\":\"bool\",\"name\":\"obsolete\",\"type\":\"bool\"},{\"internalType\":\"bytes32\",\"name\":\"genesis\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"newBatchFee\",\"type\":\"uint256\"}],\"name\":\"setBatchFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint16\",\"name\":\"newMultiplierBatchFee\",\"type\":\"uint16\"}],\"name\":\"setMultiplierBatchFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"newPendingStateTimeout\",\"type\":\"uint64\"}],\"name\":\"setPendingStateTimeout\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"newTrustedAggregatorTimeout\",\"type\":\"uint64\"}],\"name\":\"setTrustedAggregatorTimeout\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"newVerifyBatchTimeTarget\",\"type\":\"uint64\"}],\"name\":\"setVerifyBatchTimeTarget\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSequencedBatches\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalVerifiedBatches\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"trustedAggregatorTimeout\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"contractITransparentUpgradeableProxy\",\"name\":\"rollupContract\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"newRollupTypeID\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"upgradeData\",\"type\":\"bytes\"}],\"name\":\"updateRollup\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"verifyBatchTimeTarget\",\"outputs\":[{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"pendingStateNum\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"initNumBatch\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"finalNewBatch\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"newLocalExitRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"beneficiary\",\"type\":\"address\"},{\"internalType\":\"bytes32[24]\",\"name\":\"proof\",\"type\":\"bytes32[24]\"}],\"name\":\"verifyBatches\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"rollupID\",\"type\":\"uint32\"},{\"internalType\":\"uint64\",\"name\":\"pendingStateNum\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"initNumBatch\",\"type\":\"uint64\"},{\"internalType\":\"uint64\",\"name\":\"finalNewBatch\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"newLocalExitRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"newStateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"beneficiary\",\"type\":\"address\"},{\"internalType\":\"bytes32[24]\",\"name\":\"proof\",\"type\":\"bytes32[24]\"}],\"name\":\"verifyBatchesTrustedAggregator\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x60e060405234801562000010575f80fd5b5060405162006fae38038062006fae83398101604081905262000033916200013c565b6001600160a01b0380841660805280831660c052811660a0528282826200005962000065565b5050505050506200018d565b5f54610100900460ff1615620000d15760405162461bcd60e51b815260206004820152602760248201527f496e697469616c697a61626c653a20636f6e747261637420697320696e697469604482015266616c697a696e6760c81b606482015260840160405180910390fd5b5f5460ff908116101562000122575f805460ff191660ff9081179091556040519081527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b565b6001600160a01b038116811462000139575f80fd5b50565b5f805f606084860312156200014f575f80fd5b83516200015c8162000124565b60208501519093506200016f8162000124565b6040850151909250620001828162000124565b809150509250925092565b60805160a05160c051616dbe620001f05f395f8181610b8a0152818161273d015261444d01525f818161090f015281816135c5015261481e01525f8181610add01528181611358015281816115780152818161242b01526146f80152616dbe5ff3fe608060405234801562000010575f80fd5b50600436106200038c575f3560e01c8063841b24d711620001e3578063c1acbc341162000113578063dbc1697611620000ab578063e46761c41162000083578063e46761c41462000b84578063f34eb8eb1462000bac578063f4e926751462000bc3578063f9c4c2ae1462000bd4575f80fd5b8063dbc169761462000b42578063dde0ff771462000b4c578063e0bfd3d21462000b6d575f80fd5b8063d02103ca11620000eb578063d02103ca1462000ad7578063d5073f6f1462000aff578063d547741f1462000b16578063d939b3151462000b2d575f80fd5b8063c1acbc341462000a6f578063c4c928c21462000a98578063ceee281d1462000aaf575f80fd5b80639c9f3dfe1162000187578063a2967d99116200015f578063a2967d9914620008ff578063a3c573eb1462000909578063afd23cbe1462000957578063b99d0ad7146200098e575f80fd5b80639c9f3dfe14620008c9578063a066215c14620008e0578063a217fddf14620008f7575f80fd5b806391d1485411620001bb57806391d14854146200086057806399f5634e14620008a85780639a908e7314620008b2575f80fd5b8063841b24d7146200080157806387c20c0114620008325780638bd4f0711462000849575f80fd5b80632528016911620002bf57806355a71ee011620002635780637222020f116200023b5780637222020f1462000785578063727885e9146200079c5780637975fcfe14620007b35780637fb6e76a14620007d9575f80fd5b806355a71ee0146200062957806360469169146200066d57806365c0504d1462000677575f80fd5b806336568abe116200029757806336568abe14620005f2578063394218e91462000609578063477fa2701462000620575f80fd5b806325280169146200050b5780632f2ff15d14620005c657806330c27dde14620005dd575f80fd5b80631489ed1011620003335780631796a1ae116200030b5780631796a1ae146200048f5780631816b7e514620004b65780632072f6c514620004cd578063248a9ca314620004d7575f80fd5b80631489ed10146200045357806315064c96146200046a5780631608859c1462000478575f80fd5b80630a0d9fbe11620003675780630a0d9fbe146200040457806311f6b287146200042557806312b86e19146200043c575f80fd5b80630645af091462000390578063066ec01214620003a9578063080b311114620003dc575b5f80fd5b620003a7620003a136600462005780565b62000d3e565b005b608454620003be9067ffffffffffffffff1681565b60405167ffffffffffffffff90911681526020015b60405180910390f35b620003f3620003ed36600462005869565b62001224565b6040519015158152602001620003d3565b608554620003be9068010000000000000000900467ffffffffffffffff1681565b620003be620004363660046200589f565b6200124d565b620003a76200044d366004620058cd565b6200126c565b620003a7620004643660046200595f565b6200148c565b606f54620003f39060ff1681565b620003a76200048936600462005869565b6200166b565b607e54620004a09063ffffffff1681565b60405163ffffffff9091168152602001620003d3565b620003a7620004c7366004620059e4565b62001747565b620003a762001844565b620004fc620004e836600462005a0e565b5f9081526034602052604090206001015490565b604051908152602001620003d3565b620005916200051c36600462005869565b60408051606080820183525f808352602080840182905292840181905263ffffffff9590951685526081825282852067ffffffffffffffff9485168652600301825293829020825194850183528054855260010154808416918501919091526801000000000000000090049091169082015290565b604080518251815260208084015167ffffffffffffffff908116918301919091529282015190921690820152606001620003d3565b620003a7620005d736600462005a26565b62001958565b608754620003be9067ffffffffffffffff1681565b620003a76200060336600462005a26565b62001980565b620003a76200061a36600462005a57565b620019e0565b608654620004fc565b620004fc6200063a36600462005869565b63ffffffff82165f90815260816020908152604080832067ffffffffffffffff8516845260020190915290205492915050565b620004fc62001afb565b6200072d620006883660046200589f565b607f6020525f908152604090208054600182015460029092015473ffffffffffffffffffffffffffffffffffffffff918216929182169167ffffffffffffffff740100000000000000000000000000000000000000008204169160ff7c010000000000000000000000000000000000000000000000000000000083048116927d0100000000000000000000000000000000000000000000000000000000009004169086565b6040805173ffffffffffffffffffffffffffffffffffffffff978816815296909516602087015267ffffffffffffffff9093169385019390935260ff166060840152901515608083015260a082015260c001620003d3565b620003a7620007963660046200589f565b62001b12565b620003a7620007ad36600462005b4f565b62001c8d565b620007ca620007c436600462005c16565b620021f6565b604051620003d3919062005ce7565b620004a0620007ea36600462005a57565b60836020525f908152604090205463ffffffff1681565b608454620003be907801000000000000000000000000000000000000000000000000900467ffffffffffffffff1681565b620003a7620008433660046200595f565b62002228565b620003a76200085a366004620058cd565b62002625565b620003f36200087136600462005a26565b5f91825260346020908152604080842073ffffffffffffffffffffffffffffffffffffffff93909316845291905290205460ff1690565b620004fc620026f6565b620003be620008c336600462005cfb565b62002804565b620003a7620008da36600462005a57565b62002a54565b620003a7620008f136600462005a57565b62002b3d565b620004fc5f81565b620004fc62002c27565b620009317f000000000000000000000000000000000000000000000000000000000000000081565b60405173ffffffffffffffffffffffffffffffffffffffff9091168152602001620003d3565b6085546200097a90700100000000000000000000000000000000900461ffff1681565b60405161ffff9091168152602001620003d3565b62000a296200099f36600462005869565b60408051608080820183525f8083526020808401829052838501829052606093840182905263ffffffff9690961681526081865283812067ffffffffffffffff958616825260040186528390208351918201845280548086168352680100000000000000009004909416948101949094526001830154918401919091526002909101549082015290565b604051620003d391905f60808201905067ffffffffffffffff80845116835280602085015116602084015250604083015160408301526060830151606083015292915050565b608454620003be90700100000000000000000000000000000000900467ffffffffffffffff1681565b620003a762000aa936600462005d26565b6200300f565b620004a062000ac036600462005db8565b60826020525f908152604090205463ffffffff1681565b620009317f000000000000000000000000000000000000000000000000000000000000000081565b620003a762000b1036600462005a0e565b6200345d565b620003a762000b2736600462005a26565b62003512565b608554620003be9067ffffffffffffffff1681565b620003a76200353a565b608454620003be9068010000000000000000900467ffffffffffffffff1681565b620003a762000b7e36600462005de7565b62003644565b620009317f000000000000000000000000000000000000000000000000000000000000000081565b620003a762000bbd36600462005e5f565b62003758565b608054620004a09063ffffffff1681565b62000cb062000be53660046200589f565b60816020525f90815260409020805460018201546005830154600684015460079094015473ffffffffffffffffffffffffffffffffffffffff80851695740100000000000000000000000000000000000000009586900467ffffffffffffffff908116969286169592909204821693928282169268010000000000000000808404821693700100000000000000000000000000000000808204841694780100000000000000000000000000000000000000000000000090920484169380831693830416910460ff168c565b6040805173ffffffffffffffffffffffffffffffffffffffff9d8e16815267ffffffffffffffff9c8d1660208201529c909a16998c019990995296891660608b015260808a019590955292871660a089015290861660c0880152851660e0870152841661010086015283166101208501529190911661014083015260ff1661016082015261018001620003d3565b5f54600290610100900460ff1615801562000d5f57505f5460ff8083169116105b62000df1576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201527f647920696e697469616c697a656400000000000000000000000000000000000060648201526084015b60405180910390fd5b5f805461010060ff84167fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000090921691909117179055608580546084805477ffffffffffffffffffffffffffffffffffffffffffffffff16780100000000000000000000000000000000000000000000000067ffffffffffffffff8e8116919091029190911790915567016345785d8a00006086558c167fffffffffffffffffffffffffffffffff00000000000000000000000000000000909116176907080000000000000000177fffffffffffffffffffffffffffff0000ffffffffffffffffffffffffffffffff167103ea0000000000000000000000000000000017905562000efa6200399d565b62000f267f084e94f375e9d647f87f5b2ceffba1e062c70f6009fdbcf80291e803b5c9edd48c62003a35565b62000f325f8862003a35565b62000f5e7fac75d24dbb35ea80e25fab167da4dea46c1915260426570db84f184891f5f5908862003a35565b62000f8a7f3dfe277d2a2c04b75fb2eb3743fa00005ae3678a20c299e65fdf4df76517f68e8862003a35565b62000fb67f66156603fe29d13f97c6f3e3dff4ef71919f9aa61c555be0182d954e94221aac8862003a35565b62000fe27fab66e11c4f712cd06ab11bf9339b48bef39e12d4a22eeef71d2860a0c90482bd8962003a35565b6200100e7fa0fab074aba36a6fa69f1a83ee86e5abfb8433966eb57efb13dc2fc2f24ddd088962003a35565b6200103a7f62ba6ba2ffed8cfe316b583325ea41ac6e7ba9e5864d2bc6fabba7ac26d2f0f48962003a35565b620010667fa5c5790f581d443ed43873ab47cfb8c5d66a6db268e58b5971bb33fc66e07db18962003a35565b620010b27f084e94f375e9d647f87f5b2ceffba1e062c70f6009fdbcf80291e803b5c9edd47f73cb0569fdbea2544dae03fdb2fe10eda92a72a2e8cd2bd496e85b762505a3f062003a41565b620010de7f73cb0569fdbea2544dae03fdb2fe10eda92a72a2e8cd2bd496e85b762505a3f08962003a35565b6200110a7f8cf807f6970720f8e2c208c7c5037595982c7bd9ed93c380d09df743d0dcc3fb8962003a35565b620011567f141f8f32ce6198eee741f695cec728bfd32d289f1acf73621fb303581000545e7f9b6f082d8d3644ae2f24a3c32e356d6f2d9b2844d9b26164fbc82663ff28595162003a41565b620011827f141f8f32ce6198eee741f695cec728bfd32d289f1acf73621fb303581000545e8762003a35565b620011ae7f9b6f082d8d3644ae2f24a3c32e356d6f2d9b2844d9b26164fbc82663ff2859518762003a35565b620011ba5f3362003a35565b5f80547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff00ff16905560405160ff821681527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15050505050505050505050565b63ffffffff82165f90815260816020526040812062001244908362003a8b565b90505b92915050565b63ffffffff81165f908152608160205260408120620012479062003ad1565b7f084e94f375e9d647f87f5b2ceffba1e062c70f6009fdbcf80291e803b5c9edd4620012988162003b67565b63ffffffff89165f908152608160205260409020620012be818a8a8a8a8a8a8a62003b73565b6006810180547fffffffffffffffffffffffffffffffff0000000000000000ffffffffffffffff166801000000000000000067ffffffffffffffff8981169182029290921783555f90815260028401602052604090208690556005830187905590547001000000000000000000000000000000009004161562001356576006810180546fffffffffffffffffffffffffffffffff1690555b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166333d6247d6200139c62002c27565b6040518263ffffffff1660e01b8152600401620013bb91815260200190565b5f604051808303815f87803b158015620013d3575f80fd5b505af1158015620013e6573d5f803e3d5ffd5b50506084805477ffffffffffffffffffffffffffffffffffffffffffffffff167a093a8000000000000000000000000000000000000000000000000017905550506040805167ffffffffffffffff881681526020810186905290810186905233606082015263ffffffff8b16907f3182bd6e6f74fc1fdc88b60f3a4f4c7f79db6ae6f5b88a1b3f5a1e28ec210d5e9060800160405180910390a250505050505050505050565b7f084e94f375e9d647f87f5b2ceffba1e062c70f6009fdbcf80291e803b5c9edd4620014b88162003b67565b63ffffffff89165f908152608160205260409020620014de818a8a8a8a8a8a8a62004052565b6006810180547fffffffffffffffffffffffffffffffff0000000000000000ffffffffffffffff166801000000000000000067ffffffffffffffff8a81169182029290921783555f90815260028401602052604090208790556005830188905590547001000000000000000000000000000000009004161562001576576006810180546fffffffffffffffffffffffffffffffff1690555b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166333d6247d620015bc62002c27565b6040518263ffffffff1660e01b8152600401620015db91815260200190565b5f604051808303815f87803b158015620015f3575f80fd5b505af115801562001606573d5f803e3d5ffd5b50506040805167ffffffffffffffff8b1681526020810189905290810189905233925063ffffffff8d1691507fd1ec3a1216f08b6eff72e169ceb548b782db18a6614852618d86bb19f3f9b0d39060600160405180910390a350505050505050505050565b63ffffffff82165f9081526081602090815260408083203384527fc17b14a573f65366cdad721c7c0a0f76536bb4a86b935cdac44610e4f010b52a9092529091205460ff166200173657606f5460ff1615620016f3576040517f2f0047fc00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b620016ff818362003a8b565b62001736576040517f0ce9e4a200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b62001742818362004598565b505050565b7fa5c5790f581d443ed43873ab47cfb8c5d66a6db268e58b5971bb33fc66e07db1620017738162003b67565b6103e88261ffff1610806200178d57506103ff8261ffff16115b15620017c5576040517f4c2533c800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b608580547fffffffffffffffffffffffffffff0000ffffffffffffffffffffffffffffffff1670010000000000000000000000000000000061ffff8516908102919091179091556040519081527f7019933d795eba185c180209e8ae8bffbaa25bcef293364687702c31f4d302c5906020015b60405180910390a15050565b335f9081527f8875b94af5657a2903def9906d67a3f42d8a836d24b5602c00f00fc855339fcd602052604090205460ff166200194c57608454700100000000000000000000000000000000900467ffffffffffffffff161580620018e257506084544290620018d69062093a8090700100000000000000000000000000000000900467ffffffffffffffff1662005f1f565b67ffffffffffffffff16115b806200191457506087544290620019089062093a809067ffffffffffffffff1662005f1f565b67ffffffffffffffff16115b156200194c576040517fd257555a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b620019566200481c565b565b5f82815260346020526040902060010154620019748162003b67565b620017428383620048a3565b73ffffffffffffffffffffffffffffffffffffffff81163314620019d0576040517f5a568e6800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b620019dc82826200495f565b5050565b7fa5c5790f581d443ed43873ab47cfb8c5d66a6db268e58b5971bb33fc66e07db162001a0c8162003b67565b606f5460ff1662001a7d5760845467ffffffffffffffff780100000000000000000000000000000000000000000000000090910481169083161062001a7d576040517f401636df00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6084805477ffffffffffffffffffffffffffffffffffffffffffffffff16780100000000000000000000000000000000000000000000000067ffffffffffffffff8516908102919091179091556040519081527f1f4fa24c2e4bad19a7f3ec5c5485f70d46c798461c2e684f55bbd0fc661373a19060200162001838565b5f608654606462001b0d919062005f4a565b905090565b7fab66e11c4f712cd06ab11bf9339b48bef39e12d4a22eeef71d2860a0c90482bd62001b3e8162003b67565b63ffffffff8216158062001b5d5750607e5463ffffffff908116908316115b1562001b95576040517f7512e5cb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b63ffffffff82165f908152607f602052604090206001808201547d010000000000000000000000000000000000000000000000000000000000900460ff161515900362001c0e576040517f3b8d3d9900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001810180547fffff00ffffffffffffffffffffffffffffffffffffffffffffffffffffffffff167d01000000000000000000000000000000000000000000000000000000000017905560405163ffffffff8416907f4710d2ee567ef1ed6eb2f651dde4589524bcf7cebc62147a99b281cc836e7e44905f90a2505050565b7fa0fab074aba36a6fa69f1a83ee86e5abfb8433966eb57efb13dc2fc2f24ddd0862001cb98162003b67565b63ffffffff8816158062001cd85750607e5463ffffffff908116908916115b1562001d10576040517f7512e5cb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b63ffffffff88165f908152607f602052604090206001808201547d010000000000000000000000000000000000000000000000000000000000900460ff161515900362001d89576040517f3b8d3d9900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff88165f9081526083602052604090205463ffffffff161562001de0576040517f6f91fc1200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b608080545f9190829062001dfa9063ffffffff1662005f64565b825463ffffffff8281166101009490940a93840293021916919091179091558254604080515f808252602082019283905293945073ffffffffffffffffffffffffffffffffffffffff90921691309162001e549062005738565b62001e629392919062005f89565b604051809103905ff08015801562001e7c573d5f803e3d5ffd5b5090508160835f8c67ffffffffffffffff1667ffffffffffffffff1681526020019081526020015f205f6101000a81548163ffffffff021916908363ffffffff1602179055508160825f8373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548163ffffffff021916908363ffffffff1602179055505f60815f8463ffffffff1663ffffffff1681526020019081526020015f20905081815f015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508360010160149054906101000a900467ffffffffffffffff168160010160146101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550836001015f9054906101000a900473ffffffffffffffffffffffffffffffffffffffff16816001015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055508a815f0160146101000a81548167ffffffffffffffff021916908367ffffffffffffffff1602179055508360020154816002015f8067ffffffffffffffff1681526020019081526020015f20819055508b63ffffffff168160070160086101000a81548167ffffffffffffffff021916908367ffffffffffffffff16021790555083600101601c9054906101000a900460ff168160070160106101000a81548160ff021916908360ff1602179055508263ffffffff167f194c983456df6701c6a50830b90fe80e72b823411d0d524970c9590dc277a6418d848e8c60405162002153949392919063ffffffff94909416845273ffffffffffffffffffffffffffffffffffffffff928316602085015267ffffffffffffffff91909116604084015216606082015260800190565b60405180910390a26040517f7125702200000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff831690637125702290620021b9908d908d9088908e908e908e9060040162005fcc565b5f604051808303815f87803b158015620021d1575f80fd5b505af1158015620021e4573d5f803e3d5ffd5b50505050505050505050505050505050565b63ffffffff86165f9081526081602052604090206060906200221d90878787878762004a19565b979650505050505050565b606f5460ff161562002266576040517f2f0047fc00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b63ffffffff88165f90815260816020908152604080832060845467ffffffffffffffff8a811686526003830190945291909320600101544292620022ca92780100000000000000000000000000000000000000000000000090048116911662005f1f565b67ffffffffffffffff1611156200230d576040517f8a0704d300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6103e86200231c88886200603b565b67ffffffffffffffff1611156200235f576040517fb59f753a00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b62002371818989898989898962004052565b6200237d818762004be0565b60855467ffffffffffffffff165f03620024c3576006810180547fffffffffffffffffffffffffffffffff0000000000000000ffffffffffffffff166801000000000000000067ffffffffffffffff8981169182029290921783555f90815260028401602052604090208690556005830187905590547001000000000000000000000000000000009004161562002429576006810180546fffffffffffffffffffffffffffffffff1690555b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166333d6247d6200246f62002c27565b6040518263ffffffff1660e01b81526004016200248e91815260200190565b5f604051808303815f87803b158015620024a6575f80fd5b505af1158015620024b9573d5f803e3d5ffd5b50505050620025c5565b620024ce8162004de7565b600681018054700100000000000000000000000000000000900467ffffffffffffffff1690601062002500836200605f565b825467ffffffffffffffff9182166101009390930a92830292820219169190911790915560408051608081018252428316815289831660208083019182528284018b8152606084018b81526006890154700100000000000000000000000000000000900487165f90815260048a0190935294909120925183549251861668010000000000000000027fffffffffffffffffffffffffffffffff000000000000000000000000000000009093169516949094171781559151600183015551600290910155505b6040805167ffffffffffffffff8816815260208101869052908101869052339063ffffffff8b16907faac1e7a157b259544ebacd6e8a82ae5d6c8f174e12aa48696277bcc9a661f0b49060600160405180910390a3505050505050505050565b606f5460ff161562002663576040517f2f0047fc00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b63ffffffff88165f90815260816020526040902062002689818989898989898962003b73565b67ffffffffffffffff87165f9081526004820160209081526040918290206002015482519081529081018590527f1f44c21118c4603cfb4e1b621dbcfa2b73efcececee2b99b620b2953d33a7010910160405180910390a1620026eb6200481c565b505050505050505050565b6040517f70a082310000000000000000000000000000000000000000000000000000000081523060048201525f90819073ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016906370a0823190602401602060405180830381865afa15801562002783573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190620027a991906200607e565b6084549091505f90620027d59067ffffffffffffffff680100000000000000008204811691166200603b565b67ffffffffffffffff169050805f03620027f1575f9250505090565b620027fd8183620060c3565b9250505090565b606f545f9060ff161562002844576040517f2f0047fc00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b335f9081526082602052604081205463ffffffff169081900362002894576040517f71653c1500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8367ffffffffffffffff165f03620028d8576040517f2590ccf900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b63ffffffff81165f90815260816020526040812060848054919287926200290b90849067ffffffffffffffff1662005f1f565b82546101009290920a67ffffffffffffffff81810219909316918316021790915560068301541690505f62002941878362005f1f565b60068401805467ffffffffffffffff8084167fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000009092168217909255604080516060810182528a815242841660208083019182528886168385019081525f95865260038b019091529290932090518155915160019290920180549151841668010000000000000000027fffffffffffffffffffffffffffffffff000000000000000000000000000000009092169290931691909117179055905062002a058362004de7565b60405167ffffffffffffffff8216815263ffffffff8516907f1d9f30260051d51d70339da239ea7b080021adcaabfa71c9b0ea339a20cf9a259060200160405180910390a29695505050505050565b7fa5c5790f581d443ed43873ab47cfb8c5d66a6db268e58b5971bb33fc66e07db162002a808162003b67565b606f5460ff1662002ad55760855467ffffffffffffffff9081169083161062002ad5576040517f48a05a9000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b608580547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001667ffffffffffffffff84169081179091556040519081527fc4121f4e22c69632ebb7cf1f462be0511dc034f999b52013eddfb24aab765c759060200162001838565b7fa5c5790f581d443ed43873ab47cfb8c5d66a6db268e58b5971bb33fc66e07db162002b698162003b67565b620151808267ffffffffffffffff16111562002bb1576040517fe067dfe800000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b608580547fffffffffffffffffffffffffffffffff0000000000000000ffffffffffffffff166801000000000000000067ffffffffffffffff8516908102919091179091556040519081527f1b023231a1ab6b5d93992f168fb44498e1a7e64cef58daff6f1c216de6a68c289060200162001838565b6080545f9063ffffffff1680820362002c4157505f919050565b5f8167ffffffffffffffff81111562002c5e5762002c5e62005a73565b60405190808252806020026020018201604052801562002c88578160200160208202803683370190505b5090505f5b8281101562002cf85760815f62002ca6836001620060d9565b63ffffffff1663ffffffff1681526020019081526020015f206005015482828151811062002cd85762002cd8620060ef565b60209081029190910101528062002cef816200611c565b91505062002c8d565b505f60205b8360011462002f51575f62002d1460028662006156565b62002d21600287620060c3565b62002d2d9190620060d9565b90505f8167ffffffffffffffff81111562002d4c5762002d4c62005a73565b60405190808252806020026020018201604052801562002d76578160200160208202803683370190505b5090505f5b8281101562002efd5762002d916001846200616c565b8114801562002dac575062002da860028862006156565b6001145b1562002e34578562002dc082600262005f4a565b8151811062002dd35762002dd3620060ef565b60200260200101518560405160200162002df7929190918252602082015260400190565b6040516020818303038152906040528051906020012082828151811062002e225762002e22620060ef565b60200260200101818152505062002ee8565b8562002e4282600262005f4a565b8151811062002e555762002e55620060ef565b60200260200101518682600262002e6d919062005f4a565b62002e7a906001620060d9565b8151811062002e8d5762002e8d620060ef565b602002602001015160405160200162002eb0929190918252602082015260400190565b6040516020818303038152906040528051906020012082828151811062002edb5762002edb620060ef565b6020026020010181815250505b8062002ef4816200611c565b91505062002d7b565b50809450819550838460405160200162002f21929190918252602082015260400190565b604051602081830303815290604052805190602001209350828062002f469062006182565b935050505062002cfd565b5f835f8151811062002f675762002f67620060ef565b602002602001015190505f5b8281101562003005576040805160208101849052908101859052606001604080517fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe081840301815282825280516020918201209083018790529082018690529250606001604051602081830303815290604052805190602001209350808062002ffc906200611c565b91505062002f73565b5095945050505050565b7f66156603fe29d13f97c6f3e3dff4ef71919f9aa61c555be0182d954e94221aac6200303b8162003b67565b63ffffffff841615806200305a5750607e5463ffffffff908116908516115b1562003092576040517f7512e5cb00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff85165f9081526082602052604081205463ffffffff1690819003620030f8576040517f74a086a300000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b63ffffffff8181165f9081526081602052604090206007810154909187166801000000000000000090910467ffffffffffffffff160362003165576040517f4f61d51900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b63ffffffff86165f908152607f602052604090206001808201547d010000000000000000000000000000000000000000000000000000000000900460ff1615159003620031de576040517f3b8d3d9900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60018101546007830154700100000000000000000000000000000000900460ff9081167c010000000000000000000000000000000000000000000000000000000090920416146200325b576040517fb541abe200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6001808201805491840180547fffffffffffffffffffffffff0000000000000000000000000000000000000000811673ffffffffffffffffffffffffffffffffffffffff9094169384178255915467ffffffffffffffff740100000000000000000000000000000000000000009182900416027fffffffff000000000000000000000000000000000000000000000000000000009092169092171790556007820180546801000000000000000063ffffffff8a16027fffffffffffffffffffffffffffffffff0000000000000000ffffffffffffffff9091161790555f62003343846200124d565b6007840180547fffffffffffffffffffffffffffffffffffffffffffffffff00000000000000001667ffffffffffffffff831617905582546040517f4f1ef28600000000000000000000000000000000000000000000000000000000815291925073ffffffffffffffffffffffffffffffffffffffff8b811692634f1ef28692620033d79216908b908b90600401620061b9565b5f604051808303815f87803b158015620033ef575f80fd5b505af115801562003402573d5f803e3d5ffd5b50506040805163ffffffff8c8116825267ffffffffffffffff86166020830152881693507ff585e04c05d396901170247783d3e5f0ee9c1df23072985b50af089f5e48b19d92500160405180910390a2505050505050505050565b7f8cf807f6970720f8e2c208c7c5037595982c7bd9ed93c380d09df743d0dcc3fb620034898162003b67565b683635c9adc5dea00000821180620034a45750633b9aca0082105b15620034dc576040517f8586952500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b60868290556040518281527ffb383653f53ee079978d0c9aff7aeff04a10166ce244cca9c9f9d8d96bed45b29060200162001838565b5f828152603460205260409020600101546200352e8162003b67565b6200174283836200495f565b7f62ba6ba2ffed8cfe316b583325ea41ac6e7ba9e5864d2bc6fabba7ac26d2f0f4620035668162003b67565b608780547fffffffffffffffffffffffffffffffffffffffffffffffff0000000000000000164267ffffffffffffffff16179055604080517fdbc1697600000000000000000000000000000000000000000000000000000000815290517f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff169163dbc16976916004808301925f92919082900301818387803b15801562003620575f80fd5b505af115801562003633573d5f803e3d5ffd5b505050506200364162004ef7565b50565b7f3dfe277d2a2c04b75fb2eb3743fa00005ae3678a20c299e65fdf4df76517f68e620036708162003b67565b67ffffffffffffffff84165f9081526083602052604090205463ffffffff1615620036c7576040517f6f91fc1200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b73ffffffffffffffffffffffffffffffffffffffff87165f9081526082602052604090205463ffffffff16156200372a576040517fd409b93000000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f6200373b88888888875f62004f86565b5f8080526002909101602052604090209390935550505050505050565b7fac75d24dbb35ea80e25fab167da4dea46c1915260426570db84f184891f5f590620037848162003b67565b607e80545f919082906200379e9063ffffffff1662005f64565b91906101000a81548163ffffffff021916908363ffffffff160217905590506040518060c001604052808973ffffffffffffffffffffffffffffffffffffffff1681526020018873ffffffffffffffffffffffffffffffffffffffff1681526020018767ffffffffffffffff1681526020018660ff1681526020015f1515815260200185815250607f5f8363ffffffff1663ffffffff1681526020019081526020015f205f820151815f015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506020820151816001015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555060408201518160010160146101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550606082015181600101601c6101000a81548160ff021916908360ff160217905550608082015181600101601d6101000a81548160ff02191690831515021790555060a082015181600201559050508063ffffffff167fa2970448b3bd66ba7e524e7b2a5b9cf94fa29e32488fb942afdfe70dd4b77b528989898989896040516200398b9695949392919062006222565b60405180910390a25050505050505050565b5f54610100900460ff1662001956576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602b60248201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960448201527f6e697469616c697a696e67000000000000000000000000000000000000000000606482015260840162000de8565b620019dc8282620048a3565b5f82815260346020526040808220600101805490849055905190918391839186917fbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff9190a4505050565b60855467ffffffffffffffff8281165f9081526004850160205260408120549092429262003abe92918116911662005f1f565b67ffffffffffffffff1611159392505050565b60068101545f90700100000000000000000000000000000000900467ffffffffffffffff161562003b445750600681015467ffffffffffffffff70010000000000000000000000000000000090910481165f90815260049092016020526040909120546801000000000000000090041690565b506006015468010000000000000000900467ffffffffffffffff1690565b919050565b62003641813362005209565b60078801545f9067ffffffffffffffff908116908716101562003bc2576040517fead1340b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff88161562003ca957600689015467ffffffffffffffff7001000000000000000000000000000000009091048116908916111562003c34576040517fbb14c20500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5067ffffffffffffffff8088165f90815260048a016020526040902060028101548154909288811668010000000000000000909204161462003ca2576040517f2bd2e3e700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5062003d56565b5067ffffffffffffffff85165f9081526002890160205260409020548062003cfd576040517f4997b98600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600689015467ffffffffffffffff680100000000000000009091048116908716111562003d56576040517f1e56e9e200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600689015467ffffffffffffffff7001000000000000000000000000000000009091048116908816118062003d9f57508767ffffffffffffffff168767ffffffffffffffff1611155b8062003dd95750600689015467ffffffffffffffff7801000000000000000000000000000000000000000000000000909104811690881611155b1562003e11576040517fbfa7079f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff8781165f90815260048b01602052604090205468010000000000000000900481169086161462003e76576040517f32a2a77f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f62003e878a888888868962004a19565b90505f7f30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f000000160028360405162003ebd919062006287565b602060405180830381855afa15801562003ed9573d5f803e3d5ffd5b5050506040513d601f19601f8201168201806040525081019062003efe91906200607e565b62003f0a919062006156565b60018c01546040805160208101825283815290517f9121da8a00000000000000000000000000000000000000000000000000000000815292935073ffffffffffffffffffffffffffffffffffffffff90911691639121da8a9162003f7491889190600401620062a4565b602060405180830381865afa15801562003f90573d5f803e3d5ffd5b505050506040513d601f19601f8201168201806040525081019062003fb69190620062e0565b62003fed576040517f09bde33900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff89165f90815260048c01602052604090206002015485900362004045576040517fa47276bd00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5050505050505050505050565b5f806200405f8a62003ad1565b60078b015490915067ffffffffffffffff9081169089161015620040af576040517fead1340b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff891615620041985760068a015467ffffffffffffffff7001000000000000000000000000000000009091048116908a16111562004121576040517fbb14c20500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff808a165f90815260048c01602052604090206002810154815490945090918a811668010000000000000000909204161462004191576040517f2bd2e3e700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b506200423c565b67ffffffffffffffff88165f90815260028b016020526040902054915081620041ed576040517f4997b98600000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8067ffffffffffffffff168867ffffffffffffffff1611156200423c576040517f1e56e9e200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8067ffffffffffffffff168767ffffffffffffffff16116200428a576040517fb9b18f5700000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f6200429b8b8a8a8a878b62004a19565b90505f7f30644e72e131a029b85045b68181585d2833e84879b9709143e1f593f0000001600283604051620042d1919062006287565b602060405180830381855afa158015620042ed573d5f803e3d5ffd5b5050506040513d601f19601f820116820180604052508101906200431291906200607e565b6200431e919062006156565b60018d01546040805160208101825283815290517f9121da8a00000000000000000000000000000000000000000000000000000000815292935073ffffffffffffffffffffffffffffffffffffffff90911691639121da8a916200438891899190600401620062a4565b602060405180830381865afa158015620043a4573d5f803e3d5ffd5b505050506040513d601f19601f82011682018060405250810190620043ca9190620062e0565b62004401576040517f09bde33900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f6200440e848b6200603b565b905062004475878267ffffffffffffffff166200442a620026f6565b62004436919062005f4a565b73ffffffffffffffffffffffffffffffffffffffff7f000000000000000000000000000000000000000000000000000000000000000016919062005272565b80608460088282829054906101000a900467ffffffffffffffff166200449c919062005f1f565b82546101009290920a67ffffffffffffffff818102199093169183160217909155608480547fffffffffffffffff0000000000000000ffffffffffffffffffffffffffffffff16700100000000000000000000000000000000428416021790558e546040517f32c2d153000000000000000000000000000000000000000000000000000000008152918d166004830152602482018b905233604483015273ffffffffffffffffffffffffffffffffffffffff1691506332c2d153906064015f604051808303815f87803b15801562004572575f80fd5b505af115801562004585573d5f803e3d5ffd5b5050505050505050505050505050505050565b600682015467ffffffffffffffff78010000000000000000000000000000000000000000000000009091048116908216111580620045fb5750600682015467ffffffffffffffff7001000000000000000000000000000000009091048116908216115b1562004633576040517fd086b70b00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b67ffffffffffffffff8181165f818152600485016020908152604080832080546006890180547fffffffffffffffffffffffffffffffff0000000000000000ffffffffffffffff166801000000000000000092839004909816918202979097178755600280830154828752908a01909452919093209190915560018201546005870155835477ffffffffffffffffffffffffffffffffffffffffffffffff167801000000000000000000000000000000000000000000000000909302929092179092557f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff166333d6247d6200473c62002c27565b6040518263ffffffff1660e01b81526004016200475b91815260200190565b5f604051808303815f87803b15801562004773575f80fd5b505af115801562004786573d5f803e3d5ffd5b5050855473ffffffffffffffffffffffffffffffffffffffff165f908152608260209081526040918290205460028701546001880154845167ffffffffffffffff898116825294810192909252818501529188166060830152915163ffffffff90921693507f581910eb7a27738945c2f00a91f2284b2d6de9d4e472b12f901c2b0df045e21b925081900360800190a250505050565b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff16632072f6c56040518163ffffffff1660e01b81526004015f604051808303815f87803b15801562004882575f80fd5b505af115801562004895573d5f803e3d5ffd5b505050506200195662005301565b5f82815260346020908152604080832073ffffffffffffffffffffffffffffffffffffffff8516845290915290205460ff16620019dc575f82815260346020908152604080832073ffffffffffffffffffffffffffffffffffffffff8516808552925280832080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016600117905551339285917f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d9190a45050565b5f82815260346020908152604080832073ffffffffffffffffffffffffffffffffffffffff8516845290915290205460ff1615620019dc575f82815260346020908152604080832073ffffffffffffffffffffffffffffffffffffffff8516808552925280832080547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0016905551339285917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a45050565b67ffffffffffffffff8086165f818152600389016020526040808220549388168252902054606092911580159062004a4f575081155b1562004a87576040517f6818c29e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8062004abf576040517f66385b5100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b62004aca8462005394565b62004b01576040517f176b913c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b885460018a01546040517fffffffffffffffffffffffffffffffffffffffff0000000000000000000000003360601b16602082015260348101889052605481018590527fffffffffffffffff00000000000000000000000000000000000000000000000060c08c811b821660748401527401000000000000000000000000000000000000000094859004811b8216607c84015293909204831b82166084820152608c810187905260ac810184905260cc81018990529189901b1660ec82015260f401604051602081830303815290604052925050509695505050505050565b5f62004bec8362003ad1565b9050815f8062004bfd84846200603b565b60855467ffffffffffffffff91821692505f9162004c299168010000000000000000900416426200616c565b90505b8467ffffffffffffffff168467ffffffffffffffff161462004cbc5767ffffffffffffffff8085165f9081526003890160205260409020600181015490911682101562004c9657600181015468010000000000000000900467ffffffffffffffff16945062004cb5565b62004ca286866200603b565b67ffffffffffffffff1693505062004cbc565b5062004c2c565b5f62004cc984846200616c565b90508381101562004d2757808403600c811162004ce7578062004cea565b600c5b9050806103e80a81608560109054906101000a900461ffff1661ffff160a608654028162004d1c5762004d1c62006096565b046086555062004d9e565b838103600c811162004d3a578062004d3d565b600c5b90505f816103e80a82608560109054906101000a900461ffff1661ffff160a670de0b6b3a7640000028162004d765762004d7662006096565b04905080608654670de0b6b3a7640000028162004d975762004d9762006096565b0460865550505b683635c9adc5dea00000608654111562004dc557683635c9adc5dea0000060865562004ddd565b633b9aca00608654101562004ddd57633b9aca006086555b5050505050505050565b600681015467ffffffffffffffff780100000000000000000000000000000000000000000000000082048116700100000000000000000000000000000000909204161115620036415760068101545f9062004e6a907801000000000000000000000000000000000000000000000000900467ffffffffffffffff16600162005f1f565b905062004e78828262003a8b565b15620019dc5760068201545f9060029062004eb3908490700100000000000000000000000000000000900467ffffffffffffffff166200603b565b62004ebf919062006301565b62004ecb908362005f1f565b905062004ed9838262003a8b565b1562004eeb5762001742838262004598565b62001742838362004598565b606f5460ff1662004f34576040517f5386698100000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b606f80547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001690556040517f1e5e34eea33501aecf2ebec9fe0e884a40804275ea7fe10b2ba084c8374308b3905f90a1565b608080545f918291829062004fa19063ffffffff1662005f64565b91906101000a81548163ffffffff021916908363ffffffff160217905590508060835f8767ffffffffffffffff1667ffffffffffffffff1681526020019081526020015f205f6101000a81548163ffffffff021916908363ffffffff1602179055508060825f8a73ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020015f205f6101000a81548163ffffffff021916908363ffffffff16021790555060815f8263ffffffff1663ffffffff1681526020019081526020015f20915087825f015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550858260010160146101000a81548167ffffffffffffffff021916908367ffffffffffffffff16021790555086826001015f6101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555084825f0160146101000a81548167ffffffffffffffff021916908367ffffffffffffffff160217905550838260070160106101000a81548160ff021916908360ff1602179055508063ffffffff167fadfc7d56f7e39b08b321534f14bfb135ad27698f7d2f5ad0edc2356ea9a3f850878a888888604051620051f695949392919067ffffffffffffffff958616815273ffffffffffffffffffffffffffffffffffffffff949094166020850152918416604084015260ff166060830152909116608082015260a00190565b60405180910390a2509695505050505050565b5f82815260346020908152604080832073ffffffffffffffffffffffffffffffffffffffff8516845290915290205460ff16620019dc576040517fec2b7c3e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6040805173ffffffffffffffffffffffffffffffffffffffff8416602482015260448082018490528251808303909101815260649091019091526020810180517bffffffffffffffffffffffffffffffffffffffffffffffffffffffff167fa9059cbb00000000000000000000000000000000000000000000000000000000179052620017429084906200541b565b606f5460ff16156200533f576040517f2f0047fc00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b606f80547fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff001660011790556040517f2261efe5aef6fedc1fd1550b25facc9181745623049c7901287030b9ad1a5497905f90a1565b5f67ffffffff0000000167ffffffffffffffff8316108015620053cc575067ffffffff00000001604083901c67ffffffffffffffff16105b8015620053ee575067ffffffff00000001608083901c67ffffffffffffffff16105b801562005406575067ffffffff0000000160c083901c105b156200541457506001919050565b505f919050565b5f6200547e826040518060400160405280602081526020017f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c65648152508573ffffffffffffffffffffffffffffffffffffffff166200552d9092919063ffffffff16565b8051909150156200174257808060200190518101906200549f9190620062e0565b62001742576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602a60248201527f5361666545524332303a204552433230206f7065726174696f6e20646964206e60448201527f6f74207375636365656400000000000000000000000000000000000000000000606482015260840162000de8565b60606200553d84845f8562005545565b949350505050565b606082471015620055d9576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602660248201527f416464726573733a20696e73756666696369656e742062616c616e636520666f60448201527f722063616c6c0000000000000000000000000000000000000000000000000000606482015260840162000de8565b5f808673ffffffffffffffffffffffffffffffffffffffff16858760405162005603919062006287565b5f6040518083038185875af1925050503d805f81146200563f576040519150601f19603f3d011682016040523d82523d5f602084013e62005644565b606091505b50915091506200221d8783838760608315620056eb5782515f03620056e35773ffffffffffffffffffffffffffffffffffffffff85163b620056e3576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e7472616374000000604482015260640162000de8565b50816200553d565b6200553d8383815115620057025781518083602001fd5b806040517f08c379a000000000000000000000000000000000000000000000000000000000815260040162000de8919062005ce7565b610a5e806200632b83390190565b73ffffffffffffffffffffffffffffffffffffffff8116811462003641575f80fd5b803567ffffffffffffffff8116811462003b62575f80fd5b5f805f805f805f805f806101408b8d0312156200579b575f80fd5b8a35620057a88162005746565b9950620057b860208c0162005768565b9850620057c860408c0162005768565b975060608b0135620057da8162005746565b965060808b0135620057ec8162005746565b955060a08b0135620057fe8162005746565b945060c08b0135620058108162005746565b935060e08b0135620058228162005746565b9250620058336101008c0162005768565b9150620058446101208c0162005768565b90509295989b9194979a5092959850565b803563ffffffff8116811462003b62575f80fd5b5f80604083850312156200587b575f80fd5b620058868362005855565b9150620058966020840162005768565b90509250929050565b5f60208284031215620058b0575f80fd5b620012448262005855565b80610300810183101562001247575f80fd5b5f805f805f805f806103e0898b031215620058e6575f80fd5b620058f18962005855565b97506200590160208a0162005768565b96506200591160408a0162005768565b95506200592160608a0162005768565b94506200593160808a0162005768565b935060a0890135925060c08901359150620059508a60e08b01620058bb565b90509295985092959890939650565b5f805f805f805f806103e0898b03121562005978575f80fd5b620059838962005855565b97506200599360208a0162005768565b9650620059a360408a0162005768565b9550620059b360608a0162005768565b94506080890135935060a0890135925060c0890135620059d38162005746565b9150620059508a60e08b01620058bb565b5f60208284031215620059f5575f80fd5b813561ffff8116811462005a07575f80fd5b9392505050565b5f6020828403121562005a1f575f80fd5b5035919050565b5f806040838503121562005a38575f80fd5b82359150602083013562005a4c8162005746565b809150509250929050565b5f6020828403121562005a68575f80fd5b620012448262005768565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b5f82601f83011262005ab0575f80fd5b813567ffffffffffffffff8082111562005ace5762005ace62005a73565b604051601f83017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f0116810190828211818310171562005b175762005b1762005a73565b8160405283815286602085880101111562005b30575f80fd5b836020870160208301375f602085830101528094505050505092915050565b5f805f805f805f60e0888a03121562005b66575f80fd5b62005b718862005855565b965062005b816020890162005768565b9550604088013562005b938162005746565b9450606088013562005ba58162005746565b9350608088013562005bb78162005746565b925060a088013567ffffffffffffffff8082111562005bd4575f80fd5b62005be28b838c0162005aa0565b935060c08a013591508082111562005bf8575f80fd5b5062005c078a828b0162005aa0565b91505092959891949750929550565b5f805f805f8060c0878903121562005c2c575f80fd5b62005c378762005855565b955062005c476020880162005768565b945062005c576040880162005768565b9350606087013592506080870135915060a087013590509295509295509295565b5f5b8381101562005c9457818101518382015260200162005c7a565b50505f910152565b5f815180845262005cb581602086016020860162005c78565b601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0169290920160200192915050565b602081525f62001244602083018462005c9c565b5f806040838503121562005d0d575f80fd5b62005d188362005768565b946020939093013593505050565b5f805f806060858703121562005d3a575f80fd5b843562005d478162005746565b935062005d576020860162005855565b9250604085013567ffffffffffffffff8082111562005d74575f80fd5b818701915087601f83011262005d88575f80fd5b81358181111562005d97575f80fd5b88602082850101111562005da9575f80fd5b95989497505060200194505050565b5f6020828403121562005dc9575f80fd5b813562005a078162005746565b803560ff8116811462003b62575f80fd5b5f805f805f8060c0878903121562005dfd575f80fd5b863562005e0a8162005746565b9550602087013562005e1c8162005746565b945062005e2c6040880162005768565b935062005e3c6060880162005768565b92506080870135915062005e5360a0880162005dd6565b90509295509295509295565b5f805f805f8060c0878903121562005e75575f80fd5b863562005e828162005746565b9550602087013562005e948162005746565b945062005ea46040880162005768565b935062005eb46060880162005dd6565b92506080870135915060a087013567ffffffffffffffff81111562005ed7575f80fd5b62005ee589828a0162005aa0565b9150509295509295509295565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b67ffffffffffffffff81811683821601908082111562005f435762005f4362005ef2565b5092915050565b808202811582820484141762001247576200124762005ef2565b5f63ffffffff80831681810362005f7f5762005f7f62005ef2565b6001019392505050565b5f73ffffffffffffffffffffffffffffffffffffffff80861683528085166020840152506060604083015262005fc3606083018462005c9c565b95945050505050565b5f73ffffffffffffffffffffffffffffffffffffffff8089168352808816602084015263ffffffff8716604084015280861660608401525060c060808301526200601a60c083018562005c9c565b82810360a08401526200602e818562005c9c565b9998505050505050505050565b67ffffffffffffffff82811682821603908082111562005f435762005f4362005ef2565b5f67ffffffffffffffff80831681810362005f7f5762005f7f62005ef2565b5f602082840312156200608f575f80fd5b5051919050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601260045260245ffd5b5f82620060d457620060d462006096565b500490565b8082018082111562001247576200124762005ef2565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b5f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036200614f576200614f62005ef2565b5060010190565b5f8262006167576200616762006096565b500690565b8181038181111562001247576200124762005ef2565b5f8162006193576200619362005ef2565b507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0190565b73ffffffffffffffffffffffffffffffffffffffff8416815260406020820152816040820152818360608301375f818301606090810191909152601f9092017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe016010192915050565b5f73ffffffffffffffffffffffffffffffffffffffff808916835280881660208401525067ffffffffffffffff8616604083015260ff8516606083015283608083015260c060a08301526200627b60c083018462005c9c565b98975050505050505050565b5f82516200629a81846020870162005c78565b9190910192915050565b6103208101610300808584378201835f5b6001811015620062d6578151835260209283019290910190600101620062b5565b5050509392505050565b5f60208284031215620062f1575f80fd5b8151801515811462005a07575f80fd5b5f67ffffffffffffffff808416806200631e576200631e62006096565b9216919091049291505056fe60a060405260405162000a5e38038062000a5e833981016040819052620000269162000375565b828162000034828262000060565b50506001600160a01b038216608052620000576200005160805190565b620000c5565b5050506200046c565b6200006b8262000136565b6040516001600160a01b038316907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b905f90a2805115620000b757620000b28282620001b5565b505050565b620000c16200022e565b5050565b7f7e644d79422f17c01e4894b5f4f588d331ebfa28653d42ae832dc59e38c9798f620001065f8051602062000a3e833981519152546001600160a01b031690565b604080516001600160a01b03928316815291841660208301520160405180910390a1620001338162000250565b50565b806001600160a01b03163b5f036200017157604051634c9c8ce360e01b81526001600160a01b03821660048201526024015b60405180910390fd5b807f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc5b80546001600160a01b0319166001600160a01b039290921691909117905550565b60605f80846001600160a01b031684604051620001d391906200044f565b5f60405180830381855af49150503d805f81146200020d576040519150601f19603f3d011682016040523d82523d5f602084013e62000212565b606091505b5090925090506200022585838362000291565b95945050505050565b34156200024e5760405163b398979f60e01b815260040160405180910390fd5b565b6001600160a01b0381166200027b57604051633173bdd160e11b81525f600482015260240162000168565b805f8051602062000a3e83398151915262000194565b606082620002aa57620002a482620002f7565b620002f0565b8151158015620002c257506001600160a01b0384163b155b15620002ed57604051639996b31560e01b81526001600160a01b038516600482015260240162000168565b50805b9392505050565b805115620003085780518082602001fd5b604051630a12f52160e11b815260040160405180910390fd5b80516001600160a01b038116811462000338575f80fd5b919050565b634e487b7160e01b5f52604160045260245ffd5b5f5b838110156200036d57818101518382015260200162000353565b50505f910152565b5f805f6060848603121562000388575f80fd5b620003938462000321565b9250620003a36020850162000321565b60408501519092506001600160401b0380821115620003c0575f80fd5b818601915086601f830112620003d4575f80fd5b815181811115620003e957620003e96200033d565b604051601f8201601f19908116603f011681019083821181831017156200041457620004146200033d565b816040528281528960208487010111156200042d575f80fd5b6200044083602083016020880162000351565b80955050505050509250925092565b5f82516200046281846020870162000351565b9190910192915050565b6080516105ba620004845f395f601001526105ba5ff3fe608060405261000c61000e565b005b7f000000000000000000000000000000000000000000000000000000000000000073ffffffffffffffffffffffffffffffffffffffff1633036100a7575f357fffffffff00000000000000000000000000000000000000000000000000000000167f4f1ef286000000000000000000000000000000000000000000000000000000001461009f5761009d6100ab565b565b61009d6100bb565b61009d5b61009d6100b66100e9565b61012d565b5f806100ca3660048184610410565b8101906100d79190610464565b915091506100e5828261014b565b5050565b5f6101287f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc5473ffffffffffffffffffffffffffffffffffffffff1690565b905090565b365f80375f80365f845af43d5f803e808015610147573d5ff35b3d5ffd5b610154826101b2565b60405173ffffffffffffffffffffffffffffffffffffffff8316907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b905f90a28051156101aa576101a58282610285565b505050565b6100e5610304565b8073ffffffffffffffffffffffffffffffffffffffff163b5f0361021f576040517f4c9c8ce300000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff821660048201526024015b60405180910390fd5b7f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc80547fffffffffffffffffffffffff00000000000000000000000000000000000000001673ffffffffffffffffffffffffffffffffffffffff92909216919091179055565b60605f808473ffffffffffffffffffffffffffffffffffffffff16846040516102ae9190610558565b5f60405180830381855af49150503d805f81146102e6576040519150601f19603f3d011682016040523d82523d5f602084013e6102eb565b606091505b50915091506102fb85838361033c565b95945050505050565b341561009d576040517fb398979f00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b6060826103515761034c826103ce565b6103c7565b8151158015610375575073ffffffffffffffffffffffffffffffffffffffff84163b155b156103c4576040517f9996b31500000000000000000000000000000000000000000000000000000000815273ffffffffffffffffffffffffffffffffffffffff85166004820152602401610216565b50805b9392505050565b8051156103de5780518082602001fd5b6040517f1425ea4200000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f808585111561041e575f80fd5b8386111561042a575f80fd5b5050820193919092039150565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b5f8060408385031215610475575f80fd5b823573ffffffffffffffffffffffffffffffffffffffff81168114610498575f80fd5b9150602083013567ffffffffffffffff808211156104b4575f80fd5b818501915085601f8301126104c7575f80fd5b8135818111156104d9576104d9610437565b604051601f82017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0908116603f0116810190838211818310171561051f5761051f610437565b81604052828152886020848701011115610537575f80fd5b826020860160208301375f6020848301015280955050505050509250929050565b5f82515f5b81811015610577576020818601810151858301520161055d565b505f92019182525091905056fea26469706673582212200ca61bd1e45d482203caba1d216b11bb6992f1ce0f6427bfe86e65b2f53457a264736f6c63430008140033b53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103a2646970667358221220a0a1c440ed85d9f9eeda617b6048e00cf12863c4f97b47eb0b40c67dec4452d164736f6c63430008140033",
}

// MockfeijoapolygonrollupmanagerABI is the input ABI used to generate the binding from.
// Deprecated: Use MockfeijoapolygonrollupmanagerMetaData.ABI instead.
var MockfeijoapolygonrollupmanagerABI = MockfeijoapolygonrollupmanagerMetaData.ABI

// MockfeijoapolygonrollupmanagerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use MockfeijoapolygonrollupmanagerMetaData.Bin instead.
var MockfeijoapolygonrollupmanagerBin = MockfeijoapolygonrollupmanagerMetaData.Bin

// DeployMockfeijoapolygonrollupmanager deploys a new Ethereum contract, binding an instance of Mockfeijoapolygonrollupmanager to it.
func DeployMockfeijoapolygonrollupmanager(auth *bind.TransactOpts, backend bind.ContractBackend, _globalExitRootManager common.Address, _pol common.Address, _bridgeAddress common.Address) (common.Address, *types.Transaction, *Mockfeijoapolygonrollupmanager, error) {
	parsed, err := MockfeijoapolygonrollupmanagerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(MockfeijoapolygonrollupmanagerBin), backend, _globalExitRootManager, _pol, _bridgeAddress)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Mockfeijoapolygonrollupmanager{MockfeijoapolygonrollupmanagerCaller: MockfeijoapolygonrollupmanagerCaller{contract: contract}, MockfeijoapolygonrollupmanagerTransactor: MockfeijoapolygonrollupmanagerTransactor{contract: contract}, MockfeijoapolygonrollupmanagerFilterer: MockfeijoapolygonrollupmanagerFilterer{contract: contract}}, nil
}

// Mockfeijoapolygonrollupmanager is an auto generated Go binding around an Ethereum contract.
type Mockfeijoapolygonrollupmanager struct {
	MockfeijoapolygonrollupmanagerCaller     // Read-only binding to the contract
	MockfeijoapolygonrollupmanagerTransactor // Write-only binding to the contract
	MockfeijoapolygonrollupmanagerFilterer   // Log filterer for contract events
}

// MockfeijoapolygonrollupmanagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type MockfeijoapolygonrollupmanagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MockfeijoapolygonrollupmanagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MockfeijoapolygonrollupmanagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MockfeijoapolygonrollupmanagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MockfeijoapolygonrollupmanagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MockfeijoapolygonrollupmanagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MockfeijoapolygonrollupmanagerSession struct {
	Contract     *Mockfeijoapolygonrollupmanager // Generic contract binding to set the session for
	CallOpts     bind.CallOpts                   // Call options to use throughout this session
	TransactOpts bind.TransactOpts               // Transaction auth options to use throughout this session
}

// MockfeijoapolygonrollupmanagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MockfeijoapolygonrollupmanagerCallerSession struct {
	Contract *MockfeijoapolygonrollupmanagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                         // Call options to use throughout this session
}

// MockfeijoapolygonrollupmanagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MockfeijoapolygonrollupmanagerTransactorSession struct {
	Contract     *MockfeijoapolygonrollupmanagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                         // Transaction auth options to use throughout this session
}

// MockfeijoapolygonrollupmanagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type MockfeijoapolygonrollupmanagerRaw struct {
	Contract *Mockfeijoapolygonrollupmanager // Generic contract binding to access the raw methods on
}

// MockfeijoapolygonrollupmanagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MockfeijoapolygonrollupmanagerCallerRaw struct {
	Contract *MockfeijoapolygonrollupmanagerCaller // Generic read-only contract binding to access the raw methods on
}

// MockfeijoapolygonrollupmanagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MockfeijoapolygonrollupmanagerTransactorRaw struct {
	Contract *MockfeijoapolygonrollupmanagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMockfeijoapolygonrollupmanager creates a new instance of Mockfeijoapolygonrollupmanager, bound to a specific deployed contract.
func NewMockfeijoapolygonrollupmanager(address common.Address, backend bind.ContractBackend) (*Mockfeijoapolygonrollupmanager, error) {
	contract, err := bindMockfeijoapolygonrollupmanager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Mockfeijoapolygonrollupmanager{MockfeijoapolygonrollupmanagerCaller: MockfeijoapolygonrollupmanagerCaller{contract: contract}, MockfeijoapolygonrollupmanagerTransactor: MockfeijoapolygonrollupmanagerTransactor{contract: contract}, MockfeijoapolygonrollupmanagerFilterer: MockfeijoapolygonrollupmanagerFilterer{contract: contract}}, nil
}

// NewMockfeijoapolygonrollupmanagerCaller creates a new read-only instance of Mockfeijoapolygonrollupmanager, bound to a specific deployed contract.
func NewMockfeijoapolygonrollupmanagerCaller(address common.Address, caller bind.ContractCaller) (*MockfeijoapolygonrollupmanagerCaller, error) {
	contract, err := bindMockfeijoapolygonrollupmanager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MockfeijoapolygonrollupmanagerCaller{contract: contract}, nil
}

// NewMockfeijoapolygonrollupmanagerTransactor creates a new write-only instance of Mockfeijoapolygonrollupmanager, bound to a specific deployed contract.
func NewMockfeijoapolygonrollupmanagerTransactor(address common.Address, transactor bind.ContractTransactor) (*MockfeijoapolygonrollupmanagerTransactor, error) {
	contract, err := bindMockfeijoapolygonrollupmanager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MockfeijoapolygonrollupmanagerTransactor{contract: contract}, nil
}

// NewMockfeijoapolygonrollupmanagerFilterer creates a new log filterer instance of Mockfeijoapolygonrollupmanager, bound to a specific deployed contract.
func NewMockfeijoapolygonrollupmanagerFilterer(address common.Address, filterer bind.ContractFilterer) (*MockfeijoapolygonrollupmanagerFilterer, error) {
	contract, err := bindMockfeijoapolygonrollupmanager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MockfeijoapolygonrollupmanagerFilterer{contract: contract}, nil
}

// bindMockfeijoapolygonrollupmanager binds a generic wrapper to an already deployed contract.
func bindMockfeijoapolygonrollupmanager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MockfeijoapolygonrollupmanagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Mockfeijoapolygonrollupmanager.Contract.MockfeijoapolygonrollupmanagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.MockfeijoapolygonrollupmanagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.MockfeijoapolygonrollupmanagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Mockfeijoapolygonrollupmanager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Mockfeijoapolygonrollupmanager.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.DEFAULTADMINROLE(&_Mockfeijoapolygonrollupmanager.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.DEFAULTADMINROLE(&_Mockfeijoapolygonrollupmanager.CallOpts)
}

// BridgeAddress is a free data retrieval call binding the contract method 0xa3c573eb.
//
// Solidity: function bridgeAddress() view returns(address)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCaller) BridgeAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Mockfeijoapolygonrollupmanager.contract.Call(opts, &out, "bridgeAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// BridgeAddress is a free data retrieval call binding the contract method 0xa3c573eb.
//
// Solidity: function bridgeAddress() view returns(address)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerSession) BridgeAddress() (common.Address, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.BridgeAddress(&_Mockfeijoapolygonrollupmanager.CallOpts)
}

// BridgeAddress is a free data retrieval call binding the contract method 0xa3c573eb.
//
// Solidity: function bridgeAddress() view returns(address)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCallerSession) BridgeAddress() (common.Address, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.BridgeAddress(&_Mockfeijoapolygonrollupmanager.CallOpts)
}

// CalculateRewardPerBatch is a free data retrieval call binding the contract method 0x99f5634e.
//
// Solidity: function calculateRewardPerBatch() view returns(uint256)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCaller) CalculateRewardPerBatch(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Mockfeijoapolygonrollupmanager.contract.Call(opts, &out, "calculateRewardPerBatch")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CalculateRewardPerBatch is a free data retrieval call binding the contract method 0x99f5634e.
//
// Solidity: function calculateRewardPerBatch() view returns(uint256)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerSession) CalculateRewardPerBatch() (*big.Int, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.CalculateRewardPerBatch(&_Mockfeijoapolygonrollupmanager.CallOpts)
}

// CalculateRewardPerBatch is a free data retrieval call binding the contract method 0x99f5634e.
//
// Solidity: function calculateRewardPerBatch() view returns(uint256)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCallerSession) CalculateRewardPerBatch() (*big.Int, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.CalculateRewardPerBatch(&_Mockfeijoapolygonrollupmanager.CallOpts)
}

// ChainIDToRollupID is a free data retrieval call binding the contract method 0x7fb6e76a.
//
// Solidity: function chainIDToRollupID(uint64 chainID) view returns(uint32 rollupID)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCaller) ChainIDToRollupID(opts *bind.CallOpts, chainID uint64) (uint32, error) {
	var out []interface{}
	err := _Mockfeijoapolygonrollupmanager.contract.Call(opts, &out, "chainIDToRollupID", chainID)

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// ChainIDToRollupID is a free data retrieval call binding the contract method 0x7fb6e76a.
//
// Solidity: function chainIDToRollupID(uint64 chainID) view returns(uint32 rollupID)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerSession) ChainIDToRollupID(chainID uint64) (uint32, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.ChainIDToRollupID(&_Mockfeijoapolygonrollupmanager.CallOpts, chainID)
}

// ChainIDToRollupID is a free data retrieval call binding the contract method 0x7fb6e76a.
//
// Solidity: function chainIDToRollupID(uint64 chainID) view returns(uint32 rollupID)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCallerSession) ChainIDToRollupID(chainID uint64) (uint32, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.ChainIDToRollupID(&_Mockfeijoapolygonrollupmanager.CallOpts, chainID)
}

// GetBatchFee is a free data retrieval call binding the contract method 0x477fa270.
//
// Solidity: function getBatchFee() view returns(uint256)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCaller) GetBatchFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Mockfeijoapolygonrollupmanager.contract.Call(opts, &out, "getBatchFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBatchFee is a free data retrieval call binding the contract method 0x477fa270.
//
// Solidity: function getBatchFee() view returns(uint256)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerSession) GetBatchFee() (*big.Int, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.GetBatchFee(&_Mockfeijoapolygonrollupmanager.CallOpts)
}

// GetBatchFee is a free data retrieval call binding the contract method 0x477fa270.
//
// Solidity: function getBatchFee() view returns(uint256)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCallerSession) GetBatchFee() (*big.Int, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.GetBatchFee(&_Mockfeijoapolygonrollupmanager.CallOpts)
}

// GetForcedBatchFee is a free data retrieval call binding the contract method 0x60469169.
//
// Solidity: function getForcedBatchFee() view returns(uint256)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCaller) GetForcedBatchFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Mockfeijoapolygonrollupmanager.contract.Call(opts, &out, "getForcedBatchFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetForcedBatchFee is a free data retrieval call binding the contract method 0x60469169.
//
// Solidity: function getForcedBatchFee() view returns(uint256)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerSession) GetForcedBatchFee() (*big.Int, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.GetForcedBatchFee(&_Mockfeijoapolygonrollupmanager.CallOpts)
}

// GetForcedBatchFee is a free data retrieval call binding the contract method 0x60469169.
//
// Solidity: function getForcedBatchFee() view returns(uint256)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCallerSession) GetForcedBatchFee() (*big.Int, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.GetForcedBatchFee(&_Mockfeijoapolygonrollupmanager.CallOpts)
}

// GetInputSnarkBytes is a free data retrieval call binding the contract method 0x7975fcfe.
//
// Solidity: function getInputSnarkBytes(uint32 rollupID, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 oldStateRoot, bytes32 newStateRoot) view returns(bytes)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCaller) GetInputSnarkBytes(opts *bind.CallOpts, rollupID uint32, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, oldStateRoot [32]byte, newStateRoot [32]byte) ([]byte, error) {
	var out []interface{}
	err := _Mockfeijoapolygonrollupmanager.contract.Call(opts, &out, "getInputSnarkBytes", rollupID, initNumBatch, finalNewBatch, newLocalExitRoot, oldStateRoot, newStateRoot)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetInputSnarkBytes is a free data retrieval call binding the contract method 0x7975fcfe.
//
// Solidity: function getInputSnarkBytes(uint32 rollupID, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 oldStateRoot, bytes32 newStateRoot) view returns(bytes)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerSession) GetInputSnarkBytes(rollupID uint32, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, oldStateRoot [32]byte, newStateRoot [32]byte) ([]byte, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.GetInputSnarkBytes(&_Mockfeijoapolygonrollupmanager.CallOpts, rollupID, initNumBatch, finalNewBatch, newLocalExitRoot, oldStateRoot, newStateRoot)
}

// GetInputSnarkBytes is a free data retrieval call binding the contract method 0x7975fcfe.
//
// Solidity: function getInputSnarkBytes(uint32 rollupID, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 oldStateRoot, bytes32 newStateRoot) view returns(bytes)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCallerSession) GetInputSnarkBytes(rollupID uint32, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, oldStateRoot [32]byte, newStateRoot [32]byte) ([]byte, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.GetInputSnarkBytes(&_Mockfeijoapolygonrollupmanager.CallOpts, rollupID, initNumBatch, finalNewBatch, newLocalExitRoot, oldStateRoot, newStateRoot)
}

// GetLastVerifiedBatch is a free data retrieval call binding the contract method 0x11f6b287.
//
// Solidity: function getLastVerifiedBatch(uint32 rollupID) view returns(uint64)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCaller) GetLastVerifiedBatch(opts *bind.CallOpts, rollupID uint32) (uint64, error) {
	var out []interface{}
	err := _Mockfeijoapolygonrollupmanager.contract.Call(opts, &out, "getLastVerifiedBatch", rollupID)

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// GetLastVerifiedBatch is a free data retrieval call binding the contract method 0x11f6b287.
//
// Solidity: function getLastVerifiedBatch(uint32 rollupID) view returns(uint64)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerSession) GetLastVerifiedBatch(rollupID uint32) (uint64, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.GetLastVerifiedBatch(&_Mockfeijoapolygonrollupmanager.CallOpts, rollupID)
}

// GetLastVerifiedBatch is a free data retrieval call binding the contract method 0x11f6b287.
//
// Solidity: function getLastVerifiedBatch(uint32 rollupID) view returns(uint64)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCallerSession) GetLastVerifiedBatch(rollupID uint32) (uint64, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.GetLastVerifiedBatch(&_Mockfeijoapolygonrollupmanager.CallOpts, rollupID)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _Mockfeijoapolygonrollupmanager.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.GetRoleAdmin(&_Mockfeijoapolygonrollupmanager.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.GetRoleAdmin(&_Mockfeijoapolygonrollupmanager.CallOpts, role)
}

// GetRollupBatchNumToStateRoot is a free data retrieval call binding the contract method 0x55a71ee0.
//
// Solidity: function getRollupBatchNumToStateRoot(uint32 rollupID, uint64 batchNum) view returns(bytes32)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCaller) GetRollupBatchNumToStateRoot(opts *bind.CallOpts, rollupID uint32, batchNum uint64) ([32]byte, error) {
	var out []interface{}
	err := _Mockfeijoapolygonrollupmanager.contract.Call(opts, &out, "getRollupBatchNumToStateRoot", rollupID, batchNum)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRollupBatchNumToStateRoot is a free data retrieval call binding the contract method 0x55a71ee0.
//
// Solidity: function getRollupBatchNumToStateRoot(uint32 rollupID, uint64 batchNum) view returns(bytes32)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerSession) GetRollupBatchNumToStateRoot(rollupID uint32, batchNum uint64) ([32]byte, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.GetRollupBatchNumToStateRoot(&_Mockfeijoapolygonrollupmanager.CallOpts, rollupID, batchNum)
}

// GetRollupBatchNumToStateRoot is a free data retrieval call binding the contract method 0x55a71ee0.
//
// Solidity: function getRollupBatchNumToStateRoot(uint32 rollupID, uint64 batchNum) view returns(bytes32)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCallerSession) GetRollupBatchNumToStateRoot(rollupID uint32, batchNum uint64) ([32]byte, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.GetRollupBatchNumToStateRoot(&_Mockfeijoapolygonrollupmanager.CallOpts, rollupID, batchNum)
}

// GetRollupExitRoot is a free data retrieval call binding the contract method 0xa2967d99.
//
// Solidity: function getRollupExitRoot() view returns(bytes32)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCaller) GetRollupExitRoot(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _Mockfeijoapolygonrollupmanager.contract.Call(opts, &out, "getRollupExitRoot")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRollupExitRoot is a free data retrieval call binding the contract method 0xa2967d99.
//
// Solidity: function getRollupExitRoot() view returns(bytes32)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerSession) GetRollupExitRoot() ([32]byte, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.GetRollupExitRoot(&_Mockfeijoapolygonrollupmanager.CallOpts)
}

// GetRollupExitRoot is a free data retrieval call binding the contract method 0xa2967d99.
//
// Solidity: function getRollupExitRoot() view returns(bytes32)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCallerSession) GetRollupExitRoot() ([32]byte, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.GetRollupExitRoot(&_Mockfeijoapolygonrollupmanager.CallOpts)
}

// GetRollupPendingStateTransitions is a free data retrieval call binding the contract method 0xb99d0ad7.
//
// Solidity: function getRollupPendingStateTransitions(uint32 rollupID, uint64 batchNum) view returns((uint64,uint64,bytes32,bytes32))
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCaller) GetRollupPendingStateTransitions(opts *bind.CallOpts, rollupID uint32, batchNum uint64) (LegacyZKEVMStateVariablesPendingState, error) {
	var out []interface{}
	err := _Mockfeijoapolygonrollupmanager.contract.Call(opts, &out, "getRollupPendingStateTransitions", rollupID, batchNum)

	if err != nil {
		return *new(LegacyZKEVMStateVariablesPendingState), err
	}

	out0 := *abi.ConvertType(out[0], new(LegacyZKEVMStateVariablesPendingState)).(*LegacyZKEVMStateVariablesPendingState)

	return out0, err

}

// GetRollupPendingStateTransitions is a free data retrieval call binding the contract method 0xb99d0ad7.
//
// Solidity: function getRollupPendingStateTransitions(uint32 rollupID, uint64 batchNum) view returns((uint64,uint64,bytes32,bytes32))
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerSession) GetRollupPendingStateTransitions(rollupID uint32, batchNum uint64) (LegacyZKEVMStateVariablesPendingState, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.GetRollupPendingStateTransitions(&_Mockfeijoapolygonrollupmanager.CallOpts, rollupID, batchNum)
}

// GetRollupPendingStateTransitions is a free data retrieval call binding the contract method 0xb99d0ad7.
//
// Solidity: function getRollupPendingStateTransitions(uint32 rollupID, uint64 batchNum) view returns((uint64,uint64,bytes32,bytes32))
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCallerSession) GetRollupPendingStateTransitions(rollupID uint32, batchNum uint64) (LegacyZKEVMStateVariablesPendingState, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.GetRollupPendingStateTransitions(&_Mockfeijoapolygonrollupmanager.CallOpts, rollupID, batchNum)
}

// GetRollupSequencedBatches is a free data retrieval call binding the contract method 0x25280169.
//
// Solidity: function getRollupSequencedBatches(uint32 rollupID, uint64 batchNum) view returns((bytes32,uint64,uint64))
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCaller) GetRollupSequencedBatches(opts *bind.CallOpts, rollupID uint32, batchNum uint64) (LegacyZKEVMStateVariablesSequencedBatchData, error) {
	var out []interface{}
	err := _Mockfeijoapolygonrollupmanager.contract.Call(opts, &out, "getRollupSequencedBatches", rollupID, batchNum)

	if err != nil {
		return *new(LegacyZKEVMStateVariablesSequencedBatchData), err
	}

	out0 := *abi.ConvertType(out[0], new(LegacyZKEVMStateVariablesSequencedBatchData)).(*LegacyZKEVMStateVariablesSequencedBatchData)

	return out0, err

}

// GetRollupSequencedBatches is a free data retrieval call binding the contract method 0x25280169.
//
// Solidity: function getRollupSequencedBatches(uint32 rollupID, uint64 batchNum) view returns((bytes32,uint64,uint64))
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerSession) GetRollupSequencedBatches(rollupID uint32, batchNum uint64) (LegacyZKEVMStateVariablesSequencedBatchData, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.GetRollupSequencedBatches(&_Mockfeijoapolygonrollupmanager.CallOpts, rollupID, batchNum)
}

// GetRollupSequencedBatches is a free data retrieval call binding the contract method 0x25280169.
//
// Solidity: function getRollupSequencedBatches(uint32 rollupID, uint64 batchNum) view returns((bytes32,uint64,uint64))
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCallerSession) GetRollupSequencedBatches(rollupID uint32, batchNum uint64) (LegacyZKEVMStateVariablesSequencedBatchData, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.GetRollupSequencedBatches(&_Mockfeijoapolygonrollupmanager.CallOpts, rollupID, batchNum)
}

// GlobalExitRootManager is a free data retrieval call binding the contract method 0xd02103ca.
//
// Solidity: function globalExitRootManager() view returns(address)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCaller) GlobalExitRootManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Mockfeijoapolygonrollupmanager.contract.Call(opts, &out, "globalExitRootManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GlobalExitRootManager is a free data retrieval call binding the contract method 0xd02103ca.
//
// Solidity: function globalExitRootManager() view returns(address)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerSession) GlobalExitRootManager() (common.Address, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.GlobalExitRootManager(&_Mockfeijoapolygonrollupmanager.CallOpts)
}

// GlobalExitRootManager is a free data retrieval call binding the contract method 0xd02103ca.
//
// Solidity: function globalExitRootManager() view returns(address)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCallerSession) GlobalExitRootManager() (common.Address, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.GlobalExitRootManager(&_Mockfeijoapolygonrollupmanager.CallOpts)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _Mockfeijoapolygonrollupmanager.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.HasRole(&_Mockfeijoapolygonrollupmanager.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.HasRole(&_Mockfeijoapolygonrollupmanager.CallOpts, role, account)
}

// IsEmergencyState is a free data retrieval call binding the contract method 0x15064c96.
//
// Solidity: function isEmergencyState() view returns(bool)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCaller) IsEmergencyState(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Mockfeijoapolygonrollupmanager.contract.Call(opts, &out, "isEmergencyState")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsEmergencyState is a free data retrieval call binding the contract method 0x15064c96.
//
// Solidity: function isEmergencyState() view returns(bool)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerSession) IsEmergencyState() (bool, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.IsEmergencyState(&_Mockfeijoapolygonrollupmanager.CallOpts)
}

// IsEmergencyState is a free data retrieval call binding the contract method 0x15064c96.
//
// Solidity: function isEmergencyState() view returns(bool)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCallerSession) IsEmergencyState() (bool, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.IsEmergencyState(&_Mockfeijoapolygonrollupmanager.CallOpts)
}

// IsPendingStateConsolidable is a free data retrieval call binding the contract method 0x080b3111.
//
// Solidity: function isPendingStateConsolidable(uint32 rollupID, uint64 pendingStateNum) view returns(bool)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCaller) IsPendingStateConsolidable(opts *bind.CallOpts, rollupID uint32, pendingStateNum uint64) (bool, error) {
	var out []interface{}
	err := _Mockfeijoapolygonrollupmanager.contract.Call(opts, &out, "isPendingStateConsolidable", rollupID, pendingStateNum)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsPendingStateConsolidable is a free data retrieval call binding the contract method 0x080b3111.
//
// Solidity: function isPendingStateConsolidable(uint32 rollupID, uint64 pendingStateNum) view returns(bool)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerSession) IsPendingStateConsolidable(rollupID uint32, pendingStateNum uint64) (bool, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.IsPendingStateConsolidable(&_Mockfeijoapolygonrollupmanager.CallOpts, rollupID, pendingStateNum)
}

// IsPendingStateConsolidable is a free data retrieval call binding the contract method 0x080b3111.
//
// Solidity: function isPendingStateConsolidable(uint32 rollupID, uint64 pendingStateNum) view returns(bool)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCallerSession) IsPendingStateConsolidable(rollupID uint32, pendingStateNum uint64) (bool, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.IsPendingStateConsolidable(&_Mockfeijoapolygonrollupmanager.CallOpts, rollupID, pendingStateNum)
}

// LastAggregationTimestamp is a free data retrieval call binding the contract method 0xc1acbc34.
//
// Solidity: function lastAggregationTimestamp() view returns(uint64)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCaller) LastAggregationTimestamp(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Mockfeijoapolygonrollupmanager.contract.Call(opts, &out, "lastAggregationTimestamp")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// LastAggregationTimestamp is a free data retrieval call binding the contract method 0xc1acbc34.
//
// Solidity: function lastAggregationTimestamp() view returns(uint64)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerSession) LastAggregationTimestamp() (uint64, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.LastAggregationTimestamp(&_Mockfeijoapolygonrollupmanager.CallOpts)
}

// LastAggregationTimestamp is a free data retrieval call binding the contract method 0xc1acbc34.
//
// Solidity: function lastAggregationTimestamp() view returns(uint64)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCallerSession) LastAggregationTimestamp() (uint64, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.LastAggregationTimestamp(&_Mockfeijoapolygonrollupmanager.CallOpts)
}

// LastDeactivatedEmergencyStateTimestamp is a free data retrieval call binding the contract method 0x30c27dde.
//
// Solidity: function lastDeactivatedEmergencyStateTimestamp() view returns(uint64)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCaller) LastDeactivatedEmergencyStateTimestamp(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Mockfeijoapolygonrollupmanager.contract.Call(opts, &out, "lastDeactivatedEmergencyStateTimestamp")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// LastDeactivatedEmergencyStateTimestamp is a free data retrieval call binding the contract method 0x30c27dde.
//
// Solidity: function lastDeactivatedEmergencyStateTimestamp() view returns(uint64)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerSession) LastDeactivatedEmergencyStateTimestamp() (uint64, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.LastDeactivatedEmergencyStateTimestamp(&_Mockfeijoapolygonrollupmanager.CallOpts)
}

// LastDeactivatedEmergencyStateTimestamp is a free data retrieval call binding the contract method 0x30c27dde.
//
// Solidity: function lastDeactivatedEmergencyStateTimestamp() view returns(uint64)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCallerSession) LastDeactivatedEmergencyStateTimestamp() (uint64, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.LastDeactivatedEmergencyStateTimestamp(&_Mockfeijoapolygonrollupmanager.CallOpts)
}

// MultiplierBatchFee is a free data retrieval call binding the contract method 0xafd23cbe.
//
// Solidity: function multiplierBatchFee() view returns(uint16)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCaller) MultiplierBatchFee(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _Mockfeijoapolygonrollupmanager.contract.Call(opts, &out, "multiplierBatchFee")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// MultiplierBatchFee is a free data retrieval call binding the contract method 0xafd23cbe.
//
// Solidity: function multiplierBatchFee() view returns(uint16)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerSession) MultiplierBatchFee() (uint16, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.MultiplierBatchFee(&_Mockfeijoapolygonrollupmanager.CallOpts)
}

// MultiplierBatchFee is a free data retrieval call binding the contract method 0xafd23cbe.
//
// Solidity: function multiplierBatchFee() view returns(uint16)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCallerSession) MultiplierBatchFee() (uint16, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.MultiplierBatchFee(&_Mockfeijoapolygonrollupmanager.CallOpts)
}

// PendingStateTimeout is a free data retrieval call binding the contract method 0xd939b315.
//
// Solidity: function pendingStateTimeout() view returns(uint64)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCaller) PendingStateTimeout(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Mockfeijoapolygonrollupmanager.contract.Call(opts, &out, "pendingStateTimeout")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// PendingStateTimeout is a free data retrieval call binding the contract method 0xd939b315.
//
// Solidity: function pendingStateTimeout() view returns(uint64)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerSession) PendingStateTimeout() (uint64, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.PendingStateTimeout(&_Mockfeijoapolygonrollupmanager.CallOpts)
}

// PendingStateTimeout is a free data retrieval call binding the contract method 0xd939b315.
//
// Solidity: function pendingStateTimeout() view returns(uint64)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCallerSession) PendingStateTimeout() (uint64, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.PendingStateTimeout(&_Mockfeijoapolygonrollupmanager.CallOpts)
}

// Pol is a free data retrieval call binding the contract method 0xe46761c4.
//
// Solidity: function pol() view returns(address)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCaller) Pol(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Mockfeijoapolygonrollupmanager.contract.Call(opts, &out, "pol")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Pol is a free data retrieval call binding the contract method 0xe46761c4.
//
// Solidity: function pol() view returns(address)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerSession) Pol() (common.Address, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.Pol(&_Mockfeijoapolygonrollupmanager.CallOpts)
}

// Pol is a free data retrieval call binding the contract method 0xe46761c4.
//
// Solidity: function pol() view returns(address)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCallerSession) Pol() (common.Address, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.Pol(&_Mockfeijoapolygonrollupmanager.CallOpts)
}

// RollupAddressToID is a free data retrieval call binding the contract method 0xceee281d.
//
// Solidity: function rollupAddressToID(address rollupAddress) view returns(uint32 rollupID)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCaller) RollupAddressToID(opts *bind.CallOpts, rollupAddress common.Address) (uint32, error) {
	var out []interface{}
	err := _Mockfeijoapolygonrollupmanager.contract.Call(opts, &out, "rollupAddressToID", rollupAddress)

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// RollupAddressToID is a free data retrieval call binding the contract method 0xceee281d.
//
// Solidity: function rollupAddressToID(address rollupAddress) view returns(uint32 rollupID)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerSession) RollupAddressToID(rollupAddress common.Address) (uint32, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.RollupAddressToID(&_Mockfeijoapolygonrollupmanager.CallOpts, rollupAddress)
}

// RollupAddressToID is a free data retrieval call binding the contract method 0xceee281d.
//
// Solidity: function rollupAddressToID(address rollupAddress) view returns(uint32 rollupID)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCallerSession) RollupAddressToID(rollupAddress common.Address) (uint32, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.RollupAddressToID(&_Mockfeijoapolygonrollupmanager.CallOpts, rollupAddress)
}

// RollupCount is a free data retrieval call binding the contract method 0xf4e92675.
//
// Solidity: function rollupCount() view returns(uint32)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCaller) RollupCount(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _Mockfeijoapolygonrollupmanager.contract.Call(opts, &out, "rollupCount")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// RollupCount is a free data retrieval call binding the contract method 0xf4e92675.
//
// Solidity: function rollupCount() view returns(uint32)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerSession) RollupCount() (uint32, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.RollupCount(&_Mockfeijoapolygonrollupmanager.CallOpts)
}

// RollupCount is a free data retrieval call binding the contract method 0xf4e92675.
//
// Solidity: function rollupCount() view returns(uint32)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCallerSession) RollupCount() (uint32, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.RollupCount(&_Mockfeijoapolygonrollupmanager.CallOpts)
}

// RollupIDToRollupData is a free data retrieval call binding the contract method 0xf9c4c2ae.
//
// Solidity: function rollupIDToRollupData(uint32 rollupID) view returns(address rollupContract, uint64 chainID, address verifier, uint64 forkID, bytes32 lastLocalExitRoot, uint64 lastBatchSequenced, uint64 lastVerifiedBatch, uint64 lastPendingState, uint64 lastPendingStateConsolidated, uint64 lastVerifiedBatchBeforeUpgrade, uint64 rollupTypeID, uint8 rollupCompatibilityID)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCaller) RollupIDToRollupData(opts *bind.CallOpts, rollupID uint32) (struct {
	RollupContract                 common.Address
	ChainID                        uint64
	Verifier                       common.Address
	ForkID                         uint64
	LastLocalExitRoot              [32]byte
	LastBatchSequenced             uint64
	LastVerifiedBatch              uint64
	LastPendingState               uint64
	LastPendingStateConsolidated   uint64
	LastVerifiedBatchBeforeUpgrade uint64
	RollupTypeID                   uint64
	RollupCompatibilityID          uint8
}, error) {
	var out []interface{}
	err := _Mockfeijoapolygonrollupmanager.contract.Call(opts, &out, "rollupIDToRollupData", rollupID)

	outstruct := new(struct {
		RollupContract                 common.Address
		ChainID                        uint64
		Verifier                       common.Address
		ForkID                         uint64
		LastLocalExitRoot              [32]byte
		LastBatchSequenced             uint64
		LastVerifiedBatch              uint64
		LastPendingState               uint64
		LastPendingStateConsolidated   uint64
		LastVerifiedBatchBeforeUpgrade uint64
		RollupTypeID                   uint64
		RollupCompatibilityID          uint8
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.RollupContract = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.ChainID = *abi.ConvertType(out[1], new(uint64)).(*uint64)
	outstruct.Verifier = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)
	outstruct.ForkID = *abi.ConvertType(out[3], new(uint64)).(*uint64)
	outstruct.LastLocalExitRoot = *abi.ConvertType(out[4], new([32]byte)).(*[32]byte)
	outstruct.LastBatchSequenced = *abi.ConvertType(out[5], new(uint64)).(*uint64)
	outstruct.LastVerifiedBatch = *abi.ConvertType(out[6], new(uint64)).(*uint64)
	outstruct.LastPendingState = *abi.ConvertType(out[7], new(uint64)).(*uint64)
	outstruct.LastPendingStateConsolidated = *abi.ConvertType(out[8], new(uint64)).(*uint64)
	outstruct.LastVerifiedBatchBeforeUpgrade = *abi.ConvertType(out[9], new(uint64)).(*uint64)
	outstruct.RollupTypeID = *abi.ConvertType(out[10], new(uint64)).(*uint64)
	outstruct.RollupCompatibilityID = *abi.ConvertType(out[11], new(uint8)).(*uint8)

	return *outstruct, err

}

// RollupIDToRollupData is a free data retrieval call binding the contract method 0xf9c4c2ae.
//
// Solidity: function rollupIDToRollupData(uint32 rollupID) view returns(address rollupContract, uint64 chainID, address verifier, uint64 forkID, bytes32 lastLocalExitRoot, uint64 lastBatchSequenced, uint64 lastVerifiedBatch, uint64 lastPendingState, uint64 lastPendingStateConsolidated, uint64 lastVerifiedBatchBeforeUpgrade, uint64 rollupTypeID, uint8 rollupCompatibilityID)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerSession) RollupIDToRollupData(rollupID uint32) (struct {
	RollupContract                 common.Address
	ChainID                        uint64
	Verifier                       common.Address
	ForkID                         uint64
	LastLocalExitRoot              [32]byte
	LastBatchSequenced             uint64
	LastVerifiedBatch              uint64
	LastPendingState               uint64
	LastPendingStateConsolidated   uint64
	LastVerifiedBatchBeforeUpgrade uint64
	RollupTypeID                   uint64
	RollupCompatibilityID          uint8
}, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.RollupIDToRollupData(&_Mockfeijoapolygonrollupmanager.CallOpts, rollupID)
}

// RollupIDToRollupData is a free data retrieval call binding the contract method 0xf9c4c2ae.
//
// Solidity: function rollupIDToRollupData(uint32 rollupID) view returns(address rollupContract, uint64 chainID, address verifier, uint64 forkID, bytes32 lastLocalExitRoot, uint64 lastBatchSequenced, uint64 lastVerifiedBatch, uint64 lastPendingState, uint64 lastPendingStateConsolidated, uint64 lastVerifiedBatchBeforeUpgrade, uint64 rollupTypeID, uint8 rollupCompatibilityID)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCallerSession) RollupIDToRollupData(rollupID uint32) (struct {
	RollupContract                 common.Address
	ChainID                        uint64
	Verifier                       common.Address
	ForkID                         uint64
	LastLocalExitRoot              [32]byte
	LastBatchSequenced             uint64
	LastVerifiedBatch              uint64
	LastPendingState               uint64
	LastPendingStateConsolidated   uint64
	LastVerifiedBatchBeforeUpgrade uint64
	RollupTypeID                   uint64
	RollupCompatibilityID          uint8
}, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.RollupIDToRollupData(&_Mockfeijoapolygonrollupmanager.CallOpts, rollupID)
}

// RollupTypeCount is a free data retrieval call binding the contract method 0x1796a1ae.
//
// Solidity: function rollupTypeCount() view returns(uint32)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCaller) RollupTypeCount(opts *bind.CallOpts) (uint32, error) {
	var out []interface{}
	err := _Mockfeijoapolygonrollupmanager.contract.Call(opts, &out, "rollupTypeCount")

	if err != nil {
		return *new(uint32), err
	}

	out0 := *abi.ConvertType(out[0], new(uint32)).(*uint32)

	return out0, err

}

// RollupTypeCount is a free data retrieval call binding the contract method 0x1796a1ae.
//
// Solidity: function rollupTypeCount() view returns(uint32)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerSession) RollupTypeCount() (uint32, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.RollupTypeCount(&_Mockfeijoapolygonrollupmanager.CallOpts)
}

// RollupTypeCount is a free data retrieval call binding the contract method 0x1796a1ae.
//
// Solidity: function rollupTypeCount() view returns(uint32)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCallerSession) RollupTypeCount() (uint32, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.RollupTypeCount(&_Mockfeijoapolygonrollupmanager.CallOpts)
}

// RollupTypeMap is a free data retrieval call binding the contract method 0x65c0504d.
//
// Solidity: function rollupTypeMap(uint32 rollupTypeID) view returns(address consensusImplementation, address verifier, uint64 forkID, uint8 rollupCompatibilityID, bool obsolete, bytes32 genesis)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCaller) RollupTypeMap(opts *bind.CallOpts, rollupTypeID uint32) (struct {
	ConsensusImplementation common.Address
	Verifier                common.Address
	ForkID                  uint64
	RollupCompatibilityID   uint8
	Obsolete                bool
	Genesis                 [32]byte
}, error) {
	var out []interface{}
	err := _Mockfeijoapolygonrollupmanager.contract.Call(opts, &out, "rollupTypeMap", rollupTypeID)

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
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerSession) RollupTypeMap(rollupTypeID uint32) (struct {
	ConsensusImplementation common.Address
	Verifier                common.Address
	ForkID                  uint64
	RollupCompatibilityID   uint8
	Obsolete                bool
	Genesis                 [32]byte
}, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.RollupTypeMap(&_Mockfeijoapolygonrollupmanager.CallOpts, rollupTypeID)
}

// RollupTypeMap is a free data retrieval call binding the contract method 0x65c0504d.
//
// Solidity: function rollupTypeMap(uint32 rollupTypeID) view returns(address consensusImplementation, address verifier, uint64 forkID, uint8 rollupCompatibilityID, bool obsolete, bytes32 genesis)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCallerSession) RollupTypeMap(rollupTypeID uint32) (struct {
	ConsensusImplementation common.Address
	Verifier                common.Address
	ForkID                  uint64
	RollupCompatibilityID   uint8
	Obsolete                bool
	Genesis                 [32]byte
}, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.RollupTypeMap(&_Mockfeijoapolygonrollupmanager.CallOpts, rollupTypeID)
}

// TotalSequencedBatches is a free data retrieval call binding the contract method 0x066ec012.
//
// Solidity: function totalSequencedBatches() view returns(uint64)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCaller) TotalSequencedBatches(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Mockfeijoapolygonrollupmanager.contract.Call(opts, &out, "totalSequencedBatches")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// TotalSequencedBatches is a free data retrieval call binding the contract method 0x066ec012.
//
// Solidity: function totalSequencedBatches() view returns(uint64)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerSession) TotalSequencedBatches() (uint64, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.TotalSequencedBatches(&_Mockfeijoapolygonrollupmanager.CallOpts)
}

// TotalSequencedBatches is a free data retrieval call binding the contract method 0x066ec012.
//
// Solidity: function totalSequencedBatches() view returns(uint64)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCallerSession) TotalSequencedBatches() (uint64, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.TotalSequencedBatches(&_Mockfeijoapolygonrollupmanager.CallOpts)
}

// TotalVerifiedBatches is a free data retrieval call binding the contract method 0xdde0ff77.
//
// Solidity: function totalVerifiedBatches() view returns(uint64)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCaller) TotalVerifiedBatches(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Mockfeijoapolygonrollupmanager.contract.Call(opts, &out, "totalVerifiedBatches")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// TotalVerifiedBatches is a free data retrieval call binding the contract method 0xdde0ff77.
//
// Solidity: function totalVerifiedBatches() view returns(uint64)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerSession) TotalVerifiedBatches() (uint64, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.TotalVerifiedBatches(&_Mockfeijoapolygonrollupmanager.CallOpts)
}

// TotalVerifiedBatches is a free data retrieval call binding the contract method 0xdde0ff77.
//
// Solidity: function totalVerifiedBatches() view returns(uint64)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCallerSession) TotalVerifiedBatches() (uint64, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.TotalVerifiedBatches(&_Mockfeijoapolygonrollupmanager.CallOpts)
}

// TrustedAggregatorTimeout is a free data retrieval call binding the contract method 0x841b24d7.
//
// Solidity: function trustedAggregatorTimeout() view returns(uint64)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCaller) TrustedAggregatorTimeout(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Mockfeijoapolygonrollupmanager.contract.Call(opts, &out, "trustedAggregatorTimeout")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// TrustedAggregatorTimeout is a free data retrieval call binding the contract method 0x841b24d7.
//
// Solidity: function trustedAggregatorTimeout() view returns(uint64)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerSession) TrustedAggregatorTimeout() (uint64, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.TrustedAggregatorTimeout(&_Mockfeijoapolygonrollupmanager.CallOpts)
}

// TrustedAggregatorTimeout is a free data retrieval call binding the contract method 0x841b24d7.
//
// Solidity: function trustedAggregatorTimeout() view returns(uint64)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCallerSession) TrustedAggregatorTimeout() (uint64, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.TrustedAggregatorTimeout(&_Mockfeijoapolygonrollupmanager.CallOpts)
}

// VerifyBatchTimeTarget is a free data retrieval call binding the contract method 0x0a0d9fbe.
//
// Solidity: function verifyBatchTimeTarget() view returns(uint64)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCaller) VerifyBatchTimeTarget(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _Mockfeijoapolygonrollupmanager.contract.Call(opts, &out, "verifyBatchTimeTarget")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// VerifyBatchTimeTarget is a free data retrieval call binding the contract method 0x0a0d9fbe.
//
// Solidity: function verifyBatchTimeTarget() view returns(uint64)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerSession) VerifyBatchTimeTarget() (uint64, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.VerifyBatchTimeTarget(&_Mockfeijoapolygonrollupmanager.CallOpts)
}

// VerifyBatchTimeTarget is a free data retrieval call binding the contract method 0x0a0d9fbe.
//
// Solidity: function verifyBatchTimeTarget() view returns(uint64)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerCallerSession) VerifyBatchTimeTarget() (uint64, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.VerifyBatchTimeTarget(&_Mockfeijoapolygonrollupmanager.CallOpts)
}

// ActivateEmergencyState is a paid mutator transaction binding the contract method 0x2072f6c5.
//
// Solidity: function activateEmergencyState() returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerTransactor) ActivateEmergencyState(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.contract.Transact(opts, "activateEmergencyState")
}

// ActivateEmergencyState is a paid mutator transaction binding the contract method 0x2072f6c5.
//
// Solidity: function activateEmergencyState() returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerSession) ActivateEmergencyState() (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.ActivateEmergencyState(&_Mockfeijoapolygonrollupmanager.TransactOpts)
}

// ActivateEmergencyState is a paid mutator transaction binding the contract method 0x2072f6c5.
//
// Solidity: function activateEmergencyState() returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerTransactorSession) ActivateEmergencyState() (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.ActivateEmergencyState(&_Mockfeijoapolygonrollupmanager.TransactOpts)
}

// AddExistingRollup is a paid mutator transaction binding the contract method 0xe0bfd3d2.
//
// Solidity: function addExistingRollup(address rollupAddress, address verifier, uint64 forkID, uint64 chainID, bytes32 genesis, uint8 rollupCompatibilityID) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerTransactor) AddExistingRollup(opts *bind.TransactOpts, rollupAddress common.Address, verifier common.Address, forkID uint64, chainID uint64, genesis [32]byte, rollupCompatibilityID uint8) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.contract.Transact(opts, "addExistingRollup", rollupAddress, verifier, forkID, chainID, genesis, rollupCompatibilityID)
}

// AddExistingRollup is a paid mutator transaction binding the contract method 0xe0bfd3d2.
//
// Solidity: function addExistingRollup(address rollupAddress, address verifier, uint64 forkID, uint64 chainID, bytes32 genesis, uint8 rollupCompatibilityID) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerSession) AddExistingRollup(rollupAddress common.Address, verifier common.Address, forkID uint64, chainID uint64, genesis [32]byte, rollupCompatibilityID uint8) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.AddExistingRollup(&_Mockfeijoapolygonrollupmanager.TransactOpts, rollupAddress, verifier, forkID, chainID, genesis, rollupCompatibilityID)
}

// AddExistingRollup is a paid mutator transaction binding the contract method 0xe0bfd3d2.
//
// Solidity: function addExistingRollup(address rollupAddress, address verifier, uint64 forkID, uint64 chainID, bytes32 genesis, uint8 rollupCompatibilityID) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerTransactorSession) AddExistingRollup(rollupAddress common.Address, verifier common.Address, forkID uint64, chainID uint64, genesis [32]byte, rollupCompatibilityID uint8) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.AddExistingRollup(&_Mockfeijoapolygonrollupmanager.TransactOpts, rollupAddress, verifier, forkID, chainID, genesis, rollupCompatibilityID)
}

// AddNewRollupType is a paid mutator transaction binding the contract method 0xf34eb8eb.
//
// Solidity: function addNewRollupType(address consensusImplementation, address verifier, uint64 forkID, uint8 rollupCompatibilityID, bytes32 genesis, string description) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerTransactor) AddNewRollupType(opts *bind.TransactOpts, consensusImplementation common.Address, verifier common.Address, forkID uint64, rollupCompatibilityID uint8, genesis [32]byte, description string) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.contract.Transact(opts, "addNewRollupType", consensusImplementation, verifier, forkID, rollupCompatibilityID, genesis, description)
}

// AddNewRollupType is a paid mutator transaction binding the contract method 0xf34eb8eb.
//
// Solidity: function addNewRollupType(address consensusImplementation, address verifier, uint64 forkID, uint8 rollupCompatibilityID, bytes32 genesis, string description) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerSession) AddNewRollupType(consensusImplementation common.Address, verifier common.Address, forkID uint64, rollupCompatibilityID uint8, genesis [32]byte, description string) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.AddNewRollupType(&_Mockfeijoapolygonrollupmanager.TransactOpts, consensusImplementation, verifier, forkID, rollupCompatibilityID, genesis, description)
}

// AddNewRollupType is a paid mutator transaction binding the contract method 0xf34eb8eb.
//
// Solidity: function addNewRollupType(address consensusImplementation, address verifier, uint64 forkID, uint8 rollupCompatibilityID, bytes32 genesis, string description) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerTransactorSession) AddNewRollupType(consensusImplementation common.Address, verifier common.Address, forkID uint64, rollupCompatibilityID uint8, genesis [32]byte, description string) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.AddNewRollupType(&_Mockfeijoapolygonrollupmanager.TransactOpts, consensusImplementation, verifier, forkID, rollupCompatibilityID, genesis, description)
}

// ConsolidatePendingState is a paid mutator transaction binding the contract method 0x1608859c.
//
// Solidity: function consolidatePendingState(uint32 rollupID, uint64 pendingStateNum) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerTransactor) ConsolidatePendingState(opts *bind.TransactOpts, rollupID uint32, pendingStateNum uint64) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.contract.Transact(opts, "consolidatePendingState", rollupID, pendingStateNum)
}

// ConsolidatePendingState is a paid mutator transaction binding the contract method 0x1608859c.
//
// Solidity: function consolidatePendingState(uint32 rollupID, uint64 pendingStateNum) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerSession) ConsolidatePendingState(rollupID uint32, pendingStateNum uint64) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.ConsolidatePendingState(&_Mockfeijoapolygonrollupmanager.TransactOpts, rollupID, pendingStateNum)
}

// ConsolidatePendingState is a paid mutator transaction binding the contract method 0x1608859c.
//
// Solidity: function consolidatePendingState(uint32 rollupID, uint64 pendingStateNum) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerTransactorSession) ConsolidatePendingState(rollupID uint32, pendingStateNum uint64) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.ConsolidatePendingState(&_Mockfeijoapolygonrollupmanager.TransactOpts, rollupID, pendingStateNum)
}

// CreateNewRollup is a paid mutator transaction binding the contract method 0x727885e9.
//
// Solidity: function createNewRollup(uint32 rollupTypeID, uint64 chainID, address admin, address sequencer, address gasTokenAddress, string sequencerURL, string networkName) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerTransactor) CreateNewRollup(opts *bind.TransactOpts, rollupTypeID uint32, chainID uint64, admin common.Address, sequencer common.Address, gasTokenAddress common.Address, sequencerURL string, networkName string) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.contract.Transact(opts, "createNewRollup", rollupTypeID, chainID, admin, sequencer, gasTokenAddress, sequencerURL, networkName)
}

// CreateNewRollup is a paid mutator transaction binding the contract method 0x727885e9.
//
// Solidity: function createNewRollup(uint32 rollupTypeID, uint64 chainID, address admin, address sequencer, address gasTokenAddress, string sequencerURL, string networkName) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerSession) CreateNewRollup(rollupTypeID uint32, chainID uint64, admin common.Address, sequencer common.Address, gasTokenAddress common.Address, sequencerURL string, networkName string) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.CreateNewRollup(&_Mockfeijoapolygonrollupmanager.TransactOpts, rollupTypeID, chainID, admin, sequencer, gasTokenAddress, sequencerURL, networkName)
}

// CreateNewRollup is a paid mutator transaction binding the contract method 0x727885e9.
//
// Solidity: function createNewRollup(uint32 rollupTypeID, uint64 chainID, address admin, address sequencer, address gasTokenAddress, string sequencerURL, string networkName) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerTransactorSession) CreateNewRollup(rollupTypeID uint32, chainID uint64, admin common.Address, sequencer common.Address, gasTokenAddress common.Address, sequencerURL string, networkName string) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.CreateNewRollup(&_Mockfeijoapolygonrollupmanager.TransactOpts, rollupTypeID, chainID, admin, sequencer, gasTokenAddress, sequencerURL, networkName)
}

// DeactivateEmergencyState is a paid mutator transaction binding the contract method 0xdbc16976.
//
// Solidity: function deactivateEmergencyState() returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerTransactor) DeactivateEmergencyState(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.contract.Transact(opts, "deactivateEmergencyState")
}

// DeactivateEmergencyState is a paid mutator transaction binding the contract method 0xdbc16976.
//
// Solidity: function deactivateEmergencyState() returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerSession) DeactivateEmergencyState() (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.DeactivateEmergencyState(&_Mockfeijoapolygonrollupmanager.TransactOpts)
}

// DeactivateEmergencyState is a paid mutator transaction binding the contract method 0xdbc16976.
//
// Solidity: function deactivateEmergencyState() returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerTransactorSession) DeactivateEmergencyState() (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.DeactivateEmergencyState(&_Mockfeijoapolygonrollupmanager.TransactOpts)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.GrantRole(&_Mockfeijoapolygonrollupmanager.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.GrantRole(&_Mockfeijoapolygonrollupmanager.TransactOpts, role, account)
}

// Initialize is a paid mutator transaction binding the contract method 0x0645af09.
//
// Solidity: function initialize(address trustedAggregator, uint64 _pendingStateTimeout, uint64 _trustedAggregatorTimeout, address admin, address timelock, address emergencyCouncil, address polygonZkEVM, address zkEVMVerifier, uint64 zkEVMForkID, uint64 zkEVMChainID) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerTransactor) Initialize(opts *bind.TransactOpts, trustedAggregator common.Address, _pendingStateTimeout uint64, _trustedAggregatorTimeout uint64, admin common.Address, timelock common.Address, emergencyCouncil common.Address, polygonZkEVM common.Address, zkEVMVerifier common.Address, zkEVMForkID uint64, zkEVMChainID uint64) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.contract.Transact(opts, "initialize", trustedAggregator, _pendingStateTimeout, _trustedAggregatorTimeout, admin, timelock, emergencyCouncil, polygonZkEVM, zkEVMVerifier, zkEVMForkID, zkEVMChainID)
}

// Initialize is a paid mutator transaction binding the contract method 0x0645af09.
//
// Solidity: function initialize(address trustedAggregator, uint64 _pendingStateTimeout, uint64 _trustedAggregatorTimeout, address admin, address timelock, address emergencyCouncil, address polygonZkEVM, address zkEVMVerifier, uint64 zkEVMForkID, uint64 zkEVMChainID) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerSession) Initialize(trustedAggregator common.Address, _pendingStateTimeout uint64, _trustedAggregatorTimeout uint64, admin common.Address, timelock common.Address, emergencyCouncil common.Address, polygonZkEVM common.Address, zkEVMVerifier common.Address, zkEVMForkID uint64, zkEVMChainID uint64) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.Initialize(&_Mockfeijoapolygonrollupmanager.TransactOpts, trustedAggregator, _pendingStateTimeout, _trustedAggregatorTimeout, admin, timelock, emergencyCouncil, polygonZkEVM, zkEVMVerifier, zkEVMForkID, zkEVMChainID)
}

// Initialize is a paid mutator transaction binding the contract method 0x0645af09.
//
// Solidity: function initialize(address trustedAggregator, uint64 _pendingStateTimeout, uint64 _trustedAggregatorTimeout, address admin, address timelock, address emergencyCouncil, address polygonZkEVM, address zkEVMVerifier, uint64 zkEVMForkID, uint64 zkEVMChainID) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerTransactorSession) Initialize(trustedAggregator common.Address, _pendingStateTimeout uint64, _trustedAggregatorTimeout uint64, admin common.Address, timelock common.Address, emergencyCouncil common.Address, polygonZkEVM common.Address, zkEVMVerifier common.Address, zkEVMForkID uint64, zkEVMChainID uint64) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.Initialize(&_Mockfeijoapolygonrollupmanager.TransactOpts, trustedAggregator, _pendingStateTimeout, _trustedAggregatorTimeout, admin, timelock, emergencyCouncil, polygonZkEVM, zkEVMVerifier, zkEVMForkID, zkEVMChainID)
}

// ObsoleteRollupType is a paid mutator transaction binding the contract method 0x7222020f.
//
// Solidity: function obsoleteRollupType(uint32 rollupTypeID) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerTransactor) ObsoleteRollupType(opts *bind.TransactOpts, rollupTypeID uint32) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.contract.Transact(opts, "obsoleteRollupType", rollupTypeID)
}

// ObsoleteRollupType is a paid mutator transaction binding the contract method 0x7222020f.
//
// Solidity: function obsoleteRollupType(uint32 rollupTypeID) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerSession) ObsoleteRollupType(rollupTypeID uint32) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.ObsoleteRollupType(&_Mockfeijoapolygonrollupmanager.TransactOpts, rollupTypeID)
}

// ObsoleteRollupType is a paid mutator transaction binding the contract method 0x7222020f.
//
// Solidity: function obsoleteRollupType(uint32 rollupTypeID) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerTransactorSession) ObsoleteRollupType(rollupTypeID uint32) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.ObsoleteRollupType(&_Mockfeijoapolygonrollupmanager.TransactOpts, rollupTypeID)
}

// OnSequenceBatches is a paid mutator transaction binding the contract method 0x9a908e73.
//
// Solidity: function onSequenceBatches(uint64 newSequencedBatches, bytes32 newAccInputHash) returns(uint64)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerTransactor) OnSequenceBatches(opts *bind.TransactOpts, newSequencedBatches uint64, newAccInputHash [32]byte) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.contract.Transact(opts, "onSequenceBatches", newSequencedBatches, newAccInputHash)
}

// OnSequenceBatches is a paid mutator transaction binding the contract method 0x9a908e73.
//
// Solidity: function onSequenceBatches(uint64 newSequencedBatches, bytes32 newAccInputHash) returns(uint64)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerSession) OnSequenceBatches(newSequencedBatches uint64, newAccInputHash [32]byte) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.OnSequenceBatches(&_Mockfeijoapolygonrollupmanager.TransactOpts, newSequencedBatches, newAccInputHash)
}

// OnSequenceBatches is a paid mutator transaction binding the contract method 0x9a908e73.
//
// Solidity: function onSequenceBatches(uint64 newSequencedBatches, bytes32 newAccInputHash) returns(uint64)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerTransactorSession) OnSequenceBatches(newSequencedBatches uint64, newAccInputHash [32]byte) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.OnSequenceBatches(&_Mockfeijoapolygonrollupmanager.TransactOpts, newSequencedBatches, newAccInputHash)
}

// OverridePendingState is a paid mutator transaction binding the contract method 0x12b86e19.
//
// Solidity: function overridePendingState(uint32 rollupID, uint64 initPendingStateNum, uint64 finalPendingStateNum, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 newStateRoot, bytes32[24] proof) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerTransactor) OverridePendingState(opts *bind.TransactOpts, rollupID uint32, initPendingStateNum uint64, finalPendingStateNum uint64, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, newStateRoot [32]byte, proof [24][32]byte) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.contract.Transact(opts, "overridePendingState", rollupID, initPendingStateNum, finalPendingStateNum, initNumBatch, finalNewBatch, newLocalExitRoot, newStateRoot, proof)
}

// OverridePendingState is a paid mutator transaction binding the contract method 0x12b86e19.
//
// Solidity: function overridePendingState(uint32 rollupID, uint64 initPendingStateNum, uint64 finalPendingStateNum, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 newStateRoot, bytes32[24] proof) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerSession) OverridePendingState(rollupID uint32, initPendingStateNum uint64, finalPendingStateNum uint64, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, newStateRoot [32]byte, proof [24][32]byte) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.OverridePendingState(&_Mockfeijoapolygonrollupmanager.TransactOpts, rollupID, initPendingStateNum, finalPendingStateNum, initNumBatch, finalNewBatch, newLocalExitRoot, newStateRoot, proof)
}

// OverridePendingState is a paid mutator transaction binding the contract method 0x12b86e19.
//
// Solidity: function overridePendingState(uint32 rollupID, uint64 initPendingStateNum, uint64 finalPendingStateNum, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 newStateRoot, bytes32[24] proof) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerTransactorSession) OverridePendingState(rollupID uint32, initPendingStateNum uint64, finalPendingStateNum uint64, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, newStateRoot [32]byte, proof [24][32]byte) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.OverridePendingState(&_Mockfeijoapolygonrollupmanager.TransactOpts, rollupID, initPendingStateNum, finalPendingStateNum, initNumBatch, finalNewBatch, newLocalExitRoot, newStateRoot, proof)
}

// ProveNonDeterministicPendingState is a paid mutator transaction binding the contract method 0x8bd4f071.
//
// Solidity: function proveNonDeterministicPendingState(uint32 rollupID, uint64 initPendingStateNum, uint64 finalPendingStateNum, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 newStateRoot, bytes32[24] proof) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerTransactor) ProveNonDeterministicPendingState(opts *bind.TransactOpts, rollupID uint32, initPendingStateNum uint64, finalPendingStateNum uint64, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, newStateRoot [32]byte, proof [24][32]byte) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.contract.Transact(opts, "proveNonDeterministicPendingState", rollupID, initPendingStateNum, finalPendingStateNum, initNumBatch, finalNewBatch, newLocalExitRoot, newStateRoot, proof)
}

// ProveNonDeterministicPendingState is a paid mutator transaction binding the contract method 0x8bd4f071.
//
// Solidity: function proveNonDeterministicPendingState(uint32 rollupID, uint64 initPendingStateNum, uint64 finalPendingStateNum, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 newStateRoot, bytes32[24] proof) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerSession) ProveNonDeterministicPendingState(rollupID uint32, initPendingStateNum uint64, finalPendingStateNum uint64, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, newStateRoot [32]byte, proof [24][32]byte) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.ProveNonDeterministicPendingState(&_Mockfeijoapolygonrollupmanager.TransactOpts, rollupID, initPendingStateNum, finalPendingStateNum, initNumBatch, finalNewBatch, newLocalExitRoot, newStateRoot, proof)
}

// ProveNonDeterministicPendingState is a paid mutator transaction binding the contract method 0x8bd4f071.
//
// Solidity: function proveNonDeterministicPendingState(uint32 rollupID, uint64 initPendingStateNum, uint64 finalPendingStateNum, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 newStateRoot, bytes32[24] proof) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerTransactorSession) ProveNonDeterministicPendingState(rollupID uint32, initPendingStateNum uint64, finalPendingStateNum uint64, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, newStateRoot [32]byte, proof [24][32]byte) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.ProveNonDeterministicPendingState(&_Mockfeijoapolygonrollupmanager.TransactOpts, rollupID, initPendingStateNum, finalPendingStateNum, initNumBatch, finalNewBatch, newLocalExitRoot, newStateRoot, proof)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.contract.Transact(opts, "renounceRole", role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.RenounceRole(&_Mockfeijoapolygonrollupmanager.TransactOpts, role, account)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address account) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerTransactorSession) RenounceRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.RenounceRole(&_Mockfeijoapolygonrollupmanager.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.RevokeRole(&_Mockfeijoapolygonrollupmanager.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.RevokeRole(&_Mockfeijoapolygonrollupmanager.TransactOpts, role, account)
}

// SetBatchFee is a paid mutator transaction binding the contract method 0xd5073f6f.
//
// Solidity: function setBatchFee(uint256 newBatchFee) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerTransactor) SetBatchFee(opts *bind.TransactOpts, newBatchFee *big.Int) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.contract.Transact(opts, "setBatchFee", newBatchFee)
}

// SetBatchFee is a paid mutator transaction binding the contract method 0xd5073f6f.
//
// Solidity: function setBatchFee(uint256 newBatchFee) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerSession) SetBatchFee(newBatchFee *big.Int) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.SetBatchFee(&_Mockfeijoapolygonrollupmanager.TransactOpts, newBatchFee)
}

// SetBatchFee is a paid mutator transaction binding the contract method 0xd5073f6f.
//
// Solidity: function setBatchFee(uint256 newBatchFee) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerTransactorSession) SetBatchFee(newBatchFee *big.Int) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.SetBatchFee(&_Mockfeijoapolygonrollupmanager.TransactOpts, newBatchFee)
}

// SetMultiplierBatchFee is a paid mutator transaction binding the contract method 0x1816b7e5.
//
// Solidity: function setMultiplierBatchFee(uint16 newMultiplierBatchFee) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerTransactor) SetMultiplierBatchFee(opts *bind.TransactOpts, newMultiplierBatchFee uint16) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.contract.Transact(opts, "setMultiplierBatchFee", newMultiplierBatchFee)
}

// SetMultiplierBatchFee is a paid mutator transaction binding the contract method 0x1816b7e5.
//
// Solidity: function setMultiplierBatchFee(uint16 newMultiplierBatchFee) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerSession) SetMultiplierBatchFee(newMultiplierBatchFee uint16) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.SetMultiplierBatchFee(&_Mockfeijoapolygonrollupmanager.TransactOpts, newMultiplierBatchFee)
}

// SetMultiplierBatchFee is a paid mutator transaction binding the contract method 0x1816b7e5.
//
// Solidity: function setMultiplierBatchFee(uint16 newMultiplierBatchFee) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerTransactorSession) SetMultiplierBatchFee(newMultiplierBatchFee uint16) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.SetMultiplierBatchFee(&_Mockfeijoapolygonrollupmanager.TransactOpts, newMultiplierBatchFee)
}

// SetPendingStateTimeout is a paid mutator transaction binding the contract method 0x9c9f3dfe.
//
// Solidity: function setPendingStateTimeout(uint64 newPendingStateTimeout) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerTransactor) SetPendingStateTimeout(opts *bind.TransactOpts, newPendingStateTimeout uint64) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.contract.Transact(opts, "setPendingStateTimeout", newPendingStateTimeout)
}

// SetPendingStateTimeout is a paid mutator transaction binding the contract method 0x9c9f3dfe.
//
// Solidity: function setPendingStateTimeout(uint64 newPendingStateTimeout) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerSession) SetPendingStateTimeout(newPendingStateTimeout uint64) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.SetPendingStateTimeout(&_Mockfeijoapolygonrollupmanager.TransactOpts, newPendingStateTimeout)
}

// SetPendingStateTimeout is a paid mutator transaction binding the contract method 0x9c9f3dfe.
//
// Solidity: function setPendingStateTimeout(uint64 newPendingStateTimeout) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerTransactorSession) SetPendingStateTimeout(newPendingStateTimeout uint64) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.SetPendingStateTimeout(&_Mockfeijoapolygonrollupmanager.TransactOpts, newPendingStateTimeout)
}

// SetTrustedAggregatorTimeout is a paid mutator transaction binding the contract method 0x394218e9.
//
// Solidity: function setTrustedAggregatorTimeout(uint64 newTrustedAggregatorTimeout) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerTransactor) SetTrustedAggregatorTimeout(opts *bind.TransactOpts, newTrustedAggregatorTimeout uint64) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.contract.Transact(opts, "setTrustedAggregatorTimeout", newTrustedAggregatorTimeout)
}

// SetTrustedAggregatorTimeout is a paid mutator transaction binding the contract method 0x394218e9.
//
// Solidity: function setTrustedAggregatorTimeout(uint64 newTrustedAggregatorTimeout) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerSession) SetTrustedAggregatorTimeout(newTrustedAggregatorTimeout uint64) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.SetTrustedAggregatorTimeout(&_Mockfeijoapolygonrollupmanager.TransactOpts, newTrustedAggregatorTimeout)
}

// SetTrustedAggregatorTimeout is a paid mutator transaction binding the contract method 0x394218e9.
//
// Solidity: function setTrustedAggregatorTimeout(uint64 newTrustedAggregatorTimeout) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerTransactorSession) SetTrustedAggregatorTimeout(newTrustedAggregatorTimeout uint64) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.SetTrustedAggregatorTimeout(&_Mockfeijoapolygonrollupmanager.TransactOpts, newTrustedAggregatorTimeout)
}

// SetVerifyBatchTimeTarget is a paid mutator transaction binding the contract method 0xa066215c.
//
// Solidity: function setVerifyBatchTimeTarget(uint64 newVerifyBatchTimeTarget) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerTransactor) SetVerifyBatchTimeTarget(opts *bind.TransactOpts, newVerifyBatchTimeTarget uint64) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.contract.Transact(opts, "setVerifyBatchTimeTarget", newVerifyBatchTimeTarget)
}

// SetVerifyBatchTimeTarget is a paid mutator transaction binding the contract method 0xa066215c.
//
// Solidity: function setVerifyBatchTimeTarget(uint64 newVerifyBatchTimeTarget) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerSession) SetVerifyBatchTimeTarget(newVerifyBatchTimeTarget uint64) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.SetVerifyBatchTimeTarget(&_Mockfeijoapolygonrollupmanager.TransactOpts, newVerifyBatchTimeTarget)
}

// SetVerifyBatchTimeTarget is a paid mutator transaction binding the contract method 0xa066215c.
//
// Solidity: function setVerifyBatchTimeTarget(uint64 newVerifyBatchTimeTarget) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerTransactorSession) SetVerifyBatchTimeTarget(newVerifyBatchTimeTarget uint64) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.SetVerifyBatchTimeTarget(&_Mockfeijoapolygonrollupmanager.TransactOpts, newVerifyBatchTimeTarget)
}

// UpdateRollup is a paid mutator transaction binding the contract method 0xc4c928c2.
//
// Solidity: function updateRollup(address rollupContract, uint32 newRollupTypeID, bytes upgradeData) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerTransactor) UpdateRollup(opts *bind.TransactOpts, rollupContract common.Address, newRollupTypeID uint32, upgradeData []byte) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.contract.Transact(opts, "updateRollup", rollupContract, newRollupTypeID, upgradeData)
}

// UpdateRollup is a paid mutator transaction binding the contract method 0xc4c928c2.
//
// Solidity: function updateRollup(address rollupContract, uint32 newRollupTypeID, bytes upgradeData) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerSession) UpdateRollup(rollupContract common.Address, newRollupTypeID uint32, upgradeData []byte) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.UpdateRollup(&_Mockfeijoapolygonrollupmanager.TransactOpts, rollupContract, newRollupTypeID, upgradeData)
}

// UpdateRollup is a paid mutator transaction binding the contract method 0xc4c928c2.
//
// Solidity: function updateRollup(address rollupContract, uint32 newRollupTypeID, bytes upgradeData) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerTransactorSession) UpdateRollup(rollupContract common.Address, newRollupTypeID uint32, upgradeData []byte) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.UpdateRollup(&_Mockfeijoapolygonrollupmanager.TransactOpts, rollupContract, newRollupTypeID, upgradeData)
}

// VerifyBatches is a paid mutator transaction binding the contract method 0x87c20c01.
//
// Solidity: function verifyBatches(uint32 rollupID, uint64 pendingStateNum, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 newStateRoot, address beneficiary, bytes32[24] proof) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerTransactor) VerifyBatches(opts *bind.TransactOpts, rollupID uint32, pendingStateNum uint64, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, newStateRoot [32]byte, beneficiary common.Address, proof [24][32]byte) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.contract.Transact(opts, "verifyBatches", rollupID, pendingStateNum, initNumBatch, finalNewBatch, newLocalExitRoot, newStateRoot, beneficiary, proof)
}

// VerifyBatches is a paid mutator transaction binding the contract method 0x87c20c01.
//
// Solidity: function verifyBatches(uint32 rollupID, uint64 pendingStateNum, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 newStateRoot, address beneficiary, bytes32[24] proof) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerSession) VerifyBatches(rollupID uint32, pendingStateNum uint64, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, newStateRoot [32]byte, beneficiary common.Address, proof [24][32]byte) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.VerifyBatches(&_Mockfeijoapolygonrollupmanager.TransactOpts, rollupID, pendingStateNum, initNumBatch, finalNewBatch, newLocalExitRoot, newStateRoot, beneficiary, proof)
}

// VerifyBatches is a paid mutator transaction binding the contract method 0x87c20c01.
//
// Solidity: function verifyBatches(uint32 rollupID, uint64 pendingStateNum, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 newStateRoot, address beneficiary, bytes32[24] proof) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerTransactorSession) VerifyBatches(rollupID uint32, pendingStateNum uint64, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, newStateRoot [32]byte, beneficiary common.Address, proof [24][32]byte) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.VerifyBatches(&_Mockfeijoapolygonrollupmanager.TransactOpts, rollupID, pendingStateNum, initNumBatch, finalNewBatch, newLocalExitRoot, newStateRoot, beneficiary, proof)
}

// VerifyBatchesTrustedAggregator is a paid mutator transaction binding the contract method 0x1489ed10.
//
// Solidity: function verifyBatchesTrustedAggregator(uint32 rollupID, uint64 pendingStateNum, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 newStateRoot, address beneficiary, bytes32[24] proof) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerTransactor) VerifyBatchesTrustedAggregator(opts *bind.TransactOpts, rollupID uint32, pendingStateNum uint64, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, newStateRoot [32]byte, beneficiary common.Address, proof [24][32]byte) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.contract.Transact(opts, "verifyBatchesTrustedAggregator", rollupID, pendingStateNum, initNumBatch, finalNewBatch, newLocalExitRoot, newStateRoot, beneficiary, proof)
}

// VerifyBatchesTrustedAggregator is a paid mutator transaction binding the contract method 0x1489ed10.
//
// Solidity: function verifyBatchesTrustedAggregator(uint32 rollupID, uint64 pendingStateNum, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 newStateRoot, address beneficiary, bytes32[24] proof) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerSession) VerifyBatchesTrustedAggregator(rollupID uint32, pendingStateNum uint64, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, newStateRoot [32]byte, beneficiary common.Address, proof [24][32]byte) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.VerifyBatchesTrustedAggregator(&_Mockfeijoapolygonrollupmanager.TransactOpts, rollupID, pendingStateNum, initNumBatch, finalNewBatch, newLocalExitRoot, newStateRoot, beneficiary, proof)
}

// VerifyBatchesTrustedAggregator is a paid mutator transaction binding the contract method 0x1489ed10.
//
// Solidity: function verifyBatchesTrustedAggregator(uint32 rollupID, uint64 pendingStateNum, uint64 initNumBatch, uint64 finalNewBatch, bytes32 newLocalExitRoot, bytes32 newStateRoot, address beneficiary, bytes32[24] proof) returns()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerTransactorSession) VerifyBatchesTrustedAggregator(rollupID uint32, pendingStateNum uint64, initNumBatch uint64, finalNewBatch uint64, newLocalExitRoot [32]byte, newStateRoot [32]byte, beneficiary common.Address, proof [24][32]byte) (*types.Transaction, error) {
	return _Mockfeijoapolygonrollupmanager.Contract.VerifyBatchesTrustedAggregator(&_Mockfeijoapolygonrollupmanager.TransactOpts, rollupID, pendingStateNum, initNumBatch, finalNewBatch, newLocalExitRoot, newStateRoot, beneficiary, proof)
}

// MockfeijoapolygonrollupmanagerAddExistingRollupIterator is returned from FilterAddExistingRollup and is used to iterate over the raw logs and unpacked data for AddExistingRollup events raised by the Mockfeijoapolygonrollupmanager contract.
type MockfeijoapolygonrollupmanagerAddExistingRollupIterator struct {
	Event *MockfeijoapolygonrollupmanagerAddExistingRollup // Event containing the contract specifics and raw log

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
func (it *MockfeijoapolygonrollupmanagerAddExistingRollupIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockfeijoapolygonrollupmanagerAddExistingRollup)
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
		it.Event = new(MockfeijoapolygonrollupmanagerAddExistingRollup)
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
func (it *MockfeijoapolygonrollupmanagerAddExistingRollupIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockfeijoapolygonrollupmanagerAddExistingRollupIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockfeijoapolygonrollupmanagerAddExistingRollup represents a AddExistingRollup event raised by the Mockfeijoapolygonrollupmanager contract.
type MockfeijoapolygonrollupmanagerAddExistingRollup struct {
	RollupID                       uint32
	ForkID                         uint64
	RollupAddress                  common.Address
	ChainID                        uint64
	RollupCompatibilityID          uint8
	LastVerifiedBatchBeforeUpgrade uint64
	Raw                            types.Log // Blockchain specific contextual infos
}

// FilterAddExistingRollup is a free log retrieval operation binding the contract event 0xadfc7d56f7e39b08b321534f14bfb135ad27698f7d2f5ad0edc2356ea9a3f850.
//
// Solidity: event AddExistingRollup(uint32 indexed rollupID, uint64 forkID, address rollupAddress, uint64 chainID, uint8 rollupCompatibilityID, uint64 lastVerifiedBatchBeforeUpgrade)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) FilterAddExistingRollup(opts *bind.FilterOpts, rollupID []uint32) (*MockfeijoapolygonrollupmanagerAddExistingRollupIterator, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	logs, sub, err := _Mockfeijoapolygonrollupmanager.contract.FilterLogs(opts, "AddExistingRollup", rollupIDRule)
	if err != nil {
		return nil, err
	}
	return &MockfeijoapolygonrollupmanagerAddExistingRollupIterator{contract: _Mockfeijoapolygonrollupmanager.contract, event: "AddExistingRollup", logs: logs, sub: sub}, nil
}

// WatchAddExistingRollup is a free log subscription operation binding the contract event 0xadfc7d56f7e39b08b321534f14bfb135ad27698f7d2f5ad0edc2356ea9a3f850.
//
// Solidity: event AddExistingRollup(uint32 indexed rollupID, uint64 forkID, address rollupAddress, uint64 chainID, uint8 rollupCompatibilityID, uint64 lastVerifiedBatchBeforeUpgrade)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) WatchAddExistingRollup(opts *bind.WatchOpts, sink chan<- *MockfeijoapolygonrollupmanagerAddExistingRollup, rollupID []uint32) (event.Subscription, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	logs, sub, err := _Mockfeijoapolygonrollupmanager.contract.WatchLogs(opts, "AddExistingRollup", rollupIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockfeijoapolygonrollupmanagerAddExistingRollup)
				if err := _Mockfeijoapolygonrollupmanager.contract.UnpackLog(event, "AddExistingRollup", log); err != nil {
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
// Solidity: event AddExistingRollup(uint32 indexed rollupID, uint64 forkID, address rollupAddress, uint64 chainID, uint8 rollupCompatibilityID, uint64 lastVerifiedBatchBeforeUpgrade)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) ParseAddExistingRollup(log types.Log) (*MockfeijoapolygonrollupmanagerAddExistingRollup, error) {
	event := new(MockfeijoapolygonrollupmanagerAddExistingRollup)
	if err := _Mockfeijoapolygonrollupmanager.contract.UnpackLog(event, "AddExistingRollup", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockfeijoapolygonrollupmanagerAddNewRollupTypeIterator is returned from FilterAddNewRollupType and is used to iterate over the raw logs and unpacked data for AddNewRollupType events raised by the Mockfeijoapolygonrollupmanager contract.
type MockfeijoapolygonrollupmanagerAddNewRollupTypeIterator struct {
	Event *MockfeijoapolygonrollupmanagerAddNewRollupType // Event containing the contract specifics and raw log

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
func (it *MockfeijoapolygonrollupmanagerAddNewRollupTypeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockfeijoapolygonrollupmanagerAddNewRollupType)
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
		it.Event = new(MockfeijoapolygonrollupmanagerAddNewRollupType)
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
func (it *MockfeijoapolygonrollupmanagerAddNewRollupTypeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockfeijoapolygonrollupmanagerAddNewRollupTypeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockfeijoapolygonrollupmanagerAddNewRollupType represents a AddNewRollupType event raised by the Mockfeijoapolygonrollupmanager contract.
type MockfeijoapolygonrollupmanagerAddNewRollupType struct {
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
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) FilterAddNewRollupType(opts *bind.FilterOpts, rollupTypeID []uint32) (*MockfeijoapolygonrollupmanagerAddNewRollupTypeIterator, error) {

	var rollupTypeIDRule []interface{}
	for _, rollupTypeIDItem := range rollupTypeID {
		rollupTypeIDRule = append(rollupTypeIDRule, rollupTypeIDItem)
	}

	logs, sub, err := _Mockfeijoapolygonrollupmanager.contract.FilterLogs(opts, "AddNewRollupType", rollupTypeIDRule)
	if err != nil {
		return nil, err
	}
	return &MockfeijoapolygonrollupmanagerAddNewRollupTypeIterator{contract: _Mockfeijoapolygonrollupmanager.contract, event: "AddNewRollupType", logs: logs, sub: sub}, nil
}

// WatchAddNewRollupType is a free log subscription operation binding the contract event 0xa2970448b3bd66ba7e524e7b2a5b9cf94fa29e32488fb942afdfe70dd4b77b52.
//
// Solidity: event AddNewRollupType(uint32 indexed rollupTypeID, address consensusImplementation, address verifier, uint64 forkID, uint8 rollupCompatibilityID, bytes32 genesis, string description)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) WatchAddNewRollupType(opts *bind.WatchOpts, sink chan<- *MockfeijoapolygonrollupmanagerAddNewRollupType, rollupTypeID []uint32) (event.Subscription, error) {

	var rollupTypeIDRule []interface{}
	for _, rollupTypeIDItem := range rollupTypeID {
		rollupTypeIDRule = append(rollupTypeIDRule, rollupTypeIDItem)
	}

	logs, sub, err := _Mockfeijoapolygonrollupmanager.contract.WatchLogs(opts, "AddNewRollupType", rollupTypeIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockfeijoapolygonrollupmanagerAddNewRollupType)
				if err := _Mockfeijoapolygonrollupmanager.contract.UnpackLog(event, "AddNewRollupType", log); err != nil {
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
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) ParseAddNewRollupType(log types.Log) (*MockfeijoapolygonrollupmanagerAddNewRollupType, error) {
	event := new(MockfeijoapolygonrollupmanagerAddNewRollupType)
	if err := _Mockfeijoapolygonrollupmanager.contract.UnpackLog(event, "AddNewRollupType", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockfeijoapolygonrollupmanagerConsolidatePendingStateIterator is returned from FilterConsolidatePendingState and is used to iterate over the raw logs and unpacked data for ConsolidatePendingState events raised by the Mockfeijoapolygonrollupmanager contract.
type MockfeijoapolygonrollupmanagerConsolidatePendingStateIterator struct {
	Event *MockfeijoapolygonrollupmanagerConsolidatePendingState // Event containing the contract specifics and raw log

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
func (it *MockfeijoapolygonrollupmanagerConsolidatePendingStateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockfeijoapolygonrollupmanagerConsolidatePendingState)
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
		it.Event = new(MockfeijoapolygonrollupmanagerConsolidatePendingState)
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
func (it *MockfeijoapolygonrollupmanagerConsolidatePendingStateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockfeijoapolygonrollupmanagerConsolidatePendingStateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockfeijoapolygonrollupmanagerConsolidatePendingState represents a ConsolidatePendingState event raised by the Mockfeijoapolygonrollupmanager contract.
type MockfeijoapolygonrollupmanagerConsolidatePendingState struct {
	RollupID        uint32
	NumBatch        uint64
	StateRoot       [32]byte
	ExitRoot        [32]byte
	PendingStateNum uint64
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterConsolidatePendingState is a free log retrieval operation binding the contract event 0x581910eb7a27738945c2f00a91f2284b2d6de9d4e472b12f901c2b0df045e21b.
//
// Solidity: event ConsolidatePendingState(uint32 indexed rollupID, uint64 numBatch, bytes32 stateRoot, bytes32 exitRoot, uint64 pendingStateNum)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) FilterConsolidatePendingState(opts *bind.FilterOpts, rollupID []uint32) (*MockfeijoapolygonrollupmanagerConsolidatePendingStateIterator, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	logs, sub, err := _Mockfeijoapolygonrollupmanager.contract.FilterLogs(opts, "ConsolidatePendingState", rollupIDRule)
	if err != nil {
		return nil, err
	}
	return &MockfeijoapolygonrollupmanagerConsolidatePendingStateIterator{contract: _Mockfeijoapolygonrollupmanager.contract, event: "ConsolidatePendingState", logs: logs, sub: sub}, nil
}

// WatchConsolidatePendingState is a free log subscription operation binding the contract event 0x581910eb7a27738945c2f00a91f2284b2d6de9d4e472b12f901c2b0df045e21b.
//
// Solidity: event ConsolidatePendingState(uint32 indexed rollupID, uint64 numBatch, bytes32 stateRoot, bytes32 exitRoot, uint64 pendingStateNum)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) WatchConsolidatePendingState(opts *bind.WatchOpts, sink chan<- *MockfeijoapolygonrollupmanagerConsolidatePendingState, rollupID []uint32) (event.Subscription, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	logs, sub, err := _Mockfeijoapolygonrollupmanager.contract.WatchLogs(opts, "ConsolidatePendingState", rollupIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockfeijoapolygonrollupmanagerConsolidatePendingState)
				if err := _Mockfeijoapolygonrollupmanager.contract.UnpackLog(event, "ConsolidatePendingState", log); err != nil {
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
// Solidity: event ConsolidatePendingState(uint32 indexed rollupID, uint64 numBatch, bytes32 stateRoot, bytes32 exitRoot, uint64 pendingStateNum)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) ParseConsolidatePendingState(log types.Log) (*MockfeijoapolygonrollupmanagerConsolidatePendingState, error) {
	event := new(MockfeijoapolygonrollupmanagerConsolidatePendingState)
	if err := _Mockfeijoapolygonrollupmanager.contract.UnpackLog(event, "ConsolidatePendingState", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockfeijoapolygonrollupmanagerCreateNewRollupIterator is returned from FilterCreateNewRollup and is used to iterate over the raw logs and unpacked data for CreateNewRollup events raised by the Mockfeijoapolygonrollupmanager contract.
type MockfeijoapolygonrollupmanagerCreateNewRollupIterator struct {
	Event *MockfeijoapolygonrollupmanagerCreateNewRollup // Event containing the contract specifics and raw log

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
func (it *MockfeijoapolygonrollupmanagerCreateNewRollupIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockfeijoapolygonrollupmanagerCreateNewRollup)
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
		it.Event = new(MockfeijoapolygonrollupmanagerCreateNewRollup)
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
func (it *MockfeijoapolygonrollupmanagerCreateNewRollupIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockfeijoapolygonrollupmanagerCreateNewRollupIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockfeijoapolygonrollupmanagerCreateNewRollup represents a CreateNewRollup event raised by the Mockfeijoapolygonrollupmanager contract.
type MockfeijoapolygonrollupmanagerCreateNewRollup struct {
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
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) FilterCreateNewRollup(opts *bind.FilterOpts, rollupID []uint32) (*MockfeijoapolygonrollupmanagerCreateNewRollupIterator, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	logs, sub, err := _Mockfeijoapolygonrollupmanager.contract.FilterLogs(opts, "CreateNewRollup", rollupIDRule)
	if err != nil {
		return nil, err
	}
	return &MockfeijoapolygonrollupmanagerCreateNewRollupIterator{contract: _Mockfeijoapolygonrollupmanager.contract, event: "CreateNewRollup", logs: logs, sub: sub}, nil
}

// WatchCreateNewRollup is a free log subscription operation binding the contract event 0x194c983456df6701c6a50830b90fe80e72b823411d0d524970c9590dc277a641.
//
// Solidity: event CreateNewRollup(uint32 indexed rollupID, uint32 rollupTypeID, address rollupAddress, uint64 chainID, address gasTokenAddress)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) WatchCreateNewRollup(opts *bind.WatchOpts, sink chan<- *MockfeijoapolygonrollupmanagerCreateNewRollup, rollupID []uint32) (event.Subscription, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	logs, sub, err := _Mockfeijoapolygonrollupmanager.contract.WatchLogs(opts, "CreateNewRollup", rollupIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockfeijoapolygonrollupmanagerCreateNewRollup)
				if err := _Mockfeijoapolygonrollupmanager.contract.UnpackLog(event, "CreateNewRollup", log); err != nil {
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
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) ParseCreateNewRollup(log types.Log) (*MockfeijoapolygonrollupmanagerCreateNewRollup, error) {
	event := new(MockfeijoapolygonrollupmanagerCreateNewRollup)
	if err := _Mockfeijoapolygonrollupmanager.contract.UnpackLog(event, "CreateNewRollup", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockfeijoapolygonrollupmanagerEmergencyStateActivatedIterator is returned from FilterEmergencyStateActivated and is used to iterate over the raw logs and unpacked data for EmergencyStateActivated events raised by the Mockfeijoapolygonrollupmanager contract.
type MockfeijoapolygonrollupmanagerEmergencyStateActivatedIterator struct {
	Event *MockfeijoapolygonrollupmanagerEmergencyStateActivated // Event containing the contract specifics and raw log

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
func (it *MockfeijoapolygonrollupmanagerEmergencyStateActivatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockfeijoapolygonrollupmanagerEmergencyStateActivated)
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
		it.Event = new(MockfeijoapolygonrollupmanagerEmergencyStateActivated)
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
func (it *MockfeijoapolygonrollupmanagerEmergencyStateActivatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockfeijoapolygonrollupmanagerEmergencyStateActivatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockfeijoapolygonrollupmanagerEmergencyStateActivated represents a EmergencyStateActivated event raised by the Mockfeijoapolygonrollupmanager contract.
type MockfeijoapolygonrollupmanagerEmergencyStateActivated struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterEmergencyStateActivated is a free log retrieval operation binding the contract event 0x2261efe5aef6fedc1fd1550b25facc9181745623049c7901287030b9ad1a5497.
//
// Solidity: event EmergencyStateActivated()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) FilterEmergencyStateActivated(opts *bind.FilterOpts) (*MockfeijoapolygonrollupmanagerEmergencyStateActivatedIterator, error) {

	logs, sub, err := _Mockfeijoapolygonrollupmanager.contract.FilterLogs(opts, "EmergencyStateActivated")
	if err != nil {
		return nil, err
	}
	return &MockfeijoapolygonrollupmanagerEmergencyStateActivatedIterator{contract: _Mockfeijoapolygonrollupmanager.contract, event: "EmergencyStateActivated", logs: logs, sub: sub}, nil
}

// WatchEmergencyStateActivated is a free log subscription operation binding the contract event 0x2261efe5aef6fedc1fd1550b25facc9181745623049c7901287030b9ad1a5497.
//
// Solidity: event EmergencyStateActivated()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) WatchEmergencyStateActivated(opts *bind.WatchOpts, sink chan<- *MockfeijoapolygonrollupmanagerEmergencyStateActivated) (event.Subscription, error) {

	logs, sub, err := _Mockfeijoapolygonrollupmanager.contract.WatchLogs(opts, "EmergencyStateActivated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockfeijoapolygonrollupmanagerEmergencyStateActivated)
				if err := _Mockfeijoapolygonrollupmanager.contract.UnpackLog(event, "EmergencyStateActivated", log); err != nil {
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
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) ParseEmergencyStateActivated(log types.Log) (*MockfeijoapolygonrollupmanagerEmergencyStateActivated, error) {
	event := new(MockfeijoapolygonrollupmanagerEmergencyStateActivated)
	if err := _Mockfeijoapolygonrollupmanager.contract.UnpackLog(event, "EmergencyStateActivated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockfeijoapolygonrollupmanagerEmergencyStateDeactivatedIterator is returned from FilterEmergencyStateDeactivated and is used to iterate over the raw logs and unpacked data for EmergencyStateDeactivated events raised by the Mockfeijoapolygonrollupmanager contract.
type MockfeijoapolygonrollupmanagerEmergencyStateDeactivatedIterator struct {
	Event *MockfeijoapolygonrollupmanagerEmergencyStateDeactivated // Event containing the contract specifics and raw log

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
func (it *MockfeijoapolygonrollupmanagerEmergencyStateDeactivatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockfeijoapolygonrollupmanagerEmergencyStateDeactivated)
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
		it.Event = new(MockfeijoapolygonrollupmanagerEmergencyStateDeactivated)
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
func (it *MockfeijoapolygonrollupmanagerEmergencyStateDeactivatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockfeijoapolygonrollupmanagerEmergencyStateDeactivatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockfeijoapolygonrollupmanagerEmergencyStateDeactivated represents a EmergencyStateDeactivated event raised by the Mockfeijoapolygonrollupmanager contract.
type MockfeijoapolygonrollupmanagerEmergencyStateDeactivated struct {
	Raw types.Log // Blockchain specific contextual infos
}

// FilterEmergencyStateDeactivated is a free log retrieval operation binding the contract event 0x1e5e34eea33501aecf2ebec9fe0e884a40804275ea7fe10b2ba084c8374308b3.
//
// Solidity: event EmergencyStateDeactivated()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) FilterEmergencyStateDeactivated(opts *bind.FilterOpts) (*MockfeijoapolygonrollupmanagerEmergencyStateDeactivatedIterator, error) {

	logs, sub, err := _Mockfeijoapolygonrollupmanager.contract.FilterLogs(opts, "EmergencyStateDeactivated")
	if err != nil {
		return nil, err
	}
	return &MockfeijoapolygonrollupmanagerEmergencyStateDeactivatedIterator{contract: _Mockfeijoapolygonrollupmanager.contract, event: "EmergencyStateDeactivated", logs: logs, sub: sub}, nil
}

// WatchEmergencyStateDeactivated is a free log subscription operation binding the contract event 0x1e5e34eea33501aecf2ebec9fe0e884a40804275ea7fe10b2ba084c8374308b3.
//
// Solidity: event EmergencyStateDeactivated()
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) WatchEmergencyStateDeactivated(opts *bind.WatchOpts, sink chan<- *MockfeijoapolygonrollupmanagerEmergencyStateDeactivated) (event.Subscription, error) {

	logs, sub, err := _Mockfeijoapolygonrollupmanager.contract.WatchLogs(opts, "EmergencyStateDeactivated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockfeijoapolygonrollupmanagerEmergencyStateDeactivated)
				if err := _Mockfeijoapolygonrollupmanager.contract.UnpackLog(event, "EmergencyStateDeactivated", log); err != nil {
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
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) ParseEmergencyStateDeactivated(log types.Log) (*MockfeijoapolygonrollupmanagerEmergencyStateDeactivated, error) {
	event := new(MockfeijoapolygonrollupmanagerEmergencyStateDeactivated)
	if err := _Mockfeijoapolygonrollupmanager.contract.UnpackLog(event, "EmergencyStateDeactivated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockfeijoapolygonrollupmanagerInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the Mockfeijoapolygonrollupmanager contract.
type MockfeijoapolygonrollupmanagerInitializedIterator struct {
	Event *MockfeijoapolygonrollupmanagerInitialized // Event containing the contract specifics and raw log

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
func (it *MockfeijoapolygonrollupmanagerInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockfeijoapolygonrollupmanagerInitialized)
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
		it.Event = new(MockfeijoapolygonrollupmanagerInitialized)
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
func (it *MockfeijoapolygonrollupmanagerInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockfeijoapolygonrollupmanagerInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockfeijoapolygonrollupmanagerInitialized represents a Initialized event raised by the Mockfeijoapolygonrollupmanager contract.
type MockfeijoapolygonrollupmanagerInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) FilterInitialized(opts *bind.FilterOpts) (*MockfeijoapolygonrollupmanagerInitializedIterator, error) {

	logs, sub, err := _Mockfeijoapolygonrollupmanager.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &MockfeijoapolygonrollupmanagerInitializedIterator{contract: _Mockfeijoapolygonrollupmanager.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *MockfeijoapolygonrollupmanagerInitialized) (event.Subscription, error) {

	logs, sub, err := _Mockfeijoapolygonrollupmanager.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockfeijoapolygonrollupmanagerInitialized)
				if err := _Mockfeijoapolygonrollupmanager.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) ParseInitialized(log types.Log) (*MockfeijoapolygonrollupmanagerInitialized, error) {
	event := new(MockfeijoapolygonrollupmanagerInitialized)
	if err := _Mockfeijoapolygonrollupmanager.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockfeijoapolygonrollupmanagerObsoleteRollupTypeIterator is returned from FilterObsoleteRollupType and is used to iterate over the raw logs and unpacked data for ObsoleteRollupType events raised by the Mockfeijoapolygonrollupmanager contract.
type MockfeijoapolygonrollupmanagerObsoleteRollupTypeIterator struct {
	Event *MockfeijoapolygonrollupmanagerObsoleteRollupType // Event containing the contract specifics and raw log

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
func (it *MockfeijoapolygonrollupmanagerObsoleteRollupTypeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockfeijoapolygonrollupmanagerObsoleteRollupType)
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
		it.Event = new(MockfeijoapolygonrollupmanagerObsoleteRollupType)
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
func (it *MockfeijoapolygonrollupmanagerObsoleteRollupTypeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockfeijoapolygonrollupmanagerObsoleteRollupTypeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockfeijoapolygonrollupmanagerObsoleteRollupType represents a ObsoleteRollupType event raised by the Mockfeijoapolygonrollupmanager contract.
type MockfeijoapolygonrollupmanagerObsoleteRollupType struct {
	RollupTypeID uint32
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterObsoleteRollupType is a free log retrieval operation binding the contract event 0x4710d2ee567ef1ed6eb2f651dde4589524bcf7cebc62147a99b281cc836e7e44.
//
// Solidity: event ObsoleteRollupType(uint32 indexed rollupTypeID)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) FilterObsoleteRollupType(opts *bind.FilterOpts, rollupTypeID []uint32) (*MockfeijoapolygonrollupmanagerObsoleteRollupTypeIterator, error) {

	var rollupTypeIDRule []interface{}
	for _, rollupTypeIDItem := range rollupTypeID {
		rollupTypeIDRule = append(rollupTypeIDRule, rollupTypeIDItem)
	}

	logs, sub, err := _Mockfeijoapolygonrollupmanager.contract.FilterLogs(opts, "ObsoleteRollupType", rollupTypeIDRule)
	if err != nil {
		return nil, err
	}
	return &MockfeijoapolygonrollupmanagerObsoleteRollupTypeIterator{contract: _Mockfeijoapolygonrollupmanager.contract, event: "ObsoleteRollupType", logs: logs, sub: sub}, nil
}

// WatchObsoleteRollupType is a free log subscription operation binding the contract event 0x4710d2ee567ef1ed6eb2f651dde4589524bcf7cebc62147a99b281cc836e7e44.
//
// Solidity: event ObsoleteRollupType(uint32 indexed rollupTypeID)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) WatchObsoleteRollupType(opts *bind.WatchOpts, sink chan<- *MockfeijoapolygonrollupmanagerObsoleteRollupType, rollupTypeID []uint32) (event.Subscription, error) {

	var rollupTypeIDRule []interface{}
	for _, rollupTypeIDItem := range rollupTypeID {
		rollupTypeIDRule = append(rollupTypeIDRule, rollupTypeIDItem)
	}

	logs, sub, err := _Mockfeijoapolygonrollupmanager.contract.WatchLogs(opts, "ObsoleteRollupType", rollupTypeIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockfeijoapolygonrollupmanagerObsoleteRollupType)
				if err := _Mockfeijoapolygonrollupmanager.contract.UnpackLog(event, "ObsoleteRollupType", log); err != nil {
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
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) ParseObsoleteRollupType(log types.Log) (*MockfeijoapolygonrollupmanagerObsoleteRollupType, error) {
	event := new(MockfeijoapolygonrollupmanagerObsoleteRollupType)
	if err := _Mockfeijoapolygonrollupmanager.contract.UnpackLog(event, "ObsoleteRollupType", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockfeijoapolygonrollupmanagerOnSequenceBatchesIterator is returned from FilterOnSequenceBatches and is used to iterate over the raw logs and unpacked data for OnSequenceBatches events raised by the Mockfeijoapolygonrollupmanager contract.
type MockfeijoapolygonrollupmanagerOnSequenceBatchesIterator struct {
	Event *MockfeijoapolygonrollupmanagerOnSequenceBatches // Event containing the contract specifics and raw log

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
func (it *MockfeijoapolygonrollupmanagerOnSequenceBatchesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockfeijoapolygonrollupmanagerOnSequenceBatches)
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
		it.Event = new(MockfeijoapolygonrollupmanagerOnSequenceBatches)
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
func (it *MockfeijoapolygonrollupmanagerOnSequenceBatchesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockfeijoapolygonrollupmanagerOnSequenceBatchesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockfeijoapolygonrollupmanagerOnSequenceBatches represents a OnSequenceBatches event raised by the Mockfeijoapolygonrollupmanager contract.
type MockfeijoapolygonrollupmanagerOnSequenceBatches struct {
	RollupID           uint32
	LastBatchSequenced uint64
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterOnSequenceBatches is a free log retrieval operation binding the contract event 0x1d9f30260051d51d70339da239ea7b080021adcaabfa71c9b0ea339a20cf9a25.
//
// Solidity: event OnSequenceBatches(uint32 indexed rollupID, uint64 lastBatchSequenced)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) FilterOnSequenceBatches(opts *bind.FilterOpts, rollupID []uint32) (*MockfeijoapolygonrollupmanagerOnSequenceBatchesIterator, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	logs, sub, err := _Mockfeijoapolygonrollupmanager.contract.FilterLogs(opts, "OnSequenceBatches", rollupIDRule)
	if err != nil {
		return nil, err
	}
	return &MockfeijoapolygonrollupmanagerOnSequenceBatchesIterator{contract: _Mockfeijoapolygonrollupmanager.contract, event: "OnSequenceBatches", logs: logs, sub: sub}, nil
}

// WatchOnSequenceBatches is a free log subscription operation binding the contract event 0x1d9f30260051d51d70339da239ea7b080021adcaabfa71c9b0ea339a20cf9a25.
//
// Solidity: event OnSequenceBatches(uint32 indexed rollupID, uint64 lastBatchSequenced)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) WatchOnSequenceBatches(opts *bind.WatchOpts, sink chan<- *MockfeijoapolygonrollupmanagerOnSequenceBatches, rollupID []uint32) (event.Subscription, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	logs, sub, err := _Mockfeijoapolygonrollupmanager.contract.WatchLogs(opts, "OnSequenceBatches", rollupIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockfeijoapolygonrollupmanagerOnSequenceBatches)
				if err := _Mockfeijoapolygonrollupmanager.contract.UnpackLog(event, "OnSequenceBatches", log); err != nil {
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

// ParseOnSequenceBatches is a log parse operation binding the contract event 0x1d9f30260051d51d70339da239ea7b080021adcaabfa71c9b0ea339a20cf9a25.
//
// Solidity: event OnSequenceBatches(uint32 indexed rollupID, uint64 lastBatchSequenced)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) ParseOnSequenceBatches(log types.Log) (*MockfeijoapolygonrollupmanagerOnSequenceBatches, error) {
	event := new(MockfeijoapolygonrollupmanagerOnSequenceBatches)
	if err := _Mockfeijoapolygonrollupmanager.contract.UnpackLog(event, "OnSequenceBatches", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockfeijoapolygonrollupmanagerOverridePendingStateIterator is returned from FilterOverridePendingState and is used to iterate over the raw logs and unpacked data for OverridePendingState events raised by the Mockfeijoapolygonrollupmanager contract.
type MockfeijoapolygonrollupmanagerOverridePendingStateIterator struct {
	Event *MockfeijoapolygonrollupmanagerOverridePendingState // Event containing the contract specifics and raw log

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
func (it *MockfeijoapolygonrollupmanagerOverridePendingStateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockfeijoapolygonrollupmanagerOverridePendingState)
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
		it.Event = new(MockfeijoapolygonrollupmanagerOverridePendingState)
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
func (it *MockfeijoapolygonrollupmanagerOverridePendingStateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockfeijoapolygonrollupmanagerOverridePendingStateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockfeijoapolygonrollupmanagerOverridePendingState represents a OverridePendingState event raised by the Mockfeijoapolygonrollupmanager contract.
type MockfeijoapolygonrollupmanagerOverridePendingState struct {
	RollupID   uint32
	NumBatch   uint64
	StateRoot  [32]byte
	ExitRoot   [32]byte
	Aggregator common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterOverridePendingState is a free log retrieval operation binding the contract event 0x3182bd6e6f74fc1fdc88b60f3a4f4c7f79db6ae6f5b88a1b3f5a1e28ec210d5e.
//
// Solidity: event OverridePendingState(uint32 indexed rollupID, uint64 numBatch, bytes32 stateRoot, bytes32 exitRoot, address aggregator)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) FilterOverridePendingState(opts *bind.FilterOpts, rollupID []uint32) (*MockfeijoapolygonrollupmanagerOverridePendingStateIterator, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	logs, sub, err := _Mockfeijoapolygonrollupmanager.contract.FilterLogs(opts, "OverridePendingState", rollupIDRule)
	if err != nil {
		return nil, err
	}
	return &MockfeijoapolygonrollupmanagerOverridePendingStateIterator{contract: _Mockfeijoapolygonrollupmanager.contract, event: "OverridePendingState", logs: logs, sub: sub}, nil
}

// WatchOverridePendingState is a free log subscription operation binding the contract event 0x3182bd6e6f74fc1fdc88b60f3a4f4c7f79db6ae6f5b88a1b3f5a1e28ec210d5e.
//
// Solidity: event OverridePendingState(uint32 indexed rollupID, uint64 numBatch, bytes32 stateRoot, bytes32 exitRoot, address aggregator)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) WatchOverridePendingState(opts *bind.WatchOpts, sink chan<- *MockfeijoapolygonrollupmanagerOverridePendingState, rollupID []uint32) (event.Subscription, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	logs, sub, err := _Mockfeijoapolygonrollupmanager.contract.WatchLogs(opts, "OverridePendingState", rollupIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockfeijoapolygonrollupmanagerOverridePendingState)
				if err := _Mockfeijoapolygonrollupmanager.contract.UnpackLog(event, "OverridePendingState", log); err != nil {
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
// Solidity: event OverridePendingState(uint32 indexed rollupID, uint64 numBatch, bytes32 stateRoot, bytes32 exitRoot, address aggregator)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) ParseOverridePendingState(log types.Log) (*MockfeijoapolygonrollupmanagerOverridePendingState, error) {
	event := new(MockfeijoapolygonrollupmanagerOverridePendingState)
	if err := _Mockfeijoapolygonrollupmanager.contract.UnpackLog(event, "OverridePendingState", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockfeijoapolygonrollupmanagerProveNonDeterministicPendingStateIterator is returned from FilterProveNonDeterministicPendingState and is used to iterate over the raw logs and unpacked data for ProveNonDeterministicPendingState events raised by the Mockfeijoapolygonrollupmanager contract.
type MockfeijoapolygonrollupmanagerProveNonDeterministicPendingStateIterator struct {
	Event *MockfeijoapolygonrollupmanagerProveNonDeterministicPendingState // Event containing the contract specifics and raw log

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
func (it *MockfeijoapolygonrollupmanagerProveNonDeterministicPendingStateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockfeijoapolygonrollupmanagerProveNonDeterministicPendingState)
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
		it.Event = new(MockfeijoapolygonrollupmanagerProveNonDeterministicPendingState)
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
func (it *MockfeijoapolygonrollupmanagerProveNonDeterministicPendingStateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockfeijoapolygonrollupmanagerProveNonDeterministicPendingStateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockfeijoapolygonrollupmanagerProveNonDeterministicPendingState represents a ProveNonDeterministicPendingState event raised by the Mockfeijoapolygonrollupmanager contract.
type MockfeijoapolygonrollupmanagerProveNonDeterministicPendingState struct {
	StoredStateRoot [32]byte
	ProvedStateRoot [32]byte
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterProveNonDeterministicPendingState is a free log retrieval operation binding the contract event 0x1f44c21118c4603cfb4e1b621dbcfa2b73efcececee2b99b620b2953d33a7010.
//
// Solidity: event ProveNonDeterministicPendingState(bytes32 storedStateRoot, bytes32 provedStateRoot)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) FilterProveNonDeterministicPendingState(opts *bind.FilterOpts) (*MockfeijoapolygonrollupmanagerProveNonDeterministicPendingStateIterator, error) {

	logs, sub, err := _Mockfeijoapolygonrollupmanager.contract.FilterLogs(opts, "ProveNonDeterministicPendingState")
	if err != nil {
		return nil, err
	}
	return &MockfeijoapolygonrollupmanagerProveNonDeterministicPendingStateIterator{contract: _Mockfeijoapolygonrollupmanager.contract, event: "ProveNonDeterministicPendingState", logs: logs, sub: sub}, nil
}

// WatchProveNonDeterministicPendingState is a free log subscription operation binding the contract event 0x1f44c21118c4603cfb4e1b621dbcfa2b73efcececee2b99b620b2953d33a7010.
//
// Solidity: event ProveNonDeterministicPendingState(bytes32 storedStateRoot, bytes32 provedStateRoot)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) WatchProveNonDeterministicPendingState(opts *bind.WatchOpts, sink chan<- *MockfeijoapolygonrollupmanagerProveNonDeterministicPendingState) (event.Subscription, error) {

	logs, sub, err := _Mockfeijoapolygonrollupmanager.contract.WatchLogs(opts, "ProveNonDeterministicPendingState")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockfeijoapolygonrollupmanagerProveNonDeterministicPendingState)
				if err := _Mockfeijoapolygonrollupmanager.contract.UnpackLog(event, "ProveNonDeterministicPendingState", log); err != nil {
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
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) ParseProveNonDeterministicPendingState(log types.Log) (*MockfeijoapolygonrollupmanagerProveNonDeterministicPendingState, error) {
	event := new(MockfeijoapolygonrollupmanagerProveNonDeterministicPendingState)
	if err := _Mockfeijoapolygonrollupmanager.contract.UnpackLog(event, "ProveNonDeterministicPendingState", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockfeijoapolygonrollupmanagerRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the Mockfeijoapolygonrollupmanager contract.
type MockfeijoapolygonrollupmanagerRoleAdminChangedIterator struct {
	Event *MockfeijoapolygonrollupmanagerRoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *MockfeijoapolygonrollupmanagerRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockfeijoapolygonrollupmanagerRoleAdminChanged)
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
		it.Event = new(MockfeijoapolygonrollupmanagerRoleAdminChanged)
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
func (it *MockfeijoapolygonrollupmanagerRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockfeijoapolygonrollupmanagerRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockfeijoapolygonrollupmanagerRoleAdminChanged represents a RoleAdminChanged event raised by the Mockfeijoapolygonrollupmanager contract.
type MockfeijoapolygonrollupmanagerRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*MockfeijoapolygonrollupmanagerRoleAdminChangedIterator, error) {

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

	logs, sub, err := _Mockfeijoapolygonrollupmanager.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &MockfeijoapolygonrollupmanagerRoleAdminChangedIterator{contract: _Mockfeijoapolygonrollupmanager.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *MockfeijoapolygonrollupmanagerRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _Mockfeijoapolygonrollupmanager.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockfeijoapolygonrollupmanagerRoleAdminChanged)
				if err := _Mockfeijoapolygonrollupmanager.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) ParseRoleAdminChanged(log types.Log) (*MockfeijoapolygonrollupmanagerRoleAdminChanged, error) {
	event := new(MockfeijoapolygonrollupmanagerRoleAdminChanged)
	if err := _Mockfeijoapolygonrollupmanager.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockfeijoapolygonrollupmanagerRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the Mockfeijoapolygonrollupmanager contract.
type MockfeijoapolygonrollupmanagerRoleGrantedIterator struct {
	Event *MockfeijoapolygonrollupmanagerRoleGranted // Event containing the contract specifics and raw log

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
func (it *MockfeijoapolygonrollupmanagerRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockfeijoapolygonrollupmanagerRoleGranted)
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
		it.Event = new(MockfeijoapolygonrollupmanagerRoleGranted)
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
func (it *MockfeijoapolygonrollupmanagerRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockfeijoapolygonrollupmanagerRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockfeijoapolygonrollupmanagerRoleGranted represents a RoleGranted event raised by the Mockfeijoapolygonrollupmanager contract.
type MockfeijoapolygonrollupmanagerRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*MockfeijoapolygonrollupmanagerRoleGrantedIterator, error) {

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

	logs, sub, err := _Mockfeijoapolygonrollupmanager.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &MockfeijoapolygonrollupmanagerRoleGrantedIterator{contract: _Mockfeijoapolygonrollupmanager.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *MockfeijoapolygonrollupmanagerRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _Mockfeijoapolygonrollupmanager.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockfeijoapolygonrollupmanagerRoleGranted)
				if err := _Mockfeijoapolygonrollupmanager.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) ParseRoleGranted(log types.Log) (*MockfeijoapolygonrollupmanagerRoleGranted, error) {
	event := new(MockfeijoapolygonrollupmanagerRoleGranted)
	if err := _Mockfeijoapolygonrollupmanager.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockfeijoapolygonrollupmanagerRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the Mockfeijoapolygonrollupmanager contract.
type MockfeijoapolygonrollupmanagerRoleRevokedIterator struct {
	Event *MockfeijoapolygonrollupmanagerRoleRevoked // Event containing the contract specifics and raw log

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
func (it *MockfeijoapolygonrollupmanagerRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockfeijoapolygonrollupmanagerRoleRevoked)
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
		it.Event = new(MockfeijoapolygonrollupmanagerRoleRevoked)
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
func (it *MockfeijoapolygonrollupmanagerRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockfeijoapolygonrollupmanagerRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockfeijoapolygonrollupmanagerRoleRevoked represents a RoleRevoked event raised by the Mockfeijoapolygonrollupmanager contract.
type MockfeijoapolygonrollupmanagerRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*MockfeijoapolygonrollupmanagerRoleRevokedIterator, error) {

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

	logs, sub, err := _Mockfeijoapolygonrollupmanager.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &MockfeijoapolygonrollupmanagerRoleRevokedIterator{contract: _Mockfeijoapolygonrollupmanager.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *MockfeijoapolygonrollupmanagerRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _Mockfeijoapolygonrollupmanager.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockfeijoapolygonrollupmanagerRoleRevoked)
				if err := _Mockfeijoapolygonrollupmanager.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) ParseRoleRevoked(log types.Log) (*MockfeijoapolygonrollupmanagerRoleRevoked, error) {
	event := new(MockfeijoapolygonrollupmanagerRoleRevoked)
	if err := _Mockfeijoapolygonrollupmanager.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockfeijoapolygonrollupmanagerSetBatchFeeIterator is returned from FilterSetBatchFee and is used to iterate over the raw logs and unpacked data for SetBatchFee events raised by the Mockfeijoapolygonrollupmanager contract.
type MockfeijoapolygonrollupmanagerSetBatchFeeIterator struct {
	Event *MockfeijoapolygonrollupmanagerSetBatchFee // Event containing the contract specifics and raw log

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
func (it *MockfeijoapolygonrollupmanagerSetBatchFeeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockfeijoapolygonrollupmanagerSetBatchFee)
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
		it.Event = new(MockfeijoapolygonrollupmanagerSetBatchFee)
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
func (it *MockfeijoapolygonrollupmanagerSetBatchFeeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockfeijoapolygonrollupmanagerSetBatchFeeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockfeijoapolygonrollupmanagerSetBatchFee represents a SetBatchFee event raised by the Mockfeijoapolygonrollupmanager contract.
type MockfeijoapolygonrollupmanagerSetBatchFee struct {
	NewBatchFee *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterSetBatchFee is a free log retrieval operation binding the contract event 0xfb383653f53ee079978d0c9aff7aeff04a10166ce244cca9c9f9d8d96bed45b2.
//
// Solidity: event SetBatchFee(uint256 newBatchFee)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) FilterSetBatchFee(opts *bind.FilterOpts) (*MockfeijoapolygonrollupmanagerSetBatchFeeIterator, error) {

	logs, sub, err := _Mockfeijoapolygonrollupmanager.contract.FilterLogs(opts, "SetBatchFee")
	if err != nil {
		return nil, err
	}
	return &MockfeijoapolygonrollupmanagerSetBatchFeeIterator{contract: _Mockfeijoapolygonrollupmanager.contract, event: "SetBatchFee", logs: logs, sub: sub}, nil
}

// WatchSetBatchFee is a free log subscription operation binding the contract event 0xfb383653f53ee079978d0c9aff7aeff04a10166ce244cca9c9f9d8d96bed45b2.
//
// Solidity: event SetBatchFee(uint256 newBatchFee)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) WatchSetBatchFee(opts *bind.WatchOpts, sink chan<- *MockfeijoapolygonrollupmanagerSetBatchFee) (event.Subscription, error) {

	logs, sub, err := _Mockfeijoapolygonrollupmanager.contract.WatchLogs(opts, "SetBatchFee")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockfeijoapolygonrollupmanagerSetBatchFee)
				if err := _Mockfeijoapolygonrollupmanager.contract.UnpackLog(event, "SetBatchFee", log); err != nil {
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

// ParseSetBatchFee is a log parse operation binding the contract event 0xfb383653f53ee079978d0c9aff7aeff04a10166ce244cca9c9f9d8d96bed45b2.
//
// Solidity: event SetBatchFee(uint256 newBatchFee)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) ParseSetBatchFee(log types.Log) (*MockfeijoapolygonrollupmanagerSetBatchFee, error) {
	event := new(MockfeijoapolygonrollupmanagerSetBatchFee)
	if err := _Mockfeijoapolygonrollupmanager.contract.UnpackLog(event, "SetBatchFee", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockfeijoapolygonrollupmanagerSetMultiplierBatchFeeIterator is returned from FilterSetMultiplierBatchFee and is used to iterate over the raw logs and unpacked data for SetMultiplierBatchFee events raised by the Mockfeijoapolygonrollupmanager contract.
type MockfeijoapolygonrollupmanagerSetMultiplierBatchFeeIterator struct {
	Event *MockfeijoapolygonrollupmanagerSetMultiplierBatchFee // Event containing the contract specifics and raw log

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
func (it *MockfeijoapolygonrollupmanagerSetMultiplierBatchFeeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockfeijoapolygonrollupmanagerSetMultiplierBatchFee)
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
		it.Event = new(MockfeijoapolygonrollupmanagerSetMultiplierBatchFee)
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
func (it *MockfeijoapolygonrollupmanagerSetMultiplierBatchFeeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockfeijoapolygonrollupmanagerSetMultiplierBatchFeeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockfeijoapolygonrollupmanagerSetMultiplierBatchFee represents a SetMultiplierBatchFee event raised by the Mockfeijoapolygonrollupmanager contract.
type MockfeijoapolygonrollupmanagerSetMultiplierBatchFee struct {
	NewMultiplierBatchFee uint16
	Raw                   types.Log // Blockchain specific contextual infos
}

// FilterSetMultiplierBatchFee is a free log retrieval operation binding the contract event 0x7019933d795eba185c180209e8ae8bffbaa25bcef293364687702c31f4d302c5.
//
// Solidity: event SetMultiplierBatchFee(uint16 newMultiplierBatchFee)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) FilterSetMultiplierBatchFee(opts *bind.FilterOpts) (*MockfeijoapolygonrollupmanagerSetMultiplierBatchFeeIterator, error) {

	logs, sub, err := _Mockfeijoapolygonrollupmanager.contract.FilterLogs(opts, "SetMultiplierBatchFee")
	if err != nil {
		return nil, err
	}
	return &MockfeijoapolygonrollupmanagerSetMultiplierBatchFeeIterator{contract: _Mockfeijoapolygonrollupmanager.contract, event: "SetMultiplierBatchFee", logs: logs, sub: sub}, nil
}

// WatchSetMultiplierBatchFee is a free log subscription operation binding the contract event 0x7019933d795eba185c180209e8ae8bffbaa25bcef293364687702c31f4d302c5.
//
// Solidity: event SetMultiplierBatchFee(uint16 newMultiplierBatchFee)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) WatchSetMultiplierBatchFee(opts *bind.WatchOpts, sink chan<- *MockfeijoapolygonrollupmanagerSetMultiplierBatchFee) (event.Subscription, error) {

	logs, sub, err := _Mockfeijoapolygonrollupmanager.contract.WatchLogs(opts, "SetMultiplierBatchFee")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockfeijoapolygonrollupmanagerSetMultiplierBatchFee)
				if err := _Mockfeijoapolygonrollupmanager.contract.UnpackLog(event, "SetMultiplierBatchFee", log); err != nil {
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

// ParseSetMultiplierBatchFee is a log parse operation binding the contract event 0x7019933d795eba185c180209e8ae8bffbaa25bcef293364687702c31f4d302c5.
//
// Solidity: event SetMultiplierBatchFee(uint16 newMultiplierBatchFee)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) ParseSetMultiplierBatchFee(log types.Log) (*MockfeijoapolygonrollupmanagerSetMultiplierBatchFee, error) {
	event := new(MockfeijoapolygonrollupmanagerSetMultiplierBatchFee)
	if err := _Mockfeijoapolygonrollupmanager.contract.UnpackLog(event, "SetMultiplierBatchFee", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockfeijoapolygonrollupmanagerSetPendingStateTimeoutIterator is returned from FilterSetPendingStateTimeout and is used to iterate over the raw logs and unpacked data for SetPendingStateTimeout events raised by the Mockfeijoapolygonrollupmanager contract.
type MockfeijoapolygonrollupmanagerSetPendingStateTimeoutIterator struct {
	Event *MockfeijoapolygonrollupmanagerSetPendingStateTimeout // Event containing the contract specifics and raw log

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
func (it *MockfeijoapolygonrollupmanagerSetPendingStateTimeoutIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockfeijoapolygonrollupmanagerSetPendingStateTimeout)
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
		it.Event = new(MockfeijoapolygonrollupmanagerSetPendingStateTimeout)
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
func (it *MockfeijoapolygonrollupmanagerSetPendingStateTimeoutIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockfeijoapolygonrollupmanagerSetPendingStateTimeoutIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockfeijoapolygonrollupmanagerSetPendingStateTimeout represents a SetPendingStateTimeout event raised by the Mockfeijoapolygonrollupmanager contract.
type MockfeijoapolygonrollupmanagerSetPendingStateTimeout struct {
	NewPendingStateTimeout uint64
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterSetPendingStateTimeout is a free log retrieval operation binding the contract event 0xc4121f4e22c69632ebb7cf1f462be0511dc034f999b52013eddfb24aab765c75.
//
// Solidity: event SetPendingStateTimeout(uint64 newPendingStateTimeout)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) FilterSetPendingStateTimeout(opts *bind.FilterOpts) (*MockfeijoapolygonrollupmanagerSetPendingStateTimeoutIterator, error) {

	logs, sub, err := _Mockfeijoapolygonrollupmanager.contract.FilterLogs(opts, "SetPendingStateTimeout")
	if err != nil {
		return nil, err
	}
	return &MockfeijoapolygonrollupmanagerSetPendingStateTimeoutIterator{contract: _Mockfeijoapolygonrollupmanager.contract, event: "SetPendingStateTimeout", logs: logs, sub: sub}, nil
}

// WatchSetPendingStateTimeout is a free log subscription operation binding the contract event 0xc4121f4e22c69632ebb7cf1f462be0511dc034f999b52013eddfb24aab765c75.
//
// Solidity: event SetPendingStateTimeout(uint64 newPendingStateTimeout)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) WatchSetPendingStateTimeout(opts *bind.WatchOpts, sink chan<- *MockfeijoapolygonrollupmanagerSetPendingStateTimeout) (event.Subscription, error) {

	logs, sub, err := _Mockfeijoapolygonrollupmanager.contract.WatchLogs(opts, "SetPendingStateTimeout")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockfeijoapolygonrollupmanagerSetPendingStateTimeout)
				if err := _Mockfeijoapolygonrollupmanager.contract.UnpackLog(event, "SetPendingStateTimeout", log); err != nil {
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
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) ParseSetPendingStateTimeout(log types.Log) (*MockfeijoapolygonrollupmanagerSetPendingStateTimeout, error) {
	event := new(MockfeijoapolygonrollupmanagerSetPendingStateTimeout)
	if err := _Mockfeijoapolygonrollupmanager.contract.UnpackLog(event, "SetPendingStateTimeout", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockfeijoapolygonrollupmanagerSetTrustedAggregatorIterator is returned from FilterSetTrustedAggregator and is used to iterate over the raw logs and unpacked data for SetTrustedAggregator events raised by the Mockfeijoapolygonrollupmanager contract.
type MockfeijoapolygonrollupmanagerSetTrustedAggregatorIterator struct {
	Event *MockfeijoapolygonrollupmanagerSetTrustedAggregator // Event containing the contract specifics and raw log

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
func (it *MockfeijoapolygonrollupmanagerSetTrustedAggregatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockfeijoapolygonrollupmanagerSetTrustedAggregator)
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
		it.Event = new(MockfeijoapolygonrollupmanagerSetTrustedAggregator)
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
func (it *MockfeijoapolygonrollupmanagerSetTrustedAggregatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockfeijoapolygonrollupmanagerSetTrustedAggregatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockfeijoapolygonrollupmanagerSetTrustedAggregator represents a SetTrustedAggregator event raised by the Mockfeijoapolygonrollupmanager contract.
type MockfeijoapolygonrollupmanagerSetTrustedAggregator struct {
	NewTrustedAggregator common.Address
	Raw                  types.Log // Blockchain specific contextual infos
}

// FilterSetTrustedAggregator is a free log retrieval operation binding the contract event 0x61f8fec29495a3078e9271456f05fb0707fd4e41f7661865f80fc437d06681ca.
//
// Solidity: event SetTrustedAggregator(address newTrustedAggregator)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) FilterSetTrustedAggregator(opts *bind.FilterOpts) (*MockfeijoapolygonrollupmanagerSetTrustedAggregatorIterator, error) {

	logs, sub, err := _Mockfeijoapolygonrollupmanager.contract.FilterLogs(opts, "SetTrustedAggregator")
	if err != nil {
		return nil, err
	}
	return &MockfeijoapolygonrollupmanagerSetTrustedAggregatorIterator{contract: _Mockfeijoapolygonrollupmanager.contract, event: "SetTrustedAggregator", logs: logs, sub: sub}, nil
}

// WatchSetTrustedAggregator is a free log subscription operation binding the contract event 0x61f8fec29495a3078e9271456f05fb0707fd4e41f7661865f80fc437d06681ca.
//
// Solidity: event SetTrustedAggregator(address newTrustedAggregator)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) WatchSetTrustedAggregator(opts *bind.WatchOpts, sink chan<- *MockfeijoapolygonrollupmanagerSetTrustedAggregator) (event.Subscription, error) {

	logs, sub, err := _Mockfeijoapolygonrollupmanager.contract.WatchLogs(opts, "SetTrustedAggregator")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockfeijoapolygonrollupmanagerSetTrustedAggregator)
				if err := _Mockfeijoapolygonrollupmanager.contract.UnpackLog(event, "SetTrustedAggregator", log); err != nil {
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
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) ParseSetTrustedAggregator(log types.Log) (*MockfeijoapolygonrollupmanagerSetTrustedAggregator, error) {
	event := new(MockfeijoapolygonrollupmanagerSetTrustedAggregator)
	if err := _Mockfeijoapolygonrollupmanager.contract.UnpackLog(event, "SetTrustedAggregator", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockfeijoapolygonrollupmanagerSetTrustedAggregatorTimeoutIterator is returned from FilterSetTrustedAggregatorTimeout and is used to iterate over the raw logs and unpacked data for SetTrustedAggregatorTimeout events raised by the Mockfeijoapolygonrollupmanager contract.
type MockfeijoapolygonrollupmanagerSetTrustedAggregatorTimeoutIterator struct {
	Event *MockfeijoapolygonrollupmanagerSetTrustedAggregatorTimeout // Event containing the contract specifics and raw log

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
func (it *MockfeijoapolygonrollupmanagerSetTrustedAggregatorTimeoutIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockfeijoapolygonrollupmanagerSetTrustedAggregatorTimeout)
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
		it.Event = new(MockfeijoapolygonrollupmanagerSetTrustedAggregatorTimeout)
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
func (it *MockfeijoapolygonrollupmanagerSetTrustedAggregatorTimeoutIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockfeijoapolygonrollupmanagerSetTrustedAggregatorTimeoutIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockfeijoapolygonrollupmanagerSetTrustedAggregatorTimeout represents a SetTrustedAggregatorTimeout event raised by the Mockfeijoapolygonrollupmanager contract.
type MockfeijoapolygonrollupmanagerSetTrustedAggregatorTimeout struct {
	NewTrustedAggregatorTimeout uint64
	Raw                         types.Log // Blockchain specific contextual infos
}

// FilterSetTrustedAggregatorTimeout is a free log retrieval operation binding the contract event 0x1f4fa24c2e4bad19a7f3ec5c5485f70d46c798461c2e684f55bbd0fc661373a1.
//
// Solidity: event SetTrustedAggregatorTimeout(uint64 newTrustedAggregatorTimeout)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) FilterSetTrustedAggregatorTimeout(opts *bind.FilterOpts) (*MockfeijoapolygonrollupmanagerSetTrustedAggregatorTimeoutIterator, error) {

	logs, sub, err := _Mockfeijoapolygonrollupmanager.contract.FilterLogs(opts, "SetTrustedAggregatorTimeout")
	if err != nil {
		return nil, err
	}
	return &MockfeijoapolygonrollupmanagerSetTrustedAggregatorTimeoutIterator{contract: _Mockfeijoapolygonrollupmanager.contract, event: "SetTrustedAggregatorTimeout", logs: logs, sub: sub}, nil
}

// WatchSetTrustedAggregatorTimeout is a free log subscription operation binding the contract event 0x1f4fa24c2e4bad19a7f3ec5c5485f70d46c798461c2e684f55bbd0fc661373a1.
//
// Solidity: event SetTrustedAggregatorTimeout(uint64 newTrustedAggregatorTimeout)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) WatchSetTrustedAggregatorTimeout(opts *bind.WatchOpts, sink chan<- *MockfeijoapolygonrollupmanagerSetTrustedAggregatorTimeout) (event.Subscription, error) {

	logs, sub, err := _Mockfeijoapolygonrollupmanager.contract.WatchLogs(opts, "SetTrustedAggregatorTimeout")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockfeijoapolygonrollupmanagerSetTrustedAggregatorTimeout)
				if err := _Mockfeijoapolygonrollupmanager.contract.UnpackLog(event, "SetTrustedAggregatorTimeout", log); err != nil {
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
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) ParseSetTrustedAggregatorTimeout(log types.Log) (*MockfeijoapolygonrollupmanagerSetTrustedAggregatorTimeout, error) {
	event := new(MockfeijoapolygonrollupmanagerSetTrustedAggregatorTimeout)
	if err := _Mockfeijoapolygonrollupmanager.contract.UnpackLog(event, "SetTrustedAggregatorTimeout", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockfeijoapolygonrollupmanagerSetVerifyBatchTimeTargetIterator is returned from FilterSetVerifyBatchTimeTarget and is used to iterate over the raw logs and unpacked data for SetVerifyBatchTimeTarget events raised by the Mockfeijoapolygonrollupmanager contract.
type MockfeijoapolygonrollupmanagerSetVerifyBatchTimeTargetIterator struct {
	Event *MockfeijoapolygonrollupmanagerSetVerifyBatchTimeTarget // Event containing the contract specifics and raw log

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
func (it *MockfeijoapolygonrollupmanagerSetVerifyBatchTimeTargetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockfeijoapolygonrollupmanagerSetVerifyBatchTimeTarget)
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
		it.Event = new(MockfeijoapolygonrollupmanagerSetVerifyBatchTimeTarget)
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
func (it *MockfeijoapolygonrollupmanagerSetVerifyBatchTimeTargetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockfeijoapolygonrollupmanagerSetVerifyBatchTimeTargetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockfeijoapolygonrollupmanagerSetVerifyBatchTimeTarget represents a SetVerifyBatchTimeTarget event raised by the Mockfeijoapolygonrollupmanager contract.
type MockfeijoapolygonrollupmanagerSetVerifyBatchTimeTarget struct {
	NewVerifyBatchTimeTarget uint64
	Raw                      types.Log // Blockchain specific contextual infos
}

// FilterSetVerifyBatchTimeTarget is a free log retrieval operation binding the contract event 0x1b023231a1ab6b5d93992f168fb44498e1a7e64cef58daff6f1c216de6a68c28.
//
// Solidity: event SetVerifyBatchTimeTarget(uint64 newVerifyBatchTimeTarget)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) FilterSetVerifyBatchTimeTarget(opts *bind.FilterOpts) (*MockfeijoapolygonrollupmanagerSetVerifyBatchTimeTargetIterator, error) {

	logs, sub, err := _Mockfeijoapolygonrollupmanager.contract.FilterLogs(opts, "SetVerifyBatchTimeTarget")
	if err != nil {
		return nil, err
	}
	return &MockfeijoapolygonrollupmanagerSetVerifyBatchTimeTargetIterator{contract: _Mockfeijoapolygonrollupmanager.contract, event: "SetVerifyBatchTimeTarget", logs: logs, sub: sub}, nil
}

// WatchSetVerifyBatchTimeTarget is a free log subscription operation binding the contract event 0x1b023231a1ab6b5d93992f168fb44498e1a7e64cef58daff6f1c216de6a68c28.
//
// Solidity: event SetVerifyBatchTimeTarget(uint64 newVerifyBatchTimeTarget)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) WatchSetVerifyBatchTimeTarget(opts *bind.WatchOpts, sink chan<- *MockfeijoapolygonrollupmanagerSetVerifyBatchTimeTarget) (event.Subscription, error) {

	logs, sub, err := _Mockfeijoapolygonrollupmanager.contract.WatchLogs(opts, "SetVerifyBatchTimeTarget")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockfeijoapolygonrollupmanagerSetVerifyBatchTimeTarget)
				if err := _Mockfeijoapolygonrollupmanager.contract.UnpackLog(event, "SetVerifyBatchTimeTarget", log); err != nil {
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

// ParseSetVerifyBatchTimeTarget is a log parse operation binding the contract event 0x1b023231a1ab6b5d93992f168fb44498e1a7e64cef58daff6f1c216de6a68c28.
//
// Solidity: event SetVerifyBatchTimeTarget(uint64 newVerifyBatchTimeTarget)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) ParseSetVerifyBatchTimeTarget(log types.Log) (*MockfeijoapolygonrollupmanagerSetVerifyBatchTimeTarget, error) {
	event := new(MockfeijoapolygonrollupmanagerSetVerifyBatchTimeTarget)
	if err := _Mockfeijoapolygonrollupmanager.contract.UnpackLog(event, "SetVerifyBatchTimeTarget", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockfeijoapolygonrollupmanagerUpdateRollupIterator is returned from FilterUpdateRollup and is used to iterate over the raw logs and unpacked data for UpdateRollup events raised by the Mockfeijoapolygonrollupmanager contract.
type MockfeijoapolygonrollupmanagerUpdateRollupIterator struct {
	Event *MockfeijoapolygonrollupmanagerUpdateRollup // Event containing the contract specifics and raw log

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
func (it *MockfeijoapolygonrollupmanagerUpdateRollupIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockfeijoapolygonrollupmanagerUpdateRollup)
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
		it.Event = new(MockfeijoapolygonrollupmanagerUpdateRollup)
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
func (it *MockfeijoapolygonrollupmanagerUpdateRollupIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockfeijoapolygonrollupmanagerUpdateRollupIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockfeijoapolygonrollupmanagerUpdateRollup represents a UpdateRollup event raised by the Mockfeijoapolygonrollupmanager contract.
type MockfeijoapolygonrollupmanagerUpdateRollup struct {
	RollupID                       uint32
	NewRollupTypeID                uint32
	LastVerifiedBatchBeforeUpgrade uint64
	Raw                            types.Log // Blockchain specific contextual infos
}

// FilterUpdateRollup is a free log retrieval operation binding the contract event 0xf585e04c05d396901170247783d3e5f0ee9c1df23072985b50af089f5e48b19d.
//
// Solidity: event UpdateRollup(uint32 indexed rollupID, uint32 newRollupTypeID, uint64 lastVerifiedBatchBeforeUpgrade)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) FilterUpdateRollup(opts *bind.FilterOpts, rollupID []uint32) (*MockfeijoapolygonrollupmanagerUpdateRollupIterator, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	logs, sub, err := _Mockfeijoapolygonrollupmanager.contract.FilterLogs(opts, "UpdateRollup", rollupIDRule)
	if err != nil {
		return nil, err
	}
	return &MockfeijoapolygonrollupmanagerUpdateRollupIterator{contract: _Mockfeijoapolygonrollupmanager.contract, event: "UpdateRollup", logs: logs, sub: sub}, nil
}

// WatchUpdateRollup is a free log subscription operation binding the contract event 0xf585e04c05d396901170247783d3e5f0ee9c1df23072985b50af089f5e48b19d.
//
// Solidity: event UpdateRollup(uint32 indexed rollupID, uint32 newRollupTypeID, uint64 lastVerifiedBatchBeforeUpgrade)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) WatchUpdateRollup(opts *bind.WatchOpts, sink chan<- *MockfeijoapolygonrollupmanagerUpdateRollup, rollupID []uint32) (event.Subscription, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	logs, sub, err := _Mockfeijoapolygonrollupmanager.contract.WatchLogs(opts, "UpdateRollup", rollupIDRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockfeijoapolygonrollupmanagerUpdateRollup)
				if err := _Mockfeijoapolygonrollupmanager.contract.UnpackLog(event, "UpdateRollup", log); err != nil {
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
// Solidity: event UpdateRollup(uint32 indexed rollupID, uint32 newRollupTypeID, uint64 lastVerifiedBatchBeforeUpgrade)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) ParseUpdateRollup(log types.Log) (*MockfeijoapolygonrollupmanagerUpdateRollup, error) {
	event := new(MockfeijoapolygonrollupmanagerUpdateRollup)
	if err := _Mockfeijoapolygonrollupmanager.contract.UnpackLog(event, "UpdateRollup", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockfeijoapolygonrollupmanagerVerifyBatchesIterator is returned from FilterVerifyBatches and is used to iterate over the raw logs and unpacked data for VerifyBatches events raised by the Mockfeijoapolygonrollupmanager contract.
type MockfeijoapolygonrollupmanagerVerifyBatchesIterator struct {
	Event *MockfeijoapolygonrollupmanagerVerifyBatches // Event containing the contract specifics and raw log

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
func (it *MockfeijoapolygonrollupmanagerVerifyBatchesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockfeijoapolygonrollupmanagerVerifyBatches)
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
		it.Event = new(MockfeijoapolygonrollupmanagerVerifyBatches)
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
func (it *MockfeijoapolygonrollupmanagerVerifyBatchesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockfeijoapolygonrollupmanagerVerifyBatchesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockfeijoapolygonrollupmanagerVerifyBatches represents a VerifyBatches event raised by the Mockfeijoapolygonrollupmanager contract.
type MockfeijoapolygonrollupmanagerVerifyBatches struct {
	RollupID   uint32
	NumBatch   uint64
	StateRoot  [32]byte
	ExitRoot   [32]byte
	Aggregator common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterVerifyBatches is a free log retrieval operation binding the contract event 0xaac1e7a157b259544ebacd6e8a82ae5d6c8f174e12aa48696277bcc9a661f0b4.
//
// Solidity: event VerifyBatches(uint32 indexed rollupID, uint64 numBatch, bytes32 stateRoot, bytes32 exitRoot, address indexed aggregator)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) FilterVerifyBatches(opts *bind.FilterOpts, rollupID []uint32, aggregator []common.Address) (*MockfeijoapolygonrollupmanagerVerifyBatchesIterator, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	var aggregatorRule []interface{}
	for _, aggregatorItem := range aggregator {
		aggregatorRule = append(aggregatorRule, aggregatorItem)
	}

	logs, sub, err := _Mockfeijoapolygonrollupmanager.contract.FilterLogs(opts, "VerifyBatches", rollupIDRule, aggregatorRule)
	if err != nil {
		return nil, err
	}
	return &MockfeijoapolygonrollupmanagerVerifyBatchesIterator{contract: _Mockfeijoapolygonrollupmanager.contract, event: "VerifyBatches", logs: logs, sub: sub}, nil
}

// WatchVerifyBatches is a free log subscription operation binding the contract event 0xaac1e7a157b259544ebacd6e8a82ae5d6c8f174e12aa48696277bcc9a661f0b4.
//
// Solidity: event VerifyBatches(uint32 indexed rollupID, uint64 numBatch, bytes32 stateRoot, bytes32 exitRoot, address indexed aggregator)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) WatchVerifyBatches(opts *bind.WatchOpts, sink chan<- *MockfeijoapolygonrollupmanagerVerifyBatches, rollupID []uint32, aggregator []common.Address) (event.Subscription, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	var aggregatorRule []interface{}
	for _, aggregatorItem := range aggregator {
		aggregatorRule = append(aggregatorRule, aggregatorItem)
	}

	logs, sub, err := _Mockfeijoapolygonrollupmanager.contract.WatchLogs(opts, "VerifyBatches", rollupIDRule, aggregatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockfeijoapolygonrollupmanagerVerifyBatches)
				if err := _Mockfeijoapolygonrollupmanager.contract.UnpackLog(event, "VerifyBatches", log); err != nil {
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

// ParseVerifyBatches is a log parse operation binding the contract event 0xaac1e7a157b259544ebacd6e8a82ae5d6c8f174e12aa48696277bcc9a661f0b4.
//
// Solidity: event VerifyBatches(uint32 indexed rollupID, uint64 numBatch, bytes32 stateRoot, bytes32 exitRoot, address indexed aggregator)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) ParseVerifyBatches(log types.Log) (*MockfeijoapolygonrollupmanagerVerifyBatches, error) {
	event := new(MockfeijoapolygonrollupmanagerVerifyBatches)
	if err := _Mockfeijoapolygonrollupmanager.contract.UnpackLog(event, "VerifyBatches", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// MockfeijoapolygonrollupmanagerVerifyBatchesTrustedAggregatorIterator is returned from FilterVerifyBatchesTrustedAggregator and is used to iterate over the raw logs and unpacked data for VerifyBatchesTrustedAggregator events raised by the Mockfeijoapolygonrollupmanager contract.
type MockfeijoapolygonrollupmanagerVerifyBatchesTrustedAggregatorIterator struct {
	Event *MockfeijoapolygonrollupmanagerVerifyBatchesTrustedAggregator // Event containing the contract specifics and raw log

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
func (it *MockfeijoapolygonrollupmanagerVerifyBatchesTrustedAggregatorIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(MockfeijoapolygonrollupmanagerVerifyBatchesTrustedAggregator)
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
		it.Event = new(MockfeijoapolygonrollupmanagerVerifyBatchesTrustedAggregator)
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
func (it *MockfeijoapolygonrollupmanagerVerifyBatchesTrustedAggregatorIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *MockfeijoapolygonrollupmanagerVerifyBatchesTrustedAggregatorIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// MockfeijoapolygonrollupmanagerVerifyBatchesTrustedAggregator represents a VerifyBatchesTrustedAggregator event raised by the Mockfeijoapolygonrollupmanager contract.
type MockfeijoapolygonrollupmanagerVerifyBatchesTrustedAggregator struct {
	RollupID   uint32
	NumBatch   uint64
	StateRoot  [32]byte
	ExitRoot   [32]byte
	Aggregator common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterVerifyBatchesTrustedAggregator is a free log retrieval operation binding the contract event 0xd1ec3a1216f08b6eff72e169ceb548b782db18a6614852618d86bb19f3f9b0d3.
//
// Solidity: event VerifyBatchesTrustedAggregator(uint32 indexed rollupID, uint64 numBatch, bytes32 stateRoot, bytes32 exitRoot, address indexed aggregator)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) FilterVerifyBatchesTrustedAggregator(opts *bind.FilterOpts, rollupID []uint32, aggregator []common.Address) (*MockfeijoapolygonrollupmanagerVerifyBatchesTrustedAggregatorIterator, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	var aggregatorRule []interface{}
	for _, aggregatorItem := range aggregator {
		aggregatorRule = append(aggregatorRule, aggregatorItem)
	}

	logs, sub, err := _Mockfeijoapolygonrollupmanager.contract.FilterLogs(opts, "VerifyBatchesTrustedAggregator", rollupIDRule, aggregatorRule)
	if err != nil {
		return nil, err
	}
	return &MockfeijoapolygonrollupmanagerVerifyBatchesTrustedAggregatorIterator{contract: _Mockfeijoapolygonrollupmanager.contract, event: "VerifyBatchesTrustedAggregator", logs: logs, sub: sub}, nil
}

// WatchVerifyBatchesTrustedAggregator is a free log subscription operation binding the contract event 0xd1ec3a1216f08b6eff72e169ceb548b782db18a6614852618d86bb19f3f9b0d3.
//
// Solidity: event VerifyBatchesTrustedAggregator(uint32 indexed rollupID, uint64 numBatch, bytes32 stateRoot, bytes32 exitRoot, address indexed aggregator)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) WatchVerifyBatchesTrustedAggregator(opts *bind.WatchOpts, sink chan<- *MockfeijoapolygonrollupmanagerVerifyBatchesTrustedAggregator, rollupID []uint32, aggregator []common.Address) (event.Subscription, error) {

	var rollupIDRule []interface{}
	for _, rollupIDItem := range rollupID {
		rollupIDRule = append(rollupIDRule, rollupIDItem)
	}

	var aggregatorRule []interface{}
	for _, aggregatorItem := range aggregator {
		aggregatorRule = append(aggregatorRule, aggregatorItem)
	}

	logs, sub, err := _Mockfeijoapolygonrollupmanager.contract.WatchLogs(opts, "VerifyBatchesTrustedAggregator", rollupIDRule, aggregatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(MockfeijoapolygonrollupmanagerVerifyBatchesTrustedAggregator)
				if err := _Mockfeijoapolygonrollupmanager.contract.UnpackLog(event, "VerifyBatchesTrustedAggregator", log); err != nil {
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

// ParseVerifyBatchesTrustedAggregator is a log parse operation binding the contract event 0xd1ec3a1216f08b6eff72e169ceb548b782db18a6614852618d86bb19f3f9b0d3.
//
// Solidity: event VerifyBatchesTrustedAggregator(uint32 indexed rollupID, uint64 numBatch, bytes32 stateRoot, bytes32 exitRoot, address indexed aggregator)
func (_Mockfeijoapolygonrollupmanager *MockfeijoapolygonrollupmanagerFilterer) ParseVerifyBatchesTrustedAggregator(log types.Log) (*MockfeijoapolygonrollupmanagerVerifyBatchesTrustedAggregator, error) {
	event := new(MockfeijoapolygonrollupmanagerVerifyBatchesTrustedAggregator)
	if err := _Mockfeijoapolygonrollupmanager.contract.UnpackLog(event, "VerifyBatchesTrustedAggregator", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
