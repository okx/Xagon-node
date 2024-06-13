package etrog

import (
	"context"
	"math"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/db"
	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/l1infotree"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/pgstatestorage"
	"github.com/0xPolygonHermez/zkevm-node/test/dbutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

var (
	stateDBCfg = dbutils.NewStateConfigFromEnv()
	stateCfg   = state.Config{
		MaxCumulativeGasUsed: 800000,
		ChainID:              1000,
		MaxLogsCount:         10000,
		MaxLogsBlockRange:    10000,
		ForkIDIntervals: []state.ForkIDInterval{{
			FromBatchNumber: 0,
			ToBatchNumber:   math.MaxUint64,
			ForkId:          5,
			Version:         "",
		}},
	}
)

func TestProcessorL1InfoTreeUpdate_Process(t *testing.T) {
	ctx := context.Background()
	if err := dbutils.InitOrResetState(stateDBCfg); err != nil {
		panic(err)
	}
	stateDb, err := db.NewSQLDB(stateDBCfg)
	require.NoError(t, err)
	defer stateDb.Close()

	mt, err := l1infotree.NewL1InfoTree(32, [][32]byte{})
	if err != nil {
		panic(err)
	}
	testState := state.NewState(stateCfg, pgstatestorage.NewPostgresStorage(stateCfg, stateDb), nil, nil, nil, mt, nil)

	sut := NewProcessorL1InfoTreeUpdate(testState)
	l1infotree := etherman.GlobalExitRoot{
		BlockNumber:       123,
		MainnetExitRoot:   common.HexToHash("abc"),
		RollupExitRoot:    common.HexToHash("abc"),
		GlobalExitRoot:    common.HexToHash("abc"),
		PreviousBlockHash: common.HexToHash("abc"),
		Timestamp:         time.Now(),
	}
	l1Block := &etherman.Block{
		BlockNumber: 123,
		L1InfoTree:  []etherman.GlobalExitRoot{l1infotree},
	}

	stateBlock := state.Block{
		BlockNumber: l1Block.BlockNumber,
		BlockHash:   l1Block.BlockHash,
		ParentHash:  l1Block.ParentHash,
		ReceivedAt:  l1Block.ReceivedAt,
	}
	dbTx, err := testState.BeginStateTransaction(ctx)
	require.NoError(t, err)
	// Add block information
	err = testState.AddBlock(ctx, &stateBlock, dbTx)
	require.NoError(t, err)

	// Test invalid call, no sequenced batches
	err = sut.Process(ctx, etherman.Order{Name: sut.SupportedEvents()[0], Pos: 0}, l1Block, dbTx)
	require.NoError(t, err)

	err = dbTx.Rollback(ctx)
	require.NoError(t, err)
}
