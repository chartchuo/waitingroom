package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	concurrentUserMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "ccwait",
			Subsystem: "worker",
			Name:      "concurrent_users",
			Help:      "Number of users request target host.",
		},
		[]string{"host"})

	avgResponseTimeMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "ccwait",
			Subsystem: "worker",
			Name:      "avg_response_time",
			Help:      "Average response time request to target host. (microsecond)",
		},
		[]string{"host"})

	maxResponseTimeMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "ccwait",
			Subsystem: "worker",
			Name:      "max_response_time",
			Help:      "Max response time request to target host. (microsecond)",
		},
		[]string{"host"})

	p95ResponseTimeMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "ccwait",
			Subsystem: "worker",
			Name:      "p95_response_time",
			Help:      "95 percentile response time request to target host. (microsecond)",
		},
		[]string{"host"})

	requestRateMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "ccwait",
			Subsystem: "worker",
			Name:      "request_rate",
			Help:      "Number of request to target host per second",
		},
		[]string{"host"})

	avgSessionTimeMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "ccwait",
			Subsystem: "worker",
			Name:      "avg_session_time",
			Help:      "Average session time",
		},
		[]string{"host"})
)

func init() {
	prometheus.MustRegister(concurrentUserMetric)
	prometheus.MustRegister(avgResponseTimeMetric)
	prometheus.MustRegister(maxResponseTimeMetric)
	prometheus.MustRegister(p95ResponseTimeMetric)
	prometheus.MustRegister(requestRateMetric)
	prometheus.MustRegister(avgSessionTimeMetric)
}
