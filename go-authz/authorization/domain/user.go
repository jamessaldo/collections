package domain

import (
	"authorization/config"
	"authorization/domain/dto"
	"authorization/util"
	"fmt"
	"html"
	"math/rand"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

type User struct {
	ID          uuid.UUID
	FirstName   string
	LastName    string
	Email       string
	Username    string
	Password    string
	PhoneNumber string
	AvatarURL   string

	IsActive bool
	Verified bool
	Provider string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Users []User

// So that we dont expose the user's email address and password to the world
func (users Users) PublicUsers() []interface{} {
	result := make([]interface{}, len(users))
	for index, user := range users {
		result[index] = user.PublicUser()
	}
	return result
}

// So that we dont expose the user's  password to the world
func (u *User) PublicUser() *dto.PublicUser {
	return &dto.PublicUser{
		ID:        u.ID,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Name:      u.FullName(),
		Username:  u.Username,
		Email:     u.Email,
		AvatarURL: u.AvatarURL,
	}
}

func (u *User) ProfileUser() *dto.ProfileUser {
	return &dto.ProfileUser{
		ID:          u.ID,
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		Name:        u.FullName(),
		Username:    u.Username,
		Email:       u.Email,
		AvatarURL:   u.AvatarURL,
		PhoneNumber: u.PhoneNumber,
	}
}

func (u *User) Prepare() {
	u.FirstName = html.EscapeString(strings.TrimSpace(u.FirstName))
	u.LastName = html.EscapeString(strings.TrimSpace(u.LastName))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))
}

func (u *User) FullName() string {
	fullName := u.FirstName + " " + u.LastName
	words := strings.Fields(fullName)
	fullName = strings.Join(words, " ")
	return fullName
}

func (u *User) Update(payload map[string]any) {
	if val, ok := payload["firstName"]; ok && val != "" {
		u.FirstName = payload["firstName"].(string)
	}
	if val, ok := payload["lastName"]; ok && val != "" {
		u.LastName = payload["lastName"].(string)
	}
	if val, ok := payload["phoneNumber"]; ok && val != "" {
		u.PhoneNumber = payload["phoneNumber"].(string)
	}
	if val, ok := payload["avatarURL"]; ok && val != "" {
		u.AvatarURL = payload["avatarURL"].(string)
	}
}

func (u *User) RegenerateUsername() {
	number := rand.Intn(99999-10000) + 10000
	u.Username = fmt.Sprintf("%s%d", u.Username, number)
}

func (user *User) GenerateTokens() (*util.TokenDetails, *util.TokenDetails, error) {
	token, err := util.CreateToken(user.ID, config.AppConfig.AccessTokenExpiresIn, config.AppConfig.AccessTokenPrivateKey)
	if err != nil {
		return nil, nil, err
	}

	refreshToken, err := util.CreateToken(user.ID, config.AppConfig.RefreshTokenExpiresIn, config.AppConfig.RefreshTokenPrivateKey)
	if err != nil {
		return nil, nil, err
	}

	return token, refreshToken, nil
}

func NewUser(firstName, lastName, email, avatarURL, provider string, isVerified bool) User {
	now := util.GetTimestampUTC()
	username := strings.ToLower(strings.Split(email, "@")[0])

	return User{
		ID:          uuid.NewV4(),
		FirstName:   firstName,
		LastName:    lastName,
		Username:    username,
		Email:       email,
		Password:    "",
		PhoneNumber: "",
		AvatarURL:   avatarURL,
		Provider:    provider,
		Verified:    isVerified,
		IsActive:    true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
