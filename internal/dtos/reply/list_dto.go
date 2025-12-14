package reply

type ReplyResponseListDTO struct {
	Data       []ReplyResponseDTO `json:"data"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	PageSize   int                `json:"pageSize"`
	TotalPages int                `json:"totalPages"`
}
