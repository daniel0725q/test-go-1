package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/truora/microservice/internal/domain"
)

type JobRepository interface {
	Create(ctx context.Context, job *domain.Job) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Job, error)
	Update(ctx context.Context, job *domain.Job) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status domain.JobStatus, progress int, totalItems int) error
	MarkCompleted(ctx context.Context, id uuid.UUID) error
	MarkFailed(ctx context.Context, id uuid.UUID, errorMessage string) error
}

type jobRepository struct {
	db *gorm.DB
}

func NewJobRepository(db *gorm.DB) JobRepository {
	return &jobRepository{db: db}
}

func (r *jobRepository) Create(ctx context.Context, job *domain.Job) error {
	result := r.db.WithContext(ctx).Create(job)
	return result.Error
}

func (r *jobRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Job, error) {
	var job domain.Job
	result := r.db.WithContext(ctx).First(&job, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &job, nil
}

func (r *jobRepository) Update(ctx context.Context, job *domain.Job) error {
	result := r.db.WithContext(ctx).Save(job)
	return result.Error
}

func (r *jobRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.JobStatus, progress int, totalItems int) error {
	result := r.db.WithContext(ctx).Model(&domain.Job{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":      status,
			"progress":    progress,
			"total_items": totalItems,
			"updated_at":  time.Now(),
		})
	return result.Error
}

func (r *jobRepository) MarkCompleted(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	result := r.db.WithContext(ctx).Model(&domain.Job{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       domain.JobStatusCompleted,
			"completed_at": &now,
			"updated_at":   now,
		})
	return result.Error
}

func (r *jobRepository) MarkFailed(ctx context.Context, id uuid.UUID, errorMessage string) error {
	now := time.Now()
	result := r.db.WithContext(ctx).Model(&domain.Job{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":        domain.JobStatusFailed,
			"error_message": errorMessage,
			"updated_at":    now,
		})
	return result.Error
}
