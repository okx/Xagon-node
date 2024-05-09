package pgpoolstorage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/db"
	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// PostgresPoolStorage is an implementation of the Pool interface
// that uses a postgres database to store the data
type PostgresPoolStorage struct {
	db *pgxpool.Pool
}

// NewPostgresPoolStorage creates and initializes an instance of PostgresPoolStorage
func NewPostgresPoolStorage(cfg db.Config) (*PostgresPoolStorage, error) {
	poolDB, err := db.NewSQLDB(cfg)
	if err != nil {
		return nil, err
	}

	return &PostgresPoolStorage{
		db: poolDB,
	}, nil
}

// AddTx adds a transaction to the pool table with the provided status
func (p *PostgresPoolStorage) AddTx(ctx context.Context, tx pool.Transaction) error {
	hash := tx.Hash().Hex()

	b, err := tx.MarshalBinary()
	if err != nil {
		return err
	}
	encoded := hex.EncodeToHex(b)

	b, err = tx.MarshalJSON()
	if err != nil {
		return err
	}
	decoded := string(b)

	gasPrice := tx.GasPrice().Uint64()
	nonce := tx.Nonce()

	sql := `
		INSERT INTO pool.transaction 
		(
			hash,
			encoded,
			decoded,
			status,
			gas_price,
			nonce,
			cumulative_gas_used,
			used_keccak_hashes,
			used_poseidon_hashes,
			used_poseidon_paddings,
			used_mem_aligns,
			used_arithmetics,
			used_binaries,
			used_steps,
			used_sha256_hashes,
			received_at,
			from_address,
			is_wip,
			ip,
			failed_reason,
			reserved_zkcounters
		) 
		VALUES 
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, NULL, $20)
			ON CONFLICT (hash) DO UPDATE SET 
			encoded = $2,
			decoded = $3,
			status = $4,
			gas_price = $5,
			nonce = $6,
			cumulative_gas_used = $7,
			used_keccak_hashes = $8, 
			used_poseidon_hashes = $9,
			used_poseidon_paddings = $10,
			used_mem_aligns = $11,
			used_arithmetics = $12,
			used_binaries = $13,
			used_steps = $14,
			used_sha256_hashes = $15,
			received_at = $16,
			from_address = $17,
			is_wip = $18,
			ip = $19,
			failed_reason = NULL,
			reserved_zkcounters = $20
	`

	// Get FromAddress from the JSON data
	data, err := state.GetSender(tx.Transaction)
	if err != nil {
		return err
	}
	fromAddress := data.String()

	if _, err := p.db.Exec(ctx, sql,
		hash,
		encoded,
		decoded,
		tx.Status,
		gasPrice,
		nonce,
		tx.GasUsed,
		tx.KeccakHashes,
		tx.PoseidonHashes,
		tx.PoseidonPaddings,
		tx.MemAligns,
		tx.Arithmetics,
		tx.Binaries,
		tx.Steps,
		tx.Sha256Hashes_V2,
		tx.ReceivedAt,
		fromAddress,
		tx.IsWIP,
		tx.IP,
		tx.ReservedZKCounters); err != nil {
		return err
	}
	return nil
}

// GetTxsByStatus returns an array of transactions filtered by status
// limit parameter is used to limit amount txs from the db,
// if limit = 0, then there is no limit
func (p *PostgresPoolStorage) GetTxsByStatus(ctx context.Context, status pool.TxStatus, limit uint64) ([]pool.Transaction, error) {
	var (
		rows pgx.Rows
		err  error
		sql  string
	)
	if limit == 0 {
		sql = `SELECT encoded, status, received_at, is_wip, ip, cumulative_gas_used, used_keccak_hashes, used_poseidon_hashes, used_poseidon_paddings, used_mem_aligns,
				used_arithmetics, used_binaries, used_steps, used_sha256_hashes, failed_reason, reserved_zkcounters FROM pool.transaction WHERE status = $1 ORDER BY gas_price DESC`
		rows, err = p.db.Query(ctx, sql, status.String())
	} else {
		sql = `SELECT encoded, status, received_at, is_wip, ip, cumulative_gas_used, used_keccak_hashes, used_poseidon_hashes, used_poseidon_paddings, used_mem_aligns,
				used_arithmetics, used_binaries, used_steps, used_sha256_hashes, failed_reason, reserved_zkcounters FROM pool.transaction WHERE status = $1 ORDER BY gas_price DESC LIMIT $2`
		rows, err = p.db.Query(ctx, sql, status.String(), limit)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	txs := make([]pool.Transaction, 0, len(rows.RawValues()))
	for rows.Next() {
		tx, err := scanTx(rows)
		if err != nil {
			return nil, err
		}
		txs = append(txs, *tx)
	}

	return txs, nil
}

// GetNonWIPPendingTxs returns an array of transactions
func (p *PostgresPoolStorage) GetNonWIPPendingTxs(ctx context.Context) ([]pool.Transaction, error) {
	var (
		rows pgx.Rows
		err  error
		sql  string
	)

	sql = `SELECT encoded, status, received_at, is_wip, ip, cumulative_gas_used, used_keccak_hashes, used_poseidon_hashes, used_poseidon_paddings, used_mem_aligns,
		used_arithmetics, used_binaries, used_steps, used_sha256_hashes, failed_reason, reserved_zkcounters FROM pool.transaction WHERE is_wip IS FALSE and status = $1`
	rows, err = p.db.Query(ctx, sql, pool.TxStatusPending)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	txs := make([]pool.Transaction, 0, len(rows.RawValues()))
	for rows.Next() {
		tx, err := scanTx(rows)
		if err != nil {
			return nil, err
		}
		txs = append(txs, *tx)
	}

	return txs, nil
}

// GetPendingTxHashesSince returns the pending tx since the given time.
func (p *PostgresPoolStorage) GetPendingTxHashesSince(ctx context.Context, since time.Time) ([]common.Hash, error) {
	sql := "SELECT hash FROM pool.transaction WHERE status = $1 AND received_at >= $2"
	rows, err := p.db.Query(ctx, sql, pool.TxStatusPending, since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	hashes := make([]common.Hash, 0, len(rows.RawValues()))
	for rows.Next() {
		var hash string
		if err := rows.Scan(&hash); err != nil {
			return nil, err
		}
		hashes = append(hashes, common.HexToHash(hash))
	}

	return hashes, nil
}

// GetTxs gets txs with the lowest nonce
func (p *PostgresPoolStorage) GetTxs(ctx context.Context, filterStatus pool.TxStatus, minGasPrice, limit uint64) ([]*pool.Transaction, error) {
	query := `
		SELECT
			encoded,
			status,
			cumulative_gas_used,
			used_keccak_hashes,
			used_poseidon_hashes,
			used_poseidon_paddings,
			used_mem_aligns,
			used_arithmetics,
			used_binaries,
			used_steps,
			used_sha256_hashes,
			received_at,
			nonce,
			is_wip,
			ip,
			reserved_zkcounters
		FROM
			pool.transaction p1
		WHERE 
			status = $1 AND
			gas_price >= $2
		ORDER BY 
			nonce ASC
		LIMIT $3
	`

	if filterStatus == pool.TxStatusFailed {
		query = `
		SELECT * FROM (
			SELECT
				encoded,
				status,
				cumulative_gas_used,
				used_keccak_hashes,
				used_poseidon_hashes,
				used_poseidon_paddings,
				used_mem_aligns,
				used_arithmetics,
				used_binaries,
				used_steps,
				used_sha256_hashes,
				received_at,
				nonce,
				is_wip,
				ip,
				reserved_zkcounters
			FROM
				pool.transaction p1
			WHERE
				status = $1 AND
				gas_price >= $2 
			ORDER BY 
				nonce ASC
			LIMIT $3
			) as tmp
		ORDER BY nonce ASC
		`
	}

	var (
		encoded, status, ip string
		receivedAt          time.Time
		cumulativeGasUsed   uint64
		usedKeccakHashes, usedPoseidonHashes, usedPoseidonPaddings,
		usedMemAligns, usedArithmetics, usedBinaries, usedSteps, usedSHA256Hashes uint32
		nonce              uint64
		isWIP              bool
		reservedZKCounters state.ZKCounters
	)

	args := []interface{}{filterStatus, minGasPrice, limit}

	rows, err := p.db.Query(ctx, query, args...)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, pool.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	txs := make([]*pool.Transaction, 0, len(rows.RawValues()))
	for rows.Next() {
		err := rows.Scan(
			&encoded,
			&status,
			&cumulativeGasUsed,
			&usedKeccakHashes,
			&usedPoseidonHashes,
			&usedPoseidonPaddings,
			&usedMemAligns,
			&usedArithmetics,
			&usedBinaries,
			&usedSteps,
			&usedSHA256Hashes,
			&receivedAt,
			&nonce,
			&isWIP,
			&ip,
			&reservedZKCounters,
		)

		if err != nil {
			return nil, err
		}
		tx := new(pool.Transaction)
		b, err := hex.DecodeHex(encoded)
		if err != nil {
			return nil, err
		}
		if err := tx.UnmarshalBinary(b); err != nil {
			return nil, err
		}

		tx.Status = pool.TxStatus(status)
		tx.ReceivedAt = receivedAt
		tx.ZKCounters = state.ZKCounters{
			GasUsed:          cumulativeGasUsed,
			KeccakHashes:     usedKeccakHashes,
			PoseidonHashes:   usedPoseidonHashes,
			PoseidonPaddings: usedPoseidonPaddings,
			MemAligns:        usedMemAligns,
			Arithmetics:      usedArithmetics,
			Binaries:         usedBinaries,
			Steps:            usedSteps,
			Sha256Hashes_V2:  usedSHA256Hashes,
		}
		tx.IsWIP = isWIP
		tx.IP = ip
		tx.ReservedZKCounters = reservedZKCounters
		txs = append(txs, tx)
	}

	return txs, nil
}

// CountTransactionsByStatus get number of transactions
// accordingly to the provided statuses
func (p *PostgresPoolStorage) CountTransactionsByStatus(ctx context.Context, status ...pool.TxStatus) (uint64, error) {
	sql := "SELECT COUNT(*) FROM pool.transaction WHERE status = ANY ($1)"
	var counter uint64
	err := p.db.QueryRow(ctx, sql, status).Scan(&counter)
	if err != nil {
		return 0, err
	}
	return counter, nil
}

// CountTransactionsByFromAndStatus get number of transactions
// accordingly to the from address and provided statuses
func (p *PostgresPoolStorage) CountTransactionsByFromAndStatus(ctx context.Context, from common.Address, status ...pool.TxStatus) (uint64, error) {
	sql := "SELECT COUNT(*) FROM pool.transaction WHERE from_address = $1 AND status = ANY ($2)"
	var counter uint64
	err := p.db.QueryRow(ctx, sql, from.String(), status).Scan(&counter)
	if err != nil {
		return 0, err
	}
	return counter, nil
}

// UpdateTxStatus updates a transaction status accordingly to the
// provided status and hash
func (p *PostgresPoolStorage) UpdateTxStatus(ctx context.Context, updateInfo pool.TxStatusUpdateInfo) error {
	sql := "UPDATE pool.transaction SET status = $1, is_wip = $2"
	args := []interface{}{updateInfo.NewStatus, updateInfo.IsWIP}

	if updateInfo.FailedReason != nil {
		sql += ", failed_reason = $3"
		args = append(args, *updateInfo.FailedReason)
		sql += " WHERE hash = $4"
	} else {
		sql += " WHERE hash = $3"
	}

	args = append(args, updateInfo.Hash.Hex())

	if _, err := p.db.Exec(ctx, sql, args...); err != nil {
		return err
	}

	return nil
}

// UpdateTxsStatus updates transactions status accordingly to the provided status and hashes
func (p *PostgresPoolStorage) UpdateTxsStatus(ctx context.Context, updateInfos []pool.TxStatusUpdateInfo) error {
	for _, updateInfo := range updateInfos {
		if err := p.UpdateTxStatus(ctx, updateInfo); err != nil {
			return err
		}
	}
	return nil
}

// DeleteTransactionsByHashes deletes txs by their hashes
func (p *PostgresPoolStorage) DeleteTransactionsByHashes(ctx context.Context, hashes []common.Hash) error {
	hh := make([]string, 0, len(hashes))
	for _, h := range hashes {
		hh = append(hh, h.Hex())
	}

	query := "DELETE FROM pool.transaction WHERE hash = ANY ($1)"
	if _, err := p.db.Exec(ctx, query, hh); err != nil {
		return err
	}
	return nil
}

// DeleteFailedTransactionsOlderThan deletes all failed transactions older than the given date
func (p *PostgresPoolStorage) DeleteFailedTransactionsOlderThan(ctx context.Context, date time.Time) error {
	sql := `DELETE FROM pool.transaction WHERE status = 'failed' and received_at < $1`

	if _, err := p.db.Exec(ctx, sql, date); err != nil {
		return err
	}
	return nil
}

// SetGasPrices sets the latest l2 and l1 gas prices
func (p *PostgresPoolStorage) SetGasPrices(ctx context.Context, l2GasPrice, l1GasPrice uint64) error {
	sql := "INSERT INTO pool.gas_price (price, l1_price, timestamp) VALUES ($1, $2, $3)"
	if _, err := p.db.Exec(ctx, sql, l2GasPrice, l1GasPrice, time.Now().UTC()); err != nil {
		return err
	}
	return nil
}

// GetGasPrices returns the latest l2 and l1 gas prices
func (p *PostgresPoolStorage) GetGasPrices(ctx context.Context) (uint64, uint64, error) {
	sql := "SELECT price, l1_price FROM pool.gas_price ORDER BY item_id DESC LIMIT 1"
	rows, err := p.db.Query(ctx, sql)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, 0, state.ErrNotFound
	} else if err != nil {
		return 0, 0, err
	}

	defer rows.Close()

	l2GasPrice := uint64(0)
	l1GasPrice := uint64(0)

	for rows.Next() {
		err := rows.Scan(&l2GasPrice, &l1GasPrice)
		if err != nil {
			return 0, 0, err
		}
	}

	return l2GasPrice, l1GasPrice, nil
}

// DeleteGasPricesHistoryOlderThan deletes all gas prices older than the given date except the last one
func (p *PostgresPoolStorage) DeleteGasPricesHistoryOlderThan(ctx context.Context, date time.Time) error {
	sql := `DELETE FROM pool.gas_price
		WHERE timestamp < $1 AND item_id NOT IN (
			SELECT item_id
			FROM pool.gas_price
			ORDER BY item_id DESC
			LIMIT 1
		)`
	if _, err := p.db.Exec(ctx, sql, date); err != nil {
		return err
	}
	return nil
}

// MinL2GasPriceSince returns the min L2 gas price after given timestamp
func (p *PostgresPoolStorage) MinL2GasPriceSince(ctx context.Context, timestamp time.Time) (uint64, error) {
	sql := "SELECT COALESCE(MIN(price), 0) FROM pool.gas_price WHERE \"timestamp\" >= $1 LIMIT 1"
	var gasPrice uint64
	err := p.db.QueryRow(ctx, sql, timestamp).Scan(&gasPrice)
	if gasPrice == 0 || errors.Is(err, pgx.ErrNoRows) {
		return 0, state.ErrNotFound
	} else if err != nil {
		return 0, err
	}

	return gasPrice, nil
}

// IsTxPending determines if the tx associated to the given hash is pending or
// not.
func (p *PostgresPoolStorage) IsTxPending(ctx context.Context, hash common.Hash) (bool, error) {
	var exists bool
	req := "SELECT exists (SELECT 1 FROM pool.transaction WHERE hash = $1 AND status = $2)"
	err := p.db.QueryRow(ctx, req, hash.Hex(), pool.TxStatusPending).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}

	return exists, nil
}

// GetTxsByFromAndNonce get all the transactions from the pool with the same from and nonce
func (p *PostgresPoolStorage) GetTxsByFromAndNonce(ctx context.Context, from common.Address, nonce uint64) ([]pool.Transaction, error) {
	sql := `SELECT encoded, status, received_at, is_wip, ip, cumulative_gas_used, used_keccak_hashes, used_poseidon_hashes, 
				   used_poseidon_paddings, used_mem_aligns,	used_arithmetics, used_binaries, used_steps, used_sha256_hashes, failed_reason, reserved_zkcounters
	          FROM pool.transaction
			 WHERE from_address = $1
			   AND nonce = $2`
	rows, err := p.db.Query(ctx, sql, from.String(), nonce)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	defer rows.Close()

	txs := make([]pool.Transaction, 0, len(rows.RawValues()))
	for rows.Next() {
		tx, err := scanTx(rows)
		if err != nil {
			return nil, err
		}
		txs = append(txs, *tx)
	}

	return txs, nil
}

// GetTxFromAddressFromByHash gets tx from address by hash
func (p *PostgresPoolStorage) GetTxFromAddressFromByHash(ctx context.Context, hash common.Hash) (common.Address, uint64, error) {
	query := `SELECT from_address, nonce
			  FROM pool.transaction
			  WHERE hash = $1
	`

	var (
		fromAddr string
		nonce    uint64
	)
	err := p.db.QueryRow(ctx, query, hash.String()).Scan(&fromAddr, &nonce)
	if err != nil {
		return common.Address{}, 0, err
	}

	return common.HexToAddress(fromAddr), nonce, nil
}

// GetNonce gets the nonce to the provided address accordingly to the txs in the pool
func (p *PostgresPoolStorage) GetNonce(ctx context.Context, address common.Address) (uint64, error) {
	sql := `SELECT MAX(nonce)
              FROM pool.transaction
             WHERE from_address = $1
               AND status IN  ($2, $3)`
	rows, err := p.db.Query(ctx, sql, address.String(), pool.TxStatusPending, pool.TxStatusSelected)
	if errors.Is(err, pgx.ErrNoRows) {
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	defer rows.Close()

	var nonce *uint64
	for rows.Next() {
		err := rows.Scan(&nonce)
		if err != nil {
			return 0, err
		} else if rows.Err() != nil {
			return 0, rows.Err()
		}
	}

	if nonce == nil {
		n := uint64(0)
		nonce = &n
	} else {
		n := *nonce + 1
		nonce = &n
	}

	return *nonce, nil
}

// GetTransactionByHash gets a transaction in the pool by its hash
func (p *PostgresPoolStorage) GetTransactionByHash(ctx context.Context, hash common.Hash) (*pool.Transaction, error) {
	var (
		encoded, status, ip string
		receivedAt          time.Time
		isWIP               bool
	)

	sql := `SELECT encoded, status, received_at, is_wip, ip
	          FROM pool.transaction
			 WHERE hash = $1`
	err := p.db.QueryRow(ctx, sql, hash.String()).Scan(&encoded, &status, &receivedAt, &isWIP, &ip)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, pool.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	b, err := hex.DecodeHex(encoded)
	if err != nil {
		return nil, err
	}

	tx := new(types.Transaction)
	if err := tx.UnmarshalBinary(b); err != nil {
		return nil, err
	}

	poolTx := &pool.Transaction{
		ReceivedAt:  receivedAt,
		Status:      pool.TxStatus(status),
		Transaction: *tx,
		IsWIP:       isWIP,
		IP:          ip,
	}

	return poolTx, nil
}

// GetTransactionByL2Hash gets a transaction in the pool by its l2 hash
func (p *PostgresPoolStorage) GetTransactionByL2Hash(ctx context.Context, hash common.Hash) (*pool.Transaction, error) {
	var (
		encoded, status, ip string
		receivedAt          time.Time
		isWIP               bool
	)

	sql := `SELECT encoded, status, received_at, is_wip, ip
	          FROM pool.transaction
			 WHERE l2_hash = $1`
	err := p.db.QueryRow(ctx, sql, hash.String()).Scan(&encoded, &status, &receivedAt, &isWIP, &ip)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, pool.ErrNotFound
	} else if err != nil {
		return nil, err
	}

	b, err := hex.DecodeHex(encoded)
	if err != nil {
		return nil, err
	}

	tx := new(types.Transaction)
	if err := tx.UnmarshalBinary(b); err != nil {
		return nil, err
	}

	poolTx := &pool.Transaction{
		ReceivedAt:  receivedAt,
		Status:      pool.TxStatus(status),
		Transaction: *tx,
		IsWIP:       isWIP,
		IP:          ip,
	}

	return poolTx, nil
}

func scanTx(rows pgx.Rows) (*pool.Transaction, error) {
	var (
		encoded, status, ip  string
		receivedAt           time.Time
		isWIP                bool
		cumulativeGasUsed    uint64
		usedKeccakHashes     uint32
		usedPoseidonHashes   uint32
		usedPoseidonPaddings uint32
		usedMemAligns        uint32
		usedArithmetics      uint32
		usedBinaries         uint32
		usedSteps            uint32
		usedSHA256Hashes     uint32
		failedReason         *string
		reservedZKCounters   state.ZKCounters
	)

	if err := rows.Scan(&encoded, &status, &receivedAt, &isWIP, &ip, &cumulativeGasUsed, &usedKeccakHashes, &usedPoseidonHashes,
		&usedPoseidonPaddings, &usedMemAligns, &usedArithmetics, &usedBinaries, &usedSteps, &usedSHA256Hashes, &failedReason, &reservedZKCounters); err != nil {
		return nil, err
	}

	tx := new(pool.Transaction)

	b, err := hex.DecodeHex(encoded)
	if err != nil {
		return nil, err
	}

	if err := tx.UnmarshalBinary(b); err != nil {
		return nil, err
	}

	tx.Status = pool.TxStatus(status)
	tx.ReceivedAt = receivedAt
	tx.IsWIP = isWIP
	tx.IP = ip
	tx.ZKCounters.GasUsed = cumulativeGasUsed
	tx.ZKCounters.KeccakHashes = usedKeccakHashes
	tx.ZKCounters.PoseidonHashes = usedPoseidonHashes
	tx.ZKCounters.PoseidonPaddings = usedPoseidonPaddings
	tx.ZKCounters.MemAligns = usedMemAligns
	tx.ZKCounters.Arithmetics = usedArithmetics
	tx.ZKCounters.Binaries = usedBinaries
	tx.ZKCounters.Steps = usedSteps
	tx.ZKCounters.Sha256Hashes_V2 = usedSHA256Hashes
	tx.FailedReason = failedReason
	tx.ReservedZKCounters = reservedZKCounters

	return tx, nil
}

// DeleteTransactionByHash deletes tx by its hash
func (p *PostgresPoolStorage) DeleteTransactionByHash(ctx context.Context, hash common.Hash) error {
	query := "DELETE FROM pool.transaction WHERE hash = $1"
	if _, err := p.db.Exec(ctx, query, hash); err != nil {
		return err
	}
	return nil
}

// GetTxZkCountersByHash gets a transaction zkcounters by its hash
func (p *PostgresPoolStorage) GetTxZkCountersByHash(ctx context.Context, hash common.Hash) (*state.ZKCounters, *state.ZKCounters, error) {
	var usedZKCounters state.ZKCounters
	var reservedZKCounters state.ZKCounters

	sql := `SELECT cumulative_gas_used, used_keccak_hashes, used_poseidon_hashes, used_poseidon_paddings, used_mem_aligns,
			used_arithmetics, used_binaries, used_steps, used_sha256_hashes, reserved_zkcounters FROM pool.transaction WHERE hash = $1`
	err := p.db.QueryRow(ctx, sql, hash.String()).Scan(&usedZKCounters.GasUsed, &usedZKCounters.KeccakHashes,
		&usedZKCounters.PoseidonHashes, &usedZKCounters.PoseidonPaddings,
		&usedZKCounters.MemAligns, &usedZKCounters.Arithmetics, &usedZKCounters.Binaries, &usedZKCounters.Steps, &usedZKCounters.Sha256Hashes_V2, &reservedZKCounters)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil, pool.ErrNotFound
	} else if err != nil {
		return nil, nil, err
	}

	return &usedZKCounters, &reservedZKCounters, nil
}

// MarkWIPTxsAsPending updates WIP status to non WIP
func (p *PostgresPoolStorage) MarkWIPTxsAsPending(ctx context.Context) error {
	const query = `UPDATE pool.transaction SET is_wip = false WHERE is_wip = true`
	if _, err := p.db.Exec(ctx, query); err != nil {
		return err
	}
	return nil
}

// UpdateTxWIPStatus updates a transaction wip status accordingly to the
// provided WIP status and hash
func (p *PostgresPoolStorage) UpdateTxWIPStatus(ctx context.Context, hash common.Hash, isWIP bool) error {
	sql := "UPDATE pool.transaction SET is_wip = $1 WHERE hash = $2"
	if _, err := p.db.Exec(ctx, sql, isWIP, hash.Hex()); err != nil {
		return err
	}
	return nil
}

// GetAllAddressesBlocked get all addresses blocked
func (p *PostgresPoolStorage) GetAllAddressesBlocked(ctx context.Context) ([]common.Address, error) {
	sql := `SELECT addr FROM pool.blocked`

	rows, err := p.db.Query(ctx, sql)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	defer rows.Close()

	var addrs []common.Address
	for rows.Next() {
		var addr string
		err := rows.Scan(&addr)
		if err != nil {
			return nil, err
		}
		addrs = append(addrs, common.HexToAddress(addr))
	}

	return addrs, nil
}

// GetEarliestProcessedTx gets the earliest processed tx from the pool. Mainly used for cleanup
func (p *PostgresPoolStorage) GetEarliestProcessedTx(ctx context.Context) (common.Hash, error) {
	const getEarliestProcessedTxnFromTxnPool = `SELECT hash
	FROM pool.transaction
	WHERE
		status = 'selected'
	ORDER BY received_at ASC
	LIMIT 1`

	var txnHash string
	err := p.db.QueryRow(ctx, getEarliestProcessedTxnFromTxnPool).Scan(&txnHash)
	if errors.Is(err, pgx.ErrNoRows) {
		return common.Hash{}, nil
	} else if err != nil {
		return common.Hash{}, err
	}

	return common.HexToHash(txnHash), nil
}
