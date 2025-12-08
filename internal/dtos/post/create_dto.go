package post

import "github.com/google/uuid"

type PostCreateDTO struct {
	UserId uuid.UUID `json:"user_id" binding:"required"`
	Body   string    `json:"body" binding:"required,max=280"`
}
