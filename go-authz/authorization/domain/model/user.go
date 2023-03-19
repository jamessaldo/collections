package model

import (
	"auth/domain/dto"
	"html"
	"strings"
	"time"

	"github.com/badoux/checkmail"
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
func (users Users) PublicUsers() []*dto.PublicUser {
	result := make([]*dto.PublicUser, len(users))
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

func (u *User) Validate(action string) map[string]string {
	var errorMessages = make(map[string]string)
	var err error

	switch strings.ToLower(action) {
	case "update":
		if u.Email == "" {
			errorMessages["email_required"] = "email required"
		}
		if u.Email != "" {
			if err = checkmail.ValidateFormat(u.Email); err != nil {
				errorMessages["invalid_email"] = "email email"
			}
		}
	case "login":
		if u.Password == "" {
			errorMessages["password_required"] = "password is required"
		}
		if u.Email == "" {
			errorMessages["email_required"] = "email is required"
		}
		if u.Email != "" {
			if err = checkmail.ValidateFormat(u.Email); err != nil {
				errorMessages["invalid_email"] = "please provide a valid email"
			}
		}
	case "forgotpassword":
		if u.Email == "" {
			errorMessages["email_required"] = "email required"
		}
		if u.Email != "" {
			if err = checkmail.ValidateFormat(u.Email); err != nil {
				errorMessages["invalid_email"] = "please provide a valid email"
			}
		}
	default:
		if u.FirstName == "" {
			errorMessages["firstname_required"] = "first name is required"
		}
		if u.LastName == "" {
			errorMessages["lastname_required"] = "last name is required"
		}
		if u.Password == "" {
			errorMessages["password_required"] = "password is required"
		}
		if u.Password != "" && len(u.Password) < 6 {
			errorMessages["invalid_password"] = "password should be at least 6 characters"
		}
		if u.Username == "" {
			errorMessages["username_required"] = "username is required"
		}
		if u.Email == "" {
			errorMessages["email_required"] = "email is required"
		}
		if u.Email != "" {
			if err = checkmail.ValidateFormat(u.Email); err != nil {
				errorMessages["invalid_email"] = "please provide a valid email"
			}
		}
	}
	return errorMessages
}

func (u *User) FullName() string {
	fullName := u.FirstName + " " + u.LastName
	words := strings.Fields(fullName)
	fullName = strings.Join(words, " ")
	return fullName
}
