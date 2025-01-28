package main

import (
	"net"

	"github.com/myacey/jxgercorp-banking/shared/backconfig"
	"github.com/myacey/jxgercorp-banking/shared/logging"
	tokenpb "github.com/myacey/jxgercorp-banking/shared/proto/token"
	"github.com/myacey/jxgercorp-banking/token/internal/repository/redisrepo"
	"github.com/myacey/jxgercorp-banking/token/internal/service"
	"github.com/myacey/jxgercorp-banking/token/internal/tokenmaker"
	"google.golang.org/grpc"
)

func main() {
	lg, err := logging.ConfigureLogger()
	if err != nil {
		panic(err)
	}

	tokenMaker, err := tokenmaker.NewPaseto("sYHS6QnCtR2KxyJkPR4mKubZh2HLuJQF")
	if err != nil {
		panic(err)
	}

	config, err := backconfig.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	rdb, err := redisrepo.ConfigureRedisClient(&config)
	if err != nil {
		panic(err)
	}
	tokenRepo := redisrepo.NewRedisTokenRepo(rdb, lg)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()
	tokenService := service.NewTokenService(tokenRepo, tokenMaker, lg)
	tokenpb.RegisterTokenServiceServer(server, tokenService)

	lg.Info("start gRPC service")
	if err := server.Serve(listener); err != nil {
		panic(err)
	}
}
