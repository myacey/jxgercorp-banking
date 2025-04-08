package telemetry

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// StartTracer initialize tracer exporter (to localhost:4318)
func StartTracer(serviceName, serviceVersion string) (*trace.TracerProvider, *metric.MeterProvider, error) {
	headers := map[string]string{
		"content-type": "application/json",
	}

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(serviceName),
		semconv.ServiceVersionKey.String(serviceVersion),
	)

	// trace
	traceExporter, err := otlptrace.New(
		context.Background(),
		otlptracehttp.NewClient(
			otlptracehttp.WithEndpoint("localhost:4318"), // OTLP collector

			// for production:
			// otlptracehttp.WithEndpoint("otel-collector:4318"),

			otlptracehttp.WithHeaders(headers),
			otlptracehttp.WithInsecure(),
		),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create trace exporter: %v", err)
	}
	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter),
		trace.WithResource(res),
	)

	// metric
	metricExporter, err := otlpmetrichttp.New(
		context.Background(),
		otlpmetrichttp.WithEndpoint("localhost:4318"),
		otlpmetrichttp.WithInsecure(),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create metric exporter: %v", err)
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(metricExporter)),
		metric.WithResource(res),
	)

	// set trace provider as a global, so all components of the service can access it
	otel.SetTracerProvider(tracerProvider)

	// set metric provider as global
	otel.SetMeterProvider(meterProvider)

	// set global propagator to tracecontext, it can help us to wire services which calls each others
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tracerProvider, meterProvider, nil
}
