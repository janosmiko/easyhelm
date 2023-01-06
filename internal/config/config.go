package config

import (
	"strings"

	"github.com/spf13/viper"
)

// Init initializes Viper.
func Init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	_ = viper.ReadInConfig()
	viper.AutomaticEnv()

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
}

func Unmarshal() map[string]interface{} {
	var c map[string]interface{}
	_ = viper.Unmarshal(&c)

	return c
}
