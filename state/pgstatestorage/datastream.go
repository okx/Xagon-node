package pgstatestorage

import (
	"context"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

// GetDSGenesisBlock returns the genesis block
func (p *PostgresStorage) GetDSGenesisBlock(ctx context.Context, dbTx pgx.Tx) (*state.DSL2Block, error) {
	const genesisL2BlockSQL = `SELECT 0 as batch_num, l2b.block_num, l2b.received_at, '0x0000000000000000000000000000000000000000' as global_exit_root, '0x0000000000000000000000000000000000000000' as block_global_exit_root, l2b.header->>'miner' AS coinbase, 0 as fork_id, l2b.block_hash, l2b.state_root, '0x0000000000000000000000000000000000000000' as block_info_root
							FROM state.l2block l2b
							WHERE l2b.block_num  = 0`

	e := p.getExecQuerier(dbTx)

	row := e.QueryRow(ctx, genesisL2BlockSQL)

	l2block, err := scanL2Block(row)
	if err != nil {
		return nil, err
	}

	return l2block, nil
}

// GetDSL2Blocks returns the L2 blocks
func (p *PostgresStorage) GetDSL2Blocks(ctx context.Context, firstBatchNumber, lastBatchNumber uint64, dbTx pgx.Tx) ([]*state.DSL2Block, error) {
	const l2BlockSQL = `SELECT l2b.batch_num, l2b.block_num, l2b.received_at, b.global_exit_root, COALESCE(l2b.header->>'globalExitRoot', '') AS block_global_exit_root, l2b.header->>'miner' AS coinbase, f.fork_id, l2b.block_hash, l2b.state_root, COALESCE(l2b.header->>'blockInfoRoot', '') AS block_info_root
						FROM state.l2block l2b, state.batch b, state.fork_id f
						WHERE l2b.batch_num BETWEEN $1 AND $2 AND l2b.batch_num = b.batch_num AND l2b.batch_num between f.from_batch_num AND f.to_batch_num
						ORDER BY l2b.block_num ASC`
	e := p.getExecQuerier(dbTx)
	rows, err := e.Query(ctx, l2BlockSQL, firstBatchNumber, lastBatchNumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	l2blocks := make([]*state.DSL2Block, 0, len(rows.RawValues()))

	for rows.Next() {
		l2block, err := scanL2Block(rows)
		if err != nil {
			return nil, err
		}
		l2blocks = append(l2blocks, l2block)
	}

	return l2blocks, nil
}

func scanL2Block(row pgx.Row) (*state.DSL2Block, error) {
	l2Block := state.DSL2Block{}
	var (
		gerStr       string
		blockGERStr  string
		coinbaseStr  string
		timestamp    time.Time
		blockHashStr string
		stateRootStr string
		blockInfoStr string
	)
	if err := row.Scan(
		&l2Block.BatchNumber,
		&l2Block.L2BlockNumber,
		&timestamp,
		&gerStr,
		&blockGERStr,
		&coinbaseStr,
		&l2Block.ForkID,
		&blockHashStr,
		&stateRootStr,
		&blockInfoStr,
	); err != nil {
		return &l2Block, err
	}
	l2Block.GlobalExitRoot = common.HexToHash(gerStr)
	l2Block.Coinbase = common.HexToAddress(coinbaseStr)
	l2Block.Timestamp = uint64(timestamp.Unix())
	l2Block.BlockHash = common.HexToHash(blockHashStr)
	l2Block.StateRoot = common.HexToHash(stateRootStr)
	l2Block.BlockInfoRoot = common.HexToHash(blockInfoStr)

	if l2Block.ForkID >= state.FORKID_ETROG {
		l2Block.GlobalExitRoot = common.HexToHash(blockGERStr)
	}

	return &l2Block, nil
}

// GetDSL2Transactions returns the L2 transactions
func (p *PostgresStorage) GetDSL2Transactions(ctx context.Context, firstL2Block, lastL2Block uint64, dbTx pgx.Tx) ([]*state.DSL2Transaction, error) {
	const l2TxSQL = `SELECT l2_block_num, t.effective_percentage, t.encoded, r.post_state, r.im_state_root, r.tx_index
					 FROM state.transaction t, state.receipt r
					 WHERE l2_block_num BETWEEN $1 AND $2 AND r.tx_hash = t.hash
					 ORDER BY t.l2_block_num ASC, r.tx_index ASC`

	e := p.getExecQuerier(dbTx)
	rows, err := e.Query(ctx, l2TxSQL, firstL2Block, lastL2Block)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	l2Txs := make([]*state.DSL2Transaction, 0, len(rows.RawValues()))

	for rows.Next() {
		l2Tx, err := scanDSL2Transaction(rows)
		if err != nil {
			return nil, err
		}
		l2Txs = append(l2Txs, l2Tx)
	}

	return l2Txs, nil
}

func scanDSL2Transaction(row pgx.Row) (*state.DSL2Transaction, error) {
	l2Transaction := state.DSL2Transaction{}
	encoded := []byte{}
	postState := []byte{}
	imStateRoot := []byte{}

	if err := row.Scan(
		&l2Transaction.L2BlockNumber,
		&l2Transaction.EffectiveGasPricePercentage,
		&encoded,
		&postState,
		&imStateRoot,
		&l2Transaction.Index,
	); err != nil {
		return nil, err
	}
	tx, err := state.DecodeTx(string(encoded))
	if err != nil {
		return nil, err
	}

	binaryTxData, err := tx.MarshalBinary()
	if err != nil {
		return nil, err
	}

	l2Transaction.Encoded = binaryTxData
	l2Transaction.EncodedLength = uint32(len(l2Transaction.Encoded))
	l2Transaction.IsValid = 1
	l2Transaction.StateRoot = common.BytesToHash(postState)
	l2Transaction.ImStateRoot = common.BytesToHash(imStateRoot)
	return &l2Transaction, nil
}

// GetDSBatches returns the DS batches
func (p *PostgresStorage) GetDSBatches(ctx context.Context, firstBatchNumber, lastBatchNumber uint64, readWIPBatch bool, dbTx pgx.Tx) ([]*state.DSBatch, error) {
	var getBatchByNumberSQL = `
		SELECT b.batch_num, b.global_exit_root, b.local_exit_root, b.acc_input_hash, b.state_root, b.timestamp, b.coinbase, b.raw_txs_data, b.forced_batch_num, b.wip, f.fork_id, vb.timestamp_batch_etrog
		 FROM state.batch b
		 LEFT JOIN 
    		state.fork_id f ON b.batch_num BETWEEN f.from_batch_num AND f.to_batch_num
		 LEFT JOIN 
    		state.virtual_batch vb ON b.batch_num = vb.batch_num
		 WHERE b.batch_num >= $1 AND b.batch_num <= $2`

	if !readWIPBatch {
		getBatchByNumberSQL += " AND b.wip is false"
	}

	getBatchByNumberSQL += " ORDER BY b.batch_num ASC"

	e := p.getExecQuerier(dbTx)
	rows, err := e.Query(ctx, getBatchByNumberSQL, firstBatchNumber, lastBatchNumber)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	batches := make([]*state.DSBatch, 0, len(rows.RawValues()))

	for rows.Next() {
		batch, err := scanDSBatch(rows)
		if err != nil {
			return nil, err
		}
		batches = append(batches, &batch)
	}

	return batches, nil
}

func scanDSBatch(row pgx.Row) (state.DSBatch, error) {
	batch := state.DSBatch{}
	var (
		gerStr      string
		lerStr      *string
		aihStr      *string
		stateStr    *string
		coinbaseStr string
	)
	err := row.Scan(
		&batch.BatchNumber,
		&gerStr,
		&lerStr,
		&aihStr,
		&stateStr,
		&batch.Timestamp,
		&coinbaseStr,
		&batch.BatchL2Data,
		&batch.ForcedBatchNum,
		&batch.WIP,
		&batch.ForkID,
		&batch.EtrogTimestamp,
	)
	if err != nil {
		return batch, err
	}
	batch.GlobalExitRoot = common.HexToHash(gerStr)
	if lerStr != nil {
		batch.LocalExitRoot = common.HexToHash(*lerStr)
	}
	if stateStr != nil {
		batch.StateRoot = common.HexToHash(*stateStr)
	}
	if aihStr != nil {
		batch.AccInputHash = common.HexToHash(*aihStr)
	}

	batch.Coinbase = common.HexToAddress(coinbaseStr)
	return batch, nil
}
