package command

import "authorization/util"

type LoginByGoogle struct {
	Code       string `json:"code"`
	PathURL    string `json:"path_url"`
	GoogleUser util.GoogleUserResult
	Command
}
