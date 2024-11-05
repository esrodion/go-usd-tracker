package traces

import (
	"context"
	"go-usdtrub/pkg/logger"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

// Init returns an instance of Jaeger Tracer.
func Init(ctx context.Context, service string) error {
	client := otlptracegrpc.NewClient(
		otlptracegrpc.WithInsecure(),
	)
	exporter, err := otlptrace.New(ctx, client)
	if err != nil {
		return err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(newResource(service)),
	)

	tracer = tp.Tracer(service)
	return nil
}

func newResource(service string) *resource.Resource {
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(service),
	)
}

func Start(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	if tracer == nil {
		log := logger.Logger().Sugar().Named("Otel")
		log.Error("attempt to access uninitialized tracer")
		err := Init(context.Background(), "test")
		if err != nil {
			log.Panic(err)
		}
	}
	return tracer.Start(ctx, name, opts...)
}
