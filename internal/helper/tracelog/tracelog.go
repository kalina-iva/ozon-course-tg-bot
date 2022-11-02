package tracelog

import (
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"go.uber.org/zap"
)

func Info(span opentracing.Span, msg string) {
	if spanContext, ok := span.Context().(jaeger.SpanContext); ok {
		zap.L().Info(msg, zap.String("id", spanContext.TraceID().String()))
	}
}
