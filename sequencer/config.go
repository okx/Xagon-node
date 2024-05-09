package sequencer

import (
	"github.com/0xPolygonHermez/zkevm-data-streamer/log"
	"github.com/0xPolygonHermez/zkevm-node/config/types"
)

// Config represents the configuration of a sequencer
type Config struct {
	// DeletePoolTxsL1BlockConfirmations is blocks amount after which txs will be deleted from the pool
	DeletePoolTxsL1BlockConfirmations uint64 `mapstructure:"DeletePoolTxsL1BlockConfirmations"`

	// DeletePoolTxsCheckInterval is frequency with which txs will be checked for deleting
	DeletePoolTxsCheckInterval types.Duration `mapstructure:"DeletePoolTxsCheckInterval"`

	// TxLifetimeCheckInterval is the time the sequencer waits to check txs lifetime
	TxLifetimeCheckInterval types.Duration `mapstructure:"TxLifetimeCheckInterval"`

	// TxLifetimeMax is the time a tx can be in the sequencer/worker memory
	TxLifetimeMax types.Duration `mapstructure:"TxLifetimeMax"`

	// LoadPoolTxsCheckInterval is the time the sequencer waits to check in there are new txs in the pool
	LoadPoolTxsCheckInterval types.Duration `mapstructure:"LoadPoolTxsCheckInterval"`

	// StateConsistencyCheckInterval is the time the sequencer waits to check if a state inconsistency has happened
	StateConsistencyCheckInterval types.Duration `mapstructure:"StateConsistencyCheckInterval"`

	// Finalizer's specific config properties
	Finalizer FinalizerCfg `mapstructure:"Finalizer"`

	// StreamServerCfg is the config for the stream server
	StreamServer StreamServerCfg `mapstructure:"StreamServer"`
}

// StreamServerCfg contains the data streamer's configuration properties
type StreamServerCfg struct {
	// Port to listen on
	Port uint16 `mapstructure:"Port"`
	// Filename of the binary data file
	Filename string `mapstructure:"Filename"`
	// Version of the binary data file
	Version uint8 `mapstructure:"Version"`
	// ChainID is the chain ID
	ChainID uint64 `mapstructure:"ChainID"`
	// Enabled is a flag to enable/disable the data streamer
	Enabled bool `mapstructure:"Enabled"`
	// Log is the log configuration
	Log log.Config `mapstructure:"Log"`
	// UpgradeEtrogBatchNumber is the batch number of the upgrade etrog
	UpgradeEtrogBatchNumber uint64 `mapstructure:"UpgradeEtrogBatchNumber"`
}

// FinalizerCfg contains the finalizer's configuration properties
type FinalizerCfg struct {
	// ForcedBatchesTimeout is the time the finalizer waits after receiving closing signal to process Forced Batches
	ForcedBatchesTimeout types.Duration `mapstructure:"ForcedBatchesTimeout"`

	// NewTxsWaitInterval is the time the finalizer sleeps between each iteration, if there are no transactions to be processed
	NewTxsWaitInterval types.Duration `mapstructure:"NewTxsWaitInterval"`

	// ResourceExhaustedMarginPct is the percentage window of the resource left out for the batch to be closed
	ResourceExhaustedMarginPct uint32 `mapstructure:"ResourceExhaustedMarginPct"`

	// ForcedBatchesL1BlockConfirmations is number of blocks to consider GER final
	ForcedBatchesL1BlockConfirmations uint64 `mapstructure:"ForcedBatchesL1BlockConfirmations"`

	// L1InfoTreeL1BlockConfirmations is number of blocks to consider L1InfoRoot final
	L1InfoTreeL1BlockConfirmations uint64 `mapstructure:"L1InfoTreeL1BlockConfirmations"`

	// ForcedBatchesCheckInterval is used by the closing signals manager to wait for its operation
	ForcedBatchesCheckInterval types.Duration `mapstructure:"ForcedBatchesCheckInterval"`

	// L1InfoTreeCheckInterval is the wait time to check if the L1InfoRoot has been updated
	L1InfoTreeCheckInterval types.Duration `mapstructure:"L1InfoTreeCheckInterval"`

	// BatchMaxDeltaTimestamp is the resolution of the timestamp used to close a batch
	BatchMaxDeltaTimestamp types.Duration `mapstructure:"BatchMaxDeltaTimestamp"`

	// L2BlockMaxDeltaTimestamp is the resolution of the timestamp used to close a L2 block
	L2BlockMaxDeltaTimestamp types.Duration `mapstructure:"L2BlockMaxDeltaTimestamp"`

	// HaltOnBatchNumber specifies the batch number where the Sequencer will stop to process more transactions and generate new batches.
	// The Sequencer will halt after it closes the batch equal to this number
	HaltOnBatchNumber uint64 `mapstructure:"HaltOnBatchNumber"`

	// SequentialBatchSanityCheck indicates if the reprocess of a closed batch (sanity check) must be done in a
	// sequential way (instead than in parallel)
	SequentialBatchSanityCheck bool `mapstructure:"SequentialBatchSanityCheck"`

	// SequentialProcessL2Block indicates if the processing of a L2 Block must be done in the same finalizer go func instead
	// in the processPendingL2Blocks go func
	SequentialProcessL2Block bool `mapstructure:"SequentialProcessL2Block"`

	// Metrics is the config for the sequencer metrics
	Metrics MetricsCfg `mapstructure:"Metrics"`
}

// MetricsCfg contains the sequencer metrics configuration properties
type MetricsCfg struct {
	// Interval is the interval of time to calculate sequencer metrics
	Interval types.Duration `mapstructure:"Interval"`

	// EnableLog is a flag to enable/disable metrics logs
	EnableLog bool `mapstructure:"EnableLog"`
}
