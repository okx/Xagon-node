package pgstatestorage

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"
)

func (p *PostgresStorage) UpdateBatchTimestamp(ctx context.Context, batchNumber uint64, batchTime time.Time, dbTx pgx.Tx) error {
	const updateL2DataSQL = "UPDATE state.batch SET timestamp = $1 WHERE batch_num = $2"

	e := p.getExecQuerier(dbTx)
	_, err := e.Exec(ctx, updateL2DataSQL, batchTime.UTC(), batchNumber)
	return err
}
