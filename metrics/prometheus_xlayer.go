package metrics

import (
	"github.com/0xPolygonHermez/zkevm-node/log"
	"github.com/prometheus/client_golang/prometheus"
)

// GaugeVecOpts holds options for the GaugeVec type.
type GaugeVecOpts struct {
	prometheus.GaugeOpts
	Labels []string
}

// RegisterGaugeVecs registers the provided gauge vec metrics to the Prometheus
// registerer.
func RegisterGaugeVecs(opts ...GaugeVecOpts) {
	if !initialized {
		return
	}

	storageMutex.Lock()
	defer storageMutex.Unlock()

	for _, options := range opts {
		registerGaugeVecIfNotExists(options)
	}
}

// UnregisterGaugeVecs unregisters the provided gauge vec metrics from the Prometheus
// registerer.
func UnregisterGaugeVecs(names ...string) {
	if !initialized {
		return
	}

	storageMutex.Lock()
	defer storageMutex.Unlock()

	for _, name := range names {
		unregisterGaugeVecIfExists(name)
	}
}

// GaugeVec retrieves gauge ver metric by name
func GaugeVec(name string) (gaugeVec *prometheus.GaugeVec, exist bool) {
	if !initialized {
		return
	}

	storageMutex.RLock()
	defer storageMutex.RUnlock()

	gaugeVec, exist = gaugeVecs[name]

	return gaugeVec, exist
}

// GaugeVecSet sets the value for gauge vec with the given name and label.
// name and label.
func GaugeVecSet(name string, label string, value float64) {
	if !initialized {
		return
	}

	if cv, ok := GaugeVec(name); ok {
		cv.WithLabelValues(label).Add(value)
	}
}

// registerGaugeVecIfNotExists registers single gauge vec metric if not exists
func registerGaugeVecIfNotExists(opts GaugeVecOpts) {
	log := log.WithFields("metricName", opts.Name)
	if _, exist := gaugeVecs[opts.Name]; exist {
		log.Warn("Gauge vec metric already exists.")
		return
	}

	log.Debug("Creating Gauge Vec Metric...")
	gauge := prometheus.NewGaugeVec(opts.GaugeOpts, opts.Labels)
	log.Debugf("Gauge Vec Metric successfully created! Labels: %p", opts.ConstLabels)

	log.Debug("Registering Gauge Vec Metric...")
	registerer.MustRegister(gauge)
	log.Debug("Gauge Vec Metric successfully registered!")

	gaugeVecs[opts.Name] = gauge
}

// unregisterGaugeVecIfExists unregisters single gauge vec metric if exists
func unregisterGaugeVecIfExists(name string) {
	var (
		gauge *prometheus.GaugeVec
		ok    bool
	)

	log := log.WithFields("metricName", name)
	if gauge, ok = gaugeVecs[name]; !ok {
		log.Warn("Trying to delete non-existing Gauge metrics.")
		return
	}

	log.Debug("Unregistering Gauge Vec Metric...")
	ok = registerer.Unregister(gauge)
	if !ok {
		log.Error("Failed to unregister Gauge Vec Metric.")
		return
	}
	delete(gauges, name)
	log.Debug("Gauge Vec Metric successfully unregistered!")
}
