package repository

import (
	"authorization/domain/model"

	"github.com/oklog/ulid/v2"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type invitationRepository struct {
	db *gorm.DB
}

type InvitationRepository interface {
	Add(*model.Invitation, *gorm.DB) (*model.Invitation, error)
	AddBatch([]model.Invitation) error
	Update(*model.Invitation, *gorm.DB) (*model.Invitation, error)
	Get(ulid.ULID) (*model.Invitation, error)
	List(opts *model.InvitationOptions) ([]model.Invitation, error)
	Delete(ulid.ULID, *gorm.DB) error
}

// invitationRepository implements the InvitationRepository interface
func NewInvitationRepository(db *gorm.DB) InvitationRepository {
	return &invitationRepository{db: db}
}

func (repo *invitationRepository) Add(invitation *model.Invitation, tx *gorm.DB) (*model.Invitation, error) {
	err := tx.Create(&invitation).Error
	if err != nil {
		return nil, err
	}
	return invitation, nil
}

// add batch gorm
func (repo *invitationRepository) AddBatch(invitations []model.Invitation) error {
	err := repo.db.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(invitations, 1000).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *invitationRepository) Update(invitation *model.Invitation, tx *gorm.DB) (*model.Invitation, error) {
	err := tx.Save(&invitation).Error
	if err != nil {
		return nil, err
	}
	return invitation, nil
}

func (repo *invitationRepository) Get(id ulid.ULID) (*model.Invitation, error) {
	var invitation model.Invitation
	err := repo.db.Where("id = ?", id).Preload("Role").Preload("Team").First(&invitation).Error
	if err != nil {
		return nil, err
	}
	return &invitation, nil
}

func (repo *invitationRepository) List(opts *model.InvitationOptions) ([]model.Invitation, error) {
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

	var invitations []model.Invitation
	err := db.Find(&invitations).Error
	if err != nil {
		return nil, err
	}
	return invitations, nil
}

func (repo *invitationRepository) Delete(id ulid.ULID, tx *gorm.DB) error {
	err := tx.Where("id = ?", id).Delete(&model.Invitation{}).Error
	if err != nil {
		return err
	}
	return nil
}
