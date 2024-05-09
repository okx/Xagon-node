package sequencer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/event"
	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	stateMetrics "github.com/0xPolygonHermez/zkevm-node/state/metrics"
	"github.com/ethereum/go-ethereum/common"
)

// Batch represents a wip or processed batch.
type Batch struct {
	batchNumber             uint64
	coinbase                common.Address
	timestamp               time.Time
	initialStateRoot        common.Hash // initial stateRoot of the batch
	imStateRoot             common.Hash // intermediate stateRoot when processing tx-by-tx
	finalStateRoot          common.Hash // final stateroot of the batch when a L2 block is processed
	countOfTxs              int
	countOfL2Blocks         int
	imRemainingResources    state.BatchResources // remaining batch resources when processing tx-by-tx
	finalRemainingResources state.BatchResources // remaining batch resources when a L2 block is processed
	closingReason           state.ClosingReason
}

func (w *Batch) isEmpty() bool {
	return w.countOfL2Blocks == 0
}

// processBatchesPendingtoCheck performs a sanity check for batches closed but pending to be checked
func (f *finalizer) processBatchesPendingtoCheck(ctx context.Context) {
	notCheckedBatches, err := f.stateIntf.GetNotCheckedBatches(ctx, nil)
	if err != nil && err != state.ErrNotFound {
		log.Fatalf("failed to get batches not checked, error: ", err)
	}

	if len(notCheckedBatches) == 0 {
		return
	}

	log.Infof("executing sanity check for not checked batches")

	prevBatchNumber := notCheckedBatches[0].BatchNumber - 1
	prevBatch, err := f.stateIntf.GetBatchByNumber(ctx, prevBatchNumber, nil)
	if err != nil {
		log.Fatalf("failed to get batch %d, error: ", prevBatchNumber, err)
	}
	oldStateRoot := prevBatch.StateRoot

	for _, notCheckedBatch := range notCheckedBatches {
		_, _ = f.batchSanityCheck(ctx, notCheckedBatch.BatchNumber, oldStateRoot, notCheckedBatch.StateRoot)
		oldStateRoot = notCheckedBatch.StateRoot
	}
}

// setWIPBatch sets finalizer wip batch to the state batch passed as parameter
func (f *finalizer) setWIPBatch(ctx context.Context, wipStateBatch *state.Batch) (*Batch, error) {
	// Retrieve prevStateBatch to init the initialStateRoot of the wip batch
	prevStateBatch, err := f.stateIntf.GetBatchByNumber(ctx, wipStateBatch.BatchNumber-1, nil)
	if err != nil {
		return nil, err
	}

	wipStateBatchBlocks, err := state.DecodeBatchV2(wipStateBatch.BatchL2Data)
	if err != nil {
		return nil, err
	}

	// Count the number of txs in the wip state batch
	wipStateBatchCountOfTxs := 0
	for _, rawBlock := range wipStateBatchBlocks.Blocks {
		wipStateBatchCountOfTxs = wipStateBatchCountOfTxs + len(rawBlock.Transactions)
	}

	remainingResources := getMaxRemainingResources(f.batchConstraints)
	overflow, overflowResource := remainingResources.Sub(wipStateBatch.Resources)
	if overflow {
		return nil, fmt.Errorf("failed to subtract used resources when setting the WIP batch to the state batch %d, overflow resource: %s", wipStateBatch.BatchNumber, overflowResource)
	}

	wipBatch := &Batch{
		batchNumber:             wipStateBatch.BatchNumber,
		coinbase:                wipStateBatch.Coinbase,
		imStateRoot:             wipStateBatch.StateRoot,
		initialStateRoot:        prevStateBatch.StateRoot,
		finalStateRoot:          wipStateBatch.StateRoot,
		timestamp:               wipStateBatch.Timestamp,
		countOfL2Blocks:         len(wipStateBatchBlocks.Blocks),
		countOfTxs:              wipStateBatchCountOfTxs,
		imRemainingResources:    remainingResources,
		finalRemainingResources: remainingResources,
	}

	return wipBatch, nil
}

// initWIPBatch inits the wip batch
func (f *finalizer) initWIPBatch(ctx context.Context) {
	for !f.isSynced(ctx) {
		log.Info("wait for synchronizer to sync last batch")
		time.Sleep(time.Second)
	}

	lastBatchNum, err := f.stateIntf.GetLastBatchNumber(ctx, nil)
	if err != nil {
		log.Fatalf("failed to get last batch number, error: %v", err)
	}

	// Get the last batch in trusted state
	lastStateBatch, err := f.stateIntf.GetBatchByNumber(ctx, lastBatchNum, nil)
	if err != nil {
		log.Fatalf("failed to get last batch %d, error: %v", lastBatchNum, err)
	}

	isClosed := !lastStateBatch.WIP

	log.Infof("batch %d isClosed: %v", lastBatchNum, isClosed)

	if isClosed { //if the last batch is close then open a new wip batch
		if lastStateBatch.BatchNumber+1 == f.cfg.HaltOnBatchNumber {
			f.Halt(ctx, fmt.Errorf("finalizer reached stop sequencer on batch number: %d", f.cfg.HaltOnBatchNumber), false)
		}

		f.wipBatch, err = f.openNewWIPBatch(ctx, lastStateBatch.BatchNumber+1, lastStateBatch.StateRoot)
		if err != nil {
			log.Fatalf("failed to open new wip batch, error: %v", err)
		}
	} else { /// if it's not closed, it is the wip state batch, set it as wip batch in the finalizer
		f.wipBatch, err = f.setWIPBatch(ctx, lastStateBatch)
		if err != nil {
			log.Fatalf("failed to set wip batch, error: %v", err)
		}
	}

	log.Infof("initial batch: %d, initialStateRoot: %s, stateRoot: %s, coinbase: %s",
		f.wipBatch.batchNumber, f.wipBatch.initialStateRoot, f.wipBatch.finalStateRoot, f.wipBatch.coinbase)
}

// finalizeWIPBatch closes the current batch and opens a new one, potentially processing forced batches between the batch is closed and the resulting new empty batch
func (f *finalizer) finalizeWIPBatch(ctx context.Context, closeReason state.ClosingReason) {
	prevTimestamp := f.wipL2Block.timestamp
	prevL1InfoTreeIndex := f.wipL2Block.l1InfoTreeExitRoot.L1InfoTreeIndex

	// Close the wip L2 block if it has transactions, otherwise we keep the wip L2 block to store it in the new wip batch
	if !f.wipL2Block.isEmpty() {
		f.closeWIPL2Block(ctx)
	}

	err := f.closeAndOpenNewWIPBatch(ctx, closeReason)
	if err != nil {
		f.Halt(ctx, fmt.Errorf("failed to create new WIP batch, error: %v", err), true)
	}

	// If we have closed the wipL2Block then we open a new one
	if f.wipL2Block == nil {
		f.openNewWIPL2Block(ctx, prevTimestamp, &prevL1InfoTreeIndex)
	}
}

// closeAndOpenNewWIPBatch closes the current batch and opens a new one, potentially processing forced batches between the batch is closed and the resulting new wip batch
func (f *finalizer) closeAndOpenNewWIPBatch(ctx context.Context, closeReason state.ClosingReason) error {
	f.nextForcedBatchesMux.Lock()
	processForcedBatches := len(f.nextForcedBatches) > 0
	f.nextForcedBatchesMux.Unlock()

	// If we will process forced batches after we close the wip batch then we must close the current wip L2 block,
	// since the processForcedBatches function needs to create new L2 blocks (cannot "reuse" the current wip L2 block if it's empty)
	if processForcedBatches {
		f.closeWIPL2Block(ctx)
	}

	// Wait until all L2 blocks are processed by the executor
	startWait := time.Now()
	f.pendingL2BlocksToProcessWG.Wait()
	elapsed := time.Since(startWait)
	log.Debugf("waiting for pending L2 blocks to be processed took: %v", elapsed)

	// Wait until all L2 blocks are store
	startWait = time.Now()
	f.pendingL2BlocksToStoreWG.Wait()
	log.Debugf("waiting for pending L2 blocks to be stored took: %v", time.Since(startWait))

	f.wipBatch.closingReason = closeReason

	// Close the wip batch
	var err error
	err = f.closeWIPBatch(ctx)
	if err != nil {
		return fmt.Errorf("failed to close batch, error: %v", err)
	}

	log.Infof("batch %d closed, closing reason: %s", f.wipBatch.batchNumber, closeReason)

	// Reprocess full batch as sanity check
	if f.cfg.SequentialBatchSanityCheck {
		// Do the full batch reprocess now
		_, _ = f.batchSanityCheck(ctx, f.wipBatch.batchNumber, f.wipBatch.initialStateRoot, f.wipBatch.finalStateRoot)
	} else {
		// Do the full batch reprocess in parallel
		go func() {
			_, _ = f.batchSanityCheck(ctx, f.wipBatch.batchNumber, f.wipBatch.initialStateRoot, f.wipBatch.finalStateRoot)
		}()
	}

	if f.wipBatch.batchNumber+1 == f.cfg.HaltOnBatchNumber {
		f.Halt(ctx, fmt.Errorf("finalizer reached stop sequencer on batch number: %d", f.cfg.HaltOnBatchNumber), false)
	}

	// Metadata for the next batch
	stateRoot := f.wipBatch.finalStateRoot
	lastBatchNumber := f.wipBatch.batchNumber

	// Process forced batches
	if processForcedBatches {
		lastBatchNumber, stateRoot = f.processForcedBatches(ctx, lastBatchNumber, stateRoot)
		// We must init/reset the wip L2 block from the state since processForcedBatches can created new L2 blocks
		f.initWIPL2Block(ctx)
	}

	f.wipBatch, err = f.openNewWIPBatch(ctx, lastBatchNumber+1, stateRoot)
	if err != nil {
		return fmt.Errorf("failed to open new wip batch, error: %v", err)
	}

	if f.wipL2Block != nil {
		f.wipBatch.imStateRoot = f.wipL2Block.imStateRoot
		// Subtract the WIP L2 block used resources to batch
		overflow, overflowResource := f.wipBatch.imRemainingResources.Sub(state.BatchResources{ZKCounters: f.wipL2Block.usedZKCounters, Bytes: f.wipL2Block.bytes})
		if overflow {
			return fmt.Errorf("failed to subtract L2 block [%d] used resources to new wip batch %d, overflow resource: %s",
				f.wipL2Block.trackingNum, f.wipBatch.batchNumber, overflowResource)
		}
	}

	log.Infof("new WIP batch %d", f.wipBatch.batchNumber)

	return nil
}

// openNewWIPBatch opens a new batch in the state and returns it as WipBatch
func (f *finalizer) openNewWIPBatch(ctx context.Context, batchNumber uint64, stateRoot common.Hash) (*Batch, error) {
	// open next batch
	newStateBatch := state.Batch{
		BatchNumber:    batchNumber,
		Coinbase:       f.sequencerAddress,
		Timestamp:      now(),
		StateRoot:      stateRoot,
		GlobalExitRoot: state.ZeroHash,
		LocalExitRoot:  state.ZeroHash,
	}

	dbTx, err := f.stateIntf.BeginStateTransaction(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin state transaction to open batch, error: %v", err)
	}

	// OpenBatch opens a new wip batch in the state
	err = f.stateIntf.OpenWIPBatch(ctx, newStateBatch, dbTx)
	if err != nil {
		if rollbackErr := dbTx.Rollback(ctx); rollbackErr != nil {
			return nil, fmt.Errorf("failed to rollback due to error when open a new wip batch, rollback error: %v, error: %v", rollbackErr, err)
		}
		return nil, fmt.Errorf("failed to open new wip batch, error: %v", err)
	}

	if err := dbTx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit database transaction for opening a wip batch, error: %v", err)
	}

	// Send batch bookmark to the datastream
	f.DSSendBatchBookmark(batchNumber)

	// Check if synchronizer is up-to-date
	for !f.isSynced(ctx) {
		log.Info("wait for synchronizer to sync last batch")
		time.Sleep(time.Second)
	}

	maxRemainingResources := getMaxRemainingResources(f.batchConstraints)

	return &Batch{
		batchNumber:             newStateBatch.BatchNumber,
		coinbase:                newStateBatch.Coinbase,
		initialStateRoot:        newStateBatch.StateRoot,
		imStateRoot:             newStateBatch.StateRoot,
		finalStateRoot:          newStateBatch.StateRoot,
		timestamp:               newStateBatch.Timestamp,
		imRemainingResources:    maxRemainingResources,
		finalRemainingResources: maxRemainingResources,
		closingReason:           state.EmptyClosingReason,
	}, err
}

// closeWIPBatch closes the current batch in the state
func (f *finalizer) closeWIPBatch(ctx context.Context) error {
	// Sanity check: batch must not be empty (should have L2 blocks)
	if f.wipBatch.isEmpty() {
		f.Halt(ctx, fmt.Errorf("closing WIP batch %d without L2 blocks and should have at least 1", f.wipBatch.batchNumber), false)
	}

	usedResources := getUsedBatchResources(f.batchConstraints, f.wipBatch.imRemainingResources)
	receipt := state.ProcessingReceipt{
		BatchNumber:    f.wipBatch.batchNumber,
		BatchResources: usedResources,
		ClosingReason:  f.wipBatch.closingReason,
	}

	dbTx, err := f.stateIntf.BeginStateTransaction(ctx)
	if err != nil {
		return err
	}

	err = f.stateIntf.CloseWIPBatch(ctx, receipt, dbTx)
	if err != nil {
		rollbackErr := dbTx.Rollback(ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back due to error when closing wip batch, rollback error: %v, error: %v", rollbackErr, err)
		}
		return err
	} else {
		err := dbTx.Commit(ctx)
		if err != nil {
			log.Errorf("error committing close wip batch, error: %v", err)
			return err
		}
	}

	return nil
}

// batchSanityCheck reprocesses a batch used as sanity check
func (f *finalizer) batchSanityCheck(ctx context.Context, batchNum uint64, initialStateRoot common.Hash, expectedNewStateRoot common.Hash) (*state.ProcessBatchResponse, error) {
	reprocessError := func(batch *state.Batch) {
		rawL2Blocks, err := state.DecodeBatchV2(batch.BatchL2Data)
		if err != nil {
			log.Errorf("error decoding BatchL2Data for batch %d, error: %v", batch.BatchNumber, err)
			return
		}

		// Log batch detailed info
		log.Errorf("batch %d sanity check error: initialStateRoot: %s, expectedNewStateRoot: %s", batch.BatchNumber, initialStateRoot, expectedNewStateRoot)
		for i, rawL2block := range rawL2Blocks.Blocks {
			log.Infof("block[%d], txs: %d, deltaTimestamp: %d, l1InfoTreeIndex: %d", i, len(rawL2block.Transactions), rawL2block.DeltaTimestamp, rawL2block.IndexL1InfoTree)
			for j, rawTx := range rawL2block.Transactions {
				log.Infof("block[%d].tx[%d]: %s, egpPct: %d, data: %s", batch.BatchNumber, i, j, rawTx.Tx.Hash(), rawTx.EfficiencyPercentage, hex.EncodeToHex(rawTx.Data))
			}
		}

		f.Halt(ctx, fmt.Errorf("batch sanity check error. Check previous errors in logs to know which was the cause"), false)
	}

	log.Debugf("batch %d sanity check: initialStateRoot: %s, expectedNewStateRoot: %s", batchNum, initialStateRoot, expectedNewStateRoot)

	batch, err := f.stateIntf.GetBatchByNumber(ctx, batchNum, nil)
	if err != nil {
		log.Errorf("failed to get batch %d, error: %v", batchNum, err)
		return nil, ErrGetBatchByNumber
	}

	batchRequest := state.ProcessRequest{
		BatchNumber:             batch.BatchNumber,
		L1InfoRoot_V2:           state.GetMockL1InfoRoot(),
		OldStateRoot:            initialStateRoot,
		Transactions:            batch.BatchL2Data,
		Coinbase:                batch.Coinbase,
		TimestampLimit_V2:       uint64(time.Now().Unix()),
		ForkID:                  f.stateIntf.GetForkIDByBatchNumber(batch.BatchNumber),
		SkipVerifyL1InfoRoot_V2: true,
		Caller:                  stateMetrics.DiscardCallerLabel,
	}
	batchRequest.L1InfoTreeData_V2, _, _, err = f.stateIntf.GetL1InfoTreeDataFromBatchL2Data(ctx, batch.BatchL2Data, nil)
	if err != nil {
		log.Errorf("failed to get L1InfoTreeData for batch %d, error: %v", batch.BatchNumber, err)
		reprocessError(nil)
		return nil, ErrGetBatchByNumber
	}

	var batchResponse *state.ProcessBatchResponse

	startProcessing := time.Now()
	batchResponse, err = f.stateIntf.ProcessBatchV2(ctx, batchRequest, false)
	endProcessing := time.Now()

	if err != nil {
		log.Errorf("failed to process batch %d, error: %v", batch.BatchNumber, err)
		reprocessError(batch)
		return nil, ErrProcessBatch
	}

	if batchResponse.ExecutorError != nil {
		log.Errorf("executor error when reprocessing batch %d, error: %v", batch.BatchNumber, batchResponse.ExecutorError)
		reprocessError(batch)
		return nil, ErrExecutorError
	}

	if batchResponse.IsRomOOCError {
		log.Errorf("failed to process batch %d because OutOfCounters", batch.BatchNumber)
		reprocessError(batch)

		payload, err := json.Marshal(batchRequest)
		if err != nil {
			log.Errorf("error marshaling payload, error: %v", err)
		} else {
			f.LogEvent(ctx, event.Level_Critical, event.EventID_ReprocessFullBatchOOC, string(payload), batchRequest)
		}

		return nil, ErrProcessBatchOOC
	}

	if batchResponse.NewStateRoot != expectedNewStateRoot {
		log.Errorf("new state root mismatch for batch %d, expected: %s, got: %s", batch.BatchNumber, expectedNewStateRoot.String(), batchResponse.NewStateRoot.String())
		reprocessError(batch)
		return nil, ErrStateRootNoMatch
	}

	err = f.stateIntf.UpdateBatchAsChecked(ctx, batch.BatchNumber, nil)
	if err != nil {
		log.Errorf("failed to update batch %d as checked, error: %v", batch.BatchNumber, err)
		reprocessError(batch)
		return nil, ErrUpdateBatchAsChecked
	}

	log.Infof("successful sanity check for batch %d, initialStateRoot: %s, stateRoot: %s, l2Blocks: %d, time: %v, used counters: %s",
		batch.BatchNumber, initialStateRoot, batchResponse.NewStateRoot.String(), len(batchResponse.BlockResponses),
		endProcessing.Sub(startProcessing), f.logZKCounters(batchResponse.UsedZkCounters))

	return batchResponse, nil
}

// maxTxsPerBatchReached checks if the batch has reached the maximum number of txs per batch
func (f *finalizer) maxTxsPerBatchReached(batch *Batch) bool {
	return (f.batchConstraints.MaxTxsPerBatch != 0) && (batch.countOfTxs >= int(f.batchConstraints.MaxTxsPerBatch))
}

// isBatchResourcesMarginExhausted checks if one of resources of the batch has reached the exhausted margin and returns the name of the exhausted resource
func (f *finalizer) isBatchResourcesMarginExhausted(resources state.BatchResources) (bool, string) {
	zkCounters := resources.ZKCounters
	result := false
	resourceName := ""
	if resources.Bytes <= f.getConstraintThresholdUint64(f.batchConstraints.MaxBatchBytesSize) {
		resourceName = "Bytes"
		result = true
	} else if zkCounters.Steps <= f.getConstraintThresholdUint32(f.batchConstraints.MaxSteps) {
		resourceName = "Steps"
		result = true
	} else if zkCounters.PoseidonPaddings <= f.getConstraintThresholdUint32(f.batchConstraints.MaxPoseidonPaddings) {
		resourceName = "PoseidonPaddings"
		result = true
	} else if zkCounters.PoseidonHashes <= f.getConstraintThresholdUint32(f.batchConstraints.MaxPoseidonHashes) {
		resourceName = "PoseidonHashes"
		result = true
	} else if zkCounters.Binaries <= f.getConstraintThresholdUint32(f.batchConstraints.MaxBinaries) {
		resourceName = "Binaries"
		result = true
	} else if zkCounters.KeccakHashes <= f.getConstraintThresholdUint32(f.batchConstraints.MaxKeccakHashes) {
		resourceName = "KeccakHashes"
		result = true
	} else if zkCounters.Arithmetics <= f.getConstraintThresholdUint32(f.batchConstraints.MaxArithmetics) {
		resourceName = "Arithmetics"
		result = true
	} else if zkCounters.MemAligns <= f.getConstraintThresholdUint32(f.batchConstraints.MaxMemAligns) {
		resourceName = "MemAligns"
		result = true
	} else if zkCounters.GasUsed <= f.getConstraintThresholdUint64(f.batchConstraints.MaxCumulativeGasUsed) {
		resourceName = "CumulativeGas"
		result = true
	} else if zkCounters.Sha256Hashes_V2 <= f.getConstraintThresholdUint32(f.batchConstraints.MaxSHA256Hashes) {
		resourceName = "SHA256Hashes"
		result = true
	}

	return result, resourceName
}

// getConstraintThresholdUint64 returns the threshold for the given input
func (f *finalizer) getConstraintThresholdUint64(input uint64) uint64 {
	return input * uint64(f.cfg.ResourceExhaustedMarginPct) / 100 //nolint:gomnd
}

// getConstraintThresholdUint32 returns the threshold for the given input
func (f *finalizer) getConstraintThresholdUint32(input uint32) uint32 {
	return input * f.cfg.ResourceExhaustedMarginPct / 100 //nolint:gomnd
}

// getUsedBatchResources calculates and returns the used resources of a batch from remaining resources
func getUsedBatchResources(constraints state.BatchConstraintsCfg, remainingResources state.BatchResources) state.BatchResources {
	return state.BatchResources{
		ZKCounters: state.ZKCounters{
			GasUsed:          constraints.MaxCumulativeGasUsed - remainingResources.ZKCounters.GasUsed,
			KeccakHashes:     constraints.MaxKeccakHashes - remainingResources.ZKCounters.KeccakHashes,
			PoseidonHashes:   constraints.MaxPoseidonHashes - remainingResources.ZKCounters.PoseidonHashes,
			PoseidonPaddings: constraints.MaxPoseidonPaddings - remainingResources.ZKCounters.PoseidonPaddings,
			MemAligns:        constraints.MaxMemAligns - remainingResources.ZKCounters.MemAligns,
			Arithmetics:      constraints.MaxArithmetics - remainingResources.ZKCounters.Arithmetics,
			Binaries:         constraints.MaxBinaries - remainingResources.ZKCounters.Binaries,
			Steps:            constraints.MaxSteps - remainingResources.ZKCounters.Steps,
			Sha256Hashes_V2:  constraints.MaxSHA256Hashes - remainingResources.ZKCounters.Sha256Hashes_V2,
		},
		Bytes: constraints.MaxBatchBytesSize - remainingResources.Bytes,
	}
}

// getMaxRemainingResources returns the max resources that can be used in a batch
func getMaxRemainingResources(constraints state.BatchConstraintsCfg) state.BatchResources {
	return state.BatchResources{
		ZKCounters: state.ZKCounters{
			GasUsed:          constraints.MaxCumulativeGasUsed,
			KeccakHashes:     constraints.MaxKeccakHashes,
			PoseidonHashes:   constraints.MaxPoseidonHashes,
			PoseidonPaddings: constraints.MaxPoseidonPaddings,
			MemAligns:        constraints.MaxMemAligns,
			Arithmetics:      constraints.MaxArithmetics,
			Binaries:         constraints.MaxBinaries,
			Steps:            constraints.MaxSteps,
			Sha256Hashes_V2:  constraints.MaxSHA256Hashes,
		},
		Bytes: constraints.MaxBatchBytesSize,
	}
}

// checkIfFinalizeBatch returns true if the batch must be closed due to a closing reason, also it returns the description of the close reason
func (f *finalizer) checkIfFinalizeBatch() (bool, state.ClosingReason) {
	// Max txs per batch
	if f.maxTxsPerBatchReached(f.wipBatch) {
		log.Infof("closing batch %d, because it reached the maximum number of txs", f.wipBatch.batchNumber)
		return true, state.MaxTxsClosingReason
	}

	// Batch resource (zkCounters or batch bytes) margin exhausted
	exhausted, resourceDesc := f.isBatchResourcesMarginExhausted(f.wipBatch.imRemainingResources)
	if exhausted {
		log.Infof("closing batch %d because it exhausted margin for %s batch resource", f.wipBatch.batchNumber, resourceDesc)
		return true, state.ResourceMarginExhaustedClosingReason
	}

	// Forced batch deadline
	if f.nextForcedBatchDeadline != 0 && now().Unix() >= f.nextForcedBatchDeadline {
		log.Infof("closing batch %d, forced batch deadline encountered", f.wipBatch.batchNumber)
		return true, state.ForcedBatchDeadlineClosingReason
	}

	// Batch timestamp resolution
	if !f.wipBatch.isEmpty() && f.wipBatch.timestamp.Add(f.cfg.BatchMaxDeltaTimestamp.Duration).Before(time.Now()) {
		log.Infof("closing batch %d, because of batch max delta timestamp reached", f.wipBatch.batchNumber)
		return true, state.MaxDeltaTimestampClosingReason
	}

	return false, ""
}
