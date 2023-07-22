package repository

import (
	"authorization/controller/exception"
	"authorization/domain"
	"errors"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type teamRepository struct {
	db *gorm.DB
}

// teamRepository implements the TeamRepository interface
type TeamRepository interface {
	Add(*domain.Team, *gorm.DB) (*domain.Team, error)
	Update(*domain.Team, *gorm.DB) (*domain.Team, error)
	Get(uuid.UUID) (*domain.Team, error)
}

func NewTeamRepository(db *gorm.DB) TeamRepository {
	return &teamRepository{db: db}
}

func (repo *teamRepository) Add(team *domain.Team, tx *gorm.DB) (*domain.Team, error) {
	err := tx.Preload("Creator").Create(&team).Error
	if err != nil {
		return nil, err
	}
	return team, nil
}

func (repo *teamRepository) Update(team *domain.Team, tx *gorm.DB) (*domain.Team, error) {
	err := tx.Save(&team).Error
	if err != nil {
		return nil, err
	}
	return team, nil
}

func (repo *teamRepository) Get(id uuid.UUID) (*domain.Team, error) {
	var team domain.Team
	err := repo.db.Preload("Creator").Where("id = ?", id).First(&team).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.NewNotFoundException(err.Error())
		}
		return nil, err
	}
	return &team, nil
}
