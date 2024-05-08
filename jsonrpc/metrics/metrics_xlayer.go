package metrics

import (
	"time"

	"github.com/0xPolygonHermez/zkevm-node/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	requestMethodName           = requestPrefix + "method"
	requestMethodDurationName   = requestPrefix + "method_duration"
	requestInnerTxCachedName    = requestPrefix + "inner_tx_cached"
	requestInnerTxExecutedName  = requestPrefix + "inner_tx_executed"
	requestInnerTxAddErrorCount = requestPrefix + "inner_tx_error_count"

	wsRequestPrefix             = prefix + "ws_request_"
	requestWsMethodName         = wsRequestPrefix + "method"
	requestWsMethodDurationName = wsRequestPrefix + "method_duration"
	requestMethodLabelName      = "method"

	start = 0.1
	width = 0.1
	count = 10

	lastDynamicGasPriceName = prefix + "dynamic_gas_price"
	lastRawGasPriceName     = prefix + "raw_gas_price"

	gauges = []prometheus.GaugeOpts{
		{
			Name: lastDynamicGasPriceName,
			Help: "[JSONRPC] dynamic gas price",
		},
		{
			Name: lastRawGasPriceName,
			Help: "[JSONRPC] raw gas price",
		},
	}

	histogramVecs = []metrics.HistogramVecOpts{
		{
			HistogramOpts: prometheus.HistogramOpts{
				Name:    requestMethodDurationName,
				Help:    "[JSONRPC] Histogram for the runtime of requests",
				Buckets: prometheus.LinearBuckets(start, width, count),
			},
			Labels: []string{requestMethodLabelName},
		},
		{
			HistogramOpts: prometheus.HistogramOpts{
				Name:    requestWsMethodDurationName,
				Help:    "[JSONRPC] Histogram for the runtime of ws requests",
				Buckets: prometheus.LinearBuckets(start, width, count),
			},
			Labels: []string{requestMethodLabelName},
		},
	}
	counterVecsXLayer = []metrics.CounterVecOpts{
		{
			CounterOpts: prometheus.CounterOpts{
				Name: requestMethodName,
				Help: "[JSONRPC] number of requests handled by method",
			},
			Labels: []string{requestMethodLabelName},
		},
		{
			CounterOpts: prometheus.CounterOpts{
				Name: requestWsMethodName,
				Help: "[JSONRPC] number of ws requests handled by method",
			},
			Labels: []string{requestMethodLabelName},
		},
		{
			CounterOpts: prometheus.CounterOpts{
				Name: requestInnerTxCachedName,
				Help: "[JSONRPC] number of cached innertx requests",
			},
			Labels: []string{"type"},
		},
		{
			CounterOpts: prometheus.CounterOpts{
				Name: requestInnerTxExecutedName,
				Help: "[JSONRPC] number of executed innertx requests",
			},
			Labels: []string{"type"},
		},
		{
			CounterOpts: prometheus.CounterOpts{
				Name: requestInnerTxAddErrorCount,
				Help: "[JSONRPC] number of add innertx count",
			},
			Labels: []string{"type"},
		},
	}
)

// WsRequestMethodDuration observes (histogram) the duration of a ws request from the
// provided starting time.
func WsRequestMethodDuration(method string, start time.Time) {
	metrics.HistogramVecObserve(requestWsMethodDurationName, method, time.Since(start).Seconds())
}

// WsRequestMethodCount increments the ws requests handled counter vector by one for
// the given method.
func WsRequestMethodCount(method string) {
	metrics.CounterVecInc(requestWsMethodName, method)
}

// RequestMethodDuration observes (histogram) the duration of a request from the
// provided starting time.
func RequestMethodDuration(method string, start time.Time) {
	metrics.HistogramVecObserve(requestMethodDurationName, method, time.Since(start).Seconds())
}

// RequestMethodCount increments the requests handled counter vector by one for
// the given method.
func RequestMethodCount(method string) {
	metrics.CounterVecInc(requestMethodName, method)
}

// RequestInnerTxExecutedCount increments the inner tx executed counter vector by one.
func RequestInnerTxExecutedCount() {
	metrics.CounterVecInc(requestInnerTxExecutedName, "executed")
}

// RequestInnerTxCachedCount increments the inner tx cached counter vector by one.
func RequestInnerTxCachedCount() {
	metrics.CounterVecInc(requestInnerTxCachedName, "cached")
}

// RequestInnerTxAddErrorCount increments the inner tx add error counter vector by one.
func RequestInnerTxAddErrorCount() {
	metrics.CounterVecInc(requestInnerTxAddErrorCount, "add_error")
}

// DynamicGasPrice sets the gauge vector to the given batch number and dynamic gas price.
func DynamicGasPrice(dgp int64) {
	metrics.GaugeSet(lastDynamicGasPriceName, float64(dgp))
}

// RawGasPrice sets the gauge vector to the given batch number and raw gas price.
func RawGasPrice(gp int64) {
	metrics.GaugeSet(lastRawGasPriceName, float64(gp))
}
