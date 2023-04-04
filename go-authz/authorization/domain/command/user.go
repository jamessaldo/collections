package command

import "auth/domain/model"

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
