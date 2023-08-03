package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/0xPolygonHermez/zkevm-node"
	"github.com/0xPolygonHermez/zkevm-node/aggregator"
	"github.com/0xPolygonHermez/zkevm-node/config"
	"github.com/0xPolygonHermez/zkevm-node/db"
	"github.com/0xPolygonHermez/zkevm-node/etherman"
	"github.com/0xPolygonHermez/zkevm-node/ethtxmanager"
	"github.com/0xPolygonHermez/zkevm-node/event"
	"github.com/0xPolygonHermez/zkevm-node/event/nileventstorage"
	"github.com/0xPolygonHermez/zkevm-node/event/pgeventstorage"
	"github.com/0xPolygonHermez/zkevm-node/gasprice"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc"
	"github.com/0xPolygonHermez/zkevm-node/jsonrpc/client"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/merkletree"
	"github.com/0xPolygonHermez/zkevm-node/metrics"
	"github.com/0xPolygonHermez/zkevm-node/pool"
	"github.com/0xPolygonHermez/zkevm-node/pool/pgpoolstorage"
	"github.com/0xPolygonHermez/zkevm-node/sequencer"
	"github.com/0xPolygonHermez/zkevm-node/sequencesender"
	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/0xPolygonHermez/zkevm-node/state/runtime/executor"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/cli/v2"
)

func start(cliCtx *cli.Context) error {
	c, err := config.Load(cliCtx, true)
	if err != nil {
		return err
	}
	setupLog(c.Log)

	if c.Log.Environment == log.EnvironmentDevelopment {
		zkevm.PrintVersion(os.Stdout)
		log.Info("Starting application")
	} else if c.Log.Environment == log.EnvironmentProduction {
		logVersion()
	}

	if c.Metrics.Enabled {
		metrics.Init()
	}
	components := cliCtx.StringSlice(config.FlagComponents)

	// Only runs migration if the component is the synchronizer and if the flag is deactivated
	if !cliCtx.Bool(config.FlagMigrations) {
		for _, comp := range components {
			if comp == SYNCHRONIZER {
				runStateMigrations(c.StateDB)
			}
		}
	}
	checkStateMigrations(c.StateDB)

	// Decide if this node instance needs an executor and/or a state tree
	var needsExecutor, needsStateTree bool

	for _, component := range components {
		switch component {
		case SEQUENCER, RPC, SYNCHRONIZER:
			needsExecutor = true
			needsStateTree = true
		}
	}

	// Event log
	var eventLog *event.EventLog
	var eventStorage event.Storage

	if c.EventLog.DB.Name != "" {
		eventStorage, err = pgeventstorage.NewPostgresEventStorage(c.EventLog.DB)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		eventStorage, err = nileventstorage.NewNilEventStorage()
		if err != nil {
			log.Fatal(err)
		}
	}
	eventLog = event.NewEventLog(c.EventLog, eventStorage)

	// Core State DB
	stateSqlDB, err := db.NewSQLDB(c.StateDB)
	if err != nil {
		log.Fatal(err)
	}

	var (
		cancelFuncs []context.CancelFunc
		etherman    *etherman.Client
	)

	etherman, err = newEtherman(*c)
	if err != nil {
		log.Fatal(err)
	}

	// READ CHAIN ID FROM POE SC
	l2ChainID, err := etherman.GetL2ChainID()
	if err != nil {
		log.Fatal(err)
	}
	// Read Fork ID FROM POE SC
	forkIDIntervals, err := etherman.GetForks(cliCtx.Context, c.NetworkConfig.Genesis.GenesisBlockNum)
	if err != nil {
		log.Fatal("error getting forks. Please check the configuration. Error: ", err)
	} else if len(forkIDIntervals) == 0 {
		log.Fatal("error: no forkID received. It should receive at least one, please check the configuration...")
	}

	currentForkID := forkIDIntervals[len(forkIDIntervals)-1].ForkId
	log.Infof("Fork ID read from POE SC = %v", forkIDIntervals[len(forkIDIntervals)-1].ForkId)
	c.Aggregator.ChainID = l2ChainID
	c.Aggregator.ForkId = currentForkID
	log.Infof("Chain ID read from POE SC = %v", l2ChainID)

	ctx := context.Background()
	st := newState(ctx, c, l2ChainID, forkIDIntervals, stateSqlDB, eventLog, needsExecutor, needsStateTree)

	ethTxManagerStorage, err := ethtxmanager.NewPostgresStorage(c.StateDB)
	if err != nil {
		log.Fatal(err)
	}

	etm := ethtxmanager.New(c.EthTxManager, etherman, ethTxManagerStorage, st)

	ev := &event.Event{
		ReceivedAt: time.Now(),
		Source:     event.Source_Node,
		Level:      event.Level_Info,
		EventID:    event.EventID_NodeComponentStarted,
	}

	var poolInstance *pool.Pool

	if c.Metrics.ProfilingEnabled {
		go startProfilingHttpServer(c.Metrics)
	}
	for _, component := range components {
		switch component {
		case AGGREGATOR:
			ev.Component = event.Component_Aggregator
			ev.Description = "Running aggregator"
			err := eventLog.LogEvent(ctx, ev)
			if err != nil {
				log.Fatal(err)
			}
			go runAggregator(ctx, c.Aggregator, etherman, etm, st)
		case SEQUENCER:
			ev.Component = event.Component_Sequencer
			ev.Description = "Running sequencer"
			err := eventLog.LogEvent(ctx, ev)
			if err != nil {
				log.Fatal(err)
			}
			if poolInstance == nil {
				poolInstance = createPool(c.Pool, l2ChainID, st, eventLog)
			}
			seq := createSequencer(*c, poolInstance, ethTxManagerStorage, st, eventLog)
			go seq.Start(ctx)
		case SEQUENCE_SENDER:
			ev.Component = event.Component_Sequence_Sender
			ev.Description = "Running sequence sender"
			err := eventLog.LogEvent(ctx, ev)
			if err != nil {
				log.Fatal(err)
			}
			if poolInstance == nil {
				poolInstance = createPool(c.Pool, l2ChainID, st, eventLog)
			}
			seqSender := createSequenceSender(*c, poolInstance, ethTxManagerStorage, st, eventLog)
			go seqSender.Start(ctx)
		case RPC:

			//LEVITATION_BEGIN
			// Start CandidateBlockPullAgent
			log.Errorf("LEVITATION:Will now start new CandidateBlockPullAgent")
			if err := jsonrpc.NewCandidateBlockPullAgent(c.RPC, l2ChainID).Start(); err != nil {
				log.Fatal(err)
			}
			//LEVITATION_END

			ev.Component = event.Component_RPC
			ev.Description = "Running JSON-RPC server"
			err := eventLog.LogEvent(ctx, ev)
			if err != nil {
				log.Fatal(err)
			}
			if poolInstance == nil {
				poolInstance = createPool(c.Pool, l2ChainID, st, eventLog)
			}
			if c.RPC.EnableL2SuggestedGasPricePolling {
				// Needed for rejecting transactions with too low gas price
				poolInstance.StartPollingMinSuggestedGasPrice(ctx)
			}
			apis := map[string]bool{}
			for _, a := range cliCtx.StringSlice(config.FlagHTTPAPI) {
				apis[a] = true
			}

			go runJSONRPCServer(*c, etherman, l2ChainID, poolInstance, st, apis)
		case SYNCHRONIZER:
			ev.Component = event.Component_Synchronizer
			ev.Description = "Running synchronizer"
			err := eventLog.LogEvent(ctx, ev)
			if err != nil {
				log.Fatal(err)
			}
			if poolInstance == nil {
				poolInstance = createPool(c.Pool, l2ChainID, st, eventLog)
			}
			go runSynchronizer(*c, etherman, etm, st, poolInstance, eventLog)
		case ETHTXMANAGER:
			ev.Component = event.Component_EthTxManager
			ev.Description = "Running eth tx manager service"
			err := eventLog.LogEvent(ctx, ev)
			if err != nil {
				log.Fatal(err)
			}
			etm := createEthTxManager(*c, ethTxManagerStorage, st)
			go etm.Start()
		case L2GASPRICER:
			ev.Component = event.Component_GasPricer
			ev.Description = "Running L2 gasPricer"
			err := eventLog.LogEvent(ctx, ev)
			if err != nil {
				log.Fatal(err)
			}
			if poolInstance == nil {
				poolInstance = createPool(c.Pool, l2ChainID, st, eventLog)
			}
			go runL2GasPriceSuggester(c.L2GasPriceSuggester, st, poolInstance, etherman)
		}
	}

	if c.Metrics.Enabled {
		go startMetricsHttpServer(c.Metrics)
	}

	waitSignal(cancelFuncs)

	return nil
}

func setupLog(c log.Config) {
	log.Init(c)
}

func runStateMigrations(c db.Config) {
	runMigrations(c, db.StateMigrationName)
}

func checkStateMigrations(c db.Config) {
	err := db.CheckMigrations(c, db.StateMigrationName)
	if err != nil {
		log.Fatal(err)
	}
}

func runPoolMigrations(c db.Config) {
	runMigrations(c, db.PoolMigrationName)
}

func runMigrations(c db.Config, name string) {
	log.Infof("running migrations for %v", name)
	err := db.RunMigrationsUp(c, name)
	if err != nil {
		log.Fatal(err)
	}
}

func newEtherman(c config.Config) (*etherman.Client, error) {
	etherman, err := etherman.NewClient(c.Etherman, c.NetworkConfig.L1Config)
	if err != nil {
		return nil, err
	}
	return etherman, nil
}

func runSynchronizer(cfg config.Config, etherman *etherman.Client, ethTxManager *ethtxmanager.Client, st *state.State, pool *pool.Pool, eventLog *event.EventLog) {
	var trustedSequencerURL string
	var err error
	if !cfg.IsTrustedSequencer {
		if cfg.Synchronizer.TrustedSequencerURL != "" {
			trustedSequencerURL = cfg.Synchronizer.TrustedSequencerURL
		} else {
			log.Debug("getting trusted sequencer URL from smc")
			trustedSequencerURL, err = etherman.GetTrustedSequencerURL()
			if err != nil {
				log.Fatal("error getting trusted sequencer URI. Error: %v", err)
			}
		}
		log.Debug("trustedSequencerURL ", trustedSequencerURL)
	}
	zkEVMClient := client.NewClient(trustedSequencerURL)

	sy, err := synchronizer.NewSynchronizer(
		cfg.IsTrustedSequencer, etherman, st, pool, ethTxManager,
		zkEVMClient, eventLog, cfg.NetworkConfig.Genesis, cfg.Synchronizer,
	)
	if err != nil {
		log.Fatal(err)
	}
	if err := sy.Sync(); err != nil {
		log.Fatal(err)
	}
}

func runJSONRPCServer(c config.Config, etherman *etherman.Client, chainID uint64, pool *pool.Pool, st *state.State, apis map[string]bool) {
	var err error

	storage := jsonrpc.NewStorage()
	c.RPC.MaxCumulativeGasUsed = c.Sequencer.MaxCumulativeGasUsed
	if !c.IsTrustedSequencer {
		if c.RPC.SequencerNodeURI == "" {
			log.Debug("getting trusted sequencer URL from smc")
			c.RPC.SequencerNodeURI, err = etherman.GetTrustedSequencerURL()
			if err != nil {
				log.Fatal("error getting trusted sequencer URI. Error: %v", err)
			}
		}
		log.Debug("SequencerNodeURI ", c.RPC.SequencerNodeURI)
	}

	services := []jsonrpc.Service{}
	if _, ok := apis[jsonrpc.APIEth]; ok {
		services = append(services, jsonrpc.Service{
			Name:    jsonrpc.APIEth,
			Service: jsonrpc.NewEthEndpoints(c.RPC, chainID, pool, st, storage),
		})
	}

	if _, ok := apis[jsonrpc.APINet]; ok {
		services = append(services, jsonrpc.Service{
			Name:    jsonrpc.APINet,
			Service: jsonrpc.NewNetEndpoints(c.RPC, chainID),
		})
	}

	if _, ok := apis[jsonrpc.APIZKEVM]; ok {
		services = append(services, jsonrpc.Service{
			Name:    jsonrpc.APIZKEVM,
			Service: jsonrpc.NewZKEVMEndpoints(c.RPC, st),
		})
	}

	if _, ok := apis[jsonrpc.APITxPool]; ok {
		services = append(services, jsonrpc.Service{
			Name:    jsonrpc.APITxPool,
			Service: &jsonrpc.TxPoolEndpoints{},
		})
	}

	if _, ok := apis[jsonrpc.APIDebug]; ok {
		services = append(services, jsonrpc.Service{
			Name:    jsonrpc.APIDebug,
			Service: jsonrpc.NewDebugEndpoints(c.RPC, st),
		})
	}

	if _, ok := apis[jsonrpc.APIWeb3]; ok {
		services = append(services, jsonrpc.Service{
			Name:    jsonrpc.APIWeb3,
			Service: &jsonrpc.Web3Endpoints{},
		})
	}

	if err := jsonrpc.NewServer(c.RPC, chainID, pool, st, storage, services).Start(); err != nil {
		log.Fatal(err)
	}
}

func createSequencer(cfg config.Config, pool *pool.Pool, etmStorage *ethtxmanager.PostgresStorage, st *state.State, eventLog *event.EventLog) *sequencer.Sequencer {
	etherman, err := newEtherman(cfg)
	if err != nil {
		log.Fatal(err)
	}

	ethTxManager := ethtxmanager.New(cfg.EthTxManager, etherman, etmStorage, st)

	seq, err := sequencer.New(cfg.Sequencer, pool, st, etherman, ethTxManager, eventLog)
	if err != nil {
		log.Fatal(err)
	}
	return seq
}

func createSequenceSender(cfg config.Config, pool *pool.Pool, etmStorage *ethtxmanager.PostgresStorage, st *state.State, eventLog *event.EventLog) *sequencesender.SequenceSender {
	etherman, err := newEtherman(cfg)
	if err != nil {
		log.Fatal(err)
	}

	for _, privateKey := range cfg.SequenceSender.PrivateKeys {
		_, err := etherman.LoadAuthFromKeyStore(privateKey.Path, privateKey.Password)
		if err != nil {
			log.Fatal(err)
		}
	}

	cfg.SequenceSender.ForkUpgradeBatchNumber = cfg.ForkUpgradeBatchNumber

	ethTxManager := ethtxmanager.New(cfg.EthTxManager, etherman, etmStorage, st)

	seqSender, err := sequencesender.New(cfg.SequenceSender, st, etherman, ethTxManager, eventLog)
	if err != nil {
		log.Fatal(err)
	}

	return seqSender
}

func runAggregator(ctx context.Context, c aggregator.Config, etherman *etherman.Client, ethTxManager *ethtxmanager.Client, st *state.State) {
	agg, err := aggregator.New(c, st, ethTxManager, etherman)
	if err != nil {
		log.Fatal(err)
	}
	err = agg.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}
}

// runL2GasPriceSuggester init gas price gasPriceEstimator based on type in config.
func runL2GasPriceSuggester(cfg gasprice.Config, state *state.State, pool *pool.Pool, etherman *etherman.Client) {
	ctx := context.Background()
	gasprice.NewL2GasPriceSuggester(ctx, cfg, pool, etherman, state)
}

func waitSignal(cancelFuncs []context.CancelFunc) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	for sig := range signals {
		switch sig {
		case os.Interrupt, os.Kill:
			log.Info("terminating application gracefully...")

			exitStatus := 0
			for _, cancel := range cancelFuncs {
				cancel()
			}
			os.Exit(exitStatus)
		}
	}
}

func newState(ctx context.Context, c *config.Config, l2ChainID uint64, forkIDIntervals []state.ForkIDInterval, sqlDB *pgxpool.Pool, eventLog *event.EventLog, needsExecutor, needsStateTree bool) *state.State {
	stateDb := state.NewPostgresStorage(sqlDB)

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

	stateCfg := state.Config{
		MaxCumulativeGasUsed:         c.Sequencer.MaxCumulativeGasUsed,
		ChainID:                      l2ChainID,
		ForkIDIntervals:              forkIDIntervals,
		MaxResourceExhaustedAttempts: c.Executor.MaxResourceExhaustedAttempts,
		WaitOnResourceExhaustion:     c.Executor.WaitOnResourceExhaustion,
		ForkUpgradeBatchNumber:       c.ForkUpgradeBatchNumber,
		ForkUpgradeNewForkId:         c.ForkUpgradeNewForkId,
	}

	st := state.NewState(stateCfg, stateDb, executorClient, stateTree, eventLog)
	return st
}

func createPool(cfgPool pool.Config, l2ChainID uint64, st *state.State, eventLog *event.EventLog) *pool.Pool {
	runPoolMigrations(cfgPool.DB)
	poolStorage, err := pgpoolstorage.NewPostgresPoolStorage(cfgPool.DB)
	if err != nil {
		log.Fatal(err)
	}
	poolInstance := pool.NewPool(cfgPool, poolStorage, st, l2ChainID, eventLog)
	return poolInstance
}

func createEthTxManager(cfg config.Config, etmStorage *ethtxmanager.PostgresStorage, st *state.State) *ethtxmanager.Client {
	etherman, err := newEtherman(cfg)
	if err != nil {
		log.Fatal(err)
	}

	for _, privateKey := range cfg.EthTxManager.PrivateKeys {
		_, err := etherman.LoadAuthFromKeyStore(privateKey.Path, privateKey.Password)
		if err != nil {
			log.Fatal(err)
		}
	}
	etm := ethtxmanager.New(cfg.EthTxManager, etherman, etmStorage, st)
	return etm
}

func startProfilingHttpServer(c metrics.Config) {
	const two = 2
	mux := http.NewServeMux()
	address := fmt.Sprintf("%s:%d", c.ProfilingHost, c.ProfilingPort)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Errorf("failed to create tcp listener for profiling: %v", err)
		return
	}
	mux.HandleFunc(metrics.ProfilingIndexEndpoint, pprof.Index)
	mux.HandleFunc(metrics.ProfileEndpoint, pprof.Profile)
	mux.HandleFunc(metrics.ProfilingCmdEndpoint, pprof.Cmdline)
	mux.HandleFunc(metrics.ProfilingSymbolEndpoint, pprof.Symbol)
	mux.HandleFunc(metrics.ProfilingTraceEndpoint, pprof.Trace)
	profilingServer := &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: two * time.Minute,
		ReadTimeout:       two * time.Minute,
	}
	log.Infof("profiling server listening on port %d", c.ProfilingPort)
	if err := profilingServer.Serve(lis); err != nil {
		if err == http.ErrServerClosed {
			log.Warnf("http server for profiling stopped")
			return
		}
		log.Errorf("closed http connection for profiling server: %v", err)
		return
	}
}

func startMetricsHttpServer(c metrics.Config) {
	const ten = 10
	mux := http.NewServeMux()
	address := fmt.Sprintf("%s:%d", c.Host, c.Port)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Errorf("failed to create tcp listener for metrics: %v", err)
		return
	}
	mux.Handle(metrics.Endpoint, promhttp.Handler())

	metricsServer := &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: ten * time.Second,
		ReadTimeout:       ten * time.Second,
	}
	log.Infof("metrics server listening on port %d", c.Port)
	if err := metricsServer.Serve(lis); err != nil {
		if err == http.ErrServerClosed {
			log.Warnf("http server for metrics stopped")
			return
		}
		log.Errorf("closed http connection for metrics server: %v", err)
		return
	}
}

func logVersion() {
	log.Infow("Starting application",
		// node version is already logged by default
		"gitRevision", zkevm.GitRev,
		"gitBranch", zkevm.GitBranch,
		"goVersion", runtime.Version(),
		"built", zkevm.BuildDate,
		"os/arch", fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	)
}
