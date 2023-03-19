package worker

import (
	"encoding/json"

	"github.com/hibiken/asynq"
	log "github.com/sirupsen/logrus"
)

const (
	// TypeEmailTask is a name of the task type
	// for sending an email.
	TypeEmailTask = "email:task"

	// TypeDelayedEmail is a name of the task type
	// for sending a delayed email.
	TypeDelayedEmail = "email:delayed"
)

type Payload struct {
	UserName     string
	TemplateName string
	To           string
	Subject      string
	Data         map[string]interface{}
}

// NewEmailTask task payload for a new email.
func NewEmailTask(data *Payload) *asynq.Task {
	// Specify task payload.

	b, err := json.Marshal(data)
	if err != nil {
		log.Error(err)
	}

	// Return a new task with given type and payload.
	return asynq.NewTask(TypeEmailTask, b)
}
