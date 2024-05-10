package test_l2_shared

import (
	"context"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/common/syncinterfaces"
	mock_syncinterfaces "github.com/0xPolygonHermez/zkevm-node/synchronizer/common/syncinterfaces/mocks"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/l2_sync/l2_shared"
	"github.com/stretchr/testify/require"
)

// Use case 1:
// - Running incaberry mode no forkid7 yet
// expected:
// -

func TestExecutorSelectorFirstConfiguredExecutor(t *testing.T) {
	mockIncaberry := mock_syncinterfaces.NewSyncTrustedStateExecutor(t)
	mock1Etrog := mock_syncinterfaces.NewSyncTrustedStateExecutor(t)
	mockState := mock_syncinterfaces.NewStateFullInterface(t)
	mockState.EXPECT().GetForkIDByBatchNumber(uint64(1 + 1)).Return(uint64(6))
	forkIdInterval := state.ForkIDInterval{
		FromBatchNumber: 0,
		ToBatchNumber:   ^uint64(0),
	}
	mockState.EXPECT().GetForkIDInMemory(uint64(6)).Return(&forkIdInterval)
	sut := l2_shared.NewSyncTrustedStateExecutorSelector(map[uint64]syncinterfaces.SyncTrustedStateExecutor{
		uint64(6): mockIncaberry,
		uint64(7): mock1Etrog,
	}, mockState)

	executor, maxBatch := sut.GetExecutor(1, 200)
	require.Equal(t, mockIncaberry, executor)
	require.Equal(t, uint64(200), maxBatch)
}

func TestExecutorSelectorFirstExecutorCapped(t *testing.T) {
	mockIncaberry := mock_syncinterfaces.NewSyncTrustedStateExecutor(t)
	mock1Etrog := mock_syncinterfaces.NewSyncTrustedStateExecutor(t)
	mockState := mock_syncinterfaces.NewStateFullInterface(t)
	interval := state.ForkIDInterval{
		FromBatchNumber: 1,
		ToBatchNumber:   99,
		ForkId:          6,
	}
	mockState.EXPECT().GetForkIDByBatchNumber(uint64(1 + 1)).Return(uint64(6))
	mockState.EXPECT().GetForkIDInMemory(uint64(6)).Return(&interval)
	sut := l2_shared.NewSyncTrustedStateExecutorSelector(map[uint64]syncinterfaces.SyncTrustedStateExecutor{
		uint64(6): mockIncaberry,
		uint64(7): mock1Etrog,
	}, mockState)

	executor, maxBatch := sut.GetExecutor(1, 200)
	require.Equal(t, mockIncaberry, executor)
	require.Equal(t, uint64(99), maxBatch)
}

func TestExecutorSelectorEtrogBatchForkId7(t *testing.T) {
	mockIncaberry := mock_syncinterfaces.NewSyncTrustedStateExecutor(t)
	mock1Etrog := mock_syncinterfaces.NewSyncTrustedStateExecutor(t)
	mockState := mock_syncinterfaces.NewStateFullInterface(t)
	interval := state.ForkIDInterval{
		FromBatchNumber: 100,
		ToBatchNumber:   300,
		ForkId:          7,
	}
	mockState.EXPECT().GetForkIDByBatchNumber(uint64(100 + 1)).Return(uint64(7))
	mockState.EXPECT().GetForkIDInMemory(uint64(7)).Return(&interval)
	sut := l2_shared.NewSyncTrustedStateExecutorSelector(map[uint64]syncinterfaces.SyncTrustedStateExecutor{
		uint64(6): mockIncaberry,
		uint64(7): mock1Etrog,
	}, mockState)

	executor, maxBatch := sut.GetExecutor(100, 200)
	require.Equal(t, mockIncaberry, executor)
	require.Equal(t, uint64(200), maxBatch)
}

func TestUnsupportedForkId(t *testing.T) {
	mockIncaberry := mock_syncinterfaces.NewSyncTrustedStateExecutor(t)
	mock1Etrog := mock_syncinterfaces.NewSyncTrustedStateExecutor(t)
	mockState := mock_syncinterfaces.NewStateFullInterface(t)

	mockState.EXPECT().GetForkIDByBatchNumber(uint64(100 + 1)).Return(uint64(8))

	sut := l2_shared.NewSyncTrustedStateExecutorSelector(map[uint64]syncinterfaces.SyncTrustedStateExecutor{
		uint64(6): mockIncaberry,
		uint64(7): mock1Etrog,
	}, mockState)

	executor, _ := sut.GetExecutor(100, 200)
	require.Equal(t, nil, executor)

	err := sut.SyncTrustedState(context.Background(), 100, 200)
	require.ErrorIs(t, err, syncinterfaces.ErrCantSyncFromL2)

}
