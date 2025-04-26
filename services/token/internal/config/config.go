package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/myacey/jxgercorp-banking/services/token/internal/pkg/grpcserver"
	"github.com/spf13/viper"
)

type AppConfig struct {
	GrpcConfig grpcserver.Config `mapstructure:"grpcserver"`

	JwtSecretKey string
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

	config.JwtSecretKey = os.Getenv("JWT_SECRET_KEY")

	log.Println("config:", config)

	return
}
