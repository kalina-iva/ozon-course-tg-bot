package tracing

import (
	"io"

	"github.com/pkg/errors"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

var closer io.Closer

func InitTracing(service string, param float64) error {
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: param,
		},
	}

	var err error
	closer, err = cfg.InitGlobalTracer(service)
	if err != nil {
		return errors.Wrap(err, "cannot init tracing")
	}
	return nil
}

func Close() {
	closer.Close()
}
