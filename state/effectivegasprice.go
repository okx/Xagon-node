package state

import (
	"errors"
	"math/big"
)

const (
	// MaxEffectivePercentage is the maximum value that can be used as effective percentage
	MaxEffectivePercentage = uint8(255)
)

var (
	// ErrEffectiveGasPriceEmpty happens when the effectiveGasPrice or gasPrice is nil or zero
	ErrEffectiveGasPriceEmpty = errors.New("effectiveGasPrice or gasPrice cannot be nil or zero")

	// ErrEffectiveGasPriceIsZero happens when the calculated EffectiveGasPrice is zero
	ErrEffectiveGasPriceIsZero = errors.New("effectiveGasPrice cannot be zero")
)

// CalculateEffectiveGasPricePercentage calculates the gas price's effective percentage
func CalculateEffectiveGasPricePercentage(gasPrice *big.Int, effectiveGasPrice *big.Int) (uint8, error) {
	const bits = 256
	var bitsBigInt = big.NewInt(bits)

	if effectiveGasPrice == nil || gasPrice == nil ||
		gasPrice.Cmp(big.NewInt(0)) == 0 || effectiveGasPrice.Cmp(big.NewInt(0)) == 0 {
		return 0, ErrEffectiveGasPriceEmpty
	}

	if gasPrice.Cmp(effectiveGasPrice) <= 0 {
		return MaxEffectivePercentage, nil
	}

	// Simulate Ceil with integer division
	b := new(big.Int).Mul(effectiveGasPrice, bitsBigInt)
	b = b.Add(b, gasPrice)
	b = b.Sub(b, big.NewInt(1)) //nolint:gomnd
	b = b.Div(b, gasPrice)
	// At this point we have a percentage between 1-256, we need to sub 1 to have it between 0-255 (byte)
	b = b.Sub(b, big.NewInt(1)) //nolint:gomnd

	return uint8(b.Uint64()), nil
}
