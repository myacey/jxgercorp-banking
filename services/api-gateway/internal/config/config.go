package config

import (
	"log"
	"os"
	"reflect"

	"github.com/go-viper/mapstructure/v2"
	"github.com/joho/godotenv"
	"github.com/myacey/jxgercorp-banking/services/api-gateway/internal/pkg/grpcclient"
	"github.com/myacey/jxgercorp-banking/services/libs/web"
	"github.com/spf13/viper"
)

type AppConfig struct {
	Services map[string]string `mapstructure:"services"`

	GrpcCfg       grpcclient.Config `mapstructure:"grpcclient"`
	HTTPServerCfg web.ServerConfig  `mapstructure:"httpserver"`
}

func LoadConfig(cfgPath string) (AppConfig, error) {
	godotenv.Load("../../.env")

	viper.SetConfigFile(cfgPath)
	if err := viper.ReadInConfig(); err != nil {
		return AppConfig{}, err
	}

	hook := mapstructure.DecodeHookFunc(func(
		from reflect.Type, to reflect.Type, data interface{},
	) (interface{}, error) {
		if from.Kind() == reflect.String && to.Kind() == reflect.String {
			return os.ExpandEnv(data.(string)), nil
		}
		return data, nil
	})

	var cfg AppConfig
	if err := viper.Unmarshal(&cfg, viper.DecodeHook(hook)); err != nil {
		return AppConfig{}, err
	}

	log.Printf("config: %+v\n", cfg)
	return cfg, nil
}
