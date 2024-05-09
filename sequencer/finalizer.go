package sequencer

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	"github.com/0xPolygonHermez/zkevm-data-streamer/datastreamer"
	ethermanTypes "github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/event"
	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/state"
	stateMetrics "github.com/0xPolygonHermez/zkevm-node/state/metrics"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/ethereum/go-ethereum/common"
)

const (
	pendingL2BlocksBufferSize = 100
	changeL2BlockSize         = 9 //1 byte (tx type = 0B) + 4 bytes for deltaTimestamp + 4 for l1InfoTreeIndex
)

var (
	now = time.Now
)

// finalizer represents the finalizer component of the sequencer.
type finalizer struct {
	cfg              FinalizerCfg
	isSynced         func(ctx context.Context) bool
	sequencerAddress common.Address
	workerIntf       workerInterface
	poolIntf         txPool
	stateIntf        stateInterface
	etherman         ethermanInterface
	wipBatch         *Batch
	wipL2Block       *L2Block
	batchConstraints state.BatchConstraintsCfg
	haltFinalizer    atomic.Bool
	// forced batches
	nextForcedBatches       []state.ForcedBatch
	nextForcedBatchDeadline int64
	nextForcedBatchesMux    *sync.Mutex
	lastForcedBatchNum      uint64
	// L1InfoTree
	lastL1InfoTreeValid bool
	lastL1InfoTree      state.L1InfoTreeExitRootStorageEntry
	lastL1InfoTreeMux   *sync.Mutex
	lastL1InfoTreeCond  *sync.Cond
	// event log
	eventLog *event.EventLog
	// effective gas price calculation instance
	effectiveGasPrice *pool.EffectiveGasPrice
	// pending L2 blocks to process (executor)
	pendingL2BlocksToProcess   chan *L2Block
	pendingL2BlocksToProcessWG *sync.WaitGroup
	// pending L2 blocks to store in the state
	pendingL2BlocksToStore   chan *L2Block
	pendingL2BlocksToStoreWG *sync.WaitGroup
	// L2 block counter for tracking purposes
	l2BlockCounter uint64
	// executor flushid control
	proverID           string
	storedFlushID      uint64
	storedFlushIDCond  *sync.Cond //Condition to wait until storedFlushID has been updated
	lastPendingFlushID uint64
	pendingFlushIDCond *sync.Cond
	// worker ready txs condition
	workerReadyTxsCond *timeoutCond
	// interval metrics
	metrics *intervalMetrics
	// stream server
	streamServer *datastreamer.StreamServer
	dataToStream chan interface{}
}

// newFinalizer returns a new instance of Finalizer.
func newFinalizer(
	cfg FinalizerCfg,
	poolCfg pool.Config,
	workerIntf workerInterface,
	poolIntf txPool,
	stateIntf stateInterface,
	etherman ethermanInterface,
	sequencerAddr common.Address,
	isSynced func(ctx context.Context) bool,
	batchConstraints state.BatchConstraintsCfg,
	eventLog *event.EventLog,
	streamServer *datastreamer.StreamServer,
	workerReadyTxsCond *timeoutCond,
	dataToStream chan interface{},
) *finalizer {
	f := finalizer{
		cfg:              cfg,
		isSynced:         isSynced,
		sequencerAddress: sequencerAddr,
		workerIntf:       workerIntf,
		poolIntf:         poolIntf,
		stateIntf:        stateIntf,
		etherman:         etherman,
		batchConstraints: batchConstraints,
		// forced batches
		nextForcedBatches:       make([]state.ForcedBatch, 0),
		nextForcedBatchDeadline: 0,
		nextForcedBatchesMux:    new(sync.Mutex),
		// L1InfoTree
		lastL1InfoTreeValid: false,
		lastL1InfoTreeMux:   new(sync.Mutex),
		lastL1InfoTreeCond:  sync.NewCond(&sync.Mutex{}),
		// event log
		eventLog: eventLog,
		// effective gas price calculation instance
		effectiveGasPrice: pool.NewEffectiveGasPrice(poolCfg.EffectiveGasPrice),
		// pending L2 blocks to process (executor)
		pendingL2BlocksToProcess:   make(chan *L2Block, pendingL2BlocksBufferSize),
		pendingL2BlocksToProcessWG: new(sync.WaitGroup),
		// pending L2 blocks to store in the state
		pendingL2BlocksToStore:   make(chan *L2Block, pendingL2BlocksBufferSize),
		pendingL2BlocksToStoreWG: new(sync.WaitGroup),
		storedFlushID:            0,
		// executor flushid control
		proverID:           "",
		storedFlushIDCond:  sync.NewCond(&sync.Mutex{}),
		lastPendingFlushID: 0,
		pendingFlushIDCond: sync.NewCond(&sync.Mutex{}),
		// worker ready txs condition
		workerReadyTxsCond: workerReadyTxsCond,
		// metrics
		metrics: newIntervalMetrics(cfg.Metrics.Interval.Duration),
		// stream server
		streamServer: streamServer,
		dataToStream: dataToStream,
	}

	f.haltFinalizer.Store(false)

	return &f
}

// Start starts the finalizer.
func (f *finalizer) Start(ctx context.Context) {
	// Do sanity check for batches closed but pending to be checked
	f.processBatchesPendingtoCheck(ctx)

	// Update L1InfoRoot
	go f.checkL1InfoTreeUpdate(ctx)

	// Get the last batch if still wip or opens a new one
	f.initWIPBatch(ctx)

	// Initializes the wip L2 block
	f.initWIPL2Block(ctx)

	// Update the prover id and flush id
	go f.updateProverIdAndFlushId(ctx)

	// Process L2 Blocks
	go f.processPendingL2Blocks(ctx)

	// Store L2 Blocks
	go f.storePendingL2Blocks(ctx)

	// Foced batches checking
	go f.checkForcedBatches(ctx)

	// Processing transactions and finalizing batches
	f.finalizeBatches(ctx)
}

// updateProverIdAndFlushId updates the prover id and flush id
func (f *finalizer) updateProverIdAndFlushId(ctx context.Context) {
	for {
		f.pendingFlushIDCond.L.Lock()
		// f.storedFlushID is >= than f.lastPendingFlushID, this means all pending txs (flushid) are stored by the executor.
		// We are "synced" with the flush id, therefore we need to wait for new tx (new pending flush id to be stored by the executor)
		for f.storedFlushID >= f.lastPendingFlushID {
			f.pendingFlushIDCond.Wait()
		}
		f.pendingFlushIDCond.L.Unlock()

		for f.storedFlushID < f.lastPendingFlushID { //TODO: review this loop as could be is pulling all the time, no sleep
			storedFlushID, proverID, err := f.stateIntf.GetStoredFlushID(ctx)
			if err != nil {
				log.Errorf("failed to get stored flush id, error: %v", err)
			} else {
				if storedFlushID != f.storedFlushID {
					// Check if prover/Executor has been restarted
					f.checkIfProverRestarted(proverID)

					// Update f.storeFlushID and signal condition f.storedFlushIDCond
					f.storedFlushIDCond.L.Lock()
					f.storedFlushID = storedFlushID
					f.storedFlushIDCond.Broadcast()
					f.storedFlushIDCond.L.Unlock()
				}
			}
		}
	}
}

// updateFlushIDs updates f.lastPendingFLushID and f.storedFlushID with newPendingFlushID and newStoredFlushID values (it they have changed)
// and sends the signals conditions f.pendingFlushIDCond and f.storedFlushIDCond to notify other go funcs that the values have changed
func (f *finalizer) updateFlushIDs(newPendingFlushID, newStoredFlushID uint64) {
	if newPendingFlushID > f.lastPendingFlushID {
		f.lastPendingFlushID = newPendingFlushID
		f.pendingFlushIDCond.Broadcast()
	}

	f.storedFlushIDCond.L.Lock()
	if newStoredFlushID > f.storedFlushID {
		f.storedFlushID = newStoredFlushID
		f.storedFlushIDCond.Broadcast()
	}
	f.storedFlushIDCond.L.Unlock()
}

func (f *finalizer) checkValidL1InfoRoot(ctx context.Context, l1InfoRoot state.L1InfoTreeExitRootStorageEntry) (bool, error) {
	// Check L1 block hash matches
	l1BlockState, err := f.stateIntf.GetBlockByNumber(ctx, l1InfoRoot.BlockNumber, nil)
	if err != nil {
		return false, fmt.Errorf("error getting L1 block %d from the state, error: %v", l1InfoRoot.BlockNumber, err)
	}

	l1BlockEth, err := f.etherman.HeaderByNumber(ctx, new(big.Int).SetUint64(l1InfoRoot.BlockNumber))
	if err != nil {
		return false, fmt.Errorf("error getting L1 block %d from ethereum, error: %v", l1InfoRoot.BlockNumber, err)
	}

	if l1BlockState.BlockHash != l1BlockEth.Hash() {
		warnmsg := fmt.Sprintf("invalid l1InfoRoot %s, index: %d, GER: %s, l1Block: %d. L1 block hash %s doesn't match block hash on ethereum %s (L1 reorg?)",
			l1InfoRoot.L1InfoTreeRoot, l1InfoRoot.L1InfoTreeIndex, l1InfoRoot.GlobalExitRoot.GlobalExitRoot, l1InfoRoot.BlockNumber, l1BlockState.BlockHash, l1BlockEth.Hash())
		log.Warnf(warnmsg)
		f.LogEvent(ctx, event.Level_Critical, event.EventID_InvalidInfoRoot, warnmsg, nil)

		return false, nil
	}

	// Check l1InfoRootIndex and GER matches
	// We retrieve first the info of the last l1InfoTree event in the block
	log.Debugf("getting l1InfoRoot events for L1 block %d, hash: %s", l1InfoRoot.BlockNumber, l1BlockState.BlockHash)
	blocks, eventsOrder, err := f.etherman.GetRollupInfoByBlockRange(ctx, l1InfoRoot.BlockNumber, &l1InfoRoot.BlockNumber)
	if err != nil {
		return false, err
	}

	// Since in the case we have several l1InfoTree events in the same block, we retrieve only the GER of last one and skips the others
	lastGER := state.ZeroHash
	for _, block := range blocks {
		blockEventsOrder := eventsOrder[block.BlockHash]
		for _, order := range blockEventsOrder {
			if order.Name == ethermanTypes.L1InfoTreeOrder {
				lastGER = block.L1InfoTree[order.Pos].GlobalExitRoot
				log.Debugf("l1InfoTree event, pos: %d, GER: %s", order.Pos, lastGER)
			}
		}
	}

	// Get the deposit count in the moment when the L1InfoRoot was synced
	depositCount, err := f.etherman.DepositCount(ctx, &l1InfoRoot.BlockNumber)
	if err != nil {
		return false, err
	}
	// l1InfoTree index starts at 0, therefore we need to subtract 1 to the depositCount to get the last index at that moment
	index := uint32(depositCount.Uint64())
	if index > 0 { // we check this as protection, but depositCount should be greater that 0 in this context
		index--
	} else {
		warnmsg := fmt.Sprintf("invalid l1InfoRoot %s, index: %d, GER: %s, blockNum: %d. DepositCount value returned by the smartcontrat is 0 and that isn't possible in this context",
			l1InfoRoot.L1InfoTreeRoot, l1InfoRoot.L1InfoTreeIndex, l1InfoRoot.GlobalExitRoot.GlobalExitRoot, l1InfoRoot.BlockNumber)
		log.Warn(warnmsg)
		f.LogEvent(ctx, event.Level_Critical, event.EventID_InvalidInfoRoot, warnmsg, nil)

		return false, nil
	}

	log.Debugf("checking valid l1InfoRoot, index: %d, GER: %s, l1Block: %d, scIndex: %d, scGER: %s",
		l1InfoRoot.BlockNumber, l1InfoRoot.L1InfoTreeIndex, l1InfoRoot.GlobalExitRoot.GlobalExitRoot, index, lastGER)

	if (l1InfoRoot.GlobalExitRoot.GlobalExitRoot != lastGER) || (l1InfoRoot.L1InfoTreeIndex != index) {
		warnmsg := fmt.Sprintf("invalid l1InfoRoot %s, index: %d, GER: %s, blockNum: %d. It doesn't match with smartcontract l1InfoRoot, index: %d, GER: %s",
			l1InfoRoot.L1InfoTreeRoot, l1InfoRoot.L1InfoTreeIndex, l1InfoRoot.GlobalExitRoot.GlobalExitRoot, l1InfoRoot.BlockNumber, index, lastGER)
		log.Warn(warnmsg)
		f.LogEvent(ctx, event.Level_Critical, event.EventID_InvalidInfoRoot, warnmsg, nil)

		return false, nil
	}

	return true, nil
}

func (f *finalizer) checkL1InfoTreeUpdate(ctx context.Context) {
	broadcastL1InfoTreeValid := func() {
		if !f.lastL1InfoTreeValid {
			f.lastL1InfoTreeCond.L.Lock()
			f.lastL1InfoTreeValid = true
			f.lastL1InfoTreeCond.Broadcast()
			f.lastL1InfoTreeCond.L.Unlock()
		}
	}

	firstL1InfoRootUpdate := true
	skipFirstSleep := true

	if f.cfg.L1InfoTreeCheckInterval.Duration.Seconds() == 0 { //nolint:gomnd
		broadcastL1InfoTreeValid()
		return
	}

	for {
		if skipFirstSleep {
			skipFirstSleep = false
		} else {
			time.Sleep(f.cfg.L1InfoTreeCheckInterval.Duration)
		}

		lastL1BlockNumber, err := f.etherman.GetLatestBlockNumber(ctx)
		if err != nil {
			log.Errorf("error getting latest L1 block number, error: %v", err)
			continue
		}

		maxBlockNumber := uint64(0)
		if f.cfg.L1InfoTreeL1BlockConfirmations <= lastL1BlockNumber {
			maxBlockNumber = lastL1BlockNumber - f.cfg.L1InfoTreeL1BlockConfirmations
		}

		l1InfoRoot, err := f.stateIntf.GetLatestL1InfoRoot(ctx, maxBlockNumber)
		if err != nil {
			log.Errorf("error getting latest l1InfoRoot, error: %v", err)
			continue
		}

		// L1InfoTreeIndex = 0 is a special case (empty tree) therefore we will set GER as zero
		if l1InfoRoot.L1InfoTreeIndex == 0 {
			l1InfoRoot.GlobalExitRoot.GlobalExitRoot = state.ZeroHash
		}

		if firstL1InfoRootUpdate || l1InfoRoot.L1InfoTreeIndex > f.lastL1InfoTree.L1InfoTreeIndex {
			log.Infof("received new l1InfoRoot %s, index: %d, l1Block: %d", l1InfoRoot.L1InfoTreeRoot, l1InfoRoot.L1InfoTreeIndex, l1InfoRoot.BlockNumber)

			// Check if new l1InfoRoot is valid. We skip it if l1InfoTreeIndex is 0 (it's a special case)
			if l1InfoRoot.L1InfoTreeIndex > 0 {
				valid, err := f.checkValidL1InfoRoot(ctx, l1InfoRoot)
				if err != nil {
					log.Errorf("error validating new l1InfoRoot, index: %d, error: %v", l1InfoRoot.L1InfoTreeIndex, err)
					continue
				}

				if !valid {
					log.Warnf("invalid l1InfoRoot %s, index: %d, l1Block: %d. Stopping syncing l1InfoTreeIndex", l1InfoRoot.L1InfoTreeRoot, l1InfoRoot.L1InfoTreeIndex, l1InfoRoot.BlockNumber)
					return
				}
			}

			firstL1InfoRootUpdate = false

			f.lastL1InfoTreeMux.Lock()
			f.lastL1InfoTree = l1InfoRoot
			f.lastL1InfoTreeMux.Unlock()

			broadcastL1InfoTreeValid()
		}
	}
}

// finalizeBatches runs the endless loop for processing transactions finalizing batches.
func (f *finalizer) finalizeBatches(ctx context.Context) {
	log.Debug("finalizer init loop")
	showNotFoundTxLog := true // used to log debug only the first message when there is no txs to process
	for {
		// We have reached the L2 block time, we need to close the current L2 block and open a new one
		if f.wipL2Block.timestamp+uint64(f.cfg.L2BlockMaxDeltaTimestamp.Seconds()) <= uint64(time.Now().Unix()) {
			f.finalizeWIPL2Block(ctx)
		}

		tx, err := f.workerIntf.GetBestFittingTx(f.wipBatch.imRemainingResources)

		// If we have txs pending to process but none of them fits into the wip batch, we close the wip batch and open a new one
		if err == ErrNoFittingTransaction {
			f.finalizeWIPBatch(ctx, state.NoTxFitsClosingReason)
			continue
		}

		if tx != nil {
			showNotFoundTxLog = true

			firstTxProcess := true

			for {
				_, err := f.processTransaction(ctx, tx, firstTxProcess)
				if err != nil {
					if err == ErrEffectiveGasPriceReprocess {
						firstTxProcess = false
						log.Infof("reprocessing tx %s because of effective gas price calculation", tx.HashStr)
						continue
					} else if err == ErrBatchResourceOverFlow {
						log.Infof("skipping tx %s due to a batch resource overflow", tx.HashStr)
						break
					} else {
						log.Errorf("failed to process tx %s, error: %v", err)
						break
					}
				}
				break
			}
		} else {
			idleTime := time.Now()

			if showNotFoundTxLog {
				log.Debug("no transactions to be processed. Waiting...")
				showNotFoundTxLog = false
			}

			// wait for new ready txs in worker
			f.workerReadyTxsCond.L.Lock()
			f.workerReadyTxsCond.WaitOrTimeout(f.cfg.NewTxsWaitInterval.Duration)
			f.workerReadyTxsCond.L.Unlock()

			// Increase idle time of the WIP L2Block
			f.wipL2Block.metrics.idleTime += time.Since(idleTime)
		}

		if f.haltFinalizer.Load() {
			// There is a fatal error and we need to halt the finalizer and stop processing new txs
			for {
				time.Sleep(5 * time.Second) //nolint:gomnd
			}
		}

		// Check if we must finalize the batch due to a closing reason (resources exhausted, max txs, timestamp resolution, forced batches deadline)
		if finalize, closeReason := f.checkIfFinalizeBatch(); finalize {
			f.finalizeWIPBatch(ctx, closeReason)
		}

		if err := ctx.Err(); err != nil {
			log.Errorf("stopping finalizer because of context, error: %v", err)
			return
		}
	}
}

// processTransaction processes a single transaction.
func (f *finalizer) processTransaction(ctx context.Context, tx *TxTracker, firstTxProcess bool) (errWg *sync.WaitGroup, err error) {
	start := time.Now()

	log.Infof("processing tx %s, batchNumber: %d, l2Block: [%d], oldStateRoot: %s, L1InfoRootIndex: %d",
		tx.HashStr, f.wipBatch.batchNumber, f.wipL2Block.trackingNum, f.wipBatch.imStateRoot, f.wipL2Block.l1InfoTreeExitRoot.L1InfoTreeIndex)

	batchRequest := state.ProcessRequest{
		BatchNumber:               f.wipBatch.batchNumber,
		OldStateRoot:              f.wipBatch.imStateRoot,
		Coinbase:                  f.wipBatch.coinbase,
		L1InfoRoot_V2:             state.GetMockL1InfoRoot(),
		TimestampLimit_V2:         f.wipL2Block.timestamp,
		Caller:                    stateMetrics.DiscardCallerLabel,
		ForkID:                    f.stateIntf.GetForkIDByBatchNumber(f.wipBatch.batchNumber),
		Transactions:              tx.RawTx,
		SkipFirstChangeL2Block_V2: true,
		SkipWriteBlockInfoRoot_V2: true,
		SkipVerifyL1InfoRoot_V2:   true,
		L1InfoTreeData_V2:         map[uint32]state.L1DataV2{},
	}

	txGasPrice := tx.GasPrice

	// If it is the first time we process this tx then we calculate the EffectiveGasPrice
	if firstTxProcess {
		// Get L1 gas price and store in txTracker to make it consistent during the lifespan of the transaction
		tx.L1GasPrice, tx.L2GasPrice = f.poolIntf.GetL1AndL2GasPrice()
		// Get the tx and l2 gas price we will use in the egp calculation. If egp is disabled we will use a "simulated" tx gas price
		txGasPrice, txL2GasPrice := f.effectiveGasPrice.GetTxAndL2GasPrice(tx.GasPrice, tx.L1GasPrice, tx.L2GasPrice)

		// Save values for later logging
		tx.EGPLog.L1GasPrice = tx.L1GasPrice
		tx.EGPLog.L2GasPrice = txL2GasPrice
		tx.EGPLog.GasUsedFirst = tx.UsedZKCounters.GasUsed
		tx.EGPLog.GasPrice.Set(txGasPrice)

		// Calculate EffectiveGasPrice
		egp, err := f.effectiveGasPrice.CalculateEffectiveGasPrice(tx.RawTx, txGasPrice, tx.UsedZKCounters.GasUsed, tx.L1GasPrice, txL2GasPrice)
		if err != nil {
			if f.effectiveGasPrice.IsEnabled() {
				return nil, err
			} else {
				log.Warnf("effectiveGasPrice is disabled, but failed to calculate effectiveGasPrice for tx %s, error: %v", tx.HashStr, err)
				tx.EGPLog.Error = fmt.Sprintf("CalculateEffectiveGasPrice#1: %s", err)
			}
		} else {
			tx.EffectiveGasPrice.Set(egp)

			// Save first EffectiveGasPrice for later logging
			tx.EGPLog.ValueFirst.Set(tx.EffectiveGasPrice)

			// If EffectiveGasPrice >= txGasPrice, we process the tx with tx.GasPrice
			if tx.EffectiveGasPrice.Cmp(txGasPrice) >= 0 {
				loss := new(big.Int).Sub(tx.EffectiveGasPrice, txGasPrice)
				// If loss > 0 the warning message indicating we loss fee for thix tx
				if loss.Cmp(new(big.Int).SetUint64(0)) == 1 {
					log.Warnf("egp-loss: gasPrice: %d, effectiveGasPrice1: %d, loss: %d, tx: %s", txGasPrice, tx.EffectiveGasPrice, loss, tx.HashStr)
				}

				tx.EffectiveGasPrice.Set(txGasPrice)
				tx.IsLastExecution = true
			}
		}
	}

	egpPercentage, err := f.effectiveGasPrice.CalculateEffectiveGasPricePercentage(txGasPrice, tx.EffectiveGasPrice)
	if err != nil {
		if f.effectiveGasPrice.IsEnabled() {
			return nil, err
		} else {
			log.Warnf("effectiveGasPrice is disabled, but failed to to calculate efftive gas price percentage (#1), error: %v", err)
			tx.EGPLog.Error = fmt.Sprintf("%s; CalculateEffectiveGasPricePercentage#1: %s", tx.EGPLog.Error, err)
		}
	} else {
		// Save percentage for later logging
		tx.EGPLog.Percentage = egpPercentage
	}

	// If EGP is disabled we use tx GasPrice (MaxEffectivePercentage=255)
	if !f.effectiveGasPrice.IsEnabled() {
		egpPercentage = state.MaxEffectivePercentage
	}

	// Assign applied EGP percentage to tx (TxTracker)
	tx.EGPPercentage = egpPercentage

	effectivePercentageAsDecodedHex, err := hex.DecodeHex(fmt.Sprintf("%x", tx.EGPPercentage))
	if err != nil {
		return nil, err
	}

	batchRequest.Transactions = append(batchRequest.Transactions, effectivePercentageAsDecodedHex...)

	executionStart := time.Now()
	batchResponse, err := f.stateIntf.ProcessBatchV2(ctx, batchRequest, false)
	executionTime := time.Since(executionStart)
	f.wipL2Block.metrics.transactionsTimes.executor += executionTime

	if err != nil && (errors.Is(err, runtime.ErrExecutorDBError) || errors.Is(err, runtime.ErrInvalidTxChangeL2BlockMinTimestamp)) {
		log.Errorf("failed to process tx %s, error: %v", tx.HashStr, err)
		return nil, err
	} else if err == nil && !batchResponse.IsRomLevelError && len(batchResponse.BlockResponses) == 0 {
		err = fmt.Errorf("executor returned no errors and no responses for tx %s", tx.HashStr)
		f.Halt(ctx, err, false)
	} else if err != nil {
		log.Errorf("error received from executor, error: %v", err)

		// Delete tx from the worker
		f.workerIntf.DeleteTx(tx.Hash, tx.From)

		// Set tx as invalid in the pool
		errMsg := err.Error()
		err = f.poolIntf.UpdateTxStatus(ctx, tx.Hash, pool.TxStatusInvalid, false, &errMsg)
		if err != nil {
			log.Errorf("failed to update status to invalid in the pool for tx %s, error: %v", tx.Hash.String(), err)
		}
		return nil, err
	}

	oldStateRoot := f.wipBatch.imStateRoot
	if len(batchResponse.BlockResponses) > 0 {
		errWg, err = f.handleProcessTransactionResponse(ctx, tx, batchResponse, oldStateRoot)
		if err != nil {
			return errWg, err
		}
	}

	// Update imStateRoot
	f.wipBatch.imStateRoot = batchResponse.NewStateRoot

	log.Infof("processed tx %s, batchNumber: %d, l2Block: [%d], newStateRoot: %s, oldStateRoot: %s, time: {process: %v, executor: %v}, used counters: %s, reserved counters: %s",
		tx.HashStr, batchRequest.BatchNumber, f.wipL2Block.trackingNum, batchResponse.NewStateRoot.String(), batchRequest.OldStateRoot.String(),
		time.Since(start), executionTime, f.logZKCounters(batchResponse.UsedZkCounters), f.logZKCounters(batchResponse.ReservedZkCounters))

	return nil, nil
}

// handleProcessTransactionResponse handles the response of transaction processing.
func (f *finalizer) handleProcessTransactionResponse(ctx context.Context, tx *TxTracker, result *state.ProcessBatchResponse, oldStateRoot common.Hash) (errWg *sync.WaitGroup, err error) {
	txResponse := result.BlockResponses[0].TransactionResponses[0]

	// Update metrics
	f.wipL2Block.metrics.processedTxsCount++

	// Handle Transaction Error
	errorCode := executor.RomErrorCode(txResponse.RomError)
	if !state.IsStateRootChanged(errorCode) {
		// If intrinsic error or OOC error, we skip adding the transaction to the batch
		errWg = f.handleProcessTransactionError(ctx, result, tx)
		return errWg, txResponse.RomError
	}

	egpEnabled := f.effectiveGasPrice.IsEnabled()

	if !tx.IsLastExecution {
		tx.IsLastExecution = true

		// Get the tx gas price we will use in the egp calculation. If egp is disabled we will use a "simulated" tx gas price
		txGasPrice, txL2GasPrice := f.effectiveGasPrice.GetTxAndL2GasPrice(tx.GasPrice, tx.L1GasPrice, tx.L2GasPrice)

		newEffectiveGasPrice, err := f.effectiveGasPrice.CalculateEffectiveGasPrice(tx.RawTx, txGasPrice, txResponse.GasUsed, tx.L1GasPrice, txL2GasPrice)
		if err != nil {
			if egpEnabled {
				log.Errorf("failed to calculate effective gas price with new gasUsed for tx %s, error: %v", tx.HashStr, err.Error())
				return nil, err
			} else {
				log.Warnf("effectiveGasPrice is disabled, but failed to calculate effective gas price with new gasUsed for tx %s, error: %v", tx.HashStr, err.Error())
				tx.EGPLog.Error = fmt.Sprintf("%s; CalculateEffectiveGasPrice#2: %s", tx.EGPLog.Error, err)
			}
		} else {
			// Save new (second) gas used and second effective gas price calculation for later logging
			tx.EGPLog.ValueSecond.Set(newEffectiveGasPrice)
			tx.EGPLog.GasUsedSecond = txResponse.GasUsed

			errCompare := f.compareTxEffectiveGasPrice(ctx, tx, newEffectiveGasPrice, txResponse.HasGaspriceOpcode, txResponse.HasBalanceOpcode)

			// If EffectiveGasPrice is disabled we will calculate the percentage and save it for later logging
			if !egpEnabled {
				effectivePercentage, err := f.effectiveGasPrice.CalculateEffectiveGasPricePercentage(txGasPrice, tx.EffectiveGasPrice)
				if err != nil {
					log.Warnf("effectiveGasPrice is disabled, but failed to calculate effective gas price percentage (#2), error: %v", err)
					tx.EGPLog.Error = fmt.Sprintf("%s, CalculateEffectiveGasPricePercentage#2: %s", tx.EGPLog.Error, err)
				} else {
					// Save percentage for later logging
					tx.EGPLog.Percentage = effectivePercentage
				}
			}

			if errCompare != nil && egpEnabled {
				return nil, errCompare
			}
		}
	}

	// Check if reserved resources of the tx fits in the remaining batch resources
	subOverflow := false
	fits, overflowResource := f.wipBatch.imRemainingResources.Fits(state.BatchResources{ZKCounters: result.ReservedZkCounters, Bytes: uint64(len(tx.RawTx))})
	if fits {
		// Subtract the used resources from the batch
		subOverflow, overflowResource = f.wipBatch.imRemainingResources.Sub(state.BatchResources{ZKCounters: result.UsedZkCounters, Bytes: uint64(len(tx.RawTx))})
		if subOverflow { // Sanity check, this cannot happen as reservedZKCounters should be >= that usedZKCounters
			sLog := fmt.Sprintf("tx %s used resources exceeds the remaining batch resources, overflow resource: %s, updating metadata for tx in worker and continuing. Batch counters: %s, tx used counters: %s",
				tx.HashStr, overflowResource, f.logZKCounters(f.wipBatch.imRemainingResources.ZKCounters), f.logZKCounters(result.UsedZkCounters))

			log.Errorf(sLog)

			f.LogEvent(ctx, event.Level_Error, event.EventID_UsedZKCountersOverflow, sLog, nil)
		}
	} else {
		log.Infof("current tx %s reserved resources exceeds the remaining batch resources, overflow resource: %s, updating metadata for tx in worker and continuing. Batch counters: %s, tx reserved counters: %s",
			tx.HashStr, overflowResource, f.logZKCounters(f.wipBatch.imRemainingResources.ZKCounters), f.logZKCounters(result.ReservedZkCounters))
		if !f.batchConstraints.IsWithinConstraints(result.ReservedZkCounters) {
			log.Infof("current tx %s reserved resources exceeds the max limit for batch resources (node OOC), setting tx as invalid in the pool", tx.HashStr)

			f.LogEvent(ctx, event.Level_Info, event.EventID_NodeOOC,
				fmt.Sprintf("tx %s exceeds node max limit batch resources (node OOC), from: %s, IP: %s", tx.HashStr, tx.FromStr, tx.IP), nil)

			// Delete the transaction from the txSorted list
			f.workerIntf.DeleteTx(tx.Hash, tx.From)

			errMsg := "node OOC"
			err = f.poolIntf.UpdateTxStatus(ctx, tx.Hash, pool.TxStatusInvalid, false, &errMsg)
			if err != nil {
				log.Errorf("failed to update status to invalid in the pool for tx %s, error: %v", tx.Hash.String(), err)
			}

			return nil, ErrBatchResourceOverFlow
		}
	}

	// If reserved tx resources don't fit in the remaining batch resources (or we got an overflow when trying to subtract the used resources)
	// we update the ZKCounters of the tx and returns ErrBatchResourceOverFlow error
	if !fits || subOverflow {
		f.workerIntf.UpdateTxZKCounters(txResponse.TxHash, tx.From, result.UsedZkCounters, result.ReservedZkCounters)
		return nil, ErrBatchResourceOverFlow
	}

	// Save Enabled, GasPriceOC, BalanceOC and final effective gas price for later logging
	tx.EGPLog.Enabled = egpEnabled
	tx.EGPLog.GasPriceOC = txResponse.HasGaspriceOpcode
	tx.EGPLog.BalanceOC = txResponse.HasBalanceOpcode
	tx.EGPLog.ValueFinal.Set(tx.EffectiveGasPrice)

	// Log here the results of EGP calculation
	log.Infof("egp-log: final: %d, first: %d, second: %d, percentage: %d, deviation: %d, maxDeviation: %d, gasUsed1: %d, gasUsed2: %d, gasPrice: %d, l1GasPrice: %d, l2GasPrice: %d, reprocess: %t, gasPriceOC: %t, balanceOC: %t, enabled: %t, txSize: %d, tx: %s, error: %s",
		tx.EGPLog.ValueFinal, tx.EGPLog.ValueFirst, tx.EGPLog.ValueSecond, tx.EGPLog.Percentage, tx.EGPLog.FinalDeviation, tx.EGPLog.MaxDeviation, tx.EGPLog.GasUsedFirst, tx.EGPLog.GasUsedSecond,
		tx.EGPLog.GasPrice, tx.EGPLog.L1GasPrice, tx.EGPLog.L2GasPrice, tx.EGPLog.Reprocess, tx.EGPLog.GasPriceOC, tx.EGPLog.BalanceOC, egpEnabled, len(tx.RawTx), tx.HashStr, tx.EGPLog.Error)

	f.wipL2Block.addTx(tx)

	f.wipBatch.countOfTxs++

	f.updateWorkerAfterSuccessfulProcessing(ctx, tx.Hash, tx.From, false, result)

	// Update metrics
	f.wipL2Block.metrics.gas += txResponse.GasUsed

	return nil, nil
}

// compareTxEffectiveGasPrice compares newEffectiveGasPrice with tx.EffectiveGasPrice.
// It returns ErrEffectiveGasPriceReprocess if the tx needs to be reprocessed with
// the tx.EffectiveGasPrice updated, otherwise it returns nil
func (f *finalizer) compareTxEffectiveGasPrice(ctx context.Context, tx *TxTracker, newEffectiveGasPrice *big.Int, hasGasPriceOC bool, hasBalanceOC bool) error {
	// Get the tx gas price we will use in the egp calculation. If egp is disabled we will use a "simulated" tx gas price
	txGasPrice, _ := f.effectiveGasPrice.GetTxAndL2GasPrice(tx.GasPrice, tx.L1GasPrice, tx.L2GasPrice)

	// Compute the absolute difference between tx.EffectiveGasPrice - newEffectiveGasPrice
	diff := new(big.Int).Abs(new(big.Int).Sub(tx.EffectiveGasPrice, newEffectiveGasPrice))
	// Compute max deviation allowed of newEffectiveGasPrice
	maxDeviation := new(big.Int).Div(new(big.Int).Mul(tx.EffectiveGasPrice, new(big.Int).SetUint64(f.effectiveGasPrice.GetFinalDeviation())), big.NewInt(100)) //nolint:gomnd

	// Save FinalDeviation (diff) and MaxDeviation for later logging
	tx.EGPLog.FinalDeviation.Set(diff)
	tx.EGPLog.MaxDeviation.Set(maxDeviation)

	// if (diff > finalDeviation)
	if diff.Cmp(maxDeviation) == 1 {
		// if newEfectiveGasPrice < txGasPrice
		if newEffectiveGasPrice.Cmp(txGasPrice) == -1 {
			if hasGasPriceOC || hasBalanceOC {
				tx.EffectiveGasPrice.Set(txGasPrice)
			} else {
				tx.EffectiveGasPrice.Set(newEffectiveGasPrice)
			}
		} else {
			tx.EffectiveGasPrice.Set(txGasPrice)

			loss := new(big.Int).Sub(newEffectiveGasPrice, txGasPrice)
			// If loss > 0 the warning message indicating we loss fee for thix tx
			if loss.Cmp(new(big.Int).SetUint64(0)) == 1 {
				log.Warnf("egp-loss: gasPrice: %d, EffectiveGasPrice2: %d, loss: %d, tx: %s", txGasPrice, newEffectiveGasPrice, loss, tx.HashStr)
			}
		}

		// Save Reprocess for later logging
		tx.EGPLog.Reprocess = true

		return ErrEffectiveGasPriceReprocess
	} // else (diff <= finalDeviation) it is ok, no reprocess of the tx is needed

	return nil
}

func (f *finalizer) updateWorkerAfterSuccessfulProcessing(ctx context.Context, txHash common.Hash, txFrom common.Address, isForced bool, result *state.ProcessBatchResponse) {
	// Delete the transaction from the worker
	if isForced {
		f.workerIntf.DeleteForcedTx(txHash, txFrom)
		log.Debugf("forced tx %s deleted from address %s", txHash.String(), txFrom.Hex())
		return
	} else {
		f.workerIntf.DeleteTx(txHash, txFrom)
		log.Debugf("tx %s deleted from address %s", txHash.String(), txFrom.Hex())
	}

	txsToDelete := f.workerIntf.UpdateAfterSingleSuccessfulTxExecution(txFrom, result.ReadWriteAddresses)
	for _, txToDelete := range txsToDelete {
		err := f.poolIntf.UpdateTxStatus(ctx, txToDelete.Hash, pool.TxStatusFailed, false, txToDelete.FailedReason)
		if err != nil {
			log.Errorf("failed to update status to failed in the pool for tx %s, error: %v", txToDelete.Hash.String(), err)
			continue
		}
	}
}

// handleProcessTransactionError handles the error of a transaction
func (f *finalizer) handleProcessTransactionError(ctx context.Context, result *state.ProcessBatchResponse, tx *TxTracker) *sync.WaitGroup {
	txResponse := result.BlockResponses[0].TransactionResponses[0]
	errorCode := executor.RomErrorCode(txResponse.RomError)
	addressInfo := result.ReadWriteAddresses[tx.From]
	log.Infof("rom error in tx %s, errorCode: %d", tx.HashStr, errorCode)
	wg := new(sync.WaitGroup)
	failedReason := executor.RomErr(errorCode).Error()
	if executor.IsROMOutOfCountersError(errorCode) {
		log.Errorf("ROM out of counters error, marking tx %s as invalid, errorCode: %d", tx.HashStr, errorCode)

		f.workerIntf.DeleteTx(tx.Hash, tx.From)

		wg.Add(1)
		go func() {
			defer wg.Done()
			err := f.poolIntf.UpdateTxStatus(ctx, tx.Hash, pool.TxStatusInvalid, false, &failedReason)
			if err != nil {
				log.Errorf("failed to update status to invalid in the pool for tx %s, error: %v", tx.HashStr, err)
			}
		}()
	} else if executor.IsInvalidNonceError(errorCode) || executor.IsInvalidBalanceError(errorCode) {
		var (
			nonce   *uint64
			balance *big.Int
		)
		if addressInfo != nil {
			nonce = addressInfo.Nonce
			balance = addressInfo.Balance
		}
		log.Errorf("intrinsic error, moving tx %s to not ready: nonce: %d, balance: %d. gasPrice: %d, error: %v", tx.Hash, nonce, balance, tx.GasPrice, txResponse.RomError)
		txsToDelete := f.workerIntf.MoveTxToNotReady(tx.Hash, tx.From, nonce, balance)
		for _, txToDelete := range txsToDelete {
			wg.Add(1)
			txToDelete := txToDelete
			go func() {
				defer wg.Done()
				err := f.poolIntf.UpdateTxStatus(ctx, txToDelete.Hash, pool.TxStatusFailed, false, &failedReason)
				if err != nil {
					log.Errorf("failed to update status to failed in the pool for tx %s, error: %v", txToDelete.Hash.String(), err)
				}
			}()
		}
	} else {
		// Delete the transaction from the txSorted list
		f.workerIntf.DeleteTx(tx.Hash, tx.From)
		log.Debugf("tx %s deleted from txSorted list", tx.HashStr)

		wg.Add(1)
		go func() {
			defer wg.Done()
			// Update the status of the transaction to failed
			err := f.poolIntf.UpdateTxStatus(ctx, tx.Hash, pool.TxStatusFailed, false, &failedReason)
			if err != nil {
				log.Errorf("failed to update status to failed in the pool for tx %s, error: %v", tx.Hash.String(), err)
			}
		}()
	}

	// Update metrics
	f.wipL2Block.metrics.gas += txResponse.GasUsed

	return wg
}

// checkIfProverRestarted checks if the proverID changed
func (f *finalizer) checkIfProverRestarted(proverID string) {
	if f.proverID != "" && f.proverID != proverID {
		f.LogEvent(context.Background(), event.Level_Critical, event.EventID_FinalizerRestart,
			fmt.Sprintf("proverID changed from %s to %s, restarting sequencer to discard current WIP batch and work with new executor", f.proverID, proverID), nil)

		log.Fatal("proverID changed from %s to %s, restarting sequencer to discard current WIP batch and work with new executor")
	}
}

// logZKCounters returns a string with all the zkCounters values
func (f *finalizer) logZKCounters(counters state.ZKCounters) string {
	return fmt.Sprintf("{gasUsed: %d, keccakHashes: %d, poseidonHashes: %d, poseidonPaddings: %d, memAligns: %d, arithmetics: %d, binaries: %d, sha256Hashes: %d, steps: %d}",
		counters.GasUsed, counters.KeccakHashes, counters.PoseidonHashes, counters.PoseidonPaddings, counters.MemAligns, counters.Arithmetics,
		counters.Binaries, counters.Sha256Hashes_V2, counters.Steps)
}

// Halt halts the finalizer
func (f *finalizer) Halt(ctx context.Context, err error, isFatal bool) {
	f.haltFinalizer.Store(true)

	f.LogEvent(ctx, event.Level_Critical, event.EventID_FinalizerHalt, fmt.Sprintf("finalizer halted due to error, error: %s", err), nil)

	if isFatal {
		log.Fatalf("fatal error on finalizer, error: %v", err)
	} else {
		for {
			log.Errorf("halting finalizer, error: %v", err)
			time.Sleep(5 * time.Second) //nolint:gomnd
		}
	}
}

// LogEvent adds an event for runtime debugging
func (f *finalizer) LogEvent(ctx context.Context, level event.Level, eventId event.EventID, description string, json interface{}) {
	event := &event.Event{
		ReceivedAt:  time.Now(),
		Source:      event.Source_Node,
		Component:   event.Component_Sequencer,
		Level:       level,
		EventID:     eventId,
		Description: description,
	}

	if json != nil {
		event.Json = json
	}

	eventErr := f.eventLog.LogEvent(ctx, event)
	if eventErr != nil {
		log.Errorf("error storing log event, error: %v", eventErr)
	}
}
