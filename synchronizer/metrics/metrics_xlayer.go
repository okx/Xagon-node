package metrics

import (
	"github.com/0xPolygonHermez/zkevm-node/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// VirtualBatchNumName is the name of the metric virtual batch number
	VirtualBatchNumName = Prefix + "virtual_batch_num"

	// VerifiedBatchNumName is the name of the metric verified batch number
	VerifiedBatchNumName = Prefix + "verified_batch_num"

	// HaltCountName is the name of the metric that counts synchronizer halt count
	HaltCountName = Prefix + "halt_count"

	gauge = []prometheus.GaugeOpts{
		{
			Name: VirtualBatchNumName,
			Help: "[SYNCHRONIZER] virtual batch num",
		},
		{
			Name: VerifiedBatchNumName,
			Help: "[SYNCHRONIZER] verified batch num",
		},
	}
	counters = []prometheus.CounterOpts{
		{
			Name: HaltCountName,
			Help: "[SYNCHRONIZER] total count of halt",
		},
	}
)

// VirtualBatchNum set the gauge to the given virtual batch num
func VirtualBatchNum(batchNum uint64) {
	metrics.GaugeSet(VirtualBatchNumName, float64(batchNum))
}

// VerifiedBatchNum set the gauge to the given verified batch num
func VerifiedBatchNum(batchNum uint64) {
	metrics.GaugeSet(VerifiedBatchNumName, float64(batchNum))
}

// HaltCount increases the counter for the synchronizer halt count.
func HaltCount() {
	metrics.CounterAdd(HaltCountName, 1)
}
