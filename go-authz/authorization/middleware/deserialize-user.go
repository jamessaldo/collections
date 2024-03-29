package middleware

import (
	"authorization/domain"
	"authorization/repository"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"authorization/config"
	"authorization/controller/exception"
	"authorization/infrastructure/persistence"
	"authorization/util"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
)

func DeserializeUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
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

		userId, err := persistence.RedisClient.Get(ctx.Request.Context(), token).Result()
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

		var user domain.User

		userBytes, err := persistence.RedisClient.Get(ctx.Request.Context(), util.UserCachePrefix+userId).Bytes()
		if err == nil {
			err = json.Unmarshal(userBytes, &user)
			if err != nil {
				log.Error().Caller().Err(err).Msg("error unmarshalling user")
				ctx.Abort()
				return
			}
		} else if err == redis.Nil {
			userId, _ := uuid.FromString(userId)
			user, err = repository.User.Get(ctx.Request.Context(), userId)
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					_ = ctx.Error(exception.NewNotFoundException("the user belonging to this token no logger exists"))
				} else {
					_ = ctx.Error(err)
				}
				ctx.Abort()
				return
			}
			json, err := json.Marshal(user)
			if err != nil {
				log.Error().Caller().Err(err).Msg("error marshalling user")
				ctx.Abort()
				return
			}

			expiredDate := util.GetTimestampUTC().Add(10 * time.Minute)
			errCache := persistence.RedisClient.Set(ctx, util.UserCachePrefix+userId.String(), json, time.Until(expiredDate)).Err()
			if errCache != nil {
				ctx.Abort()
				return
			}

			log.Info().Caller().Str("userId", user.ID.String()).Msg("Cannot find user in cache, fetching from database")
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
