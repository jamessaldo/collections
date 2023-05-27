package model

type TokenDetails struct {
	Token     *string
	TokenUuid string
	UserID    string
	ExpiresIn *int64
}
