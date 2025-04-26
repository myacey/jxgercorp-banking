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

func New(cfg config.AppConfig, grpcClient *grpcclient.ClientImpl) (*App, error) {
	app := &App{
		router: gin.Default(),
	}
	app.server = web.NewServer(cfg.HTTPServerCfg, app.router)

	err := app.initialize(cfg, grpcClient)
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

func (app *App) initialize(cfg config.AppConfig, grpcClient *grpcclient.ClientImpl) error {
	app.service = &service.Service{
		Auth: *service.NewAuthService(grpcClient),
	}

	handl := handler.NewHandler(*app.service)
	app.router.Use(handl.MetricsMiddleware())
	app.router.Use(handl.TracingMiddleware())

	public := app.router.Group("/api/v1")
	{
		// TODO: change addresses
		public.Any("/user/register", handl.ProxyHandler("http://localhost:8081"))
		public.Any("/user/login", handl.ProxyHandler("http://localhost:8081"))
		public.Any("/user/confirm", handl.ProxyHandler("http://localhost:8081"))
	}

	protected := app.router.Group("/api/v1")
	protected.Use(handl.AuthTokenMiddleware())
	{
		// TODO: change addresses
		// protected.Any("/transfer/account", handl.ProxyHandler("http://localhost:8082"))
		protected.Match(
			[]string{http.MethodOptions, http.MethodGet, http.MethodPost},
			"/transfer/account",
			handl.ProxyHandler("http://localhost:8082"),
		)

		protected.Match(
			[]string{http.MethodOptions, http.MethodGet, http.MethodPost},
			"/transfer",
			handl.ProxyHandler("http://localhost:8082"),
		)

		// protected.Any("/transfer/create", handl.ProxyHandler("http://localhost:8082"))
		// protected.Any("/transfer/search", handl.ProxyHandler("http://localhost:8082"))
		// protected.Match([]string{http.MethodOptions, http.MethodGet}, "/transfer/search", handl.ProxyHandler("http://localhost:8082"))
	}

	return nil
}
