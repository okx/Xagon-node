package pgpoolstorage

import (
	"context"
	"errors"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/pool/trace"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgx/v4"
)

// GetAllAddressesWhitelisted get all addresses whitelisted
func (p *PostgresPoolStorage) GetAllAddressesWhitelisted(ctx context.Context) ([]common.Address, error) {
	sql := `SELECT addr FROM pool.whitelisted`

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

// CREATE TABLE pool.innertx (
// hash VARCHAR(128) PRIMARY KEY NOT NULL,
// innertx text,
// created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
// );

// AddInnerTx add inner tx
func (p *PostgresPoolStorage) AddInnerTx(ctx context.Context, txHash common.Hash, innerTx []byte) error {
	sql := `INSERT INTO pool.innertx(hash, innertx) VALUES ($1, $2)`

	_, err := p.db.Exec(ctx, sql, txHash.Hex(), innerTx)
	if err != nil {
		return err
	}

	return nil
}

// GetInnerTx get inner tx
func (p *PostgresPoolStorage) GetInnerTx(ctx context.Context, txHash common.Hash) (string, error) {
	sql := `SELECT innertx FROM pool.innertx WHERE hash = $1`

	var innerTx string
	err := p.db.QueryRow(ctx, sql, txHash.Hex()).Scan(&innerTx)
	if err != nil {
		return "", err
	}

	return innerTx, nil
}

// GetPendingFromAndMinNonceBefore get pending from and min nonce before timeDuration
func (p *PostgresPoolStorage) GetPendingFromAndMinNonceBefore(ctx context.Context, timeDuration time.Duration) ([]common.Address, []uint64, error) {
	sql := `SELECT from_address, MIN(nonce) FROM pool."transaction" where status='pending' and "received_at" < $1 GROUP BY from_address`

	mLog := log.WithFields(trace.GetID(ctx))
	timeStamp := time.Now().Add(-timeDuration)
	rows, err := p.db.Query(ctx, sql, timeStamp)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			mLog.Infof("no pending transactions before %v", timeStamp)
			return nil, nil, nil
		} else {
			return nil, nil, err
		}
	}
	defer rows.Close()

	var addresses []common.Address
	var nonces []uint64
	for rows.Next() {
		var address string
		var nonce uint64
		err := rows.Scan(&address, &nonce)
		if err != nil {
			return nil, nil, err
		}
		addresses = append(addresses, common.HexToAddress(address))
		nonces = append(nonces, nonce)
	}
	mLog.Infof("pending address count %v before %v", len(addresses), timeStamp)

	return addresses, nonces, nil
}

// CREATE TABLE pool.stat (
// id INT PRIMARY KEY NOT NULL,
// total INT,
// skip_nonce INT,
// balance_issue INT,
// nonce_issue INT,
// locked INT,    // 1 for unlocked 2 for locked
// created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
// updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
// );
// insert into pool.stat(id, total, skip_nonce, balance_issue, nonce_issue, locked) values(1, 0, 0, 0, 0, 1);

// LockStat lock stat
func (p *PostgresPoolStorage) LockStat(ctx context.Context, timeDuration time.Duration) (bool, error) {
	timeStamp := time.Now().Add(-timeDuration)
	sql := `UPDATE pool.stat SET locked = 2 WHERE locked = 1 and updated_at < $1 and id=1`

	stat, err := p.db.Exec(ctx, sql, timeStamp)
	if err != nil {
		return false, err
	}
	if stat.RowsAffected() > 0 {
		return true, nil
	}

	return false, nil
}

// UnLockStat unlock stat
func (p *PostgresPoolStorage) UnLockStat(ctx context.Context) error {
	sql := `UPDATE pool.stat SET locked = 1 WHERE locked = 2 and id=1`

	_, err := p.db.Exec(ctx, sql)
	if err != nil {
		return err
	}

	return nil
}

// UpdateStatAndUnlock update stat and unlock
func (p *PostgresPoolStorage) UpdateStatAndUnlock(ctx context.Context, totoal, skip, balanceIssue, nonceIssue uint64) error {
	sql := `UPDATE pool.stat SET total = $1, skip_nonce = $2, balance_issue = $3, nonce_issue = $4, locked = 1, updated_at = CURRENT_TIMESTAMP WHERE id=1`

	_, err := p.db.Exec(ctx, sql, totoal, skip, balanceIssue, nonceIssue)
	if err != nil {
		return err
	}

	return nil
}

// GetStat get stat
func (p *PostgresPoolStorage) GetStat(ctx context.Context) (uint64, uint64, uint64, uint64, error) {
	sql := `SELECT total, skip_nonce, balance_issue, nonce_issue FROM pool.stat WHERE id=1`

	var total, skip, balanceIssue, nonceIssue uint64
	err := p.db.QueryRow(ctx, sql).Scan(&total, &skip, &balanceIssue, &nonceIssue)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	return total, skip, balanceIssue, nonceIssue, nil
}

// CountTransactionsByFromStatusAndNonce count transactions by from status and nonce
func (p *PostgresPoolStorage) CountTransactionsByFromStatusAndNonce(ctx context.Context, from common.Address, nonce uint64, status ...pool.TxStatus) (uint64, error) {
	sql := "SELECT COUNT(*) FROM pool.transaction WHERE from_address = $1 AND nonce <= $2 AND status = ANY ($3)"
	var counter uint64
	err := p.db.QueryRow(ctx, sql, from.String(), nonce, status).Scan(&counter)
	if err != nil {
		return 0, err
	}
	return counter, nil
}
