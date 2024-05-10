package l2_shared

import (
	"context"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/common"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/common/syncinterfaces"
	"github.com/jackc/pgx/v4"
)

// TrustedState is the trusted state, basically contains the batch cache for a concrete batch
type TrustedState struct {
	// LastTrustedBatches [0] -> Current  batch, [1] -> previous batch
	LastTrustedBatches []*state.Batch
}

// IsEmpty returns true if the trusted state is empty
func (ts *TrustedState) IsEmpty() bool {
	if ts == nil || len(ts.LastTrustedBatches) == 0 {
		return true
	}
	if len(ts.LastTrustedBatches) == 1 && ts.LastTrustedBatches[0] == nil {
		return true
	}
	if len(ts.LastTrustedBatches) > 1 && ts.LastTrustedBatches[0] == nil && ts.LastTrustedBatches[1] == nil {
		return true
	}
	return false
}

// GetCurrentBatch returns the current batch or nil
func (ts *TrustedState) GetCurrentBatch() *state.Batch {
	if ts == nil || len(ts.LastTrustedBatches) == 0 {
		return nil
	}
	return ts.LastTrustedBatches[0]
}

// GetPreviousBatch returns the previous batch or nil
func (ts *TrustedState) GetPreviousBatch() *state.Batch {
	if ts == nil || len(ts.LastTrustedBatches) < 2 {
		return nil
	}
	return ts.LastTrustedBatches[1]
}

// TrustedStateManager is the trusted state manager, basically contains the batch cache and create the TrustedState
type TrustedStateManager struct {
	Cache *common.Cache[uint64, *state.Batch]
}

// NewTrustedStateManager creates a new TrustedStateManager
func NewTrustedStateManager(timerProvider common.TimeProvider, timeOfLiveItems time.Duration) *TrustedStateManager {
	return &TrustedStateManager{
		Cache: common.NewCache[uint64, *state.Batch](timerProvider, timeOfLiveItems),
	}
}

// Clear clears the cache
func (ts *TrustedStateManager) Clear() {
	ts.Cache.Clear()
}

// Set sets the result batch in the cache
func (ts *TrustedStateManager) Set(resultBatch *state.Batch) {
	if resultBatch == nil {
		return
	}
	ts.Cache.Set(resultBatch.BatchNumber, resultBatch)
}

// GetStateForWorkingBatch returns the trusted state for the working batch
func (ts *TrustedStateManager) GetStateForWorkingBatch(ctx context.Context, batchNumber uint64, stateGetBatch syncinterfaces.StateGetBatchByNumberInterface, dbTx pgx.Tx) (*TrustedState, error) {
	ts.Cache.DeleteOutdated()
	res := &TrustedState{}
	var err error
	var currentBatch, previousBatch *state.Batch
	currentBatch = ts.Cache.GetOrDefault(batchNumber, nil)
	previousBatch = ts.Cache.GetOrDefault(batchNumber-1, nil)
	if currentBatch == nil {
		currentBatch, err = stateGetBatch.GetBatchByNumber(ctx, batchNumber, dbTx)
		if err != nil && err != state.ErrNotFound {
			log.Warnf("failed to get batch %v from local trusted state. Error: %v", batchNumber, err)
			return nil, err
		} else {
			ts.Cache.Set(batchNumber, currentBatch)
		}
	}
	if previousBatch == nil {
		previousBatch, err = stateGetBatch.GetBatchByNumber(ctx, batchNumber-1, dbTx)
		if err != nil && err != state.ErrNotFound {
			log.Warnf("failed to get batch %v from local trusted state. Error: %v", batchNumber-1, err)
			return nil, err
		} else {
			ts.Cache.Set(batchNumber-1, previousBatch)
		}
	}
	res.LastTrustedBatches = []*state.Batch{currentBatch, previousBatch}
	return res, nil
}
