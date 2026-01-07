package repository

import (
	"context"

	"github.com/EmersonRabelo/first-api-go/internal/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReportRepository interface {
	Create(context context.Context, report *entity.Report) error
	UpdateReportStatus(context context.Context, report *entity.Report) error
	FindById(id uuid.UUID) (*entity.Report, error)
	Update(report *entity.Report) error
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

func (r *reportRepository) FindById(id uuid.UUID) (*entity.Report, error) {
	var report entity.Report

	if err := r.db.First(&report, id).Error; err != nil {
		return nil, err
	}

	return &report, nil
}

func (r *reportRepository) Update(report *entity.Report) error {
	return r.db.Save(report).Error
}
