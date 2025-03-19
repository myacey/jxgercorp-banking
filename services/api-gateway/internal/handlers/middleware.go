package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/myacey/jxgercorp-banking/services/shared/cstmerr"
	"github.com/myacey/jxgercorp-banking/services/shared/ctxkeys"
	tokenpb "github.com/myacey/jxgercorp-banking/services/shared/proto/token"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func (h *Handler) AuthTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, span := h.tracer.Start(c.Request.Context(), "auth-middleware: AuthTokenMiddleware")
		c.Request = c.Request.WithContext(ctx)

		h.lg.Info(c.Request.Method)
		if c.Request.Method == http.MethodOptions {
			span.End()
			c.Next()
			return
		}
		authToken, err := c.Cookie("authToken")
		if err != nil {
			h.JSONError(c, cstmerr.ErrInvalidToken, http.StatusUnauthorized)
			span.End()
			return
		}

		req := &tokenpb.ValidateTokenRequest{
			Token: authToken,
		}
		resp, err := h.tokenSrv.ValidateToken(ctx, req)
		if err != nil {
			h.JSONError(c, cstmerr.ErrInvalidToken, http.StatusUnauthorized)
			span.End()
			return
		}

		username, valid := resp.Username, resp.Valid
		if !valid || username == "" {
			h.JSONError(c, cstmerr.ErrInvalidToken, http.StatusUnauthorized)
			span.End()
			return
		}

		h.lg.Debugw("auth success",
			"authToken", authToken,
			"username", username)
		c.Set(string(ctxkeys.UsernameKey), username)

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

		c.Next()
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

// TODO
// func RequestIDMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 	}
// }
