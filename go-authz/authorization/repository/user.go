package repository

import (
	"authorization/domain"

	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type userRepository struct {
	db *gorm.DB
}

// userRepository implements the UserRepository interface
type UserRepository interface {
	Add(*domain.User, *gorm.DB) (*domain.User, error)
	AddBatch([]domain.User) error
	Update(*domain.User, *gorm.DB) (*domain.User, error)
	Get(uuid.UUID) (*domain.User, error)
	List(page, pageSize int) (domain.Users, error)
	GetByEmail(string) (*domain.User, error)
	GetByUsername(string) (*domain.User, error)
	Count() (int64, error)
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (repo *userRepository) Add(user *domain.User, tx *gorm.DB) (*domain.User, error) {
	err := tx.Create(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

// add batch gorm
func (repo *userRepository) AddBatch(users []domain.User) error {
	err := repo.db.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(users, 1000).Error
	if err != nil {
		return err
	}
	return nil
}

func (repo *userRepository) Update(user *domain.User, tx *gorm.DB) (*domain.User, error) {
	err := tx.Save(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *userRepository) Get(id uuid.UUID) (*domain.User, error) {
	var user domain.User
	// get user by id and isActive = true
	err := repo.db.Where(&domain.User{ID: id, IsActive: true}).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *userRepository) List(page, pageSize int) (domain.Users, error) {
	var users []domain.User
	offset := (page - 1) * pageSize
	err := repo.db.Limit(pageSize).Offset(offset).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (repo *userRepository) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := repo.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *userRepository) GetByUsername(username string) (*domain.User, error) {
	var user domain.User
	err := repo.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repo *userRepository) Count() (int64, error) {
	var count int64
	err := repo.db.Model(&domain.User{}).Select("id").Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
