package pool

import (
	"github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/0xPolygonHermez/zkevm-node/db"
)

// Config is the pool configuration
type Config struct {
	// IntervalToRefreshBlockedAddresses is the time it takes to sync the
	// blocked address list from db to memory
	IntervalToRefreshBlockedAddresses types.Duration `mapstructure:"IntervalToRefreshBlockedAddresses"`

	// IntervalToRefreshGasPrices is the time to wait to refresh the gas prices
	IntervalToRefreshGasPrices types.Duration `mapstructure:"IntervalToRefreshGasPrices"`

	// MaxTxBytesSize is the max size of a transaction in bytes
	MaxTxBytesSize uint64 `mapstructure:"MaxTxBytesSize"`

	// MaxTxDataBytesSize is the max size of the data field of a transaction in bytes
	MaxTxDataBytesSize int `mapstructure:"MaxTxDataBytesSize"`

	// DB is the database configuration
	DB db.Config `mapstructure:"DB"`

	// DefaultMinGasPriceAllowed is the default min gas price to suggest
	DefaultMinGasPriceAllowed uint64 `mapstructure:"DefaultMinGasPriceAllowed"`

	// MinAllowedGasPriceInterval is the interval to look back of the suggested min gas price for a tx
	MinAllowedGasPriceInterval types.Duration `mapstructure:"MinAllowedGasPriceInterval"`

	// PollMinAllowedGasPriceInterval is the interval to poll the suggested min gas price for a tx
	PollMinAllowedGasPriceInterval types.Duration `mapstructure:"PollMinAllowedGasPriceInterval"`

	// AccountQueue represents the maximum number of non-executable transaction slots permitted per account
	AccountQueue uint64 `mapstructure:"AccountQueue"`

	// GlobalQueue represents the maximum number of non-executable transaction slots for all accounts
	GlobalQueue uint64 `mapstructure:"GlobalQueue"`

	// EffectiveGasPrice is the config for the effective gas price calculation
	EffectiveGasPrice EffectiveGasPriceCfg `mapstructure:"EffectiveGasPrice"`

	// ForkID is the current fork ID of the chain
	ForkID uint64 `mapstructure:"ForkID"`

	// TxFeeCap is the global transaction fee(price * gaslimit) cap for
	// send-transaction variants. The unit is ether. 0 means no cap.
	TxFeeCap float64 `mapstructure:"TxFeeCap"`
}

// EffectiveGasPriceCfg contains the configuration properties for the effective gas price
type EffectiveGasPriceCfg struct {
	// Enabled is a flag to enable/disable the effective gas price
	Enabled bool `mapstructure:"Enabled"`

	// L1GasPriceFactor is the percentage of the L1 gas price that will be used as the L2 min gas price
	L1GasPriceFactor float64 `mapstructure:"L1GasPriceFactor"`

	// ByteGasCost is the gas cost per byte that is not 0
	ByteGasCost uint64 `mapstructure:"ByteGasCost"`

	// ZeroByteGasCost is the gas cost per byte that is 0
	ZeroByteGasCost uint64 `mapstructure:"ZeroByteGasCost"`

	// NetProfit is the profit margin to apply to the calculated breakEvenGasPrice
	NetProfit float64 `mapstructure:"NetProfit"`

	// BreakEvenFactor is the factor to apply to the calculated breakevenGasPrice when comparing it with the gasPriceSigned of a tx
	BreakEvenFactor float64 `mapstructure:"BreakEvenFactor"`

	// FinalDeviationPct is the max allowed deviation percentage BreakEvenGasPrice on re-calculation
	FinalDeviationPct uint64 `mapstructure:"FinalDeviationPct"`

	// EthTransferGasPrice is the fixed gas price returned as effective gas price for txs tha are ETH transfers (0 means disabled)
	// Only one of EthTransferGasPrice or EthTransferL1GasPriceFactor params can be different than 0. If both params are set to 0, the sequencer will halt and log an error
	EthTransferGasPrice uint64 `mapstructure:"EthTransferGasPrice"`

	// EthTransferL1GasPriceFactor is the percentage of L1 gas price returned as effective gas price for txs tha are ETH transfers (0 means disabled)
	// Only one of EthTransferGasPrice or EthTransferL1GasPriceFactor params can be different than 0. If both params are set to 0, the sequencer will halt and log an error
	EthTransferL1GasPriceFactor float64 `mapstructure:"EthTransferL1GasPriceFactor"`

	// L2GasPriceSuggesterFactor is the factor to apply to L1 gas price to get the suggested L2 gas price used in the
	// calculations when the effective gas price is disabled (testing/metrics purposes)
	L2GasPriceSuggesterFactor float64 `mapstructure:"L2GasPriceSuggesterFactor"`
}
