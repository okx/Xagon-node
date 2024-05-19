package sequencer

import (
	"context"
	"sync"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/sequencer"
	"github.com/0xPolygonHermez/zkevm-node/state"
)

// Sequencer represents a sequencer
type Sequencer struct {
	cfg      sequencer.Config
	batchCfg state.BatchConfig

	pool      txPool
	worker    *Worker
	finalizer *finalizer

	workerReadyTxsCond *timeoutCond
}

// New init sequencer
func New(cfg sequencer.Config, batchCfg state.BatchConfig, txPool txPool) (*Sequencer, error) {
	return &Sequencer{
		cfg:      cfg,
		batchCfg: batchCfg,
		pool:     txPool,
	}, nil
}

// Start starts the sequencer
func (s *Sequencer) Start(ctx context.Context, levelCount, txCountPerLevel int, finishedCh chan int) {
	err := s.pool.MarkWIPTxsAsPending(ctx)
	if err != nil {
		log.Fatalf("failed to mark WIP txs as pending, error: %v", err)
	}

	s.workerReadyTxsCond = newTimeoutCond(&sync.Mutex{})
	s.worker = NewWorker(s.batchCfg.Constraints, s.workerReadyTxsCond)
	s.finalizer = newFinalizer(s.cfg.Finalizer, s.batchCfg.Constraints, s.worker, s.pool, s.workerReadyTxsCond, levelCount, txCountPerLevel, finishedCh)

	go s.loadFromPool(ctx)
	go s.finalizer.Start(ctx)
	go s.expireOldWorkerTxs(ctx)

	// Wait until context is done
	<-ctx.Done()
}

func (s *Sequencer) expireOldWorkerTxs(ctx context.Context) {
	for {
		time.Sleep(s.cfg.TxLifetimeCheckInterval.Duration)
		txTrackers := s.worker.ExpireTransactions(s.cfg.TxLifetimeMax.Duration)
		failedReason := sequencer.ErrExpiredTransaction.Error()
		for _, txTracker := range txTrackers {
			err := s.pool.UpdateTxStatus(ctx, txTracker.Hash, pool.TxStatusFailed, false, &failedReason)
			if err != nil {
				log.Errorf("failed to update tx status, error: %v", err)
			}
		}
	}
}

// loadFromPool keeps loading transactions from the pool
func (s *Sequencer) loadFromPool(ctx context.Context) {
	for {
		poolTransactions, err := s.pool.GetNonWIPPendingTxs(ctx, s.cfg.QueryPendingTxsLimit)
		if err != nil && err != pool.ErrNotFound {
			log.Errorf("error loading txs from pool, error: %v", err)
		}

		for _, tx := range poolTransactions {
			err := s.addTxToWorker(ctx, tx)
			if err != nil {
				log.Errorf("error adding transaction to worker, error: %v", err)
			}
		}

		if len(poolTransactions) == 0 {
			time.Sleep(s.cfg.LoadPoolTxsCheckInterval.Duration)
		}
	}
}

func (s *Sequencer) addTxToWorker(ctx context.Context, tx pool.Transaction) error {
	txTracker, err := s.worker.NewTxTracker(tx, tx.ZKCounters, tx.ReservedZKCounters, tx.IP)
	if err != nil {
		return err
	}

	processedCount, err := s.pool.CountTransactionsByFromAndStatus(ctx, txTracker.From, pool.TxStatusSelected, pool.TxStatusFailed)
	if err != nil {
		return err
	}
	replacedTx, dropReason := s.worker.AddTxTracker(ctx, txTracker, processedCount)
	if dropReason != nil {
		failedReason := dropReason.Error()
		return s.pool.UpdateTxStatus(ctx, txTracker.Hash, pool.TxStatusFailed, false, &failedReason)
	} else {
		if replacedTx != nil {
			failedReason := sequencer.ErrReplacedTransaction.Error()
			err := s.pool.UpdateTxStatus(ctx, replacedTx.Hash, pool.TxStatusFailed, false, &failedReason)
			if err != nil {
				log.Warnf("error when setting as failed replacedTx %s, error: %v", replacedTx.HashStr, err)
			}
		}
		return s.pool.UpdateTxWIPStatus(ctx, tx.Hash(), true)
	}
}
