package migrations_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// this migration changes length of the token name
type migrationTest0013 struct {
	blockHashValue         string
	mainExitRootValue      string
	rollupExitRootValue    string
	globalExitRootValue    string
	previousBlockHashValue string
	l1InfoRootValue        string
}

func (m migrationTest0013) insertBlock(blockNumber uint64, db *sql.DB) error {
	const addBlock = "INSERT INTO state.block (block_num, received_at, block_hash) VALUES ($1, $2, $3)"
	if _, err := db.Exec(addBlock, blockNumber, time.Now(), m.blockHashValue); err != nil {
		return err
	}
	return nil
}

func (m migrationTest0013) insertRowInOldTable(db *sql.DB, args []interface{}) error {
	insert := `
        INSERT INTO state.exit_root (block_num, timestamp, mainnet_exit_root, rollup_exit_root, global_exit_root) 
                             VALUES ($1,        $2,        $3,                $4,               $5  );`

	_, err := db.Exec(insert, args...)
	return err
}

func (m migrationTest0013) InsertData(db *sql.DB) error {
	var err error
	if err = m.insertBlock(uint64(123), db); err != nil {
		return err
	}
	if err = m.insertBlock(uint64(124), db); err != nil {
		return err
	}
	if err = m.insertRowInOldTable(db, []interface{}{123, time.Now(), m.mainExitRootValue, m.rollupExitRootValue, m.globalExitRootValue}); err != nil {
		return err
	}
	if err = m.insertRowInOldTable(db, []interface{}{124, time.Now(), m.mainExitRootValue, m.rollupExitRootValue, m.globalExitRootValue}); err != nil {
		return err
	}

	return nil
}

func (m migrationTest0013) InsertDataIntoTransactionsTable(db *sql.DB) error {
	// Insert block to respect the FKey
	const addBlock = "INSERT INTO state.block (block_num, received_at, block_hash) VALUES ($1, $2, $3)"
	if _, err := db.Exec(addBlock, 1, time.Now(), "0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"); err != nil {
		return err
	}
	if _, err := db.Exec(addBlock, 2, time.Now(), "0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f2"); err != nil {
		return err
	}
	const insertBatch = `
		INSERT INTO state.batch (batch_num, global_exit_root, local_exit_root, acc_input_hash, state_root, timestamp, coinbase, raw_txs_data, forced_batch_num, wip) 
		VALUES (0,'0x0000', '0x0000', '0x0000', '0x0000', now(), '0x0000', null, null, true)`

	// insert batch
	_, err := db.Exec(insertBatch)
	if err != nil {
		return err
	}

	const insertL2Block = `
		INSERT INTO state.l2block (block_num, block_hash, header, uncles, parent_hash, state_root, received_at, batch_num, created_at)
		VALUES (0, '0x0001', '{}', '{}', '0x0002', '0x003', now(), 0, now())`

	// insert l2 block
	_, err = db.Exec(insertL2Block)
	if err != nil {
		return err
	}

	const insertTx = `
		INSERT INTO state.transaction (hash, encoded, decoded, l2_block_num, effective_percentage, l2_hash)
		VALUES ('0x0001', 'ABCDEF', '{}', 0, 255, '0x0002')`

	// insert tx
	_, err = db.Exec(insertTx)
	if err != nil {
		return err
	}

	return nil
}

func (m migrationTest0013) insertRowInMigratedTable(db *sql.DB, args []interface{}) error {
	insert := `
        INSERT INTO state.exit_root (block_num, timestamp, mainnet_exit_root, rollup_exit_root, global_exit_root, prev_block_hash, l1_info_root, l1_info_tree_index) 
                             VALUES ($1,        $2,        $3,                $4,               $5,               $6,              $7,           $8);`

	_, err := db.Exec(insert, args...)
	return err
}

func (m migrationTest0013) RunAssertsAfterMigrationUp(t *testing.T, db *sql.DB) {
	err := m.insertBlock(uint64(125), db)
	assert.NoError(t, err)
	err = m.insertBlock(uint64(126), db)
	assert.NoError(t, err)
	err = m.insertBlock(uint64(127), db)
	assert.NoError(t, err)
	prevBlockHash := m.previousBlockHashValue
	l1InfoRoot := m.l1InfoRootValue
	err = m.insertRowInMigratedTable(db, []interface{}{125, time.Now(), m.mainExitRootValue, m.rollupExitRootValue, m.globalExitRootValue, prevBlockHash, l1InfoRoot, 1})
	assert.NoError(t, err)
	// insert duplicated l1_info_root
	err = m.insertRowInMigratedTable(db, []interface{}{126, time.Now(), m.mainExitRootValue, m.rollupExitRootValue, m.globalExitRootValue, prevBlockHash, l1InfoRoot, 1})
	assert.Error(t, err)

	// insert in the old way must work
	err = m.insertRowInOldTable(db, []interface{}{127, time.Now(), m.mainExitRootValue, m.rollupExitRootValue, m.globalExitRootValue})
	assert.NoError(t, err)

	sqlSelect := `SELECT prev_block_hash, l1_info_root FROM state.exit_root WHERE l1_info_tree_index = $1`
	currentPrevBlockHash := ""
	currentL1InfoRoot := ""
	err = db.QueryRow(sqlSelect, 1).Scan(&currentPrevBlockHash, &currentL1InfoRoot)
	assert.NoError(t, err)
	assert.Equal(t, prevBlockHash, currentPrevBlockHash)
	assert.Equal(t, l1InfoRoot, currentL1InfoRoot)

	// Check column l2_hash exists in state.transactions table
	const getL2HashColumn = `SELECT count(*) FROM information_schema.columns WHERE table_name='transaction' and column_name='l2_hash'`
	row := db.QueryRow(getL2HashColumn)
	var result int
	assert.NoError(t, row.Scan(&result))
	assert.Equal(t, 1, result)

	// Check column wip exists in state.batch table
	const getWIPColumn = `SELECT count(*) FROM information_schema.columns WHERE table_name='batch' and column_name='wip'`
	row = db.QueryRow(getWIPColumn)
	assert.NoError(t, row.Scan(&result))
	assert.Equal(t, 1, result)

	// Try to insert data into the transactions table
	err = m.InsertDataIntoTransactionsTable(db)
	assert.NoError(t, err)

	insertVirtualBatch := `INSERT INTO state.virtual_batch
	(batch_num, tx_hash, coinbase, block_num, sequencer_addr, timestamp_batch_etrog)
	VALUES(0, '0x23970ef3f8184daa93385faf802df142a521b479e8e59fbeafa11b8927eb77b1', '0x0000000000000000000000000000000000000000', 1, '0x6645F64d1cE0513bbf5E6713b7e4D0A957AC853c', '2023-12-22 16:53:00.000');`
	_, err = db.Exec(insertVirtualBatch)
	assert.NoError(t, err)
}

func (m migrationTest0013) RunAssertsAfterMigrationDown(t *testing.T, db *sql.DB) {
	sqlSelect := `SELECT count(id) FROM state.exit_root`
	count := 0
	err := db.QueryRow(sqlSelect).Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 4, count)

	// Check column l2_hash doesn't exists in state.transactions table
	const getL2HashColumn = `SELECT count(*) FROM information_schema.columns WHERE table_name='transaction' and column_name='l2_hash'`
	row := db.QueryRow(getL2HashColumn)
	var result int
	assert.NoError(t, row.Scan(&result))
	assert.Equal(t, 0, result)

	// Check column wip doesn't exists in state.batch table
	const getWIPColumn = `SELECT count(*) FROM information_schema.columns WHERE table_name='batch' and column_name='wip'`
	row = db.QueryRow(getWIPColumn)
	assert.NoError(t, row.Scan(&result))
	assert.Equal(t, 0, result)

	insertVirtualBatch := `INSERT INTO state.virtual_batch
	(batch_num, tx_hash, coinbase, block_num, sequencer_addr, timestamp_batch_etrog)
	VALUES(0, '0x23970ef3f8184daa93385faf802df142a521b479e8e59fbeafa11b8927eb77b1', '0x0000000000000000000000000000000000000000', 1, '0x6645F64d1cE0513bbf5E6713b7e4D0A957AC853c', '2023-12-22 16:53:00.000');`
	_, err = db.Exec(insertVirtualBatch)
	assert.Error(t, err)
}

func TestMigration0013(t *testing.T) {
	m := migrationTest0013{
		blockHashValue:         "0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1",
		mainExitRootValue:      "0x83fc198de31e1b2b1a8212d2430fbb7766c13d9ad305637dea3759065606475d",
		rollupExitRootValue:    "0xadb91a6a1fce56eaea561002bc9a993f4e65a7710bd72f4eee3067cbd73a743c",
		globalExitRootValue:    "0x5bf4af1a651a2a74b36e6eb208481f94c69fc959f756223dfa49608061937585",
		previousBlockHashValue: "0xe865e912b504572a4d80ad018e29797e3c11f00bf9ae2549548a25779c9d7e57",
		l1InfoRootValue:        "0x2b9484b83c6398033241865b015fb9430eb3e159182a6075d00c924845cc393e",
	}
	runMigrationTest(t, 13, m)
}
