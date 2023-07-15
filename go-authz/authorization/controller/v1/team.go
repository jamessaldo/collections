package v1

import (
	"authorization/controller/exception"
	"authorization/domain/command"
	"authorization/domain/model"
	"authorization/middleware"
	"authorization/service"
	"authorization/view"
	"net/http"
	"strconv"

	"github.com/oklog/ulid/v2"
	"github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

// TeamController : represent the team's controller contract
type TeamController interface {
	GetTeamById(*gin.Context)
	GetTeams(*gin.Context)
	CreateTeam(*gin.Context)
	UpdateTeam(*gin.Context)
	Routes(*gin.RouterGroup)
}

type teamController struct{}

// NewTeamController -> returns new team controller
func NewTeamController() TeamController {
	return &teamController{}
}

func (ctrl *teamController) Routes(route *gin.RouterGroup) {
	team := route.Group("/teams")
	team.GET("", middleware.DeserializeUser(), ctrl.GetTeams)
	team.GET("/:id", middleware.DeserializeUser(), ctrl.GetTeamById)
	team.POST("", middleware.DeserializeUser(), ctrl.CreateTeam)
	team.PUT("/:id", middleware.DeserializeUser(), ctrl.UpdateTeam)
	team.PUT("/:id/last-active", middleware.DeserializeUser(), ctrl.UpdateLastActiveTeam)
	team.DELETE("/:id/members/:membership_id", middleware.DeserializeUser(), ctrl.DeleteTeamMember)
	team.PUT("/:id/members/:membership_id", middleware.DeserializeUser(), ctrl.ChangeMemberRole)
	team.POST("/:id/invitation", middleware.DeserializeUser(), ctrl.SendInvitation)
	team.POST("/:id/invitation/:invitation_id", middleware.DeserializeUser(), ctrl.ResendInvitation)
	team.PUT("/:id/avatar", middleware.DeserializeUser(), ctrl.UpdateTeamAvatar)
	team.DELETE("/:id/avatar", middleware.DeserializeUser(), ctrl.DeleteTeamAvatar)
}

// @Summary Get team by ID
// @Schemes
// @Description Get team data by ID
// @Tags Team
// @Accept json
// @Produce json
// @Param id path string true "Team ID"
// @Success 200 {object} dto.TeamRetrievalSchema
// @Router /teams/{id} [get]
func (ctrl *teamController) GetTeamById(ctx *gin.Context) {
	bus := ctx.MustGet("bus").(*service.MessageBus)
	uow := bus.UoW
	currentUser := ctx.MustGet("currentUser").(*model.User)

	// Get team ID from request parameter
	id := ctx.Param("id")
	log.Debug().Str("id", id).Msg("Get team data by ID")

	// Get team data from database
	team, err := view.Team(uuid.FromStringOrNil(id), currentUser, uow)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get team data by ID")
		_ = ctx.Error(err)
		return
	}

	// Return success response
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"team": team}})
}

// @Summary Get all teams
// @Schemes
// @Description Get all teams data
// @Tags Team
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Param name query string false "Team name"
// @Success 200 {object} dto.Pagination
// @Router /teams [get]
func (ctrl *teamController) GetTeams(ctx *gin.Context) {
	log.Debug().Msg("Get all team data")
	bus := ctx.MustGet("bus").(*service.MessageBus)
	uow := bus.UoW
	currentUser := ctx.MustGet("currentUser").(*model.User)

	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil {
		log.Error().Err(err).Msg("failed to parse page query")
		_ = ctx.Error(err)
	}

	pageSize, err := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
	if err != nil {
		log.Error().Err(err).Msg("failed to parse page size query")
		_ = ctx.Error(err)
	}

	name := ctx.DefaultQuery("name", "")

	if page < 1 {
		page = 1
	}

	if pageSize < 10 {
		pageSize = 10
	}

	// Get team data from database
	teams, err := view.Teams(uow, currentUser, name, page, pageSize)
	if err != nil {
		log.Error().Err(err).Msg("failed to get all team data")
		_ = ctx.Error(err)
		return
	}

	// Return success response
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": teams})
}

// @Summary Create team
// @Schemes
// @Description Create team data
// @Tags Team
// @Accept json
// @Produce json
// @Param team_id body string true "Team ID"
// @Param name body string true "Team name"
// @Param is_personal body bool false "Is personal team"
// @Param description body string false "Team description"
// @Success 201 {string} string "OK"
// @Router /teams [post]
func (ctrl *teamController) CreateTeam(ctx *gin.Context) {
	log.Debug().Msg("Create team data")
	bus := ctx.MustGet("bus").(*service.MessageBus)
	currentUser := ctx.MustGet("currentUser").(*model.User)

	// Parse the request body into a User struct
	var cmd command.CreateTeam
	if err := ctx.ShouldBindJSON(&cmd); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd.User = currentUser

	err := bus.Handle(&cmd)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create team data")
		_ = ctx.Error(err)
		return
	}

	// Return success response
	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "message": "OK"})
}

// @Summary Update team
// @Schemes
// @Description Update team data
// @Tags Team
// @Accept json
// @Produce json
// @Param team_id path string true "Team ID"
// @Param name body string false "Team name"
// @Param description body string false "Team description"
// @Success 200 {string} string "OK"
// @Router /teams/{id} [put]
func (ctrl *teamController) UpdateTeam(ctx *gin.Context) {
	log.Debug().Msg("Update team data")
	bus := ctx.MustGet("bus").(*service.MessageBus)
	currentUser := ctx.MustGet("currentUser").(*model.User)

	// Get team ID from request parameter
	id := ctx.Param("id")
	log.Debug().Str("id", id).Msg("Get team data by ID")

	// Parse the request body into a User struct
	var cmd command.UpdateTeam
	if err := ctx.ShouldBindJSON(&cmd); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd.TeamID = uuid.FromStringOrNil(id)
	cmd.User = currentUser

	err := bus.Handle(&cmd)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update team data")
		_ = ctx.Error(err)
		return
	}

	// Return success response
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "OK"})
}

// @Summary Update last active team date
// @Schemes
// @Description Update last active team date
// @Tags Team
// @Accept json
// @Produce json
// @Param team_id path string true "Team ID"
// @Success 200 {string} string "OK"
// @Router /teams/{id}/last-active [put]
func (ctrl *teamController) UpdateLastActiveTeam(ctx *gin.Context) {
	bus := ctx.MustGet("bus").(*service.MessageBus)
	currentUser := ctx.MustGet("currentUser").(*model.User)

	// Get team ID from request parameter
	id := ctx.Param("id")
	log.Debug().Str("id", id).Msg("Update last active team data")

	// Parse the request body into a User struct
	var cmd command.UpdateLastActiveTeam

	cmd.TeamID = uuid.FromStringOrNil(id)
	cmd.User = currentUser

	err := bus.Handle(&cmd)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update last active team data")
		_ = ctx.Error(err)
		return
	}

	// Return success response
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "OK"})
}

// @Summary Delete team member
// @Schemes
// @Description Delete team member
// @Tags Membership
// @Accept json
// @Produce json
// @Param team_id path string true "Team ID"
// @Param membership_id path string true "Membership ID"
// @Success 200 {string} string "OK"
// @Router /teams/{id}/members/{membership_id} [delete]
func (ctrl *teamController) DeleteTeamMember(ctx *gin.Context) {
	bus := ctx.MustGet("bus").(*service.MessageBus)
	currentUser := ctx.MustGet("currentUser").(*model.User)

	// Get team ID from request parameter
	id := ctx.Param("id")
	log.Debug().Str("id", id).Msg("Delete team member data")

	// Get membership ID from request parameter
	membershipID := ctx.Param("membership_id")
	log.Debug().Str("membership_id", membershipID).Msg("Delete team member data")

	var cmd command.DeleteTeamMember

	cmd.TeamID = uuid.FromStringOrNil(id)
	cmd.MembershipID = uuid.FromStringOrNil(membershipID)
	cmd.User = currentUser

	err := bus.Handle(&cmd)
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete team member data")
		_ = ctx.Error(err)
		return
	}

	// Return success response
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "OK"})
}

// @Summary Change team member role
// @Schemes
// @Description Change team member role
// @Tags Membership
// @Accept json
// @Produce json
// @Param team_id path string true "Team ID"
// @Param membership_id path string true "Membership ID"
// @Success 200 {string} string "OK"
// @Router /teams/{id}/members/{membership_id} [put]
func (ctrl *teamController) ChangeMemberRole(ctx *gin.Context) {
	bus := ctx.MustGet("bus").(*service.MessageBus)
	currentUser := ctx.MustGet("currentUser").(*model.User)

	// Get team ID from request parameter
	id := ctx.Param("id")
	log.Debug().Str("id", id).Msg("Change team member role")

	// Get membership ID from request parameter
	membershipID := ctx.Param("membership_id")
	log.Debug().Str("membership_id", membershipID).Msg("Change team member role")

	// Parse the request body into a ChangeMemberRole struct
	var cmd command.ChangeMemberRole
	if err := ctx.ShouldBindJSON(&cmd); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd.TeamID = uuid.FromStringOrNil(id)
	cmd.MembershipID = uuid.FromStringOrNil(membershipID)
	cmd.User = currentUser

	err := bus.Handle(&cmd)
	if err != nil {
		log.Error().Err(err).Msg("Failed to change team member role")
		_ = ctx.Error(err)
		return
	}

	// Return success response
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "OK"})
}

// @Summary Send invitation
// @Schemes
// @Description Send invitation to join team
// @Tags Membership
// @Accept json
// @Produce json
// @Param team_id path string true "Team ID"
// @Success 200 {string} string "OK"
// @Router /teams/{id}/invitation [post]
func (ctrl *teamController) SendInvitation(ctx *gin.Context) {
	bus := ctx.MustGet("bus").(*service.MessageBus)
	currentUser := ctx.MustGet("currentUser").(*model.User)

	// Get team ID from request parameter
	id := ctx.Param("id")
	log.Debug().Str("id", id).Msg("Send invitation to join team")

	// Parse the request body into a User struct
	var cmd command.InviteMember
	if err := ctx.ShouldBindJSON(&cmd); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd.TeamID = uuid.FromStringOrNil(id)
	cmd.Sender = currentUser

	log.Debug().Any("Invitees data", cmd.Invitees).Msg("Send invitation to join team")

	err := bus.Handle(&cmd)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send invitation to join team")
		_ = ctx.Error(err)
		return
	}

	// Return success response
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "OK"})
}

// @Summary Resend Invitation
// @Schemes
// @Description Resend invitation to join team
// @Tags Membership
// @Accept json
// @Produce json
// @Param team_id path string true "Team ID"
// @Param invitation_id path string true "Invitation ID"
// @Success 200 {string} string "OK"
// @Router /teams/{id}/invitation/{invitation_id} [post]
func (ctrl *teamController) ResendInvitation(ctx *gin.Context) {
	bus := ctx.MustGet("bus").(*service.MessageBus)
	currentUser := ctx.MustGet("currentUser").(*model.User)

	// Get team ID from request parameter
	id := ctx.Param("id")
	log.Debug().Str("id", id).Msg("Resend invitation to join team")

	// Get invitation ID from request parameter
	invitationIDString := ctx.Param("invitation_id")
	log.Debug().Str("invitation_id", invitationIDString).Msg("Resend invitation with ID")

	// Parse the request body into a User struct
	var cmd command.ResendInvitation
	if err := ctx.ShouldBindJSON(&cmd); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	invitationID, err := ulid.Parse(invitationIDString)

	if err != nil {
		badRequest := exception.NewBadGatewayException(err.Error())
		_ = ctx.Error(badRequest)
		return
	}

	cmd.InvitationID = invitationID
	cmd.TeamID = uuid.FromStringOrNil(id)
	cmd.Sender = currentUser

	err = bus.Handle(&cmd)
	if err != nil {
		log.Error().Err(err).Msg("Failed to resend invitation to join team")
		_ = ctx.Error(err)
		return
	}

	// Return success response
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "OK"})
}

// @Summary Update team avatar
// @Schemes
// @Description Update team avatar
// @Tags Team
// @Accept json
// @Produce json
// @Param id path int true "Team ID"
// @Param avatar formData file true "avatar"
// @Success 200 {string} string "OK"
// @Router /teams/{team_id}/avatar [put]
func (ctrl *teamController) UpdateTeamAvatar(ctx *gin.Context) {
	bus := ctx.MustGet("bus").(*service.MessageBus)

	// Get team ID from request parameter
	id := ctx.Param("id")
	log.Debug().Str("id", id).Msg("Update team avatar")

	// Parse the request body into a User struct
	var cmd command.UpdateTeamAvatar
	if err := ctx.ShouldBind(&cmd); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd.TeamID = uuid.FromStringOrNil(id)

	err := bus.Handle(&cmd)
	if err != nil {
		log.Error().Err(err).Msg("Failed to update team avatar")
		_ = ctx.Error(err)
		return
	}

	// Return user data
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"message": "OK"}})
}

// @Summary Delete team avatar
// @Schemes
// @Description Delete team avatar
// @Tags Team
// @Accept json
// @Produce json
// @Param id path int true "Team ID"
// @Success 200 {string} string "OK"
// @Router /teams/{team_id}/avatar [delete]
func (ctrl *teamController) DeleteTeamAvatar(ctx *gin.Context) {
	bus := ctx.MustGet("bus").(*service.MessageBus)

	// Get team ID from request parameter
	id := ctx.Param("id")
	log.Debug().Str("id", id).Msg("Delete team avatar")

	var cmd command.DeleteTeamAvatar

	cmd.TeamID = uuid.FromStringOrNil(id)

	err := bus.Handle(&cmd)
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete team avatar")
		_ = ctx.Error(err)
		return
	}

	// Return success response
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "OK"})
}
