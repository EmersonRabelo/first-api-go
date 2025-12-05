package repository

import (
	"github.com/EmersonRabelo/first-api-go/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReplyRepository interface {
	Create(post *entity.Reply)
	FindById(id uuid.UUID) (*entity.Reply, error)
	FindAll(page, pageSize int) ([]entity.Reply, int64, error)
	Update(reply *entity.Reply) error
	Delete(id uuid.UUID) error
}

type replyRepository struct {
	db *gorm.DB
}

func NewReplyRepository(db *gorm.DB) ReplyRepository {
	return &replyRepository{db: db}
}

// Create implements ReplyRepository.
func (r *replyRepository) Create(post *entity.Reply) {
	panic("unimplemented")
}

// Delete implements ReplyRepository.
func (r *replyRepository) Delete(id uuid.UUID) error {
	panic("unimplemented")
}

// FindAll implements ReplyRepository.
func (r *replyRepository) FindAll(page int, pageSize int) ([]entity.Reply, int64, error) {
	panic("unimplemented")
}

// FindById implements ReplyRepository.
func (r *replyRepository) FindById(id uuid.UUID) (*entity.Reply, error) {
	panic("unimplemented")
}

// Update implements ReplyRepository.
func (r *replyRepository) Update(reply *entity.Reply) error {
	panic("unimplemented")
}
