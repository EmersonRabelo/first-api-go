package like

import "github.com/google/uuid"

type LikeCreateDTO struct {
	UserId uuid.UUID `json:"user_id" binding:"required"`
	PostId uuid.UUID `json:"post_id" binding:"required"`
}
