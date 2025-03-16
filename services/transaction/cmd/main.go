package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/myacey/jxgercorp-banking/services/shared/backconfig"
	"github.com/myacey/jxgercorp-banking/services/shared/logging"
	"github.com/myacey/jxgercorp-banking/services/transaction/internal/controller"
	"github.com/myacey/jxgercorp-banking/services/transaction/internal/repository/postgresrepo"
	"github.com/myacey/jxgercorp-banking/services/transaction/internal/service"
)

func main() {
	config, err := backconfig.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	lg, err := logging.ConfigureLogger()
	if err != nil {
		panic(err)
	}

	psqlQueries, conn, err := postgresrepo.ConfiurePostgres(config)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	lg.Debug("postgres conn initialized")

	trxRepo := postgresrepo.NewPostgresTransactionRepo(psqlQueries, conn, lg)
	srv := service.NewService(trxRepo, lg)
	ctrller := controller.NewController(srv, lg)

	r := gin.Default()

	// add CORS
	// cfg := cors.DefaultConfig()
	// cfg.AddAllowHeaders("Authorization")
	// cfg.AllowCredentials = true
	// // cfg.AllowAllOrigins = false
	// // I think you should whitelist a limited origins instead:
	// cfg.AllowOrigins = []string{"http://localhost:8080"}
	// // config.AllowOriginFunc = func(origin string) bool {
	// // return true
	// // }
	// r.Use(cors.New(cfg))

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080"},                            // Разрешённый фронтенд-адрес
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}, // Разрешённые методы
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},          // Разрешённые заголовки
		ExposeHeaders:    []string{"Content-Length"},                                   // Доступные клиенту заголовки
		AllowCredentials: true,                                                         // Разрешить куки и авторизационные токены
	}))

	r.POST("/api/v1/transaction/create", ctrller.CreateNewTransaction)
	r.GET("/api/v1/transaction/search", ctrller.SearchEntriesForUser)

	lg.Info("start microservice running on :8082")
	r.Run("localhost:8082")
}
