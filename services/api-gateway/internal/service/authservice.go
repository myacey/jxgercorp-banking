package service

import (
	"context"

	"github.com/myacey/jxgercorp-banking/services/libs/apperror"
)

type tokenSrv interface {
	ValidateToken(ctx context.Context, token string) (string, bool, error)
}

type Auth struct {
	tokenSrv tokenSrv // api gateway is already a bit big, so didn't move grpc client to new grpc pkg...
}

func NewAuthService(conn tokenSrv) *Auth {
	return &Auth{conn}
}

func (srv *Auth) ValidateToken(ctx context.Context, token string) (username string, valid bool, err error) {
	usrname, valid, err := srv.tokenSrv.ValidateToken(ctx, token)
	if err != nil {
		return "", valid, apperror.NewInternal("failed to validate token", err)
	}

	return usrname, valid, nil
}
