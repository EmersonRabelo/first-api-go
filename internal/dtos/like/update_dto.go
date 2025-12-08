package like

import "time"

type LikeUpdateDTO struct {
	Quantity  *uint64   `json:"quantity,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at"`
}
