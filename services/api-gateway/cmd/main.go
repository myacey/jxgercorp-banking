package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/myacey/jxgercorp-banking/services/libs/telemetry"

	"github.com/myacey/jxgercorp-banking/services/api-gateway/internal/config"
	"github.com/myacey/jxgercorp-banking/services/api-gateway/internal/httpserver"
	"github.com/myacey/jxgercorp-banking/services/api-gateway/internal/pkg/grpcclient"
)

var cfgPath = flag.String("f", "./configs/config.yaml", "path to the api-gateway's config")

func main() {
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Telemetry
	tracer, metricExporter, err := telemetry.StartTracer("api-gateway", "0.0.1")
	if err != nil {
		log.Fatal(err)
	}
	defer tracer.Shutdown(context.Background())
	defer metricExporter.Shutdown(context.Background())

	// Config
	cfg, err := config.LoadConfig(*cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	// grpc client
	grpcConn, err := grpcclient.MustInitConnection(cfg.GrpcCfg)
	if err != nil {
		log.Fatal(err)
	}
	defer grpcConn.Close()

	app := httpserver.New(cfg, grpcclient.New(grpcConn))
	go func() {
		<-ctx.Done()
		app.Stop(ctx)
	}()

	if err := app.Start(ctx); err != nil {
		log.Printf("http server returned error: %v", err)
	}
}
