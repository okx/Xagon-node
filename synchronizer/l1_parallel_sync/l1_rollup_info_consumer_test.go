package l1_parallel_sync

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	mock_syncinterfaces "github.com/0xPolygonHermez/zkevm-node/synchronizer/common/syncinterfaces/mocks"
	"github.com/ethereum/go-ethereum/common"
	types "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type consumerTestData struct {
	sut      *l1RollupInfoConsumer
	syncMock *mock_syncinterfaces.BlockRangeProcessor
	ch       chan L1SyncMessage
}

func TestGivenConsumerWhenReceiveAFullSyncAndChannelIsEmptyThenStopOk(t *testing.T) {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	data := setupConsumerTest(t)
	defer cancel()
	data.ch <- *newL1SyncMessageControl(eventProducerIsFullySynced)
	err := data.sut.Start(ctxTimeout, nil)
	require.NoError(t, err)
}
func TestGivenConsumerWhenReceiveAFullSyncAndChannelIsNotEmptyThenDontStop(t *testing.T) {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	data := setupConsumerTest(t)
	defer cancel()

	data.ch <- *newL1SyncMessageControl(eventProducerIsFullySynced)
	data.ch <- *newL1SyncMessageControl(eventNone)
	err := data.sut.Start(ctxTimeout, nil)
	require.Error(t, err)
	require.Equal(t, errContextCanceled, err)
}

func TestGivenConsumerWhenFailsToProcessRollupThenDontKnownLastEthBlock(t *testing.T) {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	data := setupConsumerTest(t)
	defer cancel()
	responseRollupInfoByBlockRange := rollupInfoByBlockRangeResult{
		blockRange: blockRange{
			fromBlock: 100,
			toBlock:   200,
		},
		blocks:           []etherman.Block{},
		order:            map[common.Hash][]etherman.Order{},
		lastBlockOfRange: types.NewBlock(&types.Header{Number: big.NewInt(123)}, nil, nil, nil, nil),
	}
	data.syncMock.
		On("ProcessBlockRange", mock.Anything, mock.Anything, mock.Anything).
		Return(errors.New("error")).
		Once()
	data.ch <- *newL1SyncMessageData(&responseRollupInfoByBlockRange)
	data.ch <- *newL1SyncMessageControl(eventProducerIsFullySynced)
	err := data.sut.Start(ctxTimeout, nil)
	require.Error(t, err)
	_, ok := data.sut.GetLastEthBlockSynced()
	require.False(t, ok)
}

func TestGivenConsumerWhenReceiveNoNextBlockThenDoNothing(t *testing.T) {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	data := setupConsumerTest(t)
	defer cancel()
	responseRollupInfoByBlockRange := rollupInfoByBlockRangeResult{
		blockRange: blockRange{
			fromBlock: 100,
			toBlock:   200,
		},
		blocks:           []etherman.Block{},
		order:            map[common.Hash][]etherman.Order{},
		lastBlockOfRange: nil,
	}
	data.sut.Reset(1234)
	// Is not going to call processBlockRange because is not expected
	data.ch <- *newL1SyncMessageData(&responseRollupInfoByBlockRange)
	data.ch <- *newL1SyncMessageControlWProducerIsFullySynced(1234)

	err := data.sut.Start(ctxTimeout, nil)
	require.NoError(t, err)
	_, ok := data.sut.GetLastEthBlockSynced()
	require.False(t, ok)
}

func TestGivenConsumerWhenNextBlockNumberIsNoSetThenAcceptAnythingAndProcess(t *testing.T) {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	data := setupConsumerTest(t)
	defer cancel()
	responseRollupInfoByBlockRange := rollupInfoByBlockRangeResult{
		blockRange: blockRange{
			fromBlock: 100,
			toBlock:   200,
		},
		blocks:           []etherman.Block{},
		order:            map[common.Hash][]etherman.Order{},
		lastBlockOfRange: types.NewBlock(&types.Header{Number: big.NewInt(123)}, nil, nil, nil, nil),
	}

	data.ch <- *newL1SyncMessageData(&responseRollupInfoByBlockRange)
	data.ch <- *newL1SyncMessageControlWProducerIsFullySynced(200)
	data.syncMock.
		On("ProcessBlockRange", mock.Anything, mock.Anything, mock.Anything).
		Return(nil).
		Once()
	err := data.sut.Start(ctxTimeout, nil)
	require.NoError(t, err)
	resultBlock, ok := data.sut.GetLastEthBlockSynced()
	require.True(t, ok)
	require.Equal(t, uint64(123), resultBlock.BlockNumber)
}

func TestGivenConsumerWhenNextBlockNumberIsNoSetThenAcceptAnythingAndProcessAndConsumerAreDesynchronizer(t *testing.T) {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	data := setupConsumerTest(t)
	defer cancel()
	responseRollupInfoByBlockRange := rollupInfoByBlockRangeResult{
		blockRange: blockRange{
			fromBlock: 100,
			toBlock:   200,
		},
		blocks:           []etherman.Block{},
		order:            map[common.Hash][]etherman.Order{},
		lastBlockOfRange: types.NewBlock(&types.Header{Number: big.NewInt(123)}, nil, nil, nil, nil),
	}

	data.ch <- *newL1SyncMessageData(&responseRollupInfoByBlockRange)
	data.ch <- *newL1SyncMessageControlWProducerIsFullySynced(300)
	data.syncMock.
		On("ProcessBlockRange", mock.Anything, mock.Anything, mock.Anything).
		Return(nil).
		Once()
	err := data.sut.Start(ctxTimeout, nil)
	require.ErrorIs(t, errConsumerAndProducerDesynchronized, err)
	resultBlock, ok := data.sut.GetLastEthBlockSynced()
	require.True(t, ok)
	require.Equal(t, uint64(123), resultBlock.BlockNumber)
}

func TestGivenConsumerWhenNextBlockNumberIsNoSetThenFirstRollupInfoSetIt(t *testing.T) {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	data := setupConsumerTest(t)
	defer cancel()
	responseRollupInfoByBlockRange := rollupInfoByBlockRangeResult{
		blockRange: blockRange{
			fromBlock: 100,
			toBlock:   200,
		},
		blocks:           []etherman.Block{},
		order:            map[common.Hash][]etherman.Order{},
		lastBlockOfRange: types.NewBlock(&types.Header{Number: big.NewInt(123)}, nil, nil, nil, nil),
	}
	// Fist package set highestBlockProcessed
	data.ch <- *newL1SyncMessageData(&responseRollupInfoByBlockRange)
	// The repeated package is ignored because is not the next BlockRange
	data.ch <- *newL1SyncMessageData(&responseRollupInfoByBlockRange)
	data.ch <- *newL1SyncMessageControlWProducerIsFullySynced(200)
	data.syncMock.
		On("ProcessBlockRange", mock.Anything, mock.Anything, mock.Anything).
		Return(nil).
		Once()
	err := data.sut.Start(ctxTimeout, nil)
	require.NoError(t, err)
	resultBlock, ok := data.sut.GetLastEthBlockSynced()
	require.True(t, ok)
	require.Equal(t, uint64(123), resultBlock.BlockNumber)
}

func TestGivenProducerDesyncrhonizedOnHeadL1(t *testing.T) {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	defer cancel()
	data := setupConsumerTest(t)
	responseRollupInfoByBlockRange := rollupInfoByBlockRangeResult{
		blockRange: blockRange{
			fromBlock: 100,
			toBlock:   200,
		},
		blocks:           []etherman.Block{},
		order:            map[common.Hash][]etherman.Order{},
		lastBlockOfRange: types.NewBlock(&types.Header{Number: big.NewInt(123)}, nil, nil, nil, nil),
	}
	// Fist package set highestBlockProcessed
	data.ch <- *newL1SyncMessageData(&responseRollupInfoByBlockRange)
	responseRollupInfoByBlockRange.blockRange.fromBlock = 300
	responseRollupInfoByBlockRange.blockRange.toBlock = 400
	data.ch <- *newL1SyncMessageData(&responseRollupInfoByBlockRange)
	data.ch <- *newL1SyncMessageControlWProducerIsFullySynced(200)
	data.syncMock.EXPECT().ProcessBlockRange(mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(1)
	err := data.sut.Start(ctxTimeout, nil)
	require.NoError(t, err)
}

func TestGivenConsumerWhenNextBlockNumberIsNoSetDontReceiveAnyBlockButAFullSyncEvent(t *testing.T) {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), time.Millisecond*100)
	data := setupConsumerTest(t)
	defer cancel()
	data.ch <- *newL1SyncMessageControlWProducerIsFullySynced(200)
	err := data.sut.Start(ctxTimeout, nil)
	require.NoError(t, err)
}

func setupConsumerTest(t *testing.T) consumerTestData {
	syncMock := mock_syncinterfaces.NewBlockRangeProcessor(t)
	ch := make(chan L1SyncMessage, 10)

	cfg := ConfigConsumer{
		ApplyAfterNumRollupReceived: minNumIterationsBeforeStartCheckingTimeWaitingForNewRollupInfoData,
		AceptableInacctivityTime:    minAcceptableTimeWaitingForNewRollupInfoData,
	}
	sut := NewL1RollupInfoConsumer(cfg, syncMock, ch)
	return consumerTestData{sut, syncMock, ch}
}
