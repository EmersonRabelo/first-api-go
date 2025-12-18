package repository

import (
	"time"

	"github.com/EmersonRabelo/first-api-go/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LikeRepository interface {
	Create(like *entity.Like) error
	FindById(id *uuid.UUID) (*entity.Like, error)
	FindAll(postId *uuid.UUID, start, end time.Time, page int, pageSize int) ([]entity.Like, int64, error)
	GetLikesCountByPostID(postID *uuid.UUID) (uint64, error)
	Update(like *entity.Like) error
	Delete(id *uuid.UUID) error
}

type likeRepository struct {
	db *gorm.DB
}

func NewLikeRepository(db *gorm.DB) LikeRepository {
	return &likeRepository{db: db}
}

func (l *likeRepository) Create(like *entity.Like) error {
	return l.db.Create(like).Error
}

func (l *likeRepository) FindAll(postId *uuid.UUID, start, end time.Time, page int, pageSize int) ([]entity.Like, int64, error) {
	var likes []entity.Like
	var amount int64

	l.db.Model(&entity.Like{}).Count(&amount)

	offset := (page - 1) * pageSize

	db := l.db.Model(&entity.Like{})

	if !start.IsZero() {
		db = db.Where("created_at >= ?", start)
	}

	if !end.IsZero() {
		db = db.Where("created_at <= ?", end)
	}

	if postId != nil {
		db = db.Where("post_id = ?", postId)
	}

	result := db.Offset(offset).Limit(pageSize).Find(&likes)

	if result.Error != nil {
		return nil, 0, result.Error
	}

	return likes, amount, nil

}

func (l *likeRepository) FindById(id *uuid.UUID) (*entity.Like, error) {
	var like entity.Like

	if err := l.db.First(&like, id).Error; err != nil {
		return nil, err
	}

	return &like, nil
}

func (l *likeRepository) Update(like *entity.Like) error {
	return l.db.Save(like).Error
}

func (l *likeRepository) Delete(id *uuid.UUID) error {
	return l.db.Delete(&entity.Like{}, id).Error
}

func (l *likeRepository) GetLikesCountByPostID(postID *uuid.UUID) (uint64, error) {
	var postLikeCount entity.PostLikesCount
	if err := l.db.
		Table("post_likes_count").
		Where("post_id = ?", postID).
		First(&postLikeCount).Error; err != nil {
		return 0, err
	}
	return postLikeCount.LikeCount, nil
}
