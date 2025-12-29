package report

import "github.com/google/uuid"

type CreateReportRequest struct {
	ReporterId uuid.UUID `json:"reporter_id" binding:"required"`
}
