package util

import (
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("capture_interval", "10s")
	viper.SetDefault("check_interval", "60s")
	viper.SetDefault("buffer_size", 2000)
	viper.SetDefault("page_size", 200)
	viper.SetConfigFile("config.toml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
