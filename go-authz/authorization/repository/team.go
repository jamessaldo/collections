package repository

import (
	"auth/domain/model"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type teamRepository struct {
	db *gorm.DB
	tx *gorm.DB
}

// teamRepository implements the TeamRepository interface
type TeamRepository interface {
	Add(*model.Team) (*model.Team, error)
	Update(*model.Team) (*model.Team, error)
	Get(uuid.UUID) (*model.Team, error)
	WithTrx(*gorm.DB) *teamRepository
}

func NewTeamRepository(db *gorm.DB) TeamRepository {
	return &teamRepository{db: db}
}

func (repo *teamRepository) Add(team *model.Team) (*model.Team, error) {
	err := repo.tx.Debug().Preload("Creator").Create(&team).Error
	if err != nil {
		return nil, err
	}
	return team, nil
}

func (repo *teamRepository) Update(team *model.Team) (*model.Team, error) {
	err := repo.tx.Debug().Save(&team).Error
	if err != nil {
		return nil, err
	}
	return team, nil
}

func (repo *teamRepository) Get(id uuid.UUID) (*model.Team, error) {
	var team model.Team
	err := repo.db.Debug().Preload("Creator").Where("id = ?", id).First(&team).Error
	if err != nil {
		return nil, err
	}
	return &team, nil
}

func (repo *teamRepository) WithTrx(tx *gorm.DB) *teamRepository {
	repo.tx = tx
	return repo
}
