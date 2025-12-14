package reply

import (
	"time"

	"github.com/google/uuid"
)

type ReplyResponseDTO struct {
	Id        uuid.UUID  `json:"id"`
	UserId    uuid.UUID  `json:"user_id"`
	PostId    uuid.UUID  `json:"post_id"`
	Body      string     `json:"body"`
	Quantity  uint64     `json:"quantity"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
