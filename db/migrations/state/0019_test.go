package migrations_test

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

type migrationTest0019 struct {
	migrationBase
}

func (m migrationTest0019) InsertData(db *sql.DB) error {
	const insertBatch1 = `
		INSERT INTO state.batch (batch_num, global_exit_root, local_exit_root, acc_input_hash, state_root, timestamp, coinbase, raw_txs_data, forced_batch_num, wip) 
		VALUES (1,'0x0001', '0x0001', '0x0001', '0x0001', now(), '0x0001', null, null, true)`

	_, err := db.Exec(insertBatch1)
	if err != nil {
		return err
	}

	const insertBatch2 = `
		INSERT INTO state.batch (batch_num, global_exit_root, local_exit_root, acc_input_hash, state_root, timestamp, coinbase, raw_txs_data, forced_batch_num, wip) 
		VALUES (2,'0x0002', '0x0002', '0x0002', '0x0002', now(), '0x0002', null, null, true)`

	_, err = db.Exec(insertBatch2)
	if err != nil {
		return err
	}

	const insertBlock1 = "INSERT INTO state.block (block_num, block_hash, parent_hash, received_at) VALUES (1,'0x0001', '0x0001', now())"

	_, err = db.Exec(insertBlock1)
	if err != nil {
		return err
	}

	const insertBlock2 = "INSERT INTO state.block (block_num, block_hash, parent_hash, received_at) VALUES (2,'0x0002', '0x0002', now())"

	_, err = db.Exec(insertBlock2)
	if err != nil {
		return err
	}

	return nil
}

func (m migrationTest0019) RunAssertsAfterMigrationUp(t *testing.T, db *sql.DB) {
	m.AssertNewAndRemovedItemsAfterMigrationUp(t, db)

	// Insert blobInner 1
	const insertBlobInner = `INSERT INTO state.blob_inner (blob_inner_num, data, block_num) VALUES (1, E'\\x1234', 1);`
	_, err := db.Exec(insertBlobInner)
	assert.NoError(t, err)

	const insertBatch1 = `
		INSERT INTO state.virtual_batch (batch_num, tx_hash, coinbase, block_num, sequencer_addr, timestamp_batch_etrog, l1_info_root, blob_inner_num, prev_l1_it_root, prev_l1_it_index) 
		VALUES (1,'0x0001', '0x0001', 1, '0x0001', now(), '0x0001', 1, '0x0001', 1)`

	_, err = db.Exec(insertBatch1)
	assert.NoError(t, err)

	const insertBatch2 = `
		INSERT INTO state.virtual_batch (batch_num, tx_hash, coinbase, block_num, sequencer_addr, timestamp_batch_etrog, l1_info_root, blob_inner_num, prev_l1_it_root, prev_l1_it_index) 
		VALUES (2,'0x0002', '0x0002', 2, '0x0002', now(), '0x0002', 1, '0x0002', 2)`

	_, err = db.Exec(insertBatch2)
	assert.NoError(t, err)
}

func (m migrationTest0019) RunAssertsAfterMigrationDown(t *testing.T, db *sql.DB) {
	var result int

	m.AssertNewAndRemovedItemsAfterMigrationDown(t, db)

	// Check column blob_inner_num doesn't exists in state.virtual_batch table
	const getBlobInnerNumColumn = `SELECT count(*) FROM information_schema.columns WHERE table_name='virtual_batch' and column_name='blob_inner_num'`
	row := db.QueryRow(getBlobInnerNumColumn)
	assert.NoError(t, row.Scan(&result))
	assert.Equal(t, 0, result)

	// Check column prev_l1_it_root doesn't exists in state.virtual_batch table
	const getPrevL1ITRootColumn = `SELECT count(*) FROM information_schema.columns WHERE table_name='virtual_batch' and column_name='prev_l1_it_root'`
	row = db.QueryRow(getPrevL1ITRootColumn)
	assert.NoError(t, row.Scan(&result))
	assert.Equal(t, 0, result)

	// Check column prev_l1_it_index doesn't exists in state.virtual_batch table
	const getPrevL1ITIndexColumn = `SELECT count(*) FROM information_schema.columns WHERE table_name='virtual_batch' and column_name='prev_l1_it_index'`
	row = db.QueryRow(getPrevL1ITIndexColumn)
	assert.NoError(t, row.Scan(&result))
	assert.Equal(t, 0, result)
}

func TestMigration0019(t *testing.T) {
	m := migrationTest0019{
		migrationBase: migrationBase{
			removedTables: []tableMetadata{
				{"state", "proof"},
			},

			newTables: []tableMetadata{
				{"state", "blob_inner"},
				{"state", "batch_proof"},
				{"state", "blob_inner_proof"},
				{"state", "blob_outer_proof"},
			},

			newColumns: []columnMetadata{
				{"state", "virtual_batch", "blob_inner_num"},
				{"state", "virtual_batch", "prev_l1_it_root"},
				{"state", "virtual_batch", "prev_l1_it_index"},
			},
		},
	}
	runMigrationTest(t, 19, m)
}
