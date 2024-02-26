package gasprice

import (
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/config/types"
)

// EstimatorType different gas estimator types.
type EstimatorType string

const (
	// DefaultType default gas price from config is set.
	DefaultType EstimatorType = "default"
	// LastNBatchesType calculate average gas tip from last n batches.
	LastNBatchesType EstimatorType = "lastnbatches"
	// FollowerType calculate the gas price basing on the L1 gasPrice.
	FollowerType EstimatorType = "follower"

	// FixedType the gas price from config that the unit is usdt, XLayer config
	FixedType EstimatorType = "fixed"
)

// Config for gas price estimator.
type Config struct {
	Type EstimatorType `mapstructure:"Type"`

	// DefaultGasPriceWei is used to set the gas price to be used by the default gas pricer or as minimim gas price by the follower gas pricer.
	DefaultGasPriceWei uint64 `mapstructure:"DefaultGasPriceWei"`
	// MaxGasPriceWei is used to limit the gas price returned by the follower gas pricer to a maximum value. It is ignored if 0.
	MaxGasPriceWei            uint64         `mapstructure:"MaxGasPriceWei"`
	MaxPrice                  *big.Int       `mapstructure:"MaxPrice"`
	IgnorePrice               *big.Int       `mapstructure:"IgnorePrice"`
	CheckBlocks               int            `mapstructure:"CheckBlocks"`
	Percentile                int            `mapstructure:"Percentile"`
	UpdatePeriod              types.Duration `mapstructure:"UpdatePeriod"`
	CleanHistoryPeriod        types.Duration `mapstructure:"CleanHistoryPeriod"`
	CleanHistoryTimeRetention types.Duration `mapstructure:"CleanHistoryTimeRetention"`

	Factor float64 `mapstructure:"Factor"`

	// XLayer config
	KafkaURL   string `mapstructure:"KafkaURL"`
	Topic      string `mapstructure:"Topic"`
	GroupID    string `mapstructure:"GroupID"`
	Username   string `mapstructure:"Username"`
	Password   string `mapstructure:"Password"`
	RootCAPath string `mapstructure:"RootCAPath"`
	L1CoinId   int    `mapstructure:"L1CoinId"`
	L2CoinId   int    `mapstructure:"L2CoinId"`
	// DefaultL1CoinPrice is the L1 token's coin price
	DefaultL1CoinPrice float64 `mapstructure:"DefaultL1CoinPrice"`
	// DefaultL2CoinPrice is the native token's coin price
	DefaultL2CoinPrice float64 `mapstructure:"DefaultL2CoinPrice"`
	GasPriceUsdt       float64 `mapstructure:"GasPriceUsdt"`

	// EnableFollowerAdjustByL2L1Price is dynamic adjust the factor through the L1 and L2 coins price in follower strategy
	EnableFollowerAdjustByL2L1Price bool `mapstructure:"EnableFollowerAdjustByL2L1Price"`
}
