package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/myacey/jxgercorp-banking/services/libs/telemetry"
	"github.com/myacey/jxgercorp-banking/services/transfer/internal/config"
	"github.com/myacey/jxgercorp-banking/services/transfer/internal/httpserver"
	"github.com/myacey/jxgercorp-banking/services/transfer/internal/repository"
)

var cfgPath = flag.String("f", "./configs/config.yaml", "path to the transfer config")

func main() {
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// config
	cfg, err := config.LoadConfig(*cfgPath)
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	// postgres
	psqlQueries, psqlPool, err := repository.ConfiurePostgres(ctx, cfg.PostgresCfg)
	if err != nil {
		log.Fatal(err)
	}
	defer psqlPool.Close()

	// Tracer for Telemetry (Jaeger)
	_, _, err = telemetry.StartTracer("transfer-microservice", "0.0.1")
	if err != nil {
		log.Fatal(err)
	}
	_ = telemetry.NewMetricsFactory("transfer-microservice")

	app, err := httpserver.New(cfg, psqlPool, psqlQueries)
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
