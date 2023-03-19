package v1

import (
	"auth/config"
	"auth/domain/command"
	"auth/infrastructure/worker"
	"auth/middleware"
	"auth/service"
	"auth/service/handlers"
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
	return authController{}
}

func (ctrl authController) Routes(route *gin.RouterGroup) {
	auth := route.Group("/auth")
	auth.GET("/logout", middleware.DeserializeUser(), ctrl.Logout)
	auth.GET("/sessions/oauth/google", ctrl.LoginByGoogle)
}

func (ctrl authController) LoginByGoogle(ctx *gin.Context) {
	uow := ctx.MustGet("uow").(*service.UnitOfWork)
	mailer := ctx.MustGet("mailer").(worker.WorkerInterface)

	code := ctx.Query("code")
	var pathUrl string = "/"

	if ctx.Query("state") != "" {
		pathUrl = ctx.Query("state")
	}

	cmd := command.LoginByGoogle{
		Code:    code,
		PathURL: pathUrl,
	}

	token, refreshToken, err := handlers.LoginByGoogle(uow, mailer, cmd)
	if err != nil {
		log.Error(err)
		ctx.Error(err)
		return
	}

	ctx.SetCookie("token", token, config.AppConfig.TokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("refreshToken", refreshToken, config.AppConfig.RefreshTokenMaxAge*60, "/", "localhost", false, true)
	ctx.Redirect(http.StatusTemporaryRedirect, fmt.Sprint(config.AppConfig.FrontEndOrigin, cmd.PathURL))
}

func (a authController) Logout(ctx *gin.Context) {
	ctx.SetCookie("token", "", -1, "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
