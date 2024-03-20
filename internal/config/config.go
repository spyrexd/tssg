package config

import "github.com/spf13/viper"

func LoadEnvConfig() {
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
}

func Get(key string) any {
	return viper.Get(key)
}
