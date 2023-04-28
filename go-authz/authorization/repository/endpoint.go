package repository

import (
	"auth/domain/model"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type endpointRepository struct {
	db *gorm.DB
}

// endpointRepository implements the EndpointRepository interface
type EndpointRepository interface {
	Add(*model.Endpoint, *gorm.DB) (*model.Endpoint, error)
	Update(*model.Endpoint, *gorm.DB) (*model.Endpoint, error)
	Get(uuid.UUID) (*model.Endpoint, error)
	List() ([]model.Endpoint, error)
}

func NewEndpointRepository(db *gorm.DB) EndpointRepository {
	return &endpointRepository{db: db}
}

func (repo *endpointRepository) Add(endpoint *model.Endpoint, tx *gorm.DB) (*model.Endpoint, error) {
	err := tx.Debug().Create(&endpoint).Error
	if err != nil {
		return nil, err
	}
	return endpoint, nil
}

func (repo *endpointRepository) Update(endpoint *model.Endpoint, tx *gorm.DB) (*model.Endpoint, error) {
	err := tx.Debug().Save(&endpoint).Error
	if err != nil {
		return nil, err
	}
	return endpoint, nil
}

func (repo *endpointRepository) Get(id uuid.UUID) (*model.Endpoint, error) {
	var endpoint model.Endpoint
	err := repo.db.Debug().Where("id = ?", id).First(&endpoint).Error
	if err != nil {
		return nil, err
	}
	return &endpoint, nil
}

func (repo *endpointRepository) List() ([]model.Endpoint, error) {
	var endpoints []model.Endpoint
	err := repo.db.Debug().Find(&endpoints).Error
	if err != nil {
		return nil, err
	}
	return endpoints, nil
}
