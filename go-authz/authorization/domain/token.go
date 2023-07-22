package domain

import (
	"github.com/oklog/ulid/v2"
	uuid "github.com/satori/go.uuid"
)

type TokenDetails struct {
	Token     *string
	TokenUlid ulid.ULID
	UserID    uuid.UUID
	ExpiresIn *int64
}

func NewTokenDetails(token string, tokenUlid ulid.ULID, userID uuid.UUID, expiresIn int64) *TokenDetails {
	return &TokenDetails{
		Token:     &token,
		TokenUlid: tokenUlid,
		UserID:    userID,
		ExpiresIn: &expiresIn,
	}
}
