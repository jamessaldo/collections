package domain

type Payload struct {
	UserName     string
	TemplateName string
	To           string
	Subject      string
	Data         map[string]interface{}
}
