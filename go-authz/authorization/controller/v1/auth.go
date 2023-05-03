package v1

import (
	"authorization/config"
	"authorization/controller/exception"
	"authorization/domain/command"
	"authorization/middleware"
	"authorization/service"
	"authorization/util"
	"authorization/view"
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
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
	auth.GET("/sessions/oauth/google", ctrl.LoginByGoogle)
}

func (ctrl *authController) LoginByGoogle(ctx *gin.Context) {
	bus := ctx.MustGet("bus").(*service.MessageBus)
	uow := bus.UoW
	code := ctx.Query("code")
	var pathUrl string = "/"

	if ctx.Query("state") != "" {
		pathUrl = ctx.Query("state")
	}

	if code == "" {
		err := exception.NewBadGatewayException("authorization code not provided")
		log.Error(err)
		_ = ctx.Error(err)
		return
	}

	tokenRes, err := util.GetGoogleOauthToken(context.Background(), code)
	if err != nil {
		err = exception.NewBadGatewayException(err.Error())
		log.Error(err)
		_ = ctx.Error(err)
		return
	}

	googleUser, err := util.GetGoogleUser(context.Background(), tokenRes.Access_token, tokenRes.Id_token)
	if err != nil {
		err = exception.NewBadGatewayException(err.Error())
		log.Error(err)
		_ = ctx.Error(err)
		return
	}

	cmd := command.LoginByGoogle{
		Code:       code,
		PathURL:    pathUrl,
		GoogleUser: *googleUser,
	}

	err = bus.Handle(&cmd)
	if err != nil {
		log.Error(err)
		_ = ctx.Error(err)
		return
	}

	token, refreshToken, err := view.LoginByGoogle(googleUser.Email, uow)
	if err != nil {
		log.Error(err)
		_ = ctx.Error(err)
		return
	}

	ctx.SetCookie("token", token, config.AppConfig.TokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("refreshToken", refreshToken, config.AppConfig.RefreshTokenMaxAge*60, "/", "localhost", false, true)
	ctx.Redirect(http.StatusTemporaryRedirect, fmt.Sprint(config.AppConfig.FrontEndOrigin, cmd.PathURL))
}

func (ctrl *authController) Logout(ctx *gin.Context) {
	ctx.SetCookie("token", "", -1, "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
