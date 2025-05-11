package config

import (
	"log"
	"os"
	"reflect"

	"github.com/go-viper/mapstructure/v2"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"

	"github.com/myacey/jxgercorp-banking/services/notification/internal/adapter/inbound/kafka"
	"github.com/myacey/jxgercorp-banking/services/notification/internal/adapter/outbound/smtp"
)

type AppConfig struct {
	KafkaConfig kafka.Config `mapstructure:"kafka"`
	SMTPConfig  smtp.Config  `mapstructure:"smtp"`
}

func LoadConfig(cfgPath string) (config AppConfig, err error) {
	_ = godotenv.Load("../../.env.private", "../../.env")

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
