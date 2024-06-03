package pgstatestorage

import (
	"context"
	"errors"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

// AddL1InfoRootToExitRoot adds a new entry in ExitRoot and returns index of L1InfoTree and error
func (p *PostgresStorage) AddL1InfoRootToExitRoot(ctx context.Context, exitRoot *state.L1InfoTreeExitRootStorageEntry, dbTx pgx.Tx) error {
	const addGlobalExitRootSQL = `
		INSERT INTO state.exit_root(block_num, timestamp, mainnet_exit_root, rollup_exit_root, global_exit_root, prev_block_hash, l1_info_root, l1_info_tree_index)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
	`
	e := p.getExecQuerier(dbTx)
	_, err := e.Exec(ctx, addGlobalExitRootSQL,
		exitRoot.BlockNumber, exitRoot.Timestamp, exitRoot.MainnetExitRoot, exitRoot.RollupExitRoot,
		exitRoot.GlobalExitRoot.GlobalExitRoot, exitRoot.PreviousBlockHash, exitRoot.L1InfoTreeRoot, exitRoot.L1InfoTreeIndex)
	return err
}

func (p *PostgresStorage) GetAllL1InfoRootEntries(ctx context.Context, dbTx pgx.Tx) ([]state.L1InfoTreeExitRootStorageEntry, error) {
	const getL1InfoRootSQL = `SELECT block_num, timestamp, mainnet_exit_root, rollup_exit_root, global_exit_root, prev_block_hash, l1_info_root, l1_info_tree_index
		FROM state.exit_root 
		WHERE l1_info_tree_index IS NOT NULL
		ORDER BY l1_info_tree_index`

	e := p.getExecQuerier(dbTx)
	rows, err := e.Query(ctx, getL1InfoRootSQL)
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
	const getL1InfoRootSQL = `SELECT block_num, timestamp, mainnet_exit_root, rollup_exit_root, global_exit_root, prev_block_hash, l1_info_root, l1_info_tree_index
		FROM state.exit_root 
		WHERE l1_info_tree_index IS NOT NULL AND block_num <= $1
		ORDER BY l1_info_tree_index DESC LIMIT 1`

	entry := state.L1InfoTreeExitRootStorageEntry{}

	e := p.getExecQuerier(nil)
	err := e.QueryRow(ctx, getL1InfoRootSQL, maxBlockNumber).Scan(&entry.BlockNumber, &entry.Timestamp, &entry.MainnetExitRoot, &entry.RollupExitRoot, &entry.GlobalExitRoot.GlobalExitRoot,
		&entry.PreviousBlockHash, &entry.L1InfoTreeRoot, &entry.L1InfoTreeIndex)

	if !errors.Is(err, pgx.ErrNoRows) {
		return entry, err
	}

	return entry, nil
}
func (p *PostgresStorage) GetLatestIndex(ctx context.Context, dbTx pgx.Tx) (uint32, error) {
	const getLatestIndexSQL = `SELECT max(l1_info_tree_index) as l1_info_tree_index FROM state.exit_root 
		WHERE l1_info_tree_index IS NOT NULL`
	var l1InfoTreeIndex *uint32
	e := p.getExecQuerier(dbTx)
	err := e.QueryRow(ctx, getLatestIndexSQL).Scan(&l1InfoTreeIndex)
	if err != nil {
		return 0, err
	}
	if l1InfoTreeIndex == nil {
		return 0, state.ErrNotFound
	}
	return *l1InfoTreeIndex, nil
}

func (p *PostgresStorage) GetL1InfoRootLeafByL1InfoRoot(ctx context.Context, l1InfoRoot common.Hash, dbTx pgx.Tx) (state.L1InfoTreeExitRootStorageEntry, error) {
	const getL1InfoRootSQL = `SELECT block_num, timestamp, mainnet_exit_root, rollup_exit_root, global_exit_root, prev_block_hash, l1_info_root, l1_info_tree_index
		FROM state.exit_root 
		WHERE l1_info_tree_index IS NOT NULL AND l1_info_root=$1`

	var entry state.L1InfoTreeExitRootStorageEntry
	e := p.getExecQuerier(dbTx)
	err := e.QueryRow(ctx, getL1InfoRootSQL, l1InfoRoot).Scan(&entry.BlockNumber, &entry.Timestamp, &entry.MainnetExitRoot, &entry.RollupExitRoot, &entry.GlobalExitRoot.GlobalExitRoot,
		&entry.PreviousBlockHash, &entry.L1InfoTreeRoot, &entry.L1InfoTreeIndex)
	if !errors.Is(err, pgx.ErrNoRows) {
		return entry, err
	}
	return entry, nil
}

func (p *PostgresStorage) GetL1InfoRootLeafByIndex(ctx context.Context, l1InfoTreeIndex uint32, dbTx pgx.Tx) (state.L1InfoTreeExitRootStorageEntry, error) {
	const getL1InfoRootByIndexSQL = `SELECT block_num, timestamp, mainnet_exit_root, rollup_exit_root, global_exit_root, prev_block_hash, l1_info_root, l1_info_tree_index
		FROM state.exit_root 
		WHERE l1_info_tree_index = $1`

	var entry state.L1InfoTreeExitRootStorageEntry
	e := p.getExecQuerier(dbTx)
	err := e.QueryRow(ctx, getL1InfoRootByIndexSQL, l1InfoTreeIndex).Scan(&entry.BlockNumber, &entry.Timestamp, &entry.MainnetExitRoot, &entry.RollupExitRoot, &entry.GlobalExitRoot.GlobalExitRoot,
		&entry.PreviousBlockHash, &entry.L1InfoTreeRoot, &entry.L1InfoTreeIndex)
	if !errors.Is(err, pgx.ErrNoRows) {
		return entry, err
	}
	return entry, nil
}

func (p *PostgresStorage) GetLeavesByL1InfoRoot(ctx context.Context, l1InfoRoot common.Hash, dbTx pgx.Tx) ([]state.L1InfoTreeExitRootStorageEntry, error) {
	// TODO: Optimize this query
	const getLeafsByL1InfoRootSQL = `SELECT block_num, timestamp, mainnet_exit_root, rollup_exit_root, global_exit_root, prev_block_hash, l1_info_root, l1_info_tree_index
		FROM state.exit_root 
		WHERE l1_info_tree_index IS NOT NULL AND l1_info_tree_index <= (SELECT l1_info_tree_index FROM state.exit_root WHERE l1_info_root=$1)
		ORDER BY l1_info_tree_index ASC`

	e := p.getExecQuerier(dbTx)
	rows, err := e.Query(ctx, getLeafsByL1InfoRootSQL, l1InfoRoot)
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
