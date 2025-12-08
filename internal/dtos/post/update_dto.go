package post

import "time"

type PostUpdateDTO struct {
	Body      *string   `json:"body,omitempty"`
	IsActive  *bool     `json:"is_active,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}
