package metrics

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	minResponseTime = 0.0001
	maxResponseTime = 2
	cntBuckets      = 16
)

var server *http.Server

var HistogramResponseTime = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Namespace: "ozon",
		Name:      "histogram_msg_processing_time_sec",
		Buckets:   prometheus.ExponentialBucketsRange(minResponseTime, maxResponseTime, cntBuckets),
	},
	[]string{"command"},
)

func InitMetrics(serverAddress string) error {
	server = &http.Server{Addr: serverAddress}
	http.Handle("/metrics", promhttp.Handler())
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return errors.Wrap(err, "metrics server failed")
	}
	return nil
}

func Close() error {
	return server.Shutdown(context.Background())
}
