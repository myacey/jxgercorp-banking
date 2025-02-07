package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/myacey/jxgercorp-banking/services/api-gateway/internal/handlers"
	"github.com/myacey/jxgercorp-banking/services/shared/logging"
	tokenpb "github.com/myacey/jxgercorp-banking/services/shared/proto/token"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
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
	// r.Use(func(c *gin.Context) {
	// 	c.Header("Access-Control-Allow-Origin", "http://localhost:8080")
	// 	c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
	// 	c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
	// 	c.Header("Access-Control-Allow-Credentials", "true")

	// 	if c.Request.Method == "OPTIONS" {
	// 		c.AbortWithStatus(http.StatusNoContent) // Пропускаем `OPTIONS` без авторизации
	// 		return
	// 	}

	// 	c.Next()
	// })

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
		// protected.("/transaction/search", handlers.ProxyHandler("http://localhost:8082"))
	}

	lg, err := logging.ConfigureLogger()
	if err != nil {
		log.Fatalf("cant initialize logger: %v", err)
	}
	lg.Info("API Gateway running on :80")
	r.Run(":80")
}
