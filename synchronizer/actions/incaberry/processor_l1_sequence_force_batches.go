package incaberry

import (
	"context"
	"fmt"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/metrics"
	stateMetrics "github.com/0xPolygonHermez/zkevm-node/state/metrics"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/actions"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

type stateProcessL1SequenceForcedBatchesInterface interface {
	GetLastVirtualBatchNum(ctx context.Context, dbTx pgx.Tx) (uint64, error)
	ResetTrustedState(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) error
	GetNextForcedBatches(ctx context.Context, nextForcedBatches int, dbTx pgx.Tx) ([]state.ForcedBatch, error)
	ProcessAndStoreClosedBatch(ctx context.Context, processingCtx state.ProcessingContext, encodedTxs []byte, dbTx pgx.Tx, caller metrics.CallerLabel) (common.Hash, uint64, string, error)
	AddVirtualBatch(ctx context.Context, virtualBatch *state.VirtualBatch, dbTx pgx.Tx) error
	AddSequence(ctx context.Context, sequence state.Sequence, dbTx pgx.Tx) error
}

type syncProcessL1SequenceForcedBatchesInterface interface {
	PendingFlushID(flushID uint64, proverID string)
	CleanTrustedState()
}

// ProcessL1SequenceForcedBatches implements L1EventProcessor
type ProcessL1SequenceForcedBatches struct {
	actions.ProcessorBase[ProcessL1SequenceForcedBatches]
	state stateProcessL1SequenceForcedBatchesInterface
	sync  syncProcessL1SequenceForcedBatchesInterface
}

// NewProcessL1SequenceForcedBatches returns instance of a processor for SequenceForceBatchesOrder
func NewProcessL1SequenceForcedBatches(state stateProcessL1SequenceForcedBatchesInterface,
	sync syncProcessL1SequenceForcedBatchesInterface) *ProcessL1SequenceForcedBatches {
	return &ProcessL1SequenceForcedBatches{
		ProcessorBase: *actions.NewProcessorBase[ProcessL1SequenceForcedBatches](
			[]etherman.EventOrder{etherman.SequenceForceBatchesOrder},
			actions.ForksIdAll),
		state: state,
		sync:  sync}
}

// Process process event
func (p *ProcessL1SequenceForcedBatches) Process(ctx context.Context, order etherman.Order, l1Block *etherman.Block, dbTx pgx.Tx) error {
	return p.processSequenceForceBatch(ctx, l1Block.SequencedForceBatches[order.Pos], *l1Block, dbTx)
}

func (s *ProcessL1SequenceForcedBatches) processSequenceForceBatch(ctx context.Context, sequenceForceBatch []etherman.SequencedForceBatch, block etherman.Block, dbTx pgx.Tx) error {
	if len(sequenceForceBatch) == 0 {
		log.Warn("Empty sequenceForceBatch array detected, ignoring...")
		return nil
	}
	// First, get last virtual batch number
	lastVirtualizedBatchNumber, err := s.state.GetLastVirtualBatchNum(ctx, dbTx)
	if err != nil {
		log.Errorf("error getting lastVirtualBatchNumber. BlockNumber: %d, error: %v", block.BlockNumber, err)
		rollbackErr := dbTx.Rollback(ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %s, error : %v", lastVirtualizedBatchNumber, block.BlockNumber, rollbackErr.Error(), err)
			return rollbackErr
		}
		log.Errorf("error getting lastVirtualBatchNumber. BlockNumber: %d, error: %v", block.BlockNumber, err)
		return err
	}
	// Clean trustedState sync variables to avoid sync the trusted state from the wrong starting point.
	// This wrong starting point would force the trusted sync to clean the virtualization of the batch reaching an inconsistency.
	s.sync.CleanTrustedState()

	// Reset trusted state
	log.Infof("ResetTrustedState: processSequenceForceBatch: Resetting trusted state. delete batch > (lastVirtualizedBatchNumber)%d, ", lastVirtualizedBatchNumber)
	err = s.state.ResetTrustedState(ctx, lastVirtualizedBatchNumber, dbTx) // This method has to reset the forced batches deleting the batchNumber for higher batchNumbers
	if err != nil {
		log.Errorf("error resetting trusted state. BatchNumber: %d, BlockNumber: %d, error: %v", lastVirtualizedBatchNumber, block.BlockNumber, err)
		rollbackErr := dbTx.Rollback(ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %s, error : %v", lastVirtualizedBatchNumber, block.BlockNumber, rollbackErr.Error(), err)
			return rollbackErr
		}
		log.Errorf("error resetting trusted state. BatchNumber: %d, BlockNumber: %d, error: %v", lastVirtualizedBatchNumber, block.BlockNumber, err)
		return err
	}
	// Read forcedBatches from db
	forcedBatches, err := s.state.GetNextForcedBatches(ctx, len(sequenceForceBatch), dbTx)
	if err != nil {
		log.Errorf("error getting forcedBatches in processSequenceForceBatch. BlockNumber: %d", block.BlockNumber)
		rollbackErr := dbTx.Rollback(ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state. BlockNumber: %d, rollbackErr: %s, error : %v", block.BlockNumber, rollbackErr.Error(), err)
			return rollbackErr
		}
		log.Errorf("error getting forcedBatches in processSequenceForceBatch. BlockNumber: %d, error: %v", block.BlockNumber, err)
		return err
	}
	if len(sequenceForceBatch) != len(forcedBatches) {
		rollbackErr := dbTx.Rollback(ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state. BlockNumber: %d, rollbackErr: %v", block.BlockNumber, rollbackErr)
			return rollbackErr
		}
		log.Error("error number of forced batches doesn't match")
		return fmt.Errorf("error number of forced batches doesn't match")
	}
	for i, fbatch := range sequenceForceBatch {
		if uint64(forcedBatches[i].ForcedAt.Unix()) != fbatch.ForcedTimestamp ||
			forcedBatches[i].GlobalExitRoot != fbatch.ForcedGlobalExitRoot ||
			common.Bytes2Hex(forcedBatches[i].RawTxsData) != common.Bytes2Hex(fbatch.Transactions) {
			log.Warnf("ForcedBatch stored: %+v", forcedBatches)
			log.Warnf("ForcedBatch sequenced received: %+v", fbatch)
			log.Errorf("error: forcedBatch received doesn't match with the next expected forcedBatch stored in db. Expected: %+v, Synced: %+v", forcedBatches[i], fbatch)
			rollbackErr := dbTx.Rollback(ctx)
			if rollbackErr != nil {
				log.Errorf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %v", fbatch.BatchNumber, block.BlockNumber, rollbackErr)
				return rollbackErr
			}
			return fmt.Errorf("error: forcedBatch received doesn't match with the next expected forcedBatch stored in db. Expected: %+v, Synced: %+v", forcedBatches[i], fbatch)
		}
		virtualBatch := state.VirtualBatch{
			BatchNumber:   fbatch.BatchNumber,
			TxHash:        fbatch.TxHash,
			Coinbase:      fbatch.Coinbase,
			SequencerAddr: fbatch.Coinbase,
			BlockNumber:   block.BlockNumber,
		}
		batch := state.ProcessingContext{
			BatchNumber:    fbatch.BatchNumber,
			GlobalExitRoot: fbatch.ForcedGlobalExitRoot,
			Timestamp:      block.ReceivedAt,
			Coinbase:       fbatch.Coinbase,
			ForcedBatchNum: &forcedBatches[i].ForcedBatchNumber,
			BatchL2Data:    &forcedBatches[i].RawTxsData,
		}
		// Process batch
		log.Infof("processSequenceFoceBatches: ProcessAndStoreClosedBatch . BatchNumber: %d, BlockNumber: %d", batch.BatchNumber, block.BlockNumber)
		_, flushID, proverID, err := s.state.ProcessAndStoreClosedBatch(ctx, batch, forcedBatches[i].RawTxsData, dbTx, stateMetrics.SynchronizerCallerLabel)
		if err != nil {
			log.Errorf("error processing batch in processSequenceForceBatch. BatchNumber: %d, BlockNumber: %d, error: %v", batch.BatchNumber, block.BlockNumber, err)
			rollbackErr := dbTx.Rollback(ctx)
			if rollbackErr != nil {
				log.Errorf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %s, error : %v", batch.BatchNumber, block.BlockNumber, rollbackErr.Error(), err)
				return rollbackErr
			}
			log.Errorf("error processing batch in processSequenceForceBatch. BatchNumber: %d, BlockNumber: %d, error: %v", batch.BatchNumber, block.BlockNumber, err)
			return err
		}
		s.sync.PendingFlushID(flushID, proverID)

		// Store virtualBatch
		log.Infof("processSequenceFoceBatches: Storing virtualBatch. BatchNumber: %d, BlockNumber: %d", virtualBatch.BatchNumber, block.BlockNumber)
		err = s.state.AddVirtualBatch(ctx, &virtualBatch, dbTx)
		if err != nil {
			log.Errorf("error storing virtualBatch in processSequenceForceBatch. BatchNumber: %d, BlockNumber: %d, error: %v", virtualBatch.BatchNumber, block.BlockNumber, err)
			rollbackErr := dbTx.Rollback(ctx)
			if rollbackErr != nil {
				log.Errorf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %s, error : %v", virtualBatch.BatchNumber, block.BlockNumber, rollbackErr.Error(), err)
				return rollbackErr
			}
			log.Errorf("error storing virtualBatch in processSequenceForceBatch. BatchNumber: %d, BlockNumber: %d, error: %v", virtualBatch.BatchNumber, block.BlockNumber, err)
			return err
		}
	}
	// Insert the sequence to allow the aggregator verify the sequence batches
	seq := state.Sequence{
		FromBatchNumber: sequenceForceBatch[0].BatchNumber,
		ToBatchNumber:   sequenceForceBatch[len(sequenceForceBatch)-1].BatchNumber,
	}
	err = s.state.AddSequence(ctx, seq, dbTx)
	if err != nil {
		log.Errorf("error adding sequence. Sequence: %+v", seq)
		rollbackErr := dbTx.Rollback(ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state. BlockNumber: %d, rollbackErr: %s, error : %v", block.BlockNumber, rollbackErr.Error(), err)
			return rollbackErr
		}
		log.Errorf("error getting adding sequence. BlockNumber: %d, error: %v", block.BlockNumber, err)
		return err
	}
	return nil
}
