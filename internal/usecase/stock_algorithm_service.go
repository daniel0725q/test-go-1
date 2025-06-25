package usecase

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/truora/microservice/internal/dto"
	"github.com/truora/microservice/internal/repository"
)

type StockAlgorithmService interface {
	BestTimeToBuyAndSell(ctx context.Context, ticker string, startDate, endDate *time.Time) (*dto.TradingRecommendation, error)
	BestTimeToBuyAndSellMultiple(ctx context.Context, tickers []string, startDate, endDate *time.Time) ([]*dto.TradingRecommendation, error)
	BestTimeToBuyAndSellGlobal(ctx context.Context, startDate, endDate *time.Time) (*dto.TradingRecommendation, error)
}

type stockAlgorithmService struct {
	stockRatingRepo repository.StockRatingRepository
}

func NewStockAlgorithmService(stockRatingRepo repository.StockRatingRepository) StockAlgorithmService {
	return &stockAlgorithmService{
		stockRatingRepo: stockRatingRepo,
	}
}

// BestTimeToBuyAndSell implements the algorithm to find the best time to buy and sell a stock
// based on target price ranges from analyst ratings
func (s *stockAlgorithmService) BestTimeToBuyAndSell(ctx context.Context, ticker string, startDate, endDate *time.Time) (*dto.TradingRecommendation, error) {
	// Get all ratings for the ticker
	ratings, err := s.stockRatingRepo.GetByTicker(ctx, ticker)
	if err != nil {
		return nil, fmt.Errorf("failed to get ratings for ticker %s: %w", ticker, err)
	}

	if len(ratings) == 0 {
		return nil, fmt.Errorf("no ratings found for ticker %s", ticker)
	}

	// Filter by date range if provided
	var filteredRatings []*dto.StockRatingResponse
	for _, rating := range ratings {
		ratingDTO := dto.FromDomain(rating)

		// Apply date filter if provided
		if startDate != nil && ratingDTO.Time.Before(*startDate) {
			continue
		}
		if endDate != nil && ratingDTO.Time.After(*endDate) {
			continue
		}

		filteredRatings = append(filteredRatings, ratingDTO)
	}

	if len(filteredRatings) == 0 {
		return nil, fmt.Errorf("no ratings found for ticker %s in the specified date range", ticker)
	}

	// Sort by time to ensure chronological order
	sort.Slice(filteredRatings, func(i, j int) bool {
		return filteredRatings[i].Time.Before(filteredRatings[j].Time)
	})

	// Extract target prices and convert to float64
	var prices []float64
	var priceData []dto.PricePoint

	for _, rating := range filteredRatings {
		// Use target_from as the price point
		if rating.TargetFrom != "" {
			if price, err := parsePrice(rating.TargetFrom); err == nil {
				prices = append(prices, price)
				priceData = append(priceData, dto.PricePoint{
					Price:     price,
					Time:      rating.Time,
					Brokerage: rating.Brokerage,
					Action:    rating.Action,
					Rating:    rating.RatingFrom,
				})
			}
		}
	}

	if len(prices) < 2 {
		return nil, fmt.Errorf("insufficient price data for ticker %s (need at least 2 price points)", ticker)
	}

	// Find best buy and sell points
	buyIndex, sellIndex, maxProfit := s.findBestBuySellPoints(prices)

	if buyIndex == -1 || sellIndex == -1 {
		return nil, fmt.Errorf("no profitable trading opportunity found for ticker %s", ticker)
	}

	// Create recommendation
	recommendation := &dto.TradingRecommendation{
		Ticker:           ticker,
		BuyPrice:         prices[buyIndex],
		SellPrice:        prices[sellIndex],
		MaxProfit:        maxProfit,
		ProfitPercentage: (maxProfit / prices[buyIndex]) * 100,
		BuyTime:          priceData[buyIndex].Time,
		SellTime:         priceData[sellIndex].Time,
		BuyBrokerage:     priceData[buyIndex].Brokerage,
		SellBrokerage:    priceData[sellIndex].Brokerage,
		BuyAction:        priceData[buyIndex].Action,
		SellAction:       priceData[sellIndex].Action,
		BuyRating:        priceData[buyIndex].Rating,
		SellRating:       priceData[sellIndex].Rating,
		TotalDataPoints:  len(prices),
		DateRange: dto.DateRange{
			StartDate: priceData[0].Time,
			EndDate:   priceData[len(priceData)-1].Time,
		},
	}

	return recommendation, nil
}

// BestTimeToBuyAndSellMultiple analyzes multiple tickers and returns recommendations for each
func (s *stockAlgorithmService) BestTimeToBuyAndSellMultiple(ctx context.Context, tickers []string, startDate, endDate *time.Time) ([]*dto.TradingRecommendation, error) {
	var recommendations []*dto.TradingRecommendation

	for _, ticker := range tickers {
		recommendation, err := s.BestTimeToBuyAndSell(ctx, ticker, startDate, endDate)
		if err != nil {
			// Log error but continue with other tickers
			fmt.Printf("Error analyzing ticker %s: %v\n", ticker, err)
			continue
		}
		recommendations = append(recommendations, recommendation)
	}

	if len(recommendations) == 0 {
		return nil, fmt.Errorf("no valid recommendations found for any ticker")
	}

	// Sort by profit percentage (highest first)
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].ProfitPercentage > recommendations[j].ProfitPercentage
	})

	return recommendations, nil
}

// BestTimeToBuyAndSellGlobal analyzes all stock ratings as if they belonged to the same ticker
// This provides a global market perspective across all available stocks
func (s *stockAlgorithmService) BestTimeToBuyAndSellGlobal(ctx context.Context, startDate, endDate *time.Time) (*dto.TradingRecommendation, error) {
	// Get all ratings from the database (we'll need to add a method to get all ratings)
	// For now, let's use a paginated approach to get all ratings
	var allRatings []*dto.StockRatingResponse
	page := 1
	pageSize := 1000

	for {
		// Get paginated ratings
		ratings, err := s.stockRatingRepo.GetPaginated(ctx, (page-1)*pageSize, pageSize)
		if err != nil {
			return nil, fmt.Errorf("failed to get ratings page %d: %w", page, err)
		}

		if len(ratings) == 0 {
			break // No more ratings
		}

		// Convert to DTOs and filter by date
		for _, rating := range ratings {
			ratingDTO := dto.FromDomain(rating)

			// Apply date filter if provided
			if startDate != nil && ratingDTO.Time.Before(*startDate) {
				continue
			}
			if endDate != nil && ratingDTO.Time.After(*endDate) {
				continue
			}

			allRatings = append(allRatings, ratingDTO)
		}

		// If we got fewer ratings than pageSize, we've reached the end
		if len(ratings) < pageSize {
			break
		}

		page++
	}

	if len(allRatings) == 0 {
		return nil, fmt.Errorf("no ratings found in the specified date range")
	}

	// Sort by time to ensure chronological order
	sort.Slice(allRatings, func(i, j int) bool {
		return allRatings[i].Time.Before(allRatings[j].Time)
	})

	// Extract target prices and convert to float64
	var prices []float64
	var priceData []dto.PricePoint

	for _, rating := range allRatings {
		// Use target_from as the price point
		if rating.TargetFrom != "" {
			if price, err := parsePrice(rating.TargetFrom); err == nil {
				prices = append(prices, price)
				priceData = append(priceData, dto.PricePoint{
					Price:     price,
					Time:      rating.Time,
					Brokerage: rating.Brokerage,
					Action:    rating.Action,
					Rating:    rating.RatingFrom,
					Ticker:    rating.Ticker, // Add ticker information
				})
			}
		}
	}

	if len(prices) < 2 {
		return nil, fmt.Errorf("insufficient price data (need at least 2 price points)")
	}

	// Find best buy and sell points
	buyIndex, sellIndex, maxProfit := s.findBestBuySellPoints(prices)

	if buyIndex == -1 || sellIndex == -1 {
		return nil, fmt.Errorf("no profitable trading opportunity found")
	}

	// Create recommendation
	recommendation := &dto.TradingRecommendation{
		Ticker:           "GLOBAL", // Indicates this is a global analysis
		BuyPrice:         prices[buyIndex],
		SellPrice:        prices[sellIndex],
		MaxProfit:        maxProfit,
		ProfitPercentage: (maxProfit / prices[buyIndex]) * 100,
		BuyTime:          priceData[buyIndex].Time,
		SellTime:         priceData[sellIndex].Time,
		BuyBrokerage:     priceData[buyIndex].Brokerage,
		SellBrokerage:    priceData[sellIndex].Brokerage,
		BuyAction:        priceData[buyIndex].Action,
		SellAction:       priceData[sellIndex].Action,
		BuyRating:        priceData[buyIndex].Rating,
		SellRating:       priceData[sellIndex].Rating,
		BuyTicker:        priceData[buyIndex].Ticker, // Add ticker information
		SellTicker:       priceData[sellIndex].Ticker,
		TotalDataPoints:  len(prices),
		DateRange: dto.DateRange{
			StartDate: priceData[0].Time,
			EndDate:   priceData[len(priceData)-1].Time,
		},
	}

	return recommendation, nil
}

// findBestBuySellPoints implements the core algorithm to find maximum profit
func (s *stockAlgorithmService) findBestBuySellPoints(prices []float64) (int, int, float64) {
	if len(prices) < 2 {
		return -1, -1, 0
	}

	minPrice := prices[0]
	minIndex := 0
	maxProfit := 0.0
	buyIndex := -1
	sellIndex := -1

	for i := 1; i < len(prices); i++ {
		// Calculate potential profit if we sell at current price
		profit := prices[i] - minPrice

		// Update max profit if current profit is higher
		if profit > maxProfit {
			maxProfit = profit
			buyIndex = minIndex
			sellIndex = i
		}

		// Update minimum price if we find a lower price
		if prices[i] < minPrice {
			minPrice = prices[i]
			minIndex = i
		}
	}

	return buyIndex, sellIndex, maxProfit
}

// parsePrice converts string price to float64
func parsePrice(priceStr string) (float64, error) {
	var price float64
	_, err := fmt.Sscanf(priceStr, "%f", &price)
	return price, err
}
