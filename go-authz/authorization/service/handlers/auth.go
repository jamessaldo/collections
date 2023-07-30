package handlers

import (
	"authorization/config"
	"authorization/controller/exception"
	"authorization/domain"
	"authorization/domain/command"
	"authorization/infrastructure/persistence"
	"authorization/infrastructure/worker"
	"authorization/repository"
	"authorization/util"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	GoogleProvider   = "Google"
	FacebookProvider = "Facebook"
)

func LoginByGoogle(ctx context.Context, cmd *command.LoginByGoogle) error {
	tx, txErr := persistence.Pool.Begin(ctx)
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
	user, userErr := repository.User.GetByEmail(ctx, email)
	if userErr != nil {
		user = domain.NewUser(googleUser.GivenName, googleUser.FamilyName, googleUser.Email, googleUser.Picture, GoogleProvider, googleUser.VerifiedEmail)

		existUser, userErr := repository.User.GetByUsername(ctx, user.Username)
		if userErr == nil {
			log.Info().Caller().Msg(fmt.Sprintf("username %s already exist, generating random username", user.Username))
			existUser.RegenerateUsername()
		}

		_, userErr = repository.User.Add(ctx, user, tx)
		if userErr != nil {
			return userErr
		}

		ownerRole, roleErr := repository.Role.GetByName(ctx, domain.Owner)
		if roleErr != nil {
			return roleErr
		}

		team := domain.NewTeam(user, ownerRole.ID, "", "", true)
		_, err := repository.Team.Add(ctx, team, tx)
		if err != nil {
			return err
		}

		errSendMail := sendWelcomeEmail(user)
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

func sendWelcomeEmail(user domain.User) error {
	data := map[string]interface{}{
		"FullName": user.FullName(),
	}
	payload := worker.Mailer.CreateEmailPayload(worker.WelcomingTemplate, user.Email, fmt.Sprintf("Selamat Datang di %s!", config.AppConfig.AppName), data)
	return worker.Mailer.SendEmail(payload)
}
