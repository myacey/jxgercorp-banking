package grpcclient

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	tokenpb "github.com/myacey/jxgercorp-banking/services/libs/proto/api/token"
)

type Config struct {
	Target string `mapstructure:"target"`
}

func MustInitConnection(cfg Config) (*grpc.ClientConn, error) {
	grpcConn, err := grpc.NewClient(
		cfg.Target,
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

func (c *ClientImpl) ValidateToken(ctx context.Context, token string) (string, bool, error) {
	resp, err := c.cli.ValidateToken(ctx, &tokenpb.ValidateTokenRequest{
		Token: token,
	})
	if err != nil {
		return "", false, err
	}

	return resp.GetUsername(), resp.GetValid(), nil
}
