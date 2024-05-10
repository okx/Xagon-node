package etherman

import "github.com/0xPolygonHermez/zkevm-node/etherman/etherscan"

// Config represents the configuration of the etherman
type Config struct {
	// URL is the URL of the Ethereum node for L1
	URL string `mapstructure:"URL"`
	// ConsensusL1URL is the URL of the consensus L1 RPC endpoint
	ConsensusL1URL string `mapstructure:"ConsensusL1URL"`

	// ForkIDChunkSize is the max interval for each call to L1 provider to get the forkIDs
	ForkIDChunkSize uint64 `mapstructure:"ForkIDChunkSize"`

	// allow that L1 gas price calculation use multiples sources
	MultiGasProvider bool `mapstructure:"MultiGasProvider"`
	// Configuration for use Etherscan as used as gas provider, basically it needs the API-KEY
	Etherscan etherscan.Config
}
