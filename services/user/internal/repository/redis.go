package repository

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	// Address     string `mapstructure:"adress"`
	Host        string `mapstructure:"host"`
	Port        string `mapstructure:"port"`
	Password    string `mapstructure:"password"`
	DBIndex     int    `mapstructure:"index"`
	PoolSize    int    `mapstructure:"pool_size"`
	MinIdleConn int    `mapstructure:"min_idle_conn"`
}

func ConfigureRedisClient(cfg RedisConfig) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Host + ":" + cfg.Port,
		Password: cfg.Password,
		DB:       cfg.DBIndex,

		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConn,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}

	return rdb, nil
}
