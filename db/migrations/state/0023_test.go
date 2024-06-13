package migrations_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type migrationTest0023 struct {
	migrationBase

	blockHashValue         string
	mainExitRootValue      string
	rollupExitRootValue    string
	globalExitRootValue    string
	previousBlockHashValue string
	l1InfoRootValue        string
}

func (m migrationTest0023) insertBlock(blockNumber uint64, db *sql.DB) error {
	const addBlock = "INSERT INTO state.block (block_num, received_at, block_hash) VALUES ($1, $2, $3)"
	if _, err := db.Exec(addBlock, blockNumber, time.Now(), m.blockHashValue); err != nil {
		return err
	}
	return nil
}

func (m migrationTest0023) insertRowInOldTable(db *sql.DB, args ...interface{}) error {
	sql := `
    INSERT INTO state.exit_root (block_num, "timestamp", mainnet_exit_root, rollup_exit_root, global_exit_root, prev_block_hash, l1_info_root, l1_info_tree_index)
                         VALUES (       $1,          $2,                $3,               $4,               $5,              $6,           $7,                 $8);`

	_, err := db.Exec(sql, args...)
	return err
}

func (m migrationTest0023) insertRowInMigratedTable(db *sql.DB, args ...interface{}) error {
	sql := `
    INSERT INTO state.exit_root (block_num, "timestamp", mainnet_exit_root, rollup_exit_root, global_exit_root, prev_block_hash, l1_info_root, l1_info_tree_index, l1_info_tree_recursive_index)
                         VALUES (       $1,          $2,                $3,               $4,               $5,              $6,           $7,                 $8,                        $9);`

	_, err := db.Exec(sql, args...)
	return err
}

func (m migrationTest0023) InsertData(db *sql.DB) error {
	var err error
	for i := uint64(1); i <= 6; i++ {
		if err = m.insertBlock(i, db); err != nil {
			return err
		}
	}

	return nil
}

func (m migrationTest0023) RunAssertsAfterMigrationUp(t *testing.T, db *sql.DB) {
	m.AssertNewAndRemovedItemsAfterMigrationUp(t, db)

	var nilL1InfoTreeIndex *uint = nil
	err := m.insertRowInOldTable(db, 1, time.Now().UTC(), m.mainExitRootValue, m.rollupExitRootValue, m.globalExitRootValue, m.previousBlockHashValue, m.l1InfoRootValue, nilL1InfoTreeIndex)
	assert.NoError(t, err)

	err = m.insertRowInOldTable(db, 2, time.Now().UTC(), m.mainExitRootValue, m.rollupExitRootValue, m.globalExitRootValue, m.previousBlockHashValue, m.l1InfoRootValue, uint(1))
	assert.NoError(t, err)

	err = m.insertRowInMigratedTable(db, 3, time.Now().UTC(), m.mainExitRootValue, m.rollupExitRootValue, m.globalExitRootValue, m.previousBlockHashValue, m.l1InfoRootValue, nilL1InfoTreeIndex, 1)
	assert.NoError(t, err)
}

func (m migrationTest0023) RunAssertsAfterMigrationDown(t *testing.T, db *sql.DB) {
	m.AssertNewAndRemovedItemsAfterMigrationDown(t, db)

	var nilL1InfoTreeIndex *uint = nil
	err := m.insertRowInOldTable(db, 4, time.Now().UTC(), m.mainExitRootValue, m.rollupExitRootValue, m.globalExitRootValue, m.previousBlockHashValue, m.l1InfoRootValue, nilL1InfoTreeIndex)
	assert.NoError(t, err)

	err = m.insertRowInOldTable(db, 5, time.Now().UTC(), m.mainExitRootValue, m.rollupExitRootValue, m.globalExitRootValue, m.previousBlockHashValue, m.l1InfoRootValue, uint(2))
	assert.NoError(t, err)

	err = m.insertRowInMigratedTable(db, 6, time.Now().UTC(), m.mainExitRootValue, m.rollupExitRootValue, m.globalExitRootValue, m.previousBlockHashValue, m.l1InfoRootValue, nilL1InfoTreeIndex, 2)
	assert.Error(t, err)
}

func TestMigration0023(t *testing.T) {
	m := migrationTest0023{
		migrationBase: migrationBase{
			newIndexes: []string{
				"idx_exit_root_l1_info_tree_recursive_index",
			},
			newColumns: []columnMetadata{
				{"state", "exit_root", "l1_info_tree_recursive_index"},
			},
		},

		blockHashValue:         "0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1",
		mainExitRootValue:      "0x83fc198de31e1b2b1a8212d2430fbb7766c13d9ad305637dea3759065606475d",
		rollupExitRootValue:    "0xadb91a6a1fce56eaea561002bc9a993f4e65a7710bd72f4eee3067cbd73a743c",
		globalExitRootValue:    "0x5bf4af1a651a2a74b36e6eb208481f94c69fc959f756223dfa49608061937585",
		previousBlockHashValue: "0xe865e912b504572a4d80ad018e29797e3c11f00bf9ae2549548a25779c9d7e57",
		l1InfoRootValue:        "0x2b9484b83c6398033241865b015fb9430eb3e159182a6075d00c924845cc393e",
	}
	runMigrationTest(t, 23, m)
}
