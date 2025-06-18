package domain

import "time"

// StockRating represents the database model for stock ratings
type StockRating struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Ticker     string    `json:"ticker" gorm:"index"`
	TargetFrom string    `json:"target_from"`
	TargetTo   string    `json:"target_to"`
	Company    string    `json:"company"`
	Action     string    `json:"action"`
	Brokerage  string    `json:"brokerage"`
	RatingFrom string    `json:"rating_from"`
	RatingTo   string    `json:"rating_to"`
	Time       time.Time `json:"time"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
