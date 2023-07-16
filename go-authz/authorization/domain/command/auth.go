package command

type LoginByGoogle struct {
	Code         string `json:"code"`
	PathURL      string `json:"path_url"`
	Email        string
	Token        string
	RefreshToken string
	Command
}
