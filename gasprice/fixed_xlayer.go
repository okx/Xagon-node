package gasprice

import (
	"context"
	"fmt"
	"math/big"
	"strconv"

	"github.com/0xPolygonHermez/zkevm-node/encoding"
	"github.com/0xPolygonHermez/zkevm-node/log"
)

const (
	// OKBWei OKB wei
	OKBWei       = 1e18
	minCoinPrice = 1e-18
)

// FixedGasPrice struct
type FixedGasPrice struct {
	cfg     Config
	pool    poolInterface
	ctx     context.Context
	eth     ethermanInterface
	ratePrc *KafkaProcessor
}

// newFixedGasPriceSuggester inits l2 fixed price suggester.
func newFixedGasPriceSuggester(ctx context.Context, cfg Config, pool poolInterface, ethMan ethermanInterface) *FixedGasPrice {
	gps := &FixedGasPrice{
		cfg:     cfg,
		pool:    pool,
		ctx:     ctx,
		eth:     ethMan,
		ratePrc: newKafkaProcessor(cfg, ctx),
	}
	gps.UpdateGasPriceAvg()
	return gps
}

// UpdateGasPriceAvg updates the gas price.
func (f *FixedGasPrice) UpdateGasPriceAvg() {
	if getApolloConfig().Enable() {
		f.cfg = getApolloConfig().get()
	}

	ctx := context.Background()
	// Get L1 gasprice
	l1GasPrice := f.eth.GetL1GasPrice(f.ctx)
	if big.NewInt(0).Cmp(l1GasPrice) == 0 {
		log.Warn("gas price 0 received. Skipping update...")
		return
	}

	l2CoinPrice := f.ratePrc.GetL2CoinPrice()
	if l2CoinPrice < minCoinPrice {
		log.Warn("the L2 native coin price too small...")
		return
	}
	res := new(big.Float).Mul(big.NewFloat(0).SetFloat64(f.cfg.GasPriceUsdt/l2CoinPrice), big.NewFloat(0).SetFloat64(OKBWei))
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
	log.Debug("Full L2 gas price value: ", result, ". Length: ", len(result.String()), ". L1 gas price value: ", l1GasPrice)

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
	log.Debugf("Storing truncated L2 gas price: %v, L2 native coin price: %v", truncateValue, l2CoinPrice)
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
