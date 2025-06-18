package domain

import (
	"time"

	"github.com/google/uuid"
)

type JobStatus string

const (
	JobStatusPending    JobStatus = "pending"
	JobStatusProcessing JobStatus = "processing"
	JobStatusCompleted  JobStatus = "completed"
	JobStatusFailed     JobStatus = "failed"
)

type Job struct {
	ID           uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Status       JobStatus  `json:"status" gorm:"type:varchar(20);not null;default:'pending'"`
	Type         string     `json:"type" gorm:"type:varchar(50);not null"`
	Progress     int        `json:"progress" gorm:"default:0"`
	TotalItems   int        `json:"total_items" gorm:"default:0"`
	ErrorMessage *string    `json:"error_message,omitempty"`
	CreatedAt    time.Time  `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
	CompletedAt  *time.Time `json:"completed_at,omitempty"`
}
