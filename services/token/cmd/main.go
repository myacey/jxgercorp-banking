package main

import (
	"context"
	"net"

	"github.com/myacey/jxgercorp-banking/services/shared/backconfig"
	"github.com/myacey/jxgercorp-banking/services/shared/logging"
	tokenpb "github.com/myacey/jxgercorp-banking/services/shared/proto/token"
	"github.com/myacey/jxgercorp-banking/services/shared/telemetry"
	"github.com/myacey/jxgercorp-banking/services/token/internal/repository/redisrepo"
	"github.com/myacey/jxgercorp-banking/services/token/internal/service"
	"github.com/myacey/jxgercorp-banking/services/token/internal/tokenmaker"
	"google.golang.org/grpc"
)

func main() {
	// Telemetry
	tracer, _, err := telemetry.StartTracer("token-service", "0.0.1")
	if err != nil {
		panic(err)
	}
	defer tracer.Shutdown(context.Background())

	// logger
	lg, err := logging.ConfigureLogger()
	if err != nil {
		panic(err)
	}

	// token maker
	tokenMaker, err := tokenmaker.NewPaseto("sYHS6QnCtR2KxyJkPR4mKubZh2HLuJQF")
	if err != nil {
		panic(err)
	}

	// config
	config, err := backconfig.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	// redis
	rdb, err := redisrepo.ConfigureRedisClient(&config)
	if err != nil {
		panic(err)
	}
	tokenRepo := redisrepo.NewRedisTokenRepo(rdb, lg, tracer.Tracer("repository"))

	// net
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	tokenService := service.NewTokenService(tokenRepo, tokenMaker, lg, tracer.Tracer("service"))
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(tokenService.TracingMiddleware, tokenService.LoggingInterceptor))
	tokenpb.RegisterTokenServiceServer(server, tokenService)

	lg.Info("start gRPC service")
	if err := server.Serve(listener); err != nil {
		panic(err)
	}
}
