package tg

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	minResponseTime = 0.0001
	maxResponseTime = 2
	cntBuckets      = 16
)

var histogramResponseTime = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Namespace: "ozon",
		Name:      "histogram_msg_processing_time_sec",
		Buckets:   prometheus.ExponentialBucketsRange(minResponseTime, maxResponseTime, cntBuckets),
	},
	[]string{"command"},
)
