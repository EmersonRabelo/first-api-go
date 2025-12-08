package user

type UpdateDTO struct {
	Name     string `json:"name"  binding:"min=3,max=50" example:"John Due"`
	Email    string `json:"email" binding:"email" example:"johndue2025@email.com"`
	IsActive bool   `json:"is_active" example:"true"`
}
