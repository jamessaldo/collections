package handlers

import (
	"authorization/config"
	"authorization/controller/exception"
	"authorization/domain/command"
	"authorization/infrastructure/worker"
	"authorization/service"
	"authorization/util"
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/oklog/ulid/v2"
)

func UpdateUserWrapper(ctx context.Context, uow *service.UnitOfWork, mailer worker.WorkerInterface, cmd interface{}) error {
	if c, ok := cmd.(*command.UpdateUser); ok {
		return UpdateUser(ctx, uow, c)
	}
	return fmt.Errorf("invalid command type, expected *command.UpdateUser, got %T", cmd)
}

func UpdateUser(ctx context.Context, uow *service.UnitOfWork, cmd *command.UpdateUser) error {
	tx, txErr := uow.Begin(ctx)
	if txErr != nil {
		return txErr
	}

	defer func() {
		tx.Rollback(ctx)
	}()

	payload := map[string]interface{}{
		"firstName":   cmd.FirstName,
		"lastName":    cmd.LastName,
		"phoneNumber": cmd.PhoneNumber,
	}

	cmd.User.Update(payload)

	_, err := uow.User.Update(ctx, cmd.User, tx)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func DeleteUserWrapper(ctx context.Context, uow *service.UnitOfWork, mailer worker.WorkerInterface, cmd interface{}) error {
	if c, ok := cmd.(*command.DeleteUser); ok {
		return DeleteUser(ctx, uow, c)
	}
	return fmt.Errorf("invalid command type, expected *command.DeleteUser, got %T", cmd)
}

func DeleteUser(ctx context.Context, uow *service.UnitOfWork, cmd *command.DeleteUser) error {
	tx, txErr := uow.Begin(ctx)
	if txErr != nil {
		return txErr
	}

	defer func() {
		tx.Rollback(ctx)
	}()

	cmd.User.IsActive = false

	_, err := uow.User.Update(ctx, cmd.User, tx)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func UpdateUserAvatarWrapper(ctx context.Context, uow *service.UnitOfWork, mailer worker.WorkerInterface, cmd interface{}) error {
	if c, ok := cmd.(*command.UpdateUserAvatar); ok {
		return UpdateUserAvatar(ctx, uow, c)
	}
	return fmt.Errorf("invalid command type, expected *command.UpdateUserAvatar, got %T", cmd)
}

func UpdateUserAvatar(ctx context.Context, uow *service.UnitOfWork, cmd *command.UpdateUserAvatar) error {
	tx, txErr := uow.Begin(ctx)
	if txErr != nil {
		return txErr
	}

	defer func() {
		tx.Rollback(ctx)
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
	avatarName := fmt.Sprintf("%s.%s", ulid.Make(), fileType)

	if err := util.SaveFileToLocal(avatarName, cmd.File); err != nil {
		return err
	}

	payload := map[string]interface{}{
		"avatarURL": config.StorageConfig.StaticPublicURL + config.StorageConfig.StaticAvatarPath + avatarName,
	}

	cmd.User.Update(payload)
	_, err := uow.User.Update(ctx, cmd.User, tx)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func DeleteUserAvatarWrapper(ctx context.Context, uow *service.UnitOfWork, mailer worker.WorkerInterface, cmd interface{}) error {
	if c, ok := cmd.(*command.DeleteUserAvatar); ok {
		return DeleteUserAvatar(ctx, uow, c)
	}
	return fmt.Errorf("invalid command type, expected *command.DeleteUserAvatar, got %T", cmd)
}

func DeleteUserAvatar(ctx context.Context, uow *service.UnitOfWork, cmd *command.DeleteUserAvatar) error {
	tx, txErr := uow.Begin(ctx)
	if txErr != nil {
		return txErr
	}

	defer func() {
		tx.Rollback(ctx)
	}()

	if cmd.User.AvatarURL != "" {
		// get avatar path after public url
		paths := strings.Split(cmd.User.AvatarURL, "/")
		path := filepath.Join(config.StorageConfig.StaticRoot, config.StorageConfig.StaticAvatarPath, paths[len(paths)-1])
		if err := util.DeleteFileInLocal(path); err != nil {
			return err
		}
		cmd.User.AvatarURL = ""
		_, err := uow.User.Update(ctx, cmd.User, tx)
		if err != nil {
			return err
		}

		err = tx.Commit(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
