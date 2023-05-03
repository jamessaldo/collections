package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"authorization/config"
	"authorization/controller/exception"
	"authorization/domain/model"
	"authorization/service"
	"authorization/util"

	"github.com/allegro/bigcache/v3"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func DeserializeUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		bus := ctx.MustGet("bus").(*service.MessageBus)
		uow := bus.UoW
		cache := ctx.MustGet("cache").(*bigcache.BigCache)

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

		var user *model.User = &model.User{}
		stringId := fmt.Sprint(sub)

		userBytes, err := cache.Get(util.UserCachePrefix + stringId)
		if err == nil {
			err = json.Unmarshal(userBytes, user)
		}
		if err != nil {
			user_id, _ := uuid.FromString(stringId)
			user, err = uow.User.Get(user_id)
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
			cache.Set(util.UserCachePrefix+stringId, json)
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
