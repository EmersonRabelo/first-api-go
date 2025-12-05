package repository

import (
	entity "github.com/EmersonRabelo/first-api-go/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostRepository interface {
	Create(post *entity.Post)
	FindById(id uuid.UUID) (*entity.Post, error)
	FindAll(page, pageSize int) ([]entity.Post, int64, error)
	Update(post *entity.Post) error
	Delete(id uuid.UUID) error
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

func (p *postRepository) Create(post *entity.Post) {
	panic("unimplemented")
}

func (p *postRepository) FindById(uuid uuid.UUID) (*entity.Post, error) {
	panic("unimplemented")
}

func (p *postRepository) FindAll(page int, pageSize int) ([]entity.Post, int64, error) {
	panic("unimplemented")
}

func (p *postRepository) Update(post *entity.Post) error {
	panic("unimplemented")
}

func (p *postRepository) Delete(id uuid.UUID) error {
	panic("unimplemented")
}
