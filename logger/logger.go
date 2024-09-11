package logger

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
)

var log *zap.Logger

func Init() {
	var err error
	log, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
}

func GetLogger() *zap.Logger {
	return log
}

func InitTracer(serviceName string) (opentracing.Tracer, error) {
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

	tracer, _, err := cfg.NewTracer()
	if err != nil {
		return nil, err
	}

	opentracing.SetGlobalTracer(tracer)
	return tracer, nil
}

func LogWithTracing(ctx context.Context, msg string, fields ...zap.Field) {
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		fields = append(fields, zap.String("trace_id", span.Context().(jaeger.SpanContext).TraceID().String()))
		fields = append(fields, zap.String("span_id", span.Context().(jaeger.SpanContext).SpanID().String()))
	}
	log.Info(msg, fields...)
}
