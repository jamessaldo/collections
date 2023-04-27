package v1

import (
	"auth/domain/command"
	"auth/domain/model"
	"auth/middleware"
	"auth/service"
	"auth/view"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// UserController : represent the user's controller contract
type UserController interface {
	GetMe(*gin.Context)
	GetUserById(*gin.Context)
	GetUsers(*gin.Context)
	Routes(*gin.RouterGroup)
}

type userController struct{}

// NewUserController -> returns new user controller
func NewUserController() UserController {
	return &userController{}
}

func (ctrl *userController) Routes(route *gin.RouterGroup) {
	user := route.Group("/users")
	user.GET("/me", middleware.DeserializeUser(), ctrl.GetMe)
	user.GET("/:id", ctrl.GetUserById)
	user.GET("", ctrl.GetUsers)
	user.PUT("", middleware.DeserializeUser(), ctrl.UpdateUser)
	user.DELETE("", middleware.DeserializeUser(), ctrl.DeleteUser)
	user.PUT("/avatar", middleware.DeserializeUser(), ctrl.UpdateUserAvatar)
	user.DELETE("/avatar", middleware.DeserializeUser(), ctrl.DeleteUserAvatar)
}

// @Summary Get current user
// @Schemes
// @Description Get current user data from context
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} dto.ProfileUser
// @Router /users/me [get]
func (ctrl *userController) GetMe(ctx *gin.Context) {
	log.Debug("Get current user data from context")
	currentUser := ctx.MustGet("currentUser").(*model.User)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"user": currentUser.ProfileUser()}})
}

// @Summary Get user by ID
// @Schemes
// @Description Get user data by ID
// @Tags User
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} dto.PublicUser
// @Router /users/{id} [get]
func (ctrl *userController) GetUserById(ctx *gin.Context) {
	uow := ctx.MustGet("uow").(*service.UnitOfWork)

	// Get user ID from request parameter
	id := ctx.Param("id")
	log.Debug("Get user data by ID = ", id)

	// Get user data from database
	user, err := view.User(uuid.FromStringOrNil(id), uow)
	if err != nil {
		log.Error(err)
		_ = ctx.Error(err)
		return
	}

	// Return user data
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"user": user}})
}

// @Summary Get all users
// @Schemes
// @Description Get all users data
// @Tags User
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Success 200 {object} dto.PublicUser
// @Router /users [get]
func (ctrl *userController) GetUsers(ctx *gin.Context) {
	log.Debug("Get all users data")
	uow := ctx.MustGet("uow").(*service.UnitOfWork)

	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil {
		log.Error(err)
		_ = ctx.Error(err)
	}

	pageSize, err := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
	if err != nil {
		log.Error(err)
		_ = ctx.Error(err)
	}

	if page < 1 {
		page = 1
	}

	if pageSize < 10 {
		pageSize = 10
	}

	// Get user data from database
	users, err := view.Users(uow, page, pageSize)
	if err != nil {
		log.Error(err)
		_ = ctx.Error(err)
		return
	}

	// Return user data
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": users})
}

// @Summary Update user
// @Schemes
// @Description Update user data
// @Tags User
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param body body command.UpdateUser true "User data"
// @Success 200 {string} string "OK"
// @Router /users [put]
func (ctrl *userController) UpdateUser(ctx *gin.Context) {
	log.Debug("Update user data")
	currentUser := ctx.MustGet("currentUser").(*model.User)
	bus := ctx.MustGet("bus").(*service.MessageBus)

	// Parse the request body into a User struct
	var cmd command.UpdateUser
	if err := ctx.ShouldBindJSON(&cmd); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd.User = currentUser

	err := bus.Handle(&cmd)
	if err != nil {
		log.Error(err)
		_ = ctx.Error(err)
		return
	}

	// Return user data
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "OK"})
}

// @Summary Delete user
// @Schemes
// @Description Delete user data
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {string} string "OK"
// @Router /users [delete]
func (ctrl *userController) DeleteUser(ctx *gin.Context) {
	log.Debug("Delete user data")
	bus := ctx.MustGet("bus").(*service.MessageBus)
	currentUser := ctx.MustGet("currentUser").(*model.User)

	cmd := command.DeleteUser{
		User: currentUser,
	}

	err := bus.Handle(&cmd)
	if err != nil {
		log.Error(err)
		_ = ctx.Error(err)
		return
	}

	// Return user data
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"message": "OK"}})
}

// @Summary Update user avatar
// @Schemes
// @Description Update user avatar
// @Tags User
// @Accept json
// @Produce json
// @Param avatar formData file true "User avatar"
// @Success 200 {string} string "OK"
// @Router /users/avatar [put]
func (ctrl *userController) UpdateUserAvatar(ctx *gin.Context) {
	log.Debug("Update user avatar")
	bus := ctx.MustGet("bus").(*service.MessageBus)
	currentUser := ctx.MustGet("currentUser").(*model.User)

	// Parse the request body into a User struct
	var cmd command.UpdateUserAvatar
	if err := ctx.ShouldBind(&cmd); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd.User = currentUser

	err := bus.Handle(&cmd)
	if err != nil {
		log.Error(err)
		_ = ctx.Error(err)
		return
	}

	// Return user data
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"message": "OK"}})
}

// @Summary Delete user avatar
// @Schemes
// @Description Delete user avatar
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {string} string "OK"
// @Router /users/avatar [delete]
func (ctrl *userController) DeleteUserAvatar(ctx *gin.Context) {
	log.Debug("Delete user avatar")
	bus := ctx.MustGet("bus").(*service.MessageBus)
	currentUser := ctx.MustGet("currentUser").(*model.User)

	cmd := command.DeleteUserAvatar{
		User: currentUser,
	}

	err := bus.Handle(&cmd)
	if err != nil {
		log.Error(err)
		_ = ctx.Error(err)
		return
	}

	// Return user data
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"message": "OK"}})
}
