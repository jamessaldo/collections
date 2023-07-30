package v1

import (
	"authorization/domain"
	"authorization/domain/command"
	"authorization/domain/dto"
	"authorization/infrastructure/persistence"
	"authorization/middleware"
	"authorization/service"
	"authorization/util"
	"authorization/view"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/rs/zerolog/log"

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
	currentUser := ctx.MustGet("currentUser").(domain.User)
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
	bus := ctx.MustGet("bus").(*service.MessageBus)
	uow := bus.UoW

	// Get user ID from request parameter
	id := ctx.Param("id")
	log.Debug().Caller().Str("id", id).Msg("Get user data by ID")

	var user *dto.PublicUser = &dto.PublicUser{}

	userBytes, err := persistence.RedisClient.Get(ctx.Request.Context(), util.UserCachePrefix+id).Bytes()
	if err == nil {
		var userCache domain.User = domain.User{}
		err = json.Unmarshal(userBytes, &userCache)
		if err == nil {
			user = userCache.PublicUser()
		}
	}

	if err != nil {
		// Get user data from database
		user, err = view.User(ctx.Request.Context(), uuid.FromStringOrNil(id), uow)
		if err != nil {
			log.Error().Caller().Err(err).Msg("Failed to get user data")
			_ = ctx.Error(err)
			return
		}
	}

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
	bus := ctx.MustGet("bus").(*service.MessageBus)
	uow := bus.UoW

	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil {
		log.Error().Caller().Err(err).Msg("Failed to get page number")
		_ = ctx.Error(err)
	}

	pageSize, err := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
	if err != nil {
		log.Error().Caller().Err(err).Msg("Failed to get page size")
		_ = ctx.Error(err)
	}

	if page < 1 {
		page = 1
	}

	if pageSize < 10 {
		pageSize = 10
	}

	// Get user data from database
	users, err := view.Users(ctx.Request.Context(), uow, page, pageSize)
	if err != nil {
		log.Error().Caller().Err(err).Msg("Failed to get users data")
		_ = ctx.Error(err)
		return
	}

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
	currentUser := ctx.MustGet("currentUser").(domain.User)
	bus := ctx.MustGet("bus").(*service.MessageBus)

	// Parse the request body into a User struct
	var cmd command.UpdateUser
	if err := ctx.ShouldBindJSON(&cmd); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd.User = currentUser

	err := bus.Handle(ctx.Request.Context(), &cmd)
	if err != nil {
		log.Error().Caller().Err(err).Msg("Failed to update user data")
		_ = ctx.Error(err)
		return
	}

	persistence.RedisClient.Del(ctx.Request.Context(), util.UserCachePrefix+currentUser.ID.String())
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
	bus := ctx.MustGet("bus").(*service.MessageBus)
	currentUser := ctx.MustGet("currentUser").(domain.User)

	cmd := command.DeleteUser{
		User: currentUser,
	}

	err := bus.Handle(ctx.Request.Context(), &cmd)
	if err != nil {
		log.Error().Caller().Err(err).Msg("Failed to delete user data")
		_ = ctx.Error(err)
		return
	}

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
	bus := ctx.MustGet("bus").(*service.MessageBus)
	currentUser := ctx.MustGet("currentUser").(domain.User)

	// Parse the request body into a User struct
	var cmd command.UpdateUserAvatar
	if err := ctx.ShouldBind(&cmd); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd.User = currentUser

	err := bus.Handle(ctx.Request.Context(), &cmd)
	if err != nil {
		log.Error().Caller().Err(err).Msg("Failed to update user avatar")
		_ = ctx.Error(err)
		return
	}

	persistence.RedisClient.Del(ctx.Request.Context(), util.UserCachePrefix+currentUser.ID.String())
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
	bus := ctx.MustGet("bus").(*service.MessageBus)
	currentUser := ctx.MustGet("currentUser").(domain.User)

	cmd := command.DeleteUserAvatar{
		User: currentUser,
	}

	err := bus.Handle(ctx.Request.Context(), &cmd)
	if err != nil {
		log.Error().Caller().Err(err).Msg("Failed to delete user avatar")
		_ = ctx.Error(err)
		return
	}

	persistence.RedisClient.Del(ctx.Request.Context(), util.UserCachePrefix+currentUser.ID.String())
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"message": "OK"}})
}
