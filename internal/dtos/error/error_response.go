package error

type ErrorResponse struct {
	Error   string `json:"error" example:"Mensagem de erro"`
	Details string `json:"details,omitempty" example:"Detalhes adicionais"`
}
