package service

import (
	"context"
	"fmt"
	"time"

	tokenpb "github.com/myacey/jxgercorp-banking/services/shared/proto/token"
	"github.com/myacey/jxgercorp-banking/services/token/internal/repository"
	"github.com/myacey/jxgercorp-banking/services/token/internal/tokenmaker"
	"go.opentelemetry.io/otel/trace"
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

	tracer trace.Tracer
}

func NewTokenService(r repository.TokenRepository, tokenMaker *tokenmaker.PasetoMaker, lg *zap.SugaredLogger, tr trace.Tracer) *TokenService {
	return &TokenService{
		TokenRepo:  r,
		TokenMaker: tokenMaker,
		lg:         lg,
		tracer:     tr,
	}
}

func (t *TokenService) GenerateToken(ctx context.Context, req *tokenpb.GenerateTokenRequest) (*tokenpb.GenerateTokenResponse, error) {
	ctx, span := t.tracer.Start(ctx, "service: GenerateToken")
	defer span.End()

	newToken, err := t.TokenMaker.CreateToken(req.Username, req.Ttl.AsTime())
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
	ctx, span := t.tracer.Start(ctx, "service: ValidateToken")
	defer span.End()

	payload, username, err := t.TokenMaker.VerifyToken(req.Token)
	if err != nil {
		t.lg.Error(err)
		return nil, err
	}

	dbToken, err := t.TokenRepo.GetToken(ctx, username)
	t.lg.Infow("validate token in progress", "db token", dbToken, "providen payload", payload, "username", username)
	if err != nil {
		return &tokenpb.ValidateTokenResponse{Valid: false}, err
	}

	if _, dbUsername, err := t.TokenMaker.VerifyToken(dbToken); err != nil || dbToken != req.Token || username != dbUsername {
		return &tokenpb.ValidateTokenResponse{Valid: false}, fmt.Errorf("invalid")
	}

	return &tokenpb.ValidateTokenResponse{
		Valid:    true,
		Username: username,
	}, nil
}
