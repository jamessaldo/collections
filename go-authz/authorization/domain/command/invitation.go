package command

import (
	"auth/domain/model"

	uuid "github.com/satori/go.uuid"
)

type InviteMember struct {
	TeamID   uuid.UUID `json:"team_id"`
	Invitees []Invitee `json:"invitees"`
	Sender   *model.User
	Command
}

// make a stuct that contains Email and Role, and use it in the InviteMember struct
type Invitee struct {
	Email string         `json:"email"`
	Role  model.RoleType `json:"role"`
}

type UpdateInvitationStatus struct {
	InvitationID string `json:"invitation_id"`
	Status       string `json:"status"`
	User         *model.User
	Command
}

type DeleteInvitation struct {
	InvitationID string `json:"invitation_id"`
	User         *model.User
	Command
}

type ResendInvitation struct {
	InvitationID string    `json:"invitation_id"`
	TeamID       uuid.UUID `json:"team_id"`
	Sender       *model.User
	Command
}
