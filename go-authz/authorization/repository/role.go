package repository

import (
	"authorization/controller/exception"
	"authorization/domain"
	"errors"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type roleRepository struct {
	db *gorm.DB
}

// roleRepository implements the RoleRepository interface
type RoleRepository interface {
	Add(*domain.Role, *gorm.DB) (*domain.Role, error)
	Update(*domain.Role, *gorm.DB) (*domain.Role, error)
	Get(domain.RoleType) (*domain.Role, error)
	List() ([]domain.Role, error)
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (repo *roleRepository) Add(role *domain.Role, tx *gorm.DB) (*domain.Role, error) {
	err := tx.Clauses(clause.OnConflict{DoNothing: true}).Create(&role).Error
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (repo *roleRepository) Update(role *domain.Role, tx *gorm.DB) (*domain.Role, error) {
	err := tx.Save(&role).Error
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (repo *roleRepository) Get(name domain.RoleType) (*domain.Role, error) {
	var role domain.Role
	err := repo.db.Where("name = ?", name).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exception.NewNotFoundException(fmt.Sprintf("Role with name %s is not exist! Detail: %s", name, err.Error()))
		}
		return nil, err
	}
	return &role, nil
}

func (repo *roleRepository) List() ([]domain.Role, error) {
	var roles []domain.Role
	err := repo.db.Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}
