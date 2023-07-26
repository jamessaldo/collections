package dto

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type TeamRetrievalSchema struct {
	ID           uuid.UUID                    `json:"id"`
	Name         string                       `json:"name"`
	Description  string                       `json:"description"`
	AvatarURL    string                       `json:"avatar_url"`
	IsPersonal   bool                         `json:"is_personal"`
	Creator      interface{}                  `json:"creator"`
	LastActiveAt time.Time                    `json:"last_active_at,omitempty"`
	NumOfMembers int64                        `json:"num_of_members,omitempty"`
	Memberships  []*MembershipRetrievalSchema `json:"memberships,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MembershipRetrievalSchema struct {
	ID   uuid.UUID   `json:"id"`
	Role string      `json:"role"`
	User interface{} `json:"user"`
}
