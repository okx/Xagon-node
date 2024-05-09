package pgstatestorage

import (
	"context"
	"errors"
	"fmt"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

const (
	l1InfoTreeIndexFieldName = "l1_info_tree_index"
)

// AddL1InfoRootToExitRoot adds a new entry in ExitRoot and returns index of L1InfoTree and error
func (p *PostgresStorage) AddL1InfoRootToExitRoot(ctx context.Context, exitRoot *state.L1InfoTreeExitRootStorageEntry, dbTx pgx.Tx) error {
	return p.addL1InfoRootToExitRootVx(ctx, exitRoot, dbTx, l1InfoTreeIndexFieldName)
}

func (p *PostgresStorage) addL1InfoRootToExitRootVx(ctx context.Context, exitRoot *state.L1InfoTreeExitRootStorageEntry, dbTx pgx.Tx, indexFieldName string) error {
	const addGlobalExitRootSQL = `
		INSERT INTO state.exit_root(block_num, timestamp, mainnet_exit_root, rollup_exit_root, global_exit_root, prev_block_hash, l1_info_root, %s)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
	`
	sql := fmt.Sprintf(addGlobalExitRootSQL, indexFieldName)
	e := p.getExecQuerier(dbTx)
	_, err := e.Exec(ctx, sql,
		exitRoot.BlockNumber, exitRoot.Timestamp, exitRoot.MainnetExitRoot, exitRoot.RollupExitRoot,
		exitRoot.GlobalExitRoot.GlobalExitRoot, exitRoot.PreviousBlockHash, exitRoot.L1InfoTreeRoot, exitRoot.L1InfoTreeIndex)
	return err
}

func (p *PostgresStorage) GetAllL1InfoRootEntries(ctx context.Context, dbTx pgx.Tx) ([]state.L1InfoTreeExitRootStorageEntry, error) {
	return p.GetAllL1InfoRootEntriesVx(ctx, dbTx, l1InfoTreeIndexFieldName)
}

func (p *PostgresStorage) GetAllL1InfoRootEntriesVx(ctx context.Context, dbTx pgx.Tx, indexFieldName string) ([]state.L1InfoTreeExitRootStorageEntry, error) {
	const getL1InfoRootSQL = `SELECT block_num, timestamp, mainnet_exit_root, rollup_exit_root, global_exit_root, prev_block_hash, l1_info_root, %s
		FROM state.exit_root 
		WHERE %s IS NOT NULL
		ORDER BY %s`

	sql := fmt.Sprintf(getL1InfoRootSQL, indexFieldName, indexFieldName, indexFieldName)
	e := p.getExecQuerier(dbTx)
	rows, err := e.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []state.L1InfoTreeExitRootStorageEntry
	for rows.Next() {
		var entry state.L1InfoTreeExitRootStorageEntry
		err := rows.Scan(&entry.BlockNumber, &entry.Timestamp, &entry.MainnetExitRoot, &entry.RollupExitRoot, &entry.GlobalExitRoot.GlobalExitRoot,
			&entry.PreviousBlockHash, &entry.L1InfoTreeRoot, &entry.L1InfoTreeIndex)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

// GetLatestL1InfoRoot is used to get the latest L1InfoRoot
func (p *PostgresStorage) GetLatestL1InfoRoot(ctx context.Context, maxBlockNumber uint64) (state.L1InfoTreeExitRootStorageEntry, error) {
	return p.GetLatestL1InfoRootVx(ctx, maxBlockNumber, nil, l1InfoTreeIndexFieldName)
}

// GetLatestL1InfoRoot is used to get the latest L1InfoRoot
func (p *PostgresStorage) GetLatestL1InfoRootVx(ctx context.Context, maxBlockNumber uint64, dbTx pgx.Tx, indexFieldName string) (state.L1InfoTreeExitRootStorageEntry, error) {
	const getL1InfoRootSQL = `SELECT block_num, timestamp, mainnet_exit_root, rollup_exit_root, global_exit_root, prev_block_hash, l1_info_root, %s
		FROM state.exit_root 
		WHERE %s IS NOT NULL AND block_num <= $1
		ORDER BY %s DESC LIMIT 1`

	sql := fmt.Sprintf(getL1InfoRootSQL, indexFieldName, indexFieldName, indexFieldName)

	entry := state.L1InfoTreeExitRootStorageEntry{}

	e := p.getExecQuerier(dbTx)
	err := e.QueryRow(ctx, sql, maxBlockNumber).Scan(&entry.BlockNumber, &entry.Timestamp, &entry.MainnetExitRoot, &entry.RollupExitRoot, &entry.GlobalExitRoot.GlobalExitRoot,
		&entry.PreviousBlockHash, &entry.L1InfoTreeRoot, &entry.L1InfoTreeIndex)

	if !errors.Is(err, pgx.ErrNoRows) {
		return entry, err
	}

	return entry, nil
}
func (p *PostgresStorage) GetLatestIndex(ctx context.Context, dbTx pgx.Tx) (uint32, error) {
	return p.GetLatestIndexVx(ctx, dbTx, l1InfoTreeIndexFieldName)
}
func (p *PostgresStorage) GetLatestIndexVx(ctx context.Context, dbTx pgx.Tx, indexFieldName string) (uint32, error) {
	const getLatestIndexSQL = `SELECT max(%s) as %s FROM state.exit_root 
		WHERE %s IS NOT NULL`
	sql := fmt.Sprintf(getLatestIndexSQL, indexFieldName, indexFieldName, indexFieldName)

	var l1InfoTreeIndex *uint32
	e := p.getExecQuerier(dbTx)
	err := e.QueryRow(ctx, sql).Scan(&l1InfoTreeIndex)
	if err != nil {
		return 0, err
	}
	if l1InfoTreeIndex == nil {
		return 0, state.ErrNotFound
	}
	return *l1InfoTreeIndex, nil
}

func (p *PostgresStorage) GetL1InfoRootLeafByL1InfoRoot(ctx context.Context, l1InfoRoot common.Hash, dbTx pgx.Tx) (state.L1InfoTreeExitRootStorageEntry, error) {
	return p.GetL1InfoRootLeafByL1InfoRootVx(ctx, l1InfoRoot, dbTx, l1InfoTreeIndexFieldName)
}

func (p *PostgresStorage) GetL1InfoRootLeafByL1InfoRootVx(ctx context.Context, l1InfoRoot common.Hash, dbTx pgx.Tx, indexFieldName string) (state.L1InfoTreeExitRootStorageEntry, error) {
	const getL1InfoRootSQL = `SELECT block_num, timestamp, mainnet_exit_root, rollup_exit_root, global_exit_root, prev_block_hash, l1_info_root, %s
		FROM state.exit_root 
		WHERE %s IS NOT NULL AND l1_info_root=$1`
	sql := fmt.Sprintf(getL1InfoRootSQL, indexFieldName, indexFieldName)
	var entry state.L1InfoTreeExitRootStorageEntry
	e := p.getExecQuerier(dbTx)
	err := e.QueryRow(ctx, sql, l1InfoRoot).Scan(&entry.BlockNumber, &entry.Timestamp, &entry.MainnetExitRoot, &entry.RollupExitRoot, &entry.GlobalExitRoot.GlobalExitRoot,
		&entry.PreviousBlockHash, &entry.L1InfoTreeRoot, &entry.L1InfoTreeIndex)
	if !errors.Is(err, pgx.ErrNoRows) {
		return entry, err
	}
	return entry, nil
}

func (p *PostgresStorage) GetL1InfoRootLeafByIndex(ctx context.Context, l1InfoTreeIndex uint32, dbTx pgx.Tx) (state.L1InfoTreeExitRootStorageEntry, error) {
	return p.GetL1InfoRootLeafByIndexVx(ctx, l1InfoTreeIndex, dbTx, l1InfoTreeIndexFieldName)
}

func (p *PostgresStorage) GetL1InfoRootLeafByIndexVx(ctx context.Context, l1InfoTreeIndex uint32, dbTx pgx.Tx, indexFieldName string) (state.L1InfoTreeExitRootStorageEntry, error) {
	const getL1InfoRootByIndexSQL = `SELECT block_num, timestamp, mainnet_exit_root, rollup_exit_root, global_exit_root, prev_block_hash, l1_info_root, %s
		FROM state.exit_root 
		WHERE %s = $1`
	sql := fmt.Sprintf(getL1InfoRootByIndexSQL, indexFieldName, indexFieldName)
	var entry state.L1InfoTreeExitRootStorageEntry
	e := p.getExecQuerier(dbTx)
	err := e.QueryRow(ctx, sql, l1InfoTreeIndex).Scan(&entry.BlockNumber, &entry.Timestamp, &entry.MainnetExitRoot, &entry.RollupExitRoot, &entry.GlobalExitRoot.GlobalExitRoot,
		&entry.PreviousBlockHash, &entry.L1InfoTreeRoot, &entry.L1InfoTreeIndex)
	if !errors.Is(err, pgx.ErrNoRows) {
		return entry, err
	}
	return entry, nil
}

func (p *PostgresStorage) GetLeavesByL1InfoRoot(ctx context.Context, l1InfoRoot common.Hash, dbTx pgx.Tx) ([]state.L1InfoTreeExitRootStorageEntry, error) {
	return p.GetLeavesByL1InfoRootVx(ctx, l1InfoRoot, dbTx, l1InfoTreeIndexFieldName)
}

func (p *PostgresStorage) GetLeavesByL1InfoRootVx(ctx context.Context, l1InfoRoot common.Hash, dbTx pgx.Tx, indexFieldName string) ([]state.L1InfoTreeExitRootStorageEntry, error) {
	// TODO: Optimize this query
	const getLeavesByL1InfoRootSQL = `SELECT block_num, timestamp, mainnet_exit_root, rollup_exit_root, global_exit_root, prev_block_hash, l1_info_root, %s
		FROM state.exit_root 
		WHERE %s IS NOT NULL AND %s <= (SELECT %s FROM state.exit_root WHERE l1_info_root=$1)
		ORDER BY %s ASC`
	sql := fmt.Sprintf(getLeavesByL1InfoRootSQL, indexFieldName, indexFieldName, indexFieldName, indexFieldName, indexFieldName)
	e := p.getExecQuerier(dbTx)
	rows, err := e.Query(ctx, sql, l1InfoRoot)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entries := make([]state.L1InfoTreeExitRootStorageEntry, 0)

	for rows.Next() {
		entry, err := scanL1InfoTreeExitRootStorageEntry(rows)
		if err != nil {
			return entries, err
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func scanL1InfoTreeExitRootStorageEntry(row pgx.Row) (state.L1InfoTreeExitRootStorageEntry, error) {
	entry := state.L1InfoTreeExitRootStorageEntry{}

	if err := row.Scan(
		&entry.BlockNumber, &entry.Timestamp, &entry.MainnetExitRoot, &entry.RollupExitRoot, &entry.GlobalExitRoot.GlobalExitRoot,
		&entry.PreviousBlockHash, &entry.L1InfoTreeRoot, &entry.L1InfoTreeIndex); err != nil {
		return entry, err
	}
	return entry, nil
}
