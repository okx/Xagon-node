package synchronizer

import (
	"fmt"

	"github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/0xPolygonHermez/zkevm-node/synchronizer/l2_sync"
)

// Config represents the configuration of the synchronizer
type Config struct {
	// SyncInterval is the delay interval between reading new rollup information
	SyncInterval types.Duration `mapstructure:"SyncInterval"`
	// SyncChunkSize is the number of blocks to sync on each chunk
	SyncChunkSize uint64 `mapstructure:"SyncChunkSize"`
	// TrustedSequencerURL is the rpc url to connect and sync the trusted state
	TrustedSequencerURL string `mapstructure:"TrustedSequencerURL"`
	// SyncBlockProtection specify the state to sync (lastest, finalized or safe)
	SyncBlockProtection string `mapstructure:"SyncBlockProtection"`

	// L1SyncCheckL2BlockHash if is true when a batch is closed is force to check  L2Block hash against trustedNode (only apply for permissionless)
	L1SyncCheckL2BlockHash bool `mapstructure:"L1SyncCheckL2BlockHash"`
	// L1SyncCheckL2BlockNumberhModulus is the modulus used to choose the l2block to check
	// a modules 5, for instance, means check all l2block multiples of 5 (10,15,20,...)
	L1SyncCheckL2BlockNumberhModulus uint64 `mapstructure:"L1SyncCheckL2BlockNumberhModulus"`

	L1BlockCheck L1BlockCheckConfig `mapstructure:"L1BlockCheck"`
	// L1SynchronizationMode define how to synchronize with L1:
	// - parallel: Request data to L1 in parallel, and process sequentially. The advantage is that executor is not blocked waiting for L1 data
	// - sequential: Request data to L1 and execute
	L1SynchronizationMode string `jsonschema:"enum=sequential,enum=parallel"`
	// L1ParallelSynchronization Configuration for parallel mode (if L1SynchronizationMode equal to 'parallel')
	L1ParallelSynchronization L1ParallelSynchronizationConfig
	// L2Synchronization Configuration for L2 synchronization
	L2Synchronization l2_sync.Config `mapstructure:"L2Synchronization"`
}

// L1BlockCheckConfig Configuration for L1 Block Checker
type L1BlockCheckConfig struct {
	// If enabled then the check l1 Block Hash is active
	Enabled bool `mapstructure:"Enabled"`
	// L1SafeBlockPoint is the point that a block is considered safe enough to be checked
	// it can be: finalized, safe,pending or latest
	L1SafeBlockPoint string `mapstructure:"L1SafeBlockPoint" jsonschema:"enum=finalized,enum=safe, enum=pending,enum=latest"`
	// L1SafeBlockOffset is the offset to add to L1SafeBlockPoint as a safe point
	// it can be positive or negative
	// Example: L1SafeBlockPoint= finalized, L1SafeBlockOffset= -10, then the safe block ten blocks before the finalized block
	L1SafeBlockOffset int `mapstructure:"L1SafeBlockOffset"`
	// ForceCheckBeforeStart if is true then the first time the system is started it will force to check all pending blocks
	ForceCheckBeforeStart bool `mapstructure:"ForceCheckBeforeStart"`

	// If enabled then the pre-check is active, will check blocks between L1SafeBlock and L1PreSafeBlock
	PreCheckEnabled bool `mapstructure:"PreCheckEnabled"`
	// L1PreSafeBlockPoint is the point that a block is considered safe enough to be checked
	// it can be: finalized, safe,pending or latest
	L1PreSafeBlockPoint string `mapstructure:"L1PreSafeBlockPoint" jsonschema:"enum=finalized,enum=safe, enum=pending,enum=latest"`
	// L1PreSafeBlockOffset is the offset to add to L1PreSafeBlockPoint as a safe point
	// it can be positive or negative
	// Example: L1PreSafeBlockPoint= finalized, L1PreSafeBlockOffset= -10, then the safe block ten blocks before the finalized block
	L1PreSafeBlockOffset int `mapstructure:"L1PreSafeBlockOffset"`
}

func (c *L1BlockCheckConfig) String() string {
	return fmt.Sprintf("Enable: %v, L1SafeBlockPoint: %s, L1SafeBlockOffset: %d, ForceCheckBeforeStart: %v", c.Enabled, c.L1SafeBlockPoint, c.L1SafeBlockOffset, c.ForceCheckBeforeStart)
}

// L1ParallelSynchronizationConfig Configuration for parallel mode (if UL1SynchronizationMode equal to 'parallel')
type L1ParallelSynchronizationConfig struct {
	// MaxClients Number of clients used to synchronize with L1
	MaxClients uint64
	// MaxPendingNoProcessedBlocks Size of the buffer used to store rollup information from L1, must be >= to NumberOfEthereumClientsToSync
	// sugested twice of NumberOfParallelOfEthereumClients
	MaxPendingNoProcessedBlocks uint64

	// RequestLastBlockPeriod is the time to wait to request the
	// last block to L1 to known if we need to retrieve more data.
	// This value only apply when the system is synchronized
	RequestLastBlockPeriod types.Duration

	// Consumer Configuration for the consumer of rollup information from L1
	PerformanceWarning L1PerformanceCheckConfig

	// RequestLastBlockTimeout Timeout for request LastBlock On L1
	RequestLastBlockTimeout types.Duration
	// RequestLastBlockMaxRetries Max number of retries to request LastBlock On L1
	RequestLastBlockMaxRetries int
	// StatisticsPeriod how ofter show a log with statistics (0 is disabled)
	StatisticsPeriod types.Duration
	// TimeOutMainLoop is the timeout for the main loop of the L1 synchronizer when is not updated
	TimeOutMainLoop types.Duration
	// RollupInfoRetriesSpacing is the minimum time between retries to request rollup info (it will sleep for fulfill this time) to avoid spamming L1
	RollupInfoRetriesSpacing types.Duration
	// FallbackToSequentialModeOnSynchronized if true switch to sequential mode if the system is synchronized
	FallbackToSequentialModeOnSynchronized bool
}

// L1PerformanceCheckConfig Configuration for the consumer of rollup information from L1
type L1PerformanceCheckConfig struct {
	// AceptableInacctivityTime is the expected maximum time that the consumer
	// could wait until new data is produced. If the time is greater it emmit a log to warn about
	// that. The idea is keep working the consumer as much as possible, so if the producer is not
	// fast enought then you could increse the number of parallel clients to sync with L1
	AceptableInacctivityTime types.Duration
	// ApplyAfterNumRollupReceived is the number of iterations to
	// start checking the time waiting for new rollup info data
	ApplyAfterNumRollupReceived int
}
