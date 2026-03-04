package telemetry

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type UserMetrics struct {
	usersCreated metric.Int64Counter
}

func NewUserMetrics(m metric.Meter) *UserMetrics {
	created, _ := m.Int64Counter(
		"users_created_total",
		metric.WithDescription("Total registered users"),
	)

	return &UserMetrics{
		usersCreated: created,
	}
}

func (m *UserMetrics) RecordUserCreated(ctx context.Context, attrs ...attribute.KeyValue) {
	m.usersCreated.Add(ctx, 1, metric.WithAttributes(attrs...))
}
