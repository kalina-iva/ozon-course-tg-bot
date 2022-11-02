package tracing

import (
	"io"

	"github.com/pkg/errors"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

var closer io.Closer

func InitTracing(serviceName string, samplingRatio float64) error {
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: samplingRatio,
		},
	}

	var err error
	closer, err = cfg.InitGlobalTracer(serviceName)
	if err != nil {
		return errors.Wrap(err, "cannot init tracing")
	}
	return nil
}

func Close() {
	closer.Close()
}
