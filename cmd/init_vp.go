package main

import "github.com/spf13/viper"

func initConfig() error {
	viper.AddConfigPath("source/configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
