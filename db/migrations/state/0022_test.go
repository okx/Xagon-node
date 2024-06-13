package migrations_test

import (
	"database/sql"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/hex"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

type migrationTest0022TestCase struct {
	Name  string
	Block migrationTest0022TestCaseBlock
}

type migrationTest0022TestCaseBlock struct {
	Transactions []migrationTest0022TestCaseTransaction
}

type migrationTest0022TestCaseTransaction struct {
	CurrentIndex uint
}

type migrationTest0022 struct {
	TestCases []migrationTest0022TestCase
}

func (m migrationTest0022) InsertData(db *sql.DB) error {
	const addBlock0 = "INSERT INTO state.block (block_num, received_at, block_hash) VALUES (0, now(), '0x0')"
	if _, err := db.Exec(addBlock0); err != nil {
		return err
	}

	const addBatch0 = `
		INSERT INTO state.batch (batch_num, global_exit_root, local_exit_root, acc_input_hash, state_root, timestamp, coinbase, raw_txs_data, forced_batch_num, wip) 
		VALUES (0,'0x0000', '0x0000', '0x0000', '0x0000', now(), '0x0000', null, null, true)`
	if _, err := db.Exec(addBatch0); err != nil {
		return err
	}

	const addL2Block = "INSERT INTO state.l2block (block_num, block_hash, header, uncles, parent_hash, state_root, received_at, batch_num, created_at) VALUES ($1, $2, '{}', '{}', '0x0', '0x0', now(), 0, now())"
	const addTransaction = "INSERT INTO state.transaction (hash, encoded, decoded, l2_block_num, effective_percentage, l2_hash) VALUES ($1, 'ABCDEF', '{}', $2, 255, $1)"
	const addReceipt = "INSERT INTO state.receipt (tx_hash, type, post_state, status, cumulative_gas_used, gas_used, effective_gas_price, block_num, tx_index, contract_address) VALUES ($1, 1, null, 1, 1234, 1234, 1, $2, $3, '')"

	txUnique := 0
	for tci, testCase := range m.TestCases {
		blockNum := uint64(tci + 1)
		blockHash := common.HexToHash(hex.EncodeUint64(blockNum)).String()
		if _, err := db.Exec(addL2Block, blockNum, blockHash); err != nil {
			return err
		}
		for _, tx := range testCase.Block.Transactions {
			txUnique++
			txHash := common.HexToHash(hex.EncodeUint64(uint64(txUnique))).String()
			if _, err := db.Exec(addTransaction, txHash, blockNum); err != nil {
				return err
			}
			if _, err := db.Exec(addReceipt, txHash, blockNum, tx.CurrentIndex); err != nil {
				return err
			}
		}
	}

	return nil
}

func (m migrationTest0022) RunAssertsAfterMigrationUp(t *testing.T, db *sql.DB) {
	const getReceiptsByBlock = "SELECT r.tx_index FROM state.receipt r WHERE r.block_num = $1 ORDER BY r.tx_index"

	for tci := range m.TestCases {
		blockNum := uint64(tci + 1)

		rows, err := db.Query(getReceiptsByBlock, blockNum)
		require.NoError(t, err)

		var expectedIndex = uint(0)
		var txIndex uint
		for rows.Next() {
			err := rows.Scan(&txIndex)
			require.NoError(t, err)
			require.Equal(t, expectedIndex, txIndex)
			expectedIndex++
		}
	}
}

func (m migrationTest0022) RunAssertsAfterMigrationDown(t *testing.T, db *sql.DB) {
	m.RunAssertsAfterMigrationUp(t, db)
}

func TestMigration0022(t *testing.T) {
	runMigrationTest(t, 22, migrationTest0022{
		TestCases: []migrationTest0022TestCase{
			{
				Name: "single tx with correct index",
				Block: migrationTest0022TestCaseBlock{
					Transactions: []migrationTest0022TestCaseTransaction{
						{CurrentIndex: 0},
					},
				},
			},
			{
				Name: "multiple txs indexes are correct",
				Block: migrationTest0022TestCaseBlock{
					Transactions: []migrationTest0022TestCaseTransaction{
						{CurrentIndex: 0},
						{CurrentIndex: 1},
						{CurrentIndex: 2},
					},
				},
			},
			{
				Name: "single tx with wrong tx index",
				Block: migrationTest0022TestCaseBlock{
					Transactions: []migrationTest0022TestCaseTransaction{
						{CurrentIndex: 3},
					},
				},
			},
			{
				Name: "multiple txs missing 0 index",
				Block: migrationTest0022TestCaseBlock{
					Transactions: []migrationTest0022TestCaseTransaction{
						{CurrentIndex: 1},
						{CurrentIndex: 2},
						{CurrentIndex: 3},
						{CurrentIndex: 4},
					},
				},
			},
			{
				Name: "multiple has index 0 but also txs index gap",
				Block: migrationTest0022TestCaseBlock{
					Transactions: []migrationTest0022TestCaseTransaction{
						{CurrentIndex: 0},
						{CurrentIndex: 2},
						{CurrentIndex: 4},
						{CurrentIndex: 6},
					},
				},
			},
		},
	})
}
