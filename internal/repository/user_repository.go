package repository

import (
	entity "github.com/EmersonRabelo/first-api-go/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *entity.User)
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

// Create implements UserRepository.
func (u *userRepository) Create(user *entity.User) {
	panic("unimplemented")
}

// FindAll implements UserRepository.
func (u *userRepository) FindAll(page int, pageSize int) ([]entity.User, int64, error) {
	panic("unimplemented")
}

// FindById implements UserRepository.
func (u *userRepository) FindById(id uuid.UUID) (*entity.User, error) {
	panic("unimplemented")
}

func (u *userRepository) Update(user *entity.User) error {
	panic("unimplemented")
}

func (u *userRepository) Delete(id uuid.UUID) error {
	panic("unimplemented")
}
