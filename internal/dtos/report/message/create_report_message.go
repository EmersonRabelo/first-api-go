package message

import (
	"time"

	"github.com/google/uuid"
)

type CreateReportMessage struct {
	Id         uuid.UUID  `json:"id"`
	PostId     uuid.UUID  `json:"post_id"`
	ReporterId uuid.UUID  `json:"reporter_id"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
	Body       string     `json:"body"`
}
