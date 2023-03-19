package repository

import (
	"auth/domain/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type roleRepository struct {
	db *gorm.DB
	tx *gorm.DB
}

// roleRepository implements the RoleRepository interface
type RoleRepository interface {
	Add(*model.Role) (*model.Role, error)
	Update(*model.Role) (*model.Role, error)
	Get(model.RoleType) (*model.Role, error)
	List() ([]model.Role, error)
	WithTrx(*gorm.DB) *roleRepository
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (repo *roleRepository) Add(role *model.Role) (*model.Role, error) {
	err := repo.tx.Debug().Clauses(clause.OnConflict{DoNothing: true}).Create(&role).Error
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (repo *roleRepository) Update(role *model.Role) (*model.Role, error) {
	err := repo.tx.Debug().Save(&role).Error
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

func (repo *roleRepository) WithTrx(tx *gorm.DB) *roleRepository {
	repo.tx = tx
	return repo
}
