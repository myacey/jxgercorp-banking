package service

import (
	"context"

	"google.golang.org/grpc"
)

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
