package config

import (
	"time"

	"github.com/mitchellh/mapstructure"
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

func DecoderErrorUnset(c *mapstructure.DecoderConfig) {
	c.ErrorUnset = true
}

type ApplicationConfiguration struct {
	FrontEndOrigin  string `mapstructure:"FRONTEND_ORIGIN"`
	AppHost         string `mapstructure:"APP_HOST"`
	AppPort         string `mapstructure:"APP_PORT"`
	AppExtAuthzPort string `mapstructure:"APP_EXT_AUTHZ_PORT"`
	AppEnv          string `mapstructure:"APP_ENV"`
	AppName         string `mapstructure:"APP_NAME"`

	// JWT
	JWTKid                 string        `mapstructure:"JWT_KID"`
	AccessTokenPrivateKey  string        `mapstructure:"ACCESS_TOKEN_PRIVATE_KEY"`
	AccessTokenPublicKey   string        `mapstructure:"ACCESS_TOKEN_PUBLIC_KEY"`
	RefreshTokenPrivateKey string        `mapstructure:"REFRESH_TOKEN_PRIVATE_KEY"`
	RefreshTokenPublicKey  string        `mapstructure:"REFRESH_TOKEN_PUBLIC_KEY"`
	AccessTokenExpiresIn   time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRED_IN"`
	RefreshTokenExpiresIn  time.Duration `mapstructure:"REFRESH_TOKEN_EXPIRED_IN"`
	AccessTokenMaxAge      int           `mapstructure:"ACCESS_TOKEN_MAXAGE"`
	RefreshTokenMaxAge     int           `mapstructure:"REFRESH_TOKEN_MAXAGE"`

	// Google OAuth
	GoogleClientID         string `mapstructure:"GOOGLE_OAUTH_CLIENT_ID"`
	GoogleClientSecret     string `mapstructure:"GOOGLE_OAUTH_CLIENT_SECRET"`
	GoogleOAuthRedirectUrl string `mapstructure:"GOOGLE_OAUTH_REDIRECT_URL"`
}

func loadAppConfig(path string) (config ApplicationConfiguration, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")
	viper.SetDefault("APP_PORT", "8888")
	viper.SetDefault("APP_EXT_AUTHZ_PORT", "8889")
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("APP_NAME", "svc-authorization")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config, DecoderErrorUnset)
	return
}
