package user

type UserResponseListDTO struct {
	Data       []UserResponseDTO `json:"data"`
	Total      int64             `json:"total"`
	PageSize   int               `json:"pageSize"`
	TotalPages int               `json:"totalPages"`
}
