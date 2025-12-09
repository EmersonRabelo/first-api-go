package user

import (
	"time"

	"github.com/google/uuid"
)

type UserResponseDTO struct {
	Id        uuid.UUID  `json:"id" example:"01890f54-d4aa-7b4a-a102-acae7b6a89e8"`
	Name      string     `json:"name" example:"John Due"`
	Email     string     `json:"email" example:"johndue2025@email.com"`
	IsActive  bool       `json:"is_active" example:"true"`
	CreatedAT time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
