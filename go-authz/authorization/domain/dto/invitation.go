package dto

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type InvitationRetreivalSchema struct {
	ID         ulid.ULID           `json:"id"`
	Email      string              `json:"email"`
	ExpiresAt  time.Time           `json:"expires_at"`
	Status     string              `json:"status"`
	Role       string              `json:"role"`
	SenderName string              `json:"sender_name"`
	Team       TeamRetrievalSchema `json:"team"`
}
