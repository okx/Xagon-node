package incaberry

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/actions"
	"github.com/jackc/pgx/v4"
)

type stateProcessL1ForcedBatchesInterface interface {
	AddForcedBatch(ctx context.Context, forcedBatch *state.ForcedBatch, dbTx pgx.Tx) error
}

// ProcessL1ForcedBatches implements L1EventProcessor
type ProcessL1ForcedBatches struct {
	actions.ProcessorBase[ProcessL1ForcedBatches]
	state stateProcessL1ForcedBatchesInterface
}

// NewProcessL1ForcedBatches returns instance of a processor for ForcedBatchesOrder
func NewProcessL1ForcedBatches(state stateProcessL1ForcedBatchesInterface) *ProcessL1ForcedBatches {
	return &ProcessL1ForcedBatches{
		ProcessorBase: *actions.NewProcessorBase[ProcessL1ForcedBatches](
			[]etherman.EventOrder{etherman.ForcedBatchesOrder},
			actions.ForksIdAll),
		state: state}
}

// Process process event
func (p *ProcessL1ForcedBatches) Process(ctx context.Context, order etherman.Order, l1Block *etherman.Block, dbTx pgx.Tx) error {
	return p.processForcedBatch(ctx, l1Block.ForcedBatches[order.Pos], dbTx)
}

func (p *ProcessL1ForcedBatches) processForcedBatch(ctx context.Context, forcedBatch etherman.ForcedBatch, dbTx pgx.Tx) error {
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
	err := p.state.AddForcedBatch(ctx, &forcedB, dbTx)
	if err != nil {
		log.Errorf("error storing the forcedBatch in processForcedBatch. BlockNumber: %d", forcedBatch.BlockNumber)
		rollbackErr := dbTx.Rollback(ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state. BlockNumber: %d, rollbackErr: %s, error : %v", forcedBatch.BlockNumber, rollbackErr.Error(), err)
			return rollbackErr
		}
		log.Errorf("error storing the forcedBatch in processForcedBatch. BlockNumber: %d, error: %v", forcedBatch.BlockNumber, err)
		return err
	}
	return nil
}
