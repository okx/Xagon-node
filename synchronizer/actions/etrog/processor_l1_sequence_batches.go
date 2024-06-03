package etrog

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/metrics"
	stateMetrics "github.com/0xPolygonHermez/zkevm-node/state/metrics"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/actions"
	syncCommon "github.com/0xPolygonHermez/zkevm-node/synchronizer/common"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/common/syncinterfaces"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

type stateProcessSequenceBatches interface {
	GetNextForcedBatches(ctx context.Context, nextForcedBatches int, dbTx pgx.Tx) ([]state.ForcedBatch, error)
	GetBatchByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*state.Batch, error)
	ProcessAndStoreClosedBatchV2(ctx context.Context, processingCtx state.ProcessingContextV2, dbTx pgx.Tx, caller metrics.CallerLabel) (common.Hash, uint64, string, error)
	ExecuteBatchV2(ctx context.Context, batch state.Batch, L1InfoTreeRoot common.Hash, l1InfoTreeData map[uint32]state.L1DataV2, timestampLimit time.Time, updateMerkleTree bool, skipVerifyL1InfoRoot uint32, forcedBlockHashL1 *common.Hash, dbTx pgx.Tx) (*executor.ProcessBatchResponseV2, error)
	AddAccumulatedInputHash(ctx context.Context, batchNum uint64, accInputHash common.Hash, dbTx pgx.Tx) error
	ResetTrustedState(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) error
	AddSequence(ctx context.Context, sequence state.Sequence, dbTx pgx.Tx) error
	AddVirtualBatch(ctx context.Context, virtualBatch *state.VirtualBatch, dbTx pgx.Tx) error
	AddTrustedReorg(ctx context.Context, trustedReorg *state.TrustedReorg, dbTx pgx.Tx) error
	GetL1InfoTreeDataFromBatchL2Data(ctx context.Context, batchL2Data []byte, dbTx pgx.Tx) (map[uint32]state.L1DataV2, common.Hash, common.Hash, error)
	UpdateBatchTimestamp(ctx context.Context, batchNumber uint64, timestamp time.Time, dbTx pgx.Tx) error
}

type syncProcessSequenceBatchesInterface interface {
	PendingFlushID(flushID uint64, proverID string)
	IsTrustedSequencer() bool
	CleanTrustedState()
}

// ProcessorL1SequenceBatchesEtrog implements L1EventProcessor
type ProcessorL1SequenceBatchesEtrog struct {
	actions.ProcessorBase[ProcessorL1SequenceBatchesEtrog]
	state        stateProcessSequenceBatches
	sync         syncProcessSequenceBatchesInterface
	timeProvider syncCommon.TimeProvider
	halter       syncinterfaces.CriticalErrorHandler
}

// NewProcessorL1SequenceBatches returns instance of a processor for SequenceBatchesOrder
func NewProcessorL1SequenceBatches(state stateProcessSequenceBatches,
	sync syncProcessSequenceBatchesInterface,
	timeProvider syncCommon.TimeProvider,
	halter syncinterfaces.CriticalErrorHandler) *ProcessorL1SequenceBatchesEtrog {
	return &ProcessorL1SequenceBatchesEtrog{
		ProcessorBase: actions.ProcessorBase[ProcessorL1SequenceBatchesEtrog]{
			SupportedEvent:    []etherman.EventOrder{etherman.SequenceBatchesOrder, etherman.InitialSequenceBatchesOrder},
			SupportedForkdIds: &actions.ForksIdOnlyEtrog},
		state:        state,
		sync:         sync,
		timeProvider: timeProvider,
		halter:       halter,
	}
}

// Process process event
func (g *ProcessorL1SequenceBatchesEtrog) Process(ctx context.Context, order etherman.Order, l1Block *etherman.Block, dbTx pgx.Tx) error {
	if l1Block == nil || len(l1Block.SequencedBatches) <= order.Pos {
		return actions.ErrInvalidParams
	}
	err := g.ProcessSequenceBatches(ctx, l1Block.SequencedBatches[order.Pos], l1Block.BlockNumber, l1Block.ReceivedAt, dbTx)
	return err
}

// ProcessSequenceBatches process sequence of batches
func (p *ProcessorL1SequenceBatchesEtrog) ProcessSequenceBatches(ctx context.Context, sequencedBatches []etherman.SequencedBatch, blockNumber uint64, l1BlockTimestamp time.Time, dbTx pgx.Tx) error {
	if len(sequencedBatches) == 0 {
		log.Warn("Empty sequencedBatches array detected, ignoring...")
		return nil
	}
	now := p.timeProvider.Now()
	for _, sbatch := range sequencedBatches {
		virtualBatch := state.VirtualBatch{
			BatchNumber:         sbatch.BatchNumber,
			TxHash:              sbatch.TxHash,
			Coinbase:            sbatch.Coinbase,
			BlockNumber:         blockNumber,
			SequencerAddr:       sbatch.SequencerAddr,
			TimestampBatchEtrog: &l1BlockTimestamp,
		}
		batch := state.Batch{
			BatchNumber: sbatch.BatchNumber,
			// This timestamp now is the timeLimit. It can't be the one virtual.BatchTimestamp
			//   because when sync from trusted we don't now the real BatchTimestamp and
			//   will fails the comparation of batch time >= than previous one.
			Timestamp:   now,
			Coinbase:    sbatch.Coinbase,
			BatchL2Data: sbatch.PolygonRollupBaseEtrogBatchData.Transactions,
		}
		var (
			processCtx        state.ProcessingContextV2
			forcedBlockHashL1 *common.Hash
			err               error
		)
		leaves := make(map[uint32]state.L1DataV2)

		// ForcedBatch must be processed
		if sbatch.PolygonRollupBaseEtrogBatchData.ForcedTimestamp > 0 && sbatch.BatchNumber != 1 { // If this is true means that the batch is forced
			log.Debug("FORCED BATCH SEQUENCED!")
			// Read forcedBatches from db
			forcedBatches, err := p.state.GetNextForcedBatches(ctx, 1, dbTx)
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
			if uint64(forcedBatches[0].ForcedAt.Unix()) != sbatch.PolygonRollupBaseEtrogBatchData.ForcedTimestamp ||
				forcedBatches[0].GlobalExitRoot != sbatch.PolygonRollupBaseEtrogBatchData.ForcedGlobalExitRoot ||
				common.Bytes2Hex(forcedBatches[0].RawTxsData) != common.Bytes2Hex(sbatch.PolygonRollupBaseEtrogBatchData.Transactions) {
				log.Warnf("ForcedBatch stored: %+v. RawTxsData: %s", forcedBatches, common.Bytes2Hex(forcedBatches[0].RawTxsData))
				log.Warnf("ForcedBatch sequenced received: %+v. RawTxsData: %s", sbatch, common.Bytes2Hex(sbatch.PolygonRollupBaseEtrogBatchData.Transactions))
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
			batch.GlobalExitRoot = sbatch.PolygonRollupBaseEtrogBatchData.ForcedGlobalExitRoot
			tstampLimit := forcedBatches[0].ForcedAt
			txs := forcedBatches[0].RawTxsData
			// The leaves are no needed for forced batches
			processCtx = state.ProcessingContextV2{
				BatchNumber:          sbatch.BatchNumber,
				Coinbase:             sbatch.SequencerAddr,
				Timestamp:            &tstampLimit,
				L1InfoRoot:           sbatch.PolygonRollupBaseEtrogBatchData.ForcedGlobalExitRoot,
				BatchL2Data:          &txs,
				ForcedBlockHashL1:    forcedBlockHashL1,
				SkipVerifyL1InfoRoot: 1,
				ClosingReason:        state.SyncL1EventSequencedForcedBatchClosingReason,
			}
		} else if sbatch.PolygonRollupBaseEtrogBatchData.ForcedTimestamp > 0 && sbatch.BatchNumber == 1 { // This is the initial batch (injected)
			log.Debug("Processing initial batch")
			batch.GlobalExitRoot = sbatch.PolygonRollupBaseEtrogBatchData.ForcedGlobalExitRoot
			var fBHL1 common.Hash = sbatch.PolygonRollupBaseEtrogBatchData.ForcedBlockHashL1
			forcedBlockHashL1 = &fBHL1
			txs := sbatch.PolygonRollupBaseEtrogBatchData.Transactions
			tstampLimit := time.Unix(int64(sbatch.PolygonRollupBaseEtrogBatchData.ForcedTimestamp), 0)
			processCtx = state.ProcessingContextV2{
				BatchNumber:          1,
				Coinbase:             sbatch.SequencerAddr,
				Timestamp:            &tstampLimit,
				L1InfoRoot:           sbatch.PolygonRollupBaseEtrogBatchData.ForcedGlobalExitRoot,
				BatchL2Data:          &txs,
				ForcedBlockHashL1:    forcedBlockHashL1,
				SkipVerifyL1InfoRoot: 1,
				ClosingReason:        state.SyncL1EventInitialBatchClosingReason,
			}
		} else {
			var maxGER common.Hash
			leaves, _, maxGER, err = p.state.GetL1InfoTreeDataFromBatchL2Data(ctx, batch.BatchL2Data, dbTx)
			if err != nil {
				log.Errorf("error getting L1InfoRootLeafByL1InfoRoot. sbatch.L1InfoRoot: %v", *sbatch.L1InfoRoot)
				rollbackErr := dbTx.Rollback(ctx)
				if rollbackErr != nil {
					log.Errorf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %s, error : %v", virtualBatch.BatchNumber, blockNumber, rollbackErr.Error(), err)
					return rollbackErr
				}
				return err
			}
			batch.GlobalExitRoot = maxGER
			processCtx = state.ProcessingContextV2{
				BatchNumber:          batch.BatchNumber,
				Coinbase:             batch.Coinbase,
				Timestamp:            &l1BlockTimestamp,
				L1InfoRoot:           *sbatch.L1InfoRoot,
				L1InfoTreeData:       leaves,
				ForcedBatchNum:       batch.ForcedBatchNum,
				BatchL2Data:          &batch.BatchL2Data,
				SkipVerifyL1InfoRoot: 1,
				GlobalExitRoot:       batch.GlobalExitRoot,
				ClosingReason:        state.SyncL1EventSequencedBatchClosingReason,
			}
			if batch.GlobalExitRoot == (common.Hash{}) {
				if len(leaves) > 0 {
					globalExitRoot := leaves[uint32(len(leaves)-1)].GlobalExitRoot
					log.Debugf("Empty GER detected for batch: %d usign GER of last leaf (%d):%s",
						batch.BatchNumber,
						uint32(len(leaves)-1),
						globalExitRoot)

					processCtx.GlobalExitRoot = globalExitRoot
					batch.GlobalExitRoot = globalExitRoot
				} else {
					log.Debugf("Empty leaves array detected for batch: %d usign GER:%s", batch.BatchNumber, processCtx.GlobalExitRoot.String())
				}
			}
		}
		virtualBatch.L1InfoRoot = &processCtx.L1InfoRoot
		var newRoot common.Hash

		// First get trusted batch from db
		tBatch, err := p.state.GetBatchByNumber(ctx, batch.BatchNumber, dbTx)
		if err != nil {
			if errors.Is(err, state.ErrNotFound) {
				log.Debugf("BatchNumber: %d, not found in trusted state. Storing it...", batch.BatchNumber)
				// If it is not found, store batch
				log.Infof("processSequenceBatches: (not found batch) ProcessAndStoreClosedBatch . BatchNumber: %d, BlockNumber: %d GER:%s", processCtx.BatchNumber, blockNumber, processCtx.GlobalExitRoot.String())
				newStateRoot, flushID, proverID, err := p.state.ProcessAndStoreClosedBatchV2(ctx, processCtx, dbTx, stateMetrics.SynchronizerCallerLabel)
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
				p.sync.PendingFlushID(flushID, proverID)

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
			// Batch already exists
			// We update the timestamp of the batch to match the timestamp
			err := p.state.UpdateBatchTimestamp(ctx, batch.BatchNumber, *processCtx.Timestamp, dbTx)
			if err != nil {
				log.Errorf("error updating batch timestamp %s. BatchNumber: %d, BlockNumber: %d, error: %v", processCtx.Timestamp, batch.BatchNumber, blockNumber, err)
				rollbackErr := dbTx.Rollback(ctx)
				if rollbackErr != nil {
					log.Errorf("error rolling back state because error updating batch timestamp. BatchNumber: %d, BlockNumber: %d, rollbackErr: %s, error : %v", batch.BatchNumber, blockNumber, rollbackErr.Error(), err)
					return rollbackErr
				}
			}
			// Reprocess batch to compare the stateRoot with tBatch.StateRoot and get accInputHash
			batchRespose, err := p.state.ExecuteBatchV2(ctx, batch, processCtx.L1InfoRoot, leaves, *processCtx.Timestamp, false, processCtx.SkipVerifyL1InfoRoot, processCtx.ForcedBlockHashL1, dbTx)
			if err != nil {
				log.Errorf("error executing L1 batch: %+v, error: %v", batch, err)
				rollbackErr := dbTx.Rollback(ctx)
				if rollbackErr != nil {
					log.Errorf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %s, error : %v", batch.BatchNumber, blockNumber, rollbackErr.Error(), err)
					return rollbackErr
				}
				return err
			}
			newRoot = common.BytesToHash(batchRespose.NewStateRoot)
			accumulatedInputHash := common.BytesToHash(batchRespose.NewAccInputHash)

			//AddAccumulatedInputHash
			err = p.state.AddAccumulatedInputHash(ctx, batch.BatchNumber, accumulatedInputHash, dbTx)
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
		status := p.checkTrustedState(ctx, batch, tBatch, newRoot, dbTx)
		if status {
			// Clean trustedState sync variables to avoid sync the trusted state from the wrong starting point.
			// This wrong starting point would force the trusted sync to clean the virtualization of the batch reaching an inconsistency.
			p.sync.CleanTrustedState()

			// Reset trusted state
			previousBatchNumber := batch.BatchNumber - 1
			if tBatch.WIP {
				log.Infof("cleaning state before inserting batch from L1. Clean until batch: %d", previousBatchNumber)
			} else {
				log.Warnf("missmatch in trusted state detected, discarding batches until batchNum %d", previousBatchNumber)
			}
			log.Infof("ResetTrustedState: Resetting trusted state. delete batch > %d, ", previousBatchNumber)
			err = p.state.ResetTrustedState(ctx, previousBatchNumber, dbTx) // This method has to reset the forced batches deleting the batchNumber for higher batchNumbers
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
			_, flushID, proverID, err := p.state.ProcessAndStoreClosedBatchV2(ctx, processCtx, dbTx, stateMetrics.SynchronizerCallerLabel)
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
			p.sync.PendingFlushID(flushID, proverID)
		}

		// Store virtualBatch
		log.Infof("processSequenceBatches: Storing virtualBatch. BatchNumber: %d, BlockNumber: %d", virtualBatch.BatchNumber, blockNumber)
		err = p.state.AddVirtualBatch(ctx, &virtualBatch, dbTx)
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
	err := p.state.AddSequence(ctx, seq, dbTx)
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

func (p *ProcessorL1SequenceBatchesEtrog) checkTrustedState(ctx context.Context, batch state.Batch, tBatch *state.Batch, newRoot common.Hash, dbTx pgx.Tx) bool {
	//Compare virtual state with trusted state
	var reorgReasons strings.Builder
	batchNumStr := fmt.Sprintf("Batch: %d.", batch.BatchNumber)
	if newRoot != tBatch.StateRoot {
		errMsg := batchNumStr + fmt.Sprintf("Different field StateRoot. Virtual: %s, Trusted: %s\n", newRoot.String(), tBatch.StateRoot.String())
		log.Warnf(errMsg)
		reorgReasons.WriteString(errMsg)
	}
	if hex.EncodeToString(batch.BatchL2Data) != hex.EncodeToString(tBatch.BatchL2Data) {
		errMsg := batchNumStr + fmt.Sprintf("Different field BatchL2Data. Virtual: %s, Trusted: %s\n", hex.EncodeToString(batch.BatchL2Data), hex.EncodeToString(tBatch.BatchL2Data))
		log.Warnf(errMsg)
		reorgReasons.WriteString(errMsg)
	}
	if batch.GlobalExitRoot.String() != tBatch.GlobalExitRoot.String() {
		errMsg := batchNumStr + fmt.Sprintf("Different field GlobalExitRoot. Virtual: %s, Trusted: %s\n", batch.GlobalExitRoot.String(), tBatch.GlobalExitRoot.String())
		log.Warnf(errMsg)
		reorgReasons.WriteString(fmt.Sprintf("Different field GlobalExitRoot. Virtual: %s, Trusted: %s\n", batch.GlobalExitRoot.String(), tBatch.GlobalExitRoot.String()))
	}
	if batch.Timestamp.Unix() < tBatch.Timestamp.Unix() { // TODO: this timestamp will be different in permissionless nodes and the trusted node
		errMsg := batchNumStr + fmt.Sprintf("Invalid timestamp. Virtual timestamp limit(%d) must be greater or equal than Trusted timestamp (%d)\n", batch.Timestamp.Unix(), tBatch.Timestamp.Unix())
		log.Warnf(errMsg)
		reorgReasons.WriteString(errMsg)
	}
	if batch.Coinbase.String() != tBatch.Coinbase.String() {
		errMsg := batchNumStr + fmt.Sprintf("Different field Coinbase. Virtual: %s, Trusted: %s\n", batch.Coinbase.String(), tBatch.Coinbase.String())
		log.Warnf(errMsg)
		reorgReasons.WriteString(errMsg)
	}
	if tBatch.WIP {
		errMsg := batchNumStr + "Trusted batch is WIP\n"
		log.Warnf(errMsg)
		reorgReasons.WriteString(errMsg)
	}

	if reorgReasons.Len() > 0 {
		reason := reorgReasons.String()

		if p.sync.IsTrustedSequencer() {
			log.Errorf("TRUSTED REORG DETECTED! Batch: %d reason:%s", batch.BatchNumber, reason)
			// Halt function never have to return! it must blocks the process
			p.halt(ctx, fmt.Errorf("TRUSTED REORG DETECTED! Batch: %d", batch.BatchNumber))
			log.Errorf("CRITICAL!!!: Never have to execute this code. Halt function never have to return! it must blocks the process")
		}
		if !tBatch.WIP {
			log.Warnf("missmatch in trusted state detected for Batch Number: %d. Reasons: %s", tBatch.BatchNumber, reason)
			// Store trusted reorg register
			tr := state.TrustedReorg{
				BatchNumber: tBatch.BatchNumber,
				Reason:      reason,
			}
			err := p.state.AddTrustedReorg(ctx, &tr, dbTx)
			if err != nil {
				log.Error("error storing trusted reorg register into the db. Error: ", err)
			}
		} else {
			log.Warnf("incomplete trusted batch %d detected. Syncing full batch from L1", tBatch.BatchNumber)
		}
		return true
	}
	return false
}

// halt halts the Synchronizer
func (p *ProcessorL1SequenceBatchesEtrog) halt(ctx context.Context, err error) {
	p.halter.CriticalError(ctx, err)
}
