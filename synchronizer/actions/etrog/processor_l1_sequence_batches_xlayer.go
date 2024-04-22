package etrog

import (
	"context"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/jackc/pgx/v4"
)

func (p *ProcessorL1SequenceBatchesEtrog) updatePermissionLessBatchTimestamp(ctx context.Context, batchNumber uint64, batchTime time.Time, dbTx pgx.Tx) error {
	if p.sync.IsTrustedSequencer() {
		return nil
	}

	err := p.state.UpdateBatchTimestamp(ctx, batchNumber, batchTime, dbTx)
	if err != nil {
		log.Errorf("error update batch timestamp for batch: %d. Error; %v", batchNumber, err)
		rollbackErr := dbTx.Rollback(ctx)
		if rollbackErr != nil {
			log.Errorf("error rolling back state. BatchNumber: %d, rollbackErr: %v", batchNumber, rollbackErr)
			return rollbackErr
		}
		return err
	}
	log.Infof("Permission less rpc updates batch timestamp for batch: %v with new timestamp:%v", batchNumber, batchTime)
	return nil
}
