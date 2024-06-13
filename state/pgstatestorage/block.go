package pgstatestorage

import (
	"context"
	"errors"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

const (
	getLastBlockNumSQL   = "SELECT block_num FROM state.block ORDER BY block_num DESC LIMIT 1"
	getBlockTimeByNumSQL = "SELECT received_at FROM state.block WHERE block_num = $1"
)

// AddBlock adds a new block to the State Store
func (p *PostgresStorage) AddBlock(ctx context.Context, block *state.Block, dbTx pgx.Tx) error {
	const addBlockSQL = "INSERT INTO state.block (block_num, block_hash, parent_hash, received_at, checked) VALUES ($1, $2, $3, $4, $5)"

	e := p.getExecQuerier(dbTx)
	_, err := e.Exec(ctx, addBlockSQL, block.BlockNumber, block.BlockHash.String(), block.ParentHash.String(), block.ReceivedAt, block.Checked)
	return err
}

// GetLastBlock returns the last L1 block.
func (p *PostgresStorage) GetLastBlock(ctx context.Context, dbTx pgx.Tx) (*state.Block, error) {
	var (
		blockHash  string
		parentHash string
		block      state.Block
	)
	const getLastBlockSQL = "SELECT block_num, block_hash, parent_hash, received_at, checked FROM state.block ORDER BY block_num DESC LIMIT 1"

	q := p.getExecQuerier(dbTx)

	err := q.QueryRow(ctx, getLastBlockSQL).Scan(&block.BlockNumber, &blockHash, &parentHash, &block.ReceivedAt, &block.Checked)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrStateNotSynchronized
	}
	block.BlockHash = common.HexToHash(blockHash)
	block.ParentHash = common.HexToHash(parentHash)
	return &block, err
}

// GetFirstUncheckedBlock returns the first L1 block that has not been checked from a given block number.
func (p *PostgresStorage) GetFirstUncheckedBlock(ctx context.Context, fromBlockNumber uint64, dbTx pgx.Tx) (*state.Block, error) {
	var (
		blockHash  string
		parentHash string
		block      state.Block
	)
	const getLastBlockSQL = "SELECT block_num, block_hash, parent_hash, received_at, checked FROM state.block  WHERE block_num>=$1 AND  checked=false ORDER BY block_num LIMIT 1"

	q := p.getExecQuerier(dbTx)

	err := q.QueryRow(ctx, getLastBlockSQL, fromBlockNumber).Scan(&block.BlockNumber, &blockHash, &parentHash, &block.ReceivedAt, &block.Checked)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrNotFound
	}
	block.BlockHash = common.HexToHash(blockHash)
	block.ParentHash = common.HexToHash(parentHash)
	return &block, err
}

func (p *PostgresStorage) GetUncheckedBlocks(ctx context.Context, fromBlockNumber uint64, toBlockNumber uint64, dbTx pgx.Tx) ([]*state.Block, error) {
	const getUncheckedBlocksSQL = "SELECT block_num, block_hash, parent_hash, received_at, checked FROM state.block WHERE block_num>=$1 AND block_num<=$2 AND checked=false ORDER BY block_num"

	q := p.getExecQuerier(dbTx)

	rows, err := q.Query(ctx, getUncheckedBlocksSQL, fromBlockNumber, toBlockNumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var blocks []*state.Block
	for rows.Next() {
		var (
			blockHash  string
			parentHash string
			block      state.Block
		)
		err := rows.Scan(&block.BlockNumber, &blockHash, &parentHash, &block.ReceivedAt, &block.Checked)
		if err != nil {
			return nil, err
		}
		block.BlockHash = common.HexToHash(blockHash)
		block.ParentHash = common.HexToHash(parentHash)
		blocks = append(blocks, &block)
	}
	return blocks, nil
}

// GetPreviousBlock gets the offset previous L1 block respect to latest.
func (p *PostgresStorage) GetPreviousBlock(ctx context.Context, offset uint64, dbTx pgx.Tx) (*state.Block, error) {
	var (
		blockHash  string
		parentHash string
		block      state.Block
	)
	const getPreviousBlockSQL = "SELECT block_num, block_hash, parent_hash, received_at,checked FROM state.block ORDER BY block_num DESC LIMIT 1 OFFSET $1"

	q := p.getExecQuerier(dbTx)

	err := q.QueryRow(ctx, getPreviousBlockSQL, offset).Scan(&block.BlockNumber, &blockHash, &parentHash, &block.ReceivedAt, &block.Checked)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrNotFound
	}
	block.BlockHash = common.HexToHash(blockHash)
	block.ParentHash = common.HexToHash(parentHash)
	return &block, err
}

// GetPreviousBlockToBlockNumber gets the previous L1 block respect blockNumber.
func (p *PostgresStorage) GetPreviousBlockToBlockNumber(ctx context.Context, blockNumber uint64, dbTx pgx.Tx) (*state.Block, error) {
	var (
		blockHash  string
		parentHash string
		block      state.Block
	)
	const getPreviousBlockSQL = "SELECT block_num, block_hash, parent_hash, received_at,checked FROM state.block WHERE block_num < $1 ORDER BY block_num DESC LIMIT 1 "

	q := p.getExecQuerier(dbTx)

	err := q.QueryRow(ctx, getPreviousBlockSQL, blockNumber).Scan(&block.BlockNumber, &blockHash, &parentHash, &block.ReceivedAt, &block.Checked)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrNotFound
	}
	block.BlockHash = common.HexToHash(blockHash)
	block.ParentHash = common.HexToHash(parentHash)
	return &block, err
}

// GetBlockByNumber returns the L1 block with the given number.
func (p *PostgresStorage) GetBlockByNumber(ctx context.Context, blockNumber uint64, dbTx pgx.Tx) (*state.Block, error) {
	var (
		blockHash  string
		parentHash string
		block      state.Block
	)
	const getBlockByNumberSQL = "SELECT block_num, block_hash, parent_hash, received_at,checked FROM state.block WHERE block_num = $1"

	q := p.getExecQuerier(dbTx)

	err := q.QueryRow(ctx, getBlockByNumberSQL, blockNumber).Scan(&block.BlockNumber, &blockHash, &parentHash, &block.ReceivedAt, &block.Checked)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrNotFound
	}
	block.BlockHash = common.HexToHash(blockHash)
	block.ParentHash = common.HexToHash(parentHash)
	return &block, err
}

// UpdateCheckedBlockByNumber update checked flag for a block
func (p *PostgresStorage) UpdateCheckedBlockByNumber(ctx context.Context, blockNumber uint64, newCheckedStatus bool, dbTx pgx.Tx) error {
	const query = `
    UPDATE state.block
       SET checked = $1 WHERE block_num = $2`

	e := p.getExecQuerier(dbTx)
	_, err := e.Exec(ctx, query, newCheckedStatus, blockNumber)
	return err
}
