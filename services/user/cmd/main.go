package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/myacey/jxgercorp-banking/services/shared/backconfig"
	"github.com/myacey/jxgercorp-banking/services/shared/logging"
	tokenpb "github.com/myacey/jxgercorp-banking/services/shared/proto/token"
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

	// user repository
	userRepo := postgresrepo.NewUserRepo(psqlQueries, logger)

	// grpc conn with token service
	grpcConn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal(err)
	}
	defer grpcConn.Close()
	tokenServiceRPC := tokenpb.NewTokenServiceClient(grpcConn)

	// CONFIRMATION
	// repo
	rdb, err := redisrepo.ConfigureRedisClient(&config)
	if err != nil {
		panic(err)
	}
	confirmRepo := redisrepo.NewConfirmationCodesRepo(rdb, logger)
	cnfrmService := confirmation.NewConfirmationService(confirmRepo, "notif.email.register.confirm", 0)

	srv := service.NewService(userRepo, tokenServiceRPC, logger, cnfrmService)

	ctrller := controller.NewController(srv, tokenServiceRPC, logger)

	r := gin.Default()

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
