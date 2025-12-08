package reply

import "github.com/google/uuid"

type ReplyCreateDTO struct {
	UserId   uuid.UUID `json:"user_id" binding:"required"`
	PostId   uuid.UUID `json:"post_id" binding:"required"`
	Body     string    `json:"body" binding:"required,max=280"`
	Quantity uint64    `json:"quantity" binding:"required"`
}
