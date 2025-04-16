package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/myacey/jxgercorp-banking/services/api-gateway/internal/config"
	"github.com/myacey/jxgercorp-banking/services/api-gateway/internal/httpserver"
	"github.com/myacey/jxgercorp-banking/services/libs/telemetry"
)

var cfgPath = flag.String("f", "./services/api-gateway/configs/config.yaml", "path to the api-gateway's config")

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

	app, err := httpserver.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		<-ctx.Done()
		app.Stop(ctx)
	}()

	if err := app.Start(ctx); err != nil {
		log.Printf("http server returned error: %v", err)
	}
}
