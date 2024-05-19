package pgstatestorage

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

const (
	getLastBatchNumberSQL = "SELECT batch_num FROM state.batch ORDER BY batch_num DESC LIMIT 1"
)

// GetTimeForLatestBatchVirtualization returns the timestamp of the latest
// virtual batch.
func (p *PostgresStorage) GetTimeForLatestBatchVirtualization(ctx context.Context, dbTx pgx.Tx) (time.Time, error) {
	var (
		blockNum  uint64
		timestamp time.Time
	)
	const getLastVirtualBatchBlockNumSQL = "SELECT block_num FROM state.virtual_batch ORDER BY batch_num DESC LIMIT 1"

	e := p.getExecQuerier(dbTx)
	err := e.QueryRow(ctx, getLastVirtualBatchBlockNumSQL).Scan(&blockNum)

	if errors.Is(err, pgx.ErrNoRows) {
		return time.Time{}, state.ErrNotFound
	} else if err != nil {
		return time.Time{}, err
	}

	err = p.QueryRow(ctx, getBlockTimeByNumSQL, blockNum).Scan(&timestamp)

	if errors.Is(err, pgx.ErrNoRows) {
		return time.Time{}, state.ErrNotFound
	} else if err != nil {
		return time.Time{}, err
	}

	return timestamp, nil
}

// AddVerifiedBatch adds a new VerifiedBatch to the db
func (p *PostgresStorage) AddVerifiedBatch(ctx context.Context, verifiedBatch *state.VerifiedBatch, dbTx pgx.Tx) error {
	e := p.getExecQuerier(dbTx)
	const addVerifiedBatchSQL = "INSERT INTO state.verified_batch (block_num, batch_num, tx_hash, aggregator, state_root, is_trusted) VALUES ($1, $2, $3, $4, $5, $6)"
	_, err := e.Exec(ctx, addVerifiedBatchSQL, verifiedBatch.BlockNumber, verifiedBatch.BatchNumber, verifiedBatch.TxHash.String(), verifiedBatch.Aggregator.String(), verifiedBatch.StateRoot.String(), verifiedBatch.IsTrusted)
	return err
}

// GetVerifiedBatch get an L1 verifiedBatch.
func (p *PostgresStorage) GetVerifiedBatch(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*state.VerifiedBatch, error) {
	var (
		verifiedBatch state.VerifiedBatch
		txHash        string
		agg           string
		sr            string
	)

	const getVerifiedBatchSQL = `
    SELECT block_num, batch_num, tx_hash, aggregator, state_root, is_trusted
      FROM state.verified_batch
     WHERE batch_num = $1`

	e := p.getExecQuerier(dbTx)
	err := e.QueryRow(ctx, getVerifiedBatchSQL, batchNumber).Scan(&verifiedBatch.BlockNumber, &verifiedBatch.BatchNumber, &txHash, &agg, &sr, &verifiedBatch.IsTrusted)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrNotFound
	} else if err != nil {
		return nil, err
	}
	verifiedBatch.Aggregator = common.HexToAddress(agg)
	verifiedBatch.TxHash = common.HexToHash(txHash)
	verifiedBatch.StateRoot = common.HexToHash(sr)
	return &verifiedBatch, nil
}

// GetLastNBatches returns the last numBatches batches.
func (p *PostgresStorage) GetLastNBatches(ctx context.Context, numBatches uint, dbTx pgx.Tx) ([]*state.Batch, error) {
	const getLastNBatchesSQL = "SELECT batch_num, global_exit_root, local_exit_root, acc_input_hash, state_root, timestamp, coinbase, raw_txs_data, forced_batch_num, batch_resources, high_reserved_counters, wip from state.batch ORDER BY batch_num DESC LIMIT $1"

	e := p.getExecQuerier(dbTx)
	rows, err := e.Query(ctx, getLastNBatchesSQL, numBatches)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrStateNotSynchronized
	} else if err != nil {
		return nil, err
	}
	defer rows.Close()

	batches := make([]*state.Batch, 0, len(rows.RawValues()))

	for rows.Next() {
		batch, err := scanBatch(rows)
		if err != nil {
			return nil, err
		}
		batches = append(batches, &batch)
	}

	return batches, nil
}

// GetLastNBatchesByL2BlockNumber returns the last numBatches batches along with the l2 block state root by l2BlockNumber
// if the l2BlockNumber parameter is nil, it means we want to get the most recent last N batches
func (p *PostgresStorage) GetLastNBatchesByL2BlockNumber(ctx context.Context, l2BlockNumber *uint64, numBatches uint, dbTx pgx.Tx) ([]*state.Batch, common.Hash, error) {
	const getLastNBatchesByBlockNumberSQL = `
        SELECT b.batch_num,
               b.global_exit_root,
               b.local_exit_root,
			   b.acc_input_hash,
               b.state_root,
               b.timestamp,
               b.coinbase,
               b.raw_txs_data,
			   b.wip,
               /* gets the state root of the l2 block with the highest number associated to the batch in the row */
               (SELECT l2b1.header->>'stateRoot'
                  FROM state.l2block l2b1
                 WHERE l2b1.block_num = (SELECT MAX(l2b2.block_num)
                                           FROM state.l2block l2b2
                                          WHERE l2b2.batch_num = b.batch_num)) as l2_block_state_root
          FROM state.batch b
               /* if there is a value for the parameter $1 (l2 block number), filter the batches with batch number
                * smaller or equal than the batch associated to the l2 block number */
         WHERE ($1::int8 IS NOT NULL AND b.batch_num <= (SELECT MAX(l2b.batch_num)
                                                           FROM state.l2block l2b
                                                          WHERE l2b.block_num = $1))
               /* OR if $1 is null, this means we want to get the most updated information from state, so it considers all the batches.
                * this is generally used by estimate gas, process unsigned transactions and it is required by claim transactions to add
                * the open batch to the result and get the most updated globalExitRoot synced from L1 and stored in the current open batch when 
                * there was not transactions yet to create a l2 block with it */
            OR $1 IS NULL
         ORDER BY b.batch_num DESC
         LIMIT $2;`

	var l2BlockStateRoot *common.Hash
	e := p.getExecQuerier(dbTx)
	rows, err := e.Query(ctx, getLastNBatchesByBlockNumberSQL, l2BlockNumber, numBatches)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, common.Hash{}, state.ErrStateNotSynchronized
	} else if err != nil {
		return nil, common.Hash{}, err
	}
	defer rows.Close()

	batches := make([]*state.Batch, 0, len(rows.RawValues()))
	emptyHash := common.Hash{}

	for rows.Next() {
		batch, _l2BlockStateRoot, err := scanBatchWithL2BlockStateRoot(rows)
		if err != nil {
			return nil, common.Hash{}, err
		}
		batches = append(batches, &batch)
		if l2BlockStateRoot == nil && _l2BlockStateRoot != nil {
			l2BlockStateRoot = _l2BlockStateRoot
		}
		// if there is no corresponding l2_block, it will use the latest batch state_root
		// it is related to https://github.com/0xPolygonHermez/zkevm-node/issues/1299
		if l2BlockStateRoot == nil && batch.StateRoot != emptyHash {
			l2BlockStateRoot = &batch.StateRoot
		}
	}

	return batches, *l2BlockStateRoot, nil
}

// GetLastBatchNumber get last trusted batch number
func (p *PostgresStorage) GetLastBatchNumber(ctx context.Context, dbTx pgx.Tx) (uint64, error) {
	var batchNumber uint64
	q := p.getExecQuerier(dbTx)

	err := q.QueryRow(ctx, getLastBatchNumberSQL).Scan(&batchNumber)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, state.ErrStateNotSynchronized
	}
	return batchNumber, err
}

// GetLastBatchTime gets last trusted batch time
func (p *PostgresStorage) GetLastBatchTime(ctx context.Context, dbTx pgx.Tx) (time.Time, error) {
	var timestamp time.Time
	const getLastBatchTimeSQL = "SELECT timestamp FROM state.batch ORDER BY batch_num DESC LIMIT 1"

	e := p.getExecQuerier(dbTx)
	err := e.QueryRow(ctx, getLastBatchTimeSQL).Scan(&timestamp)

	if errors.Is(err, pgx.ErrNoRows) {
		return time.Time{}, state.ErrStateNotSynchronized
	} else if err != nil {
		return time.Time{}, err
	}
	return timestamp, nil
}

// GetLastVirtualBatchNum gets last virtual batch num
func (p *PostgresStorage) GetLastVirtualBatchNum(ctx context.Context, dbTx pgx.Tx) (uint64, error) {
	var batchNum uint64
	const getLastVirtualBatchNumSQL = "SELECT COALESCE(MAX(batch_num), 0) FROM state.virtual_batch"

	e := p.getExecQuerier(dbTx)
	err := e.QueryRow(ctx, getLastVirtualBatchNumSQL).Scan(&batchNum)

	if errors.Is(err, pgx.ErrNoRows) {
		return 0, state.ErrNotFound
	} else if err != nil {
		return 0, err
	}
	return batchNum, nil
}

// GetLatestVirtualBatchTimestamp gets last virtual batch timestamp
func (p *PostgresStorage) GetLatestVirtualBatchTimestamp(ctx context.Context, dbTx pgx.Tx) (time.Time, error) {
	const getLastVirtualBatchTimestampSQL = `SELECT COALESCE(MAX(block.received_at), NOW()) FROM state.virtual_batch INNER JOIN state.block ON state.block.block_num = virtual_batch.block_num`
	var timestamp time.Time
	e := p.getExecQuerier(dbTx)
	err := e.QueryRow(ctx, getLastVirtualBatchTimestampSQL).Scan(&timestamp)

	if errors.Is(err, pgx.ErrNoRows) {
		return time.Unix(0, 0), state.ErrNotFound
	} else if err != nil {
		return time.Unix(0, 0), err
	}
	return timestamp, nil
}

// SetLastBatchInfoSeenOnEthereum sets the last batch number that affected
// the roll-up and the last batch number that was consolidated on ethereum
// in order to allow the components to know if the state is synchronized or not
func (p *PostgresStorage) SetLastBatchInfoSeenOnEthereum(ctx context.Context, lastBatchNumberSeen, lastBatchNumberVerified uint64, dbTx pgx.Tx) error {
	const query = `
    UPDATE state.sync_info
       SET last_batch_num_seen = $1, last_batch_num_consolidated = $2`

	e := p.getExecQuerier(dbTx)
	_, err := e.Exec(ctx, query, lastBatchNumberSeen, lastBatchNumberVerified)
	return err
}

// SetInitSyncBatch sets the initial batch number where the synchronization started
func (p *PostgresStorage) SetInitSyncBatch(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) error {
	updateInitBatchSQL := "UPDATE state.sync_info SET init_sync_batch = $1"

	e := p.getExecQuerier(dbTx)
	_, err := e.Exec(ctx, updateInitBatchSQL, batchNumber)
	return err
}

// GetBatchByNumber returns the batch with the given number.
func (p *PostgresStorage) GetBatchByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*state.Batch, error) {
	const getBatchByNumberSQL = `
		SELECT batch_num, global_exit_root, local_exit_root, acc_input_hash, state_root, timestamp, coinbase, raw_txs_data, forced_batch_num, batch_resources, high_reserved_counters, wip
		  FROM state.batch 
		 WHERE batch_num = $1`

	e := p.getExecQuerier(dbTx)
	row := e.QueryRow(ctx, getBatchByNumberSQL, batchNumber)
	batch, err := scanBatch(row)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return &batch, nil
}

// GetBatchByTxHash returns the batch including the given tx
func (p *PostgresStorage) GetBatchByTxHash(ctx context.Context, transactionHash common.Hash, dbTx pgx.Tx) (*state.Batch, error) {
	const getBatchByTxHashSQL = `
		SELECT b.batch_num, b.global_exit_root, b.local_exit_root, b.acc_input_hash, b.state_root, b.timestamp, b.coinbase, b.raw_txs_data, b.forced_batch_num, b.batch_resources, b.high_reserved_counters, b.wip
		  FROM state.transaction t, state.batch b, state.l2block l 
		  WHERE t.hash = $1 AND l.block_num = t.l2_block_num AND b.batch_num = l.batch_num`

	e := p.getExecQuerier(dbTx)
	row := e.QueryRow(ctx, getBatchByTxHashSQL, transactionHash.String())
	batch, err := scanBatch(row)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrStateNotSynchronized
	} else if err != nil {
		return nil, err
	}
	return &batch, nil
}

// GetBatchByL2BlockNumber returns the batch related to the l2 block accordingly to the provided l2 block number.
func (p *PostgresStorage) GetBatchByL2BlockNumber(ctx context.Context, l2BlockNumber uint64, dbTx pgx.Tx) (*state.Batch, error) {
	const getBatchByL2BlockNumberSQL = `
		SELECT bt.batch_num, bt.global_exit_root, bt.local_exit_root, bt.acc_input_hash, bt.state_root, bt.timestamp, bt.coinbase, bt.raw_txs_data, bt.forced_batch_num, bt.batch_resources, bt.high_reserved_counters, bt.wip
		  FROM state.batch bt
		 INNER JOIN state.l2block bl
		    ON bt.batch_num = bl.batch_num
		 WHERE bl.block_num = $1
		 LIMIT 1;`

	e := p.getExecQuerier(dbTx)
	row := e.QueryRow(ctx, getBatchByL2BlockNumberSQL, l2BlockNumber)
	batch, err := scanBatch(row)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrStateNotSynchronized
	} else if err != nil {
		return nil, err
	}
	return &batch, nil
}

// GetVirtualBatchByNumber gets batch from batch table that exists on virtual batch
func (p *PostgresStorage) GetVirtualBatchByNumber(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*state.Batch, error) {
	const query = `
		SELECT
			batch_num,
			global_exit_root,
			local_exit_root,
			acc_input_hash,
			state_root,
			timestamp,
			coinbase,
			raw_txs_data,
			forced_batch_num,
			batch_resources, 
			high_reserved_counters,
			wip
		FROM
			state.batch
		WHERE
			batch_num = $1 AND
			EXISTS (SELECT batch_num FROM state.virtual_batch WHERE batch_num = $1)
		`
	e := p.getExecQuerier(dbTx)
	row := e.QueryRow(ctx, query, batchNumber)
	batch, err := scanBatch(row)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return &batch, nil
}

// IsBatchVirtualized checks if batch is virtualized
func (p *PostgresStorage) IsBatchVirtualized(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (bool, error) {
	const query = `SELECT EXISTS (SELECT 1 FROM state.virtual_batch WHERE batch_num = $1)`
	e := p.getExecQuerier(dbTx)
	var exists bool
	err := e.QueryRow(ctx, query, batchNumber).Scan(&exists)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return exists, err
	}
	return exists, nil
}

// IsBatchConsolidated checks if batch is consolidated/verified.
func (p *PostgresStorage) IsBatchConsolidated(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (bool, error) {
	const query = `SELECT EXISTS (SELECT 1 FROM state.verified_batch WHERE batch_num = $1)`
	e := p.getExecQuerier(dbTx)
	var exists bool
	err := e.QueryRow(ctx, query, batchNumber).Scan(&exists)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return exists, err
	}
	return exists, nil
}

// IsSequencingTXSynced checks if sequencing tx has been synced into the state
func (p *PostgresStorage) IsSequencingTXSynced(ctx context.Context, transactionHash common.Hash, dbTx pgx.Tx) (bool, error) {
	const query = `SELECT EXISTS (SELECT 1 FROM state.virtual_batch WHERE tx_hash = $1)`
	e := p.getExecQuerier(dbTx)
	var exists bool
	err := e.QueryRow(ctx, query, transactionHash.String()).Scan(&exists)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return exists, err
	}
	return exists, nil
}
func scanBatch(row pgx.Row) (state.Batch, error) {
	batch := state.Batch{}
	var (
		gerStr               string
		lerStr               *string
		aihStr               *string
		stateStr             *string
		coinbaseStr          string
		resourcesData        []byte
		highReservedCounters []byte
		wip                  bool
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
		&resourcesData,
		&highReservedCounters,
		&wip,
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

	if resourcesData != nil {
		err = json.Unmarshal(resourcesData, &batch.Resources)
		if err != nil {
			return batch, err
		}
	}

	if highReservedCounters != nil {
		err = json.Unmarshal(highReservedCounters, &batch.HighReservedZKCounters)
		if err != nil {
			return batch, err
		}
	}

	batch.WIP = wip

	batch.Coinbase = common.HexToAddress(coinbaseStr)
	return batch, nil
}

func scanBatchWithL2BlockStateRoot(row pgx.Row) (state.Batch, *common.Hash, error) {
	batch := state.Batch{}
	var (
		gerStr              string
		lerStr              *string
		aihStr              *string
		stateStr            *string
		coinbaseStr         string
		l2BlockStateRootStr *string
		wip                 bool
	)
	if err := row.Scan(
		&batch.BatchNumber,
		&gerStr,
		&lerStr,
		&aihStr,
		&stateStr,
		&batch.Timestamp,
		&coinbaseStr,
		&batch.BatchL2Data,
		&wip,
		&l2BlockStateRootStr,
	); err != nil {
		return batch, nil, err
	}
	batch.GlobalExitRoot = common.HexToHash(gerStr)
	if lerStr != nil {
		batch.LocalExitRoot = common.HexToHash(*lerStr)
	}
	if stateStr != nil {
		batch.StateRoot = common.HexToHash(*stateStr)
	}
	if stateStr != nil {
		batch.AccInputHash = common.HexToHash(*aihStr)
	}
	var l2BlockStateRoot *common.Hash
	if l2BlockStateRootStr != nil {
		h := common.HexToHash(*l2BlockStateRootStr)
		l2BlockStateRoot = &h
	}
	batch.WIP = wip
	batch.Coinbase = common.HexToAddress(coinbaseStr)
	return batch, l2BlockStateRoot, nil
}

func scanForcedBatch(row pgx.Row) (state.ForcedBatch, error) {
	forcedBatch := state.ForcedBatch{}
	var (
		gerStr      string
		coinbaseStr string
		rawTxsStr   string
		err         error
	)
	if err := row.Scan(
		&forcedBatch.ForcedBatchNumber,
		&gerStr,
		&forcedBatch.ForcedAt,
		&rawTxsStr,
		&coinbaseStr,
		&forcedBatch.BlockNumber,
	); err != nil {
		return forcedBatch, err
	}
	forcedBatch.RawTxsData, err = hex.DecodeString(rawTxsStr)
	if err != nil {
		return forcedBatch, err
	}
	forcedBatch.GlobalExitRoot = common.HexToHash(gerStr)
	forcedBatch.Sequencer = common.HexToAddress(coinbaseStr)
	return forcedBatch, nil
}

// AddVirtualBatch adds a new virtual batch to the storage.
func (p *PostgresStorage) AddVirtualBatch(ctx context.Context, virtualBatch *state.VirtualBatch, dbTx pgx.Tx) error {
	if virtualBatch.TimestampBatchEtrog == nil {
		const addVirtualBatchSQL = "INSERT INTO state.virtual_batch (batch_num, tx_hash, coinbase, block_num, sequencer_addr) VALUES ($1, $2, $3, $4, $5)"
		e := p.getExecQuerier(dbTx)
		_, err := e.Exec(ctx, addVirtualBatchSQL, virtualBatch.BatchNumber, virtualBatch.TxHash.String(), virtualBatch.Coinbase.String(), virtualBatch.BlockNumber, virtualBatch.SequencerAddr.String())
		return err
	} else {
		var l1InfoRoot *string
		if virtualBatch.L1InfoRoot != nil {
			l1IR := virtualBatch.L1InfoRoot.String()
			l1InfoRoot = &l1IR
		}
		const addVirtualBatchSQL = "INSERT INTO state.virtual_batch (batch_num, tx_hash, coinbase, block_num, sequencer_addr, timestamp_batch_etrog, l1_info_root) VALUES ($1, $2, $3, $4, $5, $6, $7)"
		e := p.getExecQuerier(dbTx)
		_, err := e.Exec(ctx, addVirtualBatchSQL, virtualBatch.BatchNumber, virtualBatch.TxHash.String(), virtualBatch.Coinbase.String(), virtualBatch.BlockNumber, virtualBatch.SequencerAddr.String(),
			virtualBatch.TimestampBatchEtrog.UTC(), l1InfoRoot)
		return err
	}
}

// GetVirtualBatch get an L1 virtualBatch.
func (p *PostgresStorage) GetVirtualBatch(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*state.VirtualBatch, error) {
	var (
		virtualBatch  state.VirtualBatch
		txHash        string
		coinbase      string
		sequencerAddr string
		l1InfoRoot    *string
	)

	const getVirtualBatchSQL = `
    SELECT block_num, batch_num, tx_hash, coinbase, sequencer_addr, timestamp_batch_etrog, l1_info_root
      FROM state.virtual_batch
     WHERE batch_num = $1`

	e := p.getExecQuerier(dbTx)
	err := e.QueryRow(ctx, getVirtualBatchSQL, batchNumber).Scan(&virtualBatch.BlockNumber, &virtualBatch.BatchNumber, &txHash, &coinbase, &sequencerAddr, &virtualBatch.TimestampBatchEtrog, &l1InfoRoot)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrNotFound
	} else if err != nil {
		return nil, err
	}
	virtualBatch.Coinbase = common.HexToAddress(coinbase)
	virtualBatch.SequencerAddr = common.HexToAddress(sequencerAddr)
	virtualBatch.TxHash = common.HexToHash(txHash)
	if l1InfoRoot != nil {
		l1InfoR := common.HexToHash(*l1InfoRoot)
		virtualBatch.L1InfoRoot = &l1InfoR
	}
	return &virtualBatch, nil
}

func (p *PostgresStorage) StoreGenesisBatch(ctx context.Context, batch state.Batch, closingReason string, dbTx pgx.Tx) error {
	const addGenesisBatchSQL = "INSERT INTO state.batch (batch_num, global_exit_root, local_exit_root, acc_input_hash, state_root, timestamp, coinbase, raw_txs_data, forced_batch_num,closing_reason, wip) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9,$10, FALSE)"

	if batch.BatchNumber != 0 {
		return fmt.Errorf("%w. Got %d, should be 0", state.ErrUnexpectedBatch, batch.BatchNumber)
	}
	e := p.getExecQuerier(dbTx)
	_, err := e.Exec(
		ctx,
		addGenesisBatchSQL,
		batch.BatchNumber,
		batch.GlobalExitRoot.String(),
		batch.LocalExitRoot.String(),
		batch.AccInputHash.String(),
		batch.StateRoot.String(),
		batch.Timestamp.UTC(),
		batch.Coinbase.String(),
		batch.BatchL2Data,
		batch.ForcedBatchNum,
		closingReason,
	)

	return err
}

// OpenBatchInStorage adds a new batch into the state storage, with the necessary data to start processing transactions within it.
// It's meant to be used by sequencers, since they don't necessarily know what transactions are going to be added
// in this batch yet. In other words it's the creation of a WIP batch.
// Note that this will add a batch with batch number N + 1, where N it's the greatest batch number on the state.
func (p *PostgresStorage) OpenBatchInStorage(ctx context.Context, batchContext state.ProcessingContext, dbTx pgx.Tx) error {
	const openBatchSQL = "INSERT INTO state.batch (batch_num, global_exit_root, timestamp, coinbase, forced_batch_num, raw_txs_data, wip) VALUES ($1, $2, $3, $4, $5, $6, TRUE)"

	e := p.getExecQuerier(dbTx)
	_, err := e.Exec(
		ctx, openBatchSQL,
		batchContext.BatchNumber,
		batchContext.GlobalExitRoot.String(),
		batchContext.Timestamp.UTC(),
		batchContext.Coinbase.String(),
		batchContext.ForcedBatchNum,
		batchContext.BatchL2Data,
	)
	return err
}

// OpenWIPBatchInStorage adds a new wip batch into the state storage
func (p *PostgresStorage) OpenWIPBatchInStorage(ctx context.Context, batch state.Batch, dbTx pgx.Tx) error {
	const openBatchSQL = "INSERT INTO state.batch (batch_num, global_exit_root, state_root, local_exit_root, timestamp, coinbase, forced_batch_num, raw_txs_data, batch_resources, wip, checked) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, TRUE, FALSE)"

	resourcesData, err := json.Marshal(batch.Resources)
	if err != nil {
		return err
	}
	resources := string(resourcesData)

	e := p.getExecQuerier(dbTx)
	_, err = e.Exec(
		ctx, openBatchSQL,
		batch.BatchNumber,
		batch.GlobalExitRoot.String(),
		batch.StateRoot.String(),
		batch.LocalExitRoot.String(),
		batch.Timestamp.UTC(),
		batch.Coinbase.String(),
		batch.ForcedBatchNum,
		batch.BatchL2Data,
		resources,
	)
	return err
}

// CloseBatchInStorage closes a batch in the state storage
func (p *PostgresStorage) CloseBatchInStorage(ctx context.Context, receipt state.ProcessingReceipt, dbTx pgx.Tx) error {
	const closeBatchSQL = `UPDATE state.batch 
		SET state_root = $1, local_exit_root = $2, acc_input_hash = $3, raw_txs_data = $4, batch_resources = $5, closing_reason = $6, wip = FALSE
		  WHERE batch_num = $7`

	e := p.getExecQuerier(dbTx)
	batchResourcesJsonBytes, err := json.Marshal(receipt.BatchResources)
	if err != nil {
		return err
	}
	_, err = e.Exec(ctx, closeBatchSQL, receipt.StateRoot.String(), receipt.LocalExitRoot.String(),
		receipt.AccInputHash.String(), receipt.BatchL2Data, string(batchResourcesJsonBytes), receipt.ClosingReason, receipt.BatchNumber)

	return err
}

// CloseWIPBatchInStorage is used by sequencer to close the wip batch in the state storage
func (p *PostgresStorage) CloseWIPBatchInStorage(ctx context.Context, receipt state.ProcessingReceipt, dbTx pgx.Tx) error {
	const closeWIPBatchSQL = `UPDATE state.batch SET batch_resources = $1, closing_reason = $2, wip = FALSE WHERE batch_num = $3`

	e := p.getExecQuerier(dbTx)
	batchResourcesJsonBytes, err := json.Marshal(receipt.BatchResources)
	if err != nil {
		return err
	}
	_, err = e.Exec(ctx, closeWIPBatchSQL, string(batchResourcesJsonBytes), receipt.ClosingReason, receipt.BatchNumber)

	return err
}

// GetWIPBatchInStorage returns the wip batch in the state
func (p *PostgresStorage) GetWIPBatchInStorage(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*state.Batch, error) {
	const getWIPBatchByNumberSQL = `
		SELECT batch_num, global_exit_root, local_exit_root, acc_input_hash, state_root, timestamp, coinbase, raw_txs_data, forced_batch_num, batch_resources, high_reserved_counters, wip
		  FROM state.batch 
		 WHERE batch_num = $1 AND wip = TRUE`

	e := p.getExecQuerier(dbTx)
	row := e.QueryRow(ctx, getWIPBatchByNumberSQL, batchNumber)
	batch, err := scanBatch(row)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return &batch, nil
}

// IsBatchClosed indicates if the batch referenced by batchNum is closed or not
func (p *PostgresStorage) IsBatchClosed(ctx context.Context, batchNum uint64, dbTx pgx.Tx) (bool, error) {
	const isBatchClosedSQL = "SELECT not(wip) FROM state.batch WHERE batch_num = $1"

	q := p.getExecQuerier(dbTx)
	var isClosed bool
	err := q.QueryRow(ctx, isBatchClosedSQL, batchNum).Scan(&isClosed)
	return isClosed, err
}

// GetBatchNumberOfL2Block gets a batch number for l2 block by its number
func (p *PostgresStorage) GetBatchNumberOfL2Block(ctx context.Context, blockNumber uint64, dbTx pgx.Tx) (uint64, error) {
	getBatchNumByBlockNum := "SELECT batch_num FROM state.l2block WHERE block_num = $1"
	batchNumber := uint64(0)
	q := p.getExecQuerier(dbTx)
	err := q.QueryRow(ctx, getBatchNumByBlockNum, blockNumber).
		Scan(&batchNumber)

	if errors.Is(err, pgx.ErrNoRows) {
		return batchNumber, state.ErrNotFound
	} else if err != nil {
		return batchNumber, err
	}
	return batchNumber, nil
}

// BatchNumberByL2BlockNumber gets a batch number by a l2 block number
func (p *PostgresStorage) BatchNumberByL2BlockNumber(ctx context.Context, blockNumber uint64, dbTx pgx.Tx) (uint64, error) {
	getBatchNumByBlockNum := "SELECT batch_num FROM state.l2block WHERE block_num = $1"
	batchNumber := uint64(0)
	q := p.getExecQuerier(dbTx)
	err := q.QueryRow(ctx, getBatchNumByBlockNum, blockNumber).
		Scan(&batchNumber)

	if errors.Is(err, pgx.ErrNoRows) {
		return batchNumber, state.ErrNotFound
	} else if err != nil {
		return batchNumber, err
	}
	return batchNumber, nil
}

// GetLastVerifiedBatchNumberUntilL1Block gets the last batch number that was verified in
// or before the provided l1 block number. This is used to identify if a batch is safe or finalized.
func (p *PostgresStorage) GetLastVerifiedBatchNumberUntilL1Block(ctx context.Context, l1BlockNumber uint64, dbTx pgx.Tx) (uint64, error) {
	var batchNumber uint64
	const query = `
    SELECT vb.batch_num
      FROM state.verified_batch vb
	 WHERE vb.block_num <= $1
     ORDER BY vb.batch_num DESC LIMIT 1`

	q := p.getExecQuerier(dbTx)
	err := q.QueryRow(ctx, query, l1BlockNumber).Scan(&batchNumber)

	if errors.Is(err, pgx.ErrNoRows) {
		return 0, state.ErrNotFound
	} else if err != nil {
		return 0, err
	}

	return batchNumber, nil
}

// GetLastVerifiedBatch gets last verified batch
func (p *PostgresStorage) GetLastVerifiedBatch(ctx context.Context, dbTx pgx.Tx) (*state.VerifiedBatch, error) {
	const query = "SELECT block_num, batch_num, tx_hash, aggregator FROM state.verified_batch ORDER BY batch_num DESC LIMIT 1"
	var (
		verifiedBatch state.VerifiedBatch
		txHash, agg   string
	)
	e := p.getExecQuerier(dbTx)
	err := e.QueryRow(ctx, query).Scan(&verifiedBatch.BlockNumber, &verifiedBatch.BatchNumber, &txHash, &agg)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrNotFound
	} else if err != nil {
		return nil, err
	}
	verifiedBatch.Aggregator = common.HexToAddress(agg)
	verifiedBatch.TxHash = common.HexToHash(txHash)
	return &verifiedBatch, nil
}

// GetVirtualBatchToProve return the next batch that is not proved, neither in
// proved process.
func (p *PostgresStorage) GetVirtualBatchToProve(ctx context.Context, lastVerfiedBatchNumber uint64, maxL1Block uint64, dbTx pgx.Tx) (*state.Batch, error) {
	const query = `
		SELECT
			b.batch_num,
			b.global_exit_root,
			b.local_exit_root,
			b.acc_input_hash,
			b.state_root,
			v.timestamp_batch_etrog,
			b.coinbase,
			b.raw_txs_data,
			b.forced_batch_num,
			b.batch_resources, 
			b.high_reserved_counters,
			b.wip
		FROM
			state.batch b,
			state.virtual_batch v
		WHERE
			b.batch_num > $1 AND b.batch_num = v.batch_num AND
			v.block_num <= $2 AND
			NOT EXISTS (
				SELECT p.batch_num FROM state.proof p 
				WHERE v.batch_num >= p.batch_num AND v.batch_num <= p.batch_num_final
			)
		ORDER BY b.batch_num ASC LIMIT 1
		`
	e := p.getExecQuerier(dbTx)
	row := e.QueryRow(ctx, query, lastVerfiedBatchNumber, maxL1Block)
	batch, err := scanBatch(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return &batch, nil
}

// AddSequence stores the sequence information to allow the aggregator verify sequences.
func (p *PostgresStorage) AddSequence(ctx context.Context, sequence state.Sequence, dbTx pgx.Tx) error {
	const addSequenceSQL = "INSERT INTO state.sequences (from_batch_num, to_batch_num) VALUES($1, $2) ON CONFLICT (from_batch_num) DO UPDATE SET to_batch_num = $2"

	e := p.getExecQuerier(dbTx)
	_, err := e.Exec(ctx, addSequenceSQL, sequence.FromBatchNumber, sequence.ToBatchNumber)
	return err
}

// GetSequences get the next sequences higher than an specify batch number
func (p *PostgresStorage) GetSequences(ctx context.Context, lastVerifiedBatchNumber uint64, dbTx pgx.Tx) ([]state.Sequence, error) {
	const getSequencesSQL = "SELECT from_batch_num, to_batch_num FROM state.sequences WHERE from_batch_num >= $1 ORDER BY from_batch_num ASC"
	q := p.getExecQuerier(dbTx)

	rows, err := q.Query(ctx, getSequencesSQL, lastVerifiedBatchNumber)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrStateNotSynchronized
	} else if err != nil {
		return nil, err
	}
	defer rows.Close()

	sequences := make([]state.Sequence, 0, len(rows.RawValues()))

	for rows.Next() {
		var sequence state.Sequence
		if err := rows.Scan(
			&sequence.FromBatchNumber,
			&sequence.ToBatchNumber,
		); err != nil {
			return sequences, err
		}
		sequences = append(sequences, sequence)
	}
	return sequences, err
}

// GetLastClosedBatch returns the latest closed batch
func (p *PostgresStorage) GetLastClosedBatch(ctx context.Context, dbTx pgx.Tx) (*state.Batch, error) {
	const getLastClosedBatchSQL = `
		SELECT bt.batch_num, bt.global_exit_root, bt.local_exit_root, bt.acc_input_hash, bt.state_root, bt.timestamp, bt.coinbase, bt.raw_txs_data, bt.forced_batch_num, bt.batch_resources, bt.high_reserved_counters, bt.wip
			FROM state.batch bt
			WHERE wip = FALSE
			ORDER BY bt.batch_num DESC
			LIMIT 1;`

	e := p.getExecQuerier(dbTx)
	row := e.QueryRow(ctx, getLastClosedBatchSQL)
	batch, err := scanBatch(row)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrStateNotSynchronized
	} else if err != nil {
		return nil, err
	}
	return &batch, nil
}

// GetLastClosedBatchNumber returns the latest closed batch
func (p *PostgresStorage) GetLastClosedBatchNumber(ctx context.Context, dbTx pgx.Tx) (uint64, error) {
	const getLastClosedBatchSQL = `
		SELECT bt.batch_num
			FROM state.batch bt
			WHERE wip = FALSE
			ORDER BY bt.batch_num DESC
			LIMIT 1;`

	batchNumber := uint64(0)
	e := p.getExecQuerier(dbTx)
	err := e.QueryRow(ctx, getLastClosedBatchSQL).Scan(&batchNumber)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, state.ErrStateNotSynchronized
	} else if err != nil {
		return 0, err
	}
	return batchNumber, nil
}

// UpdateBatchL2Data updates data tx data in a batch
func (p *PostgresStorage) UpdateBatchL2Data(ctx context.Context, batchNumber uint64, batchL2Data []byte, dbTx pgx.Tx) error {
	const updateL2DataSQL = "UPDATE state.batch SET raw_txs_data = $2 WHERE batch_num = $1"

	e := p.getExecQuerier(dbTx)
	_, err := e.Exec(ctx, updateL2DataSQL, batchNumber, batchL2Data)
	return err
}

// UpdateWIPBatch updates the data in a batch
func (p *PostgresStorage) UpdateWIPBatch(ctx context.Context, receipt state.ProcessingReceipt, dbTx pgx.Tx) error {
	const updateL2DataSQL = "UPDATE state.batch SET raw_txs_data = $2, global_exit_root = $3, state_root = $4, local_exit_root = $5, batch_resources = $6, high_reserved_counters = $7 WHERE batch_num = $1"

	e := p.getExecQuerier(dbTx)
	batchResourcesJsonBytes, err := json.Marshal(receipt.BatchResources)
	if err != nil {
		return err
	}

	highReservedCounters, err := json.Marshal(receipt.HighReservedZKCounters)
	if err != nil {
		return err
	}

	_, err = e.Exec(ctx, updateL2DataSQL, receipt.BatchNumber, receipt.BatchL2Data, receipt.GlobalExitRoot.String(), receipt.StateRoot.String(), receipt.LocalExitRoot.String(), string(batchResourcesJsonBytes), string(highReservedCounters))
	return err
}

// updateBatchAsChecked updates the batch to set it as checked (sequencer sanity check was successful)
func (p *PostgresStorage) UpdateBatchAsChecked(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) error {
	const updateL2DataSQL = "UPDATE state.batch SET checked = TRUE WHERE batch_num = $1"

	e := p.getExecQuerier(dbTx)
	_, err := e.Exec(ctx, updateL2DataSQL, batchNumber)
	return err
}

// IsBatchChecked indicates if the batch is closed and checked (sequencer sanity check was successful)
func (p *PostgresStorage) IsBatchChecked(ctx context.Context, batchNum uint64, dbTx pgx.Tx) (bool, error) {
	const isBatchCheckedSQL = "SELECT not(wip) AND checked FROM state.batch WHERE batch_num = $1"

	q := p.getExecQuerier(dbTx)
	var isChecked bool
	err := q.QueryRow(ctx, isBatchCheckedSQL, batchNum).Scan(&isChecked)
	return isChecked, err
}

// AddAccumulatedInputHash adds the accumulated input hash
func (p *PostgresStorage) AddAccumulatedInputHash(ctx context.Context, batchNum uint64, accInputHash common.Hash, dbTx pgx.Tx) error {
	const addAccInputHashBatchSQL = "UPDATE state.batch SET acc_input_hash = $1 WHERE batch_num = $2"
	e := p.getExecQuerier(dbTx)
	_, err := e.Exec(ctx, addAccInputHashBatchSQL, accInputHash.String(), batchNum)
	return err
}

// GetLocalExitRootByBatchNumber get local exit root by batch number
func (p *PostgresStorage) GetLocalExitRootByBatchNumber(ctx context.Context, batchNum uint64, dbTx pgx.Tx) (common.Hash, error) {
	const query = "SELECT local_exit_root FROM state.batch WHERE batch_num = $1"
	var localExitRootStr string
	e := p.getExecQuerier(dbTx)
	err := e.QueryRow(ctx, query, batchNum).Scan(&localExitRootStr)
	if errors.Is(err, pgx.ErrNoRows) {
		return common.Hash{}, state.ErrNotFound
	} else if err != nil {
		return common.Hash{}, err
	}
	return common.HexToHash(localExitRootStr), nil
}

// GetBlockNumVirtualBatchByBatchNum get block num of virtual batch by block num
func (p *PostgresStorage) GetBlockNumVirtualBatchByBatchNum(ctx context.Context, batchNum uint64, dbTx pgx.Tx) (uint64, error) {
	const query = "SELECT block_num FROM state.virtual_batch WHERE batch_num = $1"
	var blockNum uint64
	e := p.getExecQuerier(dbTx)
	err := e.QueryRow(ctx, query, batchNum).Scan(&blockNum)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, state.ErrNotFound
	} else if err != nil {
		return 0, err
	}
	return blockNum, nil
}

// BuildChangeL2Block returns a changeL2Block tx to use in the BatchL2Data
func (p *PostgresStorage) BuildChangeL2Block(deltaTimestamp uint32, l1InfoTreeIndex uint32) []byte {
	changeL2BlockMark := []byte{0x0B}
	changeL2Block := []byte{}

	// changeL2Block transaction mark
	changeL2Block = append(changeL2Block, changeL2BlockMark...)
	// changeL2Block deltaTimeStamp
	deltaTimestampBytes := make([]byte, 4) //nolint:gomnd
	binary.BigEndian.PutUint32(deltaTimestampBytes, deltaTimestamp)
	changeL2Block = append(changeL2Block, deltaTimestampBytes...)
	// changeL2Block l1InfoTreeIndexBytes
	l1InfoTreeIndexBytes := make([]byte, 4) //nolint:gomnd
	binary.BigEndian.PutUint32(l1InfoTreeIndexBytes, l1InfoTreeIndex)
	changeL2Block = append(changeL2Block, l1InfoTreeIndexBytes...)

	return changeL2Block
}

// GetRawBatchTimestamps returns the timestamp of the batch with the given number.
// it returns batch_num.tstamp and virtual_batch.batch_timestamp
func (p *PostgresStorage) GetRawBatchTimestamps(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (*time.Time, *time.Time, error) {
	const sql = `
	SELECT b.timestamp AS batch_timestamp, v.timestamp_batch_etrog AS virtual_batch_timestamp
		FROM state.batch AS b
		LEFT JOIN state.virtual_batch AS v ON b.batch_num = v.batch_num
		WHERE b.batch_num = $1;
	`
	var batchTimestamp, virtualBatchTimestamp *time.Time
	e := p.getExecQuerier(dbTx)
	err := e.QueryRow(ctx, sql, batchNumber).Scan(&batchTimestamp, &virtualBatchTimestamp)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil, nil
	}
	return batchTimestamp, virtualBatchTimestamp, err
}

// GetVirtualBatchParentHash returns the parent hash of the virtual batch with the given number.
func (p *PostgresStorage) GetVirtualBatchParentHash(ctx context.Context, batchNumber uint64, dbTx pgx.Tx) (common.Hash, error) {
	var parentHash string

	const sql = `SELECT b.parent_hash FROM state.virtual_batch v, state.block b
     WHERE v.batch_num = $1 and b.block_num = v.block_num`

	e := p.getExecQuerier(dbTx)
	err := e.QueryRow(ctx, sql, batchNumber).Scan(&parentHash)
	if errors.Is(err, pgx.ErrNoRows) {
		return common.Hash{}, state.ErrNotFound
	} else if err != nil {
		return common.Hash{}, err
	}
	return common.HexToHash(parentHash), nil
}

// GetForcedBatchParentHash returns the parent hash of the forced batch with the given number and the globalExitRoot.
func (p *PostgresStorage) GetForcedBatchParentHash(ctx context.Context, forcedBatchNumber uint64, dbTx pgx.Tx) (common.Hash, error) {
	var (
		parentHash string
	)

	const sql = `SELECT b.parent_hash FROM state.forced_batch f, state.block b
     WHERE f.forced_batch_num = $1 and b.block_num = f.block_num`

	e := p.getExecQuerier(dbTx)
	err := e.QueryRow(ctx, sql, forcedBatchNumber).Scan(&parentHash)
	if errors.Is(err, pgx.ErrNoRows) {
		return common.Hash{}, state.ErrNotFound
	} else if err != nil {
		return common.Hash{}, err
	}
	return common.HexToHash(parentHash), nil
}

// GetLatestBatchGlobalExitRoot gets the last GER that is not zero from batches
func (p *PostgresStorage) GetLatestBatchGlobalExitRoot(ctx context.Context, dbTx pgx.Tx) (common.Hash, error) {
	var lastGER string
	const query = "SELECT global_exit_root FROM state.batch where global_exit_root != $1 ORDER BY batch_num DESC LIMIT 1"

	q := p.getExecQuerier(dbTx)
	err := q.QueryRow(ctx, query, state.ZeroHash.String()).Scan(&lastGER)

	if errors.Is(err, pgx.ErrNoRows) {
		return state.ZeroHash, nil
	} else if err != nil {
		return state.ZeroHash, err
	}

	return common.HexToHash(lastGER), nil
}

// GetNotCheckedBatches returns the batches that are closed but not checked
func (p *PostgresStorage) GetNotCheckedBatches(ctx context.Context, dbTx pgx.Tx) ([]*state.Batch, error) {
	const getBatchesNotCheckedSQL = `
		SELECT batch_num, global_exit_root, local_exit_root, acc_input_hash, state_root, timestamp, coinbase, raw_txs_data, forced_batch_num, batch_resources, high_reserved_counters, wip 
		from state.batch WHERE wip IS FALSE AND checked IS FALSE ORDER BY batch_num ASC`

	e := p.getExecQuerier(dbTx)
	rows, err := e.Query(ctx, getBatchesNotCheckedSQL)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, state.ErrNotFound
	} else if err != nil {
		return nil, err
	}
	defer rows.Close()

	batches := make([]*state.Batch, 0, len(rows.RawValues()))

	for rows.Next() {
		batch, err := scanBatch(rows)
		if err != nil {
			return nil, err
		}
		batches = append(batches, &batch)
	}

	return batches, nil
}

// UpdateBatchTimestamp updates the timestamp of the state.batch with the given number.
func (p *PostgresStorage) UpdateBatchTimestamp(ctx context.Context, batchNumber uint64, timestamp time.Time, dbTx pgx.Tx) error {
	const sql = "UPDATE state.batch SET timestamp = $1 WHERE batch_num = $2"
	e := p.getExecQuerier(dbTx)
	_, err := e.Exec(ctx, sql, timestamp.UTC(), batchNumber)
	return err
}
