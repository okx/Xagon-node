package nubit

import "time"

// NubitNamespaceBytesLength is the fixed-size bytes array.
const NubitNamespaceBytesLength = 58

// Config is the NubitDA backend configurations
type Config struct {
	NubitRpcURL             string        `mapstructure:"NubitRpcURL"`
	NubitModularAppName     string        `mapstructure:"NubitModularAppName"`
	NubitAuthKey            string        `mapstructure:"NubitAuthKey"`
	NubitNamespace          string        `mapstructure:"NubitNamespace"`
	NubitMaxBatchesSize     uint64        `mapstructure:"NubitMaxBatchesSize"`
	NubitGetProofMaxRetry   uint64        `mapstructure:"NubitGetProofMaxRetry"`
	NubitGetProofWaitPeriod time.Duration `mapstructure:"NubitGetProofWaitPeriod"`
}
