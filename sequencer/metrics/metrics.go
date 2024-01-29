package metrics

import (
	"time"

	"github.com/0xPolygonHermez/zkevm-node/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	// Prefix for the metrics of the sequencer package.
	Prefix = "sequencer_"
	// SequencesSentToL1CountName is the name of the metric that counts the sequences sent to L1.
	SequencesSentToL1CountName = Prefix + "sequences_sent_to_L1_count"
	// GasPriceEstimatedAverageName is the name of the metric that shows the average estimated gas price.
	GasPriceEstimatedAverageName = Prefix + "gas_price_estimated_average"
	// TxProcessedName is the name of the metric that counts the processed transactions.
	TxProcessedName = Prefix + "transaction_processed"
	// SequencesOversizedDataErrorName is the name of the metric that counts the sequences with oversized data error.
	SequencesOversizedDataErrorName = Prefix + "sequences_oversized_data_error"
	// EthToMaticPriceName is the name of the metric that shows the Ethereum to Matic price.
	EthToMaticPriceName = Prefix + "eth_to_matic_price"
	// SequenceRewardInMaticName is the name of the metric that shows the reward in Matic of a sequence.
	SequenceRewardInMaticName = Prefix + "sequence_reward_in_matic"
	// ProcessingTimeName is the name of the metric that shows the processing time.
	ProcessingTimeName = Prefix + "processing_time"
	// PendingTxCountName is the name of metric that shows the number of pending transactions.
	PendingTxCountName = Prefix + "pending_tx_count"
	// BatchExecuteTimeName is the name of the metric that shows the batch execution time.
	BatchExecuteTimeName = Prefix + "batch_execute_time"
	// TrustBatchNumName is the name of the metric that shows the trust batch num
	TrustBatchNumName = Prefix + "trust_batch_num"
	// WorkerPrefix is the prefix for the metrics of the worker.
	WorkerPrefix = Prefix + "worker_"
	// WorkerProcessingTimeName is the name of the metric that shows the worker processing time.
	WorkerProcessingTimeName = WorkerPrefix + "processing_time"
	// TxProcessedLabelName is the name of the label for the processed transactions.
	TxProcessedLabelName = "status"
	// BatchFinalizeTypeLabelName is the name of the label for the batch finalize type.
	BatchFinalizeTypeLabelName = "batch_type"
	// HaltCountName is the name of the metric that counts the halt count.
	HaltCountName = Prefix + "halt_count"
)

// TxProcessedLabel represents the possible values for the
// `sequencer_transaction_processed` metric `type` label.
type TxProcessedLabel string

const (
	// TxProcessedLabelSuccessful represents a successful transaction
	TxProcessedLabelSuccessful TxProcessedLabel = "successful"
	// TxProcessedLabelInvalid represents an invalid transaction
	TxProcessedLabelInvalid TxProcessedLabel = "invalid"
	// TxProcessedLabelFailed represents a failed transaction
	TxProcessedLabelFailed TxProcessedLabel = "failed"
)

// BatchFinalizeTypeLabel batch finalize type label
type BatchFinalizeTypeLabel string

const (
	// BatchFinalizeTypeLabelDeadline batch finalize type deadline label
	BatchFinalizeTypeLabelDeadline BatchFinalizeTypeLabel = "deadline"
	// BatchFinalizeTypeLabelFullBatch batch finalize type full batch label
	BatchFinalizeTypeLabelFullBatch BatchFinalizeTypeLabel = "full_batch"
)

// Register the metrics for the sequencer package.
func Register() {
	var (
		counters    []prometheus.CounterOpts
		counterVecs []metrics.CounterVecOpts
		gauges      []prometheus.GaugeOpts
		gaugeVecs   []metrics.GaugeVecOpts
		histograms  []prometheus.HistogramOpts
	)

	counters = []prometheus.CounterOpts{
		{
			Name: SequencesSentToL1CountName,
			Help: "[SEQUENCER] total count of sequences sent to L1",
		},
		{
			Name: SequencesOversizedDataErrorName,
			Help: "[SEQUENCER] total count of sequences with oversized data error",
		},
		{
			Name: HaltCountName,
			Help: "[SEQUENCER] total count of halt",
		},
	}

	counterVecs = []metrics.CounterVecOpts{
		{
			CounterOpts: prometheus.CounterOpts{
				Name: TxProcessedName,
				Help: "[SEQUENCER] number of transactions processed",
			},
			Labels: []string{TxProcessedLabelName},
		},
	}

	gauges = []prometheus.GaugeOpts{
		{
			Name: GasPriceEstimatedAverageName,
			Help: "[SEQUENCER] average gas price estimated",
		},
		{
			Name: EthToMaticPriceName,
			Help: "[SEQUENCER] eth to matic price",
		},
		{
			Name: SequenceRewardInMaticName,
			Help: "[SEQUENCER] reward for a sequence in Matic",
		},
		{
			Name: PendingTxCountName,
			Help: "[SEQUENCER] number of pending transactions",
		},
		{
			Name: TrustBatchNumName,
			Help: "[SEQUENCER] trust batch num",
		},
	}

	gaugeVecs = []metrics.GaugeVecOpts{
		{
			GaugeOpts: prometheus.GaugeOpts{
				Name: BatchExecuteTimeName,
				Help: "[SEQUENCER] batch execution time",
			},
			Labels: []string{BatchFinalizeTypeLabelName},
		},
	}

	histograms = []prometheus.HistogramOpts{
		{
			Name: ProcessingTimeName,
			Help: "[SEQUENCER] processing time",
		},
		{
			Name: WorkerProcessingTimeName,
			Help: "[SEQUENCER] worker processing time",
		},
	}

	metrics.RegisterCounters(counters...)
	metrics.RegisterCounterVecs(counterVecs...)
	metrics.RegisterGauges(gauges...)
	metrics.RegisterGaugeVecs(gaugeVecs...)
	metrics.RegisterHistograms(histograms...)
}

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

// AverageGasPrice sets the gauge to the given average gas price.
func AverageGasPrice(price float64) {
	metrics.GaugeSet(GasPriceEstimatedAverageName, price)
}

// SequencesSentToL1 increases the counter by the provided number of sequences
// sent to L1.
func SequencesSentToL1(numSequences float64) {
	metrics.CounterAdd(SequencesSentToL1CountName, numSequences)
}

// TxProcessed increases the counter vector by the provided transactions count
// and for the given label (status).
func TxProcessed(status TxProcessedLabel, count float64) {
	metrics.CounterVecAdd(TxProcessedName, string(status), count)
}

// SequencesOvesizedDataError increases the counter for sequences that
// encounter a OversizedData error.
func SequencesOvesizedDataError() {
	metrics.CounterInc(SequencesOversizedDataErrorName)
}

// EthToMaticPrice sets the gauge for the Ethereum to Matic price.
func EthToMaticPrice(price float64) {
	metrics.GaugeSet(EthToMaticPriceName, price)
}

// SequenceRewardInMatic sets the gauge for the reward in Matic of a sequence.
func SequenceRewardInMatic(reward float64) {
	metrics.GaugeSet(SequenceRewardInMaticName, reward)
}

// ProcessingTime observes the last processing time on the histogram.
func ProcessingTime(lastProcessTime time.Duration) {
	execTimeInSeconds := float64(lastProcessTime) / float64(time.Second)
	metrics.HistogramObserve(ProcessingTimeName, execTimeInSeconds)
}

// WorkerProcessingTime observes the last processing time on the histogram.
func WorkerProcessingTime(lastProcessTime time.Duration) {
	execTimeInSeconds := float64(lastProcessTime) / float64(time.Second)
	metrics.HistogramObserve(WorkerProcessingTimeName, execTimeInSeconds)
}

// HaltCount increases the counter for the sequencer halt count.
func HaltCount() {
	metrics.CounterAdd(HaltCountName, 1)
}
