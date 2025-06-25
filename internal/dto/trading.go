package dto

import "time"

// PricePoint represents a single price point with metadata
type PricePoint struct {
	Price     float64   `json:"price"`
	Time      time.Time `json:"time"`
	Brokerage string    `json:"brokerage"`
	Action    string    `json:"action"`
	Rating    string    `json:"rating"`
	Ticker    string    `json:"ticker"`
}

// DateRange represents a date range
type DateRange struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

// TradingRecommendation represents a trading recommendation from the algorithm
type TradingRecommendation struct {
	Ticker           string    `json:"ticker"`
	BuyPrice         float64   `json:"buy_price"`
	SellPrice        float64   `json:"sell_price"`
	MaxProfit        float64   `json:"max_profit"`
	ProfitPercentage float64   `json:"profit_percentage"`
	BuyTime          time.Time `json:"buy_time"`
	SellTime         time.Time `json:"sell_time"`
	BuyBrokerage     string    `json:"buy_brokerage"`
	SellBrokerage    string    `json:"sell_brokerage"`
	BuyAction        string    `json:"buy_action"`
	SellAction       string    `json:"sell_action"`
	BuyRating        string    `json:"buy_rating"`
	SellRating       string    `json:"sell_rating"`
	BuyTicker        string    `json:"buy_ticker,omitempty"`
	SellTicker       string    `json:"sell_ticker,omitempty"`
	TotalDataPoints  int       `json:"total_data_points"`
	DateRange        DateRange `json:"date_range"`
}

// TradingAnalysisRequest represents a request for trading analysis
type TradingAnalysisRequest struct {
	Tickers   []string  `json:"tickers"`
	StartDate time.Time `json:"start_date,omitempty"`
	EndDate   time.Time `json:"end_date,omitempty"`
}
