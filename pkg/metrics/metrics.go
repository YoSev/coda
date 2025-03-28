package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var gauges = map[string]prometheus.Gauge{
	"coda_runtime_total": prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "coda_runtime_total",
		Help: "Amount of total milliseconds of all executed codas",
	}),
	"coda_total": prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "coda_total",
		Help: "Amount of executed codas",
	}),
	"coda_successful_total": prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "coda_successful_total",
		Help: "Amount of executed successful codas",
	}),
	"coda_failed_total": prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "coda_failed_total",
		Help: "Amount of executed failed codas",
	}),
	"operations_runtime_total": prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "operations_runtime_total",
		Help: "Amount of total milliseconds of all executed operations",
	}),
	"operations_total": prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "operations_total",
		Help: "Amount of executed operations",
	}),
	"operations_successful_total": prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "operations_successful_total",
		Help: "Amount of executed successful operations",
	}),
	"operations_failed_total": prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "operations_failed_total",
		Help: "Amount of executed failed operations",
	}),
	"operations_blacklisted_total": prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "operations_blacklisted_total",
		Help: "Amount of executed blacklisted operations",
	}),
	"variables_total": prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "variables_total",
		Help: "Amount of resolved variables",
	}),
	"variables_failed_total": prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "variables_failed_total",
		Help: "Amount of failed resolved variables",
	}),
	"variables_successful_total": prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "variables_successful_total",
		Help: "Amount of successful resolved variables",
	}),
}
var registry = prometheus.NewRegistry()

func init() {
	for _, gauge := range gauges {
		registry.MustRegister(gauge)
	}
}

func Registry() *prometheus.Registry {
	return registry
}

func Inc(name string) {
	if gauge, ok := gauges[name]; ok {
		gauge.Inc()
	}
}

func Dec(name string) {
	if gauge, ok := gauges[name]; ok {
		gauge.Inc()
	}
}

func IncValue(name string, value float64) {
	if gauge, ok := gauges[name]; ok {
		gauge.Add(value)
	}
}

func DecValue(name string, value float64) {
	if gauge, ok := gauges[name]; ok {
		gauge.Sub(value)
	}
}
