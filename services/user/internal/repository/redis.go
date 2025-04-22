package repository

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Address     string `mapstructure:"adress"`
	Password    string // secret from env
	DBIndex     int    `mapstructure:"index"`
	PoolSize    int    `mapstructure:"pool_size"`
	MinIdleConn int    `mapstructure:"min_idle_conn"`
}

func ConfigureRedisClient(cfg RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password, // no password set
		DB:       cfg.DBIndex,

		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConn,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return rdb, nil
}
