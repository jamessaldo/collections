package infrastructure

import (
	"auth/domain/model"
	"auth/infrastructure/persistence"
	"auth/infrastructure/worker"
	"auth/service"
	"auth/service/handlers"
	"context"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type Bootstraps struct {
	Bus       *service.MessageBus
	mailer    worker.WorkerInterface
	cache     *bigcache.BigCache
	Endpoints map[string]*model.Endpoint
}

func NewBootstraps() *Bootstraps {
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

	uow, err := service.NewUnitOfWork(db)
	if err != nil {
		log.Fatal(err)
	}

	asynqClient := worker.CreateAsynqClient()
	defer asynqClient.Close()

	mailer := worker.NewMailer(asynqClient)
	if err != nil {
		log.Fatal(err)
	}

	messagebus := service.NewMessageBus(handlers.COMMAND_HANDLERS, uow, mailer)
	cache, _ := bigcache.New(context.Background(), bigcache.DefaultConfig(10*time.Minute))
	endpoints := make(map[string]*model.Endpoint)

	return &Bootstraps{
		Bus:       messagebus,
		mailer:    mailer,
		cache:     cache,
		Endpoints: endpoints,
	}
}

func (bootstrap *Bootstraps) BootstrapMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("mailer", bootstrap.mailer)
		ctx.Set("bus", bootstrap.Bus)
		ctx.Set("cache", bootstrap.cache)
		ctx.Next()
	}
}
