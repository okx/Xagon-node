package state_test

import (
	"context"
	"math"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/mocks"
	"github.com/stretchr/testify/require"
)

func TestGetSyncingInfoErrors(t *testing.T) {
	var err error
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
	testState := state.NewState(stateCfg, mockStorage, mockExecutor, nil, nil, nil)
	mockStorage.EXPECT().Begin(ctx).Return(mocks.NewDbTxMock(t), nil)
	dbTx, err := testState.BeginStateTransaction(ctx)
	require.NoError(t, err)
	mockStorage.EXPECT().GetSyncInfoData(ctx, dbTx).Return(state.SyncInfoDataOnStorage{}, state.ErrNotFound).Once()
	_, err = testState.GetSyncingInfo(ctx, dbTx)
	require.ErrorIs(t, err, state.ErrStateNotSynchronized)

	mockStorage.EXPECT().GetSyncInfoData(ctx, dbTx).Return(state.SyncInfoDataOnStorage{InitialSyncingBatch: 1}, nil).Once()
	mockStorage.EXPECT().GetFirstL2BlockNumberForBatchNumber(ctx, uint64(1), dbTx).Return(uint64(0), state.ErrNotFound).Once()

	_, err = testState.GetSyncingInfo(ctx, dbTx)
	require.ErrorIs(t, err, state.ErrStateNotSynchronized)

	mockStorage.EXPECT().GetSyncInfoData(ctx, dbTx).Return(state.SyncInfoDataOnStorage{InitialSyncingBatch: 1}, nil).Once()
	mockStorage.EXPECT().GetFirstL2BlockNumberForBatchNumber(ctx, uint64(1), dbTx).Return(uint64(123), nil).Once()
	mockStorage.EXPECT().GetLastL2BlockNumber(ctx, dbTx).Return(uint64(0), state.ErrNotFound).Once()
	_, err = testState.GetSyncingInfo(ctx, dbTx)
	require.ErrorIs(t, err, state.ErrStateNotSynchronized)

	mockStorage.EXPECT().GetSyncInfoData(ctx, dbTx).Return(state.SyncInfoDataOnStorage{InitialSyncingBatch: 1}, nil).Once()
	mockStorage.EXPECT().GetFirstL2BlockNumberForBatchNumber(ctx, uint64(1), dbTx).Return(uint64(123), nil).Once()
	mockStorage.EXPECT().GetLastL2BlockNumber(ctx, dbTx).Return(uint64(567), nil).Once()
	mockStorage.EXPECT().GetLastBatchNumber(ctx, dbTx).Return(uint64(0), state.ErrNotFound).Once()
	_, err = testState.GetSyncingInfo(ctx, dbTx)
	require.ErrorIs(t, err, state.ErrStateNotSynchronized)
}

func TestGetSyncingInfoOk(t *testing.T) {
	var err error
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
	testState := state.NewState(stateCfg, mockStorage, mockExecutor, nil, nil, nil)
	mockStorage.EXPECT().Begin(ctx).Return(mocks.NewDbTxMock(t), nil)
	dbTx, err := testState.BeginStateTransaction(ctx)
	require.NoError(t, err)

	mockStorage.EXPECT().GetSyncInfoData(ctx, dbTx).Return(state.SyncInfoDataOnStorage{InitialSyncingBatch: 1, LastBatchNumberSeen: 50, LastBatchNumberConsolidated: 42}, nil).Once()
	mockStorage.EXPECT().GetFirstL2BlockNumberForBatchNumber(ctx, uint64(1), dbTx).Return(uint64(123), nil).Once()
	mockStorage.EXPECT().GetLastL2BlockNumber(ctx, dbTx).Return(uint64(567), nil).Once()
	mockStorage.EXPECT().GetLastBatchNumber(ctx, dbTx).Return(uint64(12), nil).Once()
	res, err := testState.GetSyncingInfo(ctx, dbTx)
	require.NoError(t, err)
	require.Equal(t, state.SyncingInfo{
		InitialSyncingBlock:   uint64(123),
		CurrentBlockNumber:    uint64(567),
		EstimatedHighestBlock: uint64(597),
		IsSynchronizing:       true,
	}, res)
}
