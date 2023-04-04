package command

import "auth/util"

type LoginByGoogle struct {
	Code       string `json:"code"`
	PathURL    string `json:"path_url"`
	GoogleUser util.GoogleUserResult
	Command
}
