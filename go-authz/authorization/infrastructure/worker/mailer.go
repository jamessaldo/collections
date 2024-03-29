package worker

import (
	"authorization/config"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/hibiken/asynq"
)

var (
	Mailer MailerInterface
)

type MailerInterface interface {
	SendEmail(payload *EmailPayload) error
	CreateEmailPayload(templateName EmailTemplate, to, subject string, data map[string]interface{}) *EmailPayload
}

type AsynqClient struct {
	client *asynq.Client
}

var _ MailerInterface = &AsynqClient{}

func CreateMailer(client *asynq.Client) {
	Mailer = &AsynqClient{client: client}
}

// Enqueue task to send email
func (ac *AsynqClient) SendEmail(payload *EmailPayload) error {
	// Define tasks.
	task := newEmailTask(payload)

	// Set a delay duration to 2 minutes.
	delay := 2 * time.Second

	// Process the task immediately in critical queue.
	if _, err := ac.client.Enqueue(
		task,                    // task payload
		asynq.Queue("critical"), // set queue for task
		asynq.ProcessIn(delay),
	); err != nil {
		log.Error().Caller().Err(err).Msg("Failed to enqueue a task")
		return err
	}
	return nil
}

func (ac *AsynqClient) CreateEmailPayload(templateName EmailTemplate, to, subject string, data map[string]interface{}) *EmailPayload {
	return &EmailPayload{
		TemplateName: templateName,
		To:           to,
		Subject:      subject,
		Data:         data,
	}
}

func CreateMailerClient() *asynq.Client {
	// Create a new Redis connection for the client.
	redisConnection := asynq.RedisClientOpt{
		Addr: config.StorageConfig.RedisHost + ":" + config.StorageConfig.RedisPort, // Redis server address
	}
	// Create a new Asynq client.
	client := asynq.NewClient(redisConnection)
	return client
}
