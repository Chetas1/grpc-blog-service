package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	GrpcServer Server
	GrpcClient Client
}

type Server struct {
	Host     string
	Port     int
	Protocol string
}

type Client struct {
	ServerAddress string
}

func LoadConfig() (config Config, err error) {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	if err = viper.ReadInConfig(); err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
