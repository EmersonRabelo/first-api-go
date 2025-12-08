package reply

type UserResponseListDTO struct {
	Data       []ReplyResponseDTO `json:"data"`
	Total      int64              `json:"total"`
	PageSize   int                `json:"pageSize"`
	TotalPages int                `json:"totalPages"`
}
