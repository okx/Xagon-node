package pgstatestorage

import (
	"context"
	"errors"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/jackc/pgx/v4"
)

// GetBatchL2DataByNumber returns the batch L2 data of the given batch number.
func (p *PostgresStorage) GetBatchL2DataByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) ([]byte, error) {
	getBatchL2DataByBatchNumber := "SELECT raw_txs_data FROM state.batch WHERE batch_num = $1"
	q := p.getExecQuerier(dbTx)
	var batchL2Data []byte
	err := q.QueryRow(ctx, getBatchL2DataByBatchNumber, batchNumber).Scan(&batchL2Data)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return batchL2Data, nil
}
