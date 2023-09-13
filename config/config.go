package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	MongoURL   string `mapstructure:"MONGO_URL"`
	JwtSecret  string `mapstructure:"JWT_SECRET"`
	ServerPort string `mapstructure:"SERVER_PORT"`
}

var Cfg Config

func LoadConfig() error {
	viper.AddConfigPath("./")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	err := viper.ReadInConfig()
	viper.AutomaticEnv()
	if err != nil {
		return err
	}
	err = viper.Unmarshal(&Cfg)
	return err
}
