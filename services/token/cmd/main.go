package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/myacey/jxgercorp-banking/services/libs/telemetry"
	"github.com/myacey/jxgercorp-banking/services/token/internal/config"
	"github.com/myacey/jxgercorp-banking/services/token/internal/pkg/grpcserver"
	"github.com/myacey/jxgercorp-banking/services/token/internal/service"
)

var cfgPath = flag.String("f", "./configs/config.yaml", "path to the user service config")

func main() {
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// config
	cfg, err := config.LoadConfig(*cfgPath)
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	// Telemetry
	tracer, _, err := telemetry.StartTracer("token-microservice", "0.0.1")
	if err != nil {
		log.Fatalf("failed to start tracer: %v", err)
	}
	defer tracer.Shutdown(ctx)

	service := service.Service{*service.NewToken(cfg.JwtSecretKey)}
	grpcServer, err := grpcserver.New(cfg.GrpcConfig, &service)
	if err != nil {
		log.Fatalf("failed to craete server: %v", err)
	}

	if err := grpcServer.Start(); err != nil {
		log.Fatalf("grpc server error: %v", err)
	}

	<-ctx.Done()
	grpcServer.Stop()
}
