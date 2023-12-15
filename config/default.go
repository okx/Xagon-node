package config

// DefaultValues is the default configuration
const DefaultValues = `
IsTrustedSequencer = false
ForkUpgradeBatchNumber = 0
ForkUpgradeNewForkId = 0

[Log]
Environment = "development" # "production" or "development"
Level = "info"
Outputs = ["stderr"]

[State]
	[State.DB]
	User = "state_user"
	Password = "state_password"
	Name = "state_db"
	Host = "x1-state-db"
	Port = "5432"
	EnableLog = false	
	MaxConns = 200
	[State.Batch]
		[State.Batch.Constraints]
		MaxTxsPerBatch = 300
		MaxBatchBytesSize = 120000
		MaxCumulativeGasUsed = 30000000
		MaxKeccakHashes = 2145
		MaxPoseidonHashes = 252357
		MaxPoseidonPaddings = 135191
		MaxMemAligns = 236585
		MaxArithmetics = 236585
		MaxBinaries = 473170
		MaxSteps = 7570538

[Pool]
FreeClaimGasLimit = 150000
IntervalToRefreshBlockedAddresses = "5m"
IntervalToRefreshGasPrices = "5s"
MaxTxBytesSize=100132
MaxTxDataBytesSize=100000
DefaultMinGasPriceAllowed = 1000000000
MinAllowedGasPriceInterval = "5m"
PollMinAllowedGasPriceInterval = "15s"
AccountQueue = 64
GlobalQueue = 1024
FreeGasAddress = "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
    [Pool.EffectiveGasPrice]
	Enabled = false
	L1GasPriceFactor = 0.25
	ByteGasCost = 16
	ZeroByteGasCost = 4
	NetProfit = 1
	BreakEvenFactor = 1.1	
	FinalDeviationPct = 10
	L2GasPriceSuggesterFactor = 0.5
    [Pool.DB]
	User = "pool_user"
	Password = "pool_password"
	Name = "pool_db"
	Host = "x1-pool-db"
	Port = "5432"
	EnableLog = false
	MaxConns = 200

[Etherman]
URL = "http://localhost:8545"
ForkIDChunkSize = 20000
MultiGasProvider = false
	[Etherman.Etherscan]
		ApiKey = ""

[EthTxManager]
FrequencyToMonitorTxs = "1s"
WaitTxToBeMined = "2m"
ForcedGas = 0
GasPriceMarginFactor = 1
MaxGasPriceLimit = 0

[RPC]
Host = "0.0.0.0"
Port = 8545
ReadTimeout = "60s"
WriteTimeout = "60s"
MaxRequestsPerIPAndSecond = 500
SequencerNodeURI = ""
EnableL2SuggestedGasPricePolling = true
BatchRequestsEnabled = false
BatchRequestsLimit = 20
MaxLogsCount = 10000
MaxLogsBlockRange = 10000
MaxNativeBlockHashBlockRange = 60000
EnableHttpLog = true
	[RPC.WebSockets]
		Enabled = true
		Host = "0.0.0.0"
		Port = 8546
		ReadLimit = 104857600

[Synchronizer]
SyncInterval = "1s"
SyncChunkSize = 100
TrustedSequencerURL = "" # If it is empty or not specified, then the value is read from the smc
L1SynchronizationMode = "sequential" # "sequential" or "parallel"
	[Synchronizer.L1ParallelSynchronization]
		MaxClients = 10
		MaxPendingNoProcessedBlocks = 25
		RequestLastBlockPeriod = "5s"
		RequestLastBlockTimeout = "5s"
		RequestLastBlockMaxRetries = 3
		StatisticsPeriod = "5m"
		TimeoutMainLoop = "5m"
		RollupInfoRetriesSpacing= "5s"
		FallbackToSequentialModeOnSynchronized = false
		[Synchronizer.L1ParallelSynchronization.PerformanceWarning]
			AceptableInacctivityTime = "5s"
			ApplyAfterNumRollupReceived = 10

[Sequencer]
WaitPeriodPoolIsEmpty = "1s"
BlocksAmountForTxsToBeDeleted = 100
FrequencyToCheckTxsForDelete = "12h"
TxLifetimeCheckTimeout = "10m"
MaxTxLifetime = "3h"
	[Sequencer.Finalizer]
		GERDeadlineTimeout = "5s"
		ForcedBatchDeadlineTimeout = "60s"
		SleepDuration = "100ms"
		ResourcePercentageToCloseBatch = 10
		GERFinalityNumberOfBlocks = 64
		ClosingSignalsManagerWaitForCheckingL1Timeout = "10s"
		ClosingSignalsManagerWaitForCheckingGER = "10s"
		ClosingSignalsManagerWaitForCheckingForcedBatches = "10s"
		ForcedBatchesFinalityNumberOfBlocks = 64
		TimestampResolution = "10s"
		StopSequencerOnBatchNum = 0
		SequentialReprocessFullBatch = false
	[Sequencer.DBManager]
		PoolRetrievalInterval = "500ms"
		L2ReorgRetrievalInterval = "5s"
	[Sequencer.StreamServer]
		Port = 0
		Filename = ""
		Enabled = false

[SequenceSender]
WaitPeriodSendSequence = "5s"
LastBatchVirtualizationTimeMaxWaitPeriod = "5s"
MaxBatchesForL1 = 10
L2Coinbase = "0xf39fd6e51aad88f6f4ce6ab8827279cfffb92266"
PrivateKey = {Path = "/pk/sequencer.keystore", Password = "testonly"}
UseValidium = true
GasOffset = 80000

[Aggregator]
Host = "0.0.0.0"
Port = 50081
RetryTime = "5s"
VerifyProofInterval = "90s"
TxProfitabilityCheckerType = "acceptall"
TxProfitabilityMinReward = "1.1"
ProofStatePollingInterval = "5s"
CleanupLockedProofsInterval = "2m"
GeneratingProofCleanupThreshold = "10m"
GasOffset = 0

[L2GasPriceSuggester]
Type = "follower"
UpdatePeriod = "10s"
Factor = 0.15
DefaultGasPriceWei = 2000000000
MaxGasPriceWei = 0
CleanHistoryPeriod = "1h"
CleanHistoryTimeRetention = "5m"

[MTClient]
URI = "x1-prover:50061"

[Executor]
URI = "x1-prover:50071"
MaxResourceExhaustedAttempts = 3
WaitOnResourceExhaustion = "1s"
MaxGRPCMessageSize = 100000000

[Metrics]
Host = "0.0.0.0"
Port = 9091
Enabled = false

[HashDB]
User = "prover_user"
Password = "prover_pass"
Name = "prover_db"
Host = "x1-state-db"
Port = "5432"
EnableLog = false
MaxConns = 200
`
