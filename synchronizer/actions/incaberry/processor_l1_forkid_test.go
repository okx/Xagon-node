package incaberry_test

import (
	"context"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/actions/incaberry"
	mock_syncinterfaces "github.com/0xPolygonHermez/zkevm-node/synchronizer/common/syncinterfaces/mocks"
	syncMocks "github.com/0xPolygonHermez/zkevm-node/synchronizer/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockForkdIdTest struct {
	mockState *mock_syncinterfaces.StateFullInterface
	mockSync  *mock_syncinterfaces.SynchronizerIsTrustedSequencer
	mockDbTx  *syncMocks.DbTxMock
}

func newMockForkdIdTest(t *testing.T) *mockForkdIdTest {
	mockState := mock_syncinterfaces.NewStateFullInterface(t)
	mockSync := mock_syncinterfaces.NewSynchronizerIsTrustedSequencer(t)
	mockDbTx := syncMocks.NewDbTxMock(t)
	return &mockForkdIdTest{mockState, mockSync, mockDbTx}
}

func newL1Block(blockNumber uint64, forkId uint64, fromBatchNumber uint64, version string) *etherman.Block {
	return &etherman.Block{
		SequencedBatches: [][]etherman.SequencedBatch{},
		BlockNumber:      blockNumber,
		ForkIDs:          []etherman.ForkID{{ForkID: forkId, BatchNumber: fromBatchNumber, Version: version}},
	}
}

func TestReceiveExistingForkIdAnotherFromBatchNumber(t *testing.T) {
	mocks := newMockForkdIdTest(t)
	sut := incaberry.NewProcessorForkId(mocks.mockState, mocks.mockSync)
	forkIdsOnState := []state.ForkIDInterval{
		{FromBatchNumber: 1, ToBatchNumber: 100, ForkId: 6, Version: "1.0.0", BlockNumber: 123},
		{FromBatchNumber: 101, ToBatchNumber: 200, ForkId: 7, Version: "1.0.0", BlockNumber: 123},
	}
	mocks.mockState.EXPECT().GetForkIDs(mock.Anything, mock.Anything).Return(forkIdsOnState, nil)
	err := sut.Process(context.Background(), etherman.Order{Pos: 0}, newL1Block(123, 6, 500, "1.0.0"), mocks.mockDbTx)
	require.Error(t, err)
}

func TestReceiveExistsForkIdSameBatchNumberSameBlockNumber(t *testing.T) {
	mocks := newMockForkdIdTest(t)
	sut := incaberry.NewProcessorForkId(mocks.mockState, mocks.mockSync)
	forkIdsOnState := []state.ForkIDInterval{
		{FromBatchNumber: 1, ToBatchNumber: 100, ForkId: 6, Version: "1.0.0", BlockNumber: 123},
		{FromBatchNumber: 101, ToBatchNumber: 200, ForkId: 7, Version: "1.0.0", BlockNumber: 123},
	}
	mocks.mockState.EXPECT().GetForkIDs(mock.Anything, mock.Anything).Return(forkIdsOnState, nil)

	err := sut.Process(context.Background(), etherman.Order{Pos: 0}, newL1Block(123, 6, 0, "1.0.0"), mocks.mockDbTx)
	require.NoError(t, err)
}

func TestReceiveExistsForkIdSameBatchNumberAnotherBlockNumberAndNotLastForkId(t *testing.T) {
	mocks := newMockForkdIdTest(t)
	sut := incaberry.NewProcessorForkId(mocks.mockState, mocks.mockSync)
	forkIdsOnState := []state.ForkIDInterval{
		{FromBatchNumber: 1, ToBatchNumber: 100, ForkId: 6, Version: "1.0.0", BlockNumber: 123},
		{FromBatchNumber: 101, ToBatchNumber: 200, ForkId: 7, Version: "1.0.0", BlockNumber: 123},
	}
	mocks.mockState.EXPECT().GetForkIDs(mock.Anything, mock.Anything).Return(forkIdsOnState, nil)
	//mocks.mockDbTx.EXPECT().Rollback(mock.Anything).Return(nil)
	err := sut.Process(context.Background(), etherman.Order{Pos: 0}, newL1Block(456, 6, 0, "1.0.0"), mocks.mockDbTx)
	require.Error(t, err)
}

func TestReceiveAForkIdWithIdPreviousToCurrentOnState(t *testing.T) {
	mocks := newMockForkdIdTest(t)
	sut := incaberry.NewProcessorForkId(mocks.mockState, mocks.mockSync)
	forkIdsOnState := []state.ForkIDInterval{
		{FromBatchNumber: 100, ToBatchNumber: 200, ForkId: 6, Version: "1.0.0", BlockNumber: 123},
		{FromBatchNumber: 201, ToBatchNumber: 300, ForkId: 7, Version: "1.0.0", BlockNumber: 123},
	}
	mocks.mockState.EXPECT().GetForkIDs(mock.Anything, mock.Anything).Return(forkIdsOnState, nil)
	err := sut.Process(context.Background(), etherman.Order{Pos: 0}, newL1Block(456, 5, 0, "1.0.0"), mocks.mockDbTx)
	require.Error(t, err)
}

func TestReceiveExistsForkIdSameBatchNumberAnotherBlockNumberAndLastForkId(t *testing.T) {
	mocks := newMockForkdIdTest(t)
	sut := incaberry.NewProcessorForkId(mocks.mockState, mocks.mockSync)
	forkIdsOnState := []state.ForkIDInterval{
		{FromBatchNumber: 1, ToBatchNumber: 100, ForkId: 6, Version: "1.0.0", BlockNumber: 123},
		{FromBatchNumber: 101, ToBatchNumber: 200, ForkId: 7, Version: "1.0.0", BlockNumber: 123},
	}
	mocks.mockState.EXPECT().GetForkIDs(mock.Anything, mock.Anything).Return(forkIdsOnState, nil)
	mocks.mockState.EXPECT().UpdateForkIDBlockNumber(mock.Anything, uint64(7), uint64(456), true, mock.Anything).Return(nil)
	//mocks.mockDbTx.EXPECT().Commit(mock.Anything).Return(nil)
	err := sut.Process(context.Background(), etherman.Order{Pos: 0}, newL1Block(456, 7, 100, "1.0.0"), mocks.mockDbTx)
	require.NoError(t, err)
}

func TestReceiveNewForkIdAffectFutureBatch(t *testing.T) {
	mocks := newMockForkdIdTest(t)
	sut := incaberry.NewProcessorForkId(mocks.mockState, mocks.mockSync)
	forkIdsOnState := []state.ForkIDInterval{
		{FromBatchNumber: 1, ToBatchNumber: 100, ForkId: 6, Version: "1.0.0", BlockNumber: 123},
		{FromBatchNumber: 101, ToBatchNumber: 200, ForkId: 7, Version: "1.0.0", BlockNumber: 123},
	}
	mocks.mockState.EXPECT().GetForkIDs(mock.Anything, mock.Anything).Return(forkIdsOnState, nil)
	mocks.mockState.EXPECT().GetLastBatchNumber(mock.Anything, mock.Anything).Return(uint64(101), nil)
	mocks.mockState.EXPECT().AddForkIDInterval(mock.Anything, state.ForkIDInterval{FromBatchNumber: 102, ToBatchNumber: ^uint64(0), ForkId: 8, Version: "2.0.0", BlockNumber: 456}, mock.Anything).Return(nil)
	//mocks.mockDbTx.EXPECT().Commit(mock.Anything).Return(nil)
	err := sut.Process(context.Background(), etherman.Order{Pos: 0}, newL1Block(456, 8, 101, "2.0.0"), mocks.mockDbTx)
	require.NoError(t, err)
}

func TestReceiveNewForkIdAffectPastBatchTrustedNode(t *testing.T) {
	mocks := newMockForkdIdTest(t)
	sut := incaberry.NewProcessorForkId(mocks.mockState, mocks.mockSync)
	forkIdsOnState := []state.ForkIDInterval{
		{FromBatchNumber: 1, ToBatchNumber: 100, ForkId: 6, Version: "1.0.0", BlockNumber: 123},
		{FromBatchNumber: 101, ToBatchNumber: 200, ForkId: 7, Version: "1.0.0", BlockNumber: 123},
	}
	mocks.mockState.EXPECT().GetForkIDs(mock.Anything, mock.Anything).Return(forkIdsOnState, nil)
	mocks.mockState.EXPECT().GetLastBatchNumber(mock.Anything, mock.Anything).Return(uint64(101), nil)
	mocks.mockState.EXPECT().AddForkIDInterval(mock.Anything, state.ForkIDInterval{FromBatchNumber: 101, ToBatchNumber: ^uint64(0), ForkId: 8, Version: "2.0.0", BlockNumber: 456}, mock.Anything).Return(nil)
	mocks.mockSync.EXPECT().IsTrustedSequencer().Return(true)
	err := sut.Process(context.Background(), etherman.Order{Pos: 0}, newL1Block(456, 8, 100, "2.0.0"), mocks.mockDbTx)
	require.NoError(t, err)
}

func TestReceiveNewForkIdAffectPastBatchPermissionlessNode(t *testing.T) {
	mocks := newMockForkdIdTest(t)
	sut := incaberry.NewProcessorForkId(mocks.mockState, mocks.mockSync)
	forkIdsOnState := []state.ForkIDInterval{
		{FromBatchNumber: 1, ToBatchNumber: 100, ForkId: 6, Version: "1.0.0", BlockNumber: 123},
		{FromBatchNumber: 101, ToBatchNumber: 200, ForkId: 7, Version: "1.0.0", BlockNumber: 123},
	}
	mocks.mockState.EXPECT().GetForkIDs(mock.Anything, mock.Anything).Return(forkIdsOnState, nil)
	mocks.mockState.EXPECT().GetLastBatchNumber(mock.Anything, mock.Anything).Return(uint64(101), nil)
	mocks.mockState.EXPECT().AddForkIDInterval(mock.Anything, state.ForkIDInterval{FromBatchNumber: 101, ToBatchNumber: ^uint64(0), ForkId: 8, Version: "2.0.0", BlockNumber: 456}, mock.Anything).Return(nil)
	mocks.mockSync.EXPECT().IsTrustedSequencer().Return(false)
	mocks.mockState.EXPECT().ResetForkID(mock.Anything, uint64(101), mock.Anything).Return(nil)
	mocks.mockDbTx.EXPECT().Commit(mock.Anything).Return(nil)
	err := sut.Process(context.Background(), etherman.Order{Pos: 0}, newL1Block(456, 8, 100, "2.0.0"), mocks.mockDbTx)
	require.Error(t, err)
	require.Equal(t, "new ForkID detected, reseting synchronizarion", err.Error())
}
