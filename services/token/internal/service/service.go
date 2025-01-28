package service

import (
	"context"
	"fmt"
	"time"

	tokenpb "github.com/myacey/jxgercorp-banking/shared/proto/token"
	"github.com/myacey/jxgercorp-banking/token/internal/repository"
	"github.com/myacey/jxgercorp-banking/token/internal/tokenmaker"
	"go.uber.org/zap"
)

type TokenServiceInterface interface {
	GenerateToken(ctx context.Context, req *tokenpb.GenerateTokenRequest) (*tokenpb.GenerateTokenResponse, error)
	ValidateToken(ctx context.Context, req *tokenpb.ValidateTokenRequest) (*tokenpb.ValidateTokenResponse, error)
}

type TokenService struct {
	TokenRepo  repository.TokenRepository
	TokenMaker *tokenmaker.PasetoMaker
	lg         *zap.SugaredLogger

	tokenpb.UnimplementedTokenServiceServer
}

func NewTokenService(r repository.TokenRepository, tokenMaker *tokenmaker.PasetoMaker, lg *zap.SugaredLogger) *TokenService {
	return &TokenService{
		TokenRepo:  r,
		TokenMaker: tokenMaker,
		lg:         lg,
	}
}

func (t *TokenService) GenerateToken(ctx context.Context, req *tokenpb.GenerateTokenRequest) (*tokenpb.GenerateTokenResponse, error) {
	newToken, err := t.TokenMaker.CreateToken(req.Username, time.Since(req.Ttl.AsTime()))
	if err != nil {
		return &tokenpb.GenerateTokenResponse{Token: ""}, err
	}

	err = t.TokenRepo.CreateToken(ctx, newToken, req.Username, time.Since(req.Ttl.AsTime()))
	if err != nil {
		return &tokenpb.GenerateTokenResponse{Token: ""}, err
	}

	return &tokenpb.GenerateTokenResponse{
		Token: newToken,
	}, nil
}

func (t *TokenService) ValidateToken(ctx context.Context, req *tokenpb.ValidateTokenRequest) (*tokenpb.ValidateTokenResponse, error) {
	_, dbToken, err := t.TokenRepo.GetToken(ctx, req.Username)
	if err != nil {
		return &tokenpb.ValidateTokenResponse{Valid: false}, err
	}

	if _, err := t.TokenMaker.VerifyToken(dbToken); err != nil || dbToken != req.Token {
		return &tokenpb.ValidateTokenResponse{Valid: false}, fmt.Errorf("invalid")
	}

	return &tokenpb.ValidateTokenResponse{
		Valid: true,
	}, nil
}
