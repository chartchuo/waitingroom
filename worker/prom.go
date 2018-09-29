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

	requestRateMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "ccwait",
			Subsystem: "worker",
			Name:      "request_rate",
			Help:      "Number of request to target host per second",
		},
		[]string{"host"})

	// responseTimeHistogram = prometheus.NewHistogramVec(
	// 	prometheus.HistogramOpts{
	// 		Namespace: "ccwait",
	// 		Subsystem: "worker",
	// 		Name:      "response_time_histogram",
	// 		Help:      "Response time from target host distribution.",
	// 		Buckets:   []float64{100, 200, 500, 1000, 2000, 5000, 10000},
	// 	},
	// 	[]string{"host"},
	// )
)

func init() {
	prometheus.MustRegister(concurrentUserMetric)
	prometheus.MustRegister(avgResponseTimeMetric)
	prometheus.MustRegister(requestRateMetric)
	// prometheus.MustRegister(responseTimeHistogram)

}
