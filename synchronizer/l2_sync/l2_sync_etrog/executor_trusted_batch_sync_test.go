package l2_sync_etrog

import (
	"context"
	"encoding/hex"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	mock_syncinterfaces "github.com/0xPolygonHermez/zkevm-node/synchronizer/common/syncinterfaces/mocks"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/l2_sync/l2_shared"
	mock_l2_sync_etrog "github.com/0xPolygonHermez/zkevm-node/synchronizer/l2_sync/l2_sync_etrog/mocks"
	"github.com/ethereum/go-ethereum/common"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const (
	// changeL2Block + deltaTimeStamp + indexL1InfoTree
	codedL2BlockHeader = "0b73e6af6f00000000"
	// 2 x [ tx coded in RLP + r,s,v,efficiencyPercentage]
	codedRLP2Txs1 = "ee02843b9aca00830186a0944d5cf5032b2a844602278b01199ed191a86c93ff88016345785d8a0000808203e88080bff0e780ba7db409339fd3f71969fa2cbf1b8535f6c725a1499d3318d3ef9c2b6340ddfab84add2c188f9efddb99771db1fe621c981846394ea4f035c85bcdd51bffee03843b9aca00830186a0944d5cf5032b2a844602278b01199ed191a86c93ff88016345785d8a0000808203e880805b346aa02230b22e62f73608de9ff39a162a6c24be9822209c770e3685b92d0756d5316ef954eefc58b068231ccea001fb7ac763ebe03afd009ad71cab36861e1bff"
)

var (
	hashExamplesValues = []string{"0x723e5c4c7ee7890e1e66c2e391d553ee792d2204ecb4fe921830f12f8dcd1a92",
		"0x9c8fa7ce2e197f9f1b3c30de9f93de3c1cb290e6c118a18446f47a9e1364c3ab",
		"0x896cfc0684057d0560e950dee352189528167f4663609678d19c7a506a03fe4e",
		"0xde6d2dac4b6e0cb39ed1924db533558a23e5c56ab60fadac8c7d21e7eceb121a",
		"0x9883711e78d02992ac1bd6f19de3bf7bb3f926742d4601632da23525e33f8555"}
)

type testDataForBathExecutor struct {
	ctx       context.Context
	stateMock *mock_l2_sync_etrog.StateInterface
	syncMock  *mock_syncinterfaces.SynchronizerFlushIDManager
	sut       *SyncTrustedBatchExecutorForEtrog
}

func TestIncrementalProcessUpdateBatchL2DataOnCache(t *testing.T) {
	// Arrange
	stateMock := mock_l2_sync_etrog.NewStateInterface(t)
	syncMock := mock_syncinterfaces.NewSynchronizerFlushIDManager(t)

	sut := SyncTrustedBatchExecutorForEtrog{
		state: stateMock,
		sync:  syncMock,
	}
	ctx := context.Background()

	stateBatchL2Data, _ := hex.DecodeString(codedL2BlockHeader + codedRLP2Txs1)
	trustedBatchL2Data, _ := hex.DecodeString(codedL2BlockHeader + codedRLP2Txs1 + codedL2BlockHeader + codedRLP2Txs1)
	expectedStateRoot := common.HexToHash("0x723e5c4c7ee7890e1e66c2e391d553ee792d2204ecb4fe921830f12f8dcd1a92")
	//deltaBatchL2Data := []byte{4}
	batchNumber := uint64(123)
	data := l2_shared.ProcessData{
		BatchNumber:  batchNumber,
		OldStateRoot: common.Hash{},
		TrustedBatch: &types.Batch{
			Number:      123,
			BatchL2Data: trustedBatchL2Data,
			StateRoot:   expectedStateRoot,
		},
		StateBatch: &state.Batch{
			BatchNumber: batchNumber,
			BatchL2Data: stateBatchL2Data,
		},
	}

	stateMock.EXPECT().UpdateWIPBatch(ctx, mock.Anything, mock.Anything).Return(nil).Once()
	stateMock.EXPECT().GetL1InfoTreeDataFromBatchL2Data(ctx, mock.Anything, mock.Anything).Return(map[uint32]state.L1DataV2{}, expectedStateRoot, common.Hash{}, nil).Once()
	stateMock.EXPECT().GetForkIDByBatchNumber(batchNumber).Return(uint64(7)).Once()

	processBatchResp := &state.ProcessBatchResponse{
		NewStateRoot: expectedStateRoot,
	}
	stateMock.EXPECT().ProcessBatchV2(ctx, mock.Anything, true).Return(processBatchResp, "", nil).Once()

	syncMock.EXPECT().PendingFlushID(mock.Anything, mock.Anything).Once()
	syncMock.EXPECT().CheckFlushID(mock.Anything).Return(nil).Maybe()
	// Act
	res, err := sut.IncrementalProcess(ctx, &data, nil)
	// Assert
	log.Info(res)
	require.NoError(t, err)
	require.Equal(t, trustedBatchL2Data, res.UpdateBatch.BatchL2Data)
	require.Equal(t, false, res.ClearCache)
}

// This test check that if you process incrementally a batch that need to be close
// the GER is update because we call UpdateWIPBatch
// NOTE: CloseBatch() doesnt update GER
func TestIncrementalProcessUpdateBatchL2DataAndGER(t *testing.T) {
	// Arrange
	stateMock := mock_l2_sync_etrog.NewStateInterface(t)
	syncMock := mock_syncinterfaces.NewSynchronizerFlushIDManager(t)

	sut := SyncTrustedBatchExecutorForEtrog{
		state: stateMock,
		sync:  syncMock,
	}
	ctx := context.Background()

	stateBatchL2Data, _ := hex.DecodeString(codedL2BlockHeader + codedRLP2Txs1)
	trustedBatchL2Data, _ := hex.DecodeString(codedL2BlockHeader + codedRLP2Txs1 + codedL2BlockHeader + codedRLP2Txs1)
	expectedStateRoot := common.HexToHash("0x723e5c4c7ee7890e1e66c2e391d553ee792d2204ecb4fe921830f12f8dcd1a92")
	//deltaBatchL2Data := []byte{4}
	batchNumber := uint64(123)
	data := l2_shared.ProcessData{
		BatchNumber:       batchNumber,
		OldStateRoot:      common.Hash{},
		BatchMustBeClosed: true,
		TrustedBatch: &types.Batch{
			Number:      123,
			BatchL2Data: trustedBatchL2Data,
			StateRoot:   expectedStateRoot,
			Closed:      true,
		},
		StateBatch: &state.Batch{
			BatchNumber:    batchNumber,
			BatchL2Data:    stateBatchL2Data,
			GlobalExitRoot: common.HexToHash("0x9c8fa7ce2e197f9f1b3c30de9f93de3c1cb290e6c118a18446f47a9e1364c3ab"),
		},
	}
	expectedUpdate := state.ProcessingReceipt{
		BatchNumber:    123,
		StateRoot:      expectedStateRoot,
		LocalExitRoot:  data.TrustedBatch.LocalExitRoot,
		GlobalExitRoot: data.TrustedBatch.GlobalExitRoot,
		AccInputHash:   data.TrustedBatch.AccInputHash,
		BatchL2Data:    trustedBatchL2Data,
	}

	stateMock.EXPECT().UpdateWIPBatch(ctx, expectedUpdate, mock.Anything).Return(nil).Once()
	stateMock.EXPECT().GetL1InfoTreeDataFromBatchL2Data(ctx, mock.Anything, mock.Anything).Return(map[uint32]state.L1DataV2{}, expectedStateRoot, common.Hash{}, nil).Once()
	stateMock.EXPECT().GetForkIDByBatchNumber(batchNumber).Return(uint64(7)).Once()

	processBatchResp := &state.ProcessBatchResponse{
		NewStateRoot: expectedStateRoot,
	}
	stateMock.EXPECT().ProcessBatchV2(ctx, mock.Anything, true).Return(processBatchResp, "", nil).Once()

	syncMock.EXPECT().PendingFlushID(mock.Anything, mock.Anything).Once()
	syncMock.EXPECT().CheckFlushID(mock.Anything).Return(nil).Maybe()
	expectedUpdateClosed := expectedUpdate
	expectedUpdateClosed.GlobalExitRoot = common.Hash{}
	expectedUpdateClosed.ClosingReason = state.SyncL2TrustedBatchClosingReason
	stateMock.EXPECT().CloseBatch(ctx, expectedUpdateClosed, mock.Anything).Return(nil).Once()
	// Act
	res, err := sut.IncrementalProcess(ctx, &data, nil)
	// Assert
	log.Info(res)
	require.NoError(t, err)
	require.Equal(t, trustedBatchL2Data, res.UpdateBatch.BatchL2Data)
	require.Equal(t, false, res.ClearCache)
}

func newTestData(t *testing.T) testDataForBathExecutor {
	stateMock := mock_l2_sync_etrog.NewStateInterface(t)
	syncMock := mock_syncinterfaces.NewSynchronizerFlushIDManager(t)

	sut := SyncTrustedBatchExecutorForEtrog{
		state: stateMock,
		sync:  syncMock,
	}
	return testDataForBathExecutor{
		ctx:       context.Background(),
		stateMock: stateMock,
		syncMock:  syncMock,
		sut:       &sut,
	}
}

func newData() l2_shared.ProcessData {
	return l2_shared.ProcessData{
		BatchNumber: 123,
		Mode:        l2_shared.IncrementalProcessMode,
		DebugPrefix: "test",
		StateBatch: &state.Batch{
			BatchNumber:   123,
			StateRoot:     common.HexToHash(hashExamplesValues[0]),
			LocalExitRoot: common.HexToHash(hashExamplesValues[1]),
			AccInputHash:  common.HexToHash(hashExamplesValues[2]),
			WIP:           true,
		},
		TrustedBatch: &types.Batch{
			Number:        123,
			StateRoot:     common.HexToHash(hashExamplesValues[0]),
			LocalExitRoot: common.HexToHash(hashExamplesValues[1]),
			AccInputHash:  common.HexToHash(hashExamplesValues[2]),
			BatchL2Data:   []byte{1, 2, 3, 4},
			Closed:        false,
		},
	}
}

func TestNothingProcessDontCloseBatch(t *testing.T) {
	testData := newTestData(t)

	// Arrange
	data := l2_shared.ProcessData{
		BatchNumber:       123,
		Mode:              l2_shared.NothingProcessMode,
		BatchMustBeClosed: false,
		DebugPrefix:       "test",
		StateBatch:        &state.Batch{WIP: true},
		TrustedBatch:      &types.Batch{},
	}

	response, err := testData.sut.NothingProcess(testData.ctx, &data, nil)
	require.NoError(t, err)
	require.Equal(t, false, response.ClearCache)
	require.Equal(t, false, response.UpdateBatchWithProcessBatchResponse)
	require.Equal(t, true, data.StateBatch.WIP)
}

func TestNothingProcessDoesntMatchBatchCantProcessBecauseNoPreviousStateBatch(t *testing.T) {
	testData := newTestData(t)
	// Arrange
	data := l2_shared.ProcessData{
		BatchNumber:       123,
		Mode:              l2_shared.NothingProcessMode,
		BatchMustBeClosed: false,
		DebugPrefix:       "test",
		StateBatch: &state.Batch{
			BatchNumber: 123,
			StateRoot:   common.HexToHash(hashExamplesValues[1]),
			WIP:         true,
		},
		TrustedBatch: &types.Batch{
			Number:    123,
			StateRoot: common.HexToHash(hashExamplesValues[0]),
		},
		PreviousStateBatch: nil,
	}

	_, err := testData.sut.NothingProcess(testData.ctx, &data, nil)
	require.ErrorIs(t, err, ErrCantReprocessBatchMissingPreviousStateBatch)
}

func TestNothingProcessDoesntMatchBatchReprocess(t *testing.T) {
	testData := newTestData(t)
	// Arrange
	data := l2_shared.ProcessData{
		BatchNumber:       123,
		Mode:              l2_shared.NothingProcessMode,
		BatchMustBeClosed: false,
		DebugPrefix:       "test",
		StateBatch: &state.Batch{
			BatchNumber: 123,
			StateRoot:   common.HexToHash(hashExamplesValues[1]),
			BatchL2Data: []byte{1, 2, 3, 4},
			WIP:         true,
		},
		TrustedBatch: &types.Batch{
			Number:      123,
			StateRoot:   common.HexToHash(hashExamplesValues[0]),
			BatchL2Data: []byte{1, 2, 3, 4},
		},
		PreviousStateBatch: &state.Batch{
			BatchNumber: 122,
			StateRoot:   common.HexToHash(hashExamplesValues[2]),
		},
	}
	testData.stateMock.EXPECT().GetLastVirtualBatchNum(testData.ctx, mock.Anything).Return(uint64(122), nil).Maybe()
	testData.stateMock.EXPECT().ResetTrustedState(testData.ctx, data.BatchNumber-1, mock.Anything).Return(nil).Once()
	testData.stateMock.EXPECT().OpenBatch(testData.ctx, mock.Anything, mock.Anything).Return(nil).Once()
	testData.stateMock.EXPECT().GetL1InfoTreeDataFromBatchL2Data(testData.ctx, mock.Anything, mock.Anything).Return(map[uint32]state.L1DataV2{}, common.Hash{}, common.Hash{}, nil).Once()
	testData.stateMock.EXPECT().GetForkIDByBatchNumber(data.BatchNumber).Return(uint64(state.FORKID_ETROG)).Once()
	testData.syncMock.EXPECT().PendingFlushID(mock.Anything, mock.Anything).Once()
	testData.stateMock.EXPECT().UpdateWIPBatch(testData.ctx, mock.Anything, mock.Anything).Return(nil).Once()
	processBatchResp := &state.ProcessBatchResponse{
		NewStateRoot: data.TrustedBatch.StateRoot,
	}
	testData.stateMock.EXPECT().ProcessBatchV2(testData.ctx, mock.Anything, true).Return(processBatchResp, "", nil).Once()
	testData.stateMock.EXPECT().GetBatchByNumber(testData.ctx, data.BatchNumber, mock.Anything).Return(&state.Batch{}, nil).Once()
	_, err := testData.sut.NothingProcess(testData.ctx, &data, nil)
	require.NoError(t, err)
}

func TestReprocessRejectDeleteVirtualBatch(t *testing.T) {
	testData := newTestData(t)
	// Arrange
	data := l2_shared.ProcessData{
		BatchNumber:       123,
		Mode:              l2_shared.NothingProcessMode,
		BatchMustBeClosed: false,
		DebugPrefix:       "test",
		StateBatch: &state.Batch{
			BatchNumber: 123,
			StateRoot:   common.HexToHash(hashExamplesValues[1]),
			BatchL2Data: []byte{1, 2, 3, 4},
			WIP:         true,
		},
		TrustedBatch: &types.Batch{
			Number:      123,
			StateRoot:   common.HexToHash(hashExamplesValues[0]),
			BatchL2Data: []byte{1, 2, 3, 4},
		},
		PreviousStateBatch: &state.Batch{
			BatchNumber: 122,
			StateRoot:   common.HexToHash(hashExamplesValues[2]),
		},
	}
	testData.stateMock.EXPECT().GetLastVirtualBatchNum(testData.ctx, mock.Anything).Return(uint64(123), nil).Maybe()
	_, err := testData.sut.ReProcess(testData.ctx, &data, nil)
	require.Error(t, err)
}

func TestNothingProcessIfBatchMustBeClosedThenCloseBatch(t *testing.T) {
	testData := newTestData(t)
	// Arrange
	data := newData()
	data.StateBatch.BatchL2Data = data.TrustedBatch.BatchL2Data
	data.BatchMustBeClosed = true
	testData.stateMock.EXPECT().CloseBatch(testData.ctx, mock.Anything, mock.Anything).Return(nil).Once()

	response, err := testData.sut.NothingProcess(testData.ctx, &data, nil)
	require.NoError(t, err)
	require.Equal(t, false, response.ClearCache)
	require.Equal(t, false, response.UpdateBatchWithProcessBatchResponse)
	require.Equal(t, false, data.StateBatch.WIP)
}

func TestNothingProcessIfNotBatchMustBeClosedThenDoNothing(t *testing.T) {
	testData := newTestData(t)
	data := newData()
	data.StateBatch.BatchL2Data = data.TrustedBatch.BatchL2Data
	data.BatchMustBeClosed = false
	_, err := testData.sut.NothingProcess(testData.ctx, &data, nil)
	require.NoError(t, err)
}
func TestCloseBatchGivenAlreadyCloseAndTheBatchDataDoesntMatchExpectedThenHalt(t *testing.T) {
	testData := newTestData(t)
	data := newData()

	testData.stateMock.EXPECT().CloseBatch(testData.ctx, mock.Anything, mock.Anything).Return(state.ErrBatchAlreadyClosed).Once()
	testData.stateMock.EXPECT().GetBatchByNumber(testData.ctx, data.BatchNumber, mock.Anything).Return(&state.Batch{}, nil).Once()
	res := testData.sut.CloseBatch(testData.ctx, data.TrustedBatch, nil, "test")
	require.ErrorIs(t, res, ErrCriticalClosedBatchDontContainExpectedData)
}

func TestCloseBatchGivenAlreadyClosedAndTheDataAreRightThenNoError(t *testing.T) {
	testData := newTestData(t)
	data := newData()
	data.TrustedBatch.Closed = true
	stateBatchEqualToTrusted := &state.Batch{
		BatchNumber:    data.BatchNumber,
		GlobalExitRoot: data.TrustedBatch.GlobalExitRoot,
		LocalExitRoot:  data.TrustedBatch.LocalExitRoot,
		StateRoot:      data.TrustedBatch.StateRoot,
		AccInputHash:   data.TrustedBatch.AccInputHash,
		BatchL2Data:    data.TrustedBatch.BatchL2Data,
		WIP:            false,
		Timestamp:      time.Unix(int64(data.TrustedBatch.Timestamp+123), 0),
	}
	testData.stateMock.EXPECT().CloseBatch(testData.ctx, mock.Anything, mock.Anything).Return(state.ErrBatchAlreadyClosed).Once()
	testData.stateMock.EXPECT().GetBatchByNumber(testData.ctx, data.BatchNumber, mock.Anything).Return(stateBatchEqualToTrusted, nil).Once()
	// No call to HALT!
	res := testData.sut.CloseBatch(testData.ctx, data.TrustedBatch, nil, "test")
	require.NoError(t, res)
}

func TestEmptyWIPBatch(t *testing.T) {
	testData := newTestData(t)
	// Arrange
	expectedBatch := state.Batch{
		BatchNumber:    123,
		Coinbase:       common.HexToAddress("0x01"),
		StateRoot:      common.HexToHash("0x02"),
		GlobalExitRoot: common.HexToHash("0x03"),
		LocalExitRoot:  common.HexToHash("0x04"),
		Timestamp:      time.Now().Truncate(time.Second),
		WIP:            true,
	}
	data := l2_shared.ProcessData{
		BatchNumber:       123,
		Mode:              l2_shared.FullProcessMode,
		BatchMustBeClosed: false,
		DebugPrefix:       "test",
		StateBatch:        nil,
		TrustedBatch: &types.Batch{
			Number:         123,
			Coinbase:       expectedBatch.Coinbase,
			StateRoot:      expectedBatch.StateRoot,
			GlobalExitRoot: expectedBatch.GlobalExitRoot,
			LocalExitRoot:  expectedBatch.LocalExitRoot,
			Timestamp:      (types.ArgUint64)(expectedBatch.Timestamp.Unix()),
			Closed:         false,
		},
	}
	testData.stateMock.EXPECT().OpenBatch(testData.ctx, mock.Anything, mock.Anything).Return(nil).Once()
	testData.stateMock.EXPECT().UpdateWIPBatch(testData.ctx, mock.Anything, mock.Anything).Return(nil).Once()

	response, err := testData.sut.FullProcess(testData.ctx, &data, nil)
	require.NoError(t, err)
	require.Equal(t, false, response.ClearCache)
	require.Equal(t, false, response.UpdateBatchWithProcessBatchResponse)
	require.Equal(t, true, response.UpdateBatch.WIP)
	require.Equal(t, 0, len(response.UpdateBatch.BatchL2Data))
	require.Equal(t, expectedBatch, *response.UpdateBatch)
}

func TestEmptyBatchClosed(t *testing.T) {
	testData := newTestData(t)
	// Arrange
	expectedBatch := state.Batch{
		BatchNumber:    123,
		Coinbase:       common.HexToAddress("0x01"),
		StateRoot:      common.HexToHash("0x02"),
		GlobalExitRoot: common.HexToHash("0x03"),
		LocalExitRoot:  common.HexToHash("0x04"),
		Timestamp:      time.Now().Truncate(time.Second),
		WIP:            false,
	}
	data := l2_shared.ProcessData{
		BatchNumber:       123,
		Mode:              l2_shared.FullProcessMode,
		BatchMustBeClosed: true,
		DebugPrefix:       "test",
		StateBatch:        nil,
		TrustedBatch: &types.Batch{
			Number:         123,
			Coinbase:       expectedBatch.Coinbase,
			StateRoot:      expectedBatch.StateRoot,
			GlobalExitRoot: expectedBatch.GlobalExitRoot,
			LocalExitRoot:  expectedBatch.LocalExitRoot,
			Timestamp:      (types.ArgUint64)(expectedBatch.Timestamp.Unix()),
			Closed:         true,
		},
	}
	testData.stateMock.EXPECT().OpenBatch(testData.ctx, mock.Anything, mock.Anything).Return(nil).Once()
	testData.stateMock.EXPECT().CloseBatch(testData.ctx, mock.Anything, mock.Anything).Return(nil).Once()

	response, err := testData.sut.FullProcess(testData.ctx, &data, nil)
	require.NoError(t, err)
	require.Equal(t, false, response.ClearCache)
	require.Equal(t, false, response.UpdateBatchWithProcessBatchResponse)
	require.Equal(t, false, response.UpdateBatch.WIP)
	require.Equal(t, 0, len(response.UpdateBatch.BatchL2Data))
	require.Equal(t, expectedBatch, *response.UpdateBatch)
}
