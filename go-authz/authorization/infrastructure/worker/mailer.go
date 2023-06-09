package worker

import (
	"authorization/config"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/hibiken/asynq"
)

type WorkerInterface interface {
	SendEmail(payload *Payload) error
}

type AsynqClient struct {
	client *asynq.Client
}

var _ WorkerInterface = &AsynqClient{}

func NewMailer(client *asynq.Client) *AsynqClient {
	return &AsynqClient{client: client}
}

// Enqueue task to send email
func (ac *AsynqClient) SendEmail(payload *Payload) error {
	// Define tasks.
	task := NewEmailTask(payload)

	// Set a delay duration to 2 minutes.
	delay := 2 * time.Second

	// Process the task immediately in critical queue.
	if _, err := ac.client.Enqueue(
		task,                    // task payload
		asynq.Queue("critical"), // set queue for task
		asynq.ProcessIn(delay),
	); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func CreateAsynqClient() *asynq.Client {
	// Create a new Redis connection for the client.
	redisConnection := asynq.RedisClientOpt{
		Addr: config.StorageConfig.RedisHost + ":" + config.StorageConfig.RedisPort, // Redis server address
	}
	// Create a new Asynq client.
	client := asynq.NewClient(redisConnection)
	return client
}
