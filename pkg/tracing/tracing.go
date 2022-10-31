package tracing

import (
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/pkg/logger"
	"go.uber.org/zap"
)

func InitTracing(service string, param float64) {
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: param,
		},
	}

	_, err := cfg.InitGlobalTracer(service)
	if err != nil {
		logger.Fatal("cannot init tracing", zap.Error(err))
	}
}
