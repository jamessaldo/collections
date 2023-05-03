package handlers

import (
	"authorization/config"
	"authorization/controller/exception"
	"authorization/domain/command"
	"authorization/infrastructure/worker"
	"authorization/service"
	"authorization/util"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/segmentio/ksuid"
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

	payload := map[string]interface{}{
		"firstName":   cmd.FirstName,
		"lastName":    cmd.LastName,
		"phoneNumber": cmd.PhoneNumber,
	}

	cmd.User.Update(payload)

	_, err := uow.User.Update(cmd.User, tx)
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

	_, err := uow.User.Update(cmd.User, tx)
	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}

func UpdateUserAvatarWrapper(uow *service.UnitOfWork, mailer worker.WorkerInterface, cmd interface{}) error {
	if c, ok := cmd.(*command.UpdateUserAvatar); ok {
		return UpdateUserAvatar(uow, c)
	}
	return fmt.Errorf("invalid command type, expected *command.UpdateUserAvatar, got %T", cmd)
}

func UpdateUserAvatar(uow *service.UnitOfWork, cmd *command.UpdateUserAvatar) error {
	tx, txErr := uow.Begin(&gorm.Session{})
	if txErr != nil {
		return txErr
	}

	defer func() {
		tx.Rollback()
	}()

	fileContentType := cmd.File.Header.Get("Content-Type")

	// check file type, only allow image
	if !strings.HasPrefix(fileContentType, "image") {
		return exception.NewBadRequestException("invalid file type, only allow image")
	}

	if cmd.File.Size > config.StorageConfig.StaticMaxAvatarSize {
		return exception.NewBadRequestException("file size too large")
	}

	if cmd.User.AvatarURL != "" {
		// get avatar path after public url
		paths := strings.Split(cmd.User.AvatarURL, "/")
		path := filepath.Join(config.StorageConfig.StaticRoot, config.StorageConfig.StaticAvatarPath, paths[len(paths)-1])
		if err := util.DeleteFileInLocal(path); err != nil {
			return err
		}
	}

	fileType := strings.Split(fileContentType, "/")[1]
	avatarName := fmt.Sprintf("%s.%s", ksuid.New().String(), fileType)

	if err := util.SaveFileToLocal(avatarName, cmd.File); err != nil {
		return err
	}

	payload := map[string]interface{}{
		"avatarURL": config.StorageConfig.StaticPublicURL + config.StorageConfig.StaticAvatarPath + avatarName,
	}

	cmd.User.Update(payload)
	_, err := uow.User.Update(cmd.User, tx)
	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}

func DeleteUserAvatarWrapper(uow *service.UnitOfWork, mailer worker.WorkerInterface, cmd interface{}) error {
	if c, ok := cmd.(*command.DeleteUserAvatar); ok {
		return DeleteUserAvatar(uow, c)
	}
	return fmt.Errorf("invalid command type, expected *command.DeleteUserAvatar, got %T", cmd)
}

func DeleteUserAvatar(uow *service.UnitOfWork, cmd *command.DeleteUserAvatar) error {
	tx, txErr := uow.Begin(&gorm.Session{})
	if txErr != nil {
		return txErr
	}

	defer func() {
		tx.Rollback()
	}()

	if cmd.User.AvatarURL != "" {
		// get avatar path after public url
		paths := strings.Split(cmd.User.AvatarURL, "/")
		path := filepath.Join(config.StorageConfig.StaticRoot, config.StorageConfig.StaticAvatarPath, paths[len(paths)-1])
		if err := util.DeleteFileInLocal(path); err != nil {
			return err
		}
		cmd.User.AvatarURL = ""
		_, err := uow.User.Update(cmd.User, tx)
		if err != nil {
			return err
		}
		tx.Commit()
	}

	return nil
}
