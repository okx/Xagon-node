package sequencer

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/0xPolygonHermez/zkevm-data-streamer/datastreamer"
	"github.com/0xPolygonHermez/zkevm-node/event"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	pmetrics "github.com/0xPolygonHermez/zkevm-node/sequencer/metrics"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/datastream"
	"github.com/ethereum/go-ethereum/common"
	"google.golang.org/protobuf/proto"
)

const (
	datastreamChannelBufferSize = 50
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

	numberOfStateInconsistencies uint64
}

// New init sequencer
func New(cfg Config, batchCfg state.BatchConfig, poolCfg pool.Config, txPool txPool, stateIntf stateInterface, etherman ethermanInterface, eventLog *event.EventLog) (*Sequencer, error) {
	sequencer := &Sequencer{
		cfg:       cfg,
		batchCfg:  batchCfg,
		poolCfg:   poolCfg,
		pool:      txPool,
		stateIntf: stateIntf,
		etherman:  etherman,
		eventLog:  eventLog,
	}

	sequencer.dataToStream = make(chan interface{}, datastreamChannelBufferSize)

	return sequencer, nil
}

// Start starts the sequencer
func (s *Sequencer) Start(ctx context.Context) {
	for !s.isSynced(ctx) {
		log.Infof("waiting for synchronizer to sync...")
		time.Sleep(time.Second)
	}
	pmetrics.Register()

	err := s.pool.MarkWIPTxsAsPending(ctx)
	if err != nil {
		log.Fatalf("failed to mark wip txs as pending, error: %v", err)
	}

	// Start stream server if enabled
	if s.cfg.StreamServer.Enabled {
		s.streamServer, err = datastreamer.NewServer(s.cfg.StreamServer.Port, s.cfg.StreamServer.Version, s.cfg.StreamServer.ChainID, state.StreamTypeSequencer, s.cfg.StreamServer.Filename, s.cfg.StreamServer.WriteTimeout.Duration, s.cfg.StreamServer.InactivityTimeout.Duration, s.cfg.StreamServer.InactivityCheckInterval.Duration, &s.cfg.StreamServer.Log)
		if err != nil {
			log.Fatalf("failed to create stream server, error: %v", err)
		}

		err = s.streamServer.Start()
		if err != nil {
			log.Fatalf("failed to start stream server, error: %v", err)
		}

		s.updateDataStreamerFile(ctx, s.cfg.StreamServer.ChainID)
	}

	if s.streamServer != nil {
		go s.sendDataToStreamer(s.cfg.StreamServer.ChainID, s.cfg.StreamServer.Version)
	}

	s.workerReadyTxsCond = newTimeoutCond(&sync.Mutex{})
	s.worker = NewWorker(s.stateIntf, s.batchCfg.Constraints, s.workerReadyTxsCond)
	s.finalizer = newFinalizer(s.cfg.Finalizer, s.poolCfg, s.worker, s.pool, s.stateIntf, s.etherman, s.cfg.L2Coinbase, s.isSynced, s.batchCfg.Constraints, s.eventLog, s.streamServer, s.workerReadyTxsCond, s.dataToStream)
	go s.finalizer.Start(ctx)

	go s.loadFromPool(ctx)

	go s.deleteOldPoolTxs(ctx)

	go s.expireOldWorkerTxs(ctx)

	go s.checkStateInconsistency(ctx)

	// X Layer handler
	go s.countPendingTx()
	go s.countReadyTx()

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
	err := state.GenerateDataStreamFile(ctx, s.streamServer, s.stateIntf, true, nil, chainID, s.cfg.StreamServer.UpgradeEtrogBatchNumber, s.cfg.StreamServer.Version)
	if err != nil {
		log.Fatalf("failed to generate data streamer file, error: %v", err)
	}
	log.Info("data streamer file updated")
}

func (s *Sequencer) deleteOldPoolTxs(ctx context.Context) {
	for {
		time.Sleep(s.cfg.DeletePoolTxsCheckInterval.Duration)

		if s.finalizer.haltFinalizer.Load() {
			return
		}

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

		if s.finalizer.haltFinalizer.Load() {
			return
		}

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
		if s.finalizer.haltFinalizer.Load() {
			return
		}
		poolTransactions, err := s.pool.GetNonWIPPendingTxs(ctx, getQueryPendingTxsLimit(s.cfg.QueryPendingTxsLimit))
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

	// XLayer free gas tx
	freeGp, isClaimTx, gp := s.checkFreeGas(tx, txTracker)
	if freeGp {
		defaultGp := s.pool.GetDynamicGasPrice()
		copyBaseGp := new(big.Int)
		if isClaimTx {
			txTracker.IsClaimTx = true
			baseGp := s.worker.getBaseClaimGp(defaultGp)
			copyBaseGp.Set(baseGp)
		} else {
			copyBaseGp.Set(defaultGp)
		}
		txTracker.GasPrice = copyBaseGp.Mul(copyBaseGp, new(big.Int).SetUint64(uint64(gp)))
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
func (s *Sequencer) sendDataToStreamer(chainID uint64, version uint8) {
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

		s.finalizer.DataToStreamChannelCountAdd(-1)

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

				bookMark := &datastream.BookMark{
					Type:  datastream.BookmarkType_BOOKMARK_TYPE_L2_BLOCK,
					Value: l2Block.L2BlockNumber,
				}

				marshalledBookMark, err := proto.Marshal(bookMark)
				if err != nil {
					log.Errorf("failed to marshal bookmark for l2block %d, error: %v", l2Block.L2BlockNumber, err)
					continue
				}

				_, err = s.streamServer.AddStreamBookmark(marshalledBookMark)
				if err != nil {
					log.Errorf("failed to add stream bookmark for l2block %d, error: %v", l2Block.L2BlockNumber, err)
					continue
				}

				// Get previous block timestamp to calculate delta timestamp
				previousL2Block := datastream.L2Block{}
				if l2Block.L2BlockNumber > 0 {
					bookMark = &datastream.BookMark{
						Type:  datastream.BookmarkType_BOOKMARK_TYPE_L2_BLOCK,
						Value: l2Block.L2BlockNumber - 1,
					}

					marshalledBookMark, err := proto.Marshal(bookMark)
					if err != nil {
						log.Errorf("failed to marshal bookmark for l2block %d, error: %v", l2Block.L2BlockNumber, err)
						continue
					}

					previousL2BlockEntry, err := s.streamServer.GetFirstEventAfterBookmark(marshalledBookMark)
					if err != nil {
						log.Errorf("failed to get previous l2block %d, error: %v", l2Block.L2BlockNumber-1, err)
						continue
					}

					err = proto.Unmarshal(previousL2BlockEntry.Data, &previousL2Block)
					if err != nil {
						log.Errorf("failed to unmarshal previous l2block %d, error: %v", l2Block.L2BlockNumber-1, err)
						continue
					}
				}

				streamL2Block := &datastream.L2Block{
					Number:          l2Block.L2BlockNumber,
					BatchNumber:     l2Block.BatchNumber,
					Timestamp:       l2Block.Timestamp,
					DeltaTimestamp:  uint32(l2Block.Timestamp - previousL2Block.Timestamp),
					MinTimestamp:    l2Block.MinTimestamp,
					L1Blockhash:     l2Block.L1BlockHash.Bytes(),
					L1InfotreeIndex: l2Block.L1InfoTreeIndex,
					Hash:            l2Block.BlockHash.Bytes(),
					StateRoot:       l2Block.StateRoot.Bytes(),
					GlobalExitRoot:  l2Block.GlobalExitRoot.Bytes(),
					Coinbase:        l2Block.Coinbase.Bytes(),
				}

				marshalledL2Block, err := proto.Marshal(streamL2Block)
				if err != nil {
					log.Errorf("failed to marshal l2block %d, error: %v", l2Block.L2BlockNumber, err)
					continue
				}

				_, err = s.streamServer.AddStreamEntry(datastreamer.EntryType(datastream.EntryType_ENTRY_TYPE_L2_BLOCK), marshalledL2Block)
				if err != nil {
					log.Errorf("failed to add stream entry for l2block %d, error: %v", l2Block.L2BlockNumber, err)
					continue
				}

				for _, l2Transaction := range l2Block.Txs {
					streamL2Transaction := &datastream.Transaction{
						L2BlockNumber:               l2Transaction.L2BlockNumber,
						IsValid:                     l2Transaction.IsValid != 0,
						Encoded:                     l2Transaction.Encoded,
						EffectiveGasPricePercentage: uint32(l2Transaction.EffectiveGasPricePercentage),
						ImStateRoot:                 l2Transaction.ImStateRoot.Bytes(),
					}

					// Clear the state root if the ForkID is > ETROG
					// currently this is redundant as the current implementation of the sequencer
					// leaves the ImStateRoot empty
					if l2Block.ForkID > state.FORKID_ETROG {
						streamL2Transaction.ImStateRoot = common.Hash{}.Bytes()
					}

					marshalledL2Transaction, err := proto.Marshal(streamL2Transaction)
					if err != nil {
						log.Errorf("failed to marshal l2tx for l2block %d, error: %v", l2Block.L2BlockNumber, err)
						continue
					}

					_, err = s.streamServer.AddStreamEntry(datastreamer.EntryType(datastream.EntryType_ENTRY_TYPE_TRANSACTION), marshalledL2Transaction)
					if err != nil {
						log.Errorf("failed to add l2tx stream entry for l2block %d, error: %v", l2Block.L2BlockNumber, err)
						continue
					}
				}

				if version >= state.DSVersion4 {
					streamL2BlockEnd := &datastream.L2BlockEnd{
						Number: l2Block.L2BlockNumber,
					}

					marshalledL2BlockEnd, err := proto.Marshal(streamL2BlockEnd)
					if err != nil {
						log.Errorf("failed to marshal l2block %d, error: %v", l2Block.L2BlockNumber, err)
						continue
					}

					_, err = s.streamServer.AddStreamEntry(datastreamer.EntryType(datastream.EntryType_ENTRY_TYPE_L2_BLOCK_END), marshalledL2BlockEnd)
					if err != nil {
						log.Errorf("failed to add stream entry for l2blockEnd %d, error: %v", l2Block.L2BlockNumber, err)
						continue
					}
				}

				err = s.streamServer.CommitAtomicOp()
				if err != nil {
					log.Errorf("failed to commit atomic op for l2block %d, error: %v ", l2Block.L2BlockNumber, err)
					continue
				}
			// Stream a bookmark
			case datastream.BookMark:
				err = s.streamServer.StartAtomicOp()
				if err != nil {
					log.Errorf("failed to start atomic op for bookmark type %d, value %d, error: %v", data.Type, data.Value, err)
					continue
				}

				marshalledBookMark, err := proto.Marshal(&data)
				if err != nil {
					log.Errorf("failed to marshal bookmark type %d, value %d, error: %v", data.Type, data.Value, err)
					continue
				}

				_, err = s.streamServer.AddStreamBookmark(marshalledBookMark)
				if err != nil {
					log.Errorf("failed to add stream bookmark for bookmark type %d, value %d, error: %v", data.Type, data.Value, err)
					continue
				}

				err = s.streamServer.CommitAtomicOp()
				if err != nil {
					log.Errorf("failed to commit atomic op for bookmark type %d, value %d, error: %v", data.Type, data.Value, err)
					continue
				}
			case datastream.BatchStart:
				err = s.streamServer.StartAtomicOp()
				if err != nil {
					log.Errorf("failed to start atomic op for batch start, error: %v", err)
					continue
				}

				data.ChainId = chainID

				marshalledBatchStart, err := proto.Marshal(&data)
				if err != nil {
					log.Errorf("failed to marshal batch start error: %v", err)
					continue
				}

				_, err = s.streamServer.AddStreamEntry(datastreamer.EntryType(datastream.EntryType_ENTRY_TYPE_BATCH_START), marshalledBatchStart)
				if err != nil {
					log.Errorf("failed to add stream entry for batch start, error: %v", err)
					continue
				}

				err = s.streamServer.CommitAtomicOp()
				if err != nil {
					log.Errorf("failed to commit atomic op for batch start, error: %v", err)
					continue
				}

			case datastream.BatchEnd:
				err = s.streamServer.StartAtomicOp()
				if err != nil {
					log.Errorf("failed to start atomic op for batch end, error: %v", err)
					continue
				}

				marshalledBatchEnd, err := proto.Marshal(&data)
				if err != nil {
					log.Errorf("failed to marshal batch end, error: %v", err)
					continue
				}

				_, err = s.streamServer.AddStreamEntry(datastreamer.EntryType(datastream.EntryType_ENTRY_TYPE_BATCH_END), marshalledBatchEnd)
				if err != nil {
					log.Errorf("failed to add stream entry for batch end, error: %v", err)
					continue
				}

				err = s.streamServer.CommitAtomicOp()
				if err != nil {
					log.Errorf("failed to commit atomic op for batch end, error: %v", err)
					continue
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
