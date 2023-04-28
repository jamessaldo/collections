package service

import (
	"context"

	"auth/repository"

	"gorm.io/gorm"
)

type UnitOfWork struct {
	db         *gorm.DB
	ctx        context.Context
	User       repository.UserRepository
	Role       repository.RoleRepository
	Endpoint   repository.EndpointRepository
	Team       repository.TeamRepository
	Membership repository.MembershipRepository
	Invitation repository.InvitationRepository
}

func NewUnitOfWork(db *gorm.DB) (*UnitOfWork, error) {
	ctx := context.Background()

	return &UnitOfWork{
		db:         db,
		ctx:        ctx,
		User:       repository.NewUserRepository(db),
		Role:       repository.NewRoleRepository(db),
		Endpoint:   repository.NewEndpointRepository(db),
		Team:       repository.NewTeamRepository(db),
		Membership: repository.NewMembershipRepository(db),
		Invitation: repository.NewInvitationRepository(db),
	}, nil
}

func (u *UnitOfWork) GetDB() *gorm.DB {
	return u.db
}

func (u *UnitOfWork) Begin(sessionConfig *gorm.Session) (*gorm.DB, error) {
	tx := u.db.Session(sessionConfig).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}
