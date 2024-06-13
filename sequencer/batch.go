package sequencer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/event"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	stateMetrics "github.com/0xPolygonHermez/zkevm-node/state/metrics"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

// Batch represents a wip or processed batch.
type Batch struct {
	batchNumber                 uint64
	coinbase                    common.Address
	timestamp                   time.Time
	initialStateRoot            common.Hash // initial stateRoot of the batch
	imStateRoot                 common.Hash // intermediate stateRoot when processing tx-by-tx
	finalStateRoot              common.Hash // final stateroot of the batch when a L2 block is processed
	countOfTxs                  int
	countOfL2Blocks             int
	imRemainingResources        state.BatchResources // remaining batch resources when processing tx-by-tx
	imHighReservedZKCounters    state.ZKCounters
	finalRemainingResources     state.BatchResources // remaining batch resources when a L2 block is processed
	finalHighReservedZKCounters state.ZKCounters
	closingReason               state.ClosingReason
	finalLocalExitRoot          common.Hash
}

func (b *Batch) isEmpty() bool {
	return b.countOfL2Blocks == 0
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

	remainingResources := getMaxBatchResources(f.batchConstraints)
	overflow, overflowResource := remainingResources.Sub(wipStateBatch.Resources)
	if overflow {
		return nil, fmt.Errorf("failed to subtract used resources when setting the wip batch to the state batch %d, overflow resource: %s", wipStateBatch.BatchNumber, overflowResource)
	}

	wipBatch := &Batch{
		batchNumber:                 wipStateBatch.BatchNumber,
		coinbase:                    wipStateBatch.Coinbase,
		imStateRoot:                 wipStateBatch.StateRoot,
		initialStateRoot:            prevStateBatch.StateRoot,
		finalStateRoot:              wipStateBatch.StateRoot,
		timestamp:                   wipStateBatch.Timestamp,
		countOfL2Blocks:             len(wipStateBatchBlocks.Blocks),
		countOfTxs:                  wipStateBatchCountOfTxs,
		imRemainingResources:        remainingResources,
		finalRemainingResources:     remainingResources,
		imHighReservedZKCounters:    wipStateBatch.HighReservedZKCounters,
		finalHighReservedZKCounters: wipStateBatch.HighReservedZKCounters,
		finalLocalExitRoot:          wipStateBatch.LocalExitRoot,
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
		f.wipBatch = f.openNewWIPBatch(lastStateBatch.BatchNumber+1, lastStateBatch.StateRoot)
		f.pipBatch = nil
		f.sipBatch = nil
	} else { /// if it's not closed, it is the wip/pip/sip batch
		f.wipBatch, err = f.setWIPBatch(ctx, lastStateBatch)
		if err != nil {
			log.Fatalf("failed to set wip batch, error: %v", err)
		}
		f.pipBatch = f.wipBatch
		f.sipBatch = f.wipBatch
	}

	log.Infof("initial batch: %d, initialStateRoot: %s, stateRoot: %s, coinbase: %s",
		f.wipBatch.batchNumber, f.wipBatch.initialStateRoot, f.wipBatch.finalStateRoot, f.wipBatch.coinbase)
}

func (f *finalizer) processL2BlockReorg(ctx context.Context) error {
	f.waitPendingL2Blocks()

	if f.sipBatch != nil && f.sipBatch.batchNumber != f.wipBatch.batchNumber {
		// If the sip batch is the previous to the current wip batch and it's still open these means that the L2 block that caused
		// the reorg is the first L2 block of the wip batch, therefore we need to close sip batch before to continue.
		// If we don't close the sip batch the initWIPBatch function will load the sip batch as the initial one and when trying to reprocess
		// the first tx reorged we can have a batch resource overflow (if we have closed the sip batch for this reason) and we will return
		// the reorged tx to the worker (calling UpdateTxZKCounters) missing the order in which we need to reprocess the reorged txs

		err := f.finalizeSIPBatch(ctx)
		if err != nil {
			return fmt.Errorf("error finalizing sip batch, error: %v", err)
		}
	}

	f.workerIntf.RestoreTxsPendingToStore(ctx)

	f.initWIPBatch(ctx)

	f.initWIPL2Block(ctx)

	// Since when processing the L2 block reorg we sync the state root we can reset next state root syncing
	f.scheduleNextStateRootSync()

	f.l2BlockReorg.Store(false)

	return nil
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
		f.Halt(ctx, fmt.Errorf("failed to create new wip batch, error: %v", err), true)
	}

	// If we have closed the wipL2Block then we open a new one
	if f.wipL2Block == nil {
		f.openNewWIPL2Block(ctx, prevTimestamp, &prevL1InfoTreeIndex)
	}
}

// finalizeSIPBatch closes the current store-in-progress batch
func (f *finalizer) finalizeSIPBatch(ctx context.Context) error {
	dbTx, err := f.stateIntf.BeginStateTransaction(ctx)
	if err != nil {
		return fmt.Errorf("error creating db transaction to close sip batch %d, error: %v", f.sipBatch.batchNumber, err)
	}

	// Close sip batch (close in statedb)
	err = f.closeSIPBatch(ctx, dbTx)
	if err != nil {
		return fmt.Errorf("failed to close sip batch %d, error: %v", f.sipBatch.batchNumber, err)
	}

	if err != nil {
		rollbackErr := dbTx.Rollback(ctx)
		if rollbackErr != nil {
			return fmt.Errorf("error when rollback db transaction to close sip batch %d, error: %v", f.sipBatch.batchNumber, rollbackErr)
		}
		return err
	}

	err = dbTx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("error when commit db transaction to close sip batch %d, error: %v", f.sipBatch.batchNumber, err)
	}

	return nil
}

// closeAndOpenNewWIPBatch closes the current batch and opens a new one, potentially processing forced batches between the batch is closed and the resulting new wip batch
func (f *finalizer) closeAndOpenNewWIPBatch(ctx context.Context, closeReason state.ClosingReason) error {
	f.nextForcedBatchesMux.Lock()
	processForcedBatches := len(f.nextForcedBatches) > 0
	f.nextForcedBatchesMux.Unlock()

	f.wipBatch.closingReason = closeReason

	var lastStateRoot common.Hash

	//TODO: review forced batches implementation since is not good "idea" to check here for forced batches, maybe is better to do it on finalizeBatches loop
	if processForcedBatches {
		// If we have reach the time to sync stateroot or we will process forced batches we must close the current wip L2 block and wip batch
		f.closeWIPL2Block(ctx)
		// We need to wait that all pending L2 blocks are processed and stored
		f.waitPendingL2Blocks()

		lastStateRoot = f.sipBatch.finalStateRoot

		err := f.finalizeSIPBatch(ctx)
		if err != nil {
			return fmt.Errorf("error finalizing sip batch %d when processing forced batches, error: %v", f.sipBatch.batchNumber, err)
		}
	} else {
		lastStateRoot = f.wipBatch.imStateRoot
	}

	// Close the wip batch. After will close them f.wipBatch will be nil, therefore we store in local variables the info we need from the f.wipBatch
	lastBatchNumber := f.wipBatch.batchNumber

	f.closeWIPBatch(ctx)

	if lastBatchNumber+1 == f.cfg.HaltOnBatchNumber {
		f.waitPendingL2Blocks()

		// We finalize the current sip batch
		err := f.finalizeSIPBatch(ctx)
		if err != nil {
			return fmt.Errorf("error finalizing sip batch %d when halting on batch %d", f.sipBatch.batchNumber, f.cfg.HaltOnBatchNumber)
		}

		f.Halt(ctx, fmt.Errorf("finalizer reached stop sequencer on batch number: %d", f.cfg.HaltOnBatchNumber), false)
	}

	// Process forced batches
	if processForcedBatches {
		lastBatchNumber, lastStateRoot = f.processForcedBatches(ctx, lastBatchNumber, lastStateRoot)
	}

	f.wipBatch = f.openNewWIPBatch(lastBatchNumber+1, lastStateRoot)

	if processForcedBatches {
		// We need to init/reset the wip L2 block in case we have processed forced batches
		f.initWIPL2Block(ctx)
	} else if f.wipL2Block != nil {
		// If we are "reusing" the wip L2 block because it's empty we assign it to the new wip batch
		f.wipBatch.imStateRoot = f.wipL2Block.imStateRoot
		f.wipL2Block.batch = f.wipBatch

		// We subtract the wip L2 block used resources to the new wip batch
		overflow, overflowResource := f.wipBatch.imRemainingResources.Sub(state.BatchResources{ZKCounters: f.wipL2Block.usedZKCountersOnNew, Bytes: f.wipL2Block.bytes})
		if overflow {
			return fmt.Errorf("failed to subtract L2 block [%d] used resources to new wip batch %d, overflow resource: %s",
				f.wipL2Block.trackingNum, f.wipBatch.batchNumber, overflowResource)
		}
	}

	log.Infof("new wip batch %d", f.wipBatch.batchNumber)

	return nil
}

// openNewWIPBatch opens a new batch in the state and returns it as WipBatch
func (f *finalizer) openNewWIPBatch(batchNumber uint64, stateRoot common.Hash) *Batch {
	maxRemainingResources := getMaxBatchResources(f.batchConstraints)

	return &Batch{
		batchNumber:             batchNumber,
		coinbase:                f.l2Coinbase,
		initialStateRoot:        stateRoot,
		imStateRoot:             stateRoot,
		finalStateRoot:          stateRoot,
		timestamp:               now(),
		imRemainingResources:    maxRemainingResources,
		finalRemainingResources: maxRemainingResources,
		closingReason:           state.EmptyClosingReason,
		finalLocalExitRoot:      state.ZeroHash,
	}
}

// insertSIPBatch inserts a new state-in-progress batch in the state db
func (f *finalizer) insertSIPBatch(ctx context.Context, batchNumber uint64, stateRoot common.Hash, dbTx pgx.Tx) error {
	// open next batch
	newStateBatch := state.Batch{
		BatchNumber:    batchNumber,
		Coinbase:       f.l2Coinbase,
		Timestamp:      now(),
		StateRoot:      stateRoot,
		GlobalExitRoot: state.ZeroHash,
		LocalExitRoot:  state.ZeroHash,
	}

	// OpenBatch opens a new wip batch in the state
	//TODO: rename OpenWipBatch to InsertBatch
	err := f.stateIntf.OpenWIPBatch(ctx, newStateBatch, dbTx)
	if err != nil {
		return fmt.Errorf("failed to insert new batch in state db, error: %v", err)
	}

	// Send batch bookmark to the datastream
	f.DSSendBatchBookmark(ctx, batchNumber)
	// Send batch start to the datastream
	f.DSSendBatchStart(ctx, batchNumber, false)

	// Check if synchronizer is up-to-date
	//TODO: review if this is needed
	for !f.isSynced(ctx) {
		log.Info("wait for synchronizer to sync last batch")
		time.Sleep(time.Second)
	}

	return nil
}

// closeWIPBatch closes the current wip batch
func (f *finalizer) closeWIPBatch(ctx context.Context) {
	// Sanity check: batch must not be empty (should have L2 blocks)
	if f.wipBatch.isEmpty() {
		f.Halt(ctx, fmt.Errorf("closing wip batch %d without L2 blocks and should have at least 1", f.wipBatch.batchNumber), false)
	}

	log.Infof("wip batch %d closed, closing reason: %s", f.wipBatch.batchNumber, f.wipBatch.closingReason)

	f.wipBatch = nil
}

// closeSIPBatch closes the current sip batch in the state
func (f *finalizer) closeSIPBatch(ctx context.Context, dbTx pgx.Tx) error {
	// Sanity check: this can't happen
	if f.sipBatch == nil {
		f.Halt(ctx, fmt.Errorf("closing sip batch that is nil"), false)
	}

	// Sanity check: batch must not be empty (should have L2 blocks)
	if f.sipBatch.isEmpty() {
		f.Halt(ctx, fmt.Errorf("closing sip batch %d without L2 blocks and should have at least 1", f.sipBatch.batchNumber), false)
	}

	usedResources := getUsedBatchResources(f.batchConstraints, f.sipBatch.imRemainingResources)
	receipt := state.ProcessingReceipt{
		BatchNumber:    f.sipBatch.batchNumber,
		BatchResources: usedResources,
		ClosingReason:  f.sipBatch.closingReason,
	}

	err := f.stateIntf.CloseWIPBatch(ctx, receipt, dbTx)

	if err != nil {
		return err
	}

	// We store values needed for the batch sanity check in local variables, as we can execute the sanity check in a go func (parallel) and in this case f.sipBatch will be nil during some time
	batchNumber := f.sipBatch.batchNumber
	initialStateRoot := f.sipBatch.initialStateRoot
	finalStateRoot := f.sipBatch.finalStateRoot

	// Reprocess full batch as sanity check
	if f.cfg.SequentialBatchSanityCheck {
		// Do the full batch reprocess now
		_, _ = f.batchSanityCheck(ctx, batchNumber, initialStateRoot, finalStateRoot)
	} else {
		// Do the full batch reprocess in parallel
		go func() {
			_, _ = f.batchSanityCheck(ctx, batchNumber, initialStateRoot, finalStateRoot)
		}()
	}

	// Sent batch to DS
	f.DSSendBatchEnd(ctx, f.sipBatch.batchNumber, f.sipBatch.finalStateRoot, f.sipBatch.finalLocalExitRoot)

	log.Infof("sip batch %d closed in statedb, closing reason: %s", f.sipBatch.batchNumber, f.sipBatch.closingReason)

	f.sipBatch = nil

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
		batchLog := ""
		totalTxs := 0
		for blockIdx, rawL2block := range rawL2Blocks.Blocks {
			totalTxs += len(rawL2block.Transactions)
			batchLog += fmt.Sprintf("block[%d], txs: %d, deltaTimestamp: %d, l1InfoTreeIndex: %d\n", blockIdx, len(rawL2block.Transactions), rawL2block.DeltaTimestamp, rawL2block.IndexL1InfoTree)
			for txIdx, rawTx := range rawL2block.Transactions {
				batchLog += fmt.Sprintf("   tx[%d]: %s, egpPct: %d\n", txIdx, rawTx.Tx.Hash(), rawTx.EfficiencyPercentage)
			}
		}
		log.Infof("dump batch %d, blocks: %d, txs: %d\n%s", batch.BatchNumber, len(rawL2Blocks.Blocks), totalTxs, batchLog)

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

	startProcessing := time.Now()
	batchResponse, contextid, err := f.stateIntf.ProcessBatchV2(ctx, batchRequest, false)
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

	log.Infof("successful sanity check for batch %d, initialStateRoot: %s, stateRoot: %s, l2Blocks: %d, time: %v, used counters: %s, contextId: %s",
		batch.BatchNumber, initialStateRoot, batchResponse.NewStateRoot.String(), len(batchResponse.BlockResponses),
		endProcessing.Sub(startProcessing), f.logZKCounters(batchResponse.UsedZkCounters), contextid)

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

// getMaxBatchResources returns the max resources that can be used in a batch
func getMaxBatchResources(constraints state.BatchConstraintsCfg) state.BatchResources {
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

// getNeededZKCounters returns the needed counters to fit a tx in the wip batch. The needed counters are the counters used by the tx plus the high reserved counters.
// It will take into account the current high reserved counter got with previous txs but also checking reserved counters diff needed by this tx, since could be greater.
func getNeededZKCounters(highReservedCounters state.ZKCounters, usedCounters state.ZKCounters, reservedCounters state.ZKCounters) (state.ZKCounters, state.ZKCounters) {
	neededCounter := func(counterName string, highCounter uint32, usedCounter uint32, reservedCounter uint32) (uint32, uint32) {
		if reservedCounter < usedCounter {
			log.Warnf("%s reserved counter %d is less than used counter %d, this shouldn't be possible", counterName, reservedCounter, usedCounter)
			return usedCounter + highCounter, highCounter
		}
		diffReserved := reservedCounter - usedCounter
		if diffReserved > highCounter { // reserved counter for this tx (difference) is greater that the high reserved counter got in previous txs
			return usedCounter + diffReserved, diffReserved
		} else {
			return usedCounter + highCounter, highCounter
		}
	}

	needed := state.ZKCounters{}
	newHigh := state.ZKCounters{}

	needed.Arithmetics, newHigh.Arithmetics = neededCounter("Arithmetics", highReservedCounters.Arithmetics, usedCounters.Arithmetics, reservedCounters.Arithmetics)
	needed.Binaries, newHigh.Binaries = neededCounter("Binaries", highReservedCounters.Binaries, usedCounters.Binaries, reservedCounters.Binaries)
	needed.KeccakHashes, newHigh.KeccakHashes = neededCounter("KeccakHashes", highReservedCounters.KeccakHashes, usedCounters.KeccakHashes, reservedCounters.KeccakHashes)
	needed.MemAligns, newHigh.MemAligns = neededCounter("MemAligns", highReservedCounters.MemAligns, usedCounters.MemAligns, reservedCounters.MemAligns)
	needed.PoseidonHashes, newHigh.PoseidonHashes = neededCounter("PoseidonHashes", highReservedCounters.PoseidonHashes, usedCounters.PoseidonHashes, reservedCounters.PoseidonHashes)
	needed.PoseidonPaddings, newHigh.PoseidonPaddings = neededCounter("PoseidonPaddings", highReservedCounters.PoseidonPaddings, usedCounters.PoseidonPaddings, reservedCounters.PoseidonPaddings)
	needed.Sha256Hashes_V2, newHigh.Sha256Hashes_V2 = neededCounter("Sha256Hashes_V2", highReservedCounters.Sha256Hashes_V2, usedCounters.Sha256Hashes_V2, reservedCounters.Sha256Hashes_V2)
	needed.Steps, newHigh.Steps = neededCounter("Steps", highReservedCounters.Steps, usedCounters.Steps, reservedCounters.Steps)

	if reservedCounters.GasUsed < usedCounters.GasUsed {
		log.Warnf("gasUsed reserved counter %d is less than used counter %d, this shouldn't be possible", reservedCounters.GasUsed, usedCounters.GasUsed)
		needed.GasUsed = usedCounters.GasUsed + highReservedCounters.GasUsed
	} else {
		diffReserved := reservedCounters.GasUsed - usedCounters.GasUsed
		if diffReserved > highReservedCounters.GasUsed {
			needed.GasUsed = usedCounters.GasUsed + diffReserved
			newHigh.GasUsed = diffReserved
		} else {
			needed.GasUsed = usedCounters.GasUsed + highReservedCounters.GasUsed
			newHigh.GasUsed = highReservedCounters.GasUsed
		}
	}

	return needed, newHigh
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
