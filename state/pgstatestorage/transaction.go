package pgstatestorage

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
)

const maxTopics = 4

// GetTxsOlderThanNL1BlocksUntilTxHash get txs hashes to delete from tx pool from the oldest processed transaction to the latest
// txn that has been virtualized.
// Works like GetTxsOlderThanNL1Blocks but pulls hashes until earliestTxHash
func (p *PostgresStorage) GetTxsOlderThanNL1BlocksUntilTxHash(ctx context.Context, nL1Blocks uint64, earliestTxHash common.Hash, dbTx pgx.Tx) ([]common.Hash, error) {
	var earliestBatchNum, latestBatchNum, blockNum uint64
	const getLatestBatchNumByBlockNumFromVirtualBatch = "SELECT batch_num FROM state.virtual_batch WHERE block_num <= $1 ORDER BY batch_num DESC LIMIT 1"
	const getTxsHashesBeforeBatchNum = "SELECT hash FROM state.transaction JOIN state.l2block ON state.transaction.l2_block_num = state.l2block.block_num AND state.l2block.batch_num >= $1 AND state.l2block.batch_num <= $2"

	// Get lower bound batch_num which is the batch num from the oldest tx in txpool
	const getEarliestBatchNumByTxHashFromVirtualBatch = `SELECT batch_num
	FROM state.transaction
	JOIN state.l2block ON
		state.transaction.l2_block_num = state.l2block.block_num AND state.transaction.hash = $1`

	e := p.getExecQuerier(dbTx)

	err := e.QueryRow(ctx, getLastBlockNumSQL).Scan(&blockNum)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	blockNum = blockNum - nL1Blocks
	if blockNum <= 0 {
		return nil, errors.New("blockNumDiff is too big, there are no txs to delete")
	}

	err = e.QueryRow(ctx, getEarliestBatchNumByTxHashFromVirtualBatch, earliestTxHash.String()).Scan(&earliestBatchNum)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	err = e.QueryRow(ctx, getLatestBatchNumByBlockNumFromVirtualBatch, blockNum).Scan(&latestBatchNum)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	rows, err := e.Query(ctx, getTxsHashesBeforeBatchNum, earliestBatchNum, latestBatchNum)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrNotFound
	} else if err != nil {
		return nil, err
	}
	hashes := make([]common.Hash, 0, len(rows.RawValues()))
	for rows.Next() {
		var hash string
		err := rows.Scan(&hash)
		if err != nil {
			return nil, err
		}
		hashes = append(hashes, common.HexToHash(hash))
	}

	return hashes, nil
}

// GetTxsOlderThanNL1Blocks get txs hashes to delete from tx pool
func (p *PostgresStorage) GetTxsOlderThanNL1Blocks(ctx context.Context, nL1Blocks uint64, dbTx pgx.Tx) ([]common.Hash, error) {
	var batchNum, blockNum uint64
	const getBatchNumByBlockNumFromVirtualBatch = "SELECT batch_num FROM state.virtual_batch WHERE block_num <= $1 ORDER BY batch_num DESC LIMIT 1"
	const getTxsHashesBeforeBatchNum = "SELECT hash FROM state.transaction JOIN state.l2block ON state.transaction.l2_block_num = state.l2block.block_num AND state.l2block.batch_num <= $1"

	e := p.getExecQuerier(dbTx)

	err := e.QueryRow(ctx, getLastBlockNumSQL).Scan(&blockNum)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	blockNum = blockNum - nL1Blocks
	if blockNum <= 0 {
		return nil, errors.New("blockNumDiff is too big, there are no txs to delete")
	}

	err = e.QueryRow(ctx, getBatchNumByBlockNumFromVirtualBatch, blockNum).Scan(&batchNum)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrNotFound
	} else if err != nil {
		return nil, err
	}
	rows, err := e.Query(ctx, getTxsHashesBeforeBatchNum, batchNum)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrNotFound
	} else if err != nil {
		return nil, err
	}
	hashes := make([]common.Hash, 0, len(rows.RawValues()))
	for rows.Next() {
		var hash string
		err := rows.Scan(&hash)
		if err != nil {
			return nil, err
		}
		hashes = append(hashes, common.HexToHash(hash))
	}

	return hashes, nil
}

// GetEncodedTransactionsByBatchNumber returns the encoded field of all
// transactions in the given batch.
func (p *PostgresStorage) GetEncodedTransactionsByBatchNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (encodedTxs []string, effectivePercentages []uint8, err error) {
	const getEncodedTransactionsByBatchNumberSQL = "SELECT encoded, COALESCE(effective_percentage, 255) FROM state.transaction t INNER JOIN state.l2block b ON t.l2_block_num = b.block_num WHERE b.batch_num = $1 ORDER BY l2_block_num ASC"

	e := p.getExecQuerier(dbTx)
	rows, err := e.Query(ctx, getEncodedTransactionsByBatchNumberSQL, batchNumber)
	if !errors.Is(err, pgx.ErrNoRows) && err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	encodedTxs = make([]string, 0, len(rows.RawValues()))
	effectivePercentages = make([]uint8, 0, len(rows.RawValues()))

	for rows.Next() {
		var (
			encoded             string
			effectivePercentage uint8
		)
		err := rows.Scan(&encoded, &effectivePercentage)
		if err != nil {
			return nil, nil, err
		}

		encodedTxs = append(encodedTxs, encoded)
		effectivePercentages = append(effectivePercentages, effectivePercentage)
	}

	return encodedTxs, effectivePercentages, nil
}

// GetTransactionsByBatchNumber returns the transactions in the given batch.
func (p *PostgresStorage) GetTransactionsByBatchNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (txs []types.Transaction, effectivePercentages []uint8, err error) {
	var encodedTxs []string
	encodedTxs, effectivePercentages, err = p.GetEncodedTransactionsByBatchNumber(ctx, batchNumber, dbTx)
	if err != nil {
		return nil, nil, err
	}

	for i := 0; i < len(encodedTxs); i++ {
		tx, err := state.DecodeTx(encodedTxs[i])
		if err != nil {
			return nil, nil, err
		}
		txs = append(txs, *tx)
	}

	return txs, effectivePercentages, nil
}

// GetTxsHashesByBatchNumber returns the hashes of the transactions in the
// given batch.
func (p *PostgresStorage) GetTxsHashesByBatchNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (encoded []common.Hash, err error) {
	const getTransactionHashesByBatchNumberSQL = "SELECT hash FROM state.transaction t INNER JOIN state.l2block b ON t.l2_block_num = b.block_num WHERE b.batch_num = $1 ORDER BY l2_block_num ASC"

	e := p.getExecQuerier(dbTx)
	rows, err := e.Query(ctx, getTransactionHashesByBatchNumberSQL, batchNumber)
	if !errors.Is(err, pgx.ErrNoRows) && err != nil {
		return nil, err
	}
	defer rows.Close()

	txs := make([]common.Hash, 0, len(rows.RawValues()))

	for rows.Next() {
		var hexHash string
		err := rows.Scan(&hexHash)
		if err != nil {
			return nil, err
		}

		txs = append(txs, common.HexToHash(hexHash))
	}
	return txs, nil
}

// GetTransactionByHash gets a transaction accordingly to the provided transaction hash
func (p *PostgresStorage) GetTransactionByHash(ctx context.Context, transactionHash common.Hash, dbTx pgx.Tx) (*types.Transaction, error) {
	var encoded string
	const getTransactionByHashSQL = "SELECT transaction.encoded FROM state.transaction WHERE hash = $1"

	q := p.getExecQuerier(dbTx)
	err := q.QueryRow(ctx, getTransactionByHashSQL, transactionHash.String()).Scan(&encoded)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	tx, err := state.DecodeTx(encoded)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// GetTransactionByL2Hash gets a transaction accordingly to the provided transaction l2 hash
func (p *PostgresStorage) GetTransactionByL2Hash(ctx context.Context, l2TxHash common.Hash, dbTx pgx.Tx) (*types.Transaction, error) {
	var encoded string
	const getTransactionByHashSQL = "SELECT transaction.encoded FROM state.transaction WHERE l2_hash = $1"

	q := p.getExecQuerier(dbTx)
	err := q.QueryRow(ctx, getTransactionByHashSQL, l2TxHash.String()).Scan(&encoded)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	tx, err := state.DecodeTx(encoded)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// GetTransactionReceipt gets a transaction receipt accordingly to the provided transaction hash
func (p *PostgresStorage) GetTransactionReceipt(ctx context.Context, transactionHash common.Hash, dbTx pgx.Tx) (*types.Receipt, error) {
	var txHash, encodedTx, contractAddress, l2BlockHash string
	var l2BlockNum uint64
	var effective_gas_price *uint64

	const getReceiptSQL = `
		SELECT 
			r.tx_index,
			r.tx_hash,
		    r.type,
			r.post_state,
			r.status,
			r.cumulative_gas_used,
			r.gas_used,
			r.contract_address,
			r.effective_gas_price,
			t.encoded,
			t.l2_block_num,
			b.block_hash
	      FROM state.receipt r
		 INNER JOIN state.transaction t
		    ON t.hash = r.tx_hash
		 INNER JOIN state.l2block b
		    ON b.block_num = t.l2_block_num
		 WHERE r.tx_hash = $1`

	receipt := types.Receipt{}
	q := p.getExecQuerier(dbTx)
	err := q.QueryRow(ctx, getReceiptSQL, transactionHash.String()).
		Scan(&receipt.TransactionIndex,
			&txHash,
			&receipt.Type,
			&receipt.PostState,
			&receipt.Status,
			&receipt.CumulativeGasUsed,
			&receipt.GasUsed,
			&contractAddress,
			&effective_gas_price,
			&encodedTx,
			&l2BlockNum,
			&l2BlockHash,
		)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	receipt.TxHash = common.HexToHash(txHash)
	receipt.ContractAddress = common.HexToAddress(contractAddress)

	logs, err := p.getTransactionLogs(ctx, transactionHash, dbTx)
	if !errors.Is(err, pgx.ErrNoRows) && err != nil {
		return nil, err
	}

	receipt.BlockNumber = big.NewInt(0).SetUint64(l2BlockNum)
	receipt.BlockHash = common.HexToHash(l2BlockHash)
	if effective_gas_price != nil {
		receipt.EffectiveGasPrice = big.NewInt(0).SetUint64(*effective_gas_price)
	}
	receipt.Logs = logs
	receipt.Bloom = types.CreateBloom(types.Receipts{&receipt})

	return &receipt, nil
}

// GetTransactionByL2BlockHashAndIndex gets a transaction accordingly to the block hash and transaction index provided.
// since we only have a single transaction per l2 block, any index different from 0 will return a not found result
func (p *PostgresStorage) GetTransactionByL2BlockHashAndIndex(ctx context.Context, blockHash common.Hash, index uint64, dbTx pgx.Tx) (*types.Transaction, error) {
	var encoded string
	q := p.getExecQuerier(dbTx)
	const query = `
        SELECT t.encoded
          FROM state.transaction t
         INNER JOIN state.l2block b
            ON t.l2_block_num = b.block_num
         INNER JOIN state.receipt r
            ON r.tx_hash = t.hash
         WHERE b.block_hash = $1
           AND r.tx_index = $2`
	err := q.QueryRow(ctx, query, blockHash.String(), index).Scan(&encoded)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	tx, err := state.DecodeTx(encoded)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// GetTransactionByL2BlockNumberAndIndex gets a transaction accordingly to the block number and transaction index provided.
// since we only have a single transaction per l2 block, any index different from 0 will return a not found result
func (p *PostgresStorage) GetTransactionByL2BlockNumberAndIndex(ctx context.Context, blockNumber uint64, index uint64, dbTx pgx.Tx) (*types.Transaction, error) {
	var encoded string
	const getTransactionByL2BlockNumberAndIndexSQL = "SELECT t.encoded FROM state.transaction t WHERE t.l2_block_num = $1 AND 0 = $2"

	q := p.getExecQuerier(dbTx)
	err := q.QueryRow(ctx, getTransactionByL2BlockNumberAndIndexSQL, blockNumber, index).Scan(&encoded)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	tx, err := state.DecodeTx(encoded)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

// getTransactionLogs returns the logs of a transaction by transaction hash
func (p *PostgresStorage) getTransactionLogs(ctx context.Context, transactionHash common.Hash, dbTx pgx.Tx) ([]*types.Log, error) {
	q := p.getExecQuerier(dbTx)

	const getTransactionLogsSQL = `
	SELECT t.l2_block_num, b.block_hash, l.tx_hash, r.tx_index, l.log_index, l.address, l.data, l.topic0, l.topic1, l.topic2, l.topic3
	FROM state.log l
	INNER JOIN state.transaction t ON t.hash = l.tx_hash
	INNER JOIN state.l2block b ON b.block_num = t.l2_block_num 
	INNER JOIN state.receipt r ON r.tx_hash = t.hash
	WHERE t.hash = $1
	ORDER BY l.log_index ASC`
	rows, err := q.Query(ctx, getTransactionLogsSQL, transactionHash.String())
	if !errors.Is(err, pgx.ErrNoRows) && err != nil {
		return nil, err
	}
	return scanLogs(rows)
}

func scanLogs(rows pgx.Rows) ([]*types.Log, error) {
	defer rows.Close()

	logs := make([]*types.Log, 0, len(rows.RawValues()))

	for rows.Next() {
		if rows.Err() != nil {
			return nil, rows.Err()
		}

		var log types.Log
		var txIndex uint
		var blockHash, txHash, logAddress, logData string
		var topic0, topic1, topic2, topic3 *string

		err := rows.Scan(&log.BlockNumber, &blockHash, &txHash, &txIndex, &log.Index,
			&logAddress, &logData, &topic0, &topic1, &topic2, &topic3)
		if err != nil {
			return nil, err
		}

		log.BlockHash = common.HexToHash(blockHash)
		log.TxHash = common.HexToHash(txHash)
		log.Address = common.HexToAddress(logAddress)
		log.TxIndex = txIndex
		log.Data, err = hex.DecodeHex(logData)
		if err != nil {
			return nil, err
		}

		log.Topics = []common.Hash{}
		if topic0 != nil {
			log.Topics = append(log.Topics, common.HexToHash(*topic0))
		}

		if topic1 != nil {
			log.Topics = append(log.Topics, common.HexToHash(*topic1))
		}

		if topic2 != nil {
			log.Topics = append(log.Topics, common.HexToHash(*topic2))
		}

		if topic3 != nil {
			log.Topics = append(log.Topics, common.HexToHash(*topic3))
		}

		logs = append(logs, &log)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return logs, nil
}

// GetTxsByBlockNumber returns all the txs in a given block
func (p *PostgresStorage) GetTxsByBlockNumber(ctx context.Context, blockNumber uint64, dbTx pgx.Tx) ([]*types.Transaction, error) {
	const getTxsByBlockNumSQL = `SELECT t.encoded 
	   FROM state.transaction t
	   JOIN state.receipt r
	     ON t.hash = r.tx_hash
	  WHERE t.l2_block_num = $1
	    AND r.block_num = $1
	  ORDER by r.tx_index ASC`

	q := p.getExecQuerier(dbTx)
	rows, err := q.Query(ctx, getTxsByBlockNumSQL, blockNumber)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	defer rows.Close()

	txs := make([]*types.Transaction, 0, len(rows.RawValues()))
	var encoded string
	for rows.Next() {
		if err = rows.Scan(&encoded); err != nil {
			return nil, err
		}

		tx, err := state.DecodeTx(encoded)
		if err != nil {
			return nil, err
		}
		txs = append(txs, tx)
	}

	return txs, nil
}

// GetTxsByBatchNumber returns all the txs in a given batch
func (p *PostgresStorage) GetTxsByBatchNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) ([]*types.Transaction, error) {
	q := p.getExecQuerier(dbTx)

	const getTxsByBatchNumSQL = `
        SELECT encoded
          FROM state.transaction t
         INNER JOIN state.l2block b
            ON b.block_num = t.l2_block_num
         WHERE b.batch_num = $1`

	rows, err := q.Query(ctx, getTxsByBatchNumSQL, batchNumber)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	defer rows.Close()

	txs := make([]*types.Transaction, 0, len(rows.RawValues()))
	var encoded string
	for rows.Next() {
		if err = rows.Scan(&encoded); err != nil {
			return nil, err
		}

		tx, err := state.DecodeTx(encoded)
		if err != nil {
			return nil, err
		}
		txs = append(txs, tx)
	}

	return txs, nil
}

// AddReceipt adds a new receipt to the State Store
func (p *PostgresStorage) AddReceipt(ctx context.Context, receipt *types.Receipt, imStateRoot common.Hash, dbTx pgx.Tx) error {
	e := p.getExecQuerier(dbTx)

	var effectiveGasPrice *uint64

	if receipt.EffectiveGasPrice != nil {
		egf := receipt.EffectiveGasPrice.Uint64()
		effectiveGasPrice = &egf
	}

	const addReceiptSQL = `
        INSERT INTO state.receipt (tx_hash, type, post_state, status, cumulative_gas_used, gas_used, effective_gas_price, block_num, tx_index, contract_address, im_state_root)
                           VALUES (     $1,   $2,         $3,     $4,                  $5,       $6,        		  $7,        $8,       $9,			    $10,           $11)`
	_, err := e.Exec(ctx, addReceiptSQL, receipt.TxHash.String(), receipt.Type, receipt.PostState, receipt.Status, receipt.CumulativeGasUsed, receipt.GasUsed, effectiveGasPrice, receipt.BlockNumber.Uint64(), receipt.TransactionIndex, receipt.ContractAddress.String(), imStateRoot.Bytes())
	return err
}

// AddReceipts adds a list of receipts to the State Store
func (p *PostgresStorage) AddReceipts(ctx context.Context, receipts []*types.Receipt, imStateRoots []common.Hash, dbTx pgx.Tx) error {
	if len(receipts) == 0 {
		return nil
	}

	receiptRows := [][]interface{}{}

	for i, receipt := range receipts {
		var egp uint64
		if receipt.EffectiveGasPrice != nil {
			egp = receipt.EffectiveGasPrice.Uint64()
		}
		receiptRow := []interface{}{receipt.TxHash.String(), receipt.Type, receipt.PostState, receipt.Status, receipt.CumulativeGasUsed, receipt.GasUsed, egp, receipt.BlockNumber.Uint64(), receipt.TransactionIndex, receipt.ContractAddress.String(), imStateRoots[i].Bytes()}
		receiptRows = append(receiptRows, receiptRow)
	}

	_, err := dbTx.CopyFrom(ctx, pgx.Identifier{"state", "receipt"},
		[]string{"tx_hash", "type", "post_state", "status", "cumulative_gas_used", "gas_used", "effective_gas_price", "block_num", "tx_index", "contract_address", "im_state_root"},
		pgx.CopyFromRows(receiptRows))

	return err
}

// AddLog adds a new log to the State Store
func (p *PostgresStorage) AddLog(ctx context.Context, l *types.Log, dbTx pgx.Tx) error {
	const addLogSQL = `INSERT INTO state.log (tx_hash, log_index, address, data, topic0, topic1, topic2, topic3)
	                                  VALUES (     $1,        $2,      $3,   $4,     $5,     $6,     $7,     $8)`

	var topicsAsHex [maxTopics]*string
	for i := 0; i < len(l.Topics); i++ {
		topicHex := l.Topics[i].String()
		topicsAsHex[i] = &topicHex
	}

	e := p.getExecQuerier(dbTx)
	_, err := e.Exec(ctx, addLogSQL,
		l.TxHash.String(), l.Index, l.Address.String(), hex.EncodeToHex(l.Data),
		topicsAsHex[0], topicsAsHex[1], topicsAsHex[2], topicsAsHex[3])
	return err
}

// // AddLogs adds a list of logs to the State Store
func (p *PostgresStorage) AddLogs(ctx context.Context, logs []*types.Log, dbTx pgx.Tx) error {
	if len(logs) == 0 {
		return nil
	}

	logsRows := [][]interface{}{}

	for _, log := range logs {
		var topicsAsHex [maxTopics]*string
		for i := 0; i < len(log.Topics); i++ {
			topicHex := log.Topics[i].String()
			topicsAsHex[i] = &topicHex
		}
		logRow := []interface{}{log.TxHash.String(), log.Index, log.Address.String(), hex.EncodeToHex(log.Data), topicsAsHex[0], topicsAsHex[1], topicsAsHex[2], topicsAsHex[3]}
		logsRows = append(logsRows, logRow)
	}

	_, err := dbTx.CopyFrom(ctx, pgx.Identifier{"state", "log"},
		[]string{"tx_hash", "log_index", "address", "data", "topic0", "topic1", "topic2", "topic3"},
		pgx.CopyFromRows(logsRows))

	return err
}

// GetTransactionEGPLogByHash gets the EGP log accordingly to the provided transaction hash
func (p *PostgresStorage) GetTransactionEGPLogByHash(ctx context.Context, transactionHash common.Hash, dbTx pgx.Tx) (*state.EffectiveGasPriceLog, error) {
	var (
		egpLogData []byte
		egpLog     state.EffectiveGasPriceLog
	)
	const getTransactionByHashSQL = "SELECT egp_log FROM state.transaction WHERE hash = $1"

	q := p.getExecQuerier(dbTx)
	err := q.QueryRow(ctx, getTransactionByHashSQL, transactionHash.String()).Scan(&egpLogData)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	err = json.Unmarshal(egpLogData, &egpLog)
	if err != nil {
		return nil, err
	}

	return &egpLog, nil
}

// GetL2TxHashByTxHash gets the L2 Hash from the tx found by the provided tx hash
func (p *PostgresStorage) GetL2TxHashByTxHash(ctx context.Context, hash common.Hash, dbTx pgx.Tx) (*common.Hash, error) {
	const getTransactionByHashSQL = "SELECT transaction.l2_hash FROM state.transaction WHERE hash = $1"

	var l2HashHex *string
	q := p.getExecQuerier(dbTx)
	err := q.QueryRow(ctx, getTransactionByHashSQL, hash.String()).Scan(&l2HashHex)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	if l2HashHex == nil {
		return nil, nil
	}

	l2Hash := common.HexToHash(*l2HashHex)
	return &l2Hash, nil
}
