package config

import (
	"github.com/spf13/viper"
)

var c *viper.Viper

// Init initialize config
func Init(env string) {
	c = viper.New()
	c.SetConfigType("yaml")
	c.SetConfigName(env)
	c.AddConfigPath("config/environments/")
	c.AddConfigPath("../config/environments/")
	c.AddConfigPath("/run/secrets/")
	if err := c.ReadInConfig(); err != nil {
		panic(err)
	}
}

// GetConfig gets config
func GetConfig() *viper.Viper {
	return c
}
