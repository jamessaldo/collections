package main

import (
	"authorization/config"
	"authorization/controller"
	"authorization/infrastructure"
	"authorization/infrastructure/persistence"
	"flag"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// @title           Authorization API
// @version         1.0
// @description     This is Authorization API documentation.

// @contact.name   Jamessaldo
// @contact.url    https://github.com/jamessaldo/collections/issues
// @contact.email  ghozyghlmlaff@gmail.com

// @host      localhost:8888
// @BasePath  /api/v1

// create an enum of the environment
const (
	Production  = "production"
	Development = "development"
	Local       = "local"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	if config.AppConfig.AppEnv == Production {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}

func main() {
	// create a directory for avatars if it doesn't exist
	if config.StorageConfig.StaticDriver == "local" {
		if err := os.MkdirAll(config.StorageConfig.StaticRoot+config.StorageConfig.StaticAvatarPath, 0755); err != nil {
			log.Fatal().Caller().Err(err).Msg("Cannot start the server, reason: cannot create avatar directory")
		}
	}

	asynqClient, bootstrap := infrastructure.NewBootstraps()
	defer asynqClient.Close()

	handleArgs(bootstrap.Bus.UoW.GetDB())

	server := controller.Server{}
	server.InitializeApp(bootstrap.BootstrapMiddleware())
}

func handleArgs(pool *pgxpool.Pool) {
	flag.Parse()
	args := flag.Args()

	if len(args) >= 1 {
		switch args[0] {
		case "seed":
			persistence.Execute(pool, args[1:]...)
			os.Exit(0)
		}
	}
}
