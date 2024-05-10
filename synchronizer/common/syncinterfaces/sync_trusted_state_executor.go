package syncinterfaces

import (
	"context"
	"errors"

	"github.com/0xPolygonHermez/zkevm-node/state"
)

var (
	// ErrMissingSyncFromL1 is returned when we are behind expected L1 sync, so must be done from L1
	ErrMissingSyncFromL1 = errors.New("must sync from L1")
	// ErrFatalDesyncFromL1 is returned when trusted node and permissionless node have different data
	ErrFatalDesyncFromL1 = errors.New("fatal situation: the TrustedNode have another data!. Halt or do something")
	// ErrCantSyncFromL2 is returned when can't sync from L2, for example the forkid is not supported by L2 sync
	ErrCantSyncFromL2 = errors.New("can't sync from L2")
)

// SyncTrustedStateExecutor is the interface that class that synchronize permissionless with a trusted node
type SyncTrustedStateExecutor interface {
	// SyncTrustedState syncs the trusted state with the permissionless state
	//  maximumBatchToProcess: maximum Batchnumber of batches to process, after have to returns
	// if returns error ErrMissingSyncFromL1 then must force a L1 sync
	//
	SyncTrustedState(ctx context.Context, latestSyncedBatch uint64, maximumBatchNumberToProcess uint64) error
	// CleanTrustedState clean cache of Batches and StateRoot
	CleanTrustedState()
	// Returns the cached data for a batch
	GetCachedBatch(batchNumber uint64) *state.Batch
}
