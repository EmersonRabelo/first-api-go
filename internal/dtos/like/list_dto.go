package like

type LikeResponseListDTO struct {
	Data       []LikeResponseDTO `json:"data"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	PageSize   int               `json:"pageSize"`
	TotalPages int               `json:"totalPages"`
}
