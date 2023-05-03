package command

import (
	"authorization/domain/model"
	"mime/multipart"
)

// TODO: fix naming/usage of command
type UpdateUser struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	User        *model.User
	Command
}

type DeleteUser struct {
	User *model.User
	Command
}

type UpdateUserAvatar struct {
	File *multipart.FileHeader `form:"avatar" binding:"required"`
	User *model.User
}

type DeleteUserAvatar struct {
	User *model.User
}
