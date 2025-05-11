package httpserver

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/myacey/jxgercorp-banking/services/api-gateway/internal/config"
	"github.com/myacey/jxgercorp-banking/services/api-gateway/internal/httpserver/handler"
	"github.com/myacey/jxgercorp-banking/services/api-gateway/internal/pkg/grpcclient"
	"github.com/myacey/jxgercorp-banking/services/api-gateway/internal/service"
	"github.com/myacey/jxgercorp-banking/services/libs/web"
)

type App struct {
	server  web.Server
	router  *gin.Engine
	service *service.Service
}

func New(cfg config.AppConfig, grpcClient *grpcclient.ClientImpl) *App {
	app := &App{
		router: gin.Default(),
	}
	app.server = web.NewServer(cfg.HTTPServerCfg, app.router)

	app.initialize(cfg, grpcClient)

	return app
}

func (app *App) Start(ctx context.Context) error {
	return app.server.Run(ctx)
}

func (app *App) Stop(ctx context.Context) error {
	return app.server.Shutdown(ctx)
}

func (app *App) initialize(cfg config.AppConfig, grpcClient *grpcclient.ClientImpl) {
	app.service = &service.Service{
		Auth: *service.NewAuthService(grpcClient),
	}

	handl := handler.NewHandler(*app.service)
	app.router.Use(handl.MetricsMiddleware())
	app.router.Use(handl.TracingMiddleware())

	public := app.router.Group("/api/v1/user")
	{
		public.Match(
			[]string{http.MethodOptions, http.MethodPost},
			"/register",
			handl.ProxyHandler(cfg.Services["user-service"]),
		)
		public.Match(
			[]string{http.MethodOptions, http.MethodPost},
			"/login",
			handl.ProxyHandler(cfg.Services["user-service"]),
		)
		public.Match(
			[]string{http.MethodOptions, http.MethodGet},
			"/confirm",
			handl.ProxyHandler(cfg.Services["user-service"]),
		)
	}

	protected := app.router.Group("/api/v1/transfer")
	protected.Use(handl.AuthTokenMiddleware())
	{
		protected.Match(
			[]string{http.MethodOptions, http.MethodGet, http.MethodPost},
			"/account",
			handl.ProxyHandler(cfg.Services["transfer-service"]),
		)

		protected.Match(
			[]string{http.MethodOptions, http.MethodGet, http.MethodPost},
			"",
			handl.ProxyHandler(cfg.Services["transfer-service"]),
		)
	}
}
