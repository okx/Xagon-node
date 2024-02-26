package metrics

import (
	"github.com/0xPolygonHermez/zkevm-node/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// PendingTxCountName is the name of metric that shows the number of pending transactions.
	PendingTxCountName = Prefix + "pending_tx_count"
	// BatchExecuteTimeName is the name of the metric that shows the batch execution time.
	BatchExecuteTimeName = Prefix + "batch_execute_time"
	// TrustBatchNumName is the name of the metric that shows the trust batch num
	TrustBatchNumName = Prefix + "trust_batch_num"
	// BatchFinalizeTypeLabelName is the name of the label for the batch finalize type.
	BatchFinalizeTypeLabelName = "batch_type"
	// HaltCountName is the name of the metric that counts the halt count.
	HaltCountName = Prefix + "halt_count"

	gaugeVecs = []metrics.GaugeVecOpts{
		{
			GaugeOpts: prometheus.GaugeOpts{
				Name: BatchExecuteTimeName,
				Help: "[SEQUENCER] batch execution time",
			},
			Labels: []string{BatchFinalizeTypeLabelName},
		},
	}

	gaugesX1 = []prometheus.GaugeOpts{
		{
			Name: PendingTxCountName,
			Help: "[SEQUENCER] number of pending transactions",
		},
		{
			Name: TrustBatchNumName,
			Help: "[SEQUENCER] trust batch num",
		},
	}
	countersX1 = []prometheus.CounterOpts{
		{
			Name: HaltCountName,
			Help: "[SEQUENCER] total count of halt",
		},
	}
)

// BatchFinalizeTypeLabel batch finalize type label
type BatchFinalizeTypeLabel string

const (
	// BatchFinalizeTypeLabelDeadline batch finalize type deadline label
	BatchFinalizeTypeLabelDeadline BatchFinalizeTypeLabel = "deadline"
	// BatchFinalizeTypeLabelFullBatch batch finalize type full batch label
	BatchFinalizeTypeLabelFullBatch BatchFinalizeTypeLabel = "full_batch"
)

// PendingTxCount sets the gauge to the given number of pending transactions.
func PendingTxCount(count int) {
	metrics.GaugeSet(PendingTxCountName, float64(count))
}

// BatchExecuteTime sets the gauge vector to the given batch type and time.
func BatchExecuteTime(batchType BatchFinalizeTypeLabel, time int64) {
	metrics.GaugeVecSet(BatchExecuteTimeName, string(batchType), float64(time))
}

// TrustBatchNum set the gauge to the given trust batch num
func TrustBatchNum(batchNum uint64) {
	metrics.GaugeSet(TrustBatchNumName, float64(batchNum))
}

// HaltCount increases the counter for the sequencer halt count.
func HaltCount() {
	metrics.CounterAdd(HaltCountName, 1)
}
