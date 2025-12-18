package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Like struct {
	Id        uuid.UUID      `gorm:"type:uuid;primaryKey;column:like_id" json:"id"`
	UserId    uuid.UUID      `gorm:"type:uuid;not null;column:user_id;uniqueIndex:idx_user_post" json:"user_id"`
	User      User           `gorm:"foreignKey:UserId;references:Id"`
	PostId    uuid.UUID      `gorm:"type:uuid;not null;column:post_id;uniqueIndex:idx_user_post" json:"post_id"`
	Post      Post           `gorm:"foreignKey:PostId;references:Id"`
	CreatedAt time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP;column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index;column:deleted_at" json:"deleted_at"`
}

func (l *Like) TableName() string {
	return "likes"
}
