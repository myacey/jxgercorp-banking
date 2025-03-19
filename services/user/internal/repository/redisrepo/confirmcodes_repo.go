package redisrepo

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/myacey/jxgercorp-banking/services/shared/cstmerr"
	"github.com/myacey/jxgercorp-banking/services/user/internal/repository"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type RedisConfirmationCodesRepostory struct {
	store *redis.Client
	lg    *zap.SugaredLogger

	tracer trace.Tracer
}

func NewConfirmationCodesRepo(store *redis.Client, logger *zap.SugaredLogger, tr trace.Tracer) repository.ConfirmCodesRepository {
	return &RedisConfirmationCodesRepostory{store, logger, tr}
}

func (cc *RedisConfirmationCodesRepostory) CreateCode(ctx context.Context, username, code string) error {
	ctx, span := cc.tracer.Start(ctx, "confirmation-repository: CreateCode")
	defer span.End()

	res := cc.store.Set(ctx, username, code, 24*time.Hour)
	if res.Err() != nil {
		return cstmerr.New(http.StatusInternalServerError, cstmerr.ErrUnknown.Error(), res.Err())
	}

	return nil
}

func (cc *RedisConfirmationCodesRepostory) GetCode(ctx context.Context, username string) (string, error) {
	ctx, span := cc.tracer.Start(ctx, "confirmation-repository: GetCode")
	defer span.End()

	res := cc.store.Get(ctx, username)
	if res.Err() != nil {
		switch {
		case errors.Is(res.Err(), redis.Nil):
			return "", cstmerr.New(http.StatusBadRequest, "invalid code", nil)
		default:
			return "", cstmerr.New(http.StatusInternalServerError, cstmerr.ErrUnknown.Error(), res.Err())
		}
	}

	return res.Val(), nil
}
