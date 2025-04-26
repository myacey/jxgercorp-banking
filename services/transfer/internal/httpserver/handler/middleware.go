package handler

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.19.0"
	"go.opentelemetry.io/otel/trace"
)

func (h *Handler) TracingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fullpath := c.Request.Method + " " + c.FullPath()
		method := c.Request.Method

		ctx := otel.GetTextMapPropagator().Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))
		ctx, span := h.tracer.Start(ctx, "transfer-service: "+fullpath, trace.WithSpanKind(trace.SpanKindServer))
		defer span.End()

		// start := time.Now()

		// span.SetAttributes(
		// 	attribute.String("http.method"),
		// 	attribute.String("http.route", c.FullPath()),
		// )

		c.Request = c.Request.WithContext(ctx)
		c.Next()

		status := c.Writer.Status()
		attribues := []attribute.KeyValue{
			semconv.HTTPMethod(method),
			semconv.HTTPRoute(fullpath),
			semconv.HTTPStatusCode(status),
			semconv.ServiceName("transfer-service"),
		}
		span.SetAttributes(attribues...)

		// duration := time.Since(start)
		// h.metrics.HTTPMetrics.RecordDuration(ctx, duration.Milliseconds(), attribues...)
		// if !(200 <= c.Writer.Status() && c.Writer.Status() < 300) {
		// 	h.metrics.HTTPMetrics.RecordError(ctx, attribues...)
		// }

		span.SetAttributes(attribute.Int("http.status_code", c.Writer.Status()))
		// h.metrics.HTTPMetrics.RecordHit(ctx, attribues...)
	}
}
