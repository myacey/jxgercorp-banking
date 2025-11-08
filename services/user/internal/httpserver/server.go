package httpserver

import (
	"context"
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/myacey/jxgercorp-banking/services/libs/web"
	"github.com/myacey/jxgercorp-banking/services/user/internal/config"
	"github.com/myacey/jxgercorp-banking/services/user/internal/httpserver/handler"
	"github.com/myacey/jxgercorp-banking/services/user/internal/pkg/grpcclient"
	"github.com/myacey/jxgercorp-banking/services/user/internal/pkg/hasher"
	"github.com/myacey/jxgercorp-banking/services/user/internal/pkg/kafka"
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

func New(
	cfg config.AppConfig,
	conn *sql.DB,
	queries *db.Queries,
	store *redis.Client,
	grpcClient *grpcclient.ClientImpl,
) (*App, error) {
	app := &App{}
	err := app.initialize(cfg, conn, queries, store, grpcClient)
	app.server = web.NewServer(cfg.HTTPServerCfg, app.router)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (app *App) Start(ctx context.Context) error {
	return app.server.Run(ctx)
}

func (app *App) Stop(ctx context.Context) error {
	return app.server.Shutdown(ctx)
}

func (app *App) initialize(
	cfg config.AppConfig,
	conn *sql.DB,
	queries *db.Queries,
	store *redis.Client,
	grpcClient *grpcclient.ClientImpl,
) error {
	usrRepo := userrepo.NewUserRepo(queries)
	confirmRepo := confirmationrepo.NewConfirmationCodesRepo(store)

	kafkaProducer := kafka.NewProducer(cfg.KafkaCfg)

	confirmSrv := service.NewConfirmationService(confirmRepo, kafkaProducer, cfg.ConfirmationCfg)
	hasherSrv := hasher.NewBcrypt()
	app.service = &service.Service{
		User: *service.NewUserSrv(usrRepo, *confirmSrv, grpcClient, hasherSrv),
	}

	handlr := handler.NewHandler(&app.service.User, cfg.HandlerCfg)

	app.router = gin.Default()
	app.router.ContextWithFallback = true
	app.router.Use(handlr.TracingMiddleware())

	app.router.POST("/api/v1/user/register", handlr.CreateUser)
	app.router.POST("/api/v1/user/login", handlr.Login)
	app.router.GET("/api/v1/user/confirm", handlr.ConfirmUserEmail)

	return nil
}
