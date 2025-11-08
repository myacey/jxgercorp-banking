package config

import (
	"log"
	"os"
	"reflect"

	"github.com/go-viper/mapstructure/v2"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"

	"github.com/myacey/jxgercorp-banking/services/libs/web"
	"github.com/myacey/jxgercorp-banking/services/user/internal/httpserver/handler"
	"github.com/myacey/jxgercorp-banking/services/user/internal/pkg/grpcclient"
	"github.com/myacey/jxgercorp-banking/services/user/internal/pkg/kafka"
	"github.com/myacey/jxgercorp-banking/services/user/internal/repository"
	"github.com/myacey/jxgercorp-banking/services/user/internal/service"
)

type AppConfig struct {
	PostgresCfg repository.PostgresConfig `mapstructure:"postgres"`
	RedisCfg    repository.RedisConfig    `mapstructure:"redis"`

	KafkaCfg      kafka.Config      `mapstructure:"kafka"`
	GrpcConfig    grpcclient.Config `mapstructure:"grpcclient"`
	HTTPServerCfg web.ServerConfig  `mapstructure:"httpserver"`

	ConfirmationCfg service.ConfirmationConfig `mapstructure:"confirmation"`
	HandlerCfg      handler.Config             `mapstructure:"handler"`
}

func LoadConfig(cfgPath string) (config AppConfig, err error) {
	_ = godotenv.Load("../../.env")
	_ = godotenv.Load("../../.env.private")

	viper.SetConfigFile(cfgPath)
	if err = viper.ReadInConfig(); err != nil {
		return
	}

	hook := mapstructure.DecodeHookFunc(func(
		from reflect.Type, to reflect.Type, data interface{},
	) (interface{}, error) {
		if from.Kind() == reflect.String && to.Kind() == reflect.String {
			return os.ExpandEnv(data.(string)), nil
		}
		return data, nil
	})

	if err = viper.Unmarshal(&config, viper.DecodeHook(hook)); err != nil {
		return
	}

	log.Printf("config: %+v", config)

	return
}
