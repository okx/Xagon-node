package migrations_test

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// this migration changes length of the token name
type migrationTest0020 struct{}

func (m migrationTest0020) InsertData(db *sql.DB) error {
	addBlocks := `
	INSERT INTO state.block
	(block_num, block_hash, parent_hash, received_at, checked)
	VALUES(1, '0x013be63487a53c874614dd1ae0434cf211e393b2e386c8fde74da203b5469b20', '0x0328698ebeda498df8c63040e2a4771d24722ab2c1e8291226b9215c7eec50fe', '2024-03-11 02:52:23.000', true);
	INSERT INTO state.block
	(block_num, block_hash, parent_hash, received_at, checked)
	VALUES(2, '0x013be63487a53c874614dd1ae0434cf211e393b2e386c8fde74da203b5469b21', '0x0328698ebeda498df8c63040e2a4771d24722ab2c1e8291226b9215c7eec50f1', '2024-03-11 02:52:24.000', true);
	INSERT INTO state.block
	(block_num, block_hash, parent_hash, received_at, checked)
	VALUES(3, '0x013be63487a53c874614dd1ae0434cf211e393b2e386c8fde74da203b5469b22', '0x0328698ebeda498df8c63040e2a4771d24722ab2c1e8291226b9215c7eec50f2', '2024-03-11 02:52:25.000', false);
	INSERT INTO state.block
	(block_num, block_hash, parent_hash, received_at, checked)
	VALUES(4, '0x013be63487a53c874614dd1ae0434cf211e393b2e386c8fde74da203b5469b23', '0x0328698ebeda498df8c63040e2a4771d24722ab2c1e8291226b9215c7eec50f3', '2024-03-11 02:52:26.000', false);
	INSERT INTO state.block
	(block_num, block_hash, parent_hash, received_at, checked)
	VALUES(5, '0x013be63487a53c874614dd1ae0434cf211e393b2e386c8fde74da203b5469b24', '0x0328698ebeda498df8c63040e2a4771d24722ab2c1e8291226b9215c7eec50f4', '2024-03-11 02:52:27.000', true);
	INSERT INTO state.block
	(block_num, block_hash, parent_hash, received_at, checked)
	VALUES(6, '0x013be63487a53c874614dd1ae0434cf211e393b2e386c8fde74da203b5469b25', '0x0328698ebeda498df8c63040e2a4771d24722ab2c1e8291226b9215c7eec50f5', '2024-03-11 02:52:28.000', true);
	INSERT INTO state.block
	(block_num, block_hash, parent_hash, received_at, checked)
	VALUES(7, '0x013be63487a53c874614dd1ae0434cf211e393b2e386c8fde74da203b5469b26', '0x0328698ebeda498df8c63040e2a4771d24722ab2c1e8291226b9215c7eec50f6', '2024-03-11 02:52:29.000', true);
	INSERT INTO state.block
	(block_num, block_hash, parent_hash, received_at, checked)
	VALUES(8, '0x013be63487a53c874614dd1ae0434cf211e393b2e386c8fde74da203b5469b27', '0x0328698ebeda498df8c63040e2a4771d24722ab2c1e8291226b9215c7eec50f7', '2024-03-11 02:52:30.000', true);
	INSERT INTO state.block
	(block_num, block_hash, parent_hash, received_at, checked)
	VALUES(9, '0x013be63487a53c874614dd1ae0434cf211e393b2e386c8fde74da203b5469b28', '0x0328698ebeda498df8c63040e2a4771d24722ab2c1e8291226b9215c7eec50f8', '2024-03-11 02:52:31.000', true);
	INSERT INTO state.block
	(block_num, block_hash, parent_hash, received_at, checked)
	VALUES(10, '0x013be63487a53c874614dd1ae0434cf211e393b2e386c8fde74da203b5469b29', '0x0328698ebeda498df8c63040e2a4771d24722ab2c1e8291226b9215c7eec50f9', '2024-03-11 02:52:32.000', true);
	INSERT INTO state.block
	(block_num, block_hash, parent_hash, received_at, checked)
	VALUES(11, '0x013be63487a53c874614dd1ae0434cf211e393b2e386c8fde74da203b5469b2a', '0x0328698ebeda498df8c63040e2a4771d24722ab2c1e8291226b9215c7eec50fa', '2024-03-11 02:52:33.000', true);
	INSERT INTO state.batch
	(batch_num, global_exit_root, local_exit_root, state_root, acc_input_hash, "timestamp", coinbase, raw_txs_data, forced_batch_num, batch_resources, closing_reason, wip, checked)
	VALUES(1, '0x0000000000000000000000000000000000000000000000000000000000000000', '0x0000000000000000000000000000000000000000000000000000000000000000', '0x3f86b09b43e3e49a41fc20a07579b79eba044253367817d5c241d23c0e2bc5c9', '0xa5bd7311fe00707809dd3aa718be2ea0cb363626b9db44172098515f07acf940', '2023-03-24 16:35:27.000', '0x148Ee7dAF16574cD020aFa34CC658f8F3fbd2800', decode('','hex'), NULL, '{"Bytes": 0, "ZKCounters": {"GasUsed": 0, "UsedSteps": 0, "UsedBinaries": 0, "UsedMemAligns": 0, "UsedArithmetics": 0, "UsedKeccakHashes": 0, "UsedPoseidonHashes": 0, "UsedSha256Hashes_V2": 0, "UsedPoseidonPaddings": 0}}'::jsonb, '', false, true);
	INSERT INTO state.virtual_batch
	(batch_num, tx_hash, coinbase, block_num, sequencer_addr, timestamp_batch_etrog, l1_info_root)
	VALUES(1, '0x4314ed5d8ad4812e88895942b2b4642af176d80a97c5489a16a7a5aeb08b51a6', '0x148Ee7dAF16574cD020aFa34CC658f8F3fbd2800', 2, '0x148Ee7dAF16574cD020aFa34CC658f8F3fbd2800', '2024-04-09 16:26:45.000', '0xcdb4258d7ccd8fd41c4a26fd8d9d1fadbc9c506e64d489170525a65e2ad3580b');
	INSERT INTO state.verified_batch
	(batch_num, tx_hash, aggregator, state_root, block_num, is_trusted)
	VALUES(1, '0x28e82f15ab7bac043598623c65a838c315d00ecb5d6e013c406d6bb889680592', '0x6329Fe417621925C81c16F9F9a18c203C21Af7ab', '0x80bd488b1e150b9b42611d038c7fdfa43a3e95b3a02e5c2d57074e73b583f8fd', 3, true);
	INSERT INTO state.fork_id
	(fork_id, from_batch_num, to_batch_num, "version", block_num)
	VALUES(5, 813267, 1228916, 'v2.0.0-RC1-fork.5', 5);
	INSERT INTO state.monitored_txs
	("owner", id, from_addr, to_addr, nonce, value, "data", gas, gas_price, status, history, block_num, created_at, updated_at, gas_offset)
	VALUES('sequencer', 'sequence-from-2006249-to-2006252', '0x148Ee7dAF16574cD020aFa34CC658f8F3fbd2800', '0x519E42c24163192Dca44CD3fBDCEBF6be9130987', 58056, NULL, 'def57e540000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000006614ec3100000000000000000000000000000000000000000000000000000000001e9ce8000000000000000000000000148ee7da0000000300000000ee8306089a84ae0baa0082520894417a7ba2d8d0060ae6c54fd098590db854b9c1d58609184e72a0008082044d80802787e068e6fe23cda64eb868cefb7231a17449d508a77919f6c5408814aaab5f259d43a62eb50df0b2d5740552d3f95176a1f0e31cade590facf70b01c1129151bab0b00000003000000000b00000003000000000b00000003000000000b00000003000000000b0000000300000000000000000000000000000000000000', 1474265, 25212431373, 'done', '{0x44423d538d6fc2f2e882fcd0d1952a735d81c824827b83936e6a5e52268a7d8e}', 7, '2024-04-09 09:26:36.235', '2024-04-09 09:38:24.377', 150000);
	INSERT INTO state.exit_root
	(id, block_num, "timestamp", mainnet_exit_root, rollup_exit_root, global_exit_root, prev_block_hash, l1_info_root, l1_info_tree_index)
	VALUES(379599, 8, '2024-04-09 09:43:59.000', decode('C90DCBC69719971625800AD619E5EEEFD0378317E26F0DDE9B30B3C7C84DBD78','hex'), decode('514D72BBF7C2AD8E4D15EC1186EBF077E98208479651B1C30C5AC7DA11BAB209','hex'), decode('B20FACBED4A2774CE33A0F68D9B6F9B4D9AD553DACD73705503910B141D2102E','hex'), decode('845E01F723E5C77DBE5A4889F299860FBECD8353BFD423D366851F3A90496334','hex'), decode('EDB0EF9C80E947C411FD9B8B23318708132F8A3BD15CD366499866EF91748FC8','hex'), 8032);
	INSERT INTO state.forced_batch
	(block_num, forced_batch_num, global_exit_root, timestamp, raw_txs_data, coinbase)
	VALUES(10, 1, '0x3f86b09b43e3e49a41fc20a07579b79eba044253367817d5c241d23c0e2bc5ca', '2024-04-09 09:26:36.235', '0x3f86b09b', '0x3f86b09b43e3e49a41fc20a07579b79eba044253367817d5c241d23c0e2bc5c9');
	`
	if _, err := db.Exec(addBlocks); err != nil {
		return err
	}
	blockCount := `SELECT count(*) FROM state.block`
	var count int
	err := db.QueryRow(blockCount).Scan(&count)
	if err != nil {
		return err
	}
	if count != 11 {
		return fmt.Errorf("error: initial wrong number of blocks")
	}
	return nil
}

func (m migrationTest0020) RunAssertsAfterMigrationUp(t *testing.T, db *sql.DB) {
	blockCount := `SELECT count(*) FROM state.block`
	var count int
	err := db.QueryRow(blockCount).Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 6, count)
}

func (m migrationTest0020) RunAssertsAfterMigrationDown(t *testing.T, db *sql.DB) {
}

func TestMigration0020(t *testing.T) {
	runMigrationTest(t, 20, migrationTest0020{})
}
