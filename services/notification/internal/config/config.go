package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/myacey/jxgercorp-banking/services/notification/internal/adapter/inbound/kafka"
	"github.com/myacey/jxgercorp-banking/services/notification/internal/adapter/outbound/smtp"
	"github.com/spf13/viper"
)

type AppConfig struct {
	KafkaConfig kafka.Config `mapstructure:"kafka"`
	SMTPConfig  smtp.Config  `mapstructure:"smtp"`
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

	config.SMTPConfig.Username = os.Getenv("STMP_USERNAME")
	config.SMTPConfig.Password = os.Getenv("SMTP_PASSWORD")

	log.Println("config:", config)

	return
}
