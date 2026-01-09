package consumer

import (
	"errors"
	"fmt"
	"time"

	dto "github.com/EmersonRabelo/first-api-go/internal/dtos/report"
	contracts "github.com/EmersonRabelo/first-api-go/internal/dtos/report/message"
	"github.com/EmersonRabelo/first-api-go/internal/entity"
	"github.com/EmersonRabelo/first-api-go/internal/repository"
	"github.com/google/uuid"
)

var ErrInvalidMessage = errors.New("invalid message")

type ReportRepository interface {
	InsertIfNotExists(rep *entity.Report) error
	moderationDecision(moderationData *dto.ReportAnalysisModerationData) entity.ProcessFlag
}

type ConsumerReportService struct {
	reportRepository repository.ReportRepository
	postRepository   repository.PostRepository
}

func NewConsumerReportService(reportRepository repository.ReportRepository, postRepository repository.PostRepository) *ConsumerReportService {
	return &ConsumerReportService{
		reportRepository: reportRepository,
		postRepository:   postRepository,
	}
}

func (crs *ConsumerReportService) Create(msg contracts.ReportAnalysisResultMessage) error {

	if msg.ReportId == uuid.Nil {
		return fmt.Errorf("%w: Id is null", ErrInvalidMessage)
	}

	report, err := crs.reportRepository.FindById(msg.ReportId)

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

	if err := crs.reportRepository.Update(report); err != nil {
		return errors.New("Error persist content analysis")
	}

	post, err := crs.postRepository.FindById(&report.PostId)

	if err != nil {
		return errors.New("Post not find")
	}

	reportModerationData := dto.ReportAnalysisModerationData{
		Toxicity:       *msg.Toxicity,
		SevereToxicity: *msg.SevereToxicity,
		IdentityAttack: *msg.IdentityAttack,
		Insult:         *msg.Insult,
		Profanity:      *msg.Profanity,
		Threat:         *msg.Threat,
	}

	post.Flag = crs.moderationDecision(&reportModerationData)

	if err := crs.postRepository.Update(post); err != nil {
		return errors.New("Error updating post flag")
	}

	if post.Flag == entity.Removed || post.Flag == entity.HiddenPendingReview {
		if err := crs.postRepository.Delete(&post.Id); err != nil {
			return err
		}
	}

	return nil
}

func (crs *ConsumerReportService) moderationDecision(m *dto.ReportAnalysisModerationData) entity.ProcessFlag {
	// Hard rules (alta severidade)
	if m.Threat >= 0.70 || m.IdentityAttack >= 0.70 || m.SevereToxicity >= 0.70 {
		return entity.Removed
	}

	// Medium severity
	if m.Threat >= 0.35 || m.IdentityAttack >= 0.35 || m.SevereToxicity >= 0.45 {
		return entity.HiddenPendingReview
	}

	// Composite
	composite := 0.45*m.Toxicity + 0.35*m.Insult + 0.20*m.Profanity

	if composite >= 0.80 {
		return entity.HiddenPendingReview
	}
	if composite >= 0.65 {
		return entity.Limited
	}

	return entity.Visible
}
