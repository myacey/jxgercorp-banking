package httpserver

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/myacey/jxgercorp-banking/services/api-gateway/internal/config"
	handlers "github.com/myacey/jxgercorp-banking/services/api-gateway/internal/httpserver/handler"
	"github.com/myacey/jxgercorp-banking/services/api-gateway/internal/service"
	tokenpb "github.com/myacey/jxgercorp-banking/services/libs/proto/api/token"
	"github.com/myacey/jxgercorp-banking/services/libs/web"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	server  web.Server
	router  *gin.Engine
	service *service.Service
}

func New(cfg config.AppConfig) (*App, error) {
	app := &App{
		router: gin.Default(),
	}
	app.server = web.NewServer(cfg.HTTPServerCfg, app.router)

	err := app.initialize(cfg)
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

func (app *App) initialize(cfg config.AppConfig) error {
	grpcConn, err := grpc.NewClient(
		cfg.GrpcTarget,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(handlers.UnaryClientInterceptor),
	)
	if err != nil {
		log.Fatal("failed to init grpc conn: %w", err)
	}
	defer grpcConn.Close()

	app.service = &service.Service{
		Auth: *service.NewAuthService(tokenpb.NewTokenServiceClient(grpcConn)),
	}

	handl := handlers.NewHandler(*app.service)
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
		protected.Any("/user/balance", handl.ProxyHandler("http://localhost:8081"))
		protected.Any("/transaction/create", handl.ProxyHandler("http://localhost:8082"))
		// protected.Any("/transaction/search", handl.ProxyHandler("http://localhost:8082"))
		protected.Match([]string{http.MethodOptions, http.MethodGet}, "/transaction/search", handl.ProxyHandler("http://localhost:8082"))
	}

	return nil
}
