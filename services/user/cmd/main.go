package main

import (
	"context"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/myacey/jxgercorp-banking/services/shared/backconfig"
	"github.com/myacey/jxgercorp-banking/services/shared/logging"
	tokenpb "github.com/myacey/jxgercorp-banking/services/shared/proto/token"
	"github.com/myacey/jxgercorp-banking/services/shared/telemetry"
	"github.com/myacey/jxgercorp-banking/services/user/internal/confirmation"
	"github.com/myacey/jxgercorp-banking/services/user/internal/controller"
	"github.com/myacey/jxgercorp-banking/services/user/internal/repository/postgresrepo"
	"github.com/myacey/jxgercorp-banking/services/user/internal/repository/redisrepo"
	"github.com/myacey/jxgercorp-banking/services/user/internal/service"
)

func main() {
	// config
	config, err := backconfig.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	// logger
	logger, err := logging.ConfigureLogger()
	if err != nil {
		log.Fatal("cant initialize logger:", err)
	}
	defer logger.Sync()

	// postgres
	psqlQueries, conn, err := postgresrepo.ConfiurePostgres(config)
	if err != nil {
		logger.Fatal(err)
	}
	defer conn.Close()
	logger.Debug("postgres conn initialized")

	// Tracer for Telemetry (Jaeger)
	tp, _, err := telemetry.StartTracer("user-service", "0.0.1")
	if err != nil {
		panic(err)
	}
	defer tp.Shutdown(context.Background())
	metricsFactory := telemetry.NewMetricsFactory("user-service")
	userMetrics := metricsFactory.NewUserMetrics()

	// user repository
	userRepo := postgresrepo.NewUserRepo(psqlQueries, logger, tp.Tracer("repository"))

	// grpc conn with token service
	grpcConn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal(err)
	}
	defer grpcConn.Close()
	tokenServiceRPC := tokenpb.NewTokenServiceClient(grpcConn)

	// redis
	rdb, err := redisrepo.ConfigureRedisClient(&config)
	if err != nil {
		panic(err)
	}

	// CONFIRMATION SERVICE
	// repo
	confirmRepo := redisrepo.NewConfirmationCodesRepo(rdb, logger, tp.Tracer("confirmation-repository"))
	cnfrmService := confirmation.NewConfirmationService(confirmRepo, "notif.email.register.confirm", 0, tp.Tracer("confirmation-service"))

	srv := service.NewService(userRepo, tokenServiceRPC, logger, cnfrmService, tp.Tracer("service"))

	ctrller := controller.NewController(srv, tokenServiceRPC, logger, tp.Tracer("controller"), *userMetrics)

	r := gin.Default()
	r.ContextWithFallback = true
	r.Use(ctrller.TracingMiddleware())

	// add CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080"},                            // Разрешённый фронтенд-адрес
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}, // Разрешённые методы
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},          // Разрешённые заголовки
		ExposeHeaders:    []string{"Content-Length"},                                   // Доступные клиенту заголовки
		AllowCredentials: true,                                                         // Разрешить куки и авторизационные токены
	}))
	r.POST("/api/v1/user/register", ctrller.CreateUser)
	r.POST("/api/v1/user/login", ctrller.Login)
	r.POST("/api/v1/user/confirm", ctrller.ConfirmUserEmail)

	r.GET("/api/v1/user/balance", ctrller.GetUserBalance)

	logger.Info("User microservice running on :8081")
	r.Run("localhost:8081")
}
