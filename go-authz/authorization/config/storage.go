package config

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
)

var StorageConfig StorageConfiguration

func init() {
	var err error
	StorageConfig, err = loadStorageConfig("./env")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}
}

type StorageConfiguration struct {
	// Database
	DBDriver   string `mapstructure:"DB_DRIVER"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`

	// Redis
	RedisHost string `mapstructure:"REDIS_HOST"`
	RedisPort string `mapstructure:"REDIS_PORT"`

	// Connection Pool
	MaxIdleConns    int   `mapstructure:"MAX_IDLE_CONNS"`
	MaxOpenConns    int   `mapstructure:"MAX_OPEN_CONNS"`
	ConnMaxLifetime int64 `mapstructure:"CONN_MAX_LIFETIME"`

	// Static Storage
	StaticDriver        string `mapstructure:"STATIC_DRIVER"`
	StaticRoot          string `mapstructure:"STATIC_ROOT"`
	StaticPublicURL     string `mapstructure:"STATIC_PUBLIC_URL"`
	StaticAvatarPath    string `mapstructure:"STATIC_AVATAR_PATH"`
	StaticMaxAvatarSize int64  `mapstructure:"STATIC_MAX_AVATAR_SIZE"`
}

func loadStorageConfig(path string) (config StorageConfiguration, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("storage")

	viper.SetDefault("MAX_IDLE_CONNS", 10)
	viper.SetDefault("MAX_OPEN_CONNS", 100)
	viper.SetDefault("CONN_MAX_LIFETIME", 30)

	viper.SetDefault("STATIC_DRIVER", "local")
	viper.SetDefault("STATIC_ROOT", "static/")
	viper.SetDefault("STATIC_PUBLIC_URL", fmt.Sprintf("http://localhost:%s/static/", AppConfig.AppPort))
	viper.SetDefault("STATIC_AVATAR_PATH", "avatar/")
	viper.SetDefault("STATIC_MAX_AVATAR_SIZE", 1024*1024*2)

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config, DecoderErrorUnset)
	return
}
