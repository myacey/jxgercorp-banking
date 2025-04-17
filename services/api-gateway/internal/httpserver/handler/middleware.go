package handler

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/myacey/jxgercorp-banking/services/libs/apperror"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.18.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func (h *Handler) AuthTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := h.tracer.Start(c.Request.Context(), "auth-middleware: AuthTokenMiddleware")
		c.Request = c.Request.WithContext(ctx)

		if c.Request.Method == http.MethodOptions {
			span.End()
			c.Next()
			return
		}
		authToken, err := c.Cookie("authToken")
		if err != nil {
			wrapCtxWithError(c, apperror.NewUnauthorized("invalid token"))
			span.End()
			return
		}

		usrname, valid, err := h.srv.Auth.ValidateToken(ctx, authToken)
		if err != nil || !valid {
			wrapCtxWithError(c, err)
			span.End()
			return
		}

		c.Request.Header.Set(HeaderUsername, usrname)

		span.End()

		c.Next()
	}
}

func (h *Handler) TracingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// if there is a context inside the carrier we have to extract it and update our local context,
		// so it will add a parent span to the first service span if the service was called form the another service,
		// this is not relevant to the example, but it's a good practice to always extract carrier in the first span,
		// because it makes service wiring extremely easy in future
		ctx := otel.GetTextMapPropagator().Extract(c.Request.Context(), propagation.HeaderCarrier(c.Request.Header))

		// start a new span with the context from carrier
		ctx, span := h.tracer.Start(ctx, "tracing-middleware: "+c.Request.Method+" "+c.FullPath())
		defer span.End()

		span.SetAttributes(
			attribute.String("http.method", c.Request.Method),
			attribute.String("http.url", c.Request.RequestURI),
		)

		c.Request = c.Request.WithContext(ctx)
		c.Request.Header.Set(HeaderRequestID, span.SpanContext().TraceID().String())
		c.Next()
	}
}

// MetricsMiddleware provides request ID and adds metrics
// to prometheus
func (h *Handler) MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// create requestID
		requestID := uuid.New()
		var byteArray [16]byte
		copy(byteArray[:], requestID[:])
		ctx := trace.ContextWithRemoteSpanContext(
			c.Request.Context(),
			trace.NewSpanContext(trace.SpanContextConfig{
				TraceID: byteArray,
			}),
		)

		ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(c.Request.Header))

		// start a new span with the context from carrier
		ctx, span := h.tracer.Start(ctx, "metrics-middleware: "+c.Request.Method+" "+c.FullPath())
		defer span.End()
		c.Request = c.Request.WithContext(ctx)

		// start := time.Now()

		// h.metrics.ActiveRequestsGauge.Add(c.Request.Context(), 1)
		// defer h.metrics.ActiveRequestsGauge.Add(c.Request.Context(), -1)

		c.Next()

		// duration := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.FullPath()

		attribues := []attribute.KeyValue{
			semconv.HTTPMethod(method),
			semconv.HTTPRoute(path),
			semconv.HTTPStatusCode(status),
			semconv.ServiceName("api-gateway"),
		}

		// h.metrics.RequestCounter.Add(c.Request.Context(), 1, metric.WithAttributes(attribues...))
		// h.metrics.DurationHistogram.Record(c.Request.Context(), float64(duration.Microseconds()), metric.WithAttributes(attribues...))
		// if status >= 500 {
		// 	h.metrics.ErrorCounter.Add(c.Request.Context(), 1, metric.WithAttributes(attribues...))
		// }

		span.SetAttributes(attribues...)
	}
}

func UnaryClientInterceptor(ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption,
) error {
	// get original ctx and make it gRPC's metadata
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	} else {
		md = md.Copy()
	}

	// map carrier for injecting
	carrier := propagation.MapCarrier{}
	for k, vals := range md {
		if len(vals) > 0 {
			carrier[k] = vals[0]
		}
	}

	// injet trace data to carrier
	otel.GetTextMapPropagator().Inject(ctx, carrier)

	// move carrier data to md
	for k, v := range carrier {
		md.Set(k, v)
	}

	newCtx := metadata.NewOutgoingContext(ctx, md)
	return invoker(newCtx, method, req, reply, cc, opts...)
}
