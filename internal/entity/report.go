package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProcessStatus string

const (
	StatusPending    ProcessStatus = "pending"    // report recebido, aguardando processamento externo
	StatusProcessing ProcessStatus = "processing" // em processamento (opcional, dependendo do controle do worker)
	StatusDone       ProcessStatus = "done"       // processado com sucesso
	StatusError      ProcessStatus = "error"      // erro no processamento
)

type Report struct {
	Id                        uuid.UUID      `gorm:"type:uuid;primaryKey;column:report_id" json:"id"`
	PostId                    uuid.UUID      `gorm:"type:uuid;not null;column:post_id;index:idx_post_reports" json:"post_id"`
	Post                      Post           `gorm:"foreignKey:PostId;references:Id"`
	ReporterId                uuid.UUID      `gorm:"type:uuid;not null;column:reporter_id;index:idx_user_report" json:"reporter_id"`
	Reporter                  User           `gorm:"foreignKey:ReporterId;references:Id"`
	Status                    ProcessStatus  `gorm:"type:varchar(30);column:status;default:'pending'" json:"status"`
	CreatedAt                 time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP;column:created_at" json:"created_at"`
	UpdatedAt                 time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt                 gorm.DeletedAt `gorm:"index;column:deleted_at" json:"deleted_at"`
	PerspectiveToxicity       *float64       `gorm:"column:perspective_toxicity" json:"perspective_toxicity,omitempty"`
	PerspectiveSevereToxicity *float64       `gorm:"column:perspective_severe_toxicity" json:"perspective_severe_toxicity,omitempty"`
	PerspectiveInsult         *float64       `gorm:"column:perspective_insult" json:"perspective_insult,omitempty"`
	PerspectiveProfanity      *float64       `gorm:"column:perspective_profanity" json:"perspective_profanity,omitempty"`
	PerspectiveThreat         *float64       `gorm:"column:perspective_threat" json:"perspective_threat,omitempty"`
	PerspectiveIdentityAttack *float64       `gorm:"column:perspective_identity_attack" json:"perspective_identity_attack,omitempty"`
	PerspectiveLanguage       *string        `gorm:"column:perspective_language" json:"perspective_language,omitempty"`
	PerspectiveResponseAt     *time.Time     `gorm:"column:perspective_response_at" json:"perspective_response_at,omitempty"`
}

func (report *Report) TableName() string {
	return "reports"
}
