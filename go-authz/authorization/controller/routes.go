package controller

import (
	"authorization/config"
	v1 "authorization/controller/v1"
	"authorization/middleware"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"

	docs "authorization/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct{}

func (server *Server) InitializeApp(bootstrap gin.HandlerFunc) {
	userControllerV1 := v1.NewUserController()
	authControllerV1 := v1.NewAuthController()
	teamControllerV1 := v1.NewTeamController()

	docs.SwaggerInfo.BasePath = "/api/v1"

	router := gin.Default()
	router.MaxMultipartMemory = config.StorageConfig.StaticMaxAvatarSize
	router.Use(middleware.CORSMiddleware()) //For CORS
	router.Use(middleware.HandleCustomError())
	router.Use(bootstrap)
	router.NoRoute(NoRoute)
	router.StaticFS("/static", http.Dir("./static"))

	routerApi := router.Group("/api")
	routerV1 := routerApi.Group("/v1")

	//user routes
	userControllerV1.Routes(routerV1)

	//authentication routes
	authControllerV1.Routes(routerApi)

	//team routes
	teamControllerV1.Routes(routerV1)

	routerV1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//Starting the application
	log.Fatal(router.Run(config.AppConfig.AppHost + ":" + config.AppConfig.AppPort))
}

func NoRoute(ctx *gin.Context) {
	ctx.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "Route Not Found"})
}
