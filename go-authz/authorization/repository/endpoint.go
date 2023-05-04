package repository

import (
	"authorization/domain/model"

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
	ListFilteredBy(teamId, userId uuid.UUID) ([]model.Endpoint, error)
}

func NewEndpointRepository(db *gorm.DB) EndpointRepository {
	return &endpointRepository{db: db}
}

func (repo *endpointRepository) Add(endpoint *model.Endpoint, tx *gorm.DB) (*model.Endpoint, error) {
	err := tx.Create(&endpoint).Error
	if err != nil {
		return nil, err
	}
	return endpoint, nil
}

func (repo *endpointRepository) Update(endpoint *model.Endpoint, tx *gorm.DB) (*model.Endpoint, error) {
	err := tx.Save(&endpoint).Error
	if err != nil {
		return nil, err
	}
	return endpoint, nil
}

func (repo *endpointRepository) Get(id uuid.UUID) (*model.Endpoint, error) {
	var endpoint model.Endpoint
	err := repo.db.Where("id = ?", id).First(&endpoint).Error
	if err != nil {
		return nil, err
	}
	return &endpoint, nil
}

func (repo *endpointRepository) List() ([]model.Endpoint, error) {
	var endpoints []model.Endpoint
	err := repo.db.Find(&endpoints).Error
	if err != nil {
		return nil, err
	}
	return endpoints, nil
}

func (repo *endpointRepository) ListFilteredBy(teamId, userId uuid.UUID) ([]model.Endpoint, error) {
	var endpoints []model.Endpoint
	err := repo.db.Raw("SELECT e as endpoint FROM memberships m LEFT JOIN accesses a ON a.role_id = m.role_id LEFT JOIN endpoints e ON e.id = a.endpoint_id WHERE m.team_id = ? AND m.user_id = ?", teamId, userId).Scan(&endpoints).Error
	if err != nil {
		return nil, err
	}
	return endpoints, nil
}
