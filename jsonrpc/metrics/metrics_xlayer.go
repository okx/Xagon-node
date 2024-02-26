package metrics

import (
	"time"

	"github.com/0xPolygonHermez/zkevm-node/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	requestMethodName         = requestPrefix + "method"
	requestMethodDurationName = requestPrefix + "method_duration"

	wsRequestPrefix             = prefix + "ws_request_"
	requestWsMethodName         = wsRequestPrefix + "method"
	requestWsMethodDurationName = wsRequestPrefix + "method_duration"
	requestMethodLabelName      = "method"

	start         = 0.1
	width         = 0.1
	count         = 10
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
	counterVecsX1 = []metrics.CounterVecOpts{
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
)

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
