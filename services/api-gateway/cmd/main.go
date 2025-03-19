package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/myacey/jxgercorp-banking/services/api-gateway/internal/handlers"
	"github.com/myacey/jxgercorp-banking/services/shared/logging"
	tokenpb "github.com/myacey/jxgercorp-banking/services/shared/proto/token"
	"github.com/myacey/jxgercorp-banking/services/shared/telemetry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Telemetry
	tracer, err := telemetry.StartTracer("api-gateway", "0.0.1")
	if err != nil {
		panic(err)
	}
	defer tracer.Shutdown(context.Background())

	logger, err := logging.ConfigureLogger()
	if err != nil {
		log.Fatal("cant configure logger:", err)
	}
	defer logger.Sync()

	grpcConn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(handlers.UnaryClientInterceptor),
	)
	if err != nil {
		logger.Fatal(err)
	}
	defer grpcConn.Close()
	tokenServiceRPC := tokenpb.NewTokenServiceClient(grpcConn)

	handl := handlers.NewHandler(tokenServiceRPC, logger, tracer.Tracer("api-gateway"))

	r := gin.Default()
	r.Use(handl.TracingMiddleware())

	public := r.Group("/api/v1")
	{
		// TODO: change addresses
		public.Any("/user/register", handl.ProxyHandler("http://localhost:8081"))
		public.Any("/user/login", handl.ProxyHandler("http://localhost:8081"))
		public.Any("/user/confirm", handl.ProxyHandler("http://localhost:8081"))
	}

	protected := r.Group("/api/v1")
	protected.Use(handl.AuthTokenMiddleware())
	{
		// TODO: change addresses
		protected.Any("/user/balance", handl.ProxyHandler("http://localhost:8081"))
		protected.Any("/transaction/create", handl.ProxyHandler("http://localhost:8082"))
		// protected.Any("/transaction/search", handl.ProxyHandler("http://localhost:8082"))
		protected.Match([]string{http.MethodOptions, http.MethodGet}, "/transaction/search", handl.ProxyHandler("http://localhost:8082"))
	}

	lg, err := logging.ConfigureLogger()
	if err != nil {
		log.Fatalf("cant initialize logger: %v", err)
	}
	lg.Info("API Gateway running on :80")
	r.Run(":80")
}
