package repository

import (
	"auth/domain/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type roleRepository struct {
	db *gorm.DB
}

// roleRepository implements the RoleRepository interface
type RoleRepository interface {
	Add(*model.Role, *gorm.DB) (*model.Role, error)
	Update(*model.Role, *gorm.DB) (*model.Role, error)
	Get(model.RoleType) (*model.Role, error)
	List() ([]model.Role, error)
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (repo *roleRepository) Add(role *model.Role, tx *gorm.DB) (*model.Role, error) {
	err := tx.Debug().Clauses(clause.OnConflict{DoNothing: true}).Create(&role).Error
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (repo *roleRepository) Update(role *model.Role, tx *gorm.DB) (*model.Role, error) {
	err := tx.Debug().Save(&role).Error
	if err != nil {
		return nil, err
	}
	return role, nil
}

func (repo *roleRepository) Get(name model.RoleType) (*model.Role, error) {
	var role model.Role
	err := repo.db.Debug().Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (repo *roleRepository) List() ([]model.Role, error) {
	var roles []model.Role
	err := repo.db.Debug().Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}
