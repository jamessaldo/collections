package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"authorization/config"
	"authorization/controller/exception"
	"authorization/domain/model"
	"authorization/infrastructure/persistence"
	"authorization/service"
	"authorization/util"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func DeserializeUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bus := ctx.MustGet("bus").(*service.MessageBus)
		uow := bus.UoW

		var token string

		authorizationHeader := ctx.Request.Header.Get("Authorization")
		fields := strings.Fields(authorizationHeader)

		if len(fields) != 0 && fields[0] == "Bearer" {
			token = fields[1]
		} else {
			cookie, err := ctx.Cookie("access_token")
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

		_ctx := context.TODO()
		stringId, err := persistence.RedisClient.Get(_ctx, token).Result()
		if err == redis.Nil {
			_ = ctx.Error(exception.NewUnauthorizedException("Token is invalid or session has expired: " + err.Error()))
			ctx.Abort()
			return
		}

		_, err = util.ValidateToken(token, config.AppConfig.AccessTokenPublicKey)
		if err != nil {
			_ = ctx.Error(exception.NewUnauthorizedException(err.Error()))
			ctx.Abort()
			return
		}

		var user *model.User = &model.User{}

		userBytes, err := persistence.RedisClient.Get(_ctx, util.UserCachePrefix+stringId).Bytes()
		if err == nil {
			err = json.Unmarshal(userBytes, user)
		} else if err == redis.Nil {
			userId, _ := uuid.FromString(stringId)
			user, err = uow.User.Get(userId)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					_ = ctx.Error(exception.NewNotFoundException("the user belonging to this token no logger exists"))
				} else {
					_ = ctx.Error(err)
				}
				ctx.Abort()
				return
			}
			json, _ := json.Marshal(user)

			expiredDate := time.Now().Add(10 * time.Minute)
			errCache := persistence.RedisClient.Set(ctx, util.UserCachePrefix+stringId, json, time.Until(expiredDate)).Err()
			if errCache != nil {
				ctx.Abort()
				return
			}
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
