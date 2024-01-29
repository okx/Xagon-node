package metrics

import (
	"github.com/0xPolygonHermez/zkevm-node/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	// Prefix for the metrics of the sequencer package.
	Prefix = "ethtxmanager_"
	// HaltCountName is the name of the metric that counts the halt count.
	HaltCountName = Prefix + "halt_count"
)

// Register the metrics for the sequencer package.
func Register() {
	var counters = []prometheus.CounterOpts{
		{
			Name: HaltCountName,
			Help: "[ETHTXMANAGER] total count of halt",
		},
	}

	metrics.RegisterCounters(counters...)
}

// HaltCount increases the counter for the eth-tx-manager halt count.
func HaltCount() {
	metrics.CounterAdd(HaltCountName, 1)
}
