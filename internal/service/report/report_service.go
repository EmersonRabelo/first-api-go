package report

import (
	"context"
	"errors"
	"time"

	dto "github.com/EmersonRabelo/first-api-go/internal/dtos/report"
	contract "github.com/EmersonRabelo/first-api-go/internal/dtos/report/message"
	"github.com/EmersonRabelo/first-api-go/internal/entity"
	"github.com/EmersonRabelo/first-api-go/internal/service"
	"github.com/google/uuid"
)

type ReportRepository interface {
	Create(context context.Context, report *entity.Report) error
	UpdateReportStatus(context context.Context, report *entity.Report) error
}

type ReportProducer interface {
	Publish(ctx context.Context, msg *contract.CreateReportMessage) error
}

type ReportService struct {
	repository  ReportRepository
	producer    ReportProducer
	postService service.PostService
	userService service.UserService
}

func NewReportService(repository ReportRepository, producer ReportProducer) *ReportService {
	return &ReportService{
		repository: repository,
		producer:   producer,
	}
}

func (s *ReportService) Create(ctx context.Context, postID, reporterId uuid.UUID) (*dto.CreateReportResponse, error) {

	post, err := s.postService.FindById(&postID)

	if err != nil {
		return nil, errors.New("Postagem não encontrado")
	}

	if _, err := s.userService.FindById(&reporterId); err != nil {
		return nil, errors.New("Usuário não encontrado")
	}

	report := &entity.Report{
		Id:         uuid.New(),
		PostId:     postID,
		ReporterId: reporterId,
		Status:     entity.StatusPending,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.repository.Create(ctx, report); err != nil {
		return nil, errors.New("Erro ao persistir a denuncia")
	}

	msg := &contract.CreateReportMessage{
		Id:         report.Id,
		PostId:     report.PostId,
		ReporterId: report.ReporterId,
		CreatedAt:  report.CreatedAt,
		Body:       post.Body,
	}

	if err := s.producer.Publish(ctx, msg); err != nil {

		report.Status = entity.StatusError
		report.UpdatedAt = time.Now()

		if err := s.repository.UpdateReportStatus(ctx, report); err != nil {
			return nil, errors.New("Erro ao atualizar status para: error")
		}

		return s.toReportResponseDto(report), errors.New("Erro ao publicar mensagem")
	}

	return s.toReportResponseDto(report), nil
}

func (s *ReportService) toReportResponseDto(report *entity.Report) *dto.CreateReportResponse {
	return &dto.CreateReportResponse{
		Id:        report.Id,
		Message:   "Sua denuncia foi recebida com sucesso!",
		Status:    dto.ProcessStatus(report.Status),
		CreatedAt: report.CreatedAt,
	}
}
