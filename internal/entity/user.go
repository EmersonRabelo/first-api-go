package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id        uuid.UUID `gorm:"type:uuid;primaryKey;column:user_id" json:"id"`
	Name      string    `gorm:"size:50;not null;column:user_name" json:"name"`
	Email     string    `gorm:"size:255;not null;uniqueIndex;column:user_email" json:"email"`
	IsActive  bool      `gorm:"not null;default:true;column:is_active" json:"is_active"`
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP;column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

func (u *User) TableName() string {
	return "users"
}
