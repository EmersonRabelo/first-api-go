package like

type LikeResponseListDTO struct {
	Data       []LikeResponseDTO `json:"data"`
	Total      int64             `json:"total"`
	PageSize   int               `json:"pageSize"`
	TotalPages int               `json:"totalPages"`
}
