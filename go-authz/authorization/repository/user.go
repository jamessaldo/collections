package repository

import (
	"auth/domain/model"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userRepository struct {
	db *gorm.DB
	tx *gorm.DB
}

// userRepository implements the UserRepository interface
type UserRepository interface {
	Add(*model.User) (*model.User, error)
	AddBatch([]model.User) error
	Update(*model.User) (*model.User, error)
	Get(uuid.UUID) (*model.User, error)
	List(page, pageSize int) (model.Users, error)
	GetByEmail(string) (*model.User, error)
	GetByUsername(string) (*model.User, error)
	Count() (int64, error)
	WithTrx(*gorm.DB) *userRepository
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (repo *userRepository) Add(user *model.User) (*model.User, error) {
	err := repo.tx.Debug().Create(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

// add batch gorm
func (repo *userRepository) AddBatch(users []model.User) error {
	err := repo.db.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(users, 1000).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *userRepository) Update(user *model.User) (*model.User, error) {
	err := repo.tx.Debug().Save(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *userRepository) Get(id uuid.UUID) (*model.User, error) {
	var user model.User
	err := repo.db.Debug().Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *userRepository) List(page, pageSize int) (model.Users, error) {
	var users []model.User
	offset := (page - 1) * pageSize
	err := repo.db.Debug().Limit(pageSize).Offset(offset).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (repo *userRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User
	err := repo.db.Debug().Where("email = ?", email).Take(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *userRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	err := repo.db.Debug().Where("username = ?", username).Take(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *userRepository) Count() (int64, error) {
	var count int64
	err := repo.db.Debug().Model(&model.User{}).Select("id").Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (repo *userRepository) WithTrx(tx *gorm.DB) *userRepository {
	repo.tx = tx
	return repo
}