package config

import (
	"todo-app/app/types"

	"github.com/spf13/viper"
)

func LoadConfig(path string) (config types.Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	if err = viper.ReadInConfig(); err != nil {
		return types.Config{}, err
	}

	if err = viper.Unmarshal(&config); err != nil {
		return types.Config{}, err
	}

	return config, err
}
