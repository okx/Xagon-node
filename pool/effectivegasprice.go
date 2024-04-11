package pool

import (
	"bytes"
	"math/big"

	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/0xPolygonHermez/zkevm-node/state"
)

// EffectiveGasPrice implements the effective gas prices calculations and checks
type EffectiveGasPrice struct {
	cfg EffectiveGasPriceCfg
}

// NewEffectiveGasPrice creates and initializes an instance of EffectiveGasPrice
func NewEffectiveGasPrice(cfg EffectiveGasPriceCfg) *EffectiveGasPrice {
	if (cfg.EthTransferGasPrice != 0) && (cfg.EthTransferL1GasPriceFactor != 0) {
		log.Fatalf("configuration error. Only one of the following config params EthTransferGasPrice or EthTransferL1GasPriceFactor from Pool.effectiveGasPrice section can be set to a value different to 0")
	}
	return &EffectiveGasPrice{
		cfg: cfg,
	}
}

// IsEnabled return if effectiveGasPrice calculation is enabled
func (e *EffectiveGasPrice) IsEnabled() bool {
	return e.cfg.Enabled
}

// GetFinalDeviation return the value for the config parameter FinalDeviationPct
func (e *EffectiveGasPrice) GetFinalDeviation() uint64 {
	return e.cfg.FinalDeviationPct
}

// GetTxAndL2GasPrice return the tx gas price and l2 suggested gas price to use in egp calculations
// If egp is disabled we will use a "simulated" tx and l2 gas price, that is calculated using the L2GasPriceSuggesterFactor config param
func (e *EffectiveGasPrice) GetTxAndL2GasPrice(txGasPrice *big.Int, l1GasPrice uint64, l2GasPrice uint64) (egpTxGasPrice *big.Int, egpL2GasPrice uint64) {
	if !e.cfg.Enabled {
		// If egp is not enabled we use the L2GasPriceSuggesterFactor to calculate the "simulated" suggested L2 gas price
		gp := new(big.Int).SetUint64(uint64(e.cfg.L2GasPriceSuggesterFactor * float64(l1GasPrice)))
		return gp, gp.Uint64()
	} else {
		return txGasPrice, l2GasPrice
	}
}

// CalculateBreakEvenGasPrice calculates the break even gas price for a transaction
func (e *EffectiveGasPrice) CalculateBreakEvenGasPrice(rawTx []byte, txGasPrice *big.Int, txGasUsed uint64, l1GasPrice uint64) (*big.Int, error) {
	const ethTransferGas = 21000

	if l1GasPrice == 0 {
		return nil, ErrZeroL1GasPrice
	}

	if txGasUsed == 0 {
		// Returns tx.GasPrice as the breakEvenGasPrice
		return txGasPrice, nil
	}

	// If the tx is a ETH transfer (gas == 21000) then check if we need to return a "fix" effective gas price
	if txGasUsed == ethTransferGas {
		if e.cfg.EthTransferGasPrice != 0 {
			return new(big.Int).SetUint64(e.cfg.EthTransferGasPrice), nil
		} else if e.cfg.EthTransferL1GasPriceFactor != 0 {
			ethGasPrice := uint64(float64(l1GasPrice) * e.cfg.EthTransferL1GasPriceFactor)
			if ethGasPrice == 0 {
				ethGasPrice = 1
			}
			return new(big.Int).SetUint64(ethGasPrice), nil
		}
	}

	// Get L2 Min Gas Price
	l2MinGasPrice := uint64(float64(l1GasPrice) * e.cfg.L1GasPriceFactor)

	txZeroBytes := uint64(bytes.Count(rawTx, []byte{0}))
	txNonZeroBytes := uint64(len(rawTx)) - txZeroBytes + state.EfficiencyPercentageByteLength

	// Calculate BreakEvenGasPrice
	totalTxPrice := (txGasUsed * l2MinGasPrice) +
		((txNonZeroBytes*e.cfg.ByteGasCost)+(txZeroBytes*e.cfg.ZeroByteGasCost))*l1GasPrice
	breakEvenGasPrice := new(big.Int).SetUint64(uint64(float64(totalTxPrice/txGasUsed) * e.cfg.NetProfit))

	if breakEvenGasPrice.Cmp(new(big.Int).SetUint64(0)) == 0 {
		breakEvenGasPrice.SetUint64(1)
	}

	return breakEvenGasPrice, nil
}

// CalculateEffectiveGasPrice calculates the final effective gas price for a tx
func (e *EffectiveGasPrice) CalculateEffectiveGasPrice(rawTx []byte, txGasPrice *big.Int, txGasUsed uint64, l1GasPrice uint64, l2GasPrice uint64) (*big.Int, error) {
	breakEvenGasPrice, err := e.CalculateBreakEvenGasPrice(rawTx, txGasPrice, txGasUsed, l1GasPrice)

	if err != nil {
		return nil, err
	}

	bfL2GasPrice := new(big.Float).SetUint64(l2GasPrice)
	bfTxGasPrice := new(big.Float).SetInt(txGasPrice)

	ratioPriority := new(big.Float).SetFloat64(1.0)

	if bfL2GasPrice.Cmp(new(big.Float).SetUint64(0)) == 1 && bfTxGasPrice.Cmp(bfL2GasPrice) == 1 {
		//ratioPriority = (txGasPrice / l2GasPrice)
		ratioPriority = new(big.Float).Quo(bfTxGasPrice, bfL2GasPrice)
	}

	bfEffectiveGasPrice := new(big.Float).Mul(new(big.Float).SetInt(breakEvenGasPrice), ratioPriority)

	effectiveGasPrice := new(big.Int)
	bfEffectiveGasPrice.Int(effectiveGasPrice)

	if effectiveGasPrice.Cmp(new(big.Int).SetUint64(0)) == 0 {
		return nil, state.ErrEffectiveGasPriceIsZero
	}

	return effectiveGasPrice, nil
}
