package usecase

import (
	"context"
	"fmt"

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
	GetHello(ctx context.Context) ([]*dto.StockRatingResponse, error)
}

type stockRatingService struct {
	stockRatingRepo repository.StockRatingRepository
	externalAPIRepo repository.ExternalAPIRepository
}

func NewStockRatingService(stockRatingRepo repository.StockRatingRepository, externalAPIRepo repository.ExternalAPIRepository) StockRatingService {
	return &stockRatingService{
		stockRatingRepo: stockRatingRepo,
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

func (s *stockRatingService) GetHello(ctx context.Context) ([]*dto.StockRatingResponse, error) {
	items, err := s.externalAPIRepo.GetHello(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get items from external API: %w", err)
	}

	// Convert DTOs to domain models
	domainItems := make([]*domain.StockRating, len(items))
	for i, item := range items {
		domainItems[i] = item.ToDomain()
	}

	// Store all items in the database
	if err := s.stockRatingRepo.CreateBatch(ctx, domainItems); err != nil {
		return nil, fmt.Errorf("failed to store items in database: %w", err)
	}

	return items, nil
}
