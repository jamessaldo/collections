package handlers

import (
	"authorization/config"
	"authorization/controller/exception"
	"authorization/domain/command"
	"authorization/domain/model"
	"authorization/infrastructure/worker"
	"authorization/service"
	"authorization/util"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func CreateTeamWrapper(uow *service.UnitOfWork, mailer worker.WorkerInterface, cmd interface{}) error {
	if c, ok := cmd.(*command.CreateTeam); ok {
		return CreateTeam(uow, c)
	}
	return fmt.Errorf("invalid command type, expected *command.CreateTeam, got %T", cmd)
}

func CreateTeam(uow *service.UnitOfWork, cmd *command.CreateTeam) error {
	tx, txErr := uow.Begin(&gorm.Session{})
	if txErr != nil {
		return txErr
	}

	defer func() {
		tx.Rollback()
	}()

	_, err := uow.Team.Get(cmd.TeamID)
	if err == nil {
		return exception.NewConflictException("Team already exists")
	}

	ownerRole, err := uow.Role.Get(model.Owner)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return exception.NewNotFoundException(err.Error())
		}
		return err
	}

	team := &model.Team{
		ID:          cmd.TeamID,
		Name:        cmd.Name,
		Description: cmd.Description,
		IsPersonal:  false,
		CreatorID:   cmd.User.ID,
	}

	membership := &model.Membership{
		ID:     uuid.NewV4(),
		TeamID: cmd.TeamID,
		Team:   team,
		UserID: cmd.User.ID,
		RoleID: ownerRole.ID,
	}

	_, err = uow.Membership.Add(membership, tx)
	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}

func UpdateTeamWrapper(uow *service.UnitOfWork, mailer worker.WorkerInterface, cmd interface{}) error {
	if c, ok := cmd.(*command.UpdateTeam); ok {
		return UpdateTeam(uow, c)
	}
	return fmt.Errorf("invalid command type, expected *command.UpdateTeam, got %T", cmd)
}

func UpdateTeam(uow *service.UnitOfWork, cmd *command.UpdateTeam) error {
	tx, txErr := uow.Begin(&gorm.Session{})
	if txErr != nil {
		return txErr
	}

	defer func() {
		tx.Rollback()
	}()

	team, err := uow.Team.Get(cmd.TeamID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return exception.NewNotFoundException(err.Error())
		}
		return err
	}

	if cmd.Name != "" {
		team.Name = cmd.Name
	}
	if cmd.Description != "" {
		team.Description = cmd.Description
	}

	_, err = uow.Team.Update(team, tx)
	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}

func UpdateLastActiveTeamWrapper(uow *service.UnitOfWork, mailer worker.WorkerInterface, cmd interface{}) error {
	if c, ok := cmd.(*command.UpdateLastActiveTeam); ok {
		return UpdateLastActiveTeam(uow, c)
	}
	return fmt.Errorf("invalid command type, expected *command.UpdateLastActiveTeam, got %T", cmd)
}

func UpdateLastActiveTeam(uow *service.UnitOfWork, cmd *command.UpdateLastActiveTeam) error {
	tx, txErr := uow.Begin(&gorm.Session{})
	if txErr != nil {
		return txErr
	}

	defer func() {
		tx.Rollback()
	}()

	opts := &model.MembershipOptions{
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

	lastActiveAt := time.Now()
	membership := memberships[0]
	membership.LastActiveAt = &lastActiveAt

	_, err = uow.Membership.Update(&membership, tx)
	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}

func DeleteTeamMemberWrapper(uow *service.UnitOfWork, mailer worker.WorkerInterface, cmd interface{}) error {
	if c, ok := cmd.(*command.DeleteTeamMember); ok {
		return DeleteTeamMember(uow, c)
	}
	return fmt.Errorf("invalid command type, expected *command.DeleteTeamMember, got %T", cmd)
}

func DeleteTeamMember(uow *service.UnitOfWork, cmd *command.DeleteTeamMember) error {
	tx, txErr := uow.Begin(&gorm.Session{})
	if txErr != nil {
		return txErr
	}

	defer func() {
		tx.Rollback()
	}()

	membership, err := uow.Membership.Get(cmd.MembershipID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return exception.NewNotFoundException(err.Error())
		}
		return err
	}

	if membership.UserID == cmd.User.ID {
		return exception.NewForbiddenException("You cannot delete yourself")
	}

	if membership.TeamID != cmd.TeamID {
		return exception.NewForbiddenException(fmt.Sprintf("Team with ID %s is not match with membership-team ID", cmd.TeamID))
	}

	if membership.Role.Name == model.Owner {
		return exception.NewForbiddenException("You cannot delete owner of the team")
	}

	err = uow.Membership.Delete(cmd.MembershipID, tx)
	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}

func ChangeMemberRoleWrapper(uow *service.UnitOfWork, mailer worker.WorkerInterface, cmd interface{}) error {
	if c, ok := cmd.(*command.ChangeMemberRole); ok {
		return ChangeMemberRole(uow, c)
	}
	return fmt.Errorf("invalid command type, expected *command.ChangeMemberRole, got %T", cmd)
}

func ChangeMemberRole(uow *service.UnitOfWork, cmd *command.ChangeMemberRole) error {
	tx, txErr := uow.Begin(&gorm.Session{})
	if txErr != nil {
		return txErr
	}

	defer func() {
		tx.Rollback()
	}()

	membership, err := uow.Membership.Get(cmd.MembershipID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return exception.NewNotFoundException(err.Error())
		}
		return err
	}

	if membership.UserID == cmd.User.ID {
		return exception.NewForbiddenException("You cannot change your own role")
	}

	if membership.TeamID != cmd.TeamID {
		return exception.NewForbiddenException(fmt.Sprintf("Team with ID %s is not match with membership-team ID", cmd.TeamID))
	}

	if membership.Role.Name == model.Owner {
		return exception.NewForbiddenException("You cannot delete owner of the team")
	}

	if cmd.Role == model.Owner {
		return exception.NewForbiddenException("You cannot change role to owner")
	}

	role, err := uow.Role.Get(cmd.Role)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return exception.NewNotFoundException(err.Error())
		}
		return err
	}

	membership.RoleID = role.ID

	_, err = uow.Membership.Update(membership, tx)
	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}

func UpdateTeamAvatarWrapper(uow *service.UnitOfWork, mailer worker.WorkerInterface, cmd interface{}) error {
	if c, ok := cmd.(*command.UpdateTeamAvatar); ok {
		return UpdateTeamAvatar(uow, c)
	}
	return fmt.Errorf("invalid command type, expected *command.UpdateTeamAvatar, got %T", cmd)
}

func UpdateTeamAvatar(uow *service.UnitOfWork, cmd *command.UpdateTeamAvatar) error {
	tx, txErr := uow.Begin(&gorm.Session{})
	if txErr != nil {
		return txErr
	}

	defer func() {
		tx.Rollback()
	}()

	team, err := uow.Team.Get(cmd.TeamID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return exception.NewNotFoundException(err.Error())
		}
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
		// get avatar path after public url
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

	tx.Commit()

	return nil
}

func DeleteTeamAvatarWrapper(uow *service.UnitOfWork, mailer worker.WorkerInterface, cmd interface{}) error {
	if c, ok := cmd.(*command.DeleteTeamAvatar); ok {
		return DeleteTeamAvatar(uow, c)
	}
	return fmt.Errorf("invalid command type, expected *command.DeleteTeamAvatar, got %T", cmd)
}

func DeleteTeamAvatar(uow *service.UnitOfWork, cmd *command.DeleteTeamAvatar) error {
	tx, txErr := uow.Begin(&gorm.Session{})
	if txErr != nil {
		return txErr
	}

	defer func() {
		tx.Rollback()
	}()

	team, err := uow.Team.Get(cmd.TeamID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return exception.NewNotFoundException(err.Error())
		}
		return err
	}

	if team.AvatarURL != "" {
		// get avatar path after public url
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
		tx.Commit()
	}

	return nil
}
