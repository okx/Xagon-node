package feijoa_test

import (
	"context"
	"os"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/db"
	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/pgstatestorage"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/actions/feijoa"
	"github.com/0xPolygonHermez/zkevm-node/test/dbutils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/require"
)

// This test is a exploratory test used to develop. It use a sequencedBlob on Sepolia
// It need Database, a prover >7.x and L1 client
// TODO: Remove this test or convert to a test than can be executed
func TestProcessASequenceBlobUsingCallDataFromSepolia(t *testing.T) {
	l1url := os.Getenv("ZKEVM_NODE_ETHERMAN_URL")
	consensusl1url := os.Getenv("ZKEVM_NODE_ETHERMAN_CONSENSUSL1URL")
	if l1url == "" || consensusl1url == "" {
		// You can set un vscode editing setings.json
		//  "go.testEnvVars": {
		//	"ZKEVM_NODE_ETHERMAN_URL": "url1",
		//	"ZKEVM_NODE_ETHERMAN_CONSENSUSL1URL": "url2",
		//}
		t.Skip("ZKEVM_NODE_ETHERMAN_URL or ZKEVM_NODE_ETHERMAN_CONSENSUSL1URL not set")
	}
	cfg := etherman.Config{
		URL:            l1url,
		ConsensusL1URL: consensusl1url,
	}
	l1Config := etherman.L1Config{
		L1ChainID: 11155111,
		//ZkEVMAddr:                 common.HexToAddress("0x31A6ae85297DD0EeBD66D7556941c33Bd41d565C"),
		//ZkEVMAddr:                 common.HexToAddress("0xD23C761025306cF5038D74FEEb077Cf66DE134DA"),
		ZkEVMAddr:                 common.HexToAddress("0x5e5880098741d1fbd38eaaac51c4215f80f92d27"),
		RollupManagerAddr:         common.HexToAddress("0x9fB0B4A5d4d60aaCfa8DC20B8DF5528Ab26848d3"),
		GlobalExitRootManagerAddr: common.HexToAddress("0x76216E45Bdd20022eEcC07999e50228d7829534B"),
	}
	eth, err := etherman.NewClient(cfg, l1Config)
	require.NoError(t, err)
	ctx := context.Background()
	//toBlock := uint64(5611933)
	//toBlock := uint64(5704000)
	toBlock := uint64(5760696)
	blocks, orders, err := eth.GetRollupInfoByBlockRange(ctx, toBlock, &toBlock)
	require.NoError(t, err)
	require.Equal(t, 1, len(blocks))
	require.Equal(t, 1, len(orders))

	realState := createRealState(t)
	err = addBlock(ctx, &blocks[0], realState, nil)
	if err != nil {
		log.Error(err)
	}
	sut := feijoa.NewProcessorSequenceBlobs(realState, realState, nil)
	err = sut.Process(ctx, orders[blocks[0].BlockHash][0], &blocks[0], nil)
	require.NoError(t, err)
}

const UniqueViolationErr = "23505"

func addBlock(ctx context.Context, block *etherman.Block, storage *state.State, dbTx pgx.Tx) error {
	b := state.Block{
		BlockNumber: block.BlockNumber,
		BlockHash:   block.BlockHash,
		ParentHash:  block.ParentHash,
		ReceivedAt:  block.ReceivedAt,
	}
	// Add block information
	err := storage.AddBlock(ctx, &b, dbTx)

	if pgerr, ok := err.(*pgconn.PgError); ok && pgerr.Code == UniqueViolationErr {
		return nil
	}
	return err
}

func createRealState(t *testing.T) *state.State {
	stateDBCfg := dbutils.NewStateConfigFromEnv()
	stateCfg := state.Config{}
	err := db.RunMigrationsUp(stateDBCfg, db.StateMigrationName)
	require.NoError(t, err)
	stateSqlDB, err := db.NewSQLDB(stateDBCfg)
	stateDb := pgstatestorage.NewPostgresStorage(stateCfg, stateSqlDB)
	executorConfig := executor.Config{
		URI:                "localhost:50071",
		MaxGRPCMessageSize: 1024 * 1024 * 1024,
	}
	executorClient, _, _ := executor.NewExecutorClient(context.TODO(), executorConfig)
	require.NoError(t, err)

	return state.NewState(stateCfg, stateDb, executorClient, nil, nil, nil, nil)
}
