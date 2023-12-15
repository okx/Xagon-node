package synchronizer

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strings"
	"time"

	"github.com/0xPolygon/cdk-data-availability/client"
	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/event"
	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	stateMetrics "github.com/0xPolygonHermez/zkevm-node/state/metrics"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/metrics"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
)

const (
	forkID5 = 5
	// ParallelMode is the value for L1SynchronizationMode to run in parallel mode
	ParallelMode = "parallel"
	// SequentialMode is the value for L1SynchronizationMode to run in sequential mode
	SequentialMode = "sequential"
)

// Synchronizer connects L1 and L2
type Synchronizer interface {
	Sync() error
	Stop()
}

// ClientSynchronizer connects L1 and L2
type ClientSynchronizer struct {
	isTrustedSequencer bool
	etherMan           EthermanInterface
	latestFlushID      uint64
	// If true the lastFlushID is stored in DB and we don't need to check again
	latestFlushIDIsFulfilled bool
	etherManForL1            []EthermanInterface
	state                    stateInterface
	pool                     poolInterface
	ethTxManager             ethTxManager
	zkEVMClient              zkEVMClientInterface
	eventLog                 *event.EventLog
	ctx                      context.Context
	cancelCtx                context.CancelFunc
	genesis                  state.Genesis
	cfg                      Config
	trustedState             struct {
		lastTrustedBatches []*state.Batch
		lastStateRoot      *common.Hash
	}
	// Id of the 'process' of the executor. Each time that it starts this value changes
	// This value is obtained from the call state.GetStoredFlushID
	// It starts as an empty string and it is filled in the first call
	// later the value is checked to be the same (in function checkFlushID)
	proverID string
	// Previous value returned by state.GetStoredFlushID, is used for decide if write a log or not
	previousExecutorFlushID    uint64
	committeeMembers           []etherman.DataCommitteeMember
	selectedCommitteeMember    int
	dataCommitteeClientFactory client.ClientFactoryInterface
	l1SyncOrchestration        *l1SyncOrchestration
}

// NewSynchronizer creates and initializes an instance of Synchronizer
func NewSynchronizer(
	isTrustedSequencer bool,
	ethMan EthermanInterface,
	etherManForL1 []EthermanInterface,
	st stateInterface,
	pool poolInterface,
	ethTxManager ethTxManager,
	zkEVMClient zkEVMClientInterface,
	eventLog *event.EventLog,
	genesis state.Genesis,
	cfg Config, clientFactory client.ClientFactoryInterface,
	runInDevelopmentMode bool) (Synchronizer, error) {
	ctx, cancel := context.WithCancel(context.Background())
	metrics.Register()

	res := &ClientSynchronizer{
		isTrustedSequencer:         isTrustedSequencer,
		state:                      st,
		etherMan:                   ethMan,
		etherManForL1:              etherManForL1,
		pool:                       pool,
		ctx:                        ctx,
		cancelCtx:                  cancel,
		ethTxManager:               ethTxManager,
		zkEVMClient:                zkEVMClient,
		eventLog:                   eventLog,
		genesis:                    genesis,
		cfg:                        cfg,
		proverID:                   "",
		previousExecutorFlushID:    0,
		dataCommitteeClientFactory: clientFactory,
		l1SyncOrchestration:        nil,
	}
	switch cfg.L1SynchronizationMode {
	case ParallelMode:
		log.Info("L1SynchronizationMode is parallel")
		var err error
		res.l1SyncOrchestration, err = newL1SyncParallel(ctx, cfg, etherManForL1, res, runInDevelopmentMode)
		if err != nil {
			log.Fatalf("Can't initialize L1SyncParallel. Error: %s", err)
		}
	case SequentialMode:
		log.Info("L1SynchronizationMode is sequential")
	default:
		log.Fatalf("L1SynchronizationMode is not valid. Valid values are: %s, %s", ParallelMode, SequentialMode)
	}
	err := res.loadCommittee()
	return res, err
}

var waitDuration = time.Duration(0)

func newL1SyncParallel(ctx context.Context, cfg Config, etherManForL1 []EthermanInterface, sync *ClientSynchronizer, runExternalControl bool) (*l1SyncOrchestration, error) {
	chIncommingRollupInfo := make(chan l1SyncMessage, cfg.L1ParallelSynchronization.MaxPendingNoProcessedBlocks)
	cfgConsumer := configConsumer{
		ApplyAfterNumRollupReceived: cfg.L1ParallelSynchronization.PerformanceWarning.ApplyAfterNumRollupReceived,
		AceptableInacctivityTime:    cfg.L1ParallelSynchronization.PerformanceWarning.AceptableInacctivityTime.Duration,
	}
	L1DataProcessor := newL1RollupInfoConsumer(cfgConsumer, sync, chIncommingRollupInfo)

	cfgProducer := configProducer{
		syncChunkSize:                              cfg.SyncChunkSize,
		ttlOfLastBlockOnL1:                         cfg.L1ParallelSynchronization.RequestLastBlockPeriod.Duration,
		timeoutForRequestLastBlockOnL1:             cfg.L1ParallelSynchronization.RequestLastBlockTimeout.Duration,
		numOfAllowedRetriesForRequestLastBlockOnL1: cfg.L1ParallelSynchronization.RequestLastBlockMaxRetries,
		timeForShowUpStatisticsLog:                 cfg.L1ParallelSynchronization.StatisticsPeriod.Duration,
		timeOutMainLoop:                            cfg.L1ParallelSynchronization.TimeOutMainLoop.Duration,
		minTimeBetweenRetriesForRollupInfo:         cfg.L1ParallelSynchronization.RollupInfoRetriesSpacing.Duration,
	}
	l1DataRetriever := newL1DataRetriever(cfgProducer, etherManForL1, chIncommingRollupInfo)
	l1SyncOrchestration := newL1SyncOrchestration(ctx, l1DataRetriever, L1DataProcessor)
	if runExternalControl {
		log.Infof("Starting external control")
		externalControl := newExternalControl(l1DataRetriever, l1SyncOrchestration)
		externalControl.start()
	}
	return l1SyncOrchestration, nil
}

// Sync function will read the last state synced and will continue from that point.
// Sync() will read blockchain events to detect rollup updates
func (s *ClientSynchronizer) Sync() error {
	startInitialization := time.Now()
	// If there is no lastEthereumBlock means that sync from the beginning is necessary. If not, it continues from the retrieved ethereum block
	// Get the latest synced block. If there is no block on db, use genesis block
	log.Info("Sync started")
	dbTx, err := s.state.BeginStateTransaction(s.ctx)
	if err != nil {
		log.Errorf("error creating db transaction to get latest block. Error: %v", err)
		return err
	}
	lastEthBlockSynced, err := s.state.GetLastBlock(s.ctx, dbTx)
	if err != nil {
		if errors.Is(err, state.ErrStateNotSynchronized) {
			log.Info("State is empty, verifying genesis block")
			valid, err := s.etherMan.VerifyGenBlockNumber(s.ctx, s.genesis.GenesisBlockNum)
			if err != nil {
				log.Error("error checking genesis block number. Error: ", err)
				rollbackErr := dbTx.Rollback(s.ctx)
				if rollbackErr != nil {
					log.Errorf("error rolling back state. RollbackErr: %v, err: %s", rollbackErr, err.Error())
					return rollbackErr
				}
				return err
			} else if !valid {
				log.Error("genesis Block number configured is not valid. It is required the block number where the PolygonZkEVM smc was deployed")
				rollbackErr := dbTx.Rollback(s.ctx)
				if rollbackErr != nil {
					log.Errorf("error rolling back state. RollbackErr: %v", rollbackErr)
					return rollbackErr
				}
				return fmt.Errorf("genesis Block number configured is not valid. It is required the block number where the PolygonZkEVM smc was deployed")
			}
			log.Info("Setting genesis block")
			header, err := s.etherMan.HeaderByNumber(s.ctx, big.NewInt(0).SetUint64(s.genesis.GenesisBlockNum))
			if err != nil {
				log.Errorf("error getting l1 block header for block %d. Error: %v", s.genesis.GenesisBlockNum, err)
				rollbackErr := dbTx.Rollback(s.ctx)
				if rollbackErr != nil {
					log.Errorf("error rolling back state. RollbackErr: %v, err: %s", rollbackErr, err.Error())
					return rollbackErr
				}
				return err
			}
			lastEthBlockSynced = &state.Block{
				BlockNumber: header.Number.Uint64(),
				BlockHash:   header.Hash(),
				ParentHash:  header.ParentHash,
				ReceivedAt:  time.Unix(int64(header.Time), 0),
			}
			newRoot, err := s.state.SetGenesis(s.ctx, *lastEthBlockSynced, s.genesis, dbTx)
			if err != nil {
				log.Error("error setting genesis: ", err)
				rollbackErr := dbTx.Rollback(s.ctx)
				if rollbackErr != nil {
					log.Errorf("error rolling back state. RollbackErr: %v, err: %s", rollbackErr, err.Error())
					return rollbackErr
				}
				return err
			}
			blocks, _, err := s.etherMan.GetRollupInfoByBlockRange(s.ctx, lastEthBlockSynced.BlockNumber, &lastEthBlockSynced.BlockNumber)
			if err != nil {
				log.Fatal(err)
			}
			err = s.processForkID(blocks[0].ForkIDs[0], blocks[0].BlockNumber, dbTx)
			if err != nil {
				log.Error("error storing genesis forkID: ", err)
				return err
			}
			var root common.Hash
			root.SetBytes(newRoot)
			if root != s.genesis.Root {
				log.Errorf("Calculated newRoot should be %s instead of %s", s.genesis.Root.String(), root.String())
				rollbackErr := dbTx.Rollback(s.ctx)
				if rollbackErr != nil {
					log.Errorf("error rolling back state. RollbackErr: %v", rollbackErr)
					return rollbackErr
				}
				return fmt.Errorf("calculated newRoot should be %s instead of %s", s.genesis.Root.String(), root.String())
			}
			log.Debug("Genesis root matches!")
		} else {
			log.Error("unexpected error getting the latest ethereum block. Error: ", err)
			rollbackErr := dbTx.Rollback(s.ctx)
			if rollbackErr != nil {
				log.Errorf("error rolling back state. RollbackErr: %v, err: %s", rollbackErr, err.Error())
				return rollbackErr
			}
			return err
		}
	}
	initBatchNumber, err := s.state.GetLastBatchNumber(s.ctx, dbTx)
	if err != nil {
		log.Error("error getting latest batchNumber synced. Error: ", err)
		rollbackErr := dbTx.Rollback(s.ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state. RollbackErr: %v, err: %s", rollbackErr, err.Error())
			return rollbackErr
		}
		return err
	}
	err = s.state.SetInitSyncBatch(s.ctx, initBatchNumber, dbTx)
	if err != nil {
		log.Error("error setting initial batch number. Error: ", err)
		rollbackErr := dbTx.Rollback(s.ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state. RollbackErr: %v, err: %s", rollbackErr, err.Error())
			return rollbackErr
		}
		return err
	}
	if err := dbTx.Commit(s.ctx); err != nil {
		log.Errorf("error committing dbTx, err: %v", err)
		rollbackErr := dbTx.Rollback(s.ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state. RollbackErr: %v, err: %s", rollbackErr, err.Error())
			return rollbackErr
		}
		return err
	}
	metrics.InitializationTime(time.Since(startInitialization))

	for {
		select {
		case <-s.ctx.Done():
			return nil
		case <-time.After(waitDuration):
			start := time.Now()
			latestSequencedBatchNumber, err := s.etherMan.GetLatestBatchNumber()
			if err != nil {
				log.Warn("error getting latest sequenced batch in the rollup. Error: ", err)
				continue
			}
			latestSyncedBatch, err := s.state.GetLastBatchNumber(s.ctx, nil)
			if err != nil {
				log.Warn("error getting latest batch synced in the db. Error: ", err)
				continue
			}
			// Check the latest verified Batch number in the smc
			lastVerifiedBatchNumber, err := s.etherMan.GetLatestVerifiedBatchNum()
			if err != nil {
				log.Warn("error getting last verified batch in the rollup. Error: ", err)
				continue
			}
			err = s.state.SetLastBatchInfoSeenOnEthereum(s.ctx, latestSequencedBatchNumber, lastVerifiedBatchNumber, nil)
			if err != nil {
				log.Warn("error setting latest batch info into db. Error: ", err)
				continue
			}
			log.Infof("latestSequencedBatchNumber: %d, latestSyncedBatch: %d, lastVerifiedBatchNumber: %d", latestSequencedBatchNumber, latestSyncedBatch, lastVerifiedBatchNumber)
			// Sync trusted state
			if latestSyncedBatch >= latestSequencedBatchNumber {
				startTrusted := time.Now()
				log.Info("Syncing trusted state")
				err = s.syncTrustedState(latestSyncedBatch)
				metrics.FullTrustedSyncTime(time.Since(startTrusted))
				if err != nil {
					log.Warn("error syncing trusted state. Error: ", err)
					s.trustedState.lastTrustedBatches = nil
					s.trustedState.lastStateRoot = nil
					continue
				}
				waitDuration = s.cfg.SyncInterval.Duration
			}
			//Sync L1Blocks
			startL1 := time.Now()
			if s.l1SyncOrchestration != nil && (latestSyncedBatch < latestSequencedBatchNumber || !s.cfg.L1ParallelSynchronization.FallbackToSequentialModeOnSynchronized) {
				log.Infof("Syncing L1 blocks in parallel lastEthBlockSynced=%d", lastEthBlockSynced.BlockNumber)
				lastEthBlockSynced, err = s.syncBlocksParallel(lastEthBlockSynced)
			} else {
				if s.l1SyncOrchestration != nil {
					log.Infof("Switching to sequential mode, stopping parallel sync and deleting object")
					s.l1SyncOrchestration.abort()
					s.l1SyncOrchestration = nil
				}
				log.Infof("Syncing L1 blocks sequentially lastEthBlockSynced=%d", lastEthBlockSynced.BlockNumber)
				lastEthBlockSynced, err = s.syncBlocksSequential(lastEthBlockSynced)
			}
			metrics.FullL1SyncTime(time.Since(startL1))
			if err != nil {
				log.Warn("error syncing blocks: ", err)
				lastEthBlockSynced, err = s.state.GetLastBlock(s.ctx, nil)
				if err != nil {
					log.Fatal("error getting lastEthBlockSynced to resume the synchronization... Error: ", err)
				}
				if s.l1SyncOrchestration != nil {
					// If have failed execution and get starting point from DB, we must reset parallel sync to this point
					// producer must start requesting this block
					s.l1SyncOrchestration.reset(lastEthBlockSynced.BlockNumber)
				}
				if s.ctx.Err() != nil {
					continue
				}
			}
			metrics.FullSyncIterationTime(time.Since(start))
			log.Info("L1 state fully synchronized")
		}
	}
}

// This function syncs the node from a specific block to the latest
// lastEthBlockSynced -> last block synced in the db
func (s *ClientSynchronizer) syncBlocksParallel(lastEthBlockSynced *state.Block) (*state.Block, error) {
	// This function will read events fromBlockNum to latestEthBlock. Check reorg to be sure that everything is ok.
	block, err := s.checkReorg(lastEthBlockSynced)
	if err != nil {
		log.Errorf("error checking reorgs. Retrying... Err: %v", err)
		return lastEthBlockSynced, fmt.Errorf("error checking reorgs")
	}
	if block != nil {
		log.Infof("reorg detected. Resetting the state from block %v to block %v", lastEthBlockSynced.BlockNumber, block.BlockNumber)
		err = s.resetState(block.BlockNumber)
		if err != nil {
			log.Errorf("error resetting the state to a previous block. Retrying... Err: %v", err)
			s.l1SyncOrchestration.reset(lastEthBlockSynced.BlockNumber)
			return lastEthBlockSynced, fmt.Errorf("error resetting the state to a previous block")
		}
		return block, nil
	}
	if !s.l1SyncOrchestration.isProducerRunning() {
		log.Infof("producer is not running. Resetting the state to start from  block %v (last on DB)", lastEthBlockSynced.BlockNumber)
		s.l1SyncOrchestration.producer.Reset(lastEthBlockSynced.BlockNumber)
	}
	log.Infof("Starting L1 sync orchestrator in parallel block: %d", lastEthBlockSynced.BlockNumber)
	return s.l1SyncOrchestration.start(lastEthBlockSynced)
}

// This function syncs the node from a specific block to the latest
func (s *ClientSynchronizer) syncBlocksSequential(lastEthBlockSynced *state.Block) (*state.Block, error) {
	// This function will read events fromBlockNum to latestEthBlock. Check reorg to be sure that everything is ok.
	block, err := s.checkReorg(lastEthBlockSynced)
	if err != nil {
		log.Errorf("error checking reorgs. Retrying... Err: %v", err)
		return lastEthBlockSynced, fmt.Errorf("error checking reorgs")
	}
	if block != nil {
		err = s.resetState(block.BlockNumber)
		if err != nil {
			log.Errorf("error resetting the state to a previous block. Retrying... Err: %v", err)
			return lastEthBlockSynced, fmt.Errorf("error resetting the state to a previous block")
		}
		return block, nil
	}

	// Call the blockchain to retrieve data
	header, err := s.etherMan.HeaderByNumber(s.ctx, nil)
	if err != nil {
		return lastEthBlockSynced, err
	}
	lastKnownBlock := header.Number

	var fromBlock uint64
	if lastEthBlockSynced.BlockNumber > 0 {
		fromBlock = lastEthBlockSynced.BlockNumber + 1
	}

	for {
		toBlock := fromBlock + s.cfg.SyncChunkSize
		log.Infof("Syncing block %d of %d", fromBlock, lastKnownBlock.Uint64())
		log.Infof("Getting rollup info from block %d to block %d", fromBlock, toBlock)
		// This function returns the rollup information contained in the ethereum blocks and an extra param called order.
		// Order param is a map that contains the event order to allow the synchronizer store the info in the same order that is readed.
		// Name can be different in the order struct. For instance: Batches or Name:NewSequencers. This name is an identifier to check
		// if the next info that must be stored in the db is a new sequencer or a batch. The value pos (position) tells what is the
		// array index where this value is.
		start := time.Now()
		blocks, order, err := s.etherMan.GetRollupInfoByBlockRange(s.ctx, fromBlock, &toBlock)
		metrics.ReadL1DataTime(time.Since(start))
		if err != nil {
			return lastEthBlockSynced, err
		}
		start = time.Now()
		err = s.processBlockRange(blocks, order)
		metrics.ProcessL1DataTime(time.Since(start))
		if err != nil {
			return lastEthBlockSynced, err
		}
		if len(blocks) > 0 {
			lastEthBlockSynced = &state.Block{
				BlockNumber: blocks[len(blocks)-1].BlockNumber,
				BlockHash:   blocks[len(blocks)-1].BlockHash,
				ParentHash:  blocks[len(blocks)-1].ParentHash,
				ReceivedAt:  blocks[len(blocks)-1].ReceivedAt,
			}
			for i := range blocks {
				log.Debug("Position: ", i, ". BlockNumber: ", blocks[i].BlockNumber, ". BlockHash: ", blocks[i].BlockHash)
			}
		}
		fromBlock = toBlock + 1

		if lastKnownBlock.Cmp(new(big.Int).SetUint64(toBlock)) < 1 {
			waitDuration = s.cfg.SyncInterval.Duration
			break
		}
		if len(blocks) == 0 { // If there is no events in the checked blocks range and lastKnownBlock > fromBlock.
			// Store the latest block of the block range. Get block info and process the block
			fb, err := s.etherMan.EthBlockByNumber(s.ctx, toBlock)
			if err != nil {
				return lastEthBlockSynced, err
			}
			b := etherman.Block{
				BlockNumber: fb.NumberU64(),
				BlockHash:   fb.Hash(),
				ParentHash:  fb.ParentHash(),
				ReceivedAt:  time.Unix(int64(fb.Time()), 0),
			}
			err = s.processBlockRange([]etherman.Block{b}, order)
			if err != nil {
				return lastEthBlockSynced, err
			}
			block := state.Block{
				BlockNumber: fb.NumberU64(),
				BlockHash:   fb.Hash(),
				ParentHash:  fb.ParentHash(),
				ReceivedAt:  time.Unix(int64(fb.Time()), 0),
			}
			lastEthBlockSynced = &block
			log.Debug("Storing empty block. BlockNumber: ", b.BlockNumber, ". BlockHash: ", b.BlockHash)
		}
	}

	return lastEthBlockSynced, nil
}

// syncTrustedState synchronizes information from the trusted sequencer
// related to the trusted state when the node has all the information from
// l1 synchronized
func (s *ClientSynchronizer) syncTrustedState(latestSyncedBatch uint64) error {
	if s.isTrustedSequencer {
		return nil
	}

	log.Info("syncTrustedState: Getting trusted state info")
	start := time.Now()
	lastTrustedStateBatchNumber, err := s.zkEVMClient.BatchNumber(s.ctx)
	metrics.GetTrustedBatchNumberTime(time.Since(start))
	if err != nil {
		log.Warn("syncTrustedState: error syncing trusted state. Error: ", err)
		return err
	}

	log.Debug("syncTrustedState: lastTrustedStateBatchNumber ", lastTrustedStateBatchNumber)
	log.Debug("syncTrustedState: latestSyncedBatch ", latestSyncedBatch)
	if lastTrustedStateBatchNumber < latestSyncedBatch {
		return nil
	}

	batchNumberToSync := latestSyncedBatch
	for batchNumberToSync <= lastTrustedStateBatchNumber {
		if batchNumberToSync == 0 {
			batchNumberToSync++
			continue
		}
		start = time.Now()
		batchToSync, err := s.zkEVMClient.BatchByNumber(s.ctx, big.NewInt(0).SetUint64(batchNumberToSync))
		metrics.GetTrustedBatchInfoTime(time.Since(start))
		if err != nil {
			log.Warnf("syncTrustedState: failed to get batch %d from trusted state. Error: %v", batchNumberToSync, err)
			return err
		}

		dbTx, err := s.state.BeginStateTransaction(s.ctx)
		if err != nil {
			log.Errorf("syncTrustedState: error creating db transaction to sync trusted batch %d: %v", batchNumberToSync, err)
			return err
		}
		start = time.Now()
		cbatches, lastStateRoot, err := s.processTrustedBatch(batchToSync, dbTx)
		metrics.ProcessTrustedBatchTime(time.Since(start))
		if err != nil {
			log.Errorf("syncTrustedState: error processing trusted batch %d: %v", batchNumberToSync, err)
			rollbackErr := dbTx.Rollback(s.ctx)
			if rollbackErr != nil {
				log.Errorf("syncTrustedState: error rolling back db transaction to sync trusted batch %d: %v", batchNumberToSync, rollbackErr)
				return rollbackErr
			}
			return err
		}
		log.Debug("syncTrustedState: Checking FlushID to commit trustedState data to db")
		err = s.checkFlushID(dbTx)
		if err != nil {
			log.Errorf("syncTrustedState: error checking flushID. Error: %v", err)
			rollbackErr := dbTx.Rollback(s.ctx)
			if rollbackErr != nil {
				log.Errorf("syncTrustedState: error rolling back state. RollbackErr: %s, Error : %v", rollbackErr.Error(), err)
				return rollbackErr
			}
			return err
		}

		if err := dbTx.Commit(s.ctx); err != nil {
			log.Errorf("syncTrustedState: error committing db transaction to sync trusted batch %v: %v", batchNumberToSync, err)
			return err
		}
		s.trustedState.lastTrustedBatches = cbatches
		s.trustedState.lastStateRoot = lastStateRoot
		batchNumberToSync++
	}

	log.Info("syncTrustedState: Trusted state fully synchronized")
	return nil
}

func (s *ClientSynchronizer) processBlockRange(blocks []etherman.Block, order map[common.Hash][]etherman.Order) error {
	// New info has to be included into the db using the state
	for i := range blocks {
		// Begin db transaction
		dbTx, err := s.state.BeginStateTransaction(s.ctx)
		if err != nil {
			log.Errorf("error creating db transaction to store block. BlockNumber: %d, error: %v", blocks[i].BlockNumber, err)
			return err
		}
		b := state.Block{
			BlockNumber: blocks[i].BlockNumber,
			BlockHash:   blocks[i].BlockHash,
			ParentHash:  blocks[i].ParentHash,
			ReceivedAt:  blocks[i].ReceivedAt,
		}
		// Add block information
		err = s.state.AddBlock(s.ctx, &b, dbTx)
		if err != nil {
			log.Errorf("error storing block. BlockNumber: %d, error: %v", blocks[i].BlockNumber, err)
			rollbackErr := dbTx.Rollback(s.ctx)
			if rollbackErr != nil {
				log.Errorf("error rolling back state to store block. BlockNumber: %d, rollbackErr: %s, error : %v", blocks[i].BlockNumber, rollbackErr.Error(), err)
				return rollbackErr
			}
			return err
		}
		for _, element := range order[blocks[i].BlockHash] {
			switch element.Name {
			case etherman.SequenceBatchesOrder:
				err = s.processSequenceBatches(blocks[i].SequencedBatches[element.Pos], blocks[i].BlockNumber, dbTx)
				if err != nil {
					return err
				}
			case etherman.ForcedBatchesOrder:
				err = s.processForcedBatch(blocks[i].ForcedBatches[element.Pos], dbTx)
				if err != nil {
					return err
				}
			case etherman.GlobalExitRootsOrder:
				err = s.processGlobalExitRoot(blocks[i].GlobalExitRoots[element.Pos], dbTx)
				if err != nil {
					return err
				}
			case etherman.SequenceForceBatchesOrder:
				err = s.processSequenceForceBatch(blocks[i].SequencedForceBatches[element.Pos], blocks[i], dbTx)
				if err != nil {
					return err
				}
			case etherman.TrustedVerifyBatchOrder:
				err = s.processTrustedVerifyBatches(blocks[i].VerifiedBatches[element.Pos], dbTx)
				if err != nil {
					return err
				}
			case etherman.ForkIDsOrder:
				err = s.processForkID(blocks[i].ForkIDs[element.Pos], blocks[i].BlockNumber, dbTx)
				if err != nil {
					return err
				}
			}
		}
		log.Debug("Checking FlushID to commit L1 data to db")
		err = s.checkFlushID(dbTx)
		if err != nil {
			log.Errorf("error checking flushID. Error: %v", err)
			rollbackErr := dbTx.Rollback(s.ctx)
			if rollbackErr != nil {
				log.Errorf("error rolling back state. RollbackErr: %s, Error : %v", rollbackErr.Error(), err)
				return rollbackErr
			}
			return err
		}
		err = dbTx.Commit(s.ctx)
		if err != nil {
			log.Errorf("error committing state to store block. BlockNumber: %d, err: %v", blocks[i].BlockNumber, err)
			rollbackErr := dbTx.Rollback(s.ctx)
			if rollbackErr != nil {
				log.Errorf("error rolling back state to store block. BlockNumber: %d, rollbackErr: %s, error : %v", blocks[i].BlockNumber, rollbackErr.Error(), err)
				return rollbackErr
			}
			return err
		}
	}
	return nil
}

// This function allows reset the state until an specific ethereum block
func (s *ClientSynchronizer) resetState(blockNumber uint64) error {
	log.Info("Reverting synchronization to block: ", blockNumber)
	dbTx, err := s.state.BeginStateTransaction(s.ctx)
	if err != nil {
		log.Error("error starting a db transaction to reset the state. Error: ", err)
		return err
	}
	err = s.state.Reset(s.ctx, blockNumber, dbTx)
	if err != nil {
		rollbackErr := dbTx.Rollback(s.ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state to store block. BlockNumber: %d, rollbackErr: %s, error : %v", blockNumber, rollbackErr.Error(), err)
			return rollbackErr
		}
		log.Error("error resetting the state. Error: ", err)
		return err
	}
	err = s.ethTxManager.Reorg(s.ctx, blockNumber+1, dbTx)
	if err != nil {
		rollbackErr := dbTx.Rollback(s.ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back eth tx manager when reorg detected. BlockNumber: %d, rollbackErr: %s, error : %v", blockNumber, rollbackErr.Error(), err)
			return rollbackErr
		}
		log.Error("error processing reorg on eth tx manager. Error: ", err)
		return err
	}
	err = dbTx.Commit(s.ctx)
	if err != nil {
		rollbackErr := dbTx.Rollback(s.ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state to store block. BlockNumber: %d, rollbackErr: %s, error : %v", blockNumber, rollbackErr.Error(), err)
			return rollbackErr
		}
		log.Error("error committing the resetted state. Error: ", err)
		return err
	}
	if s.l1SyncOrchestration != nil {
		s.l1SyncOrchestration.reset(blockNumber)
	}
	return nil
}

/*
This function will check if there is a reorg.
As input param needs the last ethereum block synced. Retrieve the block info from the blockchain
to compare it with the stored info. If hash and hash parent matches, then no reorg is detected and return a nil.
If hash or hash parent don't match, reorg detected and the function will return the block until the sync process
must be reverted. Then, check the previous ethereum block synced, get block info from the blockchain and check
hash and has parent. This operation has to be done until a match is found.
*/
func (s *ClientSynchronizer) checkReorg(latestBlock *state.Block) (*state.Block, error) {
	// This function only needs to worry about reorgs if some of the reorganized blocks contained rollup info.
	latestEthBlockSynced := *latestBlock
	var depth uint64
	for {
		block, err := s.etherMan.EthBlockByNumber(s.ctx, latestBlock.BlockNumber)
		if err != nil {
			log.Errorf("error getting latest block synced from blockchain. Block: %d, error: %v", latestBlock.BlockNumber, err)
			return nil, err
		}
		if block.NumberU64() != latestBlock.BlockNumber {
			err = fmt.Errorf("wrong ethereum block retrieved from blockchain. Block numbers don't match. BlockNumber stored: %d. BlockNumber retrieved: %d",
				latestBlock.BlockNumber, block.NumberU64())
			log.Error("error: ", err)
			return nil, err
		}
		// Compare hashes
		if (block.Hash() != latestBlock.BlockHash || block.ParentHash() != latestBlock.ParentHash) && latestBlock.BlockNumber > s.genesis.GenesisBlockNum {
			log.Infof("checkReorg: Bad block %d hashOk %t parentHashOk %t", latestBlock.BlockNumber, block.Hash() == latestBlock.BlockHash, block.ParentHash() == latestBlock.ParentHash)
			log.Debug("[checkReorg function] => latestBlockNumber: ", latestBlock.BlockNumber)
			log.Debug("[checkReorg function] => latestBlockHash: ", latestBlock.BlockHash)
			log.Debug("[checkReorg function] => latestBlockHashParent: ", latestBlock.ParentHash)
			log.Debug("[checkReorg function] => BlockNumber: ", latestBlock.BlockNumber, block.NumberU64())
			log.Debug("[checkReorg function] => BlockHash: ", block.Hash())
			log.Debug("[checkReorg function] => BlockHashParent: ", block.ParentHash())
			depth++
			log.Debug("REORG: Looking for the latest correct ethereum block. Depth: ", depth)
			// Reorg detected. Getting previous block
			dbTx, err := s.state.BeginStateTransaction(s.ctx)
			if err != nil {
				log.Errorf("error creating db transaction to get prevoius blocks")
				return nil, err
			}
			latestBlock, err = s.state.GetPreviousBlock(s.ctx, depth, dbTx)
			errC := dbTx.Commit(s.ctx)
			if errC != nil {
				log.Errorf("error committing dbTx, err: %v", errC)
				rollbackErr := dbTx.Rollback(s.ctx)
				if rollbackErr != nil {
					log.Errorf("error rolling back state. RollbackErr: %v", rollbackErr)
					return nil, rollbackErr
				}
				log.Errorf("error committing dbTx, err: %v", errC)
				return nil, errC
			}
			if errors.Is(err, state.ErrNotFound) {
				log.Warn("error checking reorg: previous block not found in db: ", err)
				return &state.Block{}, nil
			} else if err != nil {
				return nil, err
			}
		} else {
			break
		}
	}
	if latestEthBlockSynced.BlockHash != latestBlock.BlockHash {
		log.Info("Reorg detected in block: ", latestEthBlockSynced.BlockNumber, " last block OK: ", latestBlock.BlockNumber)
		return latestBlock, nil
	}
	return nil, nil
}

// Stop function stops the synchronizer
func (s *ClientSynchronizer) Stop() {
	s.cancelCtx()
}

func (s *ClientSynchronizer) checkTrustedState(batch state.Batch, tBatch *state.Batch, newRoot common.Hash, dbTx pgx.Tx) bool {
	//Compare virtual state with trusted state
	var reorgReasons strings.Builder
	if newRoot != tBatch.StateRoot {
		log.Warnf("Different field StateRoot. Virtual: %s, Trusted: %s\n", newRoot.String(), tBatch.StateRoot.String())
		reorgReasons.WriteString(fmt.Sprintf("Different field StateRoot. Virtual: %s, Trusted: %s\n", newRoot.String(), tBatch.StateRoot.String()))
	}
	if hex.EncodeToString(batch.BatchL2Data) != hex.EncodeToString(tBatch.BatchL2Data) {
		log.Warnf("Different field BatchL2Data. Virtual: %s, Trusted: %s\n", hex.EncodeToString(batch.BatchL2Data), hex.EncodeToString(tBatch.BatchL2Data))
		reorgReasons.WriteString(fmt.Sprintf("Different field BatchL2Data. Virtual: %s, Trusted: %s\n", hex.EncodeToString(batch.BatchL2Data), hex.EncodeToString(tBatch.BatchL2Data)))
	}
	if batch.GlobalExitRoot.String() != tBatch.GlobalExitRoot.String() {
		log.Warnf("Different field GlobalExitRoot. Virtual: %s, Trusted: %s\n", batch.GlobalExitRoot.String(), tBatch.GlobalExitRoot.String())
		reorgReasons.WriteString(fmt.Sprintf("Different field GlobalExitRoot. Virtual: %s, Trusted: %s\n", batch.GlobalExitRoot.String(), tBatch.GlobalExitRoot.String()))
	}
	if batch.Timestamp.Unix() != tBatch.Timestamp.Unix() {
		log.Warnf("Different field Timestamp. Virtual: %d, Trusted: %d\n", batch.Timestamp.Unix(), tBatch.Timestamp.Unix())
		reorgReasons.WriteString(fmt.Sprintf("Different field Timestamp. Virtual: %d, Trusted: %d\n", batch.Timestamp.Unix(), tBatch.Timestamp.Unix()))
	}
	if batch.Coinbase.String() != tBatch.Coinbase.String() {
		log.Warnf("Different field Coinbase. Virtual: %s, Trusted: %s\n", batch.Coinbase.String(), tBatch.Coinbase.String())
		reorgReasons.WriteString(fmt.Sprintf("Different field Coinbase. Virtual: %s, Trusted: %s\n", batch.Coinbase.String(), tBatch.Coinbase.String()))
	}

	if reorgReasons.Len() > 0 {
		reason := reorgReasons.String()

		if tBatch.StateRoot == (common.Hash{}) {
			log.Warnf("incomplete trusted batch %d detected. Syncing full batch from L1", tBatch.BatchNumber)
		} else {
			log.Warnf("missmatch in trusted state detected for Batch Number: %d. Reasons: %s", tBatch.BatchNumber, reason)
		}
		if s.isTrustedSequencer {
			s.halt(s.ctx, fmt.Errorf("TRUSTED REORG DETECTED! Batch: %d", batch.BatchNumber))
		}
		// Store trusted reorg register
		tr := state.TrustedReorg{
			BatchNumber: tBatch.BatchNumber,
			Reason:      reason,
		}
		err := s.state.AddTrustedReorg(s.ctx, &tr, dbTx)
		if err != nil {
			log.Error("error storing tursted reorg register into the db. Error: ", err)
		}
		return true
	}
	return false
}

func (s *ClientSynchronizer) processForkID(forkID etherman.ForkID, blockNumber uint64, dbTx pgx.Tx) error {
	fID := state.ForkIDInterval{
		FromBatchNumber: forkID.BatchNumber + 1,
		ToBatchNumber:   math.MaxUint64,
		ForkId:          forkID.ForkID,
		Version:         forkID.Version,
		BlockNumber:     blockNumber,
	}

	// If forkID affects to a batch from the past. State must be reseted.
	log.Debugf("ForkID: %d, synchronization must use the new forkID since batch: %d", forkID.ForkID, forkID.BatchNumber+1)
	fIds, err := s.state.GetForkIDs(s.ctx, dbTx)
	if err != nil {
		log.Error("error getting ForkIDTrustedReorg. Error: ", err)
		rollbackErr := dbTx.Rollback(s.ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state get forkID trusted state. BlockNumber: %d, rollbackErr: %s, error : %v", blockNumber, rollbackErr.Error(), err)
			return rollbackErr
		}
		return err
	}
	if len(fIds) != 0 && fIds[len(fIds)-1].ForkId == fID.ForkId { // If the forkID reset was already done
		return nil
	}
	//If the forkID.batchnumber is a future batch
	latestBatchNumber, err := s.state.GetLastBatchNumber(s.ctx, dbTx)
	if err != nil && !errors.Is(err, state.ErrStateNotSynchronized) {
		log.Error("error getting last batch number. Error: ", err)
		rollbackErr := dbTx.Rollback(s.ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state. BlockNumber: %d, rollbackErr: %s, error : %v", blockNumber, rollbackErr.Error(), err)
			return rollbackErr
		}
		return err
	}
	// Add new forkID to the state
	err = s.state.AddForkIDInterval(s.ctx, fID, dbTx)
	if err != nil {
		log.Error("error adding new forkID interval to the state. Error: ", err)
		rollbackErr := dbTx.Rollback(s.ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state to store block. BlockNumber: %d, rollbackErr: %s, error : %v", blockNumber, rollbackErr.Error(), err)
			return rollbackErr
		}
		return err
	}
	if latestBatchNumber <= forkID.BatchNumber || s.isTrustedSequencer { //If the forkID will start in a future batch or isTrustedSequencer
		log.Infof("Just adding forkID. Skipping reset forkID. ForkID: %+v.", fID)
		return nil
	}

	log.Info("ForkID received in the permissionless node that affects to a batch from the past")
	//Reset DB only if permissionless node
	log.Debugf("ForkID: %d, Reverting synchronization to batch: %d", forkID.ForkID, forkID.BatchNumber+1)
	err = s.state.ResetForkID(s.ctx, forkID.BatchNumber+1, dbTx)
	if err != nil {
		log.Error("error resetting the state. Error: ", err)
		rollbackErr := dbTx.Rollback(s.ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state to store block. BlockNumber: %d, rollbackErr: %s, error : %v", blockNumber, rollbackErr.Error(), err)
			return rollbackErr
		}
		return err
	}

	// Commit because it returns an error to force the resync
	err = dbTx.Commit(s.ctx)
	if err != nil {
		log.Error("error committing the resetted state. Error: ", err)
		rollbackErr := dbTx.Rollback(s.ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state to store block. BlockNumber: %d, rollbackErr: %s, error : %v", blockNumber, rollbackErr.Error(), err)
			return rollbackErr
		}
		return err
	}

	return fmt.Errorf("new ForkID detected, reseting synchronizarion")
}

func isZeroByteArray(bytesArray [32]byte) bool {
	var zero = [32]byte{}
	return bytes.Equal(bytesArray[:], zero[:])
}

func (s *ClientSynchronizer) processSequenceBatches(sequencedBatches []etherman.SequencedBatch, blockNumber uint64, dbTx pgx.Tx) error {
	if len(sequencedBatches) == 0 {
		log.Warn("Empty sequencedBatches array detected, ignoring...")
		return nil
	}
	for _, sbatch := range sequencedBatches {
		var batchL2Data []byte
		log.Infof("sbatch.Transactions len:%d, txs hash:%s", len(sbatch.Transactions), hex.EncodeToString(sbatch.TransactionsHash[:]))
		var err error
		if len(sbatch.Transactions) > 0 || (len(sbatch.Transactions) == 0 && isZeroByteArray(sbatch.TransactionsHash)) {
			batchL2Data = sbatch.Transactions
		} else {
			batchL2Data, err = s.getBatchL2Data(sbatch.BatchNumber, sbatch.TransactionsHash)
			if err != nil {
				return err
			}
		}

		virtualBatch := state.VirtualBatch{
			BatchNumber:   sbatch.BatchNumber,
			TxHash:        sbatch.TxHash,
			Coinbase:      sbatch.Coinbase,
			BlockNumber:   blockNumber,
			SequencerAddr: sbatch.SequencerAddr,
		}
		batch := state.Batch{
			BatchNumber:    sbatch.BatchNumber,
			GlobalExitRoot: sbatch.GlobalExitRoot,
			Timestamp:      time.Unix(int64(sbatch.Timestamp), 0),
			Coinbase:       sbatch.Coinbase,
			BatchL2Data:    batchL2Data,
		}
		// ForcedBatch must be processed
		if sbatch.MinForcedTimestamp > 0 { // If this is true means that the batch is forced
			log.Debug("FORCED BATCH SEQUENCED!")
			// Read forcedBatches from db
			forcedBatches, err := s.state.GetNextForcedBatches(s.ctx, 1, dbTx)
			if err != nil {
				log.Errorf("error getting forcedBatches. BatchNumber: %d", virtualBatch.BatchNumber)
				rollbackErr := dbTx.Rollback(s.ctx)
				if rollbackErr != nil {
					log.Errorf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %s, error : %v", virtualBatch.BatchNumber, blockNumber, rollbackErr.Error(), err)
					return rollbackErr
				}
				return err
			}
			if len(forcedBatches) == 0 {
				log.Errorf("error: empty forcedBatches array read from db. BatchNumber: %d", sbatch.BatchNumber)
				rollbackErr := dbTx.Rollback(s.ctx)
				if rollbackErr != nil {
					log.Errorf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %v", sbatch.BatchNumber, blockNumber, rollbackErr)
					return rollbackErr
				}
				return fmt.Errorf("error: empty forcedBatches array read from db. BatchNumber: %d", sbatch.BatchNumber)
			}
			if uint64(forcedBatches[0].ForcedAt.Unix()) != sbatch.MinForcedTimestamp ||
				forcedBatches[0].GlobalExitRoot != sbatch.GlobalExitRoot ||
				common.Bytes2Hex(forcedBatches[0].RawTxsData) != common.Bytes2Hex(batchL2Data) {
				log.Warnf("ForcedBatch stored: %+v. RawTxsData: %s", forcedBatches, common.Bytes2Hex(forcedBatches[0].RawTxsData))
				log.Warnf("ForcedBatch sequenced received: %+v. RawTxsData: %s", sbatch, common.Bytes2Hex(batchL2Data))
				log.Errorf("error: forcedBatch received doesn't match with the next expected forcedBatch stored in db. Expected: %+v, Synced: %+v", forcedBatches, sbatch)
				rollbackErr := dbTx.Rollback(s.ctx)
				if rollbackErr != nil {
					log.Errorf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %v", virtualBatch.BatchNumber, blockNumber, rollbackErr)
					return rollbackErr
				}
				return fmt.Errorf("error: forcedBatch received doesn't match with the next expected forcedBatch stored in db. Expected: %+v, Synced: %+v", forcedBatches, sbatch)
			}
			log.Debug("Setting forcedBatchNum: ", forcedBatches[0].ForcedBatchNumber)
			batch.ForcedBatchNum = &forcedBatches[0].ForcedBatchNumber
		}

		// Now we need to check the batch. ForcedBatches should be already stored in the batch table because this is done by the sequencer
		processCtx := state.ProcessingContext{
			BatchNumber:    batch.BatchNumber,
			Coinbase:       batch.Coinbase,
			Timestamp:      batch.Timestamp,
			GlobalExitRoot: batch.GlobalExitRoot,
			ForcedBatchNum: batch.ForcedBatchNum,
			BatchL2Data:    &batch.BatchL2Data,
		}

		var newRoot common.Hash

		// First get trusted batch from db
		tBatch, err := s.state.GetBatchByNumber(s.ctx, batch.BatchNumber, dbTx)
		if err != nil {
			if errors.Is(err, state.ErrNotFound) {
				log.Debugf("BatchNumber: %d, not found in trusted state. Storing it...", batch.BatchNumber)
				// If it is not found, store batch
				log.Infof("processSequenceBatches: (not found batch) ProcessAndStoreClosedBatch . BatchNumber: %d, BlockNumber: %d", processCtx.BatchNumber, blockNumber)
				newStateRoot, flushID, proverID, err := s.state.ProcessAndStoreClosedBatch(s.ctx, processCtx, batch.BatchL2Data, dbTx, stateMetrics.SynchronizerCallerLabel)
				if err != nil {
					log.Errorf("error storing trustedBatch. BatchNumber: %d, BlockNumber: %d, error: %v", batch.BatchNumber, blockNumber, err)
					rollbackErr := dbTx.Rollback(s.ctx)
					if rollbackErr != nil {
						log.Errorf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %s, error : %v", batch.BatchNumber, blockNumber, rollbackErr.Error(), err)
						return rollbackErr
					}
					log.Errorf("error storing batch. BatchNumber: %d, BlockNumber: %d, error: %v", batch.BatchNumber, blockNumber, err)
					return err
				}
				s.pendingFlushID(flushID, proverID)

				newRoot = newStateRoot
				tBatch = &batch
				tBatch.StateRoot = newRoot
			} else {
				log.Error("error checking trusted state: ", err)
				rollbackErr := dbTx.Rollback(s.ctx)
				if rollbackErr != nil {
					log.Errorf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %v", batch.BatchNumber, blockNumber, rollbackErr)
					return rollbackErr
				}
				return err
			}
		} else {
			// Reprocess batch to compare the stateRoot with tBatch.StateRoot and get accInputHash
			p, err := s.state.ExecuteBatch(s.ctx, batch, false, dbTx)
			if err != nil {
				log.Errorf("error executing L1 batch: %+v, error: %v", batch, err)
				rollbackErr := dbTx.Rollback(s.ctx)
				if rollbackErr != nil {
					log.Errorf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %s, error : %v", batch.BatchNumber, blockNumber, rollbackErr.Error(), err)
					return rollbackErr
				}
				return err
			}
			newRoot = common.BytesToHash(p.NewStateRoot)
			accumulatedInputHash := common.BytesToHash(p.NewAccInputHash)

			//AddAccumulatedInputHash
			err = s.state.AddAccumulatedInputHash(s.ctx, batch.BatchNumber, accumulatedInputHash, dbTx)
			if err != nil {
				log.Errorf("error adding accumulatedInputHash for batch: %d. Error; %v", batch.BatchNumber, err)
				rollbackErr := dbTx.Rollback(s.ctx)
				if rollbackErr != nil {
					log.Errorf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %v", batch.BatchNumber, blockNumber, rollbackErr)
					return rollbackErr
				}
				return err
			}
		}

		// Call the check trusted state method to compare trusted and virtual state
		status := s.checkTrustedState(batch, tBatch, newRoot, dbTx)
		if status {
			// Reorg Pool
			err := s.reorgPool(dbTx)
			if err != nil {
				rollbackErr := dbTx.Rollback(s.ctx)
				if rollbackErr != nil {
					log.Errorf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %s, error : %v", tBatch.BatchNumber, blockNumber, rollbackErr.Error(), err)
					return rollbackErr
				}
				log.Errorf("error: %v. BatchNumber: %d, BlockNumber: %d", err, tBatch.BatchNumber, blockNumber)
				return err
			}

			// Clean trustedState sync variables to avoid sync the trusted state from the wrong starting point.
			// This wrong starting point would force the trusted sync to clean the virtualization of the batch reaching an inconsistency.
			s.trustedState.lastTrustedBatches = nil
			s.trustedState.lastStateRoot = nil

			// Reset trusted state
			previousBatchNumber := batch.BatchNumber - 1
			if tBatch.StateRoot == (common.Hash{}) {
				log.Warnf("cleaning state before inserting batch from L1. Clean until batch: %d", previousBatchNumber)
			} else {
				log.Warnf("missmatch in trusted state detected, discarding batches until batchNum %d", previousBatchNumber)
			}
			log.Infof("ResetTrustedState: Resetting trusted state. delete batch > %d, ", previousBatchNumber)
			err = s.state.ResetTrustedState(s.ctx, previousBatchNumber, dbTx) // This method has to reset the forced batches deleting the batchNumber for higher batchNumbers
			if err != nil {
				log.Errorf("error resetting trusted state. BatchNumber: %d, BlockNumber: %d, error: %v", batch.BatchNumber, blockNumber, err)
				rollbackErr := dbTx.Rollback(s.ctx)
				if rollbackErr != nil {
					log.Errorf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %s, error : %v", batch.BatchNumber, blockNumber, rollbackErr.Error(), err)
					return rollbackErr
				}
				log.Errorf("error resetting trusted state. BatchNumber: %d, BlockNumber: %d, error: %v", batch.BatchNumber, blockNumber, err)
				return err
			}
			log.Infof("processSequenceBatches: (deleted previous) ProcessAndStoreClosedBatch . BatchNumber: %d, BlockNumber: %d", processCtx.BatchNumber, blockNumber)
			_, flushID, proverID, err := s.state.ProcessAndStoreClosedBatch(s.ctx, processCtx, batch.BatchL2Data, dbTx, stateMetrics.SynchronizerCallerLabel)
			if err != nil {
				log.Errorf("error storing trustedBatch. BatchNumber: %d, BlockNumber: %d, error: %v", batch.BatchNumber, blockNumber, err)
				rollbackErr := dbTx.Rollback(s.ctx)
				if rollbackErr != nil {
					log.Errorf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %s, error : %v", batch.BatchNumber, blockNumber, rollbackErr.Error(), err)
					return rollbackErr
				}
				log.Errorf("error storing batch. BatchNumber: %d, BlockNumber: %d, error: %v", batch.BatchNumber, blockNumber, err)
				return err
			}
			s.pendingFlushID(flushID, proverID)
		}

		// Store virtualBatch
		log.Infof("processSequenceBatches: Storing virtualBatch. BatchNumber: %d, BlockNumber: %d", virtualBatch.BatchNumber, blockNumber)
		err = s.state.AddVirtualBatch(s.ctx, &virtualBatch, dbTx)
		if err != nil {
			log.Errorf("error storing virtualBatch. BatchNumber: %d, BlockNumber: %d, error: %v", virtualBatch.BatchNumber, blockNumber, err)
			rollbackErr := dbTx.Rollback(s.ctx)
			if rollbackErr != nil {
				log.Errorf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %s, error : %v", virtualBatch.BatchNumber, blockNumber, rollbackErr.Error(), err)
				return rollbackErr
			}
			log.Errorf("error storing virtualBatch. BatchNumber: %d, BlockNumber: %d, error: %v", virtualBatch.BatchNumber, blockNumber, err)
			return err
		}
	}
	// Insert the sequence to allow the aggregator verify the sequence batches
	seq := state.Sequence{
		FromBatchNumber: sequencedBatches[0].BatchNumber,
		ToBatchNumber:   sequencedBatches[len(sequencedBatches)-1].BatchNumber,
	}
	err := s.state.AddSequence(s.ctx, seq, dbTx)
	if err != nil {
		log.Errorf("error adding sequence. Sequence: %+v", seq)
		rollbackErr := dbTx.Rollback(s.ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state. BlockNumber: %d, rollbackErr: %s, error : %v", blockNumber, rollbackErr.Error(), err)
			return rollbackErr
		}
		log.Errorf("error getting adding sequence. BlockNumber: %d, error: %v", blockNumber, err)
		return err
	}
	return nil
}

func (s *ClientSynchronizer) processSequenceForceBatch(sequenceForceBatch []etherman.SequencedForceBatch, block etherman.Block, dbTx pgx.Tx) error {
	if len(sequenceForceBatch) == 0 {
		log.Warn("Empty sequenceForceBatch array detected, ignoring...")
		return nil
	}
	// First, get last virtual batch number
	lastVirtualizedBatchNumber, err := s.state.GetLastVirtualBatchNum(s.ctx, dbTx)
	if err != nil {
		log.Errorf("error getting lastVirtualBatchNumber. BlockNumber: %d, error: %v", block.BlockNumber, err)
		rollbackErr := dbTx.Rollback(s.ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %s, error : %v", lastVirtualizedBatchNumber, block.BlockNumber, rollbackErr.Error(), err)
			return rollbackErr
		}
		log.Errorf("error getting lastVirtualBatchNumber. BlockNumber: %d, error: %v", block.BlockNumber, err)
		return err
	}
	// Clean trustedState sync variables to avoid sync the trusted state from the wrong starting point.
	// This wrong starting point would force the trusted sync to clean the virtualization of the batch reaching an inconsistency.
	s.trustedState.lastTrustedBatches = nil
	s.trustedState.lastStateRoot = nil

	// Reset trusted state
	log.Infof("ResetTrustedState: processSequenceForceBatch: Resetting trusted state. delete batch > (lastVirtualizedBatchNumber)%d, ", lastVirtualizedBatchNumber)
	err = s.state.ResetTrustedState(s.ctx, lastVirtualizedBatchNumber, dbTx) // This method has to reset the forced batches deleting the batchNumber for higher batchNumbers
	if err != nil {
		log.Errorf("error resetting trusted state. BatchNumber: %d, BlockNumber: %d, error: %v", lastVirtualizedBatchNumber, block.BlockNumber, err)
		rollbackErr := dbTx.Rollback(s.ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %s, error : %v", lastVirtualizedBatchNumber, block.BlockNumber, rollbackErr.Error(), err)
			return rollbackErr
		}
		log.Errorf("error resetting trusted state. BatchNumber: %d, BlockNumber: %d, error: %v", lastVirtualizedBatchNumber, block.BlockNumber, err)
		return err
	}
	// Read forcedBatches from db
	forcedBatches, err := s.state.GetNextForcedBatches(s.ctx, len(sequenceForceBatch), dbTx)
	if err != nil {
		log.Errorf("error getting forcedBatches in processSequenceForceBatch. BlockNumber: %d", block.BlockNumber)
		rollbackErr := dbTx.Rollback(s.ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state. BlockNumber: %d, rollbackErr: %s, error : %v", block.BlockNumber, rollbackErr.Error(), err)
			return rollbackErr
		}
		log.Errorf("error getting forcedBatches in processSequenceForceBatch. BlockNumber: %d, error: %v", block.BlockNumber, err)
		return err
	}
	if len(sequenceForceBatch) != len(forcedBatches) {
		rollbackErr := dbTx.Rollback(s.ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state. BlockNumber: %d, rollbackErr: %v", block.BlockNumber, rollbackErr)
			return rollbackErr
		}
		log.Error("error number of forced batches doesn't match")
		return fmt.Errorf("error number of forced batches doesn't match")
	}
	for i, fbatch := range sequenceForceBatch {
		if uint64(forcedBatches[i].ForcedAt.Unix()) != fbatch.MinForcedTimestamp ||
			forcedBatches[i].GlobalExitRoot != fbatch.GlobalExitRoot ||
			common.Bytes2Hex(forcedBatches[i].RawTxsData) != common.Bytes2Hex(fbatch.Transactions) {
			log.Warnf("ForcedBatch stored: %+v", forcedBatches)
			log.Warnf("ForcedBatch sequenced received: %+v", fbatch)
			log.Errorf("error: forcedBatch received doesn't match with the next expected forcedBatch stored in db. Expected: %+v, Synced: %+v", forcedBatches[i], fbatch)
			rollbackErr := dbTx.Rollback(s.ctx)
			if rollbackErr != nil {
				log.Errorf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %v", fbatch.BatchNumber, block.BlockNumber, rollbackErr)
				return rollbackErr
			}
			return fmt.Errorf("error: forcedBatch received doesn't match with the next expected forcedBatch stored in db. Expected: %+v, Synced: %+v", forcedBatches[i], fbatch)
		}
		virtualBatch := state.VirtualBatch{
			BatchNumber:   fbatch.BatchNumber,
			TxHash:        fbatch.TxHash,
			Coinbase:      fbatch.Coinbase,
			SequencerAddr: fbatch.Coinbase,
			BlockNumber:   block.BlockNumber,
		}
		batch := state.ProcessingContext{
			BatchNumber:    fbatch.BatchNumber,
			GlobalExitRoot: fbatch.GlobalExitRoot,
			Timestamp:      block.ReceivedAt,
			Coinbase:       fbatch.Coinbase,
			ForcedBatchNum: &forcedBatches[i].ForcedBatchNumber,
			BatchL2Data:    &forcedBatches[i].RawTxsData,
		}
		// Process batch
		log.Infof("processSequenceFoceBatches: ProcessAndStoreClosedBatch . BatchNumber: %d, BlockNumber: %d", batch.BatchNumber, block.BlockNumber)
		_, flushID, proverID, err := s.state.ProcessAndStoreClosedBatch(s.ctx, batch, forcedBatches[i].RawTxsData, dbTx, stateMetrics.SynchronizerCallerLabel)
		if err != nil {
			log.Errorf("error processing batch in processSequenceForceBatch. BatchNumber: %d, BlockNumber: %d, error: %v", batch.BatchNumber, block.BlockNumber, err)
			rollbackErr := dbTx.Rollback(s.ctx)
			if rollbackErr != nil {
				log.Errorf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %s, error : %v", batch.BatchNumber, block.BlockNumber, rollbackErr.Error(), err)
				return rollbackErr
			}
			log.Errorf("error processing batch in processSequenceForceBatch. BatchNumber: %d, BlockNumber: %d, error: %v", batch.BatchNumber, block.BlockNumber, err)
			return err
		}
		s.pendingFlushID(flushID, proverID)

		// Store virtualBatch
		log.Infof("processSequenceFoceBatches: Storing virtualBatch. BatchNumber: %d, BlockNumber: %d", virtualBatch.BatchNumber, block.BlockNumber)
		err = s.state.AddVirtualBatch(s.ctx, &virtualBatch, dbTx)
		if err != nil {
			log.Errorf("error storing virtualBatch in processSequenceForceBatch. BatchNumber: %d, BlockNumber: %d, error: %v", virtualBatch.BatchNumber, block.BlockNumber, err)
			rollbackErr := dbTx.Rollback(s.ctx)
			if rollbackErr != nil {
				log.Errorf("error rolling back state. BatchNumber: %d, BlockNumber: %d, rollbackErr: %s, error : %v", virtualBatch.BatchNumber, block.BlockNumber, rollbackErr.Error(), err)
				return rollbackErr
			}
			log.Errorf("error storing virtualBatch in processSequenceForceBatch. BatchNumber: %d, BlockNumber: %d, error: %v", virtualBatch.BatchNumber, block.BlockNumber, err)
			return err
		}
	}
	// Insert the sequence to allow the aggregator verify the sequence batches
	seq := state.Sequence{
		FromBatchNumber: sequenceForceBatch[0].BatchNumber,
		ToBatchNumber:   sequenceForceBatch[len(sequenceForceBatch)-1].BatchNumber,
	}
	err = s.state.AddSequence(s.ctx, seq, dbTx)
	if err != nil {
		log.Errorf("error adding sequence. Sequence: %+v", seq)
		rollbackErr := dbTx.Rollback(s.ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state. BlockNumber: %d, rollbackErr: %s, error : %v", block.BlockNumber, rollbackErr.Error(), err)
			return rollbackErr
		}
		log.Errorf("error getting adding sequence. BlockNumber: %d, error: %v", block.BlockNumber, err)
		return err
	}
	return nil
}

func (s *ClientSynchronizer) processForcedBatch(forcedBatch etherman.ForcedBatch, dbTx pgx.Tx) error {
	// Store forced batch into the db
	forcedB := state.ForcedBatch{
		BlockNumber:       forcedBatch.BlockNumber,
		ForcedBatchNumber: forcedBatch.ForcedBatchNumber,
		Sequencer:         forcedBatch.Sequencer,
		GlobalExitRoot:    forcedBatch.GlobalExitRoot,
		RawTxsData:        forcedBatch.RawTxsData,
		ForcedAt:          forcedBatch.ForcedAt,
	}
	log.Infof("processForcedBatch: Storing forcedBatch. BatchNumber: %d  BlockNumber: %d", forcedBatch.ForcedBatchNumber, forcedBatch.BlockNumber)
	err := s.state.AddForcedBatch(s.ctx, &forcedB, dbTx)
	if err != nil {
		log.Errorf("error storing the forcedBatch in processForcedBatch. BlockNumber: %d", forcedBatch.BlockNumber)
		rollbackErr := dbTx.Rollback(s.ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state. BlockNumber: %d, rollbackErr: %s, error : %v", forcedBatch.BlockNumber, rollbackErr.Error(), err)
			return rollbackErr
		}
		log.Errorf("error storing the forcedBatch in processForcedBatch. BlockNumber: %d, error: %v", forcedBatch.BlockNumber, err)
		return err
	}
	return nil
}

func (s *ClientSynchronizer) processGlobalExitRoot(globalExitRoot etherman.GlobalExitRoot, dbTx pgx.Tx) error {
	// Store GlobalExitRoot
	ger := state.GlobalExitRoot{
		BlockNumber:     globalExitRoot.BlockNumber,
		MainnetExitRoot: globalExitRoot.MainnetExitRoot,
		RollupExitRoot:  globalExitRoot.RollupExitRoot,
		GlobalExitRoot:  globalExitRoot.GlobalExitRoot,
	}
	err := s.state.AddGlobalExitRoot(s.ctx, &ger, dbTx)
	if err != nil {
		log.Errorf("error storing the globalExitRoot in processGlobalExitRoot. BlockNumber: %d", globalExitRoot.BlockNumber)
		rollbackErr := dbTx.Rollback(s.ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state. BlockNumber: %d, rollbackErr: %s, error : %v", globalExitRoot.BlockNumber, rollbackErr.Error(), err)
			return rollbackErr
		}
		log.Errorf("error storing the GlobalExitRoot in processGlobalExitRoot. BlockNumber: %d, error: %v", globalExitRoot.BlockNumber, err)
		return err
	}
	return nil
}

func (s *ClientSynchronizer) processTrustedVerifyBatches(lastVerifiedBatch etherman.VerifiedBatch, dbTx pgx.Tx) error {
	lastVBatch, err := s.state.GetLastVerifiedBatch(s.ctx, dbTx)
	if err != nil {
		log.Errorf("error getting lastVerifiedBatch stored in db in processTrustedVerifyBatches. Processing synced blockNumber: %d", lastVerifiedBatch.BlockNumber)
		rollbackErr := dbTx.Rollback(s.ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state. Processing synced blockNumber: %d, rollbackErr: %s, error : %v", lastVerifiedBatch.BlockNumber, rollbackErr.Error(), err)
			return rollbackErr
		}
		log.Errorf("error getting lastVerifiedBatch stored in db in processTrustedVerifyBatches. Processing synced blockNumber: %d, error: %v", lastVerifiedBatch.BlockNumber, err)
		return err
	}
	nbatches := lastVerifiedBatch.BatchNumber - lastVBatch.BatchNumber
	batch, err := s.state.GetBatchByNumber(s.ctx, lastVerifiedBatch.BatchNumber, dbTx)
	if err != nil {
		log.Errorf("error getting GetBatchByNumber stored in db in processTrustedVerifyBatches. Processing batchNumber: %d", lastVerifiedBatch.BatchNumber)
		rollbackErr := dbTx.Rollback(s.ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state. Processing batchNumber: %d, rollbackErr: %s, error : %v", lastVerifiedBatch.BatchNumber, rollbackErr.Error(), err)
			return rollbackErr
		}
		log.Errorf("error getting GetBatchByNumber stored in db in processTrustedVerifyBatches. Processing batchNumber: %d, error: %v", lastVerifiedBatch.BatchNumber, err)
		return err
	}

	// Checks that calculated state root matches with the verified state root in the smc
	if batch.StateRoot != lastVerifiedBatch.StateRoot {
		log.Warn("nbatches: ", nbatches)
		log.Warnf("Batch from db: %+v", batch)
		log.Warnf("Verified Batch: %+v", lastVerifiedBatch)
		log.Errorf("error: stateRoot calculated and state root verified don't match in processTrustedVerifyBatches. Processing batchNumber: %d", lastVerifiedBatch.BatchNumber)
		rollbackErr := dbTx.Rollback(s.ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state. Processing batchNumber: %d, rollbackErr: %v", lastVerifiedBatch.BatchNumber, rollbackErr)
			return rollbackErr
		}
		log.Errorf("error: stateRoot calculated and state root verified don't match in processTrustedVerifyBatches. Processing batchNumber: %d", lastVerifiedBatch.BatchNumber)
		return fmt.Errorf("error: stateRoot calculated and state root verified don't match in processTrustedVerifyBatches. Processing batchNumber: %d", lastVerifiedBatch.BatchNumber)
	}
	var i uint64
	for i = 1; i <= nbatches; i++ {
		verifiedB := state.VerifiedBatch{
			BlockNumber: lastVerifiedBatch.BlockNumber,
			BatchNumber: lastVBatch.BatchNumber + i,
			Aggregator:  lastVerifiedBatch.Aggregator,
			StateRoot:   lastVerifiedBatch.StateRoot,
			TxHash:      lastVerifiedBatch.TxHash,
			IsTrusted:   true,
		}
		log.Infof("processTrustedVerifyBatches: Storing verifiedB. BlockNumber: %d, BatchNumber: %d", verifiedB.BlockNumber, verifiedB.BatchNumber)
		err = s.state.AddVerifiedBatch(s.ctx, &verifiedB, dbTx)
		if err != nil {
			log.Errorf("error storing the verifiedB in processTrustedVerifyBatches. verifiedBatch: %+v, lastVerifiedBatch: %+v", verifiedB, lastVerifiedBatch)
			rollbackErr := dbTx.Rollback(s.ctx)
			if rollbackErr != nil {
				log.Errorf("error rolling back state. BlockNumber: %d, rollbackErr: %s, error : %v", lastVerifiedBatch.BlockNumber, rollbackErr.Error(), err)
				return rollbackErr
			}
			log.Errorf("error storing the verifiedB in processTrustedVerifyBatches. BlockNumber: %d, error: %v", lastVerifiedBatch.BlockNumber, err)
			return err
		}
	}
	return nil
}

func (s *ClientSynchronizer) processTrustedBatch(trustedBatch *types.Batch, dbTx pgx.Tx) ([]*state.Batch, *common.Hash, error) {
	log.Debugf("Processing trusted batch: %d", uint64(trustedBatch.Number))
	trustedBatchL2Data := trustedBatch.BatchL2Data
	batches := s.trustedState.lastTrustedBatches
	log.Debug("len(batches): ", len(batches))
	batches, err := s.getCurrentBatches(batches, trustedBatch, dbTx)
	if err != nil {
		log.Error("error getting currentBatches. Error: ", err)
		return nil, nil, err
	}

	if batches[0] != nil && (((trustedBatch.StateRoot == common.Hash{}) && (batches[0].StateRoot != common.Hash{})) ||
		len(batches[0].BatchL2Data) > len(trustedBatchL2Data)) {
		log.Error("error: inconsistency in data received from trustedNode")
		log.Infof("BatchNumber. stored: %d. synced: %d", batches[0].BatchNumber, uint64(trustedBatch.Number))
		log.Infof("GlobalExitRoot. stored:  %s. synced: %s", batches[0].GlobalExitRoot.String(), trustedBatch.GlobalExitRoot.String())
		log.Infof("LocalExitRoot. stored:  %s. synced: %s", batches[0].LocalExitRoot.String(), trustedBatch.LocalExitRoot.String())
		log.Infof("StateRoot. stored:  %s. synced: %s", batches[0].StateRoot.String(), trustedBatch.StateRoot.String())
		log.Infof("Coinbase. stored:  %s. synced: %s", batches[0].Coinbase.String(), trustedBatch.Coinbase.String())
		log.Infof("Timestamp. stored:  %d. synced: %d", uint64(batches[0].Timestamp.Unix()), uint64(trustedBatch.Timestamp))
		log.Infof("BatchL2Data. stored: %s. synced: %s", common.Bytes2Hex(batches[0].BatchL2Data), common.Bytes2Hex(trustedBatchL2Data))
		return nil, nil, fmt.Errorf("error: inconsistency in data received from trustedNode")
	}

	if s.trustedState.lastStateRoot == nil && (batches[0] == nil || (batches[0].StateRoot == common.Hash{})) {
		log.Debug("Setting stateRoot of previous batch. StateRoot: ", batches[1].StateRoot)
		// Previous synchronization incomplete. Needs to reprocess all txs again
		s.trustedState.lastStateRoot = &batches[1].StateRoot
	} else if batches[0] != nil && (batches[0].StateRoot != common.Hash{}) {
		// Previous synchronization completed
		s.trustedState.lastStateRoot = &batches[0].StateRoot
	}

	request := state.ProcessRequest{
		BatchNumber:     uint64(trustedBatch.Number),
		OldStateRoot:    *s.trustedState.lastStateRoot,
		OldAccInputHash: batches[1].AccInputHash,
		Coinbase:        common.HexToAddress(trustedBatch.Coinbase.String()),
		Timestamp:       time.Unix(int64(trustedBatch.Timestamp), 0),
	}
	// check if batch needs to be synchronized
	if batches[0] != nil {
		if checkIfSynced(batches, trustedBatch) {
			log.Debugf("Batch %d already synchronized", uint64(trustedBatch.Number))
			return batches, s.trustedState.lastStateRoot, nil
		}
		log.Infof("Batch %d needs to be updated", uint64(trustedBatch.Number))

		// Find txs to be processed and included in the trusted state
		if *s.trustedState.lastStateRoot == batches[1].StateRoot {
			prevBatch := uint64(trustedBatch.Number) - 1
			log.Infof("ResetTrustedState: processTrustedBatch: %d Cleaning state until batch:%d  ", trustedBatch.Number, prevBatch)
			// Delete txs that were stored before restart. We need to reprocess all txs because the intermediary stateRoot is only stored in memory
			err := s.state.ResetTrustedState(s.ctx, prevBatch, dbTx)
			if err != nil {
				log.Error("error resetting trusted state. Error: ", err)
				return nil, nil, err
			}
			// All txs need to be processed
			request.Transactions = trustedBatchL2Data
			// Reopen batch
			err = s.openBatch(trustedBatch, dbTx)
			if err != nil {
				log.Error("error openning batch. Error: ", err)
				return nil, nil, err
			}
			request.GlobalExitRoot = trustedBatch.GlobalExitRoot
			request.Transactions = trustedBatchL2Data
		} else {
			// Only new txs need to be processed
			storedTxs, syncedTxs, _, syncedEfficiencyPercentages, err := s.decodeTxs(trustedBatchL2Data, batches)
			if err != nil {
				return nil, nil, err
			}
			if len(storedTxs) < len(syncedTxs) {
				forkID := s.state.GetForkIDByBatchNumber(batches[0].BatchNumber)
				txsToBeAdded := syncedTxs[len(storedTxs):]
				if forkID >= forkID5 {
					syncedEfficiencyPercentages = syncedEfficiencyPercentages[len(storedTxs):]
				}
				log.Infof("Processing %d new txs with forkID: %d", len(txsToBeAdded), forkID)

				request.Transactions, err = state.EncodeTransactions(txsToBeAdded, syncedEfficiencyPercentages, forkID)
				if err != nil {
					log.Error("error encoding txs (%d) to be added to the state. Error: %v", len(txsToBeAdded), err)
					return nil, nil, err
				}
				log.Debug("request.Transactions: ", common.Bytes2Hex(request.Transactions))
			} else {
				log.Info("Nothing to sync. Node updated. Checking if it is closed")
				isBatchClosed := trustedBatch.StateRoot.String() != state.ZeroHash.String()
				if isBatchClosed {
					//Sanity check
					if s.trustedState.lastStateRoot != nil && trustedBatch.StateRoot != *s.trustedState.lastStateRoot {
						log.Errorf("batch %d, different batchL2Datas (trustedBatchL2Data: %s, batches[0].BatchL2Data: %s). Decoded txs are len(storedTxs): %d, len(syncedTxs): %d", uint64(trustedBatch.Number), trustedBatchL2Data.Hex(), "0x"+common.Bytes2Hex(batches[0].BatchL2Data), len(storedTxs), len(syncedTxs))
						for _, tx := range storedTxs {
							log.Error("stored txHash : ", tx.Hash())
						}
						for _, tx := range syncedTxs {
							log.Error("synced txHash : ", tx.Hash())
						}
						log.Errorf("batch: %d, stateRoot calculated (%s) is different from the stateRoot (%s) received during the trustedState synchronization", uint64(trustedBatch.Number), *s.trustedState.lastStateRoot, trustedBatch.StateRoot)
						return nil, nil, fmt.Errorf("batch: %d, stateRoot calculated (%s) is different from the stateRoot (%s) received during the trustedState synchronization", uint64(trustedBatch.Number), *s.trustedState.lastStateRoot, trustedBatch.StateRoot)
					}
					receipt := state.ProcessingReceipt{
						BatchNumber:   uint64(trustedBatch.Number),
						StateRoot:     trustedBatch.StateRoot,
						LocalExitRoot: trustedBatch.LocalExitRoot,
						BatchL2Data:   trustedBatchL2Data,
						AccInputHash:  trustedBatch.AccInputHash,
					}
					log.Debugf("closing batch %d", uint64(trustedBatch.Number))
					if err := s.state.CloseBatch(s.ctx, receipt, dbTx); err != nil {
						// This is a workaround to avoid closing a batch that was already closed
						if err.Error() != state.ErrBatchAlreadyClosed.Error() {
							log.Errorf("error closing batch %d", uint64(trustedBatch.Number))
							return nil, nil, err
						} else {
							log.Warnf("CASE 02: the batch [%d] was already closed", uint64(trustedBatch.Number))
							log.Info("batches[0].BatchNumber: ", batches[0].BatchNumber)
							log.Info("batches[0].AccInputHash: ", batches[0].AccInputHash)
							log.Info("batches[0].StateRoot: ", batches[0].StateRoot)
							log.Info("batches[0].LocalExitRoot: ", batches[0].LocalExitRoot)
							log.Info("batches[0].GlobalExitRoot: ", batches[0].GlobalExitRoot)
							log.Info("batches[0].Coinbase: ", batches[0].Coinbase)
							log.Info("batches[0].ForcedBatchNum: ", batches[0].ForcedBatchNum)
							log.Info("####################################")
							log.Info("batches[1].BatchNumber: ", batches[1].BatchNumber)
							log.Info("batches[1].AccInputHash: ", batches[1].AccInputHash)
							log.Info("batches[1].StateRoot: ", batches[1].StateRoot)
							log.Info("batches[1].LocalExitRoot: ", batches[1].LocalExitRoot)
							log.Info("batches[1].GlobalExitRoot: ", batches[1].GlobalExitRoot)
							log.Info("batches[1].Coinbase: ", batches[1].Coinbase)
							log.Info("batches[1].ForcedBatchNum: ", batches[1].ForcedBatchNum)
							log.Info("###############################")
							log.Info("trustedBatch.BatchNumber: ", trustedBatch.Number)
							log.Info("trustedBatch.AccInputHash: ", trustedBatch.AccInputHash)
							log.Info("trustedBatch.StateRoot: ", trustedBatch.StateRoot)
							log.Info("trustedBatch.LocalExitRoot: ", trustedBatch.LocalExitRoot)
							log.Info("trustedBatch.GlobalExitRoot: ", trustedBatch.GlobalExitRoot)
							log.Info("trustedBatch.Coinbase: ", trustedBatch.Coinbase)
							log.Info("trustedBatch.ForcedBatchNum: ", trustedBatch.ForcedBatchNumber)
						}
					}
					batches[0].AccInputHash = trustedBatch.AccInputHash
					batches[0].StateRoot = trustedBatch.StateRoot
					batches[0].LocalExitRoot = trustedBatch.LocalExitRoot
				}
				return batches, &trustedBatch.StateRoot, nil
			}
		}
		// Update batchL2Data
		err := s.state.UpdateBatchL2Data(s.ctx, batches[0].BatchNumber, trustedBatchL2Data, dbTx)
		if err != nil {
			log.Errorf("error opening batch %d", uint64(trustedBatch.Number))
			return nil, nil, err
		}
		batches[0].BatchL2Data = trustedBatchL2Data
		log.Debug("BatchL2Data updated for batch: ", batches[0].BatchNumber)
	} else {
		log.Infof("Batch %d needs to be synchronized", uint64(trustedBatch.Number))
		err := s.openBatch(trustedBatch, dbTx)
		if err != nil {
			log.Error("error openning batch. Error: ", err)
			return nil, nil, err
		}
		request.GlobalExitRoot = trustedBatch.GlobalExitRoot
		request.Transactions = trustedBatchL2Data
	}

	log.Debugf("Processing sequencer for batch %d", uint64(trustedBatch.Number))

	processBatchResp, err := s.processAndStoreTxs(trustedBatch, request, dbTx)
	if err != nil {
		log.Error("error procesingAndStoringTxs. Error: ", err)
		return nil, nil, err
	}

	log.Debug("TrustedBatch.StateRoot ", trustedBatch.StateRoot)
	isBatchClosed := trustedBatch.StateRoot.String() != state.ZeroHash.String()
	if isBatchClosed {
		//Sanity check
		if trustedBatch.StateRoot != processBatchResp.NewStateRoot {
			log.Error("trustedBatchL2Data: ", trustedBatchL2Data)
			log.Error("request.Transactions: ", request.Transactions)
			log.Errorf("batch: %d after processing some txs, stateRoot calculated (%s) is different from the stateRoot (%s) received during the trustedState synchronization", uint64(trustedBatch.Number), processBatchResp.NewStateRoot.String(), trustedBatch.StateRoot.String())
			return nil, nil, fmt.Errorf("batch: %d, stateRoot calculated (%s) is different from the stateRoot (%s) received during the trustedState synchronization", uint64(trustedBatch.Number), processBatchResp.NewStateRoot.String(), trustedBatch.StateRoot.String())
		}
		receipt := state.ProcessingReceipt{
			BatchNumber:   uint64(trustedBatch.Number),
			StateRoot:     processBatchResp.NewStateRoot,
			LocalExitRoot: processBatchResp.NewLocalExitRoot,
			BatchL2Data:   trustedBatchL2Data,
			AccInputHash:  trustedBatch.AccInputHash,
		}

		log.Debugf("closing batch %d", uint64(trustedBatch.Number))
		if err := s.state.CloseBatch(s.ctx, receipt, dbTx); err != nil {
			// This is a workarround to avoid closing a batch that was already closed
			if err.Error() != state.ErrBatchAlreadyClosed.Error() {
				log.Errorf("error closing batch %d", uint64(trustedBatch.Number))
				return nil, nil, err
			} else {
				log.Warnf("CASE 01: batch [%d] was already closed", uint64(trustedBatch.Number))
			}
		}
		log.Info("Batch closed right after processing some tx")
		if batches[0] != nil {
			log.Debug("Updating batches[0] values...")
			batches[0].AccInputHash = trustedBatch.AccInputHash
			batches[0].StateRoot = trustedBatch.StateRoot
			batches[0].LocalExitRoot = trustedBatch.LocalExitRoot
			batches[0].BatchL2Data = trustedBatchL2Data
		}
	}

	log.Infof("Batch %d synchronized", uint64(trustedBatch.Number))
	return batches, &processBatchResp.NewStateRoot, nil
}

func (s *ClientSynchronizer) reorgPool(dbTx pgx.Tx) error {
	latestBatchNum, err := s.etherMan.GetLatestBatchNumber()
	if err != nil {
		log.Error("error getting the latestBatchNumber virtualized in the smc. Error: ", err)
		return err
	}
	batchNumber := latestBatchNum + 1
	// Get transactions that have to be included in the pool again
	txs, err := s.state.GetReorgedTransactions(s.ctx, batchNumber, dbTx)
	if err != nil {
		log.Errorf("error getting txs from trusted state. BatchNumber: %d, error: %v", batchNumber, err)
		return err
	}
	log.Debug("Reorged transactions: ", txs)

	// Remove txs from the pool
	err = s.pool.DeleteReorgedTransactions(s.ctx, txs)
	if err != nil {
		log.Errorf("error deleting txs from the pool. BatchNumber: %d, error: %v", batchNumber, err)
		return err
	}
	log.Debug("Delete reorged transactions")

	// Add txs to the pool
	for _, tx := range txs {
		// Insert tx in WIP status to avoid the sequencer to grab them before it gets restarted
		// When the sequencer restarts, it will update the status to pending non-wip
		err = s.pool.StoreTx(s.ctx, *tx, "", true)
		if err != nil {
			log.Errorf("error storing tx into the pool again. TxHash: %s. BatchNumber: %d, error: %v", tx.Hash().String(), batchNumber, err)
			return err
		}
		log.Debug("Reorged transactions inserted in the pool: ", tx.Hash())
	}
	return nil
}

func (s *ClientSynchronizer) processAndStoreTxs(trustedBatch *types.Batch, request state.ProcessRequest, dbTx pgx.Tx) (*state.ProcessBatchResponse, error) {
	processBatchResp, err := s.state.ProcessBatch(s.ctx, request, true)
	if err != nil {
		log.Errorf("error processing sequencer batch for batch: %v", trustedBatch.Number)
		return nil, err
	}
	s.pendingFlushID(processBatchResp.FlushID, processBatchResp.ProverID)

	log.Debugf("Storing transactions %d for batch %v", len(processBatchResp.Responses), trustedBatch.Number)
	if processBatchResp.IsExecutorLevelError {
		log.Warn("executorLevelError detected. Avoid store txs...")
		return processBatchResp, nil
	} else if processBatchResp.IsRomOOCError {
		log.Warn("romOOCError detected. Avoid store txs...")
		return processBatchResp, nil
	}
	for _, tx := range processBatchResp.Responses {
		if state.IsStateRootChanged(executor.RomErrorCode(tx.RomError)) {
			log.Infof("TrustedBatch info: %+v", processBatchResp)
			log.Infof("Storing trusted tx %+v", tx)
			if _, err = s.state.StoreTransaction(s.ctx, uint64(trustedBatch.Number), tx, trustedBatch.Coinbase, uint64(trustedBatch.Timestamp), nil, dbTx); err != nil {
				log.Errorf("failed to store transactions for batch: %v. Tx: %s", trustedBatch.Number, tx.TxHash.String())
				return nil, err
			}
		}
	}
	return processBatchResp, nil
}

func (s *ClientSynchronizer) openBatch(trustedBatch *types.Batch, dbTx pgx.Tx) error {
	log.Debugf("Opening batch %d", trustedBatch.Number)
	var batchL2Data []byte = trustedBatch.BatchL2Data
	processCtx := state.ProcessingContext{
		BatchNumber:    uint64(trustedBatch.Number),
		Coinbase:       common.HexToAddress(trustedBatch.Coinbase.String()),
		Timestamp:      time.Unix(int64(trustedBatch.Timestamp), 0),
		GlobalExitRoot: trustedBatch.GlobalExitRoot,
		BatchL2Data:    &batchL2Data,
	}
	if trustedBatch.ForcedBatchNumber != nil {
		fb := uint64(*trustedBatch.ForcedBatchNumber)
		processCtx.ForcedBatchNum = &fb
	}
	err := s.state.OpenBatch(s.ctx, processCtx, dbTx)
	if err != nil {
		log.Error("error opening batch: ", trustedBatch.Number)
		return err
	}
	return nil
}

func (s *ClientSynchronizer) decodeTxs(trustedBatchL2Data types.ArgBytes, batches []*state.Batch) ([]ethTypes.Transaction, []ethTypes.Transaction, []uint8, []uint8, error) {
	forkID := s.state.GetForkIDByBatchNumber(batches[0].BatchNumber)
	syncedTxs, _, syncedEfficiencyPercentages, err := state.DecodeTxs(trustedBatchL2Data, forkID)
	if err != nil {
		log.Errorf("error decoding synced txs from trustedstate. Error: %v, TrustedBatchL2Data: %s", err, trustedBatchL2Data.Hex())
		return nil, nil, nil, nil, err
	}
	storedTxs, _, storedEfficiencyPercentages, err := state.DecodeTxs(batches[0].BatchL2Data, forkID)
	if err != nil {
		log.Errorf("error decoding stored txs from trustedstate. Error: %v, batch.BatchL2Data: %s", err, common.Bytes2Hex(batches[0].BatchL2Data))
		return nil, nil, nil, nil, err
	}
	log.Debug("len(storedTxs): ", len(storedTxs))
	log.Debug("len(syncedTxs): ", len(syncedTxs))
	return storedTxs, syncedTxs, storedEfficiencyPercentages, syncedEfficiencyPercentages, nil
}

func checkIfSynced(batches []*state.Batch, trustedBatch *types.Batch) bool {
	matchNumber := batches[0].BatchNumber == uint64(trustedBatch.Number)
	matchGER := batches[0].GlobalExitRoot.String() == trustedBatch.GlobalExitRoot.String()
	matchLER := batches[0].LocalExitRoot.String() == trustedBatch.LocalExitRoot.String()
	matchSR := batches[0].StateRoot.String() == trustedBatch.StateRoot.String()
	matchCoinbase := batches[0].Coinbase.String() == trustedBatch.Coinbase.String()
	matchTimestamp := uint64(batches[0].Timestamp.Unix()) == uint64(trustedBatch.Timestamp)
	matchL2Data := hex.EncodeToString(batches[0].BatchL2Data) == hex.EncodeToString(trustedBatch.BatchL2Data)

	if matchNumber && matchGER && matchLER && matchSR &&
		matchCoinbase && matchTimestamp && matchL2Data {
		return true
	}
	log.Infof("matchNumber %v %d %d", matchNumber, batches[0].BatchNumber, uint64(trustedBatch.Number))
	log.Infof("matchGER %v %s %s", matchGER, batches[0].GlobalExitRoot.String(), trustedBatch.GlobalExitRoot.String())
	log.Infof("matchLER %v %s %s", matchLER, batches[0].LocalExitRoot.String(), trustedBatch.LocalExitRoot.String())
	log.Infof("matchSR %v %s %s", matchSR, batches[0].StateRoot.String(), trustedBatch.StateRoot.String())
	log.Infof("matchCoinbase %v %s %s", matchCoinbase, batches[0].Coinbase.String(), trustedBatch.Coinbase.String())
	log.Infof("matchTimestamp %v %d %d", matchTimestamp, uint64(batches[0].Timestamp.Unix()), uint64(trustedBatch.Timestamp))
	log.Infof("matchL2Data %v", matchL2Data)
	return false
}

func (s *ClientSynchronizer) getCurrentBatches(batches []*state.Batch, trustedBatch *types.Batch, dbTx pgx.Tx) ([]*state.Batch, error) {
	if len(batches) == 0 || batches[0] == nil || (batches[0] != nil && uint64(trustedBatch.Number) != batches[0].BatchNumber) {
		log.Debug("Updating batch[0] value!")
		batch, err := s.state.GetBatchByNumber(s.ctx, uint64(trustedBatch.Number), dbTx)
		if err != nil && err != state.ErrNotFound {
			log.Warnf("failed to get batch %v from local trusted state. Error: %v", trustedBatch.Number, err)
			return nil, err
		}
		var prevBatch *state.Batch
		if len(batches) == 0 || batches[0] == nil || (batches[0] != nil && uint64(trustedBatch.Number-1) != batches[0].BatchNumber) {
			log.Debug("Updating batch[1] value!")
			prevBatch, err = s.state.GetBatchByNumber(s.ctx, uint64(trustedBatch.Number-1), dbTx)
			if err != nil && err != state.ErrNotFound {
				log.Warnf("failed to get prevBatch %v from local trusted state. Error: %v", trustedBatch.Number-1, err)
				return nil, err
			}
		} else {
			prevBatch = batches[0]
		}
		log.Debug("batch: ", batch)
		log.Debug("prevBatch: ", prevBatch)
		batches = []*state.Batch{batch, prevBatch}
	}
	return batches, nil
}

func (s *ClientSynchronizer) pendingFlushID(flushID uint64, proverID string) {
	log.Infof("pending flushID: %d", flushID)
	if flushID == 0 {
		log.Fatal("flushID is 0. Please check that prover/executor config parameter dbReadOnly is false")
	}
	s.latestFlushID = flushID
	s.latestFlushIDIsFulfilled = false
	s.updateAndCheckProverID(proverID)
}

func (s *ClientSynchronizer) updateAndCheckProverID(proverID string) {
	if s.proverID == "" {
		log.Infof("Current proverID is %s", proverID)
		s.proverID = proverID
		return
	}
	if s.proverID != proverID {
		event := &event.Event{
			ReceivedAt:  time.Now(),
			Source:      event.Source_Node,
			Component:   event.Component_Synchronizer,
			Level:       event.Level_Critical,
			EventID:     event.EventID_SynchronizerRestart,
			Description: fmt.Sprintf("proverID changed from %s to %s, restarting Synchronizer ", s.proverID, proverID),
		}

		err := s.eventLog.LogEvent(context.Background(), event)
		if err != nil {
			log.Errorf("error storing event payload: %v", err)
		}

		log.Fatal("restarting synchronizer because executor has been restarted (old=%s, new=%s)", s.proverID, proverID)
	}
}

func (s *ClientSynchronizer) checkFlushID(dbTx pgx.Tx) error {
	if s.latestFlushIDIsFulfilled {
		log.Debugf("no pending flushID, nothing to do. Last pending fulfilled flushID: %d, last executor flushId received: %d", s.latestFlushID, s.latestFlushID)
		return nil
	}
	storedFlushID, proverID, err := s.state.GetStoredFlushID(s.ctx)
	if err != nil {
		log.Error("error getting stored flushID. Error: ", err)
		return err
	}
	if s.previousExecutorFlushID != storedFlushID || s.proverID != proverID {
		log.Infof("executor vs local: flushid=%d/%d, proverID=%s/%s", storedFlushID,
			s.latestFlushID, proverID, s.proverID)
	} else {
		log.Debugf("executor vs local: flushid=%d/%d, proverID=%s/%s", storedFlushID,
			s.latestFlushID, proverID, s.proverID)
	}
	s.updateAndCheckProverID(proverID)
	log.Debugf("storedFlushID (executor reported): %d, latestFlushID (pending): %d", storedFlushID, s.latestFlushID)
	if storedFlushID < s.latestFlushID {
		log.Infof("Synchronized BLOCKED!: Wating for the flushID to be stored. FlushID to be stored: %d. Latest flushID stored: %d", s.latestFlushID, storedFlushID)
		iteration := 0
		start := time.Now()
		for storedFlushID < s.latestFlushID {
			log.Debugf("Waiting for the flushID to be stored. FlushID to be stored: %d. Latest flushID stored: %d iteration:%d elpased:%s",
				s.latestFlushID, storedFlushID, iteration, time.Since(start))
			time.Sleep(100 * time.Millisecond) //nolint:gomnd
			storedFlushID, _, err = s.state.GetStoredFlushID(s.ctx)
			if err != nil {
				log.Error("error getting stored flushID. Error: ", err)
				return err
			}
			iteration++
		}
		log.Infof("Synchronizer resumed, flushID stored: %d", s.latestFlushID)
	}
	log.Infof("Pending Flushid fullfiled: %d, executor have write %d", s.latestFlushID, storedFlushID)
	s.latestFlushIDIsFulfilled = true
	s.previousExecutorFlushID = storedFlushID
	return nil
}

// halt halts the Synchronizer
func (s *ClientSynchronizer) halt(ctx context.Context, err error) {
	event := &event.Event{
		ReceivedAt:  time.Now(),
		Source:      event.Source_Node,
		Component:   event.Component_Synchronizer,
		Level:       event.Level_Critical,
		EventID:     event.EventID_SynchronizerHalt,
		Description: fmt.Sprintf("Synchronizer halted due to error: %s", err),
	}

	eventErr := s.eventLog.LogEvent(ctx, event)
	if eventErr != nil {
		log.Errorf("error storing Synchronizer halt event: %v", eventErr)
	}

	for {
		log.Errorf("fatal error: %s", err)
		log.Error("halting the Synchronizer")
		time.Sleep(5 * time.Second) //nolint:gomnd
	}
}
