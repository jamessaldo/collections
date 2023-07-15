package v1

import (
	"authorization/controller/exception"
	"authorization/domain/command"
	"authorization/domain/model"
	"authorization/middleware"
	"authorization/service"
	"authorization/view"
	"net/http"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"

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
	invitation.GET("/:id/check", ctrl.GetInvitationByID)
	invitation.DELETE("/:id", middleware.DeserializeUser(), ctrl.DeleteInvitation)
}

// @Summary Verify invitation by ID
// @Schemes
// @Description Verify invitation data by ID
// @Tags Membership
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
		log.Error().Err(err).Msg("could not verify invitation")
		_ = ctx.Error(err)
		return
	}

	// Return invitation data
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "OK"})
}

// @Summary Get invitation by ID
// @Schemes
// @Description Get invitation data by ID
// @Tags Membership
// @Accept json
// @Produce json
// @Param id path string true "Invitation ID"
// @Success 200 {object} dto.InvitationRetreivalSchema
// @Router /invitations/{id}/check [get]
func (ctrl *invitationController) GetInvitationByID(ctx *gin.Context) {
	bus := ctx.MustGet("bus").(*service.MessageBus)
	uow := bus.UoW

	// Get invitation ID from request parameter
	id := ctx.Param("id")
	log.Debug().Str("id", id).Msg("Get invitation data by ID")

	idParsed, err := ulid.Parse(id)

	if err != nil {
		badRequest := exception.NewBadGatewayException(err.Error())
		_ = ctx.Error(badRequest)
		return
	}

	// Get invitation data from database
	invitation, err := view.Invitation(idParsed, uow)
	if err != nil {
		log.Error().Err(err).Msg("could not get invitation")
		_ = ctx.Error(err)
		return
	}

	// Return invitation data
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"invitation": invitation}})
}

// @Summary Delete invitation by ID
// @Schemes
// @Description Delete invitation data by ID
// @Tags Membership
// @Accept json
// @Produce json
// @Param id path string true "Invitation ID"
// @Success 200 {string} string "OK"
// @Router /invitations/{id} [delete]
func (ctrl *invitationController) DeleteInvitation(ctx *gin.Context) {
	bus := ctx.MustGet("bus").(*service.MessageBus)
	currentUser := ctx.MustGet("currentUser").(*model.User)

	// Get invitation ID from request parameter
	invitationIDString := ctx.Param("id")
	log.Debug().Str("id", invitationIDString).Msg("Delete invitation data by ID")

	var cmd command.DeleteInvitation

	invitationID, err := ulid.Parse(invitationIDString)

	if err != nil {
		badRequest := exception.NewBadGatewayException(err.Error())
		_ = ctx.Error(badRequest)
		return
	}

	cmd.InvitationID = invitationID
	cmd.User = currentUser

	err = bus.Handle(&cmd)
	if err != nil {
		log.Error().Err(err).Msg("could not delete invitation")
		_ = ctx.Error(err)
		return
	}

	// Return invitation data
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "OK"})
}
