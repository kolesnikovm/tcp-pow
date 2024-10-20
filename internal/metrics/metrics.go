package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	RequestsTotal     prometheus.Counter
	RequestsDuration  prometheus.Histogram
	ClientConnections prometheus.Gauge
)

func init() {
	RequestsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "requests_total",
			Help: "Total number of request",
		},
	)

	RequestsDuration = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name: "requests_duration_ns",
			Help: "Time in nanoseconds that server served a request",
		},
	)

	ClientConnections = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "client_connections",
			Help: "Number of client connections",
		},
	)

	prometheus.MustRegister(
		RequestsTotal,
		RequestsDuration,
		ClientConnections,
	)
}
