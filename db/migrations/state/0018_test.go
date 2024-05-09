package migrations_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type migrationTest0018 struct{}

func (m migrationTest0018) InsertData(db *sql.DB) error {
	const addBlock = "INSERT INTO state.block (block_num, received_at, block_hash) VALUES ($1, $2, $3)"
	if _, err := db.Exec(addBlock, 1, time.Now(), "0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"); err != nil {
		return err
	}
	if _, err := db.Exec(addBlock, 50, time.Now(), "0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"); err != nil {
		return err
	}
	if _, err := db.Exec(addBlock, 1050, time.Now(), "0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"); err != nil {
		return err
	}
	return nil
}

func (m migrationTest0018) RunAssertsAfterMigrationUp(t *testing.T, db *sql.DB) {
	var checked bool
	row := db.QueryRow("SELECT checked FROM state.block WHERE block_num = $1", 1)
	assert.NoError(t, row.Scan(&checked))
	assert.Equal(t, true, checked)
	row = db.QueryRow("SELECT checked FROM state.block WHERE block_num = $1", 50)
	assert.NoError(t, row.Scan(&checked))
	assert.Equal(t, true, checked)
	row = db.QueryRow("SELECT checked FROM state.block WHERE block_num = $1", 1050)
	assert.NoError(t, row.Scan(&checked))
	assert.Equal(t, false, checked)

	const addBlock = "INSERT INTO state.block (block_num, received_at, block_hash, checked) VALUES ($1, $2, $3, $4)"
	_, err := db.Exec(addBlock, 2, time.Now(), "0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1", true)
	assert.NoError(t, err)
	_, err = db.Exec(addBlock, 3, time.Now(), "0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1", false)
	assert.NoError(t, err)
	const sql = `SELECT count(*) FROM state.block WHERE checked = true`
	row = db.QueryRow(sql)
	var result int
	assert.NoError(t, row.Scan(&result))
	assert.Equal(t, 3, result, "must be 1,50 per migration and 2 by insert")

	const sqlCheckedFalse = `SELECT count(*) FROM state.block WHERE checked = false`
	row = db.QueryRow(sqlCheckedFalse)

	assert.NoError(t, row.Scan(&result))
	assert.Equal(t, 2, result, "must be 150 by migration, and 3 by insert")
}

func (m migrationTest0018) RunAssertsAfterMigrationDown(t *testing.T, db *sql.DB) {
	var result int

	// Check column wip doesn't exists in state.batch table
	const sql = `SELECT count(*) FROM state.block`
	row := db.QueryRow(sql)
	assert.NoError(t, row.Scan(&result))
	assert.Equal(t, 5, result)
}

func TestMigration0018(t *testing.T) {
	runMigrationTest(t, 18, migrationTest0018{})
}
