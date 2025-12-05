package repository

import (
	"github.com/EmersonRabelo/first-api-go/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LikeRepository interface {
	Create(post *entity.Like)
	FindById(id uuid.UUID) (*entity.Like, error)
	FindAll(page, pageSize int) ([]entity.Like, int64, error)
	Update(like *entity.Like) error
	Delete(id uuid.UUID) error
}

type likeRepository struct {
	db *gorm.DB
}

func NewLikeRepository(db *gorm.DB) LikeRepository {
	return &likeRepository{db: db}
}

// Create implements LikeRepository.
func (l *likeRepository) Create(post *entity.Like) {
	panic("unimplemented")
}

// Delete implements LikeRepository.
func (l *likeRepository) Delete(id uuid.UUID) error {
	panic("unimplemented")
}

// FindAll implements LikeRepository.
func (l *likeRepository) FindAll(page int, pageSize int) ([]entity.Like, int64, error) {
	panic("unimplemented")
}

// FindById implements LikeRepository.
func (l *likeRepository) FindById(id uuid.UUID) (*entity.Like, error) {
	panic("unimplemented")
}

// Update implements LikeRepository.
func (l *likeRepository) Update(like *entity.Like) error {
	panic("unimplemented")
}
