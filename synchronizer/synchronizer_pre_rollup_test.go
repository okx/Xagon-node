package synchronizer

import (
	"context"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	mock_syncinterfaces "github.com/0xPolygonHermez/zkevm-node/synchronizer/common/syncinterfaces/mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSyncPreRollupProcessL1InfoRootEventsAskForAllBlocks(t *testing.T) {
	mockProcessor := mock_syncinterfaces.NewBlockRangeProcessor(t)
	mockEtherman := mock_syncinterfaces.NewEthermanFullInterface(t)
	sync := &SyncPreRollup{
		etherman:            mockEtherman,
		blockRangeProcessor: mockProcessor,
		SyncChunkSize:       10,
		GenesisBlockNumber:  1234,
	}

	ctx := context.Background()
	fromBlock := uint64(1)
	toBlock := uint64(31)
	syncChunkSize := uint64(10)
	previousBlockNumber := uint64(1)
	for _, i := range []uint64{10, 20, 30, 31} {
		// Mocking the call to GetRollupInfoByBlockRangePreviousRollupGenesis
		v := i
		mockEtherman.EXPECT().GetRollupInfoByBlockRangePreviousRollupGenesis(ctx, previousBlockNumber, &v).
			Return(getRollupTest()).Once()
		previousBlockNumber = i + 1
	}

	mockProcessor.EXPECT().ProcessBlockRange(ctx, mock.Anything, mock.Anything).Return(nil).Maybe()
	err := sync.ProcessL1InfoRootEvents(ctx, fromBlock, toBlock, syncChunkSize)
	require.NoError(t, err)
}

func getRollupTest() ([]etherman.Block, map[common.Hash][]etherman.Order, error) {
	return nil, nil, nil
}

func TestSyncPreRollupGetStartingL1Block(t *testing.T) {
	mockState := mock_syncinterfaces.NewStateFullInterface(t)
	mockEtherman := mock_syncinterfaces.NewEthermanFullInterface(t)
	sync := &SyncPreRollup{
		state:              mockState,
		etherman:           mockEtherman,
		GenesisBlockNumber: 1234,
	}

	ctx := context.Background()

	for idx, testCase := range []struct {
		name                   string
		upgradeLxLyBlockNumber uint64
		blockNumber            uint64
		expectedError          bool
		expectedNeedToUpdate   bool
		expectedBlockNumber    uint64
	}{
		{name: "mid block", upgradeLxLyBlockNumber: 1000, blockNumber: 1001, expectedError: false, expectedNeedToUpdate: true, expectedBlockNumber: 1001},
		{name: "pre block", upgradeLxLyBlockNumber: 1000, blockNumber: 999, expectedError: false, expectedNeedToUpdate: true, expectedBlockNumber: 999},
		{name: "same genesis", upgradeLxLyBlockNumber: 1000, blockNumber: sync.GenesisBlockNumber, expectedError: false, expectedNeedToUpdate: false},
		{name: "genesis-1", upgradeLxLyBlockNumber: 1000, blockNumber: 1233, expectedError: false, expectedNeedToUpdate: false},
	} {
		log.Info("Running test case ", idx+1)
		block := state.Block{
			BlockNumber: testCase.blockNumber,
		}
		mockEtherman.EXPECT().GetL1BlockUpgradeLxLy(ctx, sync.GenesisBlockNumber).Return(testCase.upgradeLxLyBlockNumber, nil).Maybe()
		mockState.EXPECT().GetLastBlock(ctx, mock.Anything).Return(&block, nil).Once()
		needToUpdate, blockNumber, err := sync.getStartingL1Block(ctx, nil)
		if testCase.expectedError {
			require.Error(t, err, testCase.name)
		} else {
			require.NoError(t, err, testCase.name)
			require.Equal(t, testCase.expectedNeedToUpdate, needToUpdate, testCase.name)
			if needToUpdate {
				require.Equal(t, testCase.blockNumber, blockNumber, testCase.name)
			}
		}
	}
}
