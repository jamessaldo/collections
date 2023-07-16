package model

import (
	"authorization/domain/dto"
	"fmt"
	"html"
	"math/rand"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

type User struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;"`
	FirstName   string    `gorm:"size:100;not null;"`
	LastName    string    `gorm:"size:100;not null;"`
	Email       string    `gorm:"size:100;not null;unique"`
	Username    string    `gorm:"size:100;not null;unique"`
	Password    string    `gorm:"size:100;not null;"`
	PhoneNumber string    `gorm:"size:20;default:''"`
	AvatarURL   string    `gorm:"default:'';"`

	IsActive bool   `gorm:"default:true;"`
	Verified bool   `gorm:"default:false;"`
	Provider string `gorm:"default:'local';"`

	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
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
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
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

func NewUser(firstName, lastName, email, avatarURL, provider string, isVerified bool) *User {
	now := time.Now()
	username := strings.ToLower(strings.Split(email, "@")[0])

	return &User{
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
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
