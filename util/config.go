package util

import (
	"log"

	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("check_period", "60s")
	viper.SetDefault("page_size", 200)
	viper.SetConfigFile("config.toml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Panicln("Failed to read config:", err)
	}
}
