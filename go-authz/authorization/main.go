package main

import (
	"authorization/config"
	"authorization/controller"
	"authorization/infrastructure/persistence"
	"authorization/infrastructure/seeder"
	"authorization/infrastructure/worker"
	"authorization/repository"
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

	persistence.ConnectDB()
	persistence.ConnectRedis()
	persistence.Migration(nil)
	defer persistence.Pool.Close()

	mailerClient := worker.CreateMailerClient()
	worker.CreateMailer(mailerClient)
	defer mailerClient.Close()

	repository.CreateRepositories()
	handleArgs(persistence.Pool)
	controller.CreateRouter()
}

func handleArgs(pool *pgxpool.Pool) {
	flag.Parse()
	args := flag.Args()

	if len(args) >= 1 {
		switch args[0] {
		case "seed":
			seeder.Execute(pool, args[1:]...)
			os.Exit(0)
		}
	}
}
