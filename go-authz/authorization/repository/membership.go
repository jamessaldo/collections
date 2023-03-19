package repository

import (
	"auth/domain/model"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type membershipRepository struct {
	db *gorm.DB
	tx *gorm.DB
}

type MembershipRepository interface {
	Add(*model.Membership) (*model.Membership, error)
	Update(*model.Membership) (*model.Membership, error)
	Get(uuid.UUID) (*model.Membership, error)
	List(opts *model.MembershipOptions) ([]model.Membership, error)
	Delete(uuid.UUID) error
	Count(opts *model.MembershipOptions) (int64, error)
	WithTrx(*gorm.DB) *membershipRepository
}

// membershipRepository implements the MembershipRepository interface
func NewMembershipRepository(db *gorm.DB) MembershipRepository {
	return &membershipRepository{db: db}
}

func (repo *membershipRepository) Add(membership *model.Membership) (*model.Membership, error) {
	err := repo.tx.Debug().Omit("Team.Creator").Create(&membership).Error
	if err != nil {
		return nil, err
	}
	return membership, nil
}

// add batch gorm
func (repo *membershipRepository) AddBatch(memberships []model.Membership) error {
	err := repo.db.Debug().Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(memberships, 1000).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *membershipRepository) Update(membership *model.Membership) (*model.Membership, error) {
	err := repo.tx.Debug().Save(&membership).Error
	if err != nil {
		return nil, err
	}
	return membership, nil
}

func (repo *membershipRepository) Get(id uuid.UUID) (*model.Membership, error) {
	var membership model.Membership
	err := repo.db.Debug().Preload("Role").Where("id = ?", id).First(&membership).Error
	if err != nil {
		return nil, err
	}
	return &membership, nil
}

func (repo *membershipRepository) List(opts *model.MembershipOptions) ([]model.Membership, error) {
	db := repo.db.Debug()

	if opts.IsSelectTeam {
		db = db.Preload("Team.Creator")
	}
	if opts.IsSelectUser {
		db = db.Preload("User")
	}
	if opts.IsSelectRole {
		db = db.Preload("Role")
	}

	if opts.TeamID != uuid.Nil && opts.Name == "" {
		db = db.Where("team_id = ?", opts.TeamID)
	}
	if opts.UserID != uuid.Nil {
		db = db.Where("user_id = ?", opts.UserID)
		if opts.Name != "" {
			db = db.Joins("JOIN teams ON memberships.team_id = teams.id AND teams.name LIKE ?", "%"+opts.Name+"%")
		}
	}
	if opts.Limit > 0 {
		db = db.Limit(opts.Limit)
	}
	if opts.Skip > 0 {
		offset := (opts.Skip - 1) * opts.Limit
		db = db.Offset(offset)
	}

	var memberships []model.Membership
	err := db.Order("last_active_at DESC").Find(&memberships).Error
	if err != nil {
		return nil, err
	}
	return memberships, nil
}

func (repo *membershipRepository) Delete(id uuid.UUID) error {
	err := repo.tx.Debug().Where("id = ?", id).Delete(&model.Membership{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *membershipRepository) Count(opts *model.MembershipOptions) (int64, error) {
	db := repo.db.Model(&model.Membership{})

	if opts.TeamID != uuid.Nil && opts.Name == "" {
		db = db.Where("team_id = ?", opts.TeamID)
	}
	if opts.UserID != uuid.Nil {
		db = db.Where("user_id = ?", opts.UserID)
		if opts.Name != "" {
			db = db.Joins("JOIN teams ON memberships.team_id = teams.id AND teams.name LIKE ?", "%"+opts.Name+"%")
		}
	}

	var count int64
	err := db.Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *membershipRepository) WithTrx(tx *gorm.DB) *membershipRepository {
	repo.tx = tx
	return repo
}
