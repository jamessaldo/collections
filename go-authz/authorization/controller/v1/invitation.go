package v1

import (
	"auth/domain/command"
	"auth/domain/model"
	"auth/middleware"
	"auth/service"
	"auth/view"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

// InvitationController : represent the invitation's controller contract
type InvitationController interface {
	VerifyInvitation(*gin.Context)
	GetInvitationByID(*gin.Context)
	DeleteInvitation(*gin.Context)
	Routes(*gin.RouterGroup)
}

type invitationController struct{}

// NewInvitationController -> returns new invitation controller
func NewInvitationController() InvitationController {
	return &invitationController{}
}

func (ctrl *invitationController) Routes(route *gin.RouterGroup) {
	invitation := route.Group("/invitations")
	invitation.POST("/verify", middleware.DeserializeUser(), ctrl.VerifyInvitation)
	invitation.GET("/:id", ctrl.GetInvitationByID)
	invitation.DELETE("/:id", middleware.DeserializeUser(), ctrl.DeleteInvitation)
}

// @Summary Get invitation by ID
// @Schemes
// @Description Get invitation data by ID
// @Tags Team
// @Accept json
// @Produce json
// @Param id path string true "Team ID"
// @Success 200 {string} string "OK"
// @Router /invitations/verify [post]
func (ctrl *invitationController) VerifyInvitation(ctx *gin.Context) {
	bus := ctx.MustGet("bus").(*service.MessageBus)
	currentUser := ctx.MustGet("currentUser").(*model.User)

	// Parse the request body into a User struct
	var cmd command.UpdateInvitationStatus
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

	// Return invitation data
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "OK"})
}

// @Summary Get invitation by ID
// @Schemes
// @Description Get invitation data by ID
// @Tags Invitation
// @Accept json
// @Produce json
// @Param id path string true "Invitation ID"
// @Success 200 {object} dto.InvitationRetreivalSchema
// @Router /invitations/{id} [get]
func (ctrl *invitationController) GetInvitationByID(ctx *gin.Context) {
	uow := ctx.MustGet("uow").(*service.UnitOfWork)

	// Get invitation ID from request parameter
	id := ctx.Param("id")
	log.Debug("Get invitation data by ID = ", id)

	// Get invitation data from database
	invitation, err := view.Invitation(id, uow)
	if err != nil {
		log.Error(err)
		_ = ctx.Error(err)
		return
	}

	// Return invitation data
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"invitation": invitation}})
}

// @Summary Delete invitation by ID
// @Schemes
// @Description Delete invitation data by ID
// @Tags Invitation
// @Accept json
// @Produce json
// @Param id path string true "Invitation ID"
// @Success 200 {string} string "OK"
// @Router /invitations/{id} [delete]
func (ctrl *invitationController) DeleteInvitation(ctx *gin.Context) {
	bus := ctx.MustGet("bus").(*service.MessageBus)
	currentUser := ctx.MustGet("currentUser").(*model.User)

	// Get invitation ID from request parameter
	id := ctx.Param("id")
	log.Debug("Get invitation data by ID = ", id)

	var cmd command.DeleteInvitation

	cmd.InvitationID = id
	cmd.User = currentUser

	err := bus.Handle(&cmd)
	if err != nil {
		log.Error(err)
		_ = ctx.Error(err)
		return
	}

	// Return invitation data
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "OK"})
}
