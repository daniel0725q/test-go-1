package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/truora/microservice/internal/domain"
)

type StockRatingRepository interface {
	Create(ctx context.Context, rating *domain.StockRating) error
	CreateBatch(ctx context.Context, ratings []*domain.StockRating) error
	GetByID(ctx context.Context, id uint) (*domain.StockRating, error)
	GetByTicker(ctx context.Context, ticker string) ([]*domain.StockRating, error)
	GetLatestByTicker(ctx context.Context, ticker string) (*domain.StockRating, error)
}

type stockRatingRepository struct {
	db *gorm.DB
}

func NewStockRatingRepository(db *gorm.DB) StockRatingRepository {
	return &stockRatingRepository{db: db}
}

func (r *stockRatingRepository) Create(ctx context.Context, rating *domain.StockRating) error {
	result := r.db.WithContext(ctx).Create(rating)
	return result.Error
}

func (r *stockRatingRepository) CreateBatch(ctx context.Context, ratings []*domain.StockRating) error {
	result := r.db.WithContext(ctx).CreateInBatches(ratings, 100)
	return result.Error
}

func (r *stockRatingRepository) GetByID(ctx context.Context, id uint) (*domain.StockRating, error) {
	var rating domain.StockRating
	result := r.db.WithContext(ctx).First(&rating, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &rating, nil
}

func (r *stockRatingRepository) GetByTicker(ctx context.Context, ticker string) ([]*domain.StockRating, error) {
	var ratings []*domain.StockRating
	result := r.db.WithContext(ctx).Where("ticker = ?", ticker).Find(&ratings)
	if result.Error != nil {
		return nil, result.Error
	}
	return ratings, nil
}

func (r *stockRatingRepository) GetLatestByTicker(ctx context.Context, ticker string) (*domain.StockRating, error) {
	var rating domain.StockRating
	result := r.db.WithContext(ctx).
		Where("ticker = ?", ticker).
		Order("time DESC").
		First(&rating)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &rating, nil
}
