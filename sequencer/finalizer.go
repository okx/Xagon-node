package sequencer

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/0xPolygonHermez/zkevm-data-streamer/datastreamer"
	"github.com/0xPolygonHermez/zkevm-node/event"
	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/sequencer/metrics"
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
	now            = time.Now
	mockL1InfoRoot = common.Hash{}
)

// finalizer represents the finalizer component of the sequencer.
type finalizer struct {
	cfg              FinalizerCfg
	isSynced         func(ctx context.Context) bool
	sequencerAddress common.Address
	workerIntf       workerInterface
	poolIntf         txPool
	stateIntf        stateInterface
	etherman         etherman
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
	etherman etherman,
	sequencerAddr common.Address,
	isSynced func(ctx context.Context) bool,
	batchConstraints state.BatchConstraintsCfg,
	eventLog *event.EventLog,
	streamServer *datastreamer.StreamServer,
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
		// stream server
		streamServer: streamServer,
		dataToStream: dataToStream,
	}

	f.haltFinalizer.Store(false)

	return &f
}

// Start starts the finalizer.
func (f *finalizer) Start(ctx context.Context) {
	// Init mockL1InfoRoot to a mock value since it must be different to {0,0,...,0}
	for i := 0; i < len(mockL1InfoRoot); i++ {
		mockL1InfoRoot[i] = byte(i)
	}

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

func (f *finalizer) checkL1InfoTreeUpdate(ctx context.Context) {
	firstL1InfoRootUpdate := true

	for {
		lastL1BlockNumber, err := f.etherman.GetLatestBlockNumber(ctx)
		if err != nil {
			log.Errorf("error getting latest L1 block number, error: %v", err)
		}

		maxBlockNumber := uint64(0)
		if f.cfg.L1InfoTreeL1BlockConfirmations <= lastL1BlockNumber {
			maxBlockNumber = lastL1BlockNumber - f.cfg.L1InfoTreeL1BlockConfirmations
		}

		l1InfoRoot, err := f.stateIntf.GetLatestL1InfoRoot(ctx, maxBlockNumber)
		if err != nil {
			log.Errorf("error checking latest L1InfoRoot, error: %v", err)
			continue
		}

		// L1InfoTreeIndex = 0 is a special case (empty tree) therefore we will set GER as zero
		if l1InfoRoot.L1InfoTreeIndex == 0 {
			l1InfoRoot.GlobalExitRoot.GlobalExitRoot = state.ZeroHash
		}

		if firstL1InfoRootUpdate || l1InfoRoot.L1InfoTreeIndex > f.lastL1InfoTree.L1InfoTreeIndex {
			firstL1InfoRootUpdate = false

			log.Debugf("received new L1InfoRoot. L1InfoTreeIndex: %d", l1InfoRoot.L1InfoTreeIndex)

			f.lastL1InfoTreeMux.Lock()
			f.lastL1InfoTree = l1InfoRoot
			f.lastL1InfoTreeMux.Unlock()

			if !f.lastL1InfoTreeValid {
				f.lastL1InfoTreeCond.L.Lock()
				f.lastL1InfoTreeValid = true
				f.lastL1InfoTreeCond.Broadcast()
				f.lastL1InfoTreeCond.L.Unlock()
			}
		}

		time.Sleep(f.cfg.L1InfoTreeCheckInterval.Duration)
	}
}

// finalizeBatches runs the endless loop for processing transactions finalizing batches.
func (f *finalizer) finalizeBatches(ctx context.Context) {
	log.Debug("finalizer init loop")
	showNotFoundTxLog := true // used to log debug only the first message when there is no txs to process
	for {
		start := now()
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

		// X1 handle
		f.tryToSleep()

		metrics.WorkerProcessingTime(time.Since(start))
		metrics.GetLogStatistics().CumulativeTiming(metrics.GetTx, time.Since(start))

		if tx != nil {
			metrics.GetLogStatistics().CumulativeCounting(metrics.TxCounter)
			log.Debugf("processing tx %s", tx.HashStr)
			showNotFoundTxLog = true

			firstTxProcess := true

			for {
				_, err := f.processTransaction(ctx, tx, firstTxProcess)
				if err != nil {
					if err == ErrEffectiveGasPriceReprocess {
						firstTxProcess = false
						metrics.GetLogStatistics().CumulativeCounting(metrics.ReprocessingTxCounter)
						log.Infof("reprocessing tx %s because of effective gas price calculation", tx.HashStr)
						continue
					} else if err == ErrBatchResourceUnderFlow {
						log.Infof("skipping tx %s due to a batch resource underflow", tx.HashStr)
						break
					} else {
						metrics.GetLogStatistics().CumulativeCounting(metrics.FailTxCounter)
						log.Errorf("failed to process tx %s, error: %v", err)
						break
					}
				}
				metrics.GetLogStatistics().CumulativeValue(metrics.BatchGas, int64(tx.Gas))
				break
			}
		} else {
			// wait for new txs
			if showNotFoundTxLog {
				log.Debug("no transactions to be processed. Waiting...")
				showNotFoundTxLog = false
			}
			if f.cfg.NewTxsWaitInterval.Duration > 0 {
				time.Sleep(f.cfg.NewTxsWaitInterval.Duration)
				metrics.GetLogStatistics().CumulativeCounting(metrics.GetTxPauseCounter)
			}
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
	metrics.GetLogStatistics().SetTag(metrics.FinalizeBatchNumber, strconv.Itoa(int(f.wipBatch.batchNumber)))
	defer func() {
		metrics.ProcessingTime(time.Since(start))
		metrics.GetLogStatistics().CumulativeTiming(metrics.FinalizeBatchTiming, time.Since(start))
	}()

	log.Infof("processing tx %s, batchNumber: %d, l2Block: [%d], oldStateRoot: %s, L1InfoRootIndex: %d",
		tx.HashStr, f.wipBatch.batchNumber, f.wipL2Block.trackingNum, f.wipBatch.imStateRoot, f.wipL2Block.l1InfoTreeExitRoot.L1InfoTreeIndex)

	batchRequest := state.ProcessRequest{
		BatchNumber:               f.wipBatch.batchNumber,
		OldStateRoot:              f.wipBatch.imStateRoot,
		Coinbase:                  f.wipBatch.coinbase,
		L1InfoRoot_V2:             mockL1InfoRoot,
		TimestampLimit_V2:         f.wipL2Block.timestamp,
		Caller:                    stateMetrics.SequencerCallerLabel,
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
		tx.EGPLog.GasUsedFirst = tx.BatchResources.ZKCounters.GasUsed
		tx.EGPLog.GasPrice.Set(txGasPrice)

		// Calculate EffectiveGasPrice
		egp, err := f.effectiveGasPrice.CalculateEffectiveGasPrice(tx.RawTx, txGasPrice, tx.BatchResources.ZKCounters.GasUsed, tx.L1GasPrice, txL2GasPrice)
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

	batchResponse, err := f.stateIntf.ProcessBatchV2(ctx, batchRequest, false)

	if err != nil && (errors.Is(err, runtime.ErrExecutorDBError) || errors.Is(err, runtime.ErrInvalidTxChangeL2BlockMinTimestamp)) {
		log.Errorf("failed to process tx %s, error: %v", tx.HashStr, err)
		return nil, err
	} else if err == nil && !batchResponse.IsRomLevelError && len(batchResponse.BlockResponses) == 0 {
		err = fmt.Errorf("executor returned no errors and no responses for tx %s", tx.HashStr)
		f.Halt(ctx, err, false)
	} else if batchResponse.IsExecutorLevelError {
		log.Errorf("error received from executor, error: %v", err)
		// Delete tx from the worker
		f.workerIntf.DeleteTx(tx.Hash, tx.From)

		// Set tx as invalid in the pool
		errMsg := batchResponse.ExecutorError.Error()
		err = f.poolIntf.UpdateTxStatus(ctx, tx.Hash, pool.TxStatusInvalid, false, &errMsg)
		if err != nil {
			log.Errorf("failed to update status to invalid in the pool for tx %s, error: %v", tx.Hash.String(), err)
		} else {
			metrics.TxProcessed(metrics.TxProcessedLabelInvalid, 1)
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

	log.Infof("processed tx %s, batchNumber: %d, l2Block: [%d], newStateRoot: %s, oldStateRoot: %s, used counters: %s",
		tx.HashStr, batchRequest.BatchNumber, f.wipL2Block.trackingNum, batchResponse.NewStateRoot.String(), batchRequest.OldStateRoot.String(), f.logZKCounters(batchResponse.UsedZkCounters))

	return nil, nil
}

// handleProcessTransactionResponse handles the response of transaction processing.
func (f *finalizer) handleProcessTransactionResponse(ctx context.Context, tx *TxTracker, result *state.ProcessBatchResponse, oldStateRoot common.Hash) (errWg *sync.WaitGroup, err error) {
	// Handle Transaction Error
	errorCode := executor.RomErrorCode(result.BlockResponses[0].TransactionResponses[0].RomError)
	if !state.IsStateRootChanged(errorCode) {
		// If intrinsic error or OOC error, we skip adding the transaction to the batch
		errWg = f.handleProcessTransactionError(ctx, result, tx)
		return errWg, result.BlockResponses[0].TransactionResponses[0].RomError
	}

	egpEnabled := f.effectiveGasPrice.IsEnabled()

	if !tx.IsLastExecution {
		tx.IsLastExecution = true

		// Get the tx gas price we will use in the egp calculation. If egp is disabled we will use a "simulated" tx gas price
		txGasPrice, txL2GasPrice := f.effectiveGasPrice.GetTxAndL2GasPrice(tx.GasPrice, tx.L1GasPrice, tx.L2GasPrice)

		newEffectiveGasPrice, err := f.effectiveGasPrice.CalculateEffectiveGasPrice(tx.RawTx, txGasPrice, result.BlockResponses[0].TransactionResponses[0].GasUsed, tx.L1GasPrice, txL2GasPrice)
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
			tx.EGPLog.GasUsedSecond = result.BlockResponses[0].TransactionResponses[0].GasUsed

			errCompare := f.compareTxEffectiveGasPrice(ctx, tx, newEffectiveGasPrice, result.BlockResponses[0].TransactionResponses[0].HasGaspriceOpcode, result.BlockResponses[0].TransactionResponses[0].HasBalanceOpcode)

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

	// Check remaining resources
	overflow, overflowResource := f.wipBatch.imRemainingResources.Sub(state.BatchResources{ZKCounters: result.UsedZkCounters, Bytes: uint64(len(tx.RawTx))})
	if overflow {
		log.Infof("current tx %s exceeds the remaining batch resources, overflow resource: %s, updating metadata for tx in worker and continuing", tx.HashStr, overflowResource)
		start := time.Now()
		f.workerIntf.UpdateTxZKCounters(result.BlockResponses[0].TransactionResponses[0].TxHash, tx.From, result.UsedZkCounters)
		metrics.WorkerProcessingTime(time.Since(start))
		return nil, ErrBatchResourceUnderFlow
	}

	// Save Enabled, GasPriceOC, BalanceOC and final effective gas price for later logging
	tx.EGPLog.Enabled = egpEnabled
	tx.EGPLog.GasPriceOC = result.BlockResponses[0].TransactionResponses[0].HasGaspriceOpcode
	tx.EGPLog.BalanceOC = result.BlockResponses[0].TransactionResponses[0].HasBalanceOpcode
	tx.EGPLog.ValueFinal.Set(tx.EffectiveGasPrice)

	// Log here the results of EGP calculation
	log.Infof("egp-log: final: %d, first: %d, second: %d, percentage: %d, deviation: %d, maxDeviation: %d, gasUsed1: %d, gasUsed2: %d, gasPrice: %d, l1GasPrice: %d, l2GasPrice: %d, reprocess: %t, gasPriceOC: %t, balanceOC: %t, enabled: %t, txSize: %d, tx: %s, error: %s",
		tx.EGPLog.ValueFinal, tx.EGPLog.ValueFirst, tx.EGPLog.ValueSecond, tx.EGPLog.Percentage, tx.EGPLog.FinalDeviation, tx.EGPLog.MaxDeviation, tx.EGPLog.GasUsedFirst, tx.EGPLog.GasUsedSecond,
		tx.EGPLog.GasPrice, tx.EGPLog.L1GasPrice, tx.EGPLog.L2GasPrice, tx.EGPLog.Reprocess, tx.EGPLog.GasPriceOC, tx.EGPLog.BalanceOC, egpEnabled, len(tx.RawTx), tx.HashStr, tx.EGPLog.Error)

	f.wipL2Block.addTx(tx)

	f.wipBatch.countOfTxs++

	f.updateWorkerAfterSuccessfulProcessing(ctx, tx.Hash, tx.From, false, result)

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

	start := time.Now()
	txsToDelete := f.workerIntf.UpdateAfterSingleSuccessfulTxExecution(txFrom, result.ReadWriteAddresses)
	for _, txToDelete := range txsToDelete {
		err := f.poolIntf.UpdateTxStatus(ctx, txToDelete.Hash, pool.TxStatusFailed, false, txToDelete.FailedReason)
		if err != nil {
			log.Errorf("failed to update status to failed in the pool for tx %s, error: %v", txToDelete.Hash.String(), err)
			continue
		}
		metrics.TxProcessed(metrics.TxProcessedLabelFailed, 1)
	}
	metrics.WorkerProcessingTime(time.Since(start))
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
		start := time.Now()
		f.workerIntf.DeleteTx(tx.Hash, tx.From)
		metrics.WorkerProcessingTime(time.Since(start))

		wg.Add(1)
		go func() {
			defer wg.Done()
			err := f.poolIntf.UpdateTxStatus(ctx, tx.Hash, pool.TxStatusInvalid, false, &failedReason)
			if err != nil {
				log.Errorf("failed to update status to invalid in the pool for tx %s, error: %v", tx.HashStr, err)
			} else {
				metrics.TxProcessed(metrics.TxProcessedLabelInvalid, 1)
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
		start := time.Now()
		log.Errorf("intrinsic error, moving tx %s to not ready: nonce: %d, balance: %d. gasPrice: %d, error: %v", tx.Hash, nonce, balance, tx.GasPrice, txResponse.RomError)
		txsToDelete := f.workerIntf.MoveTxToNotReady(tx.Hash, tx.From, nonce, balance)
		for _, txToDelete := range txsToDelete {
			wg.Add(1)
			txToDelete := txToDelete
			go func() {
				defer wg.Done()
				err := f.poolIntf.UpdateTxStatus(ctx, txToDelete.Hash, pool.TxStatusFailed, false, &failedReason)
				metrics.TxProcessed(metrics.TxProcessedLabelFailed, 1)
				if err != nil {
					log.Errorf("failed to update status to failed in the pool for tx %s, error: %v", txToDelete.Hash.String(), err)
				}
			}()
		}
		metrics.WorkerProcessingTime(time.Since(start))
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
			} else {
				metrics.TxProcessed(metrics.TxProcessedLabelFailed, 1)
			}
		}()
	}

	return wg
}

// checkIfProverRestarted checks if the proverID changed
func (f *finalizer) checkIfProverRestarted(proverID string) {
	if f.proverID != "" && f.proverID != proverID {
		event := &event.Event{
			ReceivedAt:  time.Now(),
			Source:      event.Source_Node,
			Component:   event.Component_Sequencer,
			Level:       event.Level_Critical,
			EventID:     event.EventID_FinalizerRestart,
			Description: fmt.Sprintf("proverID changed from %s to %s, restarting sequencer to discard current WIP batch and work with new executor", f.proverID, proverID),
		}

		err := f.eventLog.LogEvent(context.Background(), event)
		if err != nil {
			log.Errorf("error storing payload, error: %v", err)
		}

		log.Fatal("proverID changed from %s to %s, restarting sequencer to discard current WIP batch and work with new executor")
	}
}

// logZKCounters returns a string with all the zkCounters values
func (f *finalizer) logZKCounters(counters state.ZKCounters) string {
	return fmt.Sprintf("{gasUsed: %d, keccakHashes: %d, poseidonHashes: %d, poseidonPaddings: %d, memAligns: %d, arithmetics: %d, binaries: %d, sha256Hashes: %d, steps: %d}",
		counters.GasUsed, counters.UsedKeccakHashes, counters.UsedPoseidonHashes, counters.UsedPoseidonPaddings, counters.UsedMemAligns, counters.UsedArithmetics,
		counters.UsedBinaries, counters.UsedSha256Hashes_V2, counters.UsedSteps)
}

// Halt halts the finalizer
func (f *finalizer) Halt(ctx context.Context, err error, isFatal bool) {
	f.haltFinalizer.Store(true)

	event := &event.Event{
		ReceivedAt:  time.Now(),
		Source:      event.Source_Node,
		Component:   event.Component_Sequencer,
		Level:       event.Level_Critical,
		EventID:     event.EventID_FinalizerHalt,
		Description: fmt.Sprintf("finalizer halted due to error, error: %s", err),
	}

	eventErr := f.eventLog.LogEvent(ctx, event)
	if eventErr != nil {
		log.Errorf("error storing finalizer halt event, error: %v", eventErr)
	}

	if isFatal {
		log.Fatalf("fatal error on finalizer, error: %v", err)
	} else {
		for {
			log.Errorf("halting finalizer, error: %v", err)
			time.Sleep(5 * time.Second) //nolint:gomnd
		}
	}
}
