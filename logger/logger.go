package logger

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
)

var (
	log    *zap.Logger
	tracer opentracing.Tracer
)

func Init() error {
	var err error
	log, err = zap.NewProduction()
	if err != nil {
		return err
	}
	return nil
}

func InitTracer(serviceName string) error {
	cfg := jaegercfg.Configuration{
		ServiceName: serviceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
		},
	}

	var err error
	tracer, _, err = cfg.NewTracer()
	if err != nil {
		return err
	}

	opentracing.SetGlobalTracer(tracer)
	return nil
}

func GetLogger() *zap.Logger {
	return log
}

func GetTracer() opentracing.Tracer {
	return tracer
}

func Info(msg string, fields ...zap.Field) {
	log.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	log.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	log.Fatal(msg, fields...)
}

func WithContext(ctx context.Context) *zap.Logger {
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		return log.With(
			zap.String("trace_id", span.Context().(jaeger.SpanContext).TraceID().String()),
			zap.String("span_id", span.Context().(jaeger.SpanContext).SpanID().String()),
		)
	}
	return log
}

func LogWithTracing(ctx context.Context, msg string, fields ...zap.Field) {
	WithContext(ctx).Info(msg, fields...)
}
