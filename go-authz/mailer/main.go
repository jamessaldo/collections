package main

import (
	"fmt"
	"runtime"
	"time"

	log "github.com/sirupsen/logrus"

	"mailer/config"
	"mailer/infrastructure/persistence"
	"mailer/service"

	"github.com/hibiken/asynq"
)

const (
	// TypeEmailTask is a name of the task type
	// for sending an email.
	TypeEmailTask = "email:task"

	// TypeDelayedEmail is a name of the task type
	// for sending a delayed email.
	TypeDelayedEmail = "email:delayed"
)

const (
	Development = "development"
	Production  = "production"
	Local       = "local"
)

func init() {
	log.SetReportCaller(true)
	if config.AppConfig.AppEnv == Production {
		log.SetFormatter(&log.JSONFormatter{})
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(log.DebugLevel)
	}
}

func main() {
	runtime.GOMAXPROCS(2)

	db, err := persistence.CreateDBConnection()
	if err != nil {
		log.Fatal(err)
	}
	ctx := service.Bootstrap(db)

	log.Info(fmt.Sprintf("Starting worker server in %s environment", config.AppConfig.AppEnv))

	// Create and configuring Redis connection.
	redisConnection := asynq.RedisClientOpt{
		Addr:         config.StorageConfig.RedisHost + ":" + config.StorageConfig.RedisPort, // Redis server address
		Password:     config.StorageConfig.RedisPassword,
		DialTimeout:  time.Second * 10,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	// Create and configuring Asynq worker server.
	worker := asynq.NewServer(redisConnection, asynq.Config{
		BaseContext: ctx,
		// Specify how many concurrent workers to use.
		Concurrency: 10,
		// Specify multiple queues with different priority.
		Queues: map[string]int{
			"critical": 6, // processed 60% of the time
			"default":  3, // processed 30% of the time
			"low":      1, // processed 10% of the time
		},
	})

	// Create a new task's mux instance.
	mux := asynq.NewServeMux()

	// Define a task handler for the email task.
	mux.HandleFunc(
		TypeEmailTask,           // task type
		service.HandleEmailTask, // handler function
	)

	// Define a task handler for the delayed email task.
	mux.HandleFunc(
		TypeDelayedEmail,               // task type
		service.HandleDelayedEmailTask, // handler function
	)

	// Run worker server.
	if err := worker.Run(mux); err != nil {
		log.Fatal(err)
	}
}
