package worker

import (
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

type EmailTemplate string

const (
	// TypeEmailTask is a name of the task type
	// for sending an email.
	TypeEmailTask = "email:task"

	// TypeDelayedEmail is a name of the task type
	// for sending a delayed email.
	TypeDelayedEmail = "email:delayed"

	// Email templates
	InvitationTemplate EmailTemplate = "invitation-message.html"
	WelcomingTemplate  EmailTemplate = "welcoming-message.html"
)

type EmailPayload struct {
	UserName     string
	TemplateName EmailTemplate
	To           string
	Subject      string
	Data         map[string]interface{}
}

// newEmailTask task payload for a new email.
func newEmailTask(data *EmailPayload) *asynq.Task {
	// Specify task payload.
	b, err := json.Marshal(data)
	if err != nil {
		log.Error().Caller().Err(err).Any("data", data).Msg("Failed to marshal payload data for email task")
	}

	// Return a new task with given type and payload.
	return asynq.NewTask(TypeDelayedEmail, b)
}
