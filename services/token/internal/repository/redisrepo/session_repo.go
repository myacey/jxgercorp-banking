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

func (r *RedisTokenRepo) GetToken(c context.Context, username string) (string, error) {
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
	marshalled, err := json.Marshal(newPayload)
	if err != nil {
		return cstmerr.New(http.StatusInternalServerError, cstmerr.ErrUnknown.Error(), err)
	}

	return r.store.Set(c, username, marshalled, ttl).Err()
}
