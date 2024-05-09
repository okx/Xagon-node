package pgstatestorage

import (
	"context"
	"errors"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// PostgresStorage implements the Storage interface
type PostgresStorage struct {
	cfg state.Config
	*pgxpool.Pool
}

// NewPostgresStorage creates a new StateDB
func NewPostgresStorage(cfg state.Config, db *pgxpool.Pool) *PostgresStorage {
	return &PostgresStorage{
		cfg,
		db,
	}
}

// getExecQuerier determines which execQuerier to use, dbTx or the main pgxpool
func (p *PostgresStorage) getExecQuerier(dbTx pgx.Tx) ExecQuerier {
	if dbTx != nil {
		return dbTx
	}
	return p
}

// ResetToL1BlockNumber resets the state to a block for the given DB tx
func (p *PostgresStorage) ResetToL1BlockNumber(ctx context.Context, blockNumber uint64, dbTx pgx.Tx) error {
	e := p.getExecQuerier(dbTx)
	const resetSQL = "DELETE FROM state.block WHERE block_num > $1"
	if _, err := e.Exec(ctx, resetSQL, blockNumber); err != nil {
		return err
	}

	return nil
}

// ResetForkID resets the state to reprocess the newer batches with the correct forkID
func (p *PostgresStorage) ResetForkID(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) error {
	e := p.getExecQuerier(dbTx)
	const resetVirtualStateSQL = "delete from state.block where block_num >=(select min(block_num) from state.virtual_batch where batch_num >= $1)"
	if _, err := e.Exec(ctx, resetVirtualStateSQL, batchNumber); err != nil {
		return err
	}
	err := p.ResetTrustedState(ctx, batchNumber-1, dbTx)
	if err != nil {
		return err
	}

	// Delete proofs for higher batches
	const deleteProofsSQL = "delete from state.proof where batch_num >= $1 or (batch_num <= $1 and batch_num_final  >= $1)"
	if _, err := e.Exec(ctx, deleteProofsSQL, batchNumber); err != nil {
		return err
	}

	return nil
}

// ResetTrustedState removes the batches with number greater than the given one
// from the database.
func (p *PostgresStorage) ResetTrustedState(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) error {
	const resetTrustedStateSQL = "DELETE FROM state.batch WHERE batch_num > $1"
	e := p.getExecQuerier(dbTx)
	if _, err := e.Exec(ctx, resetTrustedStateSQL, batchNumber); err != nil {
		return err
	}
	return nil
}

// GetProcessingContext returns the processing context for the given batch.
func (p *PostgresStorage) GetProcessingContext(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*state.ProcessingContext, error) {
	const getProcessingContextSQL = "SELECT batch_num, global_exit_root, timestamp, coinbase, forced_batch_num from state.batch WHERE batch_num = $1"

	e := p.getExecQuerier(dbTx)
	row := e.QueryRow(ctx, getProcessingContextSQL, batchNumber)
	processingContext := state.ProcessingContext{}
	var (
		gerStr      string
		coinbaseStr string
	)
	if err := row.Scan(
		&processingContext.BatchNumber,
		&gerStr,
		&processingContext.Timestamp,
		&coinbaseStr,
		&processingContext.ForcedBatchNum,
	); errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrStateNotSynchronized
	} else if err != nil {
		return nil, err
	}
	processingContext.GlobalExitRoot = common.HexToHash(gerStr)
	processingContext.Coinbase = common.HexToAddress(coinbaseStr)

	return &processingContext, nil
}

// GetStateRootByBatchNumber get state root by batch number
func (p *PostgresStorage) GetStateRootByBatchNumber(ctx context.Context, batchNum uint64, dbTx pgx.Tx) (common.Hash, error) {
	const query = "SELECT state_root FROM state.batch WHERE batch_num = $1"
	var stateRootStr string
	e := p.getExecQuerier(dbTx)
	err := e.QueryRow(ctx, query, batchNum).Scan(&stateRootStr)
	if errors.Is(err, pgx.ErrNoRows) {
		return common.Hash{}, state.ErrNotFound
	} else if err != nil {
		return common.Hash{}, err
	}
	return common.HexToHash(stateRootStr), nil
}

// GetLogsByBlockNumber get all the logs from a specific block ordered by tx index and log index
func (p *PostgresStorage) GetLogsByBlockNumber(ctx context.Context, blockNumber uint64, dbTx pgx.Tx) ([]*types.Log, error) {
	const query = `
      SELECT t.l2_block_num, b.block_hash, l.tx_hash, r.tx_index, l.log_index, l.address, l.data, l.topic0, l.topic1, l.topic2, l.topic3
        FROM state.log l
       INNER JOIN state.transaction t ON t.hash = l.tx_hash
       INNER JOIN state.l2block b ON b.block_num = t.l2_block_num
       INNER JOIN state.receipt r ON r.tx_hash = t.hash
       WHERE b.block_num = $1
       ORDER BY r.tx_index ASC, l.log_index ASC`

	q := p.getExecQuerier(dbTx)
	rows, err := q.Query(ctx, query, blockNumber)
	if err != nil {
		return nil, err
	}

	return scanLogs(rows)
}

// GetLogs returns the logs that match the filter
func (p *PostgresStorage) GetLogs(ctx context.Context, fromBlock uint64, toBlock uint64, addresses []common.Address, topics [][]common.Hash, blockHash *common.Hash, since *time.Time, dbTx pgx.Tx) ([]*types.Log, error) {
	// query parts
	const queryCount = `SELECT count(*) `
	const querySelect = `SELECT t.l2_block_num, b.block_hash, l.tx_hash, r.tx_index, l.log_index, l.address, l.data, l.topic0, l.topic1, l.topic2, l.topic3 `

	const queryBody = `FROM state.log l
       INNER JOIN state.transaction t ON t.hash = l.tx_hash
       INNER JOIN state.l2block b ON b.block_num = t.l2_block_num
       INNER JOIN state.receipt r ON r.tx_hash = t.hash
       WHERE (l.address = any($1) OR $1 IS NULL)
         AND (l.topic0 = any($2) OR $2 IS NULL)
         AND (l.topic1 = any($3) OR $3 IS NULL)
         AND (l.topic2 = any($4) OR $4 IS NULL)
         AND (l.topic3 = any($5) OR $5 IS NULL)
         AND (b.created_at >= $6 OR $6 IS NULL) `

	const queryFilterByBlockHash = `AND b.block_hash = $7 `
	const queryFilterByBlockNumbers = `AND b.block_num BETWEEN $7 AND $8 `

	const queryOrder = `ORDER BY b.block_num ASC, r.tx_index ASC, l.log_index ASC`

	// count queries
	const queryToCountLogsByBlockHash = "" +
		queryCount +
		queryBody +
		queryFilterByBlockHash
	const queryToCountLogsByBlockNumbers = "" +
		queryCount +
		queryBody +
		queryFilterByBlockNumbers

	// select queries
	const queryToSelectLogsByBlockHash = "" +
		querySelect +
		queryBody +
		queryFilterByBlockHash +
		queryOrder
	const queryToSelectLogsByBlockNumbers = "" +
		querySelect +
		queryBody +
		queryFilterByBlockNumbers +
		queryOrder

	args := []interface{}{}

	// address filter
	if len(addresses) > 0 {
		args = append(args, p.addressesToHex(addresses))
	} else {
		args = append(args, nil)
	}

	// topic filters
	for i := 0; i < maxTopics; i++ {
		if len(topics) > i && len(topics[i]) > 0 {
			args = append(args, p.hashesToHex(topics[i]))
		} else {
			args = append(args, nil)
		}
	}

	// since filter
	args = append(args, since)

	// block filter
	var queryToCount string
	var queryToSelect string
	if blockHash != nil {
		args = append(args, blockHash.String())
		queryToCount = queryToCountLogsByBlockHash
		queryToSelect = queryToSelectLogsByBlockHash
	} else {
		if toBlock < fromBlock {
			return nil, state.ErrInvalidBlockRange
		}

		blockRange := toBlock - fromBlock
		if p.cfg.MaxLogsBlockRange > 0 && blockRange > p.cfg.MaxLogsBlockRange {
			return nil, state.ErrMaxLogsBlockRangeLimitExceeded
		}

		args = append(args, fromBlock, toBlock)
		queryToCount = queryToCountLogsByBlockNumbers
		queryToSelect = queryToSelectLogsByBlockNumbers
	}

	q := p.getExecQuerier(dbTx)
	if p.cfg.MaxLogsCount > 0 {
		var count uint64
		err := q.QueryRow(ctx, queryToCount, args...).Scan(&count)
		if err != nil {
			return nil, err
		}

		if count > p.cfg.MaxLogsCount {
			return nil, state.ErrMaxLogsCountLimitExceeded
		}
	}

	rows, err := q.Query(ctx, queryToSelect, args...)
	if err != nil {
		return nil, err
	}
	return scanLogs(rows)
}

func (p *PostgresStorage) addressesToHex(addresses []common.Address) []string {
	converted := make([]string, 0, len(addresses))

	for _, address := range addresses {
		converted = append(converted, address.String())
	}

	return converted
}

func (p *PostgresStorage) hashesToHex(hashes []common.Hash) []string {
	converted := make([]string, 0, len(hashes))

	for _, hash := range hashes {
		converted = append(converted, hash.String())
	}

	return converted
}

// AddTrustedReorg is used to store trusted reorgs
func (p *PostgresStorage) AddTrustedReorg(ctx context.Context, reorg *state.TrustedReorg, dbTx pgx.Tx) error {
	const insertTrustedReorgSQL = "INSERT INTO state.trusted_reorg (timestamp, batch_num, reason) VALUES (NOW(), $1, $2)"

	e := p.getExecQuerier(dbTx)
	_, err := e.Exec(ctx, insertTrustedReorgSQL, reorg.BatchNumber, reorg.Reason)
	return err
}

// CountReorgs returns the number of reorgs
func (p *PostgresStorage) CountReorgs(ctx context.Context, dbTx pgx.Tx) (uint64, error) {
	const countReorgsSQL = "SELECT COUNT(*) FROM state.trusted_reorg"

	var count uint64
	q := p.getExecQuerier(dbTx)
	err := q.QueryRow(ctx, countReorgsSQL).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetReorgedTransactions returns the transactions that were reorged
func (p *PostgresStorage) GetReorgedTransactions(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) ([]*types.Transaction, error) {
	const getReorgedTransactionsSql = "SELECT encoded FROM state.transaction t INNER JOIN state.l2block b ON t.l2_block_num = b.block_num WHERE b.batch_num >= $1 ORDER BY l2_block_num ASC"
	e := p.getExecQuerier(dbTx)
	rows, err := e.Query(ctx, getReorgedTransactionsSql, batchNumber)
	if !errors.Is(err, pgx.ErrNoRows) && err != nil {
		return nil, err
	}
	defer rows.Close()

	txs := make([]*types.Transaction, 0, len(rows.RawValues()))

	for rows.Next() {
		if rows.Err() != nil {
			return nil, rows.Err()
		}
		var encodedTx string
		err := rows.Scan(&encodedTx)
		if err != nil {
			return nil, err
		}

		tx, err := state.DecodeTx(encodedTx)
		if err != nil {
			return nil, err
		}
		txs = append(txs, tx)
	}
	return txs, nil
}

// GetNativeBlockHashesInRange return the state root for the blocks in range
func (p *PostgresStorage) GetNativeBlockHashesInRange(ctx context.Context, fromBlock, toBlock uint64, dbTx pgx.Tx) ([]common.Hash, error) {
	const l2TxSQL = `
    SELECT l2b.state_root
      FROM state.l2block l2b
     WHERE block_num BETWEEN $1 AND $2
     ORDER BY l2b.block_num ASC`

	if toBlock < fromBlock {
		return nil, state.ErrInvalidBlockRange
	}

	blockRange := toBlock - fromBlock
	if p.cfg.MaxNativeBlockHashBlockRange > 0 && blockRange > p.cfg.MaxNativeBlockHashBlockRange {
		return nil, state.ErrMaxNativeBlockHashBlockRangeLimitExceeded
	}

	e := p.getExecQuerier(dbTx)
	rows, err := e.Query(ctx, l2TxSQL, fromBlock, toBlock)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	nativeBlockHashes := []common.Hash{}

	for rows.Next() {
		var nativeBlockHash string
		err := rows.Scan(&nativeBlockHash)
		if err != nil {
			return nil, err
		}
		nativeBlockHashes = append(nativeBlockHashes, common.HexToHash(nativeBlockHash))
	}
	return nativeBlockHashes, nil
}
