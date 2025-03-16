package backconfig

import "github.com/spf13/viper"

type Config struct {
	// 4ALL
	AppDomain  string `mapstructure:"APP_DOMAIN"`
	DBPassword string `mapstructure:"DB_PASSWORD"`

	// POSTGRES
	PostgresHost   string `mapstructure:"POSTGRES_HOST"`
	PostgresUser   string `mapstructure:"POSTGRES_USER"`
	PostgresDBName string `mapstructure:"POSTGRES_DB"`
	PostgresPort   string `mapstructure:"POSTGRES_PORT"`

	// REDIS
	RedisUser   string `mapstructure:"REDIS_USER"`
	RedisAdress string `mapstructure:"REDIS_ADDRESS"`

	// GOOGLE SMTP
	GoogleMailAdress  string `mapstructure:"GOOGLE_MAIL_ADRESS"`
	GoogleAppPassword string `mapstructure:"GOOGLE_APP_PASSWORD"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
