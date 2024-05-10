package migrations_test

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

type migrationTest0016 struct{}

func (m migrationTest0016) InsertData(db *sql.DB) error {
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

func (m migrationTest0016) RunAssertsAfterMigrationUp(t *testing.T, db *sql.DB) {
	var result int

	// Check column checked exists in state.batch table
	const getCheckedColumn = `SELECT count(*) FROM information_schema.columns WHERE table_name='batch' and column_name='checked'`
	row := db.QueryRow(getCheckedColumn)
	assert.NoError(t, row.Scan(&result))
	assert.Equal(t, 1, result)

	const insertBatch0 = `
		INSERT INTO state.batch (batch_num, global_exit_root, local_exit_root, acc_input_hash, state_root, timestamp, coinbase, raw_txs_data, forced_batch_num, wip, checked) 
		VALUES (1,'0x0001', '0x0001', '0x0001', '0x0001', now(), '0x0001', null, null, true, false)`

	// insert batch 1
	_, err := db.Exec(insertBatch0)
	assert.NoError(t, err)

	const insertBatch1 = `
		INSERT INTO state.batch (batch_num, global_exit_root, local_exit_root, acc_input_hash, state_root, timestamp, coinbase, raw_txs_data, forced_batch_num, wip, checked) 
		VALUES (2,'0x0002', '0x0002', '0x0002', '0x0002', now(), '0x0002', null, null, false, true)`

	// insert batch 2
	_, err = db.Exec(insertBatch1)
	assert.NoError(t, err)
}

func (m migrationTest0016) RunAssertsAfterMigrationDown(t *testing.T, db *sql.DB) {
	var result int

	// Check column wip doesn't exists in state.batch table
	const getCheckedColumn = `SELECT count(*) FROM information_schema.columns WHERE table_name='batch' and column_name='checked'`
	row := db.QueryRow(getCheckedColumn)
	assert.NoError(t, row.Scan(&result))
	assert.Equal(t, 0, result)
}

func TestMigration0016(t *testing.T) {
	runMigrationTest(t, 16, migrationTest0016{})
}
