package v1

import (
	"authorization/config"
	"authorization/controller/exception"
	"authorization/domain/command"
	"authorization/middleware"
	"authorization/service"
	"authorization/view"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type AuthController interface {
	LoginByGoogle(*gin.Context)
	Logout(*gin.Context)
	Routes(*gin.RouterGroup)
}

type authController struct{}

// NewAuthController -> returns new auth controller
func NewAuthController() AuthController {
	return &authController{}
}

func (ctrl *authController) Routes(route *gin.RouterGroup) {
	auth := route.Group("/auth")
	auth.GET("/logout", middleware.DeserializeUser(), ctrl.Logout)
	auth.GET("/refresh", ctrl.RefreshAccessToken)
	auth.GET("/sessions/oauth/google", ctrl.LoginByGoogle)
}

func (ctrl *authController) LoginByGoogle(ctx *gin.Context) {
	bus := ctx.MustGet("bus").(*service.MessageBus)
	code := ctx.Query("code")
	var pathUrl string = "/"

	if ctx.Query("state") != "" {
		pathUrl = ctx.Query("state")
	}

	if code == "" {
		err := exception.NewBadGatewayException("authorization code not provided")
		log.Error().Caller().Err(err).Msg("authorization code not provided")
		_ = ctx.Error(err)
		return
	}

	cmd := command.LoginByGoogle{
		Code:    code,
		PathURL: pathUrl,
	}

	err := bus.Handle(ctx.Request.Context(), &cmd)
	if err != nil {
		log.Error().Caller().Err(err).Msg("could not login by google")
		_ = ctx.Error(err)
		return
	}

	ctx.SetCookie("access_token", cmd.Token, config.AppConfig.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", cmd.RefreshToken, config.AppConfig.RefreshTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", config.AppConfig.AccessTokenMaxAge*60, "/", "localhost", false, false)
	ctx.Redirect(http.StatusTemporaryRedirect, fmt.Sprint(config.AppConfig.FrontEndOrigin, cmd.PathURL))
}

func (ctrl *authController) RefreshAccessToken(ctx *gin.Context) {
	bus := ctx.MustGet("bus").(*service.MessageBus)
	uow := bus.UoW

	message := "could not refresh access token"

	refresh_token, err := ctx.Cookie("refresh_token")

	if err != nil {
		log.Error().Caller().Err(err).Msg("refresh token is required")
		_ = ctx.Error(exception.NewUnauthorizedException("refresh token is required: " + err.Error()))
		return
	}

	accessToken, err := view.RefreshAccessToken(ctx.Request.Context(), refresh_token, uow)
	if err != nil {
		log.Error().Caller().Err(err).Msg("could not refresh access token")
		_ = ctx.Error(exception.NewUnauthorizedException(message))
		return
	}

	ctx.SetCookie("access_token", accessToken, config.AppConfig.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", config.AppConfig.AccessTokenMaxAge*60, "/", "localhost", false, false)
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"access_token": accessToken}})
}

func (ctrl *authController) Logout(ctx *gin.Context) {
	ctx.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "", -1, "/", "localhost", false, false)
	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
