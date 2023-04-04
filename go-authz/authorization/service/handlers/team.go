package handlers

import (
	"auth/controller/exception"
	"auth/domain/command"
	"auth/domain/model"
	"auth/infrastructure/worker"
	"auth/service"
	"errors"
	"fmt"
	"time"

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
		IsPersonal:  cmd.IsPersonal,
		CreatorID:   cmd.User.ID,
	}

	membership := &model.Membership{
		ID:     uuid.NewV4(),
		TeamID: cmd.TeamID,
		Team:   team,
		UserID: cmd.User.ID,
		RoleID: ownerRole.ID,
	}

	_, err = uow.Membership.Add(membership)
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

	_, err = uow.Team.Update(team)
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

	_, err = uow.Membership.Update(&membership)
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

	err = uow.Membership.Delete(cmd.MembershipID)
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

	_, err = uow.Membership.Update(membership)
	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}
