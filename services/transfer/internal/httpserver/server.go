package httpserver

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/myacey/jxgercorp-banking/services/libs/web"

	"github.com/myacey/jxgercorp-banking/services/transfer/internal/config"
	"github.com/myacey/jxgercorp-banking/services/transfer/internal/httpserver/handler"
	"github.com/myacey/jxgercorp-banking/services/transfer/internal/repository/accountrepo"
	db "github.com/myacey/jxgercorp-banking/services/transfer/internal/repository/sqlc"
	"github.com/myacey/jxgercorp-banking/services/transfer/internal/repository/transferrepo"
	"github.com/myacey/jxgercorp-banking/services/transfer/internal/service"
)

type App struct {
	server  web.Server
	router  *gin.Engine
	service *service.Service
}

func New(cfg config.AppConfig, conn *pgxpool.Pool, queries *db.Queries) (*App, error) {
	app := initialize(conn, queries)

	app.server = web.NewServer(cfg.HTTPServerCfg, app.router)
	return app, nil
}

func (app *App) Start(ctx context.Context) error {
	return app.server.Run(ctx)
}

func (app *App) Stop(ctx context.Context) error {
	return app.server.Shutdown(ctx)
}

func initialize(conn *pgxpool.Pool, queries *db.Queries) *App {
	app := &App{}

	accountRepo := accountrepo.NewPostgresAccount(queries)
	transferRepo := transferrepo.NewPostgresTransfer(queries)

	accountSrv := *service.NewAccount(accountRepo)
	app.service = &service.Service{
		Transfer: *service.NewTransfer(conn, &accountSrv, transferRepo),
		Account:  accountSrv,
	}

	handlr := handler.NewHandler(&app.service.Transfer, &app.service.Account)

	app.router = gin.Default()
	app.router.ContextWithFallback = true
	app.router.Use(handlr.TracingMiddleware())

	app.router.POST("/api/v1/transfer", handlr.CreateTransfer)
	app.router.GET("/api/v1/transfer", handlr.SearchTransfersWithAccount)

	app.router.POST("/api/v1/transfer/account", handlr.CreateAccount)
	app.router.GET("/api/v1/transfer/account", handlr.SearchAccounts)

	return app
}
