package telemetry

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

func NewMeter(serviceName string) metric.Meter {
	return otel.Meter(serviceName)
}

type HTTPMetrics struct {
	requestsTotal   metric.Int64Counter
	requestDuration metric.Int64Histogram
	errorTotal      metric.Int64Counter
}

// Common HTTP metrics for all services
func NewHTTPMetrics(m metric.Meter) *HTTPMetrics {
	reqCounter, _ := m.Int64Counter(
		"http_requests_total",
		metric.WithDescription("Total HTTP requests"),
	)

	reqDuration, _ := m.Int64Histogram(
		"http_server_duration_ms",
		metric.WithDescription("HTTP request duration in ms"),
		metric.WithUnit("ms"),
	)

	errCounter, _ := m.Int64Counter(
		"http_request_errors_total",
		metric.WithDescription("Total HTTP errors"),
	)

	return &HTTPMetrics{
		requestsTotal:   reqCounter,
		requestDuration: reqDuration,
		errorTotal:      errCounter,
	}
}

// RecordHit adds 1 to RequestCounter
func (m *HTTPMetrics) RecordHit(ctx context.Context, attrs ...attribute.KeyValue) {
	// m.requestsTotal.Add(ctx, 1, metric.WithAttributes(attrs...))
}

// RecordDuration adds duration info to metric
func (m *HTTPMetrics) RecordDuration(ctx context.Context, duration int64, attrs ...attribute.KeyValue) {
	// m.requestDuration.Record(ctx, duration, metric.WithAttributes(attrs...))
}

func (m *HTTPMetrics) RecordError(ctx context.Context, attrs ...attribute.KeyValue) {
	// m.errorTotal.Add(ctx, 1, metric.WithAttributes(attrs...))
}
