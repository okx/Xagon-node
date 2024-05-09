package synchronizer

import (
	"context"
	"errors"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/actions"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/common/syncinterfaces"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/l1event_orders"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

type stateBlockRangeProcessor interface {
	BeginStateTransaction(ctx context.Context) (pgx.Tx, error)
	AddBlock(ctx context.Context, block *state.Block, dbTx pgx.Tx) error
	GetForkIDByBatchNumber(batchNumber uint64) uint64
	GetForkIDByBlockNumber(blockNumber uint64) uint64
}

type ethermanI interface {
	GetFinalizedBlockNumber(ctx context.Context) (uint64, error)
}

// BlockRangeProcess is the struct that process the block range that implements syncinterfaces.BlockRangeProcessor
type BlockRangeProcess struct {
	state             stateBlockRangeProcessor
	etherMan          ethermanI
	l1EventProcessors syncinterfaces.L1EventProcessorManager
	flushIdManager    syncinterfaces.SynchronizerFlushIDManager
}

// NewBlockRangeProcessLegacy creates a new BlockRangeProcess
func NewBlockRangeProcessLegacy(
	state stateBlockRangeProcessor,
	etherMan ethermanI,
	l1EventProcessors syncinterfaces.L1EventProcessorManager,
	flushIdManager syncinterfaces.SynchronizerFlushIDManager,
) *BlockRangeProcess {
	return &BlockRangeProcess{
		state:             state,
		etherMan:          etherMan,
		l1EventProcessors: l1EventProcessors,
		flushIdManager:    flushIdManager,
	}
}

// ProcessBlockRangeSingleDbTx process the L1 events and stores the information in the db reusing same DbTx
func (s *BlockRangeProcess) ProcessBlockRangeSingleDbTx(ctx context.Context, blocks []etherman.Block, order map[common.Hash][]etherman.Order, storeBlocks syncinterfaces.ProcessBlockRangeL1BlocksMode, dbTx pgx.Tx) error {
	return s.internalProcessBlockRange(ctx, blocks, order, storeBlocks, &dbTx)
}

// ProcessBlockRange process the L1 events and stores the information in the db
func (s *BlockRangeProcess) ProcessBlockRange(ctx context.Context, blocks []etherman.Block, order map[common.Hash][]etherman.Order) error {
	return s.internalProcessBlockRange(ctx, blocks, order, syncinterfaces.StoreL1Blocks, nil)
}

// ProcessBlockRange process the L1 events and stores the information in the db
func (s *BlockRangeProcess) internalProcessBlockRange(ctx context.Context, blocks []etherman.Block, order map[common.Hash][]etherman.Order, storeBlocks syncinterfaces.ProcessBlockRangeL1BlocksMode, dbTxExt *pgx.Tx) error {
	// Check the latest finalized block in L1
	finalizedBlockNumber, err := s.etherMan.GetFinalizedBlockNumber(ctx)
	if err != nil {
		log.Errorf("error getting finalized block number in L1. Error: %v", err)
		return err
	}
	// New info has to be included into the db using the state
	for i := range blocks {
		// Begin db transaction
		var dbTx pgx.Tx
		var err error
		if dbTxExt == nil {
			log.Debugf("Starting dbTx for BlockNumber:%d", blocks[i].BlockNumber)
			dbTx, err = s.state.BeginStateTransaction(ctx)
			if err != nil {
				return err
			}
		} else {
			dbTx = *dbTxExt
		}
		// Process event received from l1
		err = s.processBlock(ctx, blocks, i, dbTx, order, storeBlocks, finalizedBlockNumber)
		if err != nil {
			if dbTxExt == nil {
				// Rollback db transaction
				rollbackErr := dbTx.Rollback(ctx)
				if rollbackErr != nil {
					if !errors.Is(rollbackErr, pgx.ErrTxClosed) {
						log.Errorf("error rolling back state. RollbackErr: %s, Error : %v", rollbackErr.Error(), err)
						return rollbackErr
					} else {
						log.Warnf("error rolling back state because is already closed. RollbackErr: %s, Error : %v", rollbackErr.Error(), err)
						return err
					}
				}
				return err
			}
			return err
		}
		if dbTxExt == nil {
			// Commit db transaction
			err = dbTx.Commit(ctx)
			if err != nil {
				log.Errorf("error committing state. BlockNumber: %d, Error: %v", blocks[i].BlockNumber, err)
			}
		}
	}
	return nil
}

func (s *BlockRangeProcess) processBlock(ctx context.Context, blocks []etherman.Block, i int, dbTx pgx.Tx, order map[common.Hash][]etherman.Order, storeBlock syncinterfaces.ProcessBlockRangeL1BlocksMode, finalizedBlockNumber uint64) error {
	var err error
	if storeBlock == syncinterfaces.StoreL1Blocks {
		b := state.Block{
			BlockNumber: blocks[i].BlockNumber,
			BlockHash:   blocks[i].BlockHash,
			ParentHash:  blocks[i].ParentHash,
			ReceivedAt:  blocks[i].ReceivedAt,
		}
		if blocks[i].BlockNumber <= finalizedBlockNumber {
			b.Checked = true
		}
		err = s.state.AddBlock(ctx, &b, dbTx)
		if err != nil {
			log.Errorf("error adding block to db. BlockNumber: %d, error: %v", blocks[i].BlockNumber, err)
			return err
		}
	} else {
		log.Debugf("Skip storing block BlockNumber:%d", blocks[i].BlockNumber)
	}
	for _, element := range order[blocks[i].BlockHash] {
		err := s.processElement(ctx, element, blocks, i, dbTx)
		if err != nil {
			return err
		}
	}
	log.Debug("Checking FlushID to commit L1 data to db")
	err = s.flushIdManager.CheckFlushID(dbTx)
	if err != nil {
		log.Errorf("error checking flushID. BlockNumber: %d, Error: %v", blocks[i].BlockNumber, err)
		return err
	}
	return nil
}

func (s *BlockRangeProcess) processElement(ctx context.Context, element etherman.Order, blocks []etherman.Block, i int, dbTx pgx.Tx) error {
	batchSequence := l1event_orders.GetSequenceFromL1EventOrder(element.Name, &blocks[i], element.Pos)
	var forkId uint64
	if batchSequence != nil {
		forkId = s.state.GetForkIDByBatchNumber(batchSequence.FromBatchNumber)
		log.Debug("EventOrder: ", element.Name, ". Batch Sequence: ", batchSequence, "forkId: ", forkId)
	} else {
		forkId = s.state.GetForkIDByBlockNumber(blocks[i].BlockNumber)
		log.Debug("EventOrder: ", element.Name, ". BlockNumber: ", blocks[i].BlockNumber, "forkId: ", forkId)
	}
	forkIdTyped := actions.ForkIdType(forkId)

	err := s.l1EventProcessors.Process(ctx, forkIdTyped, element, &blocks[i], dbTx)
	if err != nil {
		log.Error("error l1EventProcessors.Process: ", err)
		return err
	}
	return nil
}
