package dto

import (
	"time"

	"github.com/truora/microservice/internal/domain"
)

type Response struct {
	Items    []*StockRatingResponse `json:"items"`
	NextPage string                 `json:"next_page"`
}

// StockRatingResponse represents the API response for stock ratings
type StockRatingResponse struct {
	Ticker     string    `json:"ticker"`
	TargetFrom string    `json:"target_from"`
	TargetTo   string    `json:"target_to"`
	Company    string    `json:"company"`
	Action     string    `json:"action"`
	Brokerage  string    `json:"brokerage"`
	RatingFrom string    `json:"rating_from"`
	RatingTo   string    `json:"rating_to"`
	Time       time.Time `json:"time"`
}

// ToDomain converts the DTO to a domain model
func (dto *StockRatingResponse) ToDomain() *domain.StockRating {
	return &domain.StockRating{
		Ticker:     dto.Ticker,
		TargetFrom: dto.TargetFrom,
		TargetTo:   dto.TargetTo,
		Company:    dto.Company,
		Action:     dto.Action,
		Brokerage:  dto.Brokerage,
		RatingFrom: dto.RatingFrom,
		RatingTo:   dto.RatingTo,
		Time:       dto.Time,
	}
}

// FromDomain creates a DTO from a domain model
func FromDomain(model *domain.StockRating) *StockRatingResponse {
	return &StockRatingResponse{
		Ticker:     model.Ticker,
		TargetFrom: model.TargetFrom,
		TargetTo:   model.TargetTo,
		Company:    model.Company,
		Action:     model.Action,
		Brokerage:  model.Brokerage,
		RatingFrom: model.RatingFrom,
		RatingTo:   model.RatingTo,
		Time:       model.Time,
	}
}
