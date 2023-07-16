package handlers

import (
	"strings"

	"authorization/config"
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

	email := strings.ToLower(cmd.GoogleUser.Email)

	_, userErr := uow.User.GetByEmail(email)
	if userErr != nil {
		user := model.NewUser(cmd.GoogleUser.GivenName, cmd.GoogleUser.FamilyName, cmd.GoogleUser.Email, cmd.GoogleUser.Picture, GoogleProvider, cmd.GoogleUser.VerifiedEmail)

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
			return roleErr
		}

		team := model.NewTeam(user, uuid.NewV4(), ownerRole.ID, "", "", true)
		_, err := uow.Team.Add(team, tx)
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
	data := map[string]interface{}{
		"FullName": user.FullName(),
	}
	payload := mailer.CreatePayload(worker.WelcomingTemplate, user.Email, fmt.Sprintf("Selamat Datang di %s!", config.AppConfig.AppName), data)
	return mailer.SendEmail(payload)
}
