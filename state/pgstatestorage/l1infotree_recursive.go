package pgstatestorage

import (
	"context"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/jackc/pgx/v4"
)

const (
	l1InfoTreeRecursiveIndexFieldName = "l1_info_tree_recursive_index"
)

// AddL1InfoRootToExitRoot adds a new entry in ExitRoot and returns index of L1InfoTree and error
func (p *PostgresStorage) AddL1InfoTreeRecursiveRootToExitRoot(ctx context.Context, exitRoot *state.L1InfoTreeRecursiveExitRootStorageEntry, dbTx pgx.Tx) error {
	exitRootOld := state.L1InfoTreeExitRootStorageEntry(*exitRoot)
	return p.addL1InfoRootToExitRootVx(ctx, &exitRootOld, dbTx, l1InfoTreeRecursiveIndexFieldName)
}

func (p *PostgresStorage) GetAllL1InfoTreeRecursiveRootEntries(ctx context.Context, dbTx pgx.Tx) ([]state.L1InfoTreeRecursiveExitRootStorageEntry, error) {
	res, err := p.GetAllL1InfoRootEntriesVx(ctx, dbTx, l1InfoTreeRecursiveIndexFieldName)
	if err != nil {
		return nil, err
	}
	var entries []state.L1InfoTreeRecursiveExitRootStorageEntry
	for _, entry := range res {
		entries = append(entries, state.L1InfoTreeRecursiveExitRootStorageEntry(entry))
	}
	return entries, nil
}

func (p *PostgresStorage) GetLatestL1InfoTreeRecursiveRoot(ctx context.Context, maxBlockNumber uint64, dbTx pgx.Tx) (state.L1InfoTreeRecursiveExitRootStorageEntry, error) {
	res, err := p.GetLatestL1InfoRootVx(ctx, maxBlockNumber, dbTx, l1InfoTreeRecursiveIndexFieldName)
	if err != nil {
		return state.L1InfoTreeRecursiveExitRootStorageEntry{}, err
	}
	return state.L1InfoTreeRecursiveExitRootStorageEntry(res), nil
}

func (p *PostgresStorage) GetLatestL1InfoTreeRecursiveIndex(ctx context.Context, dbTx pgx.Tx) (uint32, error) {
	return p.GetLatestIndexVx(ctx, dbTx, l1InfoTreeRecursiveIndexFieldName)
}

func (p *PostgresStorage) GetL1InfoRecursiveRootLeafByIndex(ctx context.Context, l1InfoTreeIndex uint32, dbTx pgx.Tx) (state.L1InfoTreeExitRootStorageEntry, error) {
	return p.GetL1InfoRootLeafByIndexVx(ctx, l1InfoTreeIndex, dbTx, l1InfoTreeIndexFieldName)
}