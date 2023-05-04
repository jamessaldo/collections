package repository

import (
	"authorization/domain/model"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type teamRepository struct {
	db *gorm.DB
}

// teamRepository implements the TeamRepository interface
type TeamRepository interface {
	Add(*model.Team, *gorm.DB) (*model.Team, error)
	Update(*model.Team, *gorm.DB) (*model.Team, error)
	Get(uuid.UUID) (*model.Team, error)
}

func NewTeamRepository(db *gorm.DB) TeamRepository {
	return &teamRepository{db: db}
}

func (repo *teamRepository) Add(team *model.Team, tx *gorm.DB) (*model.Team, error) {
	err := tx.Preload("Creator").Create(&team).Error
	if err != nil {
		return nil, err
	}
	return team, nil
}

func (repo *teamRepository) Update(team *model.Team, tx *gorm.DB) (*model.Team, error) {
	err := tx.Save(&team).Error
	if err != nil {
		return nil, err
	}
	return team, nil
}

func (repo *teamRepository) Get(id uuid.UUID) (*model.Team, error) {
	var team model.Team
	err := repo.db.Preload("Creator").Where("id = ?", id).First(&team).Error
	if err != nil {
		return nil, err
	}
	return &team, nil
}
