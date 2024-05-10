package test_l2_shared

import (
	"context"
	"errors"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/state"
	commonSync "github.com/0xPolygonHermez/zkevm-node/synchronizer/common"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/l2_sync"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/l2_sync/l2_shared"
	mock_l2_shared "github.com/0xPolygonHermez/zkevm-node/synchronizer/l2_sync/l2_shared/mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	hash1 = common.HexToHash("0x29e885edaf8e4b51e1d2e05f9da28161d2fb4f6b1d53827d9b80a23cf2d7d9f1")
	hash2 = common.HexToHash("0x979b141b8bcd3ba17815cd76811f1fca1cabaa9d51f7c00712606970f81d6e37")
	cfg   = l2_sync.Config{
		AcceptEmptyClosedBatches: true,
	}
)

func TestCacheEmpty(t *testing.T) {
	mockExecutor := mock_l2_shared.NewSyncTrustedBatchExecutor(t)
	mockTimer := &commonSync.MockTimerProvider{}
	mockL1SyncChecker := mock_l2_shared.NewL1SyncGlobalExitRootChecker(t)
	sut := l2_shared.NewProcessorTrustedBatchSync(mockExecutor, mockTimer, mockL1SyncChecker, cfg)

	current, previous := sut.GetCurrentAndPreviousBatchFromCache(&l2_shared.TrustedState{
		LastTrustedBatches: []*state.Batch{nil, nil},
	})
	require.Nil(t, current)
	require.Nil(t, previous)
	current, previous = sut.GetCurrentAndPreviousBatchFromCache(&l2_shared.TrustedState{
		LastTrustedBatches: []*state.Batch{nil},
	})
	require.Nil(t, current)
	require.Nil(t, previous)

	current, previous = sut.GetCurrentAndPreviousBatchFromCache(&l2_shared.TrustedState{
		LastTrustedBatches: []*state.Batch{},
	})
	require.Nil(t, current)
	require.Nil(t, previous)
}

func TestCacheJustCurrent(t *testing.T) {
	mockExecutor := mock_l2_shared.NewSyncTrustedBatchExecutor(t)
	mockTimer := &commonSync.MockTimerProvider{}
	batchA := state.Batch{
		BatchNumber: 123,
		Coinbase:    common.HexToAddress("0x123"),
	}
	status := l2_shared.TrustedState{
		LastTrustedBatches: []*state.Batch{&batchA},
	}
	sut := l2_shared.NewProcessorTrustedBatchSync(mockExecutor, mockTimer, nil, cfg)

	current, previous := sut.GetCurrentAndPreviousBatchFromCache(&status)
	require.Nil(t, previous)
	require.Equal(t, &batchA, current)
	require.NotEqual(t, &batchA, &current)
}

func TestCacheJustPrevious(t *testing.T) {
	mockExecutor := mock_l2_shared.NewSyncTrustedBatchExecutor(t)
	mockTimer := &commonSync.MockTimerProvider{}
	batchA := state.Batch{
		BatchNumber: 123,
		Coinbase:    common.HexToAddress("0x123"),
	}
	status := l2_shared.TrustedState{
		LastTrustedBatches: []*state.Batch{nil, &batchA},
	}
	sut := l2_shared.NewProcessorTrustedBatchSync(mockExecutor, mockTimer, nil, cfg)

	current, previous := sut.GetCurrentAndPreviousBatchFromCache(&status)
	require.Nil(t, current)
	require.Equal(t, &batchA, previous)
	require.NotEqual(t, &batchA, &previous)
}

type TestDataForProcessorTrustedBatchSync struct {
	mockTimer          *commonSync.MockTimerProvider
	mockExecutor       *mock_l2_shared.SyncTrustedBatchExecutor
	sut                *l2_shared.ProcessorTrustedBatchSync
	trustedNodeBatch   *types.Batch
	stateCurrentBatch  *state.Batch
	statePreviousBatch *state.Batch
}

func newTestDataForProcessorTrustedBatchSync(t *testing.T) *TestDataForProcessorTrustedBatchSync {
	mockExecutor := mock_l2_shared.NewSyncTrustedBatchExecutor(t)
	mockTimer := &commonSync.MockTimerProvider{}
	return &TestDataForProcessorTrustedBatchSync{
		mockTimer:    mockTimer,
		mockExecutor: mockExecutor,
		sut:          l2_shared.NewProcessorTrustedBatchSync(mockExecutor, mockTimer, nil, cfg),
		stateCurrentBatch: &state.Batch{
			BatchNumber: 123,
			Coinbase:    common.HexToAddress("0x1230"),
			StateRoot:   common.HexToHash("0x123410"),
			WIP:         true,
		},
		statePreviousBatch: &state.Batch{
			BatchNumber: 122,
			Coinbase:    common.HexToAddress("0x1230"),
			StateRoot:   common.HexToHash("0x1220"),
			WIP:         false,
		},
		trustedNodeBatch: &types.Batch{
			Number:    123,
			Coinbase:  common.HexToAddress("0x1230"),
			StateRoot: common.HexToHash("0x123410"),
			Closed:    true,
		},
	}
}

func TestGetModeForProcessBatchIncremental(t *testing.T) {
	testData := newTestDataForProcessorTrustedBatchSync(t)
	testData.trustedNodeBatch.Closed = true
	testData.trustedNodeBatch.BatchL2Data = []byte("test")
	processData, err := testData.sut.GetModeForProcessBatch(testData.trustedNodeBatch, testData.stateCurrentBatch, testData.statePreviousBatch, "test")
	require.NoError(t, err)
	require.Equal(t, l2_shared.IncrementalProcessMode, processData.Mode, "current batch is WIP and have a intermediate state root")
	require.Equal(t, true, processData.BatchMustBeClosed, "the trustedNode batch is closed")
	require.Equal(t, testData.stateCurrentBatch.StateRoot, processData.OldStateRoot, "the old state root is the intermediate state root (the current batch state root)")
}

func TestGetModeForProcessBatchNothingNoNewL2BatchDataChangeGER(t *testing.T) {
	testData := newTestDataForProcessorTrustedBatchSync(t)
	testData.stateCurrentBatch.BatchL2Data = []byte("test")
	testData.stateCurrentBatch.GlobalExitRoot = hash1
	testData.trustedNodeBatch.Closed = true
	testData.trustedNodeBatch.BatchL2Data = []byte("test")
	testData.stateCurrentBatch.GlobalExitRoot = hash2
	processData, err := testData.sut.GetModeForProcessBatch(testData.trustedNodeBatch, testData.stateCurrentBatch, testData.statePreviousBatch, "test")
	require.NoError(t, err)
	require.Equal(t, l2_shared.NothingProcessMode, processData.Mode, "current batch is WIP and have a intermediate state root")
	require.Equal(t, true, processData.BatchMustBeClosed, "the trustedNode batch is closed")
	require.Equal(t, common.Hash{}, processData.OldStateRoot, "the old state root is none because don't need to be process")
}

func TestGetModeForProcessBatchFullProcessMode(t *testing.T) {
	testData := newTestDataForProcessorTrustedBatchSync(t)
	testData.trustedNodeBatch.Closed = true
	testData.trustedNodeBatch.BatchL2Data = []byte("test") // We add some data
	processData, err := testData.sut.GetModeForProcessBatch(testData.trustedNodeBatch, nil, testData.statePreviousBatch, "test")
	require.NoError(t, err)
	require.Equal(t, l2_shared.FullProcessMode, processData.Mode, "there is no local batch, so it needs to full process")
	require.Equal(t, true, processData.BatchMustBeClosed, "the trustedNode batch is closed")
	require.Equal(t, testData.statePreviousBatch.StateRoot, processData.OldStateRoot, "the old state root is the previous batch SR")
}

func TestGetModeForProcessBatchReprocessMode(t *testing.T) {
	testData := newTestDataForProcessorTrustedBatchSync(t)
	testData.trustedNodeBatch.Closed = true
	testData.trustedNodeBatch.BatchL2Data = []byte("test") // We add some data
	testData.stateCurrentBatch.StateRoot = state.ZeroHash
	processData, err := testData.sut.GetModeForProcessBatch(testData.trustedNodeBatch, testData.stateCurrentBatch, testData.statePreviousBatch, "test")
	require.NoError(t, err)
	require.Equal(t, l2_shared.ReprocessProcessMode, processData.Mode, "local batch doesnt have stateRoot but exists, so  so it needs to be reprocess")
	require.Equal(t, true, processData.BatchMustBeClosed, "the trustedNode batch is closed")
	require.Equal(t, testData.statePreviousBatch.StateRoot, processData.OldStateRoot, "the old state root is the previous batch SR")
}

func TestGetModeForProcessBatchNothing(t *testing.T) {
	testData := newTestDataForProcessorTrustedBatchSync(t)
	testData.stateCurrentBatch.WIP = true
	testData.trustedNodeBatch.Closed = true
	processData, err := testData.sut.GetModeForProcessBatch(testData.trustedNodeBatch, testData.stateCurrentBatch, testData.statePreviousBatch, "test")
	require.NoError(t, err)
	require.Equal(t, l2_shared.NothingProcessMode, processData.Mode, "current batch and trusted batch are the same, just need to be closed")
	require.Equal(t, true, processData.BatchMustBeClosed, "the trustedNode batch is closed")
	require.Equal(t, state.ZeroHash, processData.OldStateRoot, "no OldStateRoot, because you dont need to process anything")

	testData.stateCurrentBatch.WIP = false
	testData.trustedNodeBatch.Closed = true
	processData, err = testData.sut.GetModeForProcessBatch(testData.trustedNodeBatch, testData.stateCurrentBatch, testData.statePreviousBatch, "test")
	require.NoError(t, err)
	require.Equal(t, l2_shared.NothingProcessMode, processData.Mode, "current batch and trusted batch are the same, just need to be closed")
	require.Equal(t, false, processData.BatchMustBeClosed, "the trustedNode batch is closed but the state batch is also closed, so nothing to do")

	testData.stateCurrentBatch.WIP = false
	testData.trustedNodeBatch.Closed = false
	processData, err = testData.sut.GetModeForProcessBatch(testData.trustedNodeBatch, testData.stateCurrentBatch, testData.statePreviousBatch, "test")
	require.NoError(t, err)
	require.Equal(t, l2_shared.NothingProcessMode, processData.Mode, "current batch and trusted batch are the same, just need to be closed")
	require.Equal(t, false, processData.BatchMustBeClosed, "nothing to do")

	testData.stateCurrentBatch.WIP = false
	testData.trustedNodeBatch.Closed = false
	processData, err = testData.sut.GetModeForProcessBatch(testData.trustedNodeBatch, nil, testData.statePreviousBatch, "test")
	require.NoError(t, err)
	require.Equal(t, l2_shared.FullProcessMode, processData.Mode, "no batch in DB, fullprocess")
	require.Equal(t, false, processData.BatchMustBeClosed, "nothing to do")

	testData.stateCurrentBatch.WIP = false
	testData.trustedNodeBatch.Closed = true
	processData, err = testData.sut.GetModeForProcessBatch(testData.trustedNodeBatch, nil, testData.statePreviousBatch, "test")
	require.NoError(t, err)
	require.Equal(t, l2_shared.FullProcessMode, processData.Mode, "no batch in DB, fullprocess")
	require.Equal(t, true, processData.BatchMustBeClosed, "must be close")

}

func TestGetModeForEmptyAndClosedBatchConfiguredToReject(t *testing.T) {
	testData := newTestDataForProcessorTrustedBatchSync(t)
	testData.sut.Cfg.AcceptEmptyClosedBatches = false
	testData.sut.Cfg.ReprocessFullBatchOnClose = true
	testData.stateCurrentBatch.WIP = true
	testData.trustedNodeBatch.Closed = true
	processData, err := testData.sut.GetModeForProcessBatch(testData.trustedNodeBatch, testData.stateCurrentBatch, testData.statePreviousBatch, "test")
	require.Error(t, err)

	testData.stateCurrentBatch.WIP = false
	testData.trustedNodeBatch.Closed = true
	processData, err = testData.sut.GetModeForProcessBatch(testData.trustedNodeBatch, testData.stateCurrentBatch, testData.statePreviousBatch, "test")
	require.Error(t, err)

	testData.stateCurrentBatch.WIP = false
	testData.trustedNodeBatch.Closed = false
	processData, err = testData.sut.GetModeForProcessBatch(testData.trustedNodeBatch, testData.stateCurrentBatch, testData.statePreviousBatch, "test")
	require.NoError(t, err)
	require.Equal(t, l2_shared.NothingProcessMode, processData.Mode, "current batch and trusted batch are the same, just need to be closed")
	require.Equal(t, false, processData.BatchMustBeClosed, "nothing to do")

	testData.stateCurrentBatch.WIP = false
	testData.trustedNodeBatch.Closed = false
	processData, err = testData.sut.GetModeForProcessBatch(testData.trustedNodeBatch, nil, testData.statePreviousBatch, "test")
	require.NoError(t, err)
	require.Equal(t, l2_shared.FullProcessMode, processData.Mode, "current batch and trusted batch are the same, just need to be closed")
	require.Equal(t, false, processData.BatchMustBeClosed, "nothing to do")

	testData.stateCurrentBatch.WIP = false
	testData.trustedNodeBatch.Closed = true
	processData, err = testData.sut.GetModeForProcessBatch(testData.trustedNodeBatch, nil, testData.statePreviousBatch, "test")
	require.Error(t, err)
}

func TestGetModeReprocessFullBatchOnCloseTrue(t *testing.T) {
	testData := newTestDataForProcessorTrustedBatchSync(t)
	testData.sut.Cfg.AcceptEmptyClosedBatches = true
	testData.sut.Cfg.ReprocessFullBatchOnClose = true
	testData.stateCurrentBatch.WIP = true
	testData.stateCurrentBatch.BatchL2Data = common.Hex2Bytes("112233")
	testData.trustedNodeBatch.BatchL2Data = common.Hex2Bytes("11223344")
	testData.trustedNodeBatch.Closed = true
	// Is a incremental converted to reprocess
	testData.sut.Cfg.ReprocessFullBatchOnClose = true
	processData, err := testData.sut.GetModeForProcessBatch(testData.trustedNodeBatch, testData.stateCurrentBatch, testData.statePreviousBatch, "test")
	require.NoError(t, err)
	require.Equal(t, l2_shared.ReprocessProcessMode, processData.Mode, "current batch and trusted batch are the same, just need to be closed")
	// Is a incremental to close
	testData.sut.Cfg.ReprocessFullBatchOnClose = false
	processData, err = testData.sut.GetModeForProcessBatch(testData.trustedNodeBatch, testData.stateCurrentBatch, testData.statePreviousBatch, "test")
	require.NoError(t, err)
	require.Equal(t, l2_shared.IncrementalProcessMode, processData.Mode, "increment of batchl2data, need to incremental execution")
	// No previous batch, is a fullprocess
	testData.sut.Cfg.ReprocessFullBatchOnClose = true
	processData, err = testData.sut.GetModeForProcessBatch(testData.trustedNodeBatch, nil, testData.statePreviousBatch, "test")
	require.NoError(t, err)
	require.Equal(t, l2_shared.FullProcessMode, processData.Mode, "no previous batch and close, fullprocess")

}

func TestGetNextStatusClear(t *testing.T) {
	testData := newTestDataForProcessorTrustedBatchSync(t)
	previousStatus := l2_shared.TrustedState{
		LastTrustedBatches: []*state.Batch{testData.statePreviousBatch, testData.statePreviousBatch},
	}
	processResponse := l2_shared.NewProcessResponse()

	processResponse.ClearCache = true
	res, err := testData.sut.GetNextStatus(previousStatus, &processResponse, false, "test")
	require.NoError(t, err)
	require.True(t, res.IsEmpty())

	processResponse.ClearCache = false
	res, err = testData.sut.GetNextStatus(l2_shared.TrustedState{}, &processResponse, false, "test")
	require.NoError(t, err)
	require.True(t, res.IsEmpty())

	processResponse.ClearCache = false
	res, err = testData.sut.GetNextStatus(l2_shared.TrustedState{
		LastTrustedBatches: []*state.Batch{nil, nil},
	}, &processResponse, false, "test")
	require.NoError(t, err)
	require.True(t, res.IsEmpty())

	processResponse.ClearCache = false
	processResponse.UpdateBatchWithProcessBatchResponse = true
	res, err = testData.sut.GetNextStatus(l2_shared.TrustedState{
		LastTrustedBatches: []*state.Batch{nil, nil},
	}, &processResponse, false, "test")
	require.NoError(t, err)
	require.True(t, res.IsEmpty())
}

func TestGetNextStatusUpdate(t *testing.T) {
	testData := newTestDataForProcessorTrustedBatchSync(t)
	previousStatus := l2_shared.TrustedState{
		LastTrustedBatches: []*state.Batch{testData.statePreviousBatch, testData.statePreviousBatch},
	}
	processBatchResp := l2_shared.NewProcessResponse()
	newBatch := state.Batch{
		BatchNumber: 123,
		Coinbase:    common.HexToAddress("0x123467"),
		StateRoot:   common.HexToHash("0x123456"),
		WIP:         true,
	}
	processBatchResp.UpdateCurrentBatch(&newBatch)
	res, err := testData.sut.GetNextStatus(previousStatus, &processBatchResp, false, "test")
	require.NoError(t, err)
	require.False(t, res.IsEmpty())
	require.Equal(t, *res.LastTrustedBatches[0], newBatch)

	res, err = testData.sut.GetNextStatus(previousStatus, &processBatchResp, true, "test")
	require.NoError(t, err)
	require.False(t, res.IsEmpty())
	require.Nil(t, res.LastTrustedBatches[0])
	require.Equal(t, newBatch, *res.LastTrustedBatches[1])

	ProcessBatchResponse := &state.ProcessBatchResponse{
		NewStateRoot:     common.HexToHash("0x123-2"),
		NewAccInputHash:  common.HexToHash("0x123-3"),
		NewLocalExitRoot: common.HexToHash("0x123-4"),
		NewBatchNumber:   123,
	}
	processBatchResp.UpdateCurrentBatchWithExecutionResult(&newBatch, ProcessBatchResponse)
	res, err = testData.sut.GetNextStatus(previousStatus, &processBatchResp, true, "test")
	require.NoError(t, err)
	require.False(t, res.IsEmpty())
	require.Nil(t, res.LastTrustedBatches[0])
	require.Equal(t, processBatchResp.ProcessBatchResponse.NewStateRoot, res.LastTrustedBatches[1].StateRoot)
}

func TestGetNextStatusUpdateNothing(t *testing.T) {
	testData := newTestDataForProcessorTrustedBatchSync(t)

	batch0 := state.Batch{
		BatchNumber: 123,
	}
	batch1 := state.Batch{
		BatchNumber: 122,
	}
	previousStatus := l2_shared.TrustedState{
		LastTrustedBatches: []*state.Batch{&batch0, &batch1},
	}
	ProcessResponse := l2_shared.NewProcessResponse()
	newStatus, err := testData.sut.GetNextStatus(previousStatus, &ProcessResponse, false, "test")
	require.NoError(t, err)
	require.Equal(t, &previousStatus, newStatus)
	// If batch is close move current batch to previous one
	newStatus, err = testData.sut.GetNextStatus(previousStatus, &ProcessResponse, true, "test")
	require.NoError(t, err)
	require.Equal(t, &l2_shared.TrustedState{
		LastTrustedBatches: []*state.Batch{nil, &batch0},
	}, newStatus)
}

func TestGetNextStatusDiscardCache(t *testing.T) {
	testData := newTestDataForProcessorTrustedBatchSync(t)
	ProcessResponse := l2_shared.NewProcessResponse()
	ProcessResponse.DiscardCache()
	newStatus, err := testData.sut.GetNextStatus(l2_shared.TrustedState{}, &ProcessResponse, false, "test")
	require.NoError(t, err)
	require.True(t, newStatus.IsEmpty())
}

func TestGetNextStatusUpdateCurrentBatch(t *testing.T) {
	testData := newTestDataForProcessorTrustedBatchSync(t)
	ProcessResponse := l2_shared.NewProcessResponse()
	batch := state.Batch{
		BatchNumber: 123,
	}
	ProcessResponse.UpdateCurrentBatch(&batch)
	newStatus, err := testData.sut.GetNextStatus(l2_shared.TrustedState{}, &ProcessResponse, false, "test")
	require.NoError(t, err)
	require.Equal(t, &l2_shared.TrustedState{
		LastTrustedBatches: []*state.Batch{&batch, nil},
	}, newStatus)
}

func TestGetNextStatusUpdateExecutionResult(t *testing.T) {
	testData := newTestDataForProcessorTrustedBatchSync(t)
	ProcessResponse := l2_shared.NewProcessResponse()
	batch := state.Batch{
		BatchNumber: 123,
	}
	previousStatus := l2_shared.TrustedState{
		LastTrustedBatches: []*state.Batch{nil, nil},
	}

	ProcessResponse.UpdateCurrentBatchWithExecutionResult(&batch, &state.ProcessBatchResponse{
		NewStateRoot: common.HexToHash("0x123"),
	})
	newStatus, err := testData.sut.GetNextStatus(previousStatus, &ProcessResponse, false, "test")
	require.NoError(t, err)
	require.Equal(t, common.HexToHash("0x123"), newStatus.LastTrustedBatches[0].StateRoot)
}

func TestExecuteProcessBatchError(t *testing.T) {
	testData := newTestDataForProcessorTrustedBatchSync(t)

	data := l2_shared.ProcessData{
		Mode:              l2_shared.NothingProcessMode,
		BatchMustBeClosed: true,
	}
	returnedError := errors.New("error")
	testData.mockExecutor.EXPECT().NothingProcess(mock.Anything, mock.Anything, mock.Anything).Return(nil, returnedError)
	_, err := testData.sut.ExecuteProcessBatch(context.Background(), &data, nil)
	require.ErrorIs(t, returnedError, err)
}
