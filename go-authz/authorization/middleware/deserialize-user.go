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
	"github.com/rs/zerolog/log"
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
		userId, err := persistence.RedisClient.Get(_ctx, token).Result()
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

		userBytes, err := persistence.RedisClient.Get(_ctx, util.UserCachePrefix+userId).Bytes()
		if err == nil {
			err = json.Unmarshal(userBytes, user)
			if err != nil {
				log.Error().Err(err).Msg("error unmarshalling user")
				ctx.Abort()
				return
			}
		} else if err == redis.Nil {
			userId, _ := uuid.FromString(userId)
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
			json, err := json.Marshal(user)
			if err != nil {
				log.Error().Err(err).Msg("error marshalling user")
				ctx.Abort()
				return
			}

			expiredDate := time.Now().Add(10 * time.Minute)
			errCache := persistence.RedisClient.Set(ctx, util.UserCachePrefix+userId.String(), json, time.Until(expiredDate)).Err()
			if errCache != nil {
				ctx.Abort()
				return
			}

			log.Info().Str("userId", user.ID.String()).Msg("ga cached!")
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
