package test_l2_shared

import (
	"context"
	"math/big"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/common"
	mock_syncinterfaces "github.com/0xPolygonHermez/zkevm-node/synchronizer/common/syncinterfaces/mocks"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/l2_sync/l2_shared"
	l2sharedmocks "github.com/0xPolygonHermez/zkevm-node/synchronizer/l2_sync/l2_shared/mocks"
	syncMocks "github.com/0xPolygonHermez/zkevm-node/synchronizer/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type testDataTrustedBatchRetrieve struct {
	mockBatchProcessor *l2sharedmocks.BatchProcessor
	mockZkEVMClient    *mock_syncinterfaces.ZKEVMClientTrustedBatchesGetter
	mockState          *l2sharedmocks.StateInterface
	mockSync           *mock_syncinterfaces.SynchronizerFlushIDManager
	mockTimer          *common.MockTimerProvider
	mockDbTx           *syncMocks.DbTxMock
	TrustedStateMngr   *l2_shared.TrustedStateManager
	sut                *l2_shared.TrustedBatchesRetrieve
	ctx                context.Context
}

func newTestDataTrustedBatchRetrieve(t *testing.T) *testDataTrustedBatchRetrieve {
	mockBatchProcessor := l2sharedmocks.NewBatchProcessor(t)
	mockZkEVMClient := mock_syncinterfaces.NewZKEVMClientTrustedBatchesGetter(t)
	mockState := l2sharedmocks.NewStateInterface(t)
	mockSync := mock_syncinterfaces.NewSynchronizerFlushIDManager(t)
	mockTimer := &common.MockTimerProvider{}
	mockDbTx := syncMocks.NewDbTxMock(t)
	TrustedStateMngr := l2_shared.NewTrustedStateManager(mockTimer, 0)
	sut := l2_shared.NewTrustedBatchesRetrieve(mockBatchProcessor, mockZkEVMClient, mockState, mockSync, *TrustedStateMngr)
	ctx := context.TODO()
	return &testDataTrustedBatchRetrieve{
		mockBatchProcessor: mockBatchProcessor,
		mockZkEVMClient:    mockZkEVMClient,
		mockState:          mockState,
		mockSync:           mockSync,
		mockTimer:          mockTimer,
		mockDbTx:           mockDbTx,
		TrustedStateMngr:   TrustedStateMngr,
		sut:                sut,
		ctx:                ctx,
	}
}

const (
	closedBatch    = true
	notClosedBatch = false
)

// This test must do from 100 to 104.
// But the batch 100 is open on TrustedNode so it stop processing
func TestSyncTrustedBatchesToFromStopAfterFirstWIPBatch(t *testing.T) {
	data := newTestDataTrustedBatchRetrieve(t)
	data.mockZkEVMClient.EXPECT().BatchNumber(data.ctx).Return(uint64(102), nil)

	expectationsForSyncTrustedStateIteration(t, 100, notClosedBatch, data)

	err := data.sut.SyncTrustedState(data.ctx, 100, 104)
	require.NoError(t, err)
}

// This must process 100 (that is closed)
// and stop processing at 101 because is not yet close this batch
func TestSyncTrustedBatchesToFromStopAfterFirstWIPBatchCase2(t *testing.T) {
	data := newTestDataTrustedBatchRetrieve(t)
	data.mockZkEVMClient.EXPECT().BatchNumber(data.ctx).Return(uint64(102), nil)

	expectationsForSyncTrustedStateIteration(t, 100, closedBatch, data)
	expectationsForSyncTrustedStateIteration(t, 101, notClosedBatch, data)

	err := data.sut.SyncTrustedState(data.ctx, 100, 104)
	require.NoError(t, err)
}

// This test must do from 100 to 102. Is for check manually that the logs
// That is not tested but must not emit the log:
//   - Batch 101 is not closed. so we break synchronization from Trusted Node because can only have 1 WIP batch on state
func TestSyncTrustedBatchesToFromStopAfterFirstWIPBatchCase3(t *testing.T) {
	data := newTestDataTrustedBatchRetrieve(t)
	data.mockZkEVMClient.EXPECT().BatchNumber(data.ctx).Return(uint64(102), nil)
	expectationsForSyncTrustedStateIteration(t, 100, closedBatch, data)
	expectationsForSyncTrustedStateIteration(t, 101, closedBatch, data)
	expectationsForSyncTrustedStateIteration(t, 102, notClosedBatch, data)

	err := data.sut.SyncTrustedState(data.ctx, 100, 102)
	require.NoError(t, err)
}

func expectationsForSyncTrustedStateIteration(t *testing.T, batchNumber uint64, closed bool, data *testDataTrustedBatchRetrieve) {
	batch100 := &types.Batch{
		Number: types.ArgUint64(batchNumber),
		Closed: closed,
	}
	data.mockZkEVMClient.EXPECT().BatchByNumber(data.ctx, big.NewInt(0).SetUint64(batchNumber)).Return(batch100, nil)
	data.mockState.EXPECT().BeginStateTransaction(data.ctx).Return(data.mockDbTx, nil)
	// Get Previous Batch 99 from State
	stateBatch99 := &state.Batch{
		BatchNumber: batchNumber - 1,
	}
	data.mockState.EXPECT().GetBatchByNumber(data.ctx, uint64(batchNumber-1), data.mockDbTx).Return(stateBatch99, nil)
	stateBatch100 := &state.Batch{
		BatchNumber: batchNumber,
	}
	data.mockState.EXPECT().GetBatchByNumber(data.ctx, uint64(batchNumber), data.mockDbTx).Return(stateBatch100, nil)
	data.mockBatchProcessor.EXPECT().ProcessTrustedBatch(data.ctx, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)
	data.mockSync.EXPECT().CheckFlushID(mock.Anything).Return(nil)
	data.mockDbTx.EXPECT().Commit(data.ctx).Return(nil)
}
