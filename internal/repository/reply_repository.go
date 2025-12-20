package repository

import (
	"time"

	"github.com/EmersonRabelo/first-api-go/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReplyRepository interface {
	Create(reply *entity.Reply) error
	FindById(id *uuid.UUID) (*entity.Reply, error)
	FindAll(postID, userID *uuid.UUID, start, end time.Time, page int, pageSize int) ([]entity.Reply, int64, error)
	GetRepliesCountByPostID(postID *uuid.UUID) (uint64, error)
	Update(reply *entity.Reply) error
	Delete(id *uuid.UUID) error
}

type replyRepository struct {
	db *gorm.DB
}

func NewReplyRepository(db *gorm.DB) ReplyRepository {
	return &replyRepository{db: db}
}

func (r *replyRepository) Create(reply *entity.Reply) error {
	return r.db.Create(reply).Error
}

func (l *replyRepository) FindAll(postID, userID *uuid.UUID, start, end time.Time, page int, pageSize int) ([]entity.Reply, int64, error) {
	var replies []entity.Reply
	var amount int64

	l.db.Model(&entity.Reply{}).Count(&amount)

	offset := (page - 1) * pageSize

	db := l.db.Model(&entity.Reply{})

	if !start.IsZero() {
		db = db.Where("created_at >= ?", start)
	}

	if !end.IsZero() {
		db = db.Where("created_at <= ?", end)
	}

	if postID != nil {
		db = db.Where("post_id = ?", postID)
	}

	if userID != nil {
		db = db.Where("user_id = ?", userID)
	}

	result := db.Offset(offset).Limit(pageSize).Find(&replies)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	return replies, amount, nil
}

func (r *replyRepository) FindById(id *uuid.UUID) (*entity.Reply, error) {
	var reply entity.Reply

	if err := r.db.First(&reply, id).Error; err != nil {
		return nil, err
	}

	return &reply, nil
}

func (r *replyRepository) Update(reply *entity.Reply) error {
	return r.db.Save(reply).Error
}

func (r *replyRepository) Delete(id *uuid.UUID) error {
	return r.db.Delete(&entity.Reply{}, id).Error
}

func (l *replyRepository) GetRepliesCountByPostID(postID *uuid.UUID) (uint64, error) {
	var reply entity.Reply
	if err := l.db.Where("post_id = ?", postID).Last(&reply).Error; err != nil {
		return 0, err
	}
	return reply.Quantity, nil
}
