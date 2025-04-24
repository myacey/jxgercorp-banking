package grpcclient

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	tokenpb "github.com/myacey/jxgercorp-banking/services/libs/proto/api/token"
)

type Config struct {
	ListenAddr string `mapstructure:"listen"`
}

func MustInitConnection(cfg Config) (*grpc.ClientConn, error) {
	grpcConn, err := grpc.NewClient(
		cfg.ListenAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to init grpc conn: %v", err)
	}

	return grpcConn, nil
}

type ClientImpl struct {
	cli tokenpb.TokenServiceClient
}

func New(conn *grpc.ClientConn) *ClientImpl {
	return &ClientImpl{
		cli: tokenpb.NewTokenServiceClient(conn),
	}
}

func (c *ClientImpl) GenerateToken(ctx context.Context, username string, ttl time.Duration) (string, error) {
	resp, err := c.cli.GenerateToken(ctx, &tokenpb.GenerateTokenRequest{
		Username: username,
		Ttl:      timestamppb.New(time.Now().Add(ttl)),
	})
	if err != nil {
		return "", err
	}

	return resp.GetToken(), nil
}

func (c *ClientImpl) ValidateToken(ctx context.Context, token string) (bool, error) {
	resp, err := c.cli.ValidateToken(ctx, &tokenpb.ValidateTokenRequest{
		Token: token,
	})
	if err != nil {
		return false, err
	}

	return resp.GetValid(), nil
}
