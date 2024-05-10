package l2_sync_etrog

import (
	"context"
	"fmt"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/types"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/common/syncinterfaces"
	mock_syncinterfaces "github.com/0xPolygonHermez/zkevm-node/synchronizer/common/syncinterfaces/mocks"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/l2_sync/l2_shared"
	mock_l2_sync_etrog "github.com/0xPolygonHermez/zkevm-node/synchronizer/l2_sync/l2_sync_etrog/mocks"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

var (
	globalExitRootNonZero = common.HexToHash("0x723e5c4c7ee7890e1e66c2e391d553ee792d2204ecb4fe921830f12f8dcd1a92")
	randomError           = fmt.Errorf("random error")
)

type testData struct {
	ctx       context.Context
	stateMock *mock_l2_sync_etrog.StateGERInteface
	zkevmMock *mock_syncinterfaces.ZKEVMClientGlobalExitRootGetter
	sut       *CheckSyncStatusToProcessBatch
}

func NewTestData(t *testing.T) *testData {
	stateMock := mock_l2_sync_etrog.NewStateGERInteface(t)
	zkevmMock := mock_syncinterfaces.NewZKEVMClientGlobalExitRootGetter(t)

	sut := NewCheckSyncStatusToProcessBatch(zkevmMock, stateMock)
	return &testData{
		ctx:       context.Background(),
		stateMock: stateMock,
		zkevmMock: zkevmMock,
		sut:       sut,
	}
}

func TestCheckL1SyncStatusEnoughToProcessBatchGerZero(t *testing.T) {
	testData := NewTestData(t)
	err := testData.sut.CheckL1SyncGlobalExitRootEnoughToProcessBatch(testData.ctx, 1, state.ZeroHash, nil)
	require.NoError(t, err)
}
func TestCheckL1SyncStatusEnoughToProcessBatchGerOnDB(t *testing.T) {
	testData := NewTestData(t)
	testData.stateMock.EXPECT().GetExitRootByGlobalExitRoot(testData.ctx, globalExitRootNonZero, nil).Return(&state.GlobalExitRoot{}, nil).Once()
	err := testData.sut.CheckL1SyncGlobalExitRootEnoughToProcessBatch(testData.ctx, 1, globalExitRootNonZero, nil)
	require.NoError(t, err)
}

func TestCheckL1SyncStatusEnoughToProcessBatchGerDatabaseFails(t *testing.T) {
	testData := NewTestData(t)

	testData.stateMock.EXPECT().GetExitRootByGlobalExitRoot(testData.ctx, globalExitRootNonZero, nil).Return(nil, randomError).Once()

	err := testData.sut.CheckL1SyncGlobalExitRootEnoughToProcessBatch(testData.ctx, 1, globalExitRootNonZero, nil)
	require.Error(t, err)
}

func TestCheckL1SyncStatusEnoughToProcessBatchGerNoOnDBFailsCallToZkevm(t *testing.T) {
	testData := NewTestData(t)

	testData.stateMock.EXPECT().GetExitRootByGlobalExitRoot(testData.ctx, globalExitRootNonZero, nil).Return(nil, state.ErrNotFound).Once()
	testData.zkevmMock.EXPECT().ExitRootsByGER(testData.ctx, globalExitRootNonZero).Return(nil, randomError).Once()

	err := testData.sut.CheckL1SyncGlobalExitRootEnoughToProcessBatch(testData.ctx, 1, globalExitRootNonZero, nil)
	require.Error(t, err)
}

func TestCheckL1SyncStatusEnoughToProcessBatchGerNoOnDBWeAre1BlockBehind(t *testing.T) {
	testData := NewTestData(t)

	l1Block := uint64(123)
	testData.stateMock.EXPECT().GetExitRootByGlobalExitRoot(testData.ctx, globalExitRootNonZero, nil).Return(nil, state.ErrNotFound).Once()
	testData.zkevmMock.EXPECT().ExitRootsByGER(testData.ctx, globalExitRootNonZero).Return(&types.ExitRoots{BlockNumber: types.ArgUint64(l1Block)}, nil).Once()
	testData.stateMock.EXPECT().GetLastBlock(testData.ctx, nil).Return(&state.Block{BlockNumber: l1Block - 1}, nil).Once()

	err := testData.sut.CheckL1SyncGlobalExitRootEnoughToProcessBatch(testData.ctx, 1, globalExitRootNonZero, nil)
	require.ErrorIs(t, err, syncinterfaces.ErrMissingSyncFromL1)
}

func TestCheckL1SyncStatusEnoughToProcessBatchGerNoOnDBWeAre1BlockBeyond(t *testing.T) {
	testData := NewTestData(t)

	l1Block := uint64(123)
	testData.stateMock.EXPECT().GetExitRootByGlobalExitRoot(testData.ctx, globalExitRootNonZero, nil).Return(nil, state.ErrNotFound).Once()
	testData.zkevmMock.EXPECT().ExitRootsByGER(testData.ctx, globalExitRootNonZero).Return(&types.ExitRoots{BlockNumber: types.ArgUint64(l1Block)}, nil).Once()
	testData.stateMock.EXPECT().GetLastBlock(testData.ctx, nil).Return(&state.Block{BlockNumber: l1Block + 1}, nil).Once()

	err := testData.sut.CheckL1SyncGlobalExitRootEnoughToProcessBatch(testData.ctx, 1, globalExitRootNonZero, nil)
	require.ErrorIs(t, err, syncinterfaces.ErrFatalDesyncFromL1)
	l1BlockNumber := err.(*l2_shared.DeSyncPermissionlessAndTrustedNodeError).L1BlockNumber
	require.Equal(t, l1Block, l1BlockNumber, "returns the block where is the discrepancy")
}

func TestCheckL1SyncStatusEnoughToProcessBatchGerNoOnDBWeAreLastBlockSynced(t *testing.T) {
	testData := NewTestData(t)

	l1Block := uint64(123)
	testData.stateMock.EXPECT().GetExitRootByGlobalExitRoot(testData.ctx, globalExitRootNonZero, nil).Return(nil, state.ErrNotFound).Once()
	testData.zkevmMock.EXPECT().ExitRootsByGER(testData.ctx, globalExitRootNonZero).Return(&types.ExitRoots{BlockNumber: types.ArgUint64(l1Block)}, nil).Once()
	testData.stateMock.EXPECT().GetLastBlock(testData.ctx, nil).Return(&state.Block{BlockNumber: l1Block}, nil).Once()

	err := testData.sut.CheckL1SyncGlobalExitRootEnoughToProcessBatch(testData.ctx, 1, globalExitRootNonZero, nil)
	require.ErrorIs(t, err, syncinterfaces.ErrFatalDesyncFromL1)
	l1BlockNumber := err.(*l2_shared.DeSyncPermissionlessAndTrustedNodeError).L1BlockNumber
	require.Equal(t, l1Block, l1BlockNumber, "returns the block where is the discrepancy")
}
