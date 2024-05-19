package l1_check_block_test

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/state"
	commonsync "github.com/0xPolygonHermez/zkevm-node/synchronizer/common"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/l1_check_block"
	mock_l1_check_block "github.com/0xPolygonHermez/zkevm-node/synchronizer/l1_check_block/mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type testData struct {
	mockL1Client         *mock_l1_check_block.L1Requester
	mockState            *mock_l1_check_block.StateInterfacer
	mockBlockNumberFetch *mock_l1_check_block.SafeL1BlockNumberFetcher
	sut                  *l1_check_block.CheckL1BlockHash
	ctx                  context.Context
	stateBlock           *state.Block
}

func newTestData(t *testing.T) *testData {
	mockL1Client := mock_l1_check_block.NewL1Requester(t)
	mockState := mock_l1_check_block.NewStateInterfacer(t)
	mockBlockNumberFetch := mock_l1_check_block.NewSafeL1BlockNumberFetcher(t)
	mockBlockNumberFetch.EXPECT().Description().Return("mock").Maybe()
	sut := l1_check_block.NewCheckL1BlockHash(mockL1Client, mockState, mockBlockNumberFetch)
	require.NotNil(t, sut)
	ctx := context.Background()
	return &testData{
		mockL1Client:         mockL1Client,
		mockState:            mockState,
		mockBlockNumberFetch: mockBlockNumberFetch,
		sut:                  sut,
		ctx:                  ctx,
		stateBlock: &state.Block{
			BlockNumber: 1234,
			BlockHash:   common.HexToHash("0xb07e1289b32edefd8f3c702d016fb73c81d5950b2ebc790ad9d2cb8219066b4c"),
		},
	}
}

func TestCheckL1BlockHashNoBlocksOnDB(t *testing.T) {
	data := newTestData(t)
	data.mockState.EXPECT().GetFirstUncheckedBlock(data.ctx, uint64(0), nil).Return(nil, state.ErrNotFound)
	res := data.sut.Step(data.ctx)
	require.NoError(t, res)
}

func TestCheckL1BlockHashErrorGettingFirstUncheckedBlockFromDB(t *testing.T) {
	data := newTestData(t)
	data.mockState.EXPECT().GetFirstUncheckedBlock(data.ctx, uint64(0), nil).Return(nil, fmt.Errorf("error"))
	res := data.sut.Step(data.ctx)
	require.Error(t, res)
}

func TestCheckL1BlockHashErrorGettingGetSafeBlockNumber(t *testing.T) {
	data := newTestData(t)

	data.mockState.EXPECT().GetFirstUncheckedBlock(data.ctx, uint64(0), nil).Return(data.stateBlock, nil)
	data.mockBlockNumberFetch.EXPECT().GetSafeBlockNumber(data.ctx, data.mockL1Client).Return(uint64(0), fmt.Errorf("error"))
	res := data.sut.Step(data.ctx)
	require.Error(t, res)
}

// The first block to check is below the safe point, nothing to do
func TestCheckL1BlockHashSafePointIsInFuture(t *testing.T) {
	data := newTestData(t)

	data.mockState.EXPECT().GetFirstUncheckedBlock(data.ctx, uint64(0), nil).Return(data.stateBlock, nil)
	data.mockBlockNumberFetch.EXPECT().GetSafeBlockNumber(data.ctx, data.mockL1Client).Return(data.stateBlock.BlockNumber-1, nil)

	res := data.sut.Step(data.ctx)
	require.NoError(t, res)
}

func TestCheckL1BlockHashL1ClientReturnsANil(t *testing.T) {
	data := newTestData(t)

	data.mockState.EXPECT().GetFirstUncheckedBlock(data.ctx, uint64(0), nil).Return(data.stateBlock, nil)
	data.mockBlockNumberFetch.EXPECT().GetSafeBlockNumber(data.ctx, data.mockL1Client).Return(data.stateBlock.BlockNumber+10, nil)
	data.mockL1Client.EXPECT().HeaderByNumber(data.ctx, big.NewInt(int64(data.stateBlock.BlockNumber))).Return(nil, nil)
	res := data.sut.Step(data.ctx)
	require.Error(t, res)
}

// Check a block that is OK
func TestCheckL1BlockHashMatchHashUpdateCheckMarkOnDB(t *testing.T) {
	data := newTestData(t)

	data.mockState.EXPECT().GetFirstUncheckedBlock(data.ctx, uint64(0), nil).Return(data.stateBlock, nil)
	data.mockBlockNumberFetch.EXPECT().Description().Return("mock")
	data.mockBlockNumberFetch.EXPECT().GetSafeBlockNumber(data.ctx, data.mockL1Client).Return(data.stateBlock.BlockNumber, nil)
	l1Block := &types.Header{
		Number: big.NewInt(100),
	}
	data.mockL1Client.EXPECT().HeaderByNumber(data.ctx, big.NewInt(int64(data.stateBlock.BlockNumber))).Return(l1Block, nil)
	data.mockState.EXPECT().UpdateCheckedBlockByNumber(data.ctx, data.stateBlock.BlockNumber, true, nil).Return(nil)
	data.mockState.EXPECT().GetFirstUncheckedBlock(data.ctx, mock.Anything, nil).Return(nil, state.ErrNotFound)

	res := data.sut.Step(data.ctx)
	require.NoError(t, res)
}

// The first block to check is equal to the safe point, must be processed
func TestCheckL1BlockHashMismatch(t *testing.T) {
	data := newTestData(t)

	data.mockState.EXPECT().GetFirstUncheckedBlock(data.ctx, uint64(0), nil).Return(data.stateBlock, nil)
	data.stateBlock.BlockHash = common.HexToHash("0x1234") // Wrong hash to trigger a mismatch
	data.mockBlockNumberFetch.EXPECT().GetSafeBlockNumber(data.ctx, data.mockL1Client).Return(data.stateBlock.BlockNumber, nil)
	l1Block := &types.Header{
		Number: big.NewInt(100),
	}
	data.mockL1Client.EXPECT().HeaderByNumber(data.ctx, big.NewInt(int64(data.stateBlock.BlockNumber))).Return(l1Block, nil)

	res := data.sut.Step(data.ctx)
	require.Error(t, res)
	resErr, ok := res.(*commonsync.ReorgError)
	require.True(t, ok)
	require.Equal(t, data.stateBlock.BlockNumber, resErr.BlockNumber)
}
