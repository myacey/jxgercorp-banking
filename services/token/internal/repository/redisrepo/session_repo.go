package redisrepo

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/myacey/jxgercorp-banking/services/shared/cstmerr"
	"github.com/myacey/jxgercorp-banking/services/token/internal/models"
	"github.com/myacey/jxgercorp-banking/services/token/internal/repository"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type RedisTokenRepo struct {
	store *redis.Client
	lg    *zap.SugaredLogger

	tracer trace.Tracer
}

func NewRedisTokenRepo(store *redis.Client, logger *zap.SugaredLogger, tr trace.Tracer) repository.TokenRepository {
	return &RedisTokenRepo{store: store, lg: logger, tracer: tr}
}

func (r *RedisTokenRepo) CreateToken(c context.Context, payload string, username string, ttl time.Duration) error {
	c, span := r.tracer.Start(c, "repository: CreateToken")
	defer span.End()
	span.SetAttributes(
		attribute.String("username", username),
		attribute.Int64("ttl", ttl.Milliseconds()),
	)

	marshalled, err := json.Marshal(payload)
	if err != nil {
		return cstmerr.New(http.StatusInternalServerError, cstmerr.ErrUnknown.Error(), err)
	}

	return r.store.Set(c, username, marshalled, ttl).Err()
}

func (r *RedisTokenRepo) GetToken(c context.Context, username string) (string, error) {
	c, span := r.tracer.Start(c, "repository: GetToken")
	defer span.End()
	span.SetAttributes(
		attribute.String("username", username),
	)

	token, err := r.store.Get(c, username).Result()

	token = strings.TrimLeft(token, "\\\"")
	token = strings.TrimRight(token, "\\\"")

	if err != nil {
		switch {
		case errors.Is(err, redis.Nil):
			return "", nil // key just dont exists
		default:
			return "", err // another error
		}
	}

	return token, nil
}

func (r *RedisTokenRepo) UpdateToken(c context.Context, newPayload *models.Payload, username string, ttl time.Duration) error {
	c, span := r.tracer.Start(c, "repository: UpdateToken")
	defer span.End()
	span.SetAttributes(
		attribute.String("username", username),
		attribute.Int64("ttl", ttl.Milliseconds()),
	)

	marshalled, err := json.Marshal(newPayload)
	if err != nil {
		return cstmerr.New(http.StatusInternalServerError, cstmerr.ErrUnknown.Error(), err)
	}

	return r.store.Set(c, username, marshalled, ttl).Err()
}
