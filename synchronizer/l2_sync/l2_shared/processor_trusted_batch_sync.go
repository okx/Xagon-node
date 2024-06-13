package l2_shared

import (
	"context"
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	syncCommon "github.com/0xPolygonHermez/zkevm-node/synchronizer/common"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/l2_sync"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

// BatchProcessMode is the mode for process a batch (full, incremental, reprocess, nothing)
type BatchProcessMode string

const (
	// FullProcessMode This batch is not on database, so is the first time we process it
	FullProcessMode BatchProcessMode = "full"
	// IncrementalProcessMode We have processed this batch before, and we have the intermediate state root, so is going to be process only new Tx
	IncrementalProcessMode BatchProcessMode = "incremental"
	// ReprocessProcessMode We have processed this batch before, but we don't have the intermediate state root, so we need to reprocess it
	ReprocessProcessMode BatchProcessMode = "reprocess"
	// NothingProcessMode The batch is already synchronized, so we don't need to process it
	NothingProcessMode BatchProcessMode = "nothing"
)

var (
	// ErrFatalBatchDesynchronized is the error when the batch is desynchronized
	ErrFatalBatchDesynchronized = fmt.Errorf("batch desynchronized")
)

// ProcessData contains the data required to process a batch
type ProcessData struct {
	BatchNumber       uint64
	Mode              BatchProcessMode
	OldStateRoot      common.Hash
	OldAccInputHash   common.Hash
	BatchMustBeClosed bool
	// TrustedBatch The batch in trusted node, it NEVER will be nil
	TrustedBatch *types.Batch
	// StateBatch Current batch in state DB, it could be nil
	StateBatch *state.Batch
	// PreviousStateBatch Previous batch in state DB (BatchNumber - 1), it could be nil
	PreviousStateBatch *state.Batch
	Now                time.Time
	Description        string
	// DebugPrefix is used to log, must prefix all logs entries
	DebugPrefix string
}

// ProcessResponse contains the response of the process of a batch
type ProcessResponse struct {
	// ProcessBatchResponse have the NewStateRoot
	ProcessBatchResponse *state.ProcessBatchResponse
	// ClearCache force to clear cache for next execution
	ClearCache bool
	// UpdateBatch  update the batch for next execution
	UpdateBatch *state.Batch
	// UpdateBatchWithProcessBatchResponse update the batch (if not nil) with the data in ProcessBatchResponse
	UpdateBatchWithProcessBatchResponse bool
}

// NewProcessResponse creates a new ProcessResponse
func NewProcessResponse() ProcessResponse {
	return ProcessResponse{
		ProcessBatchResponse:                nil,
		ClearCache:                          false,
		UpdateBatch:                         nil,
		UpdateBatchWithProcessBatchResponse: false,
	}
}

// DiscardCache set to discard cache for next execution
func (p *ProcessResponse) DiscardCache() {
	p.ClearCache = true
}

// UpdateCurrentBatch update the current batch for next execution
func (p *ProcessResponse) UpdateCurrentBatch(UpdateBatch *state.Batch) {
	p.ClearCache = false
	p.UpdateBatch = UpdateBatch
	p.UpdateBatchWithProcessBatchResponse = false
}

// UpdateCurrentBatchWithExecutionResult update the current batch for next execution with the data in ProcessBatchResponse
func (p *ProcessResponse) UpdateCurrentBatchWithExecutionResult(UpdateBatch *state.Batch, ProcessBatchResponse *state.ProcessBatchResponse) {
	p.ClearCache = false
	p.UpdateBatch = UpdateBatch
	p.UpdateBatchWithProcessBatchResponse = true
	p.ProcessBatchResponse = ProcessBatchResponse
}

// CheckSanity check the sanity of the response
func (p *ProcessResponse) CheckSanity() error {
	if p.UpdateBatchWithProcessBatchResponse {
		if p.ProcessBatchResponse == nil {
			return fmt.Errorf("UpdateBatchWithProcessBatchResponse is true but ProcessBatchResponse is nil")
		}
		if p.UpdateBatch == nil {
			return fmt.Errorf("UpdateBatchWithProcessBatchResponse is true but UpdateBatch is nil")
		}
		if p.ClearCache {
			return fmt.Errorf("UpdateBatchWithProcessBatchResponse is true but ClearCache is true")
		}
	}
	if p.UpdateBatch != nil {
		if p.ClearCache {
			return fmt.Errorf("UpdateBatch is not nil but ClearCache is true")
		}
	}
	return nil
}

// SyncTrustedBatchExecutor is the interface that known how to process a batch
type SyncTrustedBatchExecutor interface {
	// FullProcess process a batch that is not on database, so is the first time we process it
	FullProcess(ctx context.Context, data *ProcessData, dbTx pgx.Tx) (*ProcessResponse, error)
	// IncrementalProcess process a batch that we have processed before, and we have the intermediate state root, so is going to be process only new Tx
	IncrementalProcess(ctx context.Context, data *ProcessData, dbTx pgx.Tx) (*ProcessResponse, error)
	// ReProcess process a batch that we have processed before, but we don't have the intermediate state root, so we need to reprocess it
	ReProcess(ctx context.Context, data *ProcessData, dbTx pgx.Tx) (*ProcessResponse, error)
	// NothingProcess process a batch that is already synchronized, so we don't need to process it
	NothingProcess(ctx context.Context, data *ProcessData, dbTx pgx.Tx) (*ProcessResponse, error)
}

// L1SyncGlobalExitRootChecker is the interface to check if the required GlobalExitRoot is already synced from L1
type L1SyncGlobalExitRootChecker interface {
	CheckL1SyncGlobalExitRootEnoughToProcessBatch(ctx context.Context, batchNumber uint64, globalExitRoot common.Hash, dbTx pgx.Tx) error
}

// PostClosedBatchChecker is the interface to implement a checker post closed batch
type PostClosedBatchChecker interface {
	CheckPostClosedBatch(ctx context.Context, processData ProcessData, dbTx pgx.Tx) error
}

// ProcessorTrustedBatchSync is a template to sync trusted state. It classify what kind of update is needed and call to SyncTrustedStateBatchExecutorSteps
//
//	  that is the one that execute the sync process
//
//		the real implementation of the steps is in the SyncTrustedStateBatchExecutorSteps interface that known how to process a batch
type ProcessorTrustedBatchSync struct {
	Steps              SyncTrustedBatchExecutor
	timeProvider       syncCommon.TimeProvider
	l1SyncChecker      L1SyncGlobalExitRootChecker
	postClosedCheckers []PostClosedBatchChecker
	Cfg                l2_sync.Config
}

// NewProcessorTrustedBatchSync creates a new SyncTrustedStateBatchExecutorTemplate
func NewProcessorTrustedBatchSync(steps SyncTrustedBatchExecutor,
	timeProvider syncCommon.TimeProvider, l1SyncChecker L1SyncGlobalExitRootChecker, cfg l2_sync.Config) *ProcessorTrustedBatchSync {
	return &ProcessorTrustedBatchSync{
		Steps:         steps,
		timeProvider:  timeProvider,
		l1SyncChecker: l1SyncChecker,
		Cfg:           cfg,
	}
}

// AddPostChecker add a post closed batch checker
func (s *ProcessorTrustedBatchSync) AddPostChecker(checker PostClosedBatchChecker) {
	if s.postClosedCheckers == nil {
		s.postClosedCheckers = make([]PostClosedBatchChecker, 0)
	}
	s.postClosedCheckers = append(s.postClosedCheckers, checker)
}

// ProcessTrustedBatch processes a trusted batch and return the new state
func (s *ProcessorTrustedBatchSync) ProcessTrustedBatch(ctx context.Context, trustedBatch *types.Batch, status TrustedState, dbTx pgx.Tx, debugPrefix string) (*TrustedState, error) {
	if trustedBatch == nil {
		log.Errorf("%s trustedBatch is nil, it never should be nil", debugPrefix)
		return nil, fmt.Errorf("%s trustedBatch is nil, it never should be nil", debugPrefix)
	}
	log.Debugf("%s Processing trusted batch: %v", debugPrefix, trustedBatch.Number)
	stateCurrentBatch, statePreviousBatch := s.GetCurrentAndPreviousBatchFromCache(&status)
	if s.l1SyncChecker != nil {
		err := s.l1SyncChecker.CheckL1SyncGlobalExitRootEnoughToProcessBatch(ctx, uint64(trustedBatch.Number), trustedBatch.GlobalExitRoot, dbTx)
		if err != nil {
			log.Errorf("%s error checking GlobalExitRoot from TrustedBatch. Error: ", debugPrefix, err)
			return nil, err
		}
	} else {
		log.Infof("Disabled check L1 sync status for process batch")
	}
	processMode, err := s.GetModeForProcessBatch(trustedBatch, stateCurrentBatch, statePreviousBatch, debugPrefix)
	if err != nil {
		log.Error("%s error getting processMode. Error: ", debugPrefix, trustedBatch.Number, err)
		return nil, err
	}
	processBatchResp, err := s.ExecuteProcessBatch(ctx, &processMode, dbTx)
	if err != nil {
		log.Errorf("%s error processing trusted batch. Error: %s", processMode.DebugPrefix, err)
		return nil, err
	}
	return s.GetNextStatus(status, processBatchResp, processMode.BatchMustBeClosed, processMode.DebugPrefix)
}

// GetCurrentAndPreviousBatchFromCache returns the current and previous batch from cache
func (s *ProcessorTrustedBatchSync) GetCurrentAndPreviousBatchFromCache(status *TrustedState) (*state.Batch, *state.Batch) {
	if status == nil {
		return nil, nil
	}
	// Duplicate batches to avoid interferences with cache
	var stateCurrentBatch *state.Batch = nil
	var statePreviousBatch *state.Batch = nil
	if len(status.LastTrustedBatches) > 0 && status.LastTrustedBatches[0] != nil {
		tmpBatch := *status.LastTrustedBatches[0]
		stateCurrentBatch = &tmpBatch
	}
	if len(status.LastTrustedBatches) > 1 && status.LastTrustedBatches[1] != nil {
		tmpBatch := *status.LastTrustedBatches[1]
		statePreviousBatch = &tmpBatch
	}
	return stateCurrentBatch, statePreviousBatch
}

// GetNextStatus returns the next cache for use in the next run
// it could be nil, that means discard current cache
func (s *ProcessorTrustedBatchSync) GetNextStatus(status TrustedState, processBatchResp *ProcessResponse, closedBatch bool, debugPrefix string) (*TrustedState, error) {
	if processBatchResp != nil {
		err := processBatchResp.CheckSanity()
		if err != nil {
			// We dont stop the process but we log the warning to be fixed
			log.Warnf("%s error checking sanity of processBatchResp. Error: ", debugPrefix, err)
		}
	}

	newStatus := updateStatus(status, processBatchResp, closedBatch)
	log.Debugf("%s Batch synchronized, updated cache for next run", debugPrefix)
	return &newStatus, nil
}

// ExecuteProcessBatch execute the batch and process it
func (s *ProcessorTrustedBatchSync) ExecuteProcessBatch(ctx context.Context, processMode *ProcessData, dbTx pgx.Tx) (*ProcessResponse, error) {
	log.Infof("%s  Processing trusted batch: mode=%s desc=%s", processMode.DebugPrefix, processMode.Mode, processMode.Description)
	var processBatchResp *ProcessResponse = nil
	var err error
	switch processMode.Mode {
	case NothingProcessMode:
		log.Debugf("%s  no new L2BatchData", processMode.DebugPrefix, processMode.BatchNumber)
		processBatchResp, err = s.Steps.NothingProcess(ctx, processMode, dbTx)
	case FullProcessMode:
		log.Debugf("%s is not on database, so is the first time we process it", processMode.DebugPrefix)
		processBatchResp, err = s.Steps.FullProcess(ctx, processMode, dbTx)
	case IncrementalProcessMode:
		log.Debugf("%s is partially synchronized", processMode.DebugPrefix)
		processBatchResp, err = s.Steps.IncrementalProcess(ctx, processMode, dbTx)
	case ReprocessProcessMode:
		log.Debugf("%s is partially synchronized but we don't have intermediate stateRoot so it needs to be fully reprocessed", processMode.DebugPrefix)
		processBatchResp, err = s.Steps.ReProcess(ctx, processMode, dbTx)
	}
	if processBatchResp != nil && err == nil && processMode.BatchMustBeClosed {
		err = checkProcessBatchResultMatchExpected(processMode, processBatchResp.ProcessBatchResponse)
		if err != nil {
			log.Error("%s error verifying batch result!  Error: ", processMode.DebugPrefix, err)
			return nil, err
		}
		if s.postClosedCheckers != nil && len(s.postClosedCheckers) > 0 {
			for _, checker := range s.postClosedCheckers {
				err := checker.CheckPostClosedBatch(ctx, *processMode, dbTx)
				if err != nil {
					log.Errorf("%s error checking post closed batch. Error: ", processMode.DebugPrefix, err)
					return nil, err
				}
			}
		}
	}
	return processBatchResp, err
}

func updateStatus(status TrustedState, response *ProcessResponse, closedBatch bool) TrustedState {
	res := TrustedState{
		LastTrustedBatches: []*state.Batch{nil, nil},
	}
	if response == nil || response.ClearCache {
		return res
	}

	res.LastTrustedBatches[0] = status.GetCurrentBatch()
	res.LastTrustedBatches[1] = status.GetPreviousBatch()

	if response.UpdateBatch != nil {
		res.LastTrustedBatches[0] = response.UpdateBatch
	}
	if response.ProcessBatchResponse != nil && response.UpdateBatchWithProcessBatchResponse && res.LastTrustedBatches[0] != nil {
		// We copy the batch to avoid to modify the original object
		tmp := *response.UpdateBatch
		res.LastTrustedBatches[0] = &tmp
		res.LastTrustedBatches[0].StateRoot = response.ProcessBatchResponse.NewStateRoot
		res.LastTrustedBatches[0].LocalExitRoot = response.ProcessBatchResponse.NewLocalExitRoot
		res.LastTrustedBatches[0].AccInputHash = response.ProcessBatchResponse.NewAccInputHash
		res.LastTrustedBatches[0].WIP = !closedBatch
	}
	if closedBatch {
		res.LastTrustedBatches[1] = res.LastTrustedBatches[0]
		res.LastTrustedBatches[0] = nil
	}
	return res
}

// GetModeForProcessBatch returns the mode for process a batch
func (s *ProcessorTrustedBatchSync) GetModeForProcessBatch(trustedNodeBatch *types.Batch, stateBatch *state.Batch, statePreviousBatch *state.Batch, debugPrefix string) (ProcessData, error) {
	// Check parameters
	if trustedNodeBatch == nil || statePreviousBatch == nil {
		return ProcessData{}, fmt.Errorf("trustedNodeBatch and statePreviousBatch can't be nil")
	}

	var result ProcessData = ProcessData{}
	if stateBatch == nil {
		result = ProcessData{
			Mode:              FullProcessMode,
			OldStateRoot:      statePreviousBatch.StateRoot,
			BatchMustBeClosed: isTrustedBatchClosed(trustedNodeBatch),
			Description:       "Batch is not on database, so is the first time we process it",
		}
	} else {
		areBatchesExactlyEqual, strDiffsBatches := AreEqualStateBatchAndTrustedBatch(stateBatch, trustedNodeBatch, CMP_BATCH_IGNORE_TSTAMP)
		newL2DataFlag, err := ThereAreNewBatchL2Data(stateBatch.BatchL2Data, trustedNodeBatch.BatchL2Data)
		if err != nil {
			return ProcessData{}, err
		}
		if !newL2DataFlag {
			// "The batch from Node, and the one in database are the same, already synchronized",
			result = ProcessData{
				Mode:              NothingProcessMode,
				OldStateRoot:      common.Hash{},
				BatchMustBeClosed: isTrustedBatchClosed(trustedNodeBatch) && stateBatch.WIP,
				Description:       "no new data on batch. Diffs: " + strDiffsBatches,
			}
			if areBatchesExactlyEqual {
				result.BatchMustBeClosed = false
				result.Description = "exactly batches: " + strDiffsBatches
			}
		} else {
			// We have a previous batch, but in node something change
			// We have processed this batch before, and we have the intermediate state root, so is going to be process only new Tx.
			if stateBatch.StateRoot != state.ZeroHash {
				result = ProcessData{
					Mode:              IncrementalProcessMode,
					OldStateRoot:      stateBatch.StateRoot,
					BatchMustBeClosed: isTrustedBatchClosed(trustedNodeBatch),
					Description:       "batch exists + intermediateStateRoot " + strDiffsBatches,
				}
			} else {
				// We have processed this batch before, but we don't have the intermediate state root, so we need to reprocess all txs.
				result = ProcessData{
					Mode:              ReprocessProcessMode,
					OldStateRoot:      statePreviousBatch.StateRoot,
					BatchMustBeClosed: isTrustedBatchClosed(trustedNodeBatch),
					Description:       "batch exists + StateRoot==Zero" + strDiffsBatches,
				}
			}
		}
	}

	if s.Cfg.ReprocessFullBatchOnClose && result.BatchMustBeClosed {
		if result.Mode == IncrementalProcessMode || result.Mode == NothingProcessMode {
			result.Description = "forced reprocess due to batch closed and ReprocessFullBatchOnClose"
			log.Infof("%s Batch %v: Converted mode %s to %s because cfg.ReprocessFullBatchOnClose", debugPrefix, trustedNodeBatch.Number, result.Mode, ReprocessProcessMode)
			result.Mode = ReprocessProcessMode
			result.OldStateRoot = statePreviousBatch.StateRoot
			result.BatchMustBeClosed = true
		}
	}

	if result.Mode == "" {
		return result, fmt.Errorf("batch %v: failed to get mode for process ", trustedNodeBatch.Number)
	}

	result.BatchNumber = uint64(trustedNodeBatch.Number)
	result.StateBatch = stateBatch
	result.TrustedBatch = trustedNodeBatch
	result.PreviousStateBatch = statePreviousBatch
	result.OldAccInputHash = statePreviousBatch.AccInputHash
	result.Now = s.timeProvider.Now()
	result.DebugPrefix = fmt.Sprintf("%s mode %s:", debugPrefix, result.Mode)
	if result.BatchMustBeClosed {
		result.DebugPrefix += " (must_be_closed)"
	}
	if isTrustedBatchEmptyAndClosed(trustedNodeBatch) {
		if s.Cfg.AcceptEmptyClosedBatches {
			log.Infof("%s Batch %v: TrustedBatch Empty and closed, accepted due configuration", result.DebugPrefix, trustedNodeBatch.Number)
		} else {
			err := fmt.Errorf("%s Batch %v: TrustedBatch Empty and closed, rejected due configuration", result.DebugPrefix, trustedNodeBatch.Number)
			log.Infof(err.Error())
			return result, err
		}
	}

	return result, nil
}

func isTrustedBatchClosed(batch *types.Batch) bool {
	return batch.Closed
}

func isTrustedBatchEmptyAndClosed(batch *types.Batch) bool {
	return len(batch.BatchL2Data) == 0 && isTrustedBatchClosed(batch)
}

func checkStateRootAndLER(batchNumber uint64, expectedStateRoot common.Hash, expectedLER common.Hash, calculatedStateRoot common.Hash, calculatedLER common.Hash) error {
	if calculatedStateRoot != expectedStateRoot {
		return fmt.Errorf("batch %v: stareRoot calculated [%s] is different from the one in the batch [%s] err:%w", batchNumber, calculatedStateRoot, expectedStateRoot, ErrFatalBatchDesynchronized)
	}
	if calculatedLER != expectedLER {
		return fmt.Errorf("batch %v: LocalExitRoot calculated [%s] is different from the one in the batch [%s] err:%w", batchNumber, calculatedLER, expectedLER, ErrFatalBatchDesynchronized)
	}
	return nil
}

func checkProcessBatchResultMatchExpected(data *ProcessData, processBatchResp *state.ProcessBatchResponse) error {
	var err error = nil
	var trustedBatch = data.TrustedBatch
	if trustedBatch == nil {
		err = fmt.Errorf("%s trustedBatch is nil, it never should be nil", data.DebugPrefix)
		log.Error(err.Error())
		return err
	}
	if len(trustedBatch.BatchL2Data) == 0 {
		log.Warnf("Batch %v: BatchL2Data is empty, no checking", trustedBatch.Number)
		return nil
	}
	if processBatchResp == nil {
		log.Warnf("Batch %v: Can't check  processBatchResp because is nil, then check store batch in DB", trustedBatch.Number)
		err = checkStateRootAndLER(uint64(trustedBatch.Number), trustedBatch.StateRoot, trustedBatch.LocalExitRoot, data.StateBatch.StateRoot, data.StateBatch.LocalExitRoot)
	} else {
		err = checkStateRootAndLER(uint64(trustedBatch.Number), trustedBatch.StateRoot, trustedBatch.LocalExitRoot, processBatchResp.NewStateRoot, processBatchResp.NewLocalExitRoot)
	}
	if err != nil {
		log.Error(err.Error())
	}
	return err
}
