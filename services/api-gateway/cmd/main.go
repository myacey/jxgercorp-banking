package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/myacey/jxgercorp-banking/api-gateway/internal/handlers"
	"github.com/myacey/jxgercorp-banking/shared/logging"
)

func main() {
	r := gin.Default()

	api := r.Group("/api/v1")
	{
		api.Any("/user/register", handlers.ProxyHandler("http://localhost:8081"))
		api.Any("/user/login", handlers.ProxyHandler("http://localhost:8081"))
	}

	lg, err := logging.ConfigureLogger()
	if err != nil {
		log.Fatalf("cant initialize logger: %v", err)
	}
	lg.Info("API Gateway running on :80")
	r.Run(":80")
}
