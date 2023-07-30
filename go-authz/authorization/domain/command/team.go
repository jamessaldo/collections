package command

import (
	"authorization/domain"
	"mime/multipart"

	uuid "github.com/satori/go.uuid"
)

type CreateTeam struct {
	TeamID      uuid.UUID
	Name        string `json:"name"`
	Description string `json:"description"`
	User        domain.User
	Command
}

type UpdateTeam struct {
	TeamID      uuid.UUID
	Name        string `json:"name"`
	Description string `json:"description"`
	User        domain.User
	Command
}

type UpdateLastActiveTeam struct {
	TeamID uuid.UUID
	User   domain.User
	Command
}

type DeleteTeamMember struct {
	TeamID       uuid.UUID
	MembershipID uuid.UUID `json:"membership_id"`
	User         domain.User
	Command
}

type ChangeMemberRole struct {
	TeamID       uuid.UUID
	MembershipID uuid.UUID       `json:"membership_id"`
	Role         domain.RoleType `json:"role"`
	User         domain.User
	Command
}

type UpdateTeamAvatar struct {
	TeamID uuid.UUID
	File   multipart.FileHeader `form:"avatar" binding:"required"`
}

type DeleteTeamAvatar struct {
	TeamID uuid.UUID
}
