package incaberry

import (
	"context"
	"fmt"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/actions"
	"github.com/jackc/pgx/v4"
)

type stateL1VerifyBatchInterface interface {
	GetLastVerifiedBatch(ctx context.Context, dbTx pgx.Tx) (*state.VerifiedBatch, error)
	GetBatchByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*state.Batch, error)
	AddVerifiedBatch(ctx context.Context, verifiedBatch *state.VerifiedBatch, dbTx pgx.Tx) error
}

// ProcessorL1VerifyBatch implements L1EventProcessor
type ProcessorL1VerifyBatch struct {
	actions.ProcessorBase[ProcessorL1VerifyBatch]
	state stateL1VerifyBatchInterface
}

// NewProcessorL1VerifyBatch returns instance of a processor for VerifyBatchOrder
func NewProcessorL1VerifyBatch(state stateL1VerifyBatchInterface) *ProcessorL1VerifyBatch {
	return &ProcessorL1VerifyBatch{
		ProcessorBase: *actions.NewProcessorBase[ProcessorL1VerifyBatch](
			[]etherman.EventOrder{etherman.VerifyBatchOrder, etherman.TrustedVerifyBatchOrder},
			actions.ForksIdAll),
		state: state,
	}
}

// Process process event
func (p *ProcessorL1VerifyBatch) Process(ctx context.Context, order etherman.Order, l1Block *etherman.Block, dbTx pgx.Tx) error {
	var isTrusted bool
	if order.Name == etherman.VerifyBatchOrder {
		isTrusted = true
	}
	return p.processVerifyBatches(ctx, l1Block.VerifiedBatches[order.Pos], isTrusted, dbTx)
}

func (p *ProcessorL1VerifyBatch) processVerifyBatches(ctx context.Context, lastVerifiedBatch etherman.VerifiedBatch, isTrusted bool, dbTx pgx.Tx) error {
	lastVBatch, err := p.state.GetLastVerifiedBatch(ctx, dbTx)
	if err != nil {
		log.Errorf("error getting lastVerifiedBatch stored in db in processVerifyBatches. Processing synced blockNumber: %d", lastVerifiedBatch.BlockNumber)
		rollbackErr := dbTx.Rollback(ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state. Processing synced blockNumber: %d, rollbackErr: %s, error : %v", lastVerifiedBatch.BlockNumber, rollbackErr.Error(), err)
			return rollbackErr
		}
		log.Errorf("error getting lastVerifiedBatch stored in db in processVerifyBatches. Processing synced blockNumber: %d, error: %v", lastVerifiedBatch.BlockNumber, err)
		return err
	}
	nbatches := lastVerifiedBatch.BatchNumber - lastVBatch.BatchNumber
	batch, err := p.state.GetBatchByNumber(ctx, lastVerifiedBatch.BatchNumber, dbTx)
	if err != nil {
		log.Errorf("error getting GetBatchByNumber stored in db in processVerifyBatches. Processing batchNumber: %d", lastVerifiedBatch.BatchNumber)
		rollbackErr := dbTx.Rollback(ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state. Processing batchNumber: %d, rollbackErr: %s, error : %v", lastVerifiedBatch.BatchNumber, rollbackErr.Error(), err)
			return rollbackErr
		}
		log.Errorf("error getting GetBatchByNumber stored in db in processVerifyBatches. Processing batchNumber: %d, error: %v", lastVerifiedBatch.BatchNumber, err)
		return err
	}

	// Checks that calculated state root matches with the verified state root in the smc
	if batch.StateRoot != lastVerifiedBatch.StateRoot {
		log.Warn("nbatches: ", nbatches)
		log.Warnf("Batch from db: %+v", batch)
		log.Warnf("Verified Batch: %+v", lastVerifiedBatch)
		log.Errorf("error: stateRoot calculated and state root verified don't match in processVerifyBatches. Processing batchNumber: %d", lastVerifiedBatch.BatchNumber)
		rollbackErr := dbTx.Rollback(ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state. Processing batchNumber: %d, rollbackErr: %v", lastVerifiedBatch.BatchNumber, rollbackErr)
			return rollbackErr
		}
		log.Errorf("error: stateRoot calculated and state root verified don't match in processVerifyBatches. Processing batchNumber: %d", lastVerifiedBatch.BatchNumber)
		return fmt.Errorf("error: stateRoot calculated and state root verified don't match in processVerifyBatches. Processing batchNumber: %d", lastVerifiedBatch.BatchNumber)
	}
	var i uint64
	for i = 1; i <= nbatches; i++ {
		verifiedB := state.VerifiedBatch{
			BlockNumber: lastVerifiedBatch.BlockNumber,
			BatchNumber: lastVBatch.BatchNumber + i,
			Aggregator:  lastVerifiedBatch.Aggregator,
			StateRoot:   lastVerifiedBatch.StateRoot,
			TxHash:      lastVerifiedBatch.TxHash,
			IsTrusted:   isTrusted,
		}
		log.Infof("processVerifyBatches: Storing verifiedB. BlockNumber: %d, BatchNumber: %d, isTrusted: %v", verifiedB.BlockNumber, verifiedB.BatchNumber, isTrusted)
		err = p.state.AddVerifiedBatch(ctx, &verifiedB, dbTx)
		if err != nil {
			log.Errorf("error storing the verifiedB in processVerifyBatches. verifiedBatch: %+v, lastVerifiedBatch: %+v", verifiedB, lastVerifiedBatch)
			rollbackErr := dbTx.Rollback(ctx)
			if rollbackErr != nil {
				log.Errorf("error rolling back state. BlockNumber: %d, rollbackErr: %s, error : %v", lastVerifiedBatch.BlockNumber, rollbackErr.Error(), err)
				return rollbackErr
			}
			log.Errorf("error storing the verifiedB in processVerifyBatches. BlockNumber: %d, error: %v", lastVerifiedBatch.BlockNumber, err)
			return err
		}
	}
	return nil
}
