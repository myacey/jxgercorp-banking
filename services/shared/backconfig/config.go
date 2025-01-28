package backconfig

import "github.com/spf13/viper"

type Config struct {
	DBPassword string `mapstructure:"DB_PASSWORD"`

	PostgresHost   string `mapstructure:"POSTGRES_HOST"`
	PostgresUser   string `mapstructure:"POSTGRES_USER"`
	PostgresDBName string `mapstructure:"POSTGRES_DB"`
	PostgresPort   string `mapstructure:"POSTGRES_PORT"`

	RedisUser   string `mapstructure:"REDIS_USER"`
	RedisAdress string `mapstructure:"REDIS_ADDRESS"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	// viper.SetConfigName("")
	// viper.SetConfigType("env")

	viper.SetConfigFile(".env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
