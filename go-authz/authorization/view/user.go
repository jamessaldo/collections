package view

import (
	"authorization/config"
	"authorization/controller/exception"
	"authorization/domain/dto"
	"authorization/infrastructure/persistence"
	"authorization/repository"
	"authorization/util"
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	uuid "github.com/satori/go.uuid"
)

func RefreshAccessToken(ctx context.Context, refreshToken string) (string, error) {
	tokenClaims, err := util.ValidateToken(refreshToken, config.AppConfig.RefreshTokenPublicKey)
	if err != nil {
		return "", err
	}

	userId, err := persistence.RedisClient.Get(ctx, tokenClaims.TokenUlid.String()).Result()
	if err == redis.Nil {
		return "", err
	}

	user, err := repository.User.Get(ctx, uuid.FromStringOrNil(userId))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", exception.NewNotFoundException(err.Error())
		}
		return "", err
	}

	accessToken, err := util.CreateToken(user.ID, config.AppConfig.AccessTokenExpiresIn, config.AppConfig.AccessTokenPrivateKey)
	if err != nil {
		return "", err
	}

	now := util.GetTimestampUTC()

	errAccess := persistence.RedisClient.Set(ctx, *accessToken.Token, user.ID.String(), time.Unix(*accessToken.ExpiresIn, 0).Sub(now)).Err()
	if errAccess != nil {
		return "", errAccess
	}

	return *accessToken.Token, nil
}

func User(ctx context.Context, id uuid.UUID) (*dto.PublicUser, error) {
	user, err := repository.User.Get(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, exception.NewNotFoundException(err.Error())
		}
		return nil, err
	}

	return user.PublicUser(), nil
}

func Users(ctx context.Context, page, pageSize int) (dto.Pagination, error) {
	users, err := repository.User.List(ctx, page, pageSize)
	if err != nil {
		return dto.Pagination{}, err
	}

	totalData, err := repository.User.Count(ctx)
	if err != nil {
		return dto.Pagination{}, err
	}

	return dto.Paginate(page, pageSize, totalData, users.PublicUsers()), nil
}
