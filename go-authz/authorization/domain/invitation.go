package domain

import (
	"authorization/controller/exception"
	"authorization/util"
	"time"

	"github.com/badoux/checkmail"
	"github.com/oklog/ulid/v2"
	uuid "github.com/satori/go.uuid"
)

type InvitationStatus string

var (
	InvitationStatusPending  InvitationStatus = "pending"
	InvitationStatusSent     InvitationStatus = "sent"
	InvitationStatusAccepted InvitationStatus = "accepted"
	InvitationStatusDeclined InvitationStatus = "declined"
	InvitationStatusExpired  InvitationStatus = "expired"
)

type Invitation struct {
	ID        ulid.ULID
	Email     string
	ExpiresAt time.Time
	Status    InvitationStatus
	TeamID    uuid.UUID
	Team      Team
	RoleID    ulid.ULID
	Role      Role
	SenderID  uuid.UUID
	Sender    User
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (invitation *Invitation) Validate(action string) map[string]string {
	var errorMessages = make(map[string]string)
	var err error
	if invitation.Email != "" {
		if err = checkmail.ValidateFormat(invitation.Email); err != nil {
			errorMessages["invalid_email"] = "email address is not valid"
		}
	}
	return errorMessages
}

func (invitation *Invitation) ResendUpdate() error {
	// check if invitation is not expired
	if invitation.ExpiresAt.After(util.GetTimestampUTC()) || invitation.Status != InvitationStatusExpired {
		return exception.NewBadRequestException("invitation is not expired")
	}

	invitation.Status = InvitationStatusPending
	invitation.ExpiresAt = util.GetTimestampUTC().Add(time.Hour * 24 * 7)
	return nil
}

type InvitationOptions struct {
	Email     string
	TeamID    uuid.UUID
	ExpiresAt time.Time
	RoleID    ulid.ULID
	Statuses  []InvitationStatus
	Limit     int
}

func NewInvitation(email string, status InvitationStatus, teamID, senderID uuid.UUID, roleID ulid.ULID) Invitation {
	return Invitation{
		ID:        ulid.Make(),
		Email:     email,
		ExpiresAt: util.GetTimestampUTC().Add(time.Hour * 24 * 7),
		Status:    status,
		TeamID:    teamID,
		RoleID:    roleID,
		SenderID:  senderID,
	}
}
