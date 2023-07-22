package repository

import (
	"authorization/domain"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type endpointRepository struct {
	db *gorm.DB
}

// endpointRepository implements the EndpointRepository interface
type EndpointRepository interface {
	Add(*domain.Endpoint, *gorm.DB) (*domain.Endpoint, error)
	Update(*domain.Endpoint, *gorm.DB) (*domain.Endpoint, error)
	Get(uuid.UUID) (*domain.Endpoint, error)
	List() ([]domain.Endpoint, error)
	ListFilteredBy(teamId, userId uuid.UUID) ([]domain.Endpoint, error)
}

func NewEndpointRepository(db *gorm.DB) EndpointRepository {
	return &endpointRepository{db: db}
}

func (repo *endpointRepository) Add(endpoint *domain.Endpoint, tx *gorm.DB) (*domain.Endpoint, error) {
	err := tx.Create(&endpoint).Error
	if err != nil {
		return nil, err
	}
	return endpoint, nil
}

func (repo *endpointRepository) Update(endpoint *domain.Endpoint, tx *gorm.DB) (*domain.Endpoint, error) {
	err := tx.Save(&endpoint).Error
	if err != nil {
		return nil, err
	}
	return endpoint, nil
}

func (repo *endpointRepository) Get(id uuid.UUID) (*domain.Endpoint, error) {
	var endpoint domain.Endpoint
	err := repo.db.Where("id = ?", id).First(&endpoint).Error
	if err != nil {
		return nil, err
	}
	return &endpoint, nil
}

func (repo *endpointRepository) List() ([]domain.Endpoint, error) {
	var endpoints []domain.Endpoint
	err := repo.db.Find(&endpoints).Error
	if err != nil {
		return nil, err
	}
	return endpoints, nil
}

func (repo *endpointRepository) ListFilteredBy(teamId, userId uuid.UUID) ([]domain.Endpoint, error) {
	var endpoints []domain.Endpoint
	err := repo.db.Raw("SELECT e as endpoint FROM memberships m LEFT JOIN accesses a ON a.role_id = m.role_id LEFT JOIN endpoints e ON e.id = a.endpoint_id WHERE m.team_id = ? AND m.user_id = ?", teamId, userId).Scan(&endpoints).Error
	if err != nil {
		return nil, err
	}
	return endpoints, nil
}
