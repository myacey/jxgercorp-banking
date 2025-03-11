package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/myacey/jxgercorp-banking/services/api-gateway/internal/handlers"
	"github.com/myacey/jxgercorp-banking/services/shared/logging"
	tokenpb "github.com/myacey/jxgercorp-banking/services/shared/proto/token"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-lib/metrics"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
)

func main() {
	// Sample configuration for testing. Use constant sampling to sample every trace
	// and enable LogSpan to log every span via configured Logger.
	cfg := jaegercfg.Configuration{
		ServiceName: "API-GATEWAY",
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
		},
	}
	// Example logger and metrics factory. Use github.com/uber/jaeger-client-go/log
	// and github.com/uber/jaeger-lib/metrics respectively to bind to real logging and metrics
	// frameworks.
	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory

	// Initialize tracer with a logger and a metrics factory
	tracer, closer, err := cfg.NewTracer(
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
	)
	// Set the singleton opentracing.Tracer with the Jaeger tracer.
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()

	logger, err := logging.ConfigureLogger()
	if err != nil {
		log.Fatal("cant configure logger:", err)
	}
	defer logger.Sync()

	grpcConn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal(err)
	}
	defer grpcConn.Close()
	tokenServiceRPC := tokenpb.NewTokenServiceClient(grpcConn)

	r := gin.Default()
	r.Use(handlers.TracingMiddleware(tracer))

	handl := handlers.NewHandler(tokenServiceRPC, logger)

	public := r.Group("/api/v1")
	{
		public.Any("/user/register", handlers.ProxyHandler("http://localhost:8081"))
		public.Any("/user/login", handlers.ProxyHandler("http://localhost:8081"))
		public.Any("/user/confirm", handlers.ProxyHandler("http://localhost:8081"))
	}

	protected := r.Group("/api/v1")
	protected.Use(handl.AuthTokenMiddleware())
	{
		protected.Any("/user/balance", handlers.ProxyHandler("http://localhost:8081"))
		protected.Any("/transaction/create", handlers.ProxyHandler("http://localhost:8082"))
		// protected.Any("/transaction/search", handlers.ProxyHandler("http://localhost:8082"))
		protected.Match([]string{http.MethodOptions, http.MethodGet}, "/transaction/search", handlers.ProxyHandler("http://localhost:8082"))
	}

	lg, err := logging.ConfigureLogger()
	if err != nil {
		log.Fatalf("cant initialize logger: %v", err)
	}
	lg.Info("API Gateway running on :80")
	r.Run(":80")
}
