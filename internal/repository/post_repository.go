package repository

import (
	"time"

	entity "github.com/EmersonRabelo/first-api-go/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostRepository interface {
	Create(post *entity.Post) error
	FindById(id *uuid.UUID) (*entity.Post, error)
	FindAll(id *uuid.UUID, start, end time.Time, page, pageSize int) ([]entity.Post, int64, error)
	Update(post *entity.Post) error
	Delete(id *uuid.UUID) error
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

func (p *postRepository) Create(post *entity.Post) error {
	return p.db.Create(post).Error
}

func (p *postRepository) FindById(id *uuid.UUID) (*entity.Post, error) {
	var post entity.Post

	if err := p.db.First(&post, id).Error; err != nil {
		return nil, err
	}

	return &post, nil
}

func (p *postRepository) FindAll(id *uuid.UUID, start, end time.Time, page, pageSize int) ([]entity.Post, int64, error) {
	var posts []entity.Post
	var amount int64 // Total sem filtrar

	p.db.Model(&entity.Post{}).Count(&amount)

	offset := (page - 1) * pageSize

	db := p.db.Model(&entity.Post{})

	if !start.IsZero() {
		db = db.Where("created_at >= ?", start)
	}

	if !end.IsZero() {
		db = db.Where("created_at <= ?", end)
	}

	if id != nil {
		db = db.Where("user_id = ?", id)
	}

	result := db.Offset(offset).Limit(pageSize).Find(&posts)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	return posts, amount, nil
}

func (p *postRepository) Update(post *entity.Post) error {
	return p.db.Save(post).Error
}

func (p *postRepository) Delete(id *uuid.UUID) error {
	return p.db.Delete(&entity.Post{}, id).Error
}
