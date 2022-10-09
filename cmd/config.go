package main

import (
	"github.com/spf13/viper"
)

type Config struct {
	DBUri             string `mapstructure:"MONGODB_LOCAL_URI"`
	Port              string `mapstructure:"PORT"`
	Origin            string `mapstructure:"CLIENT_ORIGIN"`
	GrpcServerAddress string `mapstructure:"GRPC_SERVER_ADDRESS"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("dev")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
