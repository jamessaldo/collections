package handlers

import (
	"authorization/config"
	"authorization/controller/exception"
	"authorization/domain"
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

func CreateTeamWrapper(ctx context.Context, uow *service.UnitOfWork, mailer worker.WorkerInterface, cmd interface{}) error {
	if c, ok := cmd.(*command.CreateTeam); ok {
		return CreateTeam(ctx, uow, c)
	}
	return fmt.Errorf("invalid command type, expected *command.CreateTeam, got %T", cmd)
}

func CreateTeam(ctx context.Context, uow *service.UnitOfWork, cmd *command.CreateTeam) error {
	tx, txErr := uow.Begin(ctx)
	if txErr != nil {
		return txErr
	}

	defer func() {
		tx.Rollback(ctx)
	}()

	ownerRole, err := uow.Role.GetByName(ctx, domain.Owner)
	if err != nil {
		return err
	}

	team := domain.NewTeam(cmd.User, ownerRole.ID, cmd.Name, cmd.Description, false)
	_, err = uow.Team.Add(team, tx)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	cmd.TeamID = team.ID
	return nil
}

func UpdateTeamWrapper(ctx context.Context, uow *service.UnitOfWork, mailer worker.WorkerInterface, cmd interface{}) error {
	if c, ok := cmd.(*command.UpdateTeam); ok {
		return UpdateTeam(ctx, uow, c)
	}
	return fmt.Errorf("invalid command type, expected *command.UpdateTeam, got %T", cmd)
}

func UpdateTeam(ctx context.Context, uow *service.UnitOfWork, cmd *command.UpdateTeam) error {
	tx, txErr := uow.Begin(ctx)
	if txErr != nil {
		return txErr
	}

	defer func() {
		tx.Rollback(ctx)
	}()

	team, err := uow.Team.Get(cmd.TeamID)
	if err != nil {
		return err
	}

	team.Update(map[string]interface{}{
		"name":        cmd.Name,
		"description": cmd.Description,
	})

	_, err = uow.Team.Update(team, tx)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func UpdateLastActiveTeamWrapper(ctx context.Context, uow *service.UnitOfWork, mailer worker.WorkerInterface, cmd interface{}) error {
	if c, ok := cmd.(*command.UpdateLastActiveTeam); ok {
		return UpdateLastActiveTeam(ctx, uow, c)
	}
	return fmt.Errorf("invalid command type, expected *command.UpdateLastActiveTeam, got %T", cmd)
}

func UpdateLastActiveTeam(ctx context.Context, uow *service.UnitOfWork, cmd *command.UpdateLastActiveTeam) error {
	tx, txErr := uow.Begin(ctx)
	if txErr != nil {
		return txErr
	}

	defer func() {
		tx.Rollback(ctx)
	}()

	opts := domain.MembershipOptions{
		TeamID:       cmd.TeamID,
		UserID:       cmd.User.ID,
		Limit:        1,
		IsSelectTeam: true,
	}
	memberships, err := uow.Membership.List(opts)
	if err != nil {
		return err
	}

	if len(memberships) == 0 {
		return exception.NewNotFoundException("Team is not found")
	}

	lastActiveAt := util.GetTimestampUTC()
	membership := memberships[0]
	membership.LastActiveAt = lastActiveAt

	_, err = uow.Membership.Update(membership, tx)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func DeleteTeamMemberWrapper(ctx context.Context, uow *service.UnitOfWork, mailer worker.WorkerInterface, cmd interface{}) error {
	if c, ok := cmd.(*command.DeleteTeamMember); ok {
		return DeleteTeamMember(ctx, uow, c)
	}
	return fmt.Errorf("invalid command type, expected *command.DeleteTeamMember, got %T", cmd)
}

func DeleteTeamMember(ctx context.Context, uow *service.UnitOfWork, cmd *command.DeleteTeamMember) error {
	tx, txErr := uow.Begin(ctx)
	if txErr != nil {
		return txErr
	}

	defer func() {
		tx.Rollback(ctx)
	}()

	membership, err := uow.Membership.Get(cmd.MembershipID)
	if err != nil {
		return err
	}

	err = membership.Validation(cmd.User.ID, cmd.TeamID, "")
	if err != nil {
		return err
	}

	err = uow.Membership.Delete(cmd.MembershipID, tx)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func ChangeMemberRoleWrapper(ctx context.Context, uow *service.UnitOfWork, mailer worker.WorkerInterface, cmd interface{}) error {
	if c, ok := cmd.(*command.ChangeMemberRole); ok {
		return ChangeMemberRole(ctx, uow, c)
	}
	return fmt.Errorf("invalid command type, expected *command.ChangeMemberRole, got %T", cmd)
}

func ChangeMemberRole(ctx context.Context, uow *service.UnitOfWork, cmd *command.ChangeMemberRole) error {
	tx, txErr := uow.Begin(ctx)
	if txErr != nil {
		return txErr
	}

	defer func() {
		tx.Rollback(ctx)
	}()

	membership, err := uow.Membership.Get(cmd.MembershipID)
	if err != nil {
		return err
	}

	err = membership.Validation(cmd.User.ID, cmd.TeamID, cmd.Role)
	if err != nil {
		return err
	}

	role, err := uow.Role.GetByName(ctx, cmd.Role)
	if err != nil {
		return err
	}

	membership.RoleID = role.ID

	_, err = uow.Membership.Update(membership, tx)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func UpdateTeamAvatarWrapper(ctx context.Context, uow *service.UnitOfWork, mailer worker.WorkerInterface, cmd interface{}) error {
	if c, ok := cmd.(*command.UpdateTeamAvatar); ok {
		return UpdateTeamAvatar(ctx, uow, c)
	}
	return fmt.Errorf("invalid command type, expected *command.UpdateTeamAvatar, got %T", cmd)
}

func UpdateTeamAvatar(ctx context.Context, uow *service.UnitOfWork, cmd *command.UpdateTeamAvatar) error {
	tx, txErr := uow.Begin(ctx)
	if txErr != nil {
		return txErr
	}

	defer func() {
		tx.Rollback(ctx)
	}()

	team, err := uow.Team.Get(cmd.TeamID)
	if err != nil {
		return err
	}

	fileContentType := cmd.File.Header.Get("Content-Type")

	// check file type, only allow image
	if !strings.HasPrefix(fileContentType, "image") {
		return exception.NewBadRequestException("invalid file type, only allow image")
	}

	if cmd.File.Size > config.StorageConfig.StaticMaxAvatarSize {
		return exception.NewBadRequestException("file size too large")
	}

	if team.AvatarURL != "" {
		paths := strings.Split(team.AvatarURL, "/")
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

	team.Update(payload)
	_, err = uow.Team.Update(team, tx)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func DeleteTeamAvatarWrapper(ctx context.Context, uow *service.UnitOfWork, mailer worker.WorkerInterface, cmd interface{}) error {
	if c, ok := cmd.(*command.DeleteTeamAvatar); ok {
		return DeleteTeamAvatar(ctx, uow, c)
	}
	return fmt.Errorf("invalid command type, expected *command.DeleteTeamAvatar, got %T", cmd)
}

func DeleteTeamAvatar(ctx context.Context, uow *service.UnitOfWork, cmd *command.DeleteTeamAvatar) error {
	tx, txErr := uow.Begin(ctx)
	if txErr != nil {
		return txErr
	}

	defer func() {
		tx.Rollback(ctx)
	}()

	team, err := uow.Team.Get(cmd.TeamID)
	if err != nil {
		return err
	}

	if team.AvatarURL != "" {
		paths := strings.Split(team.AvatarURL, "/")
		path := filepath.Join(config.StorageConfig.StaticRoot, config.StorageConfig.StaticAvatarPath, paths[len(paths)-1])
		if err := util.DeleteFileInLocal(path); err != nil {
			return err
		}

		team.AvatarURL = ""
		_, err = uow.Team.Update(team, tx)
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
