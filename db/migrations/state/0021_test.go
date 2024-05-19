package migrations_test

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

type migrationTest0021 struct{}

func (m migrationTest0021) InsertData(db *sql.DB) error {
	const insertBatch0 = `
		INSERT INTO state.batch (batch_num, global_exit_root, local_exit_root, acc_input_hash, state_root, timestamp, coinbase, raw_txs_data, forced_batch_num, wip) 
		VALUES (0,'0x0000', '0x0000', '0x0000', '0x0000', now(), '0x0000', null, null, true)`

	// insert batch
	_, err := db.Exec(insertBatch0)
	if err != nil {
		return err
	}

	return nil
}

func (m migrationTest0021) RunAssertsAfterMigrationUp(t *testing.T, db *sql.DB) {
	var result int

	// Check column high_reserved_counters exists in state.batch table
	const getColumn = `SELECT count(*) FROM information_schema.columns WHERE table_name='batch' and column_name='high_reserved_counters'`
	row := db.QueryRow(getColumn)
	assert.NoError(t, row.Scan(&result))
	assert.Equal(t, 1, result)

	const insertBatch0 = `
		INSERT INTO state.batch (batch_num, global_exit_root, local_exit_root, acc_input_hash, state_root, timestamp, coinbase, raw_txs_data, forced_batch_num, wip, high_reserved_counters) 
		VALUES (1,'0x0001', '0x0001', '0x0001', '0x0001', now(), '0x0001', null, null, true, '{"Steps": 1890125}')`

	// insert batch 1
	_, err := db.Exec(insertBatch0)
	assert.NoError(t, err)

	const insertBatch1 = `
		INSERT INTO state.batch (batch_num, global_exit_root, local_exit_root, acc_input_hash, state_root, timestamp, coinbase, raw_txs_data, forced_batch_num, wip, high_reserved_counters) 
		VALUES (2,'0x0002', '0x0002', '0x0002', '0x0002', now(), '0x0002', null, null, false, '{"Steps": 1890125}')`

	// insert batch 2
	_, err = db.Exec(insertBatch1)
	assert.NoError(t, err)
}

func (m migrationTest0021) RunAssertsAfterMigrationDown(t *testing.T, db *sql.DB) {
	var result int

	// Check column high_reserved_counters doesn't exists in state.batch table
	const getCheckedColumn = `SELECT count(*) FROM information_schema.columns WHERE table_name='batch' and column_name='high_reserved_counters'`
	row := db.QueryRow(getCheckedColumn)
	assert.NoError(t, row.Scan(&result))
	assert.Equal(t, 0, result)
}

func TestMigration0021(t *testing.T) {
	runMigrationTest(t, 21, migrationTest0021{})
}
