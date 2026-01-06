package consumer

import (
	"errors"
	"fmt"

	contracts "github.com/EmersonRabelo/first-api-go/internal/dtos/report/message"
	"github.com/EmersonRabelo/first-api-go/internal/entity"
	"github.com/EmersonRabelo/first-api-go/internal/repository"
)

var ErrInvalidMessage = errors.New("invalid message")

type ReportRepository interface {
	InsertIfNotExists(rep *entity.Report) error
}

type ConsumerReportService struct {
	repository repository.ReportRepository
}

func NewConsumerReportService(repository repository.ReportRepository) *ConsumerReportService {
	return &ConsumerReportService{
		repository: repository,
	}
}

func (crs *ConsumerReportService) Create(msg contracts.ReportAnalysisResultMessage) error {
	fmt.Println("Response: ", msg)
	return nil
}
