package config

import (
	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

var MailerConfig MailerConfiguration

func init() {
	var err error
	MailerConfig, err = loadMailerConfig("./env")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}
}

type MailerConfiguration struct {
	// Mailer
	MailerHost     string `mapstructure:"MAILER_HOST"`
	MailerPort     int    `mapstructure:"MAILER_PORT"`
	MailerUsername string `mapstructure:"MAILER_USERNAME"`
	MailerPassword string `mapstructure:"MAILER_PASSWORD"`
}

func loadMailerConfig(path string) (config MailerConfiguration, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("mailer")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
