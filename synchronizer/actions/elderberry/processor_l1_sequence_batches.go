package elderberry

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/actions"
	"github.com/jackc/pgx/v4"
)

var (
	// ErrInvalidInitialBatchNumber is returned when the initial batch number is not the expected one
	ErrInvalidInitialBatchNumber = errors.New("invalid initial batch number")
)

// PreviousProcessor is the interface that the previous processor (Etrog)
type PreviousProcessor interface {
	Process(ctx context.Context, order etherman.Order, l1Block *etherman.Block, dbTx pgx.Tx) error
	ProcessSequenceBatches(ctx context.Context, sequencedBatches []etherman.SequencedBatch, blockNumber uint64, l1BlockTimestamp time.Time, dbTx pgx.Tx) error
}

// StateL1SequenceBatchesElderberry state interface
type StateL1SequenceBatchesElderberry interface {
	GetLastVirtualBatchNum(ctx context.Context, dbTx pgx.Tx) (uint64, error)
	GetLastL2BlockByBatchNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*state.L2Block, error)
}

// ProcessorL1SequenceBatchesElderberry is the processor for SequenceBatches for Elderberry
type ProcessorL1SequenceBatchesElderberry struct {
	actions.ProcessorBase[ProcessorL1SequenceBatchesElderberry]
	previousProcessor PreviousProcessor
	state             StateL1SequenceBatchesElderberry
}

// NewProcessorL1SequenceBatchesElderberry returns instance of a processor for SequenceBatchesOrder
func NewProcessorL1SequenceBatchesElderberry(previousProcessor PreviousProcessor, state StateL1SequenceBatchesElderberry) *ProcessorL1SequenceBatchesElderberry {
	return &ProcessorL1SequenceBatchesElderberry{
		ProcessorBase: actions.ProcessorBase[ProcessorL1SequenceBatchesElderberry]{
			SupportedEvent:    []etherman.EventOrder{etherman.SequenceBatchesOrder},
			SupportedForkdIds: &actions.ForksIdOnlyElderberry},
		previousProcessor: previousProcessor,
		state:             state,
	}
}

// Process process event
func (g *ProcessorL1SequenceBatchesElderberry) Process(ctx context.Context, order etherman.Order, l1Block *etherman.Block, dbTx pgx.Tx) error {
	if l1Block == nil || len(l1Block.SequencedBatches) <= order.Pos {
		return actions.ErrInvalidParams
	}
	if len(l1Block.SequencedBatches[order.Pos]) == 0 {
		log.Warnf("No sequenced batches for position")
		return nil
	}

	sbatch := l1Block.SequencedBatches[order.Pos][0]

	if sbatch.SequencedBatchElderberryData == nil {
		log.Errorf("No elderberry sequenced batch data for batch %d", sbatch.BatchNumber)
		return fmt.Errorf("no elderberry sequenced batch data for batch %d", sbatch.BatchNumber)
	}
	// We need to check that the sequence match
	err := g.sanityCheckExpectedSequence(sbatch.SequencedBatchElderberryData.InitSequencedBatchNumber, dbTx)
	if err != nil {
		return err
	}
	// We known that the MaxSequenceTimestamp is the same for all the batches so we can use the first one
	err = g.previousProcessor.ProcessSequenceBatches(ctx, l1Block.SequencedBatches[order.Pos], l1Block.BlockNumber, time.Unix(int64(sbatch.SequencedBatchElderberryData.MaxSequenceTimestamp), 0), dbTx)
	// The last L2block timestamp must match MaxSequenceTimestamp
	if err != nil {
		return err
	}
	// It checks the timestamp of the last L2 block, but it's just log an error instead of refusing the event
	_ = g.sanityCheckTstampLastL2Block(sbatch.SequencedBatchElderberryData.MaxSequenceTimestamp, dbTx)
	return nil
}

func (g *ProcessorL1SequenceBatchesElderberry) sanityCheckExpectedSequence(initialBatchNumber uint64, dbTx pgx.Tx) error {
	// We need to check that the sequence match
	lastVirtualBatchNum, err := g.state.GetLastVirtualBatchNum(context.Background(), dbTx)
	if err != nil {
		log.Errorf("Error getting last virtual batch number: %s", err)
		return err
	}
	if lastVirtualBatchNum != initialBatchNumber {
		log.Errorf("The last virtual batch number is not the expected one. Expected: %d (last on DB), got: %d (L1 event)", lastVirtualBatchNum+1, initialBatchNumber)
		return fmt.Errorf("the last virtual batch number is not the expected one. Expected: %d (last on DB), got: %d (L1 event) err:%w", lastVirtualBatchNum+1, initialBatchNumber, ErrInvalidInitialBatchNumber)
	}
	return nil
}

func (g *ProcessorL1SequenceBatchesElderberry) sanityCheckTstampLastL2Block(timeLimit uint64, dbTx pgx.Tx) error {
	lastVirtualBatchNum, err := g.state.GetLastVirtualBatchNum(context.Background(), dbTx)
	if err != nil {
		log.Errorf("Error getting last virtual batch number: %s", err)
		return err
	}
	lastL2Block, err := g.state.GetLastL2BlockByBatchNumber(context.Background(), lastVirtualBatchNum, dbTx)
	if err != nil {
		log.Errorf("Error getting last virtual batch number: %s", err)
		return err
	}
	if lastL2Block == nil {
		//TODO: find the previous batch until we find a L2 block to check the timestamp
		return nil
	}
	if uint64(lastL2Block.ReceivedAt.Unix()) > timeLimit {
		log.Errorf("The last L2 block timestamp can't be greater than timeLimit. Expected: %d (L1 event), got: %d (last L2Block)", timeLimit, lastL2Block.ReceivedAt.Unix())
		return fmt.Errorf("wrong timestamp of  last L2 block timestamp with L1 event timestamp")
	}
	return nil
}
