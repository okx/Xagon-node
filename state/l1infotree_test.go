package state_test

import (
	"context"
	"math"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/db"
	"github.com/0xPolygonHermez/zkevm-node/l1infotree"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/mocks"
	"github.com/0xPolygonHermez/zkevm-node/state/pgstatestorage"
	"github.com/0xPolygonHermez/zkevm-node/test/dbutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestFirstLeafOfL1InfoTreeIsIndex0(t *testing.T) {
	stateDBCfg := dbutils.NewStateConfigFromEnv()
	if err := dbutils.InitOrResetState(stateDBCfg); err != nil {
		panic(err)
	}

	stateDb, err := db.NewSQLDB(stateDBCfg)
	if err != nil {
		panic(err)
	}
	forkID := uint64(state.FORKID_ETROG)
	stateCfg := state.Config{
		MaxCumulativeGasUsed: 800000,
		ChainID:              1000,
		MaxLogsCount:         10000,
		MaxLogsBlockRange:    10000,
		ForkIDIntervals: []state.ForkIDInterval{{
			FromBatchNumber: 0,
			ToBatchNumber:   math.MaxUint64,
			ForkId:          forkID,
			Version:         "",
		}},
	}
	ctx := context.Background()
	storage := pgstatestorage.NewPostgresStorage(stateCfg, stateDb)
	mt, err := l1infotree.NewL1InfoTree(32, [][32]byte{})
	if err != nil {
		panic(err)
	}
	testState := state.NewState(stateCfg, storage, nil, nil, nil, mt)
	dbTx, err := testState.BeginStateTransaction(ctx)
	defer func() {
		_ = dbTx.Rollback(ctx)
	}()
	require.NoError(t, err)
	block := state.Block{BlockNumber: 123}
	err = testState.AddBlock(ctx, &block, dbTx)
	require.NoError(t, err)

	leaf := state.L1InfoTreeLeaf{
		GlobalExitRoot: state.GlobalExitRoot{
			GlobalExitRoot: common.Hash{},
			BlockNumber:    123,
		},
		PreviousBlockHash: common.Hash{},
	}
	insertedLeaf, err := testState.AddL1InfoTreeLeaf(ctx, &leaf, dbTx)
	require.NoError(t, err)
	require.Equal(t, insertedLeaf.L1InfoTreeIndex, uint32(0))
}

func TestGetCurrentL1InfoRootBuildCacheIfNil(t *testing.T) {
	mockStorage := mocks.NewStorageMock(t)
	stateCfg := state.Config{
		MaxCumulativeGasUsed: 800000,
		ChainID:              1000,
		MaxLogsCount:         10000,
		MaxLogsBlockRange:    10000,
		ForkIDIntervals: []state.ForkIDInterval{{
			FromBatchNumber: 0,
			ToBatchNumber:   math.MaxUint64,
			ForkId:          uint64(state.FORKID_ETROG),
			Version:         "",
		}},
	}
	ctx := context.Background()
	testState := state.NewState(stateCfg, mockStorage, nil, nil, nil, nil)

	mockStorage.EXPECT().GetAllL1InfoRootEntries(ctx, nil).Return([]state.L1InfoTreeExitRootStorageEntry{}, nil)

	l1InfoRoot, err := testState.GetCurrentL1InfoRoot(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, l1InfoRoot, common.HexToHash("0x27ae5ba08d7291c96c8cbddcc148bf48a6d68c7974b94356f53754ef6171d757"))
}

func TestGetCurrentL1InfoRootNoBuildCacheIfNotNil(t *testing.T) {
	mockStorage := mocks.NewStorageMock(t)
	stateCfg := state.Config{
		MaxCumulativeGasUsed: 800000,
		ChainID:              1000,
		MaxLogsCount:         10000,
		MaxLogsBlockRange:    10000,
		ForkIDIntervals: []state.ForkIDInterval{{
			FromBatchNumber: 0,
			ToBatchNumber:   math.MaxUint64,
			ForkId:          uint64(state.FORKID_ETROG),
			Version:         "",
		}},
	}
	ctx := context.Background()
	l1InfoTree, err := l1infotree.NewL1InfoTree(uint8(32), nil)
	require.NoError(t, err)
	testState := state.NewState(stateCfg, mockStorage, nil, nil, nil, l1InfoTree)

	// GetCurrentL1InfoRoot use the cache value in state.l1InfoTree
	l1InfoRoot, err := testState.GetCurrentL1InfoRoot(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, l1InfoRoot, common.HexToHash("0x27ae5ba08d7291c96c8cbddcc148bf48a6d68c7974b94356f53754ef6171d757"))
}

func TestAddL1InfoTreeLeafIfNil(t *testing.T) {
	mockStorage := mocks.NewStorageMock(t)
	stateCfg := state.Config{
		MaxCumulativeGasUsed: 800000,
		ChainID:              1000,
		MaxLogsCount:         10000,
		MaxLogsBlockRange:    10000,
		ForkIDIntervals: []state.ForkIDInterval{{
			FromBatchNumber: 0,
			ToBatchNumber:   math.MaxUint64,
			ForkId:          uint64(state.FORKID_ETROG),
			Version:         "",
		}},
	}
	ctx := context.Background()
	testState := state.NewState(stateCfg, mockStorage, nil, nil, nil, nil)

	mockStorage.EXPECT().GetLatestIndex(ctx, mock.Anything).Return(uint32(0), state.ErrNotFound)
	mockStorage.EXPECT().AddL1InfoRootToExitRoot(ctx, mock.Anything, mock.Anything).Return(nil)
	// This call is for rebuild cache
	mockStorage.EXPECT().GetAllL1InfoRootEntries(ctx, nil).Return([]state.L1InfoTreeExitRootStorageEntry{}, nil)
	leaf := state.L1InfoTreeLeaf{
		GlobalExitRoot: state.GlobalExitRoot{
			GlobalExitRoot: common.Hash{},
		},
	}
	addLeaf, err := testState.AddL1InfoTreeLeaf(ctx, &leaf, nil)
	require.NoError(t, err)
	require.Equal(t, addLeaf.L1InfoTreeRoot, common.HexToHash("0xea536769cad1a63ffb1ea52ae772983905c3f0e2f8914e6c0e2af956637e480c"))
	require.Equal(t, addLeaf.L1InfoTreeIndex, uint32(0))
}
