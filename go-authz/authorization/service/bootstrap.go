package service

import (
	"auth/infrastructure/worker"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Bootstrap(db *gorm.DB, mailer worker.WorkerInterface) gin.HandlerFunc {

	return func(ctx *gin.Context) {
		uow, err := NewUnitOfWork(db)
		if err != nil {
			log.Fatal(err)
		}

		ctx.Set("uow", uow)
		ctx.Set("mailer", mailer)
		ctx.Next()
	}
}
