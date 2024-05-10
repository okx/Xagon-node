/*
object TrustedBatchesRetrieve:
- It get all pending batches from trusted node to be synchronized

You must implements BatchProcessor with the code to process the batches
*/
package l2_shared

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/common/syncinterfaces"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/metrics"
	"github.com/jackc/pgx/v4"
)

const (
	firstTrustedBatchNumber = uint64(2)
)

// StateInterface contains the methods required to interact with the state.
type StateInterface interface {
	BeginStateTransaction(ctx context.Context) (pgx.Tx, error)
	GetBatchByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*state.Batch, error)
}

// BatchProcessor is a interface with the ProcessTrustedBatch methor
//
//	this method is responsible to process a trusted batch
type BatchProcessor interface {
	// ProcessTrustedBatch processes a trusted batch
	ProcessTrustedBatch(ctx context.Context, trustedBatch *types.Batch, status TrustedState, dbTx pgx.Tx, debugPrefix string) (*TrustedState, error)
}

// TrustedState is the trusted state, basically contains the batch cache

// TrustedBatchesRetrieve it gets pending batches from Trusted node. It calls for each batch to BatchExecutor
//
//	and for each new batch calls the ProcessTrustedBatch method of the BatchExecutor interface
type TrustedBatchesRetrieve struct {
	batchExecutor          BatchProcessor
	zkEVMClient            syncinterfaces.ZKEVMClientTrustedBatchesGetter
	state                  StateInterface
	sync                   syncinterfaces.SynchronizerFlushIDManager
	TrustedStateMngr       TrustedStateManager
	firstBatchNumberToSync uint64
}

// NewTrustedBatchesRetrieve creates a new SyncTrustedStateTemplate
func NewTrustedBatchesRetrieve(batchExecutor BatchProcessor,
	zkEVMClient syncinterfaces.ZKEVMClientTrustedBatchesGetter,
	state StateInterface,
	sync syncinterfaces.SynchronizerFlushIDManager,
	TrustedStateMngr TrustedStateManager,
) *TrustedBatchesRetrieve {
	return &TrustedBatchesRetrieve{
		batchExecutor:          batchExecutor,
		zkEVMClient:            zkEVMClient,
		state:                  state,
		sync:                   sync,
		TrustedStateMngr:       TrustedStateMngr,
		firstBatchNumberToSync: firstTrustedBatchNumber,
	}
}

// CleanTrustedState Clean cache of TrustedBatches and StateRoot
func (s *TrustedBatchesRetrieve) CleanTrustedState() {
	s.TrustedStateMngr.Clear()
}

// GetCachedBatch implements syncinterfaces.SyncTrustedStateExecutor. Returns a cached batch
func (s *TrustedBatchesRetrieve) GetCachedBatch(batchNumber uint64) *state.Batch {
	return s.TrustedStateMngr.Cache.GetOrDefault(batchNumber, nil)
}

// SyncTrustedState sync trusted state from latestSyncedBatch to lastTrustedStateBatchNumber
func (s *TrustedBatchesRetrieve) SyncTrustedState(ctx context.Context, latestSyncedBatch uint64, maximumBatchNumberToProcess uint64) error {
	log.Info("syncTrustedState: Getting trusted state info")
	if latestSyncedBatch == 0 {
		log.Info("syncTrustedState: latestSyncedBatch is 0, assuming first batch as 1")
		latestSyncedBatch = 1
	}
	lastTrustedStateBatchNumberSeen, err := s.zkEVMClient.BatchNumber(ctx)

	if err != nil {
		log.Warn("syncTrustedState: error getting last batchNumber from Trusted Node. Error: ", err)
		return err
	}
	lastTrustedStateBatchNumber := min(lastTrustedStateBatchNumberSeen, maximumBatchNumberToProcess)
	log.Infof("syncTrustedState: latestSyncedBatch:%d syncTrustedState:%d (max Batch on network: %d)", latestSyncedBatch, lastTrustedStateBatchNumber, lastTrustedStateBatchNumberSeen)

	if isSyncrhonizedTrustedState(lastTrustedStateBatchNumber, latestSyncedBatch, s.firstBatchNumberToSync) {
		log.Info("syncTrustedState: Trusted state is synchronized")
		return nil
	}
	return s.syncTrustedBatchesToFrom(ctx, latestSyncedBatch, lastTrustedStateBatchNumber)
}

func isSyncrhonizedTrustedState(lastTrustedStateBatchNumber uint64, latestSyncedBatch uint64, firstBatchNumberToSync uint64) bool {
	if lastTrustedStateBatchNumber < firstBatchNumberToSync {
		return true
	}
	return lastTrustedStateBatchNumber < latestSyncedBatch
}

func sanityCheckBatchReturnedByTrusted(batch *types.Batch, expectedBatchNumber uint64) error {
	if batch == nil {
		return fmt.Errorf("batch %d is nil", expectedBatchNumber)
	}
	if uint64(batch.Number) != expectedBatchNumber {
		return fmt.Errorf("batch %d is not the expected batch %d", batch.Number, expectedBatchNumber)
	}
	return nil
}

func (s *TrustedBatchesRetrieve) syncTrustedBatchesToFrom(ctx context.Context, latestSyncedBatch uint64, lastTrustedStateBatchNumber uint64) error {
	batchNumberToSync := max(latestSyncedBatch, s.firstBatchNumberToSync)
	for batchNumberToSync <= lastTrustedStateBatchNumber {
		debugPrefix := fmt.Sprintf("syncTrustedState: batch[%d/%d]", batchNumberToSync, lastTrustedStateBatchNumber)
		start := time.Now()
		batchToSync, err := s.zkEVMClient.BatchByNumber(ctx, big.NewInt(0).SetUint64(batchNumberToSync))
		metrics.GetTrustedBatchInfoTime(time.Since(start))
		if err != nil {
			log.Warnf("%s failed to get batch %d from trusted state. Error: %v", debugPrefix, batchNumberToSync, err)
			return err
		}
		err = sanityCheckBatchReturnedByTrusted(batchToSync, batchNumberToSync)
		if err != nil {
			log.Warnf("%s sanity check over Batch returned by Trusted-RPC failed: %v", debugPrefix, err)
			return err
		}

		dbTx, err := s.state.BeginStateTransaction(ctx)
		if err != nil {
			log.Errorf("%s error creating db transaction to sync trusted batch %d: %v", debugPrefix, batchNumberToSync, err)
			return err
		}
		start = time.Now()
		previousStatus, err := s.TrustedStateMngr.GetStateForWorkingBatch(ctx, batchNumberToSync, s.state, dbTx)
		if err != nil {
			log.Errorf("%s error getting current batches to sync trusted batch %d: %v", debugPrefix, batchNumberToSync, err)
			return rollback(ctx, dbTx, err)
		}
		log.Debugf("%s processing trusted batch %d", debugPrefix, batchNumberToSync)
		newTrustedState, err := s.batchExecutor.ProcessTrustedBatch(ctx, batchToSync, *previousStatus, dbTx, debugPrefix)
		metrics.ProcessTrustedBatchTime(time.Since(start))
		if err != nil {
			log.Errorf("%s error processing trusted batch %d: %v", debugPrefix, batchNumberToSync, err)
			s.TrustedStateMngr.Clear()
			return rollback(ctx, dbTx, err)
		}
		log.Debugf("%s Checking FlushID to commit trustedState data to db", debugPrefix)
		err = s.sync.CheckFlushID(dbTx)
		if err != nil {
			log.Errorf("%s error checking flushID. Error: %v", debugPrefix, err)
			s.TrustedStateMngr.Clear()
			return rollback(ctx, dbTx, err)
		}

		if err := dbTx.Commit(ctx); err != nil {
			log.Errorf("%s error committing db transaction to sync trusted batch %v: %v", debugPrefix, batchNumberToSync, err)
			s.TrustedStateMngr.Clear()
			return err
		}
		// Update cache with result
		if newTrustedState != nil {
			s.TrustedStateMngr.Set(newTrustedState.LastTrustedBatches[0])
			s.TrustedStateMngr.Set(newTrustedState.LastTrustedBatches[1])
		} else {
			s.TrustedStateMngr.Clear()
		}
		batchNumberToSync++
		if !batchToSync.Closed && batchNumberToSync <= lastTrustedStateBatchNumber {
			log.Infof("%s Batch %d is not closed. so we break synchronization from Trusted Node because can only have 1 WIP batch on state", debugPrefix, batchToSync.Number)
			return nil
		}
	}

	log.Infof("syncTrustedState: Trusted state fully synchronized from %d to %d", latestSyncedBatch, lastTrustedStateBatchNumber)
	return nil
}

func rollback(ctx context.Context, dbTx pgx.Tx, err error) error {
	rollbackErr := dbTx.Rollback(ctx)
	if rollbackErr != nil {
		log.Errorf("syncTrustedState: error rolling back state. RollbackErr: %s, Error : %v", rollbackErr.Error(), err)
		return rollbackErr
	}
	return err
}
