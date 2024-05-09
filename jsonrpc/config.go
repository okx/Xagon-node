package jsonrpc

import (
	"github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/ethereum/go-ethereum/common"
)

// Config represents the configuration of the json rpc
type Config struct {
	// Host defines the network adapter that will be used to serve the HTTP requests
	Host string `mapstructure:"Host"`

	// Port defines the port to serve the endpoints via HTTP
	Port int `mapstructure:"Port"`

	// ReadTimeout is the HTTP server read timeout
	// check net/http.server.ReadTimeout and net/http.server.ReadHeaderTimeout
	ReadTimeout types.Duration `mapstructure:"ReadTimeout"`

	// WriteTimeout is the HTTP server write timeout
	// check net/http.server.WriteTimeout
	WriteTimeout types.Duration `mapstructure:"WriteTimeout"`

	// MaxRequestsPerIPAndSecond defines how much requests a single IP can
	// send within a single second
	MaxRequestsPerIPAndSecond float64 `mapstructure:"MaxRequestsPerIPAndSecond"`

	// SequencerNodeURI is used allow Non-Sequencer nodes
	// to relay transactions to the Sequencer node
	SequencerNodeURI string `mapstructure:"SequencerNodeURI"`

	// MaxCumulativeGasUsed is the max gas allowed per batch
	MaxCumulativeGasUsed uint64

	// WebSockets configuration
	WebSockets WebSocketsConfig `mapstructure:"WebSockets"`

	// EnableL2SuggestedGasPricePolling enables polling of the L2 gas price to block tx in the RPC with lower gas price.
	EnableL2SuggestedGasPricePolling bool `mapstructure:"EnableL2SuggestedGasPricePolling"`

	// BatchRequestsEnabled defines if the Batch requests are enabled or disabled
	BatchRequestsEnabled bool `mapstructure:"BatchRequestsEnabled"`

	// BatchRequestsLimit defines the limit of requests that can be incorporated into each batch request
	BatchRequestsLimit uint `mapstructure:"BatchRequestsLimit"`

	// L2Coinbase defines which address is going to receive the fees
	L2Coinbase common.Address

	// MaxLogsCount is a configuration to set the max number of logs that can be returned
	// in a single call to the state, if zero it means no limit
	MaxLogsCount uint64 `mapstructure:"MaxLogsCount"`

	// MaxLogsBlockRange is a configuration to set the max range for block number when querying TXs
	// logs in a single call to the state, if zero it means no limit
	MaxLogsBlockRange uint64 `mapstructure:"MaxLogsBlockRange"`

	// MaxNativeBlockHashBlockRange is a configuration to set the max range for block number when querying
	// native block hashes in a single call to the state, if zero it means no limit
	MaxNativeBlockHashBlockRange uint64 `mapstructure:"MaxNativeBlockHashBlockRange"`

	// EnableHttpLog allows the user to enable or disable the logs related to the HTTP
	// requests to be captured by the server.
	EnableHttpLog bool `mapstructure:"EnableHttpLog"`

	// ZKCountersLimits defines the ZK Counter limits
	ZKCountersLimits ZKCountersLimits
}

// ZKCountersLimits defines the ZK Counter limits
type ZKCountersLimits struct {
	MaxKeccakHashes     uint32
	MaxPoseidonHashes   uint32
	MaxPoseidonPaddings uint32
	MaxMemAligns        uint32
	MaxArithmetics      uint32
	MaxBinaries         uint32
	MaxSteps            uint32
	MaxSHA256Hashes     uint32
}

// WebSocketsConfig has parameters to config the rpc websocket support
type WebSocketsConfig struct {
	// Enabled defines if the WebSocket requests are enabled or disabled
	Enabled bool `mapstructure:"Enabled"`

	// Host defines the network adapter that will be used to serve the WS requests
	Host string `mapstructure:"Host"`

	// Port defines the port to serve the endpoints via WS
	Port int `mapstructure:"Port"`

	// ReadLimit defines the maximum size of a message read from the client (in bytes)
	ReadLimit int64 `mapstructure:"ReadLimit"`
}
