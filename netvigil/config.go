package netvigil

import (
	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigFile("config.toml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
