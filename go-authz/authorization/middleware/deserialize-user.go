package middleware

import (
	"errors"
	"fmt"
	"strings"

	"auth/config"
	"auth/controller/exception"
	"auth/service"
	"auth/util"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func DeserializeUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		uow := ctx.MustGet("uow").(*service.UnitOfWork)

		var token string

		authorizationHeader := ctx.Request.Header.Get("Authorization")
		fields := strings.Fields(authorizationHeader)

		if len(fields) != 0 && fields[0] == "Bearer" {
			token = fields[1]
		} else {
			cookie, err := ctx.Cookie("token")
			if err != nil {
				_ = ctx.Error(exception.NewUnauthorizedException("token is required"))
				ctx.Abort()
				return
			}
			token = cookie
		}

		if token == "" {
			_ = ctx.Error(exception.NewUnauthorizedException("token is required"))
			ctx.Abort()
			return
		}

		sub, err := util.ValidateToken(token, config.AppConfig.JWTTokenSecret)
		if err != nil {
			_ = ctx.Error(exception.NewUnauthorizedException("token is invalid"))
			ctx.Abort()
			return
		}
		user_id, _ := uuid.FromString(fmt.Sprint(sub))
		user, userErr := uow.User.Get(user_id)

		if userErr != nil {
			if errors.Is(userErr, gorm.ErrRecordNotFound) {
				_ = ctx.Error(exception.NewNotFoundException("the user belonging to this token no logger exists"))
				ctx.Abort()
				return
			}
			_ = ctx.Error(userErr)
			ctx.Abort()
			return
		}

		if !user.IsActive {
			_ = ctx.Error(exception.NewUnauthorizedException("the user belonging to this token has been deactivated"))
			ctx.Abort()
			return
		}

		ctx.Set("currentUser", user)
		ctx.Next()
	}
}
