package infrastructure

import (
	"authorization/domain"
	"authorization/infrastructure/mailer"
	"authorization/infrastructure/persistence"
	"authorization/service"
	"authorization/service/handlers"

	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

type Bootstraps struct {
	Bus       *service.MessageBus
	mailer    mailer.MailerInterface
	Endpoints map[string]domain.Endpoint
}

func NewBootstraps() (*asynq.Client, *Bootstraps) {
	err := persistence.CreateDBConnection()
	if err != nil {
		log.Fatal().Caller().Err(err).Msg("Failed to connect to database")
	}

	persistence.ConnectRedis()

	persistence.Migration(nil)

	uow, err := service.NewUnitOfWork(persistence.Pool)
	if err != nil {
		log.Fatal().Caller().Err(err).Msg("Failed to create unit of work")
	}

	asynqClient := mailer.CreateAsynqClient()

	mailer := mailer.NewMailer(asynqClient)
	if err != nil {
		log.Fatal().Caller().Err(err).Msg("Failed to create mailer")
	}

	messagebus := service.NewMessageBus(handlers.COMMAND_HANDLERS, uow, mailer)
	endpoints := make(map[string]domain.Endpoint)

	return asynqClient, &Bootstraps{
		Bus:       messagebus,
		mailer:    mailer,
		Endpoints: endpoints,
	}
}

func (bootstrap *Bootstraps) BootstrapMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set("mailer", bootstrap.mailer)
		ctx.Set("bus", bootstrap.Bus)
		ctx.Next()
	}
}
