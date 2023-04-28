package main

import (
	"auth/config"
	"auth/controller"
	"auth/domain/model"
	"auth/infrastructure/persistence"
	"auth/infrastructure/worker"
	"auth/service"
	"auth/service/handlers"
	"context"
	"flag"
	"os"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/gin-gonic/gin"
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
	Production = "production"
	// Development = "development"
	// Local       = "local"
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

func Bootstrap(db *gorm.DB, mailer worker.WorkerInterface) gin.HandlerFunc {
	uow, err := service.NewUnitOfWork(db)
	if err != nil {
		log.Fatal(err)
	}
	messagebus := service.NewMessageBus(handlers.COMMAND_HANDLERS, uow, mailer)
	cache, _ := bigcache.New(context.Background(), bigcache.DefaultConfig(10*time.Minute))

	return func(ctx *gin.Context) {
		ctx.Set("uow", uow)
		ctx.Set("mailer", mailer)
		ctx.Set("bus", messagebus)
		ctx.Set("cache", cache)
		ctx.Next()
	}
}

func main() {
	db, err := persistence.CreateDBConnection()
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(
		&model.User{},
		&model.Endpoint{},
		&model.Role{},
		&model.Access{},
		&model.Membership{},
		&model.Team{},
		&model.Invitation{})
	if err != nil {
		log.Fatal(err)
	}

	handleArgs(db)

	asynqClient := worker.CreateAsynqClient()
	defer asynqClient.Close()

	mailer := worker.NewMailer(asynqClient)
	if err != nil {
		log.Fatal(err)
	}

	// create a directory for avatars if it doesn't exist
	if config.StorageConfig.StaticDriver == "local" {
		if err := os.MkdirAll(config.StorageConfig.StaticRoot+config.StorageConfig.StaticAvatarPath, 0755); err != nil {
			log.Fatal(err)
		}
	}

	bootstrap := Bootstrap(db, mailer)

	server := controller.Server{}
	server.InitializeApp(bootstrap)
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
