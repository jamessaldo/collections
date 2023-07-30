package handlers

import (
	"authorization/domain"
	"context"
	"strings"
	"time"

	"authorization/config"
	"authorization/controller/exception"
	"authorization/domain/command"
	"authorization/infrastructure/mailer"
	"authorization/infrastructure/persistence"
	"authorization/service"
	"authorization/util"

	"fmt"

	"github.com/rs/zerolog/log"
)

const (
	GoogleProvider   = "Google"
	FacebookProvider = "Facebook"
)

func LoginByGoogleWrapper(ctx context.Context, uow *service.UnitOfWork, mailer mailer.MailerInterface, cmd interface{}) error {
	if c, ok := cmd.(*command.LoginByGoogle); ok {
		return LoginByGoogle(ctx, uow, mailer, c)
	}
	return fmt.Errorf("invalid command type, expected *command.LoginByGoogle, got %T", cmd)
}

func LoginByGoogle(ctx context.Context, uow *service.UnitOfWork, mailer mailer.MailerInterface, cmd *command.LoginByGoogle) error {
	tx, txErr := uow.Begin(ctx)
	if txErr != nil {
		return txErr
	}

	defer func() {
		tx.Rollback(ctx)
	}()

	googleUser, err := util.GetGoogleUser(cmd.Code)
	if err != nil {
		err = exception.NewBadGatewayException(err.Error())
		log.Error().Caller().Err(err).Msg("could not get google user")
		return err
	}

	email := strings.ToLower(googleUser.Email)
	user, userErr := uow.User.GetByEmail(ctx, email)
	if userErr != nil {
		user = domain.NewUser(googleUser.GivenName, googleUser.FamilyName, googleUser.Email, googleUser.Picture, GoogleProvider, googleUser.VerifiedEmail)

		existUser, userErr := uow.User.GetByUsername(ctx, user.Username)
		if userErr == nil {
			log.Info().Caller().Msg(fmt.Sprintf("username %s already exist, generating random username", user.Username))
			existUser.RegenerateUsername()
		}

		_, userErr = uow.User.Add(ctx, user, tx)
		if userErr != nil {
			return userErr
		}

		ownerRole, roleErr := uow.Role.GetByName(ctx, domain.Owner)
		if roleErr != nil {
			return roleErr
		}

		team := domain.NewTeam(user, ownerRole.ID, "", "", true)
		_, err := uow.Team.Add(ctx, team, tx)
		if err != nil {
			return err
		}

		errSendMail := sendWelcomeEmail(mailer, user)
		if errSendMail != nil {
			return errSendMail
		}

		tx.Commit(ctx)
	}

	accessToken, refreshToken, err := user.GenerateTokens()
	if err != nil {
		log.Error().Caller().Err(err).Msg("could not generate token")
		return err
	}

	now := util.GetTimestampUTC()

	errAccess := persistence.RedisClient.Set(ctx, *accessToken.Token, user.ID.String(), time.Unix(*accessToken.ExpiresIn, 0).Sub(now)).Err()
	if errAccess != nil {
		log.Error().Caller().Err(err).Msg("could not set token to redis")
		return errAccess
	}

	errRefresh := persistence.RedisClient.Set(ctx, *refreshToken.Token, user.ID.String(), time.Unix(*refreshToken.ExpiresIn, 0).Sub(now)).Err()
	if errRefresh != nil {
		log.Error().Caller().Err(err).Msg("could not set refresh token to redis")
		return errRefresh
	}

	cmd.Token = *accessToken.Token
	cmd.RefreshToken = *refreshToken.Token

	return nil
}

func sendWelcomeEmail(mailerInterface mailer.MailerInterface, user domain.User) error {
	data := map[string]interface{}{
		"FullName": user.FullName(),
	}
	payload := mailerInterface.CreatePayload(mailer.WelcomingTemplate, user.Email, fmt.Sprintf("Selamat Datang di %s!", config.AppConfig.AppName), data)
	return mailerInterface.SendEmail(payload)
}
