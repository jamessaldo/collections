package command

import (
	"auth/domain/model"

	uuid "github.com/satori/go.uuid"
)

type CreateTeam struct {
	TeamID      uuid.UUID `json:"team_id"`
	Name        string    `json:"name"`
	IsPersonal  bool      `json:"is_personal"`
	Description string    `json:"description"`
	User        *model.User
	Command
}

type UpdateTeam struct {
	TeamID      uuid.UUID `json:"team_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	User        *model.User
	Command
}

type UpdateLastActiveTeam struct {
	TeamID uuid.UUID `json:"team_id"`
	User   *model.User
	Command
}

type SendInvitation struct {
	Members []Invitee `json:"members"`
	TeamID  uuid.UUID `json:"team_id"`
	User    *model.User
	Command
}

type DeleteTeamMember struct {
	TeamID       uuid.UUID `json:"team_id"`
	MembershipID uuid.UUID `json:"membership_id"`
	User         *model.User
	Command
}

type ChangeMemberRole struct {
	TeamID       uuid.UUID      `json:"team_id"`
	MembershipID uuid.UUID      `json:"membership_id"`
	Role         model.RoleType `json:"role"`
	User         *model.User
	Command
}
