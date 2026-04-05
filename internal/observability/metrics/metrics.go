// Package metrics defines Prometheus metrics for the service.
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

// Metrics holds all application Prometheus metrics.
type Metrics struct {
	GetRatesTotal  prometheus.Counter
	GetRatesErrors prometheus.Counter
	GrinexTotal    prometheus.Counter
	GrinexDuration prometheus.Histogram
}

// New creates and registers application metrics.
func New(reg prometheus.Registerer) *Metrics {
	m := &Metrics{
		GetRatesTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "grinex_rates_get_rates_requests_total",
			Help: "Total number of GetRates requests.",
		}),

		GetRatesErrors: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "grinex_rates_get_rates_errors_total",
			Help: "Total number of failed GetRates requests.",
		}),

		GrinexTotal: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "grinex_rates_grinex_requests_total",
			Help: "Total number of HTTP requests to Grinex depth API.",
		}),

		GrinexDuration: prometheus.NewHistogram(prometheus.HistogramOpts{
			Name:    "grinex_rates_grinex_request_duration_seconds",
			Help:    "Duration of HTTP requests to Grinex depth API.",
			Buckets: prometheus.DefBuckets,
		}),
	}

	reg.MustRegister(
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		collectors.NewGoCollector(),
		m.GetRatesTotal,
		m.GetRatesErrors,
		m.GrinexTotal,
		m.GrinexDuration,
	)

	return m
}
