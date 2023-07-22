package domain

import (
	"authorization/controller/exception"
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
	ID        ulid.ULID        `gorm:"type:bytea;primary_key"`
	Email     string           `gorm:"size:100;not null"`
	ExpiresAt time.Time        `gorm:"default:CURRENT_TIMESTAMP"`
	Status    InvitationStatus `gorm:"size:100;not null;"`
	TeamID    uuid.UUID        `gorm:"type:uuid;not null"`
	Team      Team             `gorm:"foreignkey:TeamID;references:ID"`
	RoleID    ulid.ULID        `gorm:"type:bytea;not null"`
	Role      Role             `gorm:"foreignkey:RoleID;references:ID"`
	SenderID  uuid.UUID        `gorm:"type:uuid;not null"`
	Sender    User             `gorm:"foreignkey:SenderID;references:ID"`
	IsActive  bool             `gorm:"default:true"`
	CreatedAt time.Time        `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time        `gorm:"default:CURRENT_TIMESTAMP"`
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
	if invitation.ExpiresAt.After(time.Now()) || invitation.Status != InvitationStatusExpired {
		return exception.NewBadRequestException("invitation is not expired")
	}

	invitation.Status = InvitationStatusPending
	invitation.ExpiresAt = time.Now().Add(time.Hour * 24 * 7)
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

func NewInvitation(email string, status InvitationStatus, teamID, senderID uuid.UUID, roleID ulid.ULID) *Invitation {
	return &Invitation{
		ID:        ulid.Make(),
		Email:     email,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
		Status:    status,
		TeamID:    teamID,
		RoleID:    roleID,
		SenderID:  senderID,
	}
}
