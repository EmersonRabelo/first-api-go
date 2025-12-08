package user

type CreateDTO struct {
	Name  string `json:"name"  binding:"required,min=3,max=50" example:"John Due"`
	Email string `json:"email" binding:"required,email" example:"johndue2025@email.com"`
}
