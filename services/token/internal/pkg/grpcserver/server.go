package grpcserver

import (
	"context"
	"log"
	"net"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"

	tokenpb "github.com/myacey/jxgercorp-banking/services/libs/proto/api/token"
	"github.com/myacey/jxgercorp-banking/services/token/internal/service"
)

type Config struct {
	Address          string        `mapstructure:"listen"`
	EnableTLS        bool          `mapstructure:"enable_tls"`
	CertFile         string        `mapstructure:"cert_file"`
	KeyFile          string        `mapstructure:"key_file"`
	KeepAliveTime    time.Duration `mapstructure:"keep_alive_time"`
	KeepAliveTimeout time.Duration `mapstructure:"keep_alive_timeout"`
}

type Server struct {
	cfg        Config
	srv        *service.Service
	grpcServer *grpc.Server
	lis        net.Listener

	tokenpb.UnimplementedTokenServiceServer

	tracer trace.Tracer
}

func New(cfg Config, service *service.Service) (*Server, error) {
	var options []grpc.ServerOption
	options = append(options, grpc.KeepaliveParams(keepalive.ServerParameters{
		Time:    cfg.KeepAliveTime,
		Timeout: cfg.KeepAliveTimeout,
	}))

	if cfg.EnableTLS {
		creds, err := credentials.NewServerTLSFromFile(cfg.CertFile, cfg.KeyFile)
		if err != nil {
			return nil, err
		}
		options = append(options, grpc.Creds(creds))
	}

	srv := &Server{
		cfg: cfg,
		srv: service,

		tracer: otel.Tracer("server"),
	}
	srv.grpcServer = grpc.NewServer(grpc.ChainUnaryInterceptor(srv.TracingMiddleware))

	tokenpb.RegisterTokenServiceServer(srv.grpcServer, srv)

	return srv, nil
}

func (s *Server) GenerateToken(ctx context.Context, req *tokenpb.GenerateTokenRequest) (*tokenpb.GenerateTokenResponse, error) {
	token, err := s.srv.GenerateToken(ctx, req.GetUsername())
	return &tokenpb.GenerateTokenResponse{
		Token: token,
	}, err
}

func (s *Server) ValidateToken(ctx context.Context, req *tokenpb.ValidateTokenRequest) (*tokenpb.ValidateTokenResponse, error) {
	m, err := s.srv.ValidateToken(ctx, req.GetToken())
	if err != nil {
		return &tokenpb.ValidateTokenResponse{
			Username: "",
			Valid:    false,
		}, err
	}
	return &tokenpb.ValidateTokenResponse{
		Valid:    true,
		Username: m["username"].(string),
	}, nil
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", s.cfg.Address)
	if err != nil {
		return err
	}

	s.lis = lis

	log.Printf("starting gRPC server on %s", s.cfg.Address)
	go func() {
		if err := s.grpcServer.Serve(lis); err != nil {
			log.Fatalf("gRPC Serve error: %v", err)
		}
	}()

	return nil
}

func (s *Server) Stop() {
	log.Println("shutting down gRPC server...")
	s.grpcServer.GracefulStop()
}
