package view

import (
	"auth/controller/exception"
	"auth/domain/dto"
	"auth/service"
	"errors"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

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
