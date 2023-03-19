package config

import (
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
	RedisHost     string `mapstructure:"REDIS_HOST"`
	RedisPort     string `mapstructure:"REDIS_PORT"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`

	// Connection Pool
	MaxIdleConns    int   `mapstructure:"MAX_IDLE_CONNS"`
	MaxOpenConns    int   `mapstructure:"MAX_OPEN_CONNS"`
	ConnMaxLifetime int64 `mapstructure:"CONN_MAX_LIFETIME"`
}

func loadStorageConfig(path string) (config StorageConfiguration, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("storage")

	viper.SetDefault("MaxIdleConns", 10)
	viper.SetDefault("MaxOpenConns", 100)
	viper.SetDefault("ConnMaxLifetime", 30)

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
