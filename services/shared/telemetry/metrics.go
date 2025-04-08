package telemetry

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type MetricsFactory struct {
	meter metric.Meter
}

func NewMetricsFactory(serviceName string) *MetricsFactory {
	return &MetricsFactory{
		meter: otel.Meter(serviceName),
	}
}

type HTTPMetrics struct {
	requestCounter       metric.Int64Counter
	requestDuration      metric.Int64Histogram
	errorCounter         metric.Int64Counter
	activeRequestCounter metric.Int64UpDownCounter
}

// Common HTTP metrics for all services
func (mf *MetricsFactory) NewHTTPMetrics() *HTTPMetrics {
	reqCounter, _ := mf.meter.Int64Counter(
		"http.server.request.count",
		metric.WithDescription("Total HTTP requests"),
	)

	reqDuration, _ := mf.meter.Int64Histogram(
		"http.server.duration",
		metric.WithDescription("HTTP request duration in ms"),
		metric.WithUnit("ms"),
	)

	errCounter, _ := mf.meter.Int64Counter(
		"http.server.errors",
		metric.WithDescription("Total HTTP errors"),
	)

	activeReqCounter, _ := mf.meter.Int64UpDownCounter(
		"http.server.active_requests",
		metric.WithDescription("Number of active HTTP requests"),
	)

	return &HTTPMetrics{
		requestCounter:       reqCounter,
		requestDuration:      reqDuration,
		errorCounter:         errCounter,
		activeRequestCounter: activeReqCounter,
	}
}

// RecordHit adds 1 to RequestCounter and ActiveRequestCounter
func (m *HTTPMetrics) RecordHit(ctx context.Context, attrs ...attribute.KeyValue) {
	m.requestCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
	m.activeRequestCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
}

// RecordDuration adds duration info to metric, also decreasing ActiveRequestCounter
func (m *HTTPMetrics) RecordDuration(ctx context.Context, duration int64, attrs ...attribute.KeyValue) {
	m.requestDuration.Record(ctx, duration, metric.WithAttributes(attrs...))
}

func (m *HTTPMetrics) RecordError(ctx context.Context, attrs ...attribute.KeyValue) {
	m.errorCounter.Add(ctx, 1, metric.WithAttributes(attrs...))
}

// Сustom metrics for user service
type UserMetrics struct {
	HTTPMetrics
	registrations metric.Int64Counter

	logins metric.Int64Counter
}

func (mf *MetricsFactory) NewUserMetrics() *UserMetrics {
	regCounter, _ := mf.meter.Int64Counter(
		"user.register.count",
		metric.WithDescription("Number of successful registrations"),
	)

	loginCounter, _ := mf.meter.Int64Counter(
		"user.login.counter",
		metric.WithDescription("Number of successful logins"),
	)

	httpMetrics := mf.NewHTTPMetrics()
	return &UserMetrics{
		HTTPMetrics:   *httpMetrics,
		registrations: regCounter,
		logins:        loginCounter,
	}
}

func (m *UserMetrics) RecordRegister(ctx context.Context, attrs ...attribute.KeyValue) {
	m.registrations.Add(ctx, 1, metric.WithAttributes(attrs...))
}

func (m *UserMetrics) RecordLogin(ctx context.Context, attr ...attribute.KeyValue) {
	m.logins.Add(ctx, 1, metric.WithAttributes(attr...))
}
