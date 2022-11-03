package tracelog

import (
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"gitlab.ozon.dev/mary.kalina/telegram-bot/pkg/logger"
	"go.uber.org/zap"
)

func Info(span opentracing.Span, msg string) {
	if spanContext, ok := span.Context().(jaeger.SpanContext); ok {
		logger.Info(msg, zap.String("id", spanContext.TraceID().String()))
	}
}
