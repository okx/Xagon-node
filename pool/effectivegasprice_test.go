package pool

import (
	"math/big"
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/state"
	"github.com/stretchr/testify/assert"
)

var (
	egpCfg = EffectiveGasPriceCfg{
		Enabled:                     true,
		L1GasPriceFactor:            0.25,
		ByteGasCost:                 16,
		ZeroByteGasCost:             4,
		NetProfit:                   1,
		BreakEvenFactor:             1.1,
		FinalDeviationPct:           10,
		EthTransferGasPrice:         0,
		EthTransferL1GasPriceFactor: 0,
		L2GasPriceSuggesterFactor:   0.5,
	}
)

func TestCalculateEffectiveGasPricePercentage(t *testing.T) {
	testCases := []struct {
		name          string
		breakEven     *big.Int
		gasPrice      *big.Int
		expectedValue uint8
		err           error
	}{

		{
			name:          "Nil breakEven or gasPrice",
			gasPrice:      big.NewInt(1),
			expectedValue: uint8(0),
			err:           state.ErrEffectiveGasPriceEmpty,
		},
		{
			name:          "Zero breakEven or gasPrice",
			breakEven:     big.NewInt(1),
			gasPrice:      big.NewInt(0),
			expectedValue: uint8(0),
			err:           state.ErrEffectiveGasPriceEmpty,
		},
		{
			name:          "Both positive, gasPrice less than breakEven",
			breakEven:     big.NewInt(22000000000),
			gasPrice:      big.NewInt(11000000000),
			expectedValue: uint8(255),
		},
		{
			name:          "Both positive, gasPrice more than breakEven",
			breakEven:     big.NewInt(19800000000),
			gasPrice:      big.NewInt(22000000000),
			expectedValue: uint8(230),
		},
		{
			name:          "100% (255) effective percentage 1",
			gasPrice:      big.NewInt(22000000000),
			breakEven:     big.NewInt(22000000000),
			expectedValue: 255,
		},
		{
			name:          "100% (255) effective percentage 2",
			gasPrice:      big.NewInt(22000000000),
			breakEven:     big.NewInt(21999999999),
			expectedValue: 255,
		},
		{
			name:          "100% (255) effective percentage 3",
			gasPrice:      big.NewInt(22000000000),
			breakEven:     big.NewInt(21900000000),
			expectedValue: 254,
		},
		{
			name:          "50% (127) effective percentage",
			gasPrice:      big.NewInt(22000000000),
			breakEven:     big.NewInt(11000000000),
			expectedValue: 127,
		},
		{
			name:          "(40) effective percentage",
			gasPrice:      big.NewInt(1000),
			breakEven:     big.NewInt(157),
			expectedValue: 40,
		},
		{
			name:          "(1) effective percentage",
			gasPrice:      big.NewInt(1000),
			breakEven:     big.NewInt(1),
			expectedValue: 0,
		},
		{
			name:          "(2) effective percentage",
			gasPrice:      big.NewInt(1000),
			breakEven:     big.NewInt(4),
			expectedValue: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := state.CalculateEffectiveGasPricePercentage(tc.gasPrice, tc.breakEven)
			assert.Equal(t, tc.err, err)
			if actual != 0 {
				assert.Equal(t, tc.expectedValue, actual)
			} else {
				assert.Zero(t, tc.expectedValue)
			}
		})
	}
}

func TestCalculateBreakEvenGasPrice(t *testing.T) {
	egp := NewEffectiveGasPrice(egpCfg)

	testCases := []struct {
		name                        string
		rawTx                       []byte
		txGasPrice                  *big.Int
		txGasUsed                   uint64
		l1GasPrice                  uint64
		EthTransferGasPrice         uint64
		EthTransferL1GasPriceFactor float64
		expectedValue               *big.Int
		err                         error
	}{

		{
			name:          "Test empty tx",
			rawTx:         []byte{},
			txGasPrice:    new(big.Int).SetUint64(1000),
			txGasUsed:     200,
			l1GasPrice:    100,
			expectedValue: new(big.Int).SetUint64(33),
		},
		{
			name:          "Test l1GasPrice=0",
			rawTx:         []byte{},
			txGasPrice:    new(big.Int).SetUint64(1000),
			txGasUsed:     200,
			l1GasPrice:    0,
			expectedValue: new(big.Int).SetUint64(553),
			err:           ErrZeroL1GasPrice,
		},
		{
			name:          "Test txGasUsed=0",
			rawTx:         []byte{},
			txGasPrice:    new(big.Int).SetUint64(1000),
			txGasUsed:     0,
			l1GasPrice:    100,
			expectedValue: new(big.Int).SetUint64(1000),
		},
		{
			name:          "Test tx len=10, zeroByte=0",
			rawTx:         []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			txGasPrice:    new(big.Int).SetUint64(1000),
			txGasUsed:     200,
			l1GasPrice:    100,
			expectedValue: new(big.Int).SetUint64(113),
		},
		{
			name:          "Test tx len=10, zeroByte=10",
			rawTx:         []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			txGasPrice:    new(big.Int).SetUint64(1000),
			txGasUsed:     200,
			l1GasPrice:    100,
			expectedValue: new(big.Int).SetUint64(53),
		},
		{
			name:          "Test tx len=10, zeroByte=5",
			rawTx:         []byte{1, 0, 2, 0, 3, 0, 4, 0, 5, 0},
			txGasPrice:    new(big.Int).SetUint64(1000),
			txGasUsed:     200,
			l1GasPrice:    100,
			expectedValue: new(big.Int).SetUint64(83),
		},
		{
			name:          "Test breakEvenGasPrice = 0 must return 1",
			rawTx:         []byte{1, 0, 2, 0, 3, 0, 4, 0, 5, 0},
			txGasPrice:    new(big.Int).SetUint64(1000),
			txGasUsed:     200000,
			l1GasPrice:    1,
			expectedValue: new(big.Int).SetUint64(1),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := egp.CalculateBreakEvenGasPrice(tc.rawTx, tc.txGasPrice, tc.txGasUsed, tc.l1GasPrice)
			assert.Equal(t, tc.err, err)
			if err == nil {
				if actual.Cmp(new(big.Int).SetUint64(0)) != 0 {
					assert.Equal(t, tc.expectedValue, actual)
				} else {
					assert.Zero(t, tc.expectedValue)
				}
			}
		})
	}
}

func TestEthTransferGasPrice(t *testing.T) {
	testCases := []struct {
		name                        string
		rawTx                       []byte
		txGasPrice                  *big.Int
		txGasUsed                   uint64
		l1GasPrice                  uint64
		EthTransferGasPrice         uint64
		EthTransferL1GasPriceFactor float64
		expectedValue               *big.Int
		err                         error
	}{
		{
			name:                "Test set EthTransferGasPrice",
			rawTx:               []byte{1, 0, 2, 0, 3, 0, 4, 0, 5, 0},
			txGasPrice:          new(big.Int).SetUint64(1000),
			txGasUsed:           21000,
			l1GasPrice:          10000,
			EthTransferGasPrice: 2000,
			expectedValue:       new(big.Int).SetUint64(2000),
		},
		{
			name:                        "Test set EthTransferL1GasPriceFactor",
			rawTx:                       []byte{1, 0, 2, 0, 3, 0, 4, 0, 5, 0},
			txGasPrice:                  new(big.Int).SetUint64(1000),
			txGasUsed:                   21000,
			l1GasPrice:                  10000,
			EthTransferL1GasPriceFactor: 0.5,
			expectedValue:               new(big.Int).SetUint64(5000),
		},
		{
			name:                        "Test set No ETHTransfer",
			rawTx:                       []byte{1, 0, 2, 0, 3, 0, 4, 0, 5, 0},
			txGasPrice:                  new(big.Int).SetUint64(1000),
			txGasUsed:                   200,
			l1GasPrice:                  100,
			EthTransferL1GasPriceFactor: 0.5,
			expectedValue:               new(big.Int).SetUint64(83),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			egpCfg.EthTransferGasPrice = tc.EthTransferGasPrice
			egpCfg.EthTransferL1GasPriceFactor = tc.EthTransferL1GasPriceFactor

			egp := NewEffectiveGasPrice(egpCfg)
			actual, err := egp.CalculateBreakEvenGasPrice(tc.rawTx, tc.txGasPrice, tc.txGasUsed, tc.l1GasPrice)

			assert.Equal(t, nil, err)
			assert.Equal(t, tc.expectedValue, actual)
		})
	}
}

func TestCalculateEffectiveGasPrice(t *testing.T) {
	egp := NewEffectiveGasPrice(egpCfg)

	testCases := []struct {
		name          string
		rawTx         []byte
		txGasPrice    *big.Int
		txGasUsed     uint64
		l1GasPrice    uint64
		l2GasPrice    uint64
		expectedValue *big.Int
		err           error
	}{
		{
			name:          "Test tx len=10, zeroByte=0",
			rawTx:         []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			txGasPrice:    new(big.Int).SetUint64(1000),
			txGasUsed:     200,
			l1GasPrice:    100,
			l2GasPrice:    1000,
			expectedValue: new(big.Int).SetUint64(113),
		},
		{
			name:          "Test tx len=10, zeroByte=10",
			rawTx:         []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			txGasPrice:    new(big.Int).SetUint64(1000),
			txGasUsed:     200,
			l1GasPrice:    100,
			l2GasPrice:    500,
			expectedValue: new(big.Int).SetUint64(53 * 2),
		},
		{
			name:          "Test tx len=10, zeroByte=5",
			rawTx:         []byte{1, 0, 2, 0, 3, 0, 4, 0, 5, 0},
			txGasPrice:    new(big.Int).SetUint64(1000),
			txGasUsed:     200,
			l1GasPrice:    100,
			l2GasPrice:    250,
			expectedValue: new(big.Int).SetUint64(83 * 4),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := egp.CalculateEffectiveGasPrice(tc.rawTx, tc.txGasPrice, tc.txGasUsed, tc.l1GasPrice, tc.l2GasPrice)
			assert.Equal(t, tc.err, err)
			if err == nil {
				if actual.Cmp(new(big.Int).SetUint64(0)) != 0 {
					assert.Equal(t, tc.expectedValue, actual)
				} else {
					assert.Zero(t, tc.expectedValue)
				}
			}
		})
	}
}
