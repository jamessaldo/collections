package model

import (
	"auth/domain/dto"
	"time"

	uuid "github.com/satori/go.uuid"
)

type Team struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;"`
	Name        string    `gorm:"size:100;not null;"`
	Description string    `gorm:"size:100;not null;"`
	IsPersonal  bool      `gorm:"default:false;"`
	AvatarURL   string    `gorm:"default:'';"`
	CreatorID   uuid.UUID `gorm:"type:uuid;not null"`
	Creator     *User     `gorm:"foreignkey:CreatorID;references:ID"`

	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

type Membership struct {
	ID     uuid.UUID `gorm:"type:uuid;primary_key;"`
	TeamID uuid.UUID `gorm:"type:uuid;not null;primaryKey;uniqueIndex:membership_idx"`
	Team   *Team     `gorm:"foreignkey:TeamID;references:ID"`
	UserID uuid.UUID `gorm:"type:uuid;not null;primaryKey;uniqueIndex:membership_idx"`
	User   *User     `gorm:"foreignkey:UserID;references:ID"`
	RoleID uuid.UUID `gorm:"type:uuid;not null;primaryKey;uniqueIndex:membership_idx"`
	Role   *Role     `gorm:"foreignkey:RoleID;references:ID"`

	LastActiveAt *time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	CreatedAt    time.Time  `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time  `gorm:"default:CURRENT_TIMESTAMP"`
}

func (m *Membership) Parse() *dto.MembershipRetrievalSchema {
	return &dto.MembershipRetrievalSchema{
		ID:   m.ID,
		User: m.User.PublicUser(),
		Role: string(m.Role.Name),
	}
}

type MembershipOptions struct {
	Limit        int
	Skip         int
	Name         string
	IsSelectTeam bool
	IsSelectUser bool
	IsSelectRole bool
	TeamID       uuid.UUID
	UserID       uuid.UUID
	RoleID       uuid.UUID
}
