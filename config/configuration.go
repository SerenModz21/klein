package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Configuration struct {
	Port  string
	Mongo string
	Redis string
}

func GetConfig() Configuration {
	configuration := Configuration{}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	if error := viper.ReadInConfig(); error != nil {
		log.Panic(error)
	}

	if error := viper.Unmarshal(&configuration); error != nil {
		log.Panic(error)
	}

	return configuration
}
