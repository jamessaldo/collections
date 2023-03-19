package handlers

import (
	"errors"
	"strings"
	"time"

	"auth/config"
	"auth/controller/exception"
	"auth/domain/command"
	"auth/domain/model"
	"auth/infrastructure/worker"
	"auth/service"
	"auth/util"

	"fmt"

	uuid "github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	GoogleProvider   = "Google"
	FacebookProvider = "Facebook"
)

func LoginByGoogle(uow *service.UnitOfWork, mailer worker.WorkerInterface, cmd command.LoginByGoogle) (string, string, error) {
	tx, txErr := uow.Begin(&gorm.Session{})
	if txErr != nil {
		return "", "", txErr
	}

	defer func() {
		tx.Rollback()
	}()

	if cmd.Code == "" {
		return "", "", exception.NewBadGatewayException("authorization code not provided")
	}

	tokenRes, err := util.GetGoogleOauthToken(cmd.Code)
	if err != nil {
		return "", "", exception.NewBadGatewayException(err.Error())
	}

	googleUser, err := util.GetGoogleUser(tokenRes.Access_token, tokenRes.Id_token)
	if err != nil {
		return "", "", exception.NewBadGatewayException(err.Error())
	}

	now := time.Now()
	email := strings.ToLower(googleUser.Email)

	user, userErr := uow.User.GetByEmail(email)
	if userErr != nil {
		userId := uuid.NewV4()
		user = &model.User{
			ID:          userId,
			FirstName:   googleUser.GivenName,
			LastName:    googleUser.FamilyName,
			Email:       googleUser.Email,
			Password:    "",
			PhoneNumber: "",
			AvatarURL:   googleUser.Picture,
			Provider:    GoogleProvider,
			Verified:    googleUser.VerifiedEmail,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		username := strings.Split(email, "@")[0]
		user.Username = strings.ToLower(username)

		_, userErr = uow.User.GetByUsername(user.Username)
		if userErr == nil {
			log.Info(fmt.Sprintf("username %s already exist, generating random username", user.Username))
			user.Username = util.RandomUsername(user.Username)
		}

		_, userErr = uow.User.Add(user)
		if userErr != nil {
			return "", "", userErr
		}

		ownerRole, roleErr := uow.Role.Get(model.Owner)
		if roleErr != nil {
			if errors.Is(roleErr, gorm.ErrRecordNotFound) {
				return "", "", exception.NewNotFoundException(fmt.Sprintf("Role with name %s is not exist! Detail: %s", model.Owner, roleErr.Error()))
			}
			return "", "", roleErr
		}

		team := &model.Team{
			ID:          uuid.NewV4(),
			Name:        fmt.Sprintf("%s's Personal Team", user.FullName()),
			Description: fmt.Sprintf("%s's Personal Team will contains your personal apps.", user.FullName()),
			IsPersonal:  true,
			CreatorID:   user.ID,
			Creator:     user,
		}

		membership := &model.Membership{
			ID:     uuid.NewV4(),
			TeamID: team.ID,
			Team:   team,
			UserID: user.ID,
			RoleID: ownerRole.ID,
		}

		_, err := uow.Membership.Add(membership)
		if err != nil {
			return "", "", err
		}

		errSendMail := sendWelcomeEmail(mailer, user)
		if errSendMail != nil {
			return "", "", errSendMail
		}

		tx.Commit()
	}

	token, refreshToken, err := generateTokens(user.ID)
	if err != nil {
		return "", "", err
	}

	return token, refreshToken, nil
}

func generateTokens(userId uuid.UUID) (string, string, error) {
	token, err := util.GenerateToken(config.AppConfig.TokenExpiresIn, userId.String(), config.AppConfig.JWTTokenSecret)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := util.GenerateRefreshToken(config.AppConfig.RefreshTokenExpiresIn, userId.String(), config.AppConfig.RefreshJWTTokenSecret)
	if err != nil {
		return "", "", err
	}

	return token, refreshToken, nil
}

func sendWelcomeEmail(mailer worker.WorkerInterface, user *model.User) error {
	payload := &worker.Payload{
		TemplateName: "welcoming-message.html",
		To:           user.Email,
		Subject:      "Selamat Datang di Wedigo!",
		Data: map[string]interface{}{
			"FullName": user.FullName(),
		},
	}

	return mailer.SendEmail(payload)
}
