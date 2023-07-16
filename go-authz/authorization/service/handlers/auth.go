package handlers

import (
	"context"
	"strings"
	"time"

	"authorization/config"
	"authorization/controller/exception"
	"authorization/domain/command"
	"authorization/domain/model"
	"authorization/infrastructure/persistence"
	"authorization/infrastructure/worker"
	"authorization/service"
	"authorization/util"

	"fmt"

	"github.com/rs/zerolog/log"
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

	googleUser, err := util.GetGoogleUser(cmd.Code)
	if err != nil {
		err = exception.NewBadGatewayException(err.Error())
		log.Error().Err(err).Msg("could not get google user")
		return err
	}

	email := strings.ToLower(googleUser.Email)
	user, userErr := uow.User.GetByEmail(email)
	if userErr != nil {
		user = model.NewUser(googleUser.GivenName, googleUser.FamilyName, googleUser.Email, googleUser.Picture, GoogleProvider, googleUser.VerifiedEmail)

		existUser, userErr := uow.User.GetByUsername(user.Username)
		if userErr == nil {
			log.Info().Msg(fmt.Sprintf("username %s already exist, generating random username", user.Username))
			existUser.RegenerateUsername()
		}

		_, userErr = uow.User.Add(user, tx)
		if userErr != nil {
			return userErr
		}

		ownerRole, roleErr := uow.Role.Get(model.Owner)
		if roleErr != nil {
			return roleErr
		}

		team := model.NewTeam(user, ownerRole.ID, "", "", true)
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

	accessToken, refreshToken, err := util.GenerateTokens(user)
	if err != nil {
		log.Error().Err(err).Msg("could not generate token")
		return err
	}

	ctx := context.TODO()
	now := time.Now()

	errAccess := persistence.RedisClient.Set(ctx, *accessToken.Token, user.ID.String(), time.Unix(*accessToken.ExpiresIn, 0).Sub(now)).Err()
	if errAccess != nil {
		log.Error().Err(err).Msg("could not set token to redis")
		return errAccess
	}

	errRefresh := persistence.RedisClient.Set(ctx, *refreshToken.Token, user.ID.String(), time.Unix(*refreshToken.ExpiresIn, 0).Sub(now)).Err()
	if errRefresh != nil {
		log.Error().Err(err).Msg("could not set refresh token to redis")
		return errRefresh
	}

	cmd.Token = *accessToken.Token
	cmd.RefreshToken = *refreshToken.Token

	return nil
}

func sendWelcomeEmail(mailer worker.WorkerInterface, user *model.User) error {
	data := map[string]interface{}{
		"FullName": user.FullName(),
	}
	payload := mailer.CreatePayload(worker.WelcomingTemplate, user.Email, fmt.Sprintf("Selamat Datang di %s!", config.AppConfig.AppName), data)
	return mailer.SendEmail(payload)
}
