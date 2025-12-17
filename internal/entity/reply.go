package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Reply struct {
	Id        uuid.UUID      `gorm:"type:uuid;primaryKey;column:reply_id" json:"id"`
	UserId    uuid.UUID      `gorm:"type:uuid;not null;column:user_id;uniqueIndex:idx_user_replies" json:"user_id"`
	User      User           `gorm:"foreignKey:UserId;references:Id"`
	PostId    uuid.UUID      `gorm:"type:uuid;not null;column:post_id;uniqueIndex:idx_post_replies" json:"post_id"`
	Post      Post           `gorm:"foreignKey:PostId;references:Id"`
	Body      string         `gorm:"size:280;column:reply_body" json:"body"`
	Quantity  uint64         `gorm:"type:bigint;not null;column:quantity" json:"quantity"`
	CreatedAt time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP;column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index;column:deleted_at" json:"deleted_at"`
}

func (re *Reply) TableName() string {
	return "replies"
}
