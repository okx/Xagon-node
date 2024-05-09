package l2_shared

/*
This class is a implementation of SyncTrustedStateExecutor that selects the executor to use.
It have a map with the forkID and the executor class to use, if none is available skip trusted sync returning a nil
*/

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/common/syncinterfaces"
)

type stateSyncTrustedStateExecutorSelector interface {
	GetForkIDInMemory(forkId uint64) *state.ForkIDInterval
	GetForkIDByBatchNumber(batchNumber uint64) uint64
}

// SyncTrustedStateExecutorSelector Implements SyncTrustedStateExecutor
type SyncTrustedStateExecutorSelector struct {
	state          stateSyncTrustedStateExecutorSelector
	supportedForks map[uint64]syncinterfaces.SyncTrustedStateExecutor
}

// NewSyncTrustedStateExecutorSelector creates a new SyncTrustedStateExecutorSelector that implements SyncTrustedStateExecutor
func NewSyncTrustedStateExecutorSelector(
	supportedForks map[uint64]syncinterfaces.SyncTrustedStateExecutor,
	state stateSyncTrustedStateExecutorSelector) *SyncTrustedStateExecutorSelector {
	return &SyncTrustedStateExecutorSelector{
		supportedForks: supportedForks,
		state:          state,
	}
}

// GetExecutor returns the executor that should be used for the given batch, could be nil
// it returns the executor and the maximum batch number that the executor can process
func (s *SyncTrustedStateExecutorSelector) GetExecutor(latestSyncedBatch uint64, maximumBatchNumberToProcess uint64) (syncinterfaces.SyncTrustedStateExecutor, uint64) {
	forkIDForNextBatch := s.state.GetForkIDByBatchNumber(latestSyncedBatch + 1)
	executor, ok := s.supportedForks[forkIDForNextBatch]
	if !ok {
		log.Warnf("No supported sync from Trusted Node for  forkID %d", forkIDForNextBatch)
		return nil, 0
	}
	fork := s.state.GetForkIDInMemory(forkIDForNextBatch)
	if fork == nil {
		log.Errorf("ForkID %d range not available! that is UB", forkIDForNextBatch)
		return nil, 0
	}

	maxCapped := min(maximumBatchNumberToProcess, fork.ToBatchNumber)
	log.Debugf("using ForkID %d, lastBatch:%d  (maxBatch original:%d  capped:%d)", forkIDForNextBatch,
		latestSyncedBatch, maximumBatchNumberToProcess, maxCapped)
	return executor, maxCapped
}

// SyncTrustedState syncs the trusted state with the permissionless state. In this case
// choose which executor must use
func (s *SyncTrustedStateExecutorSelector) SyncTrustedState(ctx context.Context, latestSyncedBatch uint64, maximumBatchNumberToProcess uint64) error {
	executor, maxBatchNumber := s.GetExecutor(latestSyncedBatch, maximumBatchNumberToProcess)
	if executor == nil {
		log.Warnf("No executor available, skipping SyncTrustedState: latestSyncedBatch:%d, maximumBatchNumberToProcess:%d",
			latestSyncedBatch, maximumBatchNumberToProcess)
		return syncinterfaces.ErrCantSyncFromL2
	}
	return executor.SyncTrustedState(ctx, latestSyncedBatch, maxBatchNumber)
}

// CleanTrustedState clean cache of Batches and StateRoot
func (s *SyncTrustedStateExecutorSelector) CleanTrustedState() {
	for _, executor := range s.supportedForks {
		executor.CleanTrustedState()
	}
}

// GetCachedBatch implements syncinterfaces.SyncTrustedStateExecutor. Returns a cached batch
func (s *SyncTrustedStateExecutorSelector) GetCachedBatch(batchNumber uint64) *state.Batch {
	executor, _ := s.GetExecutor(batchNumber, 0)
	if executor == nil {
		return nil
	}
	return executor.GetCachedBatch(min(batchNumber))
}
