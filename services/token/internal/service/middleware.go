package service

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func (srv *TokenService) TracingMiddleware(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	srv.lg.Debug("AAAA")

	// get metadata from incoming context
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}

	// copy from metadata into MapCarrier map[string]string (from map[string][]string)
	carrier := propagation.MapCarrier{}
	for k, vals := range md {
		if len(vals) > 0 {
			carrier[k] = vals[0]
		}
	}

	// get trace context, inject to ctx
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

	ctx, span := srv.tracer.Start(ctx, "middleware: "+info.FullMethod)
	defer span.End()

	return handler(ctx, req)
}

func (srv *TokenService) LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// Invoke 'handler' to use your gRPC server implementation and get
	// the response.

	srv.lg.Infow("got request",
		"req", req,
		"handler", handler,
	)

	resp, err := handler(ctx, req)
	srv.lg.Infow("created answer",
		"resp", resp,
		"err", err,
	)

	return resp, err
}
