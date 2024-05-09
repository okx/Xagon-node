package config_test

import (
	"flag"
	"math/big"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/aggregator"
	"github.com/0xPolygonHermez/zkevm-node/config"
	"github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"
)

func Test_Defaults(t *testing.T) {
	tcs := []struct {
		path          string
		expectedValue interface{}
	}{
		{
			path:          "Log.Environment",
			expectedValue: log.LogEnvironment("development"),
		},
		{
			path:          "Log.Level",
			expectedValue: "info",
		},
		{
			path:          "Log.Outputs",
			expectedValue: []string{"stderr"},
		},
		{
			path:          "Synchronizer.SyncChunkSize",
			expectedValue: uint64(100),
		},
		{
			path:          "Synchronizer.L1SynchronizationMode",
			expectedValue: "sequential",
		},
		{
			path:          "Synchronizer.L1ParallelSynchronization.MaxClients",
			expectedValue: uint64(10),
		},
		{
			path:          "Synchronizer.L1ParallelSynchronization.MaxPendingNoProcessedBlocks",
			expectedValue: uint64(25),
		},
		{
			path:          "Synchronizer.L2Synchronization.AcceptEmptyClosedBatches",
			expectedValue: false,
		},
		{
			path:          "Synchronizer.L2Synchronization.ReprocessFullBatchOnClose",
			expectedValue: true,
		},
		{
			path:          "Synchronizer.L2Synchronization.CheckLastL2BlockHashOnCloseBatch",
			expectedValue: true,
		},

		{
			path:          "Sequencer.DeletePoolTxsL1BlockConfirmations",
			expectedValue: uint64(100),
		},
		{
			path:          "Sequencer.DeletePoolTxsCheckInterval",
			expectedValue: types.NewDuration(12 * time.Hour),
		},
		{
			path:          "Sequencer.TxLifetimeCheckInterval",
			expectedValue: types.NewDuration(10 * time.Minute),
		},
		{
			path:          "Sequencer.TxLifetimeMax",
			expectedValue: types.NewDuration(3 * time.Hour),
		},
		{
			path:          "Sequencer.LoadPoolTxsCheckInterval",
			expectedValue: types.NewDuration(500 * time.Millisecond),
		},
		{
			path:          "Sequencer.StateConsistencyCheckInterval",
			expectedValue: types.NewDuration(5 * time.Second),
		},
		{
			path:          "Sequencer.Finalizer.ForcedBatchesTimeout",
			expectedValue: types.NewDuration(60 * time.Second),
		},
		{
			path:          "Sequencer.Finalizer.NewTxsWaitInterval",
			expectedValue: types.NewDuration(100 * time.Millisecond),
		},
		{
			path:          "Sequencer.Finalizer.ResourceExhaustedMarginPct",
			expectedValue: uint32(10),
		},
		{
			path:          "Sequencer.Finalizer.ForcedBatchesL1BlockConfirmations",
			expectedValue: uint64(64),
		},
		{
			path:          "Sequencer.Finalizer.L1InfoTreeL1BlockConfirmations",
			expectedValue: uint64(64),
		},
		{
			path:          "Sequencer.Finalizer.ForcedBatchesCheckInterval",
			expectedValue: types.NewDuration(10 * time.Second),
		},
		{
			path:          "Sequencer.Finalizer.L1InfoTreeCheckInterval",
			expectedValue: types.NewDuration(10 * time.Second),
		},
		{
			path:          "Sequencer.Finalizer.L2BlockMaxDeltaTimestamp",
			expectedValue: types.NewDuration(3 * time.Second),
		},
		{
			path:          "Sequencer.Finalizer.HaltOnBatchNumber",
			expectedValue: uint64(0),
		},
		{
			path:          "Sequencer.Finalizer.BatchMaxDeltaTimestamp",
			expectedValue: types.NewDuration(10 * time.Second),
		},
		{
			path:          "Sequencer.Finalizer.Metrics.Interval",
			expectedValue: types.NewDuration(60 * time.Minute),
		},
		{
			path:          "Sequencer.Finalizer.Metrics.EnableLog",
			expectedValue: true,
		},
		{
			path:          "Sequencer.StreamServer.Port",
			expectedValue: uint16(0),
		},
		{
			path:          "Sequencer.StreamServer.Filename",
			expectedValue: "",
		},
		{
			path:          "Sequencer.StreamServer.Version",
			expectedValue: uint8(0),
		},
		{
			path:          "Sequencer.StreamServer.Enabled",
			expectedValue: false,
		},
		{
			path:          "SequenceSender.WaitPeriodSendSequence",
			expectedValue: types.NewDuration(5 * time.Second),
		},
		{
			path:          "SequenceSender.LastBatchVirtualizationTimeMaxWaitPeriod",
			expectedValue: types.NewDuration(5 * time.Second),
		},
		{
			path:          "SequenceSender.L1BlockTimestampMargin",
			expectedValue: types.NewDuration(30 * time.Second),
		},
		{
			path:          "SequenceSender.MaxTxSizeForL1",
			expectedValue: uint64(131072),
		},
		{
			path:          "SequenceSender.GasOffset",
			expectedValue: uint64(80000),
		},
		{
			path:          "Etherman.URL",
			expectedValue: "http://localhost:8545",
		},
		{
			path:          "NetworkConfig.L1Config.L1ChainID",
			expectedValue: uint64(5),
		},
		{
			path:          "NetworkConfig.L1Config.ZkEVMAddr",
			expectedValue: common.HexToAddress("0xa997cfD539E703921fD1e3Cf25b4c241a27a4c7A"),
		},
		{
			path:          "NetworkConfig.L1Config.PolAddr",
			expectedValue: common.HexToAddress("0x1319D23c2F7034F52Eb07399702B040bA278Ca49"),
		},
		{
			path:          "NetworkConfig.L1Config.GlobalExitRootManagerAddr",
			expectedValue: common.HexToAddress("0x4d9427DCA0406358445bC0a8F88C26b704004f74"),
		},
		{
			path:          "Etherman.MultiGasProvider",
			expectedValue: false,
		},
		{
			path:          "EthTxManager.FrequencyToMonitorTxs",
			expectedValue: types.NewDuration(1 * time.Second),
		},
		{
			path:          "EthTxManager.WaitTxToBeMined",
			expectedValue: types.NewDuration(2 * time.Minute),
		},
		{
			path:          "EthTxManager.WaitTxToBeMined",
			expectedValue: types.NewDuration(2 * time.Minute),
		},
		{
			path:          "EthTxManager.ForcedGas",
			expectedValue: uint64(0),
		},
		{
			path:          "EthTxManager.GasPriceMarginFactor",
			expectedValue: float64(1),
		},
		{
			path:          "EthTxManager.MaxGasPriceLimit",
			expectedValue: uint64(0),
		},
		{
			path:          "L2GasPriceSuggester.DefaultGasPriceWei",
			expectedValue: uint64(2000000000),
		},
		{
			path:          "L2GasPriceSuggester.MaxGasPriceWei",
			expectedValue: uint64(0),
		},
		{
			path:          "MTClient.URI",
			expectedValue: "zkevm-prover:50061",
		},
		{
			path:          "State.DB.User",
			expectedValue: "state_user",
		},
		{
			path:          "State.DB.Password",
			expectedValue: "state_password",
		},
		{
			path:          "State.DB.Name",
			expectedValue: "state_db",
		},
		{
			path:          "State.DB.Host",
			expectedValue: "zkevm-state-db",
		},
		{
			path:          "State.DB.Port",
			expectedValue: "5432",
		},
		{
			path:          "State.DB.EnableLog",
			expectedValue: false,
		},
		{
			path:          "State.DB.MaxConns",
			expectedValue: 200,
		},
		{
			path:          "Pool.IntervalToRefreshGasPrices",
			expectedValue: types.NewDuration(5 * time.Second),
		},
		{
			path:          "Pool.MaxTxBytesSize",
			expectedValue: uint64(100132),
		},
		{
			path:          "Pool.MaxTxDataBytesSize",
			expectedValue: 100000,
		},

		{
			path:          "Pool.DefaultMinGasPriceAllowed",
			expectedValue: uint64(1000000000),
		},
		{
			path:          "Pool.MinAllowedGasPriceInterval",
			expectedValue: types.NewDuration(5 * time.Minute),
		},
		{
			path:          "Pool.PollMinAllowedGasPriceInterval",
			expectedValue: types.NewDuration(15 * time.Second),
		},
		{
			path:          "Pool.AccountQueue",
			expectedValue: uint64(64),
		},
		{
			path:          "Pool.GlobalQueue",
			expectedValue: uint64(1024),
		},
		{
			path:          "Pool.EffectiveGasPrice.Enabled",
			expectedValue: false,
		},
		{
			path:          "Pool.EffectiveGasPrice.L1GasPriceFactor",
			expectedValue: float64(0.25),
		},
		{
			path:          "Pool.EffectiveGasPrice.ByteGasCost",
			expectedValue: uint64(16),
		},
		{
			path:          "Pool.EffectiveGasPrice.ZeroByteGasCost",
			expectedValue: uint64(4),
		},
		{
			path:          "Pool.EffectiveGasPrice.NetProfit",
			expectedValue: float64(1),
		},
		{
			path:          "Pool.EffectiveGasPrice.BreakEvenFactor",
			expectedValue: float64(1.1),
		},
		{
			path:          "Pool.EffectiveGasPrice.FinalDeviationPct",
			expectedValue: uint64(10),
		},
		{
			path:          "Pool.EffectiveGasPrice.EthTransferGasPrice",
			expectedValue: uint64(0),
		},
		{
			path:          "Pool.EffectiveGasPrice.EthTransferL1GasPriceFactor",
			expectedValue: float64(0),
		},
		{
			path:          "Pool.DB.User",
			expectedValue: "pool_user",
		},
		{
			path:          "Pool.DB.Password",
			expectedValue: "pool_password",
		},
		{
			path:          "Pool.DB.Name",
			expectedValue: "pool_db",
		},
		{
			path:          "Pool.DB.Host",
			expectedValue: "zkevm-pool-db",
		},
		{
			path:          "Pool.DB.Port",
			expectedValue: "5432",
		},
		{
			path:          "Pool.DB.EnableLog",
			expectedValue: false,
		},
		{
			path:          "Pool.DB.MaxConns",
			expectedValue: 200,
		},
		{
			path:          "RPC.Host",
			expectedValue: "0.0.0.0",
		},
		{
			path:          "RPC.Port",
			expectedValue: int(8545),
		},
		{
			path:          "RPC.ReadTimeout",
			expectedValue: types.NewDuration(60 * time.Second),
		},
		{
			path:          "RPC.WriteTimeout",
			expectedValue: types.NewDuration(60 * time.Second),
		},
		{
			path:          "RPC.SequencerNodeURI",
			expectedValue: "",
		},
		{
			path:          "RPC.MaxRequestsPerIPAndSecond",
			expectedValue: float64(500),
		},
		{
			path:          "RPC.EnableL2SuggestedGasPricePolling",
			expectedValue: true,
		},
		{
			path:          "RPC.BatchRequestsEnabled",
			expectedValue: false,
		},
		{
			path:          "RPC.BatchRequestsLimit",
			expectedValue: uint(20),
		},
		{
			path:          "RPC.MaxLogsCount",
			expectedValue: uint64(10000),
		},
		{
			path:          "RPC.MaxLogsBlockRange",
			expectedValue: uint64(10000),
		},
		{
			path:          "RPC.MaxNativeBlockHashBlockRange",
			expectedValue: uint64(60000),
		},
		{
			path:          "RPC.EnableHttpLog",
			expectedValue: true,
		},
		{
			path:          "RPC.WebSockets.Enabled",
			expectedValue: true,
		},
		{
			path:          "RPC.WebSockets.Host",
			expectedValue: "0.0.0.0",
		},
		{
			path:          "RPC.WebSockets.Port",
			expectedValue: int(8546),
		},
		{
			path:          "RPC.WebSockets.ReadLimit",
			expectedValue: int64(104857600),
		},
		{
			path:          "Executor.URI",
			expectedValue: "zkevm-prover:50071",
		},
		{
			path:          "Executor.MaxResourceExhaustedAttempts",
			expectedValue: 3,
		},
		{
			path:          "Executor.WaitOnResourceExhaustion",
			expectedValue: types.NewDuration(1 * time.Second),
		},
		{
			path:          "Executor.MaxGRPCMessageSize",
			expectedValue: int(100000000),
		},
		{
			path:          "Metrics.Host",
			expectedValue: "0.0.0.0",
		},
		{
			path:          "Metrics.Port",
			expectedValue: 9091,
		},
		{
			path:          "Metrics.Enabled",
			expectedValue: false,
		},
		{
			path:          "Aggregator.Host",
			expectedValue: "0.0.0.0",
		},
		{
			path:          "Aggregator.Port",
			expectedValue: 50081,
		},
		{
			path:          "Aggregator.RetryTime",
			expectedValue: types.NewDuration(5 * time.Second),
		},
		{
			path:          "Aggregator.VerifyProofInterval",
			expectedValue: types.NewDuration(90 * time.Second),
		},
		{
			path:          "Aggregator.TxProfitabilityCheckerType",
			expectedValue: aggregator.TxProfitabilityCheckerType(aggregator.ProfitabilityAcceptAll),
		},
		{
			path:          "Aggregator.TxProfitabilityMinReward",
			expectedValue: aggregator.TokenAmountWithDecimals{Int: big.NewInt(1100000000000000000)},
		},
		{
			path:          "Aggregator.ProofStatePollingInterval",
			expectedValue: types.NewDuration(5 * time.Second),
		},
		{
			path:          "Aggregator.CleanupLockedProofsInterval",
			expectedValue: types.NewDuration(2 * time.Minute),
		},
		{
			path:          "Aggregator.GeneratingProofCleanupThreshold",
			expectedValue: "10m",
		},
		{
			path:          "Aggregator.GasOffset",
			expectedValue: uint64(0),
		},
		{
			path:          "Aggregator.UpgradeEtrogBatchNumber",
			expectedValue: uint64(0),
		},
		{
			path:          "Aggregator.BatchProofL1BlockConfirmations",
			expectedValue: uint64(2),
		},
		{
			path:          "State.Batch.Constraints.MaxTxsPerBatch",
			expectedValue: uint64(300),
		},
		{
			path:          "State.Batch.Constraints.MaxBatchBytesSize",
			expectedValue: uint64(120000),
		},
		{
			path:          "State.Batch.Constraints.MaxCumulativeGasUsed",
			expectedValue: uint64(1125899906842624),
		},
		{
			path:          "State.Batch.Constraints.MaxKeccakHashes",
			expectedValue: uint32(2145),
		},
		{
			path:          "State.Batch.Constraints.MaxPoseidonHashes",
			expectedValue: uint32(252357),
		},
		{
			path:          "State.Batch.Constraints.MaxPoseidonPaddings",
			expectedValue: uint32(135191),
		},
		{
			path:          "State.Batch.Constraints.MaxMemAligns",
			expectedValue: uint32(236585),
		},
		{
			path:          "State.Batch.Constraints.MaxArithmetics",
			expectedValue: uint32(236585),
		},
		{
			path:          "State.Batch.Constraints.MaxBinaries",
			expectedValue: uint32(473170),
		},
	}
	file, err := os.CreateTemp("", "genesisConfig")
	require.NoError(t, err)
	defer func() {
		require.NoError(t, os.Remove(file.Name()))
	}()
	require.NoError(t, os.WriteFile(file.Name(), []byte("{}"), 0600))

	flagSet := flag.NewFlagSet("", flag.PanicOnError)
	flagSet.String(config.FlagNetwork, "testnet", "")
	ctx := cli.NewContext(cli.NewApp(), flagSet, nil)
	cfg, err := config.Load(ctx, true)
	if err != nil {
		t.Fatalf("Unexpected error loading default config: %v", err)
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.path, func(t *testing.T) {
			actual := getValueFromStruct(tc.path, cfg)
			require.Equal(t, tc.expectedValue, actual)
		})
	}
}

func getValueFromStruct(path string, object interface{}) interface{} {
	keySlice := strings.Split(path, ".")
	v := reflect.ValueOf(object)

	for _, key := range keySlice {
		for v.Kind() == reflect.Ptr {
			v = v.Elem()
		}

		v = v.FieldByName(key)
	}
	return v.Interface()
}

func TestEnvVarArrayDecoding(t *testing.T) {
	file, err := os.CreateTemp("", "genesisConfig")
	require.NoError(t, err)
	defer func() {
		require.NoError(t, os.Remove(file.Name()))
	}()
	require.NoError(t, os.WriteFile(file.Name(), []byte("{}"), 0600))
	flagSet := flag.NewFlagSet("", flag.PanicOnError)
	flagSet.String(config.FlagNetwork, "testnet", "")
	ctx := cli.NewContext(cli.NewApp(), flagSet, nil)

	os.Setenv("ZKEVM_NODE_LOG_OUTPUTS", "a,b,c")
	defer func() {
		os.Unsetenv("ZKEVM_NODE_LOG_OUTPUTS")
	}()

	cfg, err := config.Load(ctx, true)
	require.NoError(t, err)

	assert.Equal(t, 3, len(cfg.Log.Outputs))
	assert.Equal(t, "a", cfg.Log.Outputs[0])
	assert.Equal(t, "b", cfg.Log.Outputs[1])
	assert.Equal(t, "c", cfg.Log.Outputs[2])
}
