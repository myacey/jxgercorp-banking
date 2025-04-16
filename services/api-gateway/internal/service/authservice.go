package service

import (
	"context"
	"time"

	"github.com/myacey/jxgercorp-banking/services/libs/apperror"
	tokenpb "github.com/myacey/jxgercorp-banking/services/libs/proto/api/token"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Auth struct {
	grpConn tokenpb.TokenServiceClient
}

func NewAuthService(conn tokenpb.TokenServiceClient) *Auth {
	return &Auth{conn}
}

func (srv *Auth) GenerateToken(ctx context.Context, username string, ttl time.Duration) (string, error) {
	res, err := srv.grpConn.GenerateToken(
		ctx,
		&tokenpb.GenerateTokenRequest{
			Username: username,
			Ttl:      timestamppb.New(time.Now().Add(ttl)),
		},
	)
	if err != nil {
		return "", apperror.NewInternal("failed to gen token", err)
	}

	return res.Token, nil
}

func (srv *Auth) ValidateToken(ctx context.Context, token string) (username string, valid bool, err error) {
	res, err := srv.grpConn.ValidateToken(ctx, &tokenpb.ValidateTokenRequest{Token: token})
	if err != nil {
		return "", valid, apperror.NewUnauthorized("failed to validate token")
	}

	return res.Username, res.GetValid(), nil
}
