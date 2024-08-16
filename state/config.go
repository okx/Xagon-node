package state

import (
	"fmt"

	"github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/0xPolygonHermez/zkevm-node/db"
)

// Config is state config
type Config struct {
	// MaxCumulativeGasUsed is the max gas allowed per batch
	MaxCumulativeGasUsed uint64

	// ChainID is the L2 ChainID provided by the Network Config
	ChainID uint64

	// ForkIdIntervals is the list of fork id intervals
	ForkIDIntervals []ForkIDInterval

	// MaxResourceExhaustedAttempts is the max number of attempts to make a transaction succeed because of resource exhaustion
	MaxResourceExhaustedAttempts int

	// WaitOnResourceExhaustion is the time to wait before retrying a transaction because of resource exhaustion
	WaitOnResourceExhaustion types.Duration

	// Batch number from which there is a forkid change (fork upgrade)
	ForkUpgradeBatchNumber uint64

	// New fork id to be used for batches greaters than ForkUpgradeBatchNumber (fork upgrade)
	ForkUpgradeNewForkId uint64

	// DB is the database configuration
	DB db.Config `mapstructure:"DB"`

	// Configuration for the batch constraints
	Batch BatchConfig `mapstructure:"Batch"`

	// MaxLogsCount is a configuration to set the max number of logs that can be returned
	// in a single call to the state, if zero it means no limit
	MaxLogsCount uint64

	// MaxLogsBlockRange is a configuration to set the max range for block number when querying TXs
	// logs in a single call to the state, if zero it means no limit
	MaxLogsBlockRange uint64

	// MaxNativeBlockHashBlockRange is a configuration to set the max range for block number when querying
	// native block hashes in a single call to the state, if zero it means no limit
	MaxNativeBlockHashBlockRange uint64

	// AvoidForkIDInMemory is a configuration that forces the ForkID information to be loaded
	// from the DB every time it's needed
	AvoidForkIDInMemory bool
}

// BatchConfig represents the configuration of the batch constraints
type BatchConfig struct {
	Constraints BatchConstraintsCfg `mapstructure:"Constraints"`
}

// BatchConstraintsCfg represents the configuration of the batch constraints
type BatchConstraintsCfg struct {
	MaxTxsPerBatch       uint64 `mapstructure:"MaxTxsPerBatch"`
	MaxBatchBytesSize    uint64 `mapstructure:"MaxBatchBytesSize"`
	MaxCumulativeGasUsed uint64 `mapstructure:"MaxCumulativeGasUsed"`
	MaxKeccakHashes      uint32 `mapstructure:"MaxKeccakHashes"`
	MaxPoseidonHashes    uint32 `mapstructure:"MaxPoseidonHashes"`
	MaxPoseidonPaddings  uint32 `mapstructure:"MaxPoseidonPaddings"`
	MaxMemAligns         uint32 `mapstructure:"MaxMemAligns"`
	MaxArithmetics       uint32 `mapstructure:"MaxArithmetics"`
	MaxBinaries          uint32 `mapstructure:"MaxBinaries"`
	MaxSteps             uint32 `mapstructure:"MaxSteps"`
	MaxSHA256Hashes      uint32 `mapstructure:"MaxSHA256Hashes"`
}

// CheckNodeLevelOOC checks if the counters are within the batch constraints
func (c BatchConstraintsCfg) CheckNodeLevelOOC(counters ZKCounters) error {
	oocList := ""

	if counters.GasUsed > c.MaxCumulativeGasUsed {
		oocList += "GasUsed, "
	}
	if counters.KeccakHashes > c.MaxKeccakHashes {
		oocList += "KeccakHashes, "
	}
	if counters.PoseidonHashes > c.MaxPoseidonHashes {
		oocList += "PoseidonHashes, "
	}
	if counters.PoseidonPaddings > c.MaxPoseidonPaddings {
		oocList += "PoseidonPaddings, "
	}
	if counters.MemAligns > c.MaxMemAligns {
		oocList += "MemAligns, "
	}
	if counters.Arithmetics > c.MaxArithmetics {
		oocList += "Arithmetics, "
	}
	if counters.Binaries > c.MaxBinaries {
		oocList += "Binaries, "
	}
	if counters.Steps > c.MaxSteps {
		oocList += "Steps, "
	}
	if counters.Sha256Hashes_V2 > c.MaxSHA256Hashes {
		oocList += "Sha256Hashes, "
	}

	if oocList != "" {
		oocList = oocList[:len(oocList)-2] // Remove last comma and blank space
		return fmt.Errorf("out of counters at node level (%s)", oocList)
	}

	return nil
}
