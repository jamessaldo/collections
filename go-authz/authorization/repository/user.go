package repository

import (
	"authorization/domain/model"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userRepository struct {
	db *gorm.DB
}

// userRepository implements the UserRepository interface
type UserRepository interface {
	Add(*model.User, *gorm.DB) (*model.User, error)
	AddBatch([]model.User) error
	Update(*model.User, *gorm.DB) (*model.User, error)
	Get(uuid.UUID) (*model.User, error)
	List(page, pageSize int) (model.Users, error)
	GetByEmail(string) (*model.User, error)
	GetByUsername(string) (*model.User, error)
	Count() (int64, error)
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (repo *userRepository) Add(user *model.User, tx *gorm.DB) (*model.User, error) {
	err := tx.Create(&user).Error
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

func (repo *userRepository) Update(user *model.User, tx *gorm.DB) (*model.User, error) {
	err := tx.Save(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *userRepository) Get(id uuid.UUID) (*model.User, error) {
	var user model.User
	// get user by id and isActive = true
	err := repo.db.Where(&model.User{ID: id, IsActive: true}).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *userRepository) List(page, pageSize int) (model.Users, error) {
	var users []model.User
	offset := (page - 1) * pageSize
	err := repo.db.Limit(pageSize).Offset(offset).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (repo *userRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User
	err := repo.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *userRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	err := repo.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *userRepository) Count() (int64, error) {
	var count int64
	err := repo.db.Model(&model.User{}).Select("id").Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
