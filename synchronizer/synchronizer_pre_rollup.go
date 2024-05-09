package synchronizer

import (
	"context"
	"errors"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/common/syncinterfaces"
	"github.com/jackc/pgx/v4"
)

const (
	pregenesisSyncLogPrefix = "sync pregenesis:"
)

// SyncPreRollup is the struct for synchronizing pre genesis rollup events.
// Implements: syncinterfaces.SyncPreRollupSyncer
type SyncPreRollup struct {
	etherman            syncinterfaces.EthermanPreRollup
	state               syncinterfaces.StateLastBlockGetter
	blockRangeProcessor syncinterfaces.BlockRangeProcessor
	SyncChunkSize       uint64
	GenesisBlockNumber  uint64
}

// NewSyncPreRollup creates a new SyncPreRollup
func NewSyncPreRollup(
	etherman syncinterfaces.EthermanPreRollup,
	state syncinterfaces.StateLastBlockGetter,
	blockRangeProcessor syncinterfaces.BlockRangeProcessor,
	syncChunkSize uint64,
	genesisBlockNumber uint64,
) *SyncPreRollup {
	return &SyncPreRollup{
		etherman:            etherman,
		state:               state,
		blockRangeProcessor: blockRangeProcessor,
		SyncChunkSize:       syncChunkSize,
		GenesisBlockNumber:  genesisBlockNumber,
	}
}

// SynchronizePreGenesisRollupEvents sync pre-rollup events
func (s *SyncPreRollup) SynchronizePreGenesisRollupEvents(ctx context.Context) error {
	// Sync events from RollupManager that happen before rollup creation
	log.Info(pregenesisSyncLogPrefix + "synchronizing events from RollupManager that happen before rollup creation")
	needToUpdate, fromBlock, err := s.getStartingL1Block(ctx, nil)
	if err != nil {
		log.Errorf(pregenesisSyncLogPrefix+"error getting starting L1 block. Error: %v", err)
		return err
	}
	if needToUpdate {
		return s.ProcessL1InfoRootEvents(ctx, fromBlock, s.GenesisBlockNumber-1, s.SyncChunkSize)
	} else {
		log.Infof(pregenesisSyncLogPrefix+"No need to process blocks before the genesis block %d", s.GenesisBlockNumber)
		return nil
	}
}

// getStartingL1Block find if need to update and if yes the starting point:
// bool -> need to process blocks
// uint64 -> first block to synchronize
// error -> error
// 1. First try to get last block on DB, if there are could be fully synced or pending blocks
// 2. If DB is empty the LxLy upgrade block as starting point
func (s *SyncPreRollup) getStartingL1Block(ctx context.Context, dbTx pgx.Tx) (bool, uint64, error) {
	lastBlock, err := s.state.GetLastBlock(ctx, dbTx)
	if err != nil && errors.Is(err, state.ErrStateNotSynchronized) {
		// No block on DB
		upgradeLxLyBlockNumber, err := s.etherman.GetL1BlockUpgradeLxLy(ctx, s.GenesisBlockNumber)
		if err != nil && errors.Is(err, etherman.ErrNotFound) {
			log.Infof(pregenesisSyncLogPrefix+"LxLy upgrade not detected before genesis block %d, it'll be sync as usual. Nothing to do yet", s.GenesisBlockNumber)
			return false, 0, nil
		} else if err != nil {
			log.Errorf(pregenesisSyncLogPrefix+"error getting LxLy upgrade block. Error: %v", err)
			return false, 0, err
		}
		log.Infof(pregenesisSyncLogPrefix+"No block on DB, starting from LxLy upgrade block %d", upgradeLxLyBlockNumber)
		return true, upgradeLxLyBlockNumber, nil
	} else if err != nil {
		log.Errorf("Error getting last Block on DB err:%v", err)
		return false, 0, err
	}
	if lastBlock.BlockNumber >= s.GenesisBlockNumber-1 {
		log.Warnf(pregenesisSyncLogPrefix+"Last block processed is %d, which is greater or equal than the previous genesis block %d", lastBlock, s.GenesisBlockNumber)
		return false, 0, nil
	}
	log.Infof(pregenesisSyncLogPrefix+"Continue processing pre-genesis blocks, last block processed on DB is %d", lastBlock.BlockNumber)
	return true, lastBlock.BlockNumber, nil
}

// ProcessL1InfoRootEvents processes the L1InfoRoot events for a range for L1 blocks
func (s *SyncPreRollup) ProcessL1InfoRootEvents(ctx context.Context, fromBlock uint64, toBlock uint64, syncChunkSize uint64) error {
	startTime := time.Now()
	log.Info(pregenesisSyncLogPrefix + "synchronizing L1InfoRoot events")
	log.Infof(pregenesisSyncLogPrefix+"Starting syncing pre genesis LxLy events from block %d to block %d (total %d blocks)",
		fromBlock, toBlock, toBlock-fromBlock+1)
	for i := fromBlock; true; i += syncChunkSize {
		toBlockReq := min(i+syncChunkSize-1, toBlock)
		percent := float32(toBlockReq-fromBlock+1) * 100.0 / float32(toBlock-fromBlock+1) // nolint:gomnd
		log.Infof(pregenesisSyncLogPrefix+"sync L1InfoTree events from %d to %d percent:%3.1f %% pending_blocks:%d", i, toBlockReq, percent, toBlock-toBlockReq)
		blocks, order, err := s.etherman.GetRollupInfoByBlockRangePreviousRollupGenesis(ctx, i, &toBlockReq)
		if err != nil {
			log.Error(pregenesisSyncLogPrefix+"error getting rollupInfoByBlockRange before rollup genesis: ", err)
			return err
		}
		err = s.blockRangeProcessor.ProcessBlockRange(ctx, blocks, order)
		if err != nil {
			log.Error(pregenesisSyncLogPrefix+"error processing blocks before the genesis: ", err)
			return err
		}
		if toBlockReq == toBlock {
			break
		}
	}
	elapsedTime := time.Since(startTime)
	log.Infof(pregenesisSyncLogPrefix+"sync L1InfoTree finish: from %d to %d total_block %d done in %s", fromBlock, toBlock, toBlock-fromBlock+1, &elapsedTime)
	return nil
}
