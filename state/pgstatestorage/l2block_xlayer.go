package pgstatestorage

import (
	"github.com/jackc/pgx/v4"
	"context"
	"time"
)

// GetLastL2BlockCreateTimeBatchNumber gets the last l2 block create time in a batch by batch number
func (p *PostgresStorage) GetLastL2BlockCreateTimeBatchNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*time.Time, error) {
	const query = "SELECT created_at FROM state.l2block b WHERE batch_num = $1 ORDER BY b.block_num DESC LIMIT 1"

	var createdAt time.Time
	q := p.getExecQuerier(dbTx)
	err := q.QueryRow(ctx, query, batchNumber).Scan(&createdAt)
	if err != nil {
		return nil, err
	}

	return &createdAt, nil
}
