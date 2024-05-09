package incaberry

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/event"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/metrics"
	stateMetrics "github.com/0xPolygonHermez/zkevm-node/state/metrics"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/actions"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/common/syncinterfaces"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
)

type stateProcessSequenceBatches interface {
	GetNextForcedBatches(ctx context.Context, nextForcedBatches int, dbTx pgx.Tx) ([]state.ForcedBatch, error)
	GetBatchByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*state.Batch, error)
	ProcessAndStoreClosedBatch(ctx context.Context, processingCtx state.ProcessingContext, encodedTxs []byte, dbTx pgx.Tx, caller metrics.CallerLabel) (common.Hash, uint64, string, error)
	ExecuteBatch(ctx context.Context, batch state.Batch, updateMerkleTree bool, dbTx pgx.Tx) (*executor.ProcessBatchResponse, error)
	AddAccumulatedInputHash(ctx context.Context, batchNum uint64, accInputHash common.Hash, dbTx pgx.Tx) error
	ResetTrustedState(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) error
	AddSequence(ctx context.Context, sequence state.Sequence, dbTx pgx.Tx) error
	AddVirtualBatch(ctx context.Context, virtualBatch *state.VirtualBatch, dbTx pgx.Tx) error
	AddTrustedReorg(ctx context.Context, trustedReorg *state.TrustedReorg, dbTx pgx.Tx) error
	GetReorgedTransactions(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) ([]*ethTypes.Transaction, error)
}

type ethermanProcessSequenceBatches interface {
	GetLatestBatchNumber() (uint64, error)
}

type poolProcessSequenceBatchesInterface interface {
	DeleteReorgedTransactions(ctx context.Context, txs []*ethTypes.Transaction) error
	StoreTx(ctx context.Context, tx ethTypes.Transaction, ip string, isWIP bool) error
}

type syncProcessSequenceBatchesInterface interface {
	PendingFlushID(flushID uint64, proverID string)
	IsTrustedSequencer() bool
	CleanTrustedState()
}

// ProcessorL1SequenceBatches implements L1EventProcessor
type ProcessorL1SequenceBatches struct {
	actions.ProcessorBase[ProcessorL1SequenceBatches]
	state    stateProcessSequenceBatches
	etherMan ethermanProcessSequenceBatches
	pool     poolProcessSequenceBatchesInterface
	eventLog syncinterfaces.EventLogInterface
	sync     syncProcessSequenceBatchesInterface
}

// NewProcessorL1SequenceBatches returns instance of a processor for SequenceBatchesOrder
func NewProcessorL1SequenceBatches(state stateProcessSequenceBatches,
	etherMan ethermanProcessSequenceBatches, pool poolProcessSequenceBatchesInterface, eventLog syncinterfaces.EventLogInterface, sync syncProcessSequenceBatchesInterface) *ProcessorL1SequenceBatches {
	return &ProcessorL1SequenceBatches{
		ProcessorBase: *actions.NewProcessorBase[ProcessorL1SequenceBatches](
			[]etherman.EventOrder{etherman.SequenceBatchesOrder},
			actions.ForksIdToIncaberry),
		state:    state,
		etherMan: etherMan,
		pool:     pool,
		eventLog: eventLog,
		sync:     sync,
	}
}

// Process process event
func (g *ProcessorL1SequenceBatches) Process(ctx context.Context, order etherman.Order, l1Block *etherman.Block, dbTx pgx.Tx) error {
	if l1Block == nil || len(l1Block.SequencedBatches) <= order.Pos {
		return actions.ErrInvalidParams
	}
	err := g.processSequenceBatches(ctx, l1Block.SequencedBatches[order.Pos], l1Block.BlockNumber, dbTx)
	return err
}

func (g *ProcessorL1SequenceBatches) processSequenceBatches(ctx context.Context, sequencedBatches []etherman.SequencedBatch, blockNumber uint64, dbTx pgx.Tx) error {
	if len(sequencedBatches) == 0 {
		log.Warn("Empty sequencedBatches array detected, ignoring...")
		return nil
	}
	for _, sbatch := range sequencedBatches {
		virtualBatch := state.VirtualBatch{
			BatchNumber:   sbatch.BatchNumber,
			TxHash:        sbatch.TxHash,
			Coinbase:      sbatch.Coinbase,
			BlockNumber:   blockNumber,
			SequencerAddr: sbatch.SequencerAddr,
		}
		batch := state.Batch{
			BatchNumber:    sbatch.BatchNumber,
			GlobalExitRoot: sbatch.PolygonZkEVMBatchData.GlobalExitRoot,
			Timestamp:      time.Unix(int64(sbatch.PolygonZkEVMBatchData.Timestamp), 0),
			Coinbase:       sbatch.Coinbase,
			BatchL2Data:    sbatch.PolygonZkEVMBatchData.Transactions,
		}
		// ForcedBatch must be processed
		if sbatch.PolygonZkEVMBatchData.MinForcedTimestamp > 0 { // If this is true means that the batch is forced
			log.Debug("FORCED BATCH SEQUENCED!")
			// Read forcedBatches from db
			forcedBatches, err := g.state.GetNextForcedBatches(ctx, 1, dbTx)
			if err != nil {
				log.Errorf("error getting forcedBatches. BatchNumber: %d", virtualBatch.BatchNumber)
				rollbackErr := dbTx.Rollback(ctx)
				if rollbackErr != nil {
					log.Errorf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %s, error : %v", virtualBatch.BatchNumber, blockNumber, rollbackErr.Error(), err)
					return rollbackErr
				}
				return err
			}
			if len(forcedBatches) == 0 {
				log.Errorf("error: empty forcedBatches array read from db. BatchNumber: %d", sbatch.BatchNumber)
				rollbackErr := dbTx.Rollback(ctx)
				if rollbackErr != nil {
					log.Errorf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %v", sbatch.BatchNumber, blockNumber, rollbackErr)
					return rollbackErr
				}
				return fmt.Errorf("error: empty forcedBatches array read from db. BatchNumber: %d", sbatch.BatchNumber)
			}
			if uint64(forcedBatches[0].ForcedAt.Unix()) != sbatch.PolygonZkEVMBatchData.MinForcedTimestamp ||
				forcedBatches[0].GlobalExitRoot != sbatch.PolygonZkEVMBatchData.GlobalExitRoot ||
				common.Bytes2Hex(forcedBatches[0].RawTxsData) != common.Bytes2Hex(sbatch.PolygonZkEVMBatchData.Transactions) {
				log.Warnf("ForcedBatch stored: %+v. RawTxsData: %s", forcedBatches, common.Bytes2Hex(forcedBatches[0].RawTxsData))
				log.Warnf("ForcedBatch sequenced received: %+v. RawTxsData: %s", sbatch, common.Bytes2Hex(sbatch.PolygonZkEVMBatchData.Transactions))
				log.Errorf("error: forcedBatch received doesn't match with the next expected forcedBatch stored in db. Expected: %+v, Synced: %+v", forcedBatches, sbatch)
				rollbackErr := dbTx.Rollback(ctx)
				if rollbackErr != nil {
					log.Errorf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %v", virtualBatch.BatchNumber, blockNumber, rollbackErr)
					return rollbackErr
				}
				return fmt.Errorf("error: forcedBatch received doesn't match with the next expected forcedBatch stored in db. Expected: %+v, Synced: %+v", forcedBatches, sbatch)
			}
			log.Debug("Setting forcedBatchNum: ", forcedBatches[0].ForcedBatchNumber)
			batch.ForcedBatchNum = &forcedBatches[0].ForcedBatchNumber
		}

		// Now we need to check the batch. ForcedBatches should be already stored in the batch table because this is done by the sequencer
		processCtx := state.ProcessingContext{
			BatchNumber:    batch.BatchNumber,
			Coinbase:       batch.Coinbase,
			Timestamp:      batch.Timestamp,
			GlobalExitRoot: batch.GlobalExitRoot,
			ForcedBatchNum: batch.ForcedBatchNum,
			BatchL2Data:    &batch.BatchL2Data,
		}

		var newRoot common.Hash

		// First get trusted batch from db
		tBatch, err := g.state.GetBatchByNumber(ctx, batch.BatchNumber, dbTx)
		if err != nil {
			if errors.Is(err, state.ErrNotFound) {
				log.Debugf("BatchNumber: %d, not found in trusted state. Storing it...", batch.BatchNumber)
				// If it is not found, store batch
				log.Infof("processSequenceBatches: (not found batch) ProcessAndStoreClosedBatch . BatchNumber: %d, BlockNumber: %d", processCtx.BatchNumber, blockNumber)
				newStateRoot, flushID, proverID, err := g.state.ProcessAndStoreClosedBatch(ctx, processCtx, batch.BatchL2Data, dbTx, stateMetrics.SynchronizerCallerLabel)
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

				newRoot = newStateRoot
				tBatch = &batch
				tBatch.StateRoot = newRoot
			} else {
				log.Error("error checking trusted state: ", err)
				rollbackErr := dbTx.Rollback(ctx)
				if rollbackErr != nil {
					log.Errorf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %v", batch.BatchNumber, blockNumber, rollbackErr)
					return rollbackErr
				}
				return err
			}
		} else {
			// Reprocess batch to compare the stateRoot with tBatch.StateRoot and get accInputHash
			p, err := g.state.ExecuteBatch(ctx, batch, false, dbTx)
			if err != nil {
				log.Errorf("error executing L1 batch: %+v, error: %v", batch, err)
				rollbackErr := dbTx.Rollback(ctx)
				if rollbackErr != nil {
					log.Errorf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %s, error : %v", batch.BatchNumber, blockNumber, rollbackErr.Error(), err)
					return rollbackErr
				}
				return err
			}
			newRoot = common.BytesToHash(p.NewStateRoot)
			accumulatedInputHash := common.BytesToHash(p.NewAccInputHash)

			//AddAccumulatedInputHash
			err = g.state.AddAccumulatedInputHash(ctx, batch.BatchNumber, accumulatedInputHash, dbTx)
			if err != nil {
				log.Errorf("error adding accumulatedInputHash for batch: %d. Error; %v", batch.BatchNumber, err)
				rollbackErr := dbTx.Rollback(ctx)
				if rollbackErr != nil {
					log.Errorf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %v", batch.BatchNumber, blockNumber, rollbackErr)
					return rollbackErr
				}
				return err
			}
		}

		// Call the check trusted state method to compare trusted and virtual state
		status := g.checkTrustedState(ctx, batch, tBatch, newRoot, dbTx)
		if status {
			// Reorg Pool
			err := g.reorgPool(ctx, dbTx)
			if err != nil {
				rollbackErr := dbTx.Rollback(ctx)
				if rollbackErr != nil {
					log.Errorf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %s, error : %v", tBatch.BatchNumber, blockNumber, rollbackErr.Error(), err)
					return rollbackErr
				}
				log.Errorf("error: %v. BatchNumber: %d, BlockNumber: %d", err, tBatch.BatchNumber, blockNumber)
				return err
			}

			// Clean trustedState sync variables to avoid sync the trusted state from the wrong starting point.
			// This wrong starting point would force the trusted sync to clean the virtualization of the batch reaching an inconsistency.
			g.sync.CleanTrustedState()

			// Reset trusted state
			previousBatchNumber := batch.BatchNumber - 1
			if tBatch.StateRoot == (common.Hash{}) {
				log.Warnf("cleaning state before inserting batch from L1. Clean until batch: %d", previousBatchNumber)
			} else {
				log.Warnf("missmatch in trusted state detected, discarding batches until batchNum %d", previousBatchNumber)
			}
			log.Infof("ResetTrustedState: Resetting trusted state. delete batch > %d, ", previousBatchNumber)
			err = g.state.ResetTrustedState(ctx, previousBatchNumber, dbTx) // This method has to reset the forced batches deleting the batchNumber for higher batchNumbers
			if err != nil {
				log.Errorf("error resetting trusted state. BatchNumber: %d, BlockNumber: %d, error: %v", batch.BatchNumber, blockNumber, err)
				rollbackErr := dbTx.Rollback(ctx)
				if rollbackErr != nil {
					log.Errorf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %s, error : %v", batch.BatchNumber, blockNumber, rollbackErr.Error(), err)
					return rollbackErr
				}
				log.Errorf("error resetting trusted state. BatchNumber: %d, BlockNumber: %d, error: %v", batch.BatchNumber, blockNumber, err)
				return err
			}
			log.Infof("processSequenceBatches: (deleted previous) ProcessAndStoreClosedBatch . BatchNumber: %d, BlockNumber: %d", processCtx.BatchNumber, blockNumber)
			_, flushID, proverID, err := g.state.ProcessAndStoreClosedBatch(ctx, processCtx, batch.BatchL2Data, dbTx, stateMetrics.SynchronizerCallerLabel)
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
		}

		// Store virtualBatch
		log.Infof("processSequenceBatches: Storing virtualBatch. BatchNumber: %d, BlockNumber: %d", virtualBatch.BatchNumber, blockNumber)
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
	}
	// Insert the sequence to allow the aggregator verify the sequence batches
	seq := state.Sequence{
		FromBatchNumber: sequencedBatches[0].BatchNumber,
		ToBatchNumber:   sequencedBatches[len(sequencedBatches)-1].BatchNumber,
	}
	err := g.state.AddSequence(ctx, seq, dbTx)
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

func (g *ProcessorL1SequenceBatches) reorgPool(ctx context.Context, dbTx pgx.Tx) error {
	latestBatchNum, err := g.etherMan.GetLatestBatchNumber()
	if err != nil {
		log.Error("error getting the latestBatchNumber virtualized in the smc. Error: ", err)
		return err
	}
	batchNumber := latestBatchNum + 1
	// Get transactions that have to be included in the pool again
	txs, err := g.state.GetReorgedTransactions(ctx, batchNumber, dbTx)
	if err != nil {
		log.Errorf("error getting txs from trusted state. BatchNumber: %d, error: %v", batchNumber, err)
		return err
	}
	log.Debug("Reorged transactions: ", txs)

	// Remove txs from the pool
	err = g.pool.DeleteReorgedTransactions(ctx, txs)
	if err != nil {
		log.Errorf("error deleting txs from the pool. BatchNumber: %d, error: %v", batchNumber, err)
		return err
	}
	log.Debug("Delete reorged transactions")

	// Add txs to the pool
	for _, tx := range txs {
		// Insert tx in WIP status to avoid the sequencer to grab them before it gets restarted
		// When the sequencer restarts, it will update the status to pending non-wip
		err = g.pool.StoreTx(ctx, *tx, "", true)
		if err != nil {
			log.Errorf("error storing tx into the pool again. TxHash: %s. BatchNumber: %d, error: %v", tx.Hash().String(), batchNumber, err)
			return err
		}
		log.Debug("Reorged transactions inserted in the pool: ", tx.Hash())
	}
	return nil
}

func (g *ProcessorL1SequenceBatches) checkTrustedState(ctx context.Context, batch state.Batch, tBatch *state.Batch, newRoot common.Hash, dbTx pgx.Tx) bool {
	//Compare virtual state with trusted state
	var reorgReasons strings.Builder
	if newRoot != tBatch.StateRoot {
		log.Warnf("Different field StateRoot. Virtual: %s, Trusted: %s\n", newRoot.String(), tBatch.StateRoot.String())
		reorgReasons.WriteString(fmt.Sprintf("Different field StateRoot. Virtual: %s, Trusted: %s\n", newRoot.String(), tBatch.StateRoot.String()))
	}
	if hex.EncodeToString(batch.BatchL2Data) != hex.EncodeToString(tBatch.BatchL2Data) {
		log.Warnf("Different field BatchL2Data. Virtual: %s, Trusted: %s\n", hex.EncodeToString(batch.BatchL2Data), hex.EncodeToString(tBatch.BatchL2Data))
		reorgReasons.WriteString(fmt.Sprintf("Different field BatchL2Data. Virtual: %s, Trusted: %s\n", hex.EncodeToString(batch.BatchL2Data), hex.EncodeToString(tBatch.BatchL2Data)))
	}
	if batch.GlobalExitRoot.String() != tBatch.GlobalExitRoot.String() {
		log.Warnf("Different field GlobalExitRoot. Virtual: %s, Trusted: %s\n", batch.GlobalExitRoot.String(), tBatch.GlobalExitRoot.String())
		reorgReasons.WriteString(fmt.Sprintf("Different field GlobalExitRoot. Virtual: %s, Trusted: %s\n", batch.GlobalExitRoot.String(), tBatch.GlobalExitRoot.String()))
	}
	if batch.Timestamp.Unix() != tBatch.Timestamp.Unix() {
		log.Warnf("Different field Timestamp. Virtual: %d, Trusted: %d\n", batch.Timestamp.Unix(), tBatch.Timestamp.Unix())
		reorgReasons.WriteString(fmt.Sprintf("Different field Timestamp. Virtual: %d, Trusted: %d\n", batch.Timestamp.Unix(), tBatch.Timestamp.Unix()))
	}
	if batch.Coinbase.String() != tBatch.Coinbase.String() {
		log.Warnf("Different field Coinbase. Virtual: %s, Trusted: %s\n", batch.Coinbase.String(), tBatch.Coinbase.String())
		reorgReasons.WriteString(fmt.Sprintf("Different field Coinbase. Virtual: %s, Trusted: %s\n", batch.Coinbase.String(), tBatch.Coinbase.String()))
	}

	if reorgReasons.Len() > 0 {
		reason := reorgReasons.String()

		if tBatch.StateRoot == (common.Hash{}) {
			log.Warnf("incomplete trusted batch %d detected. Syncing full batch from L1", tBatch.BatchNumber)
		} else {
			log.Warnf("missmatch in trusted state detected for Batch Number: %d. Reasons: %s", tBatch.BatchNumber, reason)
		}
		if g.sync.IsTrustedSequencer() {
			g.halt(ctx, fmt.Errorf("TRUSTED REORG DETECTED! Batch: %d", batch.BatchNumber))
		}
		// Store trusted reorg register
		tr := state.TrustedReorg{
			BatchNumber: tBatch.BatchNumber,
			Reason:      reason,
		}
		err := g.state.AddTrustedReorg(ctx, &tr, dbTx)
		if err != nil {
			log.Error("error storing tursted reorg register into the db. Error: ", err)
		}
		return true
	}
	return false
}

// halt halts the Synchronizer
func (g *ProcessorL1SequenceBatches) halt(ctx context.Context, err error) {
	event := &event.Event{
		ReceivedAt:  time.Now(),
		Source:      event.Source_Node,
		Component:   event.Component_Synchronizer,
		Level:       event.Level_Critical,
		EventID:     event.EventID_SynchronizerHalt,
		Description: fmt.Sprintf("Synchronizer halted due to error: %s", err),
	}

	eventErr := g.eventLog.LogEvent(ctx, event)
	if eventErr != nil {
		log.Errorf("error storing Synchronizer halt event: %v", eventErr)
	}

	for {
		log.Errorf("fatal error: %s", err)
		log.Error("halting the Synchronizer")
		time.Sleep(5 * time.Second) //nolint:gomnd
	}
}
