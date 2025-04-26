package grpcserver

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func (s *Server) TracingMiddleware(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
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

	ctx, span := s.tracer.Start(ctx, "middleware: "+info.FullMethod)
	defer span.End()

	return handler(ctx, req)
}
