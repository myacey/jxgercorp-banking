package redisrepo

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/myacey/jxgercorp-banking/shared/cstmerr"
	"github.com/myacey/jxgercorp-banking/token/internal/models"
	"github.com/myacey/jxgercorp-banking/token/internal/repository"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisTokenRepo struct {
	store *redis.Client
	lg    *zap.SugaredLogger
}

func NewRedisTokenRepo(store *redis.Client, logger *zap.SugaredLogger) repository.TokenRepository {
	return &RedisTokenRepo{store: store, lg: logger}
}

func (r *RedisTokenRepo) CreateToken(c context.Context, payload string, username string, ttl time.Duration) error {
	marshalled, err := json.Marshal(payload)
	if err != nil {
		return cstmerr.New(http.StatusInternalServerError, cstmerr.ErrUnknown.Error(), err)
	}

	return r.store.Set(c, username, marshalled, ttl).Err()
}

func (r *RedisTokenRepo) GetToken(c context.Context, username string) (*models.Payload, string, error) {
	var payload *models.Payload
	marhalledPayload, err := r.store.Get(c, username).Result()
	if err != nil {
		switch {
		case errors.Is(err, redis.Nil):
			return nil, "", nil // key just dont exists
		default:
			return nil, "", err // another error
		}
	}

	if err := json.Unmarshal([]byte(marhalledPayload), &payload); err != nil {
		return nil, "", cstmerr.New(http.StatusInternalServerError, cstmerr.ErrUnknown.Error(), err)
	}

	return payload, marhalledPayload, nil
}

func (r *RedisTokenRepo) UpdateToken(c context.Context, newPayload *models.Payload, username string, ttl time.Duration) error {
	marshalled, err := json.Marshal(newPayload)
	if err != nil {
		return cstmerr.New(http.StatusInternalServerError, cstmerr.ErrUnknown.Error(), err)
	}

	return r.store.Set(c, username, marshalled, ttl).Err()
}
