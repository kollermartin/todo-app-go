package config

import (
	"todo-app/app/constant"
	"todo-app/app/types"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func LoadConfig(path string) (config *types.Config) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		logrus.WithFields(logrus.Fields{
			"event": constant.ConfigLoadLogEventErrorKey,
			"error": err.Error(),
		}).Fatal("Failed to read config file")
	}

	if err := viper.Unmarshal(&config); err != nil {
		logrus.WithFields(logrus.Fields{
			"event": constant.ConfigLoadLogEventErrorKey,
			"error": err.Error(),
		}).Fatal("Failed to unmarshal config file")
	}

	return config
}
