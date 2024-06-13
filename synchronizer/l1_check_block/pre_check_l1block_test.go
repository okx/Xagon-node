package l1_check_block_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/state"
	commonsync "github.com/0xPolygonHermez/zkevm-node/synchronizer/common"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/l1_check_block"
	mock_l1_check_block "github.com/0xPolygonHermez/zkevm-node/synchronizer/l1_check_block/mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
)

type testPreCheckData struct {
	sut              *l1_check_block.PreCheckL1BlockHash
	mockL1Client     *mock_l1_check_block.L1Requester
	mockState        *mock_l1_check_block.StatePreCheckInterfacer
	mockInitialFetch *mock_l1_check_block.SafeL1BlockNumberFetcher
	mockEndFetch     *mock_l1_check_block.SafeL1BlockNumberFetcher
	ctx              context.Context
	stateBlocks      []*state.Block
}

func newPreCheckData(t *testing.T) *testPreCheckData {
	mockL1Client := mock_l1_check_block.NewL1Requester(t)
	mockState := mock_l1_check_block.NewStatePreCheckInterfacer(t)
	mockInitialFetch := mock_l1_check_block.NewSafeL1BlockNumberFetcher(t)
	mockEndFetch := mock_l1_check_block.NewSafeL1BlockNumberFetcher(t)
	sut := l1_check_block.NewPreCheckL1BlockHash(mockL1Client, mockState, mockInitialFetch, mockEndFetch)
	return &testPreCheckData{
		sut:              sut,
		mockL1Client:     mockL1Client,
		mockState:        mockState,
		mockInitialFetch: mockInitialFetch,
		mockEndFetch:     mockEndFetch,
		ctx:              context.Background(),
		stateBlocks: []*state.Block{
			{
				BlockNumber: 1234,
				BlockHash:   common.HexToHash("0xd77dd3a9ee6f9202ca5a75024b7d9cbd3d7436b2910d450f88c261c0089c0cd9"),
			},
			{
				BlockNumber: 1237,
				BlockHash:   common.HexToHash("0x8faffac37f561c18917c33ff3540262ecfbe11a367b4e1c48181326cd8ba347f"),
			},
		},
	}
}

// If from > to, it ignore because there are no blocks to check
func TestPreCheckL1BlockFromGreaterThanTo(t *testing.T) {
	data := newPreCheckData(t)
	data.mockInitialFetch.EXPECT().Description().Return("initial")
	data.mockEndFetch.EXPECT().Description().Return("end")
	data.mockInitialFetch.EXPECT().GetSafeBlockNumber(data.ctx, data.mockL1Client).Return(uint64(1234), nil)
	data.mockEndFetch.EXPECT().GetSafeBlockNumber(data.ctx, data.mockL1Client).Return(uint64(1230), nil)

	res := data.sut.Step(data.ctx)
	require.NoError(t, res)
}

// No blocks on state -> nothing to do
func TestPreCheckL1BlockNoBlocksOnState(t *testing.T) {
	data := newPreCheckData(t)
	data.mockInitialFetch.EXPECT().Description().Return("initial")
	data.mockEndFetch.EXPECT().Description().Return("end")
	data.mockInitialFetch.EXPECT().GetSafeBlockNumber(data.ctx, data.mockL1Client).Return(uint64(1234), nil)
	data.mockEndFetch.EXPECT().GetSafeBlockNumber(data.ctx, data.mockL1Client).Return(uint64(1250), nil)
	data.mockState.EXPECT().GetUncheckedBlocks(data.ctx, uint64(1234), uint64(1250), nil).Return(nil, nil)

	res := data.sut.Step(data.ctx)
	require.NoError(t, res)
}

func TestPreCheckL1BlockBlocksMatch(t *testing.T) {
	data := newPreCheckData(t)
	data.mockInitialFetch.EXPECT().Description().Return("initial")
	data.mockEndFetch.EXPECT().Description().Return("end")
	data.mockInitialFetch.EXPECT().GetSafeBlockNumber(data.ctx, data.mockL1Client).Return(uint64(1234), nil)
	data.mockEndFetch.EXPECT().GetSafeBlockNumber(data.ctx, data.mockL1Client).Return(uint64(1250), nil)
	data.mockState.EXPECT().GetUncheckedBlocks(data.ctx, uint64(1234), uint64(1250), nil).Return(data.stateBlocks, nil)
	l1Block1 := &types.Header{
		Number: big.NewInt(int64(data.stateBlocks[0].BlockNumber)),
	}
	data.mockL1Client.EXPECT().HeaderByNumber(data.ctx, big.NewInt(int64(data.stateBlocks[0].BlockNumber))).Return(l1Block1, nil)
	l1Block2 := &types.Header{
		Number: big.NewInt(int64(data.stateBlocks[1].BlockNumber)),
	}
	data.mockL1Client.EXPECT().HeaderByNumber(data.ctx, big.NewInt(int64(data.stateBlocks[1].BlockNumber))).Return(l1Block2, nil)
	//data.mockState.EXPECT().GetUncheckedBlocks(data.ctx, uint64(1237), uint64(1237), nil).Return(data.stateBlocks[0:1], nil)

	res := data.sut.Step(data.ctx)
	require.NoError(t, res)
}

func TestPreCheckL1BlockBlocksMismatch(t *testing.T) {
	data := newPreCheckData(t)
	data.mockInitialFetch.EXPECT().Description().Return("initial")
	data.mockEndFetch.EXPECT().Description().Return("end")
	data.mockInitialFetch.EXPECT().GetSafeBlockNumber(data.ctx, data.mockL1Client).Return(uint64(1234), nil)
	data.mockEndFetch.EXPECT().GetSafeBlockNumber(data.ctx, data.mockL1Client).Return(uint64(1250), nil)
	data.stateBlocks[1].BlockHash = common.HexToHash("0x12345678901234567890123456789012345678901234567890")
	data.mockState.EXPECT().GetUncheckedBlocks(data.ctx, uint64(1234), uint64(1250), nil).Return(data.stateBlocks, nil)
	l1Block1 := &types.Header{
		Number: big.NewInt(int64(data.stateBlocks[0].BlockNumber)),
	}
	data.mockL1Client.EXPECT().HeaderByNumber(data.ctx, big.NewInt(int64(data.stateBlocks[0].BlockNumber))).Return(l1Block1, nil)
	l1Block2 := &types.Header{
		Number: big.NewInt(int64(data.stateBlocks[1].BlockNumber)),
	}
	data.mockL1Client.EXPECT().HeaderByNumber(data.ctx, big.NewInt(int64(data.stateBlocks[1].BlockNumber))).Return(l1Block2, nil)
	data.mockState.EXPECT().GetUncheckedBlocks(data.ctx, uint64(1237), uint64(1237), nil).Return(data.stateBlocks[1:2], nil)

	res := data.sut.Step(data.ctx)
	require.Error(t, res)
	resErr, ok := res.(*commonsync.ReorgError)
	require.True(t, ok, "The error must be ReorgError")
	require.Equal(t, uint64(1237), resErr.BlockNumber)
}

func TestPreCheckL1BlockBlocksMismatchButIsNoLongerInState(t *testing.T) {
	data := newPreCheckData(t)
	data.mockInitialFetch.EXPECT().Description().Return("initial")
	data.mockEndFetch.EXPECT().Description().Return("end")
	data.mockInitialFetch.EXPECT().GetSafeBlockNumber(data.ctx, data.mockL1Client).Return(uint64(1234), nil)
	data.mockEndFetch.EXPECT().GetSafeBlockNumber(data.ctx, data.mockL1Client).Return(uint64(1250), nil)
	data.stateBlocks[1].BlockHash = common.HexToHash("0x12345678901234567890123456789012345678901234567890")
	data.mockState.EXPECT().GetUncheckedBlocks(data.ctx, uint64(1234), uint64(1250), nil).Return(data.stateBlocks, nil)
	l1Block1 := &types.Header{
		Number: big.NewInt(int64(data.stateBlocks[0].BlockNumber)),
	}
	data.mockL1Client.EXPECT().HeaderByNumber(data.ctx, big.NewInt(int64(data.stateBlocks[0].BlockNumber))).Return(l1Block1, nil)
	l1Block2 := &types.Header{
		Number: big.NewInt(int64(data.stateBlocks[1].BlockNumber)),
	}
	data.mockL1Client.EXPECT().HeaderByNumber(data.ctx, big.NewInt(int64(data.stateBlocks[1].BlockNumber))).Return(l1Block2, nil)
	data.mockState.EXPECT().GetUncheckedBlocks(data.ctx, uint64(1237), uint64(1237), nil).Return(nil, nil)

	res := data.sut.Step(data.ctx)
	require.ErrorIs(t, res, l1_check_block.ErrDeSync)
}
