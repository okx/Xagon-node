package etrog

import (
	"context"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/metrics"
	stateMetrics "github.com/0xPolygonHermez/zkevm-node/state/metrics"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/actions"
	syncCommon "github.com/0xPolygonHermez/zkevm-node/synchronizer/common"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

type stateProcessUpdateEtrogSequence interface {
	ProcessAndStoreClosedBatchV2(ctx context.Context, processingCtx state.ProcessingContextV2, dbTx pgx.Tx, caller metrics.CallerLabel) (common.Hash, uint64, string, error)
	AddSequence(ctx context.Context, sequence state.Sequence, dbTx pgx.Tx) error
	AddVirtualBatch(ctx context.Context, virtualBatch *state.VirtualBatch, dbTx pgx.Tx) error
}

type syncProcessUpdateEtrogSequenceInterface interface {
	PendingFlushID(flushID uint64, proverID string)
}

// ProcessorL1UpdateEtrogSequence implements L1EventProcessor
type ProcessorL1UpdateEtrogSequence struct {
	actions.ProcessorBase[ProcessorL1UpdateEtrogSequence]
	state        stateProcessUpdateEtrogSequence
	sync         syncProcessUpdateEtrogSequenceInterface
	timeProvider syncCommon.TimeProvider
}

// NewProcessorL1UpdateEtrogSequence returns instance of a processor for UpdateEtrogSequenceOrder
func NewProcessorL1UpdateEtrogSequence(state stateProcessUpdateEtrogSequence,
	sync syncProcessUpdateEtrogSequenceInterface,
	timeProvider syncCommon.TimeProvider) *ProcessorL1UpdateEtrogSequence {
	return &ProcessorL1UpdateEtrogSequence{
		ProcessorBase: *actions.NewProcessorBase[ProcessorL1UpdateEtrogSequence](
			[]etherman.EventOrder{etherman.UpdateEtrogSequenceOrder},
			actions.ForksIdOnlyEtrog),
		state:        state,
		sync:         sync,
		timeProvider: timeProvider,
	}
}

// Process process event
func (g *ProcessorL1UpdateEtrogSequence) Process(ctx context.Context, order etherman.Order, l1Block *etherman.Block, dbTx pgx.Tx) error {
	if l1Block == nil || l1Block.UpdateEtrogSequence.BatchNumber == 0 {
		return actions.ErrInvalidParams
	}
	err := g.processUpdateEtrogSequence(ctx, l1Block.UpdateEtrogSequence, l1Block.BlockNumber, l1Block.ReceivedAt, dbTx)
	return err
}

func (g *ProcessorL1UpdateEtrogSequence) processUpdateEtrogSequence(ctx context.Context, updateEtrogSequence etherman.UpdateEtrogSequence, blockNumber uint64, l1BlockTimestamp time.Time, dbTx pgx.Tx) error {
	now := g.timeProvider.Now()
	batch := state.Batch{
		BatchNumber:    updateEtrogSequence.BatchNumber,
		GlobalExitRoot: updateEtrogSequence.PolygonRollupBaseEtrogBatchData.ForcedGlobalExitRoot,
		// This timestamp now is the timeLimit. It can't be the one virtual.BatchTimestamp
		//   because when sync from trusted we don't now the real BatchTimestamp and
		//   will fails the comparation of batch time >= than previous one.
		Timestamp:   now,
		Coinbase:    updateEtrogSequence.SequencerAddr,
		BatchL2Data: updateEtrogSequence.PolygonRollupBaseEtrogBatchData.Transactions,
	}

	log.Debug("Processing update etrog sequence batch")
	var fBHL1 common.Hash = updateEtrogSequence.PolygonRollupBaseEtrogBatchData.ForcedBlockHashL1
	forcedBlockHashL1 := &fBHL1
	txs := updateEtrogSequence.PolygonRollupBaseEtrogBatchData.Transactions
	tstampLimit := time.Unix(int64(updateEtrogSequence.PolygonRollupBaseEtrogBatchData.ForcedTimestamp), 0)
	processCtx := state.ProcessingContextV2{
		BatchNumber:          updateEtrogSequence.BatchNumber,
		Coinbase:             updateEtrogSequence.SequencerAddr,
		Timestamp:            &tstampLimit,
		L1InfoRoot:           updateEtrogSequence.PolygonRollupBaseEtrogBatchData.ForcedGlobalExitRoot,
		BatchL2Data:          &txs,
		ForcedBlockHashL1:    forcedBlockHashL1,
		SkipVerifyL1InfoRoot: 1,
		GlobalExitRoot:       updateEtrogSequence.PolygonRollupBaseEtrogBatchData.ForcedGlobalExitRoot,
		ClosingReason:        state.SyncL1EventUpdateEtrogSequenceClosingReason,
	}

	virtualBatch := state.VirtualBatch{
		BatchNumber:         updateEtrogSequence.BatchNumber,
		TxHash:              updateEtrogSequence.TxHash,
		Coinbase:            updateEtrogSequence.SequencerAddr,
		BlockNumber:         blockNumber,
		SequencerAddr:       updateEtrogSequence.SequencerAddr,
		TimestampBatchEtrog: &l1BlockTimestamp,
		L1InfoRoot:          &processCtx.L1InfoRoot,
	}

	log.Debugf("Storing batchNumber: %d...", batch.BatchNumber)
	// If it is not found, store batch
	_, flushID, proverID, err := g.state.ProcessAndStoreClosedBatchV2(ctx, processCtx, dbTx, stateMetrics.SynchronizerCallerLabel)
	if err != nil {
		log.Errorf("error storing trustedBatch. BatchNumber: %d, BlockNumber: %d, error: %v", batch.BatchNumber, blockNumber, err)
		rollbackErr := dbTx.Rollback(ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %s, error : %v", batch.BatchNumber, blockNumber, rollbackErr.Error(), err)
			return rollbackErr
		}
		log.Errorf("error storing batch. BatchNumber: %d, BlockNumber: %d, error: %v", batch.BatchNumber, blockNumber, err)
		return err
	}
	g.sync.PendingFlushID(flushID, proverID)

	// Store virtualBatch
	log.Infof("processUpdateEtrogSequence: Storing virtualBatch. BatchNumber: %d, BlockNumber: %d GER:%s", virtualBatch.BatchNumber, blockNumber,
		common.Hash(updateEtrogSequence.PolygonRollupBaseEtrogBatchData.ForcedGlobalExitRoot).String())
	err = g.state.AddVirtualBatch(ctx, &virtualBatch, dbTx)
	if err != nil {
		log.Errorf("error storing virtualBatch. BatchNumber: %d, BlockNumber: %d, error: %v", virtualBatch.BatchNumber, blockNumber, err)
		rollbackErr := dbTx.Rollback(ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %s, error : %v", virtualBatch.BatchNumber, blockNumber, rollbackErr.Error(), err)
			return rollbackErr
		}
		log.Errorf("error storing virtualBatch. BatchNumber: %d, BlockNumber: %d, error: %v", virtualBatch.BatchNumber, blockNumber, err)
		return err
	}
	// Insert the sequence to allow the aggregator verify the sequence batches
	seq := state.Sequence{
		FromBatchNumber: updateEtrogSequence.BatchNumber,
		ToBatchNumber:   updateEtrogSequence.BatchNumber,
	}
	err = g.state.AddSequence(ctx, seq, dbTx)
	if err != nil {
		log.Errorf("error adding sequence. Sequence: %+v", seq)
		rollbackErr := dbTx.Rollback(ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state. BlockNumber: %d, rollbackErr: %s, error : %v", blockNumber, rollbackErr.Error(), err)
			return rollbackErr
		}
		log.Errorf("error getting adding sequence. BlockNumber: %d, error: %v", blockNumber, err)
		return err
	}
	return nil
}
