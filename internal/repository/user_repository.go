package repository

import (
	entity "github.com/EmersonRabelo/first-api-go/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *entity.User) error
	FindById(id uuid.UUID) (*entity.User, error)
	FindAll(page, pageSize int) ([]entity.User, int64, error)
	Update(user *entity.User) error
	Delete(id uuid.UUID) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (u *userRepository) Create(user *entity.User) error {
	return u.db.Create(user).Error
}

func (u *userRepository) FindAll(page int, pageSize int) ([]entity.User, int64, error) {
	var users []entity.User
	var amount int64

	u.db.Model(&entity.User{}).Count(&amount)

	offset := (page - 1) * pageSize

	if err := u.db.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, amount, nil

}

func (u *userRepository) FindById(id uuid.UUID) (*entity.User, error) {
	var user entity.User

	if err := u.db.First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *userRepository) Update(user *entity.User) error {
	return u.db.Save(user).Error
}

func (u *userRepository) Delete(id uuid.UUID) error {
	return u.db.Delete(&entity.User{}, id).Error
}
