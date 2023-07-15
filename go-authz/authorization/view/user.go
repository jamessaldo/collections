package view

import (
	"authorization/config"
	"authorization/controller/exception"
	"authorization/domain/dto"
	"authorization/infrastructure/persistence"
	"authorization/service"
	"authorization/util"
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func LoginByGoogle(email string, uow *service.UnitOfWork) (string, string, error) {
	user, userErr := uow.User.GetByEmail(email)
	if userErr != nil {
		return "", "", userErr
	}

	accessToken, refreshToken, err := util.GenerateTokens(user)
	if err != nil {
		return "", "", err
	}

	ctx := context.TODO()
	now := time.Now()

	errAccess := persistence.RedisClient.Set(ctx, *accessToken.Token, user.ID.String(), time.Unix(*accessToken.ExpiresIn, 0).Sub(now)).Err()
	if errAccess != nil {
		return "", "", errAccess
	}

	errRefresh := persistence.RedisClient.Set(ctx, *refreshToken.Token, user.ID.String(), time.Unix(*refreshToken.ExpiresIn, 0).Sub(now)).Err()
	if errRefresh != nil {
		return "", "", errRefresh
	}

	return *accessToken.Token, *refreshToken.Token, nil
}

func RefreshAccessToken(refreshToken string, uow *service.UnitOfWork) (string, error) {
	ctx := context.TODO()

	tokenClaims, err := util.ValidateToken(refreshToken, config.AppConfig.RefreshTokenPublicKey)
	if err != nil {
		return "", err
	}

	userId, err := persistence.RedisClient.Get(ctx, tokenClaims.TokenUlid.String()).Result()
	if err == redis.Nil {
		return "", err
	}

	user, err := uow.User.Get(uuid.FromStringOrNil(userId))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", exception.NewNotFoundException(err.Error())
		}
		return "", err
	}

	accessToken, err := util.CreateToken(user.ID, config.AppConfig.AccessTokenExpiresIn, config.AppConfig.AccessTokenPrivateKey)
	if err != nil {
		return "", err
	}

	now := time.Now()

	errAccess := persistence.RedisClient.Set(ctx, *accessToken.Token, user.ID.String(), time.Unix(*accessToken.ExpiresIn, 0).Sub(now)).Err()
	if errAccess != nil {
		return "", errAccess
	}

	return *accessToken.Token, nil
}

func User(id uuid.UUID, uow *service.UnitOfWork) (*dto.PublicUser, error) {
	user, err := uow.User.Get(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.NewNotFoundException(err.Error())
		}
		return nil, err
	}

	return user.PublicUser(), nil
}

func Users(uow *service.UnitOfWork, page, pageSize int) (dto.Pagination, error) {
	users, err := uow.User.List(page, pageSize)
	if err != nil {
		return dto.Pagination{}, err
	}

	totalData, err := uow.User.Count()
	if err != nil {
		return dto.Pagination{}, err
	}

	return dto.Paginate(page, pageSize, totalData, users.PublicUsers()), nil
}
