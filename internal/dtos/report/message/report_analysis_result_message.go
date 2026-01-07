package message

import (
	"time"

	"github.com/google/uuid"
)

type ReportAnalysisResultMessage struct {
	ReportId       uuid.UUID  `json:"report_id"`
	Toxicity       *float64   `json:"toxicity,omitempty"`
	SevereToxicity *float64   `json:"severe_toxicity,omitempty"`
	IdentityAttack *float64   `json:"identity_attack,omitempty"`
	Insult         *float64   `json:"insult,omitempty"`
	Profanity      *float64   `json:"profanity,omitempty"`
	Threat         *float64   `json:"threat,omitempty"`
	Language       *string    `json:"language,omitempty"`
	AnalyzedAt     *time.Time `json:"analyzed_at,omitempty"`
}
