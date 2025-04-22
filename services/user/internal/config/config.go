package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"

	"github.com/myacey/jxgercorp-banking/services/libs/web"
	"github.com/myacey/jxgercorp-banking/services/user/internal/repository"
	"github.com/myacey/jxgercorp-banking/services/user/internal/service"
)

type AppConfig struct {
	PostgresCfg repository.PostgresConfig `mapstructure:"postgres"`
	RedisCfg    repository.RedisConfig    `mapstructure:"redis"`

	KafkaCfg      service.ConfirmationKafkaConfig `mapstructure:"kafka"`
	GrpcTarget    string                          `mapstructure:"grpc_target"`
	HTTPServerCfg web.ServerConfig                `mapstructure:"httpserver"`
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

	config.PostgresCfg.Password = os.Getenv("POSTGRES_PASSWORD")
	config.RedisCfg.Password = os.Getenv("REDIS_PASSWORD")

	// viper.SetConfigFile(".env")
	// viper.ReadInConfig()
	// viper.AutomaticEnv()

	// viper.SetConfigFile(cfgPath)
	// viper.MergeInConfig()

	// viper.Unmarshal(&config)

	log.Println("config:", config)

	return
}
