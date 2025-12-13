package like

import (
	"time"

	"github.com/google/uuid"
)

type LikeResponseDTO struct {
	Id        uuid.UUID  `json:"id"`
	UserId    uuid.UUID  `json:"user_id"`
	PostId    uuid.UUID  `json:"post_id"`
	Quantity  uint64     `json:"quantity"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
