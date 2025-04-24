package httpserver

import (
	"context"
	"database/sql"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/myacey/jxgercorp-banking/services/libs/web"
	"github.com/myacey/jxgercorp-banking/services/user/internal/config"
	"github.com/myacey/jxgercorp-banking/services/user/internal/httpserver/handler"
	"github.com/myacey/jxgercorp-banking/services/user/internal/pkg/grpcclient"
	"github.com/myacey/jxgercorp-banking/services/user/internal/pkg/hasher"
	"github.com/myacey/jxgercorp-banking/services/user/internal/repository/confirmationrepo"
	db "github.com/myacey/jxgercorp-banking/services/user/internal/repository/sqlc"
	"github.com/myacey/jxgercorp-banking/services/user/internal/repository/userrepo"
	"github.com/myacey/jxgercorp-banking/services/user/internal/service"
)

type App struct {
	server  web.Server
	router  *gin.Engine
	service *service.Service
}

func New(cfg config.AppConfig, conn *sql.DB, queries *db.Queries, store *redis.Client) (*App, error) {
	app := &App{}
	err := app.initialize(cfg, conn, queries, store)
	if err != nil {
		return nil, err
	}

	app.server = web.NewServer(cfg.HTTPServerCfg, app.router)
	return app, nil
}

func (app *App) Start(ctx context.Context) error {
	return app.server.Run(ctx)
}

func (app *App) Stop(ctx context.Context) error {
	return app.server.Shutdown(ctx)
}

func (app *App) initialize(cfg config.AppConfig, conn *sql.DB, queries *db.Queries, store *redis.Client) error {
	grpcConn, err := grpc.NewClient(
		cfg.GrpcTarget,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatal("failed to init grpc conn: %w", err)
	}
	defer grpcConn.Close()
	grpcClientService := grpcclient.New(grpcConn)

	usrRepo := userrepo.NewUserRepo(queries)
	confirmRepo := confirmationrepo.NewConfirmationCodesRepo(store)

	confirmSrv := service.NewConfirmationService(confirmRepo, cfg.KafkaCfg)
	hasherSrv := hasher.NewBcrypt()
	app.service = &service.Service{
		User: *service.NewUserSrv(usrRepo, *confirmSrv, grpcClientService, hasherSrv),
	}

	handlr := handler.NewHandler(&app.service.User)

	app.router = gin.Default()
	app.router.ContextWithFallback = true
	app.router.Use(handlr.TracingMiddleware())

	// add CORS
	app.router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	app.router.POST("/api/v1/user/register", handlr.CreateUser)
	app.router.POST("/api/v1/user/login", handlr.Login)
	app.router.POST("/api/v1/user/confirm", handlr.ConfirmUserEmail)

	// app.router.GET("/api/v1/user/balance", handlr.)

	return nil
}
