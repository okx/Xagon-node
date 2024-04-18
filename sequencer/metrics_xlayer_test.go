package sequencer

import (
	"testing"
	"time"
)

func Test_Summary(t *testing.T) {
	tests := []struct {
		name    string
		metrics *metrics
	}{
		{"1", &metrics{
			closedAt:           time.Now(),
			processedTxsCount:  10,
			l2BlockTxsCount:    10,
			idleTime:           time.Second,
			newL2BlockTimes:    processTimes{sequencer: time.Second, executor: time.Second},
			transactionsTimes:  processTimes{sequencer: time.Second, executor: time.Second},
			l2BlockTimes:       processTimes{sequencer: time.Second, executor: time.Second},
			gas:                10,
			estimatedTxsPerSec: 10,
			estimatedGasPerSec: 10,
			closeReason:        "deadline",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Log(tt.metrics.Summary(1, 3, uint64(time.Now().Unix())))
		})
	}
}
