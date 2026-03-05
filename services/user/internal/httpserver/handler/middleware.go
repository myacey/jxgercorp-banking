package handler

import (
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"go.opentelemetry.io/otel/trace"
)

func (h *Handler) TracingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fullpath := c.Request.Method + " " + c.FullPath()
		method := c.Request.Method

		ctx := otel.GetTextMapPropagator().Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))
		ctx, span := h.tracer.Start(ctx, "user-service: "+fullpath, trace.WithSpanKind(trace.SpanKindServer))
		defer span.End()

		c.Request = c.Request.WithContext(ctx)
		c.Next()

		status := c.Writer.Status()
		attribues := []attribute.KeyValue{
			semconv.HTTPMethod(method),
			semconv.HTTPRoute(fullpath),
			semconv.HTTPStatusCode(status),
			semconv.ServiceName("user-service"),
		}
		span.SetAttributes(attribues...)

		span.SetAttributes(attribute.Int("http.status_code", c.Writer.Status()))
	}
}

func (h *Handler) MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		ctx := c.Request.Context()

		c.Next()

		duration := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		route := c.FullPath()
		if route == "" {
			route = "unknown"
		}

		attrs := []attribute.KeyValue{
			semconv.HTTPMethod(method),
			semconv.HTTPRoute(route),
			semconv.HTTPStatusCode(status),
			semconv.ServiceName("user-service"),
		}

		h.metrics.RecordHit(ctx, attrs...)
		h.metrics.RecordDuration(ctx, duration.Milliseconds(), attrs...)

		if status >= 500 {
			errData, _ := io.ReadAll(c.Request.Response.Body)
			attrs = append(attrs, attribute.String("error", string(errData)))

			h.metrics.RecordError(ctx, attrs...)
		}
	}
}
