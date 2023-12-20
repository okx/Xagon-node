package gasprice

import (
	"context"
	"fmt"
	"math/big"
	"strconv"

	"github.com/0xPolygonHermez/zkevm-node/encoding"
	"github.com/0xPolygonHermez/zkevm-node/log"
)

// FollowerGasPrice struct.
type FollowerGasPrice struct {
	cfg      Config
	pool     poolInterface
	ctx      context.Context
	eth      ethermanInterface
	kafkaPrc *KafkaProcessor
}

// newFollowerGasPriceSuggester inits l2 follower gas price suggester which is based on the l1 gas price.
func newFollowerGasPriceSuggester(ctx context.Context, cfg Config, pool poolInterface, ethMan ethermanInterface) *FollowerGasPrice {
	gps := &FollowerGasPrice{
		cfg:  cfg,
		pool: pool,
		ctx:  ctx,
		eth:  ethMan,
	}
	if cfg.EnableFollowerAdjustByL2L1Price {
		gps.kafkaPrc = newKafkaProcessor(cfg, ctx)
	}
	gps.UpdateGasPriceAvg()
	return gps
}

// UpdateGasPriceAvg updates the gas price.
func (f *FollowerGasPrice) UpdateGasPriceAvg() {
	ctx := context.Background()
	// Get L1 gasprice
	l1GasPrice := f.eth.GetL1GasPrice(f.ctx)
	if big.NewInt(0).Cmp(l1GasPrice) == 0 {
		log.Warn("gas price 0 received. Skipping update...")
		return
	}

	// Apply factor to calculate l2 gasPrice
	factor := big.NewFloat(0).SetFloat64(f.cfg.Factor)
	res := new(big.Float).Mul(factor, big.NewFloat(0).SetInt(l1GasPrice))

	// convert the eth gas price to okb gas price
	if f.cfg.EnableFollowerAdjustByL2L1Price {
		l1CoinPrice, l2CoinPrice := f.kafkaPrc.GetL1L2CoinPrice()
		if l1CoinPrice < minCoinPrice || l2CoinPrice < minCoinPrice {
			log.Warn("the L1 or L2 native coin price too small...")
			return
		}
		res = new(big.Float).Mul(big.NewFloat(0).SetFloat64(l1CoinPrice/l2CoinPrice), res)
		log.Debug("L2 pre gas price value: ", res.String(), ". L1 coin price: ", l1CoinPrice, ". L2 coin price: ", l2CoinPrice)
	}

	// Store l2 gasPrice calculated
	result := new(big.Int)
	res.Int(result)
	minGasPrice := big.NewInt(0).SetUint64(f.cfg.DefaultGasPriceWei)
	if minGasPrice.Cmp(result) == 1 { // minGasPrice > result
		log.Warn("setting DefaultGasPriceWei for L2")
		result = minGasPrice
	}
	maxGasPrice := new(big.Int).SetUint64(f.cfg.MaxGasPriceWei)
	if f.cfg.MaxGasPriceWei > 0 && result.Cmp(maxGasPrice) == 1 { // result > maxGasPrice
		log.Warn("setting MaxGasPriceWei for L2")
		result = maxGasPrice
	}
	var truncateValue *big.Int
	log.Debug("Full L2 gas price value: ", result, ". Length: ", len(result.String()))
	numLength := len(result.String())
	if numLength > 3 { //nolint:gomnd
		aux := "%0" + strconv.Itoa(numLength-3) + "d" //nolint:gomnd
		var ok bool
		value := result.String()[:3] + fmt.Sprintf(aux, 0)
		truncateValue, ok = new(big.Int).SetString(value, encoding.Base10)
		if !ok {
			log.Error("error converting: ", truncateValue)
		}
	} else {
		truncateValue = result
	}
	log.Debug("Storing truncated L2 gas price: ", truncateValue)
	if truncateValue != nil {
		log.Infof("Set gas prices, L1: %v, L2: %v", l1GasPrice.Uint64(), truncateValue.Uint64())
		err := f.pool.SetGasPrices(ctx, truncateValue.Uint64(), l1GasPrice.Uint64())
		if err != nil {
			log.Errorf("failed to update gas price in poolDB, err: %v", err)
		}
	} else {
		log.Error("nil value detected. Skipping...")
	}
}
