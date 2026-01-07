package consumer

import (
	"errors"
	"fmt"
	"time"

	contracts "github.com/EmersonRabelo/first-api-go/internal/dtos/report/message"
	"github.com/EmersonRabelo/first-api-go/internal/entity"
	"github.com/EmersonRabelo/first-api-go/internal/repository"
	"github.com/google/uuid"
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

	if msg.ReportId == uuid.Nil {
		return fmt.Errorf("%w: Id is null", ErrInvalidMessage)
	}

	report, err := crs.repository.FindById(msg.ReportId)

	if err != nil {
		return errors.New("Report not find")
	}

	report.PerspectiveToxicity = msg.Toxicity
	report.PerspectiveSevereToxicity = msg.SevereToxicity
	report.PerspectiveIdentityAttack = msg.IdentityAttack
	report.PerspectiveInsult = msg.Insult
	report.PerspectiveProfanity = msg.Profanity
	report.PerspectiveThreat = msg.Threat
	report.PerspectiveLanguage = msg.Language
	report.PerspectiveResponseAt = msg.AnalyzedAt
	report.UpdatedAt = time.Now()
	report.Status = entity.StatusDone

	if err := crs.repository.Update(report); err != nil {
		return errors.New("Error persist content analysis")
	}

	return nil
}
