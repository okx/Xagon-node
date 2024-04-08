package state_test

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/merkletree"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/metrics"
	test "github.com/0xPolygonHermez/zkevm-node/state/test/forkid_common"
	"github.com/0xPolygonHermez/zkevm-node/test/testutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	testState *state.State
	// Tests in this file should be independent of the forkID
	// so we force an invalid forkID
	forkID   = uint64(0)
	stateCfg = state.Config{
		MaxCumulativeGasUsed: 800000,
		ChainID:              1000,
		MaxLogsCount:         10000,
		MaxLogsBlockRange:    10000,
		ForkIDIntervals: []state.ForkIDInterval{{
			FromBatchNumber: 0,
			ToBatchNumber:   math.MaxUint64,
			ForkId:          forkID,
			Version:         "",
		}},
	}
)

func TestMain(m *testing.M) {
	testState = test.InitTestState(stateCfg)
	defer test.CloseTestState()
	result := m.Run()
	os.Exit(result)
}

func TestAddBlock(t *testing.T) {
	// Init database instance
	test.InitOrResetDB(test.StateDBCfg)

	ctx := context.Background()
	tx, err := testState.BeginStateTransaction(ctx)
	require.NoError(t, err)
	block := &state.Block{
		BlockNumber: 1,
		BlockHash:   common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"),
		ParentHash:  common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"),
		ReceivedAt:  time.Now(),
	}
	err = testState.AddBlock(ctx, block, tx)
	assert.NoError(t, err)
	// Add the second block
	block.BlockNumber = 2
	err = testState.AddBlock(ctx, block, tx)
	assert.NoError(t, err)
	err = tx.Commit(ctx)
	require.NoError(t, err)
	// Get the last block
	lastBlock, err := testState.GetLastBlock(ctx, nil)
	assert.NoError(t, err)
	assert.Equal(t, uint64(2), lastBlock.BlockNumber)
	assert.Equal(t, block.BlockHash, lastBlock.BlockHash)
	assert.Equal(t, block.ParentHash, lastBlock.ParentHash)
	// Get the previous block
	prevBlock, err := testState.GetPreviousBlock(ctx, 1, nil)
	assert.NoError(t, err)
	assert.Equal(t, uint64(1), prevBlock.BlockNumber)
}

func TestProcessCloseBatch(t *testing.T) {
	// Init database instance
	test.InitOrResetDB(test.StateDBCfg)

	ctx := context.Background()
	dbTx, err := testState.BeginStateTransaction(ctx)
	require.NoError(t, err)
	// Set genesis batch
	_, err = testState.SetGenesis(ctx, state.Block{}, test.Genesis, metrics.SynchronizerCallerLabel, dbTx)
	require.NoError(t, err)
	// Open batch #1
	// processingCtx1 := state.ProcessingContext{
	// 	BatchNumber:    1,
	// 	Coinbase:       common.HexToAddress("1"),
	// 	Timestamp:      time.Now().UTC(),
	// 	globalExitRoot: common.HexToHash("a"),
	// }
	// Txs for batch #1
	// rawTxs := "f84901843b9aca00827b0c945fbdb2315678afecb367f032d93f642f64180aa380a46057361d00000000000000000000000000000000000000000000000000000000000000048203e9808073efe1fa2d3e27f26f32208550ea9b0274d49050b816cadab05a771f4275d0242fd5d92b3fb89575c070e6c930587c520ee65a3aa8cfe382fcad20421bf51d621c"
	// TODO Finish and fix this test
	// err = testState.ProcessAndStoreClosedBatch(ctx, processingCtx1, common.Hex2Bytes(rawTxs), dbTx, state.SynchronizerCallerLabel)
	// require.NoError(t, err)
	require.NoError(t, dbTx.Commit(ctx))
}

// TODO: Review this test
/*
func TestOpenCloseBatch(t *testing.T) {
	var (
		batchResources = state.BatchResources{
			ZKCounters: state.ZKCounters{
				UsedKeccakHashes: 1,
			},
			Bytes: 1,
		}
		closingReason = state.GlobalExitRootDeadlineClosingReason
	)
	// Init database instance
	test.InitOrResetDB(test.StateDBCfg)

	ctx := context.Background()
	dbTx, err := testState.BeginStateTransaction(ctx)
	require.NoError(t, err)
	// Set genesis batch
	_, err = testState.SetGenesis(ctx, state.Block{}, test.Genesis, metrics.SynchronizerCallerLabel, dbTx)
	require.NoError(t, err)
	// Open batch #1
	processingCtx1 := state.ProcessingContext{
		BatchNumber:    1,
		Coinbase:       common.HexToAddress("1"),
		Timestamp:      time.Now().UTC(),
		GlobalExitRoot: common.HexToHash("a"),
	}
	err = testState.OpenBatch(ctx, processingCtx1, dbTx)
	require.NoError(t, err)
	require.NoError(t, dbTx.Commit(ctx))
	dbTx, err = testState.BeginStateTransaction(ctx)
	require.NoError(t, err)
	// Fail opening batch #2 (#1 is still open)
	processingCtx2 := state.ProcessingContext{
		BatchNumber:    2,
		Coinbase:       common.HexToAddress("2"),
		Timestamp:      time.Now().UTC(),
		GlobalExitRoot: common.HexToHash("b"),
	}
	err = testState.OpenBatch(ctx, processingCtx2, dbTx)
	assert.Equal(t, state.ErrLastBatchShouldBeClosed, err)
	// Fail closing batch #1 (it has no txs yet)
	receipt1 := state.ProcessingReceipt{
		BatchNumber:    1,
		StateRoot:      common.HexToHash("1"),
		LocalExitRoot:  common.HexToHash("1"),
		ClosingReason:  closingReason,
		BatchResources: batchResources,
	}
	err = testState.CloseBatch(ctx, receipt1, dbTx)
	require.NoError(t, err)
	require.NoError(t, dbTx.Rollback(ctx))
	dbTx, err = testState.BeginStateTransaction(ctx)
	require.NoError(t, err)
	// Add txs to batch #1
	tx1 := *types.NewTransaction(0, common.HexToAddress("0"), big.NewInt(0), 0, big.NewInt(0), []byte("aaa"))
	tx2 := *types.NewTransaction(1, common.HexToAddress("1"), big.NewInt(1), 0, big.NewInt(1), []byte("bbb"))
	txsBatch1 := []*state.ProcessTransactionResponse{
		{
			TxHash: tx1.Hash(),
			Tx:     tx1,
		},
		{
			TxHash: tx2.Hash(),
			Tx:     tx2,
		},
	}
	block1 := []*state.ProcessBlockResponse{
		{
			TransactionResponses: txsBatch1,
		},
	}

	data, err := state.EncodeTransactions([]types.Transaction{tx1, tx2}, constants.TwoEffectivePercentages, forkID)
	require.NoError(t, err)
	receipt1.BatchL2Data = data

	err = testState.StoreTransactions(ctx, 1, block1, nil, dbTx)
	require.NoError(t, err)
	// Close batch #1
	err = testState.CloseBatch(ctx, receipt1, dbTx)
	require.NoError(t, err)
	require.NoError(t, dbTx.Commit(ctx))
	dbTx, err = testState.BeginStateTransaction(ctx)
	require.NoError(t, err)
	// Fail opening batch #3 (should open batch #2)
	processingCtx3 := state.ProcessingContext{
		BatchNumber:    3,
		Coinbase:       common.HexToAddress("3"),
		Timestamp:      time.Now().UTC(),
		GlobalExitRoot: common.HexToHash("c"),
	}
	err = testState.OpenBatch(ctx, processingCtx3, dbTx)
	require.ErrorIs(t, err, state.ErrUnexpectedBatch)
	// Fail opening batch #2 (invalid timestamp)
	processingCtx2.Timestamp = processingCtx1.Timestamp.Add(-1 * time.Second)
	err = testState.OpenBatch(ctx, processingCtx2, dbTx)
	require.Equal(t, state.ErrTimestampGE, err)
	processingCtx2.Timestamp = time.Now()
	require.NoError(t, dbTx.Rollback(ctx))
	dbTx, err = testState.BeginStateTransaction(ctx)
	require.NoError(t, err)
	// Open batch #2
	err = testState.OpenBatch(ctx, processingCtx2, dbTx)
	require.NoError(t, err)
	// Get batch #2 from DB and compare with on memory batch
	actualBatch, err := testState.GetBatchByNumber(ctx, 1, dbTx)
	require.NoError(t, err)
	batchL2Data, err := state.EncodeTransactions([]types.Transaction{tx1, tx2}, constants.TwoEffectivePercentages, forkID)
	require.NoError(t, err)
	assertBatch(t, state.Batch{
		BatchNumber:    1,
		Coinbase:       processingCtx1.Coinbase,
		BatchL2Data:    batchL2Data,
		StateRoot:      receipt1.StateRoot,
		LocalExitRoot:  receipt1.LocalExitRoot,
		Timestamp:      processingCtx1.Timestamp,
		GlobalExitRoot: processingCtx1.GlobalExitRoot,
	}, *actualBatch)
	require.NoError(t, dbTx.Commit(ctx))
}

func assertBatch(t *testing.T, expected, actual state.Batch) {
	assert.Equal(t, expected.Timestamp.Unix(), actual.Timestamp.Unix())
	actual.Timestamp = expected.Timestamp
	assert.Equal(t, expected, actual)
}
*/
func TestAddForcedBatch(t *testing.T) {
	// Init database instance
	test.InitOrResetDB(test.StateDBCfg)

	ctx := context.Background()
	tx, err := testState.BeginStateTransaction(ctx)
	require.NoError(t, err)
	block := &state.Block{
		BlockNumber: 1,
		BlockHash:   common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"),
		ParentHash:  common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"),
		ReceivedAt:  time.Now(),
	}
	err = testState.AddBlock(ctx, block, tx)
	assert.NoError(t, err)
	b := common.Hex2Bytes("0x617b3a3528F9")
	assert.NoError(t, err)
	forcedBatch := state.ForcedBatch{
		BlockNumber:       1,
		ForcedBatchNumber: 2,
		GlobalExitRoot:    common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"),
		Sequencer:         common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D"),
		RawTxsData:        b,
		ForcedAt:          time.Now(),
	}
	err = testState.AddForcedBatch(ctx, &forcedBatch, tx)
	require.NoError(t, err)
	fb, err := testState.GetForcedBatch(ctx, 2, tx)
	require.NoError(t, err)
	err = tx.Commit(ctx)
	require.NoError(t, err)
	assert.Equal(t, forcedBatch.BlockNumber, fb.BlockNumber)
	assert.Equal(t, forcedBatch.ForcedBatchNumber, fb.ForcedBatchNumber)
	assert.NotEqual(t, time.Time{}, fb.ForcedAt)
	assert.Equal(t, forcedBatch.GlobalExitRoot, fb.GlobalExitRoot)
	assert.Equal(t, forcedBatch.RawTxsData, fb.RawTxsData)
	// Test GetNextForcedBatches
	tx, err = testState.BeginStateTransaction(ctx)
	require.NoError(t, err)
	forcedBatch = state.ForcedBatch{
		BlockNumber:       1,
		ForcedBatchNumber: 3,
		GlobalExitRoot:    common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"),
		Sequencer:         common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D"),
		RawTxsData:        b,
		ForcedAt:          time.Now(),
	}
	err = testState.AddForcedBatch(ctx, &forcedBatch, tx)
	require.NoError(t, err)

	_, err = testState.Exec(ctx, "INSERT INTO state.batch (batch_num, forced_batch_num, WIP) VALUES (2, 2, FALSE)")
	assert.NoError(t, err)
	virtualBatch := state.VirtualBatch{
		BlockNumber: 1,
		BatchNumber: 2,
		TxHash:      common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"),
		Coinbase:    common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D"),
	}
	err = testState.AddVirtualBatch(ctx, &virtualBatch, tx)
	require.NoError(t, err)

	batches, err := testState.GetNextForcedBatches(ctx, 1, tx)
	require.NoError(t, err)
	assert.Equal(t, forcedBatch.BlockNumber, batches[0].BlockNumber)
	assert.Equal(t, forcedBatch.ForcedBatchNumber, batches[0].ForcedBatchNumber)
	assert.NotEqual(t, time.Time{}, batches[0].ForcedAt)
	assert.Equal(t, forcedBatch.GlobalExitRoot, batches[0].GlobalExitRoot)
	assert.Equal(t, forcedBatch.RawTxsData, batches[0].RawTxsData)
	require.NoError(t, tx.Commit(ctx))
}

func TestAddVirtualBatch(t *testing.T) {
	// Init database instance
	test.InitOrResetDB(test.StateDBCfg)

	ctx := context.Background()
	tx, err := testState.BeginStateTransaction(ctx)
	require.NoError(t, err)
	block := &state.Block{
		BlockNumber: 1,
		BlockHash:   common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"),
		ParentHash:  common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"),
		ReceivedAt:  time.Now(),
	}
	err = testState.AddBlock(ctx, block, tx)
	assert.NoError(t, err)
	_, err = testState.Exec(ctx, "INSERT INTO state.batch (batch_num, WIP) VALUES (1, FALSE)")
	assert.NoError(t, err)
	virtualBatch := state.VirtualBatch{
		BlockNumber: 1,
		BatchNumber: 1,
		TxHash:      common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"),
		Coinbase:    common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D"),
	}
	err = testState.AddVirtualBatch(ctx, &virtualBatch, tx)
	require.NoError(t, err)
	require.NoError(t, tx.Commit(ctx))
}

func TestGetTxsHashesToDelete(t *testing.T) {
	test.InitOrResetDB(test.StateDBCfg)

	ctx := context.Background()
	tx, err := testState.BeginStateTransaction(ctx)
	require.NoError(t, err)
	block1 := &state.Block{
		BlockNumber: 1,
		BlockHash:   common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"),
		ParentHash:  common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"),
		ReceivedAt:  time.Now(),
	}
	err = testState.AddBlock(ctx, block1, tx)
	assert.NoError(t, err)
	block2 := &state.Block{
		BlockNumber: 2,
		BlockHash:   common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"),
		ParentHash:  common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"),
		ReceivedAt:  time.Now(),
	}
	err = testState.AddBlock(ctx, block2, tx)
	assert.NoError(t, err)

	_, err = testState.Exec(ctx, "INSERT INTO state.batch (batch_num, WIP) VALUES (1, FALSE)")
	assert.NoError(t, err)
	require.NoError(t, err)
	virtualBatch1 := state.VirtualBatch{
		BlockNumber: 1,
		BatchNumber: 1,
		TxHash:      common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"),
		Coinbase:    common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D"),
	}

	_, err = testState.Exec(ctx, "INSERT INTO state.batch (batch_num, WIP) VALUES (2, FALSE)")
	assert.NoError(t, err)
	virtualBatch2 := state.VirtualBatch{
		BlockNumber: 1,
		BatchNumber: 2,
		TxHash:      common.HexToHash("0x132"),
		Coinbase:    common.HexToAddress("0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D"),
	}
	err = testState.AddVirtualBatch(ctx, &virtualBatch1, tx)
	require.NoError(t, err)
	err = testState.AddVirtualBatch(ctx, &virtualBatch2, tx)
	require.NoError(t, err)
	require.NoError(t, tx.Commit(ctx))

	_, err = testState.Exec(ctx, "INSERT INTO state.l2block (block_num, block_hash, received_at, batch_num, created_at) VALUES ($1, $2, $3, $4, $5)", 1, "0x423", time.Now(), 1, time.Now().UTC())
	require.NoError(t, err)
	l2Tx1 := types.NewTransaction(1, common.Address{}, big.NewInt(10), 21000, big.NewInt(1), []byte{})
	_, err = testState.Exec(ctx, "INSERT INTO state.transaction (l2_block_num, encoded, hash) VALUES ($1, $2, $3)",
		virtualBatch1.BatchNumber, fmt.Sprintf("encoded-%d", virtualBatch1.BatchNumber), l2Tx1.Hash().Hex())
	require.NoError(t, err)

	_, err = testState.Exec(ctx, "INSERT INTO state.l2block (block_num, block_hash, received_at, batch_num, created_at) VALUES ($1, $2, $3, $4, $5)", 2, "0x423", time.Now(), 2, time.Now().UTC())
	require.NoError(t, err)
	l2Tx2 := types.NewTransaction(2, common.Address{}, big.NewInt(10), 21000, big.NewInt(1), []byte{})
	_, err = testState.Exec(ctx, "INSERT INTO state.transaction (l2_block_num, encoded, hash) VALUES ($1, $2, $3)",
		virtualBatch2.BatchNumber, fmt.Sprintf("encoded-%d", virtualBatch2.BatchNumber), l2Tx2.Hash().Hex())
	require.NoError(t, err)
	txHashes, err := testState.GetTxsOlderThanNL1Blocks(ctx, 1, nil)
	require.NoError(t, err)
	require.Equal(t, l2Tx1.Hash().Hex(), txHashes[0].Hex())
}

func TestCheckSupersetBatchTransactions(t *testing.T) {
	tcs := []struct {
		description      string
		existingTxHashes []common.Hash
		processedTxs     []*state.ProcessTransactionResponse
		expectedError    bool
		expectedErrorMsg string
	}{
		{
			description:      "empty existingTxHashes and processedTx is successful",
			existingTxHashes: []common.Hash{},
			processedTxs:     []*state.ProcessTransactionResponse{},
		},
		{
			description: "happy path",
			existingTxHashes: []common.Hash{
				common.HexToHash("0x8a84686634729c57532b9ffa4e632e241b2de5c880c771c5c214d5e7ec465b1c"),
				common.HexToHash("0x30c6a361ba88906ef2085d05a2aeac15e793caff2bdc1deaaae2f4910d83de52"),
				common.HexToHash("0x0d3453b6d17841b541d4f79f78d5fa22fff281551ed4012c7590b560b2969e7f"),
			},
			processedTxs: []*state.ProcessTransactionResponse{
				{TxHash: common.HexToHash("0x8a84686634729c57532b9ffa4e632e241b2de5c880c771c5c214d5e7ec465b1c")},
				{TxHash: common.HexToHash("0x30c6a361ba88906ef2085d05a2aeac15e793caff2bdc1deaaae2f4910d83de52")},
				{TxHash: common.HexToHash("0x0d3453b6d17841b541d4f79f78d5fa22fff281551ed4012c7590b560b2969e7f")},
			},
		},
		{
			description:      "existingTxHashes bigger than processedTx gives error",
			existingTxHashes: []common.Hash{common.HexToHash(""), common.HexToHash("")},
			processedTxs:     []*state.ProcessTransactionResponse{{}},
			expectedError:    true,
			expectedErrorMsg: state.ErrExistingTxGreaterThanProcessedTx.Error(),
		},
		{
			description: "processedTx not present in existingTxHashes gives error",
			existingTxHashes: []common.Hash{
				common.HexToHash("0x8a84686634729c57532b9ffa4e632e241b2de5c880c771c5c214d5e7ec465b1c"),
				common.HexToHash("0x30c6a361ba88906ef2085d05a2aeac15e793caff2bdc1deaaae2f4910d83de52"),
			},
			processedTxs: []*state.ProcessTransactionResponse{
				{TxHash: common.HexToHash("0x8a84686634729c57532b9ffa4e632e241b2de5c880c771c5c214d5e7ec465b1c")},
				{TxHash: common.HexToHash("0x0d3453b6d17841b541d4f79f78d5fa22fff281551ed4012c7590b560b2969e7f")},
			},
			expectedError:    true,
			expectedErrorMsg: state.ErrOutOfOrderProcessedTx.Error(),
		},
		{
			description: "out of order processedTx gives error",
			existingTxHashes: []common.Hash{
				common.HexToHash("0x8a84686634729c57532b9ffa4e632e241b2de5c880c771c5c214d5e7ec465b1c"),
				common.HexToHash("0x30c6a361ba88906ef2085d05a2aeac15e793caff2bdc1deaaae2f4910d83de52"),
				common.HexToHash("0x0d3453b6d17841b541d4f79f78d5fa22fff281551ed4012c7590b560b2969e7f"),
			},
			processedTxs: []*state.ProcessTransactionResponse{
				{TxHash: common.HexToHash("0x8a84686634729c57532b9ffa4e632e241b2de5c880c771c5c214d5e7ec465b1c")},
				{TxHash: common.HexToHash("0x0d3453b6d17841b541d4f79f78d5fa22fff281551ed4012c7590b560b2969e7f")},
				{TxHash: common.HexToHash("0x30c6a361ba88906ef2085d05a2aeac15e793caff2bdc1deaaae2f4910d83de52")},
			},
			expectedError:    true,
			expectedErrorMsg: state.ErrOutOfOrderProcessedTx.Error(),
		},
	}
	for _, tc := range tcs {
		// tc := tc
		t.Run(tc.description, func(t *testing.T) {
			require.NoError(t, testutils.CheckError(
				state.CheckSupersetBatchTransactions(tc.existingTxHashes, tc.processedTxs),
				tc.expectedError,
				tc.expectedErrorMsg,
			))
		})
	}
}

func TestGetTxsHashesByBatchNumber(t *testing.T) {
	// Init database instance
	test.InitOrResetDB(test.StateDBCfg)

	ctx := context.Background()
	dbTx, err := testState.BeginStateTransaction(ctx)
	require.NoError(t, err)
	// Set genesis batch
	_, err = testState.SetGenesis(ctx, state.Block{}, test.Genesis, metrics.SynchronizerCallerLabel, dbTx)
	require.NoError(t, err)
	// Open batch #1
	processingCtx1 := state.ProcessingContext{
		BatchNumber:    1,
		Coinbase:       common.HexToAddress("1"),
		Timestamp:      time.Now().UTC(),
		GlobalExitRoot: common.HexToHash("a"),
	}
	err = testState.OpenBatch(ctx, processingCtx1, dbTx)
	require.NoError(t, err)

	// Add txs to batch #1
	tx1 := *types.NewTransaction(0, common.HexToAddress("0"), big.NewInt(0), 0, big.NewInt(0), []byte("aaa"))
	tx2 := *types.NewTransaction(1, common.HexToAddress("1"), big.NewInt(1), 0, big.NewInt(1), []byte("bbb"))
	txsBatch1 := []*state.ProcessTransactionResponse{
		{
			TxHash: tx1.Hash(),
			Tx:     tx1,
		},
		{
			TxHash: tx2.Hash(),
			Tx:     tx2,
		},
	}
	block1 := []*state.ProcessBlockResponse{
		{
			TransactionResponses: txsBatch1,
		},
	}

	err = testState.StoreTransactions(ctx, 1, block1, nil, dbTx)
	require.NoError(t, err)

	txs, err := testState.GetTxsHashesByBatchNumber(ctx, 1, dbTx)
	require.NoError(t, err)

	require.Equal(t, len(txsBatch1), len(txs))
	for i := range txsBatch1 {
		require.Equal(t, txsBatch1[i].TxHash, txs[i])
	}
	require.NoError(t, dbTx.Commit(ctx))
}

func TestGenesisNewLeafType(t *testing.T) {
	ctx := context.Background()
	// Set Genesis
	block := state.Block{
		BlockNumber: 0,
		BlockHash:   state.ZeroHash,
		ParentHash:  state.ZeroHash,
		ReceivedAt:  time.Now(),
	}

	test.Genesis.Actions = []*state.GenesisAction{
		{
			Address: "0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D",
			Type:    int(merkletree.LeafTypeBalance),
			Value:   "100000000000000000000",
		},
		{
			Address: "0x617b3a3528F9cDd6630fd3301B9c8911F7Bf063D",
			Type:    int(merkletree.LeafTypeNonce),
			Value:   "0",
		},
		{
			Address: "0x4d5Cf5032B2a844602278b01199ED191A86c93ff",
			Type:    int(merkletree.LeafTypeBalance),
			Value:   "200000000000000000000",
		},
		{
			Address: "0x4d5Cf5032B2a844602278b01199ED191A86c93ff",
			Type:    int(merkletree.LeafTypeNonce),
			Value:   "0",
		},
		{
			Address: "0x03e75d7dd38cce2e20ffee35ec914c57780a8e29",
			Type:    int(merkletree.LeafTypeBalance),
			Value:   "0",
		},
		{
			Address: "0x03e75d7dd38cce2e20ffee35ec914c57780a8e29",
			Type:    int(merkletree.LeafTypeNonce),
			Value:   "0",
		},
		{
			Address:  "0x03e75d7dd38cce2e20ffee35ec914c57780a8e29",
			Type:     int(merkletree.LeafTypeCode),
			Bytecode: "60606040525b600080fd00a165627a7a7230582012c9bd00152fa1c480f6827f81515bb19c3e63bf7ed9ffbb5fda0265983ac7980029",
		},
	}

	test.InitOrResetDB(test.StateDBCfg)
	dbTx, err := testState.BeginStateTransaction(ctx)
	require.NoError(t, err)
	stateRoot, err := testState.SetGenesis(ctx, block, test.Genesis, metrics.SynchronizerCallerLabel, dbTx)
	require.NoError(t, err)
	require.NoError(t, dbTx.Commit(ctx))

	log.Debug(string(stateRoot.Bytes()))
	log.Debug(common.BytesToHash(stateRoot.Bytes()))
	log.Debug(common.BytesToHash(stateRoot.Bytes()).String())
	log.Debug(new(big.Int).SetBytes(stateRoot.Bytes()))
	log.Debug(common.Bytes2Hex(stateRoot.Bytes()))

	require.Equal(t, "49461512068930131501252998918674096186707801477301326632372959001738876161218", new(big.Int).SetBytes(stateRoot.Bytes()).String())
}

func TestAddGetL2Block(t *testing.T) {
	// Init database instance
	test.InitOrResetDB(test.StateDBCfg)

	ctx := context.Background()
	dbTx, err := testState.BeginStateTransaction(ctx)
	require.NoError(t, err)
	block := &state.Block{
		BlockNumber: 1,
		BlockHash:   common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"),
		ParentHash:  common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1"),
		ReceivedAt:  time.Now(),
	}
	err = testState.AddBlock(ctx, block, dbTx)
	assert.NoError(t, err)

	batchNumber := uint64(1)
	_, err = testState.Exec(ctx, "INSERT INTO state.batch (batch_num, WIP) VALUES ($1, FALSE)", batchNumber)
	assert.NoError(t, err)

	time := time.Now()
	blockNumber := big.NewInt(1)

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    0,
		To:       nil,
		Value:    new(big.Int),
		Gas:      0,
		GasPrice: big.NewInt(0),
	})

	receipt := &types.Receipt{
		Type:              tx.Type(),
		PostState:         state.ZeroHash.Bytes(),
		CumulativeGasUsed: 0,
		BlockNumber:       blockNumber,
		GasUsed:           tx.Gas(),
		TxHash:            tx.Hash(),
		TransactionIndex:  0,
		Status:            types.ReceiptStatusSuccessful,
	}

	header := state.NewL2Header(&types.Header{
		Number:     big.NewInt(1),
		ParentHash: state.ZeroHash,
		Coinbase:   state.ZeroAddress,
		Root:       state.ZeroHash,
		GasUsed:    1,
		GasLimit:   10,
		Time:       uint64(time.Unix()),
	})
	transactions := []*types.Transaction{tx}

	receipts := []*types.Receipt{receipt}
	imStateRoots := []common.Hash{state.ZeroHash}

	// Create block to be able to calculate its hash
	st := trie.NewStackTrie(nil)
	l2Block := state.NewL2Block(header, transactions, []*state.L2Header{}, receipts, st)
	l2Block.ReceivedAt = time

	receipt.BlockHash = l2Block.Hash()

	numTxs := len(transactions)
	storeTxsEGPData := make([]state.StoreTxEGPData, numTxs)
	txsL2Hash := make([]common.Hash, numTxs)
	for i := range transactions {
		storeTxsEGPData[i] = state.StoreTxEGPData{EGPLog: nil, EffectivePercentage: state.MaxEffectivePercentage}
		txsL2Hash[i] = common.HexToHash(fmt.Sprintf("0x%d", i))
	}

	err = testState.AddL2Block(ctx, batchNumber, l2Block, receipts, txsL2Hash, storeTxsEGPData, imStateRoots, dbTx)
	require.NoError(t, err)
	result, err := testState.GetL2BlockByHash(ctx, l2Block.Hash(), dbTx)
	require.NoError(t, err)

	assert.Equal(t, l2Block.Hash(), result.Hash())

	result, err = testState.GetL2BlockByNumber(ctx, l2Block.NumberU64(), dbTx)
	require.NoError(t, err)

	assert.Equal(t, l2Block.Hash(), result.Hash())
	assert.Equal(t, l2Block.ReceivedAt.Unix(), result.ReceivedAt.Unix())
	assert.Equal(t, l2Block.Time(), result.Time())

	require.NoError(t, dbTx.Commit(ctx))
}

func TestGenesis(t *testing.T) {
	ctx := context.Background()
	block := state.Block{
		BlockNumber: 1,
		BlockHash:   state.ZeroHash,
		ParentHash:  state.ZeroHash,
		ReceivedAt:  time.Now(),
	}

	actions := []*state.GenesisAction{
		{
			Address: "0xb1D0Dc8E2Ce3a93EB2b32f4C7c3fD9dDAf1211FA",
			Type:    int(merkletree.LeafTypeBalance),
			Value:   "1000",
		},
		{
			Address: "0xb1D0Dc8E2Ce3a93EB2b32f4C7c3fD9dDAf1211FB",
			Type:    int(merkletree.LeafTypeBalance),
			Value:   "2000",
		},
		{
			Address: "0xb1D0Dc8E2Ce3a93EB2b32f4C7c3fD9dDAf1211FA",
			Type:    int(merkletree.LeafTypeNonce),
			Value:   "1",
		},
		{
			Address: "0xb1D0Dc8E2Ce3a93EB2b32f4C7c3fD9dDAf1211FB",
			Type:    int(merkletree.LeafTypeNonce),
			Value:   "1",
		},
		{
			Address:  "0xae4bb80be56b819606589de61d5ec3b522eeb032",
			Type:     int(merkletree.LeafTypeCode),
			Bytecode: "608060405234801561001057600080fd5b50600436106100675760003560e01c806333d6247d1161005057806333d6247d146100a85780633ed691ef146100bd578063a3c573eb146100d257600080fd5b806301fd90441461006c5780633381fe9014610088575b600080fd5b61007560015481565b6040519081526020015b60405180910390f35b6100756100963660046101c7565b60006020819052908152604090205481565b6100bb6100b63660046101c7565b610117565b005b43600090815260208190526040902054610075565b6002546100f29073ffffffffffffffffffffffffffffffffffffffff1681565b60405173ffffffffffffffffffffffffffffffffffffffff909116815260200161007f565b60025473ffffffffffffffffffffffffffffffffffffffff1633146101c2576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152603460248201527f476c6f62616c45786974526f6f744d616e616765724c323a3a7570646174654560448201527f786974526f6f743a204f4e4c595f425249444745000000000000000000000000606482015260840160405180910390fd5b600155565b6000602082840312156101d957600080fd5b503591905056fea2646970667358221220d6ed73b81f538d38669b0b750b93be08ca365978fae900eedc9ca93131c97ca664736f6c63430008090033",
		},
		{
			Address:         "0xae4bb80be56b819606589de61d5ec3b522eeb032",
			Type:            int(merkletree.LeafTypeStorage),
			StoragePosition: "0x0000000000000000000000000000000000000000000000000000000000000002",
			Value:           "0x9d98deabc42dd696deb9e40b4f1cab7ddbf55988",
		},
	}

	test.InitOrResetDB(test.StateDBCfg)

	dbTx, err := testState.BeginStateTransaction(ctx)
	require.NoError(t, err)

	test.Genesis.Actions = actions
	stateRoot, err := testState.SetGenesis(ctx, block, test.Genesis, metrics.SynchronizerCallerLabel, dbTx)
	require.NoError(t, err)
	require.NoError(t, dbTx.Commit(ctx))

	stateTree := testState.GetTree()

	for _, action := range actions {
		address := common.HexToAddress(action.Address)
		switch action.Type {
		case int(merkletree.LeafTypeBalance):
			balance, err := stateTree.GetBalance(ctx, address, stateRoot.Bytes())
			require.NoError(t, err)
			require.Equal(t, action.Value, balance.String())
		case int(merkletree.LeafTypeNonce):
			nonce, err := stateTree.GetNonce(ctx, address, stateRoot.Bytes())
			require.NoError(t, err)
			require.Equal(t, action.Value, nonce.String())
		case int(merkletree.LeafTypeCode):
			sc, err := stateTree.GetCode(ctx, address, stateRoot.Bytes())
			require.NoError(t, err)
			require.Equal(t, common.Hex2Bytes(action.Bytecode), sc)
		case int(merkletree.LeafTypeStorage):
			st, err := stateTree.GetStorageAt(ctx, address, new(big.Int).SetBytes(common.Hex2Bytes(action.StoragePosition)), stateRoot.Bytes())
			require.NoError(t, err)
			require.Equal(t, new(big.Int).SetBytes(common.Hex2Bytes(action.Value)), st)
		}
	}

	err = testState.GetTree().Flush(ctx, stateRoot, "")
	require.NoError(t, err)
}

func TestGetForkIDforGenesisBatch(t *testing.T) {
	type testCase struct {
		name           string
		cfg            state.Config
		expectedForkID uint64
	}

	testCases := []testCase{
		{
			name: "fork ID for batch 0 is defined",
			cfg: state.Config{
				ForkIDIntervals: []state.ForkIDInterval{
					{ForkId: 2, FromBatchNumber: 0, ToBatchNumber: 10},
					{ForkId: 4, FromBatchNumber: 11, ToBatchNumber: 20},
					{ForkId: 6, FromBatchNumber: 21, ToBatchNumber: math.MaxUint64},
				},
			},
			expectedForkID: 2,
		},
		{
			name: "fork ID for batch 0 is NOT defined",
			cfg: state.Config{
				ForkIDIntervals: []state.ForkIDInterval{
					{ForkId: 7, FromBatchNumber: 1, ToBatchNumber: 10},
					{ForkId: 8, FromBatchNumber: 11, ToBatchNumber: 20},
					{ForkId: 9, FromBatchNumber: 21, ToBatchNumber: math.MaxUint64},
				},
			},
			expectedForkID: 7,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			test.InitOrResetDB(test.StateDBCfg)

			st := test.InitTestState(testCase.cfg)

			forkID := st.GetForkIDByBatchNumber(0)
			assert.Equal(t, testCase.expectedForkID, forkID)

			test.CloseTestState()
		})
	}
}
