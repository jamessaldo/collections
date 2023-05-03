package main

import (
	"authorization/config"
	"authorization/controller"
	"authorization/infrastructure"
	"authorization/infrastructure/persistence"
	"flag"
	"os"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
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
	log.SetReportCaller(true)
	if config.AppConfig.AppEnv == Production {
		log.SetFormatter(&log.JSONFormatter{})
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(log.DebugLevel)
	}
}

func main() {
	// create a directory for avatars if it doesn't exist
	if config.StorageConfig.StaticDriver == "local" {
		if err := os.MkdirAll(config.StorageConfig.StaticRoot+config.StorageConfig.StaticAvatarPath, 0755); err != nil {
			log.Fatal(err)
		}
	}

	bootstrap := infrastructure.NewBootstraps()
	handleArgs(bootstrap.Bus.UoW.DB)

	server := controller.Server{}
	server.InitializeApp(bootstrap.BootstrapMiddleware())
}

func handleArgs(db *gorm.DB) {
	flag.Parse()
	args := flag.Args()

	if len(args) >= 1 {
		switch args[0] {
		case "seed":
			persistence.Execute(db, args[1:]...)
			os.Exit(0)
		}
	}
}
