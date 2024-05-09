package l1_check_block_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/common/syncinterfaces"
	mock_syncinterfaces "github.com/0xPolygonHermez/zkevm-node/synchronizer/common/syncinterfaces/mocks"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/l1_check_block"
	mock_l1_check_block "github.com/0xPolygonHermez/zkevm-node/synchronizer/l1_check_block/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	genericErrorToTest = fmt.Errorf("error")
)

type testDataIntegration struct {
	mockChecker    *mock_syncinterfaces.AsyncL1BlockChecker
	mockPreChecker *mock_syncinterfaces.AsyncL1BlockChecker
	mockState      *mock_l1_check_block.StateForL1BlockCheckerIntegration
	mockSync       *mock_l1_check_block.SyncCheckReorger
	sut            *l1_check_block.L1BlockCheckerIntegration
	ctx            context.Context
	resultOk       syncinterfaces.IterationResult
	resultError    syncinterfaces.IterationResult
	resultReorg    syncinterfaces.IterationResult
}

func newDataIntegration(t *testing.T, forceCheckOnStart bool) *testDataIntegration {
	return newDataIntegrationOnlyMainChecker(t, forceCheckOnStart)
}

func newDataIntegrationWithPreChecker(t *testing.T, forceCheckOnStart bool) *testDataIntegration {
	res := newDataIntegrationOnlyMainChecker(t, forceCheckOnStart)
	res.mockPreChecker = mock_syncinterfaces.NewAsyncL1BlockChecker(t)
	res.sut = l1_check_block.NewL1BlockCheckerIntegration(res.mockChecker, res.mockPreChecker, res.mockState, res.mockSync, forceCheckOnStart, time.Millisecond)
	return res
}

func newDataIntegrationOnlyMainChecker(t *testing.T, forceCheckOnStart bool) *testDataIntegration {
	mockChecker := mock_syncinterfaces.NewAsyncL1BlockChecker(t)
	mockSync := mock_l1_check_block.NewSyncCheckReorger(t)
	mockState := mock_l1_check_block.NewStateForL1BlockCheckerIntegration(t)
	sut := l1_check_block.NewL1BlockCheckerIntegration(mockChecker, nil, mockState, mockSync, forceCheckOnStart, time.Millisecond)
	return &testDataIntegration{
		mockChecker:    mockChecker,
		mockPreChecker: nil,
		mockSync:       mockSync,
		mockState:      mockState,
		sut:            sut,
		ctx:            context.Background(),
		resultReorg: syncinterfaces.IterationResult{
			ReorgDetected: true,
			BlockNumber:   1234,
		},
		resultOk: syncinterfaces.IterationResult{
			ReorgDetected: false,
		},
		resultError: syncinterfaces.IterationResult{
			Err:           genericErrorToTest,
			ReorgDetected: false,
		},
	}
}

func TestIntegrationIfNoForceCheckOnlyLaunchBackgroudChecker(t *testing.T) {
	data := newDataIntegration(t, false)
	data.mockChecker.EXPECT().Run(data.ctx, mock.Anything).Return()
	err := data.sut.OnStart(data.ctx)
	require.NoError(t, err)
}

func TestIntegrationIfForceCheckRunsSynchronousOneTimeAndAfterLaunchBackgroudChecker(t *testing.T) {
	data := newDataIntegration(t, true)
	data.mockChecker.EXPECT().RunSynchronous(data.ctx).Return(data.resultOk)
	data.mockChecker.EXPECT().Run(data.ctx, mock.Anything).Return()
	err := data.sut.OnStart(data.ctx)
	require.NoError(t, err)
}

func TestIntegrationIfSyncCheckReturnsReorgExecuteIt(t *testing.T) {
	data := newDataIntegration(t, true)
	data.mockChecker.EXPECT().RunSynchronous(data.ctx).Return(data.resultReorg)
	data.mockSync.EXPECT().ExecuteReorgFromMismatchBlock(uint64(1234), "").Return(nil)
	data.mockChecker.EXPECT().Run(data.ctx, mock.Anything).Return()
	err := data.sut.OnStart(data.ctx)
	require.NoError(t, err)
}

func TestIntegrationIfSyncCheckReturnErrorRetry(t *testing.T) {
	data := newDataIntegration(t, true)
	data.mockChecker.EXPECT().RunSynchronous(data.ctx).Return(data.resultError).Once()
	data.mockChecker.EXPECT().RunSynchronous(data.ctx).Return(data.resultOk).Once()
	data.mockChecker.EXPECT().Run(data.ctx, mock.Anything).Return()
	err := data.sut.OnStart(data.ctx)
	require.NoError(t, err)
}

func TestIntegrationIfSyncCheckReturnsReorgExecuteItAndFailsRetry(t *testing.T) {
	data := newDataIntegration(t, true)
	data.mockChecker.EXPECT().RunSynchronous(data.ctx).Return(data.resultReorg)
	data.mockSync.EXPECT().ExecuteReorgFromMismatchBlock(uint64(1234), mock.Anything).Return(genericErrorToTest).Once()
	data.mockSync.EXPECT().ExecuteReorgFromMismatchBlock(uint64(1234), mock.Anything).Return(nil).Once()
	data.mockChecker.EXPECT().Run(data.ctx, mock.Anything).Return()
	err := data.sut.OnStart(data.ctx)
	require.NoError(t, err)
}

// OnStart if check and preCheck execute both, and launch both in background
func TestIntegrationCheckAndPreCheckOnStartForceCheck(t *testing.T) {
	data := newDataIntegrationWithPreChecker(t, true)
	data.mockChecker.EXPECT().RunSynchronous(data.ctx).Return(data.resultOk)
	data.mockPreChecker.EXPECT().RunSynchronous(data.ctx).Return(data.resultOk)
	data.mockChecker.EXPECT().Run(data.ctx, mock.Anything).Return()
	data.mockPreChecker.EXPECT().Run(data.ctx, mock.Anything).Return()
	err := data.sut.OnStart(data.ctx)
	require.NoError(t, err)
}

// OnStart if mainChecker returns reorg doesnt need to run preCheck
func TestIntegrationCheckAndPreCheckOnStartMainCheckerReturnReorg(t *testing.T) {
	data := newDataIntegrationWithPreChecker(t, true)
	data.mockChecker.EXPECT().RunSynchronous(data.ctx).Return(data.resultReorg)
	data.mockSync.EXPECT().ExecuteReorgFromMismatchBlock(uint64(1234), mock.Anything).Return(nil).Once()
	data.mockChecker.EXPECT().Run(data.ctx, mock.Anything).Return()
	data.mockPreChecker.EXPECT().Run(data.ctx, mock.Anything).Return()
	err := data.sut.OnStart(data.ctx)
	require.NoError(t, err)
}

// If mainCheck is OK, but preCheck returns reorg, it should execute reorg
func TestIntegrationCheckAndPreCheckOnStartPreCheckerReturnReorg(t *testing.T) {
	data := newDataIntegrationWithPreChecker(t, true)
	data.mockChecker.EXPECT().RunSynchronous(data.ctx).Return(data.resultOk)
	data.mockPreChecker.EXPECT().RunSynchronous(data.ctx).Return(data.resultReorg)
	data.mockSync.EXPECT().ExecuteReorgFromMismatchBlock(uint64(1234), mock.Anything).Return(nil).Once()
	data.mockChecker.EXPECT().Run(data.ctx, mock.Anything).Return()
	data.mockPreChecker.EXPECT().Run(data.ctx, mock.Anything).Return()
	err := data.sut.OnStart(data.ctx)
	require.NoError(t, err)
}

// The process is running on background, no results yet
func TestIntegrationCheckAndPreCheckOnOnCheckReorgRunningOnBackground(t *testing.T) {
	data := newDataIntegrationWithPreChecker(t, true)
	data.mockChecker.EXPECT().GetResult().Return(nil)
	data.mockPreChecker.EXPECT().GetResult().Return(nil)
	block, err := data.sut.CheckReorgWrapper(data.ctx, nil, nil)
	require.Nil(t, block)
	require.NoError(t, err)
}

func TestIntegrationCheckAndPreCheckOnOnCheckReorgOneProcessHaveResultOK(t *testing.T) {
	data := newDataIntegrationWithPreChecker(t, true)
	data.mockChecker.EXPECT().GetResult().Return(&data.resultOk)
	data.mockPreChecker.EXPECT().GetResult().Return(nil)
	// One have been stopped, so must relaunch both
	data.mockChecker.EXPECT().Run(data.ctx, mock.Anything).Return()
	data.mockPreChecker.EXPECT().Run(data.ctx, mock.Anything).Return()
	block, err := data.sut.CheckReorgWrapper(data.ctx, nil, nil)
	require.Nil(t, block)
	require.NoError(t, err)
}

func TestIntegrationCheckAndPreCheckOnOnCheckReorgMainCheckerReorg(t *testing.T) {
	data := newDataIntegrationWithPreChecker(t, true)
	data.mockChecker.EXPECT().GetResult().Return(&data.resultReorg)
	data.mockPreChecker.EXPECT().GetResult().Return(nil)
	data.mockState.EXPECT().GetPreviousBlockToBlockNumber(data.ctx, uint64(1234), nil).Return(&state.Block{
		BlockNumber: data.resultReorg.BlockNumber - 1,
	}, nil)
	// One have been stopped,but is going to be launched OnResetState call after the reset
	block, err := data.sut.CheckReorgWrapper(data.ctx, nil, nil)
	require.NotNil(t, block)
	require.Equal(t, data.resultReorg.BlockNumber-1, block.BlockNumber)
	require.NoError(t, err)
}

func TestIntegrationCheckAndPreCheckOnOnCheckReorgPreCheckerReorg(t *testing.T) {
	data := newDataIntegrationWithPreChecker(t, true)
	data.mockChecker.EXPECT().GetResult().Return(nil)
	data.mockPreChecker.EXPECT().GetResult().Return(&data.resultReorg)
	data.mockState.EXPECT().GetPreviousBlockToBlockNumber(data.ctx, uint64(1234), nil).Return(&state.Block{
		BlockNumber: data.resultReorg.BlockNumber - 1,
	}, nil)
	// One have been stopped,but is going to be launched OnResetState call after the reset

	block, err := data.sut.CheckReorgWrapper(data.ctx, nil, nil)
	require.NotNil(t, block)
	require.Equal(t, data.resultReorg.BlockNumber-1, block.BlockNumber)
	require.NoError(t, err)
}

func TestIntegrationCheckAndPreCheckOnOnCheckReorgBothReorgWinOldest1(t *testing.T) {
	data := newDataIntegrationWithPreChecker(t, true)
	reorgMain := data.resultReorg
	reorgMain.BlockNumber = 1235
	data.mockChecker.EXPECT().GetResult().Return(&reorgMain)
	reorgPre := data.resultReorg
	reorgPre.BlockNumber = 1236
	data.mockPreChecker.EXPECT().GetResult().Return(&reorgPre)
	data.mockState.EXPECT().GetPreviousBlockToBlockNumber(data.ctx, uint64(1235), nil).Return(&state.Block{
		BlockNumber: 1234,
	}, nil)

	// Both have been stopped,but is going to be launched OnResetState call after the reset

	block, err := data.sut.CheckReorgWrapper(data.ctx, nil, nil)
	require.NotNil(t, block)
	require.Equal(t, uint64(1234), block.BlockNumber)
	require.NoError(t, err)
}

func TestIntegrationCheckAndPreCheckOnOnCheckReorgBothReorgWinOldest2(t *testing.T) {
	data := newDataIntegrationWithPreChecker(t, true)
	reorgMain := data.resultReorg
	reorgMain.BlockNumber = 1236
	data.mockChecker.EXPECT().GetResult().Return(&reorgMain)
	reorgPre := data.resultReorg
	reorgPre.BlockNumber = 1235
	data.mockPreChecker.EXPECT().GetResult().Return(&reorgPre)
	data.mockState.EXPECT().GetPreviousBlockToBlockNumber(data.ctx, uint64(1235), nil).Return(&state.Block{
		BlockNumber: 1234,
	}, nil)
	// Both have been stopped,but is going to be launched OnResetState call after the reset

	block, err := data.sut.CheckReorgWrapper(data.ctx, nil, nil)
	require.NotNil(t, block)
	require.Equal(t, uint64(1234), block.BlockNumber)
	require.NoError(t, err)
}

func TestIntegrationCheckReorgWrapperBypassReorgFuncIfNoBackgroundData(t *testing.T) {
	data := newDataIntegrationWithPreChecker(t, true)
	data.mockChecker.EXPECT().GetResult().Return(nil)
	data.mockPreChecker.EXPECT().GetResult().Return(nil)
	reorgFuncBlock := &state.Block{
		BlockNumber: 1234,
	}
	reorgFuncErr := fmt.Errorf("error")
	block, err := data.sut.CheckReorgWrapper(data.ctx, reorgFuncBlock, reorgFuncErr)
	require.Equal(t, reorgFuncBlock, block)
	require.Equal(t, reorgFuncErr, err)
}

func TestIntegrationCheckReorgWrapperChooseOldestReorgFunc(t *testing.T) {
	data := newDataIntegrationWithPreChecker(t, true)
	data.mockChecker.EXPECT().GetResult().Return(nil)
	data.mockPreChecker.EXPECT().GetResult().Return(&data.resultReorg)
	data.mockState.EXPECT().GetPreviousBlockToBlockNumber(data.ctx, uint64(1234), nil).Return(&state.Block{
		BlockNumber: 1233,
	}, nil)

	reorgFuncBlock := &state.Block{
		BlockNumber: 1230,
	}
	block, err := data.sut.CheckReorgWrapper(data.ctx, reorgFuncBlock, nil)
	require.Equal(t, reorgFuncBlock, block)
	require.NoError(t, err)
}

func TestIntegrationCheckReorgWrapperChooseOldestBackgroundCheck(t *testing.T) {
	data := newDataIntegrationWithPreChecker(t, true)
	data.mockChecker.EXPECT().GetResult().Return(nil)
	data.mockPreChecker.EXPECT().GetResult().Return(&data.resultReorg)
	data.mockState.EXPECT().GetPreviousBlockToBlockNumber(data.ctx, uint64(1234), nil).Return(&state.Block{
		BlockNumber: 1233,
	}, nil)

	reorgFuncBlock := &state.Block{
		BlockNumber: 1240,
	}
	block, err := data.sut.CheckReorgWrapper(data.ctx, reorgFuncBlock, nil)
	require.Equal(t, uint64(1233), block.BlockNumber)
	require.NoError(t, err)
}

func TestIntegrationCheckReorgWrapperIgnoreReorgFuncIfError(t *testing.T) {
	data := newDataIntegrationWithPreChecker(t, true)
	data.mockChecker.EXPECT().GetResult().Return(nil)
	data.mockPreChecker.EXPECT().GetResult().Return(&data.resultReorg)
	data.mockState.EXPECT().GetPreviousBlockToBlockNumber(data.ctx, uint64(1234), nil).Return(&state.Block{
		BlockNumber: 1233,
	}, nil)

	reorgFuncBlock := &state.Block{
		BlockNumber: 1230,
	}
	reorgFuncErr := fmt.Errorf("error")
	block, err := data.sut.CheckReorgWrapper(data.ctx, reorgFuncBlock, reorgFuncErr)
	require.Equal(t, uint64(1233), block.BlockNumber)
	require.NoError(t, err)
}
