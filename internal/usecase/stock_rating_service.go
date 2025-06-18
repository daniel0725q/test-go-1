package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/truora/microservice/internal/domain"
	"github.com/truora/microservice/internal/dto"
	"github.com/truora/microservice/internal/repository"
)

type StockRatingService interface {
	CreateStockRating(ctx context.Context, rating *dto.StockRatingResponse) error
	CreateStockRatingBatch(ctx context.Context, ratings []*dto.StockRatingResponse) error
	GetStockRatingByID(ctx context.Context, id uint) (*dto.StockRatingResponse, error)
	GetStockRatingsByTicker(ctx context.Context, ticker string) ([]*dto.StockRatingResponse, error)
	GetLatestStockRatingByTicker(ctx context.Context, ticker string) (*dto.StockRatingResponse, error)
	GetHello(ctx context.Context) (*domain.Job, error)
	GetJobByID(ctx context.Context, jobID uuid.UUID) (*domain.Job, error)
}

type stockRatingService struct {
	stockRatingRepo repository.StockRatingRepository
	jobRepo         repository.JobRepository
	externalAPIRepo repository.ExternalAPIRepository
}

func NewStockRatingService(stockRatingRepo repository.StockRatingRepository, jobRepo repository.JobRepository, externalAPIRepo repository.ExternalAPIRepository) StockRatingService {
	return &stockRatingService{
		stockRatingRepo: stockRatingRepo,
		jobRepo:         jobRepo,
		externalAPIRepo: externalAPIRepo,
	}
}

func (s *stockRatingService) CreateStockRating(ctx context.Context, rating *dto.StockRatingResponse) error {
	domainRating := rating.ToDomain()
	return s.stockRatingRepo.Create(ctx, domainRating)
}

func (s *stockRatingService) CreateStockRatingBatch(ctx context.Context, ratings []*dto.StockRatingResponse) error {
	domainRatings := make([]*domain.StockRating, len(ratings))
	for i, rating := range ratings {
		domainRatings[i] = rating.ToDomain()
	}
	return s.stockRatingRepo.CreateBatch(ctx, domainRatings)
}

func (s *stockRatingService) GetStockRatingByID(ctx context.Context, id uint) (*dto.StockRatingResponse, error) {
	rating, err := s.stockRatingRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if rating == nil {
		return nil, nil
	}
	return dto.FromDomain(rating), nil
}

func (s *stockRatingService) GetStockRatingsByTicker(ctx context.Context, ticker string) ([]*dto.StockRatingResponse, error) {
	ratings, err := s.stockRatingRepo.GetByTicker(ctx, ticker)
	if err != nil {
		return nil, err
	}

	responseRatings := make([]*dto.StockRatingResponse, len(ratings))
	for i, rating := range ratings {
		responseRatings[i] = dto.FromDomain(rating)
	}
	return responseRatings, nil
}

func (s *stockRatingService) GetLatestStockRatingByTicker(ctx context.Context, ticker string) (*dto.StockRatingResponse, error) {
	rating, err := s.stockRatingRepo.GetLatestByTicker(ctx, ticker)
	if err != nil {
		return nil, err
	}
	if rating == nil {
		return nil, nil
	}
	return dto.FromDomain(rating), nil
}

func (s *stockRatingService) GetHello(ctx context.Context) (*domain.Job, error) {
	// Create a new job
	job := &domain.Job{
		ID:     uuid.New(),
		Status: domain.JobStatusPending,
		Type:   "external_api_sync",
	}

	if err := s.jobRepo.Create(ctx, job); err != nil {
		return nil, fmt.Errorf("failed to create job: %w", err)
	}

	// Start async processing
	go s.processExternalAPISync(job.ID)

	return job, nil
}

func (s *stockRatingService) GetJobByID(ctx context.Context, jobID uuid.UUID) (*domain.Job, error) {
	return s.jobRepo.GetByID(ctx, jobID)
}

func (s *stockRatingService) processExternalAPISync(jobID uuid.UUID) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	// Update job status to processing
	if err := s.jobRepo.UpdateStatus(ctx, jobID, domain.JobStatusProcessing, 0, 0); err != nil {
		s.jobRepo.MarkFailed(ctx, jobID, fmt.Sprintf("Failed to update job status: %v", err))
		return
	}

	// Get items from external API
	items, err := s.externalAPIRepo.GetHello(ctx)
	if err != nil {
		s.jobRepo.MarkFailed(ctx, jobID, fmt.Sprintf("Failed to get items from external API: %v", err))
		return
	}

	// Update total items
	if err := s.jobRepo.UpdateStatus(ctx, jobID, domain.JobStatusProcessing, 0, len(items)); err != nil {
		s.jobRepo.MarkFailed(ctx, jobID, fmt.Sprintf("Failed to update job progress: %v", err))
		return
	}

	// Process items in chunks
	chunkSize := 100
	totalProcessed := 0

	for i := 0; i < len(items); i += chunkSize {
		end := i + chunkSize
		if end > len(items) {
			end = len(items)
		}

		chunk := items[i:end]
		domainItems := make([]*domain.StockRating, len(chunk))
		for j, item := range chunk {
			domainItems[j] = item.ToDomain()
		}

		if err := s.stockRatingRepo.CreateBatch(ctx, domainItems); err != nil {
			s.jobRepo.MarkFailed(ctx, jobID, fmt.Sprintf("Failed to store chunk %d-%d: %v", i, end, err))
			return
		}

		totalProcessed += len(chunk)
		if err := s.jobRepo.UpdateStatus(ctx, jobID, domain.JobStatusProcessing, totalProcessed, len(items)); err != nil {
			s.jobRepo.MarkFailed(ctx, jobID, fmt.Sprintf("Failed to update job progress: %v", err))
			return
		}
	}

	// Mark job as completed
	if err := s.jobRepo.MarkCompleted(ctx, jobID); err != nil {
		s.jobRepo.MarkFailed(ctx, jobID, fmt.Sprintf("Failed to mark job as completed: %v", err))
		return
	}
}
