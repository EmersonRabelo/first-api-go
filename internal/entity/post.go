package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProcessFlag string

const (
	Visible             ProcessFlag = "visible"
	Limited             ProcessFlag = "limited"
	HiddenPendingReview ProcessFlag = "hidden_pending_review"
	Removed             ProcessFlag = "removed"
)

type Post struct {
	Id        uuid.UUID      `gorm:"type:uuid;primaryKey;column:post_id" json:"id"`
	UserId    uuid.UUID      `gorm:"type:uuid;not null;column:user_id;uniqueIndex:idx_posts_user" json:"user_id"`
	User      User           `gorm:"foreignKey:UserId;references:Id"`
	Body      string         `gorm:"size:280;column:post_body" json:"body"`
	Flag      ProcessFlag    `gorm:"type:varchar(48);column:flag;default:'visible'" json:"flag"`
	IsActive  bool           `gorm:"not null;default:true;column:is_active" json:"is_active"`
	CreatedAt time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP;column:created_at;uniqueIndex:idx_posts_created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index;column:deleted_at" json:"deleted_at"`
}

func (p *Post) TableName() string {
	return "posts"
}
