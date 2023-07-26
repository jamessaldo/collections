package command

import (
	"authorization/domain"
	"mime/multipart"
)

// TODO: fix naming/usage of command
type UpdateUser struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	User        domain.User
	Command
}

type DeleteUser struct {
	User domain.User
	Command
}

type UpdateUserAvatar struct {
	File multipart.FileHeader `form:"avatar" binding:"required"`
	User domain.User
}

type DeleteUserAvatar struct {
	User domain.User
}
