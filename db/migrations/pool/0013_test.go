package pool_migrations_test

import (
	"database/sql"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/stretchr/testify/require"
)

// this migration adds reserved_zkcounters to the transaction
type migrationTest0013 struct{}

func (m migrationTest0013) InsertData(db *sql.DB) error {
	return nil
}

func (m migrationTest0013) RunAssertsAfterMigrationUp(t *testing.T, db *sql.DB) {
	var reserved_zkcounters = state.ZKCounters{
		GasUsed:        0,
		KeccakHashes:   1,
		PoseidonHashes: 2,
	}

	const insertTx = `
		INSERT INTO pool.transaction (hash, ip, received_at, from_address, reserved_zkcounters)
		VALUES ('0x0001', '127.0.0.1', '2023-12-07', '0x0011', $1)`

	_, err := db.Exec(insertTx, reserved_zkcounters)
	require.NoError(t, err)
}

func (m migrationTest0013) RunAssertsAfterMigrationDown(t *testing.T, db *sql.DB) {
	var reserved_zkcounters = state.ZKCounters{
		GasUsed:        0,
		KeccakHashes:   1,
		PoseidonHashes: 2,
	}

	const insertTx = `
		INSERT INTO pool.transaction (hash, ip, received_at, from_address, reserved_zkcounters)
		VALUES ('0x0001', '127.0.0.1', '2023-12-07', '0x0011', $1)`

	_, err := db.Exec(insertTx, reserved_zkcounters)
	require.Error(t, err)
}

func TestMigration0013(t *testing.T) {
	runMigrationTest(t, 13, migrationTest0013{})
}
