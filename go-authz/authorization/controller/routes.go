package controller

import (
	"authorization/config"
	v1 "authorization/controller/v1"
	"authorization/middleware"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin"

	docs "authorization/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func CreateRouter() {
	userControllerV1 := v1.NewUserController()
	authControllerV1 := v1.NewAuthController()
	teamControllerV1 := v1.NewTeamController()
	invitationControllerV1 := v1.NewInvitationController()

	docs.SwaggerInfo.BasePath = "/api/v1"

	router := gin.Default()
	router.MaxMultipartMemory = config.StorageConfig.StaticMaxAvatarSize
	router.Use(middleware.CORSMiddleware()) //For CORS
	router.Use(middleware.HandleCustomError())
	router.NoRoute(noRouteHandler)
	router.StaticFS("/static", http.Dir("./static"))

	routerApi := router.Group("/api")
	routerV1 := routerApi.Group("/v1")

	//authentication routes
	authControllerV1.Routes(routerApi)

	//invitation routes
	invitationControllerV1.Routes(routerV1)

	//team routes
	teamControllerV1.Routes(routerV1)

	//user routes
	userControllerV1.Routes(routerV1)

	routerV1.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//Starting the application
	log.Fatal().Caller().Err(router.Run(config.AppConfig.AppHost + ":" + config.AppConfig.AppPort)).Msg("Cannot start the server")
}

func noRouteHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Route Not Found"})
}
