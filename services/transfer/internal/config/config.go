package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/myacey/jxgercorp-banking/services/libs/web"
	"github.com/myacey/jxgercorp-banking/services/transfer/internal/repository"
	"github.com/spf13/viper"
)

type AppConfig struct {
	PostgresCfg repository.PostgresConfig `mapstructure:"postgres"`

	HTTPServerCfg web.ServerConfig `mapstructure:"httpserver"`
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

	log.Println("config:", config)

	return
}
