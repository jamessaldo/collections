package repository

import (
	"authorization/controller/exception"
	"authorization/domain"
	"errors"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type membershipRepository struct {
	db *gorm.DB
}

type MembershipRepository interface {
	Add(*domain.Membership, *gorm.DB) (*domain.Membership, error)
	AddBatch([]domain.Membership) error
	Update(*domain.Membership, *gorm.DB) (*domain.Membership, error)
	Get(uuid.UUID) (*domain.Membership, error)
	List(opts *domain.MembershipOptions) ([]domain.Membership, error)
	Delete(uuid.UUID, *gorm.DB) error
	Count(opts *domain.MembershipOptions) (int64, error)
}

// membershipRepository implements the MembershipRepository interface
func NewMembershipRepository(db *gorm.DB) MembershipRepository {
	return &membershipRepository{db: db}
}

func (repo *membershipRepository) Add(membership *domain.Membership, tx *gorm.DB) (*domain.Membership, error) {
	err := tx.Omit("Team.Creator").Create(&membership).Error
	if err != nil {
		return nil, err
	}
	return membership, nil
}

// add batch gorm
func (repo *membershipRepository) AddBatch(memberships []domain.Membership) error {
	err := repo.db.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(memberships, 1000).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *membershipRepository) Update(membership *domain.Membership, tx *gorm.DB) (*domain.Membership, error) {
	err := tx.Save(&membership).Error
	if err != nil {
		return nil, err
	}
	return membership, nil
}

func (repo *membershipRepository) Get(id uuid.UUID) (*domain.Membership, error) {
	var membership domain.Membership
	err := repo.db.Preload("Role").Where("id = ?", id).First(&membership).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.NewNotFoundException(err.Error())
		}
		return nil, err
	}
	return &membership, nil
}

func (repo *membershipRepository) List(opts *domain.MembershipOptions) ([]domain.Membership, error) {
	db := repo.db

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

	var memberships []domain.Membership
	err := db.Order("last_active_at DESC").Find(&memberships).Error
	if err != nil {
		return nil, err
	}
	return memberships, nil
}

func (repo *membershipRepository) Delete(id uuid.UUID, tx *gorm.DB) error {
	err := tx.Where("id = ?", id).Delete(&domain.Membership{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *membershipRepository) Count(opts *domain.MembershipOptions) (int64, error) {
	db := repo.db.Model(&domain.Membership{})

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
