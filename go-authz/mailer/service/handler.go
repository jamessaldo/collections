package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mailer/config"
	"mailer/domain/model"
	"text/template"

	"mailer/domain"

	"github.com/hibiken/asynq"
	log "github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

func ParseTemplate(templateFileName string, data interface{}) (string, error) {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		log.Error(err)
		return "", err
	}
	return buf.String(), nil
}

func SendEmail(to string, subject string, data interface{}, templateName string) error {
	result, err := ParseTemplate(fmt.Sprintf("%s%s", "./templates/", templateName), data)
	if err != nil {
		log.Error(err)
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s Team <%s>", config.AppConfig.AppName, config.MailerConfig.MailerUsername))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", result)
	d := gomail.NewDialer(config.MailerConfig.MailerHost, config.MailerConfig.MailerPort, config.MailerConfig.MailerUsername, config.MailerConfig.MailerPassword)
	err = d.DialAndSend(m)
	if err != nil {
		log.Error(err)
	}
	return err
}

func SendEmailTask(to, subject, templateName string, data interface{}) {
	err := SendEmail(to, subject, data, templateName)
	if err == nil {
		log.Info(fmt.Sprintf("Email sent to %s", to))
	} else {
		log.Error(err)
	}
}

// HandleEmailTask handler for email task.
func HandleEmailTask(ctx context.Context, t *asynq.Task) error {
	uow, err := NewUnitOfWork(ctx.Value(ContextDBkey).(*gorm.DB))
	if err != nil {
		log.Error(err)
		return err
	}

	// Get user ID from given task.
	var payload domain.Payload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}

	data := payload.Data

	invitation, err := uow.Invitation.List(&model.InvitationOptions{})
	if err != nil {
		log.Error(err)
	}
	log.Info(invitation)

	to := payload.To
	subject := payload.Subject
	templateName := payload.TemplateName

	log.Info(fmt.Sprintf("Sending Email to %s\n", payload.To))
	go SendEmailTask(to, subject, templateName, data)

	return nil
}

// HandleDelayedEmailTask for delayed email task.
func HandleDelayedEmailTask(ctx context.Context, t *asynq.Task) error {
	uow, err := NewUnitOfWork(ctx.Value(ContextDBkey).(*gorm.DB))
	if err != nil {
		log.Error(err)
		return err
	}

	// Get user ID from given task.
	var payload domain.Payload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}

	data := payload.Data

	invitation, err := uow.Invitation.List(&model.InvitationOptions{})
	if err != nil {
		log.Error(err)
	}
	log.Info(invitation)

	to := payload.To
	subject := payload.Subject
	templateName := payload.TemplateName

	log.Info(fmt.Sprintf("Sending Email to %s\n", payload.To))
	go SendEmailTask(to, subject, templateName, data)

	return nil
}
