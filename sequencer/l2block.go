package sequencer

import (
	"context"
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/event"
	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/state"
	stateMetrics "github.com/0xPolygonHermez/zkevm-node/state/metrics"
	"github.com/ethereum/go-ethereum/common"
)

// L2Block represents a wip or processed L2 block
type L2Block struct {
	createdAt                 time.Time
	trackingNum               uint64
	timestamp                 uint64
	deltaTimestamp            uint32
	imStateRoot               common.Hash
	l1InfoTreeExitRoot        state.L1InfoTreeExitRootStorageEntry
	l1InfoTreeExitRootChanged bool
	bytes                     uint64
	usedZKCountersOnNew       state.ZKCounters
	reservedZKCountersOnNew   state.ZKCounters
	highReservedZKCounters    state.ZKCounters
	transactions              []*TxTracker
	batch                     *Batch
	batchResponse             *state.ProcessBatchResponse
	metrics                   metrics
}

func (b *L2Block) isEmpty() bool {
	return len(b.transactions) == 0
}

// addTx adds a tx to the L2 block
func (b *L2Block) addTx(tx *TxTracker) {
	b.transactions = append(b.transactions, tx)
}

// getL1InfoTreeIndex returns the L1InfoTreeIndex that must be used when processing/storing the block
func (b *L2Block) getL1InfoTreeIndex() uint32 {
	// If the L1InfoTreeIndex has changed in this block then we return the new index, otherwise we return 0
	if b.l1InfoTreeExitRootChanged {
		return b.l1InfoTreeExitRoot.L1InfoTreeIndex
	} else {
		return 0
	}
}

// initWIPL2Block inits the wip L2 block
func (f *finalizer) initWIPL2Block(ctx context.Context) {
	// Wait to l1InfoTree to be updated for first time
	f.lastL1InfoTreeCond.L.Lock()
	for !f.lastL1InfoTreeValid {
		log.Infof("waiting for L1InfoTree to be updated")
		f.lastL1InfoTreeCond.Wait()
	}
	f.lastL1InfoTreeCond.L.Unlock()

	lastL2Block, err := f.stateIntf.GetLastL2Block(ctx, nil)
	if err != nil {
		log.Fatalf("failed to get last L2 block number, error: %v", err)
	}

	f.openNewWIPL2Block(ctx, uint64(lastL2Block.ReceivedAt.Unix()), nil)
}

// addPendingL2BlockToProcess adds a pending L2 block that is closed and ready to be processed by the executor
func (f *finalizer) addPendingL2BlockToProcess(ctx context.Context, l2Block *L2Block) {
	f.pendingL2BlocksToProcessWG.Add(1)

	select {
	case f.pendingL2BlocksToProcess <- l2Block:
	case <-ctx.Done():
		// If context is cancelled before we can send to the channel, we must decrement the WaitGroup count and
		// delete the pending TxToStore added in the worker
		f.pendingL2BlocksToProcessWG.Done()
	}
}

// addPendingL2BlockToStore adds a L2 block that is ready to be stored in the state DB once its flushid has been stored by the executor
func (f *finalizer) addPendingL2BlockToStore(ctx context.Context, l2Block *L2Block) {
	f.pendingL2BlocksToStoreWG.Add(1)

	select {
	case f.pendingL2BlocksToStore <- l2Block:
	case <-ctx.Done():
		// If context is cancelled before we can send to the channel, we must decrement the WaitGroup count and
		// delete the pending TxToStore added in the worker
		f.pendingL2BlocksToStoreWG.Done()
		for _, tx := range l2Block.transactions {
			f.workerIntf.DeleteTxPendingToStore(tx.Hash, tx.From)
		}
	}
}

// processPendingL2Blocks processes (executor) the pending to process L2 blocks
func (f *finalizer) processPendingL2Blocks(ctx context.Context) {
	//rand.Seed(time.Now().UnixNano())

	for {
		select {
		case l2Block, ok := <-f.pendingL2BlocksToProcess:
			if !ok {
				// Channel is closed
				return
			}

			// if l2BlockReorg we need to "flush" the channel to discard pending L2Blocks
			if f.l2BlockReorg.Load() {
				f.pendingL2BlocksToProcessWG.Done()
				continue
			}

			err := f.processL2Block(ctx, l2Block)

			if err != nil {
				halt := false
				if f.lastL2BlockWasReorg {
					// We had 2 consecutives reorg in the same L2 block, we halt after log/dump the info
					halt = true
				} else {
					f.l2BlockReorg.Store(true)
					f.lastL2BlockWasReorg = true
				}

				warnmsg := fmt.Sprintf("sequencer L2 block [%d] reorg detected, batch: %d, processing it...", l2Block.trackingNum, l2Block.batch.batchNumber)
				log.Warnf(warnmsg)
				f.LogEvent(ctx, event.Level_Critical, event.EventID_L2BlockReorg, warnmsg, nil)

				// Dump L2Block info
				f.dumpL2Block(l2Block)

				if halt {
					f.Halt(ctx, fmt.Errorf("consecutives L2 block reorgs in the same L2 block [%d]", l2Block.trackingNum), false)
				}
			} else {
				f.lastL2BlockWasReorg = false
			}

			f.pendingL2BlocksToProcessWG.Done()

		case <-ctx.Done():
			// The context was cancelled from outside, Wait for all goroutines to finish, cleanup and exit
			f.pendingL2BlocksToProcessWG.Wait()
			return
		default:
			time.Sleep(100 * time.Millisecond) //nolint:gomnd
		}
	}
}

// storePendingTransactions stores the pending L2 blocks in the database
func (f *finalizer) storePendingL2Blocks(ctx context.Context) {
	for {
		select {
		case l2Block, ok := <-f.pendingL2BlocksToStore:
			if !ok {
				// Channel is closed
				return
			}

			err := f.storeL2Block(ctx, l2Block)

			if err != nil {
				// Dump L2Block info
				f.dumpL2Block(l2Block)
				f.Halt(ctx, fmt.Errorf("error storing L2 block %d [%d], error: %v", l2Block.batchResponse.BlockResponses[0].BlockNumber, l2Block.trackingNum, err), true)
			}

			f.pendingL2BlocksToStoreWG.Done()
		case <-ctx.Done():
			// The context was cancelled from outside, Wait for all goroutines to finish, cleanup and exit
			f.pendingL2BlocksToStoreWG.Wait()
			return
		default:
			time.Sleep(100 * time.Millisecond) //nolint:gomnd
		}
	}
}

// processL2Block process a L2 Block and adds it to the pendingL2BlocksToStore channel
func (f *finalizer) processL2Block(ctx context.Context, l2Block *L2Block) error {
	processStart := time.Now()

	if f.pipBatch == nil {
		f.pipBatch = l2Block.batch
	} else if f.pipBatch.batchNumber != l2Block.batch.batchNumber {
		// We have received the first L2 block of the next batch to process
		// We need to "propagate" finalStateRoot to the new batch as initalStateRoot/finalStateRoot and set it as the current pipBatch
		l2Block.batch.initialStateRoot = f.pipBatch.finalStateRoot
		l2Block.batch.finalStateRoot = f.pipBatch.finalStateRoot
		f.pipBatch = l2Block.batch
	}

	initialStateRoot := f.pipBatch.finalStateRoot

	log.Infof("processing L2 block [%d], batch: %d, deltaTimestamp: %d, timestamp: %d, l1InfoTreeIndex: %d, l1InfoTreeIndexChanged: %v, initialStateRoot: %s txs: %d",
		l2Block.trackingNum, l2Block.batch.batchNumber, l2Block.deltaTimestamp, l2Block.timestamp, l2Block.l1InfoTreeExitRoot.L1InfoTreeIndex,
		l2Block.l1InfoTreeExitRootChanged, initialStateRoot, len(l2Block.transactions))

	batchResponse, batchL2DataSize, contextId, err := f.executeL2Block(ctx, initialStateRoot, l2Block)

	if err != nil {
		return fmt.Errorf("failed to execute L2 block [%d], error: %v", l2Block.trackingNum, err)
	}

	if len(batchResponse.BlockResponses) != 1 {
		return fmt.Errorf("length of batchResponse.BlockRespones returned by the executor is %d and must be 1", len(batchResponse.BlockResponses))
	}

	blockResponse := batchResponse.BlockResponses[0]

	// Sanity check. Check blockResponse.TransactionsReponses match l2Block.Transactions length, order and tx hashes
	if len(blockResponse.TransactionResponses) != len(l2Block.transactions) {
		return fmt.Errorf("length of TransactionsResponses %d doesn't match length of l2Block.transactions %d", len(blockResponse.TransactionResponses), len(l2Block.transactions))
	}
	for i, txResponse := range blockResponse.TransactionResponses {
		if txResponse.TxHash != l2Block.transactions[i].Hash {
			return fmt.Errorf("blockResponse.TransactionsResponses[%d] hash %s doesn't match l2Block.transactions[%d] hash %s", i, txResponse.TxHash.String(), i, l2Block.transactions[i].Hash)
		}
	}

	// Sanity check. Check blockResponse.timestamp matches l2block.timestamp
	if blockResponse.Timestamp != l2Block.timestamp {
		return fmt.Errorf("blockResponse.Timestamp %d doesn't match l2Block.timestamp %d", blockResponse.Timestamp, l2Block.timestamp)
	}

	l2Block.batchResponse = batchResponse

	// Check if needed resources of the L2 block fits in the remaining batch resources
	// Needed resources are the used resources plus the max difference between used and reserved of all the L2 blocks (including this) in the batch
	neededZKCounters, newHighZKCounters := getNeededZKCounters(l2Block.batch.finalHighReservedZKCounters, batchResponse.UsedZkCounters, batchResponse.ReservedZkCounters)

	// Update finalRemainingResources of the batch
	fits, overflowResource := l2Block.batch.finalRemainingResources.Fits(state.BatchResources{ZKCounters: neededZKCounters, Bytes: batchL2DataSize})
	if fits {
		subOverflow, overflowResource := l2Block.batch.finalRemainingResources.Sub(state.BatchResources{ZKCounters: batchResponse.UsedZkCounters, Bytes: batchL2DataSize})
		if subOverflow { // Sanity check, this cannot happen as reservedZKCounters should be >= that usedZKCounters
			return fmt.Errorf("error subtracting L2 block %d [%d] needed resources from the batch %d, overflow resource: %s, batch bytes: %d, L2 block bytes: %d, counters: {batch: %s, used: %s, reserved: %s, needed: %s, high: %s}",
				blockResponse.BlockNumber, l2Block.trackingNum, l2Block.batch.batchNumber, overflowResource, l2Block.batch.finalRemainingResources.Bytes, batchL2DataSize,
				f.logZKCounters(l2Block.batch.finalRemainingResources.ZKCounters), f.logZKCounters(batchResponse.UsedZkCounters), f.logZKCounters(batchResponse.ReservedZkCounters), f.logZKCounters(neededZKCounters), f.logZKCounters(l2Block.batch.finalHighReservedZKCounters))
		}

		l2Block.batch.finalHighReservedZKCounters = newHighZKCounters
		l2Block.highReservedZKCounters = l2Block.batch.finalHighReservedZKCounters
	} else {
		overflowLog := fmt.Sprintf("L2 block %d [%d] needed resources exceeds the remaining batch %d resources, overflow resource: %s, batch bytes: %d, L2 block bytes: %d, counters: {batch: %s, used: %s, reserved: %s, needed: %s, high: %s}",
			blockResponse.BlockNumber, l2Block.trackingNum, l2Block.batch.batchNumber, overflowResource, l2Block.batch.finalRemainingResources.Bytes, batchL2DataSize,
			f.logZKCounters(l2Block.batch.finalRemainingResources.ZKCounters), f.logZKCounters(batchResponse.UsedZkCounters), f.logZKCounters(batchResponse.ReservedZkCounters), f.logZKCounters(neededZKCounters), f.logZKCounters(l2Block.batch.finalHighReservedZKCounters))

		f.LogEvent(ctx, event.Level_Warning, event.EventID_ReservedZKCountersOverflow, overflowLog, nil)

		return fmt.Errorf(overflowLog)
	}

	// Update finalStateRoot/finalLocalExitRoot of the batch to the newStateRoot/newLocalExitRoot for the L2 block
	l2Block.batch.finalStateRoot = l2Block.batchResponse.NewStateRoot
	l2Block.batch.finalLocalExitRoot = l2Block.batchResponse.NewLocalExitRoot

	f.updateFlushIDs(batchResponse.FlushID, batchResponse.StoredFlushID)

	var waitStoreL2Block time.Duration
	if f.pendingL2BlocksToStoreWG.Count() > 0 {
		startWait := time.Now()
		f.pendingL2BlocksToStoreWG.Wait()
		waitStoreL2Block = time.Since(startWait)
		log.Debugf("waiting for previous L2 block to be stored took: %v", waitStoreL2Block)
	}
	f.addPendingL2BlockToStore(ctx, l2Block)

	// metrics
	l2Block.metrics.l2BlockTimes.sequencer = time.Since(processStart) - l2Block.metrics.l2BlockTimes.executor
	if f.cfg.SequentialProcessL2Block {
		l2Block.metrics.close(l2Block.createdAt, int64(len(l2Block.transactions)), f.cfg.SequentialProcessL2Block)
	}
	f.metrics.addL2BlockMetrics(l2Block.metrics)

	log.Infof("processed L2 block %d [%d], batch: %d, deltaTimestamp: %d, timestamp: %d, l1InfoTreeIndex: %d, l1InfoTreeIndexChanged: %v, initialStateRoot: %s, newStateRoot: %s, txs: %d/%d, blockHash: %s, infoRoot: %s, waitStoreL2Block: %v, counters: {used: %s, reserved: %s, needed: %s, high: %s}, contextId: %s",
		blockResponse.BlockNumber, l2Block.trackingNum, l2Block.batch.batchNumber, l2Block.deltaTimestamp, l2Block.timestamp, l2Block.l1InfoTreeExitRoot.L1InfoTreeIndex, l2Block.l1InfoTreeExitRootChanged, initialStateRoot, l2Block.batchResponse.NewStateRoot,
		len(l2Block.transactions), len(blockResponse.TransactionResponses), blockResponse.BlockHash, blockResponse.BlockInfoRoot, waitStoreL2Block,
		f.logZKCounters(batchResponse.UsedZkCounters), f.logZKCounters(batchResponse.ReservedZkCounters), f.logZKCounters(neededZKCounters), f.logZKCounters(l2Block.batch.finalHighReservedZKCounters), contextId)

	if f.cfg.Metrics.EnableLog {
		log.Infof("metrics-log: {l2block: {num: %d, trackingNum: %d, metrics: {%s}}, interval: {startAt: %d, metrics: {%s}}}",
			blockResponse.BlockNumber, l2Block.trackingNum, l2Block.metrics.log(), f.metrics.startsAt().Unix(), f.metrics.log())
	}

	return nil
}

// executeL2Block executes a L2 Block in the executor and returns the batch response from the executor and the batchL2Data size
func (f *finalizer) executeL2Block(ctx context.Context, initialStateRoot common.Hash, l2Block *L2Block) (*state.ProcessBatchResponse, uint64, string, error) {
	executeL2BLockError := func(err error) {
		log.Errorf("execute L2 block [%d] error %v, batch: %d, initialStateRoot: %s", l2Block.trackingNum, err, l2Block.batch.batchNumber, initialStateRoot)
		// Log batch detailed info
		for i, tx := range l2Block.transactions {
			log.Infof("batch: %d, block: [%d], tx position: %d, tx hash: %s", l2Block.batch.batchNumber, l2Block.trackingNum, i, tx.HashStr)
		}
	}

	batchL2Data := []byte{}

	// Add changeL2Block to batchL2Data
	changeL2BlockBytes := f.stateIntf.BuildChangeL2Block(l2Block.deltaTimestamp, l2Block.getL1InfoTreeIndex())
	batchL2Data = append(batchL2Data, changeL2BlockBytes...)

	// Add transactions data to batchL2Data
	for _, tx := range l2Block.transactions {
		epHex, err := hex.DecodeHex(fmt.Sprintf("%x", tx.EGPPercentage))
		if err != nil {
			log.Errorf("error decoding hex value for effective gas price percentage for tx %s, error: %v", tx.HashStr, err)
			return nil, 0, "", err
		}

		txData := append(tx.RawTx, epHex...)

		batchL2Data = append(batchL2Data, txData...)
	}

	batchRequest := state.ProcessRequest{
		BatchNumber:               l2Block.batch.batchNumber,
		OldStateRoot:              initialStateRoot,
		Coinbase:                  l2Block.batch.coinbase,
		L1InfoRoot_V2:             state.GetMockL1InfoRoot(),
		TimestampLimit_V2:         l2Block.timestamp,
		Transactions:              batchL2Data,
		SkipFirstChangeL2Block_V2: false,
		SkipWriteBlockInfoRoot_V2: false,
		Caller:                    stateMetrics.DiscardCallerLabel,
		ForkID:                    f.stateIntf.GetForkIDByBatchNumber(l2Block.batch.batchNumber),
		SkipVerifyL1InfoRoot_V2:   true,
		L1InfoTreeData_V2:         map[uint32]state.L1DataV2{},
	}
	batchRequest.L1InfoTreeData_V2[l2Block.l1InfoTreeExitRoot.L1InfoTreeIndex] = state.L1DataV2{
		GlobalExitRoot: l2Block.l1InfoTreeExitRoot.GlobalExitRoot.GlobalExitRoot,
		BlockHashL1:    l2Block.l1InfoTreeExitRoot.PreviousBlockHash,
		MinTimestamp:   uint64(l2Block.l1InfoTreeExitRoot.GlobalExitRoot.Timestamp.Unix()),
	}

	executionStart := time.Now()
	batchResponse, contextId, err := f.stateIntf.ProcessBatchV2(ctx, batchRequest, true)
	l2Block.metrics.l2BlockTimes.executor = time.Since(executionStart)

	if err != nil {
		executeL2BLockError(err)
		return nil, 0, contextId, err
	}

	if batchResponse.ExecutorError != nil {
		executeL2BLockError(batchResponse.ExecutorError)
		return nil, 0, contextId, ErrExecutorError
	}

	if batchResponse.IsRomOOCError {
		executeL2BLockError(batchResponse.RomError_V2)
		return nil, 0, contextId, ErrProcessBatchOOC
	}

	return batchResponse, uint64(len(batchL2Data)), contextId, nil
}

// storeL2Block stores the L2 block in the state and updates the related batch and transactions
func (f *finalizer) storeL2Block(ctx context.Context, l2Block *L2Block) error {
	startStoring := time.Now()

	blockResponse := l2Block.batchResponse.BlockResponses[0]
	log.Infof("storing L2 block %d [%d], batch: %d, deltaTimestamp: %d, timestamp: %d, l1InfoTreeIndex: %d, l1InfoTreeIndexChanged: %v, txs: %d/%d, blockHash: %s, infoRoot: %s",
		blockResponse.BlockNumber, l2Block.trackingNum, l2Block.batch.batchNumber, l2Block.deltaTimestamp, l2Block.timestamp, l2Block.l1InfoTreeExitRoot.L1InfoTreeIndex,
		l2Block.l1InfoTreeExitRootChanged, len(l2Block.transactions), len(blockResponse.TransactionResponses), blockResponse.BlockHash, blockResponse.BlockInfoRoot.String())

	// Wait until L2 block has been flushed/stored by the executor
	startWaitFlushId := time.Now()
	f.storedFlushIDCond.L.Lock()
	for f.storedFlushID < l2Block.batchResponse.FlushID {
		f.storedFlushIDCond.Wait()
	}
	f.storedFlushIDCond.L.Unlock()
	waitFlushId := time.Since(startWaitFlushId)

	// If the L2 block has txs now f.storedFlushID >= l2BlockToStore.flushId, we can store tx
	dbTx, err := f.stateIntf.BeginStateTransaction(ctx)
	if err != nil {
		return fmt.Errorf("error creating db transaction to store L2 block %d [%d], error: %v", blockResponse.BlockNumber, l2Block.trackingNum, err)
	}

	rollbackOnError := func(retError error) error {
		err := dbTx.Rollback(ctx)
		if err != nil {
			return fmt.Errorf("rollback error due to error %v, error: %v", retError, err)
		}
		return retError
	}

	if (f.sipBatch == nil) || (f.sipBatch.batchNumber != l2Block.batch.batchNumber) {
		// We have l2 blocks to store from a new batch, therefore we insert this new batch in the statedb
		// First we need to close the current sipBatch
		if f.sipBatch != nil {
			err := f.closeSIPBatch(ctx, dbTx)
			if err != nil {
				return rollbackOnError(fmt.Errorf("error when closing sip batch %d, initialStateRoot: %s, error: %v", f.sipBatch.batchNumber, f.sipBatch.initialStateRoot, err))
			}
		}
		// We insert new SIP batch in the statedb
		err := f.insertSIPBatch(ctx, l2Block.batch.batchNumber, l2Block.batch.initialStateRoot, dbTx)
		if err != nil {
			return rollbackOnError(fmt.Errorf("error when inserting new sip batch %d, initialStateRoot: %s, error: %v", l2Block.batch.batchNumber, l2Block.batch.initialStateRoot, err))
		}
		f.sipBatch = l2Block.batch
	}

	forkID := f.stateIntf.GetForkIDByBatchNumber(l2Block.batch.batchNumber)

	txsEGPLog := []*state.EffectiveGasPriceLog{}
	for _, tx := range l2Block.transactions {
		egpLog := tx.EGPLog
		txsEGPLog = append(txsEGPLog, &egpLog)
	}

	// Store L2 block in the state
	blockHash, err := f.stateIntf.StoreL2Block(ctx, l2Block.batch.batchNumber, blockResponse, txsEGPLog, dbTx)
	if err != nil {
		return rollbackOnError(fmt.Errorf("database error on storing L2 block %d [%d], error: %v", blockResponse.BlockNumber, l2Block.trackingNum, err))
	}

	// Now we need to update de BatchL2Data of the wip batch and also update the status of the L2 block txs in the pool

	batch, err := f.stateIntf.GetBatchByNumber(ctx, l2Block.batch.batchNumber, dbTx)
	if err != nil {
		return rollbackOnError(fmt.Errorf("error when getting batch %d from the state, error: %v", l2Block.batch.batchNumber, err))
	}

	// Add changeL2Block to batch.BatchL2Data
	blockL2Data := []byte{}
	changeL2BlockBytes := f.stateIntf.BuildChangeL2Block(l2Block.deltaTimestamp, l2Block.getL1InfoTreeIndex())
	blockL2Data = append(blockL2Data, changeL2BlockBytes...)

	// Add transactions data to batch.BatchL2Data
	for _, txResponse := range blockResponse.TransactionResponses {
		txData, err := state.EncodeTransaction(txResponse.Tx, uint8(txResponse.EffectivePercentage), forkID)
		if err != nil {
			return rollbackOnError(fmt.Errorf("error when encoding tx %s, error: %v", txResponse.TxHash.String(), err))
		}
		blockL2Data = append(blockL2Data, txData...)
	}

	batch.BatchL2Data = append(batch.BatchL2Data, blockL2Data...)
	batch.Resources.SumUp(state.BatchResources{ZKCounters: l2Block.batchResponse.UsedZkCounters, Bytes: uint64(len(blockL2Data))})
	batch.HighReservedZKCounters = l2Block.highReservedZKCounters

	receipt := state.ProcessingReceipt{
		BatchNumber:            l2Block.batch.batchNumber,
		StateRoot:              l2Block.batchResponse.NewStateRoot,
		LocalExitRoot:          l2Block.batchResponse.NewLocalExitRoot,
		BatchL2Data:            batch.BatchL2Data,
		BatchResources:         batch.Resources,
		HighReservedZKCounters: batch.HighReservedZKCounters,
	}

	// We need to update the batch GER only in the GER of the block (response) is not zero, since the final GER stored in the batch
	// must be the last GER from the blocks that is not zero (last L1InfoRootIndex change)
	if blockResponse.GlobalExitRoot != state.ZeroHash {
		receipt.GlobalExitRoot = blockResponse.GlobalExitRoot
	} else {
		receipt.GlobalExitRoot = batch.GlobalExitRoot
	}

	err = f.stateIntf.UpdateWIPBatch(ctx, receipt, dbTx)
	if err != nil {
		return rollbackOnError(fmt.Errorf("error when updating wip batch %d, error: %v", l2Block.batch.batchNumber, err))
	}

	err = dbTx.Commit(ctx)
	if err != nil {
		return err
	}

	// Update txs status in the pool
	for _, txResponse := range blockResponse.TransactionResponses {
		// Change Tx status to selected
		err = f.poolIntf.UpdateTxStatus(ctx, txResponse.TxHash, pool.TxStatusSelected, false, nil)
		if err != nil {
			return err
		}
	}

	// Send L2 block to data streamer
	err = f.DSSendL2Block(ctx, l2Block.batch.batchNumber, blockResponse, l2Block.getL1InfoTreeIndex(), l2Block.timestamp, blockHash)
	if err != nil {
		//TODO: we need to halt/rollback the L2 block if we had an error sending to the data streamer?
		log.Errorf("error sending L2 block %d [%d] to data streamer, error: %v", blockResponse.BlockNumber, l2Block.trackingNum, err)
	}

	for _, tx := range l2Block.transactions {
		// Delete the tx from the pending list in the worker (addrQueue)
		f.workerIntf.DeleteTxPendingToStore(tx.Hash, tx.From)
	}

	log.Infof("stored L2 block %d [%d], batch: %d, deltaTimestamp: %d, timestamp: %d, l1InfoTreeIndex: %d, l1InfoTreeIndexChanged: %v, txs: %d/%d, blockHash: %s, infoRoot: %s, time: %v, waitFlushId: %v",
		blockResponse.BlockNumber, l2Block.trackingNum, l2Block.batch.batchNumber, l2Block.deltaTimestamp, l2Block.timestamp, l2Block.l1InfoTreeExitRoot.L1InfoTreeIndex,
		l2Block.l1InfoTreeExitRootChanged, len(l2Block.transactions), len(blockResponse.TransactionResponses), blockResponse.BlockHash, blockResponse.BlockInfoRoot.String(), time.Since(startStoring), waitFlushId)

	return nil
}

// finalizeWIPL2Block closes the wip L2 block and opens a new one
func (f *finalizer) finalizeWIPL2Block(ctx context.Context) {
	log.Debugf("finalizing wip L2 block [%d]", f.wipL2Block.trackingNum)

	prevTimestamp := f.wipL2Block.timestamp
	prevL1InfoTreeIndex := f.wipL2Block.l1InfoTreeExitRoot.L1InfoTreeIndex

	f.closeWIPL2Block(ctx)

	f.openNewWIPL2Block(ctx, prevTimestamp, &prevL1InfoTreeIndex)
}

// closeWIPL2Block closes the wip L2 block
func (f *finalizer) closeWIPL2Block(ctx context.Context) {
	log.Debugf("closing wip L2 block [%d]", f.wipL2Block.trackingNum)

	f.wipBatch.countOfL2Blocks++

	if f.cfg.SequentialProcessL2Block {
		err := f.processL2Block(ctx, f.wipL2Block)
		if err != nil {
			// Dump L2Block info
			f.dumpL2Block(f.wipL2Block)
			f.Halt(ctx, fmt.Errorf("error processing L2 block [%d], error: %v", f.wipL2Block.trackingNum, err), false)
		}
		// We update imStateRoot (used in tx-by-tx execution) to the finalStateRoot that has been updated after process the WIP L2 Block
		f.wipBatch.imStateRoot = f.wipBatch.finalStateRoot
	} else {
		if f.pendingL2BlocksToProcessWG.Count() > 0 {
			startWait := time.Now()
			f.pendingL2BlocksToProcessWG.Wait()
			waitTime := time.Since(startWait)
			log.Debugf("waiting for previous L2 block to be processed took: %v", waitTime)
			f.wipL2Block.metrics.waitl2BlockTime = waitTime
		}

		f.addPendingL2BlockToProcess(ctx, f.wipL2Block)

		f.wipL2Block.metrics.close(f.wipL2Block.createdAt, int64(len(f.wipL2Block.transactions)), f.cfg.SequentialProcessL2Block)

		l2BlockResourcesUsed := state.BatchResources{}
		l2BlockResourcesReserved := state.BatchResources{}

		for _, tx := range f.wipL2Block.transactions {
			l2BlockResourcesUsed.ZKCounters.SumUp(tx.UsedZKCounters)
			l2BlockResourcesReserved.ZKCounters.SumUp(tx.ReservedZKCounters)
		}
		l2BlockResourcesUsed.ZKCounters.SumUp(f.wipL2Block.usedZKCountersOnNew)
		l2BlockResourcesReserved.ZKCounters.SumUp(f.wipL2Block.reservedZKCountersOnNew)

		log.Infof("closed wip L2 block [%d], batch: %d, deltaTimestamp: %d, timestamp: %d, l1InfoTreeIndex: %d, l1InfoTreeIndexChanged: %v, txs: %d, used counters: %s, reserved counters: %s",
			f.wipL2Block.trackingNum, f.wipL2Block.batch.batchNumber, f.wipL2Block.deltaTimestamp, f.wipL2Block.timestamp, f.wipL2Block.l1InfoTreeExitRoot.L1InfoTreeIndex,
			f.wipL2Block.l1InfoTreeExitRootChanged, len(f.wipL2Block.transactions), f.logZKCounters(l2BlockResourcesUsed.ZKCounters), f.logZKCounters(l2BlockResourcesReserved.ZKCounters))

		if f.nextStateRootSync.Before(time.Now()) {
			log.Debug("sync stateroot time reached")
			f.waitPendingL2Blocks()

			// Sanity-check: At this point f.sipBatch should be the same as the batch of the last L2 block processed
			// (only if we haven't had a L2 block reorg just in the last block and it's the first one of the wipBatch)
			if f.wipBatch.batchNumber != f.sipBatch.batchNumber && !(f.l2BlockReorg.Load() && f.wipBatch.countOfL2Blocks <= 2) {
				f.Halt(ctx, fmt.Errorf("wipBatch %d doesn't match sipBatch %d after all pending L2 blocks has been processed/stored", f.wipBatch.batchNumber, f.sipBatch.batchNumber), false)
			}

			f.wipBatch.imStateRoot = f.wipBatch.finalStateRoot
			f.scheduleNextStateRootSync()
			log.Infof("stateroot synced on L2 block [%d] to %s, next sync at %v", f.wipL2Block.trackingNum, f.wipBatch.imStateRoot, f.nextStateRootSync)
		}
	}

	f.wipL2Block = nil
}

// openNewWIPL2Block opens a new wip L2 block
func (f *finalizer) openNewWIPL2Block(ctx context.Context, prevTimestamp uint64, prevL1InfoTreeIndex *uint32) {
	processStart := time.Now()

	newL2Block := &L2Block{}
	now := time.Now()
	newL2Block.createdAt = now
	newL2Block.deltaTimestamp = uint32(uint64(now.Unix()) - prevTimestamp)
	newL2Block.timestamp = prevTimestamp + uint64(newL2Block.deltaTimestamp)

	// Tracking number
	f.l2BlockCounter++
	newL2Block.trackingNum = f.l2BlockCounter

	newL2Block.transactions = []*TxTracker{}

	f.lastL1InfoTreeMux.Lock()
	newL2Block.l1InfoTreeExitRoot = f.lastL1InfoTree
	f.lastL1InfoTreeMux.Unlock()

	// Check if L1InfoTreeIndex has changed, in this case we need to use this index in the changeL2block instead of zero
	// If it's the first wip L2 block after starting sequencer (prevL1InfoTreeIndex == nil) then we retrieve the last GER and we check if it's
	// different from the GER of the current L1InfoTreeIndex (if the GER is different this means that the index also is different)
	if prevL1InfoTreeIndex == nil {
		lastGER, err := f.stateIntf.GetLatestBatchGlobalExitRoot(ctx, nil)
		if err == nil {
			newL2Block.l1InfoTreeExitRootChanged = (newL2Block.l1InfoTreeExitRoot.GlobalExitRoot.GlobalExitRoot != lastGER)
		} else {
			// If we got an error when getting the latest GER then we consider that the index has not changed and it will be updated the next time we have a new L1InfoTreeIndex
			log.Warnf("failed to get the latest CER when initializing the WIP L2 block, assuming L1InfoTreeIndex has not changed, error: %v", err)
		}
	} else {
		newL2Block.l1InfoTreeExitRootChanged = (newL2Block.l1InfoTreeExitRoot.L1InfoTreeIndex != *prevL1InfoTreeIndex)
	}

	f.wipL2Block = newL2Block

	log.Debugf("creating new wip L2 block [%d], batch: %d, deltaTimestamp: %d, timestamp: %d, l1InfoTreeIndex: %d, l1InfoTreeIndexChanged: %v",
		f.wipL2Block.trackingNum, f.wipBatch.batchNumber, f.wipL2Block.deltaTimestamp, f.wipL2Block.timestamp, f.wipL2Block.l1InfoTreeExitRoot.L1InfoTreeIndex, f.wipL2Block.l1InfoTreeExitRootChanged)

	// We process (execute) the new wip L2 block to update the imStateRoot and also get the counters used by the wip l2block
	batchResponse, contextId, err := f.executeNewWIPL2Block(ctx)
	if err != nil {
		f.Halt(ctx, fmt.Errorf("failed to execute new wip L2 block [%d], error: %v ", f.wipL2Block.trackingNum, err), false)
	}

	if len(batchResponse.BlockResponses) != 1 {
		f.Halt(ctx, fmt.Errorf("number of L2 block [%d] responses returned by the executor is %d and must be 1", f.wipL2Block.trackingNum, len(batchResponse.BlockResponses)), false)
	}

	// Update imStateRoot
	oldIMStateRoot := f.wipBatch.imStateRoot
	f.wipL2Block.imStateRoot = batchResponse.NewStateRoot
	f.wipBatch.imStateRoot = f.wipL2Block.imStateRoot

	// Save the resources used/reserved and subtract the ZKCounters reserved by the new WIP L2 block from the WIP batch
	// We need to increase the poseidon hashes to reserve in the batch the hashes needed to write the L1InfoRoot when processing the final L2 Block (SkipWriteBlockInfoRoot_V2=false)
	f.wipL2Block.usedZKCountersOnNew = batchResponse.UsedZkCounters
	f.wipL2Block.usedZKCountersOnNew.PoseidonHashes = (batchResponse.UsedZkCounters.PoseidonHashes * 2) + 2 // nolint:gomnd
	f.wipL2Block.reservedZKCountersOnNew = batchResponse.ReservedZkCounters
	f.wipL2Block.reservedZKCountersOnNew.PoseidonHashes = (batchResponse.ReservedZkCounters.PoseidonHashes * 2) + 2 // nolint:gomnd
	f.wipL2Block.bytes = changeL2BlockSize

	neededZKCounters, newHighZKCounters := getNeededZKCounters(f.wipBatch.imHighReservedZKCounters, f.wipL2Block.usedZKCountersOnNew, f.wipL2Block.reservedZKCountersOnNew)
	subOverflow := false
	fits, overflowResource := f.wipBatch.imRemainingResources.Fits(state.BatchResources{ZKCounters: neededZKCounters, Bytes: f.wipL2Block.bytes})
	if fits {
		subOverflow, overflowResource = f.wipBatch.imRemainingResources.Sub(state.BatchResources{ZKCounters: f.wipL2Block.usedZKCountersOnNew, Bytes: f.wipL2Block.bytes})
		if subOverflow { // Sanity check, this cannot happen as reservedZKCounters should be >= that usedZKCounters
			log.Infof("new wip L2 block [%d] used resources exceeds the remaining batch resources, overflow resource: %s, closing WIP batch and creating new one. counters: {batch: %s, used: %s, reserved: %s, needed: %s, high: %s}",
				f.wipL2Block.trackingNum, overflowResource,
				f.logZKCounters(f.wipBatch.imRemainingResources.ZKCounters), f.logZKCounters(f.wipL2Block.usedZKCountersOnNew), f.logZKCounters(f.wipL2Block.reservedZKCountersOnNew), f.logZKCounters(neededZKCounters), f.logZKCounters(f.wipBatch.imHighReservedZKCounters))
		}

		f.wipBatch.imHighReservedZKCounters = newHighZKCounters
	} else {
		log.Infof("new wip L2 block [%d] reserved resources exceeds the remaining batch resources, overflow resource: %s, closing WIP batch and creating new one. counters: {batch: %s, used: %s, reserved: %s, needed: %s, high: %s}",
			f.wipL2Block.trackingNum, overflowResource,
			f.logZKCounters(f.wipBatch.imRemainingResources.ZKCounters), f.logZKCounters(f.wipL2Block.usedZKCountersOnNew), f.logZKCounters(f.wipL2Block.reservedZKCountersOnNew), f.logZKCounters(neededZKCounters), f.logZKCounters(f.wipBatch.imHighReservedZKCounters))
	}

	// If reserved WIP L2 block resources don't fit in the remaining batch resources (or we got an overflow when trying to subtract the used resources)
	// we close the WIP batch and we create a new one
	if !fits || subOverflow {
		err := f.closeAndOpenNewWIPBatch(ctx, state.ResourceExhaustedClosingReason)
		if err != nil {
			f.Halt(ctx, fmt.Errorf("failed to create new wip batch [%d], error: %v", f.wipL2Block.trackingNum, err), true)
		}
	}

	// We assign the wipBatch as the batch where this wipL2Block belongs
	f.wipL2Block.batch = f.wipBatch

	f.wipL2Block.metrics.newL2BlockTimes.sequencer = time.Since(processStart) - f.wipL2Block.metrics.newL2BlockTimes.executor

	log.Infof("created new wip L2 block [%d], batch: %d, deltaTimestamp: %d, timestamp: %d, l1InfoTreeIndex: %d, l1InfoTreeIndexChanged: %v, oldStateRoot: %s, imStateRoot: %s, counters: {used: %s, reserved: %s, needed: %s, high: %s}, contextId: %s",
		f.wipL2Block.trackingNum, f.wipBatch.batchNumber, f.wipL2Block.deltaTimestamp, f.wipL2Block.timestamp, f.wipL2Block.l1InfoTreeExitRoot.L1InfoTreeIndex, f.wipL2Block.l1InfoTreeExitRootChanged, oldIMStateRoot, f.wipL2Block.imStateRoot,
		f.logZKCounters(f.wipL2Block.usedZKCountersOnNew), f.logZKCounters(f.wipL2Block.usedZKCountersOnNew), f.logZKCounters(f.wipL2Block.reservedZKCountersOnNew), f.logZKCounters(f.wipBatch.imHighReservedZKCounters), contextId)
}

// executeNewWIPL2Block executes an empty L2 Block in the executor and returns the batch response from the executor
func (f *finalizer) executeNewWIPL2Block(ctx context.Context) (*state.ProcessBatchResponse, string, error) {
	batchRequest := state.ProcessRequest{
		BatchNumber:               f.wipBatch.batchNumber,
		OldStateRoot:              f.wipBatch.imStateRoot,
		Coinbase:                  f.wipBatch.coinbase,
		L1InfoRoot_V2:             state.GetMockL1InfoRoot(),
		TimestampLimit_V2:         f.wipL2Block.timestamp,
		Caller:                    stateMetrics.DiscardCallerLabel,
		ForkID:                    f.stateIntf.GetForkIDByBatchNumber(f.wipBatch.batchNumber),
		SkipWriteBlockInfoRoot_V2: true,
		SkipVerifyL1InfoRoot_V2:   true,
		SkipFirstChangeL2Block_V2: false,
		Transactions:              f.stateIntf.BuildChangeL2Block(f.wipL2Block.deltaTimestamp, f.wipL2Block.getL1InfoTreeIndex()),
		L1InfoTreeData_V2:         map[uint32]state.L1DataV2{},
	}

	batchRequest.L1InfoTreeData_V2[f.wipL2Block.l1InfoTreeExitRoot.L1InfoTreeIndex] = state.L1DataV2{
		GlobalExitRoot: f.wipL2Block.l1InfoTreeExitRoot.GlobalExitRoot.GlobalExitRoot,
		BlockHashL1:    f.wipL2Block.l1InfoTreeExitRoot.PreviousBlockHash,
		MinTimestamp:   uint64(f.wipL2Block.l1InfoTreeExitRoot.GlobalExitRoot.Timestamp.Unix()),
	}

	executorTime := time.Now()
	batchResponse, contextId, err := f.stateIntf.ProcessBatchV2(ctx, batchRequest, false)
	f.wipL2Block.metrics.newL2BlockTimes.executor = time.Since(executorTime)

	if err != nil {
		return nil, contextId, err
	}

	if batchResponse.ExecutorError != nil {
		return nil, contextId, ErrExecutorError
	}

	if batchResponse.IsRomOOCError {
		return nil, contextId, ErrProcessBatchOOC
	}

	return batchResponse, contextId, nil
}

func (f *finalizer) scheduleNextStateRootSync() {
	f.nextStateRootSync = time.Now().Add(f.cfg.StateRootSyncInterval.Duration)
}

func (f *finalizer) waitPendingL2Blocks() {
	// Wait until all L2 blocks are processed/discarded
	startWait := time.Now()
	f.pendingL2BlocksToProcessWG.Wait()
	log.Debugf("waiting for pending L2 blocks to be processed took: %v", time.Since(startWait))

	// Wait until all L2 blocks are stored
	startWait = time.Now()
	f.pendingL2BlocksToStoreWG.Wait()
	log.Debugf("waiting for pending L2 blocks to be stored took: %v", time.Since(startWait))
}

func (f *finalizer) dumpL2Block(l2Block *L2Block) {
	var blockResp *state.ProcessBlockResponse
	if l2Block.batchResponse != nil {
		if len(l2Block.batchResponse.BlockResponses) > 0 {
			blockResp = l2Block.batchResponse.BlockResponses[0]
		}
	}

	sLog := ""
	for i, tx := range l2Block.transactions {
		sLog += fmt.Sprintf("  tx[%d] hash: %s, from: %s, nonce: %d, gas: %d, gasPrice: %d, bytes: %d, egpPct: %d, countersOnNew: {used: %s, reserved: %s}\n",
			i, tx.HashStr, tx.FromStr, tx.Nonce, tx.Gas, tx.GasPrice, tx.Bytes, tx.EGPPercentage, f.logZKCounters(tx.UsedZKCounters), f.logZKCounters(tx.ReservedZKCounters))
	}
	log.Infof("dump L2 block [%d], timestamp: %d, deltaTimestamp: %d, imStateRoot: %s, l1InfoTreeIndex: %d, bytes: %d, used counters: %s, reserved counters: %s\n%s",
		l2Block.trackingNum, l2Block.timestamp, l2Block.deltaTimestamp, l2Block.imStateRoot, l2Block.l1InfoTreeExitRoot.L1InfoTreeIndex, l2Block.bytes,
		f.logZKCounters(l2Block.usedZKCountersOnNew), f.logZKCounters(l2Block.reservedZKCountersOnNew), sLog)

	sLog = ""
	if blockResp != nil {
		for i, txResp := range blockResp.TransactionResponses {
			sLog += fmt.Sprintf("  tx[%d] hash: %s, hashL2: %s, stateRoot: %s, type: %d, gasLeft: %d, gasUsed: %d, gasRefund: %d, createAddress: %s, changesStateRoot: %v, egp: %s, egpPct: %d, hasGaspriceOpcode: %v, hasBalanceOpcode: %v\n",
				i, txResp.TxHash, txResp.TxHashL2_V2, txResp.StateRoot, txResp.Type, txResp.GasLeft, txResp.GasUsed, txResp.GasRefunded, txResp.CreateAddress, txResp.ChangesStateRoot, txResp.EffectiveGasPrice,
				txResp.EffectivePercentage, txResp.HasGaspriceOpcode, txResp.HasBalanceOpcode)
		}

		log.Infof("dump L2 block %d [%d] response, timestamp: %d, parentHash: %s, coinbase: %s, ger: %s, blockHashL1: %s, gasUsed: %d, blockInfoRoot: %s, blockHash: %s, counters: {used: %s, reserved: %s}\n%s",
			blockResp.BlockNumber, l2Block.trackingNum, blockResp.Timestamp, blockResp.ParentHash, blockResp.Coinbase, blockResp.GlobalExitRoot, blockResp.BlockHashL1,
			blockResp.GasUsed, blockResp.BlockInfoRoot, blockResp.BlockHash, f.logZKCounters(l2Block.batchResponse.UsedZkCounters), f.logZKCounters(l2Block.batchResponse.ReservedZkCounters), sLog)
	}
}
