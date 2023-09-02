package utils

import (
	"github.com/spf13/viper"
)

func MustLoad(configFile *string, c any) {
	viper.SetConfigFile(*configFile)
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(c)
	if err != nil {
		panic(err)
	}
}
