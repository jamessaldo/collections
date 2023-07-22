package repository

import (
	"authorization/controller/exception"
	"authorization/domain"
	"errors"

	"github.com/oklog/ulid/v2"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type invitationRepository struct {
	db *gorm.DB
}

type InvitationRepository interface {
	Add(*domain.Invitation, *gorm.DB) (*domain.Invitation, error)
	AddBatch([]domain.Invitation) error
	Update(*domain.Invitation, *gorm.DB) (*domain.Invitation, error)
	Get(ulid.ULID) (*domain.Invitation, error)
	List(opts *domain.InvitationOptions) ([]domain.Invitation, error)
	Delete(ulid.ULID, *gorm.DB) error
}

// invitationRepository implements the InvitationRepository interface
func NewInvitationRepository(db *gorm.DB) InvitationRepository {
	return &invitationRepository{db: db}
}

func (repo *invitationRepository) Add(invitation *domain.Invitation, tx *gorm.DB) (*domain.Invitation, error) {
	err := tx.Create(&invitation).Error
	if err != nil {
		return nil, err
	}
	return invitation, nil
}

// add batch gorm
func (repo *invitationRepository) AddBatch(invitations []domain.Invitation) error {
	err := repo.db.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(invitations, 1000).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *invitationRepository) Update(invitation *domain.Invitation, tx *gorm.DB) (*domain.Invitation, error) {
	err := tx.Save(&invitation).Error
	if err != nil {
		return nil, err
	}
	return invitation, nil
}

func (repo *invitationRepository) Get(id ulid.ULID) (*domain.Invitation, error) {
	var invitation domain.Invitation
	err := repo.db.Where("id = ?", id).Preload("Role").Preload("Team").First(&invitation).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.NewNotFoundException(err.Error())
		}
		return nil, err
	}
	return &invitation, nil
}

func (repo *invitationRepository) List(opts *domain.InvitationOptions) ([]domain.Invitation, error) {
	db := repo.db.Preload("Role")

	if len(opts.Statuses) > 0 {
		db = db.Where("status IN (?)", opts.Statuses)
	}
	if opts.Email != "" {
		db = db.Where("email = ?", opts.Email)
	}
	if opts.TeamID != uuid.Nil {
		db = db.Where("team_id = ?", opts.TeamID)
	}
	if !opts.ExpiresAt.IsZero() {
		db = db.Where("expires_at <= ?", opts.ExpiresAt)
	}
	if opts.RoleID.String() != "" {
		db = db.Where("role_id = ?", opts.RoleID)
	}
	if opts.Limit > 0 {
		db = db.Limit(opts.Limit)
	}

	var invitations []domain.Invitation
	err := db.Find(&invitations).Error
	if err != nil {
		return nil, err
	}
	return invitations, nil
}

func (repo *invitationRepository) Delete(id ulid.ULID, tx *gorm.DB) error {
	err := tx.Where("id = ?", id).Delete(&domain.Invitation{}).Error
	if err != nil {
		return err
	}
	return nil
}
