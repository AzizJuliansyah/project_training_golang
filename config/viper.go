package config

import (
	"log"

	"github.com/spf13/viper"
)

func InitViper() {
	viper.SetConfigName("app.conf.json")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln("Error viper configuration", err.Error())
	}
}