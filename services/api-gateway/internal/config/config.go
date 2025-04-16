package config

import (
	"log"

	"github.com/myacey/jxgercorp-banking/services/libs/web"
	"github.com/spf13/viper"
)

type AppConfig struct {
	GrpcTarget    string           `mapstructure:"grpc_target"`
	HTTPServerCfg web.ServerConfig `mapstructure:"httpserver"`
}

func LoadConfig(cfgPath string) (config AppConfig, err error) {
	viper.SetConfigFile(cfgPath)
	err = viper.ReadInConfig()
	viper.Unmarshal(&config)
	log.Println("config:", config)

	return
}
