package state_test

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/metrics"
	"github.com/0xPolygonHermez/zkevm-node/state/mocks"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	addr1 = common.HexToAddress("0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266")
	hash1 = common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1")
	hash2 = common.HexToHash("0x979b141b8bcd3ba17815cd76811f1fca1cabaa9d51f7c00712606970f81d6e37")
	hash3 = common.HexToHash("3276a200a5fb45f69a4964484d6e677aefaa820924d0896e3ad1ccacfc0971ff")
	hash4 = common.HexToHash("157cd228e43abd9c0f655e08066809106b914be67dacb6efa28a24203a68b1c4")
	hash5 = common.HexToHash("33027547537d35728a741470df1ccf65de10b454ca0def7c5c20b257b7b8d161")
	time1 = time.Unix(1610000000, 0)
	time2 = time.Unix(1620000000, 0)
	data1 = []byte("data1")
)

func TestProcessAndStoreClosedBatchV2(t *testing.T) {
	stateCfg := state.Config{
		MaxCumulativeGasUsed: 800000,
		ChainID:              1000,
		MaxLogsCount:         10000,
		MaxLogsBlockRange:    10000,
		ForkIDIntervals: []state.ForkIDInterval{{
			FromBatchNumber: 0,
			ToBatchNumber:   math.MaxUint64,
			ForkId:          state.FORKID_ETROG,
			Version:         "",
		}},
	}

	ctx := context.Background()
	mockStorage := mocks.NewStorageMock(t)
	mockExecutor := mocks.NewExecutorServiceClientMock(t)
	testState := state.NewState(stateCfg, mockStorage, mockExecutor, nil, nil, nil, nil)
	mockStorage.EXPECT().Begin(ctx).Return(mocks.NewDbTxMock(t), nil)
	dbTx, err := testState.BeginStateTransaction(ctx)
	require.NoError(t, err)

	processingCtx := state.ProcessingContextV2{
		BatchNumber:    128,
		Coinbase:       addr1,
		Timestamp:      &time2,
		L1InfoRoot:     hash1,
		BatchL2Data:    &data1,
		GlobalExitRoot: hash2,
	}
	batchContext := state.ProcessingContext{
		BatchNumber:    processingCtx.BatchNumber,
		Coinbase:       processingCtx.Coinbase,
		Timestamp:      *processingCtx.Timestamp,
		GlobalExitRoot: processingCtx.GlobalExitRoot,
		ForcedBatchNum: processingCtx.ForcedBatchNum,
		BatchL2Data:    processingCtx.BatchL2Data,
	}
	latestBatch := state.Batch{
		BatchNumber: 128,
	}
	previousBatch := state.Batch{
		BatchNumber: 127,
	}

	executorResponse := executor.ProcessBatchResponseV2{
		Error:            executor.ExecutorError_EXECUTOR_ERROR_NO_ERROR,
		ErrorRom:         executor.RomError_ROM_ERROR_NO_ERROR,
		NewStateRoot:     hash3.Bytes(),
		NewLocalExitRoot: hash4.Bytes(),
		NewAccInputHash:  hash5.Bytes(),
	}
	// IMPORTANT: GlobalExitRoot is not stored in the close call
	closingReceipt := state.ProcessingReceipt{
		BatchNumber:   processingCtx.BatchNumber,
		StateRoot:     hash3,
		LocalExitRoot: hash4,
		AccInputHash:  hash5,
		BatchL2Data:   *processingCtx.BatchL2Data,
	}
	// Call the function under test
	mockStorage.EXPECT().GetLastBatchNumber(ctx, dbTx).Return(uint64(127), nil)
	mockStorage.EXPECT().IsBatchClosed(ctx, uint64(127), dbTx).Return(true, nil)
	mockStorage.EXPECT().GetLastBatchTime(ctx, dbTx).Return(time1, nil)
	// When calls to OpenBatch doesnt store the BatchL2Data yet
	batchContext.BatchL2Data = nil
	mockStorage.EXPECT().OpenBatchInStorage(ctx, batchContext, dbTx).Return(nil)
	mockStorage.EXPECT().GetLastNBatches(ctx, uint(2), dbTx).Return([]*state.Batch{&latestBatch, &previousBatch}, nil)
	mockStorage.EXPECT().IsBatchClosed(ctx, uint64(128), dbTx).Return(false, nil)
	mockStorage.EXPECT().GetForkIDByBatchNumber(uint64(128)).Return(uint64(state.FORKID_ETROG))
	mockExecutor.EXPECT().ProcessBatchV2(ctx, mock.Anything, mock.Anything).Return(&executorResponse, nil)
	mockStorage.EXPECT().CloseBatchInStorage(ctx, closingReceipt, dbTx).Return(nil)
	_, _, _, err = testState.ProcessAndStoreClosedBatchV2(ctx, processingCtx, dbTx, metrics.CallerLabel("test"))
	require.NoError(t, err)

	// Add assertions as needed
}

func TestProcessAndStoreClosedBatchV2ErrorOOC(t *testing.T) {
	stateCfg := state.Config{
		MaxCumulativeGasUsed: 800000,
		ChainID:              1000,
		MaxLogsCount:         10000,
		MaxLogsBlockRange:    10000,
		ForkIDIntervals: []state.ForkIDInterval{{
			FromBatchNumber: 0,
			ToBatchNumber:   math.MaxUint64,
			ForkId:          state.FORKID_ETROG,
			Version:         "",
		}},
	}

	ctx := context.Background()
	mockStorage := mocks.NewStorageMock(t)
	mockExecutor := mocks.NewExecutorServiceClientMock(t)
	testState := state.NewState(stateCfg, mockStorage, mockExecutor, nil, nil, nil, nil)
	mockStorage.EXPECT().Begin(ctx).Return(mocks.NewDbTxMock(t), nil)
	dbTx, err := testState.BeginStateTransaction(ctx)
	require.NoError(t, err)

	processingCtx := state.ProcessingContextV2{
		BatchNumber:    128,
		Coinbase:       addr1,
		Timestamp:      &time2,
		L1InfoRoot:     hash1,
		BatchL2Data:    &data1,
		GlobalExitRoot: hash2,
	}
	batchContext := state.ProcessingContext{
		BatchNumber:    processingCtx.BatchNumber,
		Coinbase:       processingCtx.Coinbase,
		Timestamp:      *processingCtx.Timestamp,
		GlobalExitRoot: processingCtx.GlobalExitRoot,
		ForcedBatchNum: processingCtx.ForcedBatchNum,
		BatchL2Data:    processingCtx.BatchL2Data,
	}
	latestBatch := state.Batch{
		BatchNumber: 128,
	}
	previousBatch := state.Batch{
		BatchNumber: 127,
	}

	executorResponse := executor.ProcessBatchResponseV2{
		Error:            executor.ExecutorError_EXECUTOR_ERROR_NO_ERROR,
		ErrorRom:         executor.RomError_ROM_ERROR_OUT_OF_COUNTERS_KECCAK,
		NewStateRoot:     hash3.Bytes(),
		NewLocalExitRoot: hash4.Bytes(),
		NewAccInputHash:  hash5.Bytes(),
	}
	// IMPORTANT: GlobalExitRoot is not stored in the close call
	closingReceipt := state.ProcessingReceipt{
		BatchNumber:   processingCtx.BatchNumber,
		StateRoot:     hash3,
		LocalExitRoot: hash4,
		AccInputHash:  hash5,
		BatchL2Data:   *processingCtx.BatchL2Data,
	}
	// Call the function under test
	mockStorage.EXPECT().GetLastBatchNumber(ctx, dbTx).Return(uint64(127), nil)
	mockStorage.EXPECT().IsBatchClosed(ctx, uint64(127), dbTx).Return(true, nil)
	mockStorage.EXPECT().GetLastBatchTime(ctx, dbTx).Return(time1, nil)
	// When calls to OpenBatch doesnt store the BatchL2Data yet
	batchContext.BatchL2Data = nil
	mockStorage.EXPECT().OpenBatchInStorage(ctx, batchContext, dbTx).Return(nil)
	mockStorage.EXPECT().GetLastNBatches(ctx, uint(2), dbTx).Return([]*state.Batch{&latestBatch, &previousBatch}, nil)
	mockStorage.EXPECT().IsBatchClosed(ctx, uint64(128), dbTx).Return(false, nil)
	mockStorage.EXPECT().GetForkIDByBatchNumber(uint64(128)).Return(uint64(state.FORKID_ETROG))
	mockExecutor.EXPECT().ProcessBatchV2(ctx, mock.Anything, mock.Anything).Return(&executorResponse, nil)
	mockStorage.EXPECT().CloseBatchInStorage(ctx, closingReceipt, dbTx).Return(nil)
	_, _, _, err = testState.ProcessAndStoreClosedBatchV2(ctx, processingCtx, dbTx, metrics.CallerLabel("test"))
	require.NoError(t, err)

	// Add assertions as needed
}
