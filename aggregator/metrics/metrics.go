package metrics

import (
	"github.com/0xPolygonHermez/zkevm-node/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	prefix                      = "aggregator_"
	currentConnectedProversName = prefix + "current_connected_provers"
	currentWorkingProversName   = prefix + "current_working_provers"
	sendTxFailedCount = prefix + "send_tx_failed_count"
	waitTxToBeMinedFailedCount = prefix + "wait_tx_mined_failed_count"
)

// Register the metrics for the sequencer package.
func Register() {
	gauges := []prometheus.GaugeOpts{
		{
			Name: currentConnectedProversName,
			Help: "[AGGREGATOR] current connected provers",
		},
		{
			Name: currentWorkingProversName,
			Help: "[AGGREGATOR] current working provers",
		},
		{
			Name: sendTxFailedCount,
			Help: "[AGGREGATOR] agglayerClient.SendTx failed counter",
		},
		{
			Name: waitTxToBeMinedFailedCount,
			Help: "[AGGREGATOR] agglayerClient.WaitTxToBeMined failed counter",
		},
	}

	metrics.RegisterGauges(gauges...)
}

// ConnectedProver increments the gauge for the current number of connected
// provers.
func ConnectedProver() {
	metrics.GaugeInc(currentConnectedProversName)
}

// DisconnectedProver decrements the gauge for the current number of connected
// provers.
func DisconnectedProver() {
	metrics.GaugeDec(currentConnectedProversName)
}

// WorkingProver increments the gauge for the current number of working
// provers.
func WorkingProver() {
	metrics.GaugeInc(currentWorkingProversName)
}

// IdlingProver decrements the gauge for the current number of working provers.
func IdlingProver() {
	metrics.GaugeDec(currentWorkingProversName)
}

// SendTxFailedInc increment the gauge for sendTx Failed Counter
func SendTxFailedInc() {
	metrics.GaugeInc(sendTxFailedCount)
}

// SendTxFailedReset when a proof is successfully sent to agglayer
func SendTxFailedReset() {
	metrics.GaugeSet(sendTxFailedCount, 0)
}

// WaitMinedFailedInc increment the gauge for sendTx Failed Counter
func WaitMinedFailedInc() {
	metrics.GaugeInc(waitTxToBeMinedFailedCount)
}

// WaitMinedFailedReset when a proof is successfully sent to agglayer
func WaitMinedFailedReset() {
	metrics.GaugeSet(waitTxToBeMinedFailedCount, 0)
}

