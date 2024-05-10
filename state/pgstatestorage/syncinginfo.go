package pgstatestorage

import (
	"context"
	"errors"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/jackc/pgx/v4"
)

func (p *PostgresStorage) GetSyncInfoData(ctx context.Context, dbTx pgx.Tx) (state.SyncInfoDataOnStorage, error) {
	var info state.SyncInfoDataOnStorage
	const getSyncTableSQL = `
	select last_batch_num_seen, last_batch_num_consolidated, init_sync_batch from state.sync_info;
	`
	q := p.getExecQuerier(dbTx)
	row := q.QueryRow(ctx, getSyncTableSQL)
	err := row.Scan(&info.LastBatchNumberSeen, &info.LastBatchNumberConsolidated, &info.InitialSyncingBatch)
	if errors.Is(err, pgx.ErrNoRows) {
		return state.SyncInfoDataOnStorage{}, state.ErrNotFound
	} else if err != nil {
		return state.SyncInfoDataOnStorage{}, err
	}
	return info, nil
}
