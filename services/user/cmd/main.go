package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"

	"github.com/myacey/jxgercorp-banking/services/libs/telemetry"
	"github.com/myacey/jxgercorp-banking/services/user/internal/config"
	"github.com/myacey/jxgercorp-banking/services/user/internal/httpserver"
	"github.com/myacey/jxgercorp-banking/services/user/internal/repository"
)

var cfgPath = flag.String("f", "./configs/config.yaml", "path to the api-gateway's config")

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
	psqlQueries, conn, err := repository.ConfiurePostgres(cfg.PostgresCfg)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// redis
	rdb, err := repository.ConfigureRedisClient(cfg.RedisCfg)
	if err != nil {
		log.Fatal(err)
	}

	// Tracer for Telemetry (Jaeger)
	_, _, err = telemetry.StartTracer("user-microservice", "0.0.1")
	if err != nil {
		log.Fatal(err)
	}
	_ = telemetry.NewMetricsFactory("user-microservice")
	// userMetrics := metricsFactory.NewUserMetrics()

	app, err := httpserver.New(cfg, conn, psqlQueries, rdb)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		<-ctx.Done()
		app.Stop(ctx)
		// tp.Shutdown(ctx)
	}()

	if err := app.Start(ctx); err != nil {
		log.Printf("http server returned error: %v", err)
	}
}
