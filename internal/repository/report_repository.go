package repository

import (
	"context"

	"github.com/EmersonRabelo/first-api-go/internal/entity"
	"gorm.io/gorm"
)

type ReportRepository interface {
	Create(context context.Context, report *entity.Report) error
	UpdateReportStatus(context context.Context, report *entity.Report) error
}

type reportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepository{db: db}
}

func (r *reportRepository) Create(context context.Context, report *entity.Report) error {
	return r.db.WithContext(context).Create(report).Error
}

func (r *reportRepository) UpdateReportStatus(context context.Context, report *entity.Report) error {
	return r.db.WithContext(context).Save(report).Error
}
