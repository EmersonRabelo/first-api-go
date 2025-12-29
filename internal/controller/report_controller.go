package controller

import (
	"net/http"

	dto "github.com/EmersonRabelo/first-api-go/internal/dtos/report"
	errorDTO "github.com/EmersonRabelo/first-api-go/internal/dtos/shared/error"
	report "github.com/EmersonRabelo/first-api-go/internal/service/report"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ReportHandler struct {
	service *report.ReportService
}

func NewReportHandler(service *report.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

func (r *ReportHandler) Create(context *gin.Context) {
	var req dto.CreateReportRequest

	postID, err := uuid.Parse(context.Param("id"))

	if err != nil {
		context.JSON(http.StatusBadRequest, errorDTO.ErrorResponse{Error: "Id inválido"})
		return
	}

	// Bind e Validacao automática do JSON
	if err := context.ShouldBindJSON(&req); err != nil {
		context.JSON(http.StatusBadRequest, errorDTO.ErrorResponse{Error: "Dados inválidos", Details: err.Error()})
		return
	}

	report, err := r.service.Create(context.Request.Context(), postID, req.ReporterId)

	if err != nil {
		context.JSON(http.StatusInternalServerError, errorDTO.ErrorResponse{Error: "Não foi possivel processar sua denuncia, tente mais tarde", Details: err.Error()})
		return
	}

	context.JSON(http.StatusCreated, report)
}
