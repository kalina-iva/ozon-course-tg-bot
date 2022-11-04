package metrics

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var server *http.Server

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
