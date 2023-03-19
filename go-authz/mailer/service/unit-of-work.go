package service

import (
	"context"

	"mailer/repository"

	"gorm.io/gorm"
)

type UoW interface {
	GetDB() *gorm.DB
	Begin() (*gorm.DB, error)
}

type UnitOfWork struct {
	db         *gorm.DB
	ctx        context.Context
	Invitation repository.InvitationRepository
}

func NewUnitOfWork(db *gorm.DB) (*UnitOfWork, error) {
	ctx := context.Background()

	return &UnitOfWork{
		db:         db,
		ctx:        ctx,
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
