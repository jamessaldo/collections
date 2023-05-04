package service

import (
	"context"

	"mailer/repository"

	"gorm.io/gorm"
)

type UnitOfWork struct {
	DB         *gorm.DB
	ctx        context.Context
	Invitation repository.InvitationRepository
}

func NewUnitOfWork(db *gorm.DB) (*UnitOfWork, error) {
	ctx := context.Background()

	return &UnitOfWork{
		DB:         db,
		ctx:        ctx,
		Invitation: repository.NewInvitationRepository(db),
	}, nil
}

func (u *UnitOfWork) GetDB() *gorm.DB {
	return u.DB
}

func (u *UnitOfWork) Begin(sessionConfig *gorm.Session) (*gorm.DB, error) {
	tx := u.DB.Session(sessionConfig).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}
