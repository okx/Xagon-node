package pgstatestorage

import (
	"github.com/jackc/pgx/v4"
	"context"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"errors"
)

// GetLastL2BlockTimeByBatchNumber gets the last l2 block time in a batch by batch number
func (p *PostgresStorage) GetLastL2BlockTimeByBatchNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (uint64, error) {
	const query = "SELECT header FROM state.l2block b WHERE batch_num = $1 ORDER BY b.block_num DESC LIMIT 1"

	header := &state.L2Header{}
	q := p.getExecQuerier(dbTx)
	err := q.QueryRow(ctx, query, batchNumber).Scan(&header)

	if errors.Is(err, pgx.ErrNoRows) {
		return 0, state.ErrNotFound
	} else if err != nil {
		return 0, err
	}

	return header.Time, nil
}
