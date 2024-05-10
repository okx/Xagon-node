package pool

import (
	"testing"

	"github.com/0xPolygonHermez/zkevm-node/state"
)

func TestIsWithinConstraints(t *testing.T) {
	cfg := state.BatchConstraintsCfg{
		MaxCumulativeGasUsed: 500,
		MaxKeccakHashes:      100,
		MaxPoseidonHashes:    200,
		MaxPoseidonPaddings:  150,
		MaxMemAligns:         1000,
		MaxArithmetics:       2000,
		MaxBinaries:          3000,
		MaxSteps:             4000,
		MaxSHA256Hashes:      5000,
	}

	testCases := []struct {
		desc     string
		counters state.ZKCounters
		expected bool
	}{
		{
			desc: "All constraints within limits",
			counters: state.ZKCounters{
				GasUsed:          300,
				KeccakHashes:     50,
				PoseidonHashes:   100,
				PoseidonPaddings: 75,
				MemAligns:        500,
				Arithmetics:      1000,
				Binaries:         2000,
				Steps:            2000,
				Sha256Hashes_V2:  4000,
			},
			expected: true,
		},
		{
			desc: "All constraints exceed limits",
			counters: state.ZKCounters{
				GasUsed:          600,
				KeccakHashes:     150,
				PoseidonHashes:   300,
				PoseidonPaddings: 200,
				MemAligns:        2000,
				Arithmetics:      3000,
				Binaries:         4000,
				Steps:            5000,
				Sha256Hashes_V2:  6000,
			},
			expected: false,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			if got := cfg.IsWithinConstraints(tC.counters); got != tC.expected {
				t.Errorf("Expected %v, got %v", tC.expected, got)
			}
		})
	}
}
