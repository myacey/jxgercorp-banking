package grpcclient

import (
	"context"
	"time"

	tokenpb "github.com/myacey/jxgercorp-banking/services/libs/proto/api/token"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type clientImpl struct {
	cli tokenpb.TokenServiceClient
}

func New(conn *grpc.ClientConn) *clientImpl {
	return &clientImpl{
		cli: tokenpb.NewTokenServiceClient(conn),
	}
}

func (c *clientImpl) GenerateToken(ctx context.Context, username string, ttl time.Duration) (string, error) {
	resp, err := c.cli.GenerateToken(ctx, &tokenpb.GenerateTokenRequest{
		Username: username,
		Ttl:      timestamppb.New(time.Now().Add(ttl)),
	})
	if err != nil {
		return "", err
	}

	return resp.GetToken(), nil
}

func (c *clientImpl) ValidateToken(ctx context.Context, token string) (bool, error) {
	resp, err := c.cli.ValidateToken(ctx, &tokenpb.ValidateTokenRequest{
		Token: token,
	})
	if err != nil {
		return false, err
	}

	return resp.GetValid(), nil
}
