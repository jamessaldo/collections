package command

import (
	"auth/domain/model"
	"mime/multipart"

	uuid "github.com/satori/go.uuid"
)

type CreateTeam struct {
	TeamID      uuid.UUID
	Name        string `json:"name"`
	Description string `json:"description"`
	User        *model.User
	Command
}

type UpdateTeam struct {
	TeamID      uuid.UUID
	Name        string `json:"name"`
	Description string `json:"description"`
	User        *model.User
	Command
}

type UpdateLastActiveTeam struct {
	TeamID uuid.UUID
	User   *model.User
	Command
}

type SendInvitation struct {
	Members []Invitee `json:"members"`
	TeamID  uuid.UUID
	User    *model.User
	Command
}

type DeleteTeamMember struct {
	TeamID       uuid.UUID
	MembershipID uuid.UUID `json:"membership_id"`
	User         *model.User
	Command
}

type ChangeMemberRole struct {
	TeamID       uuid.UUID
	MembershipID uuid.UUID      `json:"membership_id"`
	Role         model.RoleType `json:"role"`
	User         *model.User
	Command
}

type UpdateTeamAvatar struct {
	TeamID uuid.UUID
	File   *multipart.FileHeader `form:"avatar" binding:"required"`
}

type DeleteTeamAvatar struct {
	TeamID uuid.UUID
}
