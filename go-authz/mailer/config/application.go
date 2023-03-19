package config

import (
	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

var AppConfig ApplicationConfiguration

func init() {
	var err error
	AppConfig, err = loadAppConfig("./env")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}
}

type ApplicationConfiguration struct {
	AppEnv         string `mapstructure:"APP_ENV"`
	AppPort        string `mapstructure:"APP_PORT"`
	FrontEndOrigin string `mapstructure:"FRONTEND_ORIGIN"`
}

func loadAppConfig(path string) (config ApplicationConfiguration, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")
	viper.SetDefault("APP_PORT", "8888")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
