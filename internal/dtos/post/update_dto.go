package post

type PostUpdateDTO struct {
	Body     *string `json:"body,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
}
