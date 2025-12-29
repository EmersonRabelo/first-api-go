package report

import (
	"time"

	"github.com/google/uuid"
)

type ProcessStatus string

const (
	StatusPending    ProcessStatus = "pending"    // report recebido, aguardando processamento externo
	StatusProcessing ProcessStatus = "processing" // em processamento (opcional, dependendo do controle do worker)
	StatusDone       ProcessStatus = "done"       // processado com sucesso
	StatusError      ProcessStatus = "error"      // erro no processamento
)

type CreateReportResponse struct {
	Id        uuid.UUID     `json:"id"`
	Message   string        `json:"msg"`
	Status    ProcessStatus `json:"status"`
	CreatedAt time.Time     `json:"created_at"`
}
