package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"

	"github.com/myacey/jxgercorp-banking/services/libs/web"
	"github.com/myacey/jxgercorp-banking/services/user/internal/pkg/grpcclient"
	"github.com/myacey/jxgercorp-banking/services/user/internal/pkg/kafka"
	"github.com/myacey/jxgercorp-banking/services/user/internal/repository"
)

type AppConfig struct {
	PostgresCfg repository.PostgresConfig `mapstructure:"postgres"`
	RedisCfg    repository.RedisConfig    `mapstructure:"redis"`

	KafkaCfg      kafka.Config      `mapstructure:"kafka"`
	GrpcConfig    grpcclient.Config `mapstructure:"grpcclient"`
	HTTPServerCfg web.ServerConfig  `mapstructure:"httpserver"`
}

func LoadConfig(cfgPath string) (config AppConfig, err error) {
	_ = godotenv.Load("../../.env")

	viper.SetConfigFile(cfgPath)
	if err = viper.ReadInConfig(); err != nil {
		return
	}
	if err = viper.Unmarshal(&config); err != nil {
		return
	}

	config.PostgresCfg.Password = os.Getenv("TRANSFER_POSTGRES_PASSWORD")
	config.RedisCfg.Password = os.Getenv("USER_REDIS_PASSWORD")

	log.Println("config:", config)

	return
}
