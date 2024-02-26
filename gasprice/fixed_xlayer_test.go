package gasprice

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/0xPolygonHermez/zkevm-node/config/types"
	"github.com/0xPolygonHermez/zkevm-node/log"
)

func init() {
	log.Init(log.Config{
		Level:   "debug",
		Outputs: []string{"stdout"},
	})
}

func TestUpdateGasPriceFixed(t *testing.T) {
	ctx := context.Background()
	var d time.Duration = 1000000000

	cfg := Config{
		Type:               FixedType,
		DefaultGasPriceWei: 1000000000,
		UpdatePeriod:       types.NewDuration(d),
		Factor:             0.5,
		KafkaURL:           "127.0.0.1:9092",
		Topic:              "middle_coinPrice_push",
		DefaultL2CoinPrice: 40,
		GasPriceUsdt:       0.001,
	}
	l1GasPrice := big.NewInt(10000000000)
	l2GasPrice := uint64(25000000000000)
	poolM := new(poolMock)
	ethM := new(ethermanMock)
	ethM.On("GetL1GasPrice", ctx).Return(l1GasPrice).Once()
	poolM.On("SetGasPrices", ctx, l2GasPrice, l1GasPrice.Uint64()).Return(nil).Once()
	f := newFixedGasPriceSuggester(ctx, cfg, poolM, ethM)

	ethM.On("GetL1GasPrice", ctx).Return(l1GasPrice, l1GasPrice).Once()
	poolM.On("SetGasPrices", ctx, l2GasPrice, l1GasPrice.Uint64()).Return(nil).Once()
	f.UpdateGasPriceAvg()
}

func TestUpdateGasPriceAvgCases(t *testing.T) {
	var d time.Duration = 1000000000
	testcases := []struct {
		cfg        Config
		l1GasPrice *big.Int
		l2GasPrice uint64
	}{
		{
			cfg: Config{
				Type:               FixedType,
				DefaultGasPriceWei: 1000000000,
				UpdatePeriod:       types.NewDuration(d),
				KafkaURL:           "127.0.0.1:9092",
				Topic:              "middle_coinPrice_push",
				DefaultL2CoinPrice: 40,
				GasPriceUsdt:       0.001,
			},
			l1GasPrice: big.NewInt(10000000000),
			l2GasPrice: uint64(25000000000000),
		},
		{
			cfg: Config{
				Type:               FixedType,
				DefaultGasPriceWei: 1000000000,
				UpdatePeriod:       types.NewDuration(d),
				KafkaURL:           "127.0.0.1:9092",
				Topic:              "middle_coinPrice_push",
				DefaultL2CoinPrice: 1e-19,
				GasPriceUsdt:       0.001,
			},
			l1GasPrice: big.NewInt(10000000000),
			l2GasPrice: uint64(25000000000000),
		},
		{ // the gas price small than the min gas price
			cfg: Config{
				Type:               FixedType,
				DefaultGasPriceWei: 26000000000000,
				UpdatePeriod:       types.NewDuration(d),
				KafkaURL:           "127.0.0.1:9092",
				Topic:              "middle_coinPrice_push",
				DefaultL2CoinPrice: 40,
				GasPriceUsdt:       0.001,
			},
			l1GasPrice: big.NewInt(10000000000),
			l2GasPrice: uint64(26000000000000),
		},
		{ // the gas price bigger than the max gas price
			cfg: Config{
				Type:               FixedType,
				DefaultGasPriceWei: 1000000000000,
				MaxGasPriceWei:     23000000000000,
				UpdatePeriod:       types.NewDuration(d),
				KafkaURL:           "127.0.0.1:9092",
				Topic:              "middle_coinPrice_push",
				DefaultL2CoinPrice: 40,
				GasPriceUsdt:       0.001,
			},
			l1GasPrice: big.NewInt(10000000000),
			l2GasPrice: uint64(23000000000000),
		},
		{
			cfg: Config{
				Type:               FixedType,
				DefaultGasPriceWei: 1000000000,
				UpdatePeriod:       types.NewDuration(d),
				KafkaURL:           "127.0.0.1:9092",
				Topic:              "middle_coinPrice_push",
				DefaultL2CoinPrice: 30,
				GasPriceUsdt:       0.001,
			},
			l1GasPrice: big.NewInt(10000000000),
			l2GasPrice: uint64(33300000000000),
		},
		{
			cfg: Config{
				Type:               FixedType,
				DefaultGasPriceWei: 10,
				UpdatePeriod:       types.NewDuration(d),
				KafkaURL:           "127.0.0.1:9092",
				Topic:              "middle_coinPrice_push",
				DefaultL2CoinPrice: 30,
				GasPriceUsdt:       1e-15,
			},
			l1GasPrice: big.NewInt(10000000000),
			l2GasPrice: uint64(33),
		},
	}

	for _, tc := range testcases {
		ctx := context.Background()
		poolM := new(poolMock)
		ethM := new(ethermanMock)
		ethM.On("GetL1GasPrice", ctx).Return(tc.l1GasPrice).Twice()
		poolM.On("SetGasPrices", ctx, tc.l2GasPrice, tc.l1GasPrice.Uint64()).Return(nil).Twice()
		f := newFixedGasPriceSuggester(ctx, tc.cfg, poolM, ethM)
		f.UpdateGasPriceAvg()
	}
}
