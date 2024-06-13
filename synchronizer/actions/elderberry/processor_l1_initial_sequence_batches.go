package elderberry

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/actions"
	"github.com/jackc/pgx/v4"
)

// ProcessorL1InitialSequenceBatchesElderberry is the processor for SequenceBatches for Elderberry
// intialSequence is process in ETROG by the same class, this is just a wrapper to pass directly to ETROG
type ProcessorL1InitialSequenceBatchesElderberry struct {
	actions.ProcessorBase[ProcessorL1InitialSequenceBatchesElderberry]
	previousProcessor actions.L1EventProcessor
}

// NewProcessorL1InitialSequenceBatchesElderberry returns instance of a processor for InitialSequenceBatchesOrder
func NewProcessorL1InitialSequenceBatchesElderberry(previousProcessor actions.L1EventProcessor) *ProcessorL1InitialSequenceBatchesElderberry {
	return &ProcessorL1InitialSequenceBatchesElderberry{
		ProcessorBase: *actions.NewProcessorBase[ProcessorL1InitialSequenceBatchesElderberry](
			[]etherman.EventOrder{etherman.InitialSequenceBatchesOrder},
			actions.ForksIdOnlyElderberry),
		previousProcessor: previousProcessor,
	}
}

// Process process event
func (g *ProcessorL1InitialSequenceBatchesElderberry) Process(ctx context.Context, order etherman.Order, l1Block *etherman.Block, dbTx pgx.Tx) error {
	log.Infof("Elderberry: Executing initialSequenceBatch(%s). Processing with previous processor", g.previousProcessor.Name())
	return g.previousProcessor.Process(ctx, order, l1Block, dbTx)
}
