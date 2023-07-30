package handlers

import (
	"authorization/config"
	"authorization/controller/exception"
	"authorization/domain"
	"authorization/domain/command"
	"authorization/infrastructure/persistence"
	"authorization/repository"
	"authorization/util"
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/oklog/ulid/v2"
)

func CreateTeam(ctx context.Context, cmd *command.CreateTeam) error {
	tx, txErr := persistence.Pool.Begin(ctx)
	if txErr != nil {
		return txErr
	}

	defer func() {
		tx.Rollback(ctx)
	}()

	ownerRole, err := repository.Role.GetByName(ctx, domain.Owner)
	if err != nil {
		return err
	}

	team := domain.NewTeam(cmd.User, ownerRole.ID, cmd.Name, cmd.Description, false)
	_, err = repository.Team.Add(ctx, team, tx)
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

func UpdateTeam(ctx context.Context, cmd *command.UpdateTeam) error {
	tx, txErr := persistence.Pool.Begin(ctx)
	if txErr != nil {
		return txErr
	}

	defer func() {
		tx.Rollback(ctx)
	}()

	team, err := repository.Team.Get(ctx, cmd.TeamID)
	if err != nil {
		return err
	}

	team.Update(map[string]interface{}{
		"name":        cmd.Name,
		"description": cmd.Description,
	})

	_, err = repository.Team.Update(ctx, team, tx)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func UpdateLastActiveTeam(ctx context.Context, cmd *command.UpdateLastActiveTeam) error {
	tx, txErr := persistence.Pool.Begin(ctx)
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
	memberships, err := repository.Membership.List(ctx, opts)
	if err != nil {
		return err
	}

	if len(memberships) == 0 {
		return exception.NewNotFoundException("Team is not found")
	}

	lastActiveAt := util.GetTimestampUTC()
	membership := memberships[0]
	membership.LastActiveAt = lastActiveAt

	_, err = repository.Membership.Update(ctx, membership, tx)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func DeleteTeamMember(ctx context.Context, cmd *command.DeleteTeamMember) error {
	tx, txErr := persistence.Pool.Begin(ctx)
	if txErr != nil {
		return txErr
	}

	defer func() {
		tx.Rollback(ctx)
	}()

	membership, err := repository.Membership.Get(ctx, cmd.MembershipID)
	if err != nil {
		return err
	}

	err = membership.Validation(cmd.User.ID, cmd.TeamID, "")
	if err != nil {
		return err
	}

	err = repository.Membership.Delete(ctx, cmd.MembershipID, tx)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func ChangeMemberRole(ctx context.Context, cmd *command.ChangeMemberRole) error {
	tx, txErr := persistence.Pool.Begin(ctx)
	if txErr != nil {
		return txErr
	}

	defer func() {
		tx.Rollback(ctx)
	}()

	membership, err := repository.Membership.Get(ctx, cmd.MembershipID)
	if err != nil {
		return err
	}

	err = membership.Validation(cmd.User.ID, cmd.TeamID, cmd.Role)
	if err != nil {
		return err
	}

	role, err := repository.Role.GetByName(ctx, cmd.Role)
	if err != nil {
		return err
	}

	membership.RoleID = role.ID

	_, err = repository.Membership.Update(ctx, membership, tx)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func UpdateTeamAvatar(ctx context.Context, cmd *command.UpdateTeamAvatar) error {
	tx, txErr := persistence.Pool.Begin(ctx)
	if txErr != nil {
		return txErr
	}

	defer func() {
		tx.Rollback(ctx)
	}()

	team, err := repository.Team.Get(ctx, cmd.TeamID)
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
	_, err = repository.Team.Update(ctx, team, tx)
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func DeleteTeamAvatar(ctx context.Context, cmd *command.DeleteTeamAvatar) error {
	tx, txErr := persistence.Pool.Begin(ctx)
	if txErr != nil {
		return txErr
	}

	defer func() {
		tx.Rollback(ctx)
	}()

	team, err := repository.Team.Get(ctx, cmd.TeamID)
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
		_, err = repository.Team.Update(ctx, team, tx)
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
