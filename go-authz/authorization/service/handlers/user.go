package handlers

import (
	"auth/domain/command"
	"auth/domain/dto"
	"auth/service"

	"gorm.io/gorm"
)

func UpdateUser(uow *service.UnitOfWork, cmd *command.UpdateUser) (*dto.ProfileUser, error) {
	tx, txErr := uow.Begin(&gorm.Session{})
	if txErr != nil {
		return nil, txErr
	}

	defer func() {
		tx.Rollback()
	}()

	if cmd.FirstName != "" {
		cmd.User.FirstName = cmd.FirstName
	}
	if cmd.LastName != "" {
		cmd.User.LastName = cmd.LastName
	}
	if cmd.PhoneNumber != "" {
		cmd.User.PhoneNumber = cmd.PhoneNumber
	}

	user, err := uow.User.Update(cmd.User)
	if err != nil {
		return nil, err
	}

	tx.Commit()

	return user.ProfileUser(), nil
}

func DeleteUser(uow *service.UnitOfWork, cmd *command.DeleteUser) error {
	tx, txErr := uow.Begin(&gorm.Session{})
	if txErr != nil {
		return txErr
	}

	defer func() {
		tx.Rollback()
	}()

	cmd.User.IsActive = false

	_, err := uow.User.Update(cmd.User)
	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}
