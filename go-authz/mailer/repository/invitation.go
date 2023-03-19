package repository

import (
	"mailer/domain/model"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type invitationRepository struct {
	db *gorm.DB
	tx *gorm.DB
}

type InvitationRepository interface {
	Add(*model.Invitation) (*model.Invitation, error)
	Update(*model.Invitation) (*model.Invitation, error)
	Get(uuid.UUID) (*model.Invitation, error)
	List(opts *model.InvitationOptions) ([]model.Invitation, error)
	Delete(uuid.UUID) error
	WithTrx(*gorm.DB) *invitationRepository
}

// invitationRepository implements the InvitationRepository interface
func NewInvitationRepository(db *gorm.DB) InvitationRepository {
	return &invitationRepository{db: db}
}

func (repo *invitationRepository) Add(invitation *model.Invitation) (*model.Invitation, error) {
	err := repo.tx.Debug().Omit("Team.Creator").Create(&invitation).Error
	if err != nil {
		return nil, err
	}
	return invitation, nil
}

// add batch gorm
func (repo *invitationRepository) AddBatch(invitations []model.Invitation) error {
	err := repo.db.Debug().Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(invitations, 1000).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *invitationRepository) Update(invitation *model.Invitation) (*model.Invitation, error) {
	err := repo.tx.Debug().Save(&invitation).Error
	if err != nil {
		return nil, err
	}
	return invitation, nil
}

func (repo *invitationRepository) Get(id uuid.UUID) (*model.Invitation, error) {
	var invitation model.Invitation
	err := repo.db.Debug().Where("id = ?", id).First(&invitation).Error
	if err != nil {
		return nil, err
	}
	return &invitation, nil
}

func (repo *invitationRepository) List(opts *model.InvitationOptions) ([]model.Invitation, error) {
	db := repo.db.Debug()

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
	if opts.RoleID != uuid.Nil {
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

func (repo *invitationRepository) Delete(id uuid.UUID) error {
	err := repo.tx.Debug().Where("id = ?", id).Delete(&model.Invitation{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *invitationRepository) WithTrx(tx *gorm.DB) *invitationRepository {

	repo.tx = tx
	return repo
}
