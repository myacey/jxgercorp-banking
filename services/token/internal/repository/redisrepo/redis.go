package redisrepo

import (
	"context"
	"time"

	"github.com/myacey/jxgercorp-banking/shared/backconfig"
	"github.com/redis/go-redis/v9"
)

func ConfigureRedisClient(config *backconfig.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: config.DBPassword, // no password set
		DB:       0,

		PoolSize:        10,
		MinIdleConns:    5,
		ConnMaxIdleTime: 30 * time.Minute,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return rdb, nil
}
