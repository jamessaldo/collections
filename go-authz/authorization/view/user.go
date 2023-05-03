package view

import (
	"authorization/controller/exception"
	"authorization/domain/dto"
	"authorization/service"
	"authorization/util"
	"errors"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func LoginByGoogle(email string, uow *service.UnitOfWork) (string, string, error) {
	user, userErr := uow.User.GetByEmail(email)
	if userErr != nil {
		return "", "", userErr
	}

	token, refreshToken, err := util.GenerateTokens(user)
	if err != nil {
		return "", "", err
	}

	return token, refreshToken, nil
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
