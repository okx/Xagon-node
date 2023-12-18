package metrics

import (
	"time"

	"github.com/0xPolygonHermez/zkevm-node/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	prefix                    = "jsonrpc_"
	requestPrefix             = prefix + "request_"
	requestsHandledName       = requestPrefix + "handled"
	requestDurationName       = requestPrefix + "duration"
	requestMethodName         = requestPrefix + "method"
	requestMethodDurationName = requestPrefix + "method_duration"

	wsRequestPrefix             = prefix + "ws_request_"
	requestWsMethodName         = wsRequestPrefix + "method"
	requestWsMethodDurationName = wsRequestPrefix + "method_duration"

	requestHandledTypeLabelName = "type"
	requestMethodLabelName      = "method"
)

// RequestHandledLabel represents the possible values for the
// `jsonrpc_request_handled` metric `type` label.
type RequestHandledLabel string

const (
	// RequestHandledLabelInvalid represents an request of type invalid
	RequestHandledLabelInvalid RequestHandledLabel = "invalid"
	// RequestHandledLabelError represents an request of type error
	RequestHandledLabelError RequestHandledLabel = "error"
	// RequestHandledLabelSingle represents an request of type single
	RequestHandledLabelSingle RequestHandledLabel = "single"
	// RequestHandledLabelBatch represents an request of type batch
	RequestHandledLabelBatch RequestHandledLabel = "batch"
)

// Register the metrics for the jsonrpc package.
func Register() {
	var (
		counterVecs   []metrics.CounterVecOpts
		histograms    []prometheus.HistogramOpts
		histogramVecs []metrics.HistogramVecOpts
	)

	counterVecs = []metrics.CounterVecOpts{
		{
			CounterOpts: prometheus.CounterOpts{
				Name: requestsHandledName,
				Help: "[JSONRPC] number of requests handled",
			},
			Labels: []string{requestHandledTypeLabelName},
		},
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
	}

	start := 0.1
	width := 0.1
	count := 10
	histograms = []prometheus.HistogramOpts{
		{
			Name:    requestDurationName,
			Help:    "[JSONRPC] Histogram for the runtime of requests",
			Buckets: prometheus.LinearBuckets(start, width, count),
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

	metrics.RegisterCounterVecs(counterVecs...)
	metrics.RegisterHistograms(histograms...)
	metrics.RegisterHistogramVecs(histogramVecs...)
}

// WsRequestMethodDuration observes (histogram) the duration of a ws request from the
// provided starting time.
func WsRequestMethodDuration(method string, start time.Time) {
	metrics.HistogramVecObserve(requestMethodDurationName, method, time.Since(start).Seconds())
}

// WsRequestMethodCount increments the ws requests handled counter vector by one for
// the given method.
func WsRequestMethodCount(method string) {
	metrics.CounterVecInc(requestMethodName, method)
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

// RequestHandled increments the requests handled counter vector by one for the
// given label.
func RequestHandled(label RequestHandledLabel) {
	metrics.CounterVecInc(requestsHandledName, string(label))
}

// RequestDuration observes (histogram) the duration of a request from the
// provided starting time.
func RequestDuration(start time.Time) {
	metrics.HistogramObserve(requestDurationName, time.Since(start).Seconds())
}
