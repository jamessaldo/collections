package handlers

import (
	"auth/domain/command"
	"auth/infrastructure/worker"
	"auth/service"
	"fmt"

	"gorm.io/gorm"
)

func UpdateUserWrapper(uow *service.UnitOfWork, mailer worker.WorkerInterface, cmd interface{}) error {
	if c, ok := cmd.(*command.UpdateUser); ok {
		return UpdateUser(uow, c)
	}
	return fmt.Errorf("invalid command type, expected *command.UpdateUser, got %T", cmd)
}

func UpdateUser(uow *service.UnitOfWork, cmd *command.UpdateUser) error {
	tx, txErr := uow.Begin(&gorm.Session{})
	if txErr != nil {
		return txErr
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

	_, err := uow.User.Update(cmd.User)
	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}

func DeleteUserWrapper(uow *service.UnitOfWork, mailer worker.WorkerInterface, cmd interface{}) error {
	if c, ok := cmd.(*command.DeleteUser); ok {
		return DeleteUser(uow, c)
	}
	return fmt.Errorf("invalid command type, expected *command.DeleteUser, got %T", cmd)
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
