package controller

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func (h *Controller) TracingMiddleware(serviceName string) gin.HandlerFunc {
	tracer := otel.Tracer(serviceName)

	return func(c *gin.Context) {
		ctx, span := tracer.Start(c.Request.Context(), c.FullPath())
		// Добавляем базовые атрибуты
		span.SetAttributes(attribute.String("http.method", c.Request.Method))
		c.Request = c.Request.WithContext(ctx)

		start := time.Now()
		c.Next()
		duration := time.Since(start)
		span.SetAttributes(attribute.Float64("http.duration_ms", float64(duration.Milliseconds())))
		span.End()
	}
}
