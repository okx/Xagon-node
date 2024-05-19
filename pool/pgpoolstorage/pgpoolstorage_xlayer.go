package pgpoolstorage

import (
	"context"
	"errors"
	"time"

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

// CREATE TABLE pool.readytx(
// id SERIAL PRIMARY KEY NOT NULL,
// count INT,
// updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
// );
// insert into pool.readytx(id, count) values(1, 0);

// UpdateReadyTxCount update ready tx count
func (p *PostgresPoolStorage) UpdateReadyTxCount(ctx context.Context, count uint64) error {
	sql := `UPDATE pool.readytx SET count = $1, updated_at = $2 WHERE id=1`

	_, err := p.db.Exec(ctx, sql, count, time.Now())
	if err != nil {
		return err
	}

	return nil
}

// GetReadyTxCount get ready tx count
func (p *PostgresPoolStorage) GetReadyTxCount(ctx context.Context) (uint64, error) {
	sql := `SELECT count FROM pool.readytx where id=1`

	var count uint64
	err := p.db.QueryRow(ctx, sql).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
