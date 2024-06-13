package sequencer

import (
	"context"
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	stateMetrics "github.com/0xPolygonHermez/zkevm-node/state/metrics"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

// processForcedBatches processes all the forced batches that are pending to be processed
func (f *finalizer) processForcedBatches(ctx context.Context, lastBatchNumber uint64, stateRoot common.Hash) (newLastBatchNumber uint64, newStateRoot common.Hash) {
	f.nextForcedBatchesMux.Lock()
	defer f.nextForcedBatchesMux.Unlock()
	f.nextForcedBatchDeadline = 0

	lastForcedBatchNumber, err := f.stateIntf.GetLastTrustedForcedBatchNumber(ctx, nil)
	if err != nil {
		log.Errorf("failed to get last trusted forced batch number, error: %v", err)
		return lastBatchNumber, stateRoot
	}
	nextForcedBatchNumber := lastForcedBatchNumber + 1

	for _, forcedBatch := range f.nextForcedBatches {
		forcedBatchToProcess := forcedBatch
		// Skip already processed forced batches
		if forcedBatchToProcess.ForcedBatchNumber < nextForcedBatchNumber {
			continue
		} else if forcedBatch.ForcedBatchNumber > nextForcedBatchNumber {
			// We have a gap in the f.nextForcedBatches slice, we get the missing forced batch from the state
			missingForcedBatch, err := f.stateIntf.GetForcedBatch(ctx, nextForcedBatchNumber, nil)
			if err != nil {
				log.Errorf("failed to get missing forced batch %d, error: %v", nextForcedBatchNumber, err)
				return lastBatchNumber, stateRoot
			}
			forcedBatchToProcess = *missingForcedBatch
		}

		var contextId string
		log.Infof("processing forced batch %d, lastBatchNumber: %d, stateRoot: %s", forcedBatchToProcess.ForcedBatchNumber, lastBatchNumber, stateRoot.String())
		lastBatchNumber, stateRoot, contextId, err = f.processForcedBatch(ctx, forcedBatchToProcess, lastBatchNumber, stateRoot)

		if err != nil {
			log.Errorf("error when processing forced batch %d, error: %v", forcedBatchToProcess.ForcedBatchNumber, err)
			return lastBatchNumber, stateRoot
		}

		log.Infof("processed forced batch %d, batchNumber: %d, newStateRoot: %s, contextId: %s", forcedBatchToProcess.ForcedBatchNumber, lastBatchNumber, stateRoot.String(), contextId)

		nextForcedBatchNumber += 1
	}
	f.nextForcedBatches = make([]state.ForcedBatch, 0)

	return lastBatchNumber, stateRoot
}

func (f *finalizer) processForcedBatch(ctx context.Context, forcedBatch state.ForcedBatch, lastBatchNumber uint64, stateRoot common.Hash) (newLastBatchNumber uint64, newStateRoot common.Hash, ctxId string, retErr error) {
	dbTx, err := f.stateIntf.BeginStateTransaction(ctx)
	if err != nil {
		log.Errorf("failed to begin state transaction for process forced batch %d, error: %v", forcedBatch.ForcedBatchNumber, err)
		return lastBatchNumber, stateRoot, "", err
	}

	// Helper function in case we get an error when processing the forced batch
	rollbackOnError := func(retError error) (newLastBatchNumber uint64, newStateRoot common.Hash, ctxId string, retErr error) {
		err := dbTx.Rollback(ctx)
		if err != nil {
			return lastBatchNumber, stateRoot, "", fmt.Errorf("rollback error due to error %v, error: %v", retError, err)
		}
		return lastBatchNumber, stateRoot, "", retError
	}

	// Get L1 block for the forced batch
	fbL1Block, err := f.stateIntf.GetBlockByNumber(ctx, forcedBatch.BlockNumber, dbTx)
	if err != nil {
		return lastBatchNumber, stateRoot, "", fmt.Errorf("error getting L1 block number %d for forced batch %d, error: %v", forcedBatch.ForcedBatchNumber, forcedBatch.ForcedBatchNumber, err)
	}

	newBatchNumber := lastBatchNumber + 1

	// Open new batch on state for the forced batch
	processingCtx := state.ProcessingContext{
		BatchNumber:    newBatchNumber,
		Coinbase:       f.l2Coinbase,
		Timestamp:      time.Now(),
		GlobalExitRoot: forcedBatch.GlobalExitRoot,
		ForcedBatchNum: &forcedBatch.ForcedBatchNumber,
	}
	err = f.stateIntf.OpenBatch(ctx, processingCtx, dbTx)
	if err != nil {
		return rollbackOnError(fmt.Errorf("error opening state batch %d for forced batch %d, error: %v", newBatchNumber, forcedBatch.ForcedBatchNumber, err))
	}

	batchRequest := state.ProcessRequest{
		BatchNumber:             newBatchNumber,
		L1InfoRoot_V2:           forcedBatch.GlobalExitRoot,
		ForcedBlockHashL1:       fbL1Block.ParentHash,
		OldStateRoot:            stateRoot,
		Transactions:            forcedBatch.RawTxsData,
		Coinbase:                f.l2Coinbase,
		TimestampLimit_V2:       uint64(forcedBatch.ForcedAt.Unix()),
		ForkID:                  f.stateIntf.GetForkIDByBatchNumber(lastBatchNumber),
		SkipVerifyL1InfoRoot_V2: true,
		Caller:                  stateMetrics.DiscardCallerLabel,
	}

	batchResponse, contextId, err := f.stateIntf.ProcessBatchV2(ctx, batchRequest, true)
	if err != nil {
		return rollbackOnError(fmt.Errorf("failed to process/execute forced batch %d, error: %v", forcedBatch.ForcedBatchNumber, err))
	}

	// Close state batch
	processingReceipt := state.ProcessingReceipt{
		BatchNumber:   newBatchNumber,
		StateRoot:     batchResponse.NewStateRoot,
		LocalExitRoot: batchResponse.NewLocalExitRoot,
		BatchL2Data:   forcedBatch.RawTxsData,
		BatchResources: state.BatchResources{
			ZKCounters: batchResponse.UsedZkCounters,
			Bytes:      uint64(len(forcedBatch.RawTxsData)),
		},
		ClosingReason: state.ForcedBatchClosingReason,
	}
	err = f.stateIntf.CloseBatch(ctx, processingReceipt, dbTx)
	if err != nil {
		return rollbackOnError(fmt.Errorf("error closing state batch %d for forced batch %d, error: %v", newBatchNumber, forcedBatch.ForcedBatchNumber, err))
	}

	if len(batchResponse.BlockResponses) > 0 && !batchResponse.IsRomOOCError {
		err = f.handleProcessForcedBatchResponse(ctx, newBatchNumber, batchResponse, dbTx)
		if err != nil {
			return rollbackOnError(fmt.Errorf("error when handling batch response for forced batch %d, error: %v", forcedBatch.ForcedBatchNumber, err))
		}
	}

	err = dbTx.Commit(ctx)
	if err != nil {
		return rollbackOnError(fmt.Errorf("error when commit dbTx when processing forced batch %d, error: %v", forcedBatch.ForcedBatchNumber, err))
	}

	return newBatchNumber, batchResponse.NewStateRoot, contextId, nil
}

// addForcedTxToWorker adds the txs of the forced batch to the worker
func (f *finalizer) addForcedTxToWorker(forcedBatchResponse *state.ProcessBatchResponse) {
	for _, blockResponse := range forcedBatchResponse.BlockResponses {
		for _, txResponse := range blockResponse.TransactionResponses {
			from, err := state.GetSender(txResponse.Tx)
			if err != nil {
				log.Warnf("failed to get sender for tx %s, error: %v", txResponse.TxHash, err)
				continue
			}
			f.workerIntf.AddForcedTx(txResponse.TxHash, from)
		}
	}
}

// handleProcessForcedTxsResponse handles the block/transactions responses for the processed forced batch.
func (f *finalizer) handleProcessForcedBatchResponse(ctx context.Context, newBatchNumber uint64, batchResponse *state.ProcessBatchResponse, dbTx pgx.Tx) error {
	f.addForcedTxToWorker(batchResponse)

	f.updateFlushIDs(batchResponse.FlushID, batchResponse.StoredFlushID)

	// Wait until forced batch has been flushed/stored by the executor
	f.storedFlushIDCond.L.Lock()
	for f.storedFlushID < batchResponse.FlushID {
		f.storedFlushIDCond.Wait()
		// check if context is done after waking up
		if ctx.Err() != nil {
			f.storedFlushIDCond.L.Unlock()
			return nil
		}
	}
	f.storedFlushIDCond.L.Unlock()

	// process L2 blocks responses for the forced batch
	for _, forcedL2BlockResponse := range batchResponse.BlockResponses {
		// Store forced L2 blocks in the state
		blockHash, err := f.stateIntf.StoreL2Block(ctx, newBatchNumber, forcedL2BlockResponse, nil, dbTx)
		if err != nil {
			return fmt.Errorf("database error on storing L2 block %d, error: %v", forcedL2BlockResponse.BlockNumber, err)
		}

		// Update worker with info from the transaction responses
		for _, txResponse := range forcedL2BlockResponse.TransactionResponses {
			from, err := state.GetSender(txResponse.Tx)
			if err != nil {
				log.Warnf("failed to get sender for tx %s, error: %v", txResponse.TxHash, err)
			}

			if err == nil {
				f.updateWorkerAfterSuccessfulProcessing(ctx, txResponse.TxHash, from, true, batchResponse)
			}
		}

		// Send L2 block to data streamer
		err = f.DSSendL2Block(ctx, newBatchNumber, forcedL2BlockResponse, 0, forcedL2BlockResponse.Timestamp, blockHash)
		if err != nil {
			//TODO: we need to halt/rollback the L2 block if we had an error sending to the data streamer?
			log.Errorf("error sending L2 block %d to data streamer, error: %v", forcedL2BlockResponse.BlockNumber, err)
		}
	}

	return nil
}

// sortForcedBatches sorts the forced batches by ForcedBatchNumber
func (f *finalizer) sortForcedBatches(fb []state.ForcedBatch) []state.ForcedBatch {
	if len(fb) == 0 {
		return fb
	}
	// Sort by ForcedBatchNumber
	for i := 0; i < len(fb)-1; i++ {
		for j := i + 1; j < len(fb); j++ {
			if fb[i].ForcedBatchNumber > fb[j].ForcedBatchNumber {
				fb[i], fb[j] = fb[j], fb[i]
			}
		}
	}

	return fb
}

// setNextForcedBatchDeadline sets the next forced batch deadline
func (f *finalizer) setNextForcedBatchDeadline() {
	f.nextForcedBatchDeadline = now().Unix() + int64(f.cfg.ForcedBatchesTimeout.Duration.Seconds())
}

func (f *finalizer) checkForcedBatches(ctx context.Context) {
	for {
		time.Sleep(f.cfg.ForcedBatchesCheckInterval.Duration)

		if f.lastForcedBatchNum == 0 {
			lastTrustedForcedBatchNum, err := f.stateIntf.GetLastTrustedForcedBatchNumber(ctx, nil)
			if err != nil {
				log.Errorf("error getting last trusted forced batch number, error: %v", err)
				continue
			}
			if lastTrustedForcedBatchNum > 0 {
				f.lastForcedBatchNum = lastTrustedForcedBatchNum
			}
		}
		// Take into account L1 finality
		lastBlock, err := f.stateIntf.GetLastBlock(ctx, nil)
		if err != nil {
			log.Errorf("failed to get latest L1 block number, error: %v", err)
			continue
		}

		blockNumber := lastBlock.BlockNumber

		maxBlockNumber := uint64(0)
		finalityNumberOfBlocks := f.cfg.ForcedBatchesL1BlockConfirmations

		if finalityNumberOfBlocks <= blockNumber {
			maxBlockNumber = blockNumber - finalityNumberOfBlocks
		}

		forcedBatches, err := f.stateIntf.GetForcedBatchesSince(ctx, f.lastForcedBatchNum, maxBlockNumber, nil)
		if err != nil {
			log.Errorf("error checking forced batches, error: %v", err)
			continue
		}

		for _, forcedBatch := range forcedBatches {
			log.Debugf("finalizer received forced batch at block number: %d", forcedBatch.BlockNumber)

			f.nextForcedBatchesMux.Lock()
			f.nextForcedBatches = f.sortForcedBatches(append(f.nextForcedBatches, *forcedBatch))
			if f.nextForcedBatchDeadline == 0 {
				f.setNextForcedBatchDeadline()
			}
			f.nextForcedBatchesMux.Unlock()

			f.lastForcedBatchNum = forcedBatch.ForcedBatchNumber
		}
	}
}
