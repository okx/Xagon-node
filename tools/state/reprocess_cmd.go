package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/0xPolygonHermez/zkevm-node/config"
	"github.com/0xPolygonHermez/zkevm-node/db"
	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/event"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/merkletree"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/pgstatestorage"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/urfave/cli/v2"
)

func isNetworkConfigNeeded(cliCtx *cli.Context) bool {
	// The only reason if to obtain the chainId from the SMC
	return cliCtx.Uint64(flagChainID) == 0
}

func reprocessCmd(cliCtx *cli.Context) error {
	cfg, err := config.Load(cliCtx, isNetworkConfigNeeded(cliCtx))
	if err != nil {
		return err
	}
	log.Init(cfg.Log)
	// This connect to database
	stateSqlDB, err := db.NewSQLDB(cfg.State.DB)
	if err != nil {
		log.Fatal(err)
	}
	l2ChainID := getL2ChainID(cliCtx, cfg)
	needsExecutor := true
	needsStateTree := true

	st := newState(cliCtx.Context, cfg, l2ChainID, []state.ForkIDInterval{}, stateSqlDB, nil, needsExecutor, needsStateTree)

	forksIdIntervals, err := getforkIDIntervalsFromDB(context.Background(), st)
	log.Debugf("forkids: %v  err:%s", forksIdIntervals, err)
	if err != nil {
		log.Errorf("error getting forkIDs from db. Error: %v", err)
		return err
	}
	st.UpdateForkIDIntervalsInMemory(forksIdIntervals)

	action := reprocessAction{
		firstBatchNumber:         getFirstBatchNumber(cliCtx),
		lastBatchNumber:          getLastBatchNumber(cliCtx, cliCtx.Context, st),
		l2ChainId:                l2ChainID,
		updateHasbDB:             getUpdateHashDB(cliCtx),
		st:                       st,
		ctx:                      cliCtx.Context,
		output:                   &reprocessingOutputPretty{},
		flushIdCtrl:              NewFlushIDController(st, cliCtx.Context),
		stopOnError:              !cliCtx.Bool(dontStopOnErrorFlag.Name),
		preferExecutionStateRoot: cliCtx.Bool(preferExecutionStateRootFlag.Name),
	}
	action.output.start(action.firstBatchNumber, action.lastBatchNumber, l2ChainID)
	log.Infof("Reprocessing batches from %d to %d", action.firstBatchNumber, action.lastBatchNumber)
	err = action.start()
	action.output.end(err)

	if err != nil {
		log.Errorf("error reprocessing batches. Error: %v", err)
		return err
	}
	return nil
}

func getUpdateHashDB(cliCtx *cli.Context) bool {
	return cliCtx.Bool(writeOnHashDBFlag.Name)
}

func newEtherman(c config.Config) (*etherman.Client, error) {
	etherman, err := etherman.NewClient(c.Etherman, c.NetworkConfig.L1Config)
	if err != nil {
		return nil, err
	}
	return etherman, nil
}

func getFirstBatchNumber(cliCtx *cli.Context) uint64 {
	res := cliCtx.Uint64(firstBatchNumberFlag.Name)
	if res == 0 {
		return 1
	}
	return res
}
func getforkIDIntervalsFromDB(ctx context.Context, st *state.State) ([]state.ForkIDInterval, error) {
	log.Debug("getting forkIDs from db")
	forkIDIntervals, err := st.GetForkIDs(ctx, nil)
	if err != nil && !errors.Is(err, state.ErrStateNotSynchronized) {
		return []state.ForkIDInterval{}, fmt.Errorf("error getting forkIDs from db. Error: %v", err)
	}
	return forkIDIntervals, nil
}

func getLastBatchNumber(cliCtx *cli.Context, ctx context.Context, st *state.State) uint64 {
	res := cliCtx.Uint64(lastBatchNumberFlag.Name)
	if res == 0 {
		dbTx, err := st.BeginStateTransaction(ctx)
		if err != nil {
			log.Fatalf("error creating db transaction to get latest block. Error: %v", err)
		}
		lastBatch, err := st.GetLastBatchNumber(ctx, dbTx)
		if err != nil {
			log.Fatalf("no last batch. Error: %v", err)
		}
		return lastBatch
	}
	return res
}

func getL2ChainID(cliCtx *cli.Context, c *config.Config) uint64 {
	flagL2chainID := cliCtx.Uint64(flagChainID)
	if flagL2chainID != 0 {
		log.Debugf("Using L2ChainID from flag: %d", flagL2chainID)
		return flagL2chainID
	}

	etherman, err := newEtherman(*c)
	if err != nil {
		log.Fatal(err)
	}

	// READ CHAIN ID FROM POE SC
	log.Debug("Reading L2ChainID from SMC")
	l2ChainID, err := etherman.GetL2ChainID()
	if err != nil {
		log.Fatal(err)
	}
	log.Debugf("Using L2ChainID from SMC: %d", l2ChainID)
	return l2ChainID
}

func newState(ctx context.Context, c *config.Config, l2ChainID uint64, forkIDIntervals []state.ForkIDInterval, sqlDB *pgxpool.Pool, eventLog *event.EventLog, needsExecutor, needsStateTree bool) *state.State {
	stateCfg := state.Config{
		MaxCumulativeGasUsed:         c.State.Batch.Constraints.MaxCumulativeGasUsed,
		ChainID:                      l2ChainID,
		ForkIDIntervals:              forkIDIntervals,
		MaxResourceExhaustedAttempts: c.Executor.MaxResourceExhaustedAttempts,
		WaitOnResourceExhaustion:     c.Executor.WaitOnResourceExhaustion,
		ForkUpgradeBatchNumber:       c.ForkUpgradeBatchNumber,
		ForkUpgradeNewForkId:         c.ForkUpgradeNewForkId,
	}

	stateDb := pgstatestorage.NewPostgresStorage(stateCfg, sqlDB)

	// Executor
	var executorClient executor.ExecutorServiceClient
	if needsExecutor {
		executorClient, _, _ = executor.NewExecutorClient(ctx, c.Executor)
	}

	// State Tree
	var stateTree *merkletree.StateTree
	if needsStateTree {
		stateDBClient, _, _ := merkletree.NewMTDBServiceClient(ctx, c.MTClient)
		stateTree = merkletree.NewStateTree(stateDBClient)
	}

	st := state.NewState(stateCfg, stateDb, executorClient, stateTree, eventLog, nil, nil)
	return st
}

func getBatchByNumber(ctx context.Context, st *state.State, batchNumber uint64) (*state.Batch, error) {
	dbTx, err := st.BeginStateTransaction(ctx)
	if err != nil {
		log.Errorf("error creating db transaction to get latest block. Error: %v", err)
		return nil, err
	}
	batch, err := st.GetBatchByNumber(ctx, batchNumber, dbTx)
	if err != nil {
		log.Errorf("no batch %d. Error: %v", 0, err)
		return nil, err
	}
	_ = dbTx.Commit(ctx)
	return batch, nil
}
