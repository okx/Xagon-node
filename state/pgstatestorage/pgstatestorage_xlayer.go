package pgstatestorage

import (
	"context"
	"errors"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/jackc/pgx/v4"
)

// GetBatchL2DataByNumber returns the batch L2 data of the given batch number.
func (p *PostgresStorage) GetBatchL2DataByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) ([]byte, error) {
	batchData, err := p.GetBatchL2DataByNumbers(ctx, []uint64{batchNumber}, dbTx)
	if err != nil {
		return nil, err
	}
	data, ok := batchData[batchNumber]
	if !ok {
		return nil, state.ErrNotFound
	}
	return data, nil
}

// GetBatchL2DataByNumbers returns the batch L2 data of the given batch numbers. The data is a union of state.batch and state.forced_batch tables.
func (p *PostgresStorage) GetBatchL2DataByNumbers(ctx context.Context, batchNumbers []uint64, dbTx pgx.Tx) (map[uint64][]byte, error) {
	const getBatchL2DataByBatchNumber = `
	SELECT batch_num, raw_txs_data FROM state.batch WHERE batch_num = ANY($1)  
	UNION
	SELECT forced_batch_num, convert_from(decode(raw_txs_data, 'hex'), 'UTF8')::bytea FROM state.forced_batch WHERE forced_batch_num = ANY($2)
`
	q := p.getExecQuerier(dbTx)
	rows, err := q.Query(ctx, getBatchL2DataByBatchNumber, batchNumbers, batchNumbers)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrNotFound
	} else if err != nil {
		return nil, err
	}
	defer rows.Close()
	batchL2DataMap := make(map[uint64][]byte)
	for rows.Next() {
		var (
			batchNum    uint64
			batchL2Data []byte
		)
		err := rows.Scan(&batchNum, &batchL2Data)
		if err != nil {
			return nil, err
		}
		batchL2DataMap[batchNum] = batchL2Data
	}
	if len(batchL2DataMap) == 0 {
		return nil, state.ErrNotFound
	}
	return batchL2DataMap, nil
}
