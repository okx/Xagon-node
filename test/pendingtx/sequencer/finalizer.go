package sequencer

import (
	"context"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/sequencer"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"math"
	"math/big"
)

// finalizer represents the finalizer component of the sequencer.
type finalizer struct {
	cfg                sequencer.FinalizerCfg
	batchConstraints   state.BatchConstraintsCfg
	workerIntf         workerInterface
	poolIntf           txPool
	workerReadyTxsCond *timeoutCond

	speedLevelTxCount []int64
	txCountPerLevel   int64

	finishedCh chan int
}

// newFinalizer returns a new instance of Finalizer.
func newFinalizer(
	cfg sequencer.FinalizerCfg,
	batchConstraints state.BatchConstraintsCfg,
	workerIntf workerInterface,
	poolIntf txPool,
	workerReadyTxsCond *timeoutCond,
	levelCount, txCountPerLevel int,
	finishedCh chan int,
) *finalizer {
	f := finalizer{
		cfg:                cfg,
		batchConstraints:   batchConstraints,
		workerIntf:         workerIntf,
		poolIntf:           poolIntf,
		workerReadyTxsCond: workerReadyTxsCond,
		speedLevelTxCount:  make([]int64, levelCount),
		txCountPerLevel:    int64(txCountPerLevel),
		finishedCh:         finishedCh,
	}

	return &f
}

// Start starts the finalizer.
func (f *finalizer) Start(ctx context.Context) {
	// Processing transactions and finalizing batches
	f.finalizeBatches(ctx)
}

// finalizeBatches runs the endless loop for processing transactions finalizing batches.
func (f *finalizer) finalizeBatches(ctx context.Context) {
	log.Debug("finalizer init loop")
	showNotFoundTxLog := true // used to log debug only the first message when there is no txs to process
	maxRemainingResources := getMaxRemainingResources(f.batchConstraints)
	processedTxCount := 0

	for {
		tx, err := f.workerIntf.GetBestFittingTx(maxRemainingResources)

		// If we have txs pending to process but none of them fits into the wip batch, we close the wip batch and open a new one
		if err == sequencer.ErrNoFittingTransaction {
			continue
		}

		if tx != nil {
			log.Debugf("processing tx %s", tx.HashStr)
			showNotFoundTxLog = true

			// ignore processTransaction
			f.updateWorkerAfterSuccessfulProcessing(ctx, tx.Hash, tx.From, tx.Nonce, false)

			level := f.poolIntf.GetLevelByAddr(tx.To)
			f.speedLevelTxCount[level]++
			if f.speedLevelTxCount[level] == f.txCountPerLevel {
				//fmt.Println("===================")
				//fmt.Printf("Speed level [%d] finished.\nAll level: %v\n", level, f.speedLevelTxCount)
				//fmt.Println("===================")
				f.finishedCh <- level
			}

			processedTxCount++
		} else {
			if showNotFoundTxLog {
				log.Debug("no transactions to be processed. Waiting...")
				showNotFoundTxLog = false
			}

			// wait for new ready txs in worker
			f.workerReadyTxsCond.L.Lock()
			f.workerReadyTxsCond.WaitOrTimeout(f.cfg.NewTxsWaitInterval.Duration)
			f.workerReadyTxsCond.L.Unlock()
		}

		if err := ctx.Err(); err != nil {
			log.Errorf("stopping finalizer because of context, error: %v", err)
			return
		}
	}
}

func (f *finalizer) updateWorkerAfterSuccessfulProcessing(ctx context.Context, txHash common.Hash, txFrom common.Address, nonce uint64, isForced bool) {
	// Delete the transaction from the worker
	if isForced {
		f.workerIntf.DeleteForcedTx(txHash, txFrom)
		log.Debugf("forced tx %s deleted from address %s", txHash.String(), txFrom.Hex())
		return
	} else {
		f.workerIntf.DeleteTx(txHash, txFrom)
		log.Debugf("tx %s deleted from address %s", txHash.String(), txFrom.Hex())
	}

	touch := make(map[common.Address]*state.InfoReadWrite)
	newNonce := nonce + 1
	touch[txFrom] = &state.InfoReadWrite{Address: txFrom, Nonce: &newNonce, Balance: big.NewInt(math.MaxInt64)}
	txsToDelete := f.workerIntf.UpdateAfterSingleSuccessfulTxExecution(txFrom, touch)
	for _, txToDelete := range txsToDelete {
		err := f.poolIntf.UpdateTxStatus(ctx, txToDelete.Hash, pool.TxStatusFailed, false, txToDelete.FailedReason)
		if err != nil {
			log.Errorf("failed to update status to failed in the pool for tx %s, error: %v", txToDelete.Hash.String(), err)
			continue
		}
	}

	err := f.poolIntf.UpdateTxStatus(ctx, txHash, pool.TxStatusSelected, false, nil)
	if err != nil {
		log.Errorf("failed to update status to selected in the pool for tx %s, error: %v", txHash.String(), err)
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
