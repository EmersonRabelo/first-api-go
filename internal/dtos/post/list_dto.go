package post

type PostResponseListDTO struct {
	Data       []PostResponseDTO `json:"data"`
	Total      int64             `json:"total"`
	PageSize   int               `json:"pageSize"`
	TotalPages int               `json:"totalPages"`
}
