// Copyright 2015 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// A simple example exposing fictional RPC latencies with different types of
// random distributions (uniform, normal, and exponential) as Prometheus
// metrics.
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

	responseTimeHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "ccwait",
			Subsystem: "worker",
			Name:      "response_time_histogram",
			Help:      "Response time from target host distribution.",
			Buckets:   []float64{100, 200, 500, 1000, 2000, 5000, 10000},
		},
		[]string{"host"},
	)
)

func init() {
	prometheus.MustRegister(concurrentUserMetric)
	prometheus.MustRegister(avgResponseTimeMetric)
	prometheus.MustRegister(requestRateMetric)
	prometheus.MustRegister(responseTimeHistogram)

}
