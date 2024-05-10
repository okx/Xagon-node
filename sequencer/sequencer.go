package sequencer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/0xPolygonHermez/zkevm-data-streamer/datastreamer"
	"github.com/0xPolygonHermez/zkevm-node/event"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
)

const (
	datastreamChannelMultiplier = 2
)

// Sequencer represents a sequencer
type Sequencer struct {
	cfg      Config
	batchCfg state.BatchConfig
	poolCfg  pool.Config

	pool      txPool
	stateIntf stateInterface
	eventLog  *event.EventLog
	etherman  ethermanInterface
	worker    *Worker
	finalizer *finalizer

	workerReadyTxsCond *timeoutCond

	streamServer *datastreamer.StreamServer
	dataToStream chan interface{}

	address common.Address

	numberOfStateInconsistencies uint64
}

// New init sequencer
func New(cfg Config, batchCfg state.BatchConfig, poolCfg pool.Config, txPool txPool, stateIntf stateInterface, etherman ethermanInterface, eventLog *event.EventLog) (*Sequencer, error) {
	addr, err := etherman.TrustedSequencer()
	if err != nil {
		return nil, fmt.Errorf("failed to get trusted sequencer address, error: %v", err)
	}

	sequencer := &Sequencer{
		cfg:       cfg,
		batchCfg:  batchCfg,
		poolCfg:   poolCfg,
		pool:      txPool,
		stateIntf: stateIntf,
		etherman:  etherman,
		address:   addr,
		eventLog:  eventLog,
	}

	// TODO: Make configurable
	channelBufferSize := 200 * datastreamChannelMultiplier // nolint:gomnd
	sequencer.dataToStream = make(chan interface{}, channelBufferSize)

	return sequencer, nil
}

// Start starts the sequencer
func (s *Sequencer) Start(ctx context.Context) {
	for !s.isSynced(ctx) {
		log.Infof("waiting for synchronizer to sync...")
		time.Sleep(time.Second)
	}

	err := s.pool.MarkWIPTxsAsPending(ctx)
	if err != nil {
		log.Fatalf("failed to mark WIP txs as pending, error: %v", err)
	}

	// Start stream server if enabled
	if s.cfg.StreamServer.Enabled {
		s.streamServer, err = datastreamer.NewServer(s.cfg.StreamServer.Port, s.cfg.StreamServer.Version, s.cfg.StreamServer.ChainID, state.StreamTypeSequencer, s.cfg.StreamServer.Filename, &s.cfg.StreamServer.Log)
		if err != nil {
			log.Fatalf("failed to create stream server, error: %v", err)
		}

		err = s.streamServer.Start()
		if err != nil {
			log.Fatalf("failed to start stream server, error: %v", err)
		}

		s.updateDataStreamerFile(ctx, s.cfg.StreamServer.ChainID)
	}

	go s.loadFromPool(ctx)

	if s.streamServer != nil {
		go s.sendDataToStreamer(s.cfg.StreamServer.ChainID)
	}

	s.workerReadyTxsCond = newTimeoutCond(&sync.Mutex{})
	s.worker = NewWorker(s.stateIntf, s.batchCfg.Constraints, s.workerReadyTxsCond)
	s.finalizer = newFinalizer(s.cfg.Finalizer, s.poolCfg, s.worker, s.pool, s.stateIntf, s.etherman, s.address, s.isSynced, s.batchCfg.Constraints, s.eventLog, s.streamServer, s.workerReadyTxsCond, s.dataToStream)
	go s.finalizer.Start(ctx)

	go s.deleteOldPoolTxs(ctx)

	go s.expireOldWorkerTxs(ctx)

	go s.checkStateInconsistency(ctx)

	// Wait until context is done
	<-ctx.Done()
}

// checkStateInconsistency checks if state inconsistency happened
func (s *Sequencer) checkStateInconsistency(ctx context.Context) {
	var err error
	s.numberOfStateInconsistencies, err = s.stateIntf.CountReorgs(ctx, nil)
	if err != nil {
		log.Error("failed to get initial number of reorgs, error: %v", err)
	}
	for {
		stateInconsistenciesDetected, err := s.stateIntf.CountReorgs(ctx, nil)
		if err != nil {
			log.Error("failed to get number of reorgs, error: %v", err)
			return
		}

		if stateInconsistenciesDetected != s.numberOfStateInconsistencies {
			s.finalizer.Halt(ctx, fmt.Errorf("state inconsistency detected, halting finalizer"), false)
		}

		time.Sleep(s.cfg.StateConsistencyCheckInterval.Duration)
	}
}

func (s *Sequencer) updateDataStreamerFile(ctx context.Context, chainID uint64) {
	err := state.GenerateDataStreamerFile(ctx, s.streamServer, s.stateIntf, true, nil, chainID, s.cfg.StreamServer.UpgradeEtrogBatchNumber)
	if err != nil {
		log.Fatalf("failed to generate data streamer file, error: %v", err)
	}
	log.Info("data streamer file updated")
}

func (s *Sequencer) deleteOldPoolTxs(ctx context.Context) {
	for {
		time.Sleep(s.cfg.DeletePoolTxsCheckInterval.Duration)
		log.Infof("trying to get txs to delete from the pool...")
		earliestTxHash, err := s.pool.GetEarliestProcessedTx(ctx)
		if err != nil {
			log.Errorf("failed to get earliest tx hash to delete, err: %v", err)
			continue
		}

		txHashes, err := s.stateIntf.GetTxsOlderThanNL1BlocksUntilTxHash(ctx, s.cfg.DeletePoolTxsL1BlockConfirmations, earliestTxHash, nil)
		if err != nil {
			log.Errorf("failed to get txs hashes to delete, error: %v", err)
			continue
		}
		log.Infof("trying to delete %d selected txs", len(txHashes))
		err = s.pool.DeleteTransactionsByHashes(ctx, txHashes)
		if err != nil {
			log.Errorf("failed to delete selected txs from the pool, error: %v", err)
			continue
		}
		log.Infof("deleted %d selected txs from the pool", len(txHashes))

		log.Infof("trying to delete failed txs from the pool")
		// Delete failed txs older than a certain date (14 seconds per L1 block)
		err = s.pool.DeleteFailedTransactionsOlderThan(ctx, time.Now().Add(-time.Duration(s.cfg.DeletePoolTxsL1BlockConfirmations*14)*time.Second)) //nolint:gomnd
		if err != nil {
			log.Errorf("failed to delete failed txs from the pool, error: %v", err)
			continue
		}
		log.Infof("failed txs deleted from the pool")
	}
}

func (s *Sequencer) expireOldWorkerTxs(ctx context.Context) {
	for {
		time.Sleep(s.cfg.TxLifetimeCheckInterval.Duration)
		txTrackers := s.worker.ExpireTransactions(s.cfg.TxLifetimeMax.Duration)
		failedReason := ErrExpiredTransaction.Error()
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
		poolTransactions, err := s.pool.GetNonWIPPendingTxs(ctx)
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
	txTracker, err := s.worker.NewTxTracker(tx.Transaction, tx.ZKCounters, tx.ReservedZKCounters, tx.IP)
	if err != nil {
		return err
	}
	replacedTx, dropReason := s.worker.AddTxTracker(ctx, txTracker)
	if dropReason != nil {
		failedReason := dropReason.Error()
		return s.pool.UpdateTxStatus(ctx, txTracker.Hash, pool.TxStatusFailed, false, &failedReason)
	} else {
		if replacedTx != nil {
			failedReason := ErrReplacedTransaction.Error()
			err := s.pool.UpdateTxStatus(ctx, replacedTx.Hash, pool.TxStatusFailed, false, &failedReason)
			if err != nil {
				log.Warnf("error when setting as failed replacedTx %s, error: %v", replacedTx.HashStr, err)
			}
		}
		return s.pool.UpdateTxWIPStatus(ctx, tx.Hash(), true)
	}
}

// sendDataToStreamer sends data to the data stream server
func (s *Sequencer) sendDataToStreamer(chainID uint64) {
	var err error
	for {
		// Read error from previous iteration
		if err != nil {
			err = s.streamServer.RollbackAtomicOp()
			if err != nil {
				log.Errorf("failed to rollback atomic op, error: %v", err)
			}
			s.streamServer = nil
		}

		// Read data from channel
		dataStream := <-s.dataToStream

		if s.streamServer != nil {
			switch data := dataStream.(type) {
			// Stream a complete L2 block with its transactions
			case state.DSL2FullBlock:
				l2Block := data

				err = s.streamServer.StartAtomicOp()
				if err != nil {
					log.Errorf("failed to start atomic op for l2block %d, error: %v ", l2Block.L2BlockNumber, err)
					continue
				}

				bookMark := state.DSBookMark{
					Type:  state.BookMarkTypeL2Block,
					Value: l2Block.L2BlockNumber,
				}

				_, err = s.streamServer.AddStreamBookmark(bookMark.Encode())
				if err != nil {
					log.Errorf("failed to add stream bookmark for l2block %d, error: %v", l2Block.L2BlockNumber, err)
					continue
				}

				// Get previous block timestamp to calculate delta timestamp
				previousL2Block := state.DSL2BlockStart{}
				if l2Block.L2BlockNumber > 0 {
					bookMark = state.DSBookMark{
						Type:  state.BookMarkTypeL2Block,
						Value: l2Block.L2BlockNumber - 1,
					}

					previousL2BlockEntry, err := s.streamServer.GetFirstEventAfterBookmark(bookMark.Encode())
					if err != nil {
						log.Errorf("failed to get previous l2block %d, error: %v", l2Block.L2BlockNumber-1, err)
						continue
					}

					previousL2Block = state.DSL2BlockStart{}.Decode(previousL2BlockEntry.Data)
				}

				blockStart := state.DSL2BlockStart{
					BatchNumber:     l2Block.BatchNumber,
					L2BlockNumber:   l2Block.L2BlockNumber,
					Timestamp:       l2Block.Timestamp,
					DeltaTimestamp:  uint32(l2Block.Timestamp - previousL2Block.Timestamp),
					L1InfoTreeIndex: l2Block.L1InfoTreeIndex,
					L1BlockHash:     l2Block.L1BlockHash,
					GlobalExitRoot:  l2Block.GlobalExitRoot,
					Coinbase:        l2Block.Coinbase,
					ForkID:          l2Block.ForkID,
					ChainID:         uint32(chainID),
				}

				_, err = s.streamServer.AddStreamEntry(state.EntryTypeL2BlockStart, blockStart.Encode())
				if err != nil {
					log.Errorf("failed to add stream entry for l2block %d, error: %v", l2Block.L2BlockNumber, err)
					continue
				}

				for _, l2Transaction := range l2Block.Txs {
					_, err = s.streamServer.AddStreamEntry(state.EntryTypeL2Tx, l2Transaction.Encode())
					if err != nil {
						log.Errorf("failed to add l2tx stream entry for l2block %d, error: %v", l2Block.L2BlockNumber, err)
						continue
					}
				}

				blockEnd := state.DSL2BlockEnd{
					L2BlockNumber: l2Block.L2BlockNumber,
					BlockHash:     l2Block.BlockHash,
					StateRoot:     l2Block.StateRoot,
				}

				_, err = s.streamServer.AddStreamEntry(state.EntryTypeL2BlockEnd, blockEnd.Encode())
				if err != nil {
					log.Errorf("failed to add stream entry for l2block %d, error: %v", l2Block.L2BlockNumber, err)
					continue
				}

				err = s.streamServer.CommitAtomicOp()
				if err != nil {
					log.Errorf("failed to commit atomic op for l2block %d, error: %v ", l2Block.L2BlockNumber, err)
					continue
				}

			// Stream a bookmark
			case state.DSBookMark:
				bookmark := data

				err = s.streamServer.StartAtomicOp()
				if err != nil {
					log.Errorf("failed to start atomic op for bookmark type %d, value %d, error: %v", bookmark.Type, bookmark.Value, err)
					continue
				}

				_, err = s.streamServer.AddStreamBookmark(bookmark.Encode())
				if err != nil {
					log.Errorf("failed to add stream bookmark type %d, value %d, error: %v", bookmark.Type, bookmark.Value, err)
					continue
				}

				err = s.streamServer.CommitAtomicOp()
				if err != nil {
					log.Errorf("failed to commit atomic op for bookmark type %d, value %d, error: %v", bookmark.Type, bookmark.Value, err)
				}

			// Invalid stream message type
			default:
				log.Errorf("invalid stream message type received")
			}
		}
	}
}

func (s *Sequencer) isSynced(ctx context.Context) bool {
	lastVirtualBatchNum, err := s.stateIntf.GetLastVirtualBatchNum(ctx, nil)
	if err != nil && err != state.ErrNotFound {
		log.Errorf("failed to get last isSynced batch, error: %v", err)
		return false
	}
	lastTrustedBatchNum, err := s.stateIntf.GetLastBatchNumber(ctx, nil)
	if err != nil && err != state.ErrNotFound {
		log.Errorf("failed to get last batch num, error: %v", err)
		return false
	}
	if lastTrustedBatchNum > lastVirtualBatchNum {
		return true
	}
	lastEthBatchNum, err := s.etherman.GetLatestBatchNumber()
	if err != nil {
		log.Errorf("failed to get last eth batch, error: %v", err)
		return false
	}
	if lastVirtualBatchNum < lastEthBatchNum {
		log.Infof("waiting for the state to be synced, lastVirtualBatchNum: %d, lastEthBatchNum: %d", lastVirtualBatchNum, lastEthBatchNum)
		return false
	}

	return true
}
