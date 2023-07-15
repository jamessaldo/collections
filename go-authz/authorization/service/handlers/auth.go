package handlers

import (
	"errors"
	"strings"
	"time"

	"authorization/config"
	"authorization/controller/exception"
	"authorization/domain/command"
	"authorization/domain/model"
	"authorization/infrastructure/worker"
	"authorization/service"
	"authorization/util"

	"fmt"

	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

const (
	GoogleProvider   = "Google"
	FacebookProvider = "Facebook"
)

func LoginByGoogleWrapper(uow *service.UnitOfWork, mailer worker.WorkerInterface, cmd interface{}) error {
	if c, ok := cmd.(*command.LoginByGoogle); ok {
		return LoginByGoogle(uow, mailer, c)
	}
	return fmt.Errorf("invalid command type, expected *command.LoginByGoogle, got %T", cmd)
}

func LoginByGoogle(uow *service.UnitOfWork, mailer worker.WorkerInterface, cmd *command.LoginByGoogle) error {
	tx, txErr := uow.Begin(&gorm.Session{})
	if txErr != nil {
		return txErr
	}

	defer func() {
		tx.Rollback()
	}()

	now := time.Now()
	email := strings.ToLower(cmd.GoogleUser.Email)

	_, userErr := uow.User.GetByEmail(email)
	if userErr != nil {
		user := &model.User{
			ID:          uuid.NewV4(),
			FirstName:   cmd.GoogleUser.GivenName,
			LastName:    cmd.GoogleUser.FamilyName,
			Email:       cmd.GoogleUser.Email,
			Password:    "",
			PhoneNumber: "",
			AvatarURL:   cmd.GoogleUser.Picture,
			Provider:    GoogleProvider,
			Verified:    cmd.GoogleUser.VerifiedEmail,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		username := strings.Split(email, "@")[0]
		user.Username = strings.ToLower(username)

		_, userErr = uow.User.GetByUsername(user.Username)
		if userErr == nil {
			log.Info().Msg(fmt.Sprintf("username %s already exist, generating random username", user.Username))
			user.Username = util.RandomUsername(user.Username)
		}

		_, userErr = uow.User.Add(user, tx)
		if userErr != nil {
			return userErr
		}

		ownerRole, roleErr := uow.Role.Get(model.Owner)
		if roleErr != nil {
			if errors.Is(roleErr, gorm.ErrRecordNotFound) {
				return exception.NewNotFoundException(fmt.Sprintf("Role with name %s is not exist! Detail: %s", model.Owner, roleErr.Error()))
			}
			return roleErr
		}

		membership := user.AddPersonalTeam(ownerRole)

		_, err := uow.Membership.Add(membership, tx)
		if err != nil {
			return err
		}

		errSendMail := sendWelcomeEmail(mailer, user)
		if errSendMail != nil {
			return errSendMail
		}

		tx.Commit()
	}

	return nil
}

func sendWelcomeEmail(mailer worker.WorkerInterface, user *model.User) error {
	payload := &worker.Payload{
		TemplateName: "welcoming-message.html",
		To:           user.Email,
		Subject:      fmt.Sprintf("Selamat Datang di %s!", config.AppConfig.AppName),
		Data: map[string]interface{}{
			"FullName": user.FullName(),
		},
	}

	return mailer.SendEmail(payload)
}
