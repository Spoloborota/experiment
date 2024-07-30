package tracing

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

func InitTracer() (*trace.TracerProvider, error) {
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, err
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewSchemaless(
			semconv.ServiceNameKey.String("social_network_service"),
		)),
	)
	otel.SetTracerProvider(tp)
	return tp, nil
}

func ShutdownTracer(tp *trace.TracerProvider) {
	if err := tp.Shutdown(context.Background()); err != nil {
		log.Fatalf("Error shutting down tracer provider: %v", err)
	}
}
