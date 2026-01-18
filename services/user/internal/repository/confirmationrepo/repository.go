package confirmationrepo

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var ErrInvalidCode = errors.New("invalid code")

type RedisConfirmationCodes struct {
	store *redis.Client

	tracer trace.Tracer
}

func NewConfirmationCodesRepo(store *redis.Client) *RedisConfirmationCodes {
	return &RedisConfirmationCodes{store, otel.Tracer("repository-confirmation")}
}

func (cc *RedisConfirmationCodes) CreateCode(ctx context.Context, username, code string) error {
	ctx, span := cc.tracer.Start(ctx, "confirmation-repository: CreateCode")
	defer span.End()

	res := cc.store.Set(ctx, username, code, 24*time.Hour)
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

func (cc *RedisConfirmationCodes) GetCode(ctx context.Context, username string) (string, error) {
	ctx, span := cc.tracer.Start(ctx, "confirmation-repository: GetCode")
	defer span.End()

	res := cc.store.Get(ctx, username)
	if res.Err() != nil {
		switch {
		case errors.Is(res.Err(), redis.Nil):
			return "", ErrInvalidCode
		default:
			return "", res.Err()
		}
	}

	return res.Val(), nil
}
