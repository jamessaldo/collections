package handlers

import (
	"auth/controller/exception"
	"auth/domain/command"
	"auth/domain/dto"
	"auth/domain/model"
	"auth/service"
	"errors"
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func CreateTeam(uow *service.UnitOfWork, cmd *command.CreateTeam) (*dto.TeamRetrievalSchema, error) {
	tx, txErr := uow.Begin(&gorm.Session{})
	if txErr != nil {
		return nil, txErr
	}

	defer func() {
		tx.Rollback()
	}()

	_, err := uow.Team.Get(cmd.TeamID)
	if err == nil {
		return nil, exception.NewConflictException("Team already exists")
	}

	ownerRole, err := uow.Role.Get(model.Owner)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.NewNotFoundException(err.Error())
		}
		return nil, err
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

	uow.Membership.Add(membership)

	tx.Commit()

	return &dto.TeamRetrievalSchema{
		ID:          team.ID,
		Name:        team.Name,
		Description: team.Description,
		AvatarURL:   team.AvatarURL,
		IsPersonal:  team.IsPersonal,
		Creator:     cmd.User.PublicUser(),
		CreatedAt:   team.CreatedAt,
		UpdatedAt:   team.UpdatedAt,
	}, nil
}

func UpdateTeam(uow *service.UnitOfWork, cmd *command.UpdateTeam) (*dto.TeamRetrievalSchema, error) {
	tx, txErr := uow.Begin(&gorm.Session{})
	if txErr != nil {
		return nil, txErr
	}

	defer func() {
		tx.Rollback()
	}()

	team, err := uow.Team.Get(cmd.TeamID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.NewNotFoundException(err.Error())
		}
		return nil, err
	}

	if cmd.Name != "" {
		team.Name = cmd.Name
	}
	if cmd.Description != "" {
		team.Description = cmd.Description
	}

	team, err = uow.Team.Update(team)
	if err != nil {
		return nil, err
	}

	tx.Commit()

	return &dto.TeamRetrievalSchema{
		ID:          team.ID,
		Name:        team.Name,
		Description: team.Description,
		AvatarURL:   team.AvatarURL,
		IsPersonal:  team.IsPersonal,
		Creator:     team.Creator.PublicUser(),
		CreatedAt:   team.CreatedAt,
		UpdatedAt:   team.UpdatedAt,
	}, nil
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
